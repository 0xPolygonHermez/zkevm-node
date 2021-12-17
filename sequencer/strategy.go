package sequencer

import "time"

type StrategyType string

const (
	AcceptAll StrategyType = "acceptall"
	Base                   = "base"
)

type TxSorterType string

const (
	ByCostAndTime  TxSorterType = "bycostandtime"
	ByCostAndNonce              = "bycostandnonce"
)

type TxProfitabilityCheckerType string

const (
	BaseProfitability = "base"
)

type Strategy struct {
	StrategyType StrategyType
	TxSorterType
	TxProfitabilityCheckerType
	MinReward            uint64
	PossibleTimeToSendTx time.Duration
}
