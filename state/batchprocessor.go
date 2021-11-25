package state

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrInvalidSig     = errors.New("invalid transaction v, r, s values")
	ErrInvalidNonce   = errors.New("invalid transaction nonce")
	ErrInvalidBalance = errors.New("not enough balance")
	ErrInvalidGas     = errors.New("not enough gas")
)

// BatchProcessor is used to process a batch of transactions
type BatchProcessor interface {
	ProcessBatch(batch *Batch) error
	ProcessTransaction(tx *types.Transaction) error
	CheckTransaction(tx *types.Transaction) error
	CheckTransactionForRoot(tx *types.Transaction, root []byte) error
	Commit() (*common.Hash, *Proof, error)
	Rollback() error
}

// BasicBatchProcessor is used to process a batch of transactions
type BasicBatchProcessor struct {
	State     *BasicState
	stateRoot []byte
}

// ProcessBatch processes all transactions inside a batch
func (b *BasicBatchProcessor) ProcessBatch(batch *Batch) error {
	// TODO: Implement
	root, err := b.State.Tree.GetRoot()
	if err != nil {
		return err
	}
	b.stateRoot = root
	return nil
}

// ProcessTransaction processes a transaction inside a batch
func (b *BasicBatchProcessor) ProcessTransaction(tx *types.Transaction) error {
	return nil
}

// CheckTransaction checks if a transaction is valid
func (b *BasicBatchProcessor) CheckTransaction(tx *types.Transaction) error {
	root, err := b.State.Tree.GetRoot()
	if err != nil {
		return err
	}
	return b.CheckTransactionForRoot(tx, root)
}

// CheckTransactionForRoot checks if a transaction is valid inside a batch context
func (b *BasicBatchProcessor) CheckTransactionForRoot(tx *types.Transaction, root []byte) error {
	// Check Signature
	v, r, s := tx.RawSignatureValues()
	plainV := byte(v.Uint64() - 35 - 2*(tx.ChainId().Uint64()))

	if !crypto.ValidateSignatureValues(plainV, r, s, false) {
		return ErrInvalidSig
	}

	// Get Sender
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		return err
	}

	// Check nonce
	nonce, err := b.State.Tree.GetNonce(sender, root)
	if err != nil {
		return err
	}

	if nonce.Uint64() != tx.Nonce() {
		return ErrInvalidNonce
	}

	// Check balance
	balance, err := b.State.Tree.GetBalance(sender, root)
	if err != nil {
		return err
	}

	if balance.Cmp(tx.Cost()) < 0 {
		return ErrInvalidBalance
	}

	// Check gas
	if tx.Gas() < b.State.EstimateGas(tx) {
		return ErrInvalidGas
	}

	return nil
}

// Commit the batch state into state
func (b *BasicBatchProcessor) Commit() (*common.Hash, *Proof, error) {
	return nil, nil, nil
}

// Rollback does not apply batch state into state
func (b *BasicBatchProcessor) Rollback() error {
	return nil
}
