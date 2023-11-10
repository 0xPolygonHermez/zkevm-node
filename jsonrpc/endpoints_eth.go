package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4"
)

const (
	// DefaultSenderAddress is the address that jRPC will use
	// to communicate with the state for eth_EstimateGas and eth_Call when
	// the From field is not specified because it is optional
	DefaultSenderAddress = "0x1111111111111111111111111111111111111111"
)

// EthEndpoints contains implementations for the "eth" RPC endpoints
type EthEndpoints struct {
	cfg      Config
	chainID  uint64
	pool     types.PoolInterface
	state    types.StateInterface
	etherman types.EthermanInterface
	storage  storageInterface
	txMan    DBTxManager
}

// NewEthEndpoints creates an new instance of Eth
func NewEthEndpoints(cfg Config, chainID uint64, p types.PoolInterface, s types.StateInterface, etherman types.EthermanInterface, storage storageInterface) *EthEndpoints {
	e := &EthEndpoints{cfg: cfg, chainID: chainID, pool: p, state: s, etherman: etherman, storage: storage}
	s.RegisterNewL2BlockEventHandler(e.onNewL2Block)

	return e
}

// BlockNumber returns current block number
func (e *EthEndpoints) BlockNumber() (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		lastBlockNumber, err := e.state.GetLastL2BlockNumber(ctx, dbTx)
		if err != nil {
			return "0x0", types.NewRPCError(types.DefaultErrorCode, "failed to get the last block number from state")
		}

		return hex.EncodeUint64(lastBlockNumber), nil
	})
}

// Call executes a new message call immediately and returns the value of
// executed contract and potential error.
// Note, this function doesn't make any changes in the state/blockchain and is
// useful to execute view/pure methods and retrieve values.
func (e *EthEndpoints) Call(arg *types.TxArgs, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		if arg == nil {
			return RPCErrorResponse(types.InvalidParamsErrorCode, "missing value for required argument 0", nil, false)
		} else if blockArg == nil {
			return RPCErrorResponse(types.InvalidParamsErrorCode, "missing value for required argument 1", nil, false)
		}
		block, respErr := e.getBlockByArg(ctx, blockArg, dbTx)
		if respErr != nil {
			return nil, respErr
		}
		var blockToProcess *uint64
		if blockArg != nil {
			blockNumArg := blockArg.Number()
			if blockNumArg != nil && (*blockArg.Number() == types.LatestBlockNumber || *blockArg.Number() == types.PendingBlockNumber) {
				blockToProcess = nil
			} else {
				n := block.NumberU64()
				blockToProcess = &n
			}
		}

		// If the caller didn't supply the gas limit in the message, then we set it to maximum possible => block gas limit
		if arg.Gas == nil || uint64(*arg.Gas) <= 0 {
			header, err := e.state.GetL2BlockHeaderByNumber(ctx, block.NumberU64(), dbTx)
			if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, "failed to get block header", err, true)
			}

			gas := types.ArgUint64(header.GasLimit)
			arg.Gas = &gas
		}

		defaultSenderAddress := common.HexToAddress(DefaultSenderAddress)
		sender, tx, err := arg.ToTransaction(ctx, e.state, e.cfg.MaxCumulativeGasUsed, block.Root(), defaultSenderAddress, dbTx)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to convert arguments into an unsigned transaction", err, false)
		}

		result, err := e.state.ProcessUnsignedTransaction(ctx, tx, sender, blockToProcess, true, dbTx)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to execute the unsigned transaction", err, true)
		}

		if result.Reverted() {
			data := make([]byte, len(result.ReturnValue))
			copy(data, result.ReturnValue)
			return nil, types.NewRPCErrorWithData(types.RevertedErrorCode, result.Err.Error(), &data)
		} else if result.Failed() {
			return nil, types.NewRPCErrorWithData(types.DefaultErrorCode, result.Err.Error(), nil)
		}

		return types.ArgBytesPtr(result.ReturnValue), nil
	})
}

// ChainId returns the chain id of the client
func (e *EthEndpoints) ChainId() (interface{}, types.Error) { //nolint:revive
	return hex.EncodeUint64(e.chainID), nil
}

// Coinbase Returns the client coinbase address.
func (e *EthEndpoints) Coinbase() (interface{}, types.Error) { //nolint:revive
	if e.cfg.SequencerNodeURI != "" {
		return e.getCoinbaseFromSequencerNode()
	}
	return e.cfg.L2Coinbase.String(), nil
}

func (e *EthEndpoints) getCoinbaseFromSequencerNode() (interface{}, types.Error) {
	res, err := client.JSONRPCCall(e.cfg.SequencerNodeURI, "eth_coinbase")
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get coinbase from sequencer node", err, true)
	}

	if res.Error != nil {
		return RPCErrorResponse(res.Error.Code, res.Error.Message, nil, false)
	}

	var coinbaseAddress common.Address
	err = json.Unmarshal(res.Result, &coinbaseAddress)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to read coinbase from sequencer node", err, true)
	}
	return coinbaseAddress.String(), nil
}

