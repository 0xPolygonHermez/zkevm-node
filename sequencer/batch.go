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
	batchNumber         uint64
	coinbase            common.Address
	timestamp           time.Time
	initialStateRoot    common.Hash // initial stateRoot of the batch
	initialAccInputHash common.Hash // initial accInputHash of the batch
	imStateRoot         common.Hash // intermediate stateRoot that is updated each time a single tx is processed
	imAccInputHash      common.Hash // intermediate accInputHash that is updated each time a single tx is processed
	finalStateRoot      common.Hash // final stateroot of the batch when a L2 block is processed
	finalAccInputHash   common.Hash // final accInputHash of the batch when a L2 block is processed
	localExitRoot       common.Hash
	globalExitRoot      common.Hash // 0x000...0 (ZeroHash) means to not update
	countOfTxs          int
	remainingResources  state.BatchResources
	closingReason       state.ClosingReason
}

func (w *Batch) isEmpty() bool {
	return w.countOfTxs == 0
}

// getLastStateRoot gets the state root from the latest batch
func (f *finalizer) getLastStateRoot(ctx context.Context) (common.Hash, error) {
	var oldStateRoot common.Hash

	batches, err := f.state.GetLastNBatches(ctx, 2, nil) //nolint:gomnd
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

// createFirstBatch is using during genesis
func (f *finalizer) createFirstBatch(ctx context.Context, sequencerAddress common.Address) state.ProcessingContext {
	processingCtx := state.ProcessingContext{
		BatchNumber:    1,
		Coinbase:       sequencerAddress,
		Timestamp:      time.Now(),
		GlobalExitRoot: state.ZeroHash,
	}
	dbTx, err := f.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for opening a batch, err: %v", err)
		return processingCtx
	}
	err = f.state.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when opening batch that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		log.Errorf("failed to open a batch, err: %v", err)
		return processingCtx
	}
	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx when opening batch, err: %v", err)
		return processingCtx
	}
	return processingCtx
}

// GetWIPBatch returns ready WIP batch
func (f *finalizer) setWIPBatch(ctx context.Context, wipStateBatch *state.Batch) (*Batch, error) {
	dbTx, err := f.state.BeginStateTransaction(ctx)
	if err != nil {
		return nil, err
	}

	//TODO: Check wipStateBatch.BatchNumber > 0
	//TODO: We retrieve prevStateBatch only to init the initialStateRoot of the wip batch. To avoid this We can update the initialStateRoot in the state.batch table for the wip batch
	prevStateBatch, err := f.state.GetBatchByNumber(ctx, wipStateBatch.BatchNumber-1, dbTx)
	if err != nil {
		return nil, err
	}

	wipStateBatchBlocks, err := state.DecodeBatchV2(wipStateBatch.BatchL2Data)
	if err != nil {
		return nil, err
	}

	// Count the number of txs in the wip state batch
	wipStateBatchCountOfTxs := 0
	for _, rawBlock := range wipStateBatchBlocks.Blocks {
		wipStateBatchCountOfTxs = wipStateBatchCountOfTxs + len(rawBlock.Transactions)
	}

	wipBatch := &Batch{
		batchNumber:         wipStateBatch.BatchNumber,
		coinbase:            wipStateBatch.Coinbase,
		imStateRoot:         wipStateBatch.StateRoot,
		initialStateRoot:    prevStateBatch.StateRoot,
		finalStateRoot:      wipStateBatch.StateRoot,
		initialAccInputHash: wipStateBatch.AccInputHash,
		imAccInputHash:      wipStateBatch.AccInputHash,
		finalAccInputHash:   wipStateBatch.AccInputHash,
		localExitRoot:       wipStateBatch.LocalExitRoot,
		timestamp:           wipStateBatch.Timestamp,
		globalExitRoot:      wipStateBatch.GlobalExitRoot,
		countOfTxs:          wipStateBatchCountOfTxs,
	}

	// Init counters to MAX values
	var totalBytes uint64 = f.batchConstraints.MaxBatchBytesSize
	var batchZkCounters = state.ZKCounters{
		GasUsed:              f.batchConstraints.MaxCumulativeGasUsed,
		UsedKeccakHashes:     f.batchConstraints.MaxKeccakHashes,
		UsedPoseidonHashes:   f.batchConstraints.MaxPoseidonHashes,
		UsedPoseidonPaddings: f.batchConstraints.MaxPoseidonPaddings,
		UsedMemAligns:        f.batchConstraints.MaxMemAligns,
		UsedArithmetics:      f.batchConstraints.MaxArithmetics,
		UsedBinaries:         f.batchConstraints.MaxBinaries,
		UsedSteps:            f.batchConstraints.MaxSteps,
		UsedSha256Hashes_V2:  f.batchConstraints.MaxSHA256Hashes,
	}

	//TODO: We execute the batch to get the used counter. To avoid this We can update the counters in the state.batch table for the wip batch
	if wipStateBatchCountOfTxs > 0 {
		//TODO: Change wipStateBatch.GlobalExitRoot for L1InfoRoot and wipStateBatch.Timestamp for the TimeLimit
		batchResponse, err := f.state.ExecuteBatchV2(ctx, *wipStateBatch, wipStateBatch.GlobalExitRoot, wipStateBatch.Timestamp, false, dbTx)
		if err != nil {
			return nil, err
		}

		zkCounters := &state.ZKCounters{
			GasUsed:              batchResponse.GasUsed, //TODO: review this is ok
			UsedKeccakHashes:     batchResponse.CntKeccakHashes,
			UsedPoseidonHashes:   batchResponse.CntPoseidonHashes,
			UsedPoseidonPaddings: batchResponse.CntPoseidonPaddings,
			UsedMemAligns:        batchResponse.CntMemAligns,
			UsedArithmetics:      batchResponse.CntArithmetics,
			UsedBinaries:         batchResponse.CntBinaries,
			UsedSteps:            batchResponse.CntSteps,
			UsedSha256Hashes_V2:  batchResponse.CntSha256Hashes,
		}

		err = batchZkCounters.Sub(*zkCounters)
		if err != nil {
			return nil, err
		}

		totalBytes -= uint64(len(wipStateBatch.BatchL2Data))
	}

	wipBatch.remainingResources = state.BatchResources{ZKCounters: batchZkCounters, Bytes: totalBytes}

	return wipBatch, nil
}

