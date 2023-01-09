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
	ForcedBatchCh        chan state.Batch
	GERCh                chan common.Hash
	L2ReorgCh            chan L2ReorgEvent
	SendingToL1TimeoutCh chan bool
	TxsToStoreCh         chan *txToStore
	WgTxsToStore         *sync.WaitGroup
	cfg                  FinalizerCfg
	maxTxsPerBatch       uint64
	isSynced             func(ctx context.Context) bool
	sequencerAddress     common.Address
	worker               workerInterface
	dbManager            dbManagerInterface
	executor             stateInterface
	batch                WipBatch
	processRequest       state.ProcessSingleTxRequest
	sharedResourcesMux   *sync.RWMutex
	// closing signals
	nextGER                 common.Hash
	nextGERDeadline         int64
	nextGERMux              *sync.RWMutex
	nextForcedBatches       []state.Batch
	nextForcedBatchDeadline int64
	nextForcedBatchMux      *sync.RWMutex
}

// txToStore represents a transaction to store.
type txToStore struct {
	txResponse               *state.ProcessTransactionResponse
	batchNumber              uint64
	previousL2BlockStateRoot common.Hash
}

// WipBatch represents a work-in-progress batch.
type WipBatch struct {
	batchNumber        uint64
	coinbase           common.Address
	accInputHash       common.Hash
	stateRoot          common.Hash
	localExitRoot      common.Hash
	timestamp          uint64
	globalExitRoot     common.Hash // 0x000...0 (ZeroHash) means to not update
	txs                []TxTracker
	remainingResources BatchResources
}

// newFinalizer returns a new instance of Finalizer.
func newFinalizer(cfg FinalizerCfg, worker workerInterface, dbManager dbManagerInterface, executor stateInterface, sequencerAddr common.Address, isSynced func(ctx context.Context) bool, maxTxsPerBatch uint64) *finalizer {
	return &finalizer{
		ForcedBatchCh:        make(chan state.Batch),
		GERCh:                make(chan common.Hash),
		L2ReorgCh:            make(chan L2ReorgEvent),
		SendingToL1TimeoutCh: make(chan bool),
		TxsToStoreCh:         make(chan *txToStore),
		WgTxsToStore:         &sync.WaitGroup{},
		cfg:                  cfg,
		isSynced:             isSynced,
		sequencerAddress:     sequencerAddr,
		worker:               worker,
		dbManager:            dbManager,
		executor:             executor,
		maxTxsPerBatch:       maxTxsPerBatch,
		batch:                WipBatch{},
		processRequest:       state.ProcessSingleTxRequest{},
		// closing signals
		nextGER:                 common.Hash{},
		nextGERDeadline:         getNextGERDeadline(cfg),
		nextGERMux:              &sync.RWMutex{},
		nextForcedBatches:       make([]state.Batch, 0),
		nextForcedBatchDeadline: getNextForcedBatchDeadline(cfg),
		nextForcedBatchMux:      &sync.RWMutex{},
	}
}

