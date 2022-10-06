package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// Debug is the debug jsonrpc endpoint
type Debug struct {
	state stateInterface
	txMan dbTxManager
}

type traceConfig struct {
	Tracer *string `json:"tracer"`
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
	Stack         *[]argBig          `json:"stack,omitempty"`
	Memory        *argBytes          `json:"memory,omitempty"`
	Storage       *map[string]string `json:"storage,omitempty"`
	RefundCounter uint64             `json:"refund,omitempty"`
}

// TraceTransaction creates a response for debug_traceTransaction request.
// See https://geth.ethereum.org/docs/rpc/ns-debug#debug_tracetransaction
func (d *Debug) TraceTransaction(hash common.Hash, cfg *traceConfig) (interface{}, rpcError) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		tracer := ""
		if cfg != nil && cfg.Tracer != nil {
			tracer = *cfg.Tracer
		}

		result, err := d.state.DebugTransaction(ctx, hash, tracer, dbTx)
		if err != nil {
			const errorMessage = "failed to debug trace the transaction"
			log.Debugf("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		if tracer != "" && len(result.ExecutorTraceResult) > 0 {
			return result.ExecutorTraceResult, nil
		}

		failed := result.Failed()
		structLogs := make([]StructLogRes, 0, len(result.StructLogs))
		for _, structLog := range result.StructLogs {
			var stackRes *[]argBig
			if len(structLog.Stack) > 0 {
				stack := make([]argBig, 0, len(structLog.Stack))
				for _, stackItem := range structLog.Stack {
					if stackItem != nil {
						stack = append(stack, argBig(*stackItem))
					}
				}
				stackRes = &stack
			}

			var memoryRes *argBytes
			if len(structLog.Memory) > 0 {
				memory := make(argBytes, 0, len(structLog.Memory))
				for _, memoryItem := range structLog.Memory {
					memory = append(memory, memoryItem)
				}
				memoryRes = &memory
			}

			var storageRes *map[string]string
			if len(structLog.Storage) > 0 {
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

			structLogs = append(structLogs, StructLogRes{
				Pc:            structLog.Pc,
				Op:            structLog.Op,
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

		resp := traceTransactionResponse{
			Gas:         result.GasUsed,
			Failed:      failed,
			ReturnValue: result.ReturnValue,
			StructLogs:  structLogs,
		}

		return resp, nil
	})
}
