package state

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/helper"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/tree"
)

// Host implements host interface
type Host struct {
	State              *State
	stateRoot          []byte
	transactionContext transactionContext
	logs               map[common.Hash][]*types.Log
	forks              runtime.ForksInTime
	runtimes           []runtime.Runtime
	txBundleID         string
}

type transactionContext struct {
	simulationMode     bool
	currentTransaction *types.Transaction
	currentOrigin      common.Address
	coinBase           common.Address
	index              uint
	batchNumber        int64
}

// AccountExists check if the address already exists in the state
func (h *Host) AccountExists(ctx context.Context, address common.Address) bool {
	log.Debugf("AccountExists for address %v", address)
	return !h.Empty(ctx, address)
}

// GetStorage gets the value stored in a given address and key
func (h *Host) GetStorage(ctx context.Context, address common.Address, key *big.Int) common.Hash {
	storage, err := h.State.tree.GetStorageAt(ctx, address, key, h.stateRoot, h.txBundleID)

	if err != nil {
		log.Errorf("error on GetStorage for address %v", address)
	}

	log.Debugf("GetStorage for address %v", address)

	return common.BytesToHash(storage.Bytes())
}

// SetStorage sets storage for a given address
func (h *Host) SetStorage(ctx context.Context, address common.Address, key *big.Int, value *big.Int, config *runtime.ForksInTime) runtime.StorageStatus {
	// TODO: Check if we have to charge here
	root, _, err := h.State.tree.SetStorageAt(ctx, address, key, value, h.stateRoot, h.txBundleID)

	if err != nil {
		log.Errorf("error on SetStorage for address %v", address)
	} else {
		h.stateRoot = root
	}

	// TODO: calculate and return proper value
	log.Debugf("SetStorage for address %v", address)
	return runtime.StorageModified
}

// GetBalance gets balance for a given address
func (h *Host) GetBalance(ctx context.Context, address common.Address) *big.Int {
	balance, err := h.State.tree.GetBalance(ctx, address, h.stateRoot, h.txBundleID)

	if err != nil {
		log.Errorf("error on GetBalance for address %v", address)
	}

	log.Debugf("GetBalance for address %v", address)
	return balance
}

// GetCodeSize gets the size of the code at a given address
func (h *Host) GetCodeSize(ctx context.Context, address common.Address) int {
	code := h.GetCode(ctx, address)

	log.Debugf("GetCodeSize for address %v, len = %v", address, len(code))
	return len(code)
}

// GetCodeHash gets the hash for the code at a given address
func (h *Host) GetCodeHash(ctx context.Context, address common.Address) common.Hash {
	hash, err := h.State.tree.GetCodeHash(ctx, address, h.stateRoot, h.txBundleID)

	if err != nil {
		log.Errorf("error on GetCodeHash for address %v, err: %v", address, err)
	}

	log.Debugf("GetCodeHash for address %v => %v", address, common.BytesToHash(hash))
	return common.BytesToHash(hash)
}

// GetCode gets the code stored at a given address
func (h *Host) GetCode(ctx context.Context, address common.Address) []byte {
	code, err := h.State.tree.GetCode(ctx, address, h.stateRoot, h.txBundleID)

	if err != nil {
		log.Errorf("error on GetCode for address %v", address)
	}

	log.Debugf("GetCode for address %v", address)
	return code
}

// Selfdestruct deletes a contract and refunds gas
func (h *Host) Selfdestruct(ctx context.Context, address common.Address, beneficiary common.Address) {
	contractBalance := h.GetBalance(ctx, address)
	if contractBalance.Int64() != 0 {
		beneficiaryBalance := h.GetBalance(ctx, beneficiary)
		beneficiaryBalance.Add(beneficiaryBalance, contractBalance)
		root, _, err := h.State.tree.SetBalance(ctx, beneficiary, beneficiaryBalance, h.stateRoot, h.txBundleID)
		if err != nil {
			log.Errorf("error on Selfdestuct for address %v", address)
		}
		root, _, err = h.State.tree.SetBalance(ctx, beneficiary, big.NewInt(0), root, h.txBundleID)
		if err != nil {
			log.Errorf("error on Selfdestuct for address %v", address)
		}
		h.stateRoot = root
	}

	root, _, err := h.State.tree.SetCode(ctx, address, []byte{}, h.stateRoot, h.txBundleID)
	if err != nil {
		log.Errorf("error on Selfdestuct for address %v", address)
	}
	h.stateRoot = root

	// Storage not destroyed as per Protocol definition
}

