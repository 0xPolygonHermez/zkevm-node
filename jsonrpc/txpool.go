package jsonrpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state/helper"
)

// TxPool is the txpool jsonrpc endpoint
type TxPool struct {
	pool pool.Pool
}

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
func (t *TxPool) Content() (interface{}, error) {
	ctx := context.Background()

	pendingTxs, err := t.pool.GetPendingTxs(ctx)
	if err != nil {
		return nil, err
	}

	// collect pending
	pendingRPCTxs := make(map[common.Address]map[uint64]*txPoolTransaction, len(pendingTxs))

	for _, pendingTx := range pendingTxs {
		sender, err := helper.GetSender(&pendingTx.Transaction)
		if err != nil {
			return nil, err
		}

		if _, found := pendingRPCTxs[sender]; !found {
			pendingRPCTxs[sender] = make(map[uint64]*txPoolTransaction)
		}

		pendingRPCTxs[sender][pendingTx.Nonce()] = toTxPoolTransaction(sender, &pendingTx.Transaction)
	}

	resp := contentResponse{
		Pending: pendingRPCTxs,
		Queued:  make(map[common.Address]map[uint64]*txPoolTransaction),
	}

	return resp, nil
}

func toTxPoolTransaction(sender common.Address, t *types.Transaction) *txPoolTransaction {
	return &txPoolTransaction{
		Nonce:       argUint64(t.Nonce()),
		GasPrice:    argBig(*t.GasPrice()),
		Gas:         argUint64(t.Gas()),
		To:          t.To(),
		Value:       argBig(*t.Value()),
		Input:       argBytes(t.Data()),
		Hash:        t.Hash(),
		From:        sender,
		BlockHash:   common.Hash{},
		BlockNumber: nil,
		TxIndex:     nil,
	}
}
