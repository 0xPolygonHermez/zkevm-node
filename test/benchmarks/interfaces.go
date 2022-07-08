package benchmarks

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
}