// GetTxContext returns metadata related to the Tx Context
func (h *Host) GetTxContext() runtime.TxContext {
	return runtime.TxContext{
		Hash:        h.transactionContext.currentTransaction.Hash(),
		GasPrice:    common.BigToHash(h.transactionContext.currentTransaction.GasPrice()),
		Origin:      h.transactionContext.currentOrigin,
		Coinbase:    h.transactionContext.coinBase,
		Number:      int64(h.transactionContext.index),
		Timestamp:   time.Now().Unix(),
		GasLimit:    int64(h.transactionContext.currentTransaction.Gas()),
		ChainID:     h.transactionContext.currentTransaction.ChainId().Int64(),
		BatchNumber: h.transactionContext.batchNumber,
	}
}

// GetBlockHash gets the hash of a block (batch in L2)
func (h *Host) GetBlockHash(number int64) common.Hash {
	batch, err := h.State.GetBatchByNumber(context.Background(), uint64(number), h.txBundleID)

	if err != nil {
		log.Errorf("error on GetBlockHash for number %v", number)
	}

	log.Debugf("GetBlockHash for number %v", number)
	return batch.Hash()
}

// EmitLog generates logs
func (h *Host) EmitLog(address common.Address, topics []common.Hash, data []byte) {
	if !h.transactionContext.simulationMode {
		log.Debugf("EmitLog for address %v", address)

		txLog := &types.Log{
			Address: address,
			Topics:  topics,
			Data:    common.CopyBytes(data),
			TxHash:  h.transactionContext.currentTransaction.Hash(),
			TxIndex: h.transactionContext.index,
			Index:   h.getLogIndex(),
			Removed: false,
		}

		if _, found := h.logs[h.transactionContext.currentTransaction.Hash()]; !found {
			h.logs[h.transactionContext.currentTransaction.Hash()] = []*types.Log{}
		}

		h.logs[h.transactionContext.currentTransaction.Hash()] = append(h.logs[h.transactionContext.currentTransaction.Hash()], txLog)
	}
}

func (h *Host) getLogIndex() uint {
	nextIndex := 0
	for _, l := range h.logs {
		nextIndex += len(l)
	}
	return uint(nextIndex)
}

// Callx calls a SC
func (h *Host) Callx(ctx context.Context, contract *runtime.Contract, host runtime.Host) *runtime.ExecutionResult {
	log.Debugf("Callx to address %v", contract.CodeAddress)

	if contract.Type == runtime.Create {
		log.Debugf("Callx. New Contract Creation %v", contract.Address)
		return h.applyCreate(ctx, contract, host)
	}

	if contract.Type == runtime.Call && contract.Value.Uint64() != 0 {
		log.Debugf("Callx. New Transfer from %v to %v", contract.Caller, contract.Address)
		err := h.transfer(ctx, contract.Caller, contract.Address, contract.Value)
		if err != nil {
			return &runtime.ExecutionResult{
				GasLeft: contract.Gas,
				Err:     err,
			}
		}
	}

	var contract2 *runtime.Contract

	if contract.Type == runtime.DelegateCall {
		contract2 = runtime.NewContractCall(contract.Depth+1, contract.Caller, contract.CodeAddress, contract.Address, contract.Value, contract.Gas, contract.Code, contract.Input)
	} else {
		contract2 = runtime.NewContractCall(contract.Depth+1, contract.Address, contract.Caller, contract.CodeAddress, contract.Value, contract.Gas, contract.Code, contract.Input)
	}

	result := h.run(ctx, contract2)

	return result
}

// Empty check whether an address is empty
func (h *Host) Empty(ctx context.Context, address common.Address) bool {
	log.Debugf("Empty for address %v", address)
	return h.GetNonce(ctx, address) == 0 && h.GetBalance(ctx, address).Int64() == 0 && h.GetCodeHash(ctx, address) == EmptyCodeHash
}

// GetNonce gets the nonce for an account at a given address
func (h *Host) GetNonce(ctx context.Context, address common.Address) uint64 {
	nonce, err := h.State.tree.GetNonce(ctx, address, h.stateRoot, h.txBundleID)

	if err != nil {
		log.Errorf("error on GetNonce for address %v", address)
	}

	log.Debugf("GetNonce for address %v", address)

	return nonce.Uint64()
}

