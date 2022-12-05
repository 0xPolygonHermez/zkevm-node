package metrics

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prefix                         = "sequencer_"
	sequencesSentToL1CountName     = prefix + "sequences_sent_to_L1_count"
	gasPriceEstimatedAverageName   = prefix + "gas_price_estimated_average"
	txProcessed                    = prefix + "transaction_processed"
	sequencesOvesizedDataErrorName = prefix + "sequences_oversized_data_error"
	ethToMaticPriceName            = prefix + "eth_to_matic_price"
	sequenceRewardInMaticName      = prefix + "sequence_reward_in_matic"
	processingTime                 = prefix + "processing_time"

	txProcessedLabelName = "status"
)

// TxProcessedLabel represents the possible values for the
// `sequencer_transaction_processed` metric `type` label.
type TxProcessedLabel string

const (
	// TxProcessedLabelSuccessful represents a successful transaction
	TxProcessedLabelSuccessful TxProcessedLabel = "successful"
	// TxProcessedLabelInvalid represents an invalid transaction
	TxProcessedLabelInvalid TxProcessedLabel = "invalid"
	// TxProcessedLabelFailed represents a failed transaction
	TxProcessedLabelFailed TxProcessedLabel = "failed"
)

// Register the metrics for the sequencer package.
func Register() {
	var (
		counters    []prometheus.CounterOpts
		counterVecs []metrics.CounterVecOpts
		gauges      []prometheus.GaugeOpts
		histograms  []prometheus.HistogramOpts
	)

	counters = []prometheus.CounterOpts{
		{
			Name: sequencesSentToL1CountName,
			Help: "[SEQUENCER] total count of sequences sent to L1",
		},
		{
			Name: sequencesOvesizedDataErrorName,
			Help: "[SEQUENCER] total count of sequences with oversized data error",
		},
	}

	counterVecs = []metrics.CounterVecOpts{
		{
			CounterOpts: prometheus.CounterOpts{
				Name: txProcessed,
				Help: "[SEQUENCER] number of transactions processed",
			},
			Labels: []string{txProcessedLabelName},
		},
	}

	gauges = []prometheus.GaugeOpts{
		{
			Name: gasPriceEstimatedAverageName,
			Help: "[SEQUENCER] average gas price estimated",
		},
		{
			Name: ethToMaticPriceName,
			Help: "[SEQUENCER] eth to matic price",
		},
		{
			Name: sequenceRewardInMaticName,
			Help: "[SEQUENCER] reward for a sequence in Matic",
		},
	}

	histograms = []prometheus.HistogramOpts{
		{
			Name: processingTime,
			Help: "[SEQUENCER] processing time",
		},
	}

	metrics.RegisterCounters(counters...)
	metrics.RegisterCounterVecs(counterVecs...)
	metrics.RegisterGauges(gauges...)
	metrics.RegisterHistograms(histograms...)
}

// AverageGasPrice sets the gauge to the given average gas price.
func AverageGasPrice(price float64) {
	metrics.GaugeSet(gasPriceEstimatedAverageName, price)
}

// SequencesSentToL1 increases the counter by the provided number of sequences
// sent to L1.
func SequencesSentToL1(numSequences float64) {
	metrics.CounterAdd(sequencesSentToL1CountName, numSequences)
}

// TxProcessed increases the counter vector by the provided transactions count
// and for the given label.
func TxProcessed(status TxProcessedLabel, count float64) {
	metrics.CounterVecAdd(txProcessed, string(status), count)
}

// SequencesOvesizedDataError increases the counter for sequences that
// encounter a OversizedData error.
func SequencesOvesizedDataError() {
	metrics.CounterInc(sequencesOvesizedDataErrorName)
}

// EthToMaticPrice sets the gauge for the Ethereum to Matic price.
func EthToMaticPrice(price float64) {
	metrics.GaugeSet(ethToMaticPriceName, price)
}

// SequenceRewardInMatic sets the gauge for the reward in Matic of a sequence.
func SequenceRewardInMatic(reward float64) {
	metrics.GaugeSet(sequenceRewardInMaticName, reward)
}

// ProcessingTime observes the last iteration processing time on the histogram.
func ProcessingTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(processingTime, execTimeInSeconds)
}
