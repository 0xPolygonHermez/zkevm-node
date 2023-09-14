package synchronizer

import (
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/metrics"
)

type ll1RollupInfoConsumerStatistics struct {
	numProcessedRollupInfo         uint64
	numProcessedBlocks             uint64
	startTime                      time.Time
	timePreviousProcessingDuration time.Duration
	startStepTime                  time.Time
	cfg                            configConsumer
}

func (l *ll1RollupInfoConsumerStatistics) onStart() {
	l.startTime = time.Now()
	l.startStepTime = time.Time{}
}

func (l *ll1RollupInfoConsumerStatistics) onStartStep() {
	l.startStepTime = time.Now()
}

func (l *ll1RollupInfoConsumerStatistics) onStartProcessIncommingRollupInfoData(rollupInfo rollupInfoByBlockRangeResult) string {
	now := time.Now()
	// Time have have been blocked in the select statement
	waitingTimeForData := now.Sub(l.startStepTime)
	blocksPerSecond := float64(l.numProcessedBlocks) / time.Since(l.startTime).Seconds()
	if l.numProcessedRollupInfo > uint64(l.cfg.numIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfoData) && waitingTimeForData > l.cfg.acceptableTimeWaitingForNewRollupInfoData {
		msg := fmt.Sprintf("wasted waiting for new rollupInfo from L1: %s last_process: %s new range: %s block_per_second: %f",
			waitingTimeForData, l.timePreviousProcessingDuration, rollupInfo.blockRange.toString(), blocksPerSecond)
		log.Warnf("consumer:: Too much wasted time (waiting to receive a new data):%s", msg)
	}
	l.numProcessedRollupInfo++
	msg := fmt.Sprintf("wasted_time_waiting_for_data [%s] last_process_time [%s] block_per_second [%f]", waitingTimeForData, l.timePreviousProcessingDuration, blocksPerSecond)
	return msg
}

func (l *ll1RollupInfoConsumerStatistics) onFinishProcessIncommingRollupInfoData(rollupInfo rollupInfoByBlockRangeResult, executionTime time.Duration, err error) {
	l.timePreviousProcessingDuration = executionTime
	if err == nil {
		l.numProcessedBlocks += uint64(len(rollupInfo.blocks))
		metrics.ProcessL1DataTime(executionTime)
	}
}
