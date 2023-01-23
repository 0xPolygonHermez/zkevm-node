package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

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

	if processingReq != nil {
		f.processRequest = *processingReq
	} else {
		f.processRequest, err = f.prepareProcessRequestFromState(ctx, false)
		if err != nil {
			log.Fatalf("failed to prepare process request from state, Err: %s", err)
		}
	}

	// Closing signals receiver
	go f.listenForClosingSignals(ctx)

	// Processing transactions and finalizing batches
	f.finalizeBatches(ctx)
}

// listenForClosingSignals listens for signals for the batch and sets the deadline for when they need to be closed.
func (f *finalizer) listenForClosingSignals(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Infof("finalizer closing signal listener received context done, Err: %s", ctx.Err())
			return
		// Forced  batch ch
		case fb := <-f.closingSignalCh.ForcedBatchCh:
			f.nextForcedBatchesMux.Lock()
			f.nextForcedBatches = append(f.nextForcedBatches, fb) // TODO: change insert sort if not exists
			if f.nextForcedBatchDeadline == 0 {
				f.setNextForcedBatchDeadline()
			}
			f.nextForcedBatchesMux.Unlock()
		// globalExitRoot ch
		case ger := <-f.closingSignalCh.GERCh:
			f.nextGERMux.Lock()
			f.nextGER = ger
			if f.nextGERDeadline == 0 {
				f.setNextGERDeadline()
			}
			f.nextGERMux.Unlock()
		// L2Reorg ch
		case l2ReorgEvent := <-f.closingSignalCh.L2ReorgCh:
			f.sharedResourcesMux.Lock()
			f.handlingL2Reorg = true
			go f.worker.HandleL2Reorg(l2ReorgEvent.TxHashes)
			err := f.syncWithState(ctx, nil)
			if err != nil {
				log.Errorf("failed to sync with state, Err: %s", err)
			}
			f.handlingL2Reorg = false
			f.sharedResourcesMux.Unlock()
		// Too much time without batches in L1 ch
		case <-f.closingSignalCh.SendingToL1TimeoutCh:
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
		tx := f.worker.GetBestFittingTx(f.batch.remainingResources)
		if tx != nil {
			_ = f.processTransaction(ctx, tx)
		} else {
			if f.isCurrBatchAboveLimitWindow() {
				// Wait for all transactions to be stored in the DB
				f.txsStore.Wg.Wait()
				// The perfect moment to finalize the batch
				f.finalizeBatch(ctx)
			} else {
				// wait for new txs
				if f.cfg.SleepDurationInMs.Duration > 0 {
					time.Sleep(f.cfg.SleepDurationInMs.Duration)
				}
			}
		}

		if f.isDeadlineEncountered() || f.batch.countOfTxs >= int(f.batchConstraints.MaxTxsPerBatch) {
			f.finalizeBatch(ctx)
		}

		if err := ctx.Err(); err != nil {
			log.Infof("Stopping finalizer because of context, err: %s", err)
			return
		}
	}
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

// newWIPBatch closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new empty batch
func (f *finalizer) newWIPBatch(ctx context.Context) (*WipBatch, error) {
	var err error
	// Passing the batch without txs to the executor in order to update the State
	if f.batch.countOfTxs == 0 {
		// backup current sequence
		err = f.processTransaction(ctx, nil)
		for err != nil {
			log.Errorf("failed to process tx, err: %w", err)
			err = f.processTransaction(ctx, nil)
		}
	}

	if f.batch.stateRoot.String() == "" || f.batch.localExitRoot.String() == "" {
		return nil, errors.New("state root and local exit root must have value to close batch")
	}
	err = f.closeBatch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to close batch, err: %w", err)
	}
	// Reprocessing batch as sanity check
	go func() {
		err := f.reprocessBatch(ctx)
		if err != nil {
			// TODO: design error handling for reprocessing
			log.Errorf("failed to reprocess batch, err: %s", err)
			return
		}
	}()

	// Metadata for the next batch
	stateRoot := f.batch.stateRoot
	lastBatchNumber := f.batch.batchNumber

	// Process Forced Batches
	if len(f.nextForcedBatches) > 0 {
		lastBatchNumber, stateRoot, err = f.processForcedBatches(lastBatchNumber, stateRoot)
		if err != nil {
			log.Errorf("failed to process forced batch, err: %s", err)
		}
	}

	// Take into consideration the GER
	f.nextGERMux.Lock()
	ger := f.nextGER
	f.nextGER = state.ZeroHash
	f.nextGERDeadline = 0
	f.nextGERMux.Unlock()

	// Reset nextSendingToL1Deadline
	f.nextSendingToL1TimeoutMux.Lock()
	f.nextSendingToL1Deadline = 0
	f.nextSendingToL1TimeoutMux.Unlock()

	return f.openWIPBatch(ctx, lastBatchNumber+1, ger, stateRoot)
}

