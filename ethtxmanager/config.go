package ethtxmanager

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config is configuration for ethereum transaction manager
type Config struct {
	// FrequencyForResendingFailedSendBatches frequency of the resending batches
	FrequencyForResendingFailedSendBatches types.Duration `mapstructure:"FrequencyForResendingFailedSendBatches"`
	// FrequencyForResendingFailedVerifyBatch frequency of the resending verify batch function
	FrequencyForResendingFailedVerifyBatch types.Duration `mapstructure:"FrequencyForResendingFailedVerifyBatch"`

	// WaitTxToBeMined time to wait after transaction was sent to the ethereum
	WaitTxToBeMined types.Duration `mapstructure:"WaitTxToBeMined"`
	// WaitTxToBeSynced time to wait after transaction was sent to the ethereum to get into the state
	WaitTxToBeSynced types.Duration `mapstructure:"WaitTxToBeSynced"`
}
