package statev2

import (
	"github.com/ethereum/go-ethereum/common"
)

// VirtualBatch represents a batch of the virtual state.
type VirtualBatch struct {
	BatchNumber uint64
	Sequencer   common.Address
	TxHash      common.Hash
	BlockNumber uint64
}
