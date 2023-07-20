package sequencer

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const numberOfForcesBatches = 10

var (
	localStateDb                                 *pgxpool.Pool
	localTestDbManager                           *dbManager
	localCtx                                     context.Context
	localMtDBCancel, localExecutorCancel         context.CancelFunc
	localMtDBServiceClient                       mtDBclientpb.HashDBServiceClient
	localMtDBClientConn, localExecutorClientConn *grpc.ClientConn
	localState                                   *state.State
	localExecutorClient                          executor.ExecutorServiceClient
	testGER                                      = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	testAddr                                     = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	testRawData                                  = common.Hex2Bytes("0xee80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880801cee7e01dc62f69a12c3510c6d64de04ee6346d84b6a017f3e786c7d87f963e75d8cc91fa983cd6d9cf55fff80d73bd26cd333b0f098acc1e58edb1fd484ad731b")
)

type mocks struct {
	Etherman *EthermanMock
}

func setupTest(t *testing.T) {
	initOrResetDB()

	localCtx = context.Background()

	localStateDb, err = db.NewSQLDB(dbutils.NewStateConfigFromEnv())
	if err != nil {
		panic(err)
	}

	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "localhost")
	localMtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:50061", zkProverURI)}
	localExecutorServerConfig := executor.Config{URI: fmt.Sprintf("%s:50071", zkProverURI), MaxGRPCMessageSize: 100000000}

	localExecutorClient, localExecutorClientConn, localExecutorCancel = executor.NewExecutorClient(localCtx, localExecutorServerConfig)
	s := localExecutorClientConn.GetState()
	log.Infof("executorClientConn state: %s", s.String())

	localMtDBServiceClient, localMtDBClientConn, localMtDBCancel = merkletree.NewMTDBServiceClient(localCtx, localMtDBServerConfig)
	s = localMtDBClientConn.GetState()
	log.Infof("localStateDbClientConn state: %s", s.String())

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		panic(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	localStateTree := merkletree.NewStateTree(localMtDBServiceClient)
	localState = state.NewState(stateCfg, state.NewPostgresStorage(localStateDb), localExecutorClient, localStateTree, eventLog)

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

	localTestDbManager = newDBManager(localCtx, dbManagerCfg, nil, localState, nil, closingSignalCh, batchConstraints)

	// Set genesis batch
	dbTx, err := localState.BeginStateTransaction(localCtx)
	require.NoError(t, err)
	_, err = localState.SetGenesis(localCtx, state.Block{}, state.Genesis{}, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(localCtx))
}

func cleanup(t *testing.T) {
	localMtDBCancel()
	localMtDBClientConn.Close()
	localExecutorCancel()
	localExecutorClientConn.Close()
}

func prepareForcedBatches(t *testing.T) {
	// Create block
	const sql = `INSERT INTO state.forced_batch (forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase, block_num) VALUES ($1, $2, $3, $4, $5, $6)`

	for x := 0; x < numberOfForcesBatches; x++ {
		forcedBatchNum := int64(x)
		_, err := localState.PostgresStorage.Exec(localCtx, sql, forcedBatchNum, testGER.String(), time.Now(), testRawData, testAddr.String(), 0)
		assert.NoError(t, err)
	}
}

func TestClosingSignalsManager(t *testing.T) {
	m := mocks{
		Etherman: NewEthermanMock(t),
	}

	setupTest(t)
	channels := ClosingSignalCh{
		ForcedBatchCh: make(chan state.ForcedBatch),
	}

	prepareForcedBatches(t)
	closingSignalsManager := newClosingSignalsManager(localCtx, localTestDbManager, channels, cfg, m.Etherman)
	closingSignalsManager.Start()

	newCtx, cancelFunc := context.WithTimeout(localCtx, time.Second*3)
	defer cancelFunc()

	var fb *state.ForcedBatch

	for {
		select {
		case <-newCtx.Done():
			log.Infof("received context done, Err: %s", newCtx.Err())
			return
		// Forced  batch ch
		case fb := <-channels.ForcedBatchCh:
			log.Debug("Forced batch received", "forced batch", fb)
		}

		if fb != nil {
			break
		}
	}

	require.NotEqual(t, (*state.ForcedBatch)(nil), fb)
	require.Equal(t, nil, fb.BlockNumber)
	require.Equal(t, int64(1), fb.ForcedBatchNumber)
	require.Equal(t, testGER, fb.GlobalExitRoot)
	require.Equal(t, testAddr, fb.Sequencer)
	require.Equal(t, testRawData, fb.RawTxsData)

	cleanup(t)
}
