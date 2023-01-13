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
	nextGER                 common.Hash
	nextGERDeadline         int64
	nextGERMux              *sync.RWMutex
	nextForcedBatches       []state.ForcedBatch
	nextForcedBatchDeadline int64
	nextForcedBatchesMux    *sync.RWMutex
}

// WipBatch represents a work-in-progress batch.
type WipBatch struct {
	batchNumber         uint64
	coinbase            common.Address
	initialAccInputHash common.Hash
	accInputHash        common.Hash
	initialStateRoot    common.Hash
	stateRoot           common.Hash
	localExitRoot       common.Hash
	timestamp           uint64
	globalExitRoot      common.Hash // 0x000...0 (ZeroHash) means to not update
	txs                 []TxTracker
	remainingResources  batchResources
}

type batchConstraints struct {
	MaxTxsPerBatch       uint64
	MaxBatchBytesSize    uint64
	MaxCumulativeGasUsed uint64
	MaxKeccakHashes      uint32
	MaxPoseidonHashes    uint32
	MaxPoseidonPaddings  uint32
	MaxMemAligns         uint32
	MaxArithmetics       uint32
	MaxBinaries          uint32
	MaxSteps             uint32
}

// TODO: Add tests to config_test.go
type batchResourceWeights struct {
	WeightBatchBytesSize    int
	WeightCumulativeGasUsed int
	WeightKeccakHashes      int
	WeightPoseidonHashes    int
	WeightPoseidonPaddings  int
	WeightMemAligns         int
	WeightArithmetics       int
	WeightBinaries          int
	WeightSteps             int
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
		batch:            &WipBatch{},
		batchConstraints: batchConstraints,
		processRequest:   state.ProcessRequest{},
		// closing signals
		nextGER:                 common.Hash{},
		nextGERDeadline:         0,
		nextGERMux:              &sync.RWMutex{},
		nextForcedBatches:       make([]state.ForcedBatch, 0),
		nextForcedBatchDeadline: 0,
		nextForcedBatchesMux:    &sync.RWMutex{},
	}
}

// Start starts the finalizer.
func (f *finalizer) Start(ctx context.Context, batch *WipBatch, OldStateRoot, OldAccInputHash common.Hash) {
	var (
		err error
	)

	if batch != nil {
		f.batch = batch
	}
	f.processRequest = state.ProcessRequest{
		BatchNumber:      f.batch.batchNumber,
		OldStateRoot:     OldStateRoot,
		GlobalExitRoot:   f.batch.globalExitRoot,
		OldAccInputHash:  OldAccInputHash,
		SequencerAddress: f.sequencerAddress,
		Timestamp:        f.batch.timestamp,
		Caller:           state.SequencerCallerLabel,
	}

	// Closing signals receiver
	go f.listenForClosingSignals(ctx, err)

	// Finalize txs
	go func() {
		for {
			tx := f.worker.GetBestFittingTx(f.batch.remainingResources)
			if tx != nil {
				_, _ = f.processTransaction(ctx, tx)
			} else {
				if f.isCurrBatchAboveLimitWindow() {
					f.txsStore.Wg.Wait()
					f.closeAndOpenNewBatch(ctx)
					// // go (decide if we need to execute the full batch as a sanity check, DO IT IN PARALLEL) ==> if error: log this txs somewhere and remove them from the pipeline
				} else {
					// wait for new txs
					if f.cfg.SleepDurationInMs.Duration > 0 {
						time.Sleep(f.cfg.SleepDurationInMs.Duration)
					}
				}
			}

			if f.isDeadlineEncountered() {
				f.closeAndOpenNewBatch(ctx)
			}
		}
	}()
}

