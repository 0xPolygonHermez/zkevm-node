package params

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
	// NumberOfTxs is the number of transactions to send
	NumberOfTxs = 200
)
