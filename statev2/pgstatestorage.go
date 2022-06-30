package statev2

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	addGlobalExitRootSQL                   = "INSERT INTO statev2.exit_root (block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root) VALUES ($1, $2, $3, $4, $5)"
	getLatestExitRootSQL                   = "SELECT block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root FROM statev2.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
	getLatestExitRootBlockNumSQL           = "SELECT block_num FROM statev2.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
	addVirtualBatchSQL                     = "INSERT INTO statev2.virtual_batch (batch_num, tx_hash, sequencer, block_num) VALUES ($1, $2, $3, $4)"
	addForcedBatchSQL                      = "INSERT INTO statev2.forced_batch (forced_batch_num, global_exit_root, timestamp, raw_txs_data, sequencer, batch_num, block_num) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	getForcedBatchSQL                      = "SELECT forced_batch_num, global_exit_root, timestamp, raw_txs_data, sequencer, batch_num, block_num FROM statev2.forced_batch WHERE forced_batch_num = $1"
	addBlockSQL                            = "INSERT INTO statev2.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)"
	getLastBlockSQL                        = "SELECT block_num, block_hash, parent_hash, received_at FROM statev2.block ORDER BY block_num DESC LIMIT 1"
	getPreviousBlockSQL                    = "SELECT block_num, block_hash, parent_hash, received_at FROM statev2.block ORDER BY block_num DESC LIMIT 1 OFFSET $1"
	resetSQL                               = "DELETE FROM statev2.block WHERE block_num > $1"
	resetTrustedStateSQL                   = "DELETE FROM statev2.batch WHERE batch_num > $1"
	addVerifiedBatchSQL                    = "INSERT INTO statev2.verified_batch (block_num, batch_num, tx_hash, aggregator) VALUES ($1, $2, $3, $4)"
	getVerifiedBatchSQL                    = "SELECT block_num, batch_num, tx_hash, aggregator FROM statev2.verified_batch WHERE batch_num = $1"
	getLastBatchNumberSQL                  = "SELECT COALESCE(MAX(batch_num), 0) FROM statev2.batch"
	getLastNBatchesSQL                     = "SELECT batch_num, global_exit_root, timestamp from statev2.batch ORDER BY batch_num DESC LIMIT $1"
	getLastBatchTimeSQL                    = "SELECT timestamp FROM statev2.batch ORDER BY batch_num DESC LIMIT 1"
	getLastVirtualBatchNumSQL              = "SELECT batch_num FROM statev2.virtual_batch ORDER BY batch_num DESC LIMIT 1"
	getLastVirtualBatchBlockNumSQL         = "SELECT block_num FROM statev2.virtual_batch ORDER BY batch_num DESC LIMIT 1"
	getLastBlockNumSQL                     = "SELECT block_num FROM statev2.block ORDER BY block_num DESC LIMIT 1"
	getBlockTimeByNumSQL                   = "SELECT received_at FROM statev2.block WHERE block_num = $1"
	getBatchByNumberSQL                    = "SELECT batch_num, global_exit_root, timestamp from statev2.batch WHERE batch_num = $1"
	getEncodedTransactionsByBatchNumberSQL = "SELECT encoded from statev2.transaction WHERE batch_num = $1"
	getLastBatchSeenSQL                    = "SELECT last_batch_num_seen FROM statev2.sync_info LIMIT 1"
	updateLastBatchSeenSQL                 = "UPDATE statev2.sync_info SET last_batch_num_seen = $1"
	resetTrustedBatchSQL                   = "DELETE FROM statev2.batch WHERE batch_num > $1"
	storeBatchHeaderSQL                    = "INSERT INTO statev2.batch (batch_num, global_exit_root, timestamp, sequencer, raw_txs_data) VALUES ($1, $2, $3, $4, $5)"
	getNextForcedBatchesSQL                = "SELECT forced_batch_num, global_exit_root, timestamp, raw_txs_data, sequencer, batch_num, block_num FROM statev2.forced_batch WHERE batch_num IS NULL LIMIT $1"
	addBatchNumberInForcedBatchSQL         = "UPDATE statev2.forced_batch SET batch_num = $2 WHERE forced_batch_num = $1"
	getL2BlockByNumberSQL                  = "SELECT l2_block_num, encoded, header, uncles, received_at from statev2.transaction WHERE batch_num = $1"
	getTransactionByHashSQL                = "SELECT transaction.encoded FROM statev2.transaction WHERE hash = $1"
	getReceiptSQL                          = "SELECT r.tx_hash, r.type, r.post_state, r.status, r.cumulative_gas_used, r.gas_used, r.contract_address, t.encoded, t.l2_block_num, b.block_hash FROM statev2.receipt r INNER JOIN statev2.transaction t ON t.hash = r.tx_hash INNER JOIN statev2.l2block b ON b.block_num = t.l2_block_num WHERE r.tx_hash = $1"
	getTransactionByBlockHashAndIndexSQL   = "SELECT t.encoded FROM statev2.transaction t INNER JOIN statev2.l2block b ON t.l2_block_num = b.batch_num WHERE b.block_hash = $1 AND 0 = $2"
	getTransactionByBlockNumberAndIndexSQL = "SELECT t.encoded FROM statev2.transaction t WHERE t.l2_block_num = $1 AND 0 = $2"
	getBlockTransactionCountByHashSQL      = "SELECT COUNT(*) FROM statev2.transaction t INNER JOIN statev2.l2block b ON b.block_num = t.l2_block_num WHERE b.block_hash = $1"
	getBlockTransactionCountByNumberSQL    = "SELECT COUNT(*) FROM state.transaction t WHERE t.l2_block_num = $1"
	getTransactionLogsSQL                  = "SELECT t.l2_block_num, b.block_hash, l.tx_hash, l.log_index, l.address, l.data, l.topic0, l.topic1, l.topic2, l.topic3 FROM state.log l INNER JOIN statev2.transaction t ON t.hash = l.tx_hash INNER JOIN statev2.l2block b ON b.block_num = t.l2_block_num WHERE transaction_hash = $1"
	addL2BlockSQL                          = "INSERT INTO statev2.l2block (block_num, block_hash, header, uncles, parent_hash, state_root, received_at, batch_num) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
)

