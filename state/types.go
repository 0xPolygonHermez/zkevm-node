package state

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/pool"

	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TouchedAddress struct {
	Address common.Address
	Nonce   *uint64
	Balance *big.Int
}

// ProcessBatchRequest represents the request of a batch process.
type ProcessBatchRequest struct {
	BatchNumber      uint64
	StateRoot        common.Hash
	GlobalExitRoot   common.Hash
	OldAccInputHash  common.Hash
	TxData           []byte
	SequencerAddress common.Address
	Timestamp        uint64
	IsFirstTx        bool
	Caller           CallerLabel
}

// ProcessBatchResponse represents the response of a batch process.
type ProcessBatchResponse struct {
	NewStateRoot     common.Hash
	NewAccInputHash  common.Hash
	NewLocalExitRoot common.Hash
	NewBatchNumber   uint64
	UsedZkCounters   pool.ZkCounters
	Responses        []*ProcessTransactionResponse
	Error            error
	IsBatchProcessed bool
	TouchedAddresses map[common.Address]*TouchedAddress
}

// ProcessTransactionResponse represents the response of a tx process.
type ProcessTransactionResponse struct {
	// TxHash is the hash of the transaction
	TxHash common.Hash
	// Type indicates legacy transaction
	// It will be always 0 (legacy) in the executor
	Type uint32
	// ReturnValue is the returned data from the runtime (function result or data supplied with revert opcode)
	ReturnValue []byte
	// GasLeft is the total gas left as result of execution
	GasLeft uint64
	// GasUsed is the total gas used as result of execution or gas estimation
	GasUsed uint64
	// GasRefunded is the total gas refunded as result of execution
	GasRefunded uint64
	// Error represents any error encountered during the execution
	Error error
	// CreateAddress is the new SC Address in case of SC creation
	CreateAddress common.Address
	// StateRoot is the State Root
	StateRoot common.Hash
	// Logs emitted by LOG opcode
	Logs []*types.Log
	// IsProcessed indicates if this tx didn't fit into the batch
	IsProcessed bool
	// Tx is the whole transaction object
	Tx types.Transaction
	// ExecutionTrace contains the traces produced in the execution
	ExecutionTrace []instrumentation.StructLog
	// CallTrace contains the call trace.
	CallTrace instrumentation.ExecutorTrace
}
