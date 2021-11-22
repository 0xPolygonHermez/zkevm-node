package state

import (
	"github.com/ethereum/go-ethereum/common"
)
// Block
type Block struct {
	BlockNum  uint64
	BlockHash common.Hash
	Batches   []Batch
}
