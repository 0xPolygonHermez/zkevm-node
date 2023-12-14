package sequencer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	poolPackage "github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	statePackage "github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
)

const (
	pendingL2BlocksBufferSize = 100
	changeL2BlockSize         = 9 //1 byte (tx type = 0B) + 4 bytes for deltaTimestamp + 4 for l1InfoTreeIndex
)

var (
	now            = time.Now
	mockL1InfoRoot = common.Hash{}
)

// finalizer represents the finalizer component of the sequencer.
type finalizer struct {
	cfg                     FinalizerCfg
	isSynced                func(ctx context.Context) bool
	sequencerAddress        common.Address
	worker                  workerInterface
	pool                    txPool
	state                   stateInterface
	etherman                etherman
	wipBatch                *Batch
	wipL2Block              *L2Block
	batchConstraints        statePackage.BatchConstraintsCfg
	reprocessFullBatchError atomic.Bool
	// closing signals
	closingSignalCh ClosingSignalCh
	// GER
	currentGERHash  common.Hash // GER of the current WIP batch
	previousGERHash common.Hash // GER of the batch previous to the current WIP batch
	nextGER         common.Hash
	nextGERDeadline int64
	nextGERMux      *sync.Mutex
	// forced batches
	nextForcedBatches       []statePackage.ForcedBatch
	nextForcedBatchDeadline int64
	nextForcedBatchesMux    *sync.Mutex
	// L1InfoTree
	lastL1InfoTreeValid bool
	lastL1InfoTree      statePackage.L1InfoTreeExitRootStorageEntry
	lastL1InfoTreeMux   *sync.Mutex
	lastL1InfoTreeCond  *sync.Cond
	// L2 reorg
	handlingL2Reorg bool
	// event log
	eventLog *event.EventLog
	// effective gas price calculation instance
	effectiveGasPrice *poolPackage.EffectiveGasPrice
	// pending L2 blocks to be processed (executor)
	pendingL2BlocksToProcess   chan *L2Block
	pendingL2BlocksToProcessWG *sync.WaitGroup
	// pending L2 blocks to store in the state
	pendingL2BlocksToStore   chan *L2Block
	pendingL2BlocksToStoreWG *sync.WaitGroup
	// executer flushid control
	proverID           string
	storedFlushID      uint64
	storedFlushIDCond  *sync.Cond //Condition to wait until storedFlushID has been updated
	lastPendingFlushID uint64
	pendingFlushIDCond *sync.Cond
	// stream server
	streamServer *datastreamer.StreamServer
	dataToStream chan state.DSL2FullBlock
}

// newFinalizer returns a new instance of Finalizer.
func newFinalizer(
	cfg FinalizerCfg,
	poolCfg poolPackage.Config,
	worker workerInterface,
	pool txPool,
	state stateInterface,
	etherman etherman,
	sequencerAddr common.Address,
	isSynced func(ctx context.Context) bool,
	closingSignalCh ClosingSignalCh,
	batchConstraints statePackage.BatchConstraintsCfg,
	eventLog *event.EventLog,
	streamServer *datastreamer.StreamServer,
	dataToStream chan state.DSL2FullBlock,
) *finalizer {
	f := finalizer{
		cfg:              cfg,
		isSynced:         isSynced,
		sequencerAddress: sequencerAddr,
		worker:           worker,
		pool:             pool,
		state:            state,
		etherman:         etherman,
		batchConstraints: batchConstraints,
		// closing signals
		closingSignalCh: closingSignalCh,
		// GER //TODO: Delete GER updates as in ETROG it's not used
		currentGERHash:  statePackage.ZeroHash,
		previousGERHash: statePackage.ZeroHash,
		nextGER:         statePackage.ZeroHash,
		nextGERDeadline: 0,
		nextGERMux:      new(sync.Mutex),
		// forced batches
		nextForcedBatches:       make([]statePackage.ForcedBatch, 0),
		nextForcedBatchDeadline: 0,
		nextForcedBatchesMux:    new(sync.Mutex),
		// L1InfoTree
		lastL1InfoTreeValid: false,
		lastL1InfoTreeMux:   new(sync.Mutex),
		lastL1InfoTreeCond:  sync.NewCond(&sync.Mutex{}),
		// L2 reorg
		handlingL2Reorg: false,
		// event log
		eventLog: eventLog,
		// effective gas price calculation instance
		effectiveGasPrice: poolPackage.NewEffectiveGasPrice(poolCfg.EffectiveGasPrice, poolCfg.DefaultMinGasPriceAllowed),
		// pending L2 blocks to be processed (executor)
		pendingL2BlocksToProcess:   make(chan *L2Block, pendingL2BlocksBufferSize), //TODO: review buffer size
		pendingL2BlocksToProcessWG: new(sync.WaitGroup),
		// pending L2 blocks to store in the state
		pendingL2BlocksToStore:   make(chan *L2Block, pendingL2BlocksBufferSize), //TODO: review buffer size
		pendingL2BlocksToStoreWG: new(sync.WaitGroup),
		storedFlushID:            0,
		// executer flushid control
		proverID:           "",
		storedFlushIDCond:  sync.NewCond(&sync.Mutex{}),
		lastPendingFlushID: 0,
		pendingFlushIDCond: sync.NewCond(&sync.Mutex{}),
		// stream server
		streamServer: streamServer,
		dataToStream: dataToStream,
	}

	f.reprocessFullBatchError.Store(false)

	return &f
}

// Start starts the finalizer.
func (f *finalizer) Start(ctx context.Context) {
	// Init mockL1InfoRoot to a mock value since it must be different to {0,0,...,0}
	for i := 0; i < len(mockL1InfoRoot); i++ {
		mockL1InfoRoot[i] = byte(i)
	}

	// Update L1InfoRoot
	go f.checkL1InfoRootUpdate(ctx)

	// Get the last batch if still wip or opens a new one
	f.initWIPBatch(ctx)

	// Initializes the wip L2 block
	f.initWIPL2Block(ctx)

	// Closing signals receiver
	go f.listenForClosingSignals(ctx)

	// Update the prover id and flush id
	go f.updateProverIdAndFlushId(ctx)

	// Process L2 Blocks
	go f.processPendingL2Blocks(ctx)

	// Store L2 Blocks
	go f.storePendingL2Blocks(ctx)

	// Processing transactions and finalizing batches
	f.finalizeBatches(ctx)
}

// updateProverIdAndFlushId updates the prover id and flush id
func (f *finalizer) updateProverIdAndFlushId(ctx context.Context) {
	for {
		f.pendingFlushIDCond.L.Lock()
		// f.storedFlushID is >= than f.lastPendingFlushID, this means all pending txs (flushid) are stored by the executor.
		// We are "synced" with the flush id, therefore we need to wait for new tx (new pending flush id to be stored by the executor)
		for f.storedFlushID >= f.lastPendingFlushID {
			f.pendingFlushIDCond.Wait()
		}
		f.pendingFlushIDCond.L.Unlock()

		for f.storedFlushID < f.lastPendingFlushID { //TODO: review this loop as could be is pulling all the time, no sleep
			storedFlushID, proverID, err := f.state.GetStoredFlushID(ctx)
			if err != nil {
				log.Errorf("failed to get stored flush id, Err: %v", err)
			} else {
				if storedFlushID != f.storedFlushID {
					// Check if prover/Executor has been restarted
					f.checkIfProverRestarted(proverID)

					// Update f.storeFlushID and signal condition f.storedFlushIDCond
					f.storedFlushIDCond.L.Lock()
					f.storedFlushID = storedFlushID
					f.storedFlushIDCond.Broadcast()
					f.storedFlushIDCond.L.Unlock()
				}
			}
		}
	}
}

