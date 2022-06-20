package sequencerv2

import (
	"github.com/hermeznetwork/hermez-core/config/types"
)

// Config represents the configuration of a sequencer
type Config struct {
	// WaitPeriodPoolIsEmpty is the time the sequencer waits until
	// trying to propose a batch
	WaitPeriodPoolIsEmpty types.Duration `mapstructure:"WaitPeriodPoolIsEmpty"`

	// LastL1InteractionTimeWaitPeriod is time since sequences should be sent
	LastL1InteractionTimeMaxWaitPeriod types.Duration `mapstructure:"LastL1InteractionTimeMaxWaitPeriod"`

	// LastTimeGERUpdatedMaxWaitPeriod is possible wait time since last GER was updated
	LastTimeGERUpdatedMaxWaitPeriod types.Duration `mapstructure:"LastTimeGERUpdatedMaxWaitPeriod"`

	// LastTimeDepositMaxWaitPeriod is possible wait time since last deposit happened
	LastTimeDepositMaxWaitPeriod types.Duration `mapstructure:"LastTimeDepositMaxWaitPeriod"`

	// LastTimeBatchMaxWaitPeriod is time after which new batch should be sent
	LastTimeBatchMaxWaitPeriod types.Duration `mapstructure:"LastTimeBatchMaxWaitPeriod"`
}
