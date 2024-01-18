package sequencer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/ethereum/go-ethereum/common"
)

// Batch represents a wip or processed batch.
type Batch struct {
	batchNumber        uint64
	coinbase           common.Address
	timestamp          time.Time
	initialStateRoot   common.Hash // initial stateRoot of the batch
	imStateRoot        common.Hash // intermediate stateRoot that is updated each time a single tx is processed
	finalStateRoot     common.Hash // final stateroot of the batch when a L2 block is processed
	localExitRoot      common.Hash
	countOfTxs         int
	countOfL2Blocks    int
	remainingResources state.BatchResources
	closingReason      state.ClosingReason
}

func (w *Batch) isEmpty() bool {
	return w.countOfL2Blocks == 0
}

// setWIPBatch sets finalizer wip batch to the state batch passed as parameter
func (f *finalizer) setWIPBatch(ctx context.Context, wipStateBatch *state.Batch) (*Batch, error) {
	// Retrieve prevStateBatch to init the initialStateRoot of the wip batch
	prevStateBatch, err := f.stateIntf.GetBatchByNumber(ctx, wipStateBatch.BatchNumber-1, nil)
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

	remainingResources := getMaxRemainingResources(f.batchConstraints)
	err = remainingResources.Sub(wipStateBatch.Resources)
	if err != nil {
		return nil, err
	}

	wipBatch := &Batch{
		batchNumber:        wipStateBatch.BatchNumber,
		coinbase:           wipStateBatch.Coinbase,
		imStateRoot:        wipStateBatch.StateRoot,
		initialStateRoot:   prevStateBatch.StateRoot,
		finalStateRoot:     wipStateBatch.StateRoot,
		localExitRoot:      wipStateBatch.LocalExitRoot,
		timestamp:          wipStateBatch.Timestamp,
		countOfTxs:         wipStateBatchCountOfTxs,
		remainingResources: remainingResources,
	}

	return wipBatch, nil
}

// initWIPBatch inits the wip batch
func (f *finalizer) initWIPBatch(ctx context.Context) {
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	lastBatchNum, err := f.stateIntf.GetLastBatchNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get last batch number, error: %v", err)
	}

	// Get the last batch in trusted state
	lastStateBatch, err := f.stateIntf.GetBatchByNumber(ctx, lastBatchNum, nil)
	if err != nil {
		log.Fatalf("failed to get last batch, error: %v", err)
	}

	isClosed := !lastStateBatch.WIP

	log.Infof("batch %d isClosed: %v", lastBatchNum, isClosed)

	if isClosed { //if the last batch is close then open a new wip batch
		if lastStateBatch.BatchNumber+1 == f.cfg.HaltOnBatchNumber {
			f.Halt(ctx, fmt.Errorf("finalizer reached stop sequencer on batch number: %d", f.cfg.HaltOnBatchNumber))
		}

		f.wipBatch, err = f.openNewWIPBatch(ctx, lastStateBatch.BatchNumber+1, lastStateBatch.StateRoot, lastStateBatch.LocalExitRoot)
		if err != nil {
			log.Fatalf("failed to open new wip batch, error: %v", err)
		}
	} else { /// if it's not closed, it is the wip state batch, set it as wip batch in the finalizer
		f.wipBatch, err = f.setWIPBatch(ctx, lastStateBatch)
		if err != nil {
			log.Fatalf("failed to set wip batch, error: %v", err)
		}
	}

	log.Infof("initial batch: %d, initialStateRoot: %s, stateRoot: %s, coinbase: %s, LER: %s",
		f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.finalStateRoot, f.wipBatch.coinbase, f.wipBatch.localExitRoot)
}

// finalizeBatch retries until successful closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new empty batch
func (f *finalizer) finalizeBatch(ctx context.Context) {
	start := time.Now()
	defer func() {
		metrics.ProcessingTime(time.Since(start))
	}()

	// Finalize the wip L2 block if it has transactions, if not we keep it open to store it in the new wip batch
	if !f.wipL2Block.isEmpty() {
		f.finalizeL2Block(ctx)
	}

	err := f.closeAndOpenNewWIPBatch(ctx)
	if err != nil {
		f.Halt(ctx, fmt.Errorf("failed to create new WIP batch, error: %v", err))
	}
}

