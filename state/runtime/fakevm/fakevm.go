package fakevm

import (
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

type FakeEVM struct {
	// Context provides auxiliary blockchain related information
	Context vm.BlockContext
	vm.TxContext
	// StateDB gives access to the underlying state
	StateDB FakeDB
	// Depth is the current call stack
	// depth int

	// chainConfig contains information about the current chain
	chainConfig *params.ChainConfig
	// chain rules contains the chain rules for the current epoch
	chainRules params.Rules
	// virtual machine configuration options used to initialise the
	// evm.
	Config Config
	// global (to this context) ethereum virtual machine
	// used throughout the execution of the tx.
	// interpreter *vm.EVMInterpreter
	// abort is used to abort the EVM calling operations
	// NOTE: must be set atomically
	abort int32
	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	// callGasTemp uint64
}

// NewEVM returns a new EVM. The returned EVM is not thread safe and should
// only ever be used *once*.
// func NewFakeEVM(blockCtx vm.BlockContext, txCtx vm.TxContext, statedb runtime.FakeDB, chainConfig *params.ChainConfig, config Config) *FakeEVM {
func NewFakeEVM(blockCtx vm.BlockContext, txCtx vm.TxContext, chainConfig *params.ChainConfig, config Config) *FakeEVM {
	evm := &FakeEVM{
		Context:   blockCtx,
		TxContext: txCtx,
		// StateDB:     statedb,
		Config:      config,
		chainConfig: chainConfig,
		chainRules:  chainConfig.Rules(blockCtx.BlockNumber, blockCtx.Random != nil),
	}
	return evm
}

func (evm *FakeEVM) SetStateDB(stateDB FakeDB) {
	evm.StateDB = stateDB
}

// Cancel cancels any running EVM operation. This may be called concurrently and
// it's safe to be called multiple times.
func (evm *FakeEVM) Cancel() {
	atomic.StoreInt32(&evm.abort, 1)
}

// ChainConfig returns the environment's chain configuration
func (evm *FakeEVM) ChainConfig() *params.ChainConfig { return evm.chainConfig }

// ScopeContext contains the things that are per-call, such as stack and memory,
// but not transients like pc and gas
type ScopeContext struct {
	Memory   *Memory
	Stack    *Stack
	Contract *vm.Contract
}
