package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Eth contains implementations for the "eth" RPC endpoints
type Eth struct {
	cfg     Config
	pool    jsonRPCTxPool
	state   stateInterface
	gpe     gasPriceEstimator
	storage storageInterface
	txMan   dbTxManager
}

// BlockNumber returns current block number
func (e *Eth) BlockNumber() (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		lastBlockNumber, err := e.state.GetLastL2BlockNumber(ctx, dbTx)
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
func (e *Eth) Call(arg *txnArgs, number *BlockNumber) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		// If the caller didn't supply the gas limit in the message, then we set it to maximum possible => block gas limit
		if arg.Gas == nil || *arg.Gas == argUint64(0) {
			header, err := e.getBlockHeader(ctx, *number, dbTx)
			if err != nil {
				return rpcErrorResponse(defaultErrorCode, "failed to get block header", err)
			}

			gas := argUint64(header.GasLimit)
			arg.Gas = &gas
		}

		blockNumber, rpcErr := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		sender, tx, err := arg.ToUnsignedTransaction(ctx, e.state, blockNumber, e.cfg, dbTx)
		if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to convert arguments into an unsigned transaction", err)
		}

		var blockNumberToProcessTx *uint64
		if number != nil && *number != LatestBlockNumber && *number != PendingBlockNumber {
			blockNumberToProcessTx = &blockNumber
		}

		result := e.state.ProcessUnsignedTransaction(ctx, tx, sender, blockNumberToProcessTx, true, dbTx)
		if result.Failed() {
			return rpcErrorResponse(defaultErrorCode, result.Err.Error(), nil)
		}

		return argBytesPtr(result.ReturnValue), nil
	})
}

// ChainId returns the chain id of the client
func (e *Eth) ChainId() (interface{}, rpcError) { //nolint:revive
	return hex.EncodeUint64(e.cfg.ChainID), nil
}

// EstimateGas generates and returns an estimate of how much gas is necessary to
// allow the transaction to complete.
// The transaction will not be added to the blockchain.
// Note that the estimate may be significantly more than the amount of gas actually
// used by the transaction, for a variety of reasons including EVM mechanics and
// node performance.
func (e *Eth) EstimateGas(arg *txnArgs, number *BlockNumber) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		blockNumber, rpcErr := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		sender, tx, err := arg.ToUnsignedTransaction(ctx, e.state, blockNumber, e.cfg, dbTx)
		if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to convert arguments into an unsigned transaction", err)
		}

		var blockNumberToProcessTx *uint64
		if number != nil && *number != LatestBlockNumber && *number != PendingBlockNumber {
			blockNumberToProcessTx = &blockNumber
		}

		gasEstimation, err := e.state.EstimateGas(tx, sender, blockNumberToProcessTx, dbTx)
		if err != nil {
			return rpcErrorResponse(defaultErrorCode, err.Error(), nil)
		}
		return hex.EncodeUint64(gasEstimation), nil
	})
}

// GasPrice returns the average gas price based on the last x blocks
func (e *Eth) GasPrice() (interface{}, rpcError) {
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
func (e *Eth) GetBalance(address common.Address, number *BlockNumber) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		blockNumber, rpcErr := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		balance, err := e.state.GetBalance(ctx, address, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return hex.EncodeUint64(0), nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get balance from state", err)
		}

		return hex.EncodeBig(balance), nil
	})
}

// GetBlockByHash returns information about a block by hash
func (e *Eth) GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		block, err := e.state.GetL2BlockByHash(ctx, hash, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get block by hash from state", err)
		}

		rpcBlock := l2BlockToRPCBlock(block, fullTx)

		return rpcBlock, nil
	})
}

