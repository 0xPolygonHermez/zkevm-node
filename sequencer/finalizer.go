package sequencer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"

	"github.com/0xPolygonHermez/zkevm-node/log"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
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
	isEmpty            bool
	remainingResources batchResources
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
		cfg:              cfg,
		txsStore:         txsStore,
		closingSignalCh:  closingSignalCh,
		isSynced:         isSynced,
		sequencerAddress: sequencerAddr,
		worker:           worker,
		dbManager:        dbManager,
		executor:         executor,
		batch:            new(WipBatch),
		batchConstraints: batchConstraints,
		processRequest:   state.ProcessRequest{},
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
func (f *finalizer) Start(ctx context.Context, batch *WipBatch, OldStateRoot common.Hash) {
	if batch != nil {
		f.batch = batch
	} else {
		var err error
		f.batch, err = f.dbManager.GetWIPBatch(ctx)
		if err != nil {
			log.Fatalf("failed to get work-in-progress batch from DB, Err: %s", err)
		}
	}

	f.processRequest = state.ProcessRequest{
		BatchNumber:    f.batch.batchNumber,
		OldStateRoot:   OldStateRoot,
		GlobalExitRoot: f.batch.globalExitRoot,
		Coinbase:       f.sequencerAddress,
		Timestamp:      f.batch.timestamp,
		Caller:         state.SequencerCallerLabel,
	}

	// Closing signals receiver
	go f.listenForClosingSignals(ctx)

	// Processing transactions and finalizing batches
	f.finalizeBatches(ctx)
}

