package runtime

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation"
)

var (
	// ErrOutOfGas indicates there is not enough balance to continue the execution
	ErrOutOfGas = errors.New("out of gas")
	// ErrStackOverflow indicates a stack overflow has happened
	ErrStackOverflow = errors.New("stack overflow")
	// ErrStackUnderflow indicates a stack overflow has happened
	ErrStackUnderflow = errors.New("stack underflow")
	// ErrNotEnoughFunds indicates there is not enough funds to continue the execution
	ErrNotEnoughFunds = errors.New("not enough funds")
	// ErrInsufficientBalance indicates there is not enough balance to continue the execution
	ErrInsufficientBalance = errors.New("insufficient balance for transfer")
	// ErrCodeNotFound indicates the code was not found
	ErrCodeNotFound = errors.New("code not found, data is empty")
	// ErrMaxCodeSizeExceeded indicates the code size is beyond the maximum
	ErrMaxCodeSizeExceeded = errors.New("evm: max code size exceeded")
	// ErrContractAddressCollision there is a collision regarding contract addresses
	ErrContractAddressCollision = errors.New("contract address collision")
	// ErrDepth indicates the maximun call depth has been passed
	ErrDepth = errors.New("max call depth exceeded")
	// ErrExecutionReverted indicates the execution has been reverted
	ErrExecutionReverted = errors.New("execution was reverted")
	// ErrCodeStoreOutOfGas indicates there is not enough gas for the storage
	ErrCodeStoreOutOfGas = errors.New("contract creation code storage out of gas")
)

// CallType indicates the type of call to a contract
type CallType int

const (
	// Call is the default call for a contract
	Call CallType = iota
	// CallCode is the callcode call for a contract
	CallCode
	// DelegateCall is the delegate call for a contract
	DelegateCall
	// StaticCall is the static call for a contract
	StaticCall
	// Create is the creation call for a contract
	Create
	// Create2 is the creation call for a contract from a contract
	Create2
)

// Runtime can process contracts
type Runtime interface {
	Run(ctx context.Context, c *Contract, host Host, config *ForksInTime) *ExecutionResult
	CanRun(c *Contract, host Host, config *ForksInTime) bool
	Name() string
}

// TxContext is the context of the transaction
type TxContext struct {
	Hash        common.Hash
	GasPrice    common.Hash
	Origin      common.Address
	Coinbase    common.Address
	Number      int64
	Timestamp   int64
	GasLimit    int64
	ChainID     int64
	Difficulty  common.Hash
	BatchNumber int64
}

// ExecutionResult includes all output after executing given evm
// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	ReturnValue   []byte // Returned data from the runtime (function result or data supplied with revert opcode)
	GasLeft       uint64 // Total gas left as result of execution
	GasUsed       uint64 // Total gas used as result of execution
	Err           error  // Any error encountered during the execution, listed below
	CreateAddress common.Address
	StateRoot     []byte
	Trace         []instrumentation.Trace
	VMTrace       instrumentation.VMTrace
	StructLogs    []instrumentation.StructLog
	ExecutorTrace instrumentation.ExecutorTrace
}

// Succeeded indicates the execution was successful
func (r *ExecutionResult) Succeeded() bool {
	return r.Err == nil
}

// Failed indicates the execution was unsuccessful
func (r *ExecutionResult) Failed() bool {
	return r.Err != nil
}

// Reverted indicates the execution was reverted
func (r *ExecutionResult) Reverted() bool {
	return r.Err == ErrExecutionReverted
}
