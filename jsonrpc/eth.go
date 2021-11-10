package jsonrpc

import (
	"github.com/ethereum/go-ethereum/common"
)

// Eth is the eth jsonrpc endpoint
type Eth struct{}

// BlockNumber returns current block number
func (e *Eth) BlockNumber() (interface{}, error) {
	panic("Not implemented yet")
}

// ChainId returns the chain id of the client
func (e *Eth) ChainId() (interface{}, error) {
	panic("Not implemented yet")
}

// GetBalance returns the account's balance at the referenced block
func (e *Eth) GetBalance(address common.Address, number *int64) (interface{}, error) {
	panic("Not implemented yet")
}

// GetBlockByHash returns information about a block by hash
func (e *Eth) GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, error) {
	panic("Not implemented yet")
}

// GetBlockByNumber returns information about a block by block number
func (e *Eth) GetBlockByNumber(number int64, fullTx bool) (interface{}, error) {
	panic("Not implemented yet")
}

////getTransactionByBlockHashAndIndex
////getTransactionByBlockNumberAndIndex

// GetTransactionByHash returns a transaction by his hash
func (e *Eth) GetTransactionByHash(hash common.Hash) (interface{}, error) {
	panic("Not implemented yet")
}

// GetTransactionCount returns account nonce
func (e *Eth) GetTransactionCount(address common.Address, number *int64) (interface{}, error) {
	panic("Not implemented yet")
}

// GetTransactionReceipt returns a transaction receipt by his hash
func (e *Eth) GetTransactionReceipt(hash common.Hash) (interface{}, error) {
	panic("Not implemented yet")
}

// SendRawTransaction sends a raw transaction
func (e *Eth) SendRawTransaction(input string) (interface{}, error) {
	panic("Not implemented yet")
}
