package jsonrpc

import (
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
)

// Net contains implementations for the "net" RPC endpoints
type Net struct {
	cfg Config
}

// Version returns the current network id
func (n *Net) Version() (interface{}, rpcError) {
	return strconv.FormatUint(n.cfg.ChainID, encoding.Base10), nil
}