// EstimateGas generates and returns an estimate of how much gas is necessary to
// allow the transaction to complete.
// The transaction will not be added to the blockchain.
// Note that the estimate may be significantly more than the amount of gas actually
// used by the transaction, for a variety of reasons including EVM mechanics and
// node performance.
func (e *EthEndpoints) EstimateGas(arg *types.TxArgs, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		if arg == nil {
			return RPCErrorResponse(types.InvalidParamsErrorCode, "missing value for required argument 0", nil, false)
		}

		block, respErr := e.getBlockByArg(ctx, blockArg, dbTx)
		if respErr != nil {
			return nil, respErr
		}

		var blockToProcess *uint64
		if blockArg != nil {
			blockNumArg := blockArg.Number()
			if blockNumArg != nil && (*blockArg.Number() == types.LatestBlockNumber || *blockArg.Number() == types.PendingBlockNumber) {
				blockToProcess = nil
			} else {
				n := block.NumberU64()
				blockToProcess = &n
			}
		}

		defaultSenderAddress := common.HexToAddress(DefaultSenderAddress)
		sender, tx, err := arg.ToTransaction(ctx, e.state, e.cfg.MaxCumulativeGasUsed, block.Root(), defaultSenderAddress, dbTx)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to convert arguments into an unsigned transaction", err, false)
		}

		gasEstimation, returnValue, err := e.state.EstimateGas(tx, sender, blockToProcess, dbTx)
		if errors.Is(err, runtime.ErrExecutionReverted) {
			data := make([]byte, len(returnValue))
			copy(data, returnValue)
			return nil, types.NewRPCErrorWithData(types.RevertedErrorCode, err.Error(), &data)
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, err.Error(), nil, true)
		}
		return hex.EncodeUint64(gasEstimation), nil
	})
}

// GasPrice returns the average gas price based on the last x blocks
func (e *EthEndpoints) GasPrice() (interface{}, types.Error) {
	ctx := context.Background()
	if e.cfg.SequencerNodeURI != "" {
		return e.getPriceFromSequencerNode()
	}
	gasPrices, err := e.pool.GetGasPrices(ctx)
	if err != nil {
		return "0x0", nil
	}
	return hex.EncodeUint64(gasPrices.L2GasPrice), nil
}

func (e *EthEndpoints) getPriceFromSequencerNode() (interface{}, types.Error) {
	res, err := client.JSONRPCCall(e.cfg.SequencerNodeURI, "eth_gasPrice")
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get gas price from sequencer node", err, true)
	}

	if res.Error != nil {
		return RPCErrorResponse(res.Error.Code, res.Error.Message, nil, false)
	}

	var gasPrice types.ArgUint64
	err = json.Unmarshal(res.Result, &gasPrice)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to read gas price from sequencer node", err, true)
	}
	return gasPrice, nil
}

// GetBalance returns the account's balance at the referenced block
func (e *EthEndpoints) GetBalance(address types.ArgAddress, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		block, rpcErr := e.getBlockByArg(ctx, blockArg, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		balance, err := e.state.GetBalance(ctx, address.Address(), block.Root())
		if errors.Is(err, state.ErrNotFound) {
			return hex.EncodeUint64(0), nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get balance from state", err, true)
		}

		return hex.EncodeBig(balance), nil
	})
}

func (e *EthEndpoints) getBlockByArg(ctx context.Context, blockArg *types.BlockNumberOrHash, dbTx pgx.Tx) (*ethTypes.Block, types.Error) {
	// If no block argument is provided, return the latest block
	if blockArg == nil {
		block, err := e.state.GetLastL2Block(ctx, dbTx)
		if err != nil {
			return nil, types.NewRPCError(types.DefaultErrorCode, "failed to get the last block number from state")
		}
		return block, nil
	}

	// If we have a block hash, try to get the block by hash
	if blockArg.IsHash() {
		block, err := e.state.GetL2BlockByHash(ctx, blockArg.Hash().Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, types.NewRPCError(types.DefaultErrorCode, "header for hash not found")
		} else if err != nil {
			return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("failed to get block by hash %v", blockArg.Hash().Hash()))
		}
		return block, nil
	}

	// Otherwise, try to get the block by number
	blockNum, rpcErr := blockArg.Number().GetNumericBlockNumber(ctx, e.state, e.etherman, dbTx)
	if rpcErr != nil {
		return nil, rpcErr
	}
	block, err := e.state.GetL2BlockByNumber(context.Background(), blockNum, dbTx)
	if errors.Is(err, state.ErrNotFound) || block == nil {
		return nil, types.NewRPCError(types.DefaultErrorCode, "header not found")
	} else if err != nil {
		return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("failed to get block by number %v", blockNum))
	}

	return block, nil
}

