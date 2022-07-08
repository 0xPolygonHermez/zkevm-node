package strategy

import (
	"github.com/0xPolygonHermez/zkevm-node/sequencer/strategy/txprofitabilitychecker"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/strategy/txselector"
)

// Strategy holds config params
type Strategy struct {
	TxSelector             txselector.Config             `mapstructure:"TxSelector"`
	TxProfitabilityChecker txprofitabilitychecker.Config `mapstructure:"TxProfitabilityChecker"`
}
