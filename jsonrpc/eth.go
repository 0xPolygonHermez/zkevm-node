package jsonrpc

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// Eth contains implementations for the "eth" RPC endpoints
type Eth struct {
	chainID uint64
	pool    pool.Pool
	state   state.State
}

// BlockNumber returns current block number
func (e *Eth) BlockNumber() (interface{}, error) {
	lastBatchNumber, err := e.state.GetLastBatchNumber()
	if err != nil {
		return nil, err
	}

	return hex.EncodeUint64(lastBatchNumber), nil
}

// ChainId returns the chain id of the client
func (e *Eth) ChainId() (interface{}, error) { //nolint:golint
	return hex.EncodeUint64(e.chainID), nil
}

// EstimateGas generates and returns an estimate of how much gas is necessary to
// allow the transaction to complete.
// The transaction will not be added to the blockchain.
// Note that the estimate may be significantly more than the amount of gas actually
// used by the transaction, for a variety of reasons including EVM mechanics and
// node performance.
func (e *Eth) EstimateGas(arg *txnArgs, rawNum *BlockNumber) (interface{}, error) {
	tx := arg.ToTransaction()
	gasEstimation := e.state.EstimateGas(*tx)
	return hex.EncodeUint64(gasEstimation), nil
}

// GasPrice returns the average gas price based on the last x blocks
func (e *Eth) GasPrice() (interface{}, error) {
	gasPrice, err := e.pool.GetGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	return hex.EncodeUint64(gasPrice), nil
}

// GetBalance returns the account's balance at the referenced block
func (e *Eth) GetBalance(address common.Address, number *BlockNumber) (interface{}, error) {
	batchNumber, err := getNumericBlockNumber(e, *number)
	if err != nil {
		return nil, err
	}

	balance, err := e.state.GetBalance(address, batchNumber)
	if err != nil {
		return nil, err
	}

	return hex.EncodeUint64(balance.Uint64()), nil
}

// GetBlockByHash returns information about a block by hash
func (e *Eth) GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, error) {
	batch, err := e.state.GetBatchByHash(hash)
	if err != nil {
		return nil, err
	}

	return batch, nil
}

// GetBlockByNumber returns information about a block by block number
func (e *Eth) GetBlockByNumber(number BlockNumber, fullTx bool) (interface{}, error) {
	batchNumber, err := getNumericBlockNumber(e, number)
	if err != nil {
		return nil, err
	}

	batch, err := e.state.GetBatchByNumber(batchNumber)
	if err != nil {
		return nil, err
	}

	return batch, nil
}

// GetCode returns account code at given block number
func (e *Eth) GetCode(address common.Address, number *BlockNumber) (interface{}, error) {
	// we need this because Metamask is calling this method when a transfer is executed.
	return "0x", nil
}

// GetTransactionByBlockHashAndIndex returns information about a transaction by
// block hash and transaction index position.
func (e *Eth) GetTransactionByBlockHashAndIndex(hash common.Hash, index Index) (interface{}, error) {
	tx, err := e.state.GetTransactionByBatchHashAndIndex(hash, uint64(index))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionByBlockNumberAndIndex returns information about a transaction by
// block number and transaction index position.
func (e *Eth) GetTransactionByBlockNumberAndIndex(number *BlockNumber, index Index) (interface{}, error) {
	batchNumber, err := getNumericBlockNumber(e, *number)
	if err != nil {
		return nil, err
	}

	tx, err := e.state.GetTransactionByBatchNumberAndIndex(batchNumber, uint64(index))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionByHash returns a transaction by his hash
func (e *Eth) GetTransactionByHash(hash common.Hash) (interface{}, error) {
	tx, err := e.state.GetTransaction(hash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionCount returns account nonce
func (e *Eth) GetTransactionCount(address common.Address, number *BlockNumber) (interface{}, error) {
	batchNumber, err := getNumericBlockNumber(e, *number)
	if err != nil {
		return nil, err
	}

	nonce, err := e.state.GetNonce(address, batchNumber)
	if err != nil {
		return nil, err
	}

	return hex.EncodeUint64(nonce), nil
}

// GetTransactionReceipt returns a transaction receipt by his hash
func (e *Eth) GetTransactionReceipt(hash common.Hash) (interface{}, error) {
	tx, err := e.state.GetTransactionReceipt(hash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SendRawTransaction sends a raw transaction
func (e *Eth) SendRawTransaction(input string) (interface{}, error) {
	tx, err := hexToTx(input)
	if err != nil {
		return nil, err
	}

	if err := e.pool.AddTx(context.Background(), *tx); err != nil {
		return nil, err
	}

	return tx.Hash().Hex(), nil
}

func getNumericBlockNumber(e *Eth, number BlockNumber) (uint64, error) {
	switch number {
	case LatestBlockNumber:
		lastBatchNumber, err := e.state.GetLastBatchNumber()
		if err != nil {
			return 0, err
		}

		return lastBatchNumber, nil
	case EarliestBlockNumber:
		return 0, fmt.Errorf("fetching the earliest header is not supported")

	case PendingBlockNumber:
		return 0, fmt.Errorf("fetching the pending header is not supported")

	default:
		if number < 0 {
			return 0, fmt.Errorf("invalid argument 0: block number larger than int64")
		}
		return uint64(number), nil
	}
}

func hexToTx(str string) (*types.Transaction, error) {
	tx := new(types.Transaction)

	b, err := hex.DecodeHex(str)
	if err != nil {
		return nil, err
	}

	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return tx, nil
}
