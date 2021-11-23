package state

import (
	"encoding/binary"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/rlp"
)

var (
	// TODO: Calculate proper EmptyRootHash
	EmptyRootHash  = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	EmptyUncleHash = rlp.Hash([]*types.Header(nil))
)

// A BatchNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a batch.
type BatchNonce [8]byte

// EncodeNonce converts the given integer to a batch nonce.
func EncodeNonce(i uint64) BatchNonce {
	var n BatchNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a batch nonce.
func (n BatchNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BatchNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BatchNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BatchNonce", input, n[:])
}

// Hash returns the batch hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (b *Batch) Hash() common.Hash {
	return b.Header.Hash()
}

// EmptyBody returns true if there is no additional 'body' to complete the header
// that is: no transactions and no uncles.
func (b *Batch) EmptyBody() bool {
	return b.Header.TxHash == EmptyRootHash && b.Header.UncleHash == EmptyUncleHash
}

// EmptyReceipts returns true if there are no receipts for this batch.
func (b *Batch) EmptyReceipts() bool {
	return b.Header.ReceiptHash == EmptyRootHash
}

// Batch
type Batch struct {
	BatchNumber  uint64
	BlockNumber  uint64
	IsVirtual    bool
	Sequencer    common.Address
	Aggregator   common.Address
	Header       *types.Header
	Uncles       []*types.Header
	Transactions []*types.Transaction

	ReceivedAt time.Time
}

// NewBatchWithHeader creates a batch with the given header data.
func NewBatchWithHeader(header types.Header) *Batch {
	return &Batch{Header: &header}
}
