package jsonrpcv2

import (
	"context"

	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/jackc/pgx/v4"
)

// Hez contains implementations for the "hez" RPC endpoints
type Hez struct {
	state stateInterface
	txMan dbTxManager
}

// ConsolidatedBlockNumber returns current block number for consolidated batches
func (h *Hez) ConsolidatedBlockNumber() (interface{}, error) {
	return h.txMan.NewDbTxScope(h.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		lastBatchNumber, err := h.state.GetLastConsolidatedBlockNumber(ctx, dbTx)
		if err != nil {
			return nil, err
		}

		return hex.EncodeUint64(lastBatchNumber), nil
	})
}
