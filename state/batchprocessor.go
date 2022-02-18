package state

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/helper"
	"github.com/hermeznetwork/hermez-core/state/runtime"
)

const (
	spuriousDragonMaxCodeSize = 24576
	maxCallDepth              = 1024
	contractByteGasCost       = 200
)

var (
	// ErrInvalidSig indicates the signature of the transaction is not valid
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
	// ErrNonceIsBiggerThanAccountNonce indicates the nonce of the transaction is bigger than account nonce
	ErrNonceIsBiggerThanAccountNonce = errors.New("transaction nonce is bigger than account nonce")
	// ErrNonceIsSmallerThanAccountNonce indicates the nonce of the transaction is smaller than account nonce
	ErrNonceIsSmallerThanAccountNonce = errors.New("transaction nonce is smaller than account nonce")
	// ErrInvalidBalance indicates the balance of the account is not enough to process the transaction
	ErrInvalidBalance = errors.New("not enough balance")
	// ErrInvalidGas indicates the gaslimit is not enough to process the transaction
	ErrInvalidGas = errors.New("not enough gas")
	// ErrInvalidChainID indicates a mismatch between sequencer address and ChainID
	ErrInvalidChainID = errors.New("invalid chain id for sequencer")
	// ErrNotImplemented indicates this feature has not yet been implemented
	ErrNotImplemented = errors.New("feature not yet implemented")
	// ErrInvalidTxType indicates the tx type is not known
	ErrInvalidTxType = errors.New("unknown transaction type")
	// ErrInvalidCumulativeGas indicates the batch gas is bigger than the max allowed
	ErrInvalidCumulativeGas = errors.New("cumulative gas is bigger than allowed")
	// EmptyCodeHash is the hash of empty code
	EmptyCodeHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	// ZeroAddress is the address 0x0000000000000000000000000000000000000000
	ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

// InvalidTxErrors is map to spot invalid txs
var InvalidTxErrors = map[string]bool{
	ErrInvalidSig.Error(): true, ErrNonceIsSmallerThanAccountNonce.Error(): true, ErrInvalidBalance.Error(): true,
	ErrInvalidGas.Error(): true, ErrInvalidChainID.Error(): true,
}

// BatchProcessor is used to process a batch of transactions
type BatchProcessor interface {
	ProcessBatch(batch *Batch) error
	CheckTransaction(tx *types.Transaction) error
	ProcessTransaction(tx *types.Transaction, sequencerAddress common.Address) *runtime.ExecutionResult
	ProcessUnsignedTransaction(tx *types.Transaction, senderAddress, sequencerAddress common.Address) *runtime.ExecutionResult
	runtime.Host
}

// BasicBatchProcessor is used to process a batch of transactions
type BasicBatchProcessor struct {
	State                *BasicState
	stateRoot            []byte
	runtimes             []runtime.Runtime
	forks                runtime.ForksInTime
	SequencerAddress     common.Address
	SequencerChainID     uint64
	LastBatch            *Batch
	CumulativeGasUsed    uint64
	MaxCumulativeGasUsed uint64
}

// ProcessBatch processes all transactions inside a batch
func (b *BasicBatchProcessor) ProcessBatch(batch *Batch) error {
	var receipts []*Receipt
	var includedTxs []*types.Transaction
	var index uint

	b.CumulativeGasUsed = 0

	for _, tx := range batch.Transactions {
		senderAddress, err := helper.GetSender(tx)
		if err != nil {
			return err
		}

		result := b.processTransaction(tx, senderAddress, batch.Sequencer)

		if result.Err != nil {
			log.Warnf("Error processing transaction %s: %v", tx.Hash().String(), result.Err)
		} else {
			log.Infof("Successfully processed transaction %s", tx.Hash().String())
			b.CumulativeGasUsed += result.GasUsed
			includedTxs = append(includedTxs, tx)
			receipt := b.generateReceipt(batch, tx, index, &senderAddress, tx.To(), result.GasUsed)
			receipts = append(receipts, receipt)
			index++
		}
	}

	// Update batch
	batch.Transactions = includedTxs
	batch.Receipts = receipts

	// Set batch Header
	b.populateBatchHeader(batch)

	// Store batch
	return b.commit(batch)
}

// ProcessTransaction processes a transaction
func (b *BasicBatchProcessor) ProcessTransaction(tx *types.Transaction, sequencerAddress common.Address) *runtime.ExecutionResult {
	senderAddress, err := helper.GetSender(tx)
	if err != nil {
		return &runtime.ExecutionResult{Err: err}
	}

	// Keep track of consumed gas
	result := b.processTransaction(tx, senderAddress, sequencerAddress)
	b.CumulativeGasUsed += result.GasUsed

	if b.CumulativeGasUsed > b.MaxCumulativeGasUsed {
		result.Err = ErrInvalidCumulativeGas
	}

	return result
}

// ProcessUnsignedTransaction processes an unsigned transaction from the given
// sender.
func (b *BasicBatchProcessor) ProcessUnsignedTransaction(tx *types.Transaction, senderAddress, sequencerAddress common.Address) *runtime.ExecutionResult {
	return b.processTransaction(tx, senderAddress, sequencerAddress)
}

func (b *BasicBatchProcessor) processTransaction(tx *types.Transaction, senderAddress, sequencerAddress common.Address) *runtime.ExecutionResult {
	receiverAddress := tx.To()

	// SC creation?
	if *receiverAddress == ZeroAddress && len(tx.Data()) > 0 {
		log.Debug("smart contract creation")
		return b.create(tx, senderAddress, sequencerAddress)
	}

	// SC execution
	code := b.GetCode(*receiverAddress)
	if len(code) > 0 {
		log.Debugf("smart contract execution %v", receiverAddress)
		contract := runtime.NewContractCall(0, senderAddress, senderAddress, *receiverAddress, tx.Value(), tx.Gas(), code, tx.Data())
		result := b.run(contract)
		result.GasUsed = tx.Gas() - result.GasLeft
		log.Debugf("Transaction Data %v", tx.Data())
		log.Debugf("Returned value from execution: %v", "0x"+hex.EncodeToString(result.ReturnValue))
		return result
	}

	// Transfer
	if tx.Value() != new(big.Int) {
		return b.transfer(tx, senderAddress, *receiverAddress, sequencerAddress)
	}

	log.Error("unknown transaction type")
	return &runtime.ExecutionResult{Err: ErrInvalidTxType}
}

func (b *BasicBatchProcessor) populateBatchHeader(batch *Batch) {
	parentHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	if b.LastBatch != nil {
		parentHash = b.LastBatch.Hash()
	}

	rr := make([]*types.Receipt, 0, len(batch.Receipts))
	for _, receipt := range batch.Receipts {
		r := receipt.Receipt
		rr = append(rr, &r)
	}
	block := types.NewBlock(batch.Header, batch.Transactions, batch.Uncles, rr, &trie.StackTrie{})

	batch.Header.ParentHash = parentHash
	batch.Header.UncleHash = block.UncleHash()
	batch.Header.Coinbase = batch.Sequencer
	batch.Header.Root = common.BytesToHash(b.stateRoot)
	batch.Header.TxHash = block.Header().TxHash
	batch.Header.ReceiptHash = block.Header().ReceiptHash
	batch.Header.Bloom = block.Bloom()
	batch.Header.Difficulty = new(big.Int).SetUint64(0)
	batch.Header.GasLimit = 30000000
	batch.Header.GasUsed = b.CumulativeGasUsed
	batch.Header.Time = uint64(time.Now().Unix())
	batch.Header.MixDigest = block.MixDigest()
	batch.Header.Nonce = block.Header().Nonce
}

func (b *BasicBatchProcessor) generateReceipt(batch *Batch, tx *types.Transaction, index uint, senderAddress *common.Address, receiverAddress *common.Address, gasUsed uint64) *Receipt {
	receipt := &Receipt{}
	receipt.Type = tx.Type()
	receipt.PostState = b.stateRoot
	receipt.Status = types.ReceiptStatusSuccessful
	receipt.CumulativeGasUsed = b.CumulativeGasUsed
	receipt.BlockNumber = batch.Number()
	receipt.BlockHash = batch.Hash()
	receipt.GasUsed = gasUsed
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

// transfer processes a transfer transaction
func (b *BasicBatchProcessor) transfer(tx *types.Transaction, senderAddress, receiverAddress, sequencerAddress common.Address) *runtime.ExecutionResult {
	log.Debugf("processing transfer [%s]: start", tx.Hash().Hex())
	var result *runtime.ExecutionResult = &runtime.ExecutionResult{}
	var balances = make(map[common.Address]*big.Int)

	txb, err := tx.MarshalBinary()
	if err != nil {
		result.Err = err
		return result
	}
	encoded := hex.EncodeToHex(txb)
	log.Debugf("processing transfer [%s]: raw: %v", tx.Hash().Hex(), encoded)

	// save stateRoot and modify it only if transaction processing finishes successfully
	root := b.stateRoot

	// reset MT currentRoot in case it was modified by failed transaction
	b.State.tree.SetCurrentRoot(root)
	log.Debugf("processing transfer [%s]: root: %v", tx.Hash().Hex(), new(big.Int).SetBytes(root).String())

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

	balances[senderAddress] = senderBalance

	log.Debugf("processing transfer [%s]: sender: %v", tx.Hash().Hex(), senderAddress.Hex())
	log.Debugf("processing transfer [%s]: nonce: %v", tx.Hash().Hex(), tx.Nonce())
	log.Debugf("processing transfer [%s]: sender balance: %v", tx.Hash().Hex(), senderBalance.Text(encoding.Base10))

	// Increase Nonce
	senderNonce = big.NewInt(0).Add(senderNonce, big.NewInt(1))
	log.Debugf("processing transfer [%s]: new nonce: %v", tx.Hash().Hex(), senderNonce.Text(encoding.Base10))

	// Store new nonce
	_, _, err = b.State.tree.SetNonce(senderAddress, senderNonce)
	if err != nil {
		result.Err = err
		return result
	}
	log.Debugf("processing transfer [%s]: sender nonce set to: %v", tx.Hash().Hex(), senderNonce.Text(encoding.Base10))

	// Get receiver Balance
	receiverBalance, err := b.State.tree.GetBalance(receiverAddress, root)
	if err != nil {
		result.Err = err
		return result
	}
	log.Debugf("processing transfer [%s]: receiver balance: %v", tx.Hash().Hex(), receiverBalance.Text(encoding.Base10))
	balances[receiverAddress] = receiverBalance

	// Get sequencer Balance
	sequencerBalance, err := b.State.tree.GetBalance(sequencerAddress, root)
	if err != nil {
		result.Err = err
		return result
	}
	balances[sequencerAddress] = sequencerBalance

	// Calculate Gas
	usedGas := new(big.Int).SetUint64(b.State.EstimateGas(tx))
	usedGasValue := new(big.Int).Mul(usedGas, tx.GasPrice())
	gasLeft := new(big.Int).SetUint64(tx.Gas() - usedGas.Uint64())
	gasLeftValue := new(big.Int).Mul(gasLeft, tx.GasPrice())
	log.Debugf("processing transfer [%s]: used gas: %v", tx.Hash().Hex(), usedGas.Text(encoding.Base10))
	log.Debugf("processing transfer [%s]: remaining gas: %v", tx.Hash().Hex(), gasLeft.Text(encoding.Base10))

	// Calculate new balances
	cost := tx.Cost()
	log.Debugf("processing transfer [%s]: cost: %v", tx.Hash().Hex(), cost.Text(encoding.Base10))
	value := tx.Value()
	log.Debugf("processing transfer [%s]: value: %v", tx.Hash().Hex(), value.Text(encoding.Base10))

	// Sender has to pay transaction cost
	balances[senderAddress].Sub(balances[senderAddress], cost)
	log.Debugf("processing transfer [%s]: sender balance after cost charged: %v", tx.Hash().Hex(), balances[senderAddress].Text(encoding.Base10))

	// Refund unused gas to sender
	balances[senderAddress].Add(balances[senderAddress], gasLeftValue)
	log.Debugf("processing transfer [%s]: sender balance after refund: %v", tx.Hash().Hex(), balances[senderAddress].Text(encoding.Base10))

	// Add value to receiver
	balances[receiverAddress].Add(balances[receiverAddress], value)
	log.Debugf("processing transfer [%s]: receiver balance after value added: %v", tx.Hash().Hex(), balances[receiverAddress].Text(encoding.Base10))

	// Pay gas to the sequencer
	balances[sequencerAddress].Add(balances[sequencerAddress], usedGasValue)

	// Store new balances
	for address, balance := range balances {
		root, _, err = b.State.tree.SetBalance(address, balance)
		if err != nil {
			result.Err = err
			return result
		}
	}

	b.stateRoot = root
	log.Debugf("processing transfer [%s]: new root: %v", tx.Hash().Hex(), new(big.Int).SetBytes(root).String())

	result.GasUsed = usedGas.Uint64()
	result.GasLeft = gasLeft.Uint64()

	return result
}

// CheckTransaction checks if a transaction is valid
func (b *BasicBatchProcessor) CheckTransaction(tx *types.Transaction) error {
	senderAddress, err := helper.GetSender(tx)
	if err != nil {
		return err
	}

	senderNonce, err := b.State.tree.GetNonce(senderAddress, b.stateRoot)
	if err != nil {
		return err
	}

	balance, err := b.State.tree.GetBalance(senderAddress, b.stateRoot)
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
	if senderNonce.Uint64() > tx.Nonce() {
		log.Debugf("check transaction [%s]: invalid nonce, tx nonce is smaller than account nonce, expected: %d, found: %d",
			tx.Hash().Hex(), senderNonce.Uint64(), tx.Nonce())
		return ErrNonceIsSmallerThanAccountNonce
	}

	if senderNonce.Uint64() < tx.Nonce() {
		log.Debugf("check transaction [%s]: invalid nonce at this moment, tx nonce is bigger than account nonce, expected: %d, found: %d",
			tx.Hash().Hex(), senderNonce.Uint64(), tx.Nonce())
		return ErrNonceIsBiggerThanAccountNonce
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
			Number:     batch.Number(),
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
		err := b.State.AddTransaction(ctx, tx, batch.Number().Uint64(), uint(i))
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

func (b *BasicBatchProcessor) create(tx *types.Transaction, senderAddress, sequencerAddress common.Address) *runtime.ExecutionResult {
	address := helper.CreateAddress(senderAddress, tx.Nonce())
	contract := runtime.NewContractCreation(0, senderAddress, senderAddress, address, tx.Value(), tx.Gas(), tx.Data())

	log.Debugf("new contract address = %v", address)

	root := b.stateRoot
	gasLimit := contract.Gas

	senderNonce, err := b.State.tree.GetNonce(senderAddress, root)
	if err != nil {
		return &runtime.ExecutionResult{
			GasLeft: 0,
			Err:     err,
		}
	}

	senderBalance, err := b.State.tree.GetBalance(senderAddress, root)
	if err != nil {
		return &runtime.ExecutionResult{
			GasLeft: 0,
			Err:     err,
		}
	}

	err = b.checkTransaction(tx, senderNonce, senderBalance)
	if err != nil {
		return &runtime.ExecutionResult{
			GasLeft: 0,
			Err:     err,
		}
	}

	if contract.Depth > int(maxCallDepth)+1 {
		return &runtime.ExecutionResult{
			GasLeft: gasLimit,
			Err:     runtime.ErrDepth,
		}
	}

	// Check if there if there is a collision and the address already exists
	if !b.Empty(contract.Address) {
		return &runtime.ExecutionResult{
			GasLeft: 0,
			Err:     runtime.ErrContractAddressCollision,
		}
	}

	if tx.Value().Uint64() != 0 {
		log.Debugf("contract creation includes value transfer = %v", tx.Value())
		// Tansfer the value
		transferResult := b.transfer(tx, senderAddress, contract.Address, sequencerAddress)
		if transferResult.Err != nil {
			return &runtime.ExecutionResult{
				GasLeft: gasLimit,
				Err:     transferResult.Err,
			}
		}
	} else {
		// Increment nonce of the sender
		senderNonce.Add(senderNonce, big.NewInt(1))

		// Store new nonce
		_, _, err = b.State.tree.SetNonce(senderAddress, senderNonce)
		if err != nil {
			return &runtime.ExecutionResult{
				GasLeft: 0,
				Err:     err,
			}
		}
	}

	result := b.run(contract)
	if result.Failed() {
		return result
	}

	if b.forks.EIP158 && len(result.ReturnValue) > spuriousDragonMaxCodeSize {
		// Contract size exceeds 'SpuriousDragon' size limit
		return &runtime.ExecutionResult{
			GasLeft: 0,
			Err:     runtime.ErrMaxCodeSizeExceeded,
		}
	}

	gasCost := uint64(len(result.ReturnValue)) * contractByteGasCost

	if result.GasLeft < gasCost {
		result.Err = runtime.ErrCodeStoreOutOfGas
		result.ReturnValue = nil

		// Out of gas creating the contract
		if b.forks.Homestead {
			result.GasLeft = 0
		}

		return result
	}

	result.GasLeft -= gasCost
	root, _, err = b.State.tree.SetCode(address, result.ReturnValue)
	if err != nil {
		return &runtime.ExecutionResult{
			GasLeft: gasLimit,
			Err:     err,
		}
	}

	result.CreateAddress = address
	b.stateRoot = root

	return result
}

// AccountExists check if the address already exists in the state
func (b *BasicBatchProcessor) AccountExists(address common.Address) bool {
	// TODO: Implement this properly, may need to modify the MT
	log.Debugf("AccountExists for address %v", address)
	return !b.Empty(address)
}

// GetStorage gets the value stored in a given address and key
func (b *BasicBatchProcessor) GetStorage(address common.Address, key common.Hash) common.Hash {
	storage, err := b.State.tree.GetStorageAt(address, key, b.stateRoot)

	if err != nil {
		log.Errorf("error on GetStorage for address %v", address)
	}

	log.Debugf("GetStorage for address %v", address)
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
	log.Debugf("SetStorage for address %v", address)
	return runtime.StorageModified
}

// GetBalance gets balance for a given address
func (b *BasicBatchProcessor) GetBalance(address common.Address) *big.Int {
	balance, err := b.State.tree.GetBalance(address, b.stateRoot)

	if err != nil {
		log.Errorf("error on GetBalance for address %v", address)
	}

	log.Debugf("GetBalance for address %v", address)
	return balance
}

// GetCodeSize gets the size of the code at a given address
func (b *BasicBatchProcessor) GetCodeSize(address common.Address) int {
	code := b.GetCode(address)

	log.Debugf("GetCodeSize for address %v", address)
	return len(code)
}

// GetCodeHash gets the hash for the code at a given address
func (b *BasicBatchProcessor) GetCodeHash(address common.Address) common.Hash {
	hash, err := b.State.tree.GetCodeHash(address, b.stateRoot)

	if err != nil {
		log.Errorf("error on GetCodeHash for address %v", address)
	}

	log.Debugf("GetCodeHash for address %v => %v", address, common.BytesToHash(hash))
	return common.BytesToHash(hash)
}

// GetCode gets the code stored at a given address
func (b *BasicBatchProcessor) GetCode(address common.Address) []byte {
	code, err := b.State.tree.GetCode(address, b.stateRoot)

	if err != nil {
		log.Errorf("error on GetCode for address %v", address)
	}

	log.Debugf("GetCode for address %v", address)
	return code
}

// Selfdestruct deletes a contract and refunds gas
func (b *BasicBatchProcessor) Selfdestruct(address common.Address, beneficiary common.Address) {
	// TODO: Implement
	panic("not implemented")
}

// GetTxContext returns metadata related to the Tx Context
func (b *BasicBatchProcessor) GetTxContext() runtime.TxContext {
	// TODO: Implement
	panic("not implemented")
}

// GetBlockHash gets the hash of a block (batch in L2)
func (b *BasicBatchProcessor) GetBlockHash(number int64) common.Hash {
	batch, err := b.State.GetBatchByNumber(context.Background(), uint64(number))

	if err != nil {
		log.Errorf("error on GetBlockHash for number %v", number)
	}

	log.Debugf("GetBlockHash for number %v", number)
	return batch.Hash()
}

// EmitLog generates logs
func (b *BasicBatchProcessor) EmitLog(address common.Address, topics []common.Hash, data []byte) {
	// TODO: Implement
	log.Warn("batchprocessor.EmitLog not implemented")
}

// Callx calls a SC
func (b *BasicBatchProcessor) Callx(contract *runtime.Contract, host runtime.Host) *runtime.ExecutionResult {
	log.Debugf("Callx to address %v", contract.CodeAddress)
	root := b.stateRoot
	contract2 := runtime.NewContractCall(contract.Depth+1, contract.Address, contract.Caller, contract.CodeAddress, contract.Value, contract.Gas, contract.Code, contract.Input)
	result := b.run(contract2)
	if result.Reverted() {
		b.stateRoot = root
	}
	return result
}

// Empty check whether an address is empty
func (b *BasicBatchProcessor) Empty(address common.Address) bool {
	log.Debugf("Empty for address %v", address)
	return b.GetNonce(address) == 0 && b.GetBalance(address).Int64() == 0 && b.GetCodeHash(address) == EmptyCodeHash
}

// GetNonce gets the nonce for an account at a given address
func (b *BasicBatchProcessor) GetNonce(address common.Address) uint64 {
	nonce, err := b.State.tree.GetNonce(address, b.stateRoot)

	if err != nil {
		log.Errorf("error on GetNonce for address %v", address)
	}

	log.Debugf("GetNonce for address %v", address)
	return nonce.Uint64()
}
