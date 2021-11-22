package state

import (
	"github.com/ethereum/go-ethereum/common"
)

// BatchProcessor is used to process a batch of transactions
type BatchProcessor interface {
	ProcessBatch(batch Batch) error
	ProcessTransaction(tx types.LegacyTx) error
	CheckTransaction(tx types.LegacyTx) error
	Commit() (*common.Hash, *Proof, error)
	Rollback() error
}

// BatchProcessor is used to process a batch of transactions
type BasicBatchProcessor struct {
	State *BasicState
}

// ProcessBatch processes all transactions inside a batch
func (b *BasicBatchProcessor) ProcessBatch(batch Batch) error {
	return nil
}

// ProcessTransaction processes a transaction inside a batch
func (b *BasicBatchProcessor) ProcessTransaction(tx types.LegacyTx) error {
	return nil
}

// CheckTransaction checks a transaction is valid inside a batch context
func (b *BasicBatchProcessor) CheckTransaction(tx types.LegacyTx) error {

	// Check Nonce
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
