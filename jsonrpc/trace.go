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
	TraceAddress []uint64      `json:"traceAddress"`
	SubTraces    uint64        `json:"subtraces"`
	Action       txTraceAction `json:"action"`
	Result       txTraceResult `json:"result"`
	Type         string        `json:"type"`
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
	CallType string    `json:"callType"`
	From     string    `json:"from"`
	Gas      argUint64 `json:"gas"`
	Input    argUint64 `json:"input"`
	To       string    `json:"to"`
	Value    argUint64 `json:"value"`
}

type txTraceResult struct {
	GasUsed argUint64 `json:"gasUsed"`
	Output  argUint64 `json:"output"`
}

// ReplayTransaction creates a response for trace_replayTransaction request.
// See https://openethereum.github.io/JSONRPC-trace-module#trace_replayblocktransactions
func (d *Trace) ReplayBlockTransactions(number *BlockNumber, traceType []string) (interface{}, error) {
	ctx := context.Background()

	batchNumber, err := number.getNumericBlockNumber(ctx, d.state)
	if err != nil {
		return nil, err
	}

	results := d.state.ReplayBatchTransactions(batchNumber)

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
	return replayTransactionResponse{
		// Output: ,
		// StateDiff: ,
		// Trace: []txTrace{
		// 	Action: ,
		// 	Result: ,
		// 	SubTraces: ,
		// 	TraceAddress: ,
		// 	Type: ,
		// },
		// VMTrace: ,
	}
}
