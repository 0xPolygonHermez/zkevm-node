package executor

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config represents the configuration of the executor server
type Config struct {
	URI string `mapstructure:"URI"`
	// MaxResourceExhaustedAttempts is the max number of attempts to make a transaction succeed because of resource exhaustion
	MaxResourceExhaustedAttempts int `mapstructure:"MaxResourceExhaustedAttempts"`
	// WaitOnResourceExhaustion is the time to wait before retrying a transaction because of resource exhaustion
	WaitOnResourceExhaustion types.Duration `mapstructure:"WaitOnResourceExhaustion"`
	MaxGRPCMessageSize       int            `mapstructure:"MaxGRPCMessageSize"`
}
