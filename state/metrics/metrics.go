package metrics

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prefix                 = "state_"
	executorProcessingTime = prefix + "executor_processing_time"

	callerLabelName = "caller"
)

// Register the metrics for the sequencer package.
func Register() {
	histogramVecs := []metrics.HistogramVecOpts{
		{
			HistogramOpts: prometheus.HistogramOpts{
				Name: executorProcessingTime,
				Help: "[STATE] processing time in executor",
			},
			Labels: []string{callerLabelName},
		},
	}

	metrics.RegisterHistogramVecs(histogramVecs...)
}

// ExecutorProcessingTime observes the last processing time of the executor in the histogram vector by the provided elapsed time
// and for the given label.
func ExecutorProcessingTime(caller string, lastExecutionTime time.Duration) {
	execTimeInSeconds := float64(lastExecutionTime) / float64(time.Second)
	metrics.HistogramVecObserve(executorProcessingTime, string(caller), execTimeInSeconds)
}
