package sequencer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

type finalizer struct {
	ForcedBatchCh        chan state.Batch
	GERCh                chan common.Hash
	L2ReorgCh            chan struct{}
	SendingToL1TimeoutCh chan bool
	TxsToStoreCh         chan *txToStore
	WgTxsToStore         *sync.WaitGroup
	sequencerAddress     common.Address
	worker               workerInterface
	dbManager            dbManagerInterface
	executor             stateInterface
	maxTxsPerBatch       uint64
}

func newFinalizer(worker workerInterface, dbManager dbManagerInterface, executor stateInterface, sequencerAddr common.Address, maxTxsPerBatch uint64) *finalizer {
	return &finalizer{
		ForcedBatchCh:        make(chan state.Batch),
		GERCh:                make(chan common.Hash),
		L2ReorgCh:            make(chan struct{}),
		SendingToL1TimeoutCh: make(chan bool),
		TxsToStoreCh:         make(chan *txToStore, maxTxsPerBatch),
		WgTxsToStore:         &sync.WaitGroup{},
		sequencerAddress:     sequencerAddr,
		worker:               worker,
		dbManager:            dbManager,
		executor:             executor,
		maxTxsPerBatch:       maxTxsPerBatch,
	}
}

type txToStore struct {
	txResponse               *state.ProcessTransactionResponse
	batchNumber              uint64
	previousL2BlockStateRoot common.Hash
}

type wipBatch struct {
	batchNumber        uint64
	coinbase           common.Address
	accInputHash       common.Hash
	initialStateRoot   common.Hash
	stateRoot          common.Hash
	timestamp          uint64
	globalExitRoot     common.Hash // 0x000...0 (ZeroHash) means to not update
	txs                []TxTracker
	remainingResources BatchResources
}

func (f *finalizer) Start(ctx context.Context, batch wipBatch, OldStateRoot, OldAccInputHash common.Hash) {
	var (
		// Most of this will need mutex since goroutines (Finalize txs, L1 requirements) can modify concurrently
		nextGER                 common.Hash
		nextGERDeadline         int64
		nextGERMux              sync.RWMutex // TODO: Check all places where need to be used
		nextForcedBatches       []state.Batch
		nextForcedBatchDeadline int64
		nextForcedBatchesMux    sync.RWMutex // TODO: Check all places where need to be used
		//wipBatchMux             sync.RWMutex // TODO: Check all places where need to be used
		err error
	)
	fmt.Printf(nextGER.Hex())

	processRequest := state.ProcessSingleTxRequest{
		BatchNumber:      batch.batchNumber,
		StateRoot:        batch.stateRoot,
		OldStateRoot:     OldStateRoot,
		GlobalExitRoot:   batch.globalExitRoot,
		OldAccInputHash:  OldAccInputHash,
		SequencerAddress: f.sequencerAddress,
		Timestamp:        batch.timestamp,
		Caller:           state.SequencerCallerLabel,
	}

	// TODO: Finish all receivers handling!
	// Closing signals receiver
	go func() {
		for {
			select {
			// Forced  batch ch
			case fb := <-f.ForcedBatchCh:
				nextForcedBatches = append(nextForcedBatches, fb)
				if nextForcedBatchDeadline > 0 {
					nextForcedBatchDeadline = time.Now().Unix() // + configurable delay
				}
			// globalExitRoot ch
			case ger := <-f.GERCh:
				nextGER = ger
				if nextGERDeadline > 0 {
					nextGERDeadline = time.Now().Unix() // + configurable delay
				}
			}
			// L2 reorg ch: TODO: analyze how we handle L2 reorg. Needs work from Sync as well. Some considerations:
			// - Txs that have been popped from the state should go back to the pool
			// - Worker pool and efficiency list should be cleaned and updated
			// - WIP batch should be discarded (taking care to not lose txs or globalExitRoot update for later on)
			// Too many time without batches in L1 ch
			// Any other externality from the point of view of the sequencer should be captured using this pattern
		}
	}()

	// Finalize txs
	go func() {
		for {
			tx := f.worker.GetBestFittingTx(batch.remainingResources)
			if tx != nil {
				var ger common.Hash
				if len(batch.txs) == 0 {
					ger = batch.globalExitRoot
				} else {
					ger = state.ZeroHash
				}
				processRequest.GlobalExitRoot = ger
				result := f.executor.ProcessSingleTx(processRequest)

				// executionResult := execute tx (only passing to the executor this tx, starting at currentBatch.intermediaryRoot)
				if result.Error != nil {
					// // decide if we [MoveTxToNotReady, DeleteTx, UpdateTx]
				} else {
					processedTx := result.Responses[0]
					usedResources := BatchResources{
						zKCounters: result.UsedZkCounters,
						bytes:      uint64(processedTx.Tx.Size()),
						gas:        processedTx.GasUsed,
					}
					err = batch.remainingResources.Sub(usedResources)
					if err != nil {
						f.worker.UpdateTx(processedTx.TxHash, tx.From, usedResources.zKCounters)
						// (TODO: check the timeouts)
						continue
					}

					// We have a successful processing if we are here
					previousL2BlockStateRoot := batch.stateRoot
					batch.stateRoot = result.NewStateRoot
					batch.accInputHash = result.NewAccInputHash
					processRequest.StateRoot = result.NewStateRoot
					processRequest.OldAccInputHash = result.NewAccInputHash
					// We store the processed transaction, add it to the batch and delete from the pool atomically.
					f.WgTxsToStore.Add(1)
					f.TxsToStoreCh <- &txToStore{batchNumber: batch.batchNumber, txResponse: processedTx, previousL2BlockStateRoot: previousL2BlockStateRoot}
					f.worker.UpdateAfterSingleSuccessfulTxExecution(tx.From, result.TouchedAddresses)
				}
			} else {
				// TODO: check if the currentBatch is above a defined limit (almost full)
				// TODO: if yes:
				receipt := ClosingBatchParameters{
					BatchNumber:   batch.batchNumber,
					AccInputHash:  processRequest.OldAccInputHash,
					StateRoot:     batch.stateRoot,
					LocalExitRoot: processRequest.GlobalExitRoot,
					Txs:           batch.txs,
				}
				f.WgTxsToStore.Wait()
				f.dbManager.CloseBatch(ctx, receipt)
				batch.batchNumber += 1
				batch.txs = make([]TxTracker, 0, f.maxTxsPerBatch)
				// TODO: OPEN NEW BATCH IN DB

				// TODO: Maybe to check the database for new globalExitRoot and update the batch?
				// TODO:
				// // go (decide if we need to execute the full batch as a sanity check, DO IT IN PARALLEL) ==> if error: log this txs somewhere and remove them from the pipeline
				// // if there are pending forced batches, execute them
				// // open batch: check if we have a new globalExitRoot and update timestamp
				// else: activate flag
			}
			// Check deadlines
			if time.Now().Unix() >= nextForcedBatchDeadline {
				nextForcedBatchesMux.Lock()
				// close batch
				for len(nextForcedBatches) > 0 {
					// forced batch = pop forced batch
					// execute forced batch
					// send state mod log update through chan
				}
				nextForcedBatchDeadline = 0
				nextForcedBatchesMux.Unlock()
				// open batch
			}
			if time.Now().Unix() >= nextGERDeadline {
				nextGERMux.Lock()
				// close batch
				// open batch (with new globalExitRoot)
				nextGER = common.Hash{}
				nextGERDeadline = 0
				nextGERMux.Unlock()
			}
			// if sleep flag activated: sleep
			<-ctx.Done()
		}
	}()
}
