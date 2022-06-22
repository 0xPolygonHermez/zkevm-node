package sequencerv2

import (
	"github.com/hermeznetwork/hermez-core/config/types"
	"github.com/hermeznetwork/hermez-core/sequencerv2/profitabilitychecker"
)

// Config represents the configuration of a sequencer
type Config struct {
	// WaitPeriodPoolIsEmpty is the time the sequencer waits until
	// trying to propose a batch
	WaitPeriodPoolIsEmpty types.Duration `mapstructure:"WaitPeriodPoolIsEmpty"`

	// LastL1InteractionTimeWaitPeriod is time since sequences should be sent
	LastL1InteractionTimeMaxWaitPeriod types.Duration `mapstructure:"LastL1InteractionTimeMaxWaitPeriod"`

	// WaitBlocksToUpdateGER is number of blocks for sequencer to wait
	WaitBlocksToUpdateGER uint32 `mapstructure:"WaitBlocksToUpdateGER"`

	// LastTimeBatchMaxWaitPeriod is time after which new batch should be closed
	LastTimeBatchMaxWaitPeriod types.Duration `mapstructure:"LastTimeBatchMaxWaitPeriod"`

	// ProfitabilityChecker configuration
	ProfitabilityChecker profitabilitychecker.Config `mapstructure:"ProfitabilityChecker"`
}
