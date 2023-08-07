package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	f             *finalizer
	nilErr        error
	dbManagerMock = new(DbManagerMock)
	executorMock  = new(StateMock)
	workerMock    = new(WorkerMock)
	dbTxMock      = new(DbTxMock)
	bc            = batchConstraints{
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
	closingSignalCh = ClosingSignalCh{
		ForcedBatchCh: make(chan state.ForcedBatch),
		GERCh:         make(chan common.Hash),
		L2ReorgCh:     make(chan L2ReorgEvent),
	}
	effectiveGasPriceCfg = EffectiveGasPriceCfg{
		MaxBreakEvenGasPriceDeviationPercentage: 10,
		L1GasPriceFactor:                        0.25,
		ByteGasCost:                             16,
		MarginFactor:                            1,
		Enabled:                                 false,
	}
	cfg = FinalizerCfg{
		GERDeadlineTimeout: cfgTypes.Duration{
			Duration: 60,
		},
		ForcedBatchDeadlineTimeout: cfgTypes.Duration{
			Duration: 60,
		},
		SleepDuration: cfgTypes.Duration{
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
	nonce1          = uint64(1)
	nonce2          = uint64(2)
	seqAddr         = common.Address{}
	oldHash         = common.HexToHash("0x01")
	newHash         = common.HexToHash("0x02")
	newHash2        = common.HexToHash("0x03")
	stateRootHashes = []common.Hash{oldHash, newHash, newHash2}
	txHash          = common.BytesToHash([]byte("txHash"))
	txHash2         = common.BytesToHash([]byte("txHash2"))
	tx              = types.NewTransaction(nonce1, receiverAddr, big.NewInt(1), 100000, big.NewInt(1), nil)
	senderAddr      = common.HexToAddress("0x3445324")
	receiverAddr    = common.HexToAddress("0x1555324")
	isSynced        = func(ctx context.Context) bool {
		return true
	}
	testErrStr              = "some err"
	testErr                 = fmt.Errorf(testErrStr)
	openBatchError          = fmt.Errorf("failed to open new batch, err: %w", testErr)
	cumulativeGasErr        = state.GetZKCounterError("CumulativeGasUsed")
	testBatchL2DataAsString = "0xee80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e980801186622d03b6b8da7cf111d1ccba5bb185c56deae6a322cebc6dda0556f3cb9700910c26408b64b51c5da36ba2f38ef55ba1cee719d5a6c012259687999074321bff"
	decodedBatchL2Data      []byte
	done                    chan bool
	gasPrice                = big.NewInt(1000000)
	breakEvenGasPrice       = big.NewInt(1000000)
	l1GasPrice              = uint64(1000000)
)

func testNow() time.Time {
	return time.Unix(0, 0)
}

func TestNewFinalizer(t *testing.T) {
	eventStorage, err := nileventstorage.NewNilEventStorage()
	require.NoError(t, err)
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	dbManagerMock.On("GetLastSentFlushID", context.Background()).Return(uint64(0), nil)

	// arrange and act
	f = newFinalizer(cfg, effectiveGasPriceCfg, workerMock, dbManagerMock, executorMock, seqAddr, isSynced, closingSignalCh, bc, eventLog)

	// assert
	assert.NotNil(t, f)
	assert.Equal(t, f.cfg, cfg)
	assert.Equal(t, f.worker, workerMock)
	assert.Equal(t, dbManagerMock, dbManagerMock)
	assert.Equal(t, f.executor, executorMock)
	assert.Equal(t, f.sequencerAddress, seqAddr)
	assert.Equal(t, f.closingSignalCh, closingSignalCh)
	assert.Equal(t, f.batchConstraints, bc)
}

func TestFinalizer_handleProcessTransactionResponse(t *testing.T) {
	f = setupFinalizer(true)
	ctx = context.Background()
	txTracker := &TxTracker{Hash: txHash, From: senderAddr, Nonce: 1, GasPrice: gasPrice, BreakEvenGasPrice: breakEvenGasPrice, L1GasPrice: l1GasPrice, BatchResources: state.BatchResources{
		Bytes: 1000,
		ZKCounters: state.ZKCounters{
			CumulativeGasUsed: 500,
		},
	}}

	txResponse := &state.ProcessTransactionResponse{
		TxHash:    txHash,
		StateRoot: newHash2,
		RomError:  nil,
		GasUsed:   100000,
	}
	batchResponse := &state.ProcessBatchResponse{
		Responses: []*state.ProcessTransactionResponse{
			txResponse,
		},
	}
	txResponseIntrinsicError := &state.ProcessTransactionResponse{
		TxHash:    txHash,
		StateRoot: newHash2,
		RomError:  runtime.ErrIntrinsicInvalidNonce,
	}
	txResponseOOCError := &state.ProcessTransactionResponse{
		TxHash:    txHash,
		StateRoot: newHash2,
		RomError:  runtime.ErrOutOfCountersKeccak,
	}
	testCases := []struct {
		name                       string
		executorResponse           *state.ProcessBatchResponse
		oldStateRoot               common.Hash
		expectedStoredTx           transactionToStore
		expectedMoveToNotReadyCall bool
		expectedDeleteTxCall       bool
		expectedUpdateTxCall       bool
		expectedError              error
		expectedUpdateTxStatus     pool.TxStatus
	}{
		{
			name: "Successful transaction",
			executorResponse: &state.ProcessBatchResponse{
				Responses: []*state.ProcessTransactionResponse{
					txResponse,
				},
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
					senderAddr: {
						Address: senderAddr,
						Nonce:   &nonce2,
						Balance: big.NewInt(100),
					},
					receiverAddr: {
						Address: receiverAddr,
						Nonce:   nil,
						Balance: big.NewInt(100),
					},
				},
			},
			oldStateRoot: oldHash,
			expectedStoredTx: transactionToStore{
				txTracker:     txTracker,
				batchNumber:   f.batch.batchNumber,
				coinbase:      f.batch.coinbase,
				timestamp:     f.batch.timestamp,
				oldStateRoot:  oldHash,
				batchResponse: batchResponse,
				response:      txResponse,
				isForcedBatch: false,
			},
		},
		{
			name: "Batch resources underflow err",
			executorResponse: &state.ProcessBatchResponse{
				UsedZkCounters: state.ZKCounters{
					CumulativeGasUsed: f.batch.remainingResources.ZKCounters.CumulativeGasUsed + 1,
				},
				Responses: []*state.ProcessTransactionResponse{
					txResponse,
				},
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
					senderAddr: {
						Address: senderAddr,
						Nonce:   &nonce1,
						Balance: big.NewInt(100),
					},
				},
			},
			oldStateRoot:         oldHash,
			expectedUpdateTxCall: true,
			expectedError:        state.NewBatchRemainingResourcesUnderflowError(cumulativeGasErr, cumulativeGasErr.Error()),
		},
		{
			name: "Intrinsic err",
			executorResponse: &state.ProcessBatchResponse{
				IsRomOOCError: false,
				UsedZkCounters: state.ZKCounters{
					CumulativeGasUsed: 1,
				},
				Responses: []*state.ProcessTransactionResponse{
					txResponseIntrinsicError,
				},
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
					senderAddr: {
						Address: senderAddr,
						Nonce:   &nonce1,
						Balance: big.NewInt(100),
					},
				},
			},
			oldStateRoot:               oldHash,
			expectedMoveToNotReadyCall: true,
			expectedError:              txResponseIntrinsicError.RomError,
		},
		{
			name: "Out Of Counters err",
			executorResponse: &state.ProcessBatchResponse{
				IsRomOOCError: true,
				UsedZkCounters: state.ZKCounters{
					UsedKeccakHashes: bc.MaxKeccakHashes + 1,
				},
				Responses: []*state.ProcessTransactionResponse{
					txResponseOOCError,
				},
			},
			oldStateRoot:           oldHash,
			expectedError:          txResponseOOCError.RomError,
			expectedDeleteTxCall:   true,
			expectedUpdateTxStatus: pool.TxStatusInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storedTxs := make([]transactionToStore, 0)
			f.pendingTransactionsToStore = make(chan transactionToStore)

			if tc.expectedStoredTx.batchResponse != nil {
				done = make(chan bool) // init a new done channel
				go func() {
					for tx := range f.pendingTransactionsToStore {
						storedTxs = append(storedTxs, tx)
						f.pendingTransactionsToStoreWG.Done()
					}
					done <- true // signal that the goroutine is done
				}()
			}
			if tc.expectedDeleteTxCall {
				workerMock.On("DeleteTx", txTracker.Hash, txTracker.From).Return().Once()
			}
			if tc.expectedMoveToNotReadyCall {
				addressInfo := tc.executorResponse.ReadWriteAddresses[senderAddr]
				workerMock.On("MoveTxToNotReady", txHash, senderAddr, addressInfo.Nonce, addressInfo.Balance).Return([]*TxTracker{}).Once()
			}
			if tc.expectedUpdateTxCall {
				workerMock.On("UpdateTx", txTracker.Hash, txTracker.From, tc.executorResponse.UsedZkCounters).Return().Once()
			}
			if tc.expectedError == nil {
				//dbManagerMock.On("GetGasPrices", ctx).Return(pool.GasPrices{L1GasPrice: 0, L2GasPrice: 0}, nilErr).Once()
				workerMock.On("DeleteTx", txTracker.Hash, txTracker.From).Return().Once()
				workerMock.On("UpdateAfterSingleSuccessfulTxExecution", txTracker.From, tc.executorResponse.ReadWriteAddresses).Return([]*TxTracker{}).Once()
			}
			if tc.expectedUpdateTxStatus != "" {
				dbManagerMock.On("UpdateTxStatus", ctx, txHash, tc.expectedUpdateTxStatus, false, mock.Anything).Return(nil).Once()
			}

			errWg, err := f.handleProcessTransactionResponse(ctx, txTracker, tc.executorResponse, tc.oldStateRoot)
			if errWg != nil {
				errWg.Wait()
			}

			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
			} else {
				require.Nil(t, err)
			}

			if tc.expectedStoredTx.batchResponse != nil {
				close(f.pendingTransactionsToStore) // close the channel
				<-done                              // wait for the goroutine to finish
				f.pendingTransactionsToStoreWG.Wait()
				require.Len(t, storedTxs, 1)
				actualTx := storedTxs[0]
				assertEqualTransactionToStore(t, tc.expectedStoredTx, actualTx)
			} else {
				require.Empty(t, storedTxs)
			}

			workerMock.AssertExpectations(t)
			dbManagerMock.AssertExpectations(t)
		})
	}
}

