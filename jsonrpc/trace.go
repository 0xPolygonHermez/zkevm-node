package jsonrpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation"
)

// Trace is the trace jsonrpc endpoint
type Trace struct {
	state stateInterface
}

type replayTransactionResponse struct {
	Output    *argBytes                        `json:"output"`
	Trace     []txTrace                        `json:"trace"`
	VMTrace   txVMTrace                        `json:"vmTrace"`
	StateDiff map[common.Address]txAccountDiff `json:"stateDiff"`
}

type blockResponse struct {
	Action              txTraceAction `json:"action"`
	Result              txTraceResult `json:"result"`
	TraceAddress        []uint64      `json:"traceAddress"`
	SubTraces           uint64        `json:"subTraces"`
	TransactionPosition argUint64     `json:"transactionPosition"`
	TransactionHash     common.Hash   `json:"transactionHash"`
	BlockNumber         uint64        `json:"blockNumber"`
	BlockHash           common.Hash   `json:"blockHash"`
	Type                string        `json:"type"`
}

type txTrace struct {
	TraceAddress []uint64       `json:"traceAddress"`
	SubTraces    uint64         `json:"subtraces"`
	Action       txTraceAction  `json:"action"`
	Result       *txTraceResult `json:"result,omitempty"`
	Error        *string        `json:"error,omitempty"`
	Type         string         `json:"type"`
}

type txVMTrace struct {
	Code       argBytes             `json:"code"`
	Operations []txVMTraceOperation `json:"ops"`
}

type txAccountDiff struct{}

type txVMTraceOperation struct {
	PC                uint64                     `json:"pc"`
	Cost              uint64                     `json:"cost"`
	ExecutedOperation txVMTraceExecutedOperation `json:"ex"`
	Sub               *txVMTrace                 `json:"sub"`
}

type txVMTraceExecutedOperation struct {
	Used        argUint64             `json:"used"`
	Push        []argUint64           `json:"push"`
	MemoryDiff  *txVMTraceMemoryDiff  `json:"mem"`
	StorageDiff *txVMTraceStorageDiff `json:"store"`
}

type txVMTraceMemoryDiff struct {
	Off  argUint64 `json:"off"`
	Data []byte    `json:"data"`
}

type txVMTraceStorageDiff struct {
	Key   argUint64 `json:"key"`
	Value argUint64 `json:"val"`
}

type txTraceAction struct {
	From     string    `json:"from"`
	To       string    `json:"to"`
	Value    argUint64 `json:"value"`
	Gas      argUint64 `json:"gas"`
	Input    argBytes  `json:"input"`
	CallType string    `json:"callType"`
}

type txTraceResult struct {
	GasUsed argUint64 `json:"gasUsed"`
	Output  argBytes  `json:"output"`
}

// ReplayBlockTransactions creates a response for trace_replayBlockTransactions request.
// See https://openethereum.github.io/JSONRPC-trace-module#trace_replayblocktransactions
func (t *Trace) ReplayBlockTransactions(number *BlockNumber, traceMode []string) (interface{}, error) {
	ctx := context.Background()

	var err error
	batchNumber, err := number.getNumericBlockNumber(ctx, t.state)
	if err != nil {
		return nil, err
	}

	results, err := t.state.ReplayBatchTransactions(batchNumber, traceMode)
	if err != nil {
		return nil, err
	}

	response := make([]replayTransactionResponse, 0, len(results))

	for _, result := range results {
		response = append(response, t.executionResultToReplayTransactionResponse(result))
	}

	return response, nil
}

// ReplayTransaction creates a response for trace_replayTransaction request.
// See https://openethereum.github.io/JSONRPC-trace-module#trace_replaytransaction
func (t *Trace) ReplayTransaction(hash common.Hash, traceMode []string) (interface{}, error) {
	result := t.state.ReplayTransaction(hash, traceMode)
	response := t.executionResultToReplayTransactionResponse(result)

	return response, nil
}

