package executor

import (
	"fmt"
	"math"

	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
)

var (
	// ErrUnspecified indicates an unspecified executor error
	ErrUnspecified = fmt.Errorf("unspecified executor error")
	// ErrUnknown indicates an unknown executor error
	ErrUnknown = fmt.Errorf("unknown error")
)

// RomErr returns an instance of error related to the ExecutorError
func RomErr(errorCode RomError) error {
	switch errorCode {
	case RomError_ROM_ERROR_UNSPECIFIED:
		return fmt.Errorf("unspecified ROM error")
	case RomError_ROM_ERROR_NO_ERROR:
		return nil
	case RomError_ROM_ERROR_OUT_OF_GAS:
		return runtime.ErrOutOfGas
	case RomError_ROM_ERROR_STACK_OVERFLOW:
		return runtime.ErrStackOverflow
	case RomError_ROM_ERROR_STACK_UNDERFLOW:
		return runtime.ErrStackUnderflow
	case RomError_ROM_ERROR_MAX_CODE_SIZE_EXCEEDED:
		return runtime.ErrMaxCodeSizeExceeded
	case RomError_ROM_ERROR_CONTRACT_ADDRESS_COLLISION:
		return runtime.ErrContractAddressCollision
	case RomError_ROM_ERROR_EXECUTION_REVERTED:
		return runtime.ErrExecutionReverted
	case RomError_ROM_ERROR_OUT_OF_COUNTERS_STEP:
		return runtime.ErrOutOfCountersStep
	case RomError_ROM_ERROR_OUT_OF_COUNTERS_KECCAK:
		return runtime.ErrOutOfCountersKeccak
	case RomError_ROM_ERROR_OUT_OF_COUNTERS_BINARY:
		return runtime.ErrOutOfCountersBinary
	case RomError_ROM_ERROR_OUT_OF_COUNTERS_MEM:
		return runtime.ErrOutOfCountersMemory
	case RomError_ROM_ERROR_OUT_OF_COUNTERS_ARITH:
		return runtime.ErrOutOfCountersArith
	case RomError_ROM_ERROR_OUT_OF_COUNTERS_PADDING:
		return runtime.ErrOutOfCountersPadding
	case RomError_ROM_ERROR_OUT_OF_COUNTERS_POSEIDON:
		return runtime.ErrOutOfCountersPoseidon
	case RomError_ROM_ERROR_INVALID_JUMP:
		return runtime.ErrInvalidJump
	case RomError_ROM_ERROR_INVALID_OPCODE:
		return runtime.ErrInvalidOpCode
	case RomError_ROM_ERROR_INVALID_STATIC:
		return runtime.ErrInvalidStatic
	case RomError_ROM_ERROR_INVALID_BYTECODE_STARTS_EF:
		return runtime.ErrInvalidByteCodeStartsEF
	case RomError_ROM_ERROR_INTRINSIC_INVALID_SIGNATURE:
		return runtime.ErrIntrinsicInvalidSignature
	case RomError_ROM_ERROR_INTRINSIC_INVALID_CHAIN_ID:
		return runtime.ErrIntrinsicInvalidChainID
	case RomError_ROM_ERROR_INTRINSIC_INVALID_NONCE:
		return runtime.ErrIntrinsicInvalidNonce
	case RomError_ROM_ERROR_INTRINSIC_INVALID_GAS_LIMIT:
		return runtime.ErrIntrinsicInvalidGasLimit
	case RomError_ROM_ERROR_INTRINSIC_INVALID_BALANCE:
		return runtime.ErrIntrinsicInvalidBalance
	case RomError_ROM_ERROR_INTRINSIC_INVALID_BATCH_GAS_LIMIT:
		return runtime.ErrIntrinsicInvalidBatchGasLimit
	case RomError_ROM_ERROR_INTRINSIC_INVALID_SENDER_CODE:
		return runtime.ErrIntrinsicInvalidSenderCode
	case RomError_ROM_ERROR_INTRINSIC_TX_GAS_OVERFLOW:
		return runtime.ErrIntrinsicInvalidTxGasOverflow
	case RomError_ROM_ERROR_BATCH_DATA_TOO_BIG:
		return runtime.ErrBatchDataTooBig
	case RomError_ROM_ERROR_UNSUPPORTED_FORK_ID:
		return runtime.ErrUnsupportedForkId
	case RomError_ROM_ERROR_INVALID_RLP:
		return runtime.ErrInvalidRLP
	}
	return fmt.Errorf("unknown error")
}

