package jsonrpcv2

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
	state "github.com/hermeznetwork/hermez-core/statev2"
	"github.com/jackc/pgx/v4"
)

// Eth contains implementations for the "eth" RPC endpoints
type Eth struct {
	chainID          uint64
	pool             jsonRPCTxPool
	state            stateInterface
	sequencerAddress common.Address
	gpe              gasPriceEstimator
	storage          storageInterface
	txMan            dbTxManager
}

// BlockNumber returns current block number
func (e *Eth) BlockNumber() (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		lastBlockNumber, err := e.state.GetLastBlockNumber(ctx, dbTx)
		if err != nil {
			return "0x0", newRPCError(defaultErrorCode, "failed to get the last block number from state")
		}

		return hex.EncodeUint64(lastBlockNumber), nil
	})
}

// Call executes a new message call immediately and returns the value of
// executed contract and potential error.
// Note, this function doesn't make any changes in the state/blockchain and is
// useful to execute view/pure methods and retrieve values.
func (e *Eth) Call(arg *txnArgs, number *BlockNumber) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		// If the caller didn't supply the gas limit in the message, then we set it to maximum possible => block gas limit
		if arg.Gas == nil || *arg.Gas == argUint64(0) {
			header, err := e.getBlockHeader(ctx, *number, dbTx)
			if err != nil {
				const errorMessage = "failed to get block header"
				log.Errorf("%v: %v", errorMessage, err)
				return nil, newRPCError(defaultErrorCode, errorMessage)
			}

			gas := argUint64(header.GasLimit)
			arg.Gas = &gas
		}

		tx := arg.ToTransaction()

		var err error
		blockNumber, err := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if err != nil {
			return nil, err
		}

		result := e.state.ProcessUnsignedTransaction(ctx, tx, arg.From, e.sequencerAddress, blockNumber, dbTx)
		if result.Failed() {
			errorMessage := fmt.Sprintf("failed to execute call: %v", result.Err)
			log.Errorf("%v", errorMessage)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return argBytesPtr(result.ReturnValue), nil
	})
}

// ChainId returns the chain id of the client
func (e *Eth) ChainId() (interface{}, error) { //nolint:revive
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

	gasEstimation, err := e.state.EstimateGas(tx, arg.From)
	return hex.EncodeUint64(gasEstimation), err
}

// GasPrice returns the average gas price based on the last x blocks
func (e *Eth) GasPrice() (interface{}, error) {
	ctx := context.Background()
	gasPrice, err := e.gpe.GetAvgGasPrice(ctx)
	if err != nil {
		return "0x0", nil
	}
	if gasPrice != nil {
		return hex.EncodeUint64(gasPrice.Uint64()), nil
	}
	return hex.EncodeUint64(0), nil
}

// GetBalance returns the account's balance at the referenced block
func (e *Eth) GetBalance(address common.Address, number *BlockNumber) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		var err error
		blockNumber, err := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if err != nil {
			return nil, err
		}

		balance, err := e.state.GetBalance(ctx, address, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return hex.EncodeUint64(0), nil
		} else if err != nil {
			log.Errorf("failed to get balance from state: %v", err)
			return nil, newRPCError(defaultErrorCode, "failed to get balance from state")
		}

		return hex.EncodeBig(balance), nil
	})
}

// GetBlockByHash returns information about a block by hash
func (e *Eth) GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		block, err := e.state.GetBlockByHash(ctx, hash, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		rpcBlock := l2BlockToRPCBlock(block, fullTx)

		return rpcBlock, nil
	})
}

// GetBlockByNumber returns information about a block by block number
func (e *Eth) GetBlockByNumber(number BlockNumber, fullTx bool) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		if number == PendingBlockNumber {
			lastBlock, err := e.state.GetLastBlock(ctx, dbTx)
			if err != nil {
				const errorMessage = "couldn't load last block from state to compute the pending block"
				log.Errorf("%v: %v", errorMessage, err)
				return nil, newRPCError(defaultErrorCode, errorMessage)
			}
			header := types.CopyHeader(lastBlock.Header)
			header.ParentHash = lastBlock.Hash()
			header.Number = big.NewInt(0).SetUint64(lastBlock.Number().Uint64() + 1)
			header.TxHash = types.EmptyRootHash
			header.UncleHash = types.EmptyUncleHash
			block := &state.L2Block{Header: header}
			rpcBlock := l2BlockToRPCBlock(block, fullTx)

			return rpcBlock, nil
		}
		var err error
		blockNumber, err := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if err != nil {
			return nil, err
		}

		block, err := e.state.GetBlockByNumber(ctx, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			const errorMessage = "couldn't load block from state by number %v"
			log.Errorf("%v: %v", fmt.Sprintf(errorMessage, blockNumber), err)
			return nil, newRPCError(defaultErrorCode, fmt.Sprintf(errorMessage, blockNumber))
		}

		rpcBlock := l2BlockToRPCBlock(block, fullTx)

		return rpcBlock, nil
	})
}

