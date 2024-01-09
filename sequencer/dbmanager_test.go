package sequencer

/*import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/hashdb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

//TODO: Fix tests ETROG

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
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          5,
			Version:         "",
		}},
	}
	dbManagerCfg      = DBManagerCfg{LoadPoolTxsCheckInterval: types.NewDuration(500 * time.Millisecond)}
	executorClient    executor.ExecutorServiceClient
	mtDBServiceClient hashdb.HashDBServiceClient
	mtDBClientConn    *grpc.ClientConn
	testDbManager     *dbManager

	genesis = state.Genesis{
		FirstBatchData: &state.BatchData{
			Transactions:   "0xf8c380808401c9c380942a3dd3eb832af982ec71669e178424b10dca2ede80b8a4d3476afe000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a40d5f56745a118d0906a34e69aec8c0db1cb8fa000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005ca1ab1e0000000000000000000000000000000000000000000000000000000005ca1ab1e1bff",
			GlobalExitRoot: common.Hash{},
			Timestamp:      1697640780,
			Sequencer:      common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		},
	}
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
	executorServerConfig := executor.Config{URI: fmt.Sprintf("%s:50071", zkProverURI), MaxGRPCMessageSize: 100000000}

	executorClient, _, _ = executor.NewExecutorClient(ctx, executorServerConfig)

	mtDBServiceClient, mtDBClientConn, mtDBCancel = merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
	s := mtDBClientConn.GetState()
	log.Infof("stateDbClientConn state: %s", s.String())

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		panic(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	stateTree = merkletree.NewStateTree(mtDBServiceClient)
	testState = state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateCfg, stateDb), executorClient, stateTree, eventLog, nil)

	// DBManager
	closingSignalCh := ClosingSignalCh{
		ForcedBatchCh: make(chan state.ForcedBatch),
		GERCh:         make(chan common.Hash),
		L2ReorgCh:     make(chan L2ReorgEvent),
	}
	batchConstraints := state.BatchConstraintsCfg{
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
		MaxSHA256Hashes:            1596,
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

	_, err = testState.SetGenesis(ctx, state.Block{}, genesis, metrics.SynchronizerCallerLabel, dbTx)
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

	_, err = testState.SetGenesis(ctx, state.Block{}, genesis, metrics.SynchronizerCallerLabel, dbTx)
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

}
*/
