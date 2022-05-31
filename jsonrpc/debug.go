package jsonrpc

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

// Debug is the debug jsonrpc endpoint
type Debug struct {
	state stateInterface
}

type traceTransactionResponse struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue string         `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

type StructLogRes struct {
	Pc            uint64             `json:"pc"`
	Op            string             `json:"op"`
	Gas           uint64             `json:"gas"`
	GasCost       uint64             `json:"gasCost"`
	Depth         int                `json:"depth"`
	Error         string             `json:"error,omitempty"`
	Stack         *[]string          `json:"stack,omitempty"`
	Memory        *[]string          `json:"memory,omitempty"`
	Storage       *map[string]string `json:"storage,omitempty"`
	RefundCounter uint64             `json:"refund,omitempty"`
}

// TraceTransaction creates a response for debug_traceTransaction request.
// See https://geth.ethereum.org/docs/rpc/ns-debug#debug_tracetransaction
func (d *Debug) TraceTransaction(hash common.Hash) (interface{}, error) {
	ctx := context.Background()

	// tx, err := d.state.GetTransactionByHash(ctx, hash, "")
	// if errors.Is(err, state.ErrNotFound) {
	// 	return newGenericError("transaction not found", -32000), nil
	// }

	rcpt, err := d.state.GetTransactionReceipt(ctx, hash, "")
	if errors.Is(err, state.ErrNotFound) {
		return newRPCError(defaultErrorCode, "transaction receipt not found"), nil
	}

	failed := false
	returnValue := ""
	if rcpt.Status == types.ReceiptStatusFailed {
		failed = true
		// this is supposed to be an object value that will be parsed by blockscout
		// to identify the details of the transaction execution.
		// When this field is different from nil, it make the blockscout UI to show
		// the transaction as Failed, but as long as we don't know exactly the object
		// we need to return here, this is causing parsing errors on blockscout, which
		// is causing a side effect to send multiple requests to the core, retrying to
		// get this information.
		// TODO: we need to figure out what to return here based on geth code.
		returnValue = "{}"
	}

	resp := traceTransactionResponse{
		Gas:         rcpt.GasUsed,
		Failed:      failed,
		ReturnValue: returnValue,
		StructLogs:  []StructLogRes{},
	}

	return resp, nil
}
