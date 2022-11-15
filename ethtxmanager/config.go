package ethtxmanager

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config is configuration for ethereum transaction manager
type Config struct {
	// IntervalToReviewSendBatchTx is the time limit we wait in order to review a send batch tx sent to l1
	IntervalToReviewSendBatchTx types.Duration `mapstructure:"IntervalToReviewSendBatchTx"`

	// IntervalToReviewVerifyBatchTx is the time limit we wait in order to review a verify batch tx sent to l1
	IntervalToReviewVerifyBatchTx types.Duration `mapstructure:"IntervalToReviewVerifyBatchTx"`
}