// PostgresStorage implements the Storage interface
type PostgresStorage struct {
	*pgxpool.Pool
}

// NewPostgresStorage creates a new StateDB
func NewPostgresStorage(db *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{
		db,
	}
}

// getExecQuerier determines which execQuerier to use, dbTx or the main pgxpool
func (p *PostgresStorage) getExecQuerier(dbTx pgx.Tx) execQuerier {
	if dbTx != nil {
		return dbTx
	}
	return p
}

// Reset resets the state to a block
func (p *PostgresStorage) Reset(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	if _, err := e.Exec(ctx, resetSQL, blockNumber); err != nil {
		return err
	}
	// TODO: Remove consolidations
	return nil
}

func (p *PostgresStorage) ResetTrustedState(ctx context.Context, batchNum uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	if _, err := e.Exec(ctx, resetTrustedStateSQL, batchNum); err != nil {
		return err
	}
	// TODO: Remove consolidations
	return nil
}

// AddBlock adds a new block to the State Store
func (p *PostgresStorage) AddBlock(ctx context.Context, block *Block, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.String(), block.ParentHash.String(), block.ReceivedAt)
	return err
}

// GetLastBlock returns the last L1 block.
func (p *PostgresStorage) GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*Block, error) {
	var (
		blockHash  string
		parentHash string
		block      Block
	)
	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getLastBlockSQL).Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return &block, err
}

