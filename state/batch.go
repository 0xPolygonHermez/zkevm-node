package state

import "github.com/ethereum/go-ethereum/core/types"

// Batch
type Batch struct {
	Number       uint64
	Transactions []types.Transaction
}
