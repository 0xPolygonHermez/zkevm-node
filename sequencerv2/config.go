package sequencerv2

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/sequencerv2/profitabilitychecker"
)

// Config represents the configuration of a sequencer
type Config struct {
	// WaitPeriodPoolIsEmpty is the time the sequencer waits until
	// trying to propose a batch
	WaitPeriodPoolIsEmpty types.Duration `mapstructure:"WaitPeriodPoolIsEmpty"`

	// LastBatchVirtualizationTimeMaxWaitPeriod is time since sequences should be sent
	LastBatchVirtualizationTimeMaxWaitPeriod types.Duration `mapstructure:"LastBatchVirtualizationTimeMaxWaitPeriod"`

	// WaitBlocksToUpdateGER is number of blocks for sequencer to wait
	WaitBlocksToUpdateGER uint64 `mapstructure:"WaitBlocksToUpdateGER"`

	// LastTimeBatchMaxWaitPeriod is time after which new batch should be closed
	LastTimeBatchMaxWaitPeriod types.Duration `mapstructure:"LastTimeBatchMaxWaitPeriod"`

	// BlocksAmountForTxsToBeDeleted is blocks amount after which txs will be deleted from the pool
	BlocksAmountForTxsToBeDeleted uint64 `mapstructure:"BlocksAmountForTxsToBeDeleted"`

	// FrequencyToCheckTxsForDelete is frequency with which txs will be checked for deleting
	FrequencyToCheckTxsForDelete types.Duration `mapstructure:"FrequencyToCheckTxsForDelete"`

	// ProfitabilityChecker configuration
	ProfitabilityChecker profitabilitychecker.Config `mapstructure:"ProfitabilityChecker"`
}
