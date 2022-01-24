package state

import (
	"fmt"
	"math"
	"math/big"

	"github.com/0xPolygon/eth-state-transition/helper"
	"github.com/0xPolygon/eth-state-transition/runtime"
	"github.com/0xPolygon/eth-state-transition/runtime/evm"
	"github.com/0xPolygon/eth-state-transition/runtime/precompiled"
	"github.com/0xPolygon/eth-state-transition/types"
)

const (
	spuriousDragonMaxCodeSize = 24576

	// Per transaction not creating a contract
	TxGas uint64 = 21000

	// Per transaction that creates a contract
	TxGasContractCreation uint64 = 53000
)

var emptyCodeHashTwo = types.BytesToHash(helper.Keccak256(nil))

// GetHashByNumber returns the hash function of a block number
type GetHashByNumber = func(i uint64) types.Hash

type GetHashByNumberHelper = func(num uint64, hash types.Hash) GetHashByNumber

type Transition struct {
	runtimes []runtime.Runtime

	// forks are the enabled forks for this transition
	forks runtime.ForksInTime

	// txn is the transaction of changes
	txn *Txn

	// ctx is the block context
	ctx runtime.TxContext

	// gasPool is the current gas in the pool
	gasPool uint64

	// GetHash GetHashByNumberHelper
	getHash GetHashByNumber

	// counter on the total gas used so far
	totalGas uint64
}

// NewExecutor creates a new executor
func NewTransition(forks runtime.ForksInTime, ctx runtime.TxContext, snap Snapshot) *Transition {
	txn := NewTxn(snap)

	transition := &Transition{
		ctx:      ctx,
		txn:      txn,
		forks:    forks,
		gasPool:  uint64(ctx.GasLimit),
		totalGas: 0,
	}

	transition.setRuntime(precompiled.NewPrecompiled())
	transition.setRuntime(evm.NewEVM())

	// by default for getHash use a simple one
	transition.getHash = func(n uint64) types.Hash {
		return types.BytesToHash(helper.Keccak256([]byte(big.NewInt(int64(n)).String())))
	}

	return transition
}

func (e *Transition) Commit() []*Object {
	return e.txn.Commit()
}

func (t *Transition) TotalGas() uint64 {
	return t.totalGas
}

func (e *Transition) setRuntime(r runtime.Runtime) {
	e.runtimes = append(e.runtimes, r)
}

type BlockResult struct {
	Root     types.Hash
	Receipts []*Result
	TotalGas uint64
}

func (t *Transition) SetGetHash(helper GetHashByNumberHelper) {
	t.getHash = helper(uint64(t.ctx.Number), t.ctx.Hash)
}

func (t *Transition) subGasPool(amount uint64) error {
	if t.gasPool < amount {
		return ErrBlockLimitReached
	}
	t.gasPool -= amount
	return nil
}

func (t *Transition) addGasPool(amount uint64) {
	t.gasPool += amount
}

func (t *Transition) Txn() *Txn {
	return t.txn
}

// Write writes another transaction to the executor
func (t *Transition) Write(txn *Transaction) (*Result, error) {
	// Make a local copy and apply the transaction
	msg := txn.Copy()

	result, err := t.applyImpl(msg)
	if err != nil {
		return nil, err
	}
	t.totalGas += result.GasUsed

	logs := t.txn.Logs()

	// var root []byte

	receipt := &Result{
		// CumulativeGasUsed: t.totalGas,
		// TxHash:            txn.Hash,
		GasUsed:     result.GasUsed,
		ReturnValue: result.ReturnValue,
	}

	if t.forks.Byzantium {
		// The suicided accounts are set as deleted for the next iteration
		t.txn.CleanDeleteObjects(true)

		if result.Failed() {
			receipt.Success = false
		} else {
			receipt.Success = true
		}

	} else {
		// TODO: If byzntium is enabled you need a special step to commit the data yourself
		t.txn.CleanDeleteObjects(t.forks.EIP158)

		/*
			objs := t.txn.Commit(t.forks.EIP155)
			ss, aux := t.txn.snapshot.Commit(objs)

			t.txn = NewTxn(ss)
			root = aux
			receipt.Root = types.BytesToHash(root)
		*/
	}

	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To == nil {
		receipt.ContractAddress = helper.CreateAddress(msg.From, txn.Nonce)
	}

	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = logs

	return receipt, nil
}

