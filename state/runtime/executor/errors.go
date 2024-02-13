package executor

import (
	"fmt"
	"math"

	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
)

var (
	// ErrExecutorUnspecified indicates an unspecified executor error
	ErrExecutorUnspecified = fmt.Errorf("unspecified executor error")
	// ErrROMUnspecified indicates an unspecified ROM error
	ErrROMUnspecified = fmt.Errorf("unspecified ROM error")
	// ErrExecutorUnknown indicates an unknown executor error
	ErrExecutorUnknown = fmt.Errorf("unknown executor error")
	// ErrCodeExecutorUnknown indicates an unknown executor error
	ErrCodeExecutorUnknown = ExecutorError(math.MaxInt32)
	// ErrROMUnknown indicates an unknown ROM error
	ErrROMUnknown = fmt.Errorf("unknown ROM error")
	// ErrCodeROMUnknown indicates an unknown ROM error
	ErrCodeROMUnknown = RomError(math.MaxInt32)
)

// RomErr returns an instance of error related to the ExecutorError
func RomErr(errorCode RomError) error {
	switch errorCode {
	case RomError_ROM_ERROR_UNSPECIFIED:
		return ErrROMUnspecified
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
	case RomError_ROM_ERROR_OUT_OF_COUNTERS_SHA:
		return runtime.ErrOutOfCountersSha
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
	// Start of V2 errors
	case RomError_ROM_ERROR_INVALID_DECODE_CHANGE_L2_BLOCK:
		return runtime.ErrInvalidDecodeChangeL2Block
	case RomError_ROM_ERROR_INVALID_NOT_FIRST_TX_CHANGE_L2_BLOCK:
		return runtime.ErrInvalidNotFirstTxChangeL2Block
	case RomError_ROM_ERROR_INVALID_TX_CHANGE_L2_BLOCK_LIMIT_TIMESTAMP:
		return runtime.ErrInvalidTxChangeL2BlockLimitTimestamp
	case RomError_ROM_ERROR_INVALID_TX_CHANGE_L2_BLOCK_MIN_TIMESTAMP:
		return runtime.ErrInvalidTxChangeL2BlockMinTimestamp
	}
	return ErrROMUnknown
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
	case runtime.ErrOutOfCountersSha:
		return RomError_ROM_ERROR_OUT_OF_COUNTERS_SHA
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
	// Start of V2 errors
	case runtime.ErrInvalidDecodeChangeL2Block:
		return RomError_ROM_ERROR_INVALID_DECODE_CHANGE_L2_BLOCK
	case runtime.ErrInvalidNotFirstTxChangeL2Block:
		return RomError_ROM_ERROR_INVALID_NOT_FIRST_TX_CHANGE_L2_BLOCK
	case runtime.ErrInvalidTxChangeL2BlockLimitTimestamp:
		return RomError_ROM_ERROR_INVALID_TX_CHANGE_L2_BLOCK_LIMIT_TIMESTAMP
	case runtime.ErrInvalidTxChangeL2BlockMinTimestamp:
		return RomError_ROM_ERROR_INVALID_TX_CHANGE_L2_BLOCK_MIN_TIMESTAMP
	}
	return ErrCodeROMUnknown
}

// IsROMOutOfCountersError indicates if the error is an ROM OOC
func IsROMOutOfCountersError(error RomError) bool {
	return error >= RomError_ROM_ERROR_OUT_OF_COUNTERS_STEP && error <= RomError_ROM_ERROR_OUT_OF_COUNTERS_SHA
}

// IsROMOutOfGasError indicates if the error is an ROM OOG
func IsROMOutOfGasError(error RomError) bool {
	return error == RomError_ROM_ERROR_OUT_OF_GAS
}

