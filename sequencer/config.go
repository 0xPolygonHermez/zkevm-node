package sequencer

import (
	"time"

	"github.com/hermeznetwork/hermez-core/sequencer/strategy"
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

	// IntervalAfterWhichBatchSentAnyway this is interval for the main sequencer, that will check if there is no transactions
	IntervalAfterWhichBatchSentAnyway Duration `mapstructure:"IntervalAfterWhichBatchSentAnyway"`

	// Strategy is the configuration for the strategy
	Strategy strategy.Strategy `mapstructure:"Strategy"`

	// AllowNonRegistered determines if the sequencer will run using the default
	// chain ID
	AllowNonRegistered bool `mapstructure:"AllowNonRegistered"`

	// DefaultChainID is the common ChainID to all the sequencers
	DefaultChainID uint64 `mapstructure:"DefaultChainID"`
}
