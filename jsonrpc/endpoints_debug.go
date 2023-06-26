package jsonrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

var defaultTraceConfig = &traceConfig{
	DisableStorage:   false,
	DisableStack:     false,
	EnableMemory:     false,
	EnableReturnData: false,
	Tracer:           nil,
}

// DebugEndpoints is the debug jsonrpc endpoint
type DebugEndpoints struct {
	state types.StateInterface
	txMan dbTxManager
}

type traceConfig struct {
	DisableStorage   bool            `json:"disableStorage"`
	DisableStack     bool            `json:"disableStack"`
	EnableMemory     bool            `json:"enableMemory"`
	EnableReturnData bool            `json:"enableReturnData"`
	Tracer           *string         `json:"tracer"`
	TracerConfig     json.RawMessage `json:"tracerConfig"`
}

// StructLogRes represents the debug trace information for each opcode
type StructLogRes struct {
	Pc            uint64             `json:"pc"`
	Op            string             `json:"op"`
	Gas           uint64             `json:"gas"`
	GasCost       uint64             `json:"gasCost"`
	Depth         int                `json:"depth"`
	Error         string             `json:"error,omitempty"`
	Stack         *[]types.ArgBig    `json:"stack,omitempty"`
	Memory        *[]string          `json:"memory,omitempty"`
	Storage       *map[string]string `json:"storage,omitempty"`
	RefundCounter uint64             `json:"refund,omitempty"`
}

