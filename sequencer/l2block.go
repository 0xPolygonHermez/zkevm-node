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

var changeL2BlockMark = []byte{0x0B}

// L2Block represents a wip or processed L2 block
type L2Block struct {
	timestamp           time.Time
	deltaTimestamp      uint32
	initialStateRoot    common.Hash
	initialAccInputHash common.Hash
	batchNumber         uint64
	forcedBatch         bool
	coinbase            common.Address
	stateRoot           common.Hash
	l1InfoTreeExitRoot  state.L1InfoTreeExitRootStorageEntry
	transactions        []*TxTracker
	batchResponse       *state.ProcessBatchResponse
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
	log.Infof("L1Infotree updated. L1InfoTreeIndex: %d", f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)

	lastL2Block, err := f.dbManager.GetLastL2Block(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get last L2 block number. Error: %w", err)
	}

	f.openNewWIPL2Block(ctx, &lastL2Block.ReceivedAt)
}

// addPendingL2BlockToProcess adds a pending L2 block that is closed and ready to be processed by the executor
func (f *finalizer) addPendingL2BlockToProcess(ctx context.Context, l2Block *L2Block) {
	f.pendingL2BlocksToProcessWG.Add(1)

	for _, tx := range l2Block.transactions {
		f.worker.AddPendingTxToStore(tx.Hash, tx.From)
	}

	select {
	case f.pendingL2BlocksToProcess <- l2Block:
	case <-ctx.Done():
		// If context is cancelled before we can send to the channel, we must decrement the WaitGroup count and
		// delete the pending TxToStore added in the worker
		f.pendingL2BlocksToProcessWG.Done()
		for _, tx := range l2Block.transactions {
			f.worker.DeletePendingTxToStore(tx.Hash, tx.From)
		}
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

			log.Debugf("processing L2 block. Batch: %d, txs %d", l2Block.batchNumber, len(l2Block.transactions))
			batchResponse, err := f.processL2Block(ctx, l2Block)
			if err != nil {
				f.halt(ctx, fmt.Errorf("error processing L2 block. Error: %s", err))
			}

			if len(batchResponse.BlockResponses) == 0 {
				f.halt(ctx, fmt.Errorf("error processing L2 block. Error: BlockResponses returned by the executor is empty"))
			}

			blockResponse := batchResponse.BlockResponses[0]
			log.Infof("L2 block %d processed. Batch: %d, txs: %d/%d, blockHash: %s, infoRoot: %s",
				blockResponse.BlockNumber, l2Block.batchNumber, len(l2Block.transactions), len(blockResponse.TransactionResponses),
				blockResponse.BlockHash, blockResponse.BlockInfoRoot.String())

			l2Block.batchResponse = batchResponse

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

			// If the L2 block has txs wait until f.storedFlushID >= l2BlockToStore.flushId (this flushId is from the last tx in the L2 block)
			if len(l2Block.transactions) > 0 {
				lastFlushId := l2Block.transactions[len(l2Block.transactions)-1].FlushId
				f.storedFlushIDCond.L.Lock()
				for f.storedFlushID < lastFlushId {
					f.storedFlushIDCond.Wait()
					// check if context is done after waking up
					if ctx.Err() != nil {
						f.storedFlushIDCond.L.Unlock()
						return
					}
				}
				f.storedFlushIDCond.L.Unlock()
			}

			// If the L2 block has txs now f.storedFlushID >= l2BlockToStore.flushId, we can store tx
			blockResponse := l2Block.batchResponse.BlockResponses[0]
			log.Debugf("storing L2 block %d. Batch: %d, txs: %d/%d, blockHash: %s, infoRoot: %s",
				blockResponse.BlockNumber, l2Block.batchNumber, len(l2Block.transactions), len(blockResponse.TransactionResponses),
				blockResponse.BlockHash, blockResponse.BlockInfoRoot.String())

			err := f.storeL2Block(ctx, l2Block)
			if err != nil {
				//TODO: this doesn't halt the finalizer, review howto do it
				f.halt(ctx, fmt.Errorf("error storing L2 block %d. Error: %s", l2Block.batchResponse.BlockResponses[0].BlockNumber, err))
			}

			log.Infof("L2 block %d stored. Batch: %d, txs: %d/%d, blockHash: %s, infoRoot: %s",
				blockResponse.BlockNumber, l2Block.batchNumber, len(l2Block.transactions), len(blockResponse.TransactionResponses),
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
		log.Infof("[processL2Block] BatchNumber: %d, InitialStateRoot: %s, ExpectedNewStateRoot: %s", l2Block.batchNumber, l2Block.initialStateRoot.String(), l2Block.stateRoot.String())
		for i, tx := range l2Block.transactions {
			log.Infof("[processL2Block] BatchNumber: %d, tx position %d, tx hash: %s", l2Block.batchNumber, i, tx.HashStr)
		}
	}

	log.Debugf("[processL2Block] BatchNumber: %d, Txs: %d, InitialStateRoot: %s, ExpectedNewStateRoot: %s", l2Block.batchNumber, len(l2Block.transactions), l2Block.initialStateRoot.String(), l2Block.stateRoot.String())

	batchL2Data := []byte{}

	// Add changeL2Block to batchL2Data
	changeL2BlockBytes := f.dbManager.BuildChangeL2Block(l2Block.deltaTimestamp, l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)
	batchL2Data = append(batchL2Data, changeL2BlockBytes...)

	// Add transactions data to batchL2Data
	for _, tx := range l2Block.transactions {
		ep, err := f.effectiveGasPrice.CalculateEffectiveGasPricePercentage(tx.GasPrice, tx.EffectiveGasPrice) //TODO: store effectivePercentage in TxTracker
		if err != nil {
			log.Errorf("[processL2Block] error calculating effective gas price percentage for tx %s. Error: %s", tx.HashStr, err)
			return nil, err
		}

		//TODO: Create function to add epHex to batchL2Data as it's used in several places
		epHex, err := hex.DecodeHex(fmt.Sprintf("%x", ep))
		if err != nil {
			log.Errorf("[processL2Block] error decoding hex value for effective gas price percentage for tx %s. Error: %s", tx.HashStr, err)
			return nil, err
		}

		txData := append(tx.RawTx, epHex...)

		batchL2Data = append(batchL2Data, txData...)
	}

	// TODO: review this request
	executorBatchRequest := state.ProcessRequest{
		BatchNumber:               l2Block.batchNumber,
		OldStateRoot:              l2Block.initialStateRoot,
		OldAccInputHash:           l2Block.initialAccInputHash,
		Coinbase:                  l2Block.coinbase,
		L1InfoRoot_V2:             mockL1InfoRoot,
		TimestampLimit_V2:         uint64(l2Block.timestamp.Unix()),
		Transactions:              batchL2Data,
		SkipFirstChangeL2Block_V2: false,
		SkipWriteBlockInfoRoot_V2: false,
		Caller:                    stateMetrics.SequencerCallerLabel,
	}

	var (
		err    error
		result *state.ProcessBatchResponse
	)

	result, err = f.executor.ProcessBatchV2(ctx, executorBatchRequest, true)
	if err != nil {
		processL2BLockError()
		return nil, err
	}

	//TODO: check this error in first place?
	if result.ExecutorError != nil {
		processL2BLockError()
		return nil, ErrExecutorError
	}

	if result.IsRomOOCError {
		processL2BLockError()
		return nil, ErrProcessBatchOOC
	}

	if result.NewStateRoot != l2Block.stateRoot {
		log.Errorf("[processL2Block] new state root mismatch for L2 block %d in batch %d, expected: %s, got: %s",
			result.BlockResponses[0].BlockNumber, l2Block.batchNumber, l2Block.stateRoot.String(), result.NewStateRoot.String())
		processL2BLockError()
		return nil, ErrStateRootNoMatch
	}

	//TODO: check that result.BlockResponse is not empty

	return result, nil
}

// storeL2Block stores the L2 block in the state and updates the related batch and transactions
func (f *finalizer) storeL2Block(ctx context.Context, l2Block *L2Block) error {
	//log.Infof("storeL2Block: storing processed txToStore: %s", txToStore.response.TxHash.String())

	blockResponse := l2Block.batchResponse.BlockResponses[0]
	forkID := f.dbManager.GetForkIDByBatchNumber(l2Block.batchNumber)

	dbTx, err := f.dbManager.BeginStateTransaction(ctx)
	if err != nil {
		return fmt.Errorf("[storeL2Block] error creating db transaction. Error: %w", err)
	}

	txsEGPLog := []*state.EffectiveGasPriceLog{}
	for _, tx := range l2Block.transactions {
		egpLog := tx.EGPLog
		txsEGPLog = append(txsEGPLog, &egpLog)
	}

	// Store L2 block in the state
	err = f.dbManager.StoreL2Block(ctx, l2Block.batchNumber, l2Block.batchResponse.BlockResponses[0], txsEGPLog, dbTx)
	if err != nil {
		return fmt.Errorf("[storeL2Block] database error on storing L2 block. Error: %w", err)
	}

	// If the L2 block belongs to a regular batch (not forced) then we need to update de BatchL2Data
	// also in this case we need to update the status of the L2 block txs in the pool
	// TODO: review this
	if !l2Block.forcedBatch {
		batch, err := f.dbManager.GetBatchByNumber(ctx, l2Block.batchNumber, dbTx)
		if err != nil {
			err2 := dbTx.Rollback(ctx)
			if err2 != nil {
				log.Errorf("[storeL2Block] failed to rollback dbTx when getting batch that gave err: %s. Rollback err: %s", err, err2)
			}
			return fmt.Errorf("[storeL2Block] error when getting batch %d from the state. Error: %s", l2Block.batchNumber, err)
		}

		// Add changeL2Block to batch.BatchL2Data
		changeL2BlockBytes := f.dbManager.BuildChangeL2Block(l2Block.deltaTimestamp, l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)
		batch.BatchL2Data = append(batch.BatchL2Data, changeL2BlockBytes...)

		// Add transactions data to batch.BatchL2Data
		for _, txResponse := range blockResponse.TransactionResponses {
			txData, err := state.EncodeTransaction(txResponse.Tx, uint8(txResponse.EffectivePercentage), forkID)
			if err != nil {
				return err
			}
			batch.BatchL2Data = append(batch.BatchL2Data, txData...)
		}

		err = f.dbManager.UpdateBatch(ctx, l2Block.batchNumber, batch.BatchL2Data, l2Block.batchResponse.NewLocalExitRoot, dbTx)
		if err != nil {
			err2 := dbTx.Rollback(ctx)
			if err2 != nil {
				log.Errorf("[storeL2Block] failed to rollback dbTx when getting batch that gave err: %s. Rollback err: %s", err, err2)
			}
			return err
		}

		for _, txResponse := range blockResponse.TransactionResponses {
			// Change Tx status to selected
			err = f.dbManager.UpdateTxStatus(ctx, txResponse.TxHash, pool.TxStatusSelected, false, nil)
			if err != nil {
				return err
			}
		}
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return err
	}

	// Send L2 block to data streamer
	err = f.dbManager.DSSendL2Block(l2Block)
	if err != nil {
		return fmt.Errorf("[storeL2Block] error sending L2 block %d to data streamer", blockResponse.BlockNumber)
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

	// Update L2 block stateroot to the WIP batch stateroot
	f.wipL2Block.stateRoot = f.wipBatch.stateRoot

	f.addPendingL2BlockToProcess(ctx, f.wipL2Block)
}

func (f *finalizer) openNewWIPL2Block(ctx context.Context, prevTimestamp *time.Time) {
	//TODO: use better f.wipBatch.remainingResources.Sub() instead to subtract directly
	// Subtract the bytes needed to store the changeL2Block of the new L2 block into the WIP batch
	f.wipBatch.remainingResources.Bytes = f.wipBatch.remainingResources.Bytes - changeL2BlockSize
	// Subtract poseidon and arithmetic counters need to calculate the InfoRoot when the L2 block is closed
	f.wipBatch.remainingResources.ZKCounters.UsedPoseidonHashes = f.wipBatch.remainingResources.ZKCounters.UsedPoseidonHashes - 256 // nolint:gomnd //TODO: config param
	f.wipBatch.remainingResources.ZKCounters.UsedArithmetics = f.wipBatch.remainingResources.ZKCounters.UsedArithmetics - 1         //TODO: config param
	// After do the subtracts we need to check if we have not reached the size limit for the batch
	if f.isBatchResourcesExhausted() {
		// If we have reached the limit then close the wip batch and create a new one
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

	newL2Block.batchNumber = f.wipBatch.batchNumber
	newL2Block.forcedBatch = false
	newL2Block.initialStateRoot = f.wipBatch.stateRoot
	newL2Block.stateRoot = f.wipBatch.stateRoot
	newL2Block.initialAccInputHash = f.wipBatch.accInputHash
	newL2Block.coinbase = f.wipBatch.coinbase
	newL2Block.transactions = []*TxTracker{}

	f.lastL1InfoTreeMux.Lock()
	newL2Block.l1InfoTreeExitRoot = f.lastL1InfoTree
	f.lastL1InfoTreeMux.Unlock()

	f.wipL2Block = newL2Block

	log.Debugf("new WIP L2 block created. Batch: %d, initialStateRoot: %s, timestamp: %d", f.wipL2Block.batchNumber, f.wipL2Block.initialStateRoot, f.wipL2Block.timestamp.Unix())
}
