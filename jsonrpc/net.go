package jsonrpc

import "github.com/hermeznetwork/hermez-core/hex"

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

	return hex.EncodeUint64(chainID), nil
}