// processTransaction processes a single transaction.
func (f *finalizer) processTransaction(ctx context.Context, tx *TxTracker) (successful bool, err error) {
	f.sharedResourcesMux.Lock()
	defer f.sharedResourcesMux.Unlock()

	var ger common.Hash
	if len(f.batch.txs) == 0 {
		ger = f.batch.globalExitRoot
	} else {
		ger = state.ZeroHash
	}

	f.processRequest.GlobalExitRoot = ger
	dbTx, err := f.dbManager.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to create dbTx. Err: %w", err)
	}
	result, err := f.executor.ProcessSingleTransaction(ctx, f.processRequest)

	if result.Error != nil {
		if err == state.ErrBatchAlreadyClosed || err == state.ErrInvalidBatchNumber {
			log.Warnf("unexpected state local vs DB: %w", err)
			log.Info("reloading local sequence")
			f.batch, err = f.dbManager.GetWIPBatch(ctx)
			if err != nil {
				log.Errorf("failed to get WIP Batch from state. Err: %w", err)
			}
		}
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when processing tx that gave err: %w. Rollback err: %v",
				rollbackErr, err,
			)
			return false, rollbackErr
		}
		log.Errorf("failed processing batch, err: %w", err)
	} else {
		if tx == nil {
			return true, nil
		}

		txResponse := result.Responses[0]
		// Handle Transaction Error
		if txResponse.Error != nil {
			errorCode := executor.ErrorCode(txResponse.Error)
			addressInfo := result.TouchedAddresses[tx.From]
			if executor.IsOutOfCountersError(errorCode) {
				f.worker.DeleteTx(tx.Hash, tx.From, addressInfo.Nonce, addressInfo.Balance)
			} else if executor.IsIntrinsicError(errorCode) {
				f.worker.MoveTxToNotReady(tx.Hash, tx.From, addressInfo.Nonce, addressInfo.Balance)
			}

			return false, txResponse.Error
		}

		// Check remaining resources
		usedResources := batchResources{
			zKCounters: result.UsedZkCounters,
			bytes:      uint64(len(tx.RawTx)),
		}
		err = f.batch.remainingResources.sub(usedResources)
		if err != nil {
			f.worker.UpdateTx(txResponse.TxHash, tx.From, usedResources.zKCounters)
			return false, err
		}

		// We have a successful processing if we are here
		f.processRequest.OldStateRoot = f.batch.stateRoot
		f.processRequest.OldAccInputHash = f.batch.accInputHash
		f.batch.stateRoot = result.NewStateRoot
		f.batch.localExitRoot = result.NewLocalExitRoot
		f.batch.accInputHash = result.NewAccInputHash
		f.processRequest.OldAccInputHash = result.NewAccInputHash

		// Store the processed transaction, add it to the batch and update status in the pool atomically
		f.txsStore.Wg.Add(1)
		f.txsStore.Ch <- &txToStore{
			batchNumber:              f.batch.batchNumber,
			txResponse:               txResponse,
			previousL2BlockStateRoot: f.batch.stateRoot,
		}
		f.worker.UpdateAfterSingleSuccessfulTxExecution(tx.From, result.TouchedAddresses)
	}
	return true, nil
}

func getNextForcedBatchDeadline(cfg FinalizerCfg) int64 {
	return time.Now().Unix() + int64(cfg.NextForcedBatchDeadlineTimeoutInSec.Duration)
}

func getNextGERDeadline(cfg FinalizerCfg) int64 {
	return time.Now().Unix() + int64(cfg.NextGERDeadlineTimeoutInSec.Duration)
}

func (f *finalizer) isDeadlineEncountered() bool {

	// Forced batch deadline
	if time.Now().Unix() >= f.nextForcedBatchDeadline && f.nextForcedBatchDeadline != 0 {
		return true
	}

	// Global Exit Root deadline
	if time.Now().Unix() >= f.nextGERDeadline && f.nextGERDeadline != 0 {
		return true
	}

	return false
}

