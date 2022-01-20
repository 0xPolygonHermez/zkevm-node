package state

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/hex"
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
	// ErrInvalidChainID indicates a mismatch between sequencer address and ChainID
	ErrInvalidChainID = errors.New("invalid chain id for sequencer")
)

// BatchProcessor is used to process a batch of transactions
type BatchProcessor interface {
	ProcessBatch(batch *Batch) error
	CheckTransaction(tx *types.Transaction) error
}

// BasicBatchProcessor is used to process a batch of transactions
type BasicBatchProcessor struct {
	State            *BasicState
	stateRoot        []byte
	SequencerAddress common.Address
	SequencerChainID uint64
}

// ProcessBatch processes all transactions inside a batch
func (b *BasicBatchProcessor) ProcessBatch(batch *Batch) error {
	var receipts []*Receipt
	var cumulativeGasUsed uint64 = 0
	var gasUsed uint64
	var index uint
	var includedTxs []*types.Transaction

	for _, tx := range batch.Transactions {
		senderAddress, err := getSender(tx)
		if err != nil {
			log.Warnf("Error processing transaction %s: %v", tx.Hash().String(), err)
		}
		receiverAddress := tx.To()

		if err := b.processTransaction(tx, *senderAddress, *receiverAddress, batch.Sequencer); err != nil {
			log.Warnf("Error processing transaction %s: %v", tx.Hash().String(), err)
			// gasUsed = 0
		} else {
			log.Infof("Successfully processed transaction %s", tx.Hash().String())

			includedTxs = append(includedTxs, tx)
			gasUsed = b.State.EstimateGas(tx)
			cumulativeGasUsed += gasUsed

			// Set TX Receipt
			receipt := &Receipt{}
			receipt.Type = tx.Type()
			receipt.PostState = b.stateRoot
			receipt.Status = types.ReceiptStatusSuccessful
			receipt.CumulativeGasUsed = cumulativeGasUsed
			receipt.BlockNumber = new(big.Int).SetUint64(batch.BlockNumber)
			receipt.GasUsed = gasUsed
			receipt.TxHash = tx.Hash()
			receipt.TransactionIndex = uint(index)
			if senderAddress != nil {
				receipt.From = *senderAddress
			}
			if receiverAddress != nil {
				receipt.To = *receiverAddress
			}

			// Add receipt to the list of receipts
			receipts = append(receipts, receipt)
			index = index + 1
		}
	}

	// Update batch
	batch.Transactions = includedTxs
	batch.Receipts = receipts

	// Set batch Header
	header := &types.Header{}
	batch.Header = header
	batch.Header.ParentHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	batch.Header.UncleHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	batch.Header.Coinbase = batch.Sequencer
	batch.Header.Root = common.BytesToHash(b.stateRoot)
	batch.Header.TxHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	batch.Header.ReceiptHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	batch.Header.Bloom = types.BytesToBloom([]byte{0})
	batch.Header.Difficulty = new(big.Int).SetUint64(0)
	batch.Header.Number = new(big.Int).SetUint64(batch.BlockNumber)
	batch.Header.GasLimit = 30000000
	batch.Header.GasUsed = cumulativeGasUsed
	batch.Header.Time = 0
	// batch.Header.Extra = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	batch.Header.MixDigest = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	batch.Header.Nonce = types.BlockNonce{0, 0, 0, 0, 0, 0, 0, 0}

	// Store batch
	err := b.commit(batch)

	return err
}

