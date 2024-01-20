package sequencer

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/ethereum/go-ethereum/common"
)

// L2Block represents a wip or processed L2 block
type L2Block struct {
	timestamp                 time.Time
	deltaTimestamp            uint32
	initialStateRoot          common.Hash
	l1InfoTreeExitRoot        state.L1InfoTreeExitRootStorageEntry
	l1InfoTreeExitRootChanged bool
	transactions              []*TxTracker
	batchResponse             *state.ProcessBatchResponse
}

func (b *L2Block) isEmpty() bool {
	return len(b.transactions) == 0
}

// addTx adds a tx to the L2 block
func (b *L2Block) addTx(tx *TxTracker) {
	b.transactions = append(b.transactions, tx)
}

// getL1InfoTreeIndex returns the L1InfoTreeIndex that must be used when processing/storing the block
func (b *L2Block) getL1InfoTreeIndex() uint32 {
	// If the L1InfoTreeIndex has changed in this block then we return the new index, otherwise we return 0
	if b.l1InfoTreeExitRootChanged {
		return b.l1InfoTreeExitRoot.L1InfoTreeIndex
	} else {
		return 0
	}
}

// getUsedResources returns the estimated batch resources used to process the changeL2block tx for this block
func (b *L2Block) getUsedResources() state.BatchResources {
	// If the L1InfoTreeIndex has changed we return the resources used when the index != 0,
	// otherwise we return the used resources when the index = 0
	if b.l1InfoTreeExitRootChanged {
		return l2BlockUsedResourcesIndexNonZero
	} else {
		return l2BlockUsedResourcesIndexZero
	}
}

// getWriteInfoRootUsedResources returns the additional batch resources used when processing
// this block with the SkipWriteBlockInfoRoot_V2 flag of the executor request to false
// func (b *L2Block) getWriteInfoRootUsedResources() state.BatchResources {
// 	if b.l1InfoTreeExitRootChanged {
// 		return l2BlockUsedResourcesIndexNonZero
// 	} else {
// 		return l2BlockUsedResourcesIndexZero
// 	}
// }

// initWIPL2Block inits the wip L2 block
func (f *finalizer) initWIPL2Block(ctx context.Context) {
	// Wait to l1InfoTree to be updated for first time
	f.lastL1InfoTreeCond.L.Lock()
	for !f.lastL1InfoTreeValid {
		log.Infof("waiting for L1InfoTree to be updated")
		f.lastL1InfoTreeCond.Wait()
	}
	f.lastL1InfoTreeCond.L.Unlock()

	lastL2Block, err := f.stateIntf.GetLastL2Block(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get last L2 block number, error: %v", err)
	}

	f.openNewWIPL2Block(ctx, &lastL2Block.ReceivedAt)
}

// addPendingL2BlockToProcess adds a pending L2 block that is closed and ready to be processed by the executor
func (f *finalizer) addPendingL2BlockToProcess(ctx context.Context, l2Block *L2Block) {
	f.pendingL2BlocksToProcessWG.Add(1)

	select {
	case f.pendingL2BlocksToProcess <- l2Block:
	case <-ctx.Done():
		// If context is cancelled before we can send to the channel, we must decrement the WaitGroup count and
		// delete the pending TxToStore added in the worker
		f.pendingL2BlocksToProcessWG.Done()
	}
}

// addPendingL2BlockToStore adds a L2 block that is ready to be stored in the state DB once its flushid has been stored by the executor
func (f *finalizer) addPendingL2BlockToStore(ctx context.Context, l2Block *L2Block) {
	f.pendingL2BlocksToStoreWG.Add(1)

	for _, tx := range l2Block.transactions {
		f.workerIntf.AddPendingTxToStore(tx.Hash, tx.From)
	}

	select {
	case f.pendingL2BlocksToStore <- l2Block:
	case <-ctx.Done():
		// If context is cancelled before we can send to the channel, we must decrement the WaitGroup count and
		// delete the pending TxToStore added in the worker
		f.pendingL2BlocksToStoreWG.Done()
		for _, tx := range l2Block.transactions {
			f.workerIntf.DeletePendingTxToStore(tx.Hash, tx.From)
		}
	}
}

