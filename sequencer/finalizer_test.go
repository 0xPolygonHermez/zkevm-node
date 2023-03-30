package sequencer

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	f             *finalizer
	nilErr        error
	dbManagerMock = new(DbManagerMock)
	executorMock  = new(StateMock)
	workerMock    = new(WorkerMock)
	dbTxMock      = new(DbTxMock)
	bc            = batchConstraints{
		MaxTxsPerBatch:       150,
		MaxBatchBytesSize:    129848,
		MaxCumulativeGasUsed: 30000000,
		MaxKeccakHashes:      468,
		MaxPoseidonHashes:    279620,
		MaxPoseidonPaddings:  149796,
		MaxMemAligns:         262144,
		MaxArithmetics:       262144,
		MaxBinaries:          262144,
		MaxSteps:             8388608,
	}
	txsStore = TxsStore{
		Ch: make(chan *txToStore, 1),
		Wg: new(sync.WaitGroup),
	}
	closingSignalCh = ClosingSignalCh{
		ForcedBatchCh:        make(chan state.ForcedBatch),
		GERCh:                make(chan common.Hash),
		L2ReorgCh:            make(chan L2ReorgEvent),
		SendingToL1TimeoutCh: make(chan bool),
	}
	cfg = FinalizerCfg{
		GERDeadlineTimeoutInSec: cfgTypes.Duration{
			Duration: 60,
		},
		ForcedBatchDeadlineTimeoutInSec: cfgTypes.Duration{
			Duration: 60,
		},
		SendingToL1DeadlineTimeoutInSec: cfgTypes.Duration{
			Duration: 60,
		},
		SleepDurationInMs: cfgTypes.Duration{
			Duration: 60,
		},
		ClosingSignalsManagerWaitForCheckingL1Timeout: cfgTypes.Duration{
			Duration: 10 * time.Second,
		},
		ClosingSignalsManagerWaitForCheckingGER: cfgTypes.Duration{
			Duration: 10 * time.Second,
		},
		ClosingSignalsManagerWaitForCheckingForcedBatches: cfgTypes.Duration{
			Duration: 10 * time.Second,
		},
		ResourcePercentageToCloseBatch: 10,
		GERFinalityNumberOfBlocks:      64,
	}
	seqAddr  = common.Address{}
	oldHash  = common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f2")
	newHash  = common.HexToHash("0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	sender   = common.HexToAddress("0x3445324")
	isSynced = func(ctx context.Context) bool {
		return true
	}
	// tx1 = ethTypes.NewTransaction(0, common.HexToAddress("0"), big.NewInt(0), 0, big.NewInt(0), []byte("aaa"))
	// tx2 = ethTypes.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))

	testErr = fmt.Errorf("some error")
	// testErr2         = fmt.Errorf("some error2")
	openBatchError   = fmt.Errorf("failed to open new batch, err: %w", testErr)
	cumulativeGasErr = state.GetZKCounterError("CumulativeGasUsed")
)

func testNow() time.Time {
	return time.Unix(0, 0)
}

func TestNewFinalizer(t *testing.T) {
	// arrange and act
	f = newFinalizer(cfg, workerMock, dbManagerMock, executorMock, seqAddr, isSynced, closingSignalCh, txsStore, bc)

	// assert
	assert.NotNil(t, f)
	assert.Equal(t, f.cfg, cfg)
	assert.Equal(t, f.worker, workerMock)
	assert.Equal(t, f.dbManager, dbManagerMock)
	assert.Equal(t, f.executor, executorMock)
	assert.Equal(t, f.sequencerAddress, seqAddr)
	assert.Equal(t, f.closingSignalCh, closingSignalCh)
	assert.Equal(t, f.txsStore, txsStore)
	assert.Equal(t, f.batchConstraints, bc)
}

