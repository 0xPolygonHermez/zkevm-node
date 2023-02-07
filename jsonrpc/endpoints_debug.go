package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
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
	Stack         *[]argBig          `json:"stack"`
	Memory        *argBytes          `json:"memory"`
	Storage       *map[string]string `json:"storage,omitempty"`
	RefundCounter uint64             `json:"refund,omitempty"`
}

// TraceTransaction creates a response for debug_traceTransaction request.
// See https://geth.ethereum.org/docs/rpc/ns-debug#debug_tracetransaction
func (d *DebugEndpoints) TraceTransaction(hash common.Hash, cfg *traceConfig) (interface{}, rpcError) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		tracer := ""
		if cfg != nil && cfg.Tracer != nil {
			tracer = *cfg.Tracer
		}

		result, err := d.state.DebugTransaction(ctx, hash, tracer, dbTx)
		if err != nil {
			const errorMessage = "failed to debug trace the transaction"
			log.Infof("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		if tracer != "" && len(result.ExecutorTraceResult) > 0 {
			return result.ExecutorTraceResult, nil
		}

		failed := result.Failed()
		structLogs := make([]StructLogRes, 0, len(result.StructLogs))
		for _, structLog := range result.StructLogs {
			var stackRes *[]argBig
			if !cfg.DisableStack && len(structLog.Stack) > 0 {
				stack := make([]argBig, 0, len(structLog.Stack))
				for _, stackItem := range structLog.Stack {
					if stackItem != nil {
						stack = append(stack, argBig(*stackItem))
					}
				}
				stackRes = &stack
			}

			var memoryRes *argBytes
			if cfg.EnableMemory && len(structLog.Memory) > 0 {
				memory := make(argBytes, 0, len(structLog.Memory))
				for _, memoryItem := range structLog.Memory {
					memory = append(memory, memoryItem)
				}
				memoryRes = &memory
			}

			var storageRes *map[string]string
			if !cfg.DisableStorage && len(structLog.Storage) > 0 {
				storage := make(map[string]string, len(structLog.Storage))
				for storageKey, storageValue := range structLog.Storage {
					storage[storageKey.Hex()] = storageValue.Hex()
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
				Stack:         stackRes,
				Memory:        memoryRes,
				Storage:       storageRes,
				RefundCounter: structLog.RefundCounter,
			})
		}

		var returnValue interface{}
		if cfg.EnableReturnData {
			returnValue = common.Bytes2Hex(result.ReturnValue)
		}

		resp := traceTransactionResponse{
			Gas:         result.GasUsed,
			Failed:      failed,
			ReturnValue: returnValue,
			StructLogs:  structLogs,
		}

		return resp, nil
	})
}
