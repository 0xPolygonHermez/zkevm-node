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
