package synchronizer

import (
	"fmt"
	"time"
)

// This object keep track of the statistics of the process, to be able to estimate the ETA
type l1RollupInfoProducerStatistics struct {
	initialBlockNumber  uint64
	lastBlockNumber     uint64
	numRollupInfoOk     uint64
	numRollupInfoErrors uint64
	numRetrievedBlocks  uint64
	startTime           time.Time
	lastShowUpTime      time.Time
}

func newRollupInfoProducerStatistics(startingBlockNumber uint64) l1RollupInfoProducerStatistics {
	return l1RollupInfoProducerStatistics{
		initialBlockNumber: startingBlockNumber,
		startTime:          time.Now(),
	}
}

func (l *l1RollupInfoProducerStatistics) updateLastBlockNumber(lastBlockNumber uint64) {
	l.lastBlockNumber = lastBlockNumber
}

func (l *l1RollupInfoProducerStatistics) updateNumRollupInfoOk(numRollupInfo uint, numBlocks uint64) {
	l.numRollupInfoOk += uint64(numRollupInfo)
	l.numRetrievedBlocks += uint64(numBlocks)
}

func (l *l1RollupInfoProducerStatistics) updateNumRollupInfoErrors(numRollupInfoWithErrors uint) {
	l.numRollupInfoErrors += uint64(numRollupInfoWithErrors)
}

func (l *l1RollupInfoProducerStatistics) getETA() string {
	numTotalOfBlocks := l.lastBlockNumber - l.initialBlockNumber
	if l.numRetrievedBlocks == 0 {
		return "N/A"
	}
	elapsedTime := time.Since(l.startTime)
	eta := time.Duration(float64(elapsedTime) / float64(l.numRetrievedBlocks) * float64(numTotalOfBlocks-l.numRetrievedBlocks))
	percent := float64(l.numRetrievedBlocks) / float64(numTotalOfBlocks) * conversionFactorPercentage
	blocks_per_seconds := float64(l.numRetrievedBlocks) / float64(elapsedTime.Seconds())
	return fmt.Sprintf("ETA: %s percent:%2.2f  blocks_per_seconds:%2.2f pending_block:%v/%v num_errors:%v",
		eta, percent, blocks_per_seconds, l.numRetrievedBlocks, numTotalOfBlocks, l.numRollupInfoErrors)
}