// Start starts the finalizer.
func (f *finalizer) Start(ctx context.Context, batch *WipBatch, OldStateRoot, OldAccInputHash common.Hash) {
	var (
		err error
	)

	if batch != nil {
		f.batch = *batch
	}
	f.processRequest = state.ProcessSingleTxRequest{
		BatchNumber:      f.batch.batchNumber,
		StateRoot:        f.batch.stateRoot,
		OldStateRoot:     OldStateRoot,
		GlobalExitRoot:   f.batch.globalExitRoot,
		OldAccInputHash:  OldAccInputHash,
		SequencerAddress: f.sequencerAddress,
		Timestamp:        f.batch.timestamp,
		Caller:           state.SequencerCallerLabel,
	}

	// Closing signals receiver
	go f.handleClosingSignals(ctx, err)

	// Finalize txs
	go func() {
		for {
			tx := f.worker.GetBestFittingTx(f.batch.remainingResources)
			if tx != nil {
				if success, _ := f.processTransaction(ctx, tx); !success {
					continue
				}
			} else {
				if f.isCurrBatchAboveLimitWindow() {
					f.WgTxsToStore.Wait()
					f.reopenBatch(ctx)
					// // go (decide if we need to execute the full batch as a sanity check, DO IT IN PARALLEL) ==> if error: log this txs somewhere and remove them from the pipeline
					if len(f.nextForcedBatches) > 0 {
						// TODO: implement processing of forced batches
					}
					// // open batch: check if we have a new globalExitRoot and update timestamp
				} else {
					if f.cfg.SleepDurationInMs.Duration > 0 {
						time.Sleep(f.cfg.SleepDurationInMs.Duration)
					}
				}
			}

			f.checkDeadlines(ctx)
			if f.cfg.SleepDurationInMs.Duration > 0 {
				time.Sleep(f.cfg.SleepDurationInMs.Duration * time.Millisecond)
			}
			<-ctx.Done()
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
	result, err := f.executor.ProcessSingleTransaction(ctx, f.processRequest, dbTx)

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
		usedResources := BatchResources{
			zKCounters: result.UsedZkCounters,
			bytes:      uint64(len(tx.RawTx)),
			gas:        txResponse.GasUsed,
		}
		err = f.batch.remainingResources.Sub(usedResources)
		if err != nil {
			f.worker.UpdateTx(txResponse.TxHash, tx.From, usedResources.zKCounters)
			f.checkDeadlines(ctx)
			return false, err
		}

		// We have a successful processing if we are here
		previousL2BlockStateRoot := f.batch.stateRoot
		f.batch.stateRoot = result.NewStateRoot
		f.batch.localExitRoot = result.NewLocalExitRoot
		f.batch.accInputHash = result.NewAccInputHash
		f.processRequest.StateRoot = result.NewStateRoot
		f.processRequest.OldAccInputHash = result.NewAccInputHash

		// Store the processed transaction, add it to the batch and update status in the pool atomically
		f.WgTxsToStore.Add(1)
		f.TxsToStoreCh <- &txToStore{
			batchNumber:              f.batch.batchNumber,
			txResponse:               txResponse,
			previousL2BlockStateRoot: previousL2BlockStateRoot,
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

func (f *finalizer) checkDeadlines(ctx context.Context) {

	if time.Now().Unix() >= f.nextForcedBatchDeadline {
		f.nextForcedBatchMux.Lock()
		// TODO: Check if there are any new forced batches and pass to the channel "nextForcedBatchesCh"
		f.nextForcedBatchDeadline = getNextForcedBatchDeadline(f.cfg)
		f.nextForcedBatchMux.Unlock()
	}

	// Check deadlines
	if time.Now().Unix() >= f.nextGERDeadline {
		f.nextGERMux.Lock()
		ger, _, err := f.dbManager.GetLatestGer(ctx)
		if err != nil {
			log.Errorf("failed to get latest GER. Err: %w", err)
			return
		}
		// React only if the GER has changed
		if ger.GlobalExitRoot != f.batch.globalExitRoot {
			f.GERCh <- ger.GlobalExitRoot
		}
		f.nextGERMux.Unlock()
	}
}

func (f *finalizer) handleClosingSignals(ctx context.Context, err error) {
	for {
		select {
		// Forced  batch ch
		case fb := <-f.ForcedBatchCh:
			f.sharedResourcesMux.Lock()
			f.nextForcedBatchMux.Lock()
			f.nextForcedBatches = append(f.nextForcedBatches, fb)
			// TODO: Close current batch, process forced batch and open a new one
			f.nextForcedBatchDeadline = getNextForcedBatchDeadline(f.cfg)
			f.nextForcedBatchMux.Unlock()
			f.sharedResourcesMux.Unlock()
		// globalExitRoot ch
		case ger := <-f.GERCh:
			f.sharedResourcesMux.Lock()
			f.nextGERMux.Lock()
			f.nextGER = ger
			f.nextGERDeadline = getNextGERDeadline(f.cfg)
			f.reopenBatch(ctx)
			f.nextGERMux.Unlock()
			f.sharedResourcesMux.Unlock()
		// L2Reorg ch
		case l2ReorgEvent := <-f.L2ReorgCh:
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
		f.batch, err = f.openWIPBatch(ctx, dbTx)
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

	f.processRequest = state.ProcessSingleTxRequest{
		BatchNumber:      f.batch.batchNumber,
		StateRoot:        f.batch.stateRoot,
		OldStateRoot:     oldStateRoot,
		GlobalExitRoot:   f.batch.globalExitRoot,
		OldAccInputHash:  oldAccInputHash,
		SequencerAddress: f.sequencerAddress,
		Timestamp:        f.batch.timestamp,
		Caller:           state.SequencerCallerLabel,
	}

	return nil
}

func (f *finalizer) isCurrBatchAboveLimitWindow() bool {
	// TODO: use the new config param f.cfg.ResourcePercentageToCloseBatch for each resource to do the check
	panic("Implement me")
	return false
}

func (f *finalizer) backupWIPBatch() WipBatch {
	backup := new(WipBatch)
	*backup = f.batch
	backup.txs = make([]TxTracker, 0, len(f.batch.txs))
	backup.txs = append(backup.txs, f.batch.txs...)

	return *backup
}

func (f *finalizer) newWIPBatch(ctx context.Context) (WipBatch, error) {
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
		return WipBatch{}, errors.New("state root and local exit root must have value to close batch")
	}
	dbTx, err = f.dbManager.BeginStateTransaction(ctx)
	if err != nil {
		return WipBatch{}, fmt.Errorf("failed to begin state transaction to close batch, err: %w", err)
	}
	f.closeBatch(ctx, dbTx)

	return f.openWIPBatch(ctx, dbTx)
}

func (f *finalizer) openWIPBatch(ctx context.Context, dbTx pgx.Tx) (WipBatch, error) {
	// open next batch
	gerHash, err := f.getGERHash(ctx, dbTx)
	if err != nil {
		return WipBatch{}, err
	}

	_, err = f.openBatch(ctx, gerHash, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return WipBatch{}, fmt.Errorf(
				"failed to rollback dbTx when getting last batch num that gave err: %s. Rollback err: %s",
				rollbackErr.Error(), err.Error(),
			)
		}
		return WipBatch{}, err
	}
	if err := dbTx.Commit(ctx); err != nil {
		return WipBatch{}, err
	}

	// Check if synchronizer is up-to-date
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	return f.dbManager.GetWIPBatch(ctx)
}

func (f *finalizer) getGERHash(ctx context.Context, dbTx pgx.Tx) (gerHash common.Hash, err error) {
	if f.batch.globalExitRoot != f.nextGER {
		gerHash = f.nextGER
	} else {
		ger, _, err := f.dbManager.GetLatestGer(ctx)
		if err != nil {
			if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
				return common.Hash{}, fmt.Errorf(
					"failed to rollback dbTx when getting last globalExitRoot that gave err: %s. Rollback err: %s",
					rollbackErr.Error(), err.Error())
			}
			return common.Hash{}, err
		}
		gerHash = ger.GlobalExitRoot
	}
	return gerHash, nil
}

func (f *finalizer) reopenBatch(ctx context.Context) {
	var err error
	f.batch, err = f.newWIPBatch(ctx)
	for err != nil {
		log.Errorf("failed to create new work-in-progress batch, Err: %s", err)
		f.batch, err = f.newWIPBatch(ctx)
	}
}

func (f *finalizer) closeBatch(ctx context.Context, dbTx pgx.Tx) {
	receipt := ClosingBatchParameters{
		BatchNumber:   f.batch.batchNumber,
		AccInputHash:  f.processRequest.OldAccInputHash,
		StateRoot:     f.batch.stateRoot,
		LocalExitRoot: f.processRequest.GlobalExitRoot,
		Txs:           f.batch.txs,
	}
	f.dbManager.CloseBatch(ctx, receipt, dbTx)
}

func (f *finalizer) openBatch(ctx context.Context, gerHash common.Hash, dbTx pgx.Tx) (state.ProcessingContext, error) {
	lastBatchNum, err := f.dbManager.GetLastBatchNumber(ctx)
	if err != nil {
		return state.ProcessingContext{}, fmt.Errorf("failed to get last batch number, err: %w", err)
	}
	newBatchNum := lastBatchNum + 1
	processingCtx := state.ProcessingContext{
		BatchNumber:    newBatchNum,
		Coinbase:       f.sequencerAddress,
		Timestamp:      time.Now(),
		GlobalExitRoot: gerHash,
	}
	err = f.dbManager.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		return state.ProcessingContext{}, fmt.Errorf("failed to open new batch, err: %w", err)
	}

	return processingCtx, nil
}
