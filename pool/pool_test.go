package pool_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	bridge "github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevmbridge"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Revert"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	senderPrivateKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	senderAddress    = "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"
)

var (
	l2BridgeAddr = common.HexToAddress("0x00000000000000000000000000000001")
)

var (
	stateDBCfg = dbutils.NewStateConfigFromEnv()
	poolDBCfg  = dbutils.NewPoolConfigFromEnv()
	genesis    = state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: senderAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "90000000000000000000000000000000000000000000000000000000000",
			},
		},
	}
	cfg = pool.Config{
		FreeClaimGasLimit:              150000,
		MaxTxBytesSize:                 30132,
		MaxTxDataBytesSize:             30000,
		MinAllowedGasPriceInterval:     cfgTypes.NewDuration(5 * time.Minute),
		PollMinAllowedGasPriceInterval: cfgTypes.NewDuration(15 * time.Second),
		DefaultMinGasPriceAllowed:      1000000000,
	}
	gasPrice = big.NewInt(1000000000)
	gasLimit = uint64(21000)
	chainID  = big.NewInt(1337)
)

func TestMain(m *testing.M) {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stderr"},
	})

	code := m.Run()
	os.Exit(code)
}

func Test_AddTx(t *testing.T) {
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}

	const chainID = 2576980377
	p := setupPool(t, s, st, chainID, ctx)

	tx := new(types.Transaction)
	expectedTxEncoded := "0xf86880843b9aca008252089400000000000000000000000000000000000000008080850133333355a03ee24709870c8dbc67884c9c8acb864c1aceaaa7332b9a3db0d7a5d7c68eb8e4a0302980b070f5e3ffca3dc27b07daf69d66ab27d4df648e0b3ed059cf23aa168d"
	b, err := hex.DecodeHex(expectedTxEncoded)
	if err != nil {
		t.Error(err)
	}
	tx.UnmarshalBinary(b) //nolint:gosec,errcheck

	err = p.AddTx(ctx, *tx, "")
	if err != nil {
		t.Error(err)
	}

	rows, err := poolSqlDB.Query(ctx, "SELECT hash, encoded, decoded, status, used_steps FROM pool.transaction")
	defer rows.Close() // nolint:staticcheck
	if err != nil {
		t.Error(err)
	}

	c := 0
	for rows.Next() {
		var hash, encoded, decoded, status string
		var usedSteps int
		err := rows.Scan(&hash, &encoded, &decoded, &status, &usedSteps)
		if err != nil {
			t.Error(err)
		}
		b, _ := tx.MarshalJSON()

		assert.Equal(t, "0x3c499a6308dbf4e67bd4e949b0b609e3a0a5a7fd6a497acb23e37ae7f0a923cc", hash, "invalid hash")
		assert.Equal(t, expectedTxEncoded, encoded, "invalid encoded")
		assert.JSONEq(t, string(b), decoded, "invalid decoded")
		assert.Equal(t, string(pool.TxStatusPending), status, "invalid tx status")
		assert.Greater(t, usedSteps, 0, "invalid used steps")
		c++
	}

	assert.Equal(t, 1, c, "invalid number of txs in the pool")
}

func Test_AddTx_OversizedData(t *testing.T) {
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: senderAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "1000000000000000000000",
			},
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}

	const chainID = 2576980377
	p := pool.NewPool(cfg, s, st, common.Address{}, chainID)

	b := make([]byte, cfg.MaxTxBytesSize+1)
	to := common.HexToAddress(operations.DefaultSequencerAddress)
	tx := types.NewTransaction(0, to, big.NewInt(0), gasLimit, big.NewInt(0), b)

	// GetAuth configures and returns an auth object.
	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, chainID)
	require.NoError(t, err)
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, "")
	require.EqualError(t, err, pool.ErrOversizedData.Error())
}