// GetBlockByNumber returns information about a block by block number
func (e *Eth) GetBlockByNumber(number BlockNumber, fullTx bool) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		if number == PendingBlockNumber {
			lastBlock, err := e.state.GetLastL2Block(ctx, dbTx)
			if err != nil {
				return rpcErrorResponse(defaultErrorCode, "couldn't load last block from state to compute the pending block", err)
			}
			header := types.CopyHeader(lastBlock.Header())
			header.ParentHash = lastBlock.Hash()
			header.Number = big.NewInt(0).SetUint64(lastBlock.Number().Uint64() + 1)
			header.TxHash = types.EmptyRootHash
			header.UncleHash = types.EmptyUncleHash
			block := types.NewBlockWithHeader(header)
			rpcBlock := l2BlockToRPCBlock(block, fullTx)

			return rpcBlock, nil
		}
		var err error
		blockNumber, rpcErr := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		block, err := e.state.GetL2BlockByNumber(ctx, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, fmt.Sprintf("couldn't load block from state by number %v", blockNumber), err)
		}

		rpcBlock := l2BlockToRPCBlock(block, fullTx)

		return rpcBlock, nil
	})
}

// GetCode returns account code at given block number
func (e *Eth) GetCode(address common.Address, number *BlockNumber) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		var err error
		blockNumber, rpcErr := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		code, err := e.state.GetCode(ctx, address, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return "0x", nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get code", err)
		}

		return argBytes(code), nil
	})
}

// GetCompilers eth_getCompilers
func (e *Eth) GetCompilers() (interface{}, rpcError) {
	return []interface{}{}, nil
}

// GetFilterChanges polling method for a filter, which returns
// an array of logs which occurred since last poll.
func (e *Eth) GetFilterChanges(filterID argUint64) (interface{}, rpcError) {
	filter, err := e.storage.GetFilter(uint64(filterID))
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to get filter from storage", err)
	}

	switch filter.Type {
	case FilterTypeBlock:
		{
			return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
				res, err := e.state.GetL2BlockHashesSince(ctx, filter.LastPoll, dbTx)
				if err != nil {
					return rpcErrorResponse(defaultErrorCode, "failed to get block hashes", err)
				}
				rpcErr := e.updateFilterLastPoll(filter.ID)
				if rpcErr != nil {
					return nil, rpcErr
				}
				if len(res) == 0 {
					return nil, nil
				}
				return res, nil
			})
		}
	case FilterTypePendingTx:
		{
			res, err := e.pool.GetPendingTxHashesSince(context.Background(), filter.LastPoll)
			if err != nil {
				return rpcErrorResponse(defaultErrorCode, "failed to get pending transaction hashes", err)
			}
			rpcErr := e.updateFilterLastPoll(filter.ID)
			if rpcErr != nil {
				return nil, rpcErr
			}
			if len(res) == 0 {
				return nil, nil
			}
			return res, nil
		}
	case FilterTypeLog:
		{
			return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
				filterParameters := &LogFilter{}
				err = json.Unmarshal([]byte(filter.Parameters), filterParameters)
				if err != nil {
					return rpcErrorResponse(defaultErrorCode, "failed to read filter parameters", err)
				}

				filterParameters.Since = &filter.LastPoll

				resInterface, err := e.internalGetLogs(ctx, dbTx, filterParameters)
				if err != nil {
					return nil, err
				}
				rpcErr := e.updateFilterLastPoll(filter.ID)
				if rpcErr != nil {
					return nil, rpcErr
				}
				res := resInterface.([]rpcLog)
				if len(res) == 0 {
					return nil, nil
				}
				return res, nil
			})
		}
	default:
		return nil, nil
	}
}

// GetFilterLogs returns an array of all logs mlocking filter
// with given id.
func (e *Eth) GetFilterLogs(filterID argUint64) (interface{}, rpcError) {
	filter, err := e.storage.GetFilter(uint64(filterID))
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to get filter from storage", err)
	}

	if filter.Type != FilterTypeLog {
		return nil, nil
	}

	filterParameters := &LogFilter{}
	err = json.Unmarshal([]byte(filter.Parameters), filterParameters)
	if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to read filter parameters", err)
	}

	filterParameters.Since = nil

	return e.GetLogs(filterParameters)
}

