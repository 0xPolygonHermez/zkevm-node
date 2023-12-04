package sequencer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

// Batch represents a wip or processed batch.
type Batch struct {
	batchNumber        uint64
	coinbase           common.Address
	timestamp          time.Time
	initialStateRoot   common.Hash
	stateRoot          common.Hash
	localExitRoot      common.Hash
	globalExitRoot     common.Hash // 0x000...0 (ZeroHash) means to not update
	accInputHash       common.Hash //TODO: review use
	countOfTxs         int
	remainingResources state.BatchResources
	closingReason      state.ClosingReason
}

func (w *Batch) isEmpty() bool {
	return w.countOfTxs == 0
}

// getLastStateRoot gets the state root from the latest batch
func (f *finalizer) getLastStateRoot(ctx context.Context) (common.Hash, error) {
	var oldStateRoot common.Hash

	batches, err := f.dbManager.GetLastNBatches(ctx, 2) //nolint:gomnd
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get last %d batches, err: %w", 2, err) //nolint:gomnd
	}

	if len(batches) == 1 { //nolint:gomnd
		oldStateRoot = batches[0].StateRoot
	} else if len(batches) == 2 { //nolint:gomnd
		oldStateRoot = batches[1].StateRoot
	}

	return oldStateRoot, nil
}

// getWIPBatch gets the last batch if still wip or opens a new one
func (f *finalizer) getWIPBatch(ctx context.Context) {
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	lastBatchNum, err := f.dbManager.GetLastBatchNumber(ctx)
	if err != nil {
		log.Fatalf("failed to get last batch number. Error: %s", err)
	}

	if lastBatchNum == 0 {
		// GENESIS batch
		processingCtx := f.dbManager.CreateFirstBatch(ctx, f.sequencerAddress)
		timestamp := processingCtx.Timestamp
		oldStateRoot, err := f.getLastStateRoot(ctx)
		if err != nil {
			log.Fatalf("failed to get old state root. Error: %s", err)
		}
		f.wipBatch = &Batch{
			globalExitRoot:     processingCtx.GlobalExitRoot,
			initialStateRoot:   oldStateRoot,
			stateRoot:          oldStateRoot,
			batchNumber:        processingCtx.BatchNumber,
			coinbase:           processingCtx.Coinbase,
			timestamp:          timestamp,
			remainingResources: getMaxRemainingResources(f.batchConstraints),
		}
	} else {
		// Get the last batch if is still wip, if not open a new one
		lastBatch, err := f.dbManager.GetBatchByNumber(ctx, lastBatchNum, nil)
		if err != nil {
			log.Fatalf("failed to get last batch. Error: %s", err)
		}

		isClosed, err := f.dbManager.IsBatchClosed(ctx, lastBatchNum)
		if err != nil {
			log.Fatalf("failed to check if batch is closed. Error: %s", err)
		}

		log.Infof("batch %d isClosed: %v", lastBatchNum, isClosed)

		if isClosed { //open new wip batch
			ger, _, err := f.dbManager.GetLatestGer(ctx, f.cfg.GERFinalityNumberOfBlocks)
			if err != nil {
				log.Fatalf("failed to get latest ger. Error: %s", err)
			}

			oldStateRoot := lastBatch.StateRoot
			f.wipBatch, err = f.openNewWIPBatch(ctx, lastBatchNum+1, ger.GlobalExitRoot, oldStateRoot)
			if err != nil {
				log.Fatalf("failed to open new wip batch. Error: %s", err)
			}
		} else { /// get wip batch
			f.wipBatch, err = f.dbManager.GetWIPBatch(ctx)
			if err != nil {
				log.Fatalf("failed to get wip batch. Error: %s", err)
			}
		}
	}

	log.Infof("initial batch: %d, initialStateRoot: %s, stateRoot: %s, coinbase: %s, GER: %s, LER: %s",
		f.wipBatch.batchNumber, f.wipBatch.initialStateRoot.String(), f.wipBatch.stateRoot.String(), f.wipBatch.coinbase.String(),
		f.wipBatch.globalExitRoot.String(), f.wipBatch.localExitRoot.String())
}

// finalizeBatch retries to until successful closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new empty batch
func (f *finalizer) finalizeBatch(ctx context.Context) {
	start := time.Now()
	defer func() {
		metrics.ProcessingTime(time.Since(start))
	}()

	var err error
	f.wipBatch, err = f.closeAndOpenNewWIPBatch(ctx)
	for err != nil { //TODO: we need to review is this for loop is needed or if it's better to halt if we have an error
		log.Errorf("failed to create new WIP batch. Error: %s", err)
		f.wipBatch, err = f.closeAndOpenNewWIPBatch(ctx)
	}

	log.Infof("new WIP batch %d", f.wipBatch.batchNumber)
}

