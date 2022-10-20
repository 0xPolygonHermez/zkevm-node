package state

import (
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// SequenceStatusPending represents a sequence that has not been mined yet on L1
	SequenceStatusPending SequenceStatus = "pending"
	// SequenceStatusConfirmed represents a sequence that has been mined and the state is now virtualized
	SequenceStatusConfirmed SequenceStatus = "confirmed"
)

// SequenceStatus represents the state of a tx
type SequenceStatus string

// Sequence represents an operation sent to the PoE smart contract to be
// processed.
type Sequence struct {
	BatchNumber                              uint64
	GlobalExitRoot, StateRoot, LocalExitRoot common.Hash
	Timestamp                                time.Time
	Txs                                      []types.Transaction

	Status     SequenceStatus
	L1Tx       *types.Transaction
	SentToL1At *time.Time
}

// IsEmpty checks is sequence struct is empty
func (s Sequence) IsEmpty() bool {
	return reflect.DeepEqual(s, Sequence{})
}