// IsExecutorCountersOverflowError indicates if the error is OOC detected by the executor
func IsExecutorCountersOverflowError(error ExecutorError) bool {
	return (error >= ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_STEPS && error <= ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_POSEIDON) || error == ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_SHA256
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

// IsInvalidL2Block indicates if the error is related to L2Block and invalidate all the batch
func IsInvalidL2Block(error RomError) bool {
	return error >= RomError_ROM_ERROR_INVALID_DECODE_CHANGE_L2_BLOCK && error <= RomError_ROM_ERROR_INVALID_TX_CHANGE_L2_BLOCK_MIN_TIMESTAMP
}

// ExecutorErr returns an instance of error related to the ExecutorError
func ExecutorErr(errorCode ExecutorError) error {
	switch errorCode {
	case ExecutorError_EXECUTOR_ERROR_UNSPECIFIED:
		return ErrExecutorUnspecified
	case ExecutorError_EXECUTOR_ERROR_NO_ERROR:
		return nil
	case ExecutorError_EXECUTOR_ERROR_DB_ERROR:
		return runtime.ErrExecutorDBError
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_STEPS:
		return runtime.ErrExecutorSMMainCountersOverflowSteps
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_KECCAK:
		return runtime.ErrExecutorSMMainCountersOverflowKeccak
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_BINARY:
		return runtime.ErrExecutorSMMainCountersOverflowBinary
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_MEM:
		return runtime.ErrExecutorSMMainCountersOverflowMem
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_ARITH:
		return runtime.ErrExecutorSMMainCountersOverflowArith
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_PADDING:
		return runtime.ErrExecutorSMMainCountersOverflowPadding
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_POSEIDON:
		return runtime.ErrExecutorSMMainCountersOverflowPoseidon
	case ExecutorError_EXECUTOR_ERROR_UNSUPPORTED_FORK_ID:
		return runtime.ErrExecutorUnsupportedForkId
	case ExecutorError_EXECUTOR_ERROR_BALANCE_MISMATCH:
		return runtime.ErrExecutorBalanceMismatch
	case ExecutorError_EXECUTOR_ERROR_FEA2SCALAR:
		return runtime.ErrExecutorFEA2Scalar
	case ExecutorError_EXECUTOR_ERROR_TOS32:
		return runtime.ErrExecutorTOS32
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_INVALID_UNSIGNED_TX:
		return runtime.ErrExecutorSMMainInvalidUnsignedTx
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_INVALID_NO_COUNTERS:
		return runtime.ErrExecutorSMMainInvalidNoCounters
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_ARITH_ECRECOVER_DIVIDE_BY_ZERO:
		return runtime.ErrExecutorSMMainArithECRecoverDivideByZero
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_ADDRESS_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainAddressOutOfRange
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_ADDRESS_NEGATIVE:
		return runtime.ErrExecutorSMMainAddressNegative
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_STORAGE_INVALID_KEY:
		return runtime.ErrExecutorSMMainStorageInvalidKey
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK:
		return runtime.ErrExecutorSMMainHashK
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_SIZE_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainHashKSizeOutOfRange
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_POSITION_NEGATIVE:
		return runtime.ErrExecutorSMMainHashKPositionNegative
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_POSITION_PLUS_SIZE_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainHashKPositionPlusSizeOutOfRange
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKDIGEST_ADDRESS_NOT_FOUND:
		return runtime.ErrExecutorSMMainHashKDigestAddressNotFound
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKDIGEST_NOT_COMPLETED:
		return runtime.ErrExecutorSMMainHashKDigestNotCompleted
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP:
		return runtime.ErrExecutorSMMainHashP
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_SIZE_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainHashPSizeOutOfRange
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_POSITION_NEGATIVE:
		return runtime.ErrExecutorSMMainHashPPositionNegative
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_POSITION_PLUS_SIZE_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainHashPPositionPlusSizeOutOfRange
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPDIGEST_ADDRESS_NOT_FOUND:
		return runtime.ErrExecutorSMMainHashPDigestAddressNotFound
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPDIGEST_NOT_COMPLETED:
		return runtime.ErrExecutorSMMainHashPDigestNotCompleted
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_MEMALIGN_OFFSET_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainMemAlignOffsetOutOfRange
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_MULTIPLE_FREEIN:
		return runtime.ErrExecutorSMMainMultipleFreeIn
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_ASSERT:
		return runtime.ErrExecutorSMMainAssert
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_MEMORY:
		return runtime.ErrExecutorSMMainMemory
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_STORAGE_READ_MISMATCH:
		return runtime.ErrExecutorSMMainStorageReadMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_STORAGE_WRITE_MISMATCH:
		return runtime.ErrExecutorSMMainStorageWriteMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_VALUE_MISMATCH:
		return runtime.ErrExecutorSMMainHashKValueMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_PADDING_MISMATCH:
		return runtime.ErrExecutorSMMainHashKPaddingMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_SIZE_MISMATCH:
		return runtime.ErrExecutorSMMainHashKSizeMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKLEN_LENGTH_MISMATCH:
		return runtime.ErrExecutorSMMainHashKLenLengthMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKLEN_CALLED_TWICE:
		return runtime.ErrExecutorSMMainHashKLenCalledTwice
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKDIGEST_NOT_FOUND:
		return runtime.ErrExecutorSMMainHashKDigestNotFound
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKDIGEST_DIGEST_MISMATCH:
		return runtime.ErrExecutorSMMainHashKDigestDigestMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKDIGEST_CALLED_TWICE:
		return runtime.ErrExecutorSMMainHashKDigestCalledTwice
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_VALUE_MISMATCH:
		return runtime.ErrExecutorSMMainHashPValueMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_PADDING_MISMATCH:
		return runtime.ErrExecutorSMMainHashPPaddingMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_SIZE_MISMATCH:
		return runtime.ErrExecutorSMMainHashPSizeMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPLEN_LENGTH_MISMATCH:
		return runtime.ErrExecutorSMMainHashPLenLengthMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPLEN_CALLED_TWICE:
		return runtime.ErrExecutorSMMainHashPLenCalledTwice
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPDIGEST_DIGEST_MISMATCH:
		return runtime.ErrExecutorSMMainHashPDigestDigestMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPDIGEST_CALLED_TWICE:
		return runtime.ErrExecutorSMMainHashPDigestCalledTwice
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_ARITH_MISMATCH:
		return runtime.ErrExecutorSMMainArithMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_ARITH_ECRECOVER_MISMATCH:
		return runtime.ErrExecutorSMMainArithECRecoverMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_ADD_MISMATCH:
		return runtime.ErrExecutorSMMainBinaryAddMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_SUB_MISMATCH:
		return runtime.ErrExecutorSMMainBinarySubMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_LT_MISMATCH:
		return runtime.ErrExecutorSMMainBinaryLtMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_SLT_MISMATCH:
		return runtime.ErrExecutorSMMainBinarySLtMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_EQ_MISMATCH:
		return runtime.ErrExecutorSMMainBinaryEqMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_AND_MISMATCH:
		return runtime.ErrExecutorSMMainBinaryAndMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_OR_MISMATCH:
		return runtime.ErrExecutorSMMainBinaryOrMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_XOR_MISMATCH:
		return runtime.ErrExecutorSMMainBinaryXorMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_MEMALIGN_WRITE_MISMATCH:
		return runtime.ErrExecutorSMMainMemAlignWriteMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_MEMALIGN_WRITE8_MISMATCH:
		return runtime.ErrExecutorSMMainMemAlignWrite8Mismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_MEMALIGN_READ_MISMATCH:
		return runtime.ErrExecutorSMMainMemAlignReadMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_JMPN_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainJmpnOutOfRange
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_READ_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainHashKReadOutOfRange
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_READ_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainHashPReadOutOfRange
	case ExecutorError_EXECUTOR_ERROR_INVALID_OLD_STATE_ROOT:
		return runtime.ErrExecutorErrorInvalidOldStateRoot
	case ExecutorError_EXECUTOR_ERROR_INVALID_OLD_ACC_INPUT_HASH:
		return runtime.ErrExecutorErrorInvalidOldAccInputHash
	case ExecutorError_EXECUTOR_ERROR_INVALID_CHAIN_ID:
		return runtime.ErrExecutorErrorInvalidChainId
	case ExecutorError_EXECUTOR_ERROR_INVALID_BATCH_L2_DATA:
		return runtime.ErrExecutorErrorInvalidBatchL2Data
	case ExecutorError_EXECUTOR_ERROR_INVALID_GLOBAL_EXIT_ROOT:
		return runtime.ErrExecutorErrorInvalidGlobalExitRoot
	case ExecutorError_EXECUTOR_ERROR_INVALID_COINBASE:
		return runtime.ErrExecutorErrorInvalidCoinbase
	case ExecutorError_EXECUTOR_ERROR_INVALID_FROM:
		return runtime.ErrExecutorErrorInvalidFrom
	case ExecutorError_EXECUTOR_ERROR_INVALID_DB_KEY:
		return runtime.ErrExecutorErrorInvalidDbKey
	case ExecutorError_EXECUTOR_ERROR_INVALID_DB_VALUE:
		return runtime.ErrExecutorErrorInvalidDbValue
	case ExecutorError_EXECUTOR_ERROR_INVALID_CONTRACTS_BYTECODE_KEY:
		return runtime.ErrExecutorErrorInvalidContractsBytecodeKey
	case ExecutorError_EXECUTOR_ERROR_INVALID_CONTRACTS_BYTECODE_VALUE:
		return runtime.ErrExecutorErrorInvalidContractsBytecodeValue
	case ExecutorError_EXECUTOR_ERROR_INVALID_GET_KEY:
		return runtime.ErrExecutorErrorInvalidGetKey
		// Start of V2 errors
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_SHA256:
		return runtime.ErrExecutorSMMainCountersOverflowSha256
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS:
		return runtime.ErrExecutorSMMainHashS
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_SIZE_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainHashSSizeOutOfRange
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_POSITION_NEGATIVE:
		return runtime.ErrExecutorSMMainHashSPositionNegative
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_POSITION_PLUS_SIZE_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainHashSPositionPlusSizeOutOfRange
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSDIGEST_ADDRESS_NOT_FOUND:
		return runtime.ErrExecutorSMMainHashSDigestAddressNotFound
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSDIGEST_NOT_COMPLETED:
		return runtime.ErrExecutorSMMainHashSDigestNotCompleted
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_VALUE_MISMATCH:
		return runtime.ErrExecutorSMMainHashSValueMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_PADDING_MISMATCH:
		return runtime.ErrExecutorSMMainHashSPaddingMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_SIZE_MISMATCH:
		return runtime.ErrExecutorSMMainHashSSizeMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSLEN_LENGTH_MISMATCH:
		return runtime.ErrExecutorSMMainHashSLenLengthMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSLEN_CALLED_TWICE:
		return runtime.ErrExecutorSMMainHashSLenCalledTwice
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSDIGEST_NOT_FOUND:
		return runtime.ErrExecutorSMMainHashSDigestNotFound
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSDIGEST_DIGEST_MISMATCH:
		return runtime.ErrExecutorSMMainHashSDigestDigestMismatch
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSDIGEST_CALLED_TWICE:
		return runtime.ErrExecutorSMMainHashSDigestCalledTwice
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_READ_OUT_OF_RANGE:
		return runtime.ErrExecutorSMMainHashSReadOutOfRange
	case ExecutorError_EXECUTOR_ERROR_INVALID_L1_INFO_ROOT:
		return runtime.ErrExecutorErrorInvalidL1InfoRoot
	case ExecutorError_EXECUTOR_ERROR_INVALID_FORCED_BLOCKHASH_L1:
		return runtime.ErrExecutorErrorInvalidForcedBlockhashL1
	case ExecutorError_EXECUTOR_ERROR_INVALID_L1_DATA_V2_GLOBAL_EXIT_ROOT:
		return runtime.ErrExecutorErrorInvalidL1DataV2GlobalExitRoot
	case ExecutorError_EXECUTOR_ERROR_INVALID_L1_DATA_V2_BLOCK_HASH_L1:
		return runtime.ErrExecutorErrorInvalidL1DataV2BlockHashL1
	case ExecutorError_EXECUTOR_ERROR_INVALID_L1_SMT_PROOF:
		return runtime.ErrExecutorErrorInvalidL1SmtProof
	case ExecutorError_EXECUTOR_ERROR_INVALID_BALANCE:
		return runtime.ErrExecutorErrorInvalidBalance
	case ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_LT4_MISMATCH:
		return runtime.ErrExecutorErrorSMMainBinaryLt4Mismatch
	case ExecutorError_EXECUTOR_ERROR_INVALID_NEW_STATE_ROOT:
		return runtime.ErrExecutorErrorInvalidNewStateRoot
	case ExecutorError_EXECUTOR_ERROR_INVALID_NEW_ACC_INPUT_HASH:
		return runtime.ErrExecutorErrorInvalidNewAccInputHash
	case ExecutorError_EXECUTOR_ERROR_INVALID_NEW_LOCAL_EXIT_ROOT:
		return runtime.ErrExecutorErrorInvalidNewLocalExitRoot
	case ExecutorError_EXECUTOR_ERROR_DB_KEY_NOT_FOUND:
		return runtime.ErrExecutorErrorDBKeyNotFound
	case ExecutorError_EXECUTOR_ERROR_SMT_INVALID_DATA_SIZE:
		return runtime.ErrExecutorErrorSMTInvalidDataSize
	case ExecutorError_EXECUTOR_ERROR_HASHDB_GRPC_ERROR:
		return runtime.ErrExecutorErrorHashDBGRPCError
	case ExecutorError_EXECUTOR_ERROR_STATE_MANAGER:
		return runtime.ErrExecutorErrorStateManager
	case ExecutorError_EXECUTOR_ERROR_INVALID_L1_INFO_TREE_INDEX:
		return runtime.ErrExecutorErrorInvalidL1InfoTreeIndex
	case ExecutorError_EXECUTOR_ERROR_INVALID_L1_INFO_TREE_SMT_PROOF_VALUE:
		return runtime.ErrExecutorErrorInvalidL1InfoTreeSmtProofValue
	case ExecutorError_EXECUTOR_ERROR_INVALID_WITNESS:
		return runtime.ErrExecutorErrorInvalidWitness
	case ExecutorError_EXECUTOR_ERROR_INVALID_CBOR:
		return runtime.ErrExecutorErrorInvalidCBOR
	case ExecutorError_EXECUTOR_ERROR_INVALID_DATA_STREAM:
		return runtime.ErrExecutorErrorInvalidDataStream
	case ExecutorError_EXECUTOR_ERROR_INVALID_UPDATE_MERKLE_TREE:
		return runtime.ErrExecutorErrorInvalidUpdateMerkleTree
	case ExecutorError_EXECUTOR_ERROR_UNSUPPORTED_PRECOMPILED:
		return runtime.ErrExecutorErrorUnsupportedPrecompile
	}
	return ErrExecutorUnknown
}

// ExecutorErrorCode returns the error code for a given error
func ExecutorErrorCode(err error) ExecutorError {
	switch err {
	case nil:
		return ExecutorError_EXECUTOR_ERROR_NO_ERROR
	case runtime.ErrExecutorDBError:
		return ExecutorError_EXECUTOR_ERROR_DB_ERROR
	case runtime.ErrExecutorSMMainCountersOverflowSteps:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_STEPS
	case runtime.ErrExecutorSMMainCountersOverflowKeccak:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_KECCAK
	case runtime.ErrExecutorSMMainCountersOverflowBinary:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_BINARY
	case runtime.ErrExecutorSMMainCountersOverflowMem:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_MEM
	case runtime.ErrExecutorSMMainCountersOverflowArith:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_ARITH
	case runtime.ErrExecutorSMMainCountersOverflowPadding:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_PADDING
	case runtime.ErrExecutorSMMainCountersOverflowPoseidon:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_POSEIDON
	case runtime.ErrExecutorUnsupportedForkId:
		return ExecutorError_EXECUTOR_ERROR_UNSUPPORTED_FORK_ID
	case runtime.ErrExecutorBalanceMismatch:
		return ExecutorError_EXECUTOR_ERROR_BALANCE_MISMATCH
	case runtime.ErrExecutorFEA2Scalar:
		return ExecutorError_EXECUTOR_ERROR_FEA2SCALAR
	case runtime.ErrExecutorTOS32:
		return ExecutorError_EXECUTOR_ERROR_TOS32
	case runtime.ErrExecutorSMMainInvalidUnsignedTx:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_INVALID_UNSIGNED_TX
	case runtime.ErrExecutorSMMainInvalidNoCounters:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_INVALID_NO_COUNTERS
	case runtime.ErrExecutorSMMainArithECRecoverDivideByZero:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_ARITH_ECRECOVER_DIVIDE_BY_ZERO
	case runtime.ErrExecutorSMMainAddressOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_ADDRESS_OUT_OF_RANGE
	case runtime.ErrExecutorSMMainAddressNegative:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_ADDRESS_NEGATIVE
	case runtime.ErrExecutorSMMainStorageInvalidKey:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_STORAGE_INVALID_KEY
	case runtime.ErrExecutorSMMainHashK:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK
	case runtime.ErrExecutorSMMainHashKSizeOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_SIZE_OUT_OF_RANGE
	case runtime.ErrExecutorSMMainHashKPositionNegative:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_POSITION_NEGATIVE
	case runtime.ErrExecutorSMMainHashKPositionPlusSizeOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_POSITION_PLUS_SIZE_OUT_OF_RANGE
	case runtime.ErrExecutorSMMainHashKDigestAddressNotFound:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKDIGEST_ADDRESS_NOT_FOUND
	case runtime.ErrExecutorSMMainHashKDigestNotCompleted:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKDIGEST_NOT_COMPLETED
	case runtime.ErrExecutorSMMainHashP:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP
	case runtime.ErrExecutorSMMainHashPSizeOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_SIZE_OUT_OF_RANGE
	case runtime.ErrExecutorSMMainHashPPositionNegative:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_POSITION_NEGATIVE
	case runtime.ErrExecutorSMMainHashPPositionPlusSizeOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_POSITION_PLUS_SIZE_OUT_OF_RANGE
	case runtime.ErrExecutorSMMainHashPDigestAddressNotFound:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPDIGEST_ADDRESS_NOT_FOUND
	case runtime.ErrExecutorSMMainHashPDigestNotCompleted:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPDIGEST_NOT_COMPLETED
	case runtime.ErrExecutorSMMainMemAlignOffsetOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_MEMALIGN_OFFSET_OUT_OF_RANGE
	case runtime.ErrExecutorSMMainMultipleFreeIn:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_MULTIPLE_FREEIN
	case runtime.ErrExecutorSMMainAssert:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_ASSERT
	case runtime.ErrExecutorSMMainMemory:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_MEMORY
	case runtime.ErrExecutorSMMainStorageReadMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_STORAGE_READ_MISMATCH
	case runtime.ErrExecutorSMMainStorageWriteMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_STORAGE_WRITE_MISMATCH
	case runtime.ErrExecutorSMMainHashKValueMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_VALUE_MISMATCH
	case runtime.ErrExecutorSMMainHashKPaddingMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_PADDING_MISMATCH
	case runtime.ErrExecutorSMMainHashKSizeMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_SIZE_MISMATCH
	case runtime.ErrExecutorSMMainHashKLenLengthMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKLEN_LENGTH_MISMATCH
	case runtime.ErrExecutorSMMainHashKLenCalledTwice:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKLEN_CALLED_TWICE
	case runtime.ErrExecutorSMMainHashKDigestNotFound:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKDIGEST_NOT_FOUND
	case runtime.ErrExecutorSMMainHashKDigestDigestMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKDIGEST_DIGEST_MISMATCH
	case runtime.ErrExecutorSMMainHashKDigestCalledTwice:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHKDIGEST_CALLED_TWICE
	case runtime.ErrExecutorSMMainHashPValueMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_VALUE_MISMATCH
	case runtime.ErrExecutorSMMainHashPPaddingMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_PADDING_MISMATCH
	case runtime.ErrExecutorSMMainHashPSizeMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_SIZE_MISMATCH
	case runtime.ErrExecutorSMMainHashPLenLengthMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPLEN_LENGTH_MISMATCH
	case runtime.ErrExecutorSMMainHashPLenCalledTwice:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPLEN_CALLED_TWICE
	case runtime.ErrExecutorSMMainHashPDigestDigestMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPDIGEST_DIGEST_MISMATCH
	case runtime.ErrExecutorSMMainHashPDigestCalledTwice:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHPDIGEST_CALLED_TWICE
	case runtime.ErrExecutorSMMainArithMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_ARITH_MISMATCH
	case runtime.ErrExecutorSMMainArithECRecoverMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_ARITH_ECRECOVER_MISMATCH
	case runtime.ErrExecutorSMMainBinaryAddMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_ADD_MISMATCH
	case runtime.ErrExecutorSMMainBinarySubMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_SUB_MISMATCH
	case runtime.ErrExecutorSMMainBinaryLtMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_LT_MISMATCH
	case runtime.ErrExecutorSMMainBinarySLtMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_SLT_MISMATCH
	case runtime.ErrExecutorSMMainBinaryEqMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_EQ_MISMATCH
	case runtime.ErrExecutorSMMainBinaryAndMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_AND_MISMATCH
	case runtime.ErrExecutorSMMainBinaryOrMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_OR_MISMATCH
	case runtime.ErrExecutorSMMainBinaryXorMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_XOR_MISMATCH
	case runtime.ErrExecutorSMMainMemAlignWriteMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_MEMALIGN_WRITE_MISMATCH
	case runtime.ErrExecutorSMMainMemAlignWrite8Mismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_MEMALIGN_WRITE8_MISMATCH
	case runtime.ErrExecutorSMMainMemAlignReadMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_MEMALIGN_READ_MISMATCH
	case runtime.ErrExecutorSMMainJmpnOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_JMPN_OUT_OF_RANGE
	case runtime.ErrExecutorSMMainHashKReadOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHK_READ_OUT_OF_RANGE
	case runtime.ErrExecutorSMMainHashPReadOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHP_READ_OUT_OF_RANGE
	case runtime.ErrExecutorErrorInvalidOldStateRoot:
		return ExecutorError_EXECUTOR_ERROR_INVALID_OLD_STATE_ROOT
	case runtime.ErrExecutorErrorInvalidOldAccInputHash:
		return ExecutorError_EXECUTOR_ERROR_INVALID_OLD_ACC_INPUT_HASH
	case runtime.ErrExecutorErrorInvalidChainId:
		return ExecutorError_EXECUTOR_ERROR_INVALID_CHAIN_ID
	case runtime.ErrExecutorErrorInvalidBatchL2Data:
		return ExecutorError_EXECUTOR_ERROR_INVALID_BATCH_L2_DATA
	case runtime.ErrExecutorErrorInvalidGlobalExitRoot:
		return ExecutorError_EXECUTOR_ERROR_INVALID_GLOBAL_EXIT_ROOT
	case runtime.ErrExecutorErrorInvalidCoinbase:
		return ExecutorError_EXECUTOR_ERROR_INVALID_COINBASE
	case runtime.ErrExecutorErrorInvalidFrom:
		return ExecutorError_EXECUTOR_ERROR_INVALID_FROM
	case runtime.ErrExecutorErrorInvalidDbKey:
		return ExecutorError_EXECUTOR_ERROR_INVALID_DB_KEY
	case runtime.ErrExecutorErrorInvalidDbValue:
		return ExecutorError_EXECUTOR_ERROR_INVALID_DB_VALUE
	case runtime.ErrExecutorErrorInvalidContractsBytecodeKey:
		return ExecutorError_EXECUTOR_ERROR_INVALID_CONTRACTS_BYTECODE_KEY
	case runtime.ErrExecutorErrorInvalidContractsBytecodeValue:
		return ExecutorError_EXECUTOR_ERROR_INVALID_CONTRACTS_BYTECODE_VALUE
	case runtime.ErrExecutorErrorInvalidGetKey:
		return ExecutorError_EXECUTOR_ERROR_INVALID_GET_KEY
		// Start of V2 errors
	case runtime.ErrExecutorSMMainCountersOverflowSha256:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_COUNTERS_OVERFLOW_SHA256
	case runtime.ErrExecutorSMMainHashS:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS
	case runtime.ErrExecutorSMMainHashSSizeOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_SIZE_OUT_OF_RANGE
	case runtime.ErrExecutorSMMainHashSPositionNegative:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_POSITION_NEGATIVE
	case runtime.ErrExecutorSMMainHashSPositionPlusSizeOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_POSITION_PLUS_SIZE_OUT_OF_RANGE
	case runtime.ErrExecutorSMMainHashSDigestAddressNotFound:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSDIGEST_ADDRESS_NOT_FOUND
	case runtime.ErrExecutorSMMainHashSDigestNotCompleted:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSDIGEST_NOT_COMPLETED
	case runtime.ErrExecutorSMMainHashSValueMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_VALUE_MISMATCH
	case runtime.ErrExecutorSMMainHashSPaddingMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_PADDING_MISMATCH
	case runtime.ErrExecutorSMMainHashSSizeMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_SIZE_MISMATCH
	case runtime.ErrExecutorSMMainHashSLenLengthMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSLEN_LENGTH_MISMATCH
	case runtime.ErrExecutorSMMainHashSLenCalledTwice:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSLEN_CALLED_TWICE
	case runtime.ErrExecutorSMMainHashSDigestNotFound:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSDIGEST_NOT_FOUND
	case runtime.ErrExecutorSMMainHashSDigestDigestMismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSDIGEST_DIGEST_MISMATCH
	case runtime.ErrExecutorSMMainHashSDigestCalledTwice:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHSDIGEST_CALLED_TWICE
	case runtime.ErrExecutorSMMainHashSReadOutOfRange:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_HASHS_READ_OUT_OF_RANGE
	case runtime.ErrExecutorErrorInvalidL1InfoRoot:
		return ExecutorError_EXECUTOR_ERROR_INVALID_L1_INFO_ROOT
	case runtime.ErrExecutorErrorInvalidForcedBlockhashL1:
		return ExecutorError_EXECUTOR_ERROR_INVALID_FORCED_BLOCKHASH_L1
	case runtime.ErrExecutorErrorInvalidL1DataV2GlobalExitRoot:
		return ExecutorError_EXECUTOR_ERROR_INVALID_L1_DATA_V2_GLOBAL_EXIT_ROOT
	case runtime.ErrExecutorErrorInvalidL1DataV2BlockHashL1:
		return ExecutorError_EXECUTOR_ERROR_INVALID_L1_DATA_V2_BLOCK_HASH_L1
	case runtime.ErrExecutorErrorInvalidL1SmtProof:
		return ExecutorError_EXECUTOR_ERROR_INVALID_L1_SMT_PROOF
	case runtime.ErrExecutorErrorInvalidBalance:
		return ExecutorError_EXECUTOR_ERROR_INVALID_BALANCE
	case runtime.ErrExecutorErrorSMMainBinaryLt4Mismatch:
		return ExecutorError_EXECUTOR_ERROR_SM_MAIN_BINARY_LT4_MISMATCH
	case runtime.ErrExecutorErrorInvalidNewStateRoot:
		return ExecutorError_EXECUTOR_ERROR_INVALID_NEW_STATE_ROOT
	case runtime.ErrExecutorErrorInvalidNewAccInputHash:
		return ExecutorError_EXECUTOR_ERROR_INVALID_NEW_ACC_INPUT_HASH
	case runtime.ErrExecutorErrorInvalidNewLocalExitRoot:
		return ExecutorError_EXECUTOR_ERROR_INVALID_NEW_LOCAL_EXIT_ROOT
	case runtime.ErrExecutorErrorDBKeyNotFound:
		return ExecutorError_EXECUTOR_ERROR_DB_KEY_NOT_FOUND
	case runtime.ErrExecutorErrorSMTInvalidDataSize:
		return ExecutorError_EXECUTOR_ERROR_SMT_INVALID_DATA_SIZE
	case runtime.ErrExecutorErrorHashDBGRPCError:
		return ExecutorError_EXECUTOR_ERROR_HASHDB_GRPC_ERROR
	case runtime.ErrExecutorErrorStateManager:
		return ExecutorError_EXECUTOR_ERROR_STATE_MANAGER
	case runtime.ErrExecutorErrorInvalidL1InfoTreeIndex:
		return ExecutorError_EXECUTOR_ERROR_INVALID_L1_INFO_TREE_INDEX
	case runtime.ErrExecutorErrorInvalidL1InfoTreeSmtProofValue:
		return ExecutorError_EXECUTOR_ERROR_INVALID_L1_INFO_TREE_SMT_PROOF_VALUE
	case runtime.ErrExecutorErrorInvalidWitness:
		return ExecutorError_EXECUTOR_ERROR_INVALID_WITNESS
	case runtime.ErrExecutorErrorInvalidCBOR:
		return ExecutorError_EXECUTOR_ERROR_INVALID_CBOR
	case runtime.ErrExecutorErrorInvalidDataStream:
		return ExecutorError_EXECUTOR_ERROR_INVALID_DATA_STREAM
	case runtime.ErrExecutorErrorInvalidUpdateMerkleTree:
		return ExecutorError_EXECUTOR_ERROR_INVALID_UPDATE_MERKLE_TREE
	case runtime.ErrExecutorErrorUnsupportedPrecompile:
		return ExecutorError_EXECUTOR_ERROR_UNSUPPORTED_PRECOMPILED
	}

	return ErrCodeExecutorUnknown
}