func (f *finalizer) checkL1InfoRootUpdate(ctx context.Context) {
	var (
		firstL1InfoRootUpdate = true
		firstSleepSkipped     = false
	)

	for {
		if firstSleepSkipped {
			time.Sleep(f.cfg.WaitForCheckingL1InfoRoot.Duration)
		} else {
			firstSleepSkipped = true
		}

		lastL1BlockNumber, err := f.etherman.GetLatestBlockNumber(ctx)
		if err != nil {
			log.Errorf("error getting latest L1 block number: %v", err)
		}

		maxBlockNumber := uint64(0)
		if f.cfg.L1InfoRootFinalityNumberOfBlocks <= lastL1BlockNumber {
			maxBlockNumber = lastL1BlockNumber - f.cfg.L1InfoRootFinalityNumberOfBlocks
		}

		l1InfoRoot, err := f.state.GetLatestL1InfoRoot(ctx, maxBlockNumber)
		if err != nil {
			log.Errorf("error checking latest L1InfoRoot: %v", err)
			continue
		}

		if firstL1InfoRootUpdate || l1InfoRoot.L1InfoTreeIndex > f.lastL1InfoTree.L1InfoTreeIndex {
			firstL1InfoRootUpdate = false

			log.Debugf("received new L1InfoRoot. L1InfoTreeIndex: %d", l1InfoRoot.L1InfoTreeIndex)

			f.lastL1InfoTreeMux.Lock()
			f.lastL1InfoTree = l1InfoRoot
			f.lastL1InfoTreeMux.Unlock()

			if !f.lastL1InfoTreeValid {
				f.lastL1InfoTreeCond.L.Lock()
				f.lastL1InfoTreeValid = true
				f.lastL1InfoTreeCond.Broadcast()
				f.lastL1InfoTreeCond.L.Unlock()
			}
		}
	}
}

// listenForClosingSignals listens for signals for the batch and sets the deadline for when they need to be closed.
func (f *finalizer) listenForClosingSignals(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Infof("finalizer closing signal listener received context done, Err: %s", ctx.Err())
			return
		// ForcedBatch ch
		case fb := <-f.closingSignalCh.ForcedBatchCh:
			log.Debugf("finalizer received forced batch at block number: %v", fb.BlockNumber)

			f.nextForcedBatchesMux.Lock()
			f.nextForcedBatches = f.sortForcedBatches(append(f.nextForcedBatches, fb))
			if f.nextForcedBatchDeadline == 0 {
				f.setNextForcedBatchDeadline()
			}
			f.nextForcedBatchesMux.Unlock()
		// GlobalExitRoot ch
		case ger := <-f.closingSignalCh.GERCh:
			log.Debugf("finalizer received global exit root: %s", ger.String())
			f.nextGERMux.Lock()
			f.nextGER = ger
			if f.nextGERDeadline == 0 {
				f.setNextGERDeadline()
			}
			f.nextGERMux.Unlock()
		// L2Reorg ch
		case <-f.closingSignalCh.L2ReorgCh:
			log.Debug("finalizer received L2 reorg event")
			f.handlingL2Reorg = true
			f.halt(ctx, fmt.Errorf("L2 reorg event received"))
			return
		}
	}
}

// updateLastPendingFLushID updates f.lastPendingFLushID with newFlushID value (it it has changed) and sends
// the signal condition f.pendingFlushIDCond to notify other go funcs that the f.lastPendingFlushID value has changed
func (f *finalizer) updateLastPendingFlushID(newFlushID uint64) {
	if newFlushID > f.lastPendingFlushID {
		f.lastPendingFlushID = newFlushID
		f.pendingFlushIDCond.Broadcast()
	}
}

// finalizeBatches runs the endless loop for processing transactions finalizing batches.
func (f *finalizer) finalizeBatches(ctx context.Context) {
	log.Debug("finalizer init loop")
	showNotFoundTxLog := true // used to log debug only the first message when there is no txs to process
	for {
		start := now()
		if f.wipBatch.batchNumber == f.cfg.StopSequencerOnBatchNum {
			f.halt(ctx, fmt.Errorf("finalizer reached stop sequencer batch number: %v", f.cfg.StopSequencerOnBatchNum))
		}

		// We have reached the L2 block time, we need to close the current L2 block and open a new one
		if !f.wipL2Block.timestamp.Add(f.cfg.L2BlockTime.Duration).After(time.Now()) {
			f.finalizeL2Block(ctx)
		}

		tx, err := f.worker.GetBestFittingTx(f.wipBatch.remainingResources)

		// If we have txs pending to process but none of them fits into the wip batch, we close the wip batch and open a new one
		if err == ErrNoFittingTransaction { //TODO: review this with JEC
			f.finalizeBatch(ctx)
		}

		metrics.WorkerProcessingTime(time.Since(start))
		if tx != nil {
			log.Debugf("processing tx: %s", tx.Hash.Hex())
			showNotFoundTxLog = true

			firstTxProcess := true

			for {
				_, err := f.processTransaction(ctx, tx, firstTxProcess)
				if err != nil {
					if err == ErrEffectiveGasPriceReprocess {
						firstTxProcess = false
						log.Info("reprocessing tx because of effective gas price calculation: %s", tx.Hash.Hex())
						continue
					} else {
						log.Errorf("failed to process transaction in finalizeBatches, Err: %v", err)
						break
					}
				}
				break
			}
		} else {
			// wait for new txs
			if showNotFoundTxLog {
				log.Debug("no transactions to be processed. Waiting...")
				showNotFoundTxLog = false
			}
			if f.cfg.SleepDuration.Duration > 0 {
				time.Sleep(f.cfg.SleepDuration.Duration)
			}
		}

		if !f.cfg.SequentialReprocessFullBatch && f.reprocessFullBatchError.Load() {
			// There is an error reprocessing previous batch closed (parallel sanity check)
			// We halt the execution of the Sequencer at this point
			f.halt(ctx, fmt.Errorf("halting Sequencer because of error reprocessing full batch (sanity check). Check previous errors in logs to know which was the cause"))
		}

		if f.isDeadlineEncountered() {
			f.finalizeBatch(ctx)
		} else if f.maxTxsPerBatchReached() || f.isBatchResourcesExhausted() {
			f.finalizeBatch(ctx)
		}

		if err := ctx.Err(); err != nil {
			log.Infof("stopping finalizer because of context, err: %s", err)
			return
		}
	}
}

// sortForcedBatches sorts the forced batches by ForcedBatchNumber
func (f *finalizer) sortForcedBatches(fb []state.ForcedBatch) []state.ForcedBatch {
	if len(fb) == 0 {
		return fb
	}
	// Sort by ForcedBatchNumber
	for i := 0; i < len(fb)-1; i++ {
		for j := i + 1; j < len(fb); j++ {
			if fb[i].ForcedBatchNumber > fb[j].ForcedBatchNumber {
				fb[i], fb[j] = fb[j], fb[i]
			}
		}
	}

	return fb
}

// maxTxsPerBatchReached checks if the batch has reached the maximum number of txs per batch
func (f *finalizer) maxTxsPerBatchReached() bool {
	if f.wipBatch.countOfTxs >= int(f.batchConstraints.MaxTxsPerBatch) {
		log.Infof("closing batch: %d, because it reached the maximum number of txs.", f.wipBatch.batchNumber)
		f.wipBatch.closingReason = state.BatchFullClosingReason
		return true
	}
	return false
}