// GetCode returns account code at given block number
func (e *Eth) GetCode(address common.Address, number *BlockNumber) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		var err error
		blockNumber, err := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if err != nil {
			return nil, err
		}

		code, err := e.state.GetCode(ctx, address, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return "0x", nil
		} else if err != nil {
			const errorMessage = "failed to get code"
			log.Errorf("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return argBytes(code), nil
	})
}

// GetCompilers eth_getCompilers
func (e *Eth) GetCompilers() (interface{}, error) {
	return []interface{}{}, nil
}

// GetFilterChanges polling method for a filter, which returns
// an array of logs which occurred since last poll.
func (e *Eth) GetFilterChanges(filterID argUint64) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
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
				res, err := e.state.GetBlockHashesSince(ctx, filter.LastPoll, dbTx)
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
				res, err := e.pool.GetPendingTxHashesSince(ctx, filter.LastPoll)
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
	})
}

// GetFilterLogs returns an array of all logs mlocking filter
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
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		var err error
		fromBlock, err := filter.FromBlock.getNumericBlockNumber(ctx, e.state, dbTx)
		if err != nil {
			return nil, err
		}

		toBlock, err := filter.ToBlock.getNumericBlockNumber(ctx, e.state, dbTx)
		if err != nil {
			return nil, err
		}

		logs, err := e.state.GetLogs(ctx, fromBlock, toBlock, filter.Addresses, filter.Topics, filter.BlockHash, filter.Since, dbTx)
		if err != nil {
			return nil, err
		}

		result := make([]rpcLog, 0, len(logs))
		for _, l := range logs {
			result = append(result, logToRPCLog(*l))
		}

		return result, nil
	})
}

// GetStorageAt gets the value stored for an specific address and position
func (e *Eth) GetStorageAt(address common.Address, position common.Hash, number *BlockNumber) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		var err error
		blockNumber, err := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if err != nil {
			return nil, err
		}

		value, err := e.state.GetStorageAt(ctx, address, position.Big(), blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return argBytesPtr(common.Hash{}.Bytes()), nil
		} else if err != nil {
			return nil, err
		}

		return argBytesPtr(common.BigToHash(value).Bytes()), nil
	})
}

// GetTransactionByBlockHashAndIndex returns information about a transaction by
// block hash and transaction index position.
func (e *Eth) GetTransactionByBlockHashAndIndex(hash common.Hash, index Index) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		tx, err := e.state.GetTransactionByBlockHashAndIndex(ctx, hash, uint64(index), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			const errorMessage = "failed to get transaction"
			log.Error("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			const errorMessage = "failed to get transaction receipt"
			log.Error("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return toRPCTransaction(tx, receipt.BlockNumber, receipt.BlockHash, uint64(receipt.TransactionIndex)), nil
	})
}

// GetTransactionByBlockNumberAndIndex returns information about a transaction by
// block number and transaction index position.
func (e *Eth) GetTransactionByBlockNumberAndIndex(number *BlockNumber, index Index) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		var err error
		blockNumber, err := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if err != nil {
			return nil, err
		}

		tx, err := e.state.GetTransactionByBlockNumberAndIndex(ctx, blockNumber, uint64(index), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			const errorMessage = "failed to get transaction"
			log.Error("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			const errorMessage = "failed to get transaction receipt"
			log.Error("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return toRPCTransaction(tx, receipt.BlockNumber, receipt.BlockHash, uint64(receipt.TransactionIndex)), nil
	})
}

// GetTransactionByHash returns a transaction by his hash
func (e *Eth) GetTransactionByHash(hash common.Hash) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		tx, err := e.state.GetTransactionByHash(ctx, hash, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			const errorMessage = "failed to load transaction by hash from state"
			log.Errorf("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		receipt, err := e.state.GetTransactionReceipt(ctx, hash, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			const errorMessage = "failed to load transaction receipt from state"
			log.Errorf("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return toRPCTransaction(tx, receipt.BlockNumber, receipt.BlockHash, uint64(receipt.TransactionIndex)), nil
	})
}

// GetTransactionCount returns account nonce
func (e *Eth) GetTransactionCount(address common.Address, number *BlockNumber) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		var err error
		blockNumber, err := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if err != nil {
			return nil, err
		}

		nonce, err := e.state.GetNonce(ctx, address, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return hex.EncodeUint64(0), nil
		} else if err != nil {
			const errorMessage = "failed to count transactions"
			log.Errorf("%v:%v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return hex.EncodeUint64(nonce), nil
	})
}

// GetBlockTransactionCountByHash returns the number of transactions in a
// block from a block mlocking the given block hash.
func (e *Eth) GetBlockTransactionCountByHash(hash common.Hash) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		c, err := e.state.GetBlockTransactionCountByHash(ctx, hash, dbTx)
		if err != nil {
			const errorMessage = "failed to count transactions"
			log.Errorf("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return argUint64(c), nil
	})
}

