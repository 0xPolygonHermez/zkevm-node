/*
This file provide functions to work with ETROG batches:
- EncodeBatchV2 (equivalent to EncodeTransactions)
- DecodeBatchV2 (equivalent to DecodeTxs)
- DecodeForcedBatchV2)


// batch data format:
// 0xb                             | 1  | changeL2Block
// --------- L2 block Header ---------------------------------
// 0x73e6af6f                      | 4  | deltaTimestamp
// 0x00000012					   | 4  | indexL1InfoTree
// -------- Transaction ---------------------------------------
// 0x00...0x00					   | n  | transaction RLP coded
// 0x00...0x00					   | 32 | R
// 0x00...0x00					   | 32 | S
// 0x00							   | 32 | V
// 0x00							   | 1  | efficiencyPercentage
// Repeat Transaction or changeL2Block
// Note: RLP codification: https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp/

/ forced batch data format:
// -------- Transaction ---------------------------------------
// 0x00...0x00					   | n  | transaction RLP coded
// 0x00...0x00					   | 32 | R
// 0x00...0x00					   | 32 | S
// 0x00							   | 32 | V
// 0x00							   | 1  | efficiencyPercentage
// Repeat Transaction
*/

package state

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// L2BlockRaw is the raw representation of a L2 block.
type L2BlockRaw struct {
	DeltaTimestamp  uint32
	IndexL1InfoTree uint32
	Transactions    []L2TxRaw
}

// BatchRawV2 is the  representation of a batch of transactions.
type BatchRawV2 struct {
	Blocks []L2BlockRaw
}

// ForcedBatchRawV2 is the  representation of a forced batch of transactions.
type ForcedBatchRawV2 struct {
	Transactions []L2TxRaw
}

// L2TxRaw is the raw representation of a L2 transaction  inside a L2 block.
type L2TxRaw struct {
	Tx                   types.Transaction
	EfficiencyPercentage uint8
}

const (
	changeL2Block = uint8(0x0b)
)

var (
	// ErrBatchV2DontStartWithChangeL2Block is returned when the batch start directly with a trsansaction (without a changeL2Block)
	ErrBatchV2DontStartWithChangeL2Block = errors.New("batch v2 must start with changeL2Block before Tx (suspect a V1 Batch or a ForcedBatch?))")
	// ErrInvalidBatchV2 is returned when the batch is invalid.
	ErrInvalidBatchV2 = errors.New("invalid batch v2")
	// ErrInvalidRLP is returned when the rlp is invalid.
	ErrInvalidRLP = errors.New("invalid rlp codification")
)

func (b *BatchRawV2) String() string {
	res := ""
	nTxs := 0
	for i, block := range b.Blocks {
		res += fmt.Sprintf("Block[%d/%d]: deltaTimestamp: %d, indexL1InfoTree: %d nTxs: %d\n", i, len(b.Blocks),
			block.DeltaTimestamp, block.IndexL1InfoTree, len(block.Transactions))
		nTxs += len(block.Transactions)
	}
	res = fmt.Sprintf("BATCHv2, nBlocks: %d nTxs:%d \n", len(b.Blocks), nTxs) + res
	return res
}

// EncodeBatchV2 encodes a batch of transactions into a byte slice.
func EncodeBatchV2(batch *BatchRawV2) ([]byte, error) {
	var err error
	var batchData []byte
	if batch == nil {
		return nil, fmt.Errorf("batch is nil: %w", ErrInvalidBatchV2)
	}
	blocks := batch.Blocks
	if len(blocks) == 0 {
		return nil, fmt.Errorf("a batch need minimum a L2Block: %w", ErrInvalidBatchV2)
	}
	for _, block := range blocks {
		batchData, err = EncodeBlockHeaderV2(batchData, block)
		if err != nil {
			return nil, fmt.Errorf("can't encode block header: %w", err)
		}
		for _, tx := range block.Transactions {
			batchData, err = encodeTxRLP(batchData, tx)
			if err != nil {
				return nil, fmt.Errorf("can't encode tx: %w", err)
			}
		}
	}
	return batchData, nil
}

// EncodeBlockHeaderV2 encodes a batch of l2blocks header into a byte slice.
func EncodeBlockHeaderV2(batchData []byte, block L2BlockRaw) ([]byte, error) {
	batchData = append(batchData, changeL2Block)
	batchData = append(batchData, serializeUint32(block.DeltaTimestamp)...)
	batchData = append(batchData, serializeUint32(block.IndexL1InfoTree)...)
	return batchData, nil
}

