package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/jackc/pgx/v4"
)

type dbTxManager struct{}

type dbTxScopedFn func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error)

func (f *dbTxManager) NewDbTxScope(st types.StateInterface, scopedFn dbTxScopedFn) (interface{}, types.Error) {
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	if err != nil {
		return rpcErrorResponse(types.DefaultErrorCode, "failed to connect to the state", err)
	}

	v, rpcErr := scopedFn(ctx, dbTx)
	if rpcErr != nil {
		if txErr := dbTx.Rollback(context.Background()); txErr != nil {
			return rpcErrorResponse(types.DefaultErrorCode, "failed to rollback db transaction", txErr)
		}
		return v, rpcErr
	}

	if txErr := dbTx.Commit(context.Background()); txErr != nil {
		return rpcErrorResponse(types.DefaultErrorCode, "failed to commit db transaction", txErr)
	}
	return v, rpcErr
}