//
//func TestFinalizer_newWIPBatch(t *testing.T) {
//	// arrange
//	f = setupFinalizer(true)
//	now = testNow
//	defer func() {
//		now = time.Now
//	}()
//
//	txs := make([]types.Transaction, 0)
//	batchNum := f.batch.batchNumber + 1
//	f.processRequest.GlobalExitRoot = oldHash
//	f.processRequest.OldStateRoot = oldHash
//	f.processRequest.Transactions = []byte{}
//	expectedWipBatch := &WipBatch{
//		batchNumber:        batchNum,
//		coinbase:           f.sequencerAddress,
//		initialStateRoot:   newHash,
//		stateRoot:          newHash,
//		timestamp:          uint64(now().Unix()),
//		remainingResources: getMaxRemainingResources(f.batchConstraints),
//	}
//	closeBatchParams := ClosingBatchParameters{
//		BatchNumber:   f.batch.batchNumber,
//		StateRoot:     f.batch.stateRoot,
//		LocalExitRoot: f.batch.localExitRoot,
//		Txs:           txs,
//	}
//	batches := []*state.Batch{
//		{
//			BatchNumber: 0,
//			StateRoot:   oldHash,
//		},
//	}
//	testCases := []struct {
//		name             string
//		batches          []*state.Batch
//		closeBatchErr    error
//		closeBatchParams ClosingBatchParameters
//		openBatchErr     error
//		expectedWip      *WipBatch
//		expectedErr      error
//	}{
//		{
//			name:             "Success",
//			expectedWip:      expectedWipBatch,
//			closeBatchParams: closeBatchParams,
//			batches:          batches,
//		},
//		{
//			name:             "Close Batch Error",
//			expectedWip:      expectedWipBatch,
//			closeBatchParams: closeBatchParams,
//			batches:          batches,
//			closeBatchErr:    testErr,
//			expectedErr:      fmt.Errorf("failed to close batch, err: %w", testErr),
//		},
//		{
//			name:             "Open Batch Error",
//			expectedWip:      expectedWipBatch,
//			closeBatchParams: closeBatchParams,
//			batches:          batches,
//			openBatchErr:     testErr,
//			expectedErr:      fmt.Errorf("failed to open new batch, err: %w", testErr),
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			// arrange
//			dbManagerMock.On("CloseBatch", ctx, tc.closeBatchParams).Return(tc.closeBatchErr).Once()
//			executorMock.On("ProcessBatch", ctx, f.processRequest).Return(&state.ProcessBatchResponse{
//				IsBatchProcessed: true,
//			}, nilErr).Once()
//
//			if tc.expectedErr == nil {
//				dbManagerMock.On("GetTransactionsByBatchNumber", ctx, f.batch.batchNumber).Return(txs, nilErr).Once()
//			}
//
//			if tc.closeBatchErr == nil {
//				dbManagerMock.On("BeginStateTransaction", ctx).Return(dbTxMock, nilErr).Once()
//				dbManagerMock.On("OpenBatch", ctx, mock.Anything, dbTxMock).Return(tc.openBatchErr).Once()
//
//				// Async Calls from reprocessBatch
//				dbManagerMock.On("GetLastNBatches", ctx, uint(2)).Return(tc.batches, nilErr).Maybe()
//				dbManagerMock.On("GetTransactionsByBatchNumber", ctx, f.batch.batchNumber).Return(txs, nilErr).Maybe()
//				processRequest := f.processRequest
//				processRequest.Caller = state.DiscardCallerLabel
//				processRequest.Timestamp = f.batch.timestamp
//				executorMock.On("ProcessBatch", ctx, processRequest).Return(&state.ProcessBatchResponse{
//					NewStateRoot:     f.batch.stateRoot,
//					NewLocalExitRoot: f.batch.localExitRoot,
//
//					IsBatchProcessed: true,
//				}, nilErr).Maybe()
//
//				if tc.openBatchErr == nil {
//					dbTxMock.On("Commit", ctx).Return(nilErr).Once()
//				} else {
//					dbTxMock.On("Rollback", ctx).Return(nilErr).Once()
//				}
//			}
//
//			// act
//			wipBatch, err := f.newWIPBatch(ctx)
//
//			// assert
//			if tc.expectedErr != nil {
//				assert.Error(t, err)
//				assert.EqualError(t, err, tc.expectedErr.Error())
//				assert.Nil(t, wipBatch)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, tc.expectedWip, wipBatch)
//			}
//			dbManagerMock.AssertExpectations(t)
//			dbTxMock.AssertExpectations(t)
//		})
//	}
//}