// GetLogs returns a list of logs accordingly to the provided filter
func (e *Eth) GetLogs(filter *LogFilter) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		return e.internalGetLogs(ctx, dbTx, filter)
	})
}

func (e *Eth) internalGetLogs(ctx context.Context, dbTx pgx.Tx, filter *LogFilter) (interface{}, rpcError) {
	var err error
	fromBlock, rpcErr := filter.FromBlock.getNumericBlockNumber(ctx, e.state, dbTx)
	if rpcErr != nil {
		return nil, rpcErr
	}

	toBlock, rpcErr := filter.ToBlock.getNumericBlockNumber(ctx, e.state, dbTx)
	if rpcErr != nil {
		return nil, rpcErr
	}

	logs, err := e.state.GetLogs(ctx, fromBlock, toBlock, filter.Addresses, filter.Topics, filter.BlockHash, filter.Since, dbTx)
	if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to get logs from state", err)
	}

	result := make([]rpcLog, 0, len(logs))
	for _, l := range logs {
		result = append(result, logToRPCLog(*l))
	}

	return result, nil
}

// GetStorageAt gets the value stored for an specific address and position
func (e *Eth) GetStorageAt(address common.Address, position common.Hash, number *BlockNumber) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		var err error
		blockNumber, rpcErr := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		value, err := e.state.GetStorageAt(ctx, address, position.Big(), blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return argBytesPtr(common.Hash{}.Bytes()), nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get storage value from state", err)
		}

		return argBytesPtr(common.BigToHash(value).Bytes()), nil
	})
}

// GetTransactionByBlockHashAndIndex returns information about a transaction by
// block hash and transaction index position.
func (e *Eth) GetTransactionByBlockHashAndIndex(hash common.Hash, index Index) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		tx, err := e.state.GetTransactionByL2BlockHashAndIndex(ctx, hash, uint64(index), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get transaction", err)
		}

		receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get transaction receipt", err)
		}

		txIndex := uint64(receipt.TransactionIndex)
		return toRPCTransaction(tx, receipt.BlockNumber, &receipt.BlockHash, &txIndex), nil
	})
}

// GetTransactionByBlockNumberAndIndex returns information about a transaction by
// block number and transaction index position.
func (e *Eth) GetTransactionByBlockNumberAndIndex(number *BlockNumber, index Index) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		var err error
		blockNumber, rpcErr := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		tx, err := e.state.GetTransactionByL2BlockNumberAndIndex(ctx, blockNumber, uint64(index), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get transaction", err)
		}

		receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get transaction receipt", err)
		}

		txIndex := uint64(receipt.TransactionIndex)
		return toRPCTransaction(tx, receipt.BlockNumber, &receipt.BlockHash, &txIndex), nil
	})
}

// GetTransactionByHash returns a transaction by his hash
func (e *Eth) GetTransactionByHash(hash common.Hash) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		// try to get tx from state
		tx, err := e.state.GetTransactionByHash(ctx, hash, dbTx)
		if err != nil && !errors.Is(err, state.ErrNotFound) {
			return rpcErrorResponse(defaultErrorCode, "failed to load transaction by hash from state", err)
		}
		if tx != nil {
			receipt, err := e.state.GetTransactionReceipt(ctx, hash, dbTx)
			if errors.Is(err, state.ErrNotFound) {
				return rpcErrorResponse(defaultErrorCode, "transaction receipt not found", err)
			} else if err != nil {
				return rpcErrorResponse(defaultErrorCode, "failed to load transaction receipt from state", err)
			}

			txIndex := uint64(receipt.TransactionIndex)
			return toRPCTransaction(tx, receipt.BlockNumber, &receipt.BlockHash, &txIndex), nil
		}

		// if the tx does not exist in the state, look for it in the pool
		poolTx, err := e.pool.GetTxByHash(ctx, hash)
		if errors.Is(err, pgpoolstorage.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to load transaction by hash from pool", err)
		}
		tx = &poolTx.Transaction

		return toRPCTransaction(tx, nil, nil, nil), nil
	})
}