// checkIfProverRestarted checks if the proverID changed
func (f *finalizer) checkIfProverRestarted(proverID string) {
	if f.proverID != "" && f.proverID != proverID {
		event := &event.Event{
			ReceivedAt:  time.Now(),
			Source:      event.Source_Node,
			Component:   event.Component_Sequencer,
			Level:       event.Level_Critical,
			EventID:     event.EventID_FinalizerRestart,
			Description: fmt.Sprintf("proverID changed from %s to %s, restarting sequencer to discard current WIP batch and work with new executor", f.proverID, proverID),
		}

		err := f.eventLog.LogEvent(context.Background(), event)
		if err != nil {
			log.Errorf("error storing payload: %v", err)
		}

		log.Fatal("restarting sequencer to discard current WIP batch and work with new executor")
	}
}

// processTransaction processes a single transaction.
func (f *finalizer) processTransaction(ctx context.Context, tx *TxTracker, firstTxProcess bool) (errWg *sync.WaitGroup, err error) {
	var txHash string

	if tx != nil {
		txHash = tx.Hash.String()
	}

	log := log.WithFields("txHash", txHash, "batchNumber", f.wipBatch.batchNumber)

	start := time.Now()
	defer func() {
		metrics.ProcessingTime(time.Since(start))
	}()

	executorBatchRequest := state.ProcessRequest{
		BatchNumber:       f.wipBatch.batchNumber,
		OldStateRoot:      f.wipBatch.imStateRoot,
		OldAccInputHash:   f.wipBatch.imAccInputHash,
		Coinbase:          f.wipBatch.coinbase,
		L1InfoRoot_V2:     mockL1InfoRoot,
		TimestampLimit_V2: uint64(f.wipL2Block.timestamp.Unix()),
		Caller:            stateMetrics.SequencerCallerLabel,
		ForkID:            f.state.GetForkIDByBatchNumber(f.wipBatch.batchNumber),
	}

	if f.wipBatch.isEmpty() {
		executorBatchRequest.Transactions = f.state.BuildChangeL2Block(f.wipL2Block.deltaTimestamp, f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)
		executorBatchRequest.SkipFirstChangeL2Block_V2 = false
	} else {
		executorBatchRequest.Transactions = []byte{}
		executorBatchRequest.SkipFirstChangeL2Block_V2 = true
	}
	executorBatchRequest.SkipWriteBlockInfoRoot_V2 = true

	hashStr := "nil"
	if tx != nil {
		executorBatchRequest.Transactions = append(executorBatchRequest.Transactions, tx.RawTx...)
		hashStr = tx.HashStr

		txGasPrice := tx.GasPrice

		// If it is the first time we process this tx then we calculate the EffectiveGasPrice
		if firstTxProcess {
			// Get L1 gas price and store in txTracker to make it consistent during the lifespan of the transaction
			tx.L1GasPrice, tx.L2GasPrice = f.pool.GetL1AndL2GasPrice()
			// Get the tx and l2 gas price we will use in the egp calculation. If egp is disabled we will use a "simulated" tx gas price
			txGasPrice, txL2GasPrice := f.effectiveGasPrice.GetTxAndL2GasPrice(tx.GasPrice, tx.L1GasPrice, tx.L2GasPrice)

			// Save values for later logging
			tx.EGPLog.L1GasPrice = tx.L1GasPrice
			tx.EGPLog.L2GasPrice = txL2GasPrice
			tx.EGPLog.GasUsedFirst = tx.BatchResources.ZKCounters.GasUsed
			tx.EGPLog.GasPrice.Set(txGasPrice)

			// Calculate EffectiveGasPrice
			egp, err := f.effectiveGasPrice.CalculateEffectiveGasPrice(tx.RawTx, txGasPrice, tx.BatchResources.ZKCounters.GasUsed, tx.L1GasPrice, txL2GasPrice)
			if err != nil {
				if f.effectiveGasPrice.IsEnabled() {
					return nil, err
				} else {
					log.Warnf("EffectiveGasPrice is disabled, but failed to calculate EffectiveGasPrice: %s", err)
					tx.EGPLog.Error = fmt.Sprintf("CalculateEffectiveGasPrice#1: %s", err)
				}
			} else {
				tx.EffectiveGasPrice.Set(egp)

				// Save first EffectiveGasPrice for later logging
				tx.EGPLog.ValueFirst.Set(tx.EffectiveGasPrice)

				// If EffectiveGasPrice >= txGasPrice, we process the tx with tx.GasPrice
				if tx.EffectiveGasPrice.Cmp(txGasPrice) >= 0 {
					tx.EffectiveGasPrice.Set(txGasPrice)

					loss := new(big.Int).Sub(tx.EffectiveGasPrice, txGasPrice)
					// If loss > 0 the warning message indicating we loss fee for thix tx
					if loss.Cmp(new(big.Int).SetUint64(0)) == 1 {
						log.Warnf("egp-loss: gasPrice: %d, effectiveGasPrice1: %d, loss: %d, txHash: %s", txGasPrice, tx.EffectiveGasPrice, loss, tx.HashStr)
					}

					tx.IsLastExecution = true
				}
			}
		}

		effectivePercentage, err := f.effectiveGasPrice.CalculateEffectiveGasPricePercentage(txGasPrice, tx.EffectiveGasPrice)
		if err != nil {
			if f.effectiveGasPrice.IsEnabled() {
				return nil, err
			} else {
				log.Warnf("EffectiveGasPrice is disabled, but failed to to CalculateEffectiveGasPricePercentage#1: %s", err)
				tx.EGPLog.Error = fmt.Sprintf("%s; CalculateEffectiveGasPricePercentage#1: %s", tx.EGPLog.Error, err)
			}
		} else {
			// Save percentage for later logging
			tx.EGPLog.Percentage = effectivePercentage
		}

		// If EGP is disabled we use tx GasPrice (MaxEffectivePercentage=255)
		if !f.effectiveGasPrice.IsEnabled() {
			effectivePercentage = state.MaxEffectivePercentage
		}

		effectivePercentageAsDecodedHex, err := hex.DecodeHex(fmt.Sprintf("%x", effectivePercentage))
		if err != nil {
			return nil, err
		}

		executorBatchRequest.Transactions = append(executorBatchRequest.Transactions, effectivePercentageAsDecodedHex...)
	}

	log.Infof("processing tx: %s. Batch.BatchNumber: %d, batchNumber: %d, oldStateRoot: %s, txHash: %s, L1InfoRoot: %s", hashStr, f.wipBatch.batchNumber, executorBatchRequest.BatchNumber, executorBatchRequest.OldStateRoot, hashStr, executorBatchRequest.L1InfoRoot_V2.String())
	processBatchResponse, err := f.state.ProcessBatchV2(ctx, executorBatchRequest, false)
	if err != nil && errors.Is(err, runtime.ErrExecutorDBError) {
		log.Errorf("failed to process transaction: %s", err)
		return nil, err
	} else if err == nil && !processBatchResponse.IsRomLevelError && len(processBatchResponse.BlockResponses) == 0 && tx != nil {
		err = fmt.Errorf("executor returned no errors and no responses for tx: %s", tx.HashStr)
		f.halt(ctx, err)
	} else if processBatchResponse.IsExecutorLevelError && tx != nil {
		log.Errorf("error received from executor. Error: %v", err)
		// Delete tx from the worker
		f.worker.DeleteTx(tx.Hash, tx.From)

		// Set tx as invalid in the pool
		errMsg := processBatchResponse.ExecutorError.Error()
		err = f.pool.UpdateTxStatus(ctx, tx.Hash, poolPackage.TxStatusInvalid, false, &errMsg)
		if err != nil {
			log.Errorf("failed to update status to invalid in the pool for tx: %s, err: %s", tx.Hash.String(), err)
		} else {
			metrics.TxProcessed(metrics.TxProcessedLabelInvalid, 1)
		}
		return nil, err
	}

	oldStateRoot := f.wipBatch.imStateRoot
	if len(processBatchResponse.BlockResponses) > 0 && tx != nil {
		errWg, err = f.handleProcessTransactionResponse(ctx, tx, processBatchResponse, oldStateRoot)
		if err != nil {
			return errWg, err
		}
	}

	// Update wip batch
	f.wipBatch.imStateRoot = processBatchResponse.NewStateRoot
	f.wipBatch.localExitRoot = processBatchResponse.NewLocalExitRoot
	f.wipBatch.imAccInputHash = processBatchResponse.NewAccInputHash

	log.Infof("processed tx: %s. Batch.batchNumber: %d, batchNumber: %d, newStateRoot: %s, newLocalExitRoot: %s, oldStateRoot: %s",
		hashStr, f.wipBatch.batchNumber, executorBatchRequest.BatchNumber, processBatchResponse.NewStateRoot.String(), processBatchResponse.NewLocalExitRoot.String(), oldStateRoot.String())

	return nil, nil
}

