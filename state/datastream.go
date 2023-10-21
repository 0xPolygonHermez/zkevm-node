package state

import (
	"context"
	"encoding/binary"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

const (
	// StreamTypeSequencer represents a Sequencer stream
	StreamTypeSequencer datastreamer.StreamType = 1
	// EntryTypeBookMark represents a bookmark entry
	EntryTypeBookMark datastreamer.EntryType = datastreamer.EtBookmark
	// EntryTypeL2BlockStart represents a L2 block start
	EntryTypeL2BlockStart datastreamer.EntryType = 1
	// EntryTypeL2Tx represents a L2 transaction
	EntryTypeL2Tx datastreamer.EntryType = 2
	// EntryTypeL2BlockEnd represents a L2 block end
	EntryTypeL2BlockEnd datastreamer.EntryType = 3
	// EntryTypeUpdateGER represents a GER update
	EntryTypeUpdateGER datastreamer.EntryType = 4
	// BookMarkTypeL2Block represents a L2 block bookmark
	BookMarkTypeL2Block byte = 0
)

// DSBatch represents a data stream batch
type DSBatch struct {
	Batch
	ForkID uint16
}

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

// Decode decodes the DSL2BlockStart from a byte slice
func (b DSL2BlockStart) Decode(data []byte) DSL2BlockStart {
	b.BatchNumber = binary.LittleEndian.Uint64(data[0:8])
	b.L2BlockNumber = binary.LittleEndian.Uint64(data[8:16])
	b.Timestamp = int64(binary.LittleEndian.Uint64(data[16:24]))
	b.GlobalExitRoot = common.BytesToHash(data[24:56])
	b.Coinbase = common.BytesToAddress(data[56:76])
	b.ForkID = binary.LittleEndian.Uint16(data[76:78])
	return b
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

// Decode decodes the DSL2Transaction from a byte slice
func (l DSL2Transaction) Decode(data []byte) DSL2Transaction {
	l.EffectiveGasPricePercentage = uint8(data[0])
	l.IsValid = uint8(data[1])
	l.EncodedLength = binary.LittleEndian.Uint32(data[2:6])
	l.Encoded = data[6:]
	return l
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

// Decode decodes the DSL2BlockEnd from a byte slice
func (b DSL2BlockEnd) Decode(data []byte) DSL2BlockEnd {
	b.L2BlockNumber = binary.LittleEndian.Uint64(data[0:8])
	b.BlockHash = common.BytesToHash(data[8:40])
	b.StateRoot = common.BytesToHash(data[40:72])
	return b
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

// Decode decodes the DSBookMark from a byte slice
func (b DSBookMark) Decode(data []byte) DSBookMark {
	b.Type = byte(data[0])
	b.L2BlockNumber = binary.LittleEndian.Uint64(data[1:9])
	return b
}

// DSUpdateGER represents a data stream GER update
type DSUpdateGER struct {
	BatchNumber    uint64         // 8 bytes
	Timestamp      int64          // 8 bytes
	GlobalExitRoot common.Hash    // 32 bytes
	Coinbase       common.Address // 20 bytes
	ForkID         uint16         // 2 bytes
	StateRoot      common.Hash    // 32 bytes
}

// Encode returns the encoded DSUpdateGER as a byte slice
func (g DSUpdateGER) Encode() []byte {
	bytes := make([]byte, 0)
	bytes = binary.LittleEndian.AppendUint64(bytes, g.BatchNumber)
	bytes = binary.LittleEndian.AppendUint64(bytes, uint64(g.Timestamp))
	bytes = append(bytes, g.GlobalExitRoot[:]...)
	bytes = append(bytes, g.Coinbase[:]...)
	bytes = binary.LittleEndian.AppendUint16(bytes, g.ForkID)
	bytes = append(bytes, g.StateRoot[:]...)
	return bytes
}

// Decode decodes the DSUpdateGER from a byte slice
func (g DSUpdateGER) Decode(data []byte) DSUpdateGER {
	g.BatchNumber = binary.LittleEndian.Uint64(data[0:8])
	g.Timestamp = int64(binary.LittleEndian.Uint64(data[8:16]))
	g.GlobalExitRoot = common.BytesToHash(data[16:48])
	g.Coinbase = common.BytesToAddress(data[48:68])
	g.ForkID = binary.LittleEndian.Uint16(data[68:70])
	g.StateRoot = common.BytesToHash(data[70:102])
	return g
}

// DSState gathers the methods required to interact with the data stream state.
type DSState interface {
	GetDSGenesisBlock(ctx context.Context, dbTx pgx.Tx) (*DSL2Block, error)
	GetDSBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*DSBatch, error)
	GetDSL2Blocks(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]*DSL2Block, error)
	GetDSL2Transactions(ctx context.Context, minL2Block, maxL2Block uint64, dbTx pgx.Tx) ([]*DSL2Transaction, error)
}

// GenerateDataStreamerFile generates or resumes a data stream file
func GenerateDataStreamerFile(ctx context.Context, streamServer *datastreamer.StreamServer, stateDB DSState) error {
	header := streamServer.GetHeader()

	var currentBatchNumber uint64 = 0
	var currentL2Block uint64 = 0

	if header.TotalEntries == 0 {
		// Get Genesis block
		genesisL2Block, err := stateDB.GetDSGenesisBlock(ctx, nil)
		if err != nil {
			return err
		}

		err = streamServer.StartAtomicOp()
		if err != nil {
			return err
		}

		bookMark := DSBookMark{
			Type:          BookMarkTypeL2Block,
			L2BlockNumber: genesisL2Block.L2BlockNumber,
		}

		_, err = streamServer.AddStreamBookmark(bookMark.Encode())
		if err != nil {
			return err
		}

		genesisBlock := DSL2BlockStart{
			BatchNumber:    genesisL2Block.BatchNumber,
			L2BlockNumber:  genesisL2Block.L2BlockNumber,
			Timestamp:      genesisL2Block.Timestamp,
			GlobalExitRoot: genesisL2Block.GlobalExitRoot,
			Coinbase:       genesisL2Block.Coinbase,
			ForkID:         genesisL2Block.ForkID,
		}

		log.Infof("Genesis block: %+v", genesisBlock)

		_, err = streamServer.AddStreamEntry(1, genesisBlock.Encode())
		if err != nil {
			return err
		}

		genesisBlockEnd := DSL2BlockEnd{
			L2BlockNumber: genesisL2Block.L2BlockNumber,
			BlockHash:     genesisL2Block.BlockHash,
			StateRoot:     genesisL2Block.StateRoot,
		}

		_, err = streamServer.AddStreamEntry(EntryTypeL2BlockEnd, genesisBlockEnd.Encode())
		if err != nil {
			return err
		}

		err = streamServer.CommitAtomicOp()
		if err != nil {
			return err
		}
	} else {
		latestEntry, err := streamServer.GetEntry(header.TotalEntries - 1)
		if err != nil {
			return err
		}

		log.Infof("Latest entry: %+v", latestEntry)

		switch latestEntry.Type {
		case EntryTypeUpdateGER:
			log.Info("Latest entry type is UpdateGER")
			currentBatchNumber = binary.LittleEndian.Uint64(latestEntry.Data[0:8])
		case EntryTypeL2BlockEnd:
			log.Info("Latest entry type is L2BlockEnd")
			currentL2Block = binary.LittleEndian.Uint64(latestEntry.Data[0:8])

			bookMark := DSBookMark{
				Type:          BookMarkTypeL2Block,
				L2BlockNumber: currentL2Block,
			}

			firstEntry, err := streamServer.GetFirstEventAfterBookmark(bookMark.Encode())
			if err != nil {
				return err
			}
			currentBatchNumber = binary.LittleEndian.Uint64(firstEntry.Data[0:8])
		}
	}

	log.Infof("Current Batch number: %d", currentBatchNumber)
	log.Infof("Current L2 block number: %d", currentL2Block)

	var entry uint64 = header.TotalEntries
	var l2blocks []*DSL2Block
	var currentGER = common.Hash{}

	if entry > 0 {
		entry--
	}

	// Start on the current batch number + 1
	currentBatchNumber++

	var err error

	for err == nil {
		log.Infof("Current entry number: %d", entry)
		// Get Next Batch
		batch, err := stateDB.GetDSBatch(ctx, currentBatchNumber, nil)
		if err != nil {
			if err == ErrStateNotSynchronized {
				break
			}
			log.Errorf("Error getting batch %d: %s", currentBatchNumber, err.Error())
			return err
		}

		l2blocks, err = stateDB.GetDSL2Blocks(ctx, currentBatchNumber, nil)
		if err != nil {
			log.Errorf("Error getting L2 blocks for batch %d: %s", currentBatchNumber, err.Error())
			return err
		}

		currentBatchNumber++

		if len(l2blocks) == 0 {
			// Empty batch
			// Check if there is a GER update
			if batch.GlobalExitRoot != currentGER && batch.GlobalExitRoot != (common.Hash{}) {
				updateGer := DSUpdateGER{
					BatchNumber:    batch.BatchNumber,
					Timestamp:      batch.Timestamp.Unix(),
					GlobalExitRoot: batch.GlobalExitRoot,
					Coinbase:       batch.Coinbase,
					ForkID:         batch.ForkID,
					StateRoot:      batch.StateRoot,
				}

				err = streamServer.StartAtomicOp()
				if err != nil {
					return err
				}

				entry, err = streamServer.AddStreamEntry(EntryTypeUpdateGER, updateGer.Encode())
				if err != nil {
					return err
				}

				err = streamServer.CommitAtomicOp()
				if err != nil {
					return err
				}

				currentGER = batch.GlobalExitRoot
			}
			continue
		}

		// Get transactions for all the retrieved l2 blocks
		l2Transactions, err := stateDB.GetDSL2Transactions(ctx, l2blocks[0].L2BlockNumber, l2blocks[len(l2blocks)-1].L2BlockNumber, nil)
		if err != nil {
			return err
		}

		err = streamServer.StartAtomicOp()
		if err != nil {
			return err
		}

		for x, l2block := range l2blocks {
			blockStart := DSL2BlockStart{
				BatchNumber:    l2block.BatchNumber,
				L2BlockNumber:  l2block.L2BlockNumber,
				Timestamp:      l2block.Timestamp,
				GlobalExitRoot: l2block.GlobalExitRoot,
				Coinbase:       l2block.Coinbase,
				ForkID:         l2block.ForkID,
			}

			bookMark := DSBookMark{
				Type:          BookMarkTypeL2Block,
				L2BlockNumber: blockStart.L2BlockNumber,
			}

			_, err = streamServer.AddStreamBookmark(bookMark.Encode())
			if err != nil {
				return err
			}

			_, err = streamServer.AddStreamEntry(EntryTypeL2BlockStart, blockStart.Encode())
			if err != nil {
				return err
			}

			entry, err = streamServer.AddStreamEntry(EntryTypeL2Tx, l2Transactions[x].Encode())
			if err != nil {
				return err
			}

			blockEnd := DSL2BlockEnd{
				L2BlockNumber: l2block.L2BlockNumber,
				BlockHash:     l2block.BlockHash,
				StateRoot:     l2block.StateRoot,
			}

			_, err = streamServer.AddStreamEntry(EntryTypeL2BlockEnd, blockEnd.Encode())
			if err != nil {
				return err
			}
			currentGER = l2block.GlobalExitRoot
		}
		// Commit at the end of each batch
		err = streamServer.CommitAtomicOp()
		if err != nil {
			return err
		}
	}

	return err
}