// GetTransactionCount returns account nonce
func (e *Eth) GetTransactionCount(address common.Address, number *BlockNumber) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		var pendingNonce uint64
		var nonce uint64
		var err error
		if number != nil && *number == PendingBlockNumber {
			pendingNonce, err = e.pool.GetNonce(ctx, address)
			if err != nil {
				return rpcErrorResponse(defaultErrorCode, "failed to count pending transactions", err)
			}
		}

		blockNumber, rpcErr := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}
		nonce, err = e.state.GetNonce(ctx, address, blockNumber, dbTx)

		if errors.Is(err, state.ErrNotFound) {
			return hex.EncodeUint64(0), nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to count transactions", err)
		}

		if pendingNonce > nonce {
			nonce = pendingNonce
		}

		return hex.EncodeUint64(nonce), nil
	})
}

// GetBlockTransactionCountByHash returns the number of transactions in a
// block from a block mlocking the given block hash.
func (e *Eth) GetBlockTransactionCountByHash(hash common.Hash) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		c, err := e.state.GetL2BlockTransactionCountByHash(ctx, hash, dbTx)
		if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to count transactions", err)
		}

		return argUint64(c), nil
	})
}

// GetBlockTransactionCountByNumber returns the number of transactions in a
// block from a block mlocking the given block number.
func (e *Eth) GetBlockTransactionCountByNumber(number *BlockNumber) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		if number != nil && *number == PendingBlockNumber {
			c, err := e.pool.CountPendingTransactions(ctx)
			if err != nil {
				return rpcErrorResponse(defaultErrorCode, "failed to count pending transactions", err)
			}
			return argUint64(c), nil
		}

		var err error
		blockNumber, rpcErr := number.getNumericBlockNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		c, err := e.state.GetL2BlockTransactionCountByNumber(ctx, blockNumber, dbTx)
		if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to count transactions", err)
		}

		return argUint64(c), nil
	})
}

// GetTransactionReceipt returns a transaction receipt by his hash
func (e *Eth) GetTransactionReceipt(hash common.Hash) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		tx, err := e.state.GetTransactionByHash(ctx, hash, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get tx from state", err)
		}

		r, err := e.state.GetTransactionReceipt(ctx, hash, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get tx receipt from state", err)
		}

		receipt, err := receiptToRPCReceipt(*tx, r)
		if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to build the receipt response", err)
		}

		return receipt, nil
	})
}

// NewBlockFilter creates a filter in the node, to notify when
// a new block arrives. To check if the state has changed,
// call eth_getFilterChanges.
func (e *Eth) NewBlockFilter() (interface{}, rpcError) {
	id, err := e.storage.NewBlockFilter()
	if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to create new block filter", err)
	}

	return argUint64(id), nil
}

// NewFilter creates a filter object, based on filter options,
// to notify when the state changes (logs). To check if the state
// has changed, call eth_getFilterChanges.
func (e *Eth) NewFilter(filter *LogFilter) (interface{}, rpcError) {
	id, err := e.storage.NewLogFilter(*filter)
	if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to create new log filter", err)
	}

	return argUint64(id), nil
}

// NewPendingTransactionFilter creates a filter in the node, to
// notify when new pending transactions arrive. To check if the
// state has changed, call eth_getFilterChanges.
func (e *Eth) NewPendingTransactionFilter(filterID argUint64) (interface{}, rpcError) {
	id, err := e.storage.NewPendingTransactionFilter()
	if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to create new pending transaction filter", err)
	}

	return argUint64(id), nil
}

// SendRawTransaction has two different ways to handle new transactions:
// - for Sequencer nodes it tries to add the tx to the pool
// - for Non-Sequencer nodes it relays the Tx to the Sequencer node
func (e *Eth) SendRawTransaction(input string) (interface{}, rpcError) {
	if e.cfg.SequencerNodeURI != "" {
		return e.relayTxToSequencerNode(input)
	} else {
		return e.tryToAddTxToPool(input)
	}
}

