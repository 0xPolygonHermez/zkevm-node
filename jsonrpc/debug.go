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
	state state.State
}

type traceTransactionResponse struct {
	Gas         uint64        `json:"gas"`
	Failed      bool          `json:"failed"`
	ReturnValue string        `json:"returnValue"`
	StructLogs  []interface{} `json:"structLogs"`
}

// TraceTransaction creates a response for debug_traceTransaction request.
// See https://geth.ethereum.org/docs/rpc/ns-debug#debug_tracetransaction
func (d *Debug) TraceTransaction(hash common.Hash) (interface{}, error) {
	ctx := context.Background()

	tx, err := d.state.GetTransactionByHash(ctx, hash)
	if errors.Is(err, state.ErrNotFound) {
		return genesisIsNotTraceableError{}, nil
	}

	rcpt, err := d.state.GetTransactionReceipt(ctx, hash)
	if errors.Is(err, state.ErrNotFound) {
		return genesisIsNotTraceableError{}, nil
	}

	resp := traceTransactionResponse{
		Gas:         tx.Gas(),
		Failed:      rcpt.Status == types.ReceiptStatusFailed,
		ReturnValue: "",
		StructLogs:  []interface{}{},
	}

	return resp, nil
}