// handleProcessTransactionResponse handles the response of transaction processing.
func (f *finalizer) handleProcessTransactionResponse(ctx context.Context, tx *TxTracker, result *state.ProcessBatchResponse, oldStateRoot common.Hash) (errWg *sync.WaitGroup, err error) {
	// Handle Transaction Error
	errorCode := executor.RomErrorCode(result.BlockResponses[0].TransactionResponses[0].RomError)
	if !state.IsStateRootChanged(errorCode) {
		// If intrinsic error or OOC error, we skip adding the transaction to the batch
		errWg = f.handleProcessTransactionError(ctx, result, tx)
		return errWg, result.BlockResponses[0].TransactionResponses[0].RomError
	}

	// Check remaining resources
	err = f.checkRemainingResources(result, tx)
	if err != nil {
		return nil, err
	}

	egpEnabled := f.effectiveGasPrice.IsEnabled()

	if !tx.IsLastExecution {
		tx.IsLastExecution = true

		// Get the tx gas price we will use in the egp calculation. If egp is disabled we will use a "simulated" tx gas price
		txGasPrice, txL2GasPrice := f.effectiveGasPrice.GetTxAndL2GasPrice(tx.GasPrice, tx.L1GasPrice, tx.L2GasPrice)

		newEffectiveGasPrice, err := f.effectiveGasPrice.CalculateEffectiveGasPrice(tx.RawTx, txGasPrice, result.BlockResponses[0].TransactionResponses[0].GasUsed, tx.L1GasPrice, txL2GasPrice)
		if err != nil {
			if egpEnabled {
				log.Errorf("failed to calculate EffectiveGasPrice with new gasUsed for tx %s, error: %s", tx.HashStr, err.Error())
				return nil, err
			} else {
				log.Warnf("EffectiveGasPrice is disabled, but failed to calculate EffectiveGasPrice with new gasUsed for tx %s, error: %s", tx.HashStr, err.Error())
				tx.EGPLog.Error = fmt.Sprintf("%s; CalculateEffectiveGasPrice#2: %s", tx.EGPLog.Error, err)
			}
		} else {
			// Save new (second) gas used and second effective gas price calculation for later logging
			tx.EGPLog.ValueSecond.Set(newEffectiveGasPrice)
			tx.EGPLog.GasUsedSecond = result.BlockResponses[0].TransactionResponses[0].GasUsed

			errCompare := f.compareTxEffectiveGasPrice(ctx, tx, newEffectiveGasPrice, result.BlockResponses[0].TransactionResponses[0].HasGaspriceOpcode, result.BlockResponses[0].TransactionResponses[0].HasBalanceOpcode)

			// If EffectiveGasPrice is disabled we will calculate the percentage and save it for later logging
			if !egpEnabled {
				effectivePercentage, err := f.effectiveGasPrice.CalculateEffectiveGasPricePercentage(txGasPrice, tx.EffectiveGasPrice)
				if err != nil {
					log.Warnf("EffectiveGasPrice is disabled, but failed to CalculateEffectiveGasPricePercentage#2: %s", err)
					tx.EGPLog.Error = fmt.Sprintf("%s, CalculateEffectiveGasPricePercentage#2: %s", tx.EGPLog.Error, err)
				} else {
					// Save percentage for later logging
					tx.EGPLog.Percentage = effectivePercentage
				}
			}

			if errCompare != nil && egpEnabled {
				return nil, errCompare
			}
		}
	}

	// Save Enabled, GasPriceOC, BalanceOC and final effective gas price for later logging
	tx.EGPLog.Enabled = egpEnabled
	tx.EGPLog.GasPriceOC = result.BlockResponses[0].TransactionResponses[0].HasGaspriceOpcode
	tx.EGPLog.BalanceOC = result.BlockResponses[0].TransactionResponses[0].HasBalanceOpcode
	tx.EGPLog.ValueFinal.Set(tx.EffectiveGasPrice)

	// Log here the results of EGP calculation
	log.Infof("egp-log: final: %d, first: %d, second: %d, percentage: %d, deviation: %d, maxDeviation: %d, gasUsed1: %d, gasUsed2: %d, gasPrice: %d, l1GasPrice: %d, l2GasPrice: %d, reprocess: %t, gasPriceOC: %t, balanceOC: %t, enabled: %t, txSize: %d, txHash: %s, error: %s",
		tx.EGPLog.ValueFinal, tx.EGPLog.ValueFirst, tx.EGPLog.ValueSecond, tx.EGPLog.Percentage, tx.EGPLog.FinalDeviation, tx.EGPLog.MaxDeviation, tx.EGPLog.GasUsedFirst, tx.EGPLog.GasUsedSecond,
		tx.EGPLog.GasPrice, tx.EGPLog.L1GasPrice, tx.EGPLog.L2GasPrice, tx.EGPLog.Reprocess, tx.EGPLog.GasPriceOC, tx.EGPLog.BalanceOC, egpEnabled, len(tx.RawTx), tx.HashStr, tx.EGPLog.Error)

	tx.FlushId = result.FlushID
	f.wipL2Block.addTx(tx)

	f.updateLastPendingFlushID(result.FlushID)

	f.wipBatch.countOfTxs++

	f.updateWorkerAfterSuccessfulProcessing(ctx, tx.Hash, tx.From, false, result)

	return nil, nil
}

