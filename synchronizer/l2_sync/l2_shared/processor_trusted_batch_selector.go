package l2_shared

/*
This class is a implementation of SyncTrustedStateExecutor that selects the executor to use.
It's ready to switch between pre-etrog and etrog as soon as the forkid 7 is activated.

When ForkId7 is activated, the executor will be switched to etrog for forkid7 batches.
*/

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
)

const etrogForkId = uint64(7)

type stateSyncTrustedStateExecutorSelector interface {
	GetForkIDInMemory(forkId uint64) *state.ForkIDInterval
}

// SyncTrustedStateExecutorSelector Implements SyncTrustedStateExecutor
type SyncTrustedStateExecutorSelector struct {
	executorPreEtrog syncinterfaces.SyncTrustedStateExecutor
	executorEtrog    syncinterfaces.SyncTrustedStateExecutor
	state            stateSyncTrustedStateExecutorSelector
}

// NewSyncTrustedStateExecutorSelector creates a new SyncTrustedStateExecutorSelector that implements SyncTrustedStateExecutor
func NewSyncTrustedStateExecutorSelector(
	preEtrog syncinterfaces.SyncTrustedStateExecutor,
	etrog syncinterfaces.SyncTrustedStateExecutor,
	state stateSyncTrustedStateExecutorSelector) *SyncTrustedStateExecutorSelector {
	return &SyncTrustedStateExecutorSelector{
		executorPreEtrog: preEtrog,
		executorEtrog:    etrog,
		state:            state,
	}
}

// GetExecutor returns the executor that should be used for the given batch, could be nil
// it returns the executor and the maximum batch number that the executor can process
func (s *SyncTrustedStateExecutorSelector) GetExecutor(latestSyncedBatch uint64, maximumBatchNumberToProcess uint64) (syncinterfaces.SyncTrustedStateExecutor, uint64) {
	fork := s.state.GetForkIDInMemory(etrogForkId)
	if fork == nil {
		log.Debugf("ForkId7 not activated yet, using pre-etrog executor")
		return s.executorPreEtrog, maximumBatchNumberToProcess
	}

	if latestSyncedBatch+1 >= fork.FromBatchNumber {
		log.Debugf("ForkId7 activated, batch:%d -> etrog executor", latestSyncedBatch)
		return s.executorEtrog, maximumBatchNumberToProcess
	}
	maxCapped := min(maximumBatchNumberToProcess, fork.FromBatchNumber-1)
	log.Debugf("ForkId7 activated, batch:%d -> pre-etrog executor (maxBatch from:%d to %d)",
		latestSyncedBatch, maximumBatchNumberToProcess, maxCapped)
	return s.executorPreEtrog, maxCapped
}

// SyncTrustedState syncs the trusted state with the permissionless state. In this case
// choose which executor must use
func (s *SyncTrustedStateExecutorSelector) SyncTrustedState(ctx context.Context, latestSyncedBatch uint64, maximumBatchNumberToProcess uint64) error {
	executor, maxBatchNumber := s.GetExecutor(latestSyncedBatch, maximumBatchNumberToProcess)
	if executor == nil {
		log.Warnf("No executor selected, skipping SyncTrustedState: latestSyncedBatch:%d, maximumBatchNumberToProcess:%d",
			latestSyncedBatch, maximumBatchNumberToProcess)
		return nil
	}
	return executor.SyncTrustedState(ctx, latestSyncedBatch, maxBatchNumber)
}

// CleanTrustedState clean cache of Batches and StateRoot
func (s *SyncTrustedStateExecutorSelector) CleanTrustedState() {
	if s.executorPreEtrog != nil {
		s.executorPreEtrog.CleanTrustedState()
	}
	if s.executorEtrog != nil {
		s.executorEtrog.CleanTrustedState()
	}
}

// GetCachedBatch implements syncinterfaces.SyncTrustedStateExecutor. Returns a cached batch
func (s *SyncTrustedStateExecutorSelector) GetCachedBatch(batchNumber uint64) *state.Batch {
	executor, _ := s.GetExecutor(batchNumber, 0)
	return executor.GetCachedBatch(min(batchNumber))
}
