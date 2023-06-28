package state

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config is state config
type Config struct {
	// MaxCumulativeGasUsed is the max gas allowed per batch
	MaxCumulativeGasUsed uint64

	// ChainID is the L2 ChainID provided by the Network Config
	ChainID uint64

	// ForkIdIntervals is the list of fork id intervals
	ForkIDIntervals []ForkIDInterval

	// MaxResourceExhaustedAttempts is the max number of attempts to make a transaction succeed because of resource exhaustion
	MaxResourceExhaustedAttempts int

	// WaitOnResourceExhaustion is the time to wait before retrying a transaction because of resource exhaustion
	WaitOnResourceExhaustion types.Duration

	// Batch number where there is a forkid change (fork upgrade)
	ForkUpgradeBatchNumber uint64

	// New fork id to be used for batches greaters than ForkUpgradeBatchNumber (fork upgrade)
	ForkUpgradeNewForkId uint64
}
