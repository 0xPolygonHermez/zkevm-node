package jsonrpc

import (
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
)

// NetEndpoints contains implementations for the "net" RPC endpoints
type NetEndpoints struct {
	cfg Config
}

// Version returns the current network id
func (n *NetEndpoints) Version() (interface{}, rpcError) {
	return strconv.FormatUint(n.cfg.ChainID, encoding.Base10), nil
}
