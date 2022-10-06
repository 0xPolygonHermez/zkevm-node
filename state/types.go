package state

import (
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ProcessBatchRequest represents a request to process a batch.
type ProcessBatchRequest struct {
	BatchNum             uint64
	Coinbase             common.Address
	BatchL2Data          []byte
	OldStateRoot         common.Hash
	GlobalExitRoot       common.Hash
	OldLocalExitRoot     common.Hash
	EthTimestamp         uint64
	UpdateMerkleTree     bool
	GenerateExecuteTrace bool
	GenerateCallTrace    bool
}

// ProcessBatchResponse represents the response of a batch process.
type ProcessBatchResponse struct {
	CumulativeGasUsed   uint64
	IsBatchProcessed    bool
	Responses           []*ProcessTransactionResponse
	NewStateRoot        common.Hash
	NewLocalExitRoot    common.Hash
	CntKeccakHashes     uint32
	CntPoseidonHashes   uint32
	CntPoseidonPaddings uint32
	CntMemAligns        uint32
	CntArithmetics      uint32
	CntBinaries         uint32
	CntSteps            uint32
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
