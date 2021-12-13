package sequencer

import (
	"time"
)

// Config represents the configuration of a sequencer
type Config struct {
	// IntervalToProposeBatch is the time the sequencer waits until
	// trying to propose a batch
	IntervalToProposeBatch time.Duration

	// SyncedBlockDif is the difference, how many block left to sync. So if sequencer see, that
	// X amount of blocks are left to sync, it will start to select txs
	SyncedBlockDif uint64
}