// GetBlockByHash returns information about a block by hash
func (e *EthEndpoints) GetBlockByHash(hash types.ArgHash, fullTx bool) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		block, err := e.state.GetL2BlockByHash(ctx, hash.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get block by hash from state", err, true)
		}

		txs := block.Transactions()
		receipts := make([]ethTypes.Receipt, 0, len(txs))
		for _, tx := range txs {
			receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
			if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load receipt for tx %v", tx.Hash().String()), err, true)
			}
			receipts = append(receipts, *receipt)
		}

		rpcBlock, err := types.NewBlock(block, receipts, fullTx, false)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't build block response for block by hash %v", hash.Hash()), err, true)
		}

		return rpcBlock, nil
	})
}

// GetBlockByNumber returns information about a block by block number
func (e *EthEndpoints) GetBlockByNumber(number types.BlockNumber, fullTx bool) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		if number == types.PendingBlockNumber {
			lastBlock, err := e.state.GetLastL2Block(ctx, dbTx)
			if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, "couldn't load last block from state to compute the pending block", err, true)
			}
			header := ethTypes.CopyHeader(lastBlock.Header())
			header.ParentHash = lastBlock.Hash()
			header.Number = big.NewInt(0).SetUint64(lastBlock.Number().Uint64() + 1)
			header.TxHash = ethTypes.EmptyRootHash
			header.UncleHash = ethTypes.EmptyUncleHash
			block := ethTypes.NewBlockWithHeader(header)
			rpcBlock, err := types.NewBlock(block, nil, fullTx, false)
			if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, "couldn't build the pending block response", err, true)
			}

			return rpcBlock, nil
		}
		var err error
		blockNumber, rpcErr := number.GetNumericBlockNumber(ctx, e.state, e.etherman, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		block, err := e.state.GetL2BlockByNumber(ctx, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load block from state by number %v", blockNumber), err, true)
		}

		txs := block.Transactions()
		receipts := make([]ethTypes.Receipt, 0, len(txs))
		for _, tx := range txs {
			receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
			if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load receipt for tx %v", tx.Hash().String()), err, true)
			}
			receipts = append(receipts, *receipt)
		}

		rpcBlock, err := types.NewBlock(block, receipts, fullTx, false)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't build block response for block by number %v", blockNumber), err, true)
		}

		return rpcBlock, nil
	})
}

// GetCode returns account code at given block number
func (e *EthEndpoints) GetCode(address types.ArgAddress, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		var err error
		block, rpcErr := e.getBlockByArg(ctx, blockArg, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		code, err := e.state.GetCode(ctx, address.Address(), block.Root())
		if errors.Is(err, state.ErrNotFound) {
			return "0x", nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get code", err, true)
		}

		return types.ArgBytes(code), nil
	})
}

// GetCompilers eth_getCompilers
func (e *EthEndpoints) GetCompilers() (interface{}, types.Error) {
	return []interface{}{}, nil
}

// GetFilterChanges polling method for a filter, which returns
// an array of logs which occurred since last poll.
func (e *EthEndpoints) GetFilterChanges(filterID string) (interface{}, types.Error) {
	filter, err := e.storage.GetFilter(filterID)
	if errors.Is(err, ErrNotFound) {
		return RPCErrorResponse(types.DefaultErrorCode, "filter not found", err, false)
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get filter from storage", err, true)
	}

	switch filter.Type {
	case FilterTypeBlock:
		{
			res, err := e.state.GetL2BlockHashesSince(context.Background(), filter.LastPoll, nil)
			if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, "failed to get block hashes", err, true)
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
				return RPCErrorResponse(types.DefaultErrorCode, "failed to get pending transaction hashes", err, true)
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
			res := resInterface.([]types.Log)
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
func (e *EthEndpoints) GetFilterLogs(filterID string) (interface{}, types.Error) {
	filter, err := e.storage.GetFilter(filterID)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get filter from storage", err, true)
	}

	if filter.Type != FilterTypeLog {
		return nil, nil
	}

	filterParameters := filter.Parameters.(LogFilter)
	filterParameters.Since = nil

	return e.GetLogs(filterParameters)
}

// GetLogs returns a list of logs accordingly to the provided filter
func (e *EthEndpoints) GetLogs(filter LogFilter) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		return e.internalGetLogs(ctx, dbTx, filter)
	})
}

