package pool_test

import (
	"context"
	"crypto/rand"
	"math"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/pool/pgpoolstorage"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/pgstatestorage"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/stretchr/testify/assert"
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

	store := tree.NewPostgresStore(sqlDB)
	mt := tree.NewMerkleTree(store, 4, poseidon.Hash)
	scCodeStore := tree.NewPostgresSCCodeStore(sqlDB)
	tr := tree.NewStateTree(mt, scCodeStore)

	stateCfg := state.Config{
		DefaultChainID:       1000,
		MaxCumulativeGasUsed: 800000,
	}

	stateDB := pgstatestorage.NewPostgresStorage(sqlDB)
	st := state.NewState(stateCfg, stateDB, tr)

	s, err := pgpoolstorage.NewPostgresPoolStorage(cfg)
	if err != nil {
		t.Error(err)
	}

	st.SetGenesis(context.Background(), state.Genesis{})

	p := pool.NewPool(s, st)

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

	s, err := pgpoolstorage.NewPostgresPoolStorage(cfg)
	if err != nil {
		t.Error(err)
	}

	p := pool.NewPool(s, nil)

	const txsCount = 10

	ctx := context.Background()

	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
		err := p.AddTx(ctx, *tx)
		if err != nil {
			t.Error(err)
		}
	}

	txs, err := p.GetPendingTxs(ctx)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, txsCount, len(txs))

	for i := 0; i < txsCount; i++ {
		assert.Equal(t, pool.TxStatePending, txs[0].State)
	}
}

func Test_UpdateTxsState(t *testing.T) {
	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	ctx := context.Background()

	s, err := pgpoolstorage.NewPostgresPoolStorage(cfg)
	if err != nil {
		t.Error(err)
	}

	p := pool.NewPool(s, nil)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	err = p.AddTx(ctx, *tx1)
	if err != nil {
		t.Error(err)
	}

	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	err = p.AddTx(ctx, *tx2)
	if err != nil {
		t.Error(err)
	}

	err = p.UpdateTxsState(ctx, []common.Hash{tx1.Hash(), tx2.Hash()}, pool.TxStateInvalid)
	if err != nil {
		t.Error(err)
	}

	sqlDB, err := db.NewSQLDB(cfg)
	if err != nil {
		t.Error(err)
	}
	defer sqlDB.Close()

	var count int
	err = sqlDB.QueryRow(ctx, "SELECT COUNT(*) FROM pool.txs WHERE state = $1", pool.TxStateInvalid).Scan(&count)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, count)
}

func Test_UpdateTxState(t *testing.T) {
	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	ctx := context.Background()

	s, err := pgpoolstorage.NewPostgresPoolStorage(cfg)
	if err != nil {
		t.Error(err)
	}

	p := pool.NewPool(s, nil)

	tx := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	err = p.AddTx(ctx, *tx)
	if err != nil {
		t.Error(err)
	}

	err = p.UpdateTxState(ctx, tx.Hash(), pool.TxStateInvalid)
	if err != nil {
		t.Error(err)
	}

	sqlDB, err := db.NewSQLDB(cfg)
	if err != nil {
		t.Error(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	rows, err := sqlDB.Query(ctx, "SELECT state FROM pool.txs WHERE hash = $1", tx.Hash().Hex())
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()

	var state string
	rows.Next()
	err = rows.Scan(&state)
	if err != nil {
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

	p := pool.NewPool(s, nil)

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
