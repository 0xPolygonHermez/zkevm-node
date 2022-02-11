package state

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hermeznetwork/hermez-core/log"
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
	ChainID			   *big.Int
	GlobalExitRoot     common.Hash
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

// Size returns the true RLP encoded storage size of the batch, either by encoding
// and returning it, or returning a previsouly cached value.
func (b *Batch) Size() common.StorageSize {
	c := writeCounter(0)
	err := rlp.Encode(&c, b)
	if err != nil {
		log.Errorf("failed to compute the Size of the batch: %d", b.Number().Uint64())
	}
	return common.StorageSize(c)
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}
