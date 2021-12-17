package aggregator

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

// Config represents the configuration of the aggregator
type Config struct {
	// IntervalToConsolidateState is the time the aggregator waits until
	// trying to consolidate a new state
	IntervalToConsolidateState Duration `mapstructure:"IntervalToConsolidateState"`

	// TxProfitabilityCheckerType type for checking is it profitable for aggregator to validate batch
	TxProfitabilityCheckerType TxProfitabilityCheckerType `mapstructure:"TxProfitabilityCheckerType"`

	// TODO: understand, in which format matic collateral will be saved (10^18 or not)
	// TxProfitabilityMinReward min reward for base tx profitability checker when aggregator will validate batch
	TxProfitabilityMinReward uint64 `mapstructure:"TxProfitabilityMinReward"`
}
