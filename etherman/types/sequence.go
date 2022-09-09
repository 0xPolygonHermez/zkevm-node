package types

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Sequence represents an operation sent to the PoE smart contract to be
// processed.
type Sequence struct {
	GlobalExitRoot, StateRoot, LocalExitRoot common.Hash
	Timestamp                                int64
	Txs                                      []types.Transaction
}

// IsEmpty checks is sequence struct is empty
func (s Sequence) IsEmpty() bool {
	return reflect.DeepEqual(s, Sequence{})
}
