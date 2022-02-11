package pgstatestorage

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	getLastBlockSQL                        = "SELECT * FROM state.block ORDER BY block_num DESC LIMIT 1"
	getPreviousBlockSQL                    = "SELECT * FROM state.block ORDER BY block_num DESC LIMIT 1 OFFSET $1"
	getBlockByHashSQL                      = "SELECT * FROM state.block WHERE block_hash = $1"
	getBlockByNumberSQL                    = "SELECT * FROM state.block WHERE block_num = $1"
	getLastBlockNumberSQL                  = "SELECT COALESCE(MAX(block_num), 0) FROM state.block"
	getLastVirtualBatchSQL                 = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root FROM state.batch ORDER BY batch_num DESC LIMIT 1"
	getLastConsolidatedBatchSQL            = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root FROM state.batch WHERE consolidated_tx_hash != $1 ORDER BY batch_num DESC LIMIT 1"
	getPreviousVirtualBatchSQL             = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root FROM state.batch ORDER BY batch_num DESC LIMIT 1 OFFSET $1"
	getPreviousConsolidatedBatchSQL        = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root FROM state.batch WHERE consolidated_tx_hash != $1 ORDER BY batch_num DESC LIMIT 1 OFFSET $2"
	getBatchByHashSQL                      = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root FROM state.batch WHERE batch_hash = $1"
	getBatchByNumberSQL                    = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root FROM state.batch WHERE batch_num = $1"
	getLastVirtualBatchNumberSQL           = "SELECT COALESCE(MAX(batch_num), 0) FROM state.batch"
	getLastConsolidatedBatchNumberSQL      = "SELECT COALESCE(MAX(batch_num), 0) FROM state.batch WHERE consolidated_tx_hash != $1"
	getTransactionByHashSQL                = "SELECT transaction.encoded FROM state.transaction WHERE hash = $1"
	getTransactionByBatchHashAndIndexSQL   = "SELECT transaction.encoded FROM state.transaction inner join state.batch on (state.transaction.batch_num = state.batch.batch_num) WHERE state.batch.batch_hash = $1 and state.transaction.tx_index = $2"
	getTransactionByBatchNumberAndIndexSQL = "SELECT transaction.encoded FROM state.transaction WHERE batch_num = $1 AND tx_index = $2"
	getTransactionCountSQL                 = "SELECT COUNT(*) FROM state.transaction WHERE from_address = $1"
	consolidateBatchSQL                    = "UPDATE state.batch SET consolidated_tx_hash = $1, consolidated_at = $3, aggregator = $4 WHERE batch_num = $2"
	getTxsByBatchNumSQL                    = "SELECT transaction.encoded FROM state.transaction WHERE batch_num = $1"
	addBlockSQL                            = "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)"
	addSequencerSQL                        = "INSERT INTO state.sequencer (address, url, chain_id, block_num) VALUES ($1, $2, $3, $4) ON CONFLICT (chain_id) DO UPDATE SET address = EXCLUDED.address, url = EXCLUDED.url, block_num = EXCLUDED.block_num"
	updateLastBatchSeenSQL                 = "UPDATE state.misc SET last_batch_num_seen = $1"
	getLastBatchSeenSQL                    = "SELECT last_batch_num_seen FROM state.misc LIMIT 1"
	updateLastBatchConsolidatedSQL         = "UPDATE state.misc SET last_batch_num_consolidated = $1"
	getLastBatchConsolidatedSQL            = "SELECT last_batch_num_consolidated FROM state.misc LIMIT 1"
	getSequencerSQL                        = "SELECT * FROM state.sequencer WHERE address = $1"
	getReceiptSQL                          = "SELECT * FROM state.receipt WHERE tx_hash = $1"
	resetSQL                               = "DELETE FROM state.block WHERE block_num > $1"
	addBatchSQL                            = "INSERT INTO state.batch (batch_num, batch_hash, block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, chain_id, global_exit_root) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)"
	addTransactionSQL                      = "INSERT INTO state.transaction (hash, from_address, encoded, decoded, batch_num, tx_index) VALUES($1, $2, $3, $4, $5, $6)"
	addReceiptSQL                          = "INSERT INTO state.receipt (type, post_state, status, cumulative_gas_used, gas_used, batch_num, batch_hash, tx_hash, tx_index, tx_from, tx_to)	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
)

var (
	ten = big.NewInt(encoding.Base10)
)

