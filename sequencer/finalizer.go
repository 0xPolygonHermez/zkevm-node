package sequencer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

const oneHundred = 100

var (
	now = time.Now
)

// finalizer represents the finalizer component of the sequencer.
type finalizer struct {
	cfg                FinalizerCfg
	txsStore           TxsStore
	closingSignalCh    ClosingSignalCh
	isSynced           func(ctx context.Context) bool
	sequencerAddress   common.Address
	worker             workerInterface
	dbManager          dbManagerInterface
	executor           stateInterface
	batch              *WipBatch
	batchConstraints   batchConstraints
	processRequest     state.ProcessRequest
	sharedResourcesMux *sync.RWMutex
	lastGERHash        common.Hash
	// closing signals
	nextGER                   common.Hash
	nextGERDeadline           int64
	nextGERMux                *sync.RWMutex
	nextForcedBatches         []state.ForcedBatch
	nextForcedBatchDeadline   int64
	nextForcedBatchesMux      *sync.RWMutex
	nextSendingToL1Deadline   int64
	nextSendingToL1TimeoutMux *sync.RWMutex
	handlingL2Reorg           bool
}

// WipBatch represents a work-in-progress batch.
type WipBatch struct {
	batchNumber        uint64
	coinbase           common.Address
	initialStateRoot   common.Hash
	stateRoot          common.Hash
	localExitRoot      common.Hash
	timestamp          uint64
	globalExitRoot     common.Hash // 0x000...0 (ZeroHash) means to not update
	remainingResources batchResources
	countOfTxs         int
}

func (w *WipBatch) isEmpty() bool {
	return w.countOfTxs == 0
}

// newFinalizer returns a new instance of Finalizer.
func newFinalizer(
	cfg FinalizerCfg,
	worker workerInterface,
	dbManager dbManagerInterface,
	executor stateInterface,
	sequencerAddr common.Address,
	isSynced func(ctx context.Context) bool,
	closingSignalCh ClosingSignalCh,
	txsStore TxsStore,
	batchConstraints batchConstraints,
) *finalizer {
	return &finalizer{
		cfg:                cfg,
		txsStore:           txsStore,
		closingSignalCh:    closingSignalCh,
		isSynced:           isSynced,
		sequencerAddress:   sequencerAddr,
		worker:             worker,
		dbManager:          dbManager,
		executor:           executor,
		batch:              new(WipBatch),
		batchConstraints:   batchConstraints,
		processRequest:     state.ProcessRequest{},
		sharedResourcesMux: new(sync.RWMutex),
		lastGERHash:        state.ZeroHash,
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

// Start starts the finalizer.
func (f *finalizer) Start(ctx context.Context, batch *WipBatch, processingReq *state.ProcessRequest) {
	var err error
	if batch != nil {
		f.batch = batch
	} else {
		f.batch, err = f.dbManager.GetWIPBatch(ctx)
		if err != nil {
			log.Fatalf("failed to get work-in-progress batch from DB, Err: %s", err)
		}
	}

	if processingReq == nil {
		log.Fatal("processingReq should not be nil")
	} else {
		f.processRequest = *processingReq
	}

	// Closing signals receiver
	go f.listenForClosingSignals(ctx)

	// Processing transactions and finalizing batches
	f.finalizeBatches(ctx)
}

func (f *finalizer) SortForcedBatches(fb []state.ForcedBatch) []state.ForcedBatch {
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
			f.nextForcedBatches = f.SortForcedBatches(append(f.nextForcedBatches, fb))
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
		// Too much time without batches in L1 ch
		case <-f.closingSignalCh.SendingToL1TimeoutCh:
			log.Debug("finalizer received timeout for sending to L1")
			f.nextSendingToL1TimeoutMux.Lock()
			if f.nextSendingToL1Deadline == 0 {
				f.setNextSendingToL1Deadline()
			}
			f.nextSendingToL1TimeoutMux.Unlock()
		}
	}
}