// closeAndOpenNewWIPBatch closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new empty batch
func (f *finalizer) closeAndOpenNewWIPBatch(ctx context.Context) error {
	// Wait until all L2 blocks are processed by the executor
	startWait := time.Now()
	f.pendingL2BlocksToProcessWG.Wait()
	elapsed := time.Since(startWait)
	stateMetrics.ExecutorProcessingTime(string(stateMetrics.SequencerCallerLabel), elapsed)
	log.Debugf("waiting for pending L2 blocks to be processed took: %v", elapsed)

	// Wait until all L2 blocks are store
	startWait = time.Now()
	f.pendingL2BlocksToStoreWG.Wait()
	log.Debugf("waiting for pending L2 blocks to be stored took: %v", time.Since(startWait))

	var err error

	// Reprocess full batch as sanity check
	if f.cfg.SequentialBatchSanityCheck {
		// Do the full batch reprocess now
		_, err := f.batchSanityCheck(ctx, f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.finalStateRoot)
		if err != nil {
			// There is an error reprocessing the batch. We halt the execution of the Sequencer at this point
			return fmt.Errorf("halting sequencer because of error reprocessing full batch %d (sanity check), error: %v ", f.wipBatch.batchNumber, err)
		}
	} else {
		// Do the full batch reprocess in parallel
		go func() {
			_, _ = f.batchSanityCheck(ctx, f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.finalStateRoot)
		}()
	}

	// Close the wip batch
	err = f.closeWIPBatch(ctx)
	if err != nil {
		return fmt.Errorf("failed to close batch, error: %v", err)
	}

	log.Infof("batch %d closed", f.wipBatch.batchNumber)

	if f.wipBatch.batchNumber+1 == f.cfg.HaltOnBatchNumber {
		f.Halt(ctx, fmt.Errorf("finalizer reached stop sequencer on batch number: %d", f.cfg.HaltOnBatchNumber))
	}

	// Metadata for the next batch
	stateRoot := f.wipBatch.finalStateRoot
	lastBatchNumber := f.wipBatch.batchNumber

	// Process forced batches
	if len(f.nextForcedBatches) > 0 {
		lastBatchNumber, stateRoot = f.processForcedBatches(ctx, lastBatchNumber, stateRoot)
		// We must init/reset the wip L2 block from the state since processForcedBatches has created new L2 blocks
		f.initWIPL2Block(ctx)
	}

	batch, err := f.openNewWIPBatch(ctx, lastBatchNumber+1, stateRoot, f.wipBatch.localExitRoot)
	if err != nil {
		return fmt.Errorf("failed to open new wip batch, error: %v", err)
	}

	// Subtract the L2 block used resources to batch
	err = batch.remainingResources.Sub(f.wipL2Block.getUsedResources())
	if err != nil {
		return fmt.Errorf("failed to subtract L2 block used resources to new wip batch %d, error: %v", batch.batchNumber, err)
	}

	f.wipBatch = batch

	log.Infof("new WIP batch %d", f.wipBatch.batchNumber)

	return nil
}

// openNewWIPBatch opens a new batch in the state and returns it as WipBatch
func (f *finalizer) openNewWIPBatch(ctx context.Context, batchNumber uint64, stateRoot, LER common.Hash) (*Batch, error) {
	// open next batch
	newStateBatch := state.Batch{
		BatchNumber:    batchNumber,
		Coinbase:       f.sequencerAddress,
		Timestamp:      now(),
		GlobalExitRoot: state.ZeroHash,
		StateRoot:      stateRoot,
		LocalExitRoot:  LER,
	}

	dbTx, err := f.stateIntf.BeginStateTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin state transaction to open batch, error: %v", err)
	}

	// OpenBatch opens a new wip batch in the state
	err = f.stateIntf.OpenWIPBatch(ctx, newStateBatch, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return nil, fmt.Errorf("failed to rollback due to error when open a new wip batch, rollback error: %v, error: %v", rollbackErr, err)
		}
		return nil, fmt.Errorf("failed to open new wip batch, error: %v", err)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit database transaction for opening a wip batch, error: %v", err)
	}

	// Send batch bookmark to the datastream
	f.DSSendBatchBookmark(batchNumber)

	// Check if synchronizer is up-to-date
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	return &Batch{
		batchNumber:        newStateBatch.BatchNumber,
		coinbase:           newStateBatch.Coinbase,
		initialStateRoot:   newStateBatch.StateRoot,
		imStateRoot:        newStateBatch.StateRoot,
		finalStateRoot:     newStateBatch.StateRoot,
		timestamp:          newStateBatch.Timestamp,
		localExitRoot:      newStateBatch.LocalExitRoot,
		remainingResources: getMaxRemainingResources(f.batchConstraints),
		closingReason:      state.EmptyClosingReason,
	}, err
}

