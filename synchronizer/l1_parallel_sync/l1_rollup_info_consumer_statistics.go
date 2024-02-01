package l1_parallel_sync

import (
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/metrics"
)

type l1RollupInfoConsumerStatistics struct {
	numProcessedRollupInfo             uint64
	numProcessedRollupInfoForCheckTime uint64
	numProcessedBlocks                 uint64
	startTime                          time.Time
	timePreviousProcessingDuration     time.Duration
	startStepTime                      time.Time
	cfg                                ConfigConsumer
}

func (l *l1RollupInfoConsumerStatistics) onStart() {
	l.startTime = time.Now()
	l.startStepTime = time.Time{}
	l.numProcessedRollupInfoForCheckTime = 0
}

func (l *l1RollupInfoConsumerStatistics) onStartStep() {
	l.startStepTime = time.Now()
}

func (l *l1RollupInfoConsumerStatistics) onReset() {
	l.numProcessedRollupInfoForCheckTime = 0
	l.startStepTime = time.Time{}
}

func (l *l1RollupInfoConsumerStatistics) onStartProcessIncommingRollupInfoData(rollupInfo rollupInfoByBlockRangeResult) string {
	now := time.Now()
	// Time have have been blocked in the select statement
	waitingTimeForData := now.Sub(l.startStepTime)
	blocksPerSecond := float64(l.numProcessedBlocks) / time.Since(l.startTime).Seconds()
	generatedWarning := false
	if l.numProcessedRollupInfoForCheckTime > uint64(l.cfg.ApplyAfterNumRollupReceived) && waitingTimeForData > l.cfg.AceptableInacctivityTime {
		msg := fmt.Sprintf("wasted waiting for new rollupInfo from L1: %s last_process: %s new range: %s block_per_second: %f",
			waitingTimeForData, l.timePreviousProcessingDuration, rollupInfo.blockRange.String(), blocksPerSecond)
		log.Warnf("consumer:: Too much wasted time (waiting to receive a new data):%s", msg)
		generatedWarning = true
	}
	l.numProcessedRollupInfo++
	l.numProcessedRollupInfoForCheckTime++
	msg := fmt.Sprintf("wasted_time_waiting_for_data [%s] last_process_time [%s] block_per_second [%f]",
		waitingTimeForData.Round(time.Second).String(),
		l.timePreviousProcessingDuration,
		blocksPerSecond)
	if waitingTimeForData > l.cfg.AceptableInacctivityTime {
		msg = msg + " WASTED_TIME_EXCEED "
	}
	if generatedWarning {
		msg = msg + " WARNING_WASTED_TIME "
	}
	return msg
}

func (l *l1RollupInfoConsumerStatistics) onFinishProcessIncommingRollupInfoData(rollupInfo rollupInfoByBlockRangeResult, executionTime time.Duration, err error) {
	l.timePreviousProcessingDuration = executionTime
	if err == nil {
		l.numProcessedBlocks += uint64(len(rollupInfo.blocks))
		metrics.ProcessL1DataTime(executionTime)
	}
}
