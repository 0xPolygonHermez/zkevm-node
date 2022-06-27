package jsonrpcv2

import (
	"context"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/jackc/pgx/v4"
)

type dbTxManager struct{}

func (f *dbTxManager) NewDbTxScope(st stateInterface, scopedFn func(ctx context.Context, dbTx pgx.Tx) (interface{}, error)) (interface{}, error) {
	ctx := context.Background()
	dbTx, err := st.BeginStateTransaction(ctx)
	if err != nil {
		return nil, newRPCError(defaultErrorCode, "failed to connect to the state")
	}

	v, err := scopedFn(ctx, dbTx)
	if err != nil {
		if txErr := dbTx.Rollback(context.Background()); err != nil {
			log.Errorf("failed to roll back tx: %v", txErr)
		}
		return v, err
	}

	if txErr := dbTx.Commit(context.Background()); err != nil {
		log.Errorf("failed to commit tx: %v", txErr)
	}
	return v, err
}
