package sequencer

import (
	"time"
	"unsafe"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// StreamTypeSequencer represents a Sequencer stream
	StreamTypeSequencer datastreamer.StreamType = 1
	// EntryTypeL2Block represents a L2 block
	EntryTypeL2Block datastreamer.EntryType = 1
	// EntryTypeL2Tx represents a L2 transaction
	EntryTypeL2Tx datastreamer.EntryType = 2
)

// L2Block represents a L2 block
type L2Block struct {
	BatchNumber    uint64
	L2BlockNumber  uint64
	Timestamp      time.Time
	GlobalExitRoot common.Hash
	Coinbase       common.Address
}

// Encode returns the encoded L2Block as a byte slice
func (b L2Block) Encode() []byte {
	const size = int(unsafe.Sizeof(L2Block{}))
	return (*(*[size]byte)(unsafe.Pointer(&b)))[:]
}

// L2Transaction represents a L2 transaction
type L2Transaction struct {
	BatchNumber                 uint64
	EffectiveGasPricePercentage uint8
	IsValid                     uint8
	EncodedLength               uint32
	Encoded                     []byte
}

// Encode returns the encoded L2Transaction as a byte slice
func (l L2Transaction) Encode() []byte {
	const size = int(unsafe.Sizeof(L2Transaction{}))
	return (*(*[size]byte)(unsafe.Pointer(&l)))[:]
}
