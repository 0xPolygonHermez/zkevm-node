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
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4"
)

// EthEndpoints contains implementations for the "eth" RPC endpoints
type EthEndpoints struct {
	cfg     Config
	pool    jsonRPCTxPool
	state   stateInterface
	storage storageInterface
	txMan   dbTxManager
}

// newEthEndpoints creates an new instance of Eth
func newEthEndpoints(cfg Config, p jsonRPCTxPool, s stateInterface, storage storageInterface) *EthEndpoints {
	e := &EthEndpoints{cfg: cfg, pool: p, state: s, storage: storage}
	s.RegisterNewL2BlockEventHandler(e.onNewL2Block)

	return e
}

// BlockNumber returns current block number
func (e *EthEndpoints) BlockNumber() (interface{}, rpcError) {
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
func (e *EthEndpoints) Call(arg *txnArgs, blockNrOrHash *rpc.BlockNumberOrHash) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		ethBlockNumber, hasNum := blockNrOrHash.Number()
		bnValue := fromEthBlockNumber(ethBlockNumber)
		if !hasNum {
			if blockHash, hasHash := blockNrOrHash.Hash(); hasHash {
				if block, err := e.state.GetL2BlockByHash(ctx, blockHash, dbTx); err != nil {
					return rpcErrorResponse(defaultErrorCode, "Unknown block hash", err)
				} else {
					bnValue = BlockNumber(block.Header().Number.Int64())
				}
			}
			bnValue = LatestBlockNumber
		}
		number := &bnValue

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

		result := e.state.ProcessUnsignedTransaction(ctx, tx, sender, blockNumberToProcessTx, false, dbTx)
		if result.Failed() {
			data := make([]byte, len(result.ReturnValue))
			copy(data, result.ReturnValue)
			return rpcErrorResponseWithData(revertedErrorCode, result.Err.Error(), &data, nil)
		}

		return argBytesPtr(result.ReturnValue), nil
	})
}

// ChainId returns the chain id of the client
func (e *EthEndpoints) ChainId() (interface{}, rpcError) { //nolint:revive
	return hex.EncodeUint64(e.cfg.ChainID), nil
}

// EstimateGas generates and returns an estimate of how much gas is necessary to
// allow the transaction to complete.
// The transaction will not be added to the blockchain.
// Note that the estimate may be significantly more than the amount of gas actually
// used by the transaction, for a variety of reasons including EVM mechanics and
// node performance.
func (e *EthEndpoints) EstimateGas(arg *txnArgs, number *BlockNumber) (interface{}, rpcError) {
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
func (e *EthEndpoints) GasPrice() (interface{}, rpcError) {
	ctx := context.Background()
	gasPrice, err := e.pool.GetGasPrice(ctx)
	if err != nil {
		return "0x0", nil
	}
	return hex.EncodeUint64(gasPrice), nil
}

// GetBalance returns the account's balance at the referenced block
func (e *EthEndpoints) GetBalance(address common.Address, number *BlockNumber) (interface{}, rpcError) {
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
func (e *EthEndpoints) GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, rpcError) {
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
func (e *EthEndpoints) GetBlockByNumber(number BlockNumber, fullTx bool) (interface{}, rpcError) {
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
func (e *EthEndpoints) GetCode(address common.Address, number *BlockNumber) (interface{}, rpcError) {
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
func (e *EthEndpoints) GetCompilers() (interface{}, rpcError) {
	return []interface{}{}, nil
}

// GetFilterChanges polling method for a filter, which returns
// an array of logs which occurred since last poll.
func (e *EthEndpoints) GetFilterChanges(filterID string) (interface{}, rpcError) {
	filter, err := e.storage.GetFilter(filterID)
	if errors.Is(err, ErrNotFound) {
		return rpcErrorResponse(defaultErrorCode, "filter not found", err)
	} else if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to get filter from storage", err)
	}

	switch filter.Type {
	case FilterTypeBlock:
		{
			res, err := e.state.GetL2BlockHashesSince(context.Background(), filter.LastPoll, nil)
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
			filterParameters := filter.Parameters.(LogFilter)
			filterParameters.Since = &filter.LastPoll

			resInterface, err := e.internalGetLogs(context.Background(), nil, filterParameters)
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
		}
	default:
		return nil, nil
	}
}

// GetFilterLogs returns an array of all logs matching filter
// with given id.
func (e *EthEndpoints) GetFilterLogs(filterID string) (interface{}, rpcError) {
	filter, err := e.storage.GetFilter(filterID)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to get filter from storage", err)
	}

	if filter.Type != FilterTypeLog {
		return nil, nil
	}

	filterParameters := filter.Parameters.(LogFilter)
	filterParameters.Since = nil

	return e.GetLogs(filterParameters)
}

// GetLogs returns a list of logs accordingly to the provided filter
func (e *EthEndpoints) GetLogs(filter LogFilter) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		return e.internalGetLogs(ctx, dbTx, filter)
	})
}

