package pool

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type TxState string

const (
	TxStatePending  TxState = "pending"
	TxStateInvalid  TxState = "invalid"
	TxStateSelected TxState = "selected"
)

type Transaction struct {
	types.LegacyTx
	state TxState
}