// GetPreviousBlock gets the offset previous L1 block respect to latest.
func (p *PostgresStorage) GetPreviousBlock(ctx context.Context, offset uint64, dbTx pgx.Tx) (*Block, error) {
	var (
		blockHash  string
		parentHash string
		block      Block
	)
	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getPreviousBlockSQL, offset).Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return &block, err
}

// AddGlobalExitRoot adds a new ExitRoot to the db
func (p *PostgresStorage) AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addGlobalExitRootSQL, exitRoot.BlockNumber, exitRoot.GlobalExitRootNum.String(), exitRoot.MainnetExitRoot, exitRoot.RollupExitRoot, exitRoot.GlobalExitRoot)
	return err
}

// GetLatestExitRoot get the latest ExitRoot synced.
func (p *PostgresStorage) GetLatestGlobalExitRoot(ctx context.Context, dbTx pgx.Tx) (*GlobalExitRoot, error) {
	var (
		exitRoot  GlobalExitRoot
		globalNum uint64
		err       error
	)

	e := p.getExecQuerier(dbTx)
	err = e.QueryRow(ctx, getLatestExitRootSQL).Scan(&exitRoot.BlockNumber, &globalNum, &exitRoot.MainnetExitRoot, &exitRoot.RollupExitRoot, &exitRoot.GlobalExitRoot)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	exitRoot.GlobalExitRootNum = new(big.Int).SetUint64(globalNum)
	return &exitRoot, nil
}

// GetNumberOfBlocksSinceLastGERUpdate gets number of blocks since last global exit root update
func (p *PostgresStorage) GetNumberOfBlocksSinceLastGERUpdate(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var (
		lastBlockNum         uint64
		lastExitRootBlockNum uint64
		err                  error
	)
	e := p.getExecQuerier(dbTx)
	err = e.QueryRow(ctx, getLastBlockNumSQL).Scan(&lastBlockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	err = p.QueryRow(ctx, getLatestExitRootBlockNumSQL).Scan(&lastExitRootBlockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return lastBlockNum - lastExitRootBlockNum, nil
}

func (p *PostgresStorage) GetTimeForLatestBatchVirtualization(ctx context.Context, dbTx pgx.Tx) (time.Time, error) {
	var (
		blockNum  uint64
		timestamp time.Time
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastVirtualBatchBlockNumSQL).Scan(&blockNum)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrNotFound
	} else if err != nil {
		return time.Time{}, err
	}

	err = p.QueryRow(ctx, getBlockTimeByNumSQL, blockNum).Scan(&timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrNotFound
	} else if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

// AddForcedBatch adds a new ForcedBatch to the db
func (p *PostgresStorage) AddForcedBatch(ctx context.Context, forcedBatch *ForcedBatch, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, addForcedBatchSQL, forcedBatch.ForcedBatchNumber, forcedBatch.GlobalExitRoot.String(), forcedBatch.ForcedAt, forcedBatch.RawTxsData, forcedBatch.Sequencer.String(), forcedBatch.BatchNumber, forcedBatch.BlockNumber)
	return err
}

// GetForcedBatch get an L1 forcedBatch.
func (p *PostgresStorage) GetForcedBatch(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (*ForcedBatch, error) {
	var (
		forcedBatch    ForcedBatch
		globalExitRoot string
		rawTxs         string
		seq            string
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getForcedBatchSQL, forcedBatchNumber).Scan(&forcedBatch.ForcedBatchNumber, &globalExitRoot, &forcedBatch.ForcedAt, &rawTxs, &seq, &forcedBatch.BatchNumber, &forcedBatch.BlockNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	forcedBatch.RawTxsData, err = hex.DecodeString(rawTxs)
	if err != nil {
		return nil, err
	}
	forcedBatch.Sequencer = common.HexToAddress(seq)
	forcedBatch.GlobalExitRoot = common.HexToHash(globalExitRoot)
	return &forcedBatch, nil
}

// AddVerifiedBatch adds a new VerifiedBatch to the db
func (p *PostgresStorage) AddVerifiedBatch(ctx context.Context, verifiedBatch *VerifiedBatch, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addVerifiedBatchSQL, verifiedBatch.BlockNumber, verifiedBatch.BatchNumber, verifiedBatch.TxHash.String(), verifiedBatch.Aggregator.String())
	return err
}

// GetVerifiedBatch get an L1 verifiedBatch.
func (p *PostgresStorage) GetVerifiedBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*VerifiedBatch, error) {
	var (
		verifiedBatch VerifiedBatch
		txHash        string
		agg           string
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getVerifiedBatchSQL, batchNumber).Scan(&verifiedBatch.BlockNumber, &verifiedBatch.BatchNumber, &txHash, &agg)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	verifiedBatch.Aggregator = common.HexToAddress(agg)
	verifiedBatch.TxHash = common.HexToHash(txHash)
	return &verifiedBatch, nil
}

func (p *PostgresStorage) GetLastNBatches(ctx context.Context, numBatches uint, dbTx pgx.Tx) ([]*Batch, error) {
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getLastNBatchesSQL, numBatches)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	batches := make([]*Batch, 0, len(rows.RawValues()))

	for rows.Next() {
		var (
			batch  Batch
			gerStr string
		)
		err := rows.Scan(&batch.BatchNumber, &gerStr, &batch.Timestamp)
		if err != nil {
			return nil, err
		}

		batch.GlobalExitRootNum = new(big.Int).SetBytes(common.FromHex(gerStr))
		batches = append(batches, &batch)
	}

	return batches, nil
}

