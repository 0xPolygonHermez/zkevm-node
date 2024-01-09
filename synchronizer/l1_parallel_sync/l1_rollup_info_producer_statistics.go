package l1_parallel_sync

import (
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/metrics"
)

// This object keep track of the statistics of the process, to be able to estimate the ETA
type l1RollupInfoProducerStatistics struct {
	initialBlockNumber              uint64
	lastBlockNumber                 uint64
	numRollupInfoOk                 uint64
	numRollupInfoErrors             uint64
	numRetrievedBlocks              uint64
	startTime                       time.Time
	lastShowUpTime                  time.Time
	accumulatedTimeProcessingRollup time.Duration
	timeProvider                    common.TimeProvider
}

func newRollupInfoProducerStatistics(startingBlockNumber uint64, timeProvider common.TimeProvider) l1RollupInfoProducerStatistics {
	return l1RollupInfoProducerStatistics{
		initialBlockNumber:              startingBlockNumber,
		startTime:                       timeProvider.Now(),
		timeProvider:                    timeProvider,
		accumulatedTimeProcessingRollup: time.Duration(0),
	}
}

func (l *l1RollupInfoProducerStatistics) reset(startingBlockNumber uint64) {
	l.initialBlockNumber = startingBlockNumber
	l.startTime = l.timeProvider.Now()
	l.numRollupInfoOk = 0
	l.numRollupInfoErrors = 0
	l.numRetrievedBlocks = 0
	l.lastShowUpTime = l.timeProvider.Now()
}

func (l *l1RollupInfoProducerStatistics) updateLastBlockNumber(lastBlockNumber uint64) {
	l.lastBlockNumber = lastBlockNumber
}

func (l *l1RollupInfoProducerStatistics) onResponseRollupInfo(result responseRollupInfoByBlockRange) {
	metrics.ReadL1DataTime(result.generic.duration)
	isOk := (result.generic.err == nil)
	if isOk {
		l.numRollupInfoOk++
		l.numRetrievedBlocks += result.result.blockRange.len()
		l.accumulatedTimeProcessingRollup += result.generic.duration
	} else {
		l.numRollupInfoErrors++
	}
}

func (l *l1RollupInfoProducerStatistics) getStatisticsDebugString() string {
	numTotalOfBlocks := l.lastBlockNumber - l.initialBlockNumber
	if l.numRetrievedBlocks == 0 {
		return "N/A"
	}
	now := l.timeProvider.Now()
	elapsedTime := now.Sub(l.startTime)
	eta := l.getEstimatedTimeOfArrival()
	percent := l.getPercent()
	blocksPerSeconds := l.getBlocksPerSecond(elapsedTime)
	return fmt.Sprintf(" EstimatedTimeOfArrival: %s percent:%2.2f  blocks_per_seconds:%2.2f pending_block:%v/%v num_errors:%v",
		eta, percent, blocksPerSeconds, l.numRetrievedBlocks, numTotalOfBlocks, l.numRollupInfoErrors)
}

func (l *l1RollupInfoProducerStatistics) getEstimatedTimeOfArrival() time.Duration {
	numTotalOfBlocks := l.lastBlockNumber - l.initialBlockNumber
	if l.numRetrievedBlocks == 0 {
		return time.Duration(0)
	}
	elapsedTime := time.Since(l.startTime)
	eta := time.Duration(float64(elapsedTime) / float64(l.numRetrievedBlocks) * float64(numTotalOfBlocks-l.numRetrievedBlocks))
	return eta
}

func (l *l1RollupInfoProducerStatistics) getPercent() float64 {
	numTotalOfBlocks := l.lastBlockNumber - l.initialBlockNumber
	percent := float64(l.numRetrievedBlocks) / float64(numTotalOfBlocks) * conversionFactorPercentage
	return percent
}

func (l *l1RollupInfoProducerStatistics) getBlocksPerSecond(elapsedTime time.Duration) float64 {
	blocksPerSeconds := float64(l.numRetrievedBlocks) / elapsedTime.Seconds()
	return blocksPerSeconds
}
