package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BatchProcessor is used to process a batch of transactions
type BatchProcessor struct{}

// ProcessBatch processes all transactions inside a batch
func (b *BatchProcessor) ProcessBatch(batch Batch) error {
	return nil
}

// ProcessTransaction processes a transaction inside a batch
func (b *BatchProcessor) ProcessTransaction(tx types.Transaction) error {
	return nil
}

// CheckTransaction checks a transaction is valid inside a batch context
func (b *BatchProcessor) CheckTransaction(tx types.Transaction) error {
	return nil
}

// Commits the batch state into state
func (b *BatchProcessor) Commit() (*common.Hash, *Proof, error) {
	return nil, nil, nil
}

// Rollback does not apply batch state into state
func (b *BatchProcessor) Rollback() error {
	return nil
}
