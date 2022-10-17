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

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
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
)

var (
	stateDBCfg = dbutils.NewStateConfigFromEnv()
	poolDBCfg  = dbutils.NewPoolConfigFromEnv()
	genesis    = state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
				Type:    int(merkletree.LeafTypeBalance),
				Value:   "1000000000000000000000",
			},
		},
	}
	chainID = big.NewInt(1337)
)

func TestMain(m *testing.M) {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
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
	genesis := state.Genesis{
		Actions: []*state.GenesisAction{
			{
				Address: "0xb48cA794d49EeC406A5dD2c547717e37b5952a83",
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
	p := pool.NewPool(s, st, common.Address{}, chainID)

	txRLPHash := "0xf86e8212658082520894fd8b27a263e19f0e9592180e61f0f8c9dfeb1ff6880de0b6b3a764000080850133333355a01eac4c2defc7ed767ae36bbd02613c581b8fb87d0e4f579c9ee3a7cfdb16faa7a043ce30f43d952b9d034cf8f04fecb631192a5dbc7ee2a47f1f49c0d022a8849d"
	b, err := hex.DecodeHex(txRLPHash)
	if err != nil {
		t.Error(err)
	}
	tx := new(types.Transaction)
	tx.UnmarshalBinary(b) //nolint:gosec,errcheck

	err = p.AddTx(ctx, *tx)
	if err != nil {
		t.Error(err)
	}

	rows, err := poolSqlDB.Query(ctx, "SELECT hash, encoded, decoded, status FROM pool.txs")
	defer rows.Close() // nolint:staticcheck
	if err != nil {
		t.Error(err)
	}

	c := 0
	for rows.Next() {
		var hash, encoded, decoded, status string
		err := rows.Scan(&hash, &encoded, &decoded, &status)
		if err != nil {
			t.Error(err)
		}
		b, _ := tx.MarshalJSON()

		assert.Equal(t, "0xa3cff5abdf47d4feb8204a45c0a8c58fc9b9bb9b29c6588c1d206b746815e9cc", hash, "invalid hash")
		assert.Equal(t, txRLPHash, encoded, "invalid encoded")
		assert.JSONEq(t, string(b), decoded, "invalid decoded")
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

	p := pool.NewPool(s, st, common.Address{}, chainID.Uint64())

	const txsCount = 10
	const limit = 5

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		if err := p.AddTx(ctx, *signedTx); err != nil {
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

	p := pool.NewPool(s, st, common.Address{}, chainID.Uint64())

	const txsCount = 10
	const limit = 0

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		if err := p.AddTx(ctx, *signedTx); err != nil {
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

	p := pool.NewPool(s, st, common.Address{}, chainID.Uint64())

	const txsCount = 10

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(i)), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		if err := p.AddTx(ctx, *signedTx); err != nil {
			t.Error(err)
		}
	}

	txs, err := p.GetTxs(ctx, pool.TxStatusPending, false, 1, 10)
	require.NoError(t, err)
	// bcs it's sorted by nonce, tx with the lowest nonce is expected here
	assert.Equal(t, txs[0].Transaction.Nonce(), uint64(0))
}

func Test_GetTopFailedTxsByProfitabilityAndZkCounters(t *testing.T) {
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

	p := pool.NewPool(s, st, common.Address{}, chainID.Uint64())

	const txsCount = 10

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	txsHashes := make([]string, 0, txsCount)
	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(i)), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		if err := p.AddTx(ctx, *signedTx); err != nil {
			t.Error(err)
		}
		txsHashes = append(txsHashes, signedTx.Hash().String())
	}

	err = p.UpdateTxsStatus(ctx, txsHashes, pool.TxStatusFailed)
	require.NoError(t, err)
	err = p.IncrementFailedCounter(ctx, txsHashes[0:txsCount/2])
	require.NoError(t, err)
	txs, err := p.GetTxs(ctx, pool.TxStatusFailed, false, 1, 10)
	require.NoError(t, err)
	// bcs it's sorted by nonce, tx with the lowest nonce is expected here
	assert.Equal(t, txsCount, len(txs))
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

	p := pool.NewPool(s, st, common.Address{}, chainID.Uint64())

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx1); err != nil {
		t.Error(err)
	}

	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx2); err != nil {
		t.Error(err)
	}

	err = p.UpdateTxsStatus(ctx, []string{signedTx1.Hash().String(), signedTx2.Hash().String()}, pool.TxStatusInvalid)
	if err != nil {
		t.Error(err)
	}

	var count int
	err = poolSqlDB.QueryRow(ctx, "SELECT COUNT(*) FROM pool.txs WHERE status = $1", pool.TxStatusInvalid).Scan(&count)
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

	p := pool.NewPool(s, st, common.Address{}, chainID.Uint64())

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx); err != nil {
		t.Error(err)
	}

	err = p.UpdateTxStatus(ctx, signedTx.Hash(), pool.TxStatusInvalid)
	if err != nil {
		t.Error(err)
	}

	rows, err := poolSqlDB.Query(ctx, "SELECT status FROM pool.txs WHERE hash = $1", signedTx.Hash().Hex())
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

	p := pool.NewPool(s, nil, common.Address{}, chainID.Uint64())

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

func TestMarkReorgedTxsAsPending(t *testing.T) {
	initOrResetDB()
	ctx := context.Background()
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
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(poolDBCfg)
	if err != nil {
		t.Error(err)
	}

	p := pool.NewPool(s, st, common.Address{}, chainID.Uint64())

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx1); err != nil {
		t.Error(err)
	}

	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx2); err != nil {
		t.Error(err)
	}

	err = p.UpdateTxsStatus(ctx, []string{signedTx1.Hash().String(), signedTx2.Hash().String()}, pool.TxStatusSelected)
	if err != nil {
		t.Error(err)
	}

	err = p.MarkReorgedTxsAsPending(ctx)
	require.NoError(t, err)
	txs, err := p.GetPendingTxs(ctx, false, 100)
	require.NoError(t, err)
	require.Equal(t, signedTx1.Hash().Hex(), txs[1].Hash().Hex())
	require.Equal(t, signedTx2.Hash().Hex(), txs[0].Hash().Hex())
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

	p := pool.NewPool(s, st, common.Address{}, chainID.Uint64())

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
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		txsAddedTime = append(txsAddedTime, time.Now())
		if err := p.AddTx(ctx, *signedTx); err != nil {
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

func Test_DeleteTxsByHashes(t *testing.T) {
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

	p := pool.NewPool(s, st, common.Address{}, chainID.Uint64())

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx1); err != nil {
		t.Error(err)
	}

	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx2); err != nil {
		t.Error(err)
	}

	err = p.DeleteTxsByHashes(ctx, []common.Hash{signedTx1.Hash(), signedTx2.Hash()})
	if err != nil {
		t.Error(err)
	}

	var count int
	err = poolSqlDB.QueryRow(ctx, "SELECT COUNT(*) FROM pool.txs").Scan(&count)
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
					big.NewInt(1), uint64(1), bigIntOver256Bits, nil)
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
					bigIntOver256Bits, uint64(1), big.NewInt(1), nil)
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
					big.NewInt(1), uint64(1), big.NewInt(1), data[:])
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
					big.NewInt(1), uint64(1), big.NewInt(1), nil)
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
			p := pool.NewPool(s, st, common.Address{}, incompatibleTx.ChainId().Uint64())
			err = p.AddTx(ctx, incompatibleTx)
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
	st := state.NewState(state.Config{MaxCumulativeGasUsed: 800000}, stateDb, executorClient, stateTree)
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
