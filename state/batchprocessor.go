package state

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/runtime"
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
	// EmptyCodeHash is the hash of empty code
	EmptyCodeHash = common.Hex2Bytes("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
)

// BatchProcessor is used to process a batch of transactions
type BatchProcessor interface {
	ProcessBatch(batch *Batch) error
	CheckTransaction(tx *types.Transaction) error
	runtime.Host
}

// BasicBatchProcessor is used to process a batch of transactions
type BasicBatchProcessor struct {
	State            *BasicState
	stateRoot        []byte
	runtimes         []runtime.Runtime
	forks            runtime.ForksInTime
	SequencerAddress common.Address
	SequencerChainID uint64
}

// ProcessBatch processes all transactions inside a batch
func (b *BasicBatchProcessor) ProcessBatch(batch *Batch) error {
	var result *runtime.ExecutionResult
	var receipts []*Receipt
	var includedTxs []*types.Transaction

	var cumulativeGasUsed uint64 = 0
	var index uint

	for _, tx := range batch.Transactions {
		senderAddress, err := getSender(tx)
		if err != nil {
			log.Warnf("Error processing transaction %s: %v", tx.Hash().String(), err)
		}

		receiverAddress := tx.To()

		if receiverAddress == nil {
			result = b.run(nil)
		} else {
			result = b.transfer(tx, *senderAddress, *receiverAddress, batch.Sequencer)
		}

		if result.Err != nil {
			log.Warnf("Error processing transaction %s: %v", tx.Hash().String(), result.Err)
		} else {
			log.Infof("Successfully processed transaction %s", tx.Hash().String())

			cumulativeGasUsed += result.GasUsed
			includedTxs = append(includedTxs, tx)
			receipt := b.generateReceipt(batch.BlockNumber, tx, index, senderAddress, receiverAddress, result, cumulativeGasUsed)
			receipts = append(receipts, receipt)
			index++
		}
	}

	// Update batch
	batch.Transactions = includedTxs
	batch.Receipts = receipts

	// Set batch Header
	header := b.generateBatchHeader(batch.BlockNumber, batch.Sequencer, cumulativeGasUsed)
	batch.Header = header

	// Store batch
	err := b.commit(batch)

	return err
}

func (b *BasicBatchProcessor) generateReceipt(blockNumber uint64, tx *types.Transaction, index uint, senderAddress *common.Address, receiverAddress *common.Address, result *runtime.ExecutionResult, cumulativeGasUsed uint64) *Receipt {
	receipt := &Receipt{}
	receipt.Type = tx.Type()
	receipt.PostState = b.stateRoot
	receipt.Status = types.ReceiptStatusSuccessful
	receipt.CumulativeGasUsed = cumulativeGasUsed
	receipt.BlockNumber = new(big.Int).SetUint64(blockNumber)
	receipt.GasUsed = result.GasUsed
	receipt.TxHash = tx.Hash()
	receipt.TransactionIndex = index
	if senderAddress != nil {
		receipt.From = *senderAddress
	}
	if receiverAddress != nil {
		receipt.To = *receiverAddress
	}

	return receipt
}

func (b *BasicBatchProcessor) generateBatchHeader(blockNumber uint64, sequencerAddress common.Address, cumulativeGasUsed uint64) *types.Header {
	header := &types.Header{}
	header.ParentHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	header.UncleHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	header.Coinbase = sequencerAddress
	header.Root = common.BytesToHash(b.stateRoot)
	header.TxHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	header.ReceiptHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	header.Bloom = types.BytesToBloom([]byte{0})
	header.Difficulty = new(big.Int).SetUint64(0)
	header.Number = new(big.Int).SetUint64(blockNumber)
	header.GasLimit = 30000000
	header.GasUsed = cumulativeGasUsed
	header.Time = 0
	// header.Extra = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	header.MixDigest = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	header.Nonce = types.BlockNonce{0, 0, 0, 0, 0, 0, 0, 0}

	return header
}

