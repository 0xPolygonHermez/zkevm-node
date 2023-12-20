package sequencer

import (
	"context"
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

	//TODO: Review with Carlos which zkCounters are used when creating a new l2 block in the wip batch
	l2BlockUsedResources = statePackage.BatchResources{
		ZKCounters: statePackage.ZKCounters{
			UsedPoseidonHashes: 256, // nolint:gomnd //TODO: config param
			UsedArithmetics:    1,   // nolint:gomnd //TODO: config param
		},
		Bytes: changeL2BlockSize,
	}
)

// finalizer represents the finalizer component of the sequencer.
type finalizer struct {
	cfg              FinalizerCfg
	isSynced         func(ctx context.Context) bool
	sequencerAddress common.Address
	worker           workerInterface
	pool             txPool
	state            stateInterface
	etherman         etherman
	wipBatch         *Batch
	wipL2Block       *L2Block
	batchConstraints statePackage.BatchConstraintsCfg
	haltFinalizer    atomic.Bool
	haltError        error
	// closing signals
	closingSignalCh ClosingSignalCh
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
	dataToStream chan statePackage.DSL2FullBlock
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
		pendingL2BlocksToProcess:   make(chan *L2Block, pendingL2BlocksBufferSize),
		pendingL2BlocksToProcessWG: new(sync.WaitGroup),
		// pending L2 blocks to store in the state
		pendingL2BlocksToStore:   make(chan *L2Block, pendingL2BlocksBufferSize),
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

	f.haltFinalizer.Store(false)

	return &f
}

// Start starts the finalizer.
func (f *finalizer) Start(ctx context.Context) {
	// Init mockL1InfoRoot to a mock value since it must be different to {0,0,...,0}
	for i := 0; i < len(mockL1InfoRoot); i++ {
		mockL1InfoRoot[i] = byte(i)
	}

	// Update L1InfoRoot
	go f.checkL1InfoTreeUpdate(ctx)

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

func (f *finalizer) checkL1InfoTreeUpdate(ctx context.Context) {
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

		if f.haltFinalizer.Load() {
			// There is a fatal error and we need to halt the finalizer and stop processing new txs
			for {
				log.Errorf("halting the finalizer, fatal error: %s", f.haltError)
				time.Sleep(5 * time.Second) //nolint:gomnd
			}
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
		BatchNumber:               f.wipBatch.batchNumber,
		OldStateRoot:              f.wipBatch.imStateRoot,
		OldAccInputHash:           f.wipBatch.imAccInputHash,
		Coinbase:                  f.wipBatch.coinbase,
		L1InfoRoot_V2:             mockL1InfoRoot,
		TimestampLimit_V2:         uint64(f.wipL2Block.timestamp.Unix()),
		Caller:                    stateMetrics.SequencerCallerLabel,
		ForkID:                    f.state.GetForkIDByBatchNumber(f.wipBatch.batchNumber),
		SkipWriteBlockInfoRoot_V2: true,
		SkipVerifyL1InfoRoot_V2:   true,
	}

	executorBatchRequest.L1InfoTreeData_V2[f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex] = statePackage.L1DataV2{
		GlobalExitRoot: f.wipL2Block.l1InfoTreeExitRoot.GlobalExitRoot.GlobalExitRoot,
		BlockHashL1:    f.wipL2Block.l1InfoTreeExitRoot.PreviousBlockHash,
		MinTimestamp:   uint64(f.wipL2Block.l1InfoTreeExitRoot.GlobalExitRoot.Timestamp.Unix()),
	}

	if f.wipL2Block.isEmpty() {
		executorBatchRequest.Transactions = f.state.BuildChangeL2Block(f.wipL2Block.deltaTimestamp, f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)
		executorBatchRequest.SkipFirstChangeL2Block_V2 = false
	} else {
		executorBatchRequest.Transactions = []byte{}
		executorBatchRequest.SkipFirstChangeL2Block_V2 = true
	}

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

		egpPercentage, err := f.effectiveGasPrice.CalculateEffectiveGasPricePercentage(txGasPrice, tx.EffectiveGasPrice)
		if err != nil {
			if f.effectiveGasPrice.IsEnabled() {
				return nil, err
			} else {
				log.Warnf("EffectiveGasPrice is disabled, but failed to to CalculateEffectiveGasPricePercentage#1: %s", err)
				tx.EGPLog.Error = fmt.Sprintf("%s; CalculateEffectiveGasPricePercentage#1: %s", tx.EGPLog.Error, err)
			}
		} else {
			// Save percentage for later logging
			tx.EGPLog.Percentage = egpPercentage
		}

		// If EGP is disabled we use tx GasPrice (MaxEffectivePercentage=255)
		if !f.effectiveGasPrice.IsEnabled() {
			egpPercentage = state.MaxEffectivePercentage
		}

		// Assign applied EGP percentage to tx (TxTracker)
		tx.EGPPercentage = egpPercentage

		effectivePercentageAsDecodedHex, err := hex.DecodeHex(fmt.Sprintf("%x", tx.EGPPercentage))
		if err != nil {
			return nil, err
		}

		executorBatchRequest.Transactions = append(executorBatchRequest.Transactions, effectivePercentageAsDecodedHex...)
	}

	log.Infof("processing tx: %s. Batch.BatchNumber: %d, batchNumber: %d, oldStateRoot: %s, txHash: %s, L1InfoRootIndex: %d",
		hashStr, f.wipBatch.batchNumber, executorBatchRequest.BatchNumber, executorBatchRequest.OldStateRoot, hashStr, f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)

	processBatchResponse, err := f.state.ProcessBatchV2(ctx, executorBatchRequest, false)

	if err != nil && errors.Is(err, runtime.ErrExecutorDBError) {
		log.Errorf("failed to process transaction: %s", err)
		return nil, err
	} else if err == nil && !processBatchResponse.IsRomLevelError && len(processBatchResponse.BlockResponses) == 0 && tx != nil {
		f.halt(ctx, fmt.Errorf("executor returned no errors and no responses for tx: %s", tx.HashStr))
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

// isDeadlineEncountered returns true if any closing signal deadline is encountered
func (f *finalizer) isDeadlineEncountered() bool {
	// Forced batch deadline
	if f.nextForcedBatchDeadline != 0 && now().Unix() >= f.nextForcedBatchDeadline {
		log.Infof("closing batch %d, forced batch deadline encountered.", f.wipBatch.batchNumber)
		return true
	}
	//TODO: rename f.cfg.TimestampResolution to BatchTime or BatchMaxTime
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

// halt halts the finalizer
func (f *finalizer) halt(ctx context.Context, err error) {
	f.haltError = err
	f.haltFinalizer.Store(true)

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
}
