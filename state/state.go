package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state/db"
)

// State is the interface of the Hermez state
type State interface {
	NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) BatchProcessor
	GetStateRoot(virtual bool) (*big.Int, error)
	GetBalance(address common.Address, batchNumber uint64) (*big.Int, error)
	EstimateGas(transaction types.Transaction) uint64
	GetLastBlock() (*types.Block, error)
	GetPreviousBlock(offset uint64) (*types.Block, error)
	GetBlockByHash(hash common.Hash) (*types.Block, error)
	GetBlockByNumber(blockNumber uint64) (*types.Block, error)
	GetLastBlockNumber() (uint64, error)
	GetLastBatch(isVirtual bool) (*Batch, error)
	GetTransaction(hash common.Hash) (*types.Transaction, error)
	GetNonce(address common.Address, batchNumber uint64) (uint64, error)
	GetPreviousBatch(offset uint64) (*Batch, error)
	GetBatchByHash(hash common.Hash) (*types.Block, error)
	GetBatchByNumber(batchNumber uint64) (*types.Block, error)
	GetLastBatchNumber() (uint64, error)
	GetTransactionByBatchHashAndIndex(batchHash common.Hash, index uint64) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(batchNumber uint64, index uint64) (*types.Transaction, error)
	GetTransactionByHash(transactionHash common.Hash) (*types.Transaction, error)
	GetTransactionCount(address common.Address) (uint64, error)
	GetTransactionReceipt(transactionHash common.Hash) (*types.Receipt, error)
	Reset(blockNumber uint64) error
	ConsolidateBatch(batchNumber uint64) error
	GetTxsByBatchNum(batchNum uint64) ([]*types.Transaction, error)
	AddNewSequencer(seq Sequencer) error
}

// BasicState is a implementation of the state
type BasicState struct {
	// StateTree merkletree.Merkletree
}

// NewState creates a new State
func NewState(db db.KeyValuer) State {
	// return &State{StateTree: merkletree.NewMerkletree(db)}
	return &BasicState{}
}

// NewBatchProcessor creates a new batch processor
func (s *BasicState) NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) BatchProcessor {
	return &BasicBatchProcessor{State: s}
}

// GetStateRoot returns the root of the state tree
func (s *BasicState) GetStateRoot(virtual bool) (*big.Int, error) {
	// return s.StateTree.Root, nil
	return nil, nil
}

// GetBalance from a given address
func (s *BasicState) GetBalance(address common.Address, batchNumber uint64) (*big.Int, error) {
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
func (s *BasicState) EstimateGas(transaction types.Transaction) uint64 {
	// TODO: Calculate once we have txs that interact with SCs
	return 21000 //nolint:gomnd
}

// GetLastBlock gets the latest block
func (s *BasicState) GetLastBlock() (*types.Block, error) {
	panic("not implemented yet")
}

// GetPreviousBlock gets the offset previous block respect to latest
func (s *BasicState) GetPreviousBlock(offset uint64) (*types.Block, error) {
	return nil, nil
}

// GetBlockByHash gets the block with the required hash
func (s *BasicState) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	return nil, nil
}

// GetBlockByNumber gets the block with the required number
func (s *BasicState) GetBlockByNumber(blockNumber uint64) (*types.Block, error) {
	return nil, nil
}

// GetLastBlockNumber gets the latest block number
func (s *BasicState) GetLastBlockNumber() (uint64, error) {
	return 0, nil
}

// GetLastBatch gets the latest batch
func (s *BasicState) GetLastBatch(isVirtual bool) (*Batch, error) {
	panic("not implemented yet")
}

// GetTransaction gets a transactions by its hash
func (s *BasicState) GetTransaction(hash common.Hash) (*types.Transaction, error) {
	panic("not implemented yet")
}

// GetNonce returns the nonce of the given account at the given batch number
func (s *BasicState) GetNonce(address common.Address, batchNumber uint64) (uint64, error) {
	panic("not implemented yet")
}

// GetPreviousBatch gets the offset previous batch respect to latest
func (s *BasicState) GetPreviousBatch(offset uint64) (*Batch, error) {
	return nil, nil
}

// GetBatchByHash gets the batch with the required hash
func (s *BasicState) GetBatchByHash(hash common.Hash) (*types.Block, error) {
	return nil, nil
}

// GetBatchByNumber gets the batch with the required number
func (s *BasicState) GetBatchByNumber(batchNumber uint64) (*types.Block, error) {
	return nil, nil
}

// GetLastBatchNumber gets the latest batch number
func (s *BasicState) GetLastBatchNumber() (uint64, error) {
	return 0, nil
}

// GetTransactionByBatchHashAndIndex gets a transaction from a batch by index
func (s *BasicState) GetTransactionByBatchHashAndIndex(batchHash common.Hash, index uint64) (*types.Transaction, error) {
	return nil, nil
}

// GetTransactionByBatchNumberAndIndex gets a transaction from a batch by index
func (s *BasicState) GetTransactionByBatchNumberAndIndex(batchNumber uint64, index uint64) (*types.Transaction, error) {
	return nil, nil
}

// GetTransactionByHash gets a transaction by its hash
func (s *BasicState) GetTransactionByHash(transactionHash common.Hash) (*types.Transaction, error) {
	return nil, nil
}

// GetTransactionCount returns the number of transactions sent from an address
func (s *BasicState) GetTransactionCount(address common.Address) (uint64, error) {
	return 0, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
func (s *BasicState) GetTransactionReceipt(transactionHash common.Hash) (*types.Receipt, error) {
	return nil, nil
}

// Reset resets the state to a block
func (s *BasicState) Reset(blockNumber uint64) error {
	return nil
}

// ConsolidateBatch changes the virtual status of a batch
func (s *BasicState) ConsolidateBatch(batchNumber uint64) error {
	return nil
}

// GetTxsByBatchNum returns all the txs in a given batch
func (s *BasicState) GetTxsByBatchNum(batchNum uint64) ([]*types.Transaction, error) {
	return nil, nil
}

// AddNewSequencer stores a new sequencer
func (s *BasicState) AddNewSequencer(seq Sequencer) error {
	return nil
}