type traceTransactionResponse struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue interface{}    `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

type traceBlockTransactionResponse struct {
	Result interface{} `json:"result"`
}

// TraceTransaction creates a response for debug_traceTransaction request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtracetransaction
func (d *DebugEndpoints) TraceTransaction(hash types.ArgHash, cfg *traceConfig) (interface{}, types.Error) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		return d.buildTraceTransaction(ctx, hash.Hash(), cfg, dbTx)
	})
}

// TraceBlockByNumber creates a response for debug_traceBlockByNumber request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtraceblockbynumber
func (d *DebugEndpoints) TraceBlockByNumber(number types.BlockNumber, cfg *traceConfig) (interface{}, types.Error) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		blockNumber, rpcErr := number.GetNumericBlockNumber(ctx, d.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		block, err := d.state.GetL2BlockByNumber(ctx, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("block #%d not found", blockNumber))
		} else if err == state.ErrNotFound {
			return rpcErrorResponse(types.DefaultErrorCode, "failed to get block by number", err)
		}

		traces, rpcErr := d.buildTraceBlock(ctx, block.Transactions(), cfg, dbTx)
		if err != nil {
			return nil, rpcErr
		}

		return traces, nil
	})
}

// TraceBlockByHash creates a response for debug_traceBlockByHash request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtraceblockbyhash
func (d *DebugEndpoints) TraceBlockByHash(hash types.ArgHash, cfg *traceConfig) (interface{}, types.Error) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		block, err := d.state.GetL2BlockByHash(ctx, hash.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("block %s not found", hash.Hash().String()))
		} else if err == state.ErrNotFound {
			return rpcErrorResponse(types.DefaultErrorCode, "failed to get block by hash", err)
		}

		traces, rpcErr := d.buildTraceBlock(ctx, block.Transactions(), cfg, dbTx)
		if err != nil {
			return nil, rpcErr
		}

		return traces, nil
	})
}

func (d *DebugEndpoints) buildTraceBlock(ctx context.Context, txs []*ethTypes.Transaction, cfg *traceConfig, dbTx pgx.Tx) (interface{}, types.Error) {
	traces := []traceBlockTransactionResponse{}
	for _, tx := range txs {
		traceTransaction, err := d.buildTraceTransaction(ctx, tx.Hash(), cfg, dbTx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get trace for transaction %v", tx.Hash().String())
			return rpcErrorResponse(types.DefaultErrorCode, errMsg, err)
		}
		traceBlockTransaction := traceBlockTransactionResponse{
			Result: traceTransaction,
		}
		traces = append(traces, traceBlockTransaction)
	}

	return traces, nil
}

func (d *DebugEndpoints) buildTraceTransaction(ctx context.Context, hash common.Hash, cfg *traceConfig, dbTx pgx.Tx) (interface{}, types.Error) {
	traceCfg := cfg
	if traceCfg == nil {
		traceCfg = defaultTraceConfig
	}

	// check tracer
	if traceCfg.Tracer != nil && *traceCfg.Tracer != "" && !isBuiltInTracer(*traceCfg.Tracer) && !isJSCustomTracer(*traceCfg.Tracer) {
		return rpcErrorResponse(types.DefaultErrorCode, "invalid tracer", nil)
	}

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
		return rpcErrorResponse(types.DefaultErrorCode, "transaction not found", nil)
	} else if err != nil {
		const errorMessage = "failed to get trace"
		log.Errorf("%v: %v", errorMessage, err)
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	// if a tracer was specified, then return the trace result
	if stateTraceConfig.Tracer != nil && *stateTraceConfig.Tracer != "" && len(result.ExecutorTraceResult) > 0 {
		return result.ExecutorTraceResult, nil
	}

	receipt, err := d.state.GetTransactionReceipt(ctx, hash, dbTx)
	if err != nil {
		const errorMessage = "failed to tx receipt"
		log.Errorf("%v: %v", errorMessage, err)
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	failed := receipt.Status == ethTypes.ReceiptStatusFailed
	var returnValue interface{}
	if stateTraceConfig.EnableReturnData {
		returnValue = common.Bytes2Hex(result.ReturnValue)
	}

	structLogs := d.buildStructLogs(result.StructLogs, *traceCfg)

	resp := traceTransactionResponse{
		Gas:         result.GasUsed,
		Failed:      failed,
		ReturnValue: returnValue,
		StructLogs:  structLogs,
	}

	return resp, nil
}

func (d *DebugEndpoints) buildStructLogs(stateStructLogs []instrumentation.StructLog, cfg traceConfig) []StructLogRes {
	structLogs := make([]StructLogRes, 0, len(stateStructLogs))
	memory := fakevm.NewMemory()
	for _, structLog := range stateStructLogs {
		errRes := ""
		if structLog.Err != nil {
			errRes = structLog.Err.Error()
		}

		op := structLog.Op
		if op == "SHA3" {
			op = "KECCAK256"
		} else if op == "STOP" && structLog.Pc == 0 {
			// this stop is generated for calls with single
			// step(no depth increase) and must be ignored
			continue
		}

		structLogRes := StructLogRes{
			Pc:            structLog.Pc,
			Op:            op,
			Gas:           structLog.Gas,
			GasCost:       structLog.GasCost,
			Depth:         structLog.Depth,
			Error:         errRes,
			RefundCounter: structLog.RefundCounter,
		}

		if !cfg.DisableStack {
			stack := make([]types.ArgBig, 0, len(structLog.Stack))
			for _, stackItem := range structLog.Stack {
				if stackItem != nil {
					stack = append(stack, types.ArgBig(*stackItem))
				}
			}
			structLogRes.Stack = &stack
		}

		if cfg.EnableMemory {
			memory.Resize(uint64(structLog.MemorySize))
			if len(structLog.Memory) > 0 {
				memory.Set(uint64(structLog.MemoryOffset), uint64(len(structLog.Memory)), structLog.Memory)
			}

			if structLog.MemorySize > 0 {
				// Populate the structLog memory
				structLog.Memory = memory.Data()

				// Convert memory to string array
				const memoryChunkSize = 32
				memoryArray := make([]string, 0, len(structLog.Memory))

				for i := 0; i < len(structLog.Memory); i = i + memoryChunkSize {
					slice32Bytes := make([]byte, memoryChunkSize)
					copy(slice32Bytes, structLog.Memory[i:i+memoryChunkSize])
					memoryStringItem := hex.EncodeToString(slice32Bytes)
					memoryArray = append(memoryArray, memoryStringItem)
				}

				structLogRes.Memory = &memoryArray
			} else {
				memory = fakevm.NewMemory()
				structLogRes.Memory = &[]string{}
			}
		}

		if !cfg.DisableStorage && len(structLog.Storage) > 0 {
			storage := make(map[string]string, len(structLog.Storage))
			for storageKey, storageValue := range structLog.Storage {
				k := hex.EncodeToString(storageKey.Bytes())
				v := hex.EncodeToString(storageValue.Bytes())
				storage[k] = v
			}
			structLogRes.Storage = &storage
		}

		structLogs = append(structLogs, structLogRes)
	}
	return structLogs
}

// isBuiltInTracer checks if the tracer is one of the
// built-in tracers
func isBuiltInTracer(tracer string) bool {
	// built-in tracers
	switch tracer {
	case "callTracer", "4byteTracer", "prestateTracer", "noopTracer":
		return true
	default:
		return false
	}
}

// isJSCustomTracer checks if the tracer contains the
// functions result and fault which are required for a custom tracer
// https://geth.ethereum.org/docs/developers/evm-tracing/custom-tracer
func isJSCustomTracer(tracer string) bool {
	return strings.Contains(tracer, "result") && strings.Contains(tracer, "fault")
}
