package sequencer

import (
	"errors"
)

var (
	// ErrExpiredTransaction happens when the transaction is expired
	ErrExpiredTransaction = errors.New("transaction expired")

	//// ErrBreakEvenGasPriceEmpty happens when the breakEven or gasPrice is nil or zero
	//ErrBreakEvenGasPriceEmpty = errors.New("breakEven and gasPrice cannot be nil or zero")
	//
	//// ErrEffectiveGasPricePercentageGreaterThanMaximum happens when the effective gas price percentage is greater than 255
	//ErrEffectiveGasPricePercentageGreaterThanMaximum = fmt.Errorf("effective gas price percentage is greater than %d", effectiveGasPriceMaxPercentageValue)
)
