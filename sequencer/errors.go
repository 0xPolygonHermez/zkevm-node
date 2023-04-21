package sequencer

import "errors"

var (
	// ErrExpiredTransaction happens when the transaction is expired
	ErrExpiredTransaction = errors.New("transaction expired")
)
