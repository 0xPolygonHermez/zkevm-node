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
	// ErrDuplicatedNonce is returned when adding a new tx to the worker and there is an existing tx
	// with the same nonce and higher gasPrice (in this case we keep the existing tx)
	ErrDuplicatedNonce = errors.New("duplicated nonce")
	// ErrReplacedTransaction is returned when an existing tx is replaced by a new tx with the same nonce and higher gasPrice
	ErrReplacedTransaction = errors.New("replaced transaction")
)