// GetBlockTransactionCountByNumber returns the number of transactions in a
// block from a block mlocking the given block number.
func (e *Eth) GetBlockTransactionCountByNumber(number *BlockNumber) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		var err error
		blockNumber, err := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if err != nil {
			return nil, err
		}

		c, err := e.state.GetBlockTransactionCountByNumber(ctx, blockNumber, dbTx)
		if err != nil {
			const errorMessage = "failed to count transactions"
			log.Errorf("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return argUint64(c), nil
	})
}

// GetTransactionReceipt returns a transaction receipt by his hash
func (e *Eth) GetTransactionReceipt(hash common.Hash) (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		tx, err := e.state.GetTransactionByHash(ctx, hash, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			const errorMessage = "failed to get tx from state"
			log.Errorf("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		r, err := e.state.GetTransactionReceipt(ctx, hash, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			const errorMessage = "failed to get tx receipt from state"
			log.Errorf("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		receipt, err := receiptToRPCReceipt(*tx, r)
		if err != nil {
			const errorMessage = "failed to build the receipt response"
			log.Errorf("%v: %v", errorMessage, err)
			return nil, newRPCError(defaultErrorCode, errorMessage)
		}

		return receipt, nil
	})
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
		const errorMessage = "invalid tx input"
		log.Errorf("%v:%v", errorMessage, err)
		return nil, newRPCError(invalidParamsErrorCode, errorMessage)
	}

	log.Debugf("adding TX to the pool: %v", tx.Hash().Hex())
	if err := e.pool.AddTx(context.Background(), *tx); err != nil {
		const errorMessage = "failed to add TX to the pool"
		log.Errorf("%v[%v]:%v", errorMessage, tx.Hash().Hex(), err)
		return nil, newRPCError(defaultErrorCode, errorMessage)
	}
	log.Infof("TX added to the pool: %v", tx.Hash().Hex())

	return tx.Hash().Hex(), nil
}

// UninstallFilter uninstalls a filter with given id. Should
// always be called when wlock is no longer needed. Additionally
// Filters timeout when they aren’t requested with
// eth_getFilterChanges for a period of time.
func (e *Eth) UninstallFilter(filterID argUint64) (interface{}, error) {
	return e.storage.UninstallFilter(uint64(filterID))
}

// Syncing returns an object with data about the sync status or false.
// https://eth.wiki/json-rpc/API#eth_syncing
func (e *Eth) Syncing() (interface{}, error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, error) {
		syncInfo, err := e.state.GetSyncingInfo(ctx, dbTx)
		if err != nil {
			return nil, err
		}

		if syncInfo.CurrentBlockNumber == syncInfo.LastBlockNumberSeen {
			return false, nil
		}

		return struct {
			S argUint64 `json:"startingBlock"`
			C argUint64 `json:"currentBlock"`
			H argUint64 `json:"highestBlock"`
		}{
			S: argUint64(syncInfo.InitialSyncingBlock),
			C: argUint64(syncInfo.CurrentBlockNumber),
			H: argUint64(syncInfo.LastBlockNumberSeen),
		}, nil
	})
}

// GetUncleByBlockHashAndIndex returns information about a uncle of a
// block by hash and uncle index position
func (e *Eth) GetUncleByBlockHashAndIndex() (interface{}, error) {
	return nil, nil
}

// GetUncleByBlockHashAndIndex returns information about a uncle of a
// block by number and uncle index position
func (e *Eth) GetUncleByBlockNumberAndIndex() (interface{}, error) {
	return nil, nil
}

// GetUncleCountByBlockHash returns the number of uncles in a block
// mlocking the given block hash
func (e *Eth) GetUncleCountByBlockHash() (interface{}, error) {
	return "0x0", nil
}

// GetUncleCountByBlockNumber returns the number of uncles in a block
// mlocking the given block number
func (e *Eth) GetUncleCountByBlockNumber() (interface{}, error) {
	return "0x0", nil
}

// ProtocolVersion
func (e *Eth) ProtocolVersion() (interface{}, error) {
	return "0x0", nil
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

func (e *Eth) getBlockHeader(ctx context.Context, number BlockNumber, dbTx pgx.Tx) (*types.Header, error) {
	switch number {
	case LatestBlockNumber:
		block, err := e.state.GetLastBlock(ctx, dbTx)
		if err != nil {
			return nil, err
		}
		return block.Header, nil

	case EarliestBlockNumber:
		block, err := e.state.GetBlockByNumber(ctx, 0, dbTx)
		if err != nil {
			return nil, err
		}
		return block.Header, nil

	case PendingBlockNumber:
		lastBlock, err := e.state.GetLastBlock(ctx, dbTx)
		if err != nil {
			return nil, err
		}
		header := &types.Header{
			ParentHash: lastBlock.Hash(),
			Number:     big.NewInt(0).SetUint64(lastBlock.Number().Uint64() + 1),
			Difficulty: big.NewInt(0),
		}
		return header, nil

	default:
		return e.state.GetBlockHeader(ctx, uint64(number), dbTx)
	}
}