// RomErrorCode returns the error code for a given error
func RomErrorCode(err error) RomError {
	switch err {
	case nil:
		return RomError_ROM_ERROR_NO_ERROR
	case runtime.ErrOutOfGas:
		return RomError_ROM_ERROR_OUT_OF_GAS
	case runtime.ErrStackOverflow:
		return RomError_ROM_ERROR_STACK_OVERFLOW
	case runtime.ErrStackUnderflow:
		return RomError_ROM_ERROR_STACK_UNDERFLOW
	case runtime.ErrMaxCodeSizeExceeded:
		return RomError_ROM_ERROR_MAX_CODE_SIZE_EXCEEDED
	case runtime.ErrContractAddressCollision:
		return RomError_ROM_ERROR_CONTRACT_ADDRESS_COLLISION
	case runtime.ErrExecutionReverted:
		return RomError_ROM_ERROR_EXECUTION_REVERTED
	case runtime.ErrOutOfCountersStep:
		return RomError_ROM_ERROR_OUT_OF_COUNTERS_STEP
	case runtime.ErrOutOfCountersKeccak:
		return RomError_ROM_ERROR_OUT_OF_COUNTERS_KECCAK
	case runtime.ErrOutOfCountersBinary:
		return RomError_ROM_ERROR_OUT_OF_COUNTERS_BINARY
	case runtime.ErrOutOfCountersMemory:
		return RomError_ROM_ERROR_OUT_OF_COUNTERS_MEM
	case runtime.ErrOutOfCountersArith:
		return RomError_ROM_ERROR_OUT_OF_COUNTERS_ARITH
	case runtime.ErrOutOfCountersPadding:
		return RomError_ROM_ERROR_OUT_OF_COUNTERS_PADDING
	case runtime.ErrOutOfCountersPoseidon:
		return RomError_ROM_ERROR_OUT_OF_COUNTERS_POSEIDON
	case runtime.ErrInvalidJump:
		return RomError_ROM_ERROR_INVALID_JUMP
	case runtime.ErrInvalidOpCode:
		return RomError_ROM_ERROR_INVALID_OPCODE
	case runtime.ErrInvalidStatic:
		return RomError_ROM_ERROR_INVALID_STATIC
	case runtime.ErrInvalidByteCodeStartsEF:
		return RomError_ROM_ERROR_INVALID_BYTECODE_STARTS_EF
	case runtime.ErrIntrinsicInvalidSignature:
		return RomError_ROM_ERROR_INTRINSIC_INVALID_SIGNATURE
	case runtime.ErrIntrinsicInvalidChainID:
		return RomError_ROM_ERROR_INTRINSIC_INVALID_CHAIN_ID
	case runtime.ErrIntrinsicInvalidNonce:
		return RomError_ROM_ERROR_INTRINSIC_INVALID_NONCE
	case runtime.ErrIntrinsicInvalidGasLimit:
		return RomError_ROM_ERROR_INTRINSIC_INVALID_GAS_LIMIT
	case runtime.ErrIntrinsicInvalidBalance:
		return RomError_ROM_ERROR_INTRINSIC_INVALID_BALANCE
	case runtime.ErrIntrinsicInvalidBatchGasLimit:
		return RomError_ROM_ERROR_INTRINSIC_INVALID_BATCH_GAS_LIMIT
	case runtime.ErrIntrinsicInvalidSenderCode:
		return RomError_ROM_ERROR_INTRINSIC_INVALID_SENDER_CODE
	case runtime.ErrIntrinsicInvalidTxGasOverflow:
		return RomError_ROM_ERROR_INTRINSIC_TX_GAS_OVERFLOW
	case runtime.ErrBatchDataTooBig:
		return RomError_ROM_ERROR_BATCH_DATA_TOO_BIG
	case runtime.ErrUnsupportedForkId:
		return RomError_ROM_ERROR_UNSUPPORTED_FORK_ID
	case runtime.ErrInvalidRLP:
		return RomError_ROM_ERROR_INVALID_RLP
	}
	return math.MaxInt32
}

