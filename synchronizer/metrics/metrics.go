package metrics

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Prefix for the metrics of the synchronizer package.
	Prefix = "synchronizer_"

	// InitializationTimeName is the name of the label for the initialization of the synchronizer.
	InitializationTimeName = Prefix + "initialization_time"

	// FullTrustedSyncTimeName is the name of the label for the synchronization of the trusted state.
	FullTrustedSyncTimeName = Prefix + "full_trusted_sync_time"

	// FullL1SyncTimeName is the name of the label for the synchronization of the L1 state.
	FullL1SyncTimeName = Prefix + "full_L1_sync_time"

	// FullSyncIterationTimeName is the name of the label for a full synchronization iteration.
	FullSyncIterationTimeName = Prefix + "full_L1_sync_time"

	// ReadL1DataTimeName is the name of the label for a full synchronization iteration.
	ReadL1DataTimeName = Prefix + "read_L1_data_time"

	// ProcessL1DataTimeName is the name of the label for a full synchronization iteration.
	ProcessL1DataTimeName = Prefix + "process_L1_data_time"

	// GetTrustedBatchNumberTimeName is the name of the label for a full synchronization iteration.
	GetTrustedBatchNumberTimeName = Prefix + "get_trusted_batchNumber_time"

	// GetTrustedBatchInfoTimeName is the name of the label for a full synchronization iteration.
	GetTrustedBatchInfoTimeName = Prefix + "get_trusted_batchInfo_time"

	// ProcessTrustedBatchTimeName is the name of the label for a full synchronization iteration.
	ProcessTrustedBatchTimeName = Prefix + "process_trusted_batch_time"

	// TrustedBatchCleanCounterName is the name of the label for the counter of trusted batch resets.
	TrustedBatchCleanCounterName = Prefix + "trusted_batch_clean_counter"
)

// Register the metrics for the synchronizer package.
func Register() {
	var (
		counters    []prometheus.CounterOpts
		histograms  []prometheus.HistogramOpts
	)

	counters = []prometheus.CounterOpts{
		{
			Name: TrustedBatchCleanCounterName,
			Help: "[SYNCHRONIZER] count of trusted Batch clean",
		},
	}

	histograms = []prometheus.HistogramOpts{
		{
			Name: InitializationTimeName,
			Help: "[SYNCHRONIZER] initialization time",
		},
		{
			Name: FullTrustedSyncTimeName,
			Help: "[SYNCHRONIZER] full trusted state synchronization time",
		},
		{
			Name: FullL1SyncTimeName,
			Help: "[SYNCHRONIZER] full L1 synchronization time",
		},
		{
			Name: FullSyncIterationTimeName,
			Help: "[SYNCHRONIZER] full synchronization iteration time",
		},
		{
			Name: ReadL1DataTimeName,
			Help: "[SYNCHRONIZER] read L1 data time",
		},
		{
			Name: ProcessL1DataTimeName,
			Help: "[SYNCHRONIZER] process L1 data time",
		},
		{
			Name: GetTrustedBatchNumberTimeName,
			Help: "[SYNCHRONIZER] get trusted batchNumber time",
		},
		{
			Name: GetTrustedBatchInfoTimeName,
			Help: "[SYNCHRONIZER] get trusted batchInfo time",
		},
		{
			Name: ProcessTrustedBatchTimeName,
			Help: "[SYNCHRONIZER] process trusted batch time",
		},
	}

	metrics.RegisterCounters(counters...)
	metrics.RegisterHistograms(histograms...)
}

// InitializationTime observes the time initializing the synchronizer on the histogram.
func InitializationTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(InitializationTimeName, execTimeInSeconds)
}

// FullTrustedSyncTime observes the time for synchronize the trusted state on the histogram.
func FullTrustedSyncTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(FullTrustedSyncTimeName, execTimeInSeconds)
}

// FullL1SyncTime observes the time for synchronize the trusted state on the histogram.
func FullL1SyncTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(FullL1SyncTimeName, execTimeInSeconds)
}

// FullSyncIterationTime observes the time for synchronize the trusted state on the histogram.
func FullSyncIterationTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(FullSyncIterationTimeName, execTimeInSeconds)
}

// ReadL1DataTime observes the time for synchronize the trusted state on the histogram.
func ReadL1DataTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(ReadL1DataTimeName, execTimeInSeconds)
}

// ProcessL1DataTime observes the time for synchronize the trusted state on the histogram.
func ProcessL1DataTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(ProcessL1DataTimeName, execTimeInSeconds)
}

// GetTrustedBatchNumberTime observes the time for synchronize the trusted state on the histogram.
func GetTrustedBatchNumberTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(GetTrustedBatchNumberTimeName, execTimeInSeconds)
}

// GetTrustedBatchInfoTime observes the time for synchronize the trusted state on the histogram.
func GetTrustedBatchInfoTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(GetTrustedBatchInfoTimeName, execTimeInSeconds)
}

// ProcessTrustedBatchTime observes the time for synchronize the trusted state on the histogram.
func ProcessTrustedBatchTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(ProcessTrustedBatchTimeName, execTimeInSeconds)
}

// TrustedBatchCleanCounter increases the counter for the trusted batch clean
func TrustedBatchCleanCounter() {
	metrics.CounterInc(TrustedBatchCleanCounterName)
}