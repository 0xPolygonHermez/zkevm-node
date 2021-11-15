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
	// return s.data.Root, nil
	return nil, nil
}

// GetBalance from a given address
func (s *State) GetBalance(address common.Address) *big.Int {
	return nil
}

// EstimateGas for a transaction
func (s *State) EstimageGas(address common.Address) uint64 {
	return 0
}

// GetLastBlock gets the latest block
func (s *State) GetLastBlock() (*types.Block, error) {
	return nil, nil
}

// GetLastBatch gets the latest batch
func (s *State) GetLastBatch(isVirtual bool) (*Batch, error) {
	return nil, nil
}

// GetLastBatch gets the latest batch
func (s *State) Reset(batchnum uint64) error {
	return nil
}

// ConsolidateBatch changes the virtual status of a batch
func (s *State) ConsolidateBatch(batch Batch) error {
	return nil
}
