package state

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
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

// SetGenesis populates state with genesis information
func (s *State) SetGenesis(ctx context.Context, block Block, genesis Genesis, dbTx pgx.Tx) ([]byte, error) {
	var (
		root    common.Hash
		newRoot []byte
		err     error
	)
	if dbTx == nil {
		return newRoot, ErrDBTxNil
	}
	if s.tree == nil {
		return newRoot, ErrStateTreeNil
	}

	for _, action := range genesis.GenesisActions {
		address := common.HexToAddress(action.Address)
		switch action.Type {
		case int(merkletree.LeafTypeBalance):
			balance, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return newRoot, err
			}
			newRoot, _, err = s.tree.SetBalance(ctx, address, balance, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeNonce):
			nonce, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return newRoot, err
			}
			newRoot, _, err = s.tree.SetNonce(ctx, address, nonce, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeCode):
			code, err := hex.DecodeHex(action.Bytecode)
			if err != nil {
				return newRoot, fmt.Errorf("could not decode SC bytecode for address %q: %v", address, err)
			}
			newRoot, _, err = s.tree.SetCode(ctx, address, code, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeStorage):
			// Parse position and value
			positionBI, err := encoding.DecodeBigIntHexOrDecimal(action.StoragePosition)
			if err != nil {
				return newRoot, err
			}
			valueBI, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return newRoot, err
			}
			// Store
			newRoot, _, err = s.tree.SetStorageAt(ctx, address, positionBI, valueBI, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeSCLength):
			log.Debug("Skipped genesis action of type merkletree.LeafTypeSCLength, these actions will be handled as part of merkletree.LeafTypeCode actions")
		default:
			return newRoot, fmt.Errorf("unknown genesis action type %q", action.Type)
		}
	}

	root.SetBytes(newRoot)

	// flush state db
	err = s.tree.Flush(ctx)
	if err != nil {
		log.Errorf("error flushing state tree after genesis: %v", err)
		return newRoot, err
	}

	// store L1 block related to genesis batch
	err = s.AddBlock(ctx, &block, dbTx)
	if err != nil {
		return newRoot, err
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
		return newRoot, err
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
		return newRoot, err
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
		return newRoot, err
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

	return newRoot, s.AddL2Block(ctx, batch.BatchNumber, l2Block, receipts, MaxEffectivePercentage, dbTx)
}