func (e *EthEndpoints) internalGetLogs(ctx context.Context, dbTx pgx.Tx, filter LogFilter) (interface{}, types.Error) {
	fromBlockNumber, toBlockNumber, rpcErr := filter.GetNumericBlockNumbers(ctx, e.cfg, e.state, e.etherman, dbTx)
	if rpcErr != nil {
		return nil, rpcErr
	}

	var err error
	logs, err := e.state.GetLogs(ctx, fromBlockNumber, toBlockNumber, filter.Addresses, filter.Topics, filter.BlockHash, filter.Since, dbTx)
	if errors.Is(err, state.ErrMaxLogsCountLimitExceeded) {
		errMsg := fmt.Sprintf(state.ErrMaxLogsCountLimitExceeded.Error(), e.cfg.MaxLogsCount)
		return RPCErrorResponse(types.InvalidParamsErrorCode, errMsg, nil, false)
	} else if errors.Is(err, state.ErrMaxLogsBlockRangeLimitExceeded) {
		errMsg := fmt.Sprintf(state.ErrMaxLogsBlockRangeLimitExceeded.Error(), e.cfg.MaxLogsBlockRange)
		return RPCErrorResponse(types.InvalidParamsErrorCode, errMsg, nil, false)
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get logs from state", err, true)
	}

	result := make([]types.Log, 0, len(logs))
	for _, l := range logs {
		result = append(result, types.NewLog(*l))
	}

	return result, nil
}

// GetStorageAt gets the value stored for an specific address and position
func (e *EthEndpoints) GetStorageAt(address types.ArgAddress, storageKeyStr string, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	storageKey := types.ArgHash{}
	err := storageKey.UnmarshalText([]byte(storageKeyStr))
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "unable to decode storage key: hex string invalid", nil, false)
	}

	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		block, respErr := e.getBlockByArg(ctx, blockArg, dbTx)
		if respErr != nil {
			return nil, respErr
		}

		value, err := e.state.GetStorageAt(ctx, address.Address(), storageKey.Hash().Big(), block.Root())
		if errors.Is(err, state.ErrNotFound) {
			return types.ArgBytesPtr(common.Hash{}.Bytes()), nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get storage value from state", err, true)
		}

		return types.ArgBytesPtr(common.BigToHash(value).Bytes()), nil
	})
}

// GetTransactionByBlockHashAndIndex returns information about a transaction by
// block hash and transaction index position.
func (e *EthEndpoints) GetTransactionByBlockHashAndIndex(hash types.ArgHash, index types.Index) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		tx, err := e.state.GetTransactionByL2BlockHashAndIndex(ctx, hash.Hash(), uint64(index), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get transaction", err, true)
		}

		receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get transaction receipt", err, true)
		}

		res, err := types.NewTransaction(*tx, receipt, false)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to build transaction response", err, true)
		}

		return res, nil
	})
}

// GetTransactionByBlockNumberAndIndex returns information about a transaction by
// block number and transaction index position.
func (e *EthEndpoints) GetTransactionByBlockNumberAndIndex(number *types.BlockNumber, index types.Index) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		var err error
		blockNumber, rpcErr := number.GetNumericBlockNumber(ctx, e.state, e.etherman, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		tx, err := e.state.GetTransactionByL2BlockNumberAndIndex(ctx, blockNumber, uint64(index), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get transaction", err, true)
		}

		receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get transaction receipt", err, true)
		}

		res, err := types.NewTransaction(*tx, receipt, false)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to build transaction response", err, true)
		}

		return res, nil
	})
}

// GetTransactionByHash returns a transaction by his hash
func (e *EthEndpoints) GetTransactionByHash(hash types.ArgHash) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		// try to get tx from state
		tx, err := e.state.GetTransactionByHash(ctx, hash.Hash(), dbTx)
		if err != nil && !errors.Is(err, state.ErrNotFound) {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to load transaction by hash from state", err, true)
		}
		if tx != nil {
			receipt, err := e.state.GetTransactionReceipt(ctx, hash.Hash(), dbTx)
			if errors.Is(err, state.ErrNotFound) {
				return RPCErrorResponse(types.DefaultErrorCode, "transaction receipt not found", err, false)
			} else if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, "failed to load transaction receipt from state", err, true)
			}

			res, err := types.NewTransaction(*tx, receipt, false)
			if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, "failed to build transaction response", err, true)
			}

			return res, nil
		}

		// if the tx does not exist in the state, look for it in the pool
		if e.cfg.SequencerNodeURI != "" {
			return e.getTransactionByHashFromSequencerNode(hash.Hash())
		}
		poolTx, err := e.pool.GetTxByHash(ctx, hash.Hash())
		if errors.Is(err, pool.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to load transaction by hash from pool", err, true)
		}
		if poolTx.Status == pool.TxStatusPending {
			tx = &poolTx.Transaction
			res, err := types.NewTransaction(*tx, nil, false)
			if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, "failed to build transaction response", err, true)
			}
			return res, nil
		}
		return nil, nil
	})
}

