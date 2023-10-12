package metrics

import (
	"strconv"
	"sync"
	"time"
)

var instance *logStatisticsInstance
var once sync.Once

func GetLogStatistics() LogStatistics {
	once.Do(func() {
		instance = &logStatisticsInstance{}
		instance.init()
	})
	return instance
}

type logStatisticsInstance struct {
	timestamp  map[LogTag]time.Time
	statistics map[LogTag]int64 // value maybe the counter or time.Duration(ms)
	tags       map[LogTag]string
}

func (l *logStatisticsInstance) init() {
	l.timestamp = make(map[LogTag]time.Time)
	l.statistics = make(map[LogTag]int64)
	l.tags = make(map[LogTag]string)
}

func (l *logStatisticsInstance) CumulativeCounting(tag LogTag) {
	l.statistics[tag]++
}

func (l *logStatisticsInstance) CumulativeValue(tag LogTag, value int64) {
	l.statistics[tag] += value
}

func (l *logStatisticsInstance) CumulativeTiming(tag LogTag, duration time.Duration) {
	l.statistics[tag] += duration.Milliseconds()
}

func (l *logStatisticsInstance) SetTag(tag LogTag, value string) {
	l.tags[tag] = value
}

func (l *logStatisticsInstance) UpdateTimestamp(tag LogTag, tm time.Time) {
	l.timestamp[tag] = tm
}

func (l *logStatisticsInstance) ResetStatistics() {
	l.statistics = make(map[LogTag]int64)
	l.tags = make(map[LogTag]string)
}

func (l *logStatisticsInstance) Summary() string {
	batchTotalDuration := "-"
	if key, ok := l.timestamp[NewRound]; ok {
		batchTotalDuration = strconv.Itoa(int(time.Since(key).Milliseconds()))
	}
	processTxTiming := "ProcessTx<" + strconv.Itoa(int(l.statistics[ProcessingTxTiming])) + "ms, " +
		"Commit<" + strconv.Itoa(int(l.statistics[ProcessingTxCommit])) + "ms>, " +
		"ProcessResponse<" + strconv.Itoa(int(l.statistics[ProcessingTxResponse])) + "ms>>, "

	finalizeBatchTiming := "FinalizeBatch<" + strconv.Itoa(int(l.statistics[FinalizeBatchTiming])) + "ms, " +
		"ReprocessFullBatch<" + strconv.Itoa(int(l.statistics[FinalizeBatchReprocessFullBatch])) + "ms>, " +
		"CloseBatch<" + strconv.Itoa(int(l.statistics[FinalizeBatchCloseBatch])) + "ms>, " +
		"OpenBatch<" + strconv.Itoa(int(l.statistics[FinalizeBatchOpenBatch])) + "ms>>, "

	result := "Batch<" + l.tags[FinalizeBatchNumber] + ">, " +
		"TotalDuration<" + batchTotalDuration + "ms>, " +
		"GasUsed<" + strconv.Itoa(int(l.statistics[BatchGas])) + ">, " +
		"Tx<" + strconv.Itoa(int(l.statistics[TxCounter])) + ">, " +
		"GetTx<" + strconv.Itoa(int(l.statistics[GetTx])) + "ms>, " +
		"GetTxPause<" + strconv.Itoa(int(l.statistics[GetTxPauseCounter])) + ">, " +
		"ReprocessTx<" + strconv.Itoa(int(l.statistics[ReprocessingTxCounter])) + ">, " +
		"FailTx<" + strconv.Itoa(int(l.statistics[FailTxCounter])) + ">, " +
		"InvalidTx<" + strconv.Itoa(int(l.statistics[ProcessingInvalidTxCounter])) + ">, " +
		processTxTiming +
		finalizeBatchTiming +
		"BatchCloseReason<" + l.tags[BatchCloseReason] + ">"

	return result
}