// GetLastBatchNumber get last trusted batch number
func (p *PostgresStorage) GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var batchNumber uint64
	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getLastBatchNumberSQL).Scan(&batchNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrStateNotSynchronized
	}
	return batchNumber, err
}

// GetLastBatchTime gets last trusted batch time
func (p *PostgresStorage) GetLastBatchTime(ctx context.Context, dbTx pgx.Tx) (time.Time, error) {
	var timestamp time.Time
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastBatchTimeSQL).Scan(&timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrStateNotSynchronized
	} else if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

// GetLastVirtualBatchNum gets last virtual batch num
func (p *PostgresStorage) GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var batchNum uint64
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastVirtualBatchNumSQL).Scan(&batchNum)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}
	return batchNum, nil
}

// SetLastBatchNumberSeenOnEthereum sets the last batch number that affected
// the roll-up in order to allow the components to know if the state
// is synchronized or not
func (p *PostgresStorage) SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, updateLastBatchSeenSQL, batchNumber)
	return err
}

// GetLastBatchNumberSeenOnEthereum returns the last batch number stored
// in the state that represents the last batch number that affected the
// roll-up in the Ethereum network.
func (p *PostgresStorage) GetLastBatchNumberSeenOnEthereum(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var batchNumber uint64
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastBatchSeenSQL).Scan(&batchNumber)

	if err != nil {
		return 0, err
	}

	return batchNumber, nil
}

func (p *PostgresStorage) GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*Batch, error) {
	var (
		batch  Batch
		gerStr string
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getBatchByNumberSQL, batchNumber).Scan(&batch.BatchNumber, &gerStr, &batch.Timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	batch.GlobalExitRootNum = new(big.Int).SetBytes(common.FromHex(gerStr))
	return &batch, nil
}

func (p *PostgresStorage) GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (encoded []string, err error) {
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getEncodedTransactionsByBatchNumberSQL, batchNumber)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]string, 0, len(rows.RawValues()))

	for rows.Next() {
		var encoded string
		err := rows.Scan(&encoded)
		if err != nil {
			return nil, err
		}

		txs = append(txs, encoded)
	}
	return txs, nil
}

// ResetTrustedState resets the batches which the batch number is highter than the input.
func (p *PostgresStorage) ResetTrustedBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, resetTrustedBatchSQL, batchNumber)
	return err
}