func (e *EthEndpoints) getTransactionByHashFromSequencerNode(hash common.Hash) (interface{}, types.Error) {
	res, err := client.JSONRPCCall(e.cfg.SequencerNodeURI, "eth_getTransactionByHash", hash.String())
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get tx from sequencer node", err, true)
	}

	if res.Error != nil {
		return RPCErrorResponse(res.Error.Code, res.Error.Message, nil, false)
	}

	var tx *types.Transaction
	err = json.Unmarshal(res.Result, &tx)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to read tx from sequencer node", err, true)
	}
	return tx, nil
}

// GetTransactionCount returns account nonce
func (e *EthEndpoints) GetTransactionCount(address types.ArgAddress, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		var (
			pendingNonce uint64
			nonce        uint64
			err          error
		)

		block, respErr := e.getBlockByArg(ctx, blockArg, dbTx)
		if respErr != nil {
			return nil, respErr
		}

		if blockArg != nil {
			blockNumArg := blockArg.Number()
			if blockNumArg != nil && *blockNumArg == types.PendingBlockNumber {
				if e.cfg.SequencerNodeURI != "" {
					return e.getTransactionCountFromSequencerNode(address.Address(), blockArg.Number())
				}
				pendingNonce, err = e.pool.GetNonce(ctx, address.Address())
				if err != nil {
					return RPCErrorResponse(types.DefaultErrorCode, "failed to count pending transactions", err, true)
				}
			}
		}

		nonce, err = e.state.GetNonce(ctx, address.Address(), block.Root())

		if errors.Is(err, state.ErrNotFound) {
			return hex.EncodeUint64(0), nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to count transactions", err, true)
		}

		if pendingNonce > nonce {
			nonce = pendingNonce
		}

		return hex.EncodeUint64(nonce), nil
	})
}

func (e *EthEndpoints) getTransactionCountFromSequencerNode(address common.Address, number *types.BlockNumber) (interface{}, types.Error) {
	res, err := client.JSONRPCCall(e.cfg.SequencerNodeURI, "eth_getTransactionCount", address.String(), number.StringOrHex())
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get nonce from sequencer node", err, true)
	}

	if res.Error != nil {
		return RPCErrorResponse(res.Error.Code, res.Error.Message, nil, false)
	}

	var nonce types.ArgUint64
	err = json.Unmarshal(res.Result, &nonce)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to read nonce from sequencer node", err, true)
	}
	return nonce, nil
}

// GetBlockTransactionCountByHash returns the number of transactions in a
// block from a block matching the given block hash.
func (e *EthEndpoints) GetBlockTransactionCountByHash(hash types.ArgHash) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		c, err := e.state.GetL2BlockTransactionCountByHash(ctx, hash.Hash(), dbTx)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to count transactions", err, true)
		}

		return types.ArgUint64(c), nil
	})
}

// GetBlockTransactionCountByNumber returns the number of transactions in a
// block from a block matching the given block number.
func (e *EthEndpoints) GetBlockTransactionCountByNumber(number *types.BlockNumber) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		if number != nil && *number == types.PendingBlockNumber {
			if e.cfg.SequencerNodeURI != "" {
				return e.getBlockTransactionCountByNumberFromSequencerNode(number)
			}
			c, err := e.pool.CountPendingTransactions(ctx)
			if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, "failed to count pending transactions", err, true)
			}
			return types.ArgUint64(c), nil
		}

		var err error
		blockNumber, rpcErr := number.GetNumericBlockNumber(ctx, e.state, e.etherman, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		c, err := e.state.GetL2BlockTransactionCountByNumber(ctx, blockNumber, dbTx)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to count transactions", err, true)
		}

		return types.ArgUint64(c), nil
	})
}

func (e *EthEndpoints) getBlockTransactionCountByNumberFromSequencerNode(number *types.BlockNumber) (interface{}, types.Error) {
	res, err := client.JSONRPCCall(e.cfg.SequencerNodeURI, "eth_getBlockTransactionCountByNumber", number.StringOrHex())
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get tx count by block number from sequencer node", err, true)
	}

	if res.Error != nil {
		return RPCErrorResponse(res.Error.Code, res.Error.Message, nil, false)
	}

	var count types.ArgUint64
	err = json.Unmarshal(res.Result, &count)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to read tx count by block number from sequencer node", err, true)
	}
	return count, nil
}

// GetTransactionReceipt returns a transaction receipt by his hash
func (e *EthEndpoints) GetTransactionReceipt(hash types.ArgHash) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		tx, err := e.state.GetTransactionByHash(ctx, hash.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get tx from state", err, true)
		}

		r, err := e.state.GetTransactionReceipt(ctx, hash.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, nil
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get tx receipt from state", err, true)
		}

		receipt, err := types.NewReceipt(*tx, r)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to build the receipt response", err, true)
		}

		return receipt, nil
	})
}