func encodeTxRLP(batchData []byte, tx L2TxRaw) ([]byte, error) {
	rlpTx, err := prepareRPLTxData(tx.Tx)
	if err != nil {
		return nil, fmt.Errorf("can't encode tx to RLP: %w", err)
	}
	batchData = append(batchData, rlpTx...)
	batchData = append(batchData, tx.EfficiencyPercentage)
	return batchData, nil
}

func serializeUint32(value uint32) []byte {
	return []byte{
		byte(value >> 24), // nolint:gomnd
		byte(value >> 16), // nolint:gomnd
		byte(value >> 8),  // nolint:gomnd
		byte(value),
	} // nolint:gomnd
}

// DecodeBatchV2 decodes a batch of transactions from a byte slice.
func DecodeBatchV2(txsData []byte) (*BatchRawV2, error) {
	// The transactions is not RLP encoded. Is the raw bytes in this form: 1 byte for the transaction type (always 0b for changeL2Block) + 4 bytes for deltaTimestamp + for bytes for indexL1InfoTree
	var err error
	var blocks []L2BlockRaw
	var currentBlock *L2BlockRaw
	pos := int(0)
	for pos < len(txsData) {
		switch txsData[pos] {
		case changeL2Block:
			if currentBlock != nil {
				blocks = append(blocks, *currentBlock)
			}
			pos, currentBlock, err = decodeBlockHeader(txsData, pos+1)
			if err != nil {
				return nil, fmt.Errorf("pos: %d can't decode new BlockHeader: %w", pos, err)
			}
		// by RLP definition a tx never starts with a 0x0b. So, if is not a changeL2Block
		// is a tx
		default:
			if currentBlock == nil {
				_, _, err := decodeTxRLP(txsData, pos)
				if err == nil {
					// There is no changeL2Block but have a valid RLP transaction
					return nil, ErrBatchV2DontStartWithChangeL2Block
				} else {
					// No changeL2Block and no valid RLP transaction
					return nil, fmt.Errorf("no ChangeL2Block neither valid Tx, batch malformed : %w", ErrInvalidBatchV2)
				}
			}
			var tx *L2TxRaw
			pos, tx, err = decodeTxRLP(txsData, pos)
			if err != nil {
				return nil, fmt.Errorf("can't decode transactions: %w", err)
			}

			currentBlock.Transactions = append(currentBlock.Transactions, *tx)
		}
	}
	if currentBlock != nil {
		blocks = append(blocks, *currentBlock)
	}
	return &BatchRawV2{blocks}, nil
}

// DecodeForcedBatchV2 decodes a forced batch V2 (Etrog)
// Is forbidden changeL2Block, so are just the set of transactions
func DecodeForcedBatchV2(txsData []byte) (*ForcedBatchRawV2, error) {
	txs, _, efficiencyPercentages, err := DecodeTxs(txsData, FORKID_ETROG)
	if err != nil {
		return nil, err
	}
	// Sanity check, this should never happen
	if len(efficiencyPercentages) != len(txs) {
		return nil, fmt.Errorf("error decoding len(efficiencyPercentages) != len(txs). len(efficiencyPercentages)=%d, len(txs)=%d : %w", len(efficiencyPercentages), len(txs), ErrInvalidRLP)
	}
	forcedBatch := ForcedBatchRawV2{}
	for i, tx := range txs {
		forcedBatch.Transactions = append(forcedBatch.Transactions, L2TxRaw{
			Tx:                   tx,
			EfficiencyPercentage: efficiencyPercentages[i],
		})
	}
	return &forcedBatch, nil
}

// decodeBlockHeader decodes a block header from a byte slice.
//
//	Extract: 4 bytes for deltaTimestamp + 4 bytes for indexL1InfoTree
func decodeBlockHeader(txsData []byte, pos int) (int, *L2BlockRaw, error) {
	var err error
	currentBlock := &L2BlockRaw{}
	pos, currentBlock.DeltaTimestamp, err = deserializeUint32(txsData, pos)
	if err != nil {
		return 0, nil, fmt.Errorf("can't get deltaTimestamp: %w", err)
	}
	pos, currentBlock.IndexL1InfoTree, err = deserializeUint32(txsData, pos)
	if err != nil {
		return 0, nil, fmt.Errorf("can't get leafIndex: %w", err)
	}

	return pos, currentBlock, nil
}

