package jsonrpc

import (
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

// NetEndpoints contains implementations for the "net" RPC endpoints
type NetEndpoints struct {
	cfg     Config
	chainID uint64
}

// NewNetEndpoints returns NetEndpoints
func NewNetEndpoints(cfg Config, chainID uint64) *NetEndpoints {
	return &NetEndpoints{
		cfg:     cfg,
		chainID: chainID,
	}
}

// Version returns the current network id
func (n *NetEndpoints) Version() (interface{}, types.Error) {
	return strconv.FormatUint(n.chainID, encoding.Base10), nil
}
