package strategy

import (
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txprofitabilitychecker"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txselector"
)

// Strategy holds config params
type Strategy struct {
	TxSelector             txselector.Config             `mapstructure:"TxSelector"`
	TxProfitabilityChecker txprofitabilitychecker.Config `mapstructure:"TxProfitabilityChecker"`
}
