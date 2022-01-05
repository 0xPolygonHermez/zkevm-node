package runtime

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrOutOfGas                 = errors.New("out of gas")
	ErrStackOverflow            = errors.New("stack overflow")
	ErrStackUnderflow           = errors.New("stack underflow")
	ErrNotEnoughFunds           = errors.New("not enough funds")
	ErrInsufficientBalance      = errors.New("insufficient balance for transfer")
	ErrMaxCodeSizeExceeded      = errors.New("evm: max code size exceeded")
	ErrContractAddressCollision = errors.New("contract address collision")
	ErrDepth                    = errors.New("max call depth exceeded")
	ErrExecutionReverted        = errors.New("execution was reverted")
	ErrCodeStoreOutOfGas        = errors.New("contract creation code storage out of gas")
)

// TxContext is the context of the transaction
type TxContext struct {
	GasPrice   common.Hash
	Origin     common.Address
	Coinbase   common.Address
	Number     int64
	Timestamp  int64
	GasLimit   int64
	ChainID    int64
	Difficulty common.Hash
}

// ExecutionResult includes all output after executing given evm
// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	ReturnValue []byte // Returned data from the runtime (function result or data supplied with revert opcode)
	GasLeft     uint64 // Total gas left as result of execution
	GasUsed     uint64 // Total gas used as result of execution
	Err         error  // Any error encountered during the execution, listed below
}
