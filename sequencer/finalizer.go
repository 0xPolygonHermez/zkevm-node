package sequencer

import (
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

type finalizer struct {
	ForcedBatchCh        chan state.Batch
	GERCh                chan common.Hash
	L2ReorgCh            chan struct{}
	SendingToL1TimeoutCh chan bool
	SequencerAddress     common.Address
	worker               workerInterface
	dbManager            dbManagerInterface
	executor             executorInterface
}

type wipBatch struct {
	batchNumber           uint64
	initialStateRoot      common.Hash
	intermediaryStateRoot common.Hash
	timestamp             uint64
	GER                   common.Hash // 0x000...0 means to not update
	txs                   []TxTracker
	remainingResources    RemainingResources // Decide if we use remaining or used to keep track of the current batch resources
}

func newFinalizer(worker workerInterface, dbManager dbManagerInterface, executor executorInterface) *finalizer {
	// TODO: Init channels
	return &finalizer{worker: worker, dbManager: dbManager, executor: executor}
}

func (f *finalizer) Start() {
	// TODO: Load last batch from DB

	wipBatch := wipBatch{}

	// Most of this will need mutex since goroutines (Finalize txs, L1 requirements) can modify concurrently
	var (
		nextGER                 common.Hash
		nextGERDeadline         int64
		nextForcedBatches       []state.Batch
		nextForcedBatchDeadline int64
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
			// - Worker pool and efficicency list should be cleaned and updated
			// - WIP batch should be discarted (taking care to not lose txs or GER update for later on)
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

				// If it is not the first tx in the batch GER should be 0x00
				ger := wipBatch.GER

				// TODO: Populate processBatchRequest dinamic fields

				result := f.executor.ExecuteTransaction(processBatchRequest)
				// executionResult := execute tx (only passing to the executor this tx, starting at currentBatch.intermediaryRoot)
				// if ko:
				// // decide if we [MoveTxToNotReady, DeleteTx, UpdateTx]
				// if ok:
				// // currentBatch.remainingResources.sub(executionResult.Counters)
				// // if resourced overflows:
				// // // f.worker.UpdateTx()
				// // // continue (TODO: actually no, because we need to check the timeouts)
				// // add tx to current batch
				// // send tx to the DB (already finalized) through channel (async store)
				// // remove tx from pool (already finalized) through channel (async store)
				// // f.worker.UpdateAfterSingleSuccessfulTxExecution()
			} else {
				// check if the currentBatch is above a defined limit (almost full)
				// if yes:
				// // close batch and send batch info to the DB (already finalized) through channel (async store)
				// // go (decide if we need to execute the full batch as a sanity check, DO IT IN PARALLEL) ==> if error: log this txs somewhere and remove them from the pipeline
				// // if there are pending forced batches, execute them
				// // open batch: check if we have a new GER and update timestamp
				// else: activate flag
			}
			// Check deadlines
			if time.Now().Unix() >= nextForcedBatchDeadline {
				// close batch
				for len(nextForcedBatches) > 0 {
					// forced batch = pop forced batch
					// execute forced batch
					// send state mod log update through chan
				}
				nextForcedBatchDeadline = 0
				// open batch
			}
			if time.Now().Unix() >= nextGERDeadline {
				// close batch
				// open batch (with new GER)
				nextGER = common.Hash{}
				nextGERDeadline = 0
			}
			// if sleep flag activated: sleep
		}
	}()
}
