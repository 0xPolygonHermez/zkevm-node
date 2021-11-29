package state

import (
	"github.com/ethereum/go-ethereum/common"
)

// Block struct
type Block struct {
	BlockNumber uint64
	BlockHash   common.Hash
	ParentHash  common.Hash
	Batches     []Batch
}
