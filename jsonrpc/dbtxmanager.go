package jsonrpc

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type dbTxManager struct{}

type dbTxScopedFn func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError)

func (f *dbTxManager) NewDbTxScope(st stateInterface, scopedFn dbTxScopedFn) (interface{}, rpcError) {
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to connect to the state", err)
	}

	v, rpcErr := scopedFn(ctx, dbTx)
	if rpcErr != nil {
		if txErr := dbTx.Rollback(context.Background()); txErr != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to rollback db transaction", txErr)
		}
		return v, rpcErr
	}

	if txErr := dbTx.Commit(context.Background()); txErr != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to commit db transaction", txErr)
	}
	return v, rpcErr
}
