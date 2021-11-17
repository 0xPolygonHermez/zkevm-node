package pool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Pool interface {
	AddTx(tx types.Transaction) error
	GetPendingTxs() ([]Transaction, error)
	UpdateTxState(hash common.Hash, newState TxState) error
	CleanUpInvalidAndNonSelectedTxs() error
	EstimateGas() (uint64, error)
	GetGasPrice() (uint64, error)
}

func NewPool() Pool {
	panic("not implemented")
}
