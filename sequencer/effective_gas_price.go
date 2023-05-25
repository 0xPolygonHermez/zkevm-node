package sequencer

import (
	"errors"
	"math/big"
)

const (
	bits = 256
)

var (
	bitsBigInt           = big.NewInt(bits)
	hundredPercentInBits = big.NewInt(bits - 1)

	// ErrBreakEvenGasPriceEmpty happens when the breakEven or gasPrice is nil or zero
	ErrBreakEvenGasPriceEmpty = errors.New("breakEven and gasPrice cannot be nil or zero")
)

// CalcGasPriceEffectivePercentage calculates the gas price's effective percentage
func CalcGasPriceEffectivePercentage(breakEven *big.Int, gasPrice *big.Int) (*big.Int, error) {
	if breakEven == nil || gasPrice == nil ||
		gasPrice.Cmp(big.NewInt(0)) == 0 || breakEven.Cmp(big.NewInt(0)) == 0 {
		return nil, ErrBreakEvenGasPriceEmpty
	}

	if gasPrice.Cmp(breakEven) <= 0 {
		return hundredPercentInBits, nil
	}

	// Simulate Ceil with integer division
	b := new(big.Int).Mul(breakEven, bitsBigInt)
	b = b.Add(b, gasPrice)
	b = b.Sub(b, big.NewInt(1)) //nolint:gomnd
	b = b.Div(b, gasPrice)

	return b, nil
}
