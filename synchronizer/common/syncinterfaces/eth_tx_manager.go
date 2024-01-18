package syncinterfaces

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type EthTxManager interface {
	Reorg(ctx context.Context, fromBlockNumber uint64, dbTx pgx.Tx) error
}