// NewBlockFilter creates a filter in the node, to notify when
// a new block arrives. To check if the state has changed,
// call eth_getFilterChanges.
func (e *EthEndpoints) NewBlockFilter() (interface{}, types.Error) {
	return e.newBlockFilter(nil)
}

// internal
func (e *EthEndpoints) newBlockFilter(wsConn *atomic.Pointer[websocket.Conn]) (interface{}, types.Error) {
	id, err := e.storage.NewBlockFilter(wsConn)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to create new block filter", err, true)
	}

	return id, nil
}

// NewFilter creates a filter object, based on filter options,
// to notify when the state changes (logs). To check if the state
// has changed, call eth_getFilterChanges.
func (e *EthEndpoints) NewFilter(filter LogFilter) (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		return e.newFilter(ctx, nil, filter, dbTx)
	})
}

// internal
func (e *EthEndpoints) newFilter(ctx context.Context, wsConn *atomic.Pointer[websocket.Conn], filter LogFilter, dbTx pgx.Tx) (interface{}, types.Error) {
	if filter.ShouldFilterByBlockRange() {
		_, _, rpcErr := filter.GetNumericBlockNumbers(ctx, e.cfg, e.state, e.etherman, nil)
		if rpcErr != nil {
			return nil, rpcErr
		}
	}

	id, err := e.storage.NewLogFilter(wsConn, filter)
	if errors.Is(err, ErrFilterInvalidPayload) {
		return RPCErrorResponse(types.InvalidParamsErrorCode, err.Error(), nil, false)
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to create new log filter", err, true)
	}

	return id, nil
}

// NewPendingTransactionFilter creates a filter in the node, to
// notify when new pending transactions arrive. To check if the
// state has changed, call eth_getFilterChanges.
func (e *EthEndpoints) NewPendingTransactionFilter() (interface{}, types.Error) {
	return e.newPendingTransactionFilter(nil)
}

// internal
func (e *EthEndpoints) newPendingTransactionFilter(wsConn *atomic.Pointer[websocket.Conn]) (interface{}, types.Error) {
	return nil, types.NewRPCError(types.DefaultErrorCode, "not supported yet")
	// id, err := e.storage.NewPendingTransactionFilter(wsConn)
	// if err != nil {
	// 	return rpcErrorResponse(types.DefaultErrorCode, "failed to create new pending transaction filter", err)
	// }

	// return id, nil
}

// SendRawTransaction has two different ways to handle new transactions:
// - for Sequencer nodes it tries to add the tx to the pool
// - for Non-Sequencer nodes it relays the Tx to the Sequencer node
func (e *EthEndpoints) SendRawTransaction(httpRequest *http.Request, input string) (interface{}, types.Error) {
	if e.cfg.SequencerNodeURI != "" {
		return e.relayTxToSequencerNode(input)
	} else {
		ip := ""
		ips := httpRequest.Header.Get("X-Forwarded-For")

		// TODO: this is temporary patch remove this log
		realIp := httpRequest.Header.Get("X-Real-IP")
		log.Infof("X-Forwarded-For: %s, X-Real-IP: %s", ips, realIp)

		if ips != "" {
			ip = strings.Split(ips, ",")[0]
		}

		return e.tryToAddTxToPool(input, ip)
	}
}

func (e *EthEndpoints) relayTxToSequencerNode(input string) (interface{}, types.Error) {
	res, err := client.JSONRPCCall(e.cfg.SequencerNodeURI, "eth_sendRawTransaction", input)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to relay tx to the sequencer node", err, true)
	}

	if res.Error != nil {
		return RPCErrorResponse(res.Error.Code, res.Error.Message, nil, false)
	}

	txHash := res.Result

	return txHash, nil
}

func (e *EthEndpoints) tryToAddTxToPool(input, ip string) (interface{}, types.Error) {
	tx, err := hexToTx(input)
	if err != nil {
		return RPCErrorResponse(types.InvalidParamsErrorCode, "invalid tx input", err, false)
	}
	log.Infof("adding TX to the pool: %v", tx.Hash().Hex())
	if err := e.pool.AddTx(context.Background(), *tx, ip); err != nil {
		// it's not needed to log the error here, because we check and log if needed
		// for each specific case during the "pool.AddTx" internal steps
		return RPCErrorResponse(types.DefaultErrorCode, err.Error(), nil, false)
	}
	log.Infof("TX added to the pool: %v", tx.Hash().Hex())

	return tx.Hash().Hex(), nil
}

// UninstallFilter uninstalls a filter with given id.
func (e *EthEndpoints) UninstallFilter(filterID string) (interface{}, types.Error) {
	err := e.storage.UninstallFilter(filterID)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to uninstall filter", err, true)
	}

	return true, nil
}