func Test_AddPreEIP155Tx(t *testing.T) {
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: senderAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "1000000000000000000000",
			},
			{
				Address: "0x4d5Cf5032B2a844602278b01199ED191A86c93ff",
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "200000000000000000000",
			},
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}

	const chainID = 2576980377
	p := setupPool(t, s, st, chainID, ctx)

	batchL2Data := "0xe580843b9aca00830186a0941275fbb540c8efc58b812ba83b0d0b8b9917ae98808464fbb77c6b39bdc5f8e458aba689f2a1ff8c543a94e4817bda40f3fe34080c4ab26c1e3c2fc2cda93bc32f0a79940501fd505dcf48d94abfde932ebf1417f502cb0d9de81b"
	b, err := hex.DecodeHex(batchL2Data)
	require.NoError(t, err)
	txs, _, err := state.DecodeTxs(b)
	require.NoError(t, err)

	tx := txs[0]

	err = p.AddTx(ctx, tx, "")
	require.NoError(t, err)

	rows, err := poolSqlDB.Query(ctx, "SELECT hash, encoded, decoded, status FROM pool.transaction")
	defer rows.Close() // nolint:staticcheck
	require.NoError(t, err)

	c := 0
	for rows.Next() {
		var hash, encoded, decoded, status string
		err := rows.Scan(&hash, &encoded, &decoded, &status)
		require.NoError(t, err)

		b, err := tx.MarshalBinary()
		require.NoError(t, err)

		bJSON, err := tx.MarshalJSON()
		require.NoError(t, err)

		assert.Equal(t, tx.Hash().String(), hash, "invalid hash")
		assert.Equal(t, hex.EncodeToHex(b), encoded, "invalid encoded")
		assert.JSONEq(t, string(bJSON), decoded, "invalid decoded")
		assert.Equal(t, string(pool.TxStatusPending), status, "invalid tx status")
		c++
	}

	assert.Equal(t, 1, c, "invalid number of txs in the pool")
}

func Test_GetPendingTxs(t *testing.T) {
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := setupPool(t, s, st, chainID.Uint64(), ctx)

	const txsCount = 10
	const limit = 5

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		if err := p.AddTx(ctx, *signedTx, ""); err != nil {
			t.Error(err)
		}
	}

	txs, err := p.GetPendingTxs(ctx, false, limit)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, limit, len(txs))

	for i := 0; i < txsCount; i++ {
		assert.Equal(t, pool.TxStatusPending, txs[0].Status)
	}
}

func Test_GetPendingTxsZeroPassed(t *testing.T) {
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := setupPool(t, s, st, chainID.Uint64(), ctx)

	const txsCount = 10
	const limit = 0

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		if err := p.AddTx(ctx, *signedTx, ""); err != nil {
			t.Error(err)
		}
	}

	txs, err := p.GetPendingTxs(ctx, false, limit)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, txsCount, len(txs))

	for i := 0; i < txsCount; i++ {
		assert.Equal(t, pool.TxStatusPending, txs[0].Status)
	}
}

func Test_GetTopPendingTxByProfitabilityAndZkCounters(t *testing.T) {
	ctx := context.Background()
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer stateSqlDB.Close()

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := setupPool(t, s, st, chainID.Uint64(), ctx)

	const txsCount = 10

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), gasLimit, big.NewInt(gasPrice.Int64()+int64(i)), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		if err := p.AddTx(ctx, *signedTx, ""); err != nil {
			t.Error(err)
		}
	}

	txs, err := p.GetTxs(ctx, pool.TxStatusPending, false, 1, 10)
	require.NoError(t, err)
	// bcs it's sorted by nonce, tx with the lowest nonce is expected here
	assert.Equal(t, txs[0].Transaction.Nonce(), uint64(0))
}

func Test_UpdateTxsStatus(t *testing.T) {
	ctx := context.Background()

	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := setupPool(t, s, st, chainID.Uint64(), ctx)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx1, ""); err != nil {
		t.Error(err)
	}

	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx2, ""); err != nil {
		t.Error(err)
	}

	err = p.UpdateTxsStatus(ctx, []string{signedTx1.Hash().String(), signedTx2.Hash().String()}, pool.TxStatusInvalid)
	if err != nil {
		t.Error(err)
	}

	var count int
	err = poolSqlDB.QueryRow(ctx, "SELECT COUNT(*) FROM pool.transaction WHERE status = $1", pool.TxStatusInvalid).Scan(&count)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, count)
}

