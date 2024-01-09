package state

import (
	"context"
	"encoding/binary"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/keccak256"
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
	// SystemSC is the system smart contract address
	SystemSC = "0x000000000000000000000000000000005ca1ab1e"
	// posConstant is the constant used to compute the position of the intermediate state root
	posConstant = 1
)

// DSBatch represents a data stream batch
type DSBatch struct {
	Batch
	ForkID uint16
}

// DSFullBatch represents a data stream batch ant its L2 blocks
type DSFullBatch struct {
	DSBatch
	L2Blocks []DSL2FullBlock
}

// DSL2FullBlock represents a data stream L2 full block and its transactions
type DSL2FullBlock struct {
	DSL2Block
	Txs []DSL2Transaction
}

// DSL2Block is a full l2 block
type DSL2Block struct {
	BatchNumber    uint64         // 8 bytes
	L2BlockNumber  uint64         // 8 bytes
	Timestamp      int64          // 8 bytes
	L1BlockHash    common.Hash    // 32 bytes
	GlobalExitRoot common.Hash    // 32 bytes
	L1InfoRoot     common.Hash    // 32 bytes
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
	L1BlockHash    common.Hash    // 32 bytes
	GlobalExitRoot common.Hash    // 32 bytes
	L1InfoRoot     common.Hash    // 32 bytes
	Coinbase       common.Address // 20 bytes
	ForkID         uint16         // 2 bytes
}

// Encode returns the encoded DSL2BlockStart as a byte slice
func (b DSL2BlockStart) Encode() []byte {
	bytes := make([]byte, 0)
	bytes = binary.LittleEndian.AppendUint64(bytes, b.BatchNumber)
	bytes = binary.LittleEndian.AppendUint64(bytes, b.L2BlockNumber)
	bytes = binary.LittleEndian.AppendUint64(bytes, uint64(b.Timestamp))
	bytes = append(bytes, b.L1BlockHash.Bytes()...)
	bytes = append(bytes, b.GlobalExitRoot.Bytes()...)
	bytes = append(bytes, b.L1InfoRoot.Bytes()...)
	bytes = append(bytes, b.Coinbase.Bytes()...)
	bytes = binary.LittleEndian.AppendUint16(bytes, b.ForkID)
	return bytes
}

// Decode decodes the DSL2BlockStart from a byte slice
func (b DSL2BlockStart) Decode(data []byte) DSL2BlockStart {
	b.BatchNumber = binary.LittleEndian.Uint64(data[0:8])
	b.L2BlockNumber = binary.LittleEndian.Uint64(data[8:16])
	b.Timestamp = int64(binary.LittleEndian.Uint64(data[16:24]))
	b.L1BlockHash = common.BytesToHash(data[24:56])
	b.GlobalExitRoot = common.BytesToHash(data[56:88])
	b.L1InfoRoot = common.BytesToHash(data[88:120])
	b.Coinbase = common.BytesToAddress(data[120:140])
	b.ForkID = binary.LittleEndian.Uint16(data[140:142])
	return b
}

// DSL2Transaction represents a data stream L2 transaction
type DSL2Transaction struct {
	L2BlockNumber               uint64      // Not included in the encoded data
	EffectiveGasPricePercentage uint8       // 1 byte
	IsValid                     uint8       // 1 byte
	StateRoot                   common.Hash // 32 bytes
	EncodedLength               uint32      // 4 bytes
	Encoded                     []byte
}

// Encode returns the encoded DSL2Transaction as a byte slice
func (l DSL2Transaction) Encode() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, l.EffectiveGasPricePercentage)
	bytes = append(bytes, l.IsValid)
	bytes = append(bytes, l.StateRoot[:]...)
	bytes = binary.LittleEndian.AppendUint32(bytes, l.EncodedLength)
	bytes = append(bytes, l.Encoded...)
	return bytes
}