// handleProcessForcedTxsResponse handles the transactions responses for the processed forced batch.
func (f *finalizer) handleProcessForcedTxsResponse(ctx context.Context, request state.ProcessRequest, result *state.ProcessBatchResponse, oldStateRoot common.Hash) {
	/*log.Infof("handleForcedTxsProcessResp: batchNumber: %d, oldStateRoot: %s, newStateRoot: %s", request.BatchNumber, oldStateRoot.String(), result.NewStateRoot.String())
	parentBlockHash := f.wipL2Block.parentHash
	for _, blockResp := range result.BlockResponses {
		if blockResp.BlockNumber != f.wipL2Block.number {
			log.Fatalf("L2 block number mismatch when processing forced batch block response. blockResp.BlockNumber: %d, f.wipL2Block,number: %d", blockResp.BlockNumber, f.wipL2Block.number)
			return
		}

		l2BlockToStore := l2BlockToStore{
			l2Block: &L2Block{
				number:       blockResp.BlockNumber,
				hash:         blockResp.BlockHash,
				parentHash:   parentBlockHash,
				timestamp:    time.Unix(int64(blockResp.Timestamp), 0),
				transactions: []transactionToStore{},
			},
			batchNumber: request.BatchNumber,
			forcedBatch: true,
			coinbase:    request.Coinbase,
			stateRoot:   oldStateRoot,
			flushId:     result.FlushID,
		}
		for _, txResp := range blockResp.TransactionResponses {
			// Handle Transaction Error
			if txResp.RomError != nil {
				romErr := executor.RomErrorCode(txResp.RomError)
				if executor.IsIntrinsicError(romErr) || romErr == executor.RomError_ROM_ERROR_INVALID_RLP {
					// If we have an intrinsic error or the RLP is invalid
					// we should continue processing the batch, but skip the transaction
					log.Errorf("handleForcedTxsProcessResp: ROM error: %s", txResp.RomError)
					continue
				}
			}

			from, err := state.GetSender(txResp.Tx)
			if err != nil {
				log.Warnf("handleForcedTxsProcessResp: failed to get sender for tx (%s): %v", txResp.TxHash, err)
			}

			//TODO: How to manage L2 block for forced batch/txs
			txToStore := transactionToStore{
				hash:          txResp.TxHash,
				from:          from,
				response:      txResp,
				batchResponse: result,
				batchNumber:   request.BatchNumber,
				timestamp:     request.Timestamp_V1,
				coinbase:      request.Coinbase,
				oldStateRoot:  oldStateRoot,
				isForcedBatch: true,
				flushId:       result.FlushID,
			}

			oldStateRoot = txResp.StateRoot

			f.updateLastPendingFlushID(result.FlushID)

			f.addPendingTxToStore(ctx, txToStore)

			if err == nil {
				f.updateWorkerAfterSuccessfulProcessing(ctx, txResp.TxHash, from, true, result)
			}
		}
	}*/
}

// compareTxEffectiveGasPrice compares newEffectiveGasPrice with tx.EffectiveGasPrice.
// It returns ErrEffectiveGasPriceReprocess if the tx needs to be reprocessed with
// the tx.EffectiveGasPrice updated, otherwise it returns nil
func (f *finalizer) compareTxEffectiveGasPrice(ctx context.Context, tx *TxTracker, newEffectiveGasPrice *big.Int, hasGasPriceOC bool, hasBalanceOC bool) error {
	// Get the tx gas price we will use in the egp calculation. If egp is disabled we will use a "simulated" tx gas price
	txGasPrice, _ := f.effectiveGasPrice.GetTxAndL2GasPrice(tx.GasPrice, tx.L1GasPrice, tx.L2GasPrice)

	// Compute the absolute difference between tx.EffectiveGasPrice - newEffectiveGasPrice
	diff := new(big.Int).Abs(new(big.Int).Sub(tx.EffectiveGasPrice, newEffectiveGasPrice))
	// Compute max deviation allowed of newEffectiveGasPrice
	maxDeviation := new(big.Int).Div(new(big.Int).Mul(tx.EffectiveGasPrice, new(big.Int).SetUint64(f.effectiveGasPrice.GetFinalDeviation())), big.NewInt(100)) //nolint:gomnd

	// Save FinalDeviation (diff) and MaxDeviation for later logging
	tx.EGPLog.FinalDeviation.Set(diff)
	tx.EGPLog.MaxDeviation.Set(maxDeviation)

	// if (diff > finalDeviation)
	if diff.Cmp(maxDeviation) == 1 {
		// if newEfectiveGasPrice < txGasPrice
		if newEffectiveGasPrice.Cmp(txGasPrice) == -1 {
			if hasGasPriceOC || hasBalanceOC {
				tx.EffectiveGasPrice.Set(txGasPrice)
			} else {
				tx.EffectiveGasPrice.Set(newEffectiveGasPrice)
			}
		} else {
			tx.EffectiveGasPrice.Set(txGasPrice)

			loss := new(big.Int).Sub(newEffectiveGasPrice, txGasPrice)
			// If loss > 0 the warning message indicating we loss fee for thix tx
			if loss.Cmp(new(big.Int).SetUint64(0)) == 1 {
				log.Warnf("egp-loss: gasPrice: %d, EffectiveGasPrice2: %d, loss: %d, txHash: %s", txGasPrice, newEffectiveGasPrice, loss, tx.HashStr)
			}
		}

		// Save Reprocess for later logging
		tx.EGPLog.Reprocess = true

		return ErrEffectiveGasPriceReprocess
	} // else (diff <= finalDeviation) it is ok, no reprocess of the tx is needed

	return nil
}

func (f *finalizer) updateWorkerAfterSuccessfulProcessing(ctx context.Context, txHash common.Hash, txFrom common.Address, isForced bool, result *state.ProcessBatchResponse) {
	// Delete the transaction from the worker
	if isForced {
		f.worker.DeleteForcedTx(txHash, txFrom)
		log.Debugf("forced tx deleted from worker. txHash: %s, from: %s", txHash.String(), txFrom.Hex())
		return
	} else {
		f.worker.DeleteTx(txHash, txFrom)
		log.Debugf("tx deleted from worker. txHash: %s, from: %s", txHash.String(), txFrom.Hex())
	}

	start := time.Now()
	txsToDelete := f.worker.UpdateAfterSingleSuccessfulTxExecution(txFrom, result.ReadWriteAddresses)
	for _, txToDelete := range txsToDelete {
		err := f.pool.UpdateTxStatus(ctx, txToDelete.Hash, poolPackage.TxStatusFailed, false, txToDelete.FailedReason)
		if err != nil {
			log.Errorf("failed to update status to failed in the pool for tx: %s, err: %s", txToDelete.Hash.String(), err)
			continue
		}
		metrics.TxProcessed(metrics.TxProcessedLabelFailed, 1)
	}
	metrics.WorkerProcessingTime(time.Since(start))
}

// handleProcessTransactionError handles the error of a transaction
func (f *finalizer) handleProcessTransactionError(ctx context.Context, result *state.ProcessBatchResponse, tx *TxTracker) *sync.WaitGroup {
	txResponse := result.BlockResponses[0].TransactionResponses[0]
	errorCode := executor.RomErrorCode(txResponse.RomError)
	addressInfo := result.ReadWriteAddresses[tx.From]
	log.Infof("handleTransactionError: error in tx: %s, errorCode: %d", tx.Hash.String(), errorCode)
	wg := new(sync.WaitGroup)
	failedReason := executor.RomErr(errorCode).Error()
	if executor.IsROMOutOfCountersError(errorCode) {
		log.Errorf("ROM out of counters error, marking tx with Hash: %s as INVALID, errorCode: %s", tx.Hash.String(), errorCode.String())
		start := time.Now()
		f.worker.DeleteTx(tx.Hash, tx.From)
		metrics.WorkerProcessingTime(time.Since(start))

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := f.pool.UpdateTxStatus(ctx, tx.Hash, poolPackage.TxStatusInvalid, false, &failedReason)
			if err != nil {
				log.Errorf("failed to update status to failed in the pool for tx: %s, err: %s", tx.Hash.String(), err)
			} else {
				metrics.TxProcessed(metrics.TxProcessedLabelInvalid, 1)
			}
		}()
	} else if executor.IsInvalidNonceError(errorCode) || executor.IsInvalidBalanceError(errorCode) {
		var (
			nonce   *uint64
			balance *big.Int
		)
		if addressInfo != nil {
			nonce = addressInfo.Nonce
			balance = addressInfo.Balance
		}
		start := time.Now()
		log.Errorf("intrinsic error, moving tx with Hash: %s to NOT READY nonce(%d) balance(%d) gasPrice(%d), err: %s", tx.Hash, nonce, balance, tx.GasPrice, txResponse.RomError)
		txsToDelete := f.worker.MoveTxToNotReady(tx.Hash, tx.From, nonce, balance)
		for _, txToDelete := range txsToDelete {
			wg.Add(1)
			txToDelete := txToDelete
			go func() {
				defer wg.Done()
				err := f.pool.UpdateTxStatus(ctx, txToDelete.Hash, poolPackage.TxStatusFailed, false, &failedReason)
				metrics.TxProcessed(metrics.TxProcessedLabelFailed, 1)
				if err != nil {
					log.Errorf("failed to update status to failed in the pool for tx: %s, err: %s", txToDelete.Hash.String(), err)
				}
			}()
		}
		metrics.WorkerProcessingTime(time.Since(start))
	} else {
		// Delete the transaction from the txSorted list
		f.worker.DeleteTx(tx.Hash, tx.From)
		log.Debug("tx deleted from txSorted list", "txHash", tx.Hash.String(), "from", tx.From.Hex())

		wg.Add(1)
		go func() {
			defer wg.Done()
			// Update the status of the transaction to failed
			err := f.pool.UpdateTxStatus(ctx, tx.Hash, poolPackage.TxStatusFailed, false, &failedReason)
			if err != nil {
				log.Errorf("failed to update status to failed in the pool for tx: %s, err: %s", tx.Hash.String(), err)
			} else {
				metrics.TxProcessed(metrics.TxProcessedLabelFailed, 1)
			}
		}()
	}

	return wg
}