// initWIPBatch inits the wip batch
func (f *finalizer) initWIPBatch(ctx context.Context) {
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	lastBatchNum, err := f.state.GetLastBatchNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get last batch number. Error: %s", err)
	}

	if lastBatchNum == 0 {
		// GENESIS batch
		processingCtx := f.createFirstBatch(ctx, f.sequencerAddress)
		timestamp := processingCtx.Timestamp
		oldStateRoot, err := f.getLastStateRoot(ctx)
		if err != nil {
			log.Fatalf("failed to get old state root. Error: %s", err)
		}
		f.wipBatch = &Batch{
			globalExitRoot:   processingCtx.GlobalExitRoot,
			initialStateRoot: oldStateRoot,
			imStateRoot:      oldStateRoot,
			finalStateRoot:   oldStateRoot,
			//TODO: review init AccInputHash
			//initialAccInputHash: ,
			//imAccInputHash: ,
			//finalAccInputHash: ,
			batchNumber:        processingCtx.BatchNumber,
			coinbase:           processingCtx.Coinbase,
			timestamp:          timestamp,
			remainingResources: getMaxRemainingResources(f.batchConstraints),
		}
	} else {
		// Get the last batch in trusted state
		lastBatch, err := f.state.GetBatchByNumber(ctx, lastBatchNum, nil)
		if err != nil {
			log.Fatalf("failed to get last batch. Error: %s", err)
		}

		isClosed, err := f.state.IsBatchClosed(ctx, lastBatchNum, nil)
		if err != nil {
			log.Fatalf("failed to check if last batch is closed. Error: %s", err)
		}

		log.Infof("batch %d isClosed: %v", lastBatchNum, isClosed)

		if isClosed { //if the last batch is close then open a new wip batch
			ger, _, err := f.state.GetLatestGer(ctx, f.cfg.GERFinalityNumberOfBlocks)
			if err != nil {
				log.Fatalf("failed to get latest GER. Error: %s", err)
			}

			f.wipBatch, err = f.openNewWIPBatch(ctx, lastBatch.BatchNumber+1, ger.GlobalExitRoot, lastBatch.StateRoot, lastBatch.LocalExitRoot, lastBatch.AccInputHash)
			if err != nil {
				log.Fatalf("failed to open new wip batch. Error: %s", err)
			}
		} else { /// if it's not closed, it is the wip state batch, set it as wip batch in the finalizer
			f.wipBatch, err = f.setWIPBatch(ctx, lastBatch)
			if err != nil {
				log.Fatalf("failed to set wip batch. Error: %s", err)
			}
		}
	}

	log.Infof("initial batch: %d, initialStateRoot: %s, stateRoot: %s, coinbase: %s, GER: %s, LER: %s",
		f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.finalStateRoot, f.wipBatch.coinbase,
		f.wipBatch.globalExitRoot, f.wipBatch.localExitRoot)
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

	//TODO: we can use only 1 WG for L2 blocks pending to process/store instead to use 2
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
	if f.wipBatch.finalStateRoot == state.ZeroHash { //TODO: this check makes sense?
		return nil, errors.New("state root must have value to close batch")
	}

	// We need to process the batch to update the state root before closing the batch
	if f.wipBatch.initialStateRoot == f.wipBatch.finalStateRoot {
		log.Info("reprocessing batch because the state root has not changed...")
		_, err = f.processTransaction(ctx, nil, true)
		if err != nil {
			return nil, err
		}
	}

	// Reprocess full batch as sanity check
	if f.cfg.SequentialReprocessFullBatch {
		// Do the full batch reprocess now
		_, err := f.reprocessFullBatch(ctx, f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.initialAccInputHash, f.wipBatch.finalStateRoot)
		if err != nil {
			// There is an error reprocessing the batch. We halt the execution of the Sequencer at this point
			f.halt(ctx, fmt.Errorf("halting Sequencer because of error reprocessing full batch %d (sanity check). Error: %s ", f.wipBatch.batchNumber, err))
		}
	} else {
		// Do the full batch reprocess in parallel
		go func() {
			_, _ = f.reprocessFullBatch(ctx, f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.initialAccInputHash, f.wipBatch.finalStateRoot)
		}()
	}

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
			ForkID:         uint16(f.state.GetForkIDByBatchNumber(f.wipBatch.batchNumber)),
			StateRoot:      f.wipBatch.finalStateRoot,
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
	stateRoot := f.wipBatch.finalStateRoot
	accInputHash := f.wipBatch.finalAccInputHash
	lastBatchNumber := f.wipBatch.batchNumber

	// Process Forced Batches
	if len(f.nextForcedBatches) > 0 {
		//TODO: processForcedBatches must return the new accInputHash
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

	batch, err := f.openNewWIPBatch(ctx, lastBatchNumber+1, f.currentGERHash, stateRoot, f.wipBatch.localExitRoot, accInputHash)
	if err != nil {
		f.halt(ctx, fmt.Errorf("failed to open new wip batch. Error: %s", err))
		return nil, err
	}

	// Subtract the bytes needed to store the changeL2Block tx into the new batch
	batch.remainingResources.Bytes = batch.remainingResources.Bytes - changeL2BlockSize

	return batch, err
}