// processPendingL2Blocks processes (executor) the pending to process L2 blocks
func (f *finalizer) processPendingL2Blocks(ctx context.Context) {
	for {
		select {
		case l2Block, ok := <-f.pendingL2BlocksToProcess:
			if !ok {
				// Channel is closed
				return
			}

			l2Block.initialStateRoot = f.wipBatch.finalStateRoot

			log.Infof("processing L2 block, batch: %d, initialStateRoot: %s txs: %d, l1InfoTreeIndex: %d",
				f.wipBatch.batchNumber, l2Block.initialStateRoot, len(l2Block.transactions), l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)

			startProcessing := time.Now()
			batchResponse, err := f.processL2Block(ctx, l2Block)
			endProcessing := time.Now()

			if err != nil {
				f.Halt(ctx, fmt.Errorf("error processing L2 block, error: %v", err))
			}

			if len(batchResponse.BlockResponses) == 0 {
				f.Halt(ctx, fmt.Errorf("error processing L2 block, error: BlockResponses returned by the executor is empty"))
			}

			blockResponse := batchResponse.BlockResponses[0]

			// Sanity check. Check blockResponse.TransactionsReponses match l2Block.Transactions length, order and tx hashes
			if len(blockResponse.TransactionResponses) != len(l2Block.transactions) {
				f.Halt(ctx, fmt.Errorf("error processing L2 block, error: length of TransactionsResponses %d don't match length of l2Block.transactions %d",
					len(blockResponse.TransactionResponses), len(l2Block.transactions)))
			}
			for i, txResponse := range blockResponse.TransactionResponses {
				if txResponse.TxHash != l2Block.transactions[i].Hash {
					f.Halt(ctx, fmt.Errorf("error processing L2 block, error: TransactionsResponses hash %s in position %d don't match l2Block.transactions[%d] hash %s",
						txResponse.TxHash.String(), i, i, l2Block.transactions[i].Hash))
				}
			}

			l2Block.batchResponse = batchResponse

			// Update finalStateRoot of the batch to the newStateRoot for the L2 block
			f.wipBatch.finalStateRoot = l2Block.batchResponse.NewStateRoot

			log.Infof("processed L2 block: %d, batch: %d, initialStateRoot: %s, stateRoot: %s, txs: %d/%d, blockHash: %s, infoRoot: %s, time: %v, %s",
				blockResponse.BlockNumber, f.wipBatch.batchNumber, l2Block.initialStateRoot, l2Block.batchResponse.NewStateRoot, len(l2Block.transactions),
				len(blockResponse.TransactionResponses), blockResponse.BlockHash, blockResponse.BlockInfoRoot.String(), endProcessing.Sub(startProcessing),
				f.logZKCounters(batchResponse.UsedZkCounters))

			f.updateFlushIDs(batchResponse.FlushID, batchResponse.StoredFlushID)

			f.addPendingL2BlockToStore(ctx, l2Block)

			f.pendingL2BlocksToProcessWG.Done()
		case <-ctx.Done():
			// The context was cancelled from outside, Wait for all goroutines to finish, cleanup and exit
			f.pendingL2BlocksToProcessWG.Wait()
			return
		default:
			time.Sleep(100 * time.Millisecond) //nolint:gomnd
		}
	}
}

