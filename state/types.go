package state

import (
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ProcessRequest represents the request of a batch process.
type ProcessRequest struct {
	BatchNumber     uint64
	GlobalExitRoot  common.Hash
	OldStateRoot    common.Hash
	OldAccInputHash common.Hash
	Transactions    []byte
	Coinbase        common.Address
	Timestamp       uint64
	Caller          CallerLabel
}

// ProcessBatchResponse represents the response of a batch process.
type ProcessBatchResponse struct {
	NewStateRoot       common.Hash
	NewAccInputHash    common.Hash
	NewLocalExitRoot   common.Hash
	NewBatchNumber     uint64
	UsedZkCounters     ZKCounters
	Responses          []*ProcessTransactionResponse
	ExecutorError      error
	IsBatchProcessed   bool
	ReadWriteAddresses map[common.Address]*InfoReadWrite
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
	// RomError represents any error encountered during the execution
	RomError error
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

// ZKCounters counters for the tx
type ZKCounters struct {
	CumulativeGasUsed    uint64
	UsedKeccakHashes     uint32
	UsedPoseidonHashes   uint32
	UsedPoseidonPaddings uint32
	UsedMemAligns        uint32
	UsedArithmetics      uint32
	UsedBinaries         uint32
	UsedSteps            uint32
}

// SumUp sum ups zk counters with passed tx zk counters
func (z *ZKCounters) SumUp(other ZKCounters) {
	z.CumulativeGasUsed += other.CumulativeGasUsed
	z.UsedKeccakHashes += other.UsedKeccakHashes
	z.UsedPoseidonHashes += other.UsedPoseidonHashes
	z.UsedPoseidonPaddings += other.UsedPoseidonPaddings
	z.UsedMemAligns += other.UsedMemAligns
	z.UsedArithmetics += other.UsedArithmetics
	z.UsedBinaries += other.UsedBinaries
	z.UsedSteps += other.UsedSteps
}

// Sub subtract zk counters with passed zk counters (not safe)
func (z *ZKCounters) Sub(other ZKCounters) error {
	// ZKCounters
	if other.CumulativeGasUsed > z.CumulativeGasUsed {
		return GetZKCounterError("CumulativeGasUsed")
	}
	if other.UsedKeccakHashes > z.UsedKeccakHashes {
		return GetZKCounterError("UsedKeccakHashes")
	}
	if other.UsedPoseidonHashes > z.UsedPoseidonHashes {
		return GetZKCounterError("UsedPoseidonHashes")
	}
	if other.UsedPoseidonPaddings > z.UsedPoseidonPaddings {
		return fmt.Errorf("underflow ZKCounter: UsedPoseidonPaddings")
	}
	if other.UsedMemAligns > z.UsedMemAligns {
		return GetZKCounterError("UsedMemAligns")
	}
	if other.UsedArithmetics > z.UsedArithmetics {
		return GetZKCounterError("UsedArithmetics")
	}
	if other.UsedBinaries > z.UsedBinaries {
		return GetZKCounterError("UsedBinaries")
	}
	if other.UsedSteps > z.UsedSteps {
		return GetZKCounterError("UsedSteps")
	}

	z.CumulativeGasUsed -= other.CumulativeGasUsed
	z.UsedKeccakHashes -= other.UsedKeccakHashes
	z.UsedPoseidonHashes -= other.UsedPoseidonHashes
	z.UsedPoseidonPaddings -= other.UsedPoseidonPaddings
	z.UsedMemAligns -= other.UsedMemAligns
	z.UsedArithmetics -= other.UsedArithmetics
	z.UsedBinaries -= other.UsedBinaries
	z.UsedSteps -= other.UsedSteps

	return nil
}

// InfoReadWrite has information about modified addresses during the execution
type InfoReadWrite struct {
	Address common.Address
	Nonce   *uint64
	Balance *big.Int
}

const (
	// DebugInfoErrorType_EXECUTOR_ERROR indicates a error happened in the executor
	DebugInfoErrorType_EXECUTOR_ERROR = "EXECUTOR ERROR"
	// DebugInfoErrorType_OOC_ERROR_ON_REPROCESS_FULL_BATCH indicates and OOC error happened in the executor when reprocessing a full batch
	DebugInfoErrorType_OOC_ERROR_ON_REPROCESS_FULL_BATCH = "OOC ON REPROCESS FULL BATCH"
	// DebugInfoErrorType_EXECUTOR_RLP_ERROR indicates a error happened decoding the RLP returned by the executor
	DebugInfoErrorType_EXECUTOR_RLP_ERROR = "EXECUTOR RLP ERROR"
	// DebugInfoErrorType_FINALIZER_HALT indicates a fatal error happened in the finalizer when trying to close a batch
	DebugInfoErrorType_FINALIZER_HALT = "FINALIZER HALT"
)

// DebugInfo allows handling runtime debug info
type DebugInfo struct {
	ErrorType string
	Timestamp time.Time
	Payload   string
}

// TraceConfig sets the debug configuration for the executor
type TraceConfig struct {
	DisableStorage   bool
	DisableStack     bool
	EnableMemory     bool
	EnableReturnData bool
	Tracer           *string
}

// TrustedReorg represents a trusted reorg
type TrustedReorg struct {
	BatchNumber uint64
	Reason      string
}

const (
	// EventType_Prexecution_OOC indicates a preexecution out of couters error
	EventType_Prexecution_OOC = "PREEXECUTION OOC"
	// EventType_Prexecution_OOG indicates a preexecution out of gas error
	EventType_Prexecution_OOG = "PREEXECUTION OOG"
	// EventType_ZKCounters_Diff indicates big different in preexecution and execution regarding ZKCounters
	EventType_ZKCounters_Diff = "ZK COUNTERS DIFF"
)

// Event represents a event that may be investigated
type Event struct {
	EventType string
	Timestamp time.Time
	IP        string
	TxHash    common.Hash
	Payload   string
}

// HexToAddressPtr create an address from a hex and returns its pointer
func HexToAddressPtr(hex string) *common.Address {
	a := common.HexToAddress(hex)
	return &a
}

// HexToHashPtr create a hash from a hex and returns its pointer
func HexToHashPtr(hex string) *common.Hash {
	h := common.HexToHash(hex)
	return &h
}

// AddressPtr returns a pointer to the provided address
func AddressPtr(i common.Address) *common.Address {
	return &i
}

// HashPtr returns a pointer to the provided hash
func HashPtr(h common.Hash) *common.Hash {
	return &h
}