// Block creates a response for trace_block request.
// See https://openethereum.github.io/JSONRPC-trace-module#trace_block
func (t *Trace) Block(number *BlockNumber) (interface{}, error) {
	ctx := context.Background()

	var err error
	batchNumber, err := number.getNumericBlockNumber(ctx, t.state)
	if err != nil {
		return nil, err
	}

	batch, err := t.state.GetBatchByNumber(ctx, batchNumber, "")
	if err != nil {
		return nil, err
	}

	results, err := t.state.ReplayBatchTransactions(batchNumber, []string{"trace"})
	if err != nil {
		return nil, err
	}

	response := make([]blockResponse, 0, len(results))

	for _, result := range results {
		response = append(response, t.executionResultToBlockResponse(batch, result)...)
	}

	return response, nil
}

func (t *Trace) executionResultToReplayTransactionResponse(result *runtime.ExecutionResult) replayTransactionResponse {
	trace := make([]txTrace, 0, len(result.Trace))
	for _, tr := range result.Trace {
		trace = append(trace, t.traceToTxTrace(tr))
	}

	return replayTransactionResponse{
		// Output:  ,
		Trace:   trace,
		VMTrace: t.vmTraceToTxVMTrace(result.VMTrace),
		// StateDiff: ,
	}
}

func (t *Trace) traceToTxTrace(tr instrumentation.Trace) txTrace {
	txT := txTrace{
		TraceAddress: tr.TraceAddress[:],
		SubTraces:    tr.SubTraces,
		Action: txTraceAction{
			From:     tr.Action.From,
			To:       tr.Action.To,
			Value:    argUint64(tr.Action.Value),
			Gas:      argUint64(tr.Action.Gas),
			Input:    tr.Action.Input,
			CallType: tr.Action.CallType,
		},
		Type: tr.Type,
	}

	if tr.Result != nil {
		txT.Result = &txTraceResult{
			GasUsed: argUint64(tr.Result.GasUsed),
			Output:  argBytes(tr.Result.Output),
		}
	} else {
		txT.Error = tr.Error
	}

	return txT
}

func (t *Trace) vmTraceToTxVMTrace(vmTrace instrumentation.VMTrace) txVMTrace {
	operations := make([]txVMTraceOperation, 0, len(vmTrace.Operations))

	for _, op := range vmTrace.Operations {
		stackPush := make([]argUint64, 0, len(op.Executed.StackPush))
		for _, sp := range op.Executed.StackPush {
			stackPush = append(stackPush, argUint64(sp))
		}

		var sub *txVMTrace
		if op.Sub != nil {
			s := t.vmTraceToTxVMTrace(*op.Sub)
			sub = &s
		}

		operations = append(operations, txVMTraceOperation{
			PC:   op.Pc,
			Cost: op.GasCost,
			ExecutedOperation: txVMTraceExecutedOperation{
				Used: argUint64(op.Executed.GasUsed),
				Push: stackPush,
				MemoryDiff: &txVMTraceMemoryDiff{
					Off:  argUint64(op.Executed.MemDiff.Offset),
					Data: op.Executed.MemDiff.Data[:],
				},
				StorageDiff: &txVMTraceStorageDiff{
					Key:   argUint64(op.Executed.StoreDiff.Location),
					Value: argUint64(op.Executed.StoreDiff.Value),
				},
			},
			Sub: sub,
		})
	}

	return txVMTrace{
		Code:       vmTrace.Code,
		Operations: operations,
	}
}

func (t *Trace) executionResultToBlockResponse(b *state.Batch, result *runtime.ExecutionResult) []blockResponse {
	txMap := make(map[common.Hash]*types.Transaction, len(b.Transactions))
	for _, tx := range b.Transactions {
		txMap[tx.Hash()] = tx
	}

	blockResponses := make([]blockResponse, 0, len(result.Trace))

	for _, tr := range result.Trace {
		txTrace := t.traceToTxTrace(tr)

		blockResponses = append(blockResponses, blockResponse{
			Action:       txTrace.Action,
			Result:       *txTrace.Result,
			TraceAddress: txTrace.TraceAddress,
			SubTraces:    txTrace.SubTraces,
			// TransactionPosition: ,
			// TransactionHash:     ,
			BlockNumber: b.Number().Uint64(),
			BlockHash:   b.Hash(),
			Type:        txTrace.Type,
		})
	}

	return blockResponses
}