// ProcessTransaction processes a transaction inside a batch
func (b *BasicBatchProcessor) transfer(tx *types.Transaction, senderAddress, receiverAddress common.Address, sequencerAddress common.Address) *runtime.ExecutionResult {
	log.Debugf("processing transaction [%s]: start", tx.Hash().Hex())
	var result *runtime.ExecutionResult = &runtime.ExecutionResult{}

	txb, err := tx.MarshalBinary()
	if err != nil {
		result.Err = err
		return result
	}
	encoded := hex.EncodeToHex(txb)
	log.Debugf("processing transaction [%s]: raw: %v", tx.Hash().Hex(), encoded)

	// save stateRoot and modify it only if transaction processing finishes successfully
	root := b.stateRoot

	// reset MT currentRoot in case it was modified by failed transaction
	b.State.tree.SetCurrentRoot(root)
	log.Debugf("processing transaction [%s]: root: %v", tx.Hash().Hex(), new(big.Int).SetBytes(root).String())

	senderBalance, err := b.State.tree.GetBalance(senderAddress, root)
	if err != nil {
		result.Err = err
		return result
	}

	senderNonce, err := b.State.tree.GetNonce(senderAddress, root)
	if err != nil {
		result.Err = err
		return result
	}

	err = b.checkTransaction(tx, senderNonce, senderBalance)
	if err != nil {
		result.Err = err
		return result
	}

	log.Debugf("processing transaction [%s]: sender: %v", tx.Hash().Hex(), senderAddress.Hex())
	log.Debugf("processing transaction [%s]: nonce: %v", tx.Hash().Hex(), tx.Nonce())
	log.Debugf("processing transaction [%s]: sender balance: %v", tx.Hash().Hex(), senderBalance.Text(encoding.Base10))

	// Increase Nonce
	senderNonce = big.NewInt(0).Add(senderNonce, big.NewInt(1))
	log.Debugf("processing transaction [%s]: new nonce: %v", tx.Hash().Hex(), senderNonce.Text(encoding.Base10))

	// Store new nonce
	_, _, err = b.State.tree.SetNonce(senderAddress, senderNonce)
	if err != nil {
		result.Err = err
		return result
	}
	log.Debugf("processing transaction [%s]: sender nonce set to: %v", tx.Hash().Hex(), senderNonce.Text(encoding.Base10))

	// Calculate Gas
	usedGas := new(big.Int).SetUint64(b.State.EstimateGas(tx))
	usedGasValue := new(big.Int).Mul(usedGas, tx.GasPrice())
	gasLeft := new(big.Int).SetUint64(tx.Gas() - usedGas.Uint64())
	unusedGasValue := new(big.Int).Mul(gasLeft, tx.GasPrice())
	log.Debugf("processing transaction [%s]: used gas: %v", tx.Hash().Hex(), usedGas.Text(encoding.Base10))
	log.Debugf("processing transaction [%s]: remaining gas: %v", tx.Hash().Hex(), gasLeft.Text(encoding.Base10))

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
		result.Err = err
		return result
	}

	// Get receiver Balance
	receiverBalance, err := b.State.tree.GetBalance(receiverAddress, root)
	if err != nil {
		result.Err = err
		return result
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
			result.Err = err
			return result
		}

		// Store sequencer balance
		sequencerBalance.Add(sequencerBalance, usedGasValue)
		_, _, err = b.State.tree.SetBalance(sequencerAddress, sequencerBalance)
		if err != nil {
			result.Err = err
			return result
		}
	}

	// Store receiver balance
	root, _, err = b.State.tree.SetBalance(*tx.To(), receiverBalance)
	if err != nil {
		result.Err = err
		return result
	}

	b.stateRoot = root
	log.Debugf("processing transaction [%s]: new root: %v", tx.Hash().Hex(), new(big.Int).SetBytes(root).String())

	result.GasUsed = usedGas.Uint64()
	result.GasLeft = gasLeft.Uint64()

	return result
}

// CheckTransaction checks if a transaction is valid
func (b *BasicBatchProcessor) CheckTransaction(tx *types.Transaction) error {
	senderAddress, err := getSender(tx)
	if err != nil {
		return err
	}

	senderNonce, err := b.State.tree.GetNonce(*senderAddress, b.stateRoot)
	if err != nil {
		return err
	}

	balance, err := b.State.tree.GetBalance(*senderAddress, b.stateRoot)
	if err != nil {
		return err
	}

	return b.checkTransaction(tx, senderNonce, balance)
}