// closeWIPBatch closes the current batch in the state
func (f *finalizer) closeWIPBatch(ctx context.Context) error {
	usedResources := getUsedBatchResources(f.batchConstraints, f.wipBatch.remainingResources)
	receipt := state.ProcessingReceipt{
		BatchNumber:    f.wipBatch.batchNumber,
		BatchResources: usedResources,
		ClosingReason:  f.wipBatch.closingReason,
	}

	dbTx, err := f.stateIntf.BeginStateTransaction(ctx)
	if err != nil {
		return err
	}

	err = f.stateIntf.CloseWIPBatch(ctx, receipt, dbTx)
	if err != nil {
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back due to error when closing wip batch, rollback error: %v, error: %v", rollbackErr, err)
		}
		return err
	} else {
		err := dbTx.Commit(ctx)
		if err != nil {
			log.Errorf("error committing close wip batch, error: %v", err)
			return err
		}
	}

	return nil
}

// maxTxsPerBatchReached checks if the batch has reached the maximum number of txs per batch
func (f *finalizer) maxTxsPerBatchReached() bool {
	if f.wipBatch.countOfTxs >= int(f.batchConstraints.MaxTxsPerBatch) {
		log.Infof("closing batch %d, because it reached the maximum number of txs", f.wipBatch.batchNumber)
		f.wipBatch.closingReason = state.BatchFullClosingReason
		return true
	}
	return false
}

