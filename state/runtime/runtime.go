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
	// ErrOutOfCountersKeccak indicates there are not enough keccak counters to continue the execution
	ErrOutOfCountersKeccak = errors.New("not enough keccak counters to continue the execution")
	// ErrOutOfCountersBinary indicates there are not enough binary counters to continue the execution
	ErrOutOfCountersBinary = errors.New("not enough binary counters to continue the execution")
	// ErrOutOfCountersMemory indicates there are not enough memory align counters to continue the execution
	ErrOutOfCountersMemory = errors.New("not enough memory align counters counters to continue the execution")
	// ErrOutOfCountersArith indicates there are not enough arith counters to continue the execution
	ErrOutOfCountersArith = errors.New("not enough arith counters counters to continue the execution")
	// ErrOutOfCountersPadding indicates there are not enough padding counters to continue the execution
	ErrOutOfCountersPadding = errors.New("not enough padding counters counters to continue the execution")
	// ErrOutOfCountersPoseidon indicates there are not enough poseidon counters to continue the execution
	ErrOutOfCountersPoseidon = errors.New("not enough poseidon counters counters to continue the execution")
	// ErrInvalidTransaction indicates the transaction are invalid because of invalid jump dest, invalid opcode, invalid deploy
	// or invalid static tx
	ErrInvalidTransaction = errors.New("invalid transaction")
	// ErrIntrinsicInvalidSignature indicates the transaction is failing at the signature intrinsic check
	ErrIntrinsicInvalidSignature = errors.New("signature intrinsic error")
	// ErrIntrinsicInvalidChainID indicates the transaction is failing at the chain id intrinsic check
	ErrIntrinsicInvalidChainID = errors.New("chain id intrinsic error")
	// ErrIntrinsicInvalidNonce indicates the transaction is failing at the nonce intrinsic check
	ErrIntrinsicInvalidNonce = errors.New("nonce intrinsic error")
	// ErrIntrinsicInvalidGasLimit indicates the transaction is failing at the gas limit intrinsic check
	ErrIntrinsicInvalidGasLimit = errors.New("gas limit intrinsic error")
	// ErrIntrinsicInvalidBalance indicates the transaction is failing at balance intrinsic check
	ErrIntrinsicInvalidBalance = errors.New("balance intrinsic error")
	// ErrIntrinsicInvalidBatchGasLimit indicates the batch is exceeding the batch gas limit
	ErrIntrinsicInvalidBatchGasLimit = errors.New("batch gas limit intrinsic error")
	// ErrIntrinsicInvalidSenderCode indicates the sender code is invalid
	ErrIntrinsicInvalidSenderCode = errors.New("invalid sender code intrinsic error")
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
