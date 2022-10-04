package runtime

import (
	"encoding/json"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
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
	ErrExecutionReverted = errors.New("execution reverted")
	// ErrCodeStoreOutOfGas indicates there is not enough gas for the storage
	ErrCodeStoreOutOfGas = errors.New("contract creation code storage out of gas")
	// ErrOutOfCounters indicates the executor run out of counters while executing the transaction
	ErrOutOfCounters = errors.New("executor run out of counters")
	// ErrInvalidTransaction indicates the executor found the transaction to be invalid
	ErrInvalidTransaction = errors.New("invalid transaction")
	// ErrIntrinsicInvalidTransaction indicates the executor found the transaction to be invalid and this does not affected the state
	ErrIntrinsicInvalidTransaction = errors.New("intrinsic invalid transaction")
	// ErrBatchDataTooBig indicates the batch_l2_data is too big to be processed
	ErrBatchDataTooBig = errors.New("batch data too big")
)

// ExecutionResult includes all output after executing given evm
// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	ReturnValue         []byte // Returned data from the runtime (function result or data supplied with revert opcode)
	GasLeft             uint64 // Total gas left as result of execution
	GasUsed             uint64 // Total gas used as result of execution
	Err                 error  // Any error encountered during the execution, listed below
	CreateAddress       common.Address
	StateRoot           []byte
	StructLogs          []instrumentation.StructLog
	ExecutorTrace       instrumentation.ExecutorTrace
	ExecutorTraceResult json.RawMessage
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
	return errors.Is(r.Err, ErrExecutionReverted)
}