func (f *finalizer) listenForClosingSignals(ctx context.Context, err error) {
	for {
		select {
		// Forced  batch ch
		case fb := <-f.closingSignalCh.ForcedBatchCh:
			f.nextForcedBatchesMux.Lock()
			f.nextForcedBatches = append(f.nextForcedBatches, fb)
			if f.nextForcedBatchDeadline == 0 {
				f.nextForcedBatchDeadline = getNextForcedBatchDeadline(f.cfg)
			}
			f.nextForcedBatchesMux.Unlock()
		// globalExitRoot ch
		case ger := <-f.closingSignalCh.GERCh:
			f.nextGERMux.Lock()
			f.nextGER = ger
			if f.nextGERDeadline == 0 {
				f.nextGERDeadline = getNextGERDeadline(f.cfg)
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
		}
		// TODO: More closing signals
		// Too many time without batches in L1 ch
		// Any other externality from the point of view of the sequencer should be captured using this pattern
	}
}

func (f *finalizer) syncWIPBatchWithState(ctx context.Context) error {
	var (
		oldAccInputHash, oldStateRoot common.Hash
	)

	// Check if synchronizer is up-to-date
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	// Get data for prevBatch
	lastBatch, err := f.dbManager.GetLastBatch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last batch. err: %w", err)
	}
	isClosed, err := f.dbManager.IsBatchClosed(ctx, lastBatch.BatchNumber)
	if err != nil {
		return fmt.Errorf("failed to check is batch closed or not, err: %w", err)
	}
	if isClosed {
		dbTx, err := f.dbManager.BeginStateTransaction(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin state transaction to close batch, err: %w", err)
		}
		f.batch, err = f.openWIPBatch(ctx, 0, lastBatch.StateRoot, lastBatch.GlobalExitRoot, lastBatch.AccInputHash, dbTx)
		if err != nil {
			return fmt.Errorf("failed to recreate WIP batch from state. err: %w", err)
		}
	} else {
		if lastBatch.BatchNumber == 1 {
			oldAccInputHash = lastBatch.AccInputHash
			oldStateRoot = lastBatch.StateRoot
		} else {
			n := uint(2)
			batches, err := f.dbManager.GetLastNBatches(ctx, n)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to get last %d batches, err: %w", n, err))
			}
			oldAccInputHash = batches[1].AccInputHash
			oldStateRoot = batches[1].StateRoot
		}
	}

	f.processRequest = state.ProcessRequest{
		BatchNumber:      f.batch.batchNumber,
		OldStateRoot:     oldStateRoot,
		GlobalExitRoot:   f.batch.globalExitRoot,
		OldAccInputHash:  oldAccInputHash,
		SequencerAddress: f.sequencerAddress,
		Timestamp:        f.batch.timestamp,
		Caller:           state.SequencerCallerLabel,
	}

	return nil
}

func (f *finalizer) backupWIPBatch() *WipBatch {
	backup := &WipBatch{
		batchNumber:        f.batch.batchNumber,
		coinbase:           f.batch.coinbase,
		accInputHash:       f.batch.accInputHash,
		stateRoot:          f.batch.stateRoot,
		localExitRoot:      f.batch.localExitRoot,
		timestamp:          f.batch.timestamp,
		globalExitRoot:     f.batch.localExitRoot,
		remainingResources: f.batch.remainingResources,
	}
	backup.txs = make([]TxTracker, 0, len(f.batch.txs))
	backup.txs = append(backup.txs, f.batch.txs...)

	return backup
}

