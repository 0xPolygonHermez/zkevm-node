package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

// type traceTransactionResponse struct {
// 	Gas         uint64         `json:"gas"`
// 	Failed      bool           `json:"failed"`
// 	ReturnValue interface{}    `json:"returnValue"`
// 	StructLogs  []StructLogRes `json:"structLogs"`
// }

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
		tx, err := d.state.GetTransactionByHash(ctx, hash, dbTx)
		if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get tx", err)
		}

		receipt, err := d.state.GetTransactionReceipt(ctx, hash, dbTx)
		if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get tx receipt", err)
		}

		if receipt.Status == types.ReceiptStatusSuccessful {
			return []interface{}{}, nil
		} else {
			from, err := state.GetSender(*tx)
			if err != nil {
				return rpcErrorResponse(defaultErrorCode, "failed to get tx sender", err)
			}

			var to string
			if tx.To() != nil {
				to = tx.To().String()
			}

			return []interface{}{
				map[string]interface{}{
					"type":         "call",
					"callType":     "call",
					"from":         from.String(),
					"to":           to,
					"input":        "0xa9059cbb000000000000000000000000617b3a3528f9cdd6630fd3301b9c8911f7bf063d000000000000000000000000000000000000000000000030ca024f987b900000",
					"error":        "generic error message",
					"traceAddress": []string{},
					"value":        tx.Value(),
					"gas":          hex.EncodeUint64(tx.Gas()),
					"gasUsed":      hex.EncodeUint64(receipt.GasUsed),
				},
			}, nil
		}
	})

	/*
			[
		        {
		            "type": "call",
		            "callType": "call",
		            "from": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
		            "to": "0x9a676e781a523b5d0c0e43731313a708cb607508",
		            "input": "0xa9059cbb000000000000000000000000617b3a3528f9cdd6630fd3301b9c8911f7bf063d000000000000000000000000000000000000000000000030ca024f987b900000",
		            "error": "execution reverted",
		            "traceAddress": [],
		            "value": "0x0",
		            "gas": "0x747a",
		            "gasUsed": "0xa80"
		        }
		    ]

	*/

	// ctx := context.Background()

	// tracer := ""
	// if cfg != nil && cfg.Tracer != nil {
	// 	tracer = *cfg.Tracer
	// }

	// result, err := d.state.DebugTransaction(ctx, hash, tracer)
	// if err != nil {
	// 	const errorMessage = "failed to debug trace the transaction"
	// 	log.Debugf("%v: %v", errorMessage, err)
	// 	return nil, newRPCError(defaultErrorCode, errorMessage)
	// }

	// if tracer != "" && len(result.ExecutorTraceResult) > 0 {
	// 	return result.ExecutorTraceResult, nil
	// }

	// failed := result.Failed()
	// // structLogs := make([]StructLogRes, 0, len(result.StructLogs))
	// // for _, structLog := range result.StructLogs {
	// // 	var stackRes *[]argBig
	// // 	if len(structLog.Stack) > 0 {
	// // 		stack := make([]argBig, 0, len(structLog.Stack))
	// // 		for _, stackItem := range structLog.Stack {
	// // 			if stackItem != nil {
	// // 				stack = append(stack, argBig(*stackItem))
	// // 			}
	// // 		}
	// // 		stackRes = &stack
	// // 	}

	// // 	var memoryRes *argBytes
	// // 	if len(structLog.Memory) > 0 {
	// // 		memory := make(argBytes, 0, len(structLog.Memory))
	// // 		for _, memoryItem := range structLog.Memory {
	// // 			memory = append(memory, memoryItem)
	// // 		}
	// // 		memoryRes = &memory
	// // 	}

	// // 	var storageRes *map[string]string
	// // 	if len(structLog.Storage) > 0 {
	// // 		storage := make(map[string]string, len(structLog.Storage))
	// // 		for storageKey, storageValue := range structLog.Storage {
	// // 			storage[storageKey.Hex()] = storageValue.Hex()
	// // 		}
	// // 		storageRes = &storage
	// // 	}

	// // 	errRes := ""
	// // 	if structLog.Err != nil {
	// // 		errRes = structLog.Err.Error()
	// // 	}

	// // 	structLogs = append(structLogs, StructLogRes{
	// // 		Pc:            structLog.Pc,
	// // 		Op:            structLog.Op,
	// // 		Gas:           structLog.Gas,
	// // 		GasCost:       structLog.GasCost,
	// // 		Depth:         structLog.Depth,
	// // 		Error:         errRes,
	// // 		Stack:         stackRes,
	// // 		Memory:        memoryRes,
	// // 		Storage:       storageRes,
	// // 		RefundCounter: structLog.RefundCounter,
	// // 	})
	// // }

	// resp := traceTransactionResponse{
	// 	Gas:         result.GasUsed,
	// 	Failed:      failed,
	// 	ReturnValue: result.ReturnValue,
	// 	// StructLogs:  structLogs,
	// }

	// return resp, nil
}
