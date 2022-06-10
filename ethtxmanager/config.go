package ethtxmanager

// Config is configuration for ethereum transaction manager
type Config struct {
	// MaxSendBatchTxRetries amount of how many tries for sending sendBatch tx to the ethereum
	MaxSendTxRetries uint32 `mapstructure:"MaxSendBatchTxRetries"`
	// FrequencyForResendingFailedSendBatchesInMilliseconds frequency of the resending batches
	FrequencyForResendingFailedTxs int64 `mapstructure:"FrequencyForResendingFailedSendBatchesInMilliseconds"`
}
