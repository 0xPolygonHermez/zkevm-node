package state

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4/pgxpool"
)

// State is the interface of the Hermez state
type State interface {
	NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) BatchProcessor
	GetStateRoot(virtual bool) (*big.Int, error)
	GetBalance(address common.Address, batchNumber uint64) (*big.Int, error)
	EstimateGas(transaction types.Transaction) uint64
	GetLastBlock(ctx context.Context) (*types.Block, error)
	GetPreviousBlock(ctx context.Context, offset uint64) (*types.Block, error)
	GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error)
	GetLastBlockNumber(ctx context.Context) (uint64, error)
	GetLastBatch(ctx context.Context, isVirtual bool) (*Batch, error)
	GetTransaction(ctx context.Context, hash common.Hash) (*types.Transaction, error)
	GetNonce(address common.Address, batchNumber uint64) (uint64, error)
	GetPreviousBatch(ctx context.Context, offset uint64) (*Batch, error)
	GetBatchByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64) (*types.Block, error)
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error)
	GetTransactionCount(ctx context.Context, address common.Address) (uint64, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*types.Receipt, error)
	Reset(blockNumber uint64) error
	ConsolidateBatch(batchNumber uint64) error
	GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error)
}

const (
	getLastBlockSQL         = "SELECT * FROM block ORDER BY received_at DESC LIMIT 1"
	getPreviousBlockSQL     = "SELECT * FROM block ORDER BY received_at DESC LIMIT 1 OFFSET $1"
	getBlockByHashSQL       = "SELECT * FROM block WHERE hash = $1"
	getBlockByNumberSQL     = "SELECT * FROM block WHERE eth_block_num = $1"
	getLastBlockNumberSQL   = "SELECT eth_block_num FROM block ORDER BY received_at DESC LIMIT 1"
	getLastBatchSQL         = "SELECT * FROM batch ORDER BY batch_num DESC LIMIT 1"
	getPreviousBatchSQL     = "SELECT * FROM batch ORDER BY batch_num DESC LIMIT 1 OFFSET $1"
	getTransactionSQL       = "SELECT * FROM transaction WHERE hash = $1"
	getBatchByHashSQL       = "SELECT * FROM batch WHERE hash = $1"
	getBatchByNumberSQL     = "SELECT * FROM batch WHERE batch_num = $1"
	getLastBatchNumberSQL   = "SELECT batch_num FROM batch ORDER BY batch_num DESC LIMIT 1"
	getTransactionByHashSQL = "SELECT * FROM transaction WHERE hash = $1"
	// todo: change to jsonb
	getTransactionCountSQL = "SELECT COUNT(*) FROM transaction WHERE decoded.from = $1"
)

// BasicState is a implementation of the state
type BasicState struct {
	db *pgxpool.Pool
	// StateTree merkletree.Merkletree
}

// NewState creates a new State
func NewState(db *pgxpool.Pool) State {
	// return &State{StateTree: merkletree.NewMerkletree(db)}
	return &BasicState{db: db}
}

// NewBatchProcessor creates a new batch processor
func (s *BasicState) NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) BatchProcessor {
	return &BasicBatchProcessor{State: s}
}

// GetStateRoot returns the root of the state tree
func (s *BasicState) GetStateRoot(virtual bool) (*big.Int, error) {
	// return s.StateTree.Root, nil
	return nil, nil
}

// GetBalance from a given address
func (s *BasicState) GetBalance(address common.Address, batchNumber uint64) (*big.Int, error) {
	/*
		key, err := leafs.NewBalanceKey(common.BytesToAddress(address.Bytes()))
		if err != nil {
			return nil, err
		}

		balanceBytes, err := s.StateTree.Get(s.StateTree.Root, key)
		if err != nil {
			return nil, err
		}

		return leafs.BytesToBalance(balanceBytes), nil*/
	return nil, nil
}

// EstimateGas for a transaction
func (s *BasicState) EstimateGas(transaction types.Transaction) uint64 {
	// TODO: Calculate once we have txs that interact with SCs
	return 21000 //nolint:gomnd
}

