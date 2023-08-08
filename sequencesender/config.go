package sequencesender

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// Config represents the configuration of a sequence sender
type Config struct {
	// WaitPeriodSendSequence is the time the sequencer waits until
	// trying to send a sequence to L1
	WaitPeriodSendSequence types.Duration `mapstructure:"WaitPeriodSendSequence"`
	// LastBatchVirtualizationTimeMaxWaitPeriod is time since sequences should be sent
	LastBatchVirtualizationTimeMaxWaitPeriod types.Duration `mapstructure:"LastBatchVirtualizationTimeMaxWaitPeriod"`
	// MaxTxSizeForL1 is the maximum size a single transaction can have. This field has
	// non-trivial consequences: larger transactions than 128KB are significantly harder and
	// more expensive to propagate; larger transactions also take more resources
	// to validate whether they fit into the pool or not.
	MaxTxSizeForL1 uint64 `mapstructure:"MaxTxSizeForL1"`
	// SenderAddress defines which private key the eth tx manager needs to use
	// to sign the L1 txs
	SenderAddress string `mapstructure:"SenderAddress"`
	// PrivateKeys defines all the key store files that are going
	// to be read in order to provide the private keys to sign the L1 txs
	PrivateKeys []types.KeystoreFileConfig `mapstructure:"PrivateKeys"`
	// Batch number where there is a forkid change (fork upgrade)
	ForkUpgradeBatchNumber uint64
}
