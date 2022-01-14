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
	ProcessTransaction(tx *types.Transaction, sequencerAddress common.Address) (*common.Address, *common.Address, error)
	CheckTransaction(tx *types.Transaction) (common.Address, *big.Int, *big.Int, error)
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
	for _, tx := range batch.Transactions {
		from, to, err := b.ProcessTransaction(tx, batch.Sequencer)
		if err != nil {
			log.Warnf("Error processing transaction %s: %v", tx.Hash().String(), err)
			// gasUsed = 0
		} else {
			log.Infof("Successfully processed transaction %s", tx.Hash().String())
			gasUsed = b.State.EstimateGas(tx)
			cumulativeGasUsed += gasUsed
			// Set TX Receipt
			// For performance reasons blockhash is set outside this method
			// before adding receipt to database
			receipt := &Receipt{}
			receipt.BlockHash = batch.Hash()
			receipt.Type = tx.Type()
			receipt.PostState = b.stateRoot
			receipt.Status = types.ReceiptStatusSuccessful
			receipt.CumulativeGasUsed = cumulativeGasUsed
			receipt.BlockNumber = new(big.Int).SetUint64(batch.BlockNumber)
			receipt.GasUsed = gasUsed
			receipt.TxHash = tx.Hash()
			receipt.TransactionIndex = uint(index)
			if from != nil {
				receipt.From = *from
			}
			if to != nil {
				receipt.To = *to
			}
			receipts = append(receipts, receipt)
			index = index + 1
		}
	}

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

	batch.Receipts = receipts
	_, err := b.commit(batch)

	return err
}

// ProcessTransaction processes a transaction inside a batch
func (b *BasicBatchProcessor) ProcessTransaction(tx *types.Transaction, sequencerAddress common.Address) (*common.Address, *common.Address, error) {
	log.Debugf("processing transaction [%s]: start", tx.Hash().Hex())

	txb, err := tx.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	encoded := hex.EncodeToHex(txb)
	log.Debugf("processing transaction [%s]: raw: %v", tx.Hash().Hex(), encoded)

	// save stateRoot and modify it only if transaction processing finishes successfully
	root := b.stateRoot

	// reset MT currentRoot in case it was modified by failed transaction
	b.State.tree.SetCurrentRoot(root)
	log.Debugf("processing transaction [%s]: root: %v", tx.Hash().Hex(), new(big.Int).SetBytes(root).String())

	sender, nonce, senderBalance, err := b.CheckTransaction(tx)
	if err != nil {
		return nil, nil, err
	}
	log.Debugf("processing transaction [%s]: sender: %v", tx.Hash().Hex(), sender.Hex())
	log.Debugf("processing transaction [%s]: nonce: %v", tx.Hash().Hex(), nonce.Text(encoding.Base10))
	log.Debugf("processing transaction [%s]: sender balance: %v", tx.Hash().Hex(), senderBalance.Text(encoding.Base10))

	// Increase Nonce
	nonce.Add(nonce, big.NewInt(1))
	log.Debugf("processing transaction [%s]: new nonce: %v", tx.Hash().Hex(), nonce.Text(encoding.Base10))

	// Store new nonce
	_, _, err = b.State.tree.SetNonce(sender, nonce)
	if err != nil {
		return nil, nil, err
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

	if sequencerAddress == sender {
		senderBalance.Add(senderBalance, usedGasValue)
	}

	// Refund unused gas
	senderBalance.Add(senderBalance, unusedGasValue)
	log.Debugf("processing transaction [%s]: sender balance after refund: %v", tx.Hash().Hex(), senderBalance.Text(encoding.Base10))

	// Store new sender balances
	root, _, err = b.State.tree.SetBalance(sender, senderBalance)
	if err != nil {
		return nil, nil, err
	}

	// Get receiver Balance
	receiverBalance, err := b.State.tree.GetBalance(*tx.To(), root)
	if err != nil {
		return nil, nil, err
	}
	log.Debugf("processing transaction [%s]: receiver balance: %v", tx.Hash().Hex(), receiverBalance.Text(encoding.Base10))

	receiverBalance.Add(receiverBalance, value)
	log.Debugf("processing transaction [%s]: receiver balance after value added: %v", tx.Hash().Hex(), receiverBalance.Text(encoding.Base10))

	// Pay gas to the sequencer
	if sequencerAddress == *tx.To() && sender != *tx.To() {
		receiverBalance.Add(receiverBalance, usedGasValue)
	}

	if sequencerAddress != sender && sequencerAddress != *tx.To() {
		sequencerBalance, err := b.State.tree.GetBalance(sequencerAddress, root)
		if err != nil {
			return nil, nil, err
		}

		// Store sequencer balance
		sequencerBalance.Add(sequencerBalance, usedGasValue)
		_, _, err = b.State.tree.SetBalance(sequencerAddress, sequencerBalance)
		if err != nil {
			return nil, nil, err
		}
	}

	// Store receiver balance
	root, _, err = b.State.tree.SetBalance(*tx.To(), receiverBalance)
	if err != nil {
		return nil, nil, err
	}

	b.stateRoot = root
	log.Debugf("processing transaction [%s]: new root: %v", tx.Hash().Hex(), new(big.Int).SetBytes(root).String())

	return &sender, tx.To(), nil
}

// CheckTransaction checks if a transaction is valid
func (b *BasicBatchProcessor) CheckTransaction(tx *types.Transaction) (common.Address, *big.Int, *big.Int, error) {
	var sender = common.Address{}
	var nonce = big.NewInt(0)
	var balance = big.NewInt(0)

	// reset MT currentRoot in case it was modified by failed transaction
	b.State.tree.SetCurrentRoot(b.stateRoot)

	// Check Signature
	if err := CheckSignature(tx); err != nil {
		return sender, nonce, balance, err
	}

	// Check ChainID
	if tx.ChainId().Uint64() != b.SequencerChainID && tx.ChainId().Uint64() != b.State.cfg.DefaultChainID {
		log.Debugf("Batch ChainID: %v", b.SequencerChainID)
		log.Debugf("Transaction ChainID: %v", tx.ChainId().Uint64())
		return sender, nonce, balance, ErrInvalidChainID
	}

	// Get Sender
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		return sender, nonce, balance, err
	}

	// Check nonce
	nonce, err = b.State.tree.GetNonce(sender, b.stateRoot)
	if err != nil {
		return sender, nonce, balance, err
	}

	if nonce.Uint64() != tx.Nonce() {
		return sender, nonce, balance, ErrInvalidNonce
	}

	// Check balance
	balance, err = b.State.tree.GetBalance(sender, b.stateRoot)
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
func (b *BasicBatchProcessor) commit(batch *Batch) (*common.Hash, error) {
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
		return nil, err
	}

	// store transactions
	for i, tx := range batch.Transactions {
		err := b.State.AddTransaction(ctx, tx, batch.BatchNumber, uint(i))
		if err != nil {
			return nil, err
		}
	}

	blockHash := batch.Hash()

	// store receipts
	for _, receipt := range batch.Receipts {
		receipt.BlockHash = blockHash
		err := b.State.AddReceipt(ctx, receipt)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
