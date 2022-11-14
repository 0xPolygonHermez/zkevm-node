package executor

import (
	"fmt"
	"math"

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
	// ERROR_OUT_OF_COUNTERS_STEP indicates there is not enough step counters to continue the execution
	ERROR_OUT_OF_COUNTERS_STEP
	// ERROR_OUT_OF_COUNTERS_KECCAK indicates there is not enough keccak counters to continue the execution
	ERROR_OUT_OF_COUNTERS_KECCAK
	// ERROR_OUT_OF_COUNTERS_BINARY indicates there is not enough binary counters to continue the execution
	ERROR_OUT_OF_COUNTERS_BINARY
	// ERROR_OUT_OF_COUNTERS_MEM indicates there is not enough memory aligncounters to continue the execution
	ERROR_OUT_OF_COUNTERS_MEM
	// ERROR_OUT_OF_COUNTERS_ARITH indicates there is not enough arith counters to continue the execution
	ERROR_OUT_OF_COUNTERS_ARITH
	// ERROR_OUT_OF_COUNTERS_PADDING indicates there is not enough padding counters to continue the execution
	ERROR_OUT_OF_COUNTERS_PADDING
	// ERROR_OUT_OF_COUNTERS_POSEIDON indicates there is not enough poseidon counters to continue the execution
	ERROR_OUT_OF_COUNTERS_POSEIDON
	// ERROR_INVALID_TX indicates the transaction is invalid because of invalid jump dest, invalid opcode, invalid deploy
	// or invalid static tx
	ERROR_INVALID_TX
	// ERROR_INTRINSIC_INVALID_SIGNATURE indicates the transaction is failing at the signature intrinsic check
	ERROR_INTRINSIC_INVALID_SIGNATURE
	// ERROR_INTRINSIC_INVALID_CHAIN_ID indicates the transaction is failing at the chain id intrinsic check
	ERROR_INTRINSIC_INVALID_CHAIN_ID
	// ERROR_INTRINSIC_INVALID_NONCE indicates the transaction is failing at the nonce intrinsic check
	ERROR_INTRINSIC_INVALID_NONCE
	// ERROR_INTRINSIC_INVALID_GAS_LIMIT indicates the transaction is failing at the gas limit intrinsic check
	ERROR_INTRINSIC_INVALID_GAS_LIMIT
	// ERROR_INTRINSIC_INVALID_BALANCE indicates the transaction is failing at balance intrinsic check
	ERROR_INTRINSIC_INVALID_BALANCE
	// ERROR_INTRINSIC_INVALID_BATCH_GAS_LIMIT indicates the batch is exceeding the batch gas limit
	ERROR_INTRINSIC_INVALID_BATCH_GAS_LIMIT
	// ERROR_INTRINSIC_INVALID_SENDER_CODE indicates the batch is exceeding the batch gas limit
	ERROR_INTRINSIC_INVALID_SENDER_CODE
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
	case ERROR_OUT_OF_COUNTERS_STEP:
		return runtime.ErrOutOfCountersKeccak
	case ERROR_OUT_OF_COUNTERS_BINARY:
		return runtime.ErrOutOfCountersBinary
	case ERROR_OUT_OF_COUNTERS_MEM:
		return runtime.ErrOutOfCountersMemory
	case ERROR_OUT_OF_COUNTERS_ARITH:
		return runtime.ErrOutOfCountersArith
	case ERROR_OUT_OF_COUNTERS_PADDING:
		return runtime.ErrOutOfCountersPadding
	case ERROR_OUT_OF_COUNTERS_POSEIDON:
		return runtime.ErrOutOfCountersPoseidon
	case ERROR_INVALID_TX:
		return runtime.ErrInvalidTransaction
	case ERROR_INTRINSIC_INVALID_SIGNATURE:
		return runtime.ErrIntrinsicInvalidSignature
	case ERROR_INTRINSIC_INVALID_CHAIN_ID:
		return runtime.ErrIntrinsicInvalidChainID
	case ERROR_INTRINSIC_INVALID_NONCE:
		return runtime.ErrIntrinsicInvalidNonce
	case ERROR_INTRINSIC_INVALID_GAS_LIMIT:
		return runtime.ErrIntrinsicInvalidGasLimit
	case ERROR_INTRINSIC_INVALID_BALANCE:
		return runtime.ErrIntrinsicInvalidBalance
	case ERROR_INTRINSIC_INVALID_BATCH_GAS_LIMIT:
		return runtime.ErrIntrinsicInvalidGasLimit
	case ERROR_INTRINSIC_INVALID_SENDER_CODE:
		return runtime.ErrIntrinsicInvalidSenderCode
	case ERROR_BATCH_DATA_TOO_BIG:
		return runtime.ErrBatchDataTooBig
	}
	return fmt.Errorf("unknown error")
}

