package metrics

import (
	"time"
)

type LogTag string

type LogStatistics interface {
	CumulativeCounting(tag LogTag)
	CumulativeValue(tag LogTag, value int64)
	CumulativeTiming(tag LogTag, duration time.Duration)
	SetTag(tag LogTag, value string)
	Summary() string
	ResetStatistics()

	UpdateTimestamp(tag LogTag, tm time.Time)
}

const (
	TxCounter             LogTag = "TxCounter"
	GetTx                 LogTag = "GetTx"
	GetTxPauseCounter     LogTag = "GetTxPauseCounter"
	BatchCloseReason      LogTag = "BatchCloseReason"
	ReprocessingTxCounter LogTag = "ReProcessingTxCounter"
	FailTxCounter         LogTag = "FailTxCounter"
	NewRound              LogTag = "NewRound"
	BatchGas              LogTag = "BatchGas"

	ProcessingTxTiming         LogTag = "ProcessingTxTiming"
	ProcessingInvalidTxCounter LogTag = "ProcessingInvalidTxCounter"
	ProcessingTxCommit         LogTag = "ProcessingTxCommit"
	ProcessingTxResponse       LogTag = "ProcessingTxResponse"

	FinalizeBatchTiming             LogTag = "FinalizeBatchTiming"
	FinalizeBatchNumber             LogTag = "FinalizeBatchNumber"
	FinalizeBatchReprocessFullBatch LogTag = "FinalizeBatchReprocessFullBatch"
	FinalizeBatchCloseBatch         LogTag = "FinalizeBatchCloseBatch"
	FinalizeBatchOpenBatch          LogTag = "FinalizeBatchOpenBatch"
)
