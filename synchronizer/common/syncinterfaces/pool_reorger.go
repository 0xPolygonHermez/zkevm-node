package syncinterfaces

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type PoolReorger interface {
	ReorgPool(ctx context.Context, dbTx pgx.Tx) error
}
