package state

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Block struct
type Block struct {
	BlockNumber   uint64
	BlockHash     common.Hash
	ParentHash    common.Hash
	Batches       []Batch
	NewSequencers []Sequencer

	ReceivedAt time.Time
}

// NewBlock creates a block with the given data.
func NewBlock(blockNumber uint64) *Block {
	return &Block{BlockNumber: blockNumber}
}

// L2Block struct
type L2Block struct {
	BlockNumber  uint64
	TxHash       common.Hash
	ParentTxHash common.Hash

	ReceivedAt time.Time
}

// NewL2Block creates a L2 block with the given data.
func NewL2Block(blockNumber uint64) *L2Block {
	return &L2Block{BlockNumber: blockNumber}
}