// ErrorCode returns the error code for a given error
func ErrorCode(err error) pb.Error {
	switch err {
	case nil:
		return pb.Error(ERROR_NO_ERROR)
	case runtime.ErrOutOfGas:
		return pb.Error(ERROR_OUT_OF_GAS)
	case runtime.ErrStackOverflow:
		return pb.Error(ERROR_STACK_OVERFLOW)
	case runtime.ErrStackUnderflow:
		return pb.Error(ERROR_STACK_UNDERFLOW)
	case runtime.ErrNotEnoughFunds:
		return pb.Error(ERROR_NOT_ENOUGH_FUNDS)
	case runtime.ErrInsufficientBalance:
		return pb.Error(ERROR_INSUFFICIENT_BALANCE)
	case runtime.ErrCodeNotFound:
		return pb.Error(ERROR_CODE_NOT_FOUND)
	case runtime.ErrMaxCodeSizeExceeded:
		return pb.Error(ERROR_MAX_CODE_SIZE_EXCEEDED)
	case runtime.ErrContractAddressCollision:
		return pb.Error(ERROR_CONTRACT_ADDRESS_COLLISION)
	case runtime.ErrDepth:
		return pb.Error(ERROR_DEPTH)
	case runtime.ErrExecutionReverted:
		return pb.Error(ERROR_EXECUTION_REVERTED)
	case runtime.ErrCodeStoreOutOfGas:
		return pb.Error(ERROR_CODE_STORE_OUT_OF_GAS)
	case runtime.ErrOutOfCountersKeccak:
		return pb.Error(ERROR_OUT_OF_COUNTERS_STEP)
	case runtime.ErrOutOfCountersBinary:
		return pb.Error(ERROR_OUT_OF_COUNTERS_BINARY)
	case runtime.ErrOutOfCountersMemory:
		return pb.Error(ERROR_OUT_OF_COUNTERS_MEM)
	case runtime.ErrOutOfCountersArith:
		return pb.Error(ERROR_OUT_OF_COUNTERS_ARITH)
	case runtime.ErrOutOfCountersPadding:
		return pb.Error(ERROR_OUT_OF_COUNTERS_PADDING)
	case runtime.ErrOutOfCountersPoseidon:
		return pb.Error(ERROR_OUT_OF_COUNTERS_POSEIDON)
	case runtime.ErrInvalidTransaction:
		return pb.Error(ERROR_INVALID_TX)
	case runtime.ErrIntrinsicInvalidSignature:
		return pb.Error(ERROR_INTRINSIC_INVALID_SIGNATURE)
	case runtime.ErrIntrinsicInvalidChainID:
		return pb.Error(ERROR_INTRINSIC_INVALID_CHAIN_ID)
	case runtime.ErrIntrinsicInvalidNonce:
		return pb.Error(ERROR_INTRINSIC_INVALID_NONCE)
	case runtime.ErrIntrinsicInvalidGasLimit:
		return pb.Error(ERROR_INTRINSIC_INVALID_GAS_LIMIT)
	case runtime.ErrIntrinsicInvalidBalance:
		return pb.Error(ERROR_INTRINSIC_INVALID_BALANCE)
	case runtime.ErrIntrinsicInvalidGasLimit:
		return pb.Error(ERROR_INTRINSIC_INVALID_BATCH_GAS_LIMIT)
	case runtime.ErrIntrinsicInvalidSenderCode:
		return pb.Error(ERROR_INTRINSIC_INVALID_SENDER_CODE)
	case runtime.ErrBatchDataTooBig:
		return pb.Error(ERROR_BATCH_DATA_TOO_BIG)
	}
	return math.MaxInt32
}

// IsOutOfCountersError indicates if the error is an OOC
func IsOutOfCountersError(error pb.Error) bool {
	return int32(error) >= ERROR_OUT_OF_COUNTERS_STEP && int32(error) <= ERROR_OUT_OF_COUNTERS_POSEIDON
}

// IsIntrinsicError indicates if the error is due to a intrinsic check
func IsIntrinsicError(error pb.Error) bool {
	return int32(error) >= ERROR_INVALID_TX && int32(error) <= ERROR_INTRINSIC_INVALID_SENDER_CODE
}
