package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/jackc/pgx/v4"
)

type DBTxManager struct{}

type DBTxScopedFn func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error)

type DBTxer interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
}

func (f *DBTxManager) NewDbTxScope(db DBTxer, scopedFn DBTxScopedFn) (interface{}, types.Error) {
	ctx := context.Background()
	dbTx, err := db.BeginStateTransaction(ctx)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to connect to the state", err)
	}

	v, rpcErr := scopedFn(ctx, dbTx)
	if rpcErr != nil {
		if txErr := dbTx.Rollback(context.Background()); txErr != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to rollback db transaction", txErr)
		}
		return v, rpcErr
	}

	if txErr := dbTx.Commit(context.Background()); txErr != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to commit db transaction", txErr)
	}
	return v, rpcErr
}
