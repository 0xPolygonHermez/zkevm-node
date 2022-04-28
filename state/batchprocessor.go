package state

import (
	"context"
	"errors"
	"fmt"
	"math"
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
	nonZeroCost               = 68
	zeroCost                  = 4
)

var (
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
	// ErrIntrinsicGasOverflow indicates overflow during gas estimation
	ErrIntrinsicGasOverflow = fmt.Errorf("overflow in intrinsic gas estimation")
)

// InvalidTxErrors is map to spot invalid txs
var InvalidTxErrors = map[string]bool{
	ErrInvalidSig.Error(): true, ErrNonceIsSmallerThanAccountNonce.Error(): true, ErrInvalidBalance.Error(): true,
	ErrInvalidGas.Error(): true, ErrInvalidChainID.Error(): true,
}

// BatchProcessor is used to process a batch of transactions
type BatchProcessor struct {
	SequencerAddress     common.Address
	SequencerChainID     uint64
	LastBatch            *Batch
	CumulativeGasUsed    uint64
	MaxCumulativeGasUsed uint64
	Host                 Host
}

// SetSimulationMode allows execution without updating the state
func (b *BatchProcessor) SetSimulationMode(active bool) {
	b.Host.transactionContext.simulationMode = active
}

// ProcessBatch processes all transactions inside a batch
func (b *BatchProcessor) ProcessBatch(ctx context.Context, batch *Batch) error {
	var receipts []*Receipt
	var includedTxs []*types.Transaction
	var index uint

	b.CumulativeGasUsed = 0
	b.Host.logs = map[common.Hash][]*types.Log{}

	// Set Global Exit Root storage position
	var batchNumberBuf, storagePositionBuf [32]byte
	batchNumber := batch.Number().FillBytes(batchNumberBuf[:])
	storagePosition := new(big.Int).SetUint64(b.Host.State.cfg.GlobalExitRootStoragePosition).FillBytes(storagePositionBuf[:])
	globalExitRootPos := helper.Keccak256(batchNumber, storagePosition)

	root, _, err := b.Host.State.tree.SetStorageAt(ctx, b.Host.State.cfg.L2GlobalExitRootManagerAddr, new(big.Int).SetBytes(globalExitRootPos), new(big.Int).SetBytes(batch.GlobalExitRoot.Bytes()), b.Host.stateRoot, b.Host.txBundleID)
	if err != nil {
		return err
	}

	b.Host.stateRoot = root

	for _, tx := range batch.Transactions {
		senderAddress, err := helper.GetSender(*tx)
		if err != nil {
			return err
		}

		// Set transaction context
		b.Host.transactionContext.index = index
		b.Host.transactionContext.difficulty = batch.Header.Difficulty

		result := b.processTransaction(ctx, tx, senderAddress, batch.Sequencer)

		if result.Err != nil {
			log.Warnf("Error processing transaction %s: %v", tx.Hash().String(), result.Err)
		} else {
			log.Infof("Successfully processed transaction %s", tx.Hash().String())
		}

		if result.Succeeded() || result.Reverted() {
			b.CumulativeGasUsed += result.GasUsed
			includedTxs = append(includedTxs, tx)
			receipt := b.generateReceipt(batch, tx, index, &senderAddress, tx.To(), result)
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
	return b.commit(ctx, batch)
}

// ProcessTransaction processes a transaction
func (b *BatchProcessor) ProcessTransaction(ctx context.Context, tx *types.Transaction, sequencerAddress common.Address) *runtime.ExecutionResult {
	senderAddress, err := helper.GetSender(*tx)
	if err != nil {
		return &runtime.ExecutionResult{Err: err, StateRoot: b.Host.stateRoot}
	}

	result := b.processTransaction(ctx, tx, senderAddress, sequencerAddress)

	if !b.Host.transactionContext.simulationMode {
		b.CumulativeGasUsed += result.GasUsed

		if b.CumulativeGasUsed > b.MaxCumulativeGasUsed {
			result.Err = ErrInvalidCumulativeGas
		}
	}
	return result
}

// ProcessUnsignedTransaction processes an unsigned transaction from the given
// sender.
func (b *BatchProcessor) ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress, sequencerAddress common.Address) *runtime.ExecutionResult {
	return b.processTransaction(ctx, tx, senderAddress, sequencerAddress)
}

func (b *BatchProcessor) estimateGas(ctx context.Context, tx *types.Transaction) *runtime.ExecutionResult {
	result := &runtime.ExecutionResult{Err: nil, GasUsed: 0}

	cost := uint64(0)

	if b.IsContractCreation(tx) {
		cost += TxSmartContractCreationGas
	} else {
		cost += TxTransferGas
	}

	if !b.IsTransfer(ctx, tx) {
		payload := tx.Data()

		if len(payload) > 0 {
			zeros := uint64(0)

			for i := 0; i < len(payload); i++ {
				if payload[i] == 0 {
					zeros++
				}
			}

			nonZeros := uint64(len(payload)) - zeros

			if (math.MaxUint64-cost)/nonZeroCost < nonZeros {
				result.Err = ErrIntrinsicGasOverflow
				return result
			}

			cost += nonZeros * nonZeroCost

			if (math.MaxUint64-cost)/4 < zeros {
				result.Err = ErrIntrinsicGasOverflow
			}

			cost += zeros * zeroCost
			cost = cost + cost + cost + cost + cost // temporary solution
		}
	}

	result.GasUsed = cost

	return result
}

// IsContractCreation checks if the tx is a contract creation
func (b *BatchProcessor) IsContractCreation(tx *types.Transaction) bool {
	return tx.To() == nil && len(tx.Data()) > 0
}

// IsSmartContractExecution checks if the tx is a contract execution
func (b *BatchProcessor) IsSmartContractExecution(ctx context.Context, tx *types.Transaction) bool {
	return b.Host.GetCodeHash(ctx, *tx.To()) != EmptyCodeHash
}

// IsTransfer checks if the tx is a transfer
func (b *BatchProcessor) IsTransfer(ctx context.Context, tx *types.Transaction) bool {
	return !b.IsContractCreation(tx) && !b.IsSmartContractExecution(ctx, tx) && tx.Value().Uint64() != 0
}

func (b *BatchProcessor) processTransaction(ctx context.Context, tx *types.Transaction, senderAddress, sequencerAddress common.Address) *runtime.ExecutionResult {
	log.Debugf("Processing tx: %v", tx.Hash())

	// Set transaction context
	b.Host.transactionContext.currentTransaction = tx
	b.Host.transactionContext.currentOrigin = senderAddress
	b.Host.transactionContext.coinBase = sequencerAddress
	receiverAddress := tx.To()

	if b.IsContractCreation(tx) {
		log.Debug("smart contract creation")
		return b.create(ctx, tx, senderAddress, sequencerAddress)
	}

	if b.IsSmartContractExecution(ctx, tx) {
		return b.execute(ctx, tx, senderAddress, *receiverAddress, sequencerAddress)
	}

	if b.IsTransfer(ctx, tx) {
		return b.transfer(ctx, tx, senderAddress, *receiverAddress, sequencerAddress)
	}

	log.Error("unknown transaction type")
	return &runtime.ExecutionResult{Err: ErrInvalidTxType, StateRoot: b.Host.stateRoot}
}

func (b *BatchProcessor) populateBatchHeader(batch *Batch) {
	parentHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	if b.LastBatch != nil {
		parentHash = b.LastBatch.Hash()
	}

	rr := make([]*types.Receipt, 0, len(batch.Receipts))
	for _, receipt := range batch.Receipts {
		r := receipt.Receipt
		r.Logs = append(r.Logs, b.Host.logs[receipt.TxHash]...)
		rr = append(rr, &r)
	}
	block := types.NewBlock(batch.Header, batch.Transactions, batch.Uncles, rr, &trie.StackTrie{})

	batch.Header.ParentHash = parentHash
	batch.Header.UncleHash = block.UncleHash()
	batch.Header.Coinbase = batch.Sequencer
	batch.Header.Root = common.BytesToHash(b.Host.stateRoot)
	batch.Header.TxHash = block.TxHash()
	batch.Header.ReceiptHash = block.ReceiptHash()
	batch.Header.Bloom = block.Bloom()
	batch.Header.Difficulty = new(big.Int).SetUint64(0)
	batch.Header.GasLimit = 30000000
	batch.Header.GasUsed = b.CumulativeGasUsed
	batch.Header.Time = uint64(time.Now().Unix())
	batch.Header.MixDigest = block.MixDigest()
	batch.Header.Nonce = block.Header().Nonce
}

func (b *BatchProcessor) generateReceipt(batch *Batch, tx *types.Transaction, index uint, senderAddress *common.Address, receiverAddress *common.Address, result *runtime.ExecutionResult) *Receipt {
	receipt := &Receipt{}
	receipt.Type = tx.Type()
	receipt.PostState = b.Host.stateRoot

	if result.Succeeded() {
		receipt.Status = types.ReceiptStatusSuccessful
	} else {
		receipt.Status = types.ReceiptStatusFailed
	}

	receipt.CumulativeGasUsed = b.CumulativeGasUsed
	receipt.BlockNumber = batch.Number()
	receipt.BlockHash = batch.Hash()
	receipt.GasUsed = result.GasUsed
	receipt.TxHash = tx.Hash()
	receipt.TransactionIndex = index
	receipt.ContractAddress = result.CreateAddress
	receipt.To = receiverAddress
	if senderAddress != nil {
		receipt.From = *senderAddress
	}

	return receipt
}

// transfer processes a transfer transaction
func (b *BatchProcessor) transfer(ctx context.Context, tx *types.Transaction, senderAddress, receiverAddress, sequencerAddress common.Address) *runtime.ExecutionResult {
	log.Debugf("processing transfer [%s]: start", tx.Hash().Hex())
	var result *runtime.ExecutionResult = &runtime.ExecutionResult{}
	var balances = make(map[common.Address]*big.Int)

	txb, err := tx.MarshalBinary()
	if err != nil {
		result.Err = err
		result.StateRoot = b.Host.stateRoot
		return result
	}
	encoded := hex.EncodeToHex(txb)
	log.Debugf("processing transfer [%s]: raw: %v", tx.Hash().Hex(), encoded)

	// save stateRoot and modify it only if transaction processing finishes successfully
	root := b.Host.stateRoot

	// reset MT currentRoot in case it was modified by failed transaction
	log.Debugf("processing transfer [%s]: root: %v", tx.Hash().Hex(), new(big.Int).SetBytes(root).String())

	senderBalance, err := b.Host.State.tree.GetBalance(ctx, senderAddress, root, b.Host.txBundleID)
	if err != nil {
		if err == ErrNotFound {
			senderBalance = big.NewInt(0)
		} else {
			result.Err = err
			result.StateRoot = b.Host.stateRoot
			return result
		}
	}

	senderNonce, err := b.Host.State.tree.GetNonce(ctx, senderAddress, root, b.Host.txBundleID)
	if err != nil {
		if err == ErrNotFound {
			senderNonce = big.NewInt(0)
		} else {
			result.Err = err
			result.StateRoot = b.Host.stateRoot
			return result
		}
	}

	err = b.checkTransaction(ctx, tx, senderNonce, senderBalance)
	if err != nil {
		result.Err = err
		result.StateRoot = b.Host.stateRoot
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
	root, _, err = b.Host.State.tree.SetNonce(ctx, senderAddress, senderNonce, root, b.Host.txBundleID)
	if err != nil {
		result.Err = err
		result.StateRoot = b.Host.stateRoot
		return result
	}
	log.Debugf("processing transfer [%s]: sender nonce set to: %v", tx.Hash().Hex(), senderNonce.Text(encoding.Base10))

	// Get receiver Balance
	receiverBalance, err := b.Host.State.tree.GetBalance(ctx, receiverAddress, root, b.Host.txBundleID)
	if err != nil {
		result.Err = err
		result.StateRoot = b.Host.stateRoot
		return result
	}
	log.Debugf("processing transfer [%s]: receiver balance: %v", tx.Hash().Hex(), receiverBalance.Text(encoding.Base10))
	balances[receiverAddress] = receiverBalance

	// Get sequencer Balance
	sequencerBalance, err := b.Host.State.tree.GetBalance(ctx, sequencerAddress, root, b.Host.txBundleID)
	if err != nil {
		result.Err = err
		result.StateRoot = b.Host.stateRoot
		return result
	}
	balances[sequencerAddress] = sequencerBalance

	// Calculate Gas
	usedGas := new(big.Int).SetUint64(TxTransferGas)
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
		root, _, err = b.Host.State.tree.SetBalance(ctx, address, balance, root, b.Host.txBundleID)
		if err != nil {
			result.Err = err
			result.StateRoot = b.Host.stateRoot
			return result
		}
	}

	b.Host.stateRoot = root

	log.Debugf("processing transfer [%s]: new root: %v", tx.Hash().Hex(), new(big.Int).SetBytes(root).String())

	result.GasUsed = usedGas.Uint64()
	result.GasLeft = gasLeft.Uint64()
	result.StateRoot = b.Host.stateRoot

	return result
}

// CheckTransaction checks if a transaction is valid
func (b *BatchProcessor) CheckTransaction(ctx context.Context, tx *types.Transaction) error {
	senderAddress, err := helper.GetSender(*tx)
	if err != nil {
		return err
	}

	senderNonce, err := b.Host.State.tree.GetNonce(ctx, senderAddress, b.Host.stateRoot, b.Host.txBundleID)
	if err != nil {
		if err == ErrNotFound {
			senderNonce = big.NewInt(0)
		} else {
			return err
		}
	}

	senderBalance, err := b.Host.State.tree.GetBalance(ctx, senderAddress, b.Host.stateRoot, b.Host.txBundleID)
	if err != nil {
		if err == ErrNotFound {
			senderBalance = big.NewInt(0)
		} else {
			return err
		}
	}

	return b.checkTransaction(ctx, tx, senderNonce, senderBalance)
}

func (b *BatchProcessor) checkTransaction(ctx context.Context, tx *types.Transaction, senderNonce, senderBalance *big.Int) error {
	// Check balance
	if senderBalance.Cmp(tx.Cost()) < 0 {
		log.Debugf("check transaction [%s]: invalid balance, expected: %v, found: %v", tx.Hash().Hex(), tx.Cost().Text(encoding.Base10), senderBalance.Text(encoding.Base10))
		return ErrInvalidBalance
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

	if !b.Host.transactionContext.simulationMode {
		// Check ChainID
		if tx.ChainId().Uint64() != b.SequencerChainID && tx.ChainId().Uint64() != b.Host.State.cfg.DefaultChainID {
			log.Debugf("Batch ChainID: %v", b.SequencerChainID)
			log.Debugf("Transaction ChainID: %v", tx.ChainId().Uint64())
			return ErrInvalidChainID
		}

		// Check gas
		result := b.estimateGas(ctx, tx)
		if result.Err != nil {
			log.Debugf("check transaction [%s]: error estimating gas", tx.Hash().Hex())
			return ErrInvalidGas
		}
		if tx.Gas() < result.GasUsed {
			log.Debugf("check transaction [%s]: invalid gas, expected: %v, found: %v", tx.Hash().Hex(), tx.Gas(), result.GasUsed)
			return ErrInvalidGas
		}
	}

	return nil
}

// Commit the batch state into state
func (b *BatchProcessor) commit(ctx context.Context, batch *Batch) error {
	// Store batch into db
	var root common.Hash

	if batch.Header == nil {
		batch.Header = &types.Header{
			Root:       root,
			Difficulty: big.NewInt(0),
			Number:     batch.Number(),
		}
	}

	// set merkletree root
	if b.Host.stateRoot != nil {
		root.SetBytes(b.Host.stateRoot)
		batch.Header.Root = root

		// set local exit root
		key := new(big.Int).SetUint64(b.Host.State.cfg.LocalExitRootStoragePosition)
		localExitRoot, err := b.Host.State.tree.GetStorageAt(ctx, b.Host.State.cfg.L2GlobalExitRootManagerAddr, key, b.Host.stateRoot, b.Host.txBundleID)
		if err != nil {
			return err
		}
		batch.RollupExitRoot = common.BigToHash(localExitRoot)
	}

	err := b.Host.State.AddBatch(ctx, batch, b.Host.txBundleID)
	if err != nil {
		return err
	}

	// store transactions
	for i, tx := range batch.Transactions {
		err := b.Host.State.AddTransaction(ctx, tx, batch.Number().Uint64(), uint(i), b.Host.txBundleID)
		if err != nil {
			return err
		}
	}

	blockHash := batch.Hash()

	// store receipts
	for _, receipt := range batch.Receipts {
		receipt.BlockHash = blockHash
		err := b.Host.State.AddReceipt(ctx, receipt, b.Host.txBundleID)
		if err != nil {
			return err
		}
	}

	// store logs
	for _, tx := range batch.Transactions {
		if txLogs, found := b.Host.logs[tx.Hash()]; found {
			for _, txLog := range txLogs {
				txLog.BlockHash = blockHash
				txLog.BlockNumber = batch.Number().Uint64()
				err := b.Host.State.AddLog(ctx, *txLog, b.Host.txBundleID)
				if err != nil {
					return err
				}
			}
		}
	}

	b.LastBatch = batch

	return nil
}

func (b *BatchProcessor) execute(ctx context.Context, tx *types.Transaction, senderAddress, receiverAddress, sequencerAddress common.Address) *runtime.ExecutionResult {
	code := b.Host.GetCode(ctx, receiverAddress)
	log.Debugf("smart contract execution %v", receiverAddress)
	contract := runtime.NewContractCall(1, senderAddress, senderAddress, receiverAddress, tx.Value(), tx.Gas(), code, tx.Data())
	result := b.Host.run(ctx, contract)
	result.GasUsed = tx.Gas() - result.GasLeft

	log.Debugf("Transaction Data %v", tx.Data())
	log.Debugf("Returned value from execution: %v", "0x"+hex.EncodeToString(result.ReturnValue))

	if tx.Value().Uint64() != 0 && !result.Reverted() {
		log.Debugf("contract execution includes value transfer = %v", tx.Value())
		// Tansfer the value
		transferResult := b.transfer(ctx, tx, senderAddress, contract.Address, sequencerAddress)
		if transferResult.Err != nil {
			transferResult.StateRoot = b.Host.stateRoot
			return transferResult
		}
	} else {
		// Increment sender nonce
		senderNonce, err := b.Host.State.tree.GetNonce(ctx, senderAddress, b.Host.stateRoot, b.Host.txBundleID)
		if err != nil {
			result.Err = err
			result.StateRoot = b.Host.stateRoot
			return result
		}
		senderNonce.Add(senderNonce, big.NewInt(1))

		// Store new nonce
		root, _, err := b.Host.State.tree.SetNonce(ctx, senderAddress, senderNonce, b.Host.stateRoot, b.Host.txBundleID)
		if err != nil {
			result.Err = err
			result.StateRoot = b.Host.stateRoot
			return result
		}

		b.Host.stateRoot = root
	}

	result.StateRoot = b.Host.stateRoot

	return result
}

func (b *BatchProcessor) create(ctx context.Context, tx *types.Transaction, senderAddress, sequencerAddress common.Address) *runtime.ExecutionResult {
	root := b.Host.stateRoot

	if len(tx.Data()) <= 0 {
		return &runtime.ExecutionResult{
			GasLeft: tx.Gas(),
			Err:     runtime.ErrCodeNotFound,
		}
	}

	address := helper.CreateAddress(senderAddress, tx.Nonce())
	contract := runtime.NewContractCreation(0, senderAddress, senderAddress, address, tx.Value(), tx.Gas(), tx.Data())

	log.Debugf("new contract address = %v", address)

	gasLimit := contract.Gas

	senderNonce, err := b.Host.State.tree.GetNonce(ctx, senderAddress, root, b.Host.txBundleID)
	if err != nil {
		return &runtime.ExecutionResult{
			GasLeft:   0,
			Err:       err,
			StateRoot: b.Host.stateRoot,
		}
	}

	err = b.CheckTransaction(ctx, tx)
	if err != nil {
		return &runtime.ExecutionResult{
			GasLeft:   0,
			Err:       err,
			StateRoot: b.Host.stateRoot,
		}
	}

	if contract.Depth > int(maxCallDepth)+1 {
		return &runtime.ExecutionResult{
			GasLeft:   gasLimit,
			Err:       runtime.ErrDepth,
			StateRoot: b.Host.stateRoot,
		}
	}

	// Check if there if there is a collision and the address already exists
	if !b.Host.Empty(ctx, contract.Address) {
		return &runtime.ExecutionResult{
			GasLeft:   0,
			Err:       runtime.ErrContractAddressCollision,
			StateRoot: b.Host.stateRoot,
		}
	}

	if tx.Value().Uint64() != 0 {
		log.Debugf("contract creation includes value transfer = %v", tx.Value())
		// Tansfer the value
		transferResult := b.transfer(ctx, tx, senderAddress, contract.Address, sequencerAddress)
		if transferResult.Err != nil {
			return &runtime.ExecutionResult{
				GasLeft:   gasLimit,
				Err:       transferResult.Err,
				StateRoot: b.Host.stateRoot,
			}
		}
	} else {
		// Increment nonce of the sender
		senderNonce.Add(senderNonce, big.NewInt(1))

		// Store new nonce
		root, _, err := b.Host.State.tree.SetNonce(ctx, senderAddress, senderNonce, root, b.Host.txBundleID)
		if err != nil {
			return &runtime.ExecutionResult{
				GasLeft:   0,
				Err:       err,
				StateRoot: b.Host.stateRoot,
			}
		}

		b.Host.stateRoot = root
	}

	result := b.Host.run(ctx, contract)
	result.StateRoot = b.Host.stateRoot
	if result.Failed() {
		return result
	}
	// Update root with the result after SC Execution
	root = b.Host.stateRoot

	if b.Host.forks.EIP158 && len(result.ReturnValue) > spuriousDragonMaxCodeSize {
		// Contract size exceeds 'SpuriousDragon' size limit
		return &runtime.ExecutionResult{
			GasLeft:   0,
			Err:       runtime.ErrMaxCodeSizeExceeded,
			StateRoot: b.Host.stateRoot,
		}
	}

	gasCost := uint64(len(result.ReturnValue)) * contractByteGasCost

	if result.GasLeft < gasCost {
		result.Err = runtime.ErrCodeStoreOutOfGas
		result.ReturnValue = nil
		result.StateRoot = b.Host.stateRoot

		// Out of gas creating the contract
		if b.Host.forks.Homestead {
			result.GasLeft = 0
		}

		return result
	}

	result.GasLeft -= gasCost
	root, _, err = b.Host.State.tree.SetCode(ctx, address, result.ReturnValue, root, b.Host.txBundleID)
	if err != nil {
		return &runtime.ExecutionResult{
			GasLeft:   gasLimit,
			Err:       err,
			StateRoot: b.Host.stateRoot,
		}
	}

	result.CreateAddress = address
	result.GasUsed = gasCost
	b.Host.stateRoot = root
	result.StateRoot = root

	return result
}
