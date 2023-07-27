package sequencer

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// Config represents the configuration of a sequencer
type Config struct {
	// WaitPeriodPoolIsEmpty is the time the sequencer waits until
	// trying to add new txs to the state
	WaitPeriodPoolIsEmpty types.Duration `mapstructure:"WaitPeriodPoolIsEmpty"`

	// BlocksAmountForTxsToBeDeleted is blocks amount after which txs will be deleted from the pool
	BlocksAmountForTxsToBeDeleted uint64 `mapstructure:"BlocksAmountForTxsToBeDeleted"`

	// FrequencyToCheckTxsForDelete is frequency with which txs will be checked for deleting
	FrequencyToCheckTxsForDelete types.Duration `mapstructure:"FrequencyToCheckTxsForDelete"`

	// TxLifetimeCheckTimeout is the time the sequencer waits to check txs lifetime
	TxLifetimeCheckTimeout types.Duration `mapstructure:"TxLifetimeCheckTimeout"`

	// MaxTxLifetime is the time a tx can be in the sequencer/worker memory
	MaxTxLifetime types.Duration `mapstructure:"MaxTxLifetime"`

	// Finalizer's specific config properties
	Finalizer FinalizerCfg `mapstructure:"Finalizer"`

	// DBManager's specific config properties
	DBManager DBManagerCfg `mapstructure:"DBManager"`

	// EffectiveGasPrice is the config for the gas price
	EffectiveGasPrice EffectiveGasPriceCfg `mapstructure:"EffectiveGasPrice"`

	// Http is the config for the http server
	Http HttpConfig `mapstructure:"Http"`
}

// FinalizerCfg contains the finalizer's configuration properties
type FinalizerCfg struct {
	// GERDeadlineTimeout is the time the finalizer waits after receiving closing signal to update Global Exit Root
	GERDeadlineTimeout types.Duration `mapstructure:"GERDeadlineTimeout"`

	// ForcedBatchDeadlineTimeout is the time the finalizer waits after receiving closing signal to process Forced Batches
	ForcedBatchDeadlineTimeout types.Duration `mapstructure:"ForcedBatchDeadlineTimeout"`

	// SleepDuration is the time the finalizer sleeps between each iteration, if there are no transactions to be processed
	SleepDuration types.Duration `mapstructure:"SleepDuration"`

	// ResourcePercentageToCloseBatch is the percentage window of the resource left out for the batch to be closed
	ResourcePercentageToCloseBatch uint32 `mapstructure:"ResourcePercentageToCloseBatch"`

	// GERFinalityNumberOfBlocks is number of blocks to consider GER final
	GERFinalityNumberOfBlocks uint64 `mapstructure:"GERFinalityNumberOfBlocks"`

	// ClosingSignalsManagerWaitForCheckingL1Timeout is used by the closing signals manager to wait for its operation
	ClosingSignalsManagerWaitForCheckingL1Timeout types.Duration `mapstructure:"ClosingSignalsManagerWaitForCheckingL1Timeout"`

	// ClosingSignalsManagerWaitForCheckingGER is used by the closing signals manager to wait for its operation
	ClosingSignalsManagerWaitForCheckingGER types.Duration `mapstructure:"ClosingSignalsManagerWaitForCheckingGER"`

	// ClosingSignalsManagerWaitForCheckingL1Timeout is used by the closing signals manager to wait for its operation
	ClosingSignalsManagerWaitForCheckingForcedBatches types.Duration `mapstructure:"ClosingSignalsManagerWaitForCheckingForcedBatches"`

	// ForcedBatchesFinalityNumberOfBlocks is number of blocks to consider GER final
	ForcedBatchesFinalityNumberOfBlocks uint64 `mapstructure:"ForcedBatchesFinalityNumberOfBlocks"`

	// TimestampResolution is the resolution of the timestamp used to close a batch
	TimestampResolution types.Duration `mapstructure:"TimestampResolution"`

	// SequentialReprocessFullBatch indicates if the reprocess of a closed batch (sanity check) must be done in a
	// sequential way (instead than in parallel)
	SequentialReprocessFullBatch bool `mapstructure:"SequentialReprocessFullBatch"`
}

// DBManagerCfg contains the DBManager's configuration properties
type DBManagerCfg struct {
	PoolRetrievalInterval    types.Duration `mapstructure:"PoolRetrievalInterval"`
	L2ReorgRetrievalInterval types.Duration `mapstructure:"L2ReorgRetrievalInterval"`
}

// EffectiveGasPriceCfg contains the configuration properties for the effective gas price
type EffectiveGasPriceCfg struct {
	// MaxBreakEvenGasPriceDeviationPercentage is the max allowed deviation percentage BreakEvenGasPrice on re-calculation
	MaxBreakEvenGasPriceDeviationPercentage uint64 `mapstructure:"MaxBreakEvenGasPriceDeviationPercentage"`

	// L1GasPriceFactor is the percentage of the L1 gas price that will be used as the L2 min gas price
	L1GasPriceFactor float64 `mapstructure:"L1GasPriceFactor"`

	// ByteGasCost is the gas cost per byte
	ByteGasCost uint64 `mapstructure:"ByteGasCost"`

	// MarginFactor is the margin factor percentage to be added to the L2 min gas price
	MarginFactor float64 `mapstructure:"MarginFactor"`

	// Enabled is a flag to enable/disable the effective gas price
	Enabled bool `mapstructure:"Enabled"`

	// DefaultMinGasPriceAllowed is the default min gas price to suggest
	// This value is assigned from [Pool].DefaultMinGasPriceAllowed
	DefaultMinGasPriceAllowed uint64
}

// HttpConfig contains the http server's configuration properties
type HttpConfig struct {
	// Host defines the network adapter that will be used to serve the HTTP requests
	Host string `mapstructure:"Host"`

	// Port defines the port to serve the endpoints via HTTP. For security reasons,
	// this port should not be directly accessible from outside the localhost. If
	// remote access is necessary, ensure to establish a secure SSH tunnel. If the
	// port needs to be exposed publicly for any reason, it is critical to protect
	// it with stringent firewall rules, limiting access only to trusted sources.
	Port int `mapstructure:"Port"`

	// MaxRequestsPerIPAndSecond defines how much requests a single IP can
	// send within a single second
	MaxRequestsPerIPAndSecond float64 `mapstructure:"MaxRequestsPerIPAndSecond"`

	// ReadTimeout is the HTTP server read timeout
	// check net/http.server.ReadTimeout and net/http.server.ReadHeaderTimeout
	ReadTimeout types.Duration `mapstructure:"ReadTimeout"`

	// WriteTimeout is the HTTP server write timeout
	// check net/http.server.WriteTimeout
	WriteTimeout types.Duration `mapstructure:"WriteTimeout"`
}
