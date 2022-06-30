package statev2

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hermeznetwork/hermez-core/log"
)

// L2Block represents a block on L2
type L2Block struct {
	BlockNumber  uint64
	Header       *types.Header
	Uncles       []*types.Header
	Transactions []*types.Transaction
	RawTxsData   []byte
	Receipts     []*types.Receipt
	ReceivedAt   time.Time
}

// NewL2BlockWithHeader creates a L2 block with the given header data.
func NewL2BlockWithHeader(header types.Header) *L2Block {
	return &L2Block{Header: &header}
}

// Hash returns the L2 block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (b *L2Block) Hash() common.Hash {
	return b.Header.Hash()
}

// Number is a helper function to get the L2 block number from the header
func (b *L2Block) Number() *big.Int {
	return b.Header.Number
}

// Size returns the true RLP encoded storage size of the L2 block, either by encoding
// and returning it, or returning a previously cached value.
func (b *L2Block) Size() common.StorageSize {
	c := writeCounter(0)
	err := rlp.Encode(&c, b)
	if err != nil {
		log.Errorf("failed to compute the Size of the block: %d", b.Number().Uint64())
	}
	return common.StorageSize(c)
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}
