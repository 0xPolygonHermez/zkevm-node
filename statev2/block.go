package statev2

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Block struct
type Block struct {
	BlockNumber     uint64
	BlockHash       common.Hash
	ParentHash      common.Hash
	GlobalExitRoots []GlobalExitRoot
	ForcedBatches   []ForcedBatch
	VerifyBatch     []VerifyBatch

	ReceivedAt time.Time
}

// NewBlock creates a block with the given data.
func NewBlock(blockNumber uint64) *Block {
	return &Block{BlockNumber: blockNumber}
}