// finalizeBatches runs the endless loop for processing transactions finalizing batches.
func (f *finalizer) finalizeBatches(ctx context.Context) {
	for {
		start := now()
		log.Debug("finalizer init loop")
		tx := f.worker.GetBestFittingTx(f.batch.remainingResources)
		metrics.WorkerProcessingTime(time.Since(start))
		if tx != nil {
			f.sharedResourcesMux.Lock()
			log.Debugf("processing tx: %s", tx.Hash.Hex())
			err := f.processTransaction(ctx, tx)
			if err != nil {
				log.Errorf("failed to process transaction in finalizeBatches, Err: %v", err)
			}

			f.sharedResourcesMux.Unlock()
		} else {
			// wait for new txs
			log.Debugf("no transactions to be processed. Sleeping for %v", f.cfg.SleepDurationInMs.Duration)
			if f.cfg.SleepDurationInMs.Duration > 0 {
				time.Sleep(f.cfg.SleepDurationInMs.Duration)
			}
		}

		if f.isDeadlineEncountered() {
			log.Infof("Closing batch: %d, because deadline was encountered.", f.batch.batchNumber)
			f.finalizeBatch(ctx)
		} else if f.isBatchFull() || f.isBatchAlmostFull() {
			log.Infof("Closing batch: %d, because it's almost full.", f.batch.batchNumber)
			f.finalizeBatch(ctx)
		}

		if err := ctx.Err(); err != nil {
			log.Infof("Stopping finalizer because of context, err: %s", err)
			return
		}
	}
}

func (f *finalizer) isBatchFull() bool {
	if f.batch.countOfTxs >= int(f.batchConstraints.MaxTxsPerBatch) {
		log.Infof("Closing batch: %d, because it's full.", f.batch.batchNumber)
		return true
	}
	return false
}

// finalizeBatch retries to until successful closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new empty batch
func (f *finalizer) finalizeBatch(ctx context.Context) {
	start := time.Now()
	defer func() {
		metrics.ProcessingTime(time.Since(start))
	}()
	f.txsStore.Wg.Wait()
	var err error
	f.batch, err = f.newWIPBatch(ctx)
	for err != nil {
		log.Errorf("failed to create new work-in-progress batch, Err: %s", err)
		f.batch, err = f.newWIPBatch(ctx)
	}
}

func (f *finalizer) halt(ctx context.Context, err error) {
	debugInfo := &state.DebugInfo{
		ErrorType: state.DebugInfoErrorType_FINALIZER_HALT,
		Timestamp: time.Now(),
		Payload:   err.Error(),
	}
	debugInfoErr := f.dbManager.AddDebugInfo(ctx, debugInfo, nil)
	if debugInfoErr != nil {
		log.Errorf("error storing finalizer halt debug info: %v", debugInfoErr)
	}

	for {
		log.Errorf("fatal error: %s", err)
		log.Error("halting the finalizer")
		time.Sleep(5 * time.Second) //nolint:gomnd
	}
}

