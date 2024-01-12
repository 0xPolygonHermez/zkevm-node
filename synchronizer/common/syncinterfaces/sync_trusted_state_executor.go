package syncinterfaces

import (
	"context"
	"errors"
)

var (
	// ErrMissingSyncFromL1 is returned when we are behind expected L1 sync, so must be done from L1
	ErrMissingSyncFromL1 = errors.New("must sync from L1")
	// ErrFatalDesyncFromL1 is returned when trusted node and permissionless node have different data
	ErrFatalDesyncFromL1 = errors.New("fatal situation: the TrustedNode have another data!. Halt or do something")
)

// SyncTrustedStateExecutor is the interface that class that synchronize permissionless with a trusted node
type SyncTrustedStateExecutor interface {
	// SyncTrustedState syncs the trusted state with the permissionless state
	// if returns error ErrMissingSyncFromL1 then must force a L1 sync
	SyncTrustedState(ctx context.Context, latestSyncedBatch uint64) error
	// CleanTrustedState clean cache of Batches and StateRoot
	CleanTrustedState()
}