// storePendingTransactions stores the pending L2 blocks in the database
func (f *finalizer) storePendingL2Blocks(ctx context.Context) {
	for {
		select {
		case l2Block, ok := <-f.pendingL2BlocksToStore:
			if !ok {
				// Channel is closed
				return
			}

			// Wait until L2 block has been flushed/stored by the executor
			f.storedFlushIDCond.L.Lock()
			for f.storedFlushID < l2Block.batchResponse.FlushID {
				f.storedFlushIDCond.Wait()
				// check if context is done after waking up
				if ctx.Err() != nil {
					f.storedFlushIDCond.L.Unlock()
					return
				}
			}
			f.storedFlushIDCond.L.Unlock()

			// If the L2 block has txs now f.storedFlushID >= l2BlockToStore.flushId, we can store tx
			blockResponse := l2Block.batchResponse.BlockResponses[0]
			log.Infof("storing L2 block: %d, batch: %d, txs: %d/%d, blockHash: %s, infoRoot: %s",
				blockResponse.BlockNumber, f.wipBatch.batchNumber, len(l2Block.transactions), len(blockResponse.TransactionResponses),
				blockResponse.BlockHash, blockResponse.BlockInfoRoot.String())

			startStoring := time.Now()
			err := f.storeL2Block(ctx, l2Block)
			endStoring := time.Now()

			if err != nil {
				f.Halt(ctx, fmt.Errorf("error storing L2 block %d, error: %v", l2Block.batchResponse.BlockResponses[0].BlockNumber, err))
			}

			log.Infof("stored L2 block: %d, batch: %d, txs: %d/%d, blockHash: %s, infoRoot: %s, time: %v",
				blockResponse.BlockNumber, f.wipBatch.batchNumber, len(l2Block.transactions), len(blockResponse.TransactionResponses),
				blockResponse.BlockHash, blockResponse.BlockInfoRoot.String(), endStoring.Sub(startStoring))

			for _, tx := range l2Block.transactions {
				// Delete the tx from the pending list in the worker (addrQueue)
				f.workerIntf.DeletePendingTxToStore(tx.Hash, tx.From)
			}

			f.pendingL2BlocksToStoreWG.Done()
		case <-ctx.Done():
			// The context was cancelled from outside, Wait for all goroutines to finish, cleanup and exit
			f.pendingL2BlocksToStoreWG.Wait()
			return
		default:
			time.Sleep(100 * time.Millisecond) //nolint:gomnd
		}
	}
}

// processL2Block process (executor) a L2 Block and adds it to the pendingL2BlocksToStore channel. It returns the response block from the executor
func (f *finalizer) processL2Block(ctx context.Context, l2Block *L2Block) (*state.ProcessBatchResponse, error) {
	processL2BLockError := func(err error) {
		log.Errorf("process L2 block error %v, batch: %d, initialStateRoot: %s", err, f.wipBatch.batchNumber, l2Block.initialStateRoot.String())
		// Log batch detailed info
		for i, tx := range l2Block.transactions {
			log.Infof("batch: %d, tx position %d, tx hash: %s", f.wipBatch.batchNumber, i, tx.HashStr)
		}
	}

	batchL2Data := []byte{}

	// Add changeL2Block to batchL2Data
	changeL2BlockBytes := f.stateIntf.BuildChangeL2Block(l2Block.deltaTimestamp, l2Block.getL1InfoTreeIndex())
	batchL2Data = append(batchL2Data, changeL2BlockBytes...)

	// Add transactions data to batchL2Data
	for _, tx := range l2Block.transactions {
		epHex, err := hex.DecodeHex(fmt.Sprintf("%x", tx.EGPPercentage))
		if err != nil {
			log.Errorf("error decoding hex value for effective gas price percentage for tx %s, error: %v", tx.HashStr, err)
			return nil, err
		}

		txData := append(tx.RawTx, epHex...)

		batchL2Data = append(batchL2Data, txData...)
	}

	batchRequest := state.ProcessRequest{
		BatchNumber:               f.wipBatch.batchNumber,
		OldStateRoot:              l2Block.initialStateRoot,
		Coinbase:                  f.wipBatch.coinbase,
		L1InfoRoot_V2:             mockL1InfoRoot,
		TimestampLimit_V2:         uint64(l2Block.timestamp.Unix()),
		Transactions:              batchL2Data,
		SkipFirstChangeL2Block_V2: false,
		SkipWriteBlockInfoRoot_V2: false,
		Caller:                    stateMetrics.DiscardCallerLabel,
		ForkID:                    f.stateIntf.GetForkIDByBatchNumber(f.wipBatch.batchNumber),
		SkipVerifyL1InfoRoot_V2:   true,
		L1InfoTreeData_V2:         map[uint32]state.L1DataV2{},
	}
	batchRequest.L1InfoTreeData_V2[l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex] = state.L1DataV2{
		GlobalExitRoot: l2Block.l1InfoTreeExitRoot.GlobalExitRoot.GlobalExitRoot,
		BlockHashL1:    l2Block.l1InfoTreeExitRoot.PreviousBlockHash,
		MinTimestamp:   uint64(l2Block.l1InfoTreeExitRoot.GlobalExitRoot.Timestamp.Unix()),
	}

	var (
		err           error
		batchResponse *state.ProcessBatchResponse
	)

	batchResponse, err = f.stateIntf.ProcessBatchV2(ctx, batchRequest, true)

	if err != nil {
		processL2BLockError(err)
		return nil, err
	}

	if batchResponse.ExecutorError != nil {
		processL2BLockError(err)
		return nil, ErrExecutorError
	}

	if batchResponse.IsRomOOCError {
		processL2BLockError(err)
		return nil, ErrProcessBatchOOC
	}

	return batchResponse, nil
}