// newWIPBatch closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new empty batch
func (f *finalizer) newWIPBatch(ctx context.Context) (*WipBatch, error) {
	f.sharedResourcesMux.Lock()
	defer f.sharedResourcesMux.Unlock()

	var err error
	if f.batch.stateRoot.String() == "" || f.batch.localExitRoot.String() == "" {
		return nil, errors.New("state root and local exit root must have value to close batch")
	}

	// We need to process the batch to update the state root before closing the batch
	if f.batch.initialStateRoot == f.batch.stateRoot {
		log.Info("reprocessing batch because the state root has not changed...")
		err := f.processTransaction(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	// Reprocess full batch as sanity check
	processBatchResponse, err := f.reprocessFullBatch(ctx, f.batch.batchNumber, f.batch.stateRoot)
	if err != nil || !processBatchResponse.IsBatchProcessed {
		log.Info("halting the finalizer because of a reprocessing error")
		if err != nil {
			f.halt(ctx, fmt.Errorf("failed to reprocess batch, err: %v", err))
		} else {
			f.halt(ctx, fmt.Errorf("out of counters during reprocessFullBath"))
		}
	}

	// Close the current batch
	err = f.closeBatch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to close batch, err: %w", err)
	}

	// Metadata for the next batch
	stateRoot := f.batch.stateRoot
	lastBatchNumber := f.batch.batchNumber

	// Process Forced Batches
	if len(f.nextForcedBatches) > 0 {
		lastBatchNumber, stateRoot, err = f.processForcedBatches(ctx, lastBatchNumber, stateRoot)
		if err != nil {
			log.Warnf("failed to process forced batch, err: %s", err)
		}
	}

	// Take into consideration the GER
	f.nextGERMux.Lock()
	if f.nextGER != state.ZeroHash {
		f.lastGERHash = f.nextGER
	}
	f.nextGER = state.ZeroHash
	f.nextGERDeadline = 0
	f.nextGERMux.Unlock()

	// Reset nextSendingToL1Deadline
	f.nextSendingToL1TimeoutMux.Lock()
	f.nextSendingToL1Deadline = 0
	f.nextSendingToL1TimeoutMux.Unlock()

	batch, err := f.openWIPBatch(ctx, lastBatchNumber+1, f.lastGERHash, stateRoot)
	if err == nil {
		f.processRequest.Timestamp = batch.timestamp
		f.processRequest.BatchNumber = batch.batchNumber
		f.processRequest.OldStateRoot = stateRoot
		f.processRequest.GlobalExitRoot = batch.globalExitRoot
		f.processRequest.Transactions = make([]byte, 0, 1)
	}

	return batch, err
}

// processTransaction processes a single transaction.
func (f *finalizer) processTransaction(ctx context.Context, tx *TxTracker) error {
	var txHash string
	if tx != nil {
		txHash = tx.Hash.String()
	}
	log := log.WithFields("txHash", txHash, "batchNumber", f.processRequest.BatchNumber)
	start := time.Now()
	defer func() {
		metrics.ProcessingTime(time.Since(start))
	}()

	var ger common.Hash
	if f.batch.isEmpty() {
		ger = f.batch.globalExitRoot
	} else {
		ger = state.ZeroHash
	}

	f.processRequest.GlobalExitRoot = ger
	if tx != nil {
		f.processRequest.Transactions = tx.RawTx
	} else {
		f.processRequest.Transactions = []byte{}
	}
	hash := "nil"
	if tx != nil {
		hash = tx.HashStr
	}
	log.Infof("processTransaction: single tx. Batch.BatchNumber: %d, BatchNumber: %d, OldStateRoot: %s, txHash: %s, GER: %s", f.batch.batchNumber, f.processRequest.BatchNumber, f.processRequest.OldStateRoot, hash, f.processRequest.GlobalExitRoot.String())
	result, err := f.executor.ProcessBatch(ctx, f.processRequest, true)
	if err != nil {
		log.Errorf("failed to process transaction, isClaim: %v, err: %s", tx.IsClaim, err)
		return err
	}

	oldStateRoot := f.batch.stateRoot
	if len(result.Responses) > 0 && tx != nil {
		err = f.handleTxProcessResp(ctx, tx, result, oldStateRoot)
		if err != nil {
			return err
		}
	}

	// Update in-memory batch and processRequest
	f.processRequest.OldStateRoot = result.NewStateRoot
	f.batch.stateRoot = result.NewStateRoot
	f.batch.localExitRoot = result.NewLocalExitRoot
	log.Infof("processTransaction: data loaded in memory. batch.batchNumber: %d, batchNumber: %d, result.NewStateRoot: %s, result.NewLocalExitRoot: %s, oldStateRoot: %s", f.batch.batchNumber, f.processRequest.BatchNumber, result.NewStateRoot.String(), result.NewLocalExitRoot.String(), oldStateRoot.String())

	return nil
}

// handleTxProcessResp handles the response of a successful transaction processing.
func (f *finalizer) handleTxProcessResp(ctx context.Context, tx *TxTracker, result *state.ProcessBatchResponse, oldStateRoot common.Hash) error {
	// Handle Transaction Error
	if result.Responses[0].RomError != nil && !errors.Is(result.Responses[0].RomError, runtime.ErrExecutionReverted) {
		f.handleTransactionError(ctx, result, tx)
		return result.Responses[0].RomError
	}

	// Check remaining resources
	err := f.checkRemainingResources(ctx, result, tx)
	if err != nil {
		return err
	}

	// Store the processed transaction, add it to the batch and update status in the pool atomically
	f.storeProcessedTx(ctx, oldStateRoot, tx, result)

	return nil
}

func (f *finalizer) storeProcessedTx(ctx context.Context, previousL2BlockStateRoot common.Hash, tx *TxTracker, result *state.ProcessBatchResponse) {
	log.Infof("storeProcessedTx: storing processed tx: %s", tx.Hash.String())
	f.txsStore.Wg.Wait()
	txResponse := result.Responses[0]
	f.txsStore.Wg.Add(1)
	f.txsStore.Ch <- &txToStore{
		batchNumber:              f.batch.batchNumber,
		txResponse:               txResponse,
		coinbase:                 f.batch.coinbase,
		timestamp:                f.batch.timestamp,
		previousL2BlockStateRoot: previousL2BlockStateRoot,
	}

	// Delete the transaction from the efficiency list
	f.worker.DeleteTx(tx.Hash, tx.From)
	log.Debug("tx deleted from efficiency list", "txHash", tx.Hash.String(), "from", tx.From.Hex())

	start := time.Now()
	txsToDelete := f.worker.UpdateAfterSingleSuccessfulTxExecution(tx.From, result.ReadWriteAddresses)
	for _, txToDelete := range txsToDelete {
		err := f.dbManager.UpdateTxStatus(ctx, txToDelete.Hash, pool.TxStatusFailed, false)
		if err != nil {
			log.Errorf("failed to update status to failed in the pool for tx: %s, err: %s", txToDelete.Hash.String(), err)
		}
	}
	metrics.WorkerProcessingTime(time.Since(start))
	f.batch.countOfTxs += 1
}

// handleTransactionError handles the error of a transaction
func (f *finalizer) handleTransactionError(ctx context.Context, result *state.ProcessBatchResponse, tx *TxTracker) {
	txResponse := result.Responses[0]
	errorCode := executor.RomErrorCode(txResponse.RomError)
	addressInfo := result.ReadWriteAddresses[tx.From]
	log.Infof("handleTransactionError: error in tx: %s, errorCode: %d", tx.Hash.String(), errorCode)

	if executor.IsROMOutOfCountersError(errorCode) {
		log.Errorf("ROM out of counters error, marking tx with Hash: %s as INVALID, errorCode: %s", tx.Hash.String(), errorCode.String())
		start := time.Now()
		f.worker.DeleteTx(tx.Hash, tx.From)
		metrics.WorkerProcessingTime(time.Since(start))
		go func() {
			err := f.dbManager.UpdateTxStatus(ctx, tx.Hash, pool.TxStatusInvalid, false)
			if err != nil {
				log.Errorf("failed to update status to failed in the pool for tx: %s, err: %s", tx.Hash.String(), err)
			}
		}()
	} else if (executor.IsInvalidNonceError(errorCode) || executor.IsInvalidBalanceError(errorCode)) && !tx.IsClaim {
		var (
			nonce   *uint64
			balance *big.Int
		)
		if addressInfo != nil {
			nonce = addressInfo.Nonce
			balance = addressInfo.Balance
		}
		start := time.Now()
		log.Errorf("intrinsic error, moving tx with Hash: %s to NOT READY nonce(%d) balance(%s) cost(%s), err: %s", tx.Hash, nonce, balance.String(), tx.Cost.String(), txResponse.RomError)
		txsToDelete := f.worker.MoveTxToNotReady(tx.Hash, tx.From, nonce, balance)
		for _, txToDelete := range txsToDelete {
			err := f.dbManager.UpdateTxStatus(ctx, txToDelete.Hash, pool.TxStatusFailed, false)
			if err != nil {
				log.Errorf("failed to update status to failed in the pool for tx: %s, err: %s", txToDelete.Hash.String(), err)
			}
		}
		metrics.WorkerProcessingTime(time.Since(start))
	} else {
		// Delete the transaction from the efficiency list
		f.worker.DeleteTx(tx.Hash, tx.From)
		log.Debug("tx deleted from efficiency list", "txHash", tx.Hash.String(), "from", tx.From.Hex(), "isClaim", tx.IsClaim)

		// Update the status of the transaction to failed
		err := f.dbManager.UpdateTxStatus(ctx, tx.Hash, pool.TxStatusFailed, false)
		if err != nil {
			log.Errorf("failed to update status to failed in the pool for tx: %s, err: %s", tx.Hash.String(), err)
		}
	}
}

// syncWithState syncs the WIP batch and processRequest with the state
func (f *finalizer) syncWithState(ctx context.Context, lastBatchNum *uint64) error {
	f.sharedResourcesMux.Lock()
	defer f.sharedResourcesMux.Unlock()
	f.txsStore.Wg.Wait()

	var lastBatch *state.Batch
	var err error
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}
	if lastBatchNum == nil {
		lastBatch, err = f.dbManager.GetLastBatch(ctx)
		if err != nil {
			return fmt.Errorf("failed to get last batch, err: %w", err)
		}
	} else {
		lastBatch, err = f.dbManager.GetBatchByNumber(ctx, *lastBatchNum, nil)
		if err != nil {
			return fmt.Errorf("failed to get last batch, err: %w", err)
		}
	}

	batchNum := lastBatch.BatchNumber
	lastBatchNum = &batchNum

	isClosed, err := f.dbManager.IsBatchClosed(ctx, *lastBatchNum)
	if err != nil {
		return fmt.Errorf("failed to check if batch is closed, err: %w", err)
	}
	log.Infof("Batch %d isClosed: %v", batchNum, isClosed)
	if isClosed {
		ger, _, err := f.dbManager.GetLatestGer(ctx, f.cfg.GERFinalityNumberOfBlocks)
		if err != nil {
			return fmt.Errorf("failed to get latest ger, err: %w", err)
		}

		oldStateRoot := lastBatch.StateRoot
		f.batch, err = f.openWIPBatch(ctx, *lastBatchNum+1, ger.GlobalExitRoot, oldStateRoot)
		if err != nil {
			return err
		}
	} else {
		f.batch, err = f.dbManager.GetWIPBatch(ctx)
		if err != nil {
			return fmt.Errorf("failed to get work-in-progress batch, err: %w", err)
		}
	}
	log.Infof("Initial Batch: %+v", f.batch)
	log.Infof("Initial Batch.StateRoot: %s", f.batch.stateRoot.String())
	log.Infof("Initial Batch.GER: %s", f.batch.globalExitRoot.String())
	log.Infof("Initial Batch.Coinbase: %s", f.batch.coinbase.String())
	log.Infof("Initial Batch.InitialStateRoot: %s", f.batch.initialStateRoot.String())
	log.Infof("Initial Batch.localExitRoot: %s", f.batch.localExitRoot.String())

	f.processRequest = state.ProcessRequest{
		BatchNumber:    *lastBatchNum,
		OldStateRoot:   f.batch.stateRoot,
		GlobalExitRoot: f.batch.globalExitRoot,
		Coinbase:       f.sequencerAddress,
		Timestamp:      f.batch.timestamp,
		Transactions:   make([]byte, 0, 1),
		Caller:         state.SequencerCallerLabel,
	}

	log.Infof("synced with state, lastBatchNum: %d. State root: %s", *lastBatchNum, f.batch.initialStateRoot.Hex())

	return nil
}

