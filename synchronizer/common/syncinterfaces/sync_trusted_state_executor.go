package syncinterfaces

import "context"

// SyncTrustedStateExecutor is the interface that class that synchronize permissionless with a trusted node
type SyncTrustedStateExecutor interface {
	SyncTrustedState(ctx context.Context, latestSyncedBatch uint64) error
	CleanTrustedState()
}
