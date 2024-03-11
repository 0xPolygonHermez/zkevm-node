package jsonrpc

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

func (d *DebugEndpoints) buildInnerTransaction(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (interface{}, types.Error) {
	traceCfg := defaultTraceConfig
	tracer := "callTracer"
	traceCfg.Tracer = &tracer

	stateTraceConfig := state.TraceConfig{
		DisableStack:     traceCfg.DisableStack,
		DisableStorage:   traceCfg.DisableStorage,
		EnableMemory:     traceCfg.EnableMemory,
		EnableReturnData: traceCfg.EnableReturnData,
		Tracer:           traceCfg.Tracer,
		TracerConfig:     traceCfg.TracerConfig,
	}
	result, err := d.state.DebugTransaction(ctx, hash, stateTraceConfig, dbTx)
	if errors.Is(err, state.ErrNotFound) {
		return RPCErrorResponse(types.DefaultErrorCode, "transaction not found", nil, true)
	} else if err != nil {
		const errorMessage = "failed to get trace"
		log.Errorf("%v: %v", errorMessage, err)
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	return result.TraceResult, nil
}
