package statev2

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Batch struct
type Batch struct {
	BatchNum          uint64
	Coinbase          common.Address
	BatchL2Data       []byte
	OldStateRoot      common.Hash
	GlobalExitRootNum *big.Int
	OldLocalExitRoot  common.Hash
	EthTimestamp      time.Time
	Transactions      []types.Transaction
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
