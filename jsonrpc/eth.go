package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
)

// Eth contains implementations for the "eth" RPC endpoints
type Eth struct {
	chainIDSelector  *chainIDSelector
	pool             jsonRPCTxPool
	state            stateInterface
	sequencerAddress common.Address
	gpe              gasPriceEstimator
	storage          storageInterface
}

type blockNumberOrHash struct {
	BlockNumber *BlockNumber `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash `json:"blockHash,omitempty"`
}

// BlockNumber returns current block number
func (e *Eth) BlockNumber() (interface{}, error) {
	ctx := context.Background()

	lastBatchNumber, err := e.state.GetLastBatchNumber(ctx, "")
	if err != nil {
		return nil, err
	}

	return hex.EncodeUint64(lastBatchNumber), nil
}

// Call executes a new message call immediately and returns the value of
// executed contract and potential error.
// Note, this function doesn't make any changes in the state/blockchain and is
// useful to execute view/pure methods and retrieve values.
func (e *Eth) Call(arg *txnArgs, number *BlockNumber) (interface{}, error) {
	// If the caller didn't supply the gas limit in the message, then we set it to maximum possible => block gas limit
	if arg.Gas == nil || *arg.Gas == argUint64(0) {
		filter := blockNumberOrHash{
			BlockNumber: number,
		}

		header, err := e.getHeaderFromBlockNumberOrHash(&filter)
		if err != nil {
			log.Errorf("failed to get header from block hash or block number")
			return "0x", nil
		}

		gas := argUint64(header.GasLimit)
		arg.Gas = &gas
	}

	if arg.From == nil {
		from := state.ZeroAddress
		arg.From = &from
	}

	tx := arg.ToTransaction()

	ctx := context.Background()

	batchNumber, err := e.getNumericBlockNumber(ctx, *number)
	if err != nil {
		return "0x", nil
	}

	batch, err := e.state.GetBatchByNumber(ctx, batchNumber, "")
	if err != nil {
		return "0x", nil
	}

	bp, err := e.state.NewBatchProcessor(ctx, e.sequencerAddress, batch.Header.Root[:], "")
	if err != nil {
		return "0x", nil
	}

	result := bp.ProcessUnsignedTransaction(ctx, tx, *arg.From, e.sequencerAddress)
	if result.Failed() {
		log.Errorf("unable to execute call: %s", result.Err.Error())
		return "0x", nil
	}

	return argBytesPtr(result.ReturnValue), nil
}

// ChainId returns the chain id of the client
func (e *Eth) ChainId() (interface{}, error) { //nolint:revive
	chainID, err := e.chainIDSelector.getChainID()
	if err != nil {
		return nil, err
	}

	return hex.EncodeUint64(chainID), nil
}

// EstimateGas generates and returns an estimate of how much gas is necessary to
// allow the transaction to complete.
// The transaction will not be added to the blockchain.
// Note that the estimate may be significantly more than the amount of gas actually
// used by the transaction, for a variety of reasons including EVM mechanics and
// node performance.
func (e *Eth) EstimateGas(arg *txnArgs, rawNum *BlockNumber) (interface{}, error) {
	tx := arg.ToTransaction()

	if arg.From == nil {
		from := state.ZeroAddress
		arg.From = &from
	}

	gasEstimation, err := e.state.EstimateGas(tx, *arg.From)
	return hex.EncodeUint64(gasEstimation), err
}

// GasPrice returns the average gas price based on the last x blocks
func (e *Eth) GasPrice() (interface{}, error) {
	ctx := context.Background()
	gasPrice, err := e.gpe.GetAvgGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	if gasPrice != nil {
		return hex.EncodeUint64(gasPrice.Uint64()), nil
	}
	return hex.EncodeUint64(0), nil
}

// GetBalance returns the account's balance at the referenced block
func (e *Eth) GetBalance(address common.Address, number *BlockNumber) (interface{}, error) {
	ctx := context.Background()
	batchNumber, err := e.getNumericBlockNumber(ctx, *number)
	if err != nil {
		return nil, err
	}

	balance, err := e.state.GetBalance(ctx, address, batchNumber, "")
	if errors.Is(err, state.ErrNotFound) {
		return hex.EncodeUint64(0), nil
	} else if err != nil {
		return nil, err
	}

	return hex.EncodeBig(balance), nil
}

// GetBlockByHash returns information about a block by hash
func (e *Eth) GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, error) {
	ctx := context.Background()

	batch, err := e.state.GetBatchByHash(ctx, hash, "")
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	block := batchToRPCBlock(batch, fullTx)

	return block, nil
}

// GetBlockByNumber returns information about a block by block number
func (e *Eth) GetBlockByNumber(number BlockNumber, fullTx bool) (interface{}, error) {
	ctx := context.Background()

	if number == PendingBlockNumber {
		lastBatch, err := e.state.GetLastBatch(context.Background(), true, "")
		if err != nil {
			return nil, err
		}
		header := &types.Header{
			ParentHash: lastBatch.Hash(),
			Number:     big.NewInt(0).SetUint64(lastBatch.Number().Uint64() + 1),
			Difficulty: big.NewInt(0),
		}
		batch := &state.Batch{Header: header}
		block := batchToRPCBlock(batch, fullTx)

		return block, nil
	}

	batchNumber, err := e.getNumericBlockNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	batch, err := e.state.GetBatchByNumber(ctx, batchNumber, "")
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	block := batchToRPCBlock(batch, fullTx)

	return block, nil
}

// GetCode returns account code at given block number
func (e *Eth) GetCode(address common.Address, number *BlockNumber) (interface{}, error) {
	ctx := context.Background()

	batchNumber, err := e.getNumericBlockNumber(ctx, *number)
	if err != nil {
		return nil, err
	}

	code, err := e.state.GetCode(ctx, address, batchNumber, "")
	if errors.Is(err, state.ErrNotFound) {
		return "0x", nil
	} else if err != nil {
		return nil, err
	}

	return argBytes(code), nil
}

// GetFilterChanges polling method for a filter, which returns
// an array of logs which occurred since last poll.
func (e *Eth) GetFilterChanges(filterID argUint64) (interface{}, error) {
	filter, err := e.storage.GetFilter(uint64(filterID))
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	err = e.storage.UpdateFilterLastPoll(filter.ID)
	if err != nil {
		return nil, err
	}

	switch filter.Type {
	case FilterTypeBlock:
		{
			res, err := e.state.GetBatchHashesSince(context.Background(), filter.LastPoll, "")
			if err != nil {
				return nil, err
			}
			if len(res) == 0 {
				return nil, nil
			}
			return res, nil
		}
	case FilterTypePendingTx:
		{
			res, err := e.pool.GetPendingTxHashesSince(context.Background(), filter.LastPoll)
			if err != nil {
				return nil, err
			}
			if len(res) == 0 {
				return nil, nil
			}
			return res, nil
		}
	case FilterTypeLog:
		{
			filterParameters := &LogFilter{}
			err = json.Unmarshal([]byte(filter.Parameters), filterParameters)
			if err != nil {
				return nil, err
			}

			filterParameters.Since = &filter.LastPoll

			resInterface, err := e.GetLogs(filterParameters)
			if err != nil {
				return nil, err
			}
			res := resInterface.([]rpcLog)
			if len(res) == 0 {
				return nil, nil
			}
			return res, nil
		}
	default:
		return nil, nil
	}

}

// GetFilterLogs returns an array of all logs matching filter
// with given id.
func (e *Eth) GetFilterLogs(filterID argUint64) (interface{}, error) {
	filter, err := e.storage.GetFilter(uint64(filterID))
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if filter.Type != FilterTypeLog {
		return nil, nil
	}

	filterParameters := &LogFilter{}
	err = json.Unmarshal([]byte(filter.Parameters), filterParameters)
	if err != nil {
		return nil, err
	}

	filterParameters.Since = nil

	return e.GetLogs(filterParameters)
}

// GetLogs returns a list of logs accordingly to the provided filter
func (e *Eth) GetLogs(filter *LogFilter) (interface{}, error) {
	ctx := context.Background()

	fromBlock, err := e.getNumericBlockNumber(ctx, filter.FromBlock)
	if err != nil {
		return nil, err
	}

	toBlock, err := e.getNumericBlockNumber(ctx, filter.ToBlock)
	if err != nil {
		return nil, err
	}

	logs, err := e.state.GetLogs(ctx, fromBlock, toBlock, filter.Addresses, filter.Topics, filter.BlockHash, filter.Since, "")
	if err != nil {
		return nil, err
	}

	result := make([]rpcLog, 0, len(logs))
	for _, l := range logs {
		result = append(result, logToRPCLog(*l))
	}

	return result, nil
}

// GetStorageAt gets the value stored for an specific address and position
func (e *Eth) GetStorageAt(address common.Address, position common.Hash, number *BlockNumber) (interface{}, error) {
	ctx := context.Background()

	batchNumber, err := e.getNumericBlockNumber(ctx, *number)
	if err != nil {
		return nil, err
	}

	value, err := e.state.GetStorageAt(ctx, address, position.Big(), batchNumber, "")
	if errors.Is(err, state.ErrNotFound) {
		return argBytesPtr(common.Hash{}.Bytes()), nil
	} else if err != nil {
		return nil, err
	}

	return argBytesPtr(common.BigToHash(value).Bytes()), nil
}

// GetTransactionByBlockHashAndIndex returns information about a transaction by
// block hash and transaction index position.
func (e *Eth) GetTransactionByBlockHashAndIndex(hash common.Hash, index Index) (interface{}, error) {
	ctx := context.Background()

	tx, err := e.state.GetTransactionByBatchHashAndIndex(ctx, hash, uint64(index), "")
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), "")
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return toRPCTransaction(tx, receipt.BlockNumber, receipt.BlockHash, uint64(receipt.TransactionIndex)), nil
}

// GetTransactionByBlockNumberAndIndex returns information about a transaction by
// block number and transaction index position.
func (e *Eth) GetTransactionByBlockNumberAndIndex(number *BlockNumber, index Index) (interface{}, error) {
	ctx := context.Background()

	batchNumber, err := e.getNumericBlockNumber(ctx, *number)
	if err != nil {
		return nil, err
	}

	tx, err := e.state.GetTransactionByBatchNumberAndIndex(ctx, batchNumber, uint64(index), "")
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), "")
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return toRPCTransaction(tx, receipt.BlockNumber, receipt.BlockHash, uint64(receipt.TransactionIndex)), nil
}

// GetTransactionByHash returns a transaction by his hash
func (e *Eth) GetTransactionByHash(hash common.Hash) (interface{}, error) {
	ctx := context.Background()

	tx, err := e.state.GetTransactionByHash(ctx, hash, "")
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), "")
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return toRPCTransaction(tx, receipt.BlockNumber, receipt.BlockHash, uint64(receipt.TransactionIndex)), nil
}

// GetTransactionCount returns account nonce
func (e *Eth) GetTransactionCount(address common.Address, number *BlockNumber) (interface{}, error) {
	ctx := context.Background()
	batchNumber, err := e.getNumericBlockNumber(ctx, *number)
	if err != nil {
		return nil, err
	}

	nonce, err := e.state.GetNonce(ctx, address, batchNumber, "")
	if errors.Is(err, state.ErrNotFound) {
		return hex.EncodeUint64(0), nil
	} else if err != nil {
		return nil, err
	}

	return hex.EncodeUint64(nonce), nil
}

// GetTransactionReceipt returns a transaction receipt by his hash
func (e *Eth) GetTransactionReceipt(hash common.Hash) (interface{}, error) {
	ctx := context.Background()

	r, err := e.state.GetTransactionReceipt(ctx, hash, "")
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return stateReceiptToRPCReceipt(r), nil
}

// NewBlockFilter creates a filter in the node, to notify when
// a new block arrives. To check if the state has changed,
// call eth_getFilterChanges.
func (e *Eth) NewBlockFilter() (interface{}, error) {
	id, err := e.storage.NewBlockFilter()
	if err != nil {
		return nil, err
	}

	return argUint64(id), nil
}

// NewFilter creates a filter object, based on filter options,
// to notify when the state changes (logs). To check if the state
// has changed, call eth_getFilterChanges.
func (e *Eth) NewFilter(filter *LogFilter) (interface{}, error) {
	id, err := e.storage.NewLogFilter(*filter)
	if err != nil {
		return nil, err
	}

	return argUint64(id), nil
}

// NewPendingTransactionFilter creates a filter in the node, to
// notify when new pending transactions arrive. To check if the
// state has changed, call eth_getFilterChanges.
func (e *Eth) NewPendingTransactionFilter(filterID argUint64) (interface{}, error) {
	id, err := e.storage.NewPendingTransactionFilter()
	if err != nil {
		return nil, err
	}

	return argUint64(id), nil
}

// SendRawTransaction sends a raw transaction
func (e *Eth) SendRawTransaction(input string) (interface{}, error) {
	tx, err := hexToTx(input)
	if err != nil {
		log.Warnf("Invalid tx: %v", err)
		return nil, err
	}

	log.Debugf("adding TX to the pool: %v", tx.Hash().Hex())
	if err := e.pool.AddTx(context.Background(), *tx); err != nil {
		log.Warnf("Failed to add TX to the pool[%v]: %v", tx.Hash().Hex(), err)
		return nil, err
	}
	log.Debugf("TX added to the pool: %v", tx.Hash().Hex())

	return tx.Hash().Hex(), nil
}

// UninstallFilter uninstalls a filter with given id. Should
// always be called when watch is no longer needed. Additionally
// Filters timeout when they aren’t requested with
// eth_getFilterChanges for a period of time.
func (e *Eth) UninstallFilter(filterID argUint64) (interface{}, error) {
	return e.storage.UninstallFilter(uint64(filterID))
}

func (e *Eth) getNumericBlockNumber(ctx context.Context, number BlockNumber) (uint64, error) {
	switch number {
	case LatestBlockNumber, PendingBlockNumber:
		lastBatchNumber, err := e.state.GetLastBatchNumber(ctx, "")
		if err != nil {
			return 0, err
		}

		return lastBatchNumber, nil

	case EarliestBlockNumber:
		return 0, nil

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

func (e *Eth) getHeaderFromBlockNumberOrHash(bnh *blockNumberOrHash) (*types.Header, error) {
	var (
		header *types.Header
		err    error
	)

	if bnh.BlockNumber != nil {
		header, err = e.getBatchHeader(*bnh.BlockNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get the header of block %d: %s", *bnh.BlockNumber, err.Error())
		}
	} else if bnh.BlockHash != nil {
		block, err := e.state.GetBatchByHash(context.Background(), *bnh.BlockHash, "")
		if err != nil {
			return nil, fmt.Errorf("could not find block referenced by the hash %s, err: %v", bnh.BlockHash.String(), err)
		}

		header = block.Header
	}

	return header, nil
}

func (e *Eth) getBatchHeader(number BlockNumber) (*types.Header, error) {
	switch number {
	case LatestBlockNumber:
		batch, err := e.state.GetLastBatch(context.Background(), false, "")
		if err != nil {
			return nil, err
		}
		return batch.Header, nil

	case EarliestBlockNumber:
		batch, err := e.state.GetBatchByNumber(context.Background(), 0, "")
		if err != nil {
			return nil, err
		}
		return batch.Header, nil

	case PendingBlockNumber:
		lastBatch, err := e.state.GetLastBatch(context.Background(), true, "")
		if err != nil {
			return nil, err
		}
		header := &types.Header{
			ParentHash: lastBatch.Hash(),
			Number:     big.NewInt(0).SetUint64(lastBatch.Number().Uint64() + 1),
			Difficulty: big.NewInt(0),
		}
		return header, nil

	default:
		return e.state.GetBatchHeader(context.Background(), uint64(number), "")
	}
}
