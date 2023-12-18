package state

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
)

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

	// Batch number from which there is a forkid change (fork upgrade)
	ForkUpgradeBatchNumber uint64

	// New fork id to be used for batches greaters than ForkUpgradeBatchNumber (fork upgrade)
	ForkUpgradeNewForkId uint64

	// DB is the database configuration
	DB db.Config `mapstructure:"DB"`

	// Configuration for the batch constraints
	Batch BatchConfig `mapstructure:"Batch"`

	// MaxLogsCount is a configuration to set the max number of logs that can be returned
	// in a single call to the state, if zero it means no limit
	MaxLogsCount uint64

	// MaxLogsBlockRange is a configuration to set the max range for block number when querying TXs
	// logs in a single call to the state, if zero it means no limit
	MaxLogsBlockRange uint64

	// MaxNativeBlockHashBlockRange is a configuration to set the max range for block number when querying
	// native block hashes in a single call to the state, if zero it means no limit
	MaxNativeBlockHashBlockRange uint64
}

// BatchConfig represents the configuration of the batch constraints
type BatchConfig struct {
	Constraints BatchConstraintsCfg `mapstructure:"Constraints"`
}

// BatchConstraintsCfg represents the configuration of the batch constraints
type BatchConstraintsCfg struct {
	MaxTxsPerBatch       uint64 `mapstructure:"MaxTxsPerBatch"`
	MaxBatchBytesSize    uint64 `mapstructure:"MaxBatchBytesSize"`
	MaxCumulativeGasUsed uint64 `mapstructure:"MaxCumulativeGasUsed"`
	MaxKeccakHashes      uint32 `mapstructure:"MaxKeccakHashes"`
	MaxPoseidonHashes    uint32 `mapstructure:"MaxPoseidonHashes"`
	MaxPoseidonPaddings  uint32 `mapstructure:"MaxPoseidonPaddings"`
	MaxMemAligns         uint32 `mapstructure:"MaxMemAligns"`
	MaxArithmetics       uint32 `mapstructure:"MaxArithmetics"`
	MaxBinaries          uint32 `mapstructure:"MaxBinaries"`
	MaxSteps             uint32 `mapstructure:"MaxSteps"`
}

// IsWithinConstraints checks if the counters are within the batch constraints
func (c BatchConstraintsCfg) IsWithinConstraints(counters ZKCounters) bool {
	return counters.CumulativeGasUsed <= c.MaxCumulativeGasUsed &&
		counters.UsedKeccakHashes <= c.MaxKeccakHashes &&
		counters.UsedPoseidonHashes <= c.MaxPoseidonHashes &&
		counters.UsedPoseidonPaddings <= c.MaxPoseidonPaddings &&
		counters.UsedMemAligns <= c.MaxMemAligns &&
		counters.UsedArithmetics <= c.MaxArithmetics &&
		counters.UsedBinaries <= c.MaxBinaries &&
		counters.UsedSteps <= c.MaxSteps
}
