package state

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Hash returns the batch hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (b *Batch) Hash() common.Hash {
	return b.Header.Hash()
}

// Batch represents a batch
type Batch struct {
	BlockNumber        uint64
	Sequencer          common.Address
	Aggregator         common.Address
	ConsolidatedTxHash common.Hash
	Header             *types.Header
	Uncles             []*types.Header
	Transactions       []*types.Transaction
	RawTxsData         []byte
	Receipts           []*Receipt
	MaticCollateral    *big.Int
	ReceivedAt         time.Time
	ConsolidatedAt     *time.Time
}

// NewBatchWithHeader creates a batch with the given header data.
func NewBatchWithHeader(header types.Header) *Batch {
	return &Batch{Header: &header}
}

// Number is a helper function to get the batch number from the header
func (b *Batch) Number() *big.Int {
	return b.Header.Number
}
