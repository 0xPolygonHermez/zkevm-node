package jsonrpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state"
)

// Debug is the debug jsonrpc endpoint
type Debug struct {
	state state.State
}

// Create response for debug_traceTransaction request.
// See https://geth.ethereum.org/docs/rpc/ns-debug#debug_tracetransaction
func (t *TxPool) TraceTransaction(hash common.Hash) (interface{}, error) {
	return nil, nil
}
