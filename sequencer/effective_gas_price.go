package sequencer

import (
	"errors"
	"math/big"
)

const (
	bits = 256
)

var (
	bitsBigInt = big.NewInt(bits)

	// ErrBreakEvenGasPriceEmpty happens when the breakEven or gasPrice is nil or zero
	ErrBreakEvenGasPriceEmpty = errors.New("breakEven and gasPrice cannot be nil or zero")
	// ErrEffectiveGasPriceReprocess happens when the effective gas price requires reexecution
	ErrEffectiveGasPriceReprocess = errors.New("effective gas price requires reprocessing the transaction")
)

// CalculateEffectiveGasPricePercentage calculates the gas price's effective percentage
func CalculateEffectiveGasPricePercentage(gasPrice *big.Int, breakEven *big.Int) (uint8, error) {
	if breakEven == nil || gasPrice == nil ||
		gasPrice.Cmp(big.NewInt(0)) == 0 || breakEven.Cmp(big.NewInt(0)) == 0 {
		return 0, ErrBreakEvenGasPriceEmpty
	}

	// Simulate Ceil with integer division
	b := new(big.Int).Mul(breakEven, bitsBigInt)
	b = b.Add(b, gasPrice)
	b = b.Sub(b, big.NewInt(1)) //nolint:gomnd
	b = b.Div(b, gasPrice)

	return uint8(b.Uint64()), nil
}