// processForcedBatches processes all the forced batches that are pending to be processed
func (f *finalizer) processForcedBatches(ctx context.Context, lastBatchNumberInState uint64, stateRoot common.Hash) (uint64, common.Hash, error) {
	f.nextForcedBatchesMux.Lock()
	defer f.nextForcedBatchesMux.Unlock()
	f.nextForcedBatchDeadline = 0

	dbTx, err := f.dbManager.BeginStateTransaction(ctx)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("failed to begin state transaction, err: %w", err)
	}
	lastTrustedForcedBatchNumber, err := f.dbManager.GetLastTrustedForcedBatchNumber(ctx, dbTx)
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
			lastBatchNumberInState, stateRoot = f.processForcedBatch(lastBatchNumberInState, stateRoot, forcedBatch)
			nextForcedBatchNum += 1
		}
		// Process the current forced batch from the channel queue
		lastBatchNumberInState, stateRoot = f.processForcedBatch(lastBatchNumberInState, stateRoot, forcedBatch)
		nextForcedBatchNum += 1
	}
	f.nextForcedBatches = make([]state.ForcedBatch, 0)

	return lastBatchNumberInState, stateRoot, nil
}

func (f *finalizer) processForcedBatch(lastBatchNumberInState uint64, stateRoot common.Hash, forcedBatch state.ForcedBatch) (uint64, common.Hash) {
	processRequest := state.ProcessRequest{
		BatchNumber:    lastBatchNumberInState + 1,
		OldStateRoot:   stateRoot,
		GlobalExitRoot: forcedBatch.GlobalExitRoot,
		Transactions:   forcedBatch.RawTxsData,
		Coinbase:       f.sequencerAddress,
		Timestamp:      uint64(now().Unix()),
		Caller:         state.SequencerCallerLabel,
	}
	response, err := f.dbManager.ProcessForcedBatch(forcedBatch.ForcedBatchNumber, processRequest)
	if err != nil {
		log.Warnf("failed to process forced batch, err: %s", err)
	}

	f.nextGERMux.Lock()
	f.lastGERHash = forcedBatch.GlobalExitRoot
	f.nextGERMux.Unlock()

	stateRoot = response.NewStateRoot
	lastBatchNumberInState += 1

	return lastBatchNumberInState, stateRoot
}

