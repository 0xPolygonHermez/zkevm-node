package pool

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	//nolint:errcheck // driver for postgres DB
	_ "github.com/lib/pq"
)

type Pool struct {
	db     *sqlx.DB
	ttl    time.Duration
	maxTxs uint32
}

func NewPool(
	db *sqlx.DB,
	ttl time.Duration,
	maxTxs uint32,
) *Pool {
	return &Pool{
		db:     db,
		ttl:    ttl,
		maxTxs: maxTxs,
	}
}

func (pool *Pool) GetPendingTxs() ([]Transaction, error) {
	panic("not implemented")
}

func (pool *Pool) UpdateTxState(hash common.Hash, newState TxState) error {
	panic("not implemented")
}

func (pool *Pool) CleanUpInvalidAndNonSelectedTx() error {
	panic("not implemented")
}
