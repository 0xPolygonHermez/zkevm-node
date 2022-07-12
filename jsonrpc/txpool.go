package jsonrpc

import (
	"github.com/ethereum/go-ethereum/common"
)

// TxPool is the txpool jsonrpc endpoint
type TxPool struct{}

type contentResponse struct {
	Pending map[common.Address]map[uint64]*txPoolTransaction `json:"pending"`
	Queued  map[common.Address]map[uint64]*txPoolTransaction `json:"queued"`
}

type txPoolTransaction struct {
	Nonce       argUint64       `json:"nonce"`
	GasPrice    argBig          `json:"gasPrice"`
	Gas         argUint64       `json:"gas"`
	To          *common.Address `json:"to"`
	Value       argBig          `json:"value"`
	Input       argBytes        `json:"input"`
	Hash        common.Hash     `json:"hash"`
	From        common.Address  `json:"from"`
	BlockHash   common.Hash     `json:"blockHash"`
	BlockNumber interface{}     `json:"blockNumber"`
	TxIndex     interface{}     `json:"transactionIndex"`
}

// Content creates a response for txpool_content request.
// See https://geth.ethereum.org/docs/rpc/ns-txpool#txpool_content.
func (t *TxPool) Content() (interface{}, rpcError) {
	resp := contentResponse{
		Pending: make(map[common.Address]map[uint64]*txPoolTransaction),
		Queued:  make(map[common.Address]map[uint64]*txPoolTransaction),
	}

	return resp, nil
}