// storeL2Block stores the L2 block in the state and updates the related batch and transactions
func (f *finalizer) storeL2Block(ctx context.Context, l2Block *L2Block) error {
	//log.Infof("storeL2Block: storing processed txToStore: %s", txToStore.response.TxHash.String())
	dbTx, err := f.stateIntf.BeginStateTransaction(ctx)
	if err != nil {
		return fmt.Errorf("error creating db transaction to store L2 block, error: %v", err)
	}

	rollbackOnError := func(retError error) error {
		err := dbTx.Rollback(ctx)
		if err != nil {
			return fmt.Errorf("rollback error due to error %v, error: %v", retError, err)
		}
		return retError
	}

	blockResponse := l2Block.batchResponse.BlockResponses[0]
	forkID := f.stateIntf.GetForkIDByBatchNumber(f.wipBatch.batchNumber)

	txsEGPLog := []*state.EffectiveGasPriceLog{}
	for _, tx := range l2Block.transactions {
		egpLog := tx.EGPLog
		txsEGPLog = append(txsEGPLog, &egpLog)
	}

	// Store L2 block in the state
	err = f.stateIntf.StoreL2Block(ctx, f.wipBatch.batchNumber, blockResponse, txsEGPLog, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("database error on storing L2 block %d, error: %v", blockResponse.BlockNumber, err))
	}

	// Now we need to update de BatchL2Data of the wip batch and also update the status of the L2 block txs in the pool

	batch, err := f.stateIntf.GetBatchByNumber(ctx, f.wipBatch.batchNumber, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("error when getting batch %d from the state, error: %v", f.wipBatch.batchNumber, err))
	}

	// Add changeL2Block to batch.BatchL2Data
	blockL2Data := []byte{}
	changeL2BlockBytes := f.stateIntf.BuildChangeL2Block(l2Block.deltaTimestamp, l2Block.getL1InfoTreeIndex())
	blockL2Data = append(blockL2Data, changeL2BlockBytes...)

	// Add transactions data to batch.BatchL2Data
	for _, txResponse := range blockResponse.TransactionResponses {
		txData, err := state.EncodeTransaction(txResponse.Tx, uint8(txResponse.EffectivePercentage), forkID)
		if err != nil {
			return rollbackOnError(fmt.Errorf("error when encoding tx %s, error: %v", txResponse.TxHash.String(), err))
		}
		blockL2Data = append(blockL2Data, txData...)
	}

	batch.BatchL2Data = append(batch.BatchL2Data, blockL2Data...)
	batch.Resources.SumUp(state.BatchResources{ZKCounters: l2Block.batchResponse.UsedZkCounters, Bytes: uint64(len(blockL2Data))})

	receipt := state.ProcessingReceipt{
		BatchNumber:    f.wipBatch.batchNumber,
		StateRoot:      l2Block.batchResponse.NewStateRoot,
		LocalExitRoot:  l2Block.batchResponse.NewLocalExitRoot,
		BatchL2Data:    batch.BatchL2Data,
		BatchResources: batch.Resources,
	}

	// We need to update the batch GER only in the GER of the block (response) is not zero, since the final GER stored in the batch
	// must be the last GER from the blocks that is not zero (last L1InfoRootIndex change)
	if blockResponse.GlobalExitRoot != state.ZeroHash {
		receipt.GlobalExitRoot = blockResponse.GlobalExitRoot
	} else {
		receipt.GlobalExitRoot = batch.GlobalExitRoot
	}

	err = f.stateIntf.UpdateWIPBatch(ctx, receipt, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("error when updating wip batch %d, error: %v", f.wipBatch.batchNumber, err))
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return err
	}

	// Update txs status in the pool
	for _, txResponse := range blockResponse.TransactionResponses {
		// Change Tx status to selected
		err = f.poolIntf.UpdateTxStatus(ctx, txResponse.TxHash, pool.TxStatusSelected, false, nil)
		if err != nil {
			return err
		}
	}

	// Send L2 block to data streamer
	err = f.DSSendL2Block(f.wipBatch.batchNumber, blockResponse)
	if err != nil {
		//TODO: we need to halt/rollback the L2 block if we had an error sending to the data streamer?
		log.Errorf("error sending L2 block %d to data streamer, error: %v", blockResponse.BlockNumber, err)
	}

	return nil
}

