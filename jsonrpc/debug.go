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
	Gas         uint64        `json:"gas"`
	Failed      bool          `json:"failed"`
	ReturnValue string        `json:"returnValue"`
	StructLogs  []interface{} `json:"structLogs"`
}

type structLog struct {
}

// TraceTransaction creates a response for debug_traceTransaction request.
// See https://geth.ethereum.org/docs/rpc/ns-debug#debug_tracetransaction
func (d *Debug) TraceTransaction(hash common.Hash) (interface{}, error) {
	ctx := context.Background()

	tx, err := d.state.GetTransactionByHash(ctx, hash)
	if errors.Is(err, state.ErrNotFound) {
		return genesisIsNotTraceableError{}, nil
	}

	txStructLogs, err := d.state.TraceTransaction(hash)
	if err != nil {
		return nil, err
	}

	rcpt, err := d.state.GetTransactionReceipt(ctx, hash)
	if errors.Is(err, state.ErrNotFound) {
		return genesisIsNotTraceableError{}, nil
	}

	structLogs := make([]interface{}, 0, len(txStructLogs))
	for _, txStructLog := range txStructLogs {
		structLogs = append(structLogs, txStructLog)
	}

	resp := traceTransactionResponse{
		Gas:         tx.Gas(),
		Failed:      rcpt.Status == types.ReceiptStatusFailed,
		ReturnValue: "",
		StructLogs:  structLogs,
	}

	return resp, nil
}