func (e *EthEndpoints) internalGetLogs(ctx context.Context, dbTx pgx.Tx, filter LogFilter) (interface{}, rpcError) {
	var err error
	var fromBlock uint64 = 0
	if filter.FromBlock != nil {
		var rpcErr rpcError
		fromBlock, rpcErr = filter.FromBlock.getNumericBlockNumber(ctx, e.state, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}
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
func (e *EthEndpoints) GetStorageAt(address common.Address, position common.Hash, number *BlockNumber) (interface{}, rpcError) {
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
func (e *EthEndpoints) GetTransactionByBlockHashAndIndex(hash common.Hash, index Index) (interface{}, rpcError) {
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
		return toRPCTransaction(*tx, receipt.BlockNumber, &receipt.BlockHash, &txIndex), nil
	})
}

// GetTransactionByBlockNumberAndIndex returns information about a transaction by
// block number and transaction index position.
func (e *EthEndpoints) GetTransactionByBlockNumberAndIndex(number *BlockNumber, index Index) (interface{}, rpcError) {
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
		return toRPCTransaction(*tx, receipt.BlockNumber, &receipt.BlockHash, &txIndex), nil
	})
}

// GetTransactionByHash returns a transaction by his hash
func (e *EthEndpoints) GetTransactionByHash(hash common.Hash) (interface{}, rpcError) {
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
			return toRPCTransaction(*tx, receipt.BlockNumber, &receipt.BlockHash, &txIndex), nil
		}

		// if the tx does not exist in the state, look for it in the pool
		poolTx, err := e.pool.GetTxByHash(ctx, hash)
		if errors.Is(err, pgpoolstorage.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to load transaction by hash from pool", err)
		}
		tx = &poolTx.Transaction

		return toRPCTransaction(*tx, nil, nil, nil), nil
	})
}

// GetTransactionCount returns account nonce
func (e *EthEndpoints) GetTransactionCount(address common.Address, number *BlockNumber) (interface{}, rpcError) {
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
// block from a block matching the given block hash.
func (e *EthEndpoints) GetBlockTransactionCountByHash(hash common.Hash) (interface{}, rpcError) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, rpcError) {
		c, err := e.state.GetL2BlockTransactionCountByHash(ctx, hash, dbTx)
		if err != nil {
			return rpcErrorResponse(defaultErrorCode, "failed to count transactions", err)
		}

		return argUint64(c), nil
	})
}