// Apply applies a new transaction
func (t *Transition) applyImpl(msg *Transaction) (*runtime.ExecutionResult, error) {
	s := t.txn.Snapshot()
	result, err := t.apply(msg)
	if err != nil {
		t.txn.RevertToSnapshot(s)
	}
	return result, err
}

func (t *Transition) subGasLimitPrice(msg *Transaction) error {
	// deduct the upfront max gas cost
	upfrontGasCost := new(big.Int).Set(msg.GasPrice)
	upfrontGasCost.Mul(upfrontGasCost, new(big.Int).SetUint64(msg.Gas))

	if err := t.txn.SubBalance(msg.From, upfrontGasCost); err != nil {
		if err == runtime.ErrNotEnoughFunds {
			return ErrNotEnoughFundsForGas
		}
		return err
	}
	return nil
}

func (t *Transition) nonceCheck(msg *Transaction) error {
	nonce := t.txn.GetNonce(msg.From)

	if nonce != msg.Nonce {
		return ErrNonceIncorrect
	}
	return nil
}

var (
	ErrNonceIncorrect        = fmt.Errorf("incorrect nonce")
	ErrNotEnoughFundsForGas  = fmt.Errorf("not enough funds to cover gas costs")
	ErrBlockLimitReached     = fmt.Errorf("gas limit reached in the pool")
	ErrBlockLimitExceeded    = fmt.Errorf("transaction's gas limit exceeds block gas limit")
	ErrIntrinsicGasOverflow  = fmt.Errorf("overflow in intrinsic gas calculation")
	ErrNotEnoughIntrinsicGas = fmt.Errorf("not enough gas supplied for intrinsic gas costs")
	ErrNotEnoughFunds        = fmt.Errorf("not enough funds for transfer with given value")
)

func (t *Transition) apply(msg *Transaction) (*runtime.ExecutionResult, error) {
	txn := t.txn

	gasLeft := uint64(0)

	// First check this message satisfies all consensus rules before
	// applying the message.
	preCheck := func() error {
		// 1. the nonce of the message caller is correct
		if err := t.nonceCheck(msg); err != nil {
			return err
		}

		// 2. caller has enough balance to cover transaction fee(gaslimit * gasprice)
		if err := t.subGasLimitPrice(msg); err != nil {
			return err
		}

		// 3. the amount of gas required is available in the block
		if err := t.subGasPool(msg.Gas); err != nil {
			return err
		}

		// 4. there is no overflow when calculating intrinsic gas
		intrinsicGasCost, err := TransactionGasCost(msg, t.forks.Homestead, t.forks.Istanbul)
		if err != nil {
			return err
		}

		// 5. the purchased gas is enough to cover intrinsic usage
		gasLeft = msg.Gas - intrinsicGasCost
		// Because we are working with unsigned integers for gas, the `>` operator is used instead of the more intuitive `<`
		if gasLeft > msg.Gas {
			return ErrNotEnoughIntrinsicGas
		}

		// 6. caller has enough balance to cover asset transfer for **topmost** call
		if balance := txn.GetBalance(msg.From); balance.Cmp(msg.Value) < 0 {
			return ErrNotEnoughFunds
		}
		return nil
	}

	if err := preCheck(); err != nil {
		return nil, err
	}

	gasPrice := new(big.Int).Set(msg.GasPrice)
	value := new(big.Int).Set(msg.Value)

	// Override the context and set the specific transaction fields
	t.ctx.GasPrice = types.BytesToHash(gasPrice.Bytes())
	t.ctx.Origin = msg.From

	var result *runtime.ExecutionResult = nil
	if msg.IsContractCreation() {
		result = t.Create(msg.From, msg.Input, value, gasLeft)
	} else {
		txn.IncrNonce(msg.From)
		result = t.Call(msg.From, *msg.To, msg.Input, value, gasLeft)
	}

	// Update gas used depending on the refund.
	refund := txn.GetRefund()
	{
		result.GasUsed = msg.Gas - result.GasLeft
		maxRefund := result.GasUsed / 2
		// Refund can go up to half the gas used
		if refund > maxRefund {
			refund = maxRefund
		}

		result.GasLeft += refund
		result.GasUsed -= refund
	}

	// refund the sender
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(result.GasLeft), gasPrice)
	txn.AddBalance(msg.From, remaining)

	// pay the coinbase for the transaction
	coinbaseFee := new(big.Int).Mul(new(big.Int).SetUint64(result.GasUsed), gasPrice)
	txn.AddBalance(t.ctx.Coinbase, coinbaseFee)

	// return gas to the pool
	t.addGasPool(result.GasLeft)

	return result, nil
}

