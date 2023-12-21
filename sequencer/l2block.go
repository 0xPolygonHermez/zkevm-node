package sequencer

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	statePackage "github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/ethereum/go-ethereum/common"
)

// L2Block represents a wip or processed L2 block
type L2Block struct {
	timestamp          time.Time
	deltaTimestamp     uint32
	initialStateRoot   common.Hash
	l1InfoTreeExitRoot state.L1InfoTreeExitRootStorageEntry
	transactions       []*TxTracker
	batchResponse      *state.ProcessBatchResponse
}

func (b *L2Block) isEmpty() bool {
	return len(b.transactions) == 0
}

// addTx adds a tx to the L2 block
func (b *L2Block) addTx(tx *TxTracker) {
	b.transactions = append(b.transactions, tx)
}

// initWIPL2Block inits the wip L2 block
func (f *finalizer) initWIPL2Block(ctx context.Context) {
	f.wipL2Block = &L2Block{}

	// Wait to l1InfoTree to be updated for first time
	f.lastL1InfoTreeCond.L.Lock()
	for !f.lastL1InfoTreeValid {
		log.Infof("waiting for L1InfoTree to be updated")
		f.lastL1InfoTreeCond.Wait()
	}
	f.lastL1InfoTreeCond.L.Unlock()

	f.lastL1InfoTreeMux.Lock()
	f.wipL2Block.l1InfoTreeExitRoot = f.lastL1InfoTree
	f.lastL1InfoTreeMux.Unlock()

	lastL2Block, err := f.state.GetLastL2Block(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get last L2 block number. Error: %w", err)
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
		f.worker.AddPendingTxToStore(tx.Hash, tx.From)
	}

	select {
	case f.pendingL2BlocksToStore <- l2Block:
	case <-ctx.Done():
		// If context is cancelled before we can send to the channel, we must decrement the WaitGroup count and
		// delete the pending TxToStore added in the worker
		f.pendingL2BlocksToStoreWG.Done()
		for _, tx := range l2Block.transactions {
			f.worker.DeletePendingTxToStore(tx.Hash, tx.From)
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

			log.Debugf("processing L2 block. Batch: %d, initialStateRoot: %s txs: %d", f.wipBatch.batchNumber, l2Block.initialStateRoot, len(l2Block.transactions))
			batchResponse, err := f.processL2Block(ctx, l2Block)
			if err != nil {
				f.halt(ctx, fmt.Errorf("error processing L2 block. Error: %s", err))
			}

			if len(batchResponse.BlockResponses) == 0 {
				f.halt(ctx, fmt.Errorf("error processing L2 block. Error: BlockResponses returned by the executor is empty"))
			}

			blockResponse := batchResponse.BlockResponses[0]

			// Sanity check. Check blockResponse.TransactionsReponses match l2Block.Transactions length, order and tx hashes
			if len(blockResponse.TransactionResponses) != len(l2Block.transactions) {
				f.halt(ctx, fmt.Errorf("error processing L2 block. Error: length of TransactionsResponses %d don't match length of l2Block.transactions %d",
					len(blockResponse.TransactionResponses), len(l2Block.transactions)))
			}
			for i, txResponse := range blockResponse.TransactionResponses {
				if txResponse.TxHash != l2Block.transactions[i].Hash {
					f.halt(ctx, fmt.Errorf("error processing L2 block. Error: TransactionsResponses hash %s in position %d don't match l2Block.transactions[%d] hash %s",
						txResponse.TxHash.String(), i, i, l2Block.transactions[i].Hash))
				}
			}

			l2Block.batchResponse = batchResponse

			// Update finalStateRoot and accInputHash of the batch to the newStateRoot and NewAccInputHash for the L2 block
			f.wipBatch.finalStateRoot = l2Block.batchResponse.NewStateRoot

			log.Infof("L2 block %d processed. Batch: %d, initialStateRoot: %s, stateRoot: %s, txs: %d/%d, blockHash: %s, infoRoot: %s",
				blockResponse.BlockNumber, f.wipBatch.batchNumber, l2Block.initialStateRoot, l2Block.batchResponse.NewStateRoot,
				len(l2Block.transactions), len(blockResponse.TransactionResponses), blockResponse.BlockHash, blockResponse.BlockInfoRoot.String())

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
			log.Debugf("storing L2 block %d. Batch: %d, txs: %d/%d, blockHash: %s, infoRoot: %s",
				blockResponse.BlockNumber, f.wipBatch.batchNumber, len(l2Block.transactions), len(blockResponse.TransactionResponses),
				blockResponse.BlockHash, blockResponse.BlockInfoRoot.String())

			err := f.storeL2Block(ctx, l2Block)
			if err != nil {
				f.halt(ctx, fmt.Errorf("error storing L2 block %d. Error: %s", l2Block.batchResponse.BlockResponses[0].BlockNumber, err))
			}

			log.Infof("L2 block %d stored. Batch: %d, txs: %d/%d, blockHash: %s, infoRoot: %s",
				blockResponse.BlockNumber, f.wipBatch.batchNumber, len(l2Block.transactions), len(blockResponse.TransactionResponses),
				blockResponse.BlockHash, blockResponse.BlockInfoRoot.String())

			for _, tx := range l2Block.transactions {
				// Delete the tx from the pending list in the worker (addrQueue)
				f.worker.DeletePendingTxToStore(tx.Hash, tx.From)
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
	processL2BLockError := func() {
		// Log batch detailed info
		log.Infof("[processL2Block] BatchNumber: %d, InitialStateRoot: %s", f.wipBatch.batchNumber, l2Block.initialStateRoot.String())
		for i, tx := range l2Block.transactions {
			log.Infof("[processL2Block] BatchNumber: %d, tx position %d, tx hash: %s", f.wipBatch.batchNumber, i, tx.HashStr)
		}
	}

	batchL2Data := []byte{}

	// Add changeL2Block to batchL2Data
	changeL2BlockBytes := f.state.BuildChangeL2Block(l2Block.deltaTimestamp, l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)
	batchL2Data = append(batchL2Data, changeL2BlockBytes...)

	// Add transactions data to batchL2Data
	for _, tx := range l2Block.transactions {
		epHex, err := hex.DecodeHex(fmt.Sprintf("%x", tx.EGPPercentage))
		if err != nil {
			log.Errorf("[processL2Block] error decoding hex value for effective gas price percentage for tx %s. Error: %s", tx.HashStr, err)
			return nil, err
		}

		txData := append(tx.RawTx, epHex...)

		batchL2Data = append(batchL2Data, txData...)
	}

	executorBatchRequest := state.ProcessRequest{
		BatchNumber:               f.wipBatch.batchNumber,
		OldStateRoot:              l2Block.initialStateRoot,
		Coinbase:                  f.wipBatch.coinbase,
		L1InfoRoot_V2:             mockL1InfoRoot,
		TimestampLimit_V2:         uint64(l2Block.timestamp.Unix()),
		Transactions:              batchL2Data,
		SkipFirstChangeL2Block_V2: false,
		SkipWriteBlockInfoRoot_V2: false,
		Caller:                    stateMetrics.SequencerCallerLabel,
		ForkID:                    f.state.GetForkIDByBatchNumber(f.wipBatch.batchNumber),
		SkipVerifyL1InfoRoot_V2:   true,
		L1InfoTreeData_V2:         map[uint32]statePackage.L1DataV2{},
	}
	executorBatchRequest.L1InfoTreeData_V2[l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex] = state.L1DataV2{
		GlobalExitRoot: l2Block.l1InfoTreeExitRoot.GlobalExitRoot.GlobalExitRoot,
		BlockHashL1:    l2Block.l1InfoTreeExitRoot.PreviousBlockHash,
		MinTimestamp:   uint64(l2Block.l1InfoTreeExitRoot.GlobalExitRoot.Timestamp.Unix()),
	}

	var (
		err    error
		result *state.ProcessBatchResponse
	)

	result, err = f.state.ProcessBatchV2(ctx, executorBatchRequest, true)
	if err != nil {
		processL2BLockError()
		return nil, err
	}

	if result.ExecutorError != nil {
		processL2BLockError()
		return nil, ErrExecutorError
	}

	if result.IsRomOOCError {
		processL2BLockError()
		return nil, ErrProcessBatchOOC
	}

	return result, nil
}

// storeL2Block stores the L2 block in the state and updates the related batch and transactions
func (f *finalizer) storeL2Block(ctx context.Context, l2Block *L2Block) error {
	//log.Infof("storeL2Block: storing processed txToStore: %s", txToStore.response.TxHash.String())
	dbTx, err := f.state.BeginStateTransaction(ctx)
	if err != nil {
		return fmt.Errorf("[storeL2Block] error creating db transaction. Error: %w", err)
	}

	rollbackOnError := func(retError error) error {
		err := dbTx.Rollback(ctx)
		if err != nil {
			return fmt.Errorf("[storeL2Block] rollback error due to error %w. Error: %w", retError, err)
		}
		return retError
	}

	blockResponse := l2Block.batchResponse.BlockResponses[0]
	forkID := f.state.GetForkIDByBatchNumber(f.wipBatch.batchNumber)

	txsEGPLog := []*state.EffectiveGasPriceLog{}
	for _, tx := range l2Block.transactions {
		egpLog := tx.EGPLog
		txsEGPLog = append(txsEGPLog, &egpLog)
	}

	// Store L2 block in the state
	err = f.state.StoreL2Block(ctx, f.wipBatch.batchNumber, blockResponse, txsEGPLog, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("[storeL2Block] database error on storing L2 block %d. Error: %s", blockResponse.BlockNumber, err))
	}

	// Now we need to update de BatchL2Data of the wip batch and also update the status of the L2 block txs in the pool

	batch, err := f.state.GetBatchByNumber(ctx, f.wipBatch.batchNumber, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("[storeL2Block] error when getting batch %d from the state. Error: %w", f.wipBatch.batchNumber, err))
	}

	// Add changeL2Block to batch.BatchL2Data
	blockL2Data := []byte{}
	changeL2BlockBytes := f.state.BuildChangeL2Block(l2Block.deltaTimestamp, l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)
	blockL2Data = append(blockL2Data, changeL2BlockBytes...)

	// Add transactions data to batch.BatchL2Data
	for _, txResponse := range blockResponse.TransactionResponses {
		txData, err := state.EncodeTransaction(txResponse.Tx, uint8(txResponse.EffectivePercentage), forkID)
		if err != nil {
			return rollbackOnError(fmt.Errorf("[storeL2Block] error when encoding transaction %s. Error: %w", txResponse.TxHash.String(), err))
		}
		blockL2Data = append(blockL2Data, txData...)
	}

	batch.BatchL2Data = append(batch.BatchL2Data, blockL2Data...)
	batch.Resources.SumUp(state.BatchResources{ZKCounters: l2Block.batchResponse.UsedZkCounters, Bytes: uint64(len(blockL2Data))})

	receipt := state.ProcessingReceipt{
		BatchNumber:    f.wipBatch.batchNumber,
		GlobalExitRoot: l2Block.l1InfoTreeExitRoot.GlobalExitRoot.GlobalExitRoot,
		LocalExitRoot:  l2Block.batchResponse.NewLocalExitRoot,
		AccInputHash:   l2Block.batchResponse.NewAccInputHash,
		BatchL2Data:    batch.BatchL2Data,
		BatchResources: batch.Resources,
	}

	err = f.state.UpdateWIPBatch(ctx, receipt, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("[storeL2Block] error when updating wip batch %d. Error: %w", f.wipBatch.batchNumber, err))
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return err
	}

	// Update txs status in the pool
	for _, txResponse := range blockResponse.TransactionResponses {
		// Change Tx status to selected
		err = f.pool.UpdateTxStatus(ctx, txResponse.TxHash, pool.TxStatusSelected, false, nil)
		if err != nil {
			return err
		}
	}

	// Send L2 block to data streamer
	err = f.DSSendL2Block(f.wipBatch.batchNumber, blockResponse)
	if err != nil {
		//TODO: we need to halt/rollback the L2 block if we had an error sending to the data streamer?
		log.Errorf("[storeL2Block] error sending L2 block %d to data streamer", blockResponse.BlockNumber)
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
	// If the L2 block is empty (no txs) We need to process it to update the state root before closing it
	if f.wipL2Block.isEmpty() {
		log.Debug("processing L2 block because it is empty")
		if _, err := f.processTransaction(ctx, nil, true); err != nil {
			f.halt(ctx, fmt.Errorf("failed to process empty L2 block. Error: %s ", err))
		}
	}

	f.addPendingL2BlockToProcess(ctx, f.wipL2Block)
}

func (f *finalizer) openNewWIPL2Block(ctx context.Context, prevTimestamp *time.Time) {
	err := f.wipBatch.remainingResources.Sub(l2BlockUsedResources)

	// we finalize the wip batch if we got an error when subtracting the l2BlockUsedResources or we have exhausted some resources of the batch
	if err != nil || f.isBatchResourcesExhausted() {
		f.finalizeBatch(ctx)
	}

	// Initialize wipL2Block to a new L2 block
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

	f.wipL2Block = newL2Block

	log.Debugf("new WIP L2 block created. Batch: %d, initialStateRoot: %s, timestamp: %d", f.wipBatch.batchNumber, f.wipL2Block.initialStateRoot, f.wipL2Block.timestamp.Unix())
}