// openNewWIPBatch opens a new batch in the state and returns it as WipBatch
func (f *finalizer) openNewWIPBatch(ctx context.Context, batchNumber uint64, ger, stateRoot, ler, accInputHash common.Hash) (*Batch, error) {
	// open next batch
	newStateBatch := state.Batch{
		BatchNumber:    batchNumber,
		Coinbase:       f.sequencerAddress,
		Timestamp:      now(),
		GlobalExitRoot: ger,
		StateRoot:      stateRoot,
		LocalExitRoot:  ler,
		AccInputHash:   accInputHash,
	}

	dbTx, err := f.state.BeginStateTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin state transaction to open batch, err: %w", err)
	}

	// OpenBatch opens a new wip batch in the state
	err = f.state.OpenWIPBatch(ctx, newStateBatch, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return nil, fmt.Errorf("failed to rollback dbTx: %s. Error: %w", rollbackErr.Error(), err)
		}
		return nil, fmt.Errorf("failed to open new wip batch. Error: %w", err)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit database transaction for opening a wip batch. Error: %w", err)
	}

	// Check if synchronizer is up-to-date
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	return &Batch{
		batchNumber:         newStateBatch.BatchNumber,
		coinbase:            newStateBatch.Coinbase,
		initialStateRoot:    newStateBatch.StateRoot,
		imStateRoot:         newStateBatch.StateRoot,
		finalStateRoot:      newStateBatch.StateRoot,
		initialAccInputHash: newStateBatch.AccInputHash,
		imAccInputHash:      newStateBatch.AccInputHash,
		finalAccInputHash:   newStateBatch.AccInputHash,
		timestamp:           newStateBatch.Timestamp,
		globalExitRoot:      newStateBatch.GlobalExitRoot,
		localExitRoot:       newStateBatch.LocalExitRoot,
		remainingResources:  getMaxRemainingResources(f.batchConstraints),
		closingReason:       state.EmptyClosingReason,
	}, err
}

// closeWIPBatch closes the current batch in the state
func (f *finalizer) closeWIPBatch(ctx context.Context) error {
	/*transactions, effectivePercentages, err := f.dbManager.GetTransactionsByBatchNumber(ctx, f.wipBatch.batchNumber)
	if err != nil {
		return fmt.Errorf("failed to get transactions from transactions, err: %w", err)
	}
	for i, tx := range transactions {
		log.Debugf("[closeWIPBatch] BatchNum: %d, Tx position: %d, txHash: %s", f.wipBatch.batchNumber, i, tx.Hash().String())
	}*/
	usedResources := getUsedBatchResources(f.batchConstraints, f.wipBatch.remainingResources)
	receipt := state.ProcessingReceipt{
		BatchNumber:    f.wipBatch.batchNumber,
		BatchResources: usedResources,
		ClosingReason:  f.wipBatch.closingReason,
	}

	dbTx, err := f.state.BeginStateTransaction(ctx)
	if err != nil {
		return err
	}

	err = f.state.CloseWIPBatch(ctx, receipt, dbTx)
	if err != nil {
		err2 := dbTx.Rollback(ctx)
		if err2 != nil {
			log.Errorf("[CloseWIPBatch] error rolling back: %v", err2)
		}
		return err
	} else {
		err := dbTx.Commit(ctx)
		if err != nil {
			log.Errorf("[CloseWIPBatch] error committing: %v", err)
			return err
		}
	}

	return nil
}
