package pool

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

// TxState represents the state of a tx
type TxState string

const (
	// TxStatePending represents a tx that has not been processed
	TxStatePending TxState = "pending"
	// TxStateInvalid represents an invalid tx
	TxStateInvalid TxState = "invalid"
	// TxStateSelected represents a tx that has been selected
	TxStateSelected TxState = "selected"
)

// Transaction represents a pool tx
type Transaction struct {
	types.Transaction
	State      TxState
	ReceivedAt time.Time
}