// processForcedBatches processes all the forced batches that are pending to be processed
func (f *finalizer) processForcedBatches(ctx context.Context, lastBatchNumberInState uint64, stateRoot common.Hash) (uint64, common.Hash, error) {
	f.nextForcedBatchesMux.Lock()
	defer f.nextForcedBatchesMux.Unlock()
	f.nextForcedBatchDeadline = 0

	lastTrustedForcedBatchNumber, err := f.state.GetLastTrustedForcedBatchNumber(ctx, nil)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("failed to get last trusted forced batch number, err: %w", err)
	}
	nextForcedBatchNum := lastTrustedForcedBatchNumber + 1

	for _, forcedBatch := range f.nextForcedBatches {
		// Skip already processed forced batches
		if forcedBatch.ForcedBatchNumber < nextForcedBatchNum {
			continue
		}
		// Process in-between unprocessed forced batches
		for forcedBatch.ForcedBatchNumber > nextForcedBatchNum {
			inBetweenForcedBatch, err := f.state.GetForcedBatch(ctx, nextForcedBatchNum, nil)
			if err != nil {
				return 0, common.Hash{}, fmt.Errorf("failed to get in-between forced batch %d, err: %w", nextForcedBatchNum, err)
			}
			lastBatchNumberInState, stateRoot = f.processForcedBatch(ctx, lastBatchNumberInState, stateRoot, *inBetweenForcedBatch)
			nextForcedBatchNum += 1
		}
		// Process the current forced batch from the channel queue
		lastBatchNumberInState, stateRoot = f.processForcedBatch(ctx, lastBatchNumberInState, stateRoot, forcedBatch)
		nextForcedBatchNum += 1
	}
	f.nextForcedBatches = make([]state.ForcedBatch, 0)

	return lastBatchNumberInState, stateRoot, nil
}

// ProcessForcedBatch2 process a forced batch
func (f *finalizer) processForcedBatch2(ctx context.Context, forcedBatchNumber uint64, request state.ProcessRequest) (*state.ProcessBatchResponse, error) {
	// Open Batch
	processingCtx := state.ProcessingContext{
		BatchNumber:    request.BatchNumber,
		Coinbase:       request.Coinbase,
		Timestamp:      request.Timestamp_V1,
		GlobalExitRoot: request.GlobalExitRoot_V1,
		ForcedBatchNum: &forcedBatchNumber,
	}
	dbTx, err := f.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for opening a forced batch, err: %v", err)
		return nil, err
	}

	err = f.state.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when opening a forced batch that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		log.Errorf("failed to open a batch, err: %v", err)
		return nil, err
	}

	// Fetch Forced Batch
	forcedBatch, err := f.state.GetForcedBatch(ctx, forcedBatchNumber, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when getting forced batch err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		log.Errorf("failed to get a forced batch, err: %v", err)
		return nil, err
	}

	// Process Batch
	processBatchResponse, err := f.state.ProcessSequencerBatch(ctx, request.BatchNumber, forcedBatch.RawTxsData, request.Caller, dbTx)
	if err != nil {
		log.Errorf("failed to process a forced batch, err: %v", err)
		return nil, err
	}

	// Close Batch
	txsBytes := uint64(0)
	for _, blockResp := range processBatchResponse.BlockResponses {
		for _, resp := range blockResp.TransactionResponses {
			if !resp.ChangesStateRoot {
				continue
			}
			txsBytes += resp.Tx.Size()
		}
	}
	processingReceipt := state.ProcessingReceipt{
		BatchNumber:   request.BatchNumber,
		StateRoot:     processBatchResponse.NewStateRoot,
		LocalExitRoot: processBatchResponse.NewLocalExitRoot,
		AccInputHash:  processBatchResponse.NewAccInputHash,
		BatchL2Data:   forcedBatch.RawTxsData,
		BatchResources: state.BatchResources{
			ZKCounters: processBatchResponse.UsedZkCounters,
			Bytes:      txsBytes,
		},
		ClosingReason: state.ForcedBatchClosingReason,
	}

	isClosed := false
	tryToCloseAndCommit := true
	for tryToCloseAndCommit {
		if !isClosed {
			closingErr := f.state.CloseBatch(ctx, processingReceipt, dbTx)
			tryToCloseAndCommit = closingErr != nil
			if tryToCloseAndCommit {
				continue
			}
			isClosed = true
		}

		if err := dbTx.Commit(ctx); err != nil {
			log.Errorf("failed to commit dbTx when processing a forced batch, err: %v", err)
		}
		tryToCloseAndCommit = err != nil
	}

	return processBatchResponse, nil
}