func TestFinalizer_handleTransactionError(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	nonce := uint64(0)
	tx := &TxTracker{Hash: oldHash, From: sender, Cost: big.NewInt(0)}
	testCases := []struct {
		name               string
		error              pb.RomError
		expectedDeleteCall bool
		expectedMoveCall   bool
	}{
		{
			name:               "OutOfCountersError",
			error:              pb.RomError(executor.ROM_ERROR_OUT_OF_COUNTERS_STEP),
			expectedDeleteCall: true,
		},
		{
			name:             "IntrinsicError",
			error:            pb.RomError(executor.ROM_ERROR_INTRINSIC_INVALID_NONCE),
			expectedMoveCall: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			if tc.expectedDeleteCall {
				workerMock.On("DeleteTx", oldHash, sender).Return()
				dbManagerMock.On("UpdateTxStatus", ctx, oldHash, pool.TxStatusFailed, false).Return(nil).Once()
				dbManagerMock.On("UpdateTxStatus", ctx, oldHash, pool.TxStatusInvalid, false).Return(nil).Once()
				dbManagerMock.On("DeleteTransactionFromPool", ctx, tx.Hash).Return(nil).Once()
			}
			if tc.expectedMoveCall {
				workerMock.On("MoveTxToNotReady", oldHash, sender, &nonce, big.NewInt(0)).Return([]*TxTracker{}).Once()
			}

			result := &state.ProcessBatchResponse{
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
					sender: {Nonce: &nonce, Balance: big.NewInt(0)},
				},
				Responses: []*state.ProcessTransactionResponse{{
					RomError: executor.RomErr(tc.error),
				},
				},
			}

			// act
			f.handleTransactionError(ctx, result, tx)

			// assert
			workerMock.AssertExpectations(t)
		})
	}
}

func TestFinalizer_syncWithState(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	now = testNow
	defer func() {
		now = time.Now
	}()
	one := uint64(1)
	batches := []*state.Batch{
		{
			BatchNumber:    1,
			StateRoot:      oldHash,
			GlobalExitRoot: oldHash,
		},
	}
	testCases := []struct {
		name                  string
		batches               []*state.Batch
		lastBatchNum          *uint64
		isBatchClosed         bool
		ger                   common.Hash
		getWIPBatchErr        error
		openBatchErr          error
		isBatchClosedErr      error
		getLastBatchErr       error
		expectedProcessingCtx state.ProcessingContext
		expectedBatch         *WipBatch
		expectedErr           error
	}{
		{
			name:          "Success-Closed Batch",
			lastBatchNum:  &one,
			isBatchClosed: true,
			ger:           oldHash,
			batches:       batches,
			expectedBatch: &WipBatch{
				batchNumber:        one + 1,
				coinbase:           f.sequencerAddress,
				initialStateRoot:   oldHash,
				stateRoot:          oldHash,
				timestamp:          uint64(testNow().Unix()),
				globalExitRoot:     oldHash,
				remainingResources: getMaxRemainingResources(f.batchConstraints),
			},
			expectedProcessingCtx: state.ProcessingContext{
				BatchNumber:    one + 1,
				Coinbase:       f.sequencerAddress,
				Timestamp:      testNow(),
				GlobalExitRoot: oldHash,
			},
			expectedErr: nil,
		},
		{
			name:          "Success-Open Batch",
			lastBatchNum:  &one,
			isBatchClosed: false,
			batches:       batches,
			ger:           common.Hash{},
			expectedBatch: &WipBatch{
				batchNumber:        one,
				coinbase:           f.sequencerAddress,
				initialStateRoot:   oldHash,
				stateRoot:          oldHash,
				timestamp:          uint64(testNow().Unix()),
				globalExitRoot:     oldHash,
				remainingResources: getMaxRemainingResources(f.batchConstraints),
			},
			expectedProcessingCtx: state.ProcessingContext{
				BatchNumber:    one,
				Coinbase:       f.sequencerAddress,
				Timestamp:      testNow(),
				GlobalExitRoot: oldHash,
			},
		},
		{
			name:            "Error-Failed to get last batch",
			lastBatchNum:    nil,
			batches:         batches,
			isBatchClosed:   true,
			ger:             oldHash,
			getLastBatchErr: testErr,
			expectedErr:     fmt.Errorf("failed to get last batch, err: %w", testErr),
		},
		{
			name:             "Error-Failed to check if batch is closed",
			lastBatchNum:     &one,
			batches:          batches,
			isBatchClosed:    true,
			ger:              oldHash,
			isBatchClosedErr: testErr,
			expectedErr:      fmt.Errorf("failed to check if batch is closed, err: %w", testErr),
		},
		{
			name:           "Error-Failed to get work-in-progress batch",
			lastBatchNum:   &one,
			batches:        batches,
			isBatchClosed:  false,
			ger:            common.Hash{},
			getWIPBatchErr: testErr,
			expectedErr:    fmt.Errorf("failed to get work-in-progress batch, err: %w", testErr),
		},
		{
			name:          "Error-Failed to open new batch",
			lastBatchNum:  &one,
			batches:       batches,
			isBatchClosed: true,
			ger:           oldHash,
			openBatchErr:  testErr,
			expectedProcessingCtx: state.ProcessingContext{
				BatchNumber:    one + 1,
				Coinbase:       f.sequencerAddress,
				Timestamp:      testNow(),
				GlobalExitRoot: oldHash,
			},
			expectedErr: fmt.Errorf("failed to open new batch, err: %w", testErr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			if tc.lastBatchNum == nil {
				dbManagerMock.Mock.On("GetLastBatch", ctx).Return(tc.batches[0], tc.getLastBatchErr).Once()
			} else {
				dbManagerMock.On("GetBatchByNumber", ctx, *tc.lastBatchNum, nil).Return(tc.batches[0], nilErr).Once()
			}

			if tc.getLastBatchErr == nil {
				dbManagerMock.Mock.On("IsBatchClosed", ctx, *tc.lastBatchNum).Return(tc.isBatchClosed, tc.isBatchClosedErr).Once()
			}

			if tc.isBatchClosed {
				if tc.getLastBatchErr == nil && tc.isBatchClosedErr == nil {
					dbManagerMock.On("OpenBatch", ctx, tc.expectedProcessingCtx, dbTxMock).Return(tc.openBatchErr).Once()
				}

				if tc.getLastBatchErr == nil && tc.isBatchClosedErr == nil {
					dbManagerMock.Mock.On("GetLatestGer", ctx, f.cfg.GERFinalityNumberOfBlocks).Return(state.GlobalExitRoot{GlobalExitRoot: tc.ger}, testNow(), nil).Once()
					dbManagerMock.On("BeginStateTransaction", ctx).Return(dbTxMock, nil).Once()
					if tc.openBatchErr == nil {
						dbTxMock.On("Commit", ctx).Return(nil).Once()
					}
				}
				if tc.expectedErr != nil && tc.openBatchErr != nil {
					dbTxMock.On("Rollback", ctx).Return(nil).Once()
				}
			} else {
				dbManagerMock.Mock.On("GetWIPBatch", ctx).Return(tc.expectedBatch, tc.getWIPBatchErr).Once()
			}

			// act
			err := f.syncWithState(ctx, tc.lastBatchNum)

			// assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBatch, f.batch)
			}
			dbManagerMock.AssertExpectations(t)
		})
	}
}

