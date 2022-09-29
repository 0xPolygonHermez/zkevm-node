package executor

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
)

const (
	// ERROR_UNSPECIFIED indicates the execution ended successfully
	ERROR_UNSPECIFIED int32 = iota
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
	// ERROR_BATCH_DATA_TOO_BIG indicates the batch_l2_data is too big to be processed
	ERROR_BATCH_DATA_TOO_BIG
)

// Err returns an instance of error related to the ExecutorError
func Err(errorCode pb.Error) error {
	e := int32(errorCode)
	switch e {
	case ERROR_NO_ERROR, ERROR_UNSPECIFIED:
		return nil
	case ERROR_OUT_OF_GAS:
		return runtime.ErrOutOfGas
	case ERROR_STACK_OVERFLOW:
		return runtime.ErrStackOverflow
	case ERROR_STACK_UNDERFLOW:
		return runtime.ErrStackUnderflow
	case ERROR_NOT_ENOUGH_FUNDS:
		return runtime.ErrNotEnoughFunds
	case ERROR_INSUFFICIENT_BALANCE:
		return runtime.ErrInsufficientBalance
	case ERROR_CODE_NOT_FOUND:
		return runtime.ErrCodeNotFound
	case ERROR_MAX_CODE_SIZE_EXCEEDED:
		return runtime.ErrMaxCodeSizeExceeded
	case ERROR_CONTRACT_ADDRESS_COLLISION:
		return runtime.ErrContractAddressCollision
	case ERROR_DEPTH:
		return runtime.ErrDepth
	case ERROR_EXECUTION_REVERTED:
		return runtime.ErrExecutionReverted
	case ERROR_CODE_STORE_OUT_OF_GAS:
		return runtime.ErrCodeStoreOutOfGas
	case ERROR_OUT_OF_COUNTERS:
		return runtime.ErrOutOfCounters
	case ERROR_INVALID_TX:
		return runtime.ErrInvalidTransaction
	case ERROR_INTRINSIC_INVALID_TX:
		return runtime.ErrIntrinsicInvalidTransaction
	case ERROR_BATCH_DATA_TOO_BIG:
		return runtime.ErrBatchDataTooBig
	}
	return fmt.Errorf("unknown error")
}
