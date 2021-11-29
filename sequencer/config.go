package sequencer

import (
	"time"

	"github.com/hermeznetwork/hermez-core/etherman"
)

// Config represents the configuration of a sequencer
type Config struct {
	// IntervalToProposeBatch is the time the sequencer waits until
	// trying to propose a batch
	IntervalToProposeBatch time.Duration

	// Etherman is the configuration required by etherman to interact with L1
	Etherman etherman.Config
}
