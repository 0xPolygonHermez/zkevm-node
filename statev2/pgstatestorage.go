package statev2

import (
	"context"
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
	getBatchNumByBlockNum                  = "SELECT batch_num FROM statev2.virtual_batch WHERE block_num = $1 ORDER BY batch_num ASC LIMIT 1"
	getTxsHashesFromBatchNum               = "SELECT hash FROM statev2.transaction WHERE batch_num >= $1"
	getTxsHashesBeforeBatchNum             = "SELECT hash FROM statev2.transaction WHERE batch_num < $1"
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

// GetTxsHashesFromBlockNum get tx hashes with batch num > x
func (p *PostgresStorage) GetTxsHashesFromBlockNum(ctx context.Context, blockNum uint64, dbTx pgx.Tx) ([]common.Hash, error) {
	var batchNum uint64
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getBatchNumByBlockNum, blockNum).Scan(&batchNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	rows, err := e.Query(ctx, getTxsHashesFromBatchNum, batchNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	hashes := make([]common.Hash, 0, len(rows.RawValues()))
	for rows.Next() {
		var hash string
		err := rows.Scan(&hash)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, common.HexToHash(hash))
	}

	return hashes, nil
}

// GetTxsHashesToDelete get txs hashes to delete from tx pool
func (p *PostgresStorage) GetTxsHashesToDelete(ctx context.Context, blockNumDiff uint64, dbTx pgx.Tx) ([]common.Hash, error) {
	var batchNum, blockNum uint64
	e := p.getExecQuerier(dbTx)

	err := e.QueryRow(ctx, getLastBlockNumSQL).Scan(&blockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	blockNum = blockNum - blockNumDiff
	if blockNum <= 0 {
		return nil, errors.New("blockNumDiff is too big, there is no txs to delete")
	}

	err = e.QueryRow(ctx, getBatchNumByBlockNum, blockNum).Scan(&batchNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	rows, err := e.Query(ctx, getTxsHashesBeforeBatchNum, batchNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	hashes := make([]common.Hash, 0, len(rows.RawValues()))
	for rows.Next() {
		var hash string
		err := rows.Scan(&hash)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, common.HexToHash(hash))
	}

	return hashes, nil
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

	b, err := hex.DecodeHex(encoded)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	block.Transactions[0] = tx

	return &block, nil
}