func (b *BasicBatchProcessor) checkTransaction(tx *types.Transaction, senderNonce, senderBalance *big.Int) error {
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

	// Check nonce
	if senderNonce.Uint64() != tx.Nonce() {
		log.Debugf("check transaction [%s]: invalid nonce, expected: %d, found: %d", tx.Hash().Hex(), senderNonce.Uint64(), tx.Nonce())
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

func (b *BasicBatchProcessor) setRuntime(r runtime.Runtime) {
	b.runtimes = append(b.runtimes, r)
}

func (b *BasicBatchProcessor) run(contract *runtime.Contract) *runtime.ExecutionResult {
	for _, r := range b.runtimes {
		if r.CanRun(contract, b, &b.forks) {
			return r.Run(contract, b, &b.forks)
		}
	}

	return &runtime.ExecutionResult{
		Err: fmt.Errorf("not found"),
	}
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

// AccountExists check if the address already exists in the state
func (b *BasicBatchProcessor) AccountExists(address common.Address) bool {
	panic("not implemented")
}

// GetStorage gets the value stored in a given address and key
func (b *BasicBatchProcessor) GetStorage(address common.Address, key common.Hash) common.Hash {
	storage, err := b.State.tree.GetStorageAt(address, key, b.stateRoot)

	if err != nil {
		log.Errorf("error on GetStorage for address %v", address)
	}

	return common.BytesToHash(storage.Bytes())
}

// SetStorage sets storage for a given address
func (b *BasicBatchProcessor) SetStorage(address common.Address, key common.Hash, value common.Hash, config *runtime.ForksInTime) runtime.StorageStatus {
	// TODO: Check if we have to charge here
	root, _, err := b.State.tree.SetStorageAt(address, key, new(big.Int).SetBytes(value.Bytes()))

	if err != nil {
		log.Errorf("error on SetStorage for address %v", address)
	} else {
		b.stateRoot = root
	}

	// TODO: calculate and return proper value
	return runtime.StorageModified
}

// GetBalance gets balance for a given address
func (b *BasicBatchProcessor) GetBalance(address common.Address) *big.Int {
	balance, err := b.State.tree.GetBalance(address, b.stateRoot)

	if err != nil {
		log.Errorf("error on GetBalance for address %v", address)
	}

	return balance
}

// GetCodeSize gets the size of the code at a given address
func (b *BasicBatchProcessor) GetCodeSize(address common.Address) int {
	code := b.GetCode(address)
	return len(code)
}

// GetCodeHash gets the hash for the code at a given address
func (b *BasicBatchProcessor) GetCodeHash(address common.Address) common.Hash {
	hash, err := b.State.tree.GetCodeHash(address, b.stateRoot)

	if err != nil {
		log.Errorf("error on GetCodeHash for address %v", address)
	}

	return common.BytesToHash(hash)
}

// GetCode gets the code stored at a given address
func (b *BasicBatchProcessor) GetCode(address common.Address) []byte {
	code, err := b.State.tree.GetCode(address, b.stateRoot)

	if err != nil {
		log.Errorf("error on GetCode for address %v", address)
	}

	return code
}

// Selfdestruct deletes a contract and refunds gas
func (b *BasicBatchProcessor) Selfdestruct(address common.Address, beneficiary common.Address) {
	panic("not implemented")
}

// GetTxContext returns metadata related to the Tx Context
func (b *BasicBatchProcessor) GetTxContext() runtime.TxContext {
	panic("not implemented")
}

// GetBlockHash gets the hash of a block
func (b *BasicBatchProcessor) GetBlockHash(number int64) common.Hash {
	panic("not implemented")
}

// EmitLog generates logs
func (b *BasicBatchProcessor) EmitLog(address common.Address, topics []common.Hash, data []byte) {
	panic("not implemented")
}

// Callx calls a SC
func (b *BasicBatchProcessor) Callx(*runtime.Contract, runtime.Host) *runtime.ExecutionResult {
	panic("not implemented")
}

// Empty check whether a address is empty
func (b *BasicBatchProcessor) Empty(address common.Address) bool {
	nonce := b.GetNonce(address)
	balance := b.GetBalance(address)
	codehash := b.GetCodeHash(address)

	return nonce == 0 && balance.Int64() == 0 && codehash == EmptyCodeHash
}

// GetNonce gets the nonce for an account at a given address
func (b *BasicBatchProcessor) GetNonce(address common.Address) uint64 {
	nonce, err := b.State.tree.GetBalance(address, b.stateRoot)

	if err != nil {
		log.Errorf("error on GetNonce for address %v", address)
	}

	return nonce.Uint64()
}