func (t *Transition) Create(caller types.Address, code []byte, value *big.Int, gas uint64) *runtime.ExecutionResult {
	address := helper.CreateAddress(caller, t.txn.GetNonce(caller))
	contract := runtime.NewContractCreation(1, caller, caller, address, value, gas, code)
	return t.applyCreate(contract, t)
}

func (t *Transition) Call(caller types.Address, to types.Address, input []byte, value *big.Int, gas uint64) *runtime.ExecutionResult {
	c := runtime.NewContractCall(1, caller, caller, to, value, gas, t.txn.GetCode(to), input)
	return t.applyCall(c, runtime.Call, t)
}

func (t *Transition) run(contract *runtime.Contract, host runtime.Host) *runtime.ExecutionResult {
	for _, r := range t.runtimes {
		if r.CanRun(contract, host, &t.forks) {
			return r.Run(contract, host, &t.forks)
		}
	}

	return &runtime.ExecutionResult{
		Err: fmt.Errorf("not found"),
	}
}

func (t *Transition) transfer(from, to types.Address, amount *big.Int) error {
	if amount == nil {
		return nil
	}

	if err := t.txn.SubBalance(from, amount); err != nil {
		if err == runtime.ErrNotEnoughFunds {
			return runtime.ErrInsufficientBalance
		}
		return err
	}

	t.txn.AddBalance(to, amount)
	return nil
}

func (t *Transition) applyCall(c *runtime.Contract, callType runtime.CallType, host runtime.Host) *runtime.ExecutionResult {
	if c.Depth > int(1024)+1 {
		return &runtime.ExecutionResult{
			GasLeft: c.Gas,
			Err:     runtime.ErrDepth,
		}
	}

	snapshot := t.txn.Snapshot()
	t.txn.TouchAccount(c.Address)

	if callType == runtime.Call {
		// Transfers only allowed on calls
		if err := t.transfer(c.Caller, c.Address, c.Value); err != nil {
			return &runtime.ExecutionResult{
				GasLeft: c.Gas,
				Err:     err,
			}
		}
	}

	result := t.run(c, host)
	if result.Failed() {
		t.txn.RevertToSnapshot(snapshot)
	}

	return result
}

var emptyHash types.Hash

func (t *Transition) hasCodeOrNonce(addr types.Address) bool {
	nonce := t.txn.GetNonce(addr)
	if nonce != 0 {
		return true
	}
	codeHash := t.txn.GetCodeHash(addr)
	if codeHash != emptyCodeHashTwo && codeHash != emptyHash {
		return true
	}
	return false
}

