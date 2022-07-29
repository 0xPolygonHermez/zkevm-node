package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgx/v4"
)

// ZKEVM contains implementations for the "zkevm" RPC endpoints
type ZKEVM struct {
	config Config
	state  stateInterface
	txMan  dbTxManager
}

// ConsolidatedBlockNumber returns current block number for consolidated blocks
func (h *ZKEVM) ConsolidatedBlockNumber() (interface{}, rpcError) {
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

// IsBatchConsolidated returns the consolidation status of a provided batch ID
func (h *ZKEVM) IsBatchConsolidated(batchID int) (interface{}, rpcError) {
	return h.txMan.NewDbTxScope(h.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		isBatchConsolidated, err := h.state.IsBatchConsolidated(ctx, batchID, dbTx)
		if err != nil {
			const errorMessage = "failed to check if the batch is consolidated"
			log.Errorf("%v:%v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return isBatchConsolidated, nil
	})
}

// IsBatchVirtualized returns the virtualisation status of a provided batch ID
func (h *ZKEVM) IsBatchVirtualized(batchID int) (interface{}, rpcError) {
	return h.txMan.NewDbTxScope(h.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		isBatchVirtualized, err := h.state.IsBatchVirtualized(ctx, batchID, dbTx)
		if err != nil {
			const errorMessage = "failed to check if the batch is virtualized"
			log.Errorf("%v:%v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return isBatchVirtualized, nil
	})
}

// GetBroadcastURI returns the IP:PORT of the broadcast service provided
// by the Trusted Sequencer JSON RPC server
func (h *ZKEVM) GetBroadcastURI() (interface{}, rpcError) {
	return h.config.BroadcastURI, nil
}