// openWIPBatch opens a new batch in the state and returns it as WipBatch
func (f *finalizer) openWIPBatch(ctx context.Context, batchNum uint64, ger, stateRoot common.Hash) (*WipBatch, error) {
	dbTx, err := f.dbManager.BeginStateTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin state transaction to open batch, err: %w", err)
	}

	// open next batch
	openBatchResp, err := f.openBatch(ctx, batchNum, ger, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return nil, fmt.Errorf(
				"failed to rollback dbTx: %s. Rollback err: %w",
				rollbackErr.Error(), err,
			)
		}
		return nil, err
	}
	if err := dbTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit database transaction for opening a batch, err: %w", err)
	}

	// Check if synchronizer is up-to-date
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	return &WipBatch{
		batchNumber:        batchNum,
		coinbase:           f.sequencerAddress,
		initialStateRoot:   stateRoot,
		stateRoot:          stateRoot,
		timestamp:          uint64(openBatchResp.Timestamp.Unix()),
		globalExitRoot:     ger,
		remainingResources: getMaxRemainingResources(f.batchConstraints),
	}, err
}

// closeBatch closes the current batch in the state
func (f *finalizer) closeBatch(ctx context.Context) error {
	transactions, err := f.dbManager.GetTransactionsByBatchNumber(ctx, f.batch.batchNumber)
	if err != nil {
		return fmt.Errorf("failed to get transactions from transactions, err: %w", err)
	}
	for i, tx := range transactions {
		log.Infof("closeBatch: BatchNum: %d, Tx position: %d, txHash: %s", f.batch.batchNumber, i, tx.Hash().String())
	}
	receipt := ClosingBatchParameters{
		BatchNumber:   f.batch.batchNumber,
		StateRoot:     f.batch.stateRoot,
		LocalExitRoot: f.batch.localExitRoot,
		Txs:           transactions,
	}
	return f.dbManager.CloseBatch(ctx, receipt)
}

