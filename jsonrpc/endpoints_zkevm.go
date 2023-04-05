package jsonrpc

import (
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/context"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// ZKEVMEndpoints contains implementations for the "zkevm" RPC endpoints
type ZKEVMEndpoints struct {
	config Config
	state  types.StateInterface
	txMan  dbTxManager
}

// ConsolidatedBlockNumber returns current block number for consolidated blocks
func (e *ZKEVMEndpoints) ConsolidatedBlockNumber(ctx *context.RequestContext) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(ctx, e.state, func(ctx *context.RequestContext, dbTx pgx.Tx) (interface{}, types.Error) {
		lastBlockNumber, err := e.state.GetLastConsolidatedL2BlockNumber(ctx, dbTx)
		if err != nil {
			const errorMessage = "failed to get last consolidated block number from state"
			ctx.Logger().Errorf("%v:%v", errorMessage, err)
			return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
		}

		return hex.EncodeUint64(lastBlockNumber), nil
	})
}

// IsBlockConsolidated returns the consolidation status of a provided block number
func (e *ZKEVMEndpoints) IsBlockConsolidated(ctx *context.RequestContext, blockNumber types.ArgUint64) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(ctx, e.state, func(ctx *context.RequestContext, dbTx pgx.Tx) (interface{}, types.Error) {
		IsL2BlockConsolidated, err := e.state.IsL2BlockConsolidated(ctx, uint64(blockNumber), dbTx)
		if err != nil {
			const errorMessage = "failed to check if the block is consolidated"
			ctx.Logger().Errorf("%v: %v", errorMessage, err)
			return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
		}

		return IsL2BlockConsolidated, nil
	})
}

// IsBlockVirtualized returns the virtualization status of a provided block number
func (e *ZKEVMEndpoints) IsBlockVirtualized(ctx *context.RequestContext, blockNumber types.ArgUint64) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(ctx, e.state, func(ctx *context.RequestContext, dbTx pgx.Tx) (interface{}, types.Error) {
		IsL2BlockVirtualized, err := e.state.IsL2BlockVirtualized(ctx, uint64(blockNumber), dbTx)
		if err != nil {
			const errorMessage = "failed to check if the block is virtualized"
			ctx.Logger().Errorf("%v: %v", errorMessage, err)
			return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
		}

		return IsL2BlockVirtualized, nil
	})
}

// BatchNumberByBlockNumber returns the batch number from which the passed block number is created
func (e *ZKEVMEndpoints) BatchNumberByBlockNumber(ctx *context.RequestContext, blockNumber types.ArgUint64) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(ctx, e.state, func(ctx *context.RequestContext, dbTx pgx.Tx) (interface{}, types.Error) {
		batchNum, err := e.state.BatchNumberByL2BlockNumber(ctx, uint64(blockNumber), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			const errorMessage = "failed to get batch number from block number"
			ctx.Logger().Errorf("%v: %v", errorMessage, err.Error())
			return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
		}

		return hex.EncodeUint64(batchNum), nil
	})
}

// BatchNumber returns the latest virtualized batch number
func (e *ZKEVMEndpoints) BatchNumber(ctx *context.RequestContext) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(ctx, e.state, func(ctx *context.RequestContext, dbTx pgx.Tx) (interface{}, types.Error) {
		lastBatchNumber, err := e.state.GetLastBatchNumber(ctx, dbTx)
		if err != nil {
			return "0x0", types.NewRPCError(types.DefaultErrorCode, "failed to get the last batch number from state")
		}

		return hex.EncodeUint64(lastBatchNumber), nil
	})
}

// VirtualBatchNumber returns the latest virtualized batch number
func (e *ZKEVMEndpoints) VirtualBatchNumber(ctx *context.RequestContext) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(ctx, e.state, func(ctx *context.RequestContext, dbTx pgx.Tx) (interface{}, types.Error) {
		lastBatchNumber, err := e.state.GetLastVirtualBatchNum(ctx, dbTx)
		if err != nil {
			return "0x0", types.NewRPCError(types.DefaultErrorCode, "failed to get the last virtual batch number from state")
		}

		return hex.EncodeUint64(lastBatchNumber), nil
	})
}

// VerifiedBatchNumber returns the latest verified batch number
func (e *ZKEVMEndpoints) VerifiedBatchNumber(ctx *context.RequestContext) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(ctx, e.state, func(ctx *context.RequestContext, dbTx pgx.Tx) (interface{}, types.Error) {
		lastBatch, err := e.state.GetLastVerifiedBatch(ctx, dbTx)
		if err != nil {
			return "0x0", types.NewRPCError(types.DefaultErrorCode, "failed to get the last verified batch number from state")
		}

		return hex.EncodeUint64(lastBatch.BatchNumber), nil
	})
}

// GetBatchByNumber returns information about a batch by batch number
func (e *ZKEVMEndpoints) GetBatchByNumber(ctx *context.RequestContext, batchNumber types.BatchNumber, fullTx bool) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(ctx, e.state, func(ctx *context.RequestContext, dbTx pgx.Tx) (interface{}, types.Error) {
		var err error
		batchNumber, rpcErr := batchNumber.GetNumericBatchNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		batch, err := e.state.GetBatchByNumber(ctx, batchNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load batch from state by number %v", batchNumber), err)
		}

		txs, err := e.state.GetTransactionsByBatchNumber(ctx, batchNumber, dbTx)
		if !errors.Is(err, state.ErrNotFound) && err != nil {
			return rpcErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load batch txs from state by number %v", batchNumber), err)
		}

		receipts := make([]ethTypes.Receipt, 0, len(txs))
		for _, tx := range txs {
			receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
			if err != nil {
				return rpcErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load receipt for tx %v", tx.Hash().String()), err)
			}
			receipts = append(receipts, *receipt)
		}

		virtualBatch, err := e.state.GetVirtualBatch(ctx, batchNumber, dbTx)
		if err != nil && !errors.Is(err, state.ErrNotFound) {
			return rpcErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load virtual batch from state by number %v", batchNumber), err)
		}

		verifiedBatch, err := e.state.GetVerifiedBatch(ctx, batchNumber, dbTx)
		if err != nil && !errors.Is(err, state.ErrNotFound) {
			return rpcErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load virtual batch from state by number %v", batchNumber), err)
		}

		ger, err := e.state.GetExitRootByGlobalExitRoot(ctx, batch.GlobalExitRoot, dbTx)
		if err != nil && !errors.Is(err, state.ErrNotFound) {
			return rpcErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load full GER from state by number %v", batchNumber), err)
		} else if errors.Is(err, state.ErrNotFound) {
			ger = &state.GlobalExitRoot{}
		}

		batch.Transactions = txs
		rpcBatch := types.NewBatch(batch, virtualBatch, verifiedBatch, receipts, fullTx, ger)

		return rpcBatch, nil
	})
}