// closeAndOpenNewWIPBatch closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new empty batch
func (f *finalizer) closeAndOpenNewWIPBatch(ctx context.Context) (*Batch, error) {
	// Finalize the wip L2 block if it has transactions, if not we keep it open to store it in the new wip batch
	if !f.wipL2Block.isEmpty() {
		f.finalizeL2Block(ctx)
	}

	// Wait until all L2 blocks are processed
	startWait := time.Now()
	f.pendingL2BlocksToProcessWG.Wait()
	endWait := time.Now()
	log.Debugf("waiting for pending L2 blocks to be processed took: %s", endWait.Sub(startWait).String())

	// Wait until all L2 blocks are store
	startWait = time.Now()
	f.pendingL2BlocksToStoreWG.Wait()
	endWait = time.Now()
	log.Debugf("waiting for pending L2 blocks to be stored took: %s", endWait.Sub(startWait).String())

	var err error
	if f.wipBatch.stateRoot == state.ZeroHash {
		return nil, errors.New("state root must have value to close batch")
	}

	// We need to process the batch to update the state root before closing the batch
	if f.wipBatch.initialStateRoot == f.wipBatch.stateRoot {
		log.Info("reprocessing batch because the state root has not changed...")
		_, err = f.processTransaction(ctx, nil, true)
		if err != nil {
			return nil, err
		}
	}

	// Reprocess full batch as sanity check
	//TODO: Uncomment this
	/*if f.cfg.SequentialReprocessFullBatch {
		// Do the full batch reprocess now
		_, err := f.reprocessFullBatch(ctx, f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.stateRoot)
		if err != nil {
			// There is an error reprocessing the batch. We halt the execution of the Sequencer at this point
			f.halt(ctx, fmt.Errorf("halting Sequencer because of error reprocessing full batch %d (sanity check). Error: %s ", f.wipBatch.batchNumber, err))
		}
	} else {
		// Do the full batch reprocess in parallel
		go func() {
			_, _ = f.reprocessFullBatch(ctx, f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.stateRoot)
		}()
	}*/

	// Close the wip batch
	err = f.closeWIPBatch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to close batch, err: %w", err)
	}

	log.Infof("batch %d closed", f.wipBatch.batchNumber)

	// Check if the batch is empty and sending a GER Update to the stream is needed
	if f.streamServer != nil && f.wipBatch.isEmpty() && f.currentGERHash != f.previousGERHash {
		updateGer := state.DSUpdateGER{
			BatchNumber:    f.wipBatch.batchNumber,
			Timestamp:      f.wipBatch.timestamp.Unix(),
			GlobalExitRoot: f.wipBatch.globalExitRoot,
			Coinbase:       f.sequencerAddress,
			ForkID:         uint16(f.dbManager.GetForkIDByBatchNumber(f.wipBatch.batchNumber)),
			StateRoot:      f.wipBatch.stateRoot,
		}

		err = f.streamServer.StartAtomicOp()
		if err != nil {
			log.Errorf("failed to start atomic op for Update GER on batch %v: %v", f.wipBatch.batchNumber, err)
		}

		_, err = f.streamServer.AddStreamEntry(state.EntryTypeUpdateGER, updateGer.Encode())
		if err != nil {
			log.Errorf("failed to add stream entry for Update GER on batch %v: %v", f.wipBatch.batchNumber, err)
		}

		err = f.streamServer.CommitAtomicOp()
		if err != nil {
			log.Errorf("failed to commit atomic op for Update GER on batch  %v: %v", f.wipBatch.batchNumber, err)
		}
	}

	// Metadata for the next batch
	stateRoot := f.wipBatch.stateRoot
	lastBatchNumber := f.wipBatch.batchNumber

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
		f.previousGERHash = f.currentGERHash
		f.currentGERHash = f.nextGER
	}
	f.nextGER = state.ZeroHash
	f.nextGERDeadline = 0
	f.nextGERMux.Unlock()

	batch, err := f.openNewWIPBatch(ctx, lastBatchNumber+1, f.currentGERHash, stateRoot)

	// Substract the bytes needed to store the changeL2Block tx into the new batch
	batch.remainingResources.Bytes = batch.remainingResources.Bytes - changeL2BlockSize

	return batch, err
}

// openNewWIPBatch opens a new batch in the state and returns it as WipBatch
func (f *finalizer) openNewWIPBatch(ctx context.Context, batchNum uint64, ger, stateRoot common.Hash) (*Batch, error) {
	// open next batch
	processingCtx := state.ProcessingContext{
		BatchNumber:    batchNum,
		Coinbase:       f.sequencerAddress,
		Timestamp:      now(),
		GlobalExitRoot: ger,
	}

	dbTx, err := f.dbManager.BeginStateTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin state transaction to open batch, err: %w", err)
	}

	// OpenBatch opens a new batch in the state
	err = f.dbManager.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return nil, fmt.Errorf("failed to rollback dbTx: %s. Error: %w", rollbackErr.Error(), err)
		}
		return nil, fmt.Errorf("failed to open new batch. Error: %w", err)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit database transaction for opening a batch. Error: %w", err)
	}

	// Check if synchronizer is up-to-date
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	return &Batch{
		batchNumber:        batchNum,
		coinbase:           f.sequencerAddress,
		initialStateRoot:   stateRoot,
		stateRoot:          stateRoot,
		timestamp:          processingCtx.Timestamp,
		globalExitRoot:     ger,
		remainingResources: getMaxRemainingResources(f.batchConstraints),
		closingReason:      state.EmptyClosingReason,
	}, err
}

// closeWIPBatch closes the current batch in the state
func (f *finalizer) closeWIPBatch(ctx context.Context) error {
	transactions, effectivePercentages, err := f.dbManager.GetTransactionsByBatchNumber(ctx, f.wipBatch.batchNumber)
	if err != nil {
		return fmt.Errorf("failed to get transactions from transactions, err: %w", err)
	}
	for i, tx := range transactions {
		log.Debugf("[closeWIPBatch] BatchNum: %d, Tx position: %d, txHash: %s", f.wipBatch.batchNumber, i, tx.Hash().String())
	}
	usedResources := getUsedBatchResources(f.batchConstraints, f.wipBatch.remainingResources)
	receipt := ClosingBatchParameters{
		BatchNumber:          f.wipBatch.batchNumber,
		StateRoot:            f.wipBatch.stateRoot,
		LocalExitRoot:        f.wipBatch.localExitRoot,
		Txs:                  transactions,
		EffectivePercentages: effectivePercentages,
		BatchResources:       usedResources,
		ClosingReason:        f.wipBatch.closingReason,
	}
	return f.dbManager.CloseBatch(ctx, receipt)
}
