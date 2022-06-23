package jsonrpcv2

import (
	"context"

	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/jackc/pgx/v4"
)

// Hez contains implementations for the "hez" RPC endpoints
type Hez struct {
	defaultChainID uint64
	state          stateInterface
	txMan          txManager
}

// DefaultChainId returns the default chain id that is allowed to be used by all the sequencers
func (h *Hez) DefaultChainId() (interface{}, error) { //nolint:revive
	return hex.EncodeUint64(h.defaultChainID), nil
}

// ConsolidatedBlockNumber returns current block number for consolidated batches
func (h *Hez) ConsolidatedBlockNumber() (interface{}, error) {
	return h.txMan.NewTxScope(h.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		lastBatchNumber, err := h.state.GetLastConsolidatedBatchNumber(ctx, dbTx)
		if err != nil {
			return nil, err
		}

		return hex.EncodeUint64(lastBatchNumber), nil
	})
}
