package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BatchProcessor is used to process a batch of transactions
type BatchProcessor interface {
	ProcessBatch(batch Batch) error
	ProcessTransaction(tx types.Transaction) error
	CheckTransaction(tx types.Transaction) error
	Commit() (*common.Hash, *Proof, error)
	Rollback() error
}

// BatchProcessor is used to process a batch of transactions
type BasicBatchProcessor struct{}

// ProcessBatch processes all transactions inside a batch
func (b *BasicBatchProcessor) ProcessBatch(batch Batch) error {
	return nil
}

// ProcessTransaction processes a transaction inside a batch
func (b *BasicBatchProcessor) ProcessTransaction(tx types.Transaction) error {
	return nil
}

// CheckTransaction checks a transaction is valid inside a batch context
func (b *BasicBatchProcessor) CheckTransaction(tx types.Transaction) error {
	return nil
}

// Commits the batch state into state
func (b *BasicBatchProcessor) Commit() (*common.Hash, *Proof, error) {
	return nil, nil, nil
}

// Rollback does not apply batch state into state
func (b *BasicBatchProcessor) Rollback() error {
	return nil
}