func Test_UpdateTxStatus(t *testing.T) {
	ctx := context.Background()

	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := setupPool(t, s, st, chainID.Uint64(), ctx)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)
	tx := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx, ""); err != nil {
		t.Error(err)
	}

	err = p.UpdateTxStatus(ctx, signedTx.Hash(), pool.TxStatusInvalid, false)
	if err != nil {
		t.Error(err)
	}

	rows, err := poolSqlDB.Query(ctx, "SELECT status FROM pool.transaction WHERE hash = $1", signedTx.Hash().Hex())
	defer rows.Close() // nolint:staticcheck
	if err != nil {
		t.Error(err)
	}

	var state string
	rows.Next()
	if err := rows.Scan(&state); err != nil {
		t.Error(err)
	}

	assert.Equal(t, pool.TxStatusInvalid, pool.TxStatus(state))
}

func Test_SetAndGetGasPrice(t *testing.T) {
	initOrResetDB()

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := pool.NewPool(cfg, s, nil, common.Address{}, chainID.Uint64())

	nBig, err := rand.Int(rand.Reader, big.NewInt(0).SetUint64(math.MaxUint64))
	if err != nil {
		t.Error(err)
	}
	expectedGasPrice := nBig.Uint64()

	ctx := context.Background()

	err = p.SetGasPrice(ctx, expectedGasPrice)
	if err != nil {
		t.Error(err)
	}

	gasPrice, err := p.GetGasPrice(ctx)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expectedGasPrice, gasPrice)
}

func TestGetPendingTxSince(t *testing.T) {
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := setupPool(t, s, st, chainID.Uint64(), ctx)

	const txsCount = 10

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	txsAddedHashes := []common.Hash{}
	txsAddedTime := []time.Time{}

	timeBeforeTxs := time.Now()
	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		txsAddedTime = append(txsAddedTime, time.Now())
		if err := p.AddTx(ctx, *signedTx, ""); err != nil {
			t.Error(err)
		}
		txsAddedHashes = append(txsAddedHashes, signedTx.Hash())
		time.Sleep(1 * time.Second)
	}

	txHashes, err := p.GetPendingTxHashesSince(ctx, timeBeforeTxs)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, txsCount, len(txHashes))
	for i, txHash := range txHashes {
		assert.Equal(t, txHash.Hex(), txsAddedHashes[i].Hex())
	}

	txHashes, err = p.GetPendingTxHashesSince(ctx, txsAddedTime[5])
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 5, len(txHashes))
	assert.Equal(t, txHashes[0].Hex(), txsAddedHashes[5].Hex())
	assert.Equal(t, txHashes[1].Hex(), txsAddedHashes[6].Hex())
	assert.Equal(t, txHashes[2].Hex(), txsAddedHashes[7].Hex())
	assert.Equal(t, txHashes[3].Hex(), txsAddedHashes[8].Hex())
	assert.Equal(t, txHashes[4].Hex(), txsAddedHashes[9].Hex())

	txHashes, err = p.GetPendingTxHashesSince(ctx, txsAddedTime[8])
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, len(txHashes))
	assert.Equal(t, txHashes[0].Hex(), txsAddedHashes[8].Hex())
	assert.Equal(t, txHashes[1].Hex(), txsAddedHashes[9].Hex())

	txHashes, err = p.GetPendingTxHashesSince(ctx, txsAddedTime[9])
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(txHashes))
	assert.Equal(t, txHashes[0].Hex(), txsAddedHashes[9].Hex())

	txHashes, err = p.GetPendingTxHashesSince(ctx, txsAddedTime[9].Add(1*time.Second))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 0, len(txHashes))
}

func Test_DeleteTransactionsByHashes(t *testing.T) {
	ctx := context.Background()
	initOrResetDB()
	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)

	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := setupPool(t, s, st, chainID.Uint64(), ctx)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx1, ""); err != nil {
		t.Error(err)
	}

	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), gasLimit, gasPrice, []byte{})
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx2, ""); err != nil {
		t.Error(err)
	}

	err = p.DeleteTransactionsByHashes(ctx, []common.Hash{signedTx1.Hash(), signedTx2.Hash()})
	if err != nil {
		t.Error(err)
	}

	var count int
	err = poolSqlDB.QueryRow(ctx, "SELECT COUNT(*) FROM pool.transaction").Scan(&count)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 0, count)
}

