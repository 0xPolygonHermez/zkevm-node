package state_test

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/ci/vectors"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/hashdb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/test"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var (
	testState  *state.State
	stateTree  *merkletree.StateTree
	stateDb    *pgxpool.Pool
	err        error
	stateDBCfg = dbutils.NewStateConfigFromEnv()
	ctx        = context.Background()
	stateCfg   = state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          state.FORKID_ETROG,
			Version:         "",
		}},
	}
	forkID                             uint64 = state.FORKID_ETROG
	executorClient                     executor.ExecutorServiceClient
	mtDBServiceClient                  hashdb.HashDBServiceClient
	executorClientConn, mtDBClientConn *grpc.ClientConn
	batchResources                     = state.BatchResources{
		ZKCounters: state.ZKCounters{
			UsedKeccakHashes: 1,
		},
		Bytes: 1,
	}
	closingReason = state.GlobalExitRootDeadlineClosingReason
	genesis       = state.Genesis{
		FirstBatchData: &state.BatchData{
			Transactions:   "0xf8c380808401c9c380942a3dd3eb832af982ec71669e178424b10dca2ede80b8a4d3476afe000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a40d5f56745a118d0906a34e69aec8c0db1cb8fa000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005ca1ab1e0000000000000000000000000000000000000000000000000000000005ca1ab1e1bff",
			GlobalExitRoot: common.Hash{},
			Timestamp:      1697640780,
			Sequencer:      common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		},
	}
)

func TestMain(m *testing.M) {
	test.InitOrResetDB(stateDBCfg)

	stateDb, err = db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateDb.Close()

	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "zkevm-prover")

	executorServerConfig := executor.Config{URI: fmt.Sprintf("%s:50071", zkProverURI), MaxGRPCMessageSize: 100000000}
	var executorCancel context.CancelFunc
	executorClient, executorClientConn, executorCancel = executor.NewExecutorClient(ctx, executorServerConfig)
	s := executorClientConn.GetState()
	log.Infof("executorClientConn state: %s", s.String())
	defer func() {
		executorCancel()
		executorClientConn.Close()
	}()

	mtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:50061", zkProverURI)}
	var mtDBCancel context.CancelFunc
	mtDBServiceClient, mtDBClientConn, mtDBCancel = merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
	s = mtDBClientConn.GetState()
	log.Infof("stateDbClientConn state: %s", s.String())
	defer func() {
		mtDBCancel()
		mtDBClientConn.Close()
	}()

	stateTree = merkletree.NewStateTree(mtDBServiceClient)

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		panic(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)
	mt, err := l1infotree.NewL1InfoTree(32, [][32]byte{})
	if err != nil {
		panic(err)
	}
	testState = state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateCfg, stateDb), executorClient, stateTree, eventLog, mt)

	result := m.Run()

	os.Exit(result)
}

// TestStateTransition tests state transitions using the vector
func TestStateTransition(t *testing.T) {
	ctx := context.Background()
	// Load test vectors
	testCases, err := vectors.LoadStateTransitionTestCasesEtrog("../../../test/vectors/src/state-transition/etrog/balances.json")
	require.NoError(t, err)

	// Run test cases
	for i, testCase := range testCases {
		block := state.Block{
			BlockNumber: uint64(i + 1),
			BlockHash:   state.ZeroHash,
			ParentHash:  state.ZeroHash,
			ReceivedAt:  time.Now(),
		}

		genesisActions := vectors.GenerateGenesisActionsEtrog(testCase.Genesis)

		dbTx, err := testState.BeginStateTransaction(ctx)
		require.NoError(t, err)

		stateRoot, err := testState.SetGenesis(ctx, block, state.Genesis{Actions: genesisActions}, metrics.SynchronizerCallerLabel, dbTx)
		require.NoError(t, err)
		require.Equal(t, testCase.ExpectedOldStateRoot, stateRoot.String())
		err = dbTx.Rollback(ctx)
		require.NoError(t, err)

		// convert vector txs
		txs := make([]state.L2TxRaw, 0, len(testCase.Txs))
		for i := 0; i < len(testCase.Txs); i++ {
			vecTx := testCase.Txs[i]
			if vecTx.Type != 0x0b {
				tx, err := state.DecodeTx(vecTx.RawTx)
				require.NoError(t, err)
				l2Tx := state.L2TxRaw{
					Tx:                   *tx,
					EfficiencyPercentage: 255,
				}
				txs = append(txs, l2Tx)
			}
		}

		timestampLimit, ok := big.NewInt(0).SetString(testCase.TimestampLimit, 10)
		require.True(t, ok)

		// Generate batchdata from the txs in the test and compared with the vector
		l2block := state.L2BlockRaw{
			DeltaTimestamp:  uint32(timestampLimit.Uint64()),
			IndexL1InfoTree: testCase.Txs[0].IndexL1InfoTree,
			Transactions:    txs,
		}

		batch := state.BatchRawV2{
			Blocks: []state.L2BlockRaw{l2block},
		}

		batchData, err := state.EncodeBatchV2(&batch)
		require.NoError(t, err)

		require.Equal(t, common.FromHex(testCase.BatchL2Data), batchData)

		processRequest := state.ProcessRequest{
			BatchNumber:       uint64(i + 1),
			L1InfoRoot_V2:     common.HexToHash(testCase.L1InfoRoot),
			OldStateRoot:      stateRoot,
			OldAccInputHash:   common.HexToHash(testCase.OldAccInputHash),
			Transactions:      common.FromHex(testCase.BatchL2Data),
			TimestampLimit_V2: timestampLimit.Uint64(),
			Coinbase:          common.HexToAddress(testCase.SequencerAddress),
			ForkID:            testCase.ForkID,
		}

		processResponse, _ := testState.ProcessBatchV2(ctx, processRequest, false)
		require.Nil(t, processResponse.ExecutorError)
		require.Equal(t, testCase.ExpectedNewStateRoot, processResponse.NewStateRoot.String())
	}
}