// openBatch opens a new batch in the state
func (f *finalizer) openBatch(ctx context.Context, num uint64, ger common.Hash, dbTx pgx.Tx) (state.ProcessingContext, error) {
	processingCtx := state.ProcessingContext{
		BatchNumber:    num,
		Coinbase:       f.sequencerAddress,
		Timestamp:      now(),
		GlobalExitRoot: ger,
	}
	err := f.dbManager.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		return state.ProcessingContext{}, fmt.Errorf("failed to open new batch, err: %w", err)
	}

	return processingCtx, nil
}

// reprocessBatch reprocesses a batch used as sanity check
func (f *finalizer) reprocessFullBatch(ctx context.Context, batchNum uint64, expectedStateRoot common.Hash) (*state.ProcessBatchResponse, error) {
	batch, err := f.dbManager.GetBatchByNumber(ctx, batchNum, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch by number, err: %v", err)
	}
	processRequest := state.ProcessRequest{
		BatchNumber:    batch.BatchNumber,
		GlobalExitRoot: batch.GlobalExitRoot,
		OldStateRoot:   f.batch.initialStateRoot,
		Transactions:   batch.BatchL2Data,
		Coinbase:       batch.Coinbase,
		Timestamp:      uint64(batch.Timestamp.Unix()),
		Caller:         state.DiscardCallerLabel,
	}
	log.Infof("reprocessFullBatch: BatchNumber: %d, OldStateRoot: %s, Ger: %s", batch.BatchNumber, f.batch.initialStateRoot.String(), batch.GlobalExitRoot.String())
	txs, _, err := state.DecodeTxs(batch.BatchL2Data)
	if err != nil {
		log.Errorf("reprocessFullBatch: error decoding BatchL2Data before reprocessing full batch: %d. Error: %v", batch.BatchNumber, err)
		return nil, fmt.Errorf("reprocessFullBatch: error decoding BatchL2Data before reprocessing full batch: %d. Error: %v", batch.BatchNumber, err)
	}
	for i, tx := range txs {
		log.Infof("reprocessFullBatch: Tx position %d. TxHash: %s", i, tx.Hash())
	}

	result, err := f.executor.ProcessBatch(ctx, processRequest, false)
	if err != nil {
		log.Errorf("failed to process batch, err: %s", err)
		return nil, err
	}

	if !result.IsBatchProcessed {
		timestamp := time.Now()
		log.Errorf("failed to process batch %v because OutOfCounters", batch.BatchNumber)
		payload, err := json.Marshal(processRequest)
		if err != nil {
			log.Errorf("error marshaling payload: %v", err)
		} else {
			debugInfo := &state.DebugInfo{
				ErrorType: state.DebugInfoErrorType_OOC_ERROR_ON_REPROCESS_FULL_BATCH,
				Timestamp: timestamp,
				Payload:   string(payload),
			}
			err = f.dbManager.AddDebugInfo(ctx, debugInfo, nil)
			if err != nil {
				log.Errorf("error storing payload: %v", err)
			}
		}
		return nil, fmt.Errorf("failed to process batch because OutOfCounters error")
	}

	if result.NewStateRoot != expectedStateRoot {
		log.Errorf("batchNumber: %d, reprocessed batch has different state root, expected: %s, got: %s", batch.BatchNumber, expectedStateRoot.Hex(), result.NewStateRoot.Hex())
		return nil, fmt.Errorf("batchNumber: %d, reprocessed batch has different state root, expected: %s, got: %s", batch.BatchNumber, expectedStateRoot.Hex(), result.NewStateRoot.Hex())
	}

	return result, nil
}

