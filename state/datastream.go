package state

import (
	"encoding/binary"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// StreamTypeSequencer represents a Sequencer stream
	StreamTypeSequencer datastreamer.StreamType = 1
	// EntryTypeL2BlockStart represents a L2 block start
	EntryTypeL2BlockStart datastreamer.EntryType = 1
	// EntryTypeL2Tx represents a L2 transaction
	EntryTypeL2Tx datastreamer.EntryType = 2
	// EntryTypeL2BlockEnd represents a L2 block end
	EntryTypeL2BlockEnd datastreamer.EntryType = 3
	// BookMarkTypeL2Block represents a L2 block bookmark
	BookMarkTypeL2Block byte = 0
)

// DSL2FullBlock represents a data stream L2 full block and its transactions
type DSL2FullBlock struct {
	L2Block DSL2Block
	Txs     []DSL2Transaction
}

// DSL2Block is a full l2 block
type DSL2Block struct {
	BatchNumber    uint64         // 8 bytes
	L2BlockNumber  uint64         // 8 bytes
	Timestamp      int64          // 8 bytes
	GlobalExitRoot common.Hash    // 32 bytes
	Coinbase       common.Address // 20 bytes
	ForkID         uint16         // 2 bytes
	BlockHash      common.Hash    // 32 bytes
	StateRoot      common.Hash    // 32 bytes
}

// DSL2BlockStart represents a data stream L2 block start
type DSL2BlockStart struct {
	BatchNumber    uint64         // 8 bytes
	L2BlockNumber  uint64         // 8 bytes
	Timestamp      int64          // 8 bytes
	GlobalExitRoot common.Hash    // 32 bytes
	Coinbase       common.Address // 20 bytes
	ForkID         uint16         // 2 bytes
}

// Encode returns the encoded DSL2BlockStart as a byte slice
func (b DSL2BlockStart) Encode() []byte {
	bytes := make([]byte, 0)
	bytes = binary.LittleEndian.AppendUint64(bytes, b.BatchNumber)
	bytes = binary.LittleEndian.AppendUint64(bytes, b.L2BlockNumber)
	bytes = binary.LittleEndian.AppendUint64(bytes, uint64(b.Timestamp))
	bytes = append(bytes, b.GlobalExitRoot.Bytes()...)
	bytes = append(bytes, b.Coinbase.Bytes()...)
	bytes = binary.LittleEndian.AppendUint16(bytes, b.ForkID)
	return bytes
}

// DSL2Transaction represents a data stream L2 transaction
type DSL2Transaction struct {
	EffectiveGasPricePercentage uint8  // 1 byte
	IsValid                     uint8  // 1 byte
	EncodedLength               uint32 // 4 bytes
	Encoded                     []byte
}

// Encode returns the encoded DSL2Transaction as a byte slice
func (l DSL2Transaction) Encode() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, byte(l.EffectiveGasPricePercentage))
	bytes = append(bytes, byte(l.IsValid))
	bytes = binary.LittleEndian.AppendUint32(bytes, l.EncodedLength)
	bytes = append(bytes, l.Encoded...)
	return bytes
}

// DSL2BlockEnd represents a L2 block end
type DSL2BlockEnd struct {
	L2BlockNumber uint64      // 8 bytes
	BlockHash     common.Hash // 32 bytes
	StateRoot     common.Hash // 32 bytes
}

// Encode returns the encoded DSL2BlockEnd as a byte slice
func (b DSL2BlockEnd) Encode() []byte {
	bytes := make([]byte, 0)
	bytes = binary.LittleEndian.AppendUint64(bytes, b.L2BlockNumber)
	bytes = append(bytes, b.BlockHash[:]...)
	bytes = append(bytes, b.StateRoot[:]...)
	return bytes
}

// DSBookMark represents a data stream bookmark
type DSBookMark struct {
	Type          byte
	L2BlockNumber uint64
}

// Encode returns the encoded DSBookMark as a byte slice
func (b DSBookMark) Encode() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, b.Type)
	bytes = binary.LittleEndian.AppendUint64(bytes, b.L2BlockNumber)
	return bytes
}
