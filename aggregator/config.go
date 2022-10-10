package aggregator

import (
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
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
	IntervalToConsolidateState types.Duration `mapstructure:"IntervalToConsolidateState"`

	// IntervalFrequencyToGetProofGenerationState is the time the aggregator waits until
	// trying to get proof generation status, in case prover client returns PENDING state
	IntervalFrequencyToGetProofGenerationState types.Duration `mapstructure:"IntervalFrequencyToGetProofGenerationState"`

	// TxProfitabilityCheckerType type for checking is it profitable for aggregator to validate batch
	// possible values: base/acceptall
	TxProfitabilityCheckerType TxProfitabilityCheckerType `mapstructure:"TxProfitabilityCheckerType"`

	// TxProfitabilityMinReward min reward for base tx profitability checker when aggregator will validate batch
	// this parameter is used for the base tx profitability checker
	TxProfitabilityMinReward TokenAmountWithDecimals `mapstructure:"TxProfitabilityMinReward"`

	// IntervalAfterWhichBatchConsolidateAnyway this is interval for the main sequencer, that will check if there is no transactions
	IntervalAfterWhichBatchConsolidateAnyway types.Duration `mapstructure:"IntervalAfterWhichBatchConsolidateAnyway"`

	// ChainID is the L2 ChainID provided by the Network Config
	ChainID uint64

	// Array of Prover URIs read from configuration file
	ProverURIs []string
}