func (h *Host) transfer(ctx context.Context, senderAddress, receiverAddress common.Address, value *big.Int) error {
	var err error
	var balances = make(map[common.Address]*big.Int)

	root := h.stateRoot

	senderBalance := h.GetBalance(ctx, senderAddress)
	balances[senderAddress] = senderBalance

	receiverBalance := h.GetBalance(ctx, receiverAddress)
	balances[receiverAddress] = receiverBalance

	balances[senderAddress].Sub(balances[senderAddress], value)
	balances[receiverAddress].Add(balances[receiverAddress], value)

	// Store new balances
	for address, balance := range balances {
		root, _, err = h.State.tree.SetBalance(ctx, address, balance, root, h.txBundleID)
		if err != nil {
			return err
		}
	}

	h.stateRoot = root

	return nil
}

func (h *Host) applyCreate(ctx context.Context, contract *runtime.Contract, host runtime.Host) *runtime.ExecutionResult {
	gasLimit := contract.Gas
	root := h.stateRoot

	senderNonce, err := h.State.tree.GetNonce(ctx, contract.Caller, root, h.txBundleID)
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
	if !h.Empty(ctx, contract.Address) {
		return &runtime.ExecutionResult{
			GasLeft: 0,
			Err:     runtime.ErrContractAddressCollision,
		}
	}

	// Increment nonce of the sender
	senderNonce.Add(senderNonce, big.NewInt(1))

	// Store new nonce
	root, _, err = h.State.tree.SetNonce(ctx, contract.Caller, senderNonce, root, h.txBundleID)
	if err != nil {
		return &runtime.ExecutionResult{
			GasLeft: 0,
			Err:     err,
		}
	}

	h.stateRoot = root

	log.Debugf("apply create: run contract with address %v", contract.Address)
	result := h.run(ctx, contract)
	if result.Failed() {
		log.Debugf("apply create failed for address %v", contract.Address)
		return result
	}
	log.Debugf("apply create: run contract with address %v finished", contract.Address)
	// Update root with the result after SC Execution
	root = h.stateRoot

	if h.forks.EIP158 && len(result.ReturnValue) > spuriousDragonMaxCodeSize {
		// Contract size exceeds 'SpuriousDragon' size limit
		return &runtime.ExecutionResult{
			GasLeft: 0,
			Err:     runtime.ErrMaxCodeSizeExceeded,
		}
	}

	gasCost := uint64(len(result.ReturnValue)) * contractByteGasCost
	/*
		if result.GasLeft < gasCost {
			result.Err = runtime.ErrCodeStoreOutOfGas
			result.ReturnValue = nil

			// Out of gas creating the contract
			if b.forks.Homestead {
				result.GasLeft = 0
			}
			return result
		}
	*/
	result.GasLeft -= gasCost
	root, _, err = h.State.tree.SetCode(ctx, contract.Address, result.ReturnValue, root, h.txBundleID)
	if err != nil {
		return &runtime.ExecutionResult{
			GasLeft: gasLimit,
			Err:     err,
		}
	}

	result.CreateAddress = contract.Address
	result.GasUsed = gasCost
	h.stateRoot = root

	log.Debugf("apply create finished for address %v", contract.Address)

	return result
}

func (h *Host) setRuntime(r runtime.Runtime) {
	h.runtimes = append(h.runtimes, r)
}

func (h *Host) run(ctx context.Context, contract *runtime.Contract) *runtime.ExecutionResult {
	for _, r := range h.runtimes {
		if r.CanRun(contract, h, &h.forks) {
			result := r.Run(ctx, contract, h, &h.forks)
			if result.Err != nil {
				delete(h.logs, h.transactionContext.currentTransaction.Hash())
			}
			return result
		}
	}

	return &runtime.ExecutionResult{
		Err: fmt.Errorf("not found"),
	}
}

// GetOldStateRoot gets an old state root from the System SC
func (h *Host) GetOldStateRoot(ctx context.Context, batchNumber int64) int64 {
	// Recover old state root from System SC
	batchNumberBytes := tree.ScalarToFilledByteSlice(new(big.Int).SetInt64(batchNumber))
	storagePosition := tree.ScalarToFilledByteSlice(new(big.Int).SetUint64(h.State.cfg.OldStateRootPosition))
	oldStateRootPosition := helper.Keccak256(batchNumberBytes, storagePosition)

	oldRoot, err := h.State.tree.GetStorageAt(ctx, h.State.cfg.SystemSCAddr, new(big.Int).SetBytes(oldStateRootPosition), h.stateRoot, h.txBundleID)
	if err != nil {
		log.Errorf("error on GetOldStateRoot for batchNumber %v", batchNumber)
	}

	log.Debugf("GetOldStateRoot for bathNumber %v", batchNumber)

	return oldRoot.Int64()
}