func (f *finalizer) getLastBatchNumAndOldStateRoot(ctx context.Context) (uint64, common.Hash, error) {
	const two = 2
	var oldStateRoot common.Hash
	batches, err := f.dbManager.GetLastNBatches(ctx, two)
	if err != nil {
		return 0, common.Hash{}, fmt.Errorf("failed to get last %d batches, err: %w", two, err)
	}
	lastBatch := batches[0]

	oldStateRoot = f.getOldStateRootFromBatches(batches)
	return lastBatch.BatchNumber, oldStateRoot, nil
}

func (f *finalizer) getOldStateRootFromBatches(batches []*state.Batch) common.Hash {
	const one = 1
	const two = 2
	var oldStateRoot common.Hash
	if len(batches) == one {
		oldStateRoot = batches[0].StateRoot
	} else if len(batches) == two {
		oldStateRoot = batches[1].StateRoot
	}

	return oldStateRoot
}

// isDeadlineEncountered returns true if any closing signal deadline is encountered
func (f *finalizer) isDeadlineEncountered() bool {
	// Forced batch deadline
	if f.nextForcedBatchDeadline != 0 && now().Unix() >= f.nextForcedBatchDeadline {
		log.Infof("Closing batch: %d, forced batch deadline encountered.", f.batch.batchNumber)
		f.setNextSendingToL1Deadline()
		return true
	}
	// Global Exit Root deadline
	if f.nextGERDeadline != 0 && now().Unix() >= f.nextGERDeadline {
		log.Infof("Closing batch: %d, Global Exit Root deadline encountered.", f.batch.batchNumber)
		f.setNextSendingToL1Deadline()
		return true
	}
	// Delayed batch deadline
	if f.nextSendingToL1Deadline != 0 && now().Unix() >= f.nextSendingToL1Deadline && !f.batch.isEmpty() {
		log.Infof("Closing batch: %d, Sending to L1 deadline encountered.", f.batch.batchNumber)
		f.setNextSendingToL1Deadline()
		return true
	}
	return false
}

// checkRemainingResources checks if the transaction uses less resources than the remaining ones in the batch.
func (f *finalizer) checkRemainingResources(ctx context.Context, result *state.ProcessBatchResponse, tx *TxTracker) error {
	usedResources := batchResources{
		zKCounters: result.UsedZkCounters,
		bytes:      uint64(len(tx.RawTx)),
	}

	// Log an event in case the TX consumed more than the double of the expected for a zkCounter
	f.checkZKCounterConsumption(ctx, result.UsedZkCounters, tx)

	err := f.batch.remainingResources.sub(usedResources)
	if err != nil {
		log.Infof("current transaction exceeds the batch limit, updating metadata for tx in worker and continuing")
		start := time.Now()
		f.worker.UpdateTx(result.Responses[0].TxHash, tx.From, usedResources.zKCounters)
		metrics.WorkerProcessingTime(time.Since(start))
		return err
	}

	return nil
}