func TestFinalizer_processForcedBatches(t *testing.T) {
	// arrange
	var err error
	f = setupFinalizer(false)
	now = testNow
	defer func() {
		now = time.Now
	}()
	RawTxsData1 := make([]byte, 0, 2)
	RawTxsData1 = append(RawTxsData1, []byte("forced tx 1")...)
	RawTxsData1 = append(RawTxsData1, []byte("forced tx 2")...)
	RawTxsData2 := make([]byte, 0, 2)
	RawTxsData2 = append(RawTxsData2, []byte("forced tx 3")...)
	RawTxsData2 = append(RawTxsData2, []byte("forced tx 4")...)
	batchNumber := f.batch.batchNumber
	stateRoot := oldHash
	forcedBatch1 := state.ForcedBatch{
		ForcedBatchNumber: 2,
		GlobalExitRoot:    oldHash,
		RawTxsData:        RawTxsData1,
	}
	forcedBatch2 := state.ForcedBatch{
		ForcedBatchNumber: 3,
		GlobalExitRoot:    oldHash,
		RawTxsData:        RawTxsData2,
	}
	testCases := []struct {
		name                            string
		forcedBatch                     []state.ForcedBatch
		getLastTrustedForcedBatchNumErr error
		expectedErr                     error
	}{
		{
			name:        "Success",
			forcedBatch: []state.ForcedBatch{forcedBatch1, forcedBatch2},
		},
		{
			name:                            "GetLastTrustedForcedBatchNumber_Error",
			forcedBatch:                     []state.ForcedBatch{forcedBatch1},
			getLastTrustedForcedBatchNumErr: testErr,
			expectedErr:                     fmt.Errorf("failed to get last trusted forced batch number, err: %s", testErr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f.nextForcedBatches = tc.forcedBatch
			internalBatchNumber := batchNumber
			dbManagerMock.On("BeginStateTransaction", ctx).Return(dbTxMock, nil).Once()
			dbManagerMock.On("GetLastTrustedForcedBatchNumber", ctx, dbTxMock).Return(uint64(1), tc.getLastTrustedForcedBatchNumErr).Once()

			for _, forcedBatch := range tc.forcedBatch {
				internalBatchNumber += 1
				processRequest := state.ProcessRequest{
					BatchNumber:    internalBatchNumber,
					OldStateRoot:   stateRoot,
					GlobalExitRoot: forcedBatch.GlobalExitRoot,
					Transactions:   forcedBatch.RawTxsData,
					Coinbase:       f.sequencerAddress,
					Timestamp:      uint64(now().Unix()),
					Caller:         state.SequencerCallerLabel,
				}
				dbManagerMock.On("ProcessForcedBatch", forcedBatch.ForcedBatchNumber, processRequest).Return(&state.ProcessBatchResponse{
					NewStateRoot:   stateRoot,
					NewBatchNumber: internalBatchNumber,
				}, nilErr).Once()
			}

			// act
			batchNumber, stateRoot, err = f.processForcedBatches(ctx, batchNumber, stateRoot)

			// assert
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, tc.expectedErr)
				dbManagerMock.AssertExpectations(t)
			}
		})
	}
}