// PostgresStorage implements the Storage interface
type PostgresStorage struct {
	db *pgxpool.Pool
}

// NewPostgresStorage creates a new StateDB
func NewPostgresStorage(db *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{db: db}
}

// GetLastBlock gets the latest block
func (s *PostgresStorage) GetLastBlock(ctx context.Context) (*state.Block, error) {
	var block state.Block
	err := s.db.QueryRow(ctx, getLastBlockSQL).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetPreviousBlock gets the offset previous block respect to latest
func (s *PostgresStorage) GetPreviousBlock(ctx context.Context, offset uint64) (*state.Block, error) {
	var block state.Block
	err := s.db.QueryRow(ctx, getPreviousBlockSQL, offset).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetBlockByHash gets the block with the required hash
func (s *PostgresStorage) GetBlockByHash(ctx context.Context, hash common.Hash) (*state.Block, error) {
	var block state.Block
	err := s.db.QueryRow(ctx, getBlockByHashSQL, hash).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetBlockByNumber gets the block with the required number
func (s *PostgresStorage) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*state.Block, error) {
	var block state.Block
	err := s.db.QueryRow(ctx, getBlockByNumberSQL, blockNumber).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetLastBlockNumber gets the latest block number
func (s *PostgresStorage) GetLastBlockNumber(ctx context.Context) (uint64, error) {
	var lastBlockNum uint64
	err := s.db.QueryRow(ctx, getLastBlockNumberSQL).Scan(&lastBlockNum)

	if reflect.TypeOf(err) == reflect.TypeOf(pgx.ScanArgError{}) {
		return 0, state.ErrStateNotSynchronized
	} else if err != nil {
		return 0, err
	}

	return lastBlockNum, nil
}

// GetLastBatch gets the latest batch
func (s *PostgresStorage) GetLastBatch(ctx context.Context, isVirtual bool) (*state.Batch, error) {
	var (
		batch           state.Batch
		maticCollateral pgtype.Numeric
	)
	var err error

	if isVirtual {
		err = s.db.QueryRow(ctx, getLastVirtualBatchSQL).Scan(&batch.BlockNumber,
			&batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
			&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
			&batch.ReceivedAt, &batch.ConsolidatedAt)
	} else {
		err = s.db.QueryRow(ctx, getLastConsolidatedBatchSQL, common.Hash{}).Scan(
			&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
			&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
			&batch.ReceivedAt, &batch.ConsolidatedAt)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	batch.MaticCollateral = new(big.Int).Mul(maticCollateral.Int, big.NewInt(0).Exp(ten, big.NewInt(int64(maticCollateral.Exp)), nil))

	return &batch, nil
}

// GetPreviousBatch gets the offset previous batch respect to latest
func (s *PostgresStorage) GetPreviousBatch(ctx context.Context, isVirtual bool, offset uint64) (*state.Batch, error) {
	var (
		batch           state.Batch
		maticCollateral pgtype.Numeric
	)
	var err error

	if isVirtual {
		err = s.db.QueryRow(ctx, getPreviousVirtualBatchSQL, offset).Scan(
			&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
			&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
			&batch.ReceivedAt, &batch.ConsolidatedAt)
	} else {
		err = s.db.QueryRow(ctx, getPreviousConsolidatedBatchSQL, common.Hash{}, offset).Scan(
			&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash, &batch.Header,
			&batch.Uncles, &batch.RawTxsData, &maticCollateral,
			&batch.ReceivedAt, &batch.ConsolidatedAt)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	batch.MaticCollateral = new(big.Int).Mul(maticCollateral.Int, big.NewInt(0).Exp(ten, big.NewInt(int64(maticCollateral.Exp)), nil))
	return &batch, nil
}

// GetBatchByHash gets the batch with the required hash
func (s *PostgresStorage) GetBatchByHash(ctx context.Context, hash common.Hash) (*state.Batch, error) {
	var (
		batch           state.Batch
		maticCollateral pgtype.Numeric
	)
	err := s.db.QueryRow(ctx, getBatchByHashSQL, hash).Scan(
		&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
		&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
		&batch.ReceivedAt, &batch.ConsolidatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	batch.MaticCollateral = new(big.Int).Mul(maticCollateral.Int, big.NewInt(0).Exp(ten, big.NewInt(int64(maticCollateral.Exp)), nil))
	batch.Transactions, err = s.getBatchTransactions(ctx, batch)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return &batch, nil
}

// GetBatchByNumber gets the batch with the required number
func (s *PostgresStorage) GetBatchByNumber(ctx context.Context, batchNumber uint64) (*state.Batch, error) {
	var (
		batch           state.Batch
		maticCollateral pgtype.Numeric
	)
	err := s.db.QueryRow(ctx, getBatchByNumberSQL, batchNumber).Scan(
		&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
		&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
		&batch.ReceivedAt, &batch.ConsolidatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	batch.MaticCollateral = new(big.Int).Mul(maticCollateral.Int, big.NewInt(0).Exp(ten, big.NewInt(int64(maticCollateral.Exp)), nil))
	batch.Transactions, err = s.getBatchTransactions(ctx, batch)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return &batch, nil
}

// GetLastBatchNumber gets the latest batch number
func (s *PostgresStorage) GetLastBatchNumber(ctx context.Context) (uint64, error) {
	var lastBatchNumber uint64
	err := s.db.QueryRow(ctx, getLastVirtualBatchNumberSQL).Scan(&lastBatchNumber)

	if err != nil {
		return 0, err
	}

	return lastBatchNumber, nil
}

// GetLastConsolidatedBatchNumber gets the latest consolidated batch number
func (s *PostgresStorage) GetLastConsolidatedBatchNumber(ctx context.Context) (uint64, error) {
	var lastBatchNumber uint64
	err := s.db.QueryRow(ctx, getLastConsolidatedBatchNumberSQL, common.Hash{}).Scan(&lastBatchNumber)

	if err != nil {
		return 0, err
	}

	return lastBatchNumber, nil
}

// GetTransactionByBatchHashAndIndex gets a transaction from a batch by index
func (s *PostgresStorage) GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error) {
	var encoded string
	err := s.db.QueryRow(ctx, getTransactionByBatchHashAndIndexSQL, batchHash.Bytes(), index).Scan(&encoded)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
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

// GetTransactionByBatchNumberAndIndex gets a transaction from a batch by index
func (s *PostgresStorage) GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error) {
	var encoded string
	err := s.db.QueryRow(ctx, getTransactionByBatchNumberAndIndexSQL, batchNumber, index).Scan(&encoded)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
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

// GetTransactionByHash gets a transaction by its hash
func (s *PostgresStorage) GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error) {
	var encoded string
	err := s.db.QueryRow(ctx, getTransactionByHashSQL, transactionHash).Scan(&encoded)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
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
func (s *PostgresStorage) GetTransactionCount(ctx context.Context, fromAddress common.Address) (uint64, error) {
	var count uint64
	err := s.db.QueryRow(ctx, getTransactionCountSQL, fromAddress).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
func (s *PostgresStorage) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*state.Receipt, error) {
	var receipt state.Receipt
	var batchNumber uint64
	err := s.db.QueryRow(ctx, getReceiptSQL, transactionHash).Scan(&receipt.Type, &receipt.PostState, &receipt.Status,
		&receipt.CumulativeGasUsed, &receipt.GasUsed, &batchNumber, &receipt.BlockHash, &receipt.TxHash, &receipt.TransactionIndex, &receipt.From, &receipt.To)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	receipt.BlockNumber = new(big.Int).SetUint64(batchNumber)
	return &receipt, nil
}

// Reset resets the state to a block
func (s *PostgresStorage) Reset(ctx context.Context, blockNumber uint64) error {
	if _, err := s.db.Exec(ctx, resetSQL, blockNumber); err != nil {
		return err
	}
	return nil
}

// ConsolidateBatch changes the virtual status of a batch
func (s *PostgresStorage) ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash, consolidatedAt time.Time, aggregator common.Address) error {
	if _, err := s.db.Exec(ctx, consolidateBatchSQL, consolidatedTxHash, batchNumber, consolidatedAt, aggregator); err != nil {
		return err
	}
	return nil
}

// GetTxsByBatchNum returns all the txs in a given batch
func (s *PostgresStorage) GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error) {
	rows, err := s.db.Query(ctx, getTxsByBatchNumSQL, batchNum)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
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

// AddSequencer stores a new sequencer
func (s *PostgresStorage) AddSequencer(ctx context.Context, seq state.Sequencer) error {
	_, err := s.db.Exec(ctx, addSequencerSQL, seq.Address, seq.URL, seq.ChainID.Uint64(), seq.BlockNumber)
	return err
}

// GetSequencer gets a sequencer
func (s *PostgresStorage) GetSequencer(ctx context.Context, address common.Address) (*state.Sequencer, error) {
	var seq state.Sequencer
	var cID uint64
	err := s.db.QueryRow(ctx, getSequencerSQL, address.Bytes()).Scan(&seq.Address, &seq.URL, &cID, &seq.BlockNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	seq.ChainID = big.NewInt(0).SetUint64(cID)

	return &seq, nil
}

// AddBlock adds a new block to the State Store
func (s *PostgresStorage) AddBlock(ctx context.Context, block *state.Block) error {
	_, err := s.db.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.Bytes(), block.ParentHash.Bytes(), block.ReceivedAt)
	return err
}

// SetLastBatchNumberSeenOnEthereum sets the last batch number that affected
// the roll-up in order to allow the components to know if the state
// is synchronized or not
func (s *PostgresStorage) SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64) error {
	_, err := s.db.Exec(ctx, updateLastBatchSeenSQL, batchNumber)
	return err
}

// GetLastBatchNumberSeenOnEthereum returns the last batch number stored
// in the state that represents the last batch number that affected the
// roll-up in the Ethereum network.
func (s *PostgresStorage) GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error) {
	var batchNumber uint64
	err := s.db.QueryRow(ctx, getLastBatchSeenSQL).Scan(&batchNumber)

	if err != nil {
		return 0, err
	}

	return batchNumber, nil
}

// SetLastBatchNumberConsolidatedOnEthereum sets the last batch number that was consolidated on ethereum
func (s *PostgresStorage) SetLastBatchNumberConsolidatedOnEthereum(ctx context.Context, batchNumber uint64) error {
	_, err := s.db.Exec(ctx, updateLastBatchConsolidatedSQL, batchNumber)
	return err
}

// GetLastBatchNumberConsolidatedOnEthereum sets the last batch number that was consolidated on ethereum
func (s *PostgresStorage) GetLastBatchNumberConsolidatedOnEthereum(ctx context.Context) (uint64, error) {
	var batchNumber uint64
	err := s.db.QueryRow(ctx, getLastBatchConsolidatedSQL).Scan(&batchNumber)

	if err != nil {
		return 0, err
	}

	return batchNumber, nil
}

// AddBatch adds a new batch to the State Store
func (s *PostgresStorage) AddBatch(ctx context.Context, batch *state.Batch) error {
	_, err := s.db.Exec(ctx, addBatchSQL, batch.Number().Uint64(), batch.Hash(), batch.BlockNumber, batch.Sequencer, batch.Aggregator,
		batch.ConsolidatedTxHash, batch.Header, batch.Uncles, batch.RawTxsData, batch.MaticCollateral.String(), batch.ReceivedAt, batch.ChainID.String(), batch.GlobalExitRoot)
	return err
}

// AddTransaction adds a new transqaction to the State Store
func (s *PostgresStorage) AddTransaction(ctx context.Context, tx *types.Transaction, batchNumber uint64, index uint) error {
	binary, err := tx.MarshalBinary()
	if err != nil {
		panic(err)
	}
	encoded := hex.EncodeToHex(binary)

	binary, err = tx.MarshalJSON()
	if err != nil {
		panic(err)
	}
	decoded := string(binary)

	_, err = s.db.Exec(ctx, addTransactionSQL, tx.Hash().Bytes(), "", encoded, decoded, batchNumber, index)
	return err
}

// AddReceipt adds a new receipt to the State Store
func (s *PostgresStorage) AddReceipt(ctx context.Context, receipt *state.Receipt) error {
	_, err := s.db.Exec(ctx, addReceiptSQL, receipt.Type, receipt.PostState, receipt.Status, receipt.CumulativeGasUsed, receipt.GasUsed, receipt.BlockNumber.Uint64(), receipt.BlockHash.Bytes(), receipt.TxHash.Bytes(), receipt.TransactionIndex, receipt.From.Bytes(), receipt.To.Bytes())
	return err
}

func (s *PostgresStorage) getBatchTransactions(ctx context.Context, batch state.Batch) ([]*types.Transaction, error) {
	transactions, err := s.GetTxsByBatchNum(ctx, batch.Number().Uint64())
	if errors.Is(err, pgx.ErrNoRows) {
		transactions = []*types.Transaction{}
	} else if err != nil {
		return nil, err
	}

	return transactions, nil
}