// IsROMOutOfCountersError indicates if the error is an ROM OOC
func IsROMOutOfCountersError(error RomError) bool {
	return error >= RomError_ROM_ERROR_OUT_OF_COUNTERS_STEP && error <= RomError_ROM_ERROR_OUT_OF_COUNTERS_POSEIDON
}

// IsROMOutOfGasError indicates if the error is an ROM OOG
func IsROMOutOfGasError(error RomError) bool {
	return error == RomError_ROM_ERROR_OUT_OF_GAS
}

// IsExecutorOutOfCountersError indicates if the error is an ROM OOC
func IsExecutorOutOfCountersError(error ExecutorError) bool {
	return error >= ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_KECCAK && error <= ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_POSEIDON
}

// IsExecutorUnspecifiedError indicates an unspecified error in the executor
func IsExecutorUnspecifiedError(error ExecutorError) bool {
	return error == ExecutorError_EXECUTOR_ERROR_UNSPECIFIED
}

// IsIntrinsicError indicates if the error is due to a intrinsic check
func IsIntrinsicError(error RomError) bool {
	return error >= RomError_ROM_ERROR_INTRINSIC_INVALID_SIGNATURE && error <= RomError_ROM_ERROR_INTRINSIC_TX_GAS_OVERFLOW
}

// IsInvalidNonceError indicates if the error is due to a invalid nonce
func IsInvalidNonceError(error RomError) bool {
	return error == RomError_ROM_ERROR_INTRINSIC_INVALID_NONCE
}

// IsInvalidBalanceError indicates if the error is due to a invalid balance
func IsInvalidBalanceError(error RomError) bool {
	return error == RomError_ROM_ERROR_INTRINSIC_INVALID_BALANCE
}

// ExecutorErr returns an instance of error related to the ExecutorError
func ExecutorErr(errorCode ExecutorError) error {
	switch errorCode {
	case ExecutorError_EXECUTOR_ERROR_UNSPECIFIED:
		return ErrUnspecified
	case ExecutorError_EXECUTOR_ERROR_NO_ERROR:
		return nil
	case ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_KECCAK:
		return runtime.ErrOutOfCountersKeccak
	case ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_BINARY:
		return runtime.ErrOutOfCountersBinary
	case ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_MEM:
		return runtime.ErrOutOfCountersMemory
	case ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_ARITH:
		return runtime.ErrOutOfCountersArith
	case ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_PADDING:
		return runtime.ErrOutOfCountersPadding
	case ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_POSEIDON:
		return runtime.ErrOutOfCountersPoseidon
	case ExecutorError_EXECUTOR_ERROR_UNSUPPORTED_FORK_ID:
		return runtime.ErrUnsupportedForkId
	case ExecutorError_EXECUTOR_ERROR_BALANCE_MISMATCH:
		return runtime.ErrBalanceMismatch
	case ExecutorError_EXECUTOR_ERROR_FEA2SCALAR:
		return runtime.ErrFea2Scalar
	case ExecutorError_EXECUTOR_ERROR_TOS32:
		return runtime.ErrTos32
	}
	return ErrUnknown
}

// ExecutorErrorCode returns the error code for a given error
func ExecutorErrorCode(err error) ExecutorError {
	switch err {
	case nil:
		return ExecutorError_EXECUTOR_ERROR_NO_ERROR
	case runtime.ErrOutOfCountersKeccak:
		return ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_KECCAK
	case runtime.ErrOutOfCountersBinary:
		return ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_BINARY
	case runtime.ErrOutOfCountersMemory:
		return ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_MEM
	case runtime.ErrOutOfCountersArith:
		return ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_ARITH
	case runtime.ErrOutOfCountersPadding:
		return ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_PADDING
	case runtime.ErrOutOfCountersPoseidon:
		return ExecutorError_EXECUTOR_ERROR_COUNTERS_OVERFLOW_POSEIDON
	case runtime.ErrUnsupportedForkId:
		return ExecutorError_EXECUTOR_ERROR_UNSUPPORTED_FORK_ID
	case runtime.ErrBalanceMismatch:
		return ExecutorError_EXECUTOR_ERROR_BALANCE_MISMATCH
	case runtime.ErrFea2Scalar:
		return ExecutorError_EXECUTOR_ERROR_FEA2SCALAR
	case runtime.ErrTos32:
		return ExecutorError_EXECUTOR_ERROR_TOS32
	}
	return math.MaxInt32
}
