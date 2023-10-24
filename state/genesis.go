package state

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// Genesis contains the information to populate state on creation
type Genesis struct {
	// GenesisBlockNum is the block number where the polygonZKEVM smc was deployed on L1
	GenesisBlockNum uint64
	// Root hash of the genesis block
	Root common.Hash
	// Contracts to be deployed to L2
	GenesisActions []*GenesisAction
	// Data of the first batch after the genesis(Batch 1)
	FirstBatchData BatchData
}

// GenesisAction represents one of the values set on the SMT during genesis.
type GenesisAction struct {
	Address         string `json:"address"`
	Type            int    `json:"type"`
	StoragePosition string `json:"storagePosition"`
	Bytecode        string `json:"bytecode"`
	Key             string `json:"key"`
	Value           string `json:"value"`
	Root            string `json:"root"`
}

// BatchData represents the data of the first batch that contains initial transaction
type BatchData struct {
	Transactions   string         `json:"transactions"`
	GlobalExitRoot common.Hash    `json:"globalExitRoot"`
	Timestamp      uint64         `json:"timestamp"`
	Sequencer      common.Address `json:"sequencer"`
}

// SetGenesis populates state with genesis information
func (s *State) SetGenesis(ctx context.Context, block Block, genesis Genesis, m metrics.CallerLabel, dbTx pgx.Tx) (common.Hash, error) {
	var (
		root             common.Hash
		genesisStateRoot []byte
		err              error
	)
	if dbTx == nil {
		return common.Hash{}, ErrDBTxNil
	}
	if s.tree == nil {
		return common.Hash{}, ErrStateTreeNil
	}

	uuid := uuid.New().String()

	for _, action := range genesis.GenesisActions {
		address := common.HexToAddress(action.Address)
		switch action.Type {
		case int(merkletree.LeafTypeBalance):
			balance, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return common.Hash{}, err
			}
			genesisStateRoot, _, err = s.tree.SetBalance(ctx, address, balance, genesisStateRoot, uuid)
			if err != nil {
				return common.Hash{}, err
			}
		case int(merkletree.LeafTypeNonce):
			nonce, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return common.Hash{}, err
			}
			genesisStateRoot, _, err = s.tree.SetNonce(ctx, address, nonce, genesisStateRoot, uuid)
			if err != nil {
				return common.Hash{}, err
			}
		case int(merkletree.LeafTypeCode):
			code, err := hex.DecodeHex(action.Bytecode)
			if err != nil {
				return common.Hash{}, fmt.Errorf("could not decode SC bytecode for address %q: %v", address, err)
			}
			genesisStateRoot, _, err = s.tree.SetCode(ctx, address, code, genesisStateRoot, uuid)
			if err != nil {
				return common.Hash{}, err
			}
		case int(merkletree.LeafTypeStorage):
			// Parse position and value
			positionBI, err := encoding.DecodeBigIntHexOrDecimal(action.StoragePosition)
			if err != nil {
				return common.Hash{}, err
			}
			valueBI, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return common.Hash{}, err
			}
			// Store
			genesisStateRoot, _, err = s.tree.SetStorageAt(ctx, address, positionBI, valueBI, genesisStateRoot, uuid)
			if err != nil {
				return common.Hash{}, err
			}
		case int(merkletree.LeafTypeSCLength):
			log.Debug("Skipped genesis action of type merkletree.LeafTypeSCLength, these actions will be handled as part of merkletree.LeafTypeCode actions")
		default:
			return common.Hash{}, fmt.Errorf("unknown genesis action type %q", action.Type)
		}
	}

	root.SetBytes(genesisStateRoot)

	// flush state db
	err = s.tree.Flush(ctx, uuid)
	if err != nil {
		log.Errorf("error flushing state tree after genesis: %v", err)
		return common.Hash{}, err
	}

	// store L1 block related to genesis batch
	err = s.AddBlock(ctx, &block, dbTx)
	if err != nil {
		return common.Hash{}, err
	}

	// store genesis batch
	batch := Batch{
		BatchNumber:    0,
		Coinbase:       ZeroAddress,
		BatchL2Data:    nil,
		StateRoot:      root,
		LocalExitRoot:  ZeroHash,
		Timestamp:      block.ReceivedAt,
		Transactions:   []types.Transaction{},
		GlobalExitRoot: ZeroHash,
		ForcedBatchNum: nil,
	}

	err = s.storeGenesisBatch(ctx, batch, dbTx)
	if err != nil {
		return common.Hash{}, err
	}

	// mark the genesis batch as virtualized
	virtualBatch := &VirtualBatch{
		BatchNumber: batch.BatchNumber,
		TxHash:      ZeroHash,
		Coinbase:    ZeroAddress,
		BlockNumber: block.BlockNumber,
	}
	err = s.AddVirtualBatch(ctx, virtualBatch, dbTx)
	if err != nil {
		return common.Hash{}, err
	}

	// mark the genesis batch as verified/consolidated
	verifiedBatch := &VerifiedBatch{
		BatchNumber: batch.BatchNumber,
		TxHash:      ZeroHash,
		Aggregator:  ZeroAddress,
		BlockNumber: block.BlockNumber,
	}
	err = s.AddVerifiedBatch(ctx, verifiedBatch, dbTx)
	if err != nil {
		return common.Hash{}, err
	}

	// store L2 genesis block
	header := &types.Header{
		Number:     big.NewInt(0),
		ParentHash: ZeroHash,
		Coinbase:   ZeroAddress,
		Root:       root,
		Time:       uint64(block.ReceivedAt.Unix()),
	}
	rootHex := root.Hex()
	log.Info("Genesis root ", rootHex)

	receipts := []*types.Receipt{}
	l2Block := types.NewBlock(header, []*types.Transaction{}, []*types.Header{}, receipts, &trie.StackTrie{})
	l2Block.ReceivedAt = block.ReceivedAt

	err = s.AddL2Block(ctx, batch.BatchNumber, l2Block, receipts, MaxEffectivePercentage, dbTx)
	if err != nil {
		return common.Hash{}, err
	}

	// // Process FirstTransaction included in batch 1
	// batchL2Data := common.Hex2Bytes(genesis.FirstBatchData.Transactions[2:])
	// processCtx := ProcessingContext{
	// 	BatchNumber:    1,
	// 	Coinbase:       genesis.FirstBatchData.Sequencer,
	// 	Timestamp:      time.Unix(int64(genesis.FirstBatchData.Timestamp), 0),
	// 	GlobalExitRoot: genesis.FirstBatchData.GlobalExitRoot,
	// 	BatchL2Data:    &batchL2Data,
	// }
	// newStateRoot, flushID, proverID, err := s.ProcessAndStoreClosedBatch(ctx, processCtx, batch.BatchL2Data, dbTx, m)
	// if err != nil {
	// 	log.Error("error storing batch 1. Error: ", err)
	// 	return common.Hash{}, common.Hash{}, 0, "", err
	// }
	// var gRoot common.Hash
	// gRoot.SetBytes(genesisStateRoot)

	// // Virtualize Batch and add sequence
	// virtualBatch1 := VirtualBatch{
	// 	BatchNumber:   1,
	// 	TxHash:        ZeroHash,
	// 	Coinbase:      genesis.FirstBatchData.Sequencer,
	// 	BlockNumber:   block.BlockNumber,
	// 	SequencerAddr: genesis.FirstBatchData.Sequencer,
	// }
	// err = s.AddVirtualBatch(ctx, &virtualBatch1, dbTx)
	// if err != nil {
	// 	log.Errorf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch1.BatchNumber, genesis.GenesisBlockNum, err)
	// 	return common.Hash{}, common.Hash{}, 0, "", err
	// }
	// // Insert the sequence to allow the aggregator verify the sequence batches
	// seq := Sequence{
	// 	FromBatchNumber: 1,
	// 	ToBatchNumber:   1,
	// }
	// err = s.AddSequence(ctx, seq, dbTx)
	// if err != nil {
	// 	log.Errorf("error adding sequence. Sequence: %+v", seq)
	// 	return common.Hash{}, common.Hash{}, 0, "", err
	// }
	return root, nil
}
