package jsonrpc

import (
	"context"
	"errors"
	"fmt"

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
	state            state.State
	sequencerAddress common.Address
	gpe              gasPriceEstimator
}

type blockNumberOrHash struct {
	BlockNumber *BlockNumber `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash `json:"blockHash,omitempty"`
}

// BlockNumber returns current block number
func (e *Eth) BlockNumber() (interface{}, error) {
	ctx := context.Background()

	lastBatchNumber, err := e.state.GetLastBatchNumber(ctx)
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
			return nil, fmt.Errorf("failed to get header from block hash or block number")
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
	lastVirtualBatch, err := e.state.GetLastBatch(ctx, true)
	if err != nil {
		return nil, err
	}
	bp, err := e.state.NewBatchProcessor(ctx, e.sequencerAddress, lastVirtualBatch.Number().Uint64())
	if err != nil {
		return nil, err
	}

	result := bp.ProcessUnsignedTransaction(ctx, tx, *arg.From, e.sequencerAddress)

	if result.Failed() {
		return nil, fmt.Errorf("unable to execute call: %w", result.Err)
	}

	return argBytesPtr(result.ReturnValue), nil
}

// ChainId returns the chain id of the client
func (e *Eth) ChainId() (interface{}, error) { //nolint:golint
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
	gasEstimation := e.state.EstimateGas(tx)
	return hex.EncodeUint64(gasEstimation), nil
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

	balance, err := e.state.GetBalance(ctx, address, batchNumber)
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

	batch, err := e.state.GetBatchByHash(ctx, hash)
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

	batchNumber, err := e.getNumericBlockNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	batch, err := e.state.GetBatchByNumber(ctx, batchNumber)
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

	code, err := e.state.GetCode(ctx, address, batchNumber)
	if errors.Is(err, state.ErrNotFound) {
		return "0x", nil
	} else if err != nil {
		return nil, err
	}

	return argBytes(code), nil
}

// GetLogs returns a list of logs accordingly to the provided filter
func (e *Eth) GetLogs(filter *LogFilter) (interface{}, error) {
	ctx := context.Background()

	fromBlock, err := e.getNumericBlockNumber(ctx, filter.fromBlock)
	if err != nil {
		return nil, err
	}

	toBlock, err := e.getNumericBlockNumber(ctx, filter.toBlock)
	if err != nil {
		return nil, err
	}

	logs, err := e.state.GetLogs(ctx, fromBlock, toBlock, filter.Addresses, filter.Topics, filter.BlockHash)
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

	value, err := e.state.GetStorageAt(ctx, address, position, batchNumber)
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

	tx, err := e.state.GetTransactionByBatchHashAndIndex(ctx, hash, uint64(index))
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash())
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

	tx, err := e.state.GetTransactionByBatchNumberAndIndex(ctx, batchNumber, uint64(index))
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash())
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

	tx, err := e.state.GetTransactionByHash(ctx, hash)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	receipt, err := e.state.GetTransactionReceipt(ctx, tx.Hash())
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

	nonce, err := e.state.GetNonce(ctx, address, batchNumber)
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

	r, err := e.state.GetTransactionReceipt(ctx, hash)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return stateReceiptToRPCReceipt(r), nil
}

// SendRawTransaction sends a raw transaction
func (e *Eth) SendRawTransaction(input string) (interface{}, error) {
	tx, err := hexToTx(input)
	if err != nil {
		log.Warnf("Invalid tx: %v", err)
		return nil, err
	}

	log.Debugf("checking TX signature: %v", tx.Hash().Hex())
	if err := state.CheckSignature(tx); err != nil {
		log.Warnf("Invalid signature[%v]: %v", tx.Hash().Hex(), err)
		return nil, err
	}
	log.Debugf("TX signature OK: %v", tx.Hash().Hex())

	log.Debugf("adding TX to the pool: %v", tx.Hash().Hex())
	if err := e.pool.AddTx(context.Background(), *tx); err != nil {
		log.Warnf("Failed to add TX to the pool[%v]: %v", tx.Hash().Hex(), err)
		return nil, err
	}
	log.Debugf("TX added to the pool: %v", tx.Hash().Hex())

	return tx.Hash().Hex(), nil
}

func (e *Eth) getNumericBlockNumber(ctx context.Context, number BlockNumber) (uint64, error) {
	switch number {
	case LatestBlockNumber:
		lastBatchNumber, err := e.state.GetLastBatchNumber(ctx)
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

func (e *Eth) getHeaderFromBlockNumberOrHash(bnh *blockNumberOrHash) (*types.Header, error) {
	var (
		header *types.Header
		err    error
	)

	if bnh.BlockNumber != nil {
		header, err = e.getBatchHeader(*bnh.BlockNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get the header of block %d: %w", *bnh.BlockNumber, err)
		}
	} else if bnh.BlockHash != nil {
		block, err := e.state.GetBatchByHash(context.Background(), *bnh.BlockHash)
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
		batch, err := e.state.GetLastBatch(context.Background(), false)
		if err != nil {
			return nil, err
		}
		return batch.Header, nil

	case EarliestBlockNumber:
		return e.state.GetBatchHeader(context.Background(), uint64(0))

	case PendingBlockNumber:
		return nil, fmt.Errorf("fetching the pending header is not supported")

	default:
		return e.state.GetBatchHeader(context.Background(), uint64(number))
	}
}
