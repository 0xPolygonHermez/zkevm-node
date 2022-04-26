package jsonrpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state/runtime"
)

// Trace is the trace jsonrpc endpoint
type Trace struct {
	state stateInterface
}

type replayTransactionResponse struct {
	Output    *argBytes                        `json:"output"`
	Trace     []txTrace                        `json:"trace"`
	VMTrace   []txVMTrace                      `json:"vmTrace"`
	StateDiff map[common.Address]txAccountDiff `json:"stateDiff"`
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
	Code       []byte               `json:"code"`
	Operations []txVMTraceOperation `json:"ops"`
}

type txAccountDiff struct{}

type txVMTraceOperation struct {
	PC                uint64                      `json:"pc"`
	Cost              uint64                      `json:"cost"`
	ExecutedOperation *txVMTraceExecutedOperation `json:"ex"`
	Sub               *txVMTrace                  `json:"sub"`
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
func (d *Trace) ReplayBlockTransactions(number *BlockNumber, traceType []string) (interface{}, error) {
	ctx := context.Background()

	batchNumber, err := number.getNumericBlockNumber(ctx, d.state)
	if err != nil {
		return nil, err
	}

	results, err := d.state.ReplayBatchTransactions(batchNumber)
	if err != nil {
		return nil, err
	}

	response := make([]replayTransactionResponse, 0, len(results))

	for _, result := range results {
		response = append(response, d.executionResultToReplayTransactionResponse(result))
	}

	return response, nil
}

// ReplayTransaction creates a response for trace_replayTransaction request.
// See https://openethereum.github.io/JSONRPC-trace-module#trace_replaytransaction
func (d *Trace) ReplayTransaction(hash common.Hash, traceType []string) (interface{}, error) {
	result := d.state.ReplayTransaction(hash)
	response := d.executionResultToReplayTransactionResponse(result)

	return response, nil
}

func (d *Trace) executionResultToReplayTransactionResponse(result *runtime.ExecutionResult) replayTransactionResponse {
	// operations := []txVMTraceOperation{}

	return replayTransactionResponse{
		// Output: ,
		// StateDiff: ,
		// Trace: ,
		VMTrace: []txVMTrace{
			// Code:       result.VMTrace.Code,
			// Operations: operations,
		},
	}
}