func (t *Transition) applyCreate(c *runtime.Contract, host runtime.Host) *runtime.ExecutionResult {
	gasLimit := c.Gas

	if c.Depth > int(1024)+1 {
		return &runtime.ExecutionResult{
			GasLeft: gasLimit,
			Err:     runtime.ErrDepth,
		}
	}

	// Increment the nonce of the caller
	t.txn.IncrNonce(c.Caller)

	// Check if there if there is a collision and the address already exists
	if t.hasCodeOrNonce(c.Address) {
		return &runtime.ExecutionResult{
			GasLeft: 0,
			Err:     runtime.ErrContractAddressCollision,
		}
	}

	// Take snapshot of the current state
	snapshot := t.txn.Snapshot()

	if t.forks.EIP158 {
		// Force the creation of the account
		t.txn.CreateAccount(c.Address)
		t.txn.IncrNonce(c.Address)
	}

	// Transfer the value
	if err := t.transfer(c.Caller, c.Address, c.Value); err != nil {
		return &runtime.ExecutionResult{
			GasLeft: gasLimit,
			Err:     err,
		}
	}

	result := t.run(c, host)

	if result.Failed() {
		t.txn.RevertToSnapshot(snapshot)
		return result
	}

	if t.forks.EIP158 && len(result.ReturnValue) > spuriousDragonMaxCodeSize {
		// Contract size exceeds 'SpuriousDragon' size limit
		t.txn.RevertToSnapshot(snapshot)
		return &runtime.ExecutionResult{
			GasLeft: 0,
			Err:     runtime.ErrMaxCodeSizeExceeded,
		}
	}

	gasCost := uint64(len(result.ReturnValue)) * 200

	if result.GasLeft < gasCost {
		result.Err = runtime.ErrCodeStoreOutOfGas
		result.ReturnValue = nil

		// Out of gas creating the contract
		if t.forks.Homestead {
			t.txn.RevertToSnapshot(snapshot)
			result.GasLeft = 0
		}

		return result
	}

	result.GasLeft -= gasCost
	t.txn.SetCode(c.Address, result.ReturnValue)

	return result
}

func (t *Transition) SetStorage(addr types.Address, key types.Hash, value types.Hash, config *runtime.ForksInTime) runtime.StorageStatus {
	return t.txn.SetStorage(addr, key, value, config)
}

func (t *Transition) GetTxContext() runtime.TxContext {
	return t.ctx
}

func (t *Transition) GetBlockHash(number int64) (res types.Hash) {
	return t.getHash(uint64(number))
}

func (t *Transition) EmitLog(addr types.Address, topics []types.Hash, data []byte) {
	t.txn.EmitLog(addr, topics, data)
}

func (t *Transition) GetCodeSize(addr types.Address) int {
	return t.txn.GetCodeSize(addr)
}

func (t *Transition) GetCodeHash(addr types.Address) (res types.Hash) {
	return t.txn.GetCodeHash(addr)
}

func (t *Transition) GetCode(addr types.Address) []byte {
	return t.txn.GetCode(addr)
}

func (t *Transition) GetBalance(addr types.Address) *big.Int {
	return t.txn.GetBalance(addr)
}

func (t *Transition) GetStorage(addr types.Address, key types.Hash) types.Hash {
	return t.txn.GetState(addr, key)
}

func (t *Transition) AccountExists(addr types.Address) bool {
	return t.txn.Exist(addr)
}

func (t *Transition) Empty(addr types.Address) bool {
	return t.txn.Empty(addr)
}

func (t *Transition) GetNonce(addr types.Address) uint64 {
	return t.txn.GetNonce(addr)
}

func (t *Transition) Selfdestruct(addr types.Address, beneficiary types.Address) {
	if !t.txn.HasSuicided(addr) {
		t.txn.AddRefund(24000)
	}
	t.txn.AddBalance(beneficiary, t.txn.GetBalance(addr))
	t.txn.Suicide(addr)
}

func (t *Transition) Callx(c *runtime.Contract, h runtime.Host) *runtime.ExecutionResult {
	if c.Type == runtime.Create {
		return t.applyCreate(c, h)
	}
	return t.applyCall(c, c.Type, h)
}

func TransactionGasCost(msg *Transaction, isHomestead, isIstanbul bool) (uint64, error) {
	cost := uint64(0)

	// Contract creation is only paid on the homestead fork
	if msg.IsContractCreation() && isHomestead {
		cost += TxGasContractCreation
	} else {
		cost += TxGas
	}

	payload := msg.Input
	if len(payload) > 0 {
		zeros := uint64(0)
		for i := 0; i < len(payload); i++ {
			if payload[i] == 0 {
				zeros++
			}
		}

		nonZeros := uint64(len(payload)) - zeros
		nonZeroCost := uint64(68)
		if isIstanbul {
			nonZeroCost = 16
		}

		if (math.MaxUint64-cost)/nonZeroCost < nonZeros {
			return 0, ErrIntrinsicGasOverflow
		}

		cost += nonZeros * nonZeroCost

		if (math.MaxUint64-cost)/4 < zeros {
			return 0, ErrIntrinsicGasOverflow
		}

		cost += zeros * 4
	}

	return cost, nil
}
