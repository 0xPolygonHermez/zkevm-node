package sequencesender

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/ethereum/go-ethereum/common"
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

	// MaxBatchesForL1 is the maximum amount of batches to be sequenced in a single L1 tx
	MaxBatchesForL1 uint64 `mapstructure:"MaxBatchesForL1"`

	// SenderAddress defines which private key the eth tx manager needs to use
	// to sign the L1 txs
	SenderAddress common.Address
	// L2Coinbase defines which addess is going to receive the fees
	L2Coinbase common.Address `mapstructure:"L2Coinbase"`
	// PrivateKey defines all the key store files that are going
	// to be read in order to provide the private keys to sign the L1 txs
	PrivateKey types.KeystoreFileConfig `mapstructure:"PrivateKey"`
	// Batch number where there is a forkid change (fork upgrade)
	ForkUpgradeBatchNumber uint64

	// UseValidium is a flag to enable/disable the use of validium
	UseValidium bool `mapstructure:"UseValidium"`
}
