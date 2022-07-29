package ethtxmanager

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config is configuration for ethereum transaction manager
type Config struct {
	// MaxSendBatchTxRetries amount of how many tries for sending sendBatch tx to the ethereum
	MaxSendBatchTxRetries uint32 `mapstructure:"MaxSendBatchTxRetries"`
	// FrequencyForResendingFailedSendBatchesInMilliseconds frequency of the resending batches
	FrequencyForResendingFailedSendBatches types.Duration `mapstructure:"FrequencyForResendingFailedSendBatches"`

	// MaxVerifyBatchTxRetries amount of how many tries for sending sendBatch tx to the ethereum
	MaxVerifyBatchTxRetries uint32 `mapstructure:"MaxVerifyBatchTxRetries"`
	// FrequencyForResendingFailedVerifyBatchInMilliseconds frequency of the resending batches
	FrequencyForResendingFailedVerifyBatch types.Duration `mapstructure:"FrequencyForResendingFailedVerifyBatch"`
}
