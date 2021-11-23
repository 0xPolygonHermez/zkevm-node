package pool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Pool represents a pool of transactions
type Pool interface {
	AddTx(tx types.Transaction) error
	GetPendingTxs() ([]Transaction, error)
	UpdateTxState(hash common.Hash, newState TxState) error
	CleanUpInvalidAndNonSelectedTxs() error
	GetGasPrice() (uint64, error)
}

// NewPool creates a new pool
func NewPool() Pool {
	panic("not implemented")
}
