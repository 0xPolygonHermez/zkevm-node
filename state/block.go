package state

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

// Block is ethereum block
type Block struct {
	BlockNumber uint64
	BlockHash   common.Hash
	ParentHash  common.Hash
	ReceivedAt  time.Time
	Batches     []Batch
}