func assertEqualTransactionToStore(t *testing.T, expectedTx, actualTx transactionToStore) {
	require.Equal(t, expectedTx.txTracker, actualTx.txTracker)
	require.Equal(t, expectedTx.response, actualTx.response)
	require.Equal(t, expectedTx.batchNumber, actualTx.batchNumber)
	require.Equal(t, expectedTx.timestamp, actualTx.timestamp)
	require.Equal(t, expectedTx.coinbase, actualTx.coinbase)
	require.Equal(t, expectedTx.oldStateRoot, actualTx.oldStateRoot)
	require.Equal(t, expectedTx.isForcedBatch, actualTx.isForcedBatch)
	require.Equal(t, expectedTx.flushId, actualTx.flushId)
}

func TestFinalizer_newWIPBatch(t *testing.T) {
	// arrange
	now = testNow
	defer func() {
		now = time.Now
	}()

	f = setupFinalizer(true)
	f.processRequest.Caller = stateMetrics.SequencerCallerLabel
	f.processRequest.Timestamp = now()
	f.processRequest.Transactions = decodedBatchL2Data

	stateRootErr := errors.New("state root must have value to close batch")
	txs := []types.Transaction{*tx}
	require.NoError(t, err)
	newBatchNum := f.batch.batchNumber + 1
	expectedNewWipBatch := &WipBatch{
		batchNumber:        newBatchNum,
		coinbase:           f.sequencerAddress,
		initialStateRoot:   newHash,
		stateRoot:          newHash,
		timestamp:          now(),
		remainingResources: getMaxRemainingResources(f.batchConstraints),
	}
	closeBatchParams := ClosingBatchParameters{
		BatchNumber:          f.batch.batchNumber,
		StateRoot:            newHash,
		LocalExitRoot:        f.batch.localExitRoot,
		Txs:                  txs,
		EffectivePercentages: []uint8{255},
	}

	batches := []*state.Batch{
		{
			BatchNumber:    f.batch.batchNumber,
			StateRoot:      newHash,
			GlobalExitRoot: oldHash,
			Transactions:   txs,
			Timestamp:      now(),
			BatchL2Data:    decodedBatchL2Data,
		},
	}

	// For Empty Batch
	expectedNewWipEmptyBatch := *expectedNewWipBatch
	expectedNewWipEmptyBatch.initialStateRoot = oldHash
	expectedNewWipEmptyBatch.stateRoot = oldHash
	emptyBatch := *batches[0]
	emptyBatch.StateRoot = oldHash
	emptyBatch.Transactions = make([]types.Transaction, 0)
	emptyBatch.BatchL2Data = []byte{}
	emptyBatch.GlobalExitRoot = oldHash
	emptyBatchBatches := []*state.Batch{&emptyBatch}
	closeBatchParamsForEmptyBatch := closeBatchParams
	closeBatchParamsForEmptyBatch.StateRoot = oldHash
	closeBatchParamsForEmptyBatch.Txs = nil

	// For Forced Batch
	expectedForcedNewWipBatch := *expectedNewWipBatch
	expectedForcedNewWipBatch.batchNumber = expectedNewWipBatch.batchNumber + 1
	expectedForcedNewWipBatch.globalExitRoot = oldHash

	testCases := []struct {
		name                       string
		batches                    []*state.Batch
		closeBatchErr              error
		closeBatchParams           ClosingBatchParameters
		stateRootAndLERErr         error
		openBatchErr               error
		expectedWip                *WipBatch
		reprocessFullBatchResponse *state.ProcessBatchResponse
		expectedErr                error
		reprocessBatchErr          error
		forcedBatches              []state.ForcedBatch
	}{

		{
			name:               "Error StateRoot must have value",
			stateRootAndLERErr: stateRootErr,
			expectedErr:        stateRootErr,
		},
		{
			name:             "Error Close Batch",
			expectedWip:      expectedNewWipBatch,
			closeBatchParams: closeBatchParams,
			batches:          batches,
			closeBatchErr:    testErr,
			expectedErr:      fmt.Errorf("failed to close batch, err: %w", testErr),
			reprocessFullBatchResponse: &state.ProcessBatchResponse{
				NewStateRoot:     f.batch.stateRoot,
				NewLocalExitRoot: f.batch.localExitRoot,
				IsRomOOCError:    false,
			},
		},
		{
			name:             "Error Open Batch",
			expectedWip:      expectedNewWipBatch,
			closeBatchParams: closeBatchParams,
			batches:          batches,
			openBatchErr:     testErr,
			expectedErr:      fmt.Errorf("failed to open new batch, err: %w", testErr),
			reprocessFullBatchResponse: &state.ProcessBatchResponse{
				NewStateRoot:     f.batch.stateRoot,
				NewLocalExitRoot: f.batch.localExitRoot,
				IsRomOOCError:    false,
			},
		},
		{
			name:             "Success with closing non-empty batch",
			expectedWip:      expectedNewWipBatch,
			closeBatchParams: closeBatchParams,
			batches:          batches,
			reprocessFullBatchResponse: &state.ProcessBatchResponse{
				NewStateRoot:     f.batch.stateRoot,
				NewLocalExitRoot: f.batch.localExitRoot,
				IsRomOOCError:    false,
			},
		},
		{
			name:             "Success with closing empty batch",
			expectedWip:      &expectedNewWipEmptyBatch,
			closeBatchParams: closeBatchParamsForEmptyBatch,
			batches:          emptyBatchBatches,
			reprocessFullBatchResponse: &state.ProcessBatchResponse{
				NewStateRoot:     oldHash,
				NewLocalExitRoot: f.batch.localExitRoot,
				IsRomOOCError:    false,
			},
		},
		{
			name: "Forced Batches",
			forcedBatches: []state.ForcedBatch{{
				BlockNumber:       1,
				ForcedBatchNumber: 1,
				Sequencer:         seqAddr,
				GlobalExitRoot:    oldHash,
				RawTxsData:        nil,
				ForcedAt:          now(),
			}},
			expectedWip:      &expectedForcedNewWipBatch,
			closeBatchParams: closeBatchParams,
			batches:          batches,
			reprocessFullBatchResponse: &state.ProcessBatchResponse{
				NewStateRoot:     f.batch.stateRoot,
				NewLocalExitRoot: f.batch.localExitRoot,
				IsRomOOCError:    false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f.processRequest.GlobalExitRoot = oldHash
			f.processRequest.OldStateRoot = oldHash
			f.processRequest.BatchNumber = f.batch.batchNumber
			f.nextForcedBatches = tc.forcedBatches

			currTxs := txs
			if tc.closeBatchParams.StateRoot == oldHash {
				currTxs = nil
				f.batch.stateRoot = oldHash
				f.processRequest.Transactions = []byte{}
				defer func() {
					f.batch.stateRoot = newHash
					f.processRequest.Transactions = decodedBatchL2Data
				}()

				executorMock.On("ProcessBatch", ctx, f.processRequest, true).Return(tc.reprocessFullBatchResponse, tc.reprocessBatchErr).Once()
			}

			if tc.stateRootAndLERErr == nil {
				dbManagerMock.On("CloseBatch", ctx, tc.closeBatchParams).Return(tc.closeBatchErr).Once()
				dbManagerMock.On("GetBatchByNumber", ctx, f.batch.batchNumber, nil).Return(tc.batches[0], nilErr).Once()
				dbManagerMock.On("GetForkIDByBatchNumber", f.batch.batchNumber).Return(uint64(5)).Once()
				dbManagerMock.On("GetTransactionsByBatchNumber", ctx, f.batch.batchNumber).Return(currTxs, constants.EffectivePercentage, nilErr).Once()
				if tc.forcedBatches != nil && len(tc.forcedBatches) > 0 {
					processRequest := f.processRequest
					processRequest.BatchNumber = f.processRequest.BatchNumber + 1
					processRequest.OldStateRoot = newHash
					processRequest.Transactions = nil
					dbManagerMock.On("GetLastTrustedForcedBatchNumber", ctx, nil).Return(tc.forcedBatches[0].ForcedBatchNumber-1, nilErr).Once()
					dbManagerMock.On("ProcessForcedBatch", tc.forcedBatches[0].ForcedBatchNumber, processRequest).Return(tc.reprocessFullBatchResponse, nilErr).Once()
				}
				if tc.closeBatchErr == nil {
					dbManagerMock.On("BeginStateTransaction", ctx).Return(dbTxMock, nilErr).Once()
					dbManagerMock.On("OpenBatch", ctx, mock.Anything, dbTxMock).Return(tc.openBatchErr).Once()
					if tc.openBatchErr == nil {
						dbTxMock.On("Commit", ctx).Return(nilErr).Once()
					} else {
						dbTxMock.On("Rollback", ctx).Return(nilErr).Once()
					}
				}
				executorMock.On("ProcessBatch", ctx, f.processRequest, false).Return(tc.reprocessFullBatchResponse, tc.reprocessBatchErr).Once()
			}

			if tc.stateRootAndLERErr != nil {
				f.batch.stateRoot = state.ZeroHash
				f.batch.localExitRoot = state.ZeroHash
				defer func() {
					f.batch.stateRoot = newHash
					f.batch.localExitRoot = newHash
				}()
			}

			// act
			wipBatch, err := f.newWIPBatch(ctx)

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
			executorMock.AssertExpectations(t)
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
		name                    string
		batches                 []*state.Batch
		lastBatchNum            *uint64
		isBatchClosed           bool
		ger                     common.Hash
		getWIPBatchErr          error
		openBatchErr            error
		isBatchClosedErr        error
		getLastBatchErr         error
		expectedProcessingCtx   state.ProcessingContext
		expectedBatch           *WipBatch
		expectedErr             error
		getLastBatchByNumberErr error
		getLatestGERErr         error
	}{
		{
			name:          "Success Closed Batch",
			lastBatchNum:  &one,
			isBatchClosed: true,
			ger:           oldHash,
			batches:       batches,
			expectedBatch: &WipBatch{
				batchNumber:        one + 1,
				coinbase:           f.sequencerAddress,
				initialStateRoot:   oldHash,
				stateRoot:          oldHash,
				timestamp:          testNow(),
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
			name:          "Success Open Batch",
			lastBatchNum:  &one,
			isBatchClosed: false,
			batches:       batches,
			ger:           common.Hash{},
			expectedBatch: &WipBatch{
				batchNumber:        one,
				coinbase:           f.sequencerAddress,
				initialStateRoot:   oldHash,
				stateRoot:          oldHash,
				timestamp:          testNow(),
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
			name:            "Error Failed to get last batch",
			lastBatchNum:    nil,
			batches:         batches,
			isBatchClosed:   true,
			ger:             oldHash,
			getLastBatchErr: testErr,
			expectedErr:     fmt.Errorf("failed to get last batch, err: %w", testErr),
		},
		{
			name:             "Error Failed to check if batch is closed",
			lastBatchNum:     &one,
			batches:          batches,
			isBatchClosed:    true,
			ger:              oldHash,
			isBatchClosedErr: testErr,
			expectedErr:      fmt.Errorf("failed to check if batch is closed, err: %w", testErr),
		},
		{
			name:           "Error Failed to get work-in-progress batch",
			lastBatchNum:   &one,
			batches:        batches,
			isBatchClosed:  false,
			ger:            common.Hash{},
			getWIPBatchErr: testErr,
			expectedErr:    fmt.Errorf("failed to get work-in-progress batch, err: %w", testErr),
		},
		{
			name:          "Error Failed to open new batch",
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
		{
			name:          "Error Failed to get batch by number",
			lastBatchNum:  &one,
			batches:       batches,
			isBatchClosed: true,
			ger:           oldHash,
			expectedProcessingCtx: state.ProcessingContext{
				BatchNumber:    one + 1,
				Coinbase:       f.sequencerAddress,
				Timestamp:      testNow(),
				GlobalExitRoot: oldHash,
			},
			expectedErr:             fmt.Errorf("failed to get last batch, err: %w", testErr),
			getLastBatchByNumberErr: testErr,
		},
		{
			name:            "Error Failed to get latest GER",
			lastBatchNum:    &one,
			batches:         batches,
			isBatchClosed:   true,
			ger:             oldHash,
			expectedErr:     fmt.Errorf("failed to get latest ger, err: %w", testErr),
			getLatestGERErr: testErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			if tc.lastBatchNum == nil {
				dbManagerMock.Mock.On("GetLastBatch", ctx).Return(tc.batches[0], tc.getLastBatchErr).Once()
			} else {
				dbManagerMock.On("GetBatchByNumber", ctx, *tc.lastBatchNum, nil).Return(tc.batches[0], tc.getLastBatchByNumberErr).Once()
			}
			if tc.getLastBatchByNumberErr == nil {
				if tc.getLastBatchErr == nil {
					dbManagerMock.Mock.On("IsBatchClosed", ctx, *tc.lastBatchNum).Return(tc.isBatchClosed, tc.isBatchClosedErr).Once()
				}
				if tc.isBatchClosed {
					if tc.getLastBatchErr == nil && tc.isBatchClosedErr == nil {
						dbManagerMock.Mock.On("GetLatestGer", ctx, f.cfg.GERFinalityNumberOfBlocks).Return(state.GlobalExitRoot{GlobalExitRoot: tc.ger}, testNow(), tc.getLatestGERErr).Once()
						if tc.getLatestGERErr == nil {
							dbManagerMock.On("BeginStateTransaction", ctx).Return(dbTxMock, nil).Once()
							if tc.openBatchErr == nil {
								dbTxMock.On("Commit", ctx).Return(nil).Once()
							}
						}
					}

					if tc.getLastBatchErr == nil && tc.isBatchClosedErr == nil && tc.getLatestGERErr == nil {
						dbManagerMock.On("OpenBatch", ctx, tc.expectedProcessingCtx, dbTxMock).Return(tc.openBatchErr).Once()
					}

					if tc.expectedErr != nil && tc.openBatchErr != nil {
						dbTxMock.On("Rollback", ctx).Return(nil).Once()
					}
				} else {
					dbManagerMock.Mock.On("GetWIPBatch", ctx).Return(tc.expectedBatch, tc.getWIPBatchErr).Once()
				}
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
	ctx = context.Background()
	RawTxsData1 := make([]byte, 0, 2)
	RawTxsData1 = append(RawTxsData1, []byte(testBatchL2DataAsString)...)
	RawTxsData2 := make([]byte, 0, 2)
	RawTxsData2 = append(RawTxsData2, []byte(testBatchL2DataAsString)...)
	batchNumber := f.batch.batchNumber
	decodedBatchL2Data, err = hex.DecodeHex(testBatchL2DataAsString)
	require.NoError(t, err)

	txResp1 := &state.ProcessTransactionResponse{
		TxHash:    txHash,
		StateRoot: stateRootHashes[0],
	}

	txResp2 := &state.ProcessTransactionResponse{
		TxHash:    txHash2,
		StateRoot: stateRootHashes[1],
	}
	batchResponse1 := &state.ProcessBatchResponse{
		NewBatchNumber: f.batch.batchNumber + 1,
		Responses:      []*state.ProcessTransactionResponse{txResp1},
		NewStateRoot:   newHash,
	}

	batchResponse2 := &state.ProcessBatchResponse{
		NewBatchNumber: f.batch.batchNumber + 2,
		Responses:      []*state.ProcessTransactionResponse{txResp2},
		NewStateRoot:   newHash2,
	}
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
		forcedBatches                   []state.ForcedBatch
		getLastTrustedForcedBatchNumErr error
		expectedErr                     error
		expectedStoredTx                []transactionToStore
		processInBetweenForcedBatch     bool
		getForcedBatchError             error
	}{
		{
			name:          "Success",
			forcedBatches: []state.ForcedBatch{forcedBatch1, forcedBatch2},
			expectedStoredTx: []transactionToStore{
				{
					batchResponse: batchResponse1,
					batchNumber:   f.batch.batchNumber + 1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  stateRootHashes[0],
					isForcedBatch: true,
					response:      txResp1,
				},
				{
					batchResponse: batchResponse2,
					batchNumber:   f.batch.batchNumber + 2,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  stateRootHashes[1],
					isForcedBatch: true,
					response:      txResp2,
				},
			},
		},
		{
			name:                            "GetLastTrustedForcedBatchNumber_Error",
			forcedBatches:                   []state.ForcedBatch{forcedBatch1},
			getLastTrustedForcedBatchNumErr: testErr,
			expectedErr:                     fmt.Errorf("failed to get last trusted forced batch number, err: %s", testErr),
		},
		{
			name:          "Skip Already Processed Forced Batches",
			forcedBatches: []state.ForcedBatch{{ForcedBatchNumber: 1}},
		},
		{
			name: "Process In-Between Unprocessed Forced Batches",
			forcedBatches: []state.ForcedBatch{
				forcedBatch2,
				forcedBatch1,
			},
			expectedStoredTx: []transactionToStore{
				{
					batchResponse: batchResponse1,
					batchNumber:   f.batch.batchNumber + 1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  stateRootHashes[0],
					isForcedBatch: true,
					response:      txResp1,
				},
				{
					batchResponse: batchResponse2,
					batchNumber:   f.batch.batchNumber + 2,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  stateRootHashes[1],
					isForcedBatch: true,
					response:      txResp2,
				},
			},
			processInBetweenForcedBatch: true,
		},
		{
			name: "Error GetForcedBatch",
			forcedBatches: []state.ForcedBatch{
				forcedBatch2,
				forcedBatch1,
			},
			expectedErr:                 fmt.Errorf("failed to get in-between forced batch %d, err: %s", 2, testErr),
			getForcedBatchError:         testErr, // This will be used to simulate the error
			processInBetweenForcedBatch: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			var newStateRoot common.Hash
			stateRoot := oldHash
			storedTxs := make([]transactionToStore, 0)
			f.pendingTransactionsToStore = make(chan transactionToStore)
			if tc.expectedStoredTx != nil && len(tc.expectedStoredTx) > 0 {
				done = make(chan bool) // init a new done channel
				go func() {
					for tx := range f.pendingTransactionsToStore {
						storedTxs = append(storedTxs, tx)
						f.pendingTransactionsToStoreWG.Done()
					}
					done <- true // signal that the goroutine is done
				}()
			}
			f.nextForcedBatches = make([]state.ForcedBatch, len(tc.forcedBatches))
			copy(f.nextForcedBatches, tc.forcedBatches)
			internalBatchNumber := batchNumber
			dbManagerMock.On("GetLastTrustedForcedBatchNumber", ctx, nil).Return(uint64(1), tc.getLastTrustedForcedBatchNumErr).Once()
			tc.forcedBatches = f.sortForcedBatches(tc.forcedBatches)

			if tc.getLastTrustedForcedBatchNumErr == nil {
				for i, forcedBatch := range tc.forcedBatches {
					// Skip already processed forced batches.
					if forcedBatch.ForcedBatchNumber == 1 {
						continue
					}

					internalBatchNumber += 1
					processRequest := state.ProcessRequest{
						BatchNumber:    internalBatchNumber,
						OldStateRoot:   stateRootHashes[i],
						GlobalExitRoot: forcedBatch.GlobalExitRoot,
						Transactions:   forcedBatch.RawTxsData,
						Coinbase:       f.sequencerAddress,
						Timestamp:      now(),
						Caller:         stateMetrics.SequencerCallerLabel,
					}
					var currResp *state.ProcessBatchResponse
					if tc.expectedStoredTx == nil {
						currResp = &state.ProcessBatchResponse{
							NewStateRoot:   stateRootHashes[i+1],
							NewBatchNumber: internalBatchNumber,
						}
					} else {
						for _, storedTx := range tc.expectedStoredTx {
							if storedTx.batchNumber == internalBatchNumber {
								currResp = storedTx.batchResponse
								break
							}
						}
					}
					dbManagerMock.On("ProcessForcedBatch", forcedBatch.ForcedBatchNumber, processRequest).Return(currResp, nilErr).Once()
				}

				if tc.processInBetweenForcedBatch {
					dbManagerMock.On("GetForcedBatch", ctx, uint64(2), nil).Return(&forcedBatch1, tc.getForcedBatchError).Once()
				}
			}

			// act
			batchNumber, newStateRoot, err = f.processForcedBatches(ctx, batchNumber, stateRoot)

			// assert
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				if tc.expectedStoredTx != nil && len(tc.expectedStoredTx) > 0 {
					close(f.pendingTransactionsToStore) // ensure the channel is closed
					<-done                              // wait for the goroutine to finish
					f.pendingTransactionsToStoreWG.Wait()
					for i := range tc.expectedStoredTx {
						require.Equal(t, tc.expectedStoredTx[i], storedTxs[i])
					}
				}
				if len(tc.expectedStoredTx) > 0 {
					assert.Equal(t, stateRootHashes[len(stateRootHashes)-1], newStateRoot)
				}
				assert.Equal(t, batchNumber, internalBatchNumber)
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
		timestamp:          now(),
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
			name:        "Error BeginTransaction",
			beginTxErr:  testErr,
			expectedErr: fmt.Errorf("failed to begin state transaction to open batch, err: %w", testErr),
		},
		{
			name:         "Error OpenBatch",
			openBatchErr: testErr,
			expectedErr:  fmt.Errorf("failed to open new batch, err: %w", testErr),
		},
		{
			name:        "Error Commit",
			commitErr:   testErr,
			expectedErr: fmt.Errorf("failed to commit database transaction for opening a batch, err: %w", testErr),
		},
		{
			name:         "Error Rollback",
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
	effectivePercentages := constants.EffectivePercentage
	usedResources := getUsedBatchResources(f.batchConstraints, f.batch.remainingResources)
	receipt := ClosingBatchParameters{
		BatchNumber:          f.batch.batchNumber,
		StateRoot:            f.batch.stateRoot,
		LocalExitRoot:        f.batch.localExitRoot,
		BatchResources:       usedResources,
		Txs:                  txs,
		EffectivePercentages: effectivePercentages,
	}
	managerErr := fmt.Errorf("some err")
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
			name:        "Error Manager",
			managerErr:  managerErr,
			expectedErr: fmt.Errorf("failed to get transactions from transactions, err: %w", managerErr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			dbManagerMock.Mock.On("CloseBatch", ctx, receipt).Return(tc.managerErr).Once()
			dbManagerMock.Mock.On("GetTransactionsByBatchNumber", ctx, receipt.BatchNumber).Return(txs, effectivePercentages, tc.managerErr).Once()

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
			name:        "Error Manager",
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

func TestFinalizer_isDeadlineEncountered(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	now = testNow
	defer func() {
		now = time.Now
	}()
	testCases := []struct {
		name                        string
		nextForcedBatch             int64
		nextGER                     int64
		nextDelayedBatch            int64
		expected                    bool
		timestampResolutionDeadline bool
	}{
		{
			name:     "No deadlines",
			expected: false,
		},
		{
			name:            "Forced batch deadline",
			nextForcedBatch: now().Add(time.Second).Unix(),
			expected:        true,
		},
		{
			name:     "Global Exit Root deadline",
			nextGER:  now().Add(time.Second).Unix(),
			expected: true,
		},
		{
			name:             "Delayed batch deadline",
			nextDelayedBatch: now().Add(time.Second).Unix(),
			expected:         false,
		},
		{
			name:                        "Timestamp resolution deadline",
			timestampResolutionDeadline: true,
			expected:                    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f.nextForcedBatchDeadline = tc.nextForcedBatch
			f.nextGERDeadline = tc.nextGER
			if tc.expected == true {
				now = func() time.Time {
					return testNow().Add(time.Second * 2)
				}
			}

			// specifically for "Timestamp resolution deadline" test case
			if tc.timestampResolutionDeadline == true {
				// ensure that the batch is not empty and the timestamp is in the past
				f.batch.timestamp = now().Add(-f.cfg.TimestampResolution.Duration * 2)
				f.batch.countOfTxs = 1
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
	ctx = context.Background()
	txResponse := &state.ProcessTransactionResponse{TxHash: oldHash}
	result := &state.ProcessBatchResponse{
		UsedZkCounters: state.ZKCounters{CumulativeGasUsed: 1000},
		Responses:      []*state.ProcessTransactionResponse{txResponse},
	}
	remainingResources := state.BatchResources{
		ZKCounters: state.ZKCounters{CumulativeGasUsed: 9000},
		Bytes:      10000,
	}
	f.batch.remainingResources = remainingResources
	testCases := []struct {
		name                 string
		remaining            state.BatchResources
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
			remaining: state.BatchResources{
				Bytes: 0,
			},
			expectedErr:          state.ErrBatchResourceBytesUnderflow,
			expectedWorkerUpdate: true,
			expectedTxTracker:    &TxTracker{RawTx: []byte("test")},
		},
		{
			name: "ZkCounter Resource Exceeded",
			remaining: state.BatchResources{
				ZKCounters: state.ZKCounters{CumulativeGasUsed: 0},
			},
			expectedErr:          state.NewBatchRemainingResourcesUnderflowError(cumulativeGasErr, cumulativeGasErr.Error()),
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
			err := f.checkRemainingResources(result, tc.expectedTxTracker)

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

func TestFinalizer_handleTransactionError(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	nonce := uint64(0)
	txTracker := &TxTracker{Hash: txHash, From: senderAddr, Cost: big.NewInt(0)}
	testCases := []struct {
		name               string
		err                executor.RomError
		expectedDeleteCall bool
		updateTxStatus     pool.TxStatus
		expectedMoveCall   bool
		isRoomOOC          bool
	}{
		{
			name:               "Error OutOfCounters",
			err:                executor.RomError_ROM_ERROR_OUT_OF_COUNTERS_STEP,
			updateTxStatus:     pool.TxStatusInvalid,
			expectedDeleteCall: true,
			isRoomOOC:          true,
		},
		{
			name:             "Error IntrinsicInvalidNonce",
			err:              executor.RomError_ROM_ERROR_INTRINSIC_INVALID_NONCE,
			updateTxStatus:   pool.TxStatusFailed,
			expectedMoveCall: true,
		},
		{
			name:             "Error IntrinsicInvalidBalance",
			err:              executor.RomError_ROM_ERROR_INTRINSIC_INVALID_BALANCE,
			updateTxStatus:   pool.TxStatusFailed,
			expectedMoveCall: true,
		},
		{
			name:               "Error IntrinsicErrorChainId",
			err:                executor.RomError_ROM_ERROR_INTRINSIC_INVALID_CHAIN_ID,
			updateTxStatus:     pool.TxStatusFailed,
			expectedDeleteCall: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			if tc.expectedDeleteCall {
				workerMock.On("DeleteTx", txHash, senderAddr).Return()
				dbManagerMock.On("UpdateTxStatus", ctx, txHash, tc.updateTxStatus, false, mock.Anything).Return(nil).Once()
			}
			if tc.expectedMoveCall {
				workerMock.On("MoveTxToNotReady", txHash, senderAddr, &nonce, big.NewInt(0)).Return([]*TxTracker{
					{
						Hash: txHash2,
					},
				}).Once()

				dbManagerMock.On("UpdateTxStatus", ctx, txHash2, pool.TxStatusFailed, false, mock.Anything).Return(nil).Once()
			}

			result := &state.ProcessBatchResponse{
				IsRomOOCError: tc.isRoomOOC,
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
					senderAddr: {Nonce: &nonce, Balance: big.NewInt(0)},
				},
				Responses: []*state.ProcessTransactionResponse{
					{
						RomError: executor.RomErr(tc.err),
					},
				},
			}

			// act
			wg := f.handleProcessTransactionError(ctx, result, txTracker)
			if wg != nil {
				wg.Wait()
			}

			// assert
			workerMock.AssertExpectations(t)
		})
	}
}

func Test_processTransaction(t *testing.T) {
	f = setupFinalizer(true)
	gasUsed := uint64(100000)
	txTracker := &TxTracker{
		Hash:              txHash,
		From:              senderAddr,
		Nonce:             nonce1,
		BreakEvenGasPrice: breakEvenGasPrice,
		GasPrice:          breakEvenGasPrice,
		BatchResources: state.BatchResources{
			Bytes: 1000,
			ZKCounters: state.ZKCounters{
				CumulativeGasUsed: 500,
			},
		},
	}
	successfulTxResponse := &state.ProcessTransactionResponse{
		TxHash:    txHash,
		StateRoot: newHash,
		GasUsed:   gasUsed,
	}
	successfulBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: newHash,
		Responses: []*state.ProcessTransactionResponse{
			successfulTxResponse,
		},
		ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
			senderAddr: {
				Nonce: &nonce2,
			},
		},
	}
	outOfCountersErrBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: oldHash,
		Responses: []*state.ProcessTransactionResponse{
			{
				StateRoot: oldHash,
				RomError:  runtime.ErrOutOfCountersKeccak,
				GasUsed:   gasUsed,
			},
		},
		IsRomOOCError: true,
	}
	outOfCountersExecutorErrBatchResp := *outOfCountersErrBatchResp
	outOfCountersExecutorErrBatchResp.IsRomOOCError = false
	testCases := []struct {
		name                   string
		ctx                    context.Context
		tx                     *TxTracker
		expectedResponse       *state.ProcessBatchResponse
		executorErr            error
		expectedErr            error
		expectedStoredTx       transactionToStore
		expectedUpdateTxStatus pool.TxStatus
	}{
		{
			name:             "Successful transaction processing",
			ctx:              context.Background(),
			tx:               txTracker,
			expectedResponse: successfulBatchResp,
			expectedStoredTx: transactionToStore{
				txTracker:     txTracker,
				batchNumber:   f.batch.batchNumber,
				coinbase:      f.batch.coinbase,
				timestamp:     f.batch.timestamp,
				oldStateRoot:  newHash,
				batchResponse: successfulBatchResp,
				isForcedBatch: false,
				response:      successfulTxResponse,
			},
		},
		{
			name:                   "Out Of Counters err",
			ctx:                    context.Background(),
			tx:                     txTracker,
			expectedResponse:       outOfCountersErrBatchResp,
			expectedErr:            runtime.ErrOutOfCountersKeccak,
			expectedUpdateTxStatus: pool.TxStatusInvalid,
		},
		{
			name:             "Executor err",
			ctx:              context.Background(),
			tx:               txTracker,
			expectedResponse: &outOfCountersExecutorErrBatchResp,
			executorErr:      runtime.ErrOutOfCountersKeccak,
			expectedErr:      runtime.ErrOutOfCountersKeccak,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storedTxs := make([]transactionToStore, 0)
			f.pendingTransactionsToStore = make(chan transactionToStore, 1)
			if tc.expectedStoredTx.batchResponse != nil {
				done = make(chan bool) // init a new done channel
				go func() {
					for tx := range f.pendingTransactionsToStore {
						storedTxs = append(storedTxs, tx)
						f.pendingTransactionsToStoreWG.Done()
					}
					done <- true // signal that the goroutine is done
				}()
			}

			dbManagerMock.On("GetL1GasPrice").Return(uint64(1000000)).Once()
			executorMock.On("ProcessBatch", tc.ctx, mock.Anything, true).Return(tc.expectedResponse, tc.executorErr).Once()
			if tc.executorErr == nil {
				workerMock.On("DeleteTx", tc.tx.Hash, tc.tx.From).Return().Once()
				dbManagerMock.On("GetForkIDByBatchNumber", mock.Anything).Return(forkId5)
			}
			if tc.expectedErr == nil {
				workerMock.On("UpdateAfterSingleSuccessfulTxExecution", tc.tx.From, tc.expectedResponse.ReadWriteAddresses).Return([]*TxTracker{}).Once()
			}

			if tc.expectedUpdateTxStatus != "" {
				dbManagerMock.On("UpdateTxStatus", tc.ctx, txHash, tc.expectedUpdateTxStatus, false, mock.Anything).Return(nil).Once()
			}

			errWg, err := f.processTransaction(tc.ctx, tc.tx)

			if tc.expectedStoredTx.batchResponse != nil {
				close(f.pendingTransactionsToStore) // ensure the channel is closed
				<-done                              // wait for the goroutine to finish
				f.pendingTransactionsToStoreWG.Wait()
				require.Equal(t, tc.expectedStoredTx, storedTxs[0])
			}
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			if errWg != nil {
				errWg.Wait()
			}

			workerMock.AssertExpectations(t)
			dbManagerMock.AssertExpectations(t)
		})
	}
}

func Test_handleForcedTxsProcessResp(t *testing.T) {
	f = setupFinalizer(false)
	now = testNow
	defer func() {
		now = time.Now
	}()

	ctx = context.Background()
	txResponseOne := &state.ProcessTransactionResponse{
		TxHash:    txHash,
		StateRoot: newHash,
		RomError:  nil,
	}
	txResponseTwo := &state.ProcessTransactionResponse{
		TxHash:    common.HexToHash("0x02"),
		StateRoot: newHash2,
		RomError:  nil,
	}
	successfulBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: newHash,
		Responses: []*state.ProcessTransactionResponse{
			txResponseOne,
			txResponseTwo,
		},
	}
	txResponseReverted := &state.ProcessTransactionResponse{
		TxHash:    txHash,
		RomError:  runtime.ErrExecutionReverted,
		StateRoot: newHash,
	}
	revertedBatchResp := &state.ProcessBatchResponse{
		Responses: []*state.ProcessTransactionResponse{
			txResponseReverted,
		},
	}
	txResponseIntrinsicErr := &state.ProcessTransactionResponse{
		TxHash:    txHash,
		RomError:  runtime.ErrIntrinsicInvalidChainID,
		StateRoot: newHash,
	}
	intrinsicErrBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: newHash,
		Responses: []*state.ProcessTransactionResponse{
			txResponseOne,
			txResponseIntrinsicErr,
		},
	}

	testCases := []struct {
		name              string
		request           state.ProcessRequest
		result            *state.ProcessBatchResponse
		oldStateRoot      common.Hash
		expectedStoredTxs []transactionToStore
	}{
		{
			name: "Handle forced batch process response with successful transactions",
			request: state.ProcessRequest{
				BatchNumber:  1,
				Coinbase:     seqAddr,
				Timestamp:    now(),
				OldStateRoot: oldHash,
			},
			result:       successfulBatchResp,
			oldStateRoot: oldHash,
			expectedStoredTxs: []transactionToStore{
				{

					batchNumber:   1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  oldHash,
					response:      txResponseOne,
					isForcedBatch: true,
					batchResponse: successfulBatchResp,
				},
				{
					batchNumber:   1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  newHash,
					response:      txResponseTwo,
					isForcedBatch: true,
					batchResponse: successfulBatchResp,
				},
			},
		},
		{
			name: "Handle forced batch process response with reverted transactions",
			request: state.ProcessRequest{
				BatchNumber:  1,
				Coinbase:     seqAddr,
				Timestamp:    now(),
				OldStateRoot: oldHash,
			},
			result:       revertedBatchResp,
			oldStateRoot: oldHash,
			expectedStoredTxs: []transactionToStore{
				{
					batchNumber:   1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  oldHash,
					response:      txResponseReverted,
					isForcedBatch: true,
					batchResponse: revertedBatchResp,
				}},
		},
		{
			name: "Handle forced batch process response with intrinsic ROM err",
			request: state.ProcessRequest{
				BatchNumber:  1,
				Coinbase:     seqAddr,
				Timestamp:    now(),
				OldStateRoot: oldHash,
			},

			result:       intrinsicErrBatchResp,
			oldStateRoot: oldHash,
			expectedStoredTxs: []transactionToStore{
				{
					batchNumber:   1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  oldHash,
					response:      txResponseOne,
					isForcedBatch: true,
					batchResponse: intrinsicErrBatchResp,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storedTxs := make([]transactionToStore, 0)
			f.pendingTransactionsToStore = make(chan transactionToStore)

			// Mock storeProcessedTx to store txs into the storedTxs slice
			go func() {
				for tx := range f.pendingTransactionsToStore {
					storedTxs = append(storedTxs, tx)
					f.pendingTransactionsToStoreWG.Done()
				}
			}()

			f.handleForcedTxsProcessResp(ctx, tc.request, tc.result, tc.oldStateRoot)

			f.pendingTransactionsToStoreWG.Wait()
			require.Nil(t, err)
			require.Equal(t, len(tc.expectedStoredTxs), len(storedTxs))
			for i := 0; i < len(tc.expectedStoredTxs); i++ {
				expectedTx := tc.expectedStoredTxs[i]
				actualTx := storedTxs[i]
				require.Equal(t, expectedTx, actualTx)
			}
		})
	}
}

func TestFinalizer_storeProcessedTx(t *testing.T) {
	f = setupFinalizer(false)
	testCases := []struct {
		name              string
		batchNum          uint64
		coinbase          common.Address
		timestamp         time.Time
		previousStateRoot common.Hash
		txResponse        *state.ProcessTransactionResponse
		isForcedBatch     bool
		expectedTxToStore transactionToStore
	}{
		{
			name:              "Normal transaction",
			batchNum:          1,
			coinbase:          seqAddr,
			timestamp:         time.Now(),
			previousStateRoot: oldHash,
			txResponse: &state.ProcessTransactionResponse{
				TxHash: txHash,
			},
			isForcedBatch: false,
			expectedTxToStore: transactionToStore{
				batchNumber:  1,
				coinbase:     seqAddr,
				timestamp:    now(),
				oldStateRoot: oldHash,
				response: &state.ProcessTransactionResponse{
					TxHash: txHash,
				},
				isForcedBatch: false,
			},
		},
		{
			name:              "Forced transaction",
			batchNum:          1,
			coinbase:          seqAddr,
			timestamp:         time.Now(),
			previousStateRoot: oldHash,
			txResponse: &state.ProcessTransactionResponse{
				TxHash: txHash2,
			},
			isForcedBatch: true,
			expectedTxToStore: transactionToStore{
				batchNumber:  1,
				coinbase:     seqAddr,
				timestamp:    now(),
				oldStateRoot: oldHash,
				response: &state.ProcessTransactionResponse{
					TxHash: txHash2,
				},
				isForcedBatch: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			dbManagerMock.On("StoreProcessedTxAndDeleteFromPool", ctx, tc.expectedTxToStore).Return(nilErr)

			// act
			f.storeProcessedTx(ctx, tc.expectedTxToStore)

			// assert
			dbManagerMock.AssertExpectations(t)
		})
	}
}

func TestFinalizer_updateWorkerAfterSuccessfulProcessing(t *testing.T) {
	testCases := []struct {
		name                  string
		txTracker             *TxTracker
		processBatchResponse  *state.ProcessBatchResponse
		expectedDeleteTxCount int
		expectedUpdateCount   int
	}{
		{
			name: "Successful update with one read-write address",
			txTracker: &TxTracker{
				Hash:  oldHash,
				From:  senderAddr,
				Nonce: nonce1,
			},
			processBatchResponse: &state.ProcessBatchResponse{
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
					senderAddr: {
						Address: senderAddr,
						Nonce:   &nonce1,
					},
				},
			},
			expectedDeleteTxCount: 1,
			expectedUpdateCount:   1,
		},
		{
			name: "Successful update with multiple read-write addresses",
			txTracker: &TxTracker{
				Hash:  oldHash,
				From:  senderAddr,
				Nonce: 1,
			},
			processBatchResponse: &state.ProcessBatchResponse{
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
					senderAddr: {
						Address: senderAddr,
						Nonce:   &nonce1,
					},
					receiverAddr: {
						Address: receiverAddr,
						Nonce:   &nonce2,
					},
				},
			},
			expectedDeleteTxCount: 1,
			expectedUpdateCount:   2,
		},
		{
			name: "No update when no read-write addresses provided",
			txTracker: &TxTracker{
				Hash:  oldHash,
				From:  senderAddr,
				Nonce: 1,
			},
			processBatchResponse: &state.ProcessBatchResponse{
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{},
			},
			expectedDeleteTxCount: 1,
			expectedUpdateCount:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			finalizerInstance := setupFinalizer(false)
			workerMock.On("DeleteTx", tc.txTracker.Hash, tc.txTracker.From).Times(tc.expectedDeleteTxCount)
			txsToDelete := make([]*TxTracker, 0, len(tc.processBatchResponse.ReadWriteAddresses))
			for _, infoReadWrite := range tc.processBatchResponse.ReadWriteAddresses {
				txsToDelete = append(txsToDelete, &TxTracker{
					Hash:         oldHash,
					From:         infoReadWrite.Address,
					FailedReason: &testErrStr,
				})
			}
			workerMock.On("UpdateAfterSingleSuccessfulTxExecution", tc.txTracker.From, tc.processBatchResponse.ReadWriteAddresses).
				Return(txsToDelete)
			if tc.expectedUpdateCount > 0 {
				dbManagerMock.On("UpdateTxStatus", mock.Anything, mock.Anything, pool.TxStatusFailed, false, mock.Anything).Times(tc.expectedUpdateCount).Return(nil)
			}

			// act
			finalizerInstance.updateWorkerAfterSuccessfulProcessing(ctx, tc.txTracker, tc.processBatchResponse)

			// assert
			workerMock.AssertExpectations(t)
			dbManagerMock.AssertExpectations(t)
		})
	}
}

func TestFinalizer_reprocessFullBatch(t *testing.T) {
	successfulResult := &state.ProcessBatchResponse{
		NewStateRoot: newHash,
	}
	roomOOCErrResult := &state.ProcessBatchResponse{
		NewStateRoot:  newHash,
		IsRomOOCError: true,
	}
	testCases := []struct {
		name                     string
		batchNum                 uint64
		oldStateRoot             common.Hash
		mockGetBatchByNumber     *state.Batch
		mockGetBatchByNumberErr  error
		expectedExecutorResponse *state.ProcessBatchResponse
		expectedResult           *state.ProcessBatchResponse
		expectedDecodeErr        error
		expectedExecutorErr      error
		expectedError            error
	}{
		{
			name:     "Success",
			batchNum: 1,
			mockGetBatchByNumber: &state.Batch{
				BatchNumber:    1,
				BatchL2Data:    decodedBatchL2Data,
				GlobalExitRoot: oldHash,
				Coinbase:       common.Address{},
				Timestamp:      time.Now(),
			},
			expectedExecutorResponse: successfulResult,
			expectedResult:           successfulResult,
		},
		{
			name:                    "Error while getting batch by number",
			batchNum:                1,
			mockGetBatchByNumberErr: errors.New("database err"),
			expectedError:           fmt.Errorf("failed to get batch by number, err: database err"),
		},
		{
			name:     "Error decoding BatchL2Data",
			batchNum: 1,
			mockGetBatchByNumber: &state.Batch{
				BatchNumber:    1,
				BatchL2Data:    []byte("invalidBatchL2Data"),
				GlobalExitRoot: oldHash,
				Coinbase:       common.Address{},
				Timestamp:      time.Now(),
			},
			expectedDecodeErr: fmt.Errorf("reprocessFullBatch: error decoding BatchL2Data before reprocessing full batch: 1. Error: %v", errors.New("invalid data")),
			expectedError:     fmt.Errorf("reprocessFullBatch: error decoding BatchL2Data before reprocessing full batch: 1. Error: %v", errors.New("invalid data")),
		},
		{
			name:     "Error processing batch",
			batchNum: 1,
			mockGetBatchByNumber: &state.Batch{
				BatchNumber:    1,
				BatchL2Data:    decodedBatchL2Data,
				GlobalExitRoot: oldHash,
				Coinbase:       common.Address{},
				Timestamp:      time.Now(),
			},
			expectedExecutorErr: errors.New("processing err"),
			expectedError:       errors.New("processing err"),
		},
		{
			name:     "RomOOCError",
			batchNum: 1,
			mockGetBatchByNumber: &state.Batch{
				BatchNumber:    1,
				BatchL2Data:    decodedBatchL2Data,
				GlobalExitRoot: oldHash,
				Coinbase:       common.Address{},
				Timestamp:      time.Now(),
			},
			expectedExecutorResponse: roomOOCErrResult,
			expectedError:            fmt.Errorf("failed to process batch because OutOfCounters error"),
		},
		{
			name:     "Reprocessed batch has different state root",
			batchNum: 1,
			mockGetBatchByNumber: &state.Batch{
				BatchNumber:    1,
				BatchL2Data:    decodedBatchL2Data,
				GlobalExitRoot: oldHash,
				Coinbase:       common.Address{},
				Timestamp:      time.Now(),
			},
			expectedExecutorResponse: &state.ProcessBatchResponse{
				NewStateRoot: newHash2,
			},
			expectedError: fmt.Errorf("batchNumber: 1, reprocessed batch has different state root, expected: %s, got: %s", newHash.Hex(), newHash2.Hex()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f := setupFinalizer(true)
			dbManagerMock.On("GetBatchByNumber", context.Background(), tc.batchNum, nil).Return(tc.mockGetBatchByNumber, tc.mockGetBatchByNumberErr).Once()
			if tc.name != "Error while getting batch by number" {
				dbManagerMock.On("GetForkIDByBatchNumber", f.batch.batchNumber).Return(uint64(5)).Once()
			}
			if tc.mockGetBatchByNumberErr == nil && tc.expectedDecodeErr == nil {
				executorMock.On("ProcessBatch", context.Background(), mock.Anything, false).Return(tc.expectedExecutorResponse, tc.expectedExecutorErr)
			}

			// act
			result, err := f.reprocessFullBatch(context.Background(), tc.batchNum, newHash)

			// assert
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
			dbManagerMock.AssertExpectations(t)
			executorMock.AssertExpectations(t)
		})
	}
}

func TestFinalizer_getLastBatchNumAndOldStateRoot(t *testing.T) {
	f := setupFinalizer(false)
	testCases := []struct {
		name              string
		mockBatches       []*state.Batch
		mockError         error
		expectedBatchNum  uint64
		expectedStateRoot common.Hash
		expectedError     error
	}{
		{
			name: "Success with two batches",
			mockBatches: []*state.Batch{
				{BatchNumber: 2, StateRoot: common.BytesToHash([]byte("stateRoot2"))},
				{BatchNumber: 1, StateRoot: common.BytesToHash([]byte("stateRoot1"))},
			},
			mockError:         nil,
			expectedBatchNum:  2,
			expectedStateRoot: common.BytesToHash([]byte("stateRoot1")),
			expectedError:     nil,
		},
		{
			name: "Success with one batch",
			mockBatches: []*state.Batch{
				{BatchNumber: 1, StateRoot: common.BytesToHash([]byte("stateRoot1"))},
			},
			mockError:         nil,
			expectedBatchNum:  1,
			expectedStateRoot: common.BytesToHash([]byte("stateRoot1")),
			expectedError:     nil,
		},
		{
			name:              "Error while getting batches",
			mockBatches:       nil,
			mockError:         errors.New("database err"),
			expectedBatchNum:  0,
			expectedStateRoot: common.Hash{},
			expectedError:     errors.New("failed to get last 2 batches, err: database err"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			dbManagerMock.On("GetLastNBatches", context.Background(), uint(2)).Return(tc.mockBatches, tc.mockError).Once()

			// act
			batchNum, stateRoot, err := f.getLastBatchNumAndOldStateRoot(context.Background())

			// assert
			assert.Equal(t, tc.expectedBatchNum, batchNum)
			assert.Equal(t, tc.expectedStateRoot, stateRoot)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			dbManagerMock.AssertExpectations(t)
		})
	}
}

func TestFinalizer_getOldStateRootFromBatches(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	testCases := []struct {
		name              string
		batches           []*state.Batch
		expectedStateRoot common.Hash
	}{
		{
			name: "Success with two batches",
			batches: []*state.Batch{
				{BatchNumber: 2, StateRoot: common.BytesToHash([]byte("stateRoot2"))},
				{BatchNumber: 1, StateRoot: common.BytesToHash([]byte("stateRoot1"))},
			},
			expectedStateRoot: common.BytesToHash([]byte("stateRoot1")),
		},
		{
			name: "Success with one batch",
			batches: []*state.Batch{
				{BatchNumber: 1, StateRoot: common.BytesToHash([]byte("stateRoot1"))},
			},
			expectedStateRoot: common.BytesToHash([]byte("stateRoot1")),
		},
		{
			name:              "Success with no batches",
			batches:           []*state.Batch{},
			expectedStateRoot: common.Hash{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			stateRoot := f.getOldStateRootFromBatches(tc.batches)

			// assert
			assert.Equal(t, tc.expectedStateRoot, stateRoot)
		})
	}
}

func TestFinalizer_isBatchAlmostFull(t *testing.T) {
	// arrange
	testCases := []struct {
		name               string
		modifyResourceFunc func(resources state.BatchResources) state.BatchResources
		expectedResult     bool
	}{
		{
			name: "Is ready - MaxBatchBytesSize",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.Bytes = f.getConstraintThresholdUint64(bc.MaxBatchBytesSize) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxBatchBytesSize",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.Bytes = f.getConstraintThresholdUint64(bc.MaxBatchBytesSize) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxCumulativeGasUsed",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.CumulativeGasUsed = f.getConstraintThresholdUint64(bc.MaxCumulativeGasUsed) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxCumulativeGasUsed",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.CumulativeGasUsed = f.getConstraintThresholdUint64(bc.MaxCumulativeGasUsed) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxSteps",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedSteps = f.getConstraintThresholdUint32(bc.MaxSteps) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxSteps",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedSteps = f.getConstraintThresholdUint32(bc.MaxSteps) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxPoseidonPaddings",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedPoseidonPaddings = f.getConstraintThresholdUint32(bc.MaxPoseidonPaddings) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxPoseidonPaddings",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedPoseidonPaddings = f.getConstraintThresholdUint32(bc.MaxPoseidonPaddings) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxBinaries",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedBinaries = f.getConstraintThresholdUint32(bc.MaxBinaries) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxBinaries",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedBinaries = f.getConstraintThresholdUint32(bc.MaxBinaries) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxKeccakHashes",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedKeccakHashes = f.getConstraintThresholdUint32(bc.MaxKeccakHashes) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxKeccakHashes",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedKeccakHashes = f.getConstraintThresholdUint32(bc.MaxKeccakHashes) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxArithmetics",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedArithmetics = f.getConstraintThresholdUint32(bc.MaxArithmetics) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxArithmetics",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedArithmetics = f.getConstraintThresholdUint32(bc.MaxArithmetics) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxMemAligns",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedMemAligns = f.getConstraintThresholdUint32(bc.MaxMemAligns) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxMemAligns",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.UsedMemAligns = f.getConstraintThresholdUint32(bc.MaxMemAligns) + 1
				return resources
			},
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f = setupFinalizer(true)
			maxRemainingResource := getMaxRemainingResources(bc)
			f.batch.remainingResources = tc.modifyResourceFunc(maxRemainingResource)

			// act
			result := f.isBatchAlmostFull()

			// assert
			assert.Equal(t, tc.expectedResult, result)
			if tc.expectedResult {
				assert.Equal(t, state.BatchAlmostFullClosingReason, f.batch.closingReason)
			} else {
				assert.Equal(t, state.EmptyClosingReason, f.batch.closingReason)
			}
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
	expected := now().Unix() + int64(f.cfg.ForcedBatchDeadlineTimeout.Duration.Seconds())

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
	expected := now().Unix() + int64(f.cfg.GERDeadlineTimeout.Duration.Seconds())

	// act
	f.setNextGERDeadline()

	// assert
	assert.Equal(t, expected, f.nextGERDeadline)
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
	assert.Equal(t, remainingResources.ZKCounters.CumulativeGasUsed, bc.MaxCumulativeGasUsed)
	assert.Equal(t, remainingResources.ZKCounters.UsedKeccakHashes, bc.MaxKeccakHashes)
	assert.Equal(t, remainingResources.ZKCounters.UsedPoseidonHashes, bc.MaxPoseidonHashes)
	assert.Equal(t, remainingResources.ZKCounters.UsedPoseidonPaddings, bc.MaxPoseidonPaddings)
	assert.Equal(t, remainingResources.ZKCounters.UsedMemAligns, bc.MaxMemAligns)
	assert.Equal(t, remainingResources.ZKCounters.UsedArithmetics, bc.MaxArithmetics)
	assert.Equal(t, remainingResources.ZKCounters.UsedBinaries, bc.MaxBinaries)
	assert.Equal(t, remainingResources.ZKCounters.UsedSteps, bc.MaxSteps)
	assert.Equal(t, remainingResources.Bytes, bc.MaxBatchBytesSize)
}

func Test_isBatchFull(t *testing.T) {
	f = setupFinalizer(true)

	testCases := []struct {
		name            string
		batchCountOfTxs int
		maxTxsPerBatch  uint64
		expected        bool
	}{
		{
			name:            "Batch is not full",
			batchCountOfTxs: 5,
			maxTxsPerBatch:  10,
			expected:        false,
		},
		{
			name:            "Batch is full",
			batchCountOfTxs: 10,
			maxTxsPerBatch:  10,
			expected:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f.batch.countOfTxs = tc.batchCountOfTxs
			f.batchConstraints.MaxTxsPerBatch = tc.maxTxsPerBatch

			assert.Equal(t, tc.expected, f.isBatchFull())
			if tc.expected == true {
				assert.Equal(t, state.BatchFullClosingReason, f.batch.closingReason)
			}
		})
	}
}

func Test_sortForcedBatches(t *testing.T) {
	f = setupFinalizer(false)

	testCases := []struct {
		name     string
		input    []state.ForcedBatch
		expected []state.ForcedBatch
	}{
		{
			name:     "Empty slice",
			input:    []state.ForcedBatch{},
			expected: []state.ForcedBatch{},
		},
		{
			name:     "Single item slice",
			input:    []state.ForcedBatch{{ForcedBatchNumber: 5}},
			expected: []state.ForcedBatch{{ForcedBatchNumber: 5}},
		},
		{
			name:     "Multiple items unsorted",
			input:    []state.ForcedBatch{{ForcedBatchNumber: 5}, {ForcedBatchNumber: 3}, {ForcedBatchNumber: 9}},
			expected: []state.ForcedBatch{{ForcedBatchNumber: 3}, {ForcedBatchNumber: 5}, {ForcedBatchNumber: 9}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := f.sortForcedBatches(testCase.input)
			assert.Equal(t, testCase.expected, result, "They should be equal")
		})
	}
}

func setupFinalizer(withWipBatch bool) *finalizer {
	wipBatch := new(WipBatch)
	dbManagerMock = new(DbManagerMock)
	executorMock = new(StateMock)
	workerMock = new(WorkerMock)
	dbTxMock = new(DbTxMock)
	if withWipBatch {
		decodedBatchL2Data, err = hex.DecodeHex(testBatchL2DataAsString)
		if err != nil {
			panic(err)
		}
		wipBatch = &WipBatch{
			batchNumber:        1,
			coinbase:           seqAddr,
			initialStateRoot:   oldHash,
			stateRoot:          newHash,
			localExitRoot:      newHash,
			timestamp:          now(),
			globalExitRoot:     oldHash,
			remainingResources: getMaxRemainingResources(bc),
			closingReason:      state.EmptyClosingReason,
		}
	}
	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		panic(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)
	return &finalizer{
		cfg:                  cfg,
		effectiveGasPriceCfg: effectiveGasPriceCfg,
		closingSignalCh:      closingSignalCh,
		isSynced:             isSynced,
		sequencerAddress:     seqAddr,
		worker:               workerMock,
		dbManager:            dbManagerMock,
		executor:             executorMock,
		batch:                wipBatch,
		batchConstraints:     bc,
		processRequest:       state.ProcessRequest{},
		sharedResourcesMux:   new(sync.RWMutex),
		lastGERHash:          common.Hash{},
		// closing signals
		nextGER:                                 common.Hash{},
		nextGERDeadline:                         0,
		nextGERMux:                              new(sync.RWMutex),
		nextForcedBatches:                       make([]state.ForcedBatch, 0),
		nextForcedBatchDeadline:                 0,
		nextForcedBatchesMux:                    new(sync.RWMutex),
		handlingL2Reorg:                         false,
		eventLog:                                eventLog,
		maxBreakEvenGasPriceDeviationPercentage: big.NewInt(10),
		pendingTransactionsToStore:              make(chan transactionToStore, bc.MaxTxsPerBatch*pendingTxsBufferSizeMultiplier),
		pendingTransactionsToStoreWG:            new(sync.WaitGroup),
		pendingTransactionsToStoreMux:           new(sync.RWMutex),
		storedFlushID:                           0,
		storedFlushIDCond:                       sync.NewCond(new(sync.Mutex)),
		proverID:                                "",
		lastPendingFlushID:                      0,
		pendingFlushIDCond:                      sync.NewCond(new(sync.Mutex)),
	}
}