// ProcessTransaction processes a transaction inside a batch
func (b *BasicBatchProcessor) processTransaction(tx *types.Transaction, senderAddress, receiverAddress common.Address, sequencerAddress common.Address) error {
	log.Debugf("processing transaction [%s]: start", tx.Hash().Hex())

	txb, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	encoded := hex.EncodeToHex(txb)
	log.Debugf("processing transaction [%s]: raw: %v", tx.Hash().Hex(), encoded)

	// save stateRoot and modify it only if transaction processing finishes successfully
	root := b.stateRoot

	// reset MT currentRoot in case it was modified by failed transaction
	b.State.tree.SetCurrentRoot(root)
	log.Debugf("processing transaction [%s]: root: %v", tx.Hash().Hex(), new(big.Int).SetBytes(root).String())

	senderBalance, err := b.State.tree.GetBalance(senderAddress, b.stateRoot)
	if err != nil {
		return err
	}

	err = b.checkTransaction(tx, senderBalance)
	if err != nil {
		return err
	}

	log.Debugf("processing transaction [%s]: sender: %v", tx.Hash().Hex(), senderAddress.Hex())
	log.Debugf("processing transaction [%s]: nonce: %v", tx.Hash().Hex(), tx.Nonce())
	log.Debugf("processing transaction [%s]: sender balance: %v", tx.Hash().Hex(), senderBalance.Text(encoding.Base10))

	// Increase Nonce
	nonce := big.NewInt(0).SetUint64(tx.Nonce())
	nonce = big.NewInt(0).Add(nonce, big.NewInt(1))
	log.Debugf("processing transaction [%s]: new nonce: %v", tx.Hash().Hex(), nonce.Text(encoding.Base10))

	// Store new nonce
	_, _, err = b.State.tree.SetNonce(senderAddress, nonce)
	if err != nil {
		return err
	}

	// Calculate Gas
	usedGas := new(big.Int).SetUint64(b.State.EstimateGas(tx))
	usedGasValue := new(big.Int).Mul(usedGas, tx.GasPrice())
	remainingGas := new(big.Int).SetUint64(tx.Gas() - usedGas.Uint64())
	unusedGasValue := new(big.Int).Mul(remainingGas, tx.GasPrice())
	log.Debugf("processing transaction [%s]: used gas: %v", tx.Hash().Hex(), usedGas.Text(encoding.Base10))
	log.Debugf("processing transaction [%s]: remaining gas: %v", tx.Hash().Hex(), remainingGas.Text(encoding.Base10))

	// Calculate new balances
	cost := tx.Cost()
	log.Debugf("processing transaction [%s]: cost: %v", tx.Hash().Hex(), cost.Text(encoding.Base10))
	value := tx.Value()
	log.Debugf("processing transaction [%s]: value: %v", tx.Hash().Hex(), value.Text(encoding.Base10))

	// Sender has to pay transaction cost
	senderBalance.Sub(senderBalance, cost)
	log.Debugf("processing transaction [%s]: sender balance after cost charged: %v", tx.Hash().Hex(), senderBalance.Text(encoding.Base10))

	if sequencerAddress == senderAddress {
		senderBalance.Add(senderBalance, usedGasValue)
	}

	// Refund unused gas
	senderBalance.Add(senderBalance, unusedGasValue)
	log.Debugf("processing transaction [%s]: sender balance after refund: %v", tx.Hash().Hex(), senderBalance.Text(encoding.Base10))

	// Store new sender balances
	root, _, err = b.State.tree.SetBalance(senderAddress, senderBalance)
	if err != nil {
		return err
	}

	// Get receiver Balance
	receiverBalance, err := b.State.tree.GetBalance(receiverAddress, root)
	if err != nil {
		return err
	}
	log.Debugf("processing transaction [%s]: receiver balance: %v", tx.Hash().Hex(), receiverBalance.Text(encoding.Base10))

	receiverBalance.Add(receiverBalance, value)
	log.Debugf("processing transaction [%s]: receiver balance after value added: %v", tx.Hash().Hex(), receiverBalance.Text(encoding.Base10))

	// Pay gas to the sequencer
	if sequencerAddress == receiverAddress && senderAddress != receiverAddress {
		receiverBalance.Add(receiverBalance, usedGasValue)
	}

	if sequencerAddress != senderAddress && sequencerAddress != receiverAddress {
		sequencerBalance, err := b.State.tree.GetBalance(sequencerAddress, root)
		if err != nil {
			return err
		}

		// Store sequencer balance
		sequencerBalance.Add(sequencerBalance, usedGasValue)
		_, _, err = b.State.tree.SetBalance(sequencerAddress, sequencerBalance)
		if err != nil {
			return err
		}
	}

	// Store receiver balance
	root, _, err = b.State.tree.SetBalance(*tx.To(), receiverBalance)
	if err != nil {
		return err
	}

	b.stateRoot = root
	log.Debugf("processing transaction [%s]: new root: %v", tx.Hash().Hex(), new(big.Int).SetBytes(root).String())

	return nil
}

// CheckTransaction checks if a transaction is valid
func (b *BasicBatchProcessor) CheckTransaction(tx *types.Transaction) error {
	sender, err := getSender(tx)
	if err != nil {
		return err
	}

	balance, err := b.State.tree.GetBalance(*sender, b.stateRoot)
	if err != nil {
		return err
	}

	return b.checkTransaction(tx, balance)
}

func (b *BasicBatchProcessor) checkTransaction(tx *types.Transaction, senderBalance *big.Int) error {
	var nonce = big.NewInt(0)

	// reset MT currentRoot in case it was modified by failed transaction
	b.State.tree.SetCurrentRoot(b.stateRoot)

	// Check Signature
	if err := CheckSignature(tx); err != nil {
		return err
	}

	// Check ChainID
	if tx.ChainId().Uint64() != b.SequencerChainID && tx.ChainId().Uint64() != b.State.cfg.DefaultChainID {
		log.Debugf("Batch ChainID: %v", b.SequencerChainID)
		log.Debugf("Transaction ChainID: %v", tx.ChainId().Uint64())
		return ErrInvalidChainID
	}

	// Get Sender
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		return err
	}

	// Check nonce
	nonce, err = b.State.tree.GetNonce(sender, b.stateRoot)
	if err != nil {
		return err
	}

	if nonce.Uint64() != tx.Nonce() {
		log.Debugf("check transaction [%s]: invalid nonce, expected: %d, found: %d", tx.Hash().Hex(), nonce.Uint64(), tx.Nonce())
		return ErrInvalidNonce
	}

	// Check balance
	if senderBalance.Cmp(tx.Cost()) < 0 {
		log.Debugf("check transaction [%s]: invalid balance, expected: %v, found: %v", tx.Hash().Hex(), tx.Cost().Text(encoding.Base10), senderBalance.Text(encoding.Base10))
		return ErrInvalidBalance
	}

	// Check gas
	gasEstimation := b.State.EstimateGas(tx)
	if tx.Gas() < gasEstimation {
		log.Debugf("check transaction [%s]: invalid gas, expected: %v, found: %v", tx.Hash().Hex(), tx.Gas(), gasEstimation)
		return ErrInvalidGas
	}

	return nil
}

// Commit the batch state into state
func (b *BasicBatchProcessor) commit(batch *Batch) error {
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

	err := b.State.AddBatch(ctx, batch)
	if err != nil {
		return err
	}

	// store transactions
	for i, tx := range batch.Transactions {
		err := b.State.AddTransaction(ctx, tx, batch.BatchNumber, uint(i))
		if err != nil {
			return err
		}
	}

	blockHash := batch.Hash()

	// store receipts
	for _, receipt := range batch.Receipts {
		receipt.BlockHash = blockHash
		err := b.State.AddReceipt(ctx, receipt)
		if err != nil {
			return err
		}
	}

	return nil
}

func getSender(tx *types.Transaction) (*common.Address, error) {
	// Get Sender
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		return &common.Address{}, err
	}
	return &sender, nil
}
