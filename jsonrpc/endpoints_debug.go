package jsonrpc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// DebugEndpoints is the debug jsonrpc endpoint
type DebugEndpoints struct {
	state stateInterface
	txMan dbTxManager
}

type traceConfig struct {
	DisableStorage   bool    `json:"disableStorage"`
	DisableStack     bool    `json:"disableStack"`
	EnableMemory     bool    `json:"enableMemory"`
	EnableReturnData bool    `json:"enableReturnData"`
	Tracer           *string `json:"tracer"`
}

type traceBlockTransactionResponse struct {
	Result traceTransactionResponse `json:"result"`
}

type traceTransactionResponse struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue interface{}    `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

// StructLogRes represents a log response.
type StructLogRes struct {
	Pc            uint64             `json:"pc"`
	Op            string             `json:"op"`
	Gas           uint64             `json:"gas"`
	GasCost       uint64             `json:"gasCost"`
	Depth         int                `json:"depth"`
	Error         string             `json:"error,omitempty"`
	Stack         []argBig           `json:"stack"`
	Memory        []string           `json:"memory"`
	Storage       *map[string]string `json:"storage,omitempty"`
	RefundCounter uint64             `json:"refund,omitempty"`
}

// TraceTransaction creates a response for debug_traceTransaction request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtracetransaction
func (d *DebugEndpoints) TraceTransaction(hash common.Hash, cfg *traceConfig) (interface{}, rpcError) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		return d.buildTraceTransaction(ctx, hash, cfg, dbTx)
	})
}

// TraceTransaction creates a response for debug_traceBlockByNumber request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtraceblockbynumber
func (d *DebugEndpoints) TraceBlockByNumber(number BlockNumber, cfg *traceConfig) (interface{}, rpcError) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		blockNumber, rpcErr := number.getNumericBlockNumber(ctx, d.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		block, err := d.state.GetL2BlockByNumber(ctx, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, newRPCError(defaultErrorCode, "genesis is not traceable")
		} else if err == state.ErrNotFound {
			return rpcErrorResponse(defaultErrorCode, "failed to get block by number", err)
		}

		traces, rpcErr := d.buildTraceBlock(ctx, block.Transactions(), cfg, dbTx)
		if err != nil {
			return nil, rpcErr
		}

		return traces, nil
	})
}

// TraceTransaction creates a response for debug_traceBlockByHash request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtraceblockbyhash
func (d *DebugEndpoints) TraceBlockByHash(hash common.Hash, cfg *traceConfig) (interface{}, rpcError) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		block, err := d.state.GetL2BlockByHash(ctx, hash, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, newRPCError(defaultErrorCode, "genesis is not traceable")
		} else if err == state.ErrNotFound {
			return rpcErrorResponse(defaultErrorCode, "failed to get block by hash", err)
		}

		traces, rpcErr := d.buildTraceBlock(ctx, block.Transactions(), cfg, dbTx)
		if err != nil {
			return nil, rpcErr
		}

		return traces, nil
	})
}

func (d *DebugEndpoints) buildTraceBlock(ctx context.Context, txs []*types.Transaction, cfg *traceConfig, dbTx pgx.Tx) (interface{}, rpcError) {
	traces := []traceBlockTransactionResponse{}
	for _, tx := range txs {
		traceTransaction, err := d.buildTraceTransaction(ctx, tx.Hash(), cfg, dbTx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get trace for transaction %v", tx.Hash().String())
			return rpcErrorResponse(defaultErrorCode, errMsg, err)
		}
		traceBlockTransaction := traceBlockTransactionResponse{
			Result: traceTransaction.(traceTransactionResponse),
		}
		traces = append(traces, traceBlockTransaction)
	}

	return traces, nil
}

func (d *DebugEndpoints) buildTraceTransaction(ctx context.Context, hash common.Hash, cfg *traceConfig, dbTx pgx.Tx) (interface{}, rpcError) {
	tracer := ""
	if cfg != nil && cfg.Tracer != nil {
		tracer = *cfg.Tracer
	}

	result, err := d.state.DebugTransaction(ctx, hash, tracer, dbTx)
	if err != nil {
		const errorMessage = "failed to get trace"
		log.Infof("%v: %v", errorMessage, err)
		return nil, newRPCError(defaultErrorCode, errorMessage)
	}

	if tracer != "" && len(result.ExecutorTraceResult) > 0 {
		return result.ExecutorTraceResult, nil
	}

	failed := result.Failed()
	var returnValue interface{}
	if cfg.EnableReturnData {
		returnValue = common.Bytes2Hex(result.ReturnValue)
	}

	structLogs := d.buildStructLogs(result.StructLogs, cfg)

	resp := traceTransactionResponse{
		Gas:         result.GasUsed,
		Failed:      failed,
		ReturnValue: returnValue,
		StructLogs:  structLogs,
	}

	return resp, nil
}

func (d *DebugEndpoints) buildStructLogs(stateStructLogs []instrumentation.StructLog, cfg *traceConfig) []StructLogRes {
	structLogs := make([]StructLogRes, 0, len(stateStructLogs))
	for _, structLog := range stateStructLogs {

		stack := make([]argBig, 0, len(structLog.Stack))
		if !cfg.DisableStack && len(structLog.Stack) > 0 {
			for _, stackItem := range structLog.Stack {
				if stackItem != nil {
					stack = append(stack, argBig(*stackItem))
				}
			}
		}

		memory := make([]string, 0, len(structLog.Memory))
		if cfg.EnableMemory {
			for _, memoryItem := range structLog.Memory {
				slice32Bytes := make([]byte, 32)
				slice32Bytes[31] = memoryItem
				memoryStringItem := hex.EncodeToString(slice32Bytes)
				memory = append(memory, memoryStringItem)
			}
		}

		var storageRes *map[string]string
		if !cfg.DisableStorage && len(structLog.Storage) > 0 {
			storage := make(map[string]string, len(structLog.Storage))
			for storageKey, storageValue := range structLog.Storage {
				storage[storageKey.String()] = storageValue.String()
			}
			storageRes = &storage
		}

		errRes := ""
		if structLog.Err != nil {
			errRes = structLog.Err.Error()
		}

		op := structLog.Op
		if op == "SHA3" {
			op = "KECCAK256"
		}

		structLogs = append(structLogs, StructLogRes{
			Pc:            structLog.Pc,
			Op:            op,
			Gas:           structLog.Gas,
			GasCost:       structLog.GasCost,
			Depth:         structLog.Depth,
			Error:         errRes,
			Stack:         stack,
			Memory:        memory,
			Storage:       storageRes,
			RefundCounter: structLog.RefundCounter,
		})
	}
	return structLogs
}
