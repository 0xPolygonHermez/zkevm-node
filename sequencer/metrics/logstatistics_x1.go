package metrics

import (
	"time"
)

// LogTag is a type used for logging tags.
type logTag string

type logStatistics interface {
	CumulativeCounting(tag logTag)
	CumulativeValue(tag logTag, value int64)
	CumulativeTiming(tag logTag, duration time.Duration)
	SetTag(tag logTag, value string)
	GetTag(tag logTag) string
	GetStatistics(tag logTag) int64
	Summary() string
	ResetStatistics()

	UpdateTimestamp(tag logTag, tm time.Time)
}

const (
	// BlockCounter is a constant for logging block counts.
	BlockCounter logTag = "BlockCounter"
	// TxCounter is a constant for logging transaction counts.
	TxCounter logTag = "TxCounter"
	// GetTx is a constant for logging tx.
	GetTx logTag = "GetTx"
	// GetTxPauseCounter is used to log the transaction pause counter.
	GetTxPauseCounter logTag = "GetTxPauseCounter"
	// BatchCloseReason is used to log the batch close reason.
	BatchCloseReason logTag = "BatchCloseReason"
	// ReprocessingTxCounter is used to log the reprocessing transaction counter.
	ReprocessingTxCounter logTag = "ReProcessingTxCounter"
	// FailTxCounter is used to log the failed transaction counter.
	FailTxCounter logTag = "FailTxCounter"
	// FailTxResourceOverCounter is used to log the failed transaction resource overflow counter.
	FailTxResourceOverCounter logTag = "FailTxResourceOverCounter"

	// NewRound is used to log new round events.
	NewRound logTag = "NewRound"
	// BatchGas is used to log batch gas-related information.
	BatchGas logTag = "BatchGas"

	// ProcessingTxTiming is used to log transaction processing time.
	ProcessingTxTiming logTag = "ProcessingTxTiming"
	// ProcessingInvalidTxCounter is used to log the processing of invalid transactions counter.
	ProcessingInvalidTxCounter logTag = "ProcessingInvalidTxCounter"
	// ProcessingTxCommit is used to log transaction commit events.
	ProcessingTxCommit logTag = "ProcessingTxCommit"
	// ProcessingTxResponse is used to log transaction response events.
	ProcessingTxResponse logTag = "ProcessingTxResponse"

	// CloseWIPL2Block is used to log close WIP L2 block events.
	CloseWIPL2Block logTag = "CloseWIPL2Block"
	// OpenNewWIPL2Block is used to log open new WIP L2 block events.
	OpenNewWIPL2Block logTag = "OpenNewWIPL2Block"
	// StoreL2Block is used to log L2 block storage events.
	StoreL2Block logTag = "StoreL2Block"

	// PoolUpdateTxStatus is used to log pool update transaction status events.
	PoolUpdateTxStatus logTag = "PoolUpdateTxStatus"
	// DeletePendingTxToStore is used to log delete pending transaction to store events.
	DeletePendingTxToStore logTag = "DeletePendingTxToStore"
	// UpdateWIPBatch is used to log update WIP batch events.
	UpdateWIPBatch logTag = "UpdateWIPBatch"
	// DSSendL2Block is used to log DS send L2 block events.
	DSSendL2Block logTag = "DSSendL2Block"
	// StateStoreL2Block is used to log state store L2 block events.
	StateStoreL2Block logTag = "StateStoreL2Block"

	// FinalizeBatchTiming is used to log batch finalization time.
	FinalizeBatchTiming logTag = "FinalizeBatchTiming"
	// FinalizeBatchNumber is used to log batch numbers.
	FinalizeBatchNumber logTag = "FinalizeBatchNumber"
	// FinalizeBatchReprocessFullBatch is used to log reprocess full batch events.
	FinalizeBatchReprocessFullBatch logTag = "FinalizeBatchReprocessFullBatch"
	// FinalizeBatchCloseBatch is used to log batch close events.
	FinalizeBatchCloseBatch logTag = "FinalizeBatchCloseBatch"
	// FinalizeBatchOpenBatch is used to log batch open events.
	FinalizeBatchOpenBatch logTag = "FinalizeBatchOpenBatch"
)
