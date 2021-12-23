package strategy

import (
	"math/big"
	"time"

	"github.com/hermeznetwork/hermez-core/encoding"
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
	coin := new(big.Float).SetInt(big.NewInt(encoding.TenToThePowerOf18))
	bigval := new(big.Float).Mul(amount, coin)
	result := new(big.Int)
	bigval.Int(result)
	t.Int = result

	return nil
}

// Type different types of strategy logic
type Type string

const (
	// AcceptAll strategy accepts all txs
	AcceptAll Type = "acceptall"
	// Base strategy that have basic selection algorithm and can accept different sorting algorithms and profitability checkers
	Base = "base"
)

// Strategy holds config params
type Strategy struct {
	Type                       Type                       `mapstructure:"Type"`
	TxSorterType               TxSorterType               `mapstructure:"TxSorterType"`
	TxProfitabilityCheckerType TxProfitabilityCheckerType `mapstructure:"TxProfitabilityCheckerType"`
	MinReward                  TokenAmountWithDecimals    `mapstructure:"MinReward"`
	PossibleTimeToSendTx       Duration                   `mapstructure:"PossibleTimeToSendTx"`
}
