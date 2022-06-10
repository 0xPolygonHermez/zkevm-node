package aggregator

import (
	"fmt"
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/encoding"
	"math/big"
)

// TokenAmountWithDecimals is a wrapper type that parses token amount with decimals to big int
type TokenAmountWithDecimals struct {
	*big.Int `validate:"required"`
}

// UnmarshalText unmarshal token amount from float string to big int
func (t *TokenAmountWithDecimals) UnmarshalText(data []byte) error {
	amount, ok := new(big.Float).SetString(string(data))
	if !ok {
		return fmt.Errorf("failed to unmarshal string to float")
	}
	coin := new(big.Float).SetInt(big.NewInt(encoding.TenToThePowerOf18))
	bigval := new(big.Float).Mul(amount, coin)
	result := new(big.Int)
	bigval.Int(result)
	t.Int = result

	return nil
}

// Config represents the configuration of the aggregator
type Config struct {
	// IntervalToConsolidateState is the time the aggregator waits until
	// trying to consolidate a new state
	IntervalToConsolidateState config.Duration `mapstructure:"IntervalToConsolidateState"`

	// IntervalFrequencyToGetProofGenerationStateInSeconds is the time the aggregator waits until
	// trying to get proof generation status, in case prover client returns PENDING state
	IntervalFrequencyToGetProofGenerationStateInSeconds config.Duration `mapstructure:"IntervalFrequencyToGetProofGenerationStateInSeconds"`

	// TxProfitabilityCheckerType type for checking is it profitable for aggregator to validate batch
	// possible values: base/acceptall
	TxProfitabilityCheckerType TxProfitabilityCheckerType `mapstructure:"TxProfitabilityCheckerType"`

	// TxProfitabilityMinReward min reward for base tx profitability checker when aggregator will validate batch
	// this parameter is used for the base tx profitability checker
	TxProfitabilityMinReward TokenAmountWithDecimals `mapstructure:"TxProfitabilityMinReward"`

	// IntervalAfterWhichBatchConsolidateAnyway this is interval for the main sequencer, that will check if there is no transactions
	IntervalAfterWhichBatchConsolidateAnyway config.Duration `mapstructure:"IntervalAfterWhichBatchConsolidateAnyway"`
}
