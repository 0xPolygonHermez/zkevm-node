package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// State
type State interface {
	NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) BatchProcessor
	GetStateRoot(virtual bool) (*big.Int, error)
	GetBalance(address common.Address, batchNumber uint64) (*big.Int, error)
	EstimageGas(address common.Address) uint64
	GetLastBlock() (*types.Block, error)
	GetLastBatch(isVirtual bool) (*Batch, error)
	GetBatchByHash(hash common.Hash, withTxDetails, isVirtual bool) (*Batch, error)
	GetBatchByNumber(number uint64, withTxDetails, isVirtual bool) (*Batch, error)
	GetTransactionByBatchHashAndIndex(hash common.Hash, index uint64) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(number uint64, index uint64) (*types.Transaction, error)
	GetTransaction(hash common.Hash) (*types.Transaction, error)
	GetTransactionReceipt(hash common.Hash) (*types.Receipt, error)
	GetNonce(address common.Address, batchNumber uint64) (uint64, error)
	Reset(batchnum uint64) error
	ConsolidateBatch(batch Batch) error
}

type BasicState struct {
}

// NewState creates a new State
func NewState() State {
	return &BasicState{}
}

// NewBatchProcessor creates a new batch processor
func (s *BasicState) NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) BatchProcessor {
	return &BasicBatchProcessor{}
}

// GetStateRoot returns the root of the state tree
func (s *BasicState) GetStateRoot(virtual bool) (*big.Int, error) {
	panic("not implemented yet")
}

// GetBalance from a given address
func (s *BasicState) GetBalance(address common.Address, batchNumber uint64) (*big.Int, error) {
	panic("not implemented yet")
}

// EstimateGas for a transaction
func (s *BasicState) EstimageGas(address common.Address) uint64 {
	panic("not implemented yet")
}

// GetLastBlock gets the latest block
func (s *BasicState) GetLastBlock() (*types.Block, error) {
	panic("not implemented yet")
}

// GetLastBatch gets the latest batch
func (s *BasicState) GetLastBatch(isVirtual bool) (*Batch, error) {
	panic("not implemented yet")
}

// GetBatchByHash gets a batch by its hash
func (s *BasicState) GetBatchByHash(hash common.Hash, withTxDetails, isVirtual bool) (*Batch, error) {
	panic("not implemented yet")
}

// GetBatchByNumber gets a batch by its number
func (s *BasicState) GetBatchByNumber(number uint64, withTxDetails, isVirtual bool) (*Batch, error) {
	panic("not implemented yet")
}

// GetTransactionByBlockHashAndIndex gets a transactions by its index accordingly to the batch hash
func (s *BasicState) GetTransactionByBatchHashAndIndex(hash common.Hash, index uint64) (*types.Transaction, error) {
	panic("not implemented yet")
}

// GetTransactionByBatchNumberAndIndex gets a transactions by its index accordingly to the batch number
func (s *BasicState) GetTransactionByBatchNumberAndIndex(number uint64, index uint64) (*types.Transaction, error) {
	panic("not implemented yet")
}

// GetTransaction gets a transactions by its hash
func (s *BasicState) GetTransaction(hash common.Hash) (*types.Transaction, error) {
	panic("not implemented yet")
}

// GetTransaction gets a transactions receipt by its hash
func (s *BasicState) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	panic("not implemented yet")
}

func (s *BasicState) GetNonce(address common.Address, batchNumber uint64) (uint64, error) {
	panic("not implemented yet")
}

// GetLastBatch gets the latest batch
func (s *BasicState) Reset(batchnum uint64) error {
	panic("not implemented yet")
}

// ConsolidateBatch changes the virtual status of a batch
func (s *BasicState) ConsolidateBatch(batch Batch) error {
	panic("not implemented yet")
}
