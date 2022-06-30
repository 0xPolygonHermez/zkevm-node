package statev2

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/instrumentation"
)

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

type ProcessBatchResponse struct {
	CumulativeGasUsed   uint64
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

type ProcessTransactionResponse struct {
	// Hash of the transaction
	TxHash common.Hash
	// Type indicates legacy transaction
	// It will be always 0 (legacy) in the executor
	Type uint32
	// Returned data from the runtime (function result or data supplied with revert opcode)
	ReturnValue []byte
	// Total gas left as result of execution
	GasLeft uint64
	// Total gas used as result of execution or gas estimation
	GasUsed uint64
	// Total gas refunded as result of execution
	GasRefunded uint64
	// Any error encountered during the execution
	Error string
	// New SC Address in case of SC creation
	CreateAddress common.Address
	// State Root
	StateRoot common.Hash
	// Logs emitted by LOG opcode
	Logs []types.Log
	// Indicates if this tx didn't fit into the batch
	UnprocessedTransaction bool
	// Traces
	ExecutionTrace []instrumentation.StructLog
	CallTrace      instrumentation.ExecutorTrace
}
