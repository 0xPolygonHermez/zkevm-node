package sequencer

import (
	"fmt"
	"time"
)

//           |-----------------------------------------------------------------------------| -> totalTime
//                        |------------|    |-------------------------|                      -> transactionsTime
//           |-newL2Block-|----tx 1----|    |---tx 2---|-----tx 3-----|  |-----l2Block-----|
// sequencer |sssss     ss|sss       ss|    |sss     ss|sss         ss|  |ssss           ss| -> sequencerTime
//  executor |     xxxxx  |   xxxxxxx  |    |   xxxxx  |   xxxxxxxxx  |  |    xxxxxxxxxxx  | -> executorTime
//      idle |                         |iiii|          |              |ii|                 | -> idleTime
//

type processTimes struct {
	sequencer time.Duration
	executor  time.Duration
}

func (p *processTimes) total() time.Duration {
	return p.sequencer + p.executor
}

func (p *processTimes) sub(ptSub processTimes) {
	p.sequencer -= ptSub.sequencer
	p.executor -= ptSub.executor
}

func (p *processTimes) sumUp(ptSumUp processTimes) {
	p.sequencer += ptSumUp.sequencer
	p.executor += ptSumUp.executor
}

type metrics struct {
	closedAt           time.Time
	txsCount           int64
	idleTime           time.Duration
	newL2BlockTimes    processTimes
	transactionsTimes  processTimes
	l2BlockTimes       processTimes
	estimatedTxsPerSec float64
}

func (m *metrics) sub(mSub metrics) {
	m.txsCount -= mSub.txsCount
	m.idleTime -= mSub.idleTime
	m.newL2BlockTimes.sub(mSub.newL2BlockTimes)
	m.transactionsTimes.sub(mSub.transactionsTimes)
	m.l2BlockTimes.sub(mSub.l2BlockTimes)
}

func (m *metrics) sumUp(mSumUp metrics) {
	m.txsCount += mSumUp.txsCount
	m.idleTime += mSumUp.idleTime
	m.newL2BlockTimes.sumUp(mSumUp.newL2BlockTimes)
	m.transactionsTimes.sumUp(mSumUp.transactionsTimes)
	m.l2BlockTimes.sumUp(mSumUp.l2BlockTimes)
}

func (m *metrics) executorTime() time.Duration {
	return m.newL2BlockTimes.executor + m.transactionsTimes.executor + m.l2BlockTimes.executor
}

func (m *metrics) sequencerTime() time.Duration {
	return m.newL2BlockTimes.sequencer + m.transactionsTimes.sequencer + m.l2BlockTimes.sequencer
}

func (m *metrics) totalTime() time.Duration {
	return m.newL2BlockTimes.total() + m.transactionsTimes.total() + m.l2BlockTimes.total() + m.idleTime
}

func (m *metrics) close(createdAt time.Time, txsCount int64) {
	// Compute pending fields
	m.closedAt = time.Now()
	totalTime := time.Since(createdAt)
	m.txsCount = txsCount
	m.transactionsTimes.sequencer = totalTime - m.idleTime - m.newL2BlockTimes.total() - m.transactionsTimes.executor - m.l2BlockTimes.total()

	// Compute performance
	if m.txsCount > 0 {
		// timePerTxuS is the average time spent per tx. This includes the l2Block time since the processing time of this section is proportional to the number of txs
		timePerTxuS := (m.transactionsTimes.total() + m.l2BlockTimes.total()).Microseconds() / m.txsCount
		// estimatedTxs is the number of transactions that we estimate could have been processed in the block
		estimatedTxs := float64(totalTime.Microseconds()-m.newL2BlockTimes.total().Microseconds()) / float64(timePerTxuS)
		// estimatedTxxPerSec is the estimated transactions per second
		m.estimatedTxsPerSec = estimatedTxs / totalTime.Seconds()
	}
}

func (m *metrics) log() string {
	return fmt.Sprintf("txs: %d, estimated txs/s: %.1f, time: {total: %d, idle: %d, sequencer: {total: %d, newL2Block: %d, txs: %d, l2Block: %d}, executor: {total: %d, newL2Block: %d, txs: %d, l2Block: %d}",
		m.txsCount, m.estimatedTxsPerSec, m.totalTime().Microseconds(), m.idleTime.Microseconds(),
		m.sequencerTime().Microseconds(), m.newL2BlockTimes.sequencer.Microseconds(), m.transactionsTimes.sequencer.Microseconds(), m.l2BlockTimes.sequencer.Microseconds(),
		m.executorTime().Microseconds(), m.newL2BlockTimes.executor.Microseconds(), m.transactionsTimes.executor.Microseconds(), m.l2BlockTimes.executor.Microseconds())
}

type intervalMetrics struct {
	l2Blocks    []*metrics
	maxInterval time.Duration
	metrics
	estimatedTxsPerSecAcc   float64
	estimatedTxsPerSecCount int64
}

func newIntervalMetrics(maxInterval time.Duration) *intervalMetrics {
	return &intervalMetrics{
		l2Blocks:    []*metrics{},
		maxInterval: maxInterval,
		metrics:     metrics{},
	}
}

func (i *intervalMetrics) cleanUp() {
	now := time.Now()
	ct := 0
	for {
		if len(i.l2Blocks) == 0 {
			return
		}
		l2Block := i.l2Blocks[0]
		if l2Block.closedAt.Add(i.maxInterval).Before(now) {
			// Subtract l2Block metrics from accumulated values
			i.sub(*l2Block)
			if l2Block.txsCount > 0 {
				i.estimatedTxsPerSecAcc -= i.estimatedTxsPerSec
				i.estimatedTxsPerSecCount--
			}
			// Remove from l2Blocks
			i.l2Blocks = i.l2Blocks[1:]
			ct++
		} else {
			break
		}
	}

	if ct > 0 {
		// Compute performance
		i.computeEstimatedTxsPerSec()
	}
}

func (i *intervalMetrics) addL2BlockMetrics(l2Block metrics) {
	i.cleanUp()

	i.sumUp(l2Block)
	if l2Block.txsCount > 0 {
		i.estimatedTxsPerSecAcc += l2Block.estimatedTxsPerSec
		i.estimatedTxsPerSecCount++
		i.computeEstimatedTxsPerSec()
	}

	i.l2Blocks = append(i.l2Blocks, &l2Block)
}

func (i *intervalMetrics) computeEstimatedTxsPerSec() {
	if i.estimatedTxsPerSecCount > 0 {
		i.estimatedTxsPerSec = i.estimatedTxsPerSecAcc / float64(i.estimatedTxsPerSecCount)
	} else {
		i.estimatedTxsPerSecCount = 0
	}
}

func (i *intervalMetrics) startsAt() time.Time {
	return time.Now().Add(-i.maxInterval)
}