func TestFinalizer_openWIPBatch(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	now = testNow
	defer func() {
		now = time.Now
	}()
	batchNum := f.batch.batchNumber + 1
	expectedWipBatch := &WipBatch{
		batchNumber:        batchNum,
		coinbase:           f.sequencerAddress,
		initialStateRoot:   oldHash,
		stateRoot:          oldHash,
		timestamp:          uint64(now().Unix()),
		globalExitRoot:     oldHash,
		remainingResources: getMaxRemainingResources(f.batchConstraints),
	}
	testCases := []struct {
		name         string
		openBatchErr error
		beginTxErr   error
		commitErr    error
		rollbackErr  error
		expectedWip  *WipBatch
		expectedErr  error
	}{
		{
			name:        "Success",
			expectedWip: expectedWipBatch,
		},
		{
			name:        "BeginTransaction Error",
			beginTxErr:  testErr,
			expectedErr: fmt.Errorf("failed to begin state transaction to open batch, err: %w", testErr),
		},
		{
			name:         "OpenBatch Error",
			openBatchErr: testErr,
			expectedErr:  fmt.Errorf("failed to open new batch, err: %w", testErr),
		},
		{
			name:        "Commit Error",
			commitErr:   testErr,
			expectedErr: fmt.Errorf("failed to commit database transaction for opening a batch, err: %w", testErr),
		},
		{
			name:         "Rollback Error",
			openBatchErr: testErr,
			rollbackErr:  testErr,
			expectedErr: fmt.Errorf(
				"failed to rollback dbTx: %s. Rollback err: %w",
				testErr.Error(), openBatchError,
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			dbManagerMock.On("BeginStateTransaction", ctx).Return(dbTxMock, tc.beginTxErr).Once()
			if tc.beginTxErr == nil {
				dbManagerMock.On("OpenBatch", ctx, mock.Anything, dbTxMock).Return(tc.openBatchErr).Once()
			}

			if tc.expectedErr != nil && (tc.rollbackErr != nil || tc.openBatchErr != nil) {
				dbTxMock.On("Rollback", ctx).Return(tc.rollbackErr).Once()
			}

			if tc.expectedErr == nil || tc.commitErr != nil {
				dbTxMock.On("Commit", ctx).Return(tc.commitErr).Once()
			}

			// act
			wipBatch, err := f.openWIPBatch(ctx, batchNum, oldHash, oldHash)

			// assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Nil(t, wipBatch)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedWip, wipBatch)
			}
			dbManagerMock.AssertExpectations(t)
			dbTxMock.AssertExpectations(t)
		})
	}
}

