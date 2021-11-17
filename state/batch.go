package state

import "github.com/ethereum/go-ethereum/core/types"

// Batch
type Batch struct {
	Txs []types.Transaction
}
