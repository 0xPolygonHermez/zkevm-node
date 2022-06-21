package statev2

import "github.com/hermeznetwork/hermez-core/statev2/runtime/executor"

// Config is state config
type Config struct {
	// MaxCumulativeGasUsed is the max gas allowed per batch
	MaxCumulativeGasUsed uint64
	// Executor configuration
	ExecutorServerConfig executor.Config
}
