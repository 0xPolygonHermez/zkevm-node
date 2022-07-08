package types

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Sequence represents an operation sent to the PoE smart contract to be
// processed.
type Sequence struct {
	GlobalExitRoot  common.Hash
	Timestamp       int64
	ForceBatchesNum uint64
	Txs             []types.Transaction
}

func (s Sequence) IsEmpty() bool {
	return reflect.DeepEqual(s, Sequence{})
}
