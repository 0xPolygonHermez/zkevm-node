package sequencer

import (
	"time"
)

// Duration is a wrapper type that parses time duration from text.
type Duration struct {
	time.Duration `validate:"required"`
}

// UnmarshalText unmarshalls time duration from text.
func (d *Duration) UnmarshalText(data []byte) error {
	duration, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	d.Duration = duration
	return nil
}

// Config represents the configuration of a sequencer
type Config struct {
	// IntervalToProposeBatch is the time the sequencer waits until
	// trying to propose a batch
	IntervalToProposeBatch Duration `mapstructure:"IntervalToProposeBatch"`

	// SyncedBlockDif is the difference, how many block left to sync. So if sequencer see, that
	// X amount of blocks are left to sync, it will start to select txs
	SyncedBlockDif uint64 `mapstructure:"SyncedBlockDif"`
}