// TestFinalizer_closeBatch tests the closeBatch method.
func TestFinalizer_closeBatch(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	txs := make([]types.Transaction, 0)
	receipt := ClosingBatchParameters{
		BatchNumber:   f.batch.batchNumber,
		StateRoot:     f.batch.stateRoot,
		LocalExitRoot: f.processRequest.GlobalExitRoot,
		Txs:           txs,
	}
	managerErr := fmt.Errorf("some error")
	testCases := []struct {
		name        string
		managerErr  error
		expectedErr error
	}{
		{
			name:        "Success",
			managerErr:  nil,
			expectedErr: nil,
		},
		{
			name:        "Manager Error",
			managerErr:  managerErr,
			expectedErr: fmt.Errorf("failed to get transactions from transactions, err: %w", managerErr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			dbManagerMock.Mock.On("CloseBatch", ctx, receipt).Return(tc.managerErr).Once()
			dbManagerMock.Mock.On("GetTransactionsByBatchNumber", ctx, receipt.BatchNumber).Return(txs, tc.managerErr).Once()

			// act
			err := f.closeBatch(ctx)

			// assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.ErrorIs(t, err, tc.managerErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFinalizer_openBatch(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	now = testNow
	defer func() {
		now = time.Now
	}()
	batchNum := f.batch.batchNumber + 1
	testCases := []struct {
		name        string
		batchNum    uint64
		managerErr  error
		expectedCtx state.ProcessingContext
		expectedErr error
	}{
		{
			name:       "Success",
			batchNum:   batchNum,
			managerErr: nil,
			expectedCtx: state.ProcessingContext{
				BatchNumber:    batchNum,
				Coinbase:       f.sequencerAddress,
				Timestamp:      now(),
				GlobalExitRoot: oldHash,
			},
			expectedErr: nil,
		},
		{
			name:        "Manager Error",
			batchNum:    batchNum,
			managerErr:  testErr,
			expectedCtx: state.ProcessingContext{},
			expectedErr: openBatchError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			dbManagerMock.Mock.On("OpenBatch", mock.Anything, mock.Anything, mock.Anything).Return(tc.managerErr).Once()

			// act
			actualCtx, err := f.openBatch(ctx, tc.batchNum, oldHash, nil)

			// assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.ErrorIs(t, err, tc.managerErr)
				assert.Empty(t, actualCtx)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCtx, actualCtx)
			}
			dbManagerMock.AssertExpectations(t)
		})
	}
}

/*
// TestFinalizer_reprocessBatch is a test for reprocessBatch which tests all possible cases of reprocessBatch
func TestFinalizer_reprocessBatch(t *testing.T) {
	// arrange
	now = testNow
	defer func() {
		now = time.Now
	}()

	f = setupFinalizer(true)
	n := uint(2)
	expectedProcessBatchRequest := state.ProcessRequest{
		BatchNumber:  f.batch.batchNumber,
		OldStateRoot: oldHash,
		Coinbase:     f.sequencerAddress,
		Timestamp:    f.batch.timestamp,
		Caller:       state.DiscardCallerLabel,
	}
	batches := []*state.Batch{
		{
			BatchNumber: f.batch.batchNumber,
			StateRoot:   newHash,
			Timestamp:   testNow(),
		},
		{
			BatchNumber:    f.batch.batchNumber - 1,
			GlobalExitRoot: oldHash,
			StateRoot:      oldHash,
			Timestamp:      testNow(),
		},
	}

	// TODO: Add missing cases for this test
	testCases := []struct {
		name                       string
		getLastNBatchesErr         error
		processBatchErr            error
		batches                    []*state.Batch
		expectedErr                error
		internalErr                error
		expectedProcessRequest     state.ProcessRequest
		expectedProcessBatchResult *state.ProcessBatchResponse
	}{
		{
			name:                   "Success",
			batches:                batches,
			expectedProcessRequest: expectedProcessBatchRequest,
			expectedProcessBatchResult: &state.ProcessBatchResponse{
				NewStateRoot:     newHash,
				IsBatchProcessed: true,
			},
		},
		{
			name:               "GetLastNBatches Error",
			getLastNBatchesErr: testErr,
			internalErr:        testErr,
			expectedErr:        fmt.Errorf("failed to get old state root, err: failed to get last %d batches, err: %w", n, testErr),
		},
		{
			name:                   "ProcessBatch Error",
			processBatchErr:        testErr,
			internalErr:            testErr,
			expectedErr:            testErr,
			expectedProcessRequest: expectedProcessBatchRequest,
			batches:                batches,
		},
		{
			name:                   "ProcessBatch Result Error",
			processBatchErr:        testErr,
			internalErr:            testErr,
			expectedProcessRequest: expectedProcessBatchRequest,
			expectedErr:            testErr,
			batches:                batches,
			expectedProcessBatchResult: &state.ProcessBatchResponse{
				IsBatchProcessed: false,
				ExecutorError:    testErr2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f.processRequest = tc.expectedProcessRequest
			dbManagerMock.On("GetLastNBatches", ctx, n).Return(tc.batches, tc.getLastNBatchesErr).Once()
			if tc.getLastNBatchesErr == nil {
				executorMock.Mock.On("ProcessBatch", ctx, f.processRequest, true).Return(tc.expectedProcessBatchResult, tc.processBatchErr).Once()
				dbManagerMock.On("GetBatchByNumber", ctx, f.batch.batchNumber, nil).Return(tc.batches[0], nil).Once()
			}

			// act
			err := f.reprocessBatch(ctx, f.batch.batchNumber)

			// assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.ErrorIs(t, err, tc.internalErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
*/