// Decode decodes the DSL2Transaction from a byte slice
func (l DSL2Transaction) Decode(data []byte) DSL2Transaction {
	l.EffectiveGasPricePercentage = data[0]
	l.IsValid = data[1]
	l.StateRoot = common.BytesToHash(data[2:34])
	l.EncodedLength = binary.LittleEndian.Uint32(data[34:38])
	l.Encoded = data[38:]
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
	b.Type = data[0]
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
	GetDSBatches(ctx context.Context, firstBatchNumber, lastBatchNumber uint64, readWIPBatch bool, dbTx pgx.Tx) ([]*DSBatch, error)
	GetDSL2Blocks(ctx context.Context, firstBatchNumber, lastBatchNumber uint64, dbTx pgx.Tx) ([]*DSL2Block, error)
	GetDSL2Transactions(ctx context.Context, firstL2Block, lastL2Block uint64, dbTx pgx.Tx) ([]*DSL2Transaction, error)
	GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root common.Hash) (*big.Int, error)
	GetLastL2BlockHeader(ctx context.Context, dbTx pgx.Tx) (*L2Header, error)
	GetVirtualBatchParentHash(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (common.Hash, error)
	GetForcedBatchParentHash(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (common.Hash, error)
	GetL1InfoRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (L1InfoTreeExitRootStorageEntry, error)
}

// GenerateDataStreamerFile generates or resumes a data stream file
func GenerateDataStreamerFile(ctx context.Context, streamServer *datastreamer.StreamServer, stateDB DSState, readWIPBatch bool, imStateRoots *map[uint64][]byte) error {
	header := streamServer.GetHeader()

	var currentBatchNumber uint64 = 0
	var currentL2Block uint64 = 0
	var lastAddedL2Block uint64 = 0

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
			L1InfoRoot:     genesisL2Block.L1InfoRoot,
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
	var currentGER = common.Hash{}

	if entry > 0 {
		entry--
	}

	// Start on the current batch number + 1
	currentBatchNumber++
	var err error
	const limit = 10000

	for err == nil {
		log.Debugf("Current entry number: %d", entry)
		log.Debugf("Current batch number: %d", currentBatchNumber)

		// Get Next Batch
		batches, err := stateDB.GetDSBatches(ctx, currentBatchNumber, currentBatchNumber+limit, readWIPBatch, nil)
		if err != nil {
			if err == ErrStateNotSynchronized {
				break
			}
			log.Errorf("Error getting batch %d: %s", currentBatchNumber, err.Error())
			return err
		}

		// Finished?
		if len(batches) == 0 {
			break
		}

		l2Blocks, err := stateDB.GetDSL2Blocks(ctx, batches[0].BatchNumber, batches[len(batches)-1].BatchNumber, nil)
		if err != nil {
			log.Errorf("Error getting L2 blocks for batches starting at %d: %s", batches[0].BatchNumber, err.Error())
			return err
		}

		l2Txs := make([]*DSL2Transaction, 0)
		if len(l2Blocks) > 0 {
			l2Txs, err = stateDB.GetDSL2Transactions(ctx, l2Blocks[0].L2BlockNumber, l2Blocks[len(l2Blocks)-1].L2BlockNumber, nil)
			if err != nil {
				log.Errorf("Error getting L2 transactions for blocks starting at %d: %s", l2Blocks[0].L2BlockNumber, err.Error())
				return err
			}
		}

		// Generate full batches
		fullBatches := computeFullBatches(batches, l2Blocks, l2Txs)
		currentBatchNumber += limit

		for _, batch := range fullBatches {
			if len(batch.L2Blocks) == 0 {
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

			err = streamServer.StartAtomicOp()
			if err != nil {
				return err
			}

			for blockIndex, l2block := range batch.L2Blocks {
				if l2block.L2BlockNumber <= lastAddedL2Block && lastAddedL2Block != 0 {
					continue
				} else {
					lastAddedL2Block = l2block.L2BlockNumber
				}

				l1BlockHash := common.Hash{}

				// Get L1 block hash
				if l2block.ForkID >= FORKID_ETROG {
					isForcedBatch := false
					batchRawData := &BatchRawV2{}

					if batch.BatchNumber == 1 || batch.ForcedBatchNum != nil {
						isForcedBatch = true
					} else {
						batchRawData, err = DecodeBatchV2(batch.BatchL2Data)
						if err != nil {
							log.Errorf("Failed to decode batch data, err: %v", err)
							return err
						}
					}

					if !isForcedBatch {
						// Get current block by index
						l2blockRaw := batchRawData.Blocks[blockIndex]
						if l2block.L2BlockNumber == 1101 {
							log.Fatalf("Lo2BlockRaw: %+v", l2blockRaw)
						}
						if l2blockRaw.IndexL1InfoTree != 0 {
							l1InfoTreeExitRootStorageEntry, err := stateDB.GetL1InfoRootLeafByIndex(ctx, l2blockRaw.IndexL1InfoTree, nil)
							if err != nil {
								return err
							}
							l1BlockHash = l1InfoTreeExitRootStorageEntry.L1InfoTreeLeaf.PreviousBlockHash
						}
					} else {
						// Initial batch must be handled differently
						if batch.BatchNumber == 1 {
							l1BlockHash, err = stateDB.GetVirtualBatchParentHash(ctx, batch.BatchNumber, nil)
							if err != nil {
								return err
							}
						} else {
							l1BlockHash, err = stateDB.GetForcedBatchParentHash(ctx, *batch.ForcedBatchNum, nil)
							if err != nil {
								return err
							}
						}
					}
				}

				blockStart := DSL2BlockStart{
					BatchNumber:    l2block.BatchNumber,
					L2BlockNumber:  l2block.L2BlockNumber,
					Timestamp:      l2block.Timestamp,
					L1BlockHash:    l1BlockHash,
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

				for _, tx := range l2block.Txs {
					// Populate intermediate state root
					if imStateRoots == nil || (*imStateRoots)[blockStart.L2BlockNumber] == nil {
						position := GetSystemSCPosition(l2block.L2BlockNumber)
						imStateRoot, err := stateDB.GetStorageAt(ctx, common.HexToAddress(SystemSC), big.NewInt(0).SetBytes(position), l2block.StateRoot)
						if err != nil {
							return err
						}
						tx.StateRoot = common.BigToHash(imStateRoot)
					} else {
						tx.StateRoot = common.BytesToHash((*imStateRoots)[blockStart.L2BlockNumber])
					}

					entry, err = streamServer.AddStreamEntry(EntryTypeL2Tx, tx.Encode())
					if err != nil {
						return err
					}
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
			// Commit at the end of each batch group
			err = streamServer.CommitAtomicOp()
			if err != nil {
				return err
			}
		}
	}

	return err
}

// GetSystemSCPosition computes the position of the intermediate state root for the system smart contract
func GetSystemSCPosition(blockNumber uint64) []byte {
	v1 := big.NewInt(0).SetUint64(blockNumber).Bytes()
	v2 := big.NewInt(0).SetUint64(uint64(posConstant)).Bytes()

	// Add 0s to make v1 and v2 32 bytes long
	for len(v1) < 32 {
		v1 = append([]byte{0}, v1...)
	}
	for len(v2) < 32 {
		v2 = append([]byte{0}, v2...)
	}

	return keccak256.Hash(v1, v2)
}

// computeFullBatches computes the full batches
func computeFullBatches(batches []*DSBatch, l2Blocks []*DSL2Block, l2Txs []*DSL2Transaction) []*DSFullBatch {
	prevL2BlockNumber := uint64(0)
	currentL2Block := 0
	currentL2Tx := 0

	fullBatches := make([]*DSFullBatch, 0)

	for _, batch := range batches {
		fullBatch := &DSFullBatch{
			DSBatch: *batch,
		}

		for i := currentL2Block; i < len(l2Blocks); i++ {
			l2Block := l2Blocks[i]

			if prevL2BlockNumber != 0 && l2Block.L2BlockNumber <= prevL2BlockNumber {
				continue
			}

			if l2Block.BatchNumber == batch.BatchNumber {
				fullBlock := DSL2FullBlock{
					DSL2Block: *l2Block,
				}

				for j := currentL2Tx; j < len(l2Txs); j++ {
					l2Tx := l2Txs[j]
					if l2Tx.L2BlockNumber == l2Block.L2BlockNumber {
						fullBlock.Txs = append(fullBlock.Txs, *l2Tx)
						currentL2Tx++
					}
					if l2Tx.L2BlockNumber > l2Block.L2BlockNumber {
						break
					}
				}

				fullBatch.L2Blocks = append(fullBatch.L2Blocks, fullBlock)
				prevL2BlockNumber = l2Block.L2BlockNumber
				currentL2Block++
			} else if l2Block.BatchNumber > batch.BatchNumber {
				break
			}
		}

		fullBatches = append(fullBatches, fullBatch)
	}

	return fullBatches
}
