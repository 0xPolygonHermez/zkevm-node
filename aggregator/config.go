package aggregator

import (
	"time"

	"github.com/hermeznetwork/hermez-core/etherman"
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

// Config represents the configuration of the aggregator
type Config struct {
	// IntervalToConsolidateState is the time the aggregator waits until
	// trying to consolidate a new state
	IntervalToConsolidateState Duration `env:"HERMEZCORE_AGGREGATOR_INTERVALTOCONSOLIDATESTATE"`

	// Etherman is the configuration required by etherman to interact with L1
	Etherman etherman.Config
}