// GetBlockTransactionCountByNumber returns the number of transactions in a
// block from a block matching the given block number.
func (e *EthEndpoints) GetBlockTransactionCountByNumber(number *BlockNumber) (interface{}, rpcError) {
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
func (e *EthEndpoints) GetTransactionReceipt(hash common.Hash) (interface{}, rpcError) {
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
func (e *EthEndpoints) NewBlockFilter() (interface{}, rpcError) {
	return e.newBlockFilter(nil)
}

// internal
func (e *EthEndpoints) newBlockFilter(wsConn *websocket.Conn) (interface{}, rpcError) {
	id, err := e.storage.NewBlockFilter(wsConn)
	if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to create new block filter", err)
	}

	return id, nil
}

// NewFilter creates a filter object, based on filter options,
// to notify when the state changes (logs). To check if the state
// has changed, call eth_getFilterChanges.
func (e *EthEndpoints) NewFilter(filter LogFilter) (interface{}, rpcError) {
	return e.newFilter(nil, filter)
}

// internal
func (e *EthEndpoints) newFilter(wsConn *websocket.Conn, filter LogFilter) (interface{}, rpcError) {
	id, err := e.storage.NewLogFilter(wsConn, filter)
	if errors.Is(err, ErrFilterInvalidPayload) {
		return rpcErrorResponse(invalidParamsErrorCode, err.Error(), nil)
	} else if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to create new log filter", err)
	}

	return id, nil
}

// NewPendingTransactionFilter creates a filter in the node, to
// notify when new pending transactions arrive. To check if the
// state has changed, call eth_getFilterChanges.
func (e *EthEndpoints) NewPendingTransactionFilter() (interface{}, rpcError) {
	return e.newPendingTransactionFilter(nil)
}

// internal
func (e *EthEndpoints) newPendingTransactionFilter(wsConn *websocket.Conn) (interface{}, rpcError) {
	return nil, newRPCError(defaultErrorCode, "not supported yet")
	// id, err := e.storage.NewPendingTransactionFilter(wsConn)
	// if err != nil {
	// 	return rpcErrorResponse(defaultErrorCode, "failed to create new pending transaction filter", err)
	// }

	// return id, nil
}

// SendRawTransaction has two different ways to handle new transactions:
// - for Sequencer nodes it tries to add the tx to the pool
// - for Non-Sequencer nodes it relays the Tx to the Sequencer node
func (e *EthEndpoints) SendRawTransaction(input string) (interface{}, rpcError) {
	if e.cfg.SequencerNodeURI != "" {
		return e.relayTxToSequencerNode(input)
	} else {
		return e.tryToAddTxToPool(input)
	}
}

func (e *EthEndpoints) relayTxToSequencerNode(input string) (interface{}, rpcError) {
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

func (e *EthEndpoints) tryToAddTxToPool(input string) (interface{}, rpcError) {
	tx, err := hexToTx(input)
	if err != nil {
		return rpcErrorResponse(invalidParamsErrorCode, "invalid tx input", err)
	}

	log.Infof("adding TX to the pool: %v", tx.Hash().Hex())
	if err := e.pool.AddTx(context.Background(), *tx); err != nil {
		return rpcErrorResponse(defaultErrorCode, err.Error(), nil)
	}
	log.Infof("TX added to the pool: %v", tx.Hash().Hex())

	return tx.Hash().Hex(), nil
}

// UninstallFilter uninstalls a filter with given id.
func (e *EthEndpoints) UninstallFilter(filterID string) (interface{}, rpcError) {
	err := e.storage.UninstallFilter(filterID)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	} else if err != nil {
		return rpcErrorResponse(defaultErrorCode, "failed to uninstall filter", err)
	}

	return true, nil
}