func (f *finalizer) processForcedBatch(ctx context.Context, lastBatchNumberInState uint64, stateRoot common.Hash, forcedBatch state.ForcedBatch) (uint64, common.Hash) {
	//TODO: review this request for forced txs
	executorBatchRequest := state.ProcessRequest{
		BatchNumber:               lastBatchNumberInState + 1,
		OldStateRoot:              stateRoot,
		L1InfoRoot_V2:             forcedBatch.GlobalExitRoot,
		Transactions:              forcedBatch.RawTxsData,
		Coinbase:                  f.sequencerAddress,
		TimestampLimit_V2:         uint64(forcedBatch.ForcedAt.Unix()), //TODO: review this is the TimeStampLimit we need to use
		SkipFirstChangeL2Block_V2: false,
		SkipWriteBlockInfoRoot_V2: false,
		Caller:                    stateMetrics.SequencerCallerLabel,
	}

	response, err := f.processForcedBatch2(ctx, forcedBatch.ForcedBatchNumber, executorBatchRequest)
	if err != nil {
		// If there is EXECUTOR (Batch level) error, halt the finalizer.
		f.halt(ctx, fmt.Errorf("failed to process forced batch, Executor err: %w", err))
		return lastBatchNumberInState, stateRoot
	}

	if len(response.BlockResponses) > 0 && !response.IsRomOOCError {
		for _, blockResponse := range response.BlockResponses {
			for _, txResponse := range blockResponse.TransactionResponses {
				if !errors.Is(txResponse.RomError, executor.RomErr(executor.RomError_ROM_ERROR_INVALID_RLP)) {
					sender, err := state.GetSender(txResponse.Tx)
					if err != nil {
						log.Warnf("failed trying to add forced tx (%s) to worker. Error getting sender from tx, Err: %v", txResponse.TxHash, err)
						continue
					}
					f.worker.AddForcedTx(txResponse.TxHash, sender)
				} else {
					log.Warnf("ROM_ERROR_INVALID_RLP error received from executor for forced batch %d", forcedBatch.ForcedBatchNumber)
				}
			}
		}

		f.handleProcessForcedTxsResponse(ctx, executorBatchRequest, response, stateRoot)
	} else {
		if f.streamServer != nil && f.currentGERHash != forcedBatch.GlobalExitRoot {
			//TODO: review this datastream event
			updateGer := state.DSUpdateGER{
				BatchNumber:    executorBatchRequest.BatchNumber,
				Timestamp:      executorBatchRequest.Timestamp_V1.Unix(),
				GlobalExitRoot: executorBatchRequest.GlobalExitRoot_V1,
				Coinbase:       f.sequencerAddress,
				ForkID:         uint16(f.state.GetForkIDByBatchNumber(executorBatchRequest.BatchNumber)),
				StateRoot:      response.NewStateRoot,
			}

			err = f.streamServer.StartAtomicOp()
			if err != nil {
				log.Errorf("failed to start atomic op for forced batch %v: %v", forcedBatch.ForcedBatchNumber, err)
			}

			_, err = f.streamServer.AddStreamEntry(state.EntryTypeUpdateGER, updateGer.Encode())
			if err != nil {
				log.Errorf("failed to add stream entry for forced batch %v: %v", forcedBatch.ForcedBatchNumber, err)
			}

			err = f.streamServer.CommitAtomicOp()
			if err != nil {
				log.Errorf("failed to commit atomic op for forced batch %v: %v", forcedBatch.ForcedBatchNumber, err)
			}
		}
	}

	stateRoot = response.NewStateRoot
	lastBatchNumberInState++

	return lastBatchNumberInState, stateRoot
}

// reprocessFullBatch reprocesses a batch used as sanity check
func (f *finalizer) reprocessFullBatch(ctx context.Context, batchNum uint64, initialStateRoot common.Hash, initialAccInputHash common.Hash, expectedNewStateRoot common.Hash) (*state.ProcessBatchResponse, error) {
	reprocessError := func(batch *state.Batch) {
		f.reprocessFullBatchError.Store(true)

		rawL2Blocks, err := state.DecodeBatchV2(batch.BatchL2Data)
		if err != nil {
			log.Errorf("[reprocessFullBatch] error decoding BatchL2Data for batch %d. Error: %s", batch.BatchNumber, err)
			return
		}

		// Log batch detailed info
		log.Infof("[reprocessFullBatch] BatchNumber: %d, InitialStateRoot: %s, ExpectedNewStateRoot: %s, GER: %s", batch.BatchNumber, initialStateRoot, expectedNewStateRoot, batch.GlobalExitRoot)
		for i, rawL2block := range rawL2Blocks.Blocks {
			for j, rawTx := range rawL2block.Transactions {
				log.Infof("[reprocessFullBatch] BatchNumber: %d, block position: % d, tx position %d, tx hash: %s", batch.BatchNumber, i, j, rawTx.Tx.Hash())
			}
		}
	}

	log.Debugf("[reprocessFullBatch] reprocessing batch: %d, InitialStateRoot: %s, ExpectedNewStateRoot: %s, GER: %s", batchNum, initialStateRoot, expectedNewStateRoot)

	batch, err := f.state.GetBatchByNumber(ctx, batchNum, nil)
	if err != nil {
		log.Errorf("[reprocessFullBatch] failed to get batch %d, err: %s", batchNum, err)
		f.reprocessFullBatchError.Store(true)
		return nil, ErrGetBatchByNumber
	}

	caller := stateMetrics.DiscardCallerLabel
	if f.cfg.SequentialReprocessFullBatch {
		caller = stateMetrics.SequencerCallerLabel
	}

	// TODO: review this request for reprocess full batch
	executorBatchRequest := state.ProcessRequest{
		BatchNumber:       batch.BatchNumber,
		L1InfoRoot_V2:     mockL1InfoRoot,
		OldStateRoot:      initialStateRoot,
		OldAccInputHash:   initialAccInputHash,
		Transactions:      batch.BatchL2Data,
		Coinbase:          batch.Coinbase,
		TimestampLimit_V2: uint64(time.Now().Unix()),
		Caller:            caller,
	}

	var result *state.ProcessBatchResponse

	result, err = f.state.ProcessBatchV2(ctx, executorBatchRequest, false)
	if err != nil {
		log.Errorf("[reprocessFullBatch] failed to process batch %d. Error: %s", batch.BatchNumber, err)
		reprocessError(batch)
		return nil, ErrProcessBatch
	}

	if result.ExecutorError != nil {
		log.Errorf("[reprocessFullBatch] executor error when reprocessing batch %d, error: %s", batch.BatchNumber, result.ExecutorError)
		reprocessError(batch)
		return nil, ErrExecutorError
	}

	if result.IsRomOOCError {
		log.Errorf("[reprocessFullBatch] failed to process batch %d because OutOfCounters", batch.BatchNumber)
		reprocessError(batch)

		payload, err := json.Marshal(executorBatchRequest)
		if err != nil {
			log.Errorf("[reprocessFullBatch] error marshaling payload: %s", err)
		} else {
			event := &event.Event{
				ReceivedAt:  time.Now(),
				Source:      event.Source_Node,
				Component:   event.Component_Sequencer,
				Level:       event.Level_Critical,
				EventID:     event.EventID_ReprocessFullBatchOOC,
				Description: string(payload),
				Json:        executorBatchRequest,
			}
			err = f.eventLog.LogEvent(ctx, event)
			if err != nil {
				log.Errorf("[reprocessFullBatch] error storing payload: %s", err)
			}
		}

		return nil, ErrProcessBatchOOC
	}

	if result.NewStateRoot != expectedNewStateRoot {
		log.Errorf("[reprocessFullBatch] new state root mismatch for batch %d, expected: %s, got: %s", batch.BatchNumber, expectedNewStateRoot.String(), result.NewStateRoot.String())
		reprocessError(batch)
		return nil, ErrStateRootNoMatch
	}

	log.Infof("[reprocessFullBatch]: reprocess successfully done for batch %d", batch.BatchNumber)
	return result, nil
}

// isDeadlineEncountered returns true if any closing signal deadline is encountered
func (f *finalizer) isDeadlineEncountered() bool {
	// Forced batch deadline
	if f.nextForcedBatchDeadline != 0 && now().Unix() >= f.nextForcedBatchDeadline {
		log.Infof("closing batch %d, forced batch deadline encountered.", f.wipBatch.batchNumber)
		return true
	}
	// Global Exit Root deadline
	if f.nextGERDeadline != 0 && now().Unix() >= f.nextGERDeadline {
		log.Infof("closing batch %d, GER deadline encountered.", f.wipBatch.batchNumber)
		f.wipBatch.closingReason = state.GlobalExitRootDeadlineClosingReason
		return true
	}
	// Timestamp resolution deadline
	if !f.wipBatch.isEmpty() && f.wipBatch.timestamp.Add(f.cfg.TimestampResolution.Duration).Before(time.Now()) {
		log.Infof("closing batch %d, because of timestamp resolution.", f.wipBatch.batchNumber)
		f.wipBatch.closingReason = state.TimeoutResolutionDeadlineClosingReason
		return true
	}
	return false
}