// Syncing returns an object with data about the sync status or false.
// https://eth.wiki/json-rpc/API#eth_syncing
func (e *EthEndpoints) Syncing() (interface{}, types.Error) {
	return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		_, err := e.state.GetLastL2BlockNumber(ctx, dbTx)
		if errors.Is(err, state.ErrStateNotSynchronized) {
			return nil, types.NewRPCErrorWithData(types.DefaultErrorCode, state.ErrStateNotSynchronized.Error(), nil)
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get last block number from state", err, true)
		}

		syncInfo, err := e.state.GetSyncingInfo(ctx, dbTx)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get syncing info from state", err, true)
		}

		if syncInfo.CurrentBlockNumber >= syncInfo.LastBlockNumberSeen {
			return false, nil
		}

		return struct {
			S types.ArgUint64 `json:"startingBlock"`
			C types.ArgUint64 `json:"currentBlock"`
			H types.ArgUint64 `json:"highestBlock"`
		}{
			S: types.ArgUint64(syncInfo.InitialSyncingBlock),
			C: types.ArgUint64(syncInfo.CurrentBlockNumber),
			H: types.ArgUint64(syncInfo.LastBlockNumberSeen),
		}, nil
	})
}

// GetUncleByBlockHashAndIndex returns information about a uncle of a
// block by hash and uncle index position
func (e *EthEndpoints) GetUncleByBlockHashAndIndex(hash types.ArgHash, index types.Index) (interface{}, types.Error) {
	return nil, nil
}

// GetUncleByBlockNumberAndIndex returns information about a uncle of a
// block by number and uncle index position
func (e *EthEndpoints) GetUncleByBlockNumberAndIndex(number types.BlockNumber, index types.Index) (interface{}, types.Error) {
	return nil, nil
}

// GetUncleCountByBlockHash returns the number of uncles in a block
// matching the given block hash
func (e *EthEndpoints) GetUncleCountByBlockHash(hash types.ArgAddress) (interface{}, types.Error) {
	return "0x0", nil
}

// GetUncleCountByBlockNumber returns the number of uncles in a block
// matching the given block number
func (e *EthEndpoints) GetUncleCountByBlockNumber(number types.BlockNumber) (interface{}, types.Error) {
	return "0x0", nil
}

// ProtocolVersion returns the protocol version.
func (e *EthEndpoints) ProtocolVersion() (interface{}, types.Error) {
	return "0x0", nil
}

func hexToTx(str string) (*ethTypes.Transaction, error) {
	tx := new(ethTypes.Transaction)

	b, err := hex.DecodeHex(str)
	if err != nil {
		return nil, err
	}

	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return tx, nil
}

func (e *EthEndpoints) updateFilterLastPoll(filterID string) types.Error {
	err := e.storage.UpdateFilterLastPoll(filterID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return types.NewRPCError(types.DefaultErrorCode, "failed to update last time the filter changes were requested")
	}
	return nil
}

// Subscribe Creates a new subscription over particular events.
// The node will return a subscription id.
// For each event that matches the subscription a notification with relevant
// data is sent together with the subscription id.
func (e *EthEndpoints) Subscribe(wsConn *atomic.Pointer[websocket.Conn], name string, logFilter *LogFilter) (interface{}, types.Error) {
	switch name {
	case "newHeads":
		return e.newBlockFilter(wsConn)
	case "logs":
		return e.txMan.NewDbTxScope(e.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
			var lf LogFilter
			if logFilter != nil {
				lf = *logFilter
			}
			return e.newFilter(ctx, wsConn, lf, dbTx)
		})
	case "pendingTransactions", "newPendingTransactions":
		return e.newPendingTransactionFilter(wsConn)
	case "syncing":
		return nil, types.NewRPCError(types.DefaultErrorCode, "not supported yet")
	default:
		return nil, types.NewRPCError(types.DefaultErrorCode, "invalid filter name")
	}
}

// Unsubscribe uninstalls the filter based on the provided filterID
func (e *EthEndpoints) Unsubscribe(wsConn *websocket.Conn, filterID string) (interface{}, types.Error) {
	return e.UninstallFilter(filterID)
}

// uninstallFilterByWSConn uninstalls the filters connected to the
// provided web socket connection
func (e *EthEndpoints) uninstallFilterByWSConn(wsConn *atomic.Pointer[websocket.Conn]) error {
	return e.storage.UninstallFilterByWSConn(wsConn)
}

// onNewL2Block is triggered when the state triggers the event for a new l2 block
func (e *EthEndpoints) onNewL2Block(event state.NewL2BlockEvent) {
	log.Debugf("[onNewL2Block] new l2 block event detected for block %v", event.Block.NumberU64())
	wg := sync.WaitGroup{}

	wg.Add(1)
	go e.notifyNewHeads(&wg, event)

	wg.Add(1)
	go e.notifyNewLogs(&wg, event)

	wg.Wait()
}