// listenForClosingSignals listens for signals for the batch and sets the deadline for when they need to be closed.
func (f *finalizer) listenForClosingSignals(ctx context.Context) {
	var err error
	for {
		select {
		case <-ctx.Done():
			log.Infof("finalizer closing signal listener received context done, Err: %s", ctx.Err())
			return
		// Forced  batch ch
		case fb := <-f.closingSignalCh.ForcedBatchCh:
			f.nextForcedBatchesMux.Lock()
			f.nextForcedBatches = append(f.nextForcedBatches, fb)
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
			go f.worker.HandleL2Reorg(l2ReorgEvent.TxHashes)
			// Get current wip batch
			f.batch, err = f.dbManager.GetWIPBatch(ctx)
			for err != nil {
				log.Errorf("failed to load batch from the state, err: %s", err)
				f.batch, err = f.dbManager.GetWIPBatch(ctx)
			}
			err = f.syncWIPBatchWithState(ctx)
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

		if f.isDeadlineEncountered() {
			f.finalizeBatch(ctx)
		}

		if err := ctx.Err(); err != nil {
			log.Infof("Stopping finalizer because of context, err: %s", err)
			return
		}
	}
}

// processTransaction processes a single transaction.
func (f *finalizer) processTransaction(ctx context.Context, tx *TxTracker) error {
	f.sharedResourcesMux.Lock()
	defer f.sharedResourcesMux.Unlock()

	var ger common.Hash
	if f.batch.isEmpty {
		ger = f.batch.globalExitRoot
	} else {
		ger = state.ZeroHash
	}

	f.processRequest.GlobalExitRoot = ger
	result, err := f.executor.ProcessBatch(ctx, f.processRequest)
	if err != nil {
		log.Errorf("failed to process transaction, err: %s", err)
		return err
	}

	if result.Error != nil {
		if result.Error == state.ErrBatchAlreadyClosed || result.Error == state.ErrInvalidBatchNumber {
			log.Warnf("unexpected state local vs DB: %s", result.Error)
			log.Info("reloading local sequence")
			f.batch, err = f.dbManager.GetWIPBatch(ctx)
			if err != nil {
				log.Errorf("failed to get WIP Batch from state, err: %s", err)
			}
			return err
		}
		return fmt.Errorf("failed processing transaction, err: %w", result.Error)
	} else {
		err = f.handleSuccessfulTxProcessResp(tx, result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *finalizer) handleSuccessfulTxProcessResp(tx *TxTracker, result *state.ProcessBatchResponse) error {
	if tx == nil {
		return nil
	}

	txResponse := result.Responses[0]
	// Handle Transaction Error
	if txResponse.Error != nil {
		f.handleTransactionError(txResponse, result, tx)

		return txResponse.Error
	}

	// Check remaining resources
	err := f.checkRemainingResources(result, tx, txResponse)
	if err != nil {
		return err
	}

	// We have a successful processing if we are here, updating metadata
	f.processRequest.OldStateRoot = f.batch.stateRoot
	f.batch.stateRoot = result.NewStateRoot
	f.batch.localExitRoot = result.NewLocalExitRoot
	f.processRequest.OldAccInputHash = result.NewAccInputHash

	// Store the processed transaction, add it to the batch and update status in the pool atomically
	f.txsStore.Wg.Add(1)
	f.txsStore.Ch <- &txToStore{
		batchNumber:              f.batch.batchNumber,
		txResponse:               txResponse,
		previousL2BlockStateRoot: f.batch.stateRoot,
	}
	f.worker.UpdateAfterSingleSuccessfulTxExecution(tx.From, result.TouchedAddresses)
	f.batch.isEmpty = false

	return nil
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

// finalizeBatch closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new empty batch
func (f *finalizer) finalizeBatch(ctx context.Context) {
	var err error
	f.batch, err = f.newWIPBatch(ctx)
	for err != nil {
		log.Errorf("failed to create new work-in-progress batch, Err: %s", err)
		f.batch, err = f.newWIPBatch(ctx)
	}
}

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

func (f *finalizer) handleTransactionError(txResponse *state.ProcessTransactionResponse, result *state.ProcessBatchResponse, tx *TxTracker) {
	errorCode := executor.ErrorCode(txResponse.Error)
	addressInfo := result.TouchedAddresses[tx.From]
	if executor.IsOutOfCountersError(errorCode) {
		f.worker.DeleteTx(tx.Hash, tx.From, addressInfo.Nonce, addressInfo.Balance)
	} else if executor.IsIntrinsicError(errorCode) {
		f.worker.MoveTxToNotReady(tx.Hash, tx.From, addressInfo.Nonce, addressInfo.Balance)
	}
}

func (f *finalizer) syncWIPBatchWithState(ctx context.Context) error {
	var err error
	// Check if synchronizer is up-to-date
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	// Get data for prevBatch
	f.processRequest, err = f.prepareProcessRequestFromState(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (f *finalizer) processForcedBatches(batchNumber uint64, stateRoot common.Hash) (uint64, common.Hash) {
	f.nextForcedBatchesMux.Lock()
	for _, forcedBatch := range f.nextForcedBatches {
		batchNumber += 1
		processRequest := state.ProcessRequest{
			BatchNumber:    batchNumber,
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
			log.Errorf("failed to process forced batch, err: %s", err)
		} else {
			stateRoot = response.NewStateRoot
		}
	}
	f.nextForcedBatches = make([]state.ForcedBatch, 0)
	f.nextForcedBatchesMux.Unlock()
	return batchNumber, stateRoot
}

func (f *finalizer) newWIPBatch(ctx context.Context) (*WipBatch, error) {
	var (
		dbTx pgx.Tx
		err  error
	)

	// Passing the batch without txs to the executor in order to update the State
	if f.batch.isEmpty {
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
	dbTx, err = f.dbManager.BeginStateTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin state transaction to close batch, err: %w", err)
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
		lastBatchNumber, stateRoot = f.processForcedBatches(lastBatchNumber, stateRoot)
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

	return f.openWIPBatch(ctx, lastBatchNumber+1, ger, stateRoot, dbTx)
}

func (f *finalizer) openWIPBatch(ctx context.Context, batchNum uint64, ger, stateRoot common.Hash, dbTx pgx.Tx) (*WipBatch, error) {
	if dbTx == nil {
		var err error
		dbTx, err = f.dbManager.BeginStateTransaction(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to begin state transaction to open batch, err: %w", err)
		}
	}

	// open next batch
	openBatchResp, err := f.openBatch(ctx, batchNum, ger, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return nil, fmt.Errorf(
				"failed to rollback dbTx when getting last batch num that gave err: %s. Rollback err: %s",
				rollbackErr.Error(), err.Error(),
			)
		}
		return nil, err
	}
	if err := dbTx.Commit(ctx); err != nil {
		return nil, err
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

func (f *finalizer) closeBatch(ctx context.Context) error {
	receipt := ClosingBatchParameters{
		BatchNumber:   f.batch.batchNumber,
		AccInputHash:  f.processRequest.OldAccInputHash,
		StateRoot:     f.batch.stateRoot,
		LocalExitRoot: f.processRequest.GlobalExitRoot,
	}
	return f.dbManager.CloseBatch(ctx, receipt)
}

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
func (f *finalizer) reprocessBatch(ctx context.Context) error {
	processRequest, err := f.prepareProcessRequestFromState(ctx)
	if err != nil {
		log.Errorf("failed to prepare process request for reprocessing batch, err: %s", err)
		return err
	}
	result, err := f.executor.ProcessBatch(ctx, processRequest)
	if err != nil || result.IsBatchProcessed == false || result.Error != nil {
		if result.Error != nil {
			err = result.Error
		}
		log.Errorf("failed to reprocess batch, err: %s", err)
		return err
	}

	return nil
}

func (f *finalizer) prepareProcessRequestFromState(ctx context.Context) (state.ProcessRequest, error) {
	var (
		oldStateRoot, oldAccInputHash common.Hash
	)

	n := uint(2)
	batches, err := f.dbManager.GetLastNBatches(ctx, n)
	lastBatch := batches[0]
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get last %d batches, err: %w", n, err))
	}

	if len(batches) == 1 {
		oldAccInputHash = lastBatch.AccInputHash
		oldStateRoot = lastBatch.StateRoot
	} else if len(batches) == 2 {
		oldAccInputHash = batches[1].AccInputHash
		oldStateRoot = batches[1].StateRoot
	}

	return state.ProcessRequest{
		BatchNumber:     f.batch.batchNumber,
		OldStateRoot:    oldStateRoot,
		GlobalExitRoot:  f.batch.globalExitRoot,
		OldAccInputHash: oldAccInputHash,
		Coinbase:        f.sequencerAddress,
		Timestamp:       f.batch.timestamp,
		Caller:          state.SequencerCallerLabel,
	}, nil
}

// isCurrBatchAboveLimitWindow checks if the current batch is above the limit window for which is beneficial to close the batch
func (f *finalizer) isCurrBatchAboveLimitWindow() bool {
	resources := f.batch.remainingResources
	zkCounters := resources.zKCounters
	if resources.bytes >= f.getConstraintThresholdUint64(f.batchConstraints.MaxBatchBytesSize) {
		return true
	}
	if zkCounters.UsedSteps >= f.getConstraintThresholdUint32(f.batchConstraints.MaxSteps) {
		return true
	}
	if zkCounters.UsedPoseidonPaddings >= f.getConstraintThresholdUint32(f.batchConstraints.MaxPoseidonPaddings) {
		return true
	}
	if zkCounters.UsedBinaries >= f.getConstraintThresholdUint32(f.batchConstraints.MaxBinaries) {
		return true
	}
	if zkCounters.UsedKeccakHashes >= f.getConstraintThresholdUint32(f.batchConstraints.MaxKeccakHashes) {
		return true
	}
	if zkCounters.UsedArithmetics >= f.getConstraintThresholdUint32(f.batchConstraints.MaxArithmetics) {
		return true
	}
	if zkCounters.UsedMemAligns >= f.getConstraintThresholdUint32(f.batchConstraints.MaxMemAligns) {
		return true
	}
	if zkCounters.CumulativeGasUsed >= f.getConstraintThresholdUint64(f.batchConstraints.MaxCumulativeGasUsed) {
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
	return input * uint64(f.cfg.ResourcePercentageToCloseBatch) / 100
}

func (f *finalizer) getConstraintThresholdUint32(input uint32) uint32 {
	return uint32(input*f.cfg.ResourcePercentageToCloseBatch) / 100
}
