package statev2

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Batch struct
type Batch struct {
	BatchNumber    uint64
	Coinbase       common.Address
	BatchL2Data    []byte
	StateRoot      common.Hash
	LocalExitRoot  common.Hash
	Timestamp      time.Time
	Transactions   []types.Transaction
	GlobalExitRoot common.Hash
}

// ProcessingContext is the necessary data that a batch needs to porvide to the runtime,
// without the historical state data (processing receipt from previous batch)
type ProcessingContext struct {
	BatchNumber    uint64
	Coinbase       common.Address
	Timestamp      time.Time
	GlobalExitRoot common.Hash
}

// ProcessingReceipt indicates the outcome (StateRoot, LocalExitRoot) of processing a batch
type ProcessingReceipt struct {
	BatchNumber   uint64
	StateRoot     common.Hash
	LocalExitRoot common.Hash
}

// VerifyBatch represents a VerifyBatch
type VerifiedBatch struct {
	BlockNumber uint64
	BatchNumber uint64
	Aggregator  common.Address
	TxHash      common.Hash
}

// VirtualBatch represents a VirtualBatch
type VirtualBatch struct {
	BatchNumber uint64
	TxHash      common.Hash
	Sequencer   common.Address
	BlockNumber uint64
}