func Test_TryAddIncompatibleTxs(t *testing.T) {
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	poolSqlDB, err := db.NewSQLDB(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer poolSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	initialBalance, _ := big.NewInt(0).SetString(encoding.MaxUint256StrNumber, encoding.Base10)
	initialBalance = initialBalance.Add(initialBalance, initialBalance)
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: operations.DefaultSequencerAddress,
				Type:    int(merkletree.LeafTypeBalance),
				Value:   initialBalance.String(),
			},
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}

	type testCase struct {
		name                 string
		createIncompatibleTx func() types.Transaction
		expectedError        error
	}

	auth := operations.MustGetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(t, err)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(operations.DefaultSequencerPrivateKey, "0x"))
	require.NoError(t, err)

	chainIdOver64Bits := big.NewInt(0).SetUint64(math.MaxUint64)
	chainIdOver64Bits = chainIdOver64Bits.Add(chainIdOver64Bits, big.NewInt(1))
	authChainIdOver64Bits, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIdOver64Bits)
	require.NoError(t, err)

	bigIntOver256Bits, _ := big.NewInt(0).SetString(encoding.MaxUint256StrNumber, encoding.Base10)
	bigIntOver256Bits = bigIntOver256Bits.Add(bigIntOver256Bits, big.NewInt(1))

	testCases := []testCase{
		{
			name: "Gas price over 256 bits",
			createIncompatibleTx: func() types.Transaction {
				tx := types.NewTransaction(uint64(0),
					common.HexToAddress("0x1"),
					big.NewInt(1), gasLimit, bigIntOver256Bits, nil)
				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(t, err)
				return *signedTx
			},
			expectedError: pool.ErrInsufficientFunds,
		},
		{
			name: "Value over 256 bits",
			createIncompatibleTx: func() types.Transaction {
				tx := types.NewTransaction(uint64(0),
					common.HexToAddress("0x1"),
					bigIntOver256Bits, gasLimit, gasPrice, nil)
				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(t, err)
				return *signedTx
			},
			expectedError: pool.ErrInsufficientFunds,
		},
		{
			name: "data over 30k bytes",
			createIncompatibleTx: func() types.Transaction {
				data := [30001]byte{}
				tx := types.NewTransaction(uint64(0),
					common.HexToAddress("0x1"),
					big.NewInt(1), 141004, gasPrice, data[:])
				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(t, err)
				return *signedTx
			},
			expectedError: fmt.Errorf("data size bigger than allowed, current size is %v bytes and max allowed is %v bytes", 30001, 30000),
		},
		{
			name: "chain id over 64 bits",
			createIncompatibleTx: func() types.Transaction {
				tx := types.NewTransaction(uint64(0),
					common.HexToAddress("0x1"),
					big.NewInt(1), gasLimit, gasPrice, nil)
				signedTx, err := authChainIdOver64Bits.Signer(authChainIdOver64Bits.From, tx)
				require.NoError(t, err)
				return *signedTx
			},
			expectedError: fmt.Errorf("chain id higher than allowed, max allowed is %v", uint64(math.MaxUint64)),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			incompatibleTx := testCase.createIncompatibleTx()
			p := setupPool(t, s, st, incompatibleTx.ChainId().Uint64(), ctx)
			err = p.AddTx(ctx, incompatibleTx, "")
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func newState(sqlDB *pgxpool.Pool) *state.State {
	ctx := context.Background()
	stateDb := state.NewPostgresStorage(sqlDB)
	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "localhost")

	executorServerConfig := executor.Config{URI: fmt.Sprintf("%s:50071", zkProverURI)}
	mtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:50061", zkProverURI)}
	executorClient, _, _ := executor.NewExecutorClient(ctx, executorServerConfig)
	stateDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
	stateTree := merkletree.NewStateTree(stateDBClient)
	st := state.NewState(state.Config{MaxCumulativeGasUsed: 800000, ChainID: chainID.Uint64(), ForkIDIntervals: []state.ForkIDInterval{{
		FromBatchNumber: 0,
		ToBatchNumber:   math.MaxUint64,
		ForkId:          0,
		Version:         "",
	}}}, stateDb, executorClient, stateTree)
	return st
}

func initOrResetDB() {
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}
	if err := dbutils.InitOrResetPool(poolDBCfg); err != nil {
		panic(err)
	}
}

