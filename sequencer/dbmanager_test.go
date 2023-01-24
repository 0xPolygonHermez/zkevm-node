package sequencer

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	mtDBclientpb "github.com/0xPolygonHermez/zkevm-node/merkletree/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	executorclientpb "github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var (
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
	executorClient    executorclientpb.ExecutorServiceClient
	mtDBServiceClient mtDBclientpb.StateDBServiceClient
	mtDBClientConn    *grpc.ClientConn
	testDbManager     *dbManager
)

func TestMain(m *testing.M) {
	initOrResetDB()
	ctx = context.Background()

	stateDb, err = db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateDb.Close()

	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "34.245.104.156")
	mtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:50061", zkProverURI)}
	var mtDBCancel context.CancelFunc
	mtDBServiceClient, mtDBClientConn, mtDBCancel = merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
	s := mtDBClientConn.GetState()
	log.Infof("stateDbClientConn state: %s", s.String())
	defer func() {
		mtDBCancel()
		mtDBClientConn.Close()
	}()

	stateTree = merkletree.NewStateTree(mtDBServiceClient)
	testState = state.NewState(stateCfg, state.NewPostgresStorage(stateDb), executorClient, stateTree)

	// DBManager
	closingSignalCh := ClosingSignalCh{
		ForcedBatchCh:        make(chan state.ForcedBatch),
		GERCh:                make(chan common.Hash),
		L2ReorgCh:            make(chan L2ReorgEvent),
		SendingToL1TimeoutCh: make(chan bool),
	}

	txsStore := TxsStore{
		Ch: make(chan *txToStore),
		Wg: new(sync.WaitGroup),
	}

	batchConstraints := batchConstraints{
		MaxTxsPerBatch:       150,
		MaxBatchBytesSize:    150000,
		MaxCumulativeGasUsed: 30000000,
		MaxKeccakHashes:      468,
		MaxPoseidonHashes:    279620,
		MaxPoseidonPaddings:  149796,
		MaxMemAligns:         262144,
		MaxArithmetics:       262144,
		MaxBinaries:          262144,
		MaxSteps:             8388608,
	}

	testDbManager = newDBManager(ctx, nil, testState, nil, closingSignalCh, txsStore, batchConstraints)

	result := m.Run()
	os.Exit(result)
}

func initOrResetDB() {
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}
}

func TestOpenBatch(t *testing.T) {
	initOrResetDB()
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
}

func TestGetLastBatchNumber(t *testing.T) {
	initOrResetDB()

	TestOpenBatch(t)

	lastBatchNum, err := testDbManager.GetLastBatchNumber(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), lastBatchNum)
}

func TestCreateFirstBatch(t *testing.T) {
	initOrResetDB()

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	_, err = testState.SetGenesis(ctx, state.Block{}, state.Genesis{}, dbTx)
	require.NoError(t, err)
	err = dbTx.Commit(ctx)
	require.NoError(t, err)

	processingContext := testDbManager.CreateFirstBatch(ctx, common.Address{})
	require.Equal(t, uint64(1), processingContext.BatchNumber)
}