func TestFinalizer_isDeadlineEncountered(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	now = testNow
	defer func() {
		now = time.Now
	}()
	testCases := []struct {
		name             string
		nextForcedBatch  int64
		nextGER          int64
		nextDelayedBatch int64
		expected         bool
	}{
		{
			name:             "No deadlines",
			nextForcedBatch:  0,
			nextGER:          0,
			nextDelayedBatch: 0,
			expected:         false,
		},
		{
			name:             "Forced batch deadline",
			nextForcedBatch:  now().Add(time.Second).Unix(),
			nextGER:          0,
			nextDelayedBatch: 0,
			expected:         true,
		},
		{
			name:             "Global Exit Root deadline",
			nextForcedBatch:  0,
			nextGER:          now().Add(time.Second).Unix(),
			nextDelayedBatch: 0,
			expected:         true,
		},
		{
			name:             "Delayed batch deadline",
			nextForcedBatch:  0,
			nextGER:          0,
			nextDelayedBatch: now().Add(time.Second).Unix(),
			expected:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f.nextForcedBatchDeadline = tc.nextForcedBatch
			f.nextGERDeadline = tc.nextGER
			f.nextSendingToL1Deadline = tc.nextDelayedBatch
			if tc.expected == true {
				now = func() time.Time {
					return testNow().Add(time.Second * 2)
				}
			}

			// act
			actual := f.isDeadlineEncountered()

			// assert
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestFinalizer_checkRemainingResources(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	ctx := context.Background()
	txResponse := &state.ProcessTransactionResponse{TxHash: oldHash}
	result := &state.ProcessBatchResponse{
		UsedZkCounters: state.ZKCounters{CumulativeGasUsed: 1000},
		Responses:      []*state.ProcessTransactionResponse{txResponse},
	}
	remainingResources := batchResources{
		zKCounters: state.ZKCounters{CumulativeGasUsed: 9000},
		bytes:      10000,
	}
	f.batch.remainingResources = remainingResources
	testCases := []struct {
		name                 string
		remaining            batchResources
		expectedErr          error
		expectedWorkerUpdate bool
		expectedTxTracker    *TxTracker
	}{
		{
			name:                 "Success",
			remaining:            remainingResources,
			expectedErr:          nil,
			expectedWorkerUpdate: false,
			expectedTxTracker:    &TxTracker{RawTx: []byte("test")},
		},
		{
			name: "Bytes Resource Exceeded",
			remaining: batchResources{
				bytes: 0,
			},
			expectedErr:          ErrBatchResourceBytesUnderflow,
			expectedWorkerUpdate: true,
			expectedTxTracker:    &TxTracker{RawTx: []byte("test")},
		},
		{
			name: "ZkCounter Resource Exceeded",
			remaining: batchResources{
				zKCounters: state.ZKCounters{CumulativeGasUsed: 0},
			},
			expectedErr:          NewBatchRemainingResourcesUnderflowError(cumulativeGasErr, cumulativeGasErr.Error()),
			expectedWorkerUpdate: true,
			expectedTxTracker:    &TxTracker{RawTx: make([]byte, 0)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f.batch.remainingResources = tc.remaining
			dbManagerMock.On("AddEvent", ctx, mock.Anything, nil).Return(nil)
			if tc.expectedWorkerUpdate {
				workerMock.On("UpdateTx", txResponse.TxHash, tc.expectedTxTracker.From, result.UsedZkCounters).Return().Once()
			}

			// act
			err := f.checkRemainingResources(ctx, result, tc.expectedTxTracker)

			// assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			if tc.expectedWorkerUpdate {
				workerMock.AssertCalled(t, "UpdateTx", txResponse.TxHash, tc.expectedTxTracker.From, result.UsedZkCounters)
			} else {
				workerMock.AssertNotCalled(t, "UpdateTx", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestFinalizer_isBatchReadyToClose(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	maxRemainingResource := getMaxRemainingResources(bc)
	testCases := []struct {
		name              string
		cumulativeGasUsed uint64
		expectedResult    bool
	}{
		{
			name:              "Is ready",
			cumulativeGasUsed: f.getConstraintThresholdUint64(bc.MaxCumulativeGasUsed) - 1,
			expectedResult:    true,
		}, {
			name:              "Is NOT ready",
			cumulativeGasUsed: f.getConstraintThresholdUint64(bc.MaxCumulativeGasUsed) + 1,
			expectedResult:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			maxRemainingResource.zKCounters.CumulativeGasUsed = tc.cumulativeGasUsed
			f.batch.remainingResources = maxRemainingResource
			// act
			result := f.isBatchAlmostFull()

			// assert
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestFinalizer_setNextForcedBatchDeadline(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	now = testNow
	defer func() {
		now = time.Now
	}()
	expected := now().Unix() + int64(f.cfg.ForcedBatchDeadlineTimeoutInSec.Duration.Seconds())

	// act
	f.setNextForcedBatchDeadline()

	// assert
	assert.Equal(t, expected, f.nextForcedBatchDeadline)
}

func TestFinalizer_setNextGERDeadline(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	now = testNow
	defer func() {
		now = time.Now
	}()
	expected := now().Unix() + int64(f.cfg.GERDeadlineTimeoutInSec.Duration.Seconds())

	// act
	f.setNextGERDeadline()

	// assert
	assert.Equal(t, expected, f.nextGERDeadline)
}

func TestFinalizer_setNextSendingToL1Deadline(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	now = testNow
	defer func() {
		now = time.Now
	}()
	expected := now().Unix() + int64(f.cfg.SendingToL1DeadlineTimeoutInSec.Duration.Seconds())

	// act
	f.setNextSendingToL1Deadline()

	// assert
	assert.Equal(t, expected, f.nextSendingToL1Deadline)
}

func TestFinalizer_getConstraintThresholdUint64(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	input := uint64(100)
	expect := input * uint64(f.cfg.ResourcePercentageToCloseBatch) / 100

	// act
	result := f.getConstraintThresholdUint64(input)

	// assert
	assert.Equal(t, result, expect)
}

func TestFinalizer_getConstraintThresholdUint32(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	input := uint32(100)
	expect := uint32(input * f.cfg.ResourcePercentageToCloseBatch / 100)

	// act
	result := f.getConstraintThresholdUint32(input)

	// assert
	assert.Equal(t, result, expect)
}

func TestFinalizer_getRemainingResources(t *testing.T) {
	// act
	remainingResources := getMaxRemainingResources(bc)

	// assert
	assert.Equal(t, remainingResources.zKCounters.CumulativeGasUsed, bc.MaxCumulativeGasUsed)
	assert.Equal(t, remainingResources.zKCounters.UsedKeccakHashes, bc.MaxKeccakHashes)
	assert.Equal(t, remainingResources.zKCounters.UsedPoseidonHashes, bc.MaxPoseidonHashes)
	assert.Equal(t, remainingResources.zKCounters.UsedPoseidonPaddings, bc.MaxPoseidonPaddings)
	assert.Equal(t, remainingResources.zKCounters.UsedMemAligns, bc.MaxMemAligns)
	assert.Equal(t, remainingResources.zKCounters.UsedArithmetics, bc.MaxArithmetics)
	assert.Equal(t, remainingResources.zKCounters.UsedBinaries, bc.MaxBinaries)
	assert.Equal(t, remainingResources.zKCounters.UsedSteps, bc.MaxSteps)
	assert.Equal(t, remainingResources.bytes, bc.MaxBatchBytesSize)
}

func setupFinalizer(withWipBatch bool) *finalizer {
	wipBatch := new(WipBatch)
	dbManagerMock = new(DbManagerMock)
	executorMock = new(StateMock)
	workerMock = new(WorkerMock)
	dbTxMock = new(DbTxMock)
	if withWipBatch {
		wipBatch = &WipBatch{
			batchNumber:        1,
			coinbase:           seqAddr,
			initialStateRoot:   oldHash,
			stateRoot:          newHash,
			timestamp:          uint64(now().Unix()),
			globalExitRoot:     oldHash,
			remainingResources: getMaxRemainingResources(bc),
		}
	}
	return &finalizer{
		cfg:                cfg,
		txsStore:           txsStore,
		closingSignalCh:    closingSignalCh,
		isSynced:           isSynced,
		sequencerAddress:   seqAddr,
		worker:             workerMock,
		dbManager:          dbManagerMock,
		executor:           executorMock,
		sharedResourcesMux: new(sync.RWMutex),
		batch:              wipBatch,
		batchConstraints:   bc,
		processRequest:     state.ProcessRequest{},
		// closing signals
		nextGER:                   common.Hash{},
		nextGERDeadline:           0,
		nextGERMux:                new(sync.RWMutex),
		nextForcedBatches:         make([]state.ForcedBatch, 0),
		nextForcedBatchDeadline:   0,
		nextForcedBatchesMux:      new(sync.RWMutex),
		nextSendingToL1Deadline:   0,
		nextSendingToL1TimeoutMux: new(sync.RWMutex),
	}
}
