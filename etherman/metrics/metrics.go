package metrics

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Prefix for the metrics of the etherman package.
	Prefix = "etherman_"

	// ReadAndProcessAllEventsTimeName is the name of the label read and process all event.
	ReadAndProcessAllEventsTimeName = Prefix + "read_and_process_all_event_time"

	// ProcessAllEventTimeName is the name of the label to process all event.
	ProcessAllEventTimeName = Prefix + "process_all_event_time"

	// ProcessSingleEventTimeName is the name of the label to process a single event.
	ProcessSingleEventTimeName = Prefix + "process_single_event_time"

	// GetEventsTimeName is the name of the label to get L1 events.
	GetEventsTimeName = Prefix + "get_events_time"

	// GetForksTimeName is the name of the label to get forkIDs intervals.
	GetForksTimeName = Prefix + "get_forkIDs_time"

	// VerifyGenBlockTimeName is the name of the label to verify the genesis block.
	VerifyGenBlockTimeName = Prefix + "verify_genesisBlockNum_time"

	// EventCounterName is the name of the label to count the processed events.
	EventCounterName = Prefix + "processed_events_counter"
)

// Register the metrics for the etherman package.
func Register() {
	var (
		counters   []prometheus.CounterOpts
		histograms []prometheus.HistogramOpts
	)

	counters = []prometheus.CounterOpts{
		{
			Name: EventCounterName,
			Help: "[ETHERMAN] count processed events",
		},
	}

	histograms = []prometheus.HistogramOpts{
		{
			Name: ReadAndProcessAllEventsTimeName,
			Help: "[ETHERMAN] read and process all event time",
		},
		{
			Name: ProcessAllEventTimeName,
			Help: "[ETHERMAN] process all event time",
		},
		{
			Name: ProcessSingleEventTimeName,
			Help: "[ETHERMAN] process single event time",
		},
		{
			Name: GetEventsTimeName,
			Help: "[ETHERMAN] get L1 events time",
		},
		{
			Name: GetForksTimeName,
			Help: "[ETHERMAN] get forkIDs time",
		},
		{
			Name: VerifyGenBlockTimeName,
			Help: "[ETHERMAN] verify genesis block number time",
		},
	}

	metrics.RegisterCounters(counters...)
	metrics.RegisterHistograms(histograms...)
}

// ReadAndProcessAllEventsTime observes the time read and process all event on the histogram.
func ReadAndProcessAllEventsTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(ReadAndProcessAllEventsTimeName, execTimeInSeconds)
}

// ProcessAllEventTime observes the time to process all event on the histogram.
func ProcessAllEventTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(ProcessAllEventTimeName, execTimeInSeconds)
}

// ProcessSingleEventTime observes the time to process a single event on the histogram.
func ProcessSingleEventTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(ProcessSingleEventTimeName, execTimeInSeconds)
}

// GetEventsTime observes the time to get the events from L1 on the histogram.
func GetEventsTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(GetEventsTimeName, execTimeInSeconds)
}

// GetForksTime observes the time to get the forkIDs on the histogram.
func GetForksTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(GetForksTimeName, execTimeInSeconds)
}

// VerifyGenBlockTime observes the time for etherman to verify the genesis blocknumber on the histogram.
func VerifyGenBlockTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(VerifyGenBlockTimeName, execTimeInSeconds)
}

// EventCounter increases the counter for the processed events
func EventCounter() {
	metrics.CounterInc(EventCounterName)
}