// processTransaction processes a single transaction.
func (f *finalizer) processTransaction(ctx context.Context, tx *TxTracker) error {
	start := time.Now()
	defer func() {
		metrics.ProcessingTime(time.Since(start))
	}()
	f.sharedResourcesMux.Lock()
	defer f.sharedResourcesMux.Unlock()

	var ger common.Hash
	if f.batch.isEmpty() {
		ger = f.batch.globalExitRoot
	} else {
		ger = state.ZeroHash
	}

	f.processRequest.GlobalExitRoot = ger
	f.processRequest.Transactions = tx.RawTx
	result, err := f.executor.ProcessBatch(ctx, f.processRequest)
	if err != nil {
		log.Errorf("failed to process transaction, err: %s", err)
		return err
	}

	if result != nil && result.ExecutorError != nil {
		if result.ExecutorError == state.ErrBatchAlreadyClosed || result.ExecutorError == state.ErrInvalidBatchNumber {
			log.Warnf("unexpected state local vs DB: %s", result.ExecutorError)
			log.Info("reloading local sequence")
			f.batch, err = f.dbManager.GetWIPBatch(ctx)
			if err != nil {
				log.Errorf("failed to get WIP Batch from state, err: %s", err)
			}
			return err
		}
		return fmt.Errorf("failed processing transaction, err: %w", result.ExecutorError)
	} else {
		err = f.handleSuccessfulTxProcessResp(ctx, tx, result)
		if err != nil {
			return err
		}
	}

	return nil
}

// handleSuccessfulTxProcessResp handles the response of a successful transaction processing.
func (f *finalizer) handleSuccessfulTxProcessResp(ctx context.Context, tx *TxTracker, result *state.ProcessBatchResponse) error {
	if tx == nil {
		return nil
	}

	txResponse := result.Responses[0]
	// Handle Transaction Error
	if txResponse.RomError != nil {
		f.handleTransactionError(ctx, txResponse, result, tx)
		return txResponse.RomError
	}

	// Check remaining resources
	err := f.checkRemainingResources(result, tx, txResponse)
	if err != nil {
		return err
	}

	// We have a successful processing if we are here, updating metadata
	previousL2BlockStateRoot := f.batch.stateRoot
	f.processRequest.OldStateRoot = result.NewStateRoot
	f.batch.stateRoot = result.NewStateRoot
	f.batch.localExitRoot = result.NewLocalExitRoot

	// Store the processed transaction, add it to the batch and update status in the pool atomically
	f.txsStore.Wg.Add(1)
	f.txsStore.Ch <- &txToStore{
		batchNumber:              f.batch.batchNumber,
		txResponse:               txResponse,
		previousL2BlockStateRoot: previousL2BlockStateRoot,
	}

	f.worker.UpdateAfterSingleSuccessfulTxExecution(tx.From, result.ReadWriteAddresses)
	f.batch.countOfTxs += 1

	return nil
}

