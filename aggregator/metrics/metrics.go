package metrics

import (
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prefix                      = "aggregator_"
	currentConnectedProversName = prefix + "current_connected_provers"
	currentWorkingProversName   = prefix + "current_working_provers"
)

// Register the metrics for the sequencer package.
func Register() {
	gauges := []prometheus.GaugeOpts{
		{
			Name: currentConnectedProversName,
			Help: "[AGGREGATOR] current connected provers",
		},
		{
			Name: currentWorkingProversName,
			Help: "[AGGREGATOR] current working provers",
		},
	}

	metrics.RegisterGauges(gauges...)
}

// ConnectedProver increments the gauge for the current number of connected
// provers.
func ConnectedProver() {
	metrics.GaugeInc(currentConnectedProversName)
}

// DisconnectedProver decrements the gauge for the current number of connected
// provers.
func DisconnectedProver() {
	metrics.GaugeDec(currentConnectedProversName)
}

// WorkingProver increments the gauge for the current number of working
// provers.
func WorkingProver() {
	metrics.GaugeInc(currentWorkingProversName)
}

// IdlingProver decrements the gauge for the current number of working provers.
func IdlingProver() {
	metrics.GaugeDec(currentWorkingProversName)
}
