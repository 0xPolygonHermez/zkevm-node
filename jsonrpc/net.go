package jsonrpc

import "github.com/hermeznetwork/hermez-core/jsonrpc/hex"

// Net contains implementations for the "net" RPC endpoints
type Net struct {
	chainID uint64
}

// Version returns the current network id
func (n *Net) Version() (interface{}, error) {
	return hex.EncodeUint64(n.chainID), nil
}
