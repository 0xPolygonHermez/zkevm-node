package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Transaction represents a state tx
type Transaction struct {
	Hash     common.Hash
	BatchNum uint64
	From     common.Address
	types.Transaction
}