func Test_AddTxWithIntrinsicGasTooLow(t *testing.T) {
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := setupPool(t, s, st, chainID.Uint64(), ctx)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert transaction
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    uint64(0),
		To:       &common.Address{},
		Value:    big.NewInt(0),
		Gas:      0,
		GasPrice: gasPrice,
		Data:     []byte{},
	})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, "")
	require.Error(t, err)
	assert.Equal(t, err.Error(), pool.ErrIntrinsicGas.Error())

	tx = types.NewTx(&types.LegacyTx{
		Nonce:    uint64(0),
		To:       nil,
		Value:    big.NewInt(10),
		Gas:      0,
		GasPrice: gasPrice,
		Data:     []byte{},
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, "")
	require.Error(t, err)
	assert.Equal(t, err.Error(), pool.ErrIntrinsicGas.Error())

	tx = types.NewTx(&types.LegacyTx{
		Nonce:    uint64(0),
		To:       &common.Address{},
		Value:    big.NewInt(10),
		Gas:      uint64(21000),
		GasPrice: gasPrice,
		Data:     []byte{},
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, "")
	require.NoError(t, err)

	tx = types.NewTx(&types.LegacyTx{
		Nonce:    uint64(1),
		To:       &common.Address{},
		Value:    big.NewInt(10),
		Gas:      0,
		GasPrice: gasPrice,
		Data:     []byte("data inside tx"),
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, "")
	require.Error(t, err)
	assert.Equal(t, err.Error(), pool.ErrIntrinsicGas.Error())

	tx = types.NewTx(&types.LegacyTx{
		Nonce:    uint64(1),
		To:       &common.Address{},
		Value:    big.NewInt(10),
		Gas:      uint64(21223),
		GasPrice: gasPrice,
		Data:     []byte("data inside tx"),
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, "")
	require.Error(t, err)
	assert.Equal(t, err.Error(), pool.ErrIntrinsicGas.Error())

	tx = types.NewTx(&types.LegacyTx{
		Nonce:    uint64(1),
		To:       &common.Address{},
		Value:    big.NewInt(10),
		Gas:      uint64(21224),
		GasPrice: gasPrice,
		Data:     []byte("data inside tx"),
	})
	signedTx, err = auth.Signer(auth.From, tx)
	require.NoError(t, err)
	err = p.AddTx(ctx, *signedTx, "")
	require.NoError(t, err)

	txs, err := p.GetPendingTxs(ctx, false, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(txs))

	for i := 0; i < 2; i++ {
		assert.Equal(t, pool.TxStatusPending, txs[0].Status)
	}
}

func Test_AddTx_GasPriceErr(t *testing.T) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	bridgeSC, err := bridge.NewPolygonzkevmbridge(l2BridgeAddr, nil)
	require.NoError(t, err)

	auth.NoSend = true
	auth.GasLimit = 53000
	auth.GasPrice = big.NewInt(0)
	auth.Nonce = big.NewInt(0)

	signedTx, err := bridgeSC.ClaimAsset(auth, [32][32]byte{}, uint32(123456789), [32]byte{}, [32]byte{}, 69, common.Address{}, uint32(20), common.Address{}, big.NewInt(0), []byte{})
	require.NoError(t, err)

	claimData := signedTx.Data()

	require.NoError(t, err)
	testCases := []struct {
		name          string
		nonce         uint64
		to            *common.Address
		gasLimit      uint64
		gasPrice      *big.Int
		data          []byte
		expectedError error
	}{
		{
			name:          "GasPriceTooLowErr",
			nonce:         0,
			to:            nil,
			gasLimit:      gasLimit,
			gasPrice:      big.NewInt(0).SetUint64(gasPrice.Uint64() - uint64(1)),
			data:          []byte{},
			expectedError: pool.ErrGasPrice,
		},
		{
			name:          "NoGasPriceTooLowErr_ForClaims",
			nonce:         1,
			to:            &l2BridgeAddr,
			gasLimit:      cfg.FreeClaimGasLimit,
			gasPrice:      big.NewInt(0),
			data:          claimData,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initOrResetDB()

			stateSqlDB, err := db.NewSQLDB(stateDBCfg)
			if err != nil {
				panic(err)
			}
			defer stateSqlDB.Close() //nolint:gosec,errcheck

			poolSqlDB, err := db.NewSQLDB(poolDBCfg)
			if err != nil {
				t.Error(err)
			}
			defer poolSqlDB.Close() //nolint:gosec,errcheck

			st := newState(stateSqlDB)

			genesisBlock := state.Block{
				BlockNumber: 0,
				BlockHash:   state.ZeroHash,
				ParentHash:  state.ZeroHash,
				ReceivedAt:  time.Now(),
			}
			genesis := state.Genesis{
				Actions: []*state.GenesisAction{
					{
						Address: senderAddress,
						Type:    int(merkletree.LeafTypeBalance),
						Value:   "1000000000000000000000",
					},
				},
			}
			ctx := context.Background()
			dbTx, err := st.BeginStateTransaction(ctx)
			require.NoError(t, err)
			_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
			require.NoError(t, err)
			require.NoError(t, dbTx.Commit(ctx))

			s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
			if err != nil {
				t.Error(err)
			}

			const chainID = 2576980377
			p := setupPool(t, s, st, chainID, ctx)
			tx := types.NewTx(&types.LegacyTx{
				Nonce:    tc.nonce,
				To:       tc.to,
				Value:    big.NewInt(0),
				Gas:      tc.gasLimit,
				GasPrice: tc.gasPrice,
				Data:     tc.data,
			})
			privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
			require.NoError(t, err)

			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
			require.NoError(t, err)

			signedTx, err := auth.Signer(auth.From, tx)
			require.NoError(t, err)

			err = p.AddTx(ctx, *signedTx, "")
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_AddRevertedTx(t *testing.T) {
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := setupPool(t, s, st, chainID.Uint64(), ctx)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert transaction
	revertScData, err := hex.DecodeHex(Revert.RevertBin)
	require.NoError(t, err)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    uint64(0),
		Gas:      uint64(1000000),
		GasPrice: gasPrice,
		Data:     revertScData,
	})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, "")
	require.NoError(t, err)

	txs, err := p.GetPendingTxs(ctx, false, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(txs))

	for i := 0; i < 1; i++ {
		assert.Equal(t, pool.TxStatusPending, txs[0].Status)
	}
}

func setupPool(t *testing.T, s *pgpoolstorage.PostgresPoolStorage, st *state.State, chainID uint64, ctx context.Context) *pool.Pool {
	p := pool.NewPool(cfg, s, st, l2BridgeAddr, chainID)
	err := p.SetGasPrice(ctx, gasPrice.Uint64())
	if err != nil {
		t.Error(err)
	}
	p.StartPollingMinSuggestedGasPrice(ctx)
	return p
}

func Test_AvoidDuplicatedClaims(t *testing.T) {
	initOrResetDB()

	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		t.Error(err)
	}
	defer stateSqlDB.Close() //nolint:gosec,errcheck

	st := newState(stateSqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)

	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}
	p := setupPool(t, s, st, chainID.Uint64(), ctx)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert transaction
	bridgeSC, err := bridge.NewPolygonzkevmbridge(l2BridgeAddr, nil)
	require.NoError(t, err)

	auth.NoSend = true
	auth.GasLimit = 53000
	auth.GasPrice = big.NewInt(0)
	auth.Nonce = big.NewInt(0)

	signedTx, err := bridgeSC.ClaimAsset(auth, [32][32]byte{}, uint32(123456789), [32]byte{}, [32]byte{}, 69, common.Address{}, uint32(20), common.Address{}, big.NewInt(0), []byte{})
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, "")
	require.NoError(t, err)

	auth.Nonce = big.NewInt(1)
	signedTx, err = bridgeSC.ClaimAsset(auth, [32][32]byte{}, uint32(123456789), [32]byte{}, [32]byte{}, 69, common.Address{}, uint32(20), common.Address{}, big.NewInt(0), []byte{})
	require.NoError(t, err)

	err = p.AddTx(ctx, *signedTx, "")
	require.Equal(t, err.Error(), "deposit count already exists")
}
