package shared

import (
	"time"
)

const (
	// DefaultDeadline is the default deadline for the sequencer
	DefaultDeadline = 6000 * time.Second
	// MaxCumulativeGasUsed is the maximum cumulative gas used
	MaxCumulativeGasUsed = 80000000000
	// PrometheusPort is the port where prometheus is running
	PrometheusPort = 9092
)
