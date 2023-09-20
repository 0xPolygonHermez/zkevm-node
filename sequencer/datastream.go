package sequencer

import (
	"encoding/binary"

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

// DSL2Block represents a data stream L2 block
type DSL2Block struct {
	BatchNumber    uint64
	L2BlockNumber  uint64
	Timestamp      uint64
	GlobalExitRoot common.Hash
	Coinbase       common.Address
}

// Encode returns the encoded L2Block as a byte slice
func (b DSL2Block) Encode() []byte {
	bytes := make([]byte, 0)
	bytes = binary.LittleEndian.AppendUint64(bytes, b.BatchNumber)
	bytes = binary.LittleEndian.AppendUint64(bytes, b.L2BlockNumber)
	bytes = binary.LittleEndian.AppendUint64(bytes, uint64(b.Timestamp))
	bytes = append(bytes, b.GlobalExitRoot.Bytes()...)
	bytes = append(bytes, b.Coinbase.Bytes()...)
	return bytes
}

// DSL2Transaction represents a data stream L2 transaction
type DSL2Transaction struct {
	BatchNumber                 uint64
	EffectiveGasPricePercentage uint8
	IsValid                     uint8
	EncodedLength               uint32
	Encoded                     []byte
}

// Encode returns the encoded L2Transaction as a byte slice
func (l DSL2Transaction) Encode() []byte {
	bytes := make([]byte, 0)
	bytes = binary.LittleEndian.AppendUint64(bytes, l.BatchNumber)
	bytes = append(bytes, byte(l.EffectiveGasPricePercentage))
	bytes = append(bytes, byte(l.IsValid))
	bytes = binary.LittleEndian.AppendUint32(bytes, l.EncodedLength)
	bytes = append(bytes, l.Encoded...)
	return bytes
}