// GetLastBlock gets the latest block
func (s *BasicState) GetLastBlock(ctx context.Context) (*types.Block, error) {
	var res *types.Block
	err := s.db.QueryRow(ctx, getLastBlockSQL).Scan(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetPreviousBlock gets the offset previous block respect to latest
func (s *BasicState) GetPreviousBlock(ctx context.Context, offset uint64) (*types.Block, error) {
	var res *types.Block
	err := s.db.QueryRow(ctx, getPreviousBlockSQL, offset).Scan(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetBlockByHash gets the block with the required hash
func (s *BasicState) GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	var res *types.Block
	err := s.db.QueryRow(ctx, getBlockByHashSQL, hash).Scan(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetBlockByNumber gets the block with the required number
func (s *BasicState) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	var res *types.Block
	err := s.db.QueryRow(ctx, getBlockByNumberSQL, blockNumber).Scan(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetLastBlockNumber gets the latest block number
func (s *BasicState) GetLastBlockNumber(ctx context.Context) (uint64, error) {
	var lastBlockNum uint64
	err := s.db.QueryRow(ctx, getLastBlockNumberSQL).Scan(&lastBlockNum)
	if err != nil {
		return 0, err
	}
	return lastBlockNum, nil
}

// GetLastBatch gets the latest batch
func (s *BasicState) GetLastBatch(ctx context.Context, isVirtual bool) (*Batch, error) {
	var batch *Batch
	err := s.db.QueryRow(ctx, getLastBatchSQL).Scan(&batch)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

// GetTransaction gets a transactions by its hash
func (s *BasicState) GetTransaction(ctx context.Context, hash common.Hash) (*types.Transaction, error) {
	var tx *types.Transaction
	err := s.db.QueryRow(ctx, getTransactionSQL, hash).Scan(&tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetNonce returns the nonce of the given account at the given batch number
func (s *BasicState) GetNonce(address common.Address, batchNumber uint64) (uint64, error) {
	panic("not implemented yet")
}

// GetPreviousBatch gets the offset previous batch respect to latest
func (s *BasicState) GetPreviousBatch(ctx context.Context, offset uint64) (*Batch, error) {
	var batch *Batch
	err := s.db.QueryRow(ctx, getPreviousBatchSQL, offset).Scan(&batch)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

// GetBatchByHash gets the batch with the required hash
func (s *BasicState) GetBatchByHash(ctx context.Context, hash common.Hash) (*Batch, error) {
	var batch *Batch
	err := s.db.QueryRow(ctx, getBatchByHashSQL, hash).Scan(&batch)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

// GetBatchByNumber gets the batch with the required number
func (s *BasicState) GetBatchByNumber(ctx context.Context, batchNumber uint64) (*Batch, error) {
	var batch *Batch
	err := s.db.QueryRow(ctx, getBatchByNumberSQL, batchNumber).Scan(&batch)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

// GetLastBatchNumber gets the latest batch number
func (s *BasicState) GetLastBatchNumber(ctx context.Context) (uint64, error) {
	var lastBatchNumber uint64
	err := s.db.QueryRow(ctx, getLastBatchNumberSQL).Scan(&lastBatchNumber)
	if err != nil {
		return 0, err
	}
	return lastBatchNumber, nil
}

// GetTransactionByBatchHashAndIndex gets a transaction from a batch by index
func (s *BasicState) GetTransactionByBatchHashAndIndex(batchHash common.Hash, index uint64) (*types.Transaction, error) {
	panic("not implemented")
}

// GetTransactionByBatchNumberAndIndex gets a transaction from a batch by index
func (s *BasicState) GetTransactionByBatchNumberAndIndex(batchNumber uint64, index uint64) (*types.Transaction, error) {
	panic("not implemented")
}

// GetTransactionByHash gets a transaction by its hash
func (s *BasicState) GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error) {
	var tx *types.Transaction
	err := s.db.QueryRow(ctx, getTransactionByHashSQL, transactionHash).Scan(&tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetTransactionCount returns the number of transactions sent from an address
func (s *BasicState) GetTransactionCount(ctx context.Context, address common.Address) (uint64, error) {
	var count uint64
	err := s.db.QueryRow(ctx, getTransactionCountSQL, address).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
func (s *BasicState) GetTransactionReceipt(transactionHash common.Hash) (*types.Receipt, error) {
	return nil, nil
}

// Reset resets the state to a block
func (s *BasicState) Reset(blockNumber uint64) error {
	return nil
}

// ConsolidateBatch changes the virtual status of a batch
func (s *BasicState) ConsolidateBatch(batchNumber uint64) error {
	return nil
}

// GetTxsByBatchNum returns all the txs in a given batch
func (s *BasicState) GetTxsByBatchNum(batchNum uint64) ([]*types.Transaction, error) {
	return nil, nil
}
