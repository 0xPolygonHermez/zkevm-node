package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state/db"
)

// State
type State struct {
}

// NewState creates a new State
func NewState(db db.KeyValuer) *State {
	return &State{}
}

// NewBatchProcessor creates a new batch processor
func (s *State) NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) *BatchProcessor {
	return &BatchProcessor{}
}

// GetStateRoot returns the root of the state tree
func (s *State) GetStateRoot(virtual bool) (*big.Int, error) {
	panic("not implemented yet")
}

// GetBalance from a given address
func (s *State) GetBalance(address common.Address, batchNumber uint64) (*big.Int, error) {
	panic("not implemented yet")
}

// EstimateGas for a transaction
func (s *State) EstimageGas(address common.Address) uint64 {
	panic("not implemented yet")
}

// GetLastBlock gets the latest block
func (s *State) GetLastBlock() (*types.Block, error) {
	panic("not implemented yet")
}

// GetLastBatch gets the latest batch
func (s *State) GetLastBatch(isVirtual bool) (*Batch, error) {
	panic("not implemented yet")
}

// GetBatchByHash gets a batch by its hash
func (s *State) GetBatchByHash(hash common.Hash, withTxDetails, isVirtual bool) (*Batch, error) {
	panic("not implemented yet")
}

// GetBatchByNumber gets a batch by its number
func (s *State) GetBatchByNumber(number uint64, withTxDetails, isVirtual bool) (*Batch, error) {
	panic("not implemented yet")
}

// GetTransactionByBlockHashAndIndex gets a transactions by its index accordingly to the batch hash
func (s *State) GetTransactionByBatchHashAndIndex(hash common.Hash, index uint64) (*types.Transaction, error) {
	panic("not implemented yet")
}

// GetTransactionByBatchNumberAndIndex gets a transactions by its index accordingly to the batch number
func (s *State) GetTransactionByBatchNumberAndIndex(number uint64, index uint64) (*types.Transaction, error) {
	panic("not implemented yet")
}

// GetTransaction gets a transactions by its hash
func (s *State) GetTransaction(hash common.Hash) (*types.Transaction, error) {
	panic("not implemented yet")
}

func (s *State) GetNonce(address common.Address, batchNumber uint64) (uint64, error) {
	panic("not implemented yet")
}

// GetLastBatch gets the latest batch
func (s *State) Reset(batchnum uint64) error {
	panic("not implemented yet")
}

// ConsolidateBatch changes the virtual status of a batch
func (s *State) ConsolidateBatch(batch Batch) error {
	panic("not implemented yet")
}
