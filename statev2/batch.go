package statev2

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Batch struct
type Batch struct {
	BatchNumber       uint64
	Coinbase          common.Address
	BatchL2Data       []byte
	OldStateRoot      common.Hash
	GlobalExitRootNum *big.Int
	OldLocalExitRoot  common.Hash
	Timestamp         time.Time
	Transactions      []types.Transaction
	GlobalExitRoot    common.Hash
}

// VerifyBatch represents a VerifyBatch
type VerifiedBatch struct {
	BlockNumber uint64
	BatchNumber uint64
	Aggregator  common.Address
	TxHash      common.Hash
}
