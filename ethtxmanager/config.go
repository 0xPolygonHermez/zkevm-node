package ethtxmanager

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config is configuration for ethereum transaction manager
type Config struct {
	// FrequencyForResendingFailedTxs frequency of the resending failed txs
	FrequencyForResendingFailedTxs types.Duration `mapstructure:"FrequencyForResendingFailedTxs"`
	// WaitTxToBeMined time to wait after transaction was sent to the ethereum
	WaitTxToBeMined types.Duration `mapstructure:"WaitTxToBeMined"`
}
