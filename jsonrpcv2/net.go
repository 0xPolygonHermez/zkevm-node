package jsonrpcv2

import (
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
)

// Net contains implementations for the "net" RPC endpoints
type Net struct {
	chainID uint64
}

// Version returns the current network id
func (n *Net) Version() (interface{}, rpcError) {
	return strconv.FormatUint(n.chainID, encoding.Base10), nil
}
