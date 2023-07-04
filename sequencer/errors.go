package sequencer

import "errors"

var (
	// ErrExpiredTransaction happens when the transaction is expired
	ErrExpiredTransaction = errors.New("transaction expired")
	// ErrBreakEvenGasPriceEmpty happens when the breakEven or gasPrice is nil or zero
	ErrBreakEvenGasPriceEmpty = errors.New("breakEven and gasPrice cannot be nil or zero")
	// ErrEffectiveGasPriceReprocess happens when the effective gas price requires reexecution
	ErrEffectiveGasPriceReprocess = errors.New("effective gas price requires reprocessing the transaction")
	// ErrZeroL1GasPrice is returned if the L1 gas price is 0.
	ErrZeroL1GasPrice = errors.New("L1 gas price 0")
)
