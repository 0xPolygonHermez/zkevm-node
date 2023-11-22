package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// TxPoolEndpoints is the txpool jsonrpc endpoint
type TxPoolEndpoints struct{
	pool types.PoolInterface
}

type contentResponse struct {
	Pending map[common.Address]map[uint64]*txPoolTransaction `json:"pending"`
	Queued  map[common.Address]map[uint64]*txPoolTransaction `json:"queued"`
}

type txPoolTransaction struct {
	Nonce       types.ArgUint64 `json:"nonce"`
	GasPrice    types.ArgBig    `json:"gasPrice"`
	Gas         types.ArgUint64 `json:"gas"`
	To          *common.Address `json:"to"`
	Value       types.ArgBig    `json:"value"`
	Input       types.ArgBytes  `json:"input"`
	Hash        common.Hash     `json:"hash"`
	From        common.Address  `json:"from"`
	BlockHash   common.Hash     `json:"blockHash"`
	BlockNumber interface{}     `json:"blockNumber"`
	TxIndex     interface{}     `json:"transactionIndex"`
}

// NewTxPoolEndpoints creates an new instance of Eth
func NewTxPoolEndpoints(p types.PoolInterface) *TxPoolEndpoints {
	return &TxPoolEndpoints{pool: p}
}

// Content creates a response for txpool_content request.
// See https://geth.ethereum.org/docs/rpc/ns-txpool#txpool_content.
func (e *TxPoolEndpoints) Content() (interface{}, types.Error) {
	resp := contentResponse{
		Pending: make(map[common.Address]map[uint64]*txPoolTransaction),
		Queued:  make(map[common.Address]map[uint64]*txPoolTransaction),
	}

	return resp, nil
}

// Status creates a response for txpool_status request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-txpool#txpool-status
func (e *TxPoolEndpoints) Status() (interface{}, types.Error) {
    ctx := context.Background()
    txPendingCount, err := e.pool.CountPendingTransactions(ctx)
    if err != nil {
        log.Errorf("Failed to count pending txs from pool", err)
		return RPCErrorResponse(types.DefaultErrorCode, "Failed to count pending txs from pool", err, false)
    }

    txQueuedCount, err := e.pool.CountQueuedTransactions(ctx)
    if err != nil {
        log.Errorf("Failed to count queued txs from pool", err)
		return RPCErrorResponse(types.DefaultErrorCode, "Failed to count queued txs from pool", err, false)
    }

	resp := map[string]hexutil.Uint{
		"pending": hexutil.Uint(txPendingCount),
		"queued":  hexutil.Uint(txQueuedCount),
	}

	return resp, nil
}
