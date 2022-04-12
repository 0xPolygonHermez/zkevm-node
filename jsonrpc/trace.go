package jsonrpc

import (
	"github.com/ethereum/go-ethereum/common"
)

// Trace is the trace jsonrpc endpoint
type Trace struct {
	state stateInterface
}

type replayTransactionResponse struct {
	Output    argUint64   `json:"output"`
	StateDiff interface{} `json:"stateDiff"`
	Trace     []txTrace   `json:"trace"`
	VMTrace   interface{} `json:"vmTrace"`
}

type txTrace struct {
	Action       txTraceAction `json:"action"`
	Result       txTraceResult `json:"result"`
	SubTraces    uint          `json:"subtraces"`
	TraceAddress []interface{} `json:"traceAddress"`
	Type         string        `json:"type"`
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
// See https://openethereum.github.io/JSONRPC-trace-module#trace_replaytransaction
func (d *Trace) ReplayTransaction(hash common.Hash, traceType []string) (interface{}, error) {
	d.state.ReplayTransaction(hash)

	//TODO: translate state response to rpc response

	return replayTransactionResponse{}, nil
}
