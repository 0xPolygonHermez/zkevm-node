package sequencerv2

import (
	"github.com/hermeznetwork/hermez-core/config"
)

// Config represents the configuration of a sequencer
type Config struct {
	// WaitPeriodPoolIsEmpty is the time the sequencer waits until
	// trying to propose a batch
	WaitPeriodPoolIsEmpty config.Duration `mapstructure:"IntervalToProposeBatch"`
}
