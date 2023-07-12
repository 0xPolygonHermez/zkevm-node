package pool

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// Config is the pool configuration
type Config struct {
	// IntervalToRefreshBlockedAddresses is the time it takes to sync the
	// blocked address list from db to memory
	IntervalToRefreshBlockedAddresses types.Duration `mapstructure:"IntervalToRefreshBlockedAddresses"`

	// IntervalToRefreshGasPrices is the time to wait to refresh the gas prices
	IntervalToRefreshGasPrices types.Duration `mapstructure:"IntervalToRefreshGasPrices"`

	// MaxTxBytesSize is the max size of a transaction in bytes
	MaxTxBytesSize uint64 `mapstructure:"MaxTxBytesSize"`

	// MaxTxDataBytesSize is the max size of the data field of a transaction in bytes
	MaxTxDataBytesSize int `mapstructure:"MaxTxDataBytesSize"`

	// DB is the database configuration
	DB db.Config `mapstructure:"DB"`

	// DefaultMinGasPriceAllowed is the default min gas price to suggest
	DefaultMinGasPriceAllowed uint64 `mapstructure:"DefaultMinGasPriceAllowed"`

	// MinAllowedGasPriceInterval is the interval to look back of the suggested min gas price for a tx
	MinAllowedGasPriceInterval types.Duration `mapstructure:"MinAllowedGasPriceInterval"`

	// PollMinAllowedGasPriceInterval is the interval to poll the suggested min gas price for a tx
	PollMinAllowedGasPriceInterval types.Duration `mapstructure:"PollMinAllowedGasPriceInterval"`

	// AccountQueue represents the maximum number of non-executable transaction slots permitted per account
	AccountQueue uint64 `mapstructure:"AccountQueue"`

	// GlobalQueue represents the maximum number of non-executable transaction slots for all accounts
	GlobalQueue uint64 `mapstructure:"GlobalQueue"`
}

// EffectiveGasPrice has parameters for the effective gas price calculation.
type EffectiveGasPrice struct {
	// L1GasPriceFactor is the percentage of the L1 gas price that will be used as the L2 min gas price
	L1GasPriceFactor float64 `mapstructure:"L1GasPriceFactor"`

	// ByteGasCost is the gas cost per byte
	ByteGasCost uint64 `mapstructure:"ByteGasCost"`

	// MarginFactor is the margin factor percentage to be added to the L2 min gas price
	MarginFactor float64 `mapstructure:"MarginFactor"`
}

// BatchConfig represents the configuration of the batch constraints
type BatchConfig struct {
	Constraints     BatchConstraintsCfg  `mapstructure:"Constraints"`
	ResourceWeights BatchResourceWeights `mapstructure:"ResourceWeights"`
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
func (c BatchConstraintsCfg) IsWithinConstraints(counters state.ZKCounters) bool {
	return counters.CumulativeGasUsed <= c.MaxCumulativeGasUsed &&
		counters.UsedKeccakHashes <= c.MaxKeccakHashes &&
		counters.UsedPoseidonHashes <= c.MaxPoseidonHashes &&
		counters.UsedPoseidonPaddings <= c.MaxPoseidonPaddings &&
		counters.UsedMemAligns <= c.MaxMemAligns &&
		counters.UsedArithmetics <= c.MaxArithmetics &&
		counters.UsedBinaries <= c.MaxBinaries &&
		counters.UsedSteps <= c.MaxSteps
}

// BatchResourceWeights represents the configuration of the batch resource weights
type BatchResourceWeights struct {
	WeightBatchBytesSize    int `mapstructure:"WeightBatchBytesSize"`
	WeightCumulativeGasUsed int `mapstructure:"WeightCumulativeGasUsed"`
	WeightKeccakHashes      int `mapstructure:"WeightKeccakHashes"`
	WeightPoseidonHashes    int `mapstructure:"WeightPoseidonHashes"`
	WeightPoseidonPaddings  int `mapstructure:"WeightPoseidonPaddings"`
	WeightMemAligns         int `mapstructure:"WeightMemAligns"`
	WeightArithmetics       int `mapstructure:"WeightArithmetics"`
	WeightBinaries          int `mapstructure:"WeightBinaries"`
	WeightSteps             int `mapstructure:"WeightSteps"`
}
