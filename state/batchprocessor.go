package state

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
}

const (
	addBatchSQL = "INSERT INTO state.batch (batch_num, batch_hash, block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
)

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
			log.Errorf("Error processing transaction %s: %v", tx.Hash().String(), err)
		} else {
			log.Infof("Successfully processed transaction %s", tx.Hash().String())
		}

		// receipt := types.NewReceipt(b.stateRoot, err != nil, 0)

		// TODO: Store receipt, log it to avoid unused linter error
		// log.Debugf("%v", receipt)
	}

	_, _, err := b.commit(batch)

	return err
}

// ProcessTransaction processes a transaction inside a batch
func (b *BasicBatchProcessor) ProcessTransaction(tx *types.Transaction, sequencerAddress common.Address) error {
	// save stateRoot and modify it only if transaction processing finishes successfully
	root := b.stateRoot

	// reset MT currentRoot in case it was modified by failed transaction
	b.State.Tree.SetCurrentRoot(root)

	sender, nonce, senderBalance, err := b.CheckTransaction(tx)

	if err != nil {
		return err
	}

	// Get receiver Balance
	receiverBalance, err := b.State.Tree.GetBalance(*tx.To(), root)
	if err != nil {
		return err
	}

	// Increase Nonce
	nonce.Add(nonce, big.NewInt(1))

	// Store new nonce
	root, _, err = b.State.Tree.SetNonce(sender, nonce)
	if err != nil {
		return err
	}

	// Calculate new balances
	senderBalance.Sub(senderBalance, tx.Cost())
	receiverBalance.Add(receiverBalance, tx.Value())

	// Pay gas to the sequencer
	usedGas := new(big.Int).SetUint64(b.State.EstimateGas(tx))

	if sequencerAddress == sender {
		senderBalance.Add(senderBalance, new(big.Int).Mul(usedGas, tx.GasPrice()))
	} else if sequencerAddress == *tx.To() {
		receiverBalance.Add(receiverBalance, new(big.Int).Mul(usedGas, tx.GasPrice()))
	} else {
		sequencerBalance, err := b.State.Tree.GetBalance(sequencerAddress, root)
		if err != nil {
			return err
		}

		sequencerBalance.Add(sequencerBalance, new(big.Int).Mul(usedGas, tx.GasPrice()))
		root, _, err = b.State.Tree.SetBalance(sequencerAddress, sequencerBalance)
		if err != nil {
			return err
		}
	}

	// Refund unused gas
	remainingGas := new(big.Int).SetUint64((tx.Gas() - usedGas.Uint64()))
	senderBalance.Add(senderBalance, new(big.Int).Mul(remainingGas, tx.GasPrice()))

	// Store new balances
	root, _, err = b.State.Tree.SetBalance(sender, senderBalance)
	if err != nil {
		return err
	}

	root, _, err = b.State.Tree.SetBalance(*tx.To(), receiverBalance)
	if err != nil {
		return err
	}

	b.stateRoot = root

	return nil
}

// CheckTransaction checks if a transaction is valid
func (b *BasicBatchProcessor) CheckTransaction(tx *types.Transaction) (common.Address, *big.Int, *big.Int, error) {
	// TODO: Check ChainID when possible
	var sender = common.Address{}
	var nonce = big.NewInt(0)
	var balance = big.NewInt(0)

	// reset MT currentRoot in case it was modified by failed transaction
	b.State.Tree.SetCurrentRoot(b.stateRoot)

	// Set stateRoot if needed
	/*
		if len(b.stateRoot) == 0 {
			root, err := b.State.Tree.GetCurrentRoot()
			if err != nil {
				return sender, nonce, balance, err
			}

			b.stateRoot = root
		}
	*/
	// Check Signature
	if err := CheckSignature(tx); err != nil {
		return sender, nonce, balance, err
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
func (b *BasicBatchProcessor) commit(batch *Batch) (*common.Hash, *Proof, error) {
	// Store batch into db
	ctx := context.Background()

	var root common.Hash

	if batch.Header == nil {
		batch.Header = &types.Header{
			Root:       root,
			Difficulty: big.NewInt(0),
			Number:     new(big.Int).SetUint64(batch.BatchNumber),
		}
	}

	// set merkletree root
	if b.stateRoot != nil {
		root.SetBytes(b.stateRoot)
		batch.Header.Root = root
	}

	err := b.State.addBatch(ctx, batch)
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}

// Rollback does not apply batch state into state
// TODO: implement
/*
func (b *BasicBatchProcessor) rollback() error {
	// TODO: Implement
	return nil
}
*/

func (s *BasicState) addBatch(ctx context.Context, batch *Batch) error {
	_, err := s.db.Exec(ctx, addBatchSQL, batch.BatchNumber, batch.BatchHash, batch.BlockNumber, batch.Sequencer, batch.Aggregator,
		batch.ConsolidatedTxHash, batch.Header, batch.Uncles, batch.RawTxsData)
	return err
}
