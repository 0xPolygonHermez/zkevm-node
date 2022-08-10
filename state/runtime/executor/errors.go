package executor

import (
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
)

// ExecutorError is an error returned by the Executor
type ExecutorError int32

const (
	// ERROR_UNSPECIFIED indicates the execution ended successfully
	ERROR_UNSPECIFIED ExecutorError = iota
	// ERROR_NO_ERROR indicates the execution ended successfully
	ERROR_NO_ERROR
	// ERROR_OUT_OF_GAS indicates there is not enough balance to continue the execution
	ERROR_OUT_OF_GAS
	// ERROR_STACK_OVERFLOW indicates a stack overflow has happened
	ERROR_STACK_OVERFLOW
	// ERROR_STACK_UNDERFLOW indicates a stack overflow has happened
	ERROR_STACK_UNDERFLOW
	// ERROR_NOT_ENOUGH_FUNDS indicates there is not enough funds to continue the execution
	ERROR_NOT_ENOUGH_FUNDS
	// ERROR_INSUFFICIENT_BALANCE indicates there is not enough balance to continue the execution
	ERROR_INSUFFICIENT_BALANCE
	// ERROR_CODE_NOT_FOUND indicates the code was not found
	ERROR_CODE_NOT_FOUND
	// ERROR_MAX_CODE_SIZE_EXCEEDED indicates the code size is beyond the maximum
	ERROR_MAX_CODE_SIZE_EXCEEDED
	// ERROR_CONTRACT_ADDRESS_COLLISION there is a collision regarding contract addresses
	ERROR_CONTRACT_ADDRESS_COLLISION
	// ERROR_DEPTH indicates the maximum call depth has been passed
	ERROR_DEPTH
	// ERROR_EXECUTION_REVERTED indicates the execution has been reverted
	ERROR_EXECUTION_REVERTED
	// ERROR_CODE_STORE_OUT_OF_GAS indicates there is not enough gas for the storage
	ERROR_CODE_STORE_OUT_OF_GAS
	// ERROR_OUT_OF_COUNTERS indicates there is not enough counters to continue the execution
	ERROR_OUT_OF_COUNTERS
	// ERROR_INVALID_TX indicates the transaction is invalid
	ERROR_INVALID_TX
	// ERROR_INTRINSIC_INVALID_TX indicates the transaction is failing at the intrinsic checks
	ERROR_INTRINSIC_INVALID_TX
)

func (e ExecutorError) Error() string {
	switch e {
	case ERROR_NO_ERROR:
		return ""
	case ERROR_OUT_OF_GAS:
		return runtime.ErrOutOfGas.Error()
	case ERROR_STACK_OVERFLOW:
		return runtime.ErrStackOverflow.Error()
	case ERROR_STACK_UNDERFLOW:
		return runtime.ErrStackUnderflow.Error()
	case ERROR_NOT_ENOUGH_FUNDS:
		return runtime.ErrNotEnoughFunds.Error()
	case ERROR_INSUFFICIENT_BALANCE:
		return runtime.ErrInsufficientBalance.Error()
	case ERROR_CODE_NOT_FOUND:
		return runtime.ErrCodeNotFound.Error()
	case ERROR_MAX_CODE_SIZE_EXCEEDED:
		return runtime.ErrMaxCodeSizeExceeded.Error()
	case ERROR_CONTRACT_ADDRESS_COLLISION:
		return runtime.ErrContractAddressCollision.Error()
	case ERROR_DEPTH:
		return runtime.ErrDepth.Error()
	case ERROR_EXECUTION_REVERTED:
		return runtime.ErrExecutionReverted.Error()
	case ERROR_CODE_STORE_OUT_OF_GAS:
		return runtime.ErrCodeStoreOutOfGas.Error()
	case ERROR_OUT_OF_COUNTERS:
		return runtime.ErrOutOfCounters.Error()
	case ERROR_INVALID_TX, ERROR_INTRINSIC_INVALID_TX:
		return runtime.ErrInvalidTransaction.Error()
	}
	return "unknown error"
}