func (e *Eth) relayTxToSequencerNode(input string) (interface{}, rpcError) {
	res, err := JSONRPCCall(e.cfg.SequencerNodeURI, "eth_sendRawTransaction", input)
	if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to relay tx to the sequencer node", err)
	}

	if res.Error != nil {
		return rpcErrorResponse(res.Error.Code, res.Error.Message, nil)
	}

	txHash := res.Result

	return txHash, nil
}

func (e *Eth) tryToAddTxToPool(input string) (interface{}, rpcError) {
	tx, err := hexToTx(input)
	if err != nil {
		return rpcErrorResponse(invalidParamsErrorCode, "invalid tx input", err)
	}

	log.Debugf("adding TX to the pool: %v", tx.Hash().Hex())
	if err := e.pool.AddTx(context.Background(), *tx); err != nil {
		return rpcErrorResponse(defaultErrorCode, err.Error(), nil)
	}
	log.Infof("TX added to the pool: %v", tx.Hash().Hex())

	return tx.Hash().Hex(), nil
}

// UninstallFilter uninstalls a filter with given id. Should
// always be called when wlock is no longer needed. Additionally
// Filters timeout when they aren’t requested with
// eth_getFilterChanges for a period of time.
func (e *Eth) UninstallFilter(filterID argUint64) (interface{}, rpcError) {
	uninstalled, err := e.storage.UninstallFilter(uint64(filterID))
	if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to uninstall filter", err)
	}

	return uninstalled, nil
}

// Syncing returns an object with data about the sync status or false.
// https://eth.wiki/json-rpc/API#eth_syncing
func (e *Eth) Syncing() (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		syncInfo, err := e.state.GetSyncingInfo(ctx, dbTx)
		if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to get syncing info from state", err)
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
func (e *Eth) GetUncleByBlockHashAndIndex() (interface{}, rpcError) {
	return nil, nil
}

// GetUncleByBlockNumberAndIndex returns information about a uncle of a
// block by number and uncle index position
func (e *Eth) GetUncleByBlockNumberAndIndex() (interface{}, rpcError) {
	return nil, nil
}

// GetUncleCountByBlockHash returns the number of uncles in a block
// mlocking the given block hash
func (e *Eth) GetUncleCountByBlockHash() (interface{}, rpcError) {
	return "0x0", nil
}

// GetUncleCountByBlockNumber returns the number of uncles in a block
// mlocking the given block number
func (e *Eth) GetUncleCountByBlockNumber() (interface{}, rpcError) {
	return "0x0", nil
}

// ProtocolVersion returns the protocol version.
func (e *Eth) ProtocolVersion() (interface{}, rpcError) {
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
		block, err := e.state.GetLastL2Block(ctx, dbTx)
		if err != nil {
			return nil, err
		}
		return block.Header(), nil

	case EarliestBlockNumber:
		header, err := e.state.GetL2BlockHeaderByNumber(ctx, uint64(0), dbTx)
		if err != nil {
			return nil, err
		}
		return header, nil

	case PendingBlockNumber:
		lastBlock, err := e.state.GetLastL2Block(ctx, dbTx)
		if err != nil {
			return nil, err
		}
		parentHash := lastBlock.Hash()
		number := lastBlock.Number().Uint64() + 1

		header := &types.Header{
			ParentHash: parentHash,
			Number:     big.NewInt(0).SetUint64(number),
			Difficulty: big.NewInt(0),
			GasLimit:   lastBlock.Header().GasLimit,
		}
		return header, nil

	default:
		return e.state.GetL2BlockHeaderByNumber(ctx, uint64(number), dbTx)
	}
}

func (e *Eth) updateFilterLastPoll(filterID uint64) rpcError {
	err := e.storage.UpdateFilterLastPoll(filterID)
	if err != nil {
		return newRPCError(defaultErrorCode, "failed to update last time the filter changes were requested")
	}
	return nil
}