// AddVirtualBatch adds a new virtual batch to the storage.
func (p *PostgresStorage) AddVirtualBatch(ctx context.Context, virtualBatch *VirtualBatch, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addVirtualBatchSQL, virtualBatch.BatchNumber, virtualBatch.TxHash.String(), virtualBatch.Sequencer.String(), virtualBatch.BlockNumber)
	return err
}

// StoreBatchHeader adds a new trusted batch header to the storage.
func (p *PostgresStorage) StoreBatchHeader(ctx context.Context, batch Batch, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, storeBatchHeaderSQL, batch.BatchNumber, batch.GlobalExitRoot.String(), batch.Timestamp, batch.Coinbase.String(), batch.BatchL2Data)
	return err
}

// GetNextForcedBatches gets the next forced batches from the queue.
func (p *PostgresStorage) GetNextForcedBatches(ctx context.Context, nextForcedBatches int, dbTx pgx.Tx) ([]ForcedBatch, error) {
	q := p.getExecQuerier(dbTx)
	// Get the next forced batches
	rows, err := q.Query(ctx, getNextForcedBatchesSQL, nextForcedBatches)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	batches := make([]ForcedBatch, 0, len(rows.RawValues()))

	for rows.Next() {
		var (
			forcedBatch    ForcedBatch
			globalExitRoot string
			rawTxs         string
			seq            string
		)
		err := rows.Scan(&forcedBatch.ForcedBatchNumber, &globalExitRoot, &forcedBatch.ForcedAt, &rawTxs, &seq, &forcedBatch.BatchNumber, &forcedBatch.BlockNumber)
		if err != nil {
			return nil, err
		}
		forcedBatch.RawTxsData, err = hex.DecodeString(rawTxs)
		if err != nil {
			return nil, err
		}
		forcedBatch.Sequencer = common.HexToAddress(seq)
		forcedBatch.GlobalExitRoot = common.HexToHash(globalExitRoot)
		batches = append(batches, forcedBatch)
	}

	return batches, nil
}

// AddBatchNumberInForcedBatch updates the forced_batch table with the batchNumber.
func (p *PostgresStorage) AddBatchNumberInForcedBatch(ctx context.Context, forceBatchNumber, batchNumber uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addBatchNumberInForcedBatchSQL, forceBatchNumber, batchNumber)
	return err
}

func (p *PostgresStorage) GetL2BlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*L2Block, error) {
	var block L2Block
	var encoded string
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getL2BlockByNumberSQL, blockNumber).Scan(&block.BlockNumber, &encoded, &block.Header, &block.Uncles, &block.ReceivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	block.Transactions = make([]*types.Transaction, 1)

	tx, err := decodeTx(encoded)
	if err != nil {
		return nil, err
	}

	block.Transactions[0] = tx

	return &block, nil
}

// GetTransactionByHash gets a transaction accordingly to the provided transaction hash
func (p *PostgresStorage) GetTransactionByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error) {
	var encoded string
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByHashSQL, transactionHash).Scan(&encoded)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := decodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionReceipt gets a transaction receipt accordingly to the provided transaction hash
func (p *PostgresStorage) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Receipt, error) {
	var encodedTx string
	var l2BlockNum uint64
	var l2BlockHash string

	receipt := types.Receipt{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getReceiptSQL, transactionHash).
		Scan(&receipt.TxHash,
			&receipt.Type,
			&receipt.PostState,
			&receipt.Status,
			&receipt.CumulativeGasUsed,
			&receipt.GasUsed,
			&receipt.ContractAddress,
			&encodedTx,
			&l2BlockNum,
			&l2BlockHash,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	logs, err := p.getTransactionLogs(ctx, transactionHash, dbTx)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}

	receipt.BlockNumber = big.NewInt(0).SetUint64(l2BlockNum)
	receipt.BlockHash = common.HexToHash(l2BlockHash)
	receipt.TransactionIndex = 0

	receipt.Logs = logs
	receipt.Bloom = types.CreateBloom(types.Receipts{&receipt})

	return &receipt, nil
}

