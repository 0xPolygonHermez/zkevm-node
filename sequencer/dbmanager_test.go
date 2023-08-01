package sequencer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	mtDBclientpb "github.com/0xPolygonHermez/zkevm-node/merkletree/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var (
	mtDBCancel context.CancelFunc
	ctx        context.Context
	testState  *state.State
	stateTree  *merkletree.StateTree
	stateDb    *pgxpool.Pool
	err        error
	stateDBCfg = dbutils.NewStateConfigFromEnv()
	stateCfg   = state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
	}
	dbManagerCfg      = DBManagerCfg{PoolRetrievalInterval: types.NewDuration(500 * time.Millisecond)}
	executorClient    executor.ExecutorServiceClient
	mtDBServiceClient mtDBclientpb.HashDBServiceClient
	mtDBClientConn    *grpc.ClientConn
	testDbManager     *dbManager
)

func setupDBManager() {
	initOrResetDB()
	ctx = context.Background()

	stateDb, err = db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}

	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "localhost")
	mtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:50061", zkProverURI)}

	mtDBServiceClient, mtDBClientConn, mtDBCancel = merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
	s := mtDBClientConn.GetState()
	log.Infof("stateDbClientConn state: %s", s.String())

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		panic(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	stateTree = merkletree.NewStateTree(mtDBServiceClient)
	testState = state.NewState(stateCfg, state.NewPostgresStorage(stateDb), executorClient, stateTree, eventLog)

	// DBManager
	closingSignalCh := ClosingSignalCh{
		ForcedBatchCh: make(chan state.ForcedBatch),
		GERCh:         make(chan common.Hash),
		L2ReorgCh:     make(chan L2ReorgEvent),
	}
	batchConstraints := batchConstraints{
		MaxTxsPerBatch:       300,
		MaxBatchBytesSize:    120000,
		MaxCumulativeGasUsed: 30000000,
		MaxKeccakHashes:      2145,
		MaxPoseidonHashes:    252357,
		MaxPoseidonPaddings:  135191,
		MaxMemAligns:         236585,
		MaxArithmetics:       236585,
		MaxBinaries:          473170,
		MaxSteps:             7570538,
	}

	testDbManager = newDBManager(ctx, dbManagerCfg, nil, testState, nil, closingSignalCh, batchConstraints)
}

func initOrResetDB() {
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}
}

func cleanupDBManager() {
	mtDBCancel()
	mtDBClientConn.Close()
}

func TestOpenBatch(t *testing.T) {
	setupDBManager()
	defer stateDb.Close()

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	_, err = testState.SetGenesis(ctx, state.Block{}, state.Genesis{}, dbTx)
	require.NoError(t, err)

	processingContext := state.ProcessingContext{
		BatchNumber:    1,
		Coinbase:       common.Address{},
		Timestamp:      time.Now().UTC(),
		GlobalExitRoot: common.Hash{},
	}

	err = testDbManager.OpenBatch(ctx, processingContext, dbTx)
	require.NoError(t, err)
	err = dbTx.Commit(ctx)
	require.NoError(t, err)
	cleanupDBManager()
}

func TestGetLastBatchNumber(t *testing.T) {
	setupDBManager()
	defer stateDb.Close()

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	_, err = testState.SetGenesis(ctx, state.Block{}, state.Genesis{}, dbTx)
	require.NoError(t, err)

	processingContext := state.ProcessingContext{
		BatchNumber:    1,
		Coinbase:       common.Address{},
		Timestamp:      time.Now().UTC(),
		GlobalExitRoot: common.Hash{},
	}

	err = testDbManager.OpenBatch(ctx, processingContext, dbTx)
	require.NoError(t, err)
	err = dbTx.Commit(ctx)
	require.NoError(t, err)

	lastBatchNum, err := testDbManager.GetLastBatchNumber(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), lastBatchNum)
	cleanupDBManager()
}

func TestCreateFirstBatch(t *testing.T) {
	setupDBManager()
	defer stateDb.Close()

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = testState.SetGenesis(ctx, state.Block{}, state.Genesis{}, dbTx)
	require.NoError(t, err)
	err = dbTx.Commit(ctx)
	require.NoError(t, err)

	processingContext := testDbManager.CreateFirstBatch(ctx, common.Address{})
	require.Equal(t, uint64(1), processingContext.BatchNumber)
	cleanupDBManager()
}