func decodeTxRLP(txsData []byte, offset int) (int, *L2TxRaw, error) {
	var err error
	length, err := decodeRLPListLengthFromOffset(txsData, offset)
	if err != nil {
		return 0, nil, fmt.Errorf("can't get RLP length (offset=%d): %w", offset, err)
	}
	endPos := uint64(offset) + length + rLength + sLength + vLength + EfficiencyPercentageByteLength
	if endPos > uint64(len(txsData)) {
		return 0, nil, fmt.Errorf("can't get tx because not enough data (endPos=%d lenData=%d): %w",
			endPos, len(txsData), ErrInvalidBatchV2)
	}
	fullDataTx := txsData[offset:endPos]
	dataStart := uint64(offset) + length
	txInfo := txsData[offset:dataStart]
	rData := txsData[dataStart : dataStart+rLength]
	sData := txsData[dataStart+rLength : dataStart+rLength+sLength]
	vData := txsData[dataStart+rLength+sLength : dataStart+rLength+sLength+vLength]
	efficiencyPercentage := txsData[dataStart+rLength+sLength+vLength]
	var rlpFields [][]byte
	err = rlp.DecodeBytes(txInfo, &rlpFields)
	if err != nil {
		log.Error("error decoding tx Bytes: ", err, ". fullDataTx: ", hex.EncodeToString(fullDataTx), "\n tx: ", hex.EncodeToString(txInfo), "\n Txs received: ", hex.EncodeToString(txsData))
		return 0, nil, err
	}
	legacyTx, err := RlpFieldsToLegacyTx(rlpFields, vData, rData, sData)
	if err != nil {
		log.Debug("error creating tx from rlp fields: ", err, ". fullDataTx: ", hex.EncodeToString(fullDataTx), "\n tx: ", hex.EncodeToString(txInfo), "\n Txs received: ", hex.EncodeToString(txsData))
		return 0, nil, err
	}

	l2Tx := &L2TxRaw{
		Tx:                   *types.NewTx(legacyTx),
		EfficiencyPercentage: efficiencyPercentage,
	}

	return int(endPos), l2Tx, err
}

func deserializeUint32(txsData []byte, pos int) (int, uint32, error) {
	if len(txsData)-pos < 4 { // nolint:gomnd
		return 0, 0, fmt.Errorf("can't get u32 because not enough data: %w", ErrInvalidBatchV2)
	}
	return pos + 4, uint32(txsData[pos])<<24 | uint32(txsData[pos+1])<<16 | uint32(txsData[pos+2])<<8 | uint32(txsData[pos+3]), nil // nolint:gomnd
}

// It returns the length of data from the param offset
// ex:
// 0xc0 -> empty data -> 1 byte because it include the 0xc0
func decodeRLPListLengthFromOffset(txsData []byte, offset int) (uint64, error) {
	txDataLength := uint64(len(txsData))
	num := uint64(txsData[offset])
	if num < c0 { // c0 -> is a empty data
		log.Debugf("error num < c0 : %d, %d", num, c0)
		return 0, fmt.Errorf("first byte of tx (%x) is < 0xc0: %w", num, ErrInvalidRLP)
	}
	length := num - c0
	if length > shortRlp { // If rlp is bigger than length 55
		// n is the length of the rlp data without the header (1 byte) for example "0xf7"
		pos64 := uint64(offset)
		lengthInByteOfSize := num - f7
		if (pos64 + headerByteLength + lengthInByteOfSize) > txDataLength {
			log.Debug("error not enough data: ")
			return 0, fmt.Errorf("not enough data to get length: %w", ErrInvalidRLP)
		}

		n, err := strconv.ParseUint(hex.EncodeToString(txsData[pos64+1:pos64+1+lengthInByteOfSize]), hex.Base, hex.BitSize64) // +1 is the header. For example 0xf7
		if err != nil {
			log.Debug("error parsing length: ", err)
			return 0, fmt.Errorf("error parsing length value: %w", err)
		}
		// TODO: RLP specifications says length = n ??? that is wrong??
		length = n + num - f7 // num - f7 is the header. For example 0xf7
	}
	return length + headerByteLength, nil
}
