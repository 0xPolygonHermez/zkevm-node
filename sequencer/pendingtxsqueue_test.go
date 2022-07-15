package sequencer_test

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var senderPrivateKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"

var dbCfg = dbutils.NewConfigFromEnv()

var queueCfg = sequencer.PendingTxsQueueConfig{
	TxPendingInQueueCheckingFrequency: cfgTypes.NewDuration(1 * time.Second),
	GetPendingTxsFrequency:            cfgTypes.NewDuration(1 * time.Second),
}

func TestQueue_AddAndPopTx(t *testing.T) {
	if err := dbutils.InitOrReset(dbCfg); err != nil {
		panic(err)
	}

	sqlDB, err := db.NewSQLDB(dbCfg)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	st := newState(sqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	balance, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	genesis := state.Genesis{
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): balance,
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(dbCfg)
	if err != nil {
		panic(err)
	}

	p := pool.NewPool(s, st, common.Address{})

	const txsCount = 10

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		panic(err)
	}

	pendQueue := sequencer.NewPendingTxsQueue(queueCfg, p)
	go pendQueue.KeepPendingTxsQueue(ctx)
	go pendQueue.CleanPendTxsChan(ctx)
	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(i)), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		if err != nil {
			panic(err)
		}
		if err := p.AddTx(ctx, *signedTx); err != nil {
			panic(err)
		}
	}
	tx := pendQueue.PopPendingTx()
	assert.Equal(t, uint64(19), tx.GasPrice().Uint64())
	assert.Equal(t, 9, pendQueue.GetPendingTxsQueueLength())
	tx = pendQueue.PopPendingTx()
	assert.Equal(t, uint64(18), tx.GasPrice().Uint64())
	assert.Equal(t, 8, pendQueue.GetPendingTxsQueueLength())

	newTx := types.NewTransaction(uint64(txsCount), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(txsCount)), []byte{})
	signedTx, err := auth.Signer(auth.From, newTx)
	if err != nil {
		panic(err)
	}
	if err := p.AddTx(ctx, *signedTx); err != nil {
		panic(err)
	}

	time.Sleep(queueCfg.TxPendingInQueueCheckingFrequency.Duration * 2)
	assert.Equal(t, 9, pendQueue.GetPendingTxsQueueLength())
}

func TestQueue_AddOneTx(t *testing.T) {
	if err := dbutils.InitOrReset(dbCfg); err != nil {
		panic(err)
	}

	sqlDB, err := db.NewSQLDB(dbCfg)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	st := newState(sqlDB)

	genesisBlock := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}
	balance, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	genesis := state.Genesis{
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): balance,
		},
	}
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	require.NoError(t, err)
	err = st.SetGenesis(ctx, genesisBlock, genesis, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	s, err := pgpoolstorage.NewPostgresPoolStorage(dbCfg)
	if err != nil {
		panic(err)
	}

	p := pool.NewPool(s, st, common.Address{})

	const txsCount = 1

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		panic(err)
	}

	pendQueue := sequencer.NewPendingTxsQueue(queueCfg, p)
	go pendQueue.KeepPendingTxsQueue(ctx)
	go pendQueue.CleanPendTxsChan(ctx)
	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(i)), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		if err != nil {
			panic(err)
		}
		if err := p.AddTx(ctx, *signedTx); err != nil {
			panic(err)
		}
	}
	tx := pendQueue.PopPendingTx()
	assert.Equal(t, uint64(10), tx.GasPrice().Uint64())
	assert.Equal(t, 0, pendQueue.GetPendingTxsQueueLength())
	tx = pendQueue.PopPendingTx()
	assert.Nil(t, tx)

	newTx := types.NewTransaction(uint64(txsCount), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(txsCount)), []byte{})
	signedTx, err := auth.Signer(auth.From, newTx)
	if err != nil {
		panic(err)
	}
	if err := p.AddTx(ctx, *signedTx); err != nil {
		panic(err)
	}
	time.Sleep(queueCfg.TxPendingInQueueCheckingFrequency.Duration * 2)
	assert.Equal(t, 1, pendQueue.GetPendingTxsQueueLength())
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