// handleTransactionError handles the error of a transaction
func (f *finalizer) handleTransactionError(ctx context.Context, txResponse *state.ProcessTransactionResponse, result *state.ProcessBatchResponse, tx *TxTracker) {
	errorCode := executor.RomErrorCode(txResponse.RomError)
	addressInfo := result.ReadWriteAddresses[tx.From]

	if executor.IsROMOutOfCountersError(errorCode) {
		f.worker.DeleteTx(tx.Hash, tx.From)
		go func() {
			err := f.dbManager.UpdateTxStatus(ctx, tx.Hash, pool.TxStatusInvalid)
			if err != nil {
				log.Errorf("failed to update tx status, err: %s", err)
			}
		}()
	} else if executor.IsIntrinsicError(errorCode) {
		var (
			nonce   *uint64
			balance *big.Int
		)
		if addressInfo != nil {
			nonce = addressInfo.Nonce
			balance = addressInfo.Balance
		}
		f.worker.MoveTxToNotReady(tx.Hash, tx.From, nonce, balance)
	}
}

// syncWithState syncs the WIP batch and processRequest with the state
func (f *finalizer) syncWithState(ctx context.Context, lastBatchNum *uint64) error {
	var err error
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}
	if lastBatchNum == nil {
		batchNum, err := f.dbManager.GetLastBatchNumber(ctx)
		for err != nil {
			return fmt.Errorf("failed to get last batch number, err: %w", err)
		}
		lastBatchNum = &batchNum
	}

	isClosed, err := f.dbManager.IsBatchClosed(ctx, *lastBatchNum)
	if err != nil {
		return fmt.Errorf("failed to check if batch is closed, err: %w", err)
	}
	if isClosed {
		ger, _, err := f.dbManager.GetLatestGer(ctx, f.cfg.GERFinalityNumberOfBlocks)
		if err != nil {
			return fmt.Errorf("failed to get latest ger, err: %w", err)
		}
		_, oldStateRoot, err := f.getLastBatchNumAndStateRoot(ctx)
		if err != nil {
			return fmt.Errorf("failed to get old state root, err: %w", err)
		}
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

	f.processRequest = state.ProcessRequest{
		BatchNumber:    *lastBatchNum,
		OldStateRoot:   f.batch.initialStateRoot,
		GlobalExitRoot: f.batch.globalExitRoot,
		Coinbase:       f.sequencerAddress,
		Timestamp:      f.batch.timestamp,
		Transactions:   make([]byte, 0, 1),
		Caller:         state.SequencerCallerLabel,
	}

	return nil
}

