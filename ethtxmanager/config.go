package ethtxmanager

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config is configuration for ethereum transaction manager
type Config struct {
	// MaxSendBatchTxRetries amount of how many tries for sending sendBatch tx to the ethereum
	MaxSendBatchTxRetries uint32 `mapstructure:"MaxSendBatchTxRetries"`
	// FrequencyForResendingFailedSendBatches frequency of the resending batches
	FrequencyForResendingFailedSendBatches types.Duration `mapstructure:"FrequencyForResendingFailedSendBatches"`

	// MaxVerifyBatchTxRetries amount of how many tries for sending verifyBatch tx to the ethereum
	MaxVerifyBatchTxRetries uint32 `mapstructure:"MaxVerifyBatchTxRetries"`
	// FrequencyForResendingFailedVerifyBatch frequency of the resending verify batch function
	FrequencyForResendingFailedVerifyBatch types.Duration `mapstructure:"FrequencyForResendingFailedVerifyBatch"`
	// WaitTxToBeMined time to wait after transaction was sent to the ethereum
	WaitTxToBeMined types.Duration `mapstructure:"WaitTxToBeMined"`
	// WaitTxToBeSynced time to wait after transaction was sent to the ethereum to get into the state
	WaitTxToBeSynced types.Duration `mapstructure:"WaitTxToBeSynced"`
	// PercentageToIncreaseGasPrice when tx is failed by timeout increase gas price by this percentage
	PercentageToIncreaseGasPrice uint64 `mapstructure:"PercentageToIncreaseGasPrice"`
	// PercentageToIncreaseGasLimit when tx is failed by timeout increase gas price by this percentage
	PercentageToIncreaseGasLimit uint64 `mapstructure:"PercentageToIncreaseGasLimit"`
}
