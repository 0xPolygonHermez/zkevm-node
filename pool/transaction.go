package pool

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// TxStatePending represents a tx that has not been processed
	TxStatePending TxState = "pending"
	// TxStateInvalid represents an invalid tx
	TxStateInvalid TxState = "invalid"
	// TxStateSelected represents a tx that has been selected
	TxStateSelected TxState = "selected"
)

// TxState represents the state of a tx
type TxState string

// String returns a representation of the tx state in a string format
func (s TxState) String() string {
	return string(s)
}

// Transaction represents a pool tx
type Transaction struct {
	types.Transaction
	State      TxState
	IsClaims   bool
	ReceivedAt time.Time
}

// IsClaimTx checks, if tx is a claim tx
func (tx *Transaction) IsClaimTx(l2GlobalExitRootManagerAddr common.Address) bool {
	if *tx.To() == l2GlobalExitRootManagerAddr &&
		strings.HasPrefix("0x"+common.Bytes2Hex(tx.Data()), bridgeClaimMethodSignature) {
		return true
	}
	return false
}
