package jsonrpc

import (
	"strconv"

	"github.com/hermeznetwork/hermez-core/encoding"
)

// Net contains implementations for the "net" RPC endpoints
type Net struct {
	chainIDSelector *chainIDSelector
}

// Version returns the current network id
func (n *Net) Version() (interface{}, error) {
	chainID, err := n.chainIDSelector.getChainID()
	if err != nil {
		return nil, err
	}

	return strconv.FormatInt(int64(chainID), encoding.Base10), nil
}