// batchSanityCheck reprocesses a batch used as sanity check
func (f *finalizer) batchSanityCheck(ctx context.Context, batchNum uint64, initialStateRoot common.Hash, expectedNewStateRoot common.Hash) (*state.ProcessBatchResponse, error) {
	reprocessError := func(batch *state.Batch) {
		rawL2Blocks, err := state.DecodeBatchV2(batch.BatchL2Data)
		if err != nil {
			log.Errorf("error decoding BatchL2Data for batch %d, error: %v", batch.BatchNumber, err)
			return
		}

		// Log batch detailed info
		log.Errorf("batch %d sanity check error: initialStateRoot: %s, expectedNewStateRoot: %s", batch.BatchNumber, initialStateRoot, expectedNewStateRoot)
		for i, rawL2block := range rawL2Blocks.Blocks {
			for j, rawTx := range rawL2block.Transactions {
				log.Infof("batch %d, block position: %d, tx position: %d, tx hash: %s", batch.BatchNumber, i, j, rawTx.Tx.Hash())
			}
		}

		f.Halt(ctx, fmt.Errorf("batch sanity check error. Check previous errors in logs to know which was the cause"))
	}

	log.Debugf("batch %d sanity check: initialStateRoot: %s, expectedNewStateRoot: %s", batchNum, initialStateRoot, expectedNewStateRoot)

	batch, err := f.stateIntf.GetBatchByNumber(ctx, batchNum, nil)
	if err != nil {
		log.Errorf("failed to get batch %d, error: %v", batchNum, err)
		return nil, ErrGetBatchByNumber
	}

	caller := stateMetrics.DiscardCallerLabel
	if f.cfg.SequentialBatchSanityCheck {
		caller = stateMetrics.SequencerCallerLabel
	}

	batchRequest := state.ProcessRequest{
		BatchNumber:             batch.BatchNumber,
		L1InfoRoot_V2:           mockL1InfoRoot,
		OldStateRoot:            initialStateRoot,
		Transactions:            batch.BatchL2Data,
		Coinbase:                batch.Coinbase,
		TimestampLimit_V2:       uint64(time.Now().Unix()),
		ForkID:                  f.stateIntf.GetForkIDByBatchNumber(batch.BatchNumber),
		SkipVerifyL1InfoRoot_V2: true,
		Caller:                  caller,
	}
	batchRequest.L1InfoTreeData_V2, _, _, err = f.stateIntf.GetL1InfoTreeDataFromBatchL2Data(ctx, batch.BatchL2Data, nil)
	if err != nil {
		log.Errorf("failed to get L1InfoTreeData for batch %d, error: %v", batch.BatchNumber, err)
		reprocessError(nil)
		return nil, ErrGetBatchByNumber
	}

	var batchResponse *state.ProcessBatchResponse

	startProcessing := time.Now()
	batchResponse, err = f.stateIntf.ProcessBatchV2(ctx, batchRequest, false)
	endProcessing := time.Now()

	if err != nil {
		log.Errorf("failed to process batch %d, error: %v", batch.BatchNumber, err)
		reprocessError(batch)
		return nil, ErrProcessBatch
	}

	if batchResponse.ExecutorError != nil {
		log.Errorf("executor error when reprocessing batch %d, error: %v", batch.BatchNumber, batchResponse.ExecutorError)
		reprocessError(batch)
		return nil, ErrExecutorError
	}

	if batchResponse.IsRomOOCError {
		log.Errorf("failed to process batch %d because OutOfCounters", batch.BatchNumber)
		reprocessError(batch)

		payload, err := json.Marshal(batchRequest)
		if err != nil {
			log.Errorf("error marshaling payload, error: %v", err)
		} else {
			event := &event.Event{
				ReceivedAt:  time.Now(),
				Source:      event.Source_Node,
				Component:   event.Component_Sequencer,
				Level:       event.Level_Critical,
				EventID:     event.EventID_ReprocessFullBatchOOC,
				Description: string(payload),
				Json:        batchRequest,
			}
			err = f.eventLog.LogEvent(ctx, event)
			if err != nil {
				log.Errorf("error storing payload, error: %v", err)
			}
		}

		return nil, ErrProcessBatchOOC
	}

	if batchResponse.NewStateRoot != expectedNewStateRoot {
		log.Errorf("new state root mismatch for batch %d, expected: %s, got: %s", batch.BatchNumber, expectedNewStateRoot.String(), batchResponse.NewStateRoot.String())
		reprocessError(batch)
		return nil, ErrStateRootNoMatch
	}

	log.Infof("successful sanity check for batch %d, initialStateRoot: %s, stateRoot: %s, l2Blocks: %d, time: %v, %s",
		batch.BatchNumber, initialStateRoot, batchResponse.NewStateRoot.String(), len(batchResponse.BlockResponses),
		endProcessing.Sub(startProcessing), f.logZKCounters(batchResponse.UsedZkCounters))

	return batchResponse, nil
}

// checkRemainingResources checks if the resources passed as parameters fits in the wip batch.
func (f *finalizer) checkRemainingResources(result *state.ProcessBatchResponse, bytes uint64) error {
	usedResources := state.BatchResources{
		ZKCounters: result.UsedZkCounters,
		Bytes:      bytes,
	}

	return f.wipBatch.remainingResources.Sub(usedResources)
}

// logZKCounters returns a string with all the zkCounters values
func (f *finalizer) logZKCounters(counters state.ZKCounters) string {
	return fmt.Sprintf("gasUsed: %d, keccakHashes: %d, poseidonHashes: %d, poseidonPaddings: %d, memAligns: %d, arithmetics: %d, binaries: %d, sha256Hashes: %d, steps: %d",
		counters.GasUsed, counters.UsedKeccakHashes, counters.UsedPoseidonHashes, counters.UsedPoseidonPaddings, counters.UsedMemAligns, counters.UsedArithmetics,
		counters.UsedBinaries, counters.UsedSha256Hashes_V2, counters.UsedSteps)
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
		log.Infof("closing batch %d because it reached %s limit", f.wipBatch.batchNumber, resourceDesc)
		f.wipBatch.closingReason = state.BatchAlmostFullClosingReason
	}

	return result
}

// getConstraintThresholdUint64 returns the threshold for the given input
func (f *finalizer) getConstraintThresholdUint64(input uint64) uint64 {
	return input * uint64(f.cfg.ResourceExhaustedMarginPct) / 100 //nolint:gomnd
}

// getConstraintThresholdUint32 returns the threshold for the given input
func (f *finalizer) getConstraintThresholdUint32(input uint32) uint32 {
	return input * f.cfg.ResourceExhaustedMarginPct / 100 //nolint:gomnd
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