// processForcedBatches processes all the forced batches that are pending to be processed
func (f *finalizer) processForcedBatches(lastBatchNumberInState uint64, stateRoot common.Hash) (uint64, common.Hash, error) {
	f.nextForcedBatchesMux.Lock()
	// TODO: query for last included forced batch from database
	// Then we do integrity check - if it is lower than the last batch number in the state we skip
	// If it is higher than the last forced batch number in the state we get the one in order from database.

	for _, forcedBatch := range f.nextForcedBatches {
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
		// TODO: design error handling for forced batches
		if err != nil {
			return lastBatchNumberInState, stateRoot, err
		} else {
			stateRoot = response.NewStateRoot
			lastBatchNumberInState += 1
		}
	}
	f.nextForcedBatches = make([]state.ForcedBatch, 0)
	f.nextForcedBatchesMux.Unlock()
	return lastBatchNumberInState, stateRoot, nil
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
	receipt := ClosingBatchParameters{
		BatchNumber:   f.batch.batchNumber,
		StateRoot:     f.batch.stateRoot,
		LocalExitRoot: f.processRequest.GlobalExitRoot,
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
func (f *finalizer) reprocessBatch(ctx context.Context) error {
	processRequest, err := f.prepareProcessRequestFromState(ctx, true)
	if err != nil {
		log.Errorf("failed to prepare process request for reprocessing batch, err: %s", err)
		return err
	}
	processRequest.Caller = state.DiscardCallerLabel
	result, err := f.executor.ProcessBatch(ctx, processRequest)
	if err != nil || (result != nil && result.ExecutorError != nil) {
		if result != nil && result.ExecutorError != nil {
			err = result.ExecutorError
		}
		log.Errorf("failed to reprocess batch, err: %s", err)
		return err
	}

	return nil
}

// prepareProcessRequestFromState prepares process request from state
func (f *finalizer) prepareProcessRequestFromState(ctx context.Context, fetchTxs bool) (state.ProcessRequest, error) {
	const two = 2

	var (
		txs          []byte
		batchNum     uint64
		oldStateRoot common.Hash
		err          error
	)

	if fetchTxs {
		var lastClosedBatch *state.Batch
		batches, err := f.dbManager.GetLastNBatches(ctx, two)
		if err != nil {
			return state.ProcessRequest{}, fmt.Errorf("failed to get last %d batches, err: %w", two, err)
		}

		if len(batches) == two {
			lastClosedBatch = batches[1]
		} else {
			lastClosedBatch = batches[0]
		}

		batchNum = lastClosedBatch.BatchNumber
		oldStateRoot = f.getOldStateRootFromBatches(batches)
		txs = lastClosedBatch.BatchL2Data
		if err != nil {
			return state.ProcessRequest{}, err
		}
	} else {
		txs = make([]byte, 0, 1)
		batchNum, oldStateRoot, err = f.getLastBatchNumAndStateRoot(ctx)
		if err != nil {
			return state.ProcessRequest{}, err
		}
	}

	return state.ProcessRequest{
		BatchNumber:    batchNum,
		OldStateRoot:   oldStateRoot,
		GlobalExitRoot: f.batch.globalExitRoot,
		Coinbase:       f.sequencerAddress,
		Timestamp:      f.batch.timestamp,
		Transactions:   txs,
		Caller:         state.SequencerCallerLabel,
	}, nil
}

func (f *finalizer) getLastBatchNumAndStateRoot(ctx context.Context) (uint64, common.Hash, error) {
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
		return true
	}
	// Global Exit Root deadline
	if f.nextGERDeadline != 0 && now().Unix() >= f.nextGERDeadline {
		return true
	}
	// Delayed batch deadline
	if f.nextSendingToL1Deadline != 0 && now().Unix() >= f.nextSendingToL1Deadline {
		return true
	}

	return false
}

// checkRemainingResources checks if the transaction uses less resources than the remaining ones in the batch.
func (f *finalizer) checkRemainingResources(result *state.ProcessBatchResponse, tx *TxTracker, txResponse *state.ProcessTransactionResponse) error {
	usedResources := batchResources{
		zKCounters: result.UsedZkCounters,
		bytes:      uint64(len(tx.RawTx)),
	}
	err := f.batch.remainingResources.sub(usedResources)
	if err != nil {
		log.Infof("current transaction exceeds the batch limit, updating metadata for tx in worker and continuing")
		f.worker.UpdateTx(txResponse.TxHash, tx.From, usedResources.zKCounters)
		return err
	}

	return nil
}

// isCurrBatchAboveLimitWindow checks if the current batch is above the limit window for which is beneficial to close the batch
func (f *finalizer) isCurrBatchAboveLimitWindow() bool {
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
	f.nextForcedBatchDeadline = now().Unix() + int64(f.cfg.ForcedBatchDeadlineTimeoutInSec.Duration)
}

func (f *finalizer) setNextGERDeadline() {
	f.nextGERDeadline = now().Unix() + int64(f.cfg.GERDeadlineTimeoutInSec.Duration)
}

func (f *finalizer) setNextSendingToL1Deadline() {
	f.nextSendingToL1Deadline = now().Unix() + int64(f.cfg.SendingToL1DeadlineTimeoutInSec.Duration)
}

func (f *finalizer) getConstraintThresholdUint64(input uint64) uint64 {
	const oneHundred = 100
	return input * uint64(f.cfg.ResourcePercentageToCloseBatch) / oneHundred
}

func (f *finalizer) getConstraintThresholdUint32(input uint32) uint32 {
	const oneHundred = 100
	return uint32(input*f.cfg.ResourcePercentageToCloseBatch) / oneHundred
}