// finalizeL2Block closes the current L2 block and opens a new one
func (f *finalizer) finalizeL2Block(ctx context.Context) {
	log.Debugf("finalizing L2 block")

	f.closeWIPL2Block(ctx)

	f.openNewWIPL2Block(ctx, nil)
}

func (f *finalizer) closeWIPL2Block(ctx context.Context) {
	// If the L2 block is empty (no txs) We need to process it to update the state root and remaining batch resources before closing it
	if f.wipL2Block.isEmpty() {
		log.Debug("processing WIP L2 block because it is empty")
		if err := f.processEmptyL2Block(ctx); err != nil {
			f.Halt(ctx, fmt.Errorf("failed to process empty WIP L2 block, error: %v ", err))
		}
	}

	f.wipBatch.countOfL2Blocks++

	f.addPendingL2BlockToProcess(ctx, f.wipL2Block)
}

func (f *finalizer) openNewWIPL2Block(ctx context.Context, prevTimestamp *time.Time) {
	newL2Block := &L2Block{}

	newL2Block.timestamp = now()
	if prevTimestamp != nil {
		newL2Block.deltaTimestamp = uint32(newL2Block.timestamp.Sub(*prevTimestamp).Truncate(time.Second).Seconds())
	} else {
		newL2Block.deltaTimestamp = uint32(newL2Block.timestamp.Sub(f.wipL2Block.timestamp).Truncate(time.Second).Seconds())
	}

	newL2Block.transactions = []*TxTracker{}

	f.lastL1InfoTreeMux.Lock()
	newL2Block.l1InfoTreeExitRoot = f.lastL1InfoTree
	f.lastL1InfoTreeMux.Unlock()

	// Check if L1InfoTreeIndex has changed, in this case we need to use this index in the changeL2block instead of zero
	// If it's the first wip L2 block after starting sequencer (wipL2Block == nil) then we assume that the L1InfoTreeIndex has changed (there is no problem assuming this)
	if f.wipL2Block != nil && newL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex != f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex {
		newL2Block.l1InfoTreeExitRootChanged = true
	}

	f.wipL2Block = newL2Block

	// Check if there are enough counters in the wip batch to store the new wip L2 block
	err := f.wipBatch.remainingResources.Sub(f.wipL2Block.getUsedResources())
	// We close the wip batch and open a new one if we got an error when subtracting the getUsedResources or we have exhausted some resources of the batch
	if err != nil || f.isBatchResourcesExhausted() {
		err := f.closeAndOpenNewWIPBatch(ctx)
		if err != nil {
			f.Halt(ctx, fmt.Errorf("failed to create new WIP batch, error: %v", err))
		}
	}

	log.Debugf("new WIP L2 block created: batch: %d, initialStateRoot: %s, timestamp: %d, l1InfoTreeIndex: %d",
		f.wipBatch.batchNumber, f.wipL2Block.initialStateRoot, f.wipL2Block.timestamp.Unix(), f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)
}
