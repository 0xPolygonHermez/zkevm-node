package txprofitabilitychecker

import (
	"fmt"
	"math/big"

	"github.com/hermeznetwork/hermez-core/encoding"
)

// TxProfitabilityCheckerType for different profitability checkers types
type TxProfitabilityCheckerType string

const (
	// ProfitabilityBase type that checks sum of costs of txs against min reward
	ProfitabilityBase = "base"
	// ProfitabilityAcceptAll validate batch anyway and don't check anything
	ProfitabilityAcceptAll = "acceptall"
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

// Config for the tx profitability checker configuration
type Config struct {
	TxProfitabilityCheckerType TxProfitabilityCheckerType `mapstructure:"TxProfitabilityCheckerType"`
	MinReward                  TokenAmountWithDecimals    `mapstructure:"MinReward"`
}
