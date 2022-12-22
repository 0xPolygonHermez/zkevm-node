package sequencer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

type finalizer struct {
	ForcedBatchCh        chan state.Batch
	GERCh                chan common.Hash
	L2ReorgCh            chan struct{}
	SendingToL1TimeoutCh chan bool
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
		sequencerAddress:     sequencerAddr,
		worker:               worker,
		dbManager:            dbManager,
		executor:             executor,
		maxTxsPerBatch:       maxTxsPerBatch,
	}
}

type wipBatch struct {
	batchNumber           uint64
	initialStateRoot      common.Hash
	intermediaryStateRoot common.Hash
	timestamp             uint64
	GER                   common.Hash // 0x000...0 (ZeroHash) means to not update
	txs                   []TxTracker
	remainingResources    RemainingResources
}

func (f *finalizer) Start(ctx context.Context) {
	lastBatch, err := f.dbManager.GetLastBatch(ctx)
	if err != nil {
		log.Fatal("failed to get last batch")
	}

	wipBatch := wipBatch{
		GER:                   lastBatch.GlobalExitRoot,
		initialStateRoot:      lastBatch.StateRoot,
		intermediaryStateRoot: lastBatch.StateRoot,
		txs:                   make([]TxTracker, 0, f.maxTxsPerBatch),
	}

	// Most of this will need mutex since goroutines (Finalize txs, L1 requirements) can modify concurrently
	var (
		nextGER                 common.Hash
		nextGERDeadline         int64
		nextGERMux              sync.RWMutex // TODO: Check all places where need to be used
		nextForcedBatches       []state.Batch
		nextForcedBatchDeadline int64
		nextForcedBatchesMux    sync.RWMutex // TODO: Check all places where need to be used
		//wipBatchMux             sync.RWMutex // TODO: Check all places where need to be used
	)

	fmt.Printf(nextGER.Hex())

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
			// GER ch
			case ger := <-f.GERCh:
				nextGER = ger
				if nextGERDeadline > 0 {
					nextGERDeadline = time.Now().Unix() // + configurable delay
				}
			}
			// L2 reorg ch: TODO: analyze how we handle L2 reorg. Needs work from Sync as well. Some considerations:
			// - Txs that have been popped from the state should go back to the pool
			// - Worker pool and efficiency list should be cleaned and updated
			// - WIP batch should be discarded (taking care to not lose txs or GER update for later on)
			// Too many time without batches in L1 ch
			// Any other externality from the point of view of the sequencer should be captured using this pattern
		}
	}()

	// Finalize txs
	go func() {
		processBatchRequest := state.ProcessBatchRequest{}
		for {
			tx := f.worker.GetBestFittingTx(wipBatch.remainingResources)
			if tx != nil {
				lenOfTxs := len(wipBatch.txs)
				isFirstTx := lenOfTxs == 1

				var ger common.Hash
				if isFirstTx {
					ger = wipBatch.GER
				} else {
					ger = state.ZeroHash
				}

				// TODO: Populate processBatchRequest dynamic fields
				processBatchRequest.GlobalExitRoot = ger
				processBatchRequest.IsFirstTx = isFirstTx
				result := f.executor.ExecuteTransaction(processBatchRequest)

				// executionResult := execute tx (only passing to the executor this tx, starting at currentBatch.intermediaryRoot)
				if result.Error != nil {
					// // decide if we [MoveTxToNotReady, DeleteTx, UpdateTx]

				} else {
					processedTx := result.Responses[0]
					from, err := types.Sender(types.NewEIP155Signer(processedTx.Tx.ChainId()), &processedTx.Tx)
					if err != nil {
						from, err = types.Sender(types.HomesteadSigner{}, &processedTx.Tx)
					}

					err = wipBatch.remainingResources.Sub(RemainingResources{
						remainingZKCounters: result.UsedZkCounters,
						remainingBytes:      uint64(processedTx.Tx.Size()),
						remainingGas:        processedTx.GasUsed,
					})
					if err != nil {
						f.worker.UpdateTx(processedTx.TxHash, from, wipBatch.remainingResources.remainingZKCounters)
						// (TODO: check the timeouts)
						continue
					}

					// We have a successful processing if we are here
					wipBatch.intermediaryStateRoot = result.NewStateRoot
					processBatchRequest.StateRoot = result.NewStateRoot
					processBatchRequest.OldAccInputHash = result.NewAccInputHash
					// We store the processed transaction, add it to the batch and delete from the pool atomically.
					go f.dbManager.StoreProcessedTxAndDeleteFromPool(ctx, wipBatch.batchNumber, processedTx)
					go f.worker.UpdateAfterSingleSuccessfulTxExecution(from, result.TouchedAddresses)
				}
			} else {
				// TODO: check if the currentBatch is above a defined limit (almost full)
				// TODO: if yes:
				go func() {
					ethTxs := make([]types.Transaction, 0, len(wipBatch.txs))
					for _, txTracker := range wipBatch.txs {
						tx, err := state.DecodeTx(string(txTracker.RawTx))
						if err != nil {
							log.Fatalf("failed to decode tx")
						}
						ethTxs = append(ethTxs, *tx)
					}
					receipt := state.ProcessingReceipt{
						BatchNumber:   wipBatch.batchNumber,
						AccInputHash:  processBatchRequest.OldAccInputHash,
						StateRoot:     wipBatch.intermediaryStateRoot,
						LocalExitRoot: processBatchRequest.GlobalExitRoot,
						Txs:           ethTxs,
					}
					f.dbManager.CloseBatch(ctx, receipt)
				}()
				wipBatch.batchNumber += 1
				wipBatch.txs = make([]TxTracker, 0, f.maxTxsPerBatch)
				// TODO: Maybe to check the database for new GER and update the wipBatch?

				// TODO:
				// // go (decide if we need to execute the full batch as a sanity check, DO IT IN PARALLEL) ==> if error: log this txs somewhere and remove them from the pipeline
				// // if there are pending forced batches, execute them
				// // open batch: check if we have a new GER and update timestamp
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
				// open batch (with new GER)
				nextGER = common.Hash{}
				nextGERDeadline = 0
				nextGERMux.Unlock()
			}
			// if sleep flag activated: sleep
			<-ctx.Done()
		}
	}()
}
