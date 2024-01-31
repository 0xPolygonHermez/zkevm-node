package l2_shared

import (
	"context"

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
		return s.executorPreEtrog, maximumBatchNumberToProcess
	}

	if latestSyncedBatch+1 >= fork.FromBatchNumber {
		return s.executorEtrog, maximumBatchNumberToProcess
	}
	return s.executorPreEtrog, min(maximumBatchNumberToProcess, fork.FromBatchNumber-1)
}

// SyncTrustedState syncs the trusted state with the permissionless state. In this case
// choose which executor must use
func (s *SyncTrustedStateExecutorSelector) SyncTrustedState(ctx context.Context, latestSyncedBatch uint64, maximumBatchNumberToProcess uint64) error {
	executor, maxBatchNumber := s.GetExecutor(latestSyncedBatch, maximumBatchNumberToProcess)
	if executor == nil {
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
func (s *SyncTrustedStateExecutorSelector) GetCachedBatch(batchNumber uint64) *state.Batch {
	executor, _ := s.GetExecutor(batchNumber, 0)
	return executor.GetCachedBatch(min(batchNumber))
}