func (e *EthEndpoints) notifyNewHeads(wg *sync.WaitGroup, event state.NewL2BlockEvent) {
	defer wg.Done()
	start := time.Now()
	blockFilters, err := e.storage.GetAllBlockFiltersWithWSConn()
	if err != nil {
		log.Errorf("failed to get all block filters with web sockets connections: %v", err)
	} else {
		b, err := types.NewBlock(&event.Block, nil, false, false)
		if err != nil {
			log.Errorf("failed to build block response to subscription: %v", err)
			return
		}
		data, err := json.Marshal(b)
		if err != nil {
			log.Errorf("failed to marshal block response to subscription: %v", err)
			return
		}
		for _, filter := range blockFilters {
			filter.EnqueueSubscriptionDataToBeSent(data)
			go filter.SendEnqueuedSubscriptionData()
		}
	}
	log.Debugf("[notifyNewHeads] new l2 block event for block %v took %v to send all the messages for block filters", event.Block.NumberU64(), time.Since(start))
}

func (e *EthEndpoints) notifyNewLogs(wg *sync.WaitGroup, event state.NewL2BlockEvent) {
	defer wg.Done()
	start := time.Now()
	logFilters, err := e.storage.GetAllLogFiltersWithWSConn()
	if err != nil {
		log.Errorf("failed to get all log filters with web sockets connections: %v", err)
	} else {
		for _, filter := range logFilters {
			filterParameters := filter.Parameters.(LogFilter)
			bn := types.BlockNumber(event.Block.NumberU64())

			// if from and to blocks are new, set it to the current block to make
			// the query faster
			if filterParameters.FromBlock == nil && filterParameters.ToBlock == nil {
				filterParameters.FromBlock = &bn
				filterParameters.ToBlock = &bn
			} else {
				// if the filter has a fromBlock value set
				// and the event block number is smaller than the
				// from block, skip this filter
				if filterParameters.FromBlock != nil {
					fromBlock, rpcErr := filterParameters.FromBlock.GetNumericBlockNumber(context.Background(), e.state, e.etherman, nil)
					if rpcErr != nil {
						log.Errorf(rpcErr.Error(), filter.ID, err)
						continue
					}
					if fromBlock > event.Block.NumberU64() {
						continue
					}
					// otherwise set the from block to a fixed number
					// to avoid querying it again in the next step
					fixedFromBlock := types.BlockNumber(fromBlock)
					filterParameters.FromBlock = &fixedFromBlock
				}

				// if the filter has a toBlock value set
				// and the event block number is greater than the
				// to block, skip this filter
				if filterParameters.ToBlock != nil {
					toBlock, rpcErr := filterParameters.ToBlock.GetNumericBlockNumber(context.Background(), e.state, e.etherman, nil)
					if rpcErr != nil {
						log.Errorf(rpcErr.Error(), filter.ID, err)
						continue
					}
					if toBlock > event.Block.NumberU64() {
						continue
					}
					// otherwise set the to block to a fixed number
					// to avoid querying it again in the next step
					fixedToBlock := types.BlockNumber(toBlock)
					filterParameters.ToBlock = &fixedToBlock
				}
			}

			// get new logs for this specific filter
			changes, err := e.internalGetLogs(context.Background(), nil, filterParameters)
			if errors.Is(err, state.ErrMaxLogsCountLimitExceeded) {
				log.Infof("failed to get filters changes for filter %v, the filter seems to be returning more results than allowed and was removed: %v", filter.ID, err)
				err := e.storage.UninstallFilter(filter.ID)
				if !errors.Is(err, ErrNotFound) && err != nil {
					log.Errorf("failed to automatically uninstall filter %v: %v", filter.ID, err)
				} else {
					log.Infof("Filter %v automatically uninstalled", filter.ID)
				}
				continue
			} else if err != nil {
				log.Errorf("failed to get filters changes for filter %v with web sockets connections: %v", filter.ID, err)
				continue
			}

			// if there are new logs for the filter, send it
			if changes != nil {
				ethLogs := changes.([]types.Log)
				for _, ethLog := range ethLogs {
					data, err := json.Marshal(ethLog)
					if err != nil {
						log.Errorf("failed to marshal ethLog response to subscription: %v", err)
					}
					filter.EnqueueSubscriptionDataToBeSent(data)
					go filter.SendEnqueuedSubscriptionData()
				}
			}
		}
	}
	log.Debugf("[notifyNewLogs] new l2 block event for block %v took %v to send all the messages for log filters", event.Block.NumberU64(), time.Since(start))
}
