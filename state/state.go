package state

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// State is the interface of the Hermez state
type State interface {
	NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) BatchProcessor
	GetStateRoot(ctx context.Context, virtual bool) (*big.Int, error)
	GetBalance(address common.Address, batchNumber uint64) (*big.Int, error)
	EstimateGas(transaction *types.Transaction) uint64
	GetLastBlock(ctx context.Context) (*Block, error)
	GetPreviousBlock(ctx context.Context, offset uint64) (*Block, error)
	GetBlockByHash(ctx context.Context, hash common.Hash) (*Block, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error)
	GetLastBlockNumber(ctx context.Context) (uint64, error)
	GetNonce(address common.Address, batchNumber uint64) (uint64, error)
	GetLastBatch(ctx context.Context, isVirtual bool) (*Batch, error)
	GetPreviousBatch(ctx context.Context, isVirtual bool, offset uint64) (*Batch, error)
	GetBatchByHash(ctx context.Context, hash common.Hash) (*Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64) (*Batch, error)
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error)
	GetTransactionCount(ctx context.Context, address common.Address) (uint64, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*types.Receipt, error)
	Reset(blockNumber uint64) error
	ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash) error
	GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error)
	AddNewSequencer(seq Sequencer) error
	SetGenesis(genesis Genesis) error
	AddBlock(*Block) error
	SetLastBatchNumberSeenOnEthereum(batchNumber uint64) error
	GetLastBatchNumberSeenOnEthereum() (uint64, error)
}

const (
	getLastBlockSQL                 = "SELECT * FROM block ORDER BY block_num DESC LIMIT 1"
	getPreviousBlockSQL             = "SELECT * FROM block ORDER BY block_num DESC LIMIT 1 OFFSET $1"
	getBlockByHashSQL               = "SELECT * FROM block WHERE block_hash = $1"
	getBlockByNumberSQL             = "SELECT * FROM block WHERE block_num = $1"
	getLastBlockNumberSQL           = "SELECT MAX(block_num) FROM block"
	getLastVirtualBatchSQL          = "SELECT * FROM batch ORDER BY batch_num DESC LIMIT 1"
	getLastConsolidatedBatchSQL     = "SELECT * FROM batch WHERE consolidated_tx_hash != $1 ORDER BY batch_num DESC LIMIT 1"
	getPreviousVirtualBatchSQL      = "SELECT * FROM batch ORDER BY batch_num DESC LIMIT 1 OFFSET $1"
	getPreviousConsolidatedBatchSQL = "SELECT * FROM batch WHERE consolidated_tx_hash != $1 ORDER BY batch_num DESC LIMIT 1 OFFSET $2"
	getBatchByHashSQL               = "SELECT * FROM batch WHERE batch_hash = $1"
	getBatchByNumberSQL             = "SELECT * FROM batch WHERE batch_num = $1"
	getLastBatchNumberSQL           = "SELECT MAX(batch_num) FROM batch"
	getTransactionByHashSQL         = "SELECT transaction.encoded FROM transaction WHERE hash = $1"
	getTransactionCountSQL          = "SELECT COUNT(*) FROM transaction WHERE from_address = $1"
	consolidateBatchSQL             = "UPDATE batch SET consolidated_tx_hash = $1 WHERE batch_num = $2"
	getTxsByBatchNumSQL             = "SELECT transaction.encoded FROM transaction WHERE batch_num = $1"
)

// BasicState is a implementation of the state
type BasicState struct {
	db   *pgxpool.Pool
	Tree tree.ReadWriter
}

// NewState creates a new State
func NewState(db *pgxpool.Pool, tree tree.ReadWriter) State {
	return &BasicState{db: db, Tree: tree}
}

// NewBatchProcessor creates a new batch processor
func (s *BasicState) NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) BatchProcessor {
	return &BasicBatchProcessor{State: s}
}

// GetStateRoot returns the root of the state tree
func (s *BasicState) GetStateRoot(ctx context.Context, virtual bool) (*big.Int, error) {
	batch, err := s.GetLastBatch(ctx, virtual)
	if err != nil {
		return nil, err
	}

	root, err := s.Tree.GetRootForBatchNumber(batch.BatchNumber)
	if err != nil {
		return nil, err
	}

	return big.NewInt(0).SetBytes(root), nil
}

// GetBalance from a given address
func (s *BasicState) GetBalance(address common.Address, batchNumber uint64) (*big.Int, error) {
	root, err := s.Tree.GetRootForBatchNumber(batchNumber)
	if err != nil {
		return nil, err
	}

	return s.Tree.GetBalance(address, root)
}

// EstimateGas for a transaction
func (s *BasicState) EstimateGas(transaction *types.Transaction) uint64 {
	// TODO: Calculate once we have txs that interact with SCs
	return 21000 //nolint:gomnd
}