// GetTransactionByBlockHashAndIndex gets a transaction accordingly to the block hash and transaction index provided.
// since we only have a single transaction per l2 block, any index different from 0 will return a not found result
func (p *PostgresStorage) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index uint64, dbTx pgx.Tx) (*types.Transaction, error) {
	var encoded string
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByBlockHashAndIndexSQL, blockHash.Hex(), index).Scan(&encoded)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := decodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionByBlockNumberAndIndex gets a transaction accordingly to the block number and transaction index provided.
// since we only have a single transaction per l2 block, any index different from 0 will return a not found result
func (p *PostgresStorage) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNumber uint64, index uint64, dbTx pgx.Tx) (*types.Transaction, error) {
	var encoded string
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByBlockNumberAndIndexSQL, blockNumber, index).Scan(&encoded)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := decodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetBlockTransactionCountByHash returns the number of transactions related to the provided block hash
func (p *PostgresStorage) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash, dbTx pgx.Tx) (uint64, error) {
	var count uint64
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getBlockTransactionCountByHashSQL, blockHash.Hex()).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetBlockTransactionCountByNumber returns the number of transactions related to the provided block number
func (p *PostgresStorage) GetBlockTransactionCountByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	var count uint64
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getBlockTransactionCountByNumberSQL, blockNumber).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// getTransactionLogs returns the logs of a transaction by transaction hash
func (p *PostgresStorage) getTransactionLogs(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) ([]*types.Log, error) {
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, getTransactionLogsSQL, transactionHash)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]*types.Log, 0, len(rows.RawValues()))

	for rows.Next() {
		var log types.Log
		var txHash, logAddress, logData, topic0 string
		var topic1, topic2, topic3 *string

		err := rows.Scan(
			&log.BlockNumber,
			&log.BlockHash,
			&txHash,
			&log.Index,
			&logAddress,
			&logData,
			&topic0,
			&topic1,
			&topic2,
			&topic3)
		if err != nil {
			return nil, err
		}

		log.TxHash = common.HexToHash(txHash)
		log.Address = common.HexToAddress(logAddress)
		log.TxIndex = uint(0)
		log.Data = []byte(logData)

		log.Topics = []common.Hash{common.HexToHash(topic0)}
		if topic1 != nil {
			log.Topics = append(log.Topics, common.HexToHash(*topic1))
		}

		if topic2 != nil {
			log.Topics = append(log.Topics, common.HexToHash(*topic2))
		}

		if topic3 != nil {
			log.Topics = append(log.Topics, common.HexToHash(*topic3))
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

// decodeTx decodes a string rlp tx representation into a types.Transaction instance
func decodeTx(encodedTx string) (*types.Transaction, error) {
	b, err := hex.DecodeHex(encodedTx)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return tx, nil
}

// AddL2Block adds a new L2 block to the State Store
func (p *PostgresStorage) AddL2Block(ctx context.Context, batchNumber uint64, l2Block L2Block, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)

	var header *string
	if l2Block.Header != nil {
		headerBytes, err := json.Marshal(l2Block.Header)
		if err != nil {
			return err
		}
		*header = string(headerBytes)
	}

	var uncles *string
	if l2Block.Uncles != nil {
		unclesBytes, err := json.Marshal(l2Block.Uncles)
		if err != nil {
			return err
		}
		*uncles = string(unclesBytes)
	}

	_, err := e.Exec(ctx, addL2BlockSQL,
		l2Block.BlockNumber, l2Block.Hash().String(), header, uncles,
		l2Block.Header.ParentHash.String(), l2Block.Header.Root.String(),
		l2Block.ReceivedAt, batchNumber)
	return err
}
