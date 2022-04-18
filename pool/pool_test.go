package pool_test

import (
	"context"
	"crypto/rand"
	"math"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/pool/pgpoolstorage"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	senderPrivateKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
)

var cfg = dbutils.NewConfigFromEnv()

func TestMain(m *testing.M) {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})

	code := m.Run()
	os.Exit(code)
}

func Test_AddTx(t *testing.T) {
	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	sqlDB, err := db.NewSQLDB(cfg)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	st := newState(sqlDB)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	balance, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	genesis := state.Genesis{
		Block: genesisBlock,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0xb48cA794d49EeC406A5dD2c547717e37b5952a83"): balance,
		},
	}
	err = st.SetGenesis(context.Background(), genesis)
	if err != nil {
		t.Error(err)
	}

	s, err := pgpoolstorage.NewPostgresPoolStorage(cfg)
	if err != nil {
		t.Error(err)
	}

	p := pool.NewPool(s, st, common.Address{})

	txRLPHash := "0xf86e8212658082520894fd8b27a263e19f0e9592180e61f0f8c9dfeb1ff6880de0b6b3a764000080850133333355a01eac4c2defc7ed767ae36bbd02613c581b8fb87d0e4f579c9ee3a7cfdb16faa7a043ce30f43d952b9d034cf8f04fecb631192a5dbc7ee2a47f1f49c0d022a8849d"
	b, err := hex.DecodeHex(txRLPHash)
	if err != nil {
		t.Error(err)
	}
	tx := new(types.Transaction)
	tx.UnmarshalBinary(b) //nolint:gosec,errcheck

	ctx := context.Background()

	err = p.AddTx(ctx, *tx)
	if err != nil {
		t.Error(err)
	}

	rows, err := sqlDB.Query(ctx, "SELECT hash, encoded, decoded, state FROM pool.txs")
	defer rows.Close()
	if err != nil {
		t.Error(err)
	}

	c := 0
	for rows.Next() {
		var hash, encoded, decoded, state string
		err := rows.Scan(&hash, &encoded, &decoded, &state)
		if err != nil {
			t.Error(err)
		}
		b, _ := tx.MarshalJSON()

		assert.Equal(t, "0xa3cff5abdf47d4feb8204a45c0a8c58fc9b9bb9b29c6588c1d206b746815e9cc", hash, "invalid hash")
		assert.Equal(t, txRLPHash, encoded, "invalid encoded")
		assert.JSONEq(t, string(b), decoded, "invalid decoded")
		assert.Equal(t, string(pool.TxStatePending), state, "invalid tx state")
		c++
	}

	assert.Equal(t, 1, c, "invalid number of txs in the pool")
}

func Test_GetPendingTxs(t *testing.T) {
	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	sqlDB, err := db.NewSQLDB(cfg)
	if err != nil {
		t.Error(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	st := newState(sqlDB)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	balance, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	genesis := state.Genesis{
		Block: genesisBlock,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): balance,
		},
	}
	err = st.SetGenesis(context.Background(), genesis)
	if err != nil {
		t.Error(err)
	}

	s, err := pgpoolstorage.NewPostgresPoolStorage(cfg)
	if err != nil {
		t.Error(err)
	}

	p := pool.NewPool(s, st, common.Address{})

	const txsCount = 10
	const limit = 5

	ctx := context.Background()

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
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
		assert.Equal(t, pool.TxStatePending, txs[0].State)
	}
}

func Test_GetPendingTxsZeroPassed(t *testing.T) {
	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	sqlDB, err := db.NewSQLDB(cfg)
	if err != nil {
		t.Error(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	st := newState(sqlDB)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	balance, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	genesis := state.Genesis{
		Block: genesisBlock,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): balance,
		},
	}
	err = st.SetGenesis(context.Background(), genesis)
	if err != nil {
		t.Error(err)
	}

	s, err := pgpoolstorage.NewPostgresPoolStorage(cfg)
	if err != nil {
		t.Error(err)
	}

	p := pool.NewPool(s, st, common.Address{})

	const txsCount = 10
	const limit = 0

	ctx := context.Background()

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
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
		assert.Equal(t, pool.TxStatePending, txs[0].State)
	}
}

func Test_UpdateTxsState(t *testing.T) {
	ctx := context.Background()

	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	sqlDB, err := db.NewSQLDB(cfg)
	if err != nil {
		t.Error(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	st := newState(sqlDB)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	balance, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	genesis := state.Genesis{
		Block: genesisBlock,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): balance,
		},
	}
	err = st.SetGenesis(context.Background(), genesis)
	if err != nil {
		t.Error(err)
	}

	s, err := pgpoolstorage.NewPostgresPoolStorage(cfg)
	if err != nil {
		t.Error(err)
	}

	p := pool.NewPool(s, st, common.Address{})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
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

	err = p.UpdateTxsState(ctx, []common.Hash{signedTx1.Hash(), signedTx2.Hash()}, pool.TxStateInvalid)
	if err != nil {
		t.Error(err)
	}

	var count int
	err = sqlDB.QueryRow(ctx, "SELECT COUNT(*) FROM pool.txs WHERE state = $1", pool.TxStateInvalid).Scan(&count)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, count)
}

func Test_UpdateTxState(t *testing.T) {
	ctx := context.Background()

	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	sqlDB, err := db.NewSQLDB(cfg)
	if err != nil {
		t.Error(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	st := newState(sqlDB)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	balance, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	genesis := state.Genesis{
		Block: genesisBlock,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): balance,
		},
	}
	err = st.SetGenesis(context.Background(), genesis)
	if err != nil {
		t.Error(err)
	}

	s, err := pgpoolstorage.NewPostgresPoolStorage(cfg)
	if err != nil {
		t.Error(err)
	}

	p := pool.NewPool(s, st, common.Address{})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)

	tx := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)
	if err := p.AddTx(ctx, *signedTx); err != nil {
		t.Error(err)
	}

	err = p.UpdateTxState(ctx, signedTx.Hash(), pool.TxStateInvalid)
	if err != nil {
		t.Error(err)
	}

	rows, err := sqlDB.Query(ctx, "SELECT state FROM pool.txs WHERE hash = $1", signedTx.Hash().Hex())
	defer rows.Close()
	if err != nil {
		t.Error(err)
	}

	var state string
	rows.Next()
	if err := rows.Scan(&state); err != nil {
		t.Error(err)
	}

	assert.Equal(t, pool.TxStateInvalid, pool.TxState(state))
}

func Test_SetAndGetGasPrice(t *testing.T) {
	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	s, err := pgpoolstorage.NewPostgresPoolStorage(cfg)
	if err != nil {
		t.Error(err)
	}

	p := pool.NewPool(s, nil, common.Address{})

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

func newState(sqlDB *pgxpool.Pool) *state.State {
	store := tree.NewPostgresStore(sqlDB)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(sqlDB)
	tr := tree.NewStateTree(mt, scCodeStore)

	stateCfg := state.Config{
		DefaultChainID:       1000,
		MaxCumulativeGasUsed: 800000,
	}

	stateDB := state.NewPostgresStorage(sqlDB)
	st := state.NewState(stateCfg, stateDB, tr)

	return st
}