// setNextForcedBatchDeadline sets the next forced batch deadline
func (f *finalizer) setNextForcedBatchDeadline() {
	f.nextForcedBatchDeadline = now().Unix() + int64(f.cfg.ForcedBatchDeadlineTimeout.Duration.Seconds())
}

// setNextGERDeadline sets the next Global Exit Root deadline
func (f *finalizer) setNextGERDeadline() {
	f.nextGERDeadline = now().Unix() + int64(f.cfg.GERDeadlineTimeout.Duration.Seconds())
}

// checkRemainingResources checks if the transaction uses less resources than the remaining ones in the batch.
func (f *finalizer) checkRemainingResources(result *state.ProcessBatchResponse, tx *TxTracker) error {
	usedResources := state.BatchResources{
		ZKCounters: result.UsedZkCounters,
		Bytes:      uint64(len(tx.RawTx)),
	}

	err := f.wipBatch.remainingResources.Sub(usedResources)
	if err != nil {
		log.Infof("current transaction exceeds the remaining batch resources, updating metadata for tx in worker and continuing")
		start := time.Now()
		f.worker.UpdateTxZKCounters(result.BlockResponses[0].TransactionResponses[0].TxHash, tx.From, usedResources.ZKCounters)
		metrics.WorkerProcessingTime(time.Since(start))
		return err
	}

	return nil
}

// isBatchResourcesExhausted checks if one of resources of the wip batch has reached the max value
func (f *finalizer) isBatchResourcesExhausted() bool {
	resources := f.wipBatch.remainingResources
	zkCounters := resources.ZKCounters
	result := false
	resourceDesc := ""
	if resources.Bytes <= f.getConstraintThresholdUint64(f.batchConstraints.MaxBatchBytesSize) {
		resourceDesc = "MaxBatchBytesSize"
		result = true
	} else if zkCounters.UsedSteps <= f.getConstraintThresholdUint32(f.batchConstraints.MaxSteps) {
		resourceDesc = "MaxSteps"
		result = true
	} else if zkCounters.UsedPoseidonPaddings <= f.getConstraintThresholdUint32(f.batchConstraints.MaxPoseidonPaddings) {
		resourceDesc = "MaxPoseidonPaddings"
		result = true
	} else if zkCounters.UsedBinaries <= f.getConstraintThresholdUint32(f.batchConstraints.MaxBinaries) {
		resourceDesc = "MaxBinaries"
		result = true
	} else if zkCounters.UsedKeccakHashes <= f.getConstraintThresholdUint32(f.batchConstraints.MaxKeccakHashes) {
		resourceDesc = "MaxKeccakHashes"
		result = true
	} else if zkCounters.UsedArithmetics <= f.getConstraintThresholdUint32(f.batchConstraints.MaxArithmetics) {
		resourceDesc = "MaxArithmetics"
		result = true
	} else if zkCounters.UsedMemAligns <= f.getConstraintThresholdUint32(f.batchConstraints.MaxMemAligns) {
		resourceDesc = "MaxMemAligns"
		result = true
	} else if zkCounters.GasUsed <= f.getConstraintThresholdUint64(f.batchConstraints.MaxCumulativeGasUsed) {
		resourceDesc = "MaxCumulativeGasUsed"
		result = true
	} else if zkCounters.UsedSha256Hashes_V2 <= f.getConstraintThresholdUint32(f.batchConstraints.MaxSHA256Hashes) {
		resourceDesc = "MaxSHA256Hashes"
		result = true
	}

	if result {
		log.Infof("closing batch %d, because it reached %s limit", f.wipBatch.batchNumber, resourceDesc)
		f.wipBatch.closingReason = state.BatchAlmostFullClosingReason
	}

	return result
}

// getConstraintThresholdUint64 returns the threshold for the given input
func (f *finalizer) getConstraintThresholdUint64(input uint64) uint64 {
	return input * uint64(f.cfg.ResourcePercentageToCloseBatch) / 100 //nolint:gomnd
}

// getConstraintThresholdUint32 returns the threshold for the given input
func (f *finalizer) getConstraintThresholdUint32(input uint32) uint32 {
	return uint32(input*f.cfg.ResourcePercentageToCloseBatch) / 100 //nolint:gomnd
}

// halt halts the finalizer
func (f *finalizer) halt(ctx context.Context, err error) {
	event := &event.Event{
		ReceivedAt:  time.Now(),
		Source:      event.Source_Node,
		Component:   event.Component_Sequencer,
		Level:       event.Level_Critical,
		EventID:     event.EventID_FinalizerHalt,
		Description: fmt.Sprintf("finalizer halted due to error: %s", err),
	}

	eventErr := f.eventLog.LogEvent(ctx, event)
	if eventErr != nil {
		log.Errorf("error storing finalizer halt event: %v", eventErr)
	}

	for {
		log.Errorf("fatal error: %s", err)
		log.Error("halting the finalizer")
		time.Sleep(5 * time.Second) //nolint:gomnd
	}
}

// getUsedBatchResources returns the max resources that can be used in a batch
func getUsedBatchResources(constraints state.BatchConstraintsCfg, remainingResources state.BatchResources) state.BatchResources {
	return state.BatchResources{
		ZKCounters: state.ZKCounters{
			GasUsed:              constraints.MaxCumulativeGasUsed - remainingResources.ZKCounters.GasUsed,
			UsedKeccakHashes:     constraints.MaxKeccakHashes - remainingResources.ZKCounters.UsedKeccakHashes,
			UsedPoseidonHashes:   constraints.MaxPoseidonHashes - remainingResources.ZKCounters.UsedPoseidonHashes,
			UsedPoseidonPaddings: constraints.MaxPoseidonPaddings - remainingResources.ZKCounters.UsedPoseidonPaddings,
			UsedMemAligns:        constraints.MaxMemAligns - remainingResources.ZKCounters.UsedMemAligns,
			UsedArithmetics:      constraints.MaxArithmetics - remainingResources.ZKCounters.UsedArithmetics,
			UsedBinaries:         constraints.MaxBinaries - remainingResources.ZKCounters.UsedBinaries,
			UsedSteps:            constraints.MaxSteps - remainingResources.ZKCounters.UsedSteps,
			UsedSha256Hashes_V2:  constraints.MaxSHA256Hashes - remainingResources.ZKCounters.UsedSha256Hashes_V2,
		},
		Bytes: constraints.MaxBatchBytesSize - remainingResources.Bytes,
	}
}

// getMaxRemainingResources returns the max zkcounters that can be used in a batch
func getMaxRemainingResources(constraints state.BatchConstraintsCfg) state.BatchResources {
	return state.BatchResources{
		ZKCounters: state.ZKCounters{
			GasUsed:              constraints.MaxCumulativeGasUsed,
			UsedKeccakHashes:     constraints.MaxKeccakHashes,
			UsedPoseidonHashes:   constraints.MaxPoseidonHashes,
			UsedPoseidonPaddings: constraints.MaxPoseidonPaddings,
			UsedMemAligns:        constraints.MaxMemAligns,
			UsedArithmetics:      constraints.MaxArithmetics,
			UsedBinaries:         constraints.MaxBinaries,
			UsedSteps:            constraints.MaxSteps,
			UsedSha256Hashes_V2:  constraints.MaxSHA256Hashes,
		},
		Bytes: constraints.MaxBatchBytesSize,
	}
}
