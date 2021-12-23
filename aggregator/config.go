package aggregator

import (
	"math/big"
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

type TokenAmountWithDecimals struct {
	*big.Int `validate:"required"`
}

func (t *TokenAmountWithDecimals) UnmarshalText(data []byte) error {
	amount, _ := new(big.Float).SetString(string(data))

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))
	bigval := new(big.Float)
	bigval.Mul(amount, coin)
	result := new(big.Int)
	bigval.Int(result)
	t.Int = result

	return nil
}

// Config represents the configuration of the aggregator
type Config struct {
	// IntervalToConsolidateState is the time the aggregator waits until
	// trying to consolidate a new state
	IntervalToConsolidateState Duration `mapstructure:"IntervalToConsolidateState"`

	// TxProfitabilityCheckerType type for checking is it profitable for aggregator to validate batch
	// possible values: base and acceptall
	TxProfitabilityCheckerType TxProfitabilityCheckerType `mapstructure:"TxProfitabilityCheckerType"`

	// TODO: understand, in which format matic collateral will be saved (10^18 or not)
	// TxProfitabilityMinReward min reward for base tx profitability checker when aggregator will validate batch
	// this parameter is used for the base tx profitability checker
	TxProfitabilityMinReward TokenAmountWithDecimals `mapstructure:"TxProfitabilityMinReward"`
}
