package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgx/v4"
)

// Hez contains implementations for the "hez" RPC endpoints
type Hez struct {
	state stateInterface
	txMan dbTxManager
}

// ConsolidatedBlockNumber returns current block number for consolidated blocks
func (h *Hez) ConsolidatedBlockNumber() (interface{}, rpcError) {
	return h.txMan.NewDbTxScope(h.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		lastBlockNumber, err := h.state.GetLastConsolidatedL2BlockNumber(ctx, dbTx)
		if err != nil {
			const errorMessage = "failed to get last consolidated block number from state"
			log.Errorf("%v:%v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return hex.EncodeUint64(lastBlockNumber), nil
	})
}
