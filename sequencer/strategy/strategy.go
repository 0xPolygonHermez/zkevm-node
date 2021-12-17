package strategy

import "time"

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

// Type different types of strategy logic
type Type string

const (
	AcceptAll Type = "acceptall"
	Base           = "base"
)

// Strategy holds config params
type Strategy struct {
	Type                       Type                       `mapstructure:"Type"`
	TxSorterType               TxSorterType               `mapstructure:"TxSorterType"`
	TxProfitabilityCheckerType TxProfitabilityCheckerType `mapstructure:"TxProfitabilityCheckerType"`
	MinReward                  uint64                     `mapstructure:"MinReward"`
	PossibleTimeToSendTx       Duration                   `mapstructure:"PossibleTimeToSendTx"`
}
