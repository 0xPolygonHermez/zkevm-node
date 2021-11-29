package state

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hermeznetwork/hermez-core/log"
)

var (
	// ErrInvalidSig indicates the signature of the transaction is not valid
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
	// ErrInvalidNonce indicates the nonce of the transaction is not valid
	ErrInvalidNonce = errors.New("invalid transaction nonce")
	// ErrInvalidBalance indicates the balance of the account is not enough to process the transaction
	ErrInvalidBalance = errors.New("not enough balance")
	// ErrInvalidGas indicates the gaslimit is not enough to process the transaction
	ErrInvalidGas = errors.New("not enough gas")
)

// BatchProcessor is used to process a batch of transactions
type BatchProcessor interface {
	ProcessBatch(batch *Batch) error
	ProcessTransaction(tx *types.Transaction, sequencerAddress common.Address) error
	CheckTransaction(tx *types.Transaction) (common.Address, *big.Int, *big.Int, error)
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
	// TODO: Check if batch is virtual and process accordingly
	for _, tx := range batch.Transactions {
		err := b.ProcessTransaction(tx, batch.Sequencer)
		if err != nil {
			log.Infof("Error processing transaction %s: %v", tx.Hash().String(), err)
		}

		receipt := types.NewReceipt(b.stateRoot, err != nil, 0)

		// TODO: Store receipt, log it to avoid unused linter error
		log.Debugf("%v", receipt)
	}

	return b.State.Tree.SetRootForBatchNumber(batch.BatchNumber, b.stateRoot)
}

// ProcessTransaction processes a transaction inside a batch
func (b *BasicBatchProcessor) ProcessTransaction(tx *types.Transaction, sequencerAddress common.Address) error {
	sender, nonce, senderBalance, err := b.CheckTransaction(tx)

	if err == nil {
		// Get receiver Balance
		receiverBalance, err := b.State.Tree.GetBalance(*tx.To(), b.stateRoot)
		if err != nil {
			return err
		}

		// Increase Nonce
		nonce.Add(nonce, big.NewInt(1))

		// Store new nonce
		root, _, err := b.State.Tree.SetNonce(sender, nonce)
		if err != nil {
			return err
		}
		b.stateRoot = root

		// Calculate new balances
		senderBalance.Sub(senderBalance, tx.Cost())
		receiverBalance.Add(receiverBalance, tx.Value())

		// Pay gas to the sequencer
		sequencerBalance, err := b.State.Tree.GetBalance(sequencerAddress, b.stateRoot)
		if err != nil {
			return err
		}

		usedGas := big.NewInt(0).SetUint64(b.State.EstimateGas(tx))

		sequencerBalance.Add(sequencerBalance, usedGas.Mul(usedGas, tx.GasPrice()))

		// Refund unused gas
		remainingGas := big.NewInt(0).SetUint64((tx.Gas() - usedGas.Uint64()))
		receiverBalance.Add(receiverBalance, remainingGas.Mul(remainingGas, tx.GasPrice()))

		// Store new balances
		root, _, err = b.State.Tree.SetBalance(sender, senderBalance)
		if err != nil {
			return err
		}
		b.stateRoot = root

		root, _, err = b.State.Tree.SetBalance(*tx.To(), receiverBalance)
		if err != nil {
			return err
		}
		b.stateRoot = root

		root, _, err = b.State.Tree.SetBalance(sequencerAddress, sequencerBalance)
		if err != nil {
			return err
		}
		b.stateRoot = root
	}

	return err
}

// CheckTransaction checks if a transaction is valid
func (b *BasicBatchProcessor) CheckTransaction(tx *types.Transaction) (common.Address, *big.Int, *big.Int, error) {
	// TODO: Check ChainID when possible
	var sender = common.Address{}
	var nonce = big.NewInt(0)
	var balance = big.NewInt(0)

	// Set stateRoot if needed
	if len(b.stateRoot) == 0 {
		root, err := b.State.Tree.GetRoot()
		if err != nil {
			return sender, nonce, balance, err
		}

		b.stateRoot = root
	}

	// Check Signature
	v, r, s := tx.RawSignatureValues()
	plainV := byte(v.Uint64() - 35 - 2*(tx.ChainId().Uint64()))

	if !crypto.ValidateSignatureValues(plainV, r, s, false) {
		return sender, nonce, balance, ErrInvalidSig
	}

	// Get Sender
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		return sender, nonce, balance, err
	}

	// Check nonce
	nonce, err = b.State.Tree.GetNonce(sender, b.stateRoot)
	if err != nil {
		return sender, nonce, balance, err
	}

	if nonce.Uint64() != tx.Nonce() {
		return sender, nonce, balance, ErrInvalidNonce
	}

	// Check balance
	balance, err = b.State.Tree.GetBalance(sender, b.stateRoot)
	if err != nil {
		return sender, nonce, balance, err
	}

	if balance.Cmp(tx.Cost()) < 0 {
		return sender, nonce, balance, ErrInvalidBalance
	}

	// Check gas
	if tx.Gas() < b.State.EstimateGas(tx) {
		return sender, nonce, balance, ErrInvalidGas
	}

	return sender, nonce, balance, nil
}

// Commit the batch state into state
func (b *BasicBatchProcessor) Commit() (*common.Hash, *Proof, error) {
	// TODO: Implement
	return nil, nil, nil
}

// Rollback does not apply batch state into state
func (b *BasicBatchProcessor) Rollback() error {
	// TODO: Implement
	return nil
}
