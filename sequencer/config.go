package sequencer

import (
	"github.com/hermeznetwork/hermez-core/config/types"
	"github.com/hermeznetwork/hermez-core/pricegetter"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy"
)

// InitBatchProcessorIfDiffType let sequencer decide, how to init batch processor
type InitBatchProcessorIfDiffType string

const (
	// InitBatchProcessorIfDiffTypeSynced init batch processor from previous synced batch root
	InitBatchProcessorIfDiffTypeSynced InitBatchProcessorIfDiffType = "synced"
	// InitBatchProcessorIfDiffTypeCalculated init batch processor from previous calculated batch root
	InitBatchProcessorIfDiffTypeCalculated InitBatchProcessorIfDiffType = "calculated"
)

// Config represents the configuration of a sequencer
type Config struct {
	// IntervalToProposeBatch is the time the sequencer waits until
	// trying to propose a batch
	IntervalToProposeBatch types.Duration `mapstructure:"IntervalToProposeBatch"`

	// SyncedBlockDif is the difference, how many block left to sync. So if sequencer see, that
	// X amount of blocks are left to sync, it will start to select txs
	SyncedBlockDif uint64 `mapstructure:"SyncedBlockDif"`

	// IntervalAfterWhichBatchSentAnyway this is interval for the main sequencer, that will check if there is no transactions
	IntervalAfterWhichBatchSentAnyway types.Duration `mapstructure:"IntervalAfterWhichBatchSentAnyway"`

	// Strategy is the configuration for the strategy
	Strategy strategy.Strategy `mapstructure:"Strategy"`

	// PriceGetter config for the price getter
	PriceGetter pricegetter.Config `mapstructure:"PriceGetter"`

	// InitBatchProcessorIfDiffType is for the case, when last synchronized batch num more than latest sent batch
	// If "synced" init bp by synced batch, if "calculated" init by previous calculated root
	InitBatchProcessorIfDiffType InitBatchProcessorIfDiffType `mapstructure:"InitBatchProcessorIfDiffType"`
	// AllowNonRegistered determines if the sequencer will run using the default
	// chain ID
	AllowNonRegistered bool `mapstructure:"AllowNonRegistered"`

	// MaxSendBatchTxRetries amount of how many tries for sending sendBatch tx to the ethereum
	MaxSendBatchTxRetries uint32 `mapstructure:"MaxSendBatchTxRetries"`
	// FrequencyForResendingFailedSendBatchesInMilliseconds frequency of the resending batches
	FrequencyForResendingFailedSendBatchesInMilliseconds int64 `mapstructure:"FrequencyForResendingFailedSendBatchesInMilliseconds"`
	// DefaultChainID is the common ChainID to all the sequencers
	DefaultChainID uint64 `mapstructure:"DefaultChainID"`

	// PendingTxsQueue config for pending tx queue data structure
	PendingTxsQueue PendingTxsQueueConfig `mapstructure:"PendingTxsQueue"`
}
