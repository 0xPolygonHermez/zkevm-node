package metrics

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Prefix for the metrics of the state package.
	Prefix = "state_"
	// ExecutorProcessingTimeName is the name of the metric that shows the processing time in the executor.
	ExecutorProcessingTimeName = Prefix + "executor_processing_time"
	// CallerLabelName is the name of the label for the caller.
	CallerLabelName = "caller"
)

// Register the metrics for the sequencer package.
func Register() {
	histogramVecs := []metrics.HistogramVecOpts{
		{
			HistogramOpts: prometheus.HistogramOpts{
				Name: ExecutorProcessingTimeName,
				Help: "[STATE] processing time in executor",
			},
			Labels: []string{CallerLabelName},
		},
	}

	metrics.RegisterHistogramVecs(histogramVecs...)
}

// ExecutorProcessingTime observes the last processing time of the executor in the histogram vector by the provided elapsed time
// and for the given label.
func ExecutorProcessingTime(caller string, lastExecutionTime time.Duration) {
	execTimeInSeconds := float64(lastExecutionTime) / float64(time.Second)
	metrics.HistogramVecObserve(ExecutorProcessingTimeName, string(caller), execTimeInSeconds)
}