func (f *finalizer) checkZKCounterConsumption(ctx context.Context, zkCounters state.ZKCounters, tx *TxTracker) {
	events := ""

	if zkCounters.CumulativeGasUsed > tx.BatchResources.zKCounters.CumulativeGasUsed*2 {
		events += "CumulativeGasUsed "
	}
	if zkCounters.UsedKeccakHashes > tx.BatchResources.zKCounters.UsedKeccakHashes*2 {
		events += "UsedKeccakHashes "
	}
	if zkCounters.UsedPoseidonHashes > tx.BatchResources.zKCounters.UsedPoseidonHashes*2 {
		events += "UsedPoseidonHashes "
	}
	if zkCounters.UsedPoseidonPaddings > tx.BatchResources.zKCounters.UsedPoseidonPaddings*2 {
		events += "UsedPoseidonPaddings "
	}
	if zkCounters.UsedMemAligns > tx.BatchResources.zKCounters.UsedMemAligns*2 {
		events += "UsedMemAligns "
	}
	if zkCounters.UsedArithmetics > tx.BatchResources.zKCounters.UsedArithmetics*2 {
		events += "UsedArithmetics "
	}
	if zkCounters.UsedBinaries > tx.BatchResources.zKCounters.UsedBinaries*2 {
		events += "UsedBinaries "
	}
	if zkCounters.UsedSteps > tx.BatchResources.zKCounters.UsedSteps*2 {
		events += "UsedSteps "
	}

	if events != "" {
		event := &state.Event{
			EventType: state.EventType_ZKCounters_Diff + " " + events,
			Timestamp: time.Now(),
			IP:        tx.IP,
			TxHash:    tx.Hash,
		}

		err := f.dbManager.AddEvent(ctx, event, nil)
		if err != nil {
			log.Errorf("Error adding event: %v", err)
		}
	}
}

// isBatchAlmostFull checks if the current batch remaining resources are under the constraints threshold for most efficient moment to close a batch
func (f *finalizer) isBatchAlmostFull() bool {
	resources := f.batch.remainingResources
	zkCounters := resources.zKCounters
	if resources.bytes <= f.getConstraintThresholdUint64(f.batchConstraints.MaxBatchBytesSize) {
		return true
	}
	if zkCounters.UsedSteps <= f.getConstraintThresholdUint32(f.batchConstraints.MaxSteps) {
		return true
	}
	if zkCounters.UsedPoseidonPaddings <= f.getConstraintThresholdUint32(f.batchConstraints.MaxPoseidonPaddings) {
		return true
	}
	if zkCounters.UsedBinaries <= f.getConstraintThresholdUint32(f.batchConstraints.MaxBinaries) {
		return true
	}
	if zkCounters.UsedKeccakHashes <= f.getConstraintThresholdUint32(f.batchConstraints.MaxKeccakHashes) {
		return true
	}
	if zkCounters.UsedArithmetics <= f.getConstraintThresholdUint32(f.batchConstraints.MaxArithmetics) {
		return true
	}
	if zkCounters.UsedMemAligns <= f.getConstraintThresholdUint32(f.batchConstraints.MaxMemAligns) {
		return true
	}
	if zkCounters.CumulativeGasUsed <= f.getConstraintThresholdUint64(f.batchConstraints.MaxCumulativeGasUsed) {
		return true
	}
	return false
}

func (f *finalizer) setNextForcedBatchDeadline() {
	f.nextForcedBatchDeadline = now().Unix() + int64(f.cfg.ForcedBatchDeadlineTimeoutInSec.Duration.Seconds())
}

func (f *finalizer) setNextGERDeadline() {
	f.nextGERDeadline = now().Unix() + int64(f.cfg.GERDeadlineTimeoutInSec.Duration.Seconds())
}

func (f *finalizer) setNextSendingToL1Deadline() {
	f.nextSendingToL1Deadline = now().Unix() + int64(f.cfg.SendingToL1DeadlineTimeoutInSec.Duration.Seconds())
}

func (f *finalizer) getConstraintThresholdUint64(input uint64) uint64 {
	return input * uint64(f.cfg.ResourcePercentageToCloseBatch) / oneHundred
}

func (f *finalizer) getConstraintThresholdUint32(input uint32) uint32 {
	return uint32(input*f.cfg.ResourcePercentageToCloseBatch) / oneHundred
}
