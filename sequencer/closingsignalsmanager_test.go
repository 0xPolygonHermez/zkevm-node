package sequencer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const numberOfForcesBatches = 10

var (
	testGER     = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	testAddr    = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	testRawData = common.Hex2Bytes("0xee80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880801cee7e01dc62f69a12c3510c6d64de04ee6346d84b6a017f3e786c7d87f963e75d8cc91fa983cd6d9cf55fff80d73bd26cd333b0f098acc1e58edb1fd484ad731b")
)

func setupTest(t *testing.T) {
	initOrResetDB()
	ctx = context.Background()

	stateDb, err = db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}

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

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

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

	// Set genesis batch
	_, err = testState.SetGenesis(ctx, state.Block{}, state.Genesis{}, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))
}

func prepareForcedBatches(t *testing.T) {
	// Create block
	const sql = `INSERT INTO state.forced_batch (forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase, block_num) VALUES ($1, $2, $3, $4, $5, $6)`

	for x := 0; x < numberOfForcesBatches; x++ {
		forcedBatchNum := int64(x)
		_, err := testState.PostgresStorage.Exec(ctx, sql, forcedBatchNum, testGER.String(), time.Now(), testRawData, testAddr.String(), 0)
		assert.NoError(t, err)
	}
}

func TestClosingSignalsManager(t *testing.T) {
	setupTest(t)
	cumtomForcedBatchCh := make(chan state.ForcedBatch)
	closingSignalCh.ForcedBatchCh = cumtomForcedBatchCh

	prepareForcedBatches(t)
	closingSignalsManager := newClosingSignalsManager(ctx, testDbManager, closingSignalCh, cfg)
	closingSignalsManager.Start()

	newCtx, cancelFunc := context.WithTimeout(ctx, time.Second*3)
	defer cancelFunc()

	var fb *state.ForcedBatch

	for {
		select {
		case <-newCtx.Done():
			log.Infof("received context done, Err: %s", ctx.Err())
			return
		// Forced  batch ch
		case fb := <-closingSignalCh.ForcedBatchCh:
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
}
