package jsonrpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state"
)

// Debug is the debug jsonrpc endpoint
type Debug struct {
	state state.State
}

// TraceTransaction creates a response for debug_traceTransaction request.
// See https://geth.ethereum.org/docs/rpc/ns-debug#debug_tracetransaction
func (d *Debug) TraceTransaction(hash common.Hash) (interface{}, error) {
	return nil, nil
}
