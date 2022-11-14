package state

import (
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// SequenceGroupStatusPending represents a sequence that has not been mined yet on L1
	SequenceGroupStatusPending SequenceGroupStatus = "pending"
	// SequenceGroupStatusConfirmed represents a sequence that has been mined and the state is now virtualized
	SequenceGroupStatusConfirmed SequenceGroupStatus = "confirmed"
)

// SequenceGroupStatus represents the state of a tx
type SequenceGroupStatus string

// Sequence represents an operation sent to the PoE smart contract to be
// processed.
type Sequence struct {
	BatchNumber      uint64
	StateRoot        common.Hash
	GlobalExitRoot   common.Hash
	LocalExitRoot    common.Hash
	Timestamp        time.Time
	Txs              []types.Transaction
	IsSequenceTooBig bool
}

// SequenceGroup is a struct used to control which sequences were sent
// in the same transaction to L1
type SequenceGroup struct {
	TxHash       common.Hash
	TxNonce      uint64
	FromBatchNum uint64
	ToBatchNum   uint64
	Status       SequenceGroupStatus
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

// IsEmpty checks is sequence struct is empty
func (s Sequence) IsEmpty() bool {
	return reflect.DeepEqual(s, Sequence{})
}
