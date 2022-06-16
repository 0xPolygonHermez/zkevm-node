package statev2

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Batch represents a batch of the trusted state.
type Batch struct {
	BatchNumber    uint64
	GlobalExitRoot common.Hash
	Transactions   []*types.Transaction
	RawTxsData     []byte
	Timestamp      time.Time
}

// VirtualBatch represents a batch of the virtual state.
type VirtualBatch struct {
	BatchNumber uint64
	Sequencer   common.Address
	TxHash      common.Hash
	BlockNumber uint64
}
