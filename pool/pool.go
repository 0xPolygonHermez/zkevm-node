package pool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/db"
)

// Pool represents a pool of transactions
type Pool interface {
	AddTx(tx types.Transaction) error
	GetPendingTxs() ([]Transaction, error)
	UpdateTxState(hash common.Hash, newState TxState) error
	CleanUpInvalidAndNonSelectedTxs() error
	SetGasPrice(gasPrice uint64) error
	GetGasPrice() (uint64, error)
}

// NewPool creates a new pool
func NewPool(cfg db.Config) (Pool, error) {
	return newPostgresPool(cfg)
}
