package sequencer

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	seqMetrics "github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
)

func (f *finalizer) tryToSleep() {
	fullBatchSleepDuration := getFullBatchSleepDuration(f.cfg.FullBatchSleepDuration.Duration)
	if fullBatchSleepDuration > 0 {
		log.Infof("Slow down sequencer: %v", fullBatchSleepDuration)
		time.Sleep(fullBatchSleepDuration)
		seqMetrics.GetLogStatistics().CumulativeCounting(seqMetrics.GetTxPauseCounter)
	}
}