// GetLastBlock gets the latest block
func (s *BasicState) GetLastBlock(ctx context.Context) (*Block, error) {
	var block Block
	err := s.db.QueryRow(ctx, getLastBlockSQL).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// GetPreviousBlock gets the offset previous block respect to latest
func (s *BasicState) GetPreviousBlock(ctx context.Context, offset uint64) (*Block, error) {
	var block Block
	err := s.db.QueryRow(ctx, getPreviousBlockSQL, offset).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// GetBlockByHash gets the block with the required hash
func (s *BasicState) GetBlockByHash(ctx context.Context, hash common.Hash) (*Block, error) {
	var block Block
	err := s.db.QueryRow(ctx, getBlockByHashSQL, hash).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// GetBlockByNumber gets the block with the required number
func (s *BasicState) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error) {
	var block Block
	err := s.db.QueryRow(ctx, getBlockByNumberSQL, blockNumber).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return &block, nil
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
	var row pgx.Row

	if isVirtual {
		row = s.db.QueryRow(ctx, getLastVirtualBatchSQL)
	} else {
		row = s.db.QueryRow(ctx, getLastConsolidatedBatchSQL, common.Hash{})
	}

	var batch Batch
	err := row.Scan(
		&batch.BatchNumber, &batch.BatchHash, &batch.BlockNumber,
		&batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
		&batch.Header, &batch.Uncles, &batch.RawTxsData)

	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetPreviousBatch gets the offset previous batch respect to latest
func (s *BasicState) GetPreviousBatch(ctx context.Context, isVirtual bool, offset uint64) (*Batch, error) {
	var row pgx.Row
	if isVirtual {
		row = s.db.QueryRow(ctx, getPreviousVirtualBatchSQL, offset)
	} else {
		row = s.db.QueryRow(ctx, getPreviousConsolidatedBatchSQL, common.Hash{}, offset)
	}
	var batch Batch
	err := row.Scan(
		&batch.BatchNumber, &batch.BatchHash, &batch.BlockNumber,
		&batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash, &batch.Header,
		&batch.Uncles, &batch.RawTxsData)

	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetBatchByHash gets the batch with the required hash
func (s *BasicState) GetBatchByHash(ctx context.Context, hash common.Hash) (*Batch, error) {
	var batch Batch
	err := s.db.QueryRow(ctx, getBatchByHashSQL, hash).Scan(
		&batch.BatchNumber, &batch.BatchHash, &batch.BlockNumber, &batch.Sequencer, &batch.Aggregator,
		&batch.ConsolidatedTxHash, &batch.Header, &batch.Uncles, &batch.RawTxsData)

	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetBatchByNumber gets the batch with the required number
func (s *BasicState) GetBatchByNumber(ctx context.Context, batchNumber uint64) (*Batch, error) {
	var batch Batch
	err := s.db.QueryRow(ctx, getBatchByNumberSQL, batchNumber).Scan(
		&batch.BatchNumber, &batch.BatchHash, &batch.BlockNumber, &batch.Sequencer, &batch.Aggregator,
		&batch.ConsolidatedTxHash, &batch.Header, &batch.Uncles, &batch.RawTxsData)
	if err != nil {
		return nil, err
	}
	return &batch, nil
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

// GetNonce returns the nonce of the given account at the given batch number
func (s *BasicState) GetNonce(address common.Address, batchNumber uint64) (uint64, error) {
	panic("not implemented yet")
}

// GetTransactionByBatchHashAndIndex gets a transaction from a batch by index
func (s *BasicState) GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error) {
	panic("not implemented")
}

// GetTransactionByBatchNumberAndIndex gets a transaction from a batch by index
func (s *BasicState) GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error) {
	panic("not implemented")
}

// GetTransactionByHash gets a transaction by its hash
func (s *BasicState) GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error) {
	var encoded string
	if err := s.db.QueryRow(ctx, getTransactionByHashSQL, transactionHash).Scan(&encoded); err != nil {
		return nil, err
	}

	b, err := hex.DecodeHex(encoded)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionCount returns the number of transactions sent from an address
func (s *BasicState) GetTransactionCount(ctx context.Context, fromAddress common.Address) (uint64, error) {
	var count uint64
	err := s.db.QueryRow(ctx, getTransactionCountSQL, fromAddress).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
func (s *BasicState) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*types.Receipt, error) {
	panic("not implemented")
}

// Reset resets the state to a block
func (s *BasicState) Reset(blockNumber uint64) error {
	panic("not implemented")
}

// ConsolidateBatch changes the virtual status of a batch
func (s *BasicState) ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash) error {
	if _, err := s.db.Exec(ctx, consolidateBatchSQL, consolidatedTxHash, batchNumber); err != nil {
		return err
	}
	return nil
}

// GetTxsByBatchNum returns all the txs in a given batch
func (s *BasicState) GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error) {
	rows, err := s.db.Query(ctx, getTxsByBatchNumSQL, batchNum)
	if err != nil {
		return nil, err
	}
	txs := make([]*types.Transaction, 0, len(rows.RawValues()))
	var (
		encoded string
		tx      *types.Transaction
		b       []byte
	)
	for rows.Next() {
		if err = rows.Scan(&encoded); err != nil {
			return nil, err
		}

		tx = new(types.Transaction)

		b, err = hex.DecodeHex(encoded)
		if err != nil {
			return nil, err
		}

		if err := tx.UnmarshalBinary(b); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

// AddNewSequencer stores a new sequencer
func (s *BasicState) AddNewSequencer(seq Sequencer) error {
	return nil
}

// SetGenesis populates state with genesis information
func (s *BasicState) SetGenesis(genesis Genesis) error {
	// Genesis Balances
	for address, balance := range genesis.Balances {
		_, _, err := s.Tree.SetBalance(address, balance)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddBlock adds a new block to the State DB
func (s *BasicState) AddBlock(*Block) error {
	// TODO: Implement
	return nil
}

// SetLastBatchNumberSeenOnEthereum sets the last batch number that affected
// the roll-up in order to allow the components to know if the state
// is synchronized or not
func (s *BasicState) SetLastBatchNumberSeenOnEthereum(batchNumber uint64) error {
	return nil
}

// GetLastBatchNumberSeenOnEthereum returns the last batch number stored
// in the state that represents the last batch number that affected the
// roll-up in the Ethereum network.
func (s *BasicState) GetLastBatchNumberSeenOnEthereum() (uint64, error) {
	return 0, nil
}
