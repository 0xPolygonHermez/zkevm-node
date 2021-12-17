package aggregator

import (
	"time"

	"github.com/hermeznetwork/hermez-core/etherman"
)

// Config represents the configuration of the aggregator
type Config struct {
	// IntervalToConsolidateState is the time the aggregator waits until
	// trying to consolidate a new state
	IntervalToConsolidateState time.Duration

	// Etherman is the configuration required by etherman to interact with L1
	Etherman etherman.Config

	TxProfitabilityCheckerType TxProfitabilityCheckerType

	MinReward uint64
}