// Syncing returns an object with data about the sync status or false.
// https://eth.wiki/json-rpc/API#eth_syncing
func (e *EthEndpoints) Syncing() (interface{}, rpcError) {
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
func (e *EthEndpoints) GetUncleByBlockHashAndIndex() (interface{}, rpcError) {
	return nil, nil
}

// GetUncleByBlockNumberAndIndex returns information about a uncle of a
// block by number and uncle index position
func (e *EthEndpoints) GetUncleByBlockNumberAndIndex() (interface{}, rpcError) {
	return nil, nil
}

// GetUncleCountByBlockHash returns the number of uncles in a block
// matching the given block hash
func (e *EthEndpoints) GetUncleCountByBlockHash() (interface{}, rpcError) {
	return "0x0", nil
}

// GetUncleCountByBlockNumber returns the number of uncles in a block
// matching the given block number
func (e *EthEndpoints) GetUncleCountByBlockNumber() (interface{}, rpcError) {
	return "0x0", nil
}

// ProtocolVersion returns the protocol version.
func (e *EthEndpoints) ProtocolVersion() (interface{}, rpcError) {
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

func (e *EthEndpoints) getBlockHeader(ctx context.Context, number BlockNumber, dbTx pgx.Tx) (*types.Header, error) {
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

func (e *EthEndpoints) updateFilterLastPoll(filterID string) rpcError {
	err := e.storage.UpdateFilterLastPoll(filterID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return newRPCError(defaultErrorCode, "failed to update last time the filter changes were requested")
	}
	return nil
}

// Subscribe Creates a new subscription over particular events.
// The node will return a subscription id.
// For each event that matches the subscription a notification with relevant
// data is sent together with the subscription id.
func (e *EthEndpoints) Subscribe(wsConn *websocket.Conn, name string, logFilter *LogFilter) (interface{}, rpcError) {
	switch name {
	case "newHeads":
		return e.newBlockFilter(wsConn)
	case "logs":
		var lf LogFilter
		if logFilter != nil {
			lf = *logFilter
		}
		return e.newFilter(wsConn, lf)
	case "pendingTransactions", "newPendingTransactions":
		return e.newPendingTransactionFilter(wsConn)
	case "syncing":
		return nil, newRPCError(defaultErrorCode, "not supported yet")
	default:
		return nil, newRPCError(defaultErrorCode, "invalid filter name")
	}
}

// Unsubscribe uninstalls the filter based on the provided filterID
func (e *EthEndpoints) Unsubscribe(wsConn *websocket.Conn, filterID string) (interface{}, rpcError) {
	return e.UninstallFilter(filterID)
}

// uninstallFilterByWSConn uninstalls the filters connected to the
// provided web socket connection
func (e *EthEndpoints) uninstallFilterByWSConn(wsConn *websocket.Conn) error {
	return e.storage.UninstallFilterByWSConn(wsConn)
}

// onNewL2Block is triggered when the state triggers the event for a new l2 block
func (e *EthEndpoints) onNewL2Block(event state.NewL2BlockEvent) {
	blockFilters, err := e.storage.GetAllBlockFiltersWithWSConn()
	if err != nil {
		log.Errorf("failed to get all block filters with web sockets connections: %v", err)
	} else {
		for _, filter := range blockFilters {
			b := l2BlockToRPCBlock(&event.Block, false)
			e.sendSubscriptionResponse(filter, b)
		}
	}

	logFilters, err := e.storage.GetAllLogFiltersWithWSConn()
	if err != nil {
		log.Errorf("failed to get all log filters with web sockets connections: %v", err)
	} else {
		for _, filter := range logFilters {
			changes, err := e.GetFilterChanges(filter.ID)
			if err != nil {
				log.Errorf("failed to get filters changes for filter %v with web sockets connections: %v", filter.ID, err)
				continue
			}

			if changes != nil {
				e.sendSubscriptionResponse(filter, changes)
			}
		}
	}
}

func (e *EthEndpoints) sendSubscriptionResponse(filter *Filter, data interface{}) {
	const errMessage = "Unable to write WS message to filter %v, %s"
	result, err := json.Marshal(data)
	if err != nil {
		log.Errorf(fmt.Sprintf(errMessage, filter.ID, err.Error()))
	}

	res := SubscriptionResponse{
		JSONRPC: "2.0",
		Method:  "eth_subscription",
		Params: SubscriptionResponseParams{
			Subscription: filter.ID,
			Result:       result,
		},
	}
	message, err := json.Marshal(res)
	if err != nil {
		log.Errorf(fmt.Sprintf(errMessage, filter.ID, err.Error()))
	}

	err = filter.WsConn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Errorf(fmt.Sprintf(errMessage, filter.ID, err.Error()))
	}
}
