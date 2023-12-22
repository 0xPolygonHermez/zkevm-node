package sequencer

import (
	"github.com/0xPolygonHermez/zkevm-data-streamer/log"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// Config represents the configuration of a sequencer
type Config struct {
	// BlocksAmountForTxsToBeDeleted is blocks amount after which txs will be deleted from the pool
	BlocksAmountForTxsToBeDeleted uint64 `mapstructure:"BlocksAmountForTxsToBeDeleted"`

	// FrequencyToCheckTxsForDelete is frequency with which txs will be checked for deleting
	FrequencyToCheckTxsForDelete types.Duration `mapstructure:"FrequencyToCheckTxsForDelete"`

	// TxLifetimeCheckTimeout is the time the sequencer waits to check txs lifetime
	TxLifetimeCheckTimeout types.Duration `mapstructure:"TxLifetimeCheckTimeout"`

	// MaxTxLifetime is the time a tx can be in the sequencer/worker memory
	MaxTxLifetime types.Duration `mapstructure:"MaxTxLifetime"`

	// PoolRetrievalInteral is the time the sequencer waits to check in there are new txs in the pool
	PoolRetrievalInterval types.Duration `mapstructure:"PoolRetrievalInterval"`

	// L2ReorgRetrievalInterval is the time the sequencer waits to check if a state inconsistency has happened
	L2ReorgRetrievalInterval types.Duration `mapstructure:"L2ReorgRetrievalInterval"`

	// Finalizer's specific config properties
	Finalizer FinalizerCfg `mapstructure:"Finalizer"`

	// StreamServerCfg is the config for the stream server
	StreamServer StreamServerCfg `mapstructure:"StreamServer"`
}

// StreamServerCfg contains the data streamer's configuration properties
type StreamServerCfg struct {
	// Port to listen on
	Port uint16 `mapstructure:"Port"`
	// Filename of the binary data file
	Filename string `mapstructure:"Filename"`
	// Enabled is a flag to enable/disable the data streamer
	Enabled bool `mapstructure:"Enabled"`
	// Log is the log configuration
	Log log.Config `mapstructure:"Log"`
}

// FinalizerCfg contains the finalizer's configuration properties
type FinalizerCfg struct {
	// ForcedBatchDeadlineTimeout is the time the finalizer waits after receiving closing signal to process Forced Batches
	ForcedBatchDeadlineTimeout types.Duration `mapstructure:"ForcedBatchDeadlineTimeout"`

	// SleepDuration is the time the finalizer sleeps between each iteration, if there are no transactions to be processed
	SleepDuration types.Duration `mapstructure:"SleepDuration"`

	// ResourcePercentageToCloseBatch is the percentage window of the resource left out for the batch to be closed
	ResourcePercentageToCloseBatch uint32 `mapstructure:"ResourcePercentageToCloseBatch"`

	// ForcedBatchesFinalityNumberOfBlocks is number of blocks to consider GER final
	ForcedBatchesFinalityNumberOfBlocks uint64 `mapstructure:"ForcedBatchesFinalityNumberOfBlocks"`

	// L1InfoRootFinalityNumberOfBlocks is number of blocks to consider L1InfoRoot final
	L1InfoRootFinalityNumberOfBlocks uint64 `mapstructure:"L1InfoRootFinalityNumberOfBlocks"`

	// ClosingSignalsManagerWaitForCheckingL1Timeout is used by the closing signals manager to wait for its operation
	ClosingSignalsManagerWaitForCheckingForcedBatches types.Duration `mapstructure:"ClosingSignalsManagerWaitForCheckingForcedBatches"`

	// WaitForCheckingL1InfoRoot is the wait time to check if the L1InfoRoot has been updated
	WaitForCheckingL1InfoRoot types.Duration `mapstructure:"WaitForCheckingL1InfoRoot"`

	// TimestampResolution is the resolution of the timestamp used to close a batch
	TimestampResolution types.Duration `mapstructure:"TimestampResolution"`

	// L2BlockTime is the resolution of the timestamp used to close a L2 block
	L2BlockTime types.Duration `mapstructure:"L2BlockTime"`

	// StopSequencerOnBatchNum specifies the batch number where the Sequencer will stop to process more transactions and generate new batches. The Sequencer will halt after it closes the batch equal to this number
	StopSequencerOnBatchNum uint64 `mapstructure:"StopSequencerOnBatchNum"`

	// SequentialReprocessFullBatch indicates if the reprocess of a closed batch (sanity check) must be done in a
	// sequential way (instead than in parallel)
	SequentialReprocessFullBatch bool `mapstructure:"SequentialReprocessFullBatch"`
}