func (f *finalizer) newWIPBatch(ctx context.Context) (*WipBatch, error) {
	var (
		dbTx pgx.Tx
		err  error
	)

	// It is necessary to pass the batch without txs to the executor in order to update the State
	if len(f.batch.txs) == 0 {
		// backup current sequence
		_, err = f.processTransaction(ctx, nil)
		for err != nil {
			log.Errorf("failed to process tx, err: %w", err)
			_, err = f.processTransaction(ctx, nil)
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

	stateRoot := f.batch.stateRoot
	batchNumber := f.batch.batchNumber
	accInputHash := f.batch.accInputHash

	// Process Forced Batches
	if len(f.nextForcedBatches) > 0 {
		f.nextForcedBatchesMux.Lock()
		for _, forcedBatch := range f.nextForcedBatches {
			batchNumber += 1
			processRequest := state.ProcessRequest{
				BatchNumber:      batchNumber,
				OldStateRoot:     stateRoot,
				OldAccInputHash:  accInputHash,
				GlobalExitRoot:   forcedBatch.GlobalExitRoot,
				Transactions:     forcedBatch.RawTxsData,
				SequencerAddress: f.sequencerAddress,
				Timestamp:        uint64(time.Now().Unix()),
				Caller:           state.SequencerCallerLabel,
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
	}

	// Update metadata
	f.nextGERMux.Lock()
	ger := f.nextGER
	f.nextGER = state.ZeroHash
	f.nextGERDeadline = 0
	f.nextGERMux.Unlock()

	return f.openWIPBatch(ctx, batchNumber, ger, stateRoot, accInputHash, dbTx)
}

func (f *finalizer) openWIPBatch(ctx context.Context, batchNum uint64, ger, accInputHash, stateRoot common.Hash, dbTx pgx.Tx) (*WipBatch, error) {
	// open next batch
	openBatchResp, err := f.openBatch(ctx, ger, dbTx)
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
		batchNumber:         batchNum,
		coinbase:            f.sequencerAddress,
		initialAccInputHash: accInputHash,
		accInputHash:        accInputHash,
		initialStateRoot:    stateRoot,
		stateRoot:           stateRoot,
		timestamp:           uint64(openBatchResp.Timestamp.Unix()),
		globalExitRoot:      ger,
		txs:                 make([]TxTracker, 0, f.batchConstraints.MaxTxsPerBatch),
		remainingResources:  f.getMaxRemainingResources(),
	}, err
}

func (f *finalizer) getMaxRemainingResources() batchResources {
	return batchResources{
		zKCounters: state.ZKCounters{
			CumulativeGasUsed:    f.batchConstraints.MaxCumulativeGasUsed,
			UsedKeccakHashes:     f.batchConstraints.MaxKeccakHashes,
			UsedPoseidonHashes:   f.batchConstraints.MaxPoseidonHashes,
			UsedPoseidonPaddings: f.batchConstraints.MaxPoseidonPaddings,
			UsedMemAligns:        f.batchConstraints.MaxMemAligns,
			UsedArithmetics:      f.batchConstraints.MaxArithmetics,
			UsedBinaries:         f.batchConstraints.MaxBinaries,
			UsedSteps:            f.batchConstraints.MaxSteps,
		},
		bytes: f.batchConstraints.MaxBatchBytesSize,
	}
}

// closeAndOpenNewBatch closes the current batch and opens a new one
func (f *finalizer) closeAndOpenNewBatch(ctx context.Context) {
	var err error
	f.batch, err = f.newWIPBatch(ctx)
	for err != nil {
		log.Errorf("failed to create new work-in-progress batch, Err: %s", err)
		f.batch, err = f.newWIPBatch(ctx)
	}
}

func (f *finalizer) closeBatch(ctx context.Context) error {
	receipt := ClosingBatchParameters{
		BatchNumber:   f.batch.batchNumber,
		AccInputHash:  f.processRequest.OldAccInputHash,
		StateRoot:     f.batch.stateRoot,
		LocalExitRoot: f.processRequest.GlobalExitRoot,
		Txs:           f.batch.txs,
	}
	return f.dbManager.CloseBatch(ctx, receipt)
}

func (f *finalizer) openBatch(ctx context.Context, ger common.Hash, dbTx pgx.Tx) (state.ProcessingContext, error) {
	lastBatchNum, err := f.dbManager.GetLastBatchNumber(ctx)
	if err != nil {
		return state.ProcessingContext{}, fmt.Errorf("failed to get last batch number, err: %w", err)
	}
	newBatchNum := lastBatchNum + 1
	processingCtx := state.ProcessingContext{
		BatchNumber:    newBatchNum,
		Coinbase:       f.sequencerAddress,
		Timestamp:      time.Now(),
		GlobalExitRoot: ger,
	}
	err = f.dbManager.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		return state.ProcessingContext{}, fmt.Errorf("failed to open new batch, err: %w", err)
	}

	return processingCtx, nil
}

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

func (f *finalizer) getConstraintThresholdUint64(input uint64) uint64 {
	return input * uint64(f.cfg.ResourcePercentageToCloseBatch) / 100
}

func (f *finalizer) getConstraintThresholdUint32(input uint32) uint32 {
	return uint32(input*f.cfg.ResourcePercentageToCloseBatch) / 100
}
