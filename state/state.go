package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state/db"
	// "github.com/hermeznetwork/hermez-core/state/merkletree"
	// "github.com/hermeznetwork/hermez-core/state/merkletree/leafs"
)

// State
type State struct {
	// StateTree merkletree.Merkletree
}

// NewState creates a new State
func NewState(db db.KeyValuer) *State {
	// return &State{StateTree: merkletree.NewMerkletree(db)}
	return &State{}
}

// NewBatchProcessor creates a new batch processor
func (s *State) NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) *BatchProcessor {
	return &BatchProcessor{}
}

// GetStateRoot returns the root of the state tree
func (s *State) GetStateRoot(virtual bool) (*big.Int, error) {
	// return s.StateTree.Root, nil
	return nil, nil
}

// GetBalance from a given address
func (s *State) GetBalance(address common.Address) (*big.Int, error) {
	/*
		key, err := leafs.NewBalanceKey(common.BytesToAddress(address.Bytes()))
		if err != nil {
			return nil, err
		}

		balanceBytes, err := s.StateTree.Get(s.StateTree.Root, key)
		if err != nil {
			return nil, err
		}

		return leafs.BytesToBalance(balanceBytes), nil*/
	return nil, nil
}

// EstimateGas for a transaction
func (s *State) EstimageGas(transaction types.Transaction) uint64 {
	return 21000
}

// GetLastBlock gets the latest block
func (s *State) GetLastBlock() (*types.Block, error) {
	return nil, nil
}

// GetPreviousBlock gets the offset previous block respect to latest
func (s *State) GetPreviousBlock(offset uint64) (*types.Block, error) {
	return nil, nil
}

// GetBlockByHash gets the block with the required hash
func (s *State) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	return nil, nil
}

// GetBlockByNumber gets the block with the required number
func (s *State) GetBlockByNumber(blockNumber uint64) (*types.Block, error) {
	return nil, nil
}

// GetLastBlockNumber gets the latest block number
func (s *State) GetLastBlockNumber() (uint64, error) {
	return 0, nil
}

// GetLastBatch gets the latest batch
func (s *State) GetLastBatch(isVirtual bool) (*Batch, error) {
	return nil, nil
}

// GetPreviousBatch gets the offset previous batch respect to latest
func (s *State) GetPreviousBatch(offset uint64) (*Batch, error) {
	return nil, nil
}

// GetBatchByHash gets the batch with the required hash
func (s *State) GetBatchByHash(hash common.Hash) (*types.Block, error) {
	return nil, nil
}

// GetBatchByNumber gets the batch with the required number
func (s *State) GetBatchByNumber(batchNumber uint64) (*types.Block, error) {
	return nil, nil
}

// GetLastBatchNumber gets the latest batch number
func (s *State) GetLastBatchNumber() (uint64, error) {
	return 0, nil
}

// GetTransactionByBatchHashAndIndex gets a transaction from a batch by index
func (s *State) GetTransactionByBatchHashAndIndex(batchHash common.Hash, index uint) (*types.Transaction, error) {
	return nil, nil
}

// GetTransactionByBatchNumberAndIndex gets a transaction from a batch by index
func (s *State) GetTransactionByBatchNumberAndIndex(batchNumber uint64, index uint) (*types.Transaction, error) {
	return nil, nil
}

// GetTransactionByHash gets a transaction by its hash
func (s *State) GetTransactionByHash(transactionHash common.Hash) (*types.Transaction, error) {
	return nil, nil
}

// GetTransactionCount returns the number of transactions sent from an address
func (s *State) GetTransactionCount(address common.Address) (uint, error) {
	return 0, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
func (s *State) GetTransactionReceipt(transactionHash common.Hash) (*types.Receipt, error) {
	return nil, nil
}

// Reset resets the state to a block
func (s *State) Reset(blockNumber uint64) error {
	return nil
}

// ConsolidateBatch changes the virtual status of a batch
func (s *State) ConsolidateBatch(batchNumber uint64) error {
	return nil
}

func (s *State) GetTxsByBatchNum(batchNum uint64) ([]*types.Transaction, error) {
	return nil, nil
}
