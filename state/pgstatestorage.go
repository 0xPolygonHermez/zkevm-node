package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const maxTopics = 4

const (
	addGlobalExitRootSQL                     = "INSERT INTO state.exit_root (block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root) VALUES ($1, $2, $3, $4, $5)"
	getLatestExitRootSQL                     = "SELECT block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root FROM state.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
	getLatestExitRootBlockNumSQL             = "SELECT block_num FROM state.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
	addVirtualBatchSQL                       = "INSERT INTO state.virtual_batch (batch_num, tx_hash, coinbase, block_num) VALUES ($1, $2, $3, $4)"
	addForcedBatchSQL                        = "INSERT INTO state.forced_batch (forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase, batch_num, block_num) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	getForcedBatchSQL                        = "SELECT forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase, batch_num, block_num FROM state.forced_batch WHERE forced_batch_num = $1"
	addBlockSQL                              = "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)"
	getLastBlockSQL                          = "SELECT block_num, block_hash, parent_hash, received_at FROM state.block ORDER BY block_num DESC LIMIT 1"
	getPreviousBlockSQL                      = "SELECT block_num, block_hash, parent_hash, received_at FROM state.block ORDER BY block_num DESC LIMIT 1 OFFSET $1"
	resetSQL                                 = "DELETE FROM state.block WHERE block_num > $1"
	resetTrustedStateSQL                     = "DELETE FROM state.batch WHERE batch_num > $1"
	addVerifiedBatchSQL                      = "INSERT INTO state.verified_batch (block_num, batch_num, tx_hash, aggregator) VALUES ($1, $2, $3, $4)"
	getVerifiedBatchSQL                      = "SELECT block_num, batch_num, tx_hash, aggregator FROM state.verified_batch WHERE batch_num = $1"
	getLastBatchNumberSQL                    = "SELECT batch_num FROM state.batch ORDER BY batch_num DESC LIMIT 1"
	getLastNBatchesSQL                       = "SELECT batch_num, global_exit_root, local_exit_root, state_root, timestamp, coinbase, raw_txs_data from state.batch ORDER BY batch_num DESC LIMIT $1"
	getLastBatchTimeSQL                      = "SELECT timestamp FROM state.batch ORDER BY batch_num DESC LIMIT 1"
	getLastVirtualBatchNumSQL                = "SELECT COALESCE(MAX(batch_num), 0) FROM state.virtual_batch"
	getLastVirtualBatchBlockNumSQL           = "SELECT block_num FROM state.virtual_batch ORDER BY batch_num DESC LIMIT 1"
	getLastBlockNumSQL                       = "SELECT block_num FROM state.block ORDER BY block_num DESC LIMIT 1"
	getLastL2BlockNumber                     = "SELECT block_num FROM state.l2block ORDER BY block_num DESC LIMIT 1"
	getBlockTimeByNumSQL                     = "SELECT received_at FROM state.block WHERE block_num = $1"
	getForcedBatchByBatchNumSQL              = "SELECT forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase, batch_num, block_num from state.forced_batch WHERE batch_num = $1"
	getProcessingContextSQL                  = "SELECT batch_num, global_exit_root, timestamp, coinbase from state.batch WHERE batch_num = $1"
	getEncodedTransactionsByBatchNumberSQL   = "SELECT encoded FROM state.transaction t INNER JOIN state.l2block b ON t.l2_block_num = b.block_num WHERE b.batch_num = $1 ORDER BY l2_block_num ASC"
	getTransactionHashesByBatchNumberSQL     = "SELECT hash FROM state.transaction t INNER JOIN state.l2block b ON t.l2_block_num = b.block_num WHERE b.batch_num = $1 ORDER BY l2_block_num ASC"
	getLastBatchSeenSQL                      = "SELECT last_batch_num_seen FROM state.sync_info LIMIT 1"
	updateLastBatchSeenSQL                   = "UPDATE state.sync_info SET last_batch_num_seen = $1"
	resetTrustedBatchSQL                     = "DELETE FROM state.batch WHERE batch_num > $1"
	isBatchClosedSQL                         = "SELECT global_exit_root IS NOT NULL AND state_root IS NOT NULL FROM state.batch WHERE batch_num = $1 LIMIT 1"
	addGenesisBatchSQL                       = `INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, timestamp, coinbase, raw_txs_data) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	openBatchSQL                             = "INSERT INTO state.batch (batch_num, global_exit_root, timestamp, coinbase) VALUES ($1, $2, $3, $4)"
	closeBatchSQL                            = "UPDATE state.batch SET state_root = $1, local_exit_root = $2, raw_txs_data = $3 WHERE batch_num = $4"
	getNextForcedBatchesSQL                  = "SELECT forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase, batch_num, block_num FROM state.forced_batch WHERE batch_num IS NULL LIMIT $1"
	addBatchNumberInForcedBatchSQL           = "UPDATE state.forced_batch SET batch_num = $2 WHERE forced_batch_num = $1"
	getL2BlockByNumberSQL                    = "SELECT header, uncles, received_at FROM state.l2block b WHERE b.block_num = $1"
	getL2BlockHeaderByNumberSQL              = "SELECT header FROM state.l2block b WHERE b.block_num = $1"
	getTransactionByHashSQL                  = "SELECT transaction.encoded FROM state.transaction WHERE hash = $1"
	getTransactionByL2BlockHashAndIndexSQL   = "SELECT t.encoded FROM state.transaction t INNER JOIN state.l2block b ON t.l2_block_num = b.batch_num WHERE b.block_hash = $1 AND 0 = $2"
	getTransactionByL2BlockNumberAndIndexSQL = "SELECT t.encoded FROM state.transaction t WHERE t.l2_block_num = $1 AND 0 = $2"
	getL2BlockTransactionCountByHashSQL      = "SELECT COUNT(*) FROM state.transaction t INNER JOIN state.l2block b ON b.block_num = t.l2_block_num WHERE b.block_hash = $1"
	getL2BlockTransactionCountByNumberSQL    = "SELECT COUNT(*) FROM state.transaction t WHERE t.l2_block_num = $1"
	getLastConsolidatedBlockNumberSQL        = "SELECT b.block_num FROM state.l2block b INNER JOIN state.verified_batch vb ON vb.batch_num = b.batch_num ORDER BY b.block_num DESC LIMIT 1"
	getLastVirtualBlockHeaderSQL             = "SELECT b.header FROM state.l2block b INNER JOIN state.virtual_batch vb ON vb.batch_num = b.batch_num ORDER BY b.block_num DESC LIMIT 1"
	getL2BlockByHashSQL                      = "SELECT header, uncles, received_at FROM state.l2block b WHERE b.block_hash = $1"
	getLastL2BlockSQL                        = "SELECT header, uncles, received_at FROM state.l2block b ORDER BY b.block_num DESC LIMIT 1"
	getL2BlockHeaderByHashSQL                = "SELECT header FROM state.l2block b WHERE b.block_hash = $1"
	getTxsByBlockNumSQL                      = "SELECT encoded FROM state.transaction WHERE l2_block_num = $1"
	getL2BlockHashesSinceSQL                 = "SELECT block_hash FROM state.l2block WHERE received_at >= $1"
	getSyncingInfoSQL                        = `
		SELECT coalesce(MIN(initial_blocks.block_num), 0) as init_sync_block
			 , coalesce(MAX(virtual_blocks.block_num), 0) as last_block_num_seen
			 , coalesce(MAX(consolidated_blocks.block_num), 0) as last_block_num_consolidated
			 , coalesce(MIN(sy.init_sync_batch), 0) as init_sync_batch
			 , coalesce(MIN(sy.last_batch_num_seen), 0) as last_batch_num_seen
			 , coalesce(MIN(sy.last_batch_num_consolidated), 0) as last_batch_num_consolidated
		  FROM state.sync_info sy
		 INNER JOIN state.l2block initial_blocks
			ON initial_blocks.batch_num = sy.init_sync_batch
		 INNER JOIN state.l2block virtual_blocks
			ON virtual_blocks.batch_num = sy.last_batch_num_seen
		 INNER JOIN state.l2block consolidated_blocks
			ON consolidated_blocks.batch_num = sy.last_batch_num_consolidated;
	`
	addTransactionSQL          = "INSERT INTO state.transaction (hash, encoded, decoded, l2_block_num) VALUES($1, $2, $3, $4)"
	getBatchNumByBlockNum      = "SELECT batch_num FROM state.virtual_batch WHERE block_num <= $1 ORDER BY batch_num DESC LIMIT 1"
	getTxsHashesBeforeBatchNum = "SELECT hash FROM state.transaction JOIN state.l2block ON state.transaction.l2_block_num = state.l2block.block_num AND state.l2block.batch_num <= $1"
	isL2BlockVirtualized       = "SELECT l2b.block_num FROM state.l2block l2b INNER JOIN state.virtual_batch vb ON vb.batch_num = l2b.batch_num WHERE l2b.block_num = $1"
	isL2BlockConsolidated      = "SELECT l2b.block_num FROM state.l2block l2b INNER JOIN state.verified_batch vb ON vb.batch_num = l2b.batch_num WHERE l2b.block_num = $1"
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

// Reset resets the state to a block for the given DB tx
func (p *PostgresStorage) Reset(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	if _, err := e.Exec(ctx, resetSQL, blockNumber); err != nil {
		return err
	}
	// TODO: Remove consolidations
	return nil
}

// ResetTrustedState removes the batches with number greater than the given one
// from the database.
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

// GetTxsOlderThanNL1Blocks get txs hashes to delete from tx pool
func (p *PostgresStorage) GetTxsOlderThanNL1Blocks(ctx context.Context, nL1Blocks uint64, dbTx pgx.Tx) ([]common.Hash, error) {
	var batchNum, blockNum uint64
	e := p.getExecQuerier(dbTx)

	err := e.QueryRow(ctx, getLastBlockNumSQL).Scan(&blockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	blockNum = blockNum - nL1Blocks
	if blockNum <= 0 {
		return nil, errors.New("blockNumDiff is too big, there are no txs to delete")
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

// GetLatestGlobalExitRoot get the latest global ExitRoot synced.
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

// GetBlockNumAndMainnetExitRootByGER gets block number and mainnet exit root by the global exit root
func (p *PostgresStorage) GetBlockNumAndMainnetExitRootByGER(ctx context.Context, ger common.Hash, dbTx pgx.Tx) (uint64, common.Hash, error) {
	var (
		blockNum        uint64
		mainnetExitRoot common.Hash
	)
	e := p.getExecQuerier(dbTx)
	const getMainnetExitRoot = "SELECT block_num, mainnet_exit_root FROM state.exit_root WHERE global_exit_root = $1"
	err := e.QueryRow(ctx, getMainnetExitRoot, ger.String()).Scan(&blockNum, &mainnetExitRoot)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, common.Hash{}, ErrNotFound
	} else if err != nil {
		return 0, common.Hash{}, err
	}

	return blockNum, mainnetExitRoot, nil
}

// GetTimeForLatestBatchVirtualization returns the timestamp of the latest
// virtual batch.
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

// GetForcedBatchByBatchNumber gets an L1 forcedBatch by batch number.
func (p *PostgresStorage) GetForcedBatchByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*ForcedBatch, error) {
	var (
		forcedBatch    ForcedBatch
		globalExitRoot string
		rawTxs         string
		seq            string
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getForcedBatchByBatchNumSQL, batchNumber).Scan(&forcedBatch.ForcedBatchNumber, &globalExitRoot, &forcedBatch.ForcedAt, &rawTxs, &seq, &forcedBatch.BatchNumber, &forcedBatch.BlockNumber)
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

// GetLastNBatches returns the last numBatches batches.
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
		batch, err := scanBatch(rows)
		if err != nil {
			return nil, err
		}
		batches = append(batches, &batch)
	}

	return batches, nil
}

// GetLastNBatchesByL2BlockNumber returns the last numBatches batches along with the l2 block state root by l2BlockNumber
// if the l2BlockNumber parameter is nil, it means we want to get the most recent last N batches
func (p *PostgresStorage) GetLastNBatchesByL2BlockNumber(ctx context.Context, l2BlockNumber *uint64, numBatches uint, dbTx pgx.Tx) ([]*Batch, common.Hash, error) {
	const getLastNBatchesByBlockNumberSQL = `
        SELECT b.batch_num,
               b.global_exit_root,
               b.local_exit_root,
               b.state_root,
               b.timestamp,
               b.coinbase,
               b.raw_txs_data,
               /* gets the state root of the l2 block with the highest number associated to the batch in the row */
               (SELECT l2b1.header->>'stateRoot'
                  FROM state.l2block l2b1
                 WHERE l2b1.block_num = (SELECT MAX(l2b2.block_num)
                                           FROM state.l2block l2b2
                                          WHERE l2b2.batch_num = b.batch_num)) as l2_block_state_root
          FROM state.batch b
               /* if there is a value for the parameter $1 (l2 block number), filter the batches with batch number
                * smaller or equal than the batch associated to the l2 block number */
         WHERE ($1::int8 IS NOT NULL AND b.batch_num <= (SELECT MAX(l2b.batch_num)
                                                           FROM state.l2block l2b
                                                          WHERE l2b.block_num = $1))
               /* OR if $1 is null, this means we want to get the most updated information from state, so it considers all the batches.
                * this is generally used by estimate gas, process unsigned transactions and it is required by claim transactions to add
                * the open batch to the result and get the most updated GER synced from L1 and stored in the current open batch when 
                * there was not transactions yet to create a l2 block with it */
            OR $1 IS NULL
         ORDER BY b.batch_num DESC
         LIMIT $2;`

	var l2BlockStateRoot *common.Hash
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getLastNBatchesByBlockNumberSQL, l2BlockNumber, numBatches)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, common.Hash{}, ErrStateNotSynchronized
	} else if err != nil {
		return nil, common.Hash{}, err
	}
	defer rows.Close()

	batches := make([]*Batch, 0, len(rows.RawValues()))

	for rows.Next() {
		batch, _l2BlockStateRoot, err := scanBatchWithL2BlockStateRoot(rows)
		if err != nil {
			return nil, common.Hash{}, err
		}
		batches = append(batches, &batch)
		if l2BlockStateRoot == nil && _l2BlockStateRoot != nil {
			l2BlockStateRoot = _l2BlockStateRoot
		}
	}

	return batches, *l2BlockStateRoot, nil
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

// GetBatchByNumber returns the batch with the given number.
func (p *PostgresStorage) GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*Batch, error) {
	const getBatchByNumberSQL = `
		SELECT batch_num, global_exit_root, local_exit_root, state_root, timestamp, coinbase, raw_txs_data
		  FROM state.batch 
		 WHERE batch_num = $1`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getBatchByNumberSQL, batchNumber)
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetBatchByTxHash returns the batch including the given tx
func (p *PostgresStorage) GetBatchByTxHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*Batch, error) {
	const getBatchByTxHashSQL = `
		SELECT b.batch_num, b.global_exit_root, b.local_exit_root, b.state_root, b.timestamp, b.coinbase, b.raw_txs_data
		  FROM state.transaction t, state.batch b, state.l2block l 
		  WHERE t.hash = $1 AND l.block_num = t.l2_block_num AND b.batch_num = l.batch_num`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getBatchByTxHashSQL, transactionHash.String())
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetBatchByL2BlockNumber returns the batch related to the l2 block accordingly to the provided l2 block number.
func (p *PostgresStorage) GetBatchByL2BlockNumber(ctx context.Context, l2BlockNumber uint64, dbTx pgx.Tx) (*Batch, error) {
	const getBatchByL2BlockNumberSQL = `
		SELECT bt.batch_num, bt.global_exit_root, bt.local_exit_root, bt.state_root, bt.timestamp, bt.coinbase, bt.raw_txs_data 
		  FROM state.batch bt
		 INNER JOIN state.l2block bl
		    ON bt.batch_num = bl.batch_num
		 WHERE bl.block_num = $1
		 LIMIT 1;`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getBatchByL2BlockNumberSQL, l2BlockNumber)
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetVirtualBatchByNumber gets batch from batch table that exists on virtual batch
func (p *PostgresStorage) GetVirtualBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*Batch, error) {
	const query = `
		SELECT
			batch_num,
			global_exit_root,
			local_exit_root,
			state_root,
			timestamp,
			coinbase,
			raw_txs_data
		FROM
			state.batch
		WHERE
			batch_num = $1 AND
			EXISTS (SELECT batch_num FROM state.virtual_batch WHERE batch_num = $1)
		`
	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, query, batchNumber)
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &batch, nil
}

// IsBatchVirtualized checks if batch is virtualized
func (p *PostgresStorage) IsBatchVirtualized(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM state.virtual_batch WHERE batch_num = $1)`
	e := p.getExecQuerier(dbTx)
	var exists bool
	err := e.QueryRow(ctx, query, batchNumber).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return exists, err
	}
	return exists, nil
}

// IsSequencingTXSynced checks if sequencing tx has been synced into the state
func (p *PostgresStorage) IsSequencingTXSynced(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM state.virtual_batch WHERE tx_hash = $1)`
	e := p.getExecQuerier(dbTx)
	var exists bool
	err := e.QueryRow(ctx, query, transactionHash.String()).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return exists, err
	}
	return exists, nil
}

// GetProcessingContext returns the processing context for the given batch.
func (p *PostgresStorage) GetProcessingContext(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*ProcessingContext, error) {
	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getProcessingContextSQL, batchNumber)
	processingContext := ProcessingContext{}
	var (
		gerStr      string
		coinbaseStr string
	)
	if err := row.Scan(
		&processingContext.BatchNumber,
		&gerStr,
		&processingContext.Timestamp,
		&coinbaseStr,
	); errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	processingContext.GlobalExitRoot = common.HexToHash(gerStr)
	processingContext.Coinbase = common.HexToAddress(coinbaseStr)

	return &processingContext, nil
}

func scanBatch(row pgx.Row) (Batch, error) {
	batch := Batch{}
	var (
		gerStr      string
		lerStr      *string
		stateStr    *string
		coinbaseStr string
	)
	if err := row.Scan(
		&batch.BatchNumber,
		&gerStr,
		&lerStr,
		&stateStr,
		&batch.Timestamp,
		&coinbaseStr,
		&batch.BatchL2Data,
	); err != nil {
		return batch, err
	}
	batch.GlobalExitRoot = common.HexToHash(gerStr)
	if lerStr != nil {
		batch.LocalExitRoot = common.HexToHash(*lerStr)
	}
	if stateStr != nil {
		batch.StateRoot = common.HexToHash(*stateStr)
	}

	batch.Coinbase = common.HexToAddress(coinbaseStr)
	return batch, nil
}

func scanBatchWithL2BlockStateRoot(row pgx.Row) (Batch, *common.Hash, error) {
	batch := Batch{}
	var (
		gerStr              string
		lerStr              *string
		stateStr            *string
		coinbaseStr         string
		l2BlockStateRootStr *string
	)
	if err := row.Scan(
		&batch.BatchNumber,
		&gerStr,
		&lerStr,
		&stateStr,
		&batch.Timestamp,
		&coinbaseStr,
		&batch.BatchL2Data,
		&l2BlockStateRootStr,
	); err != nil {
		return batch, nil, err
	}
	batch.GlobalExitRoot = common.HexToHash(gerStr)
	if lerStr != nil {
		batch.LocalExitRoot = common.HexToHash(*lerStr)
	}
	if stateStr != nil {
		batch.StateRoot = common.HexToHash(*stateStr)
	}
	var l2BlockStateRoot *common.Hash
	if l2BlockStateRootStr != nil {
		h := common.HexToHash(*l2BlockStateRootStr)
		l2BlockStateRoot = &h
	}

	batch.Coinbase = common.HexToAddress(coinbaseStr)
	return batch, l2BlockStateRoot, nil
}

// GetEncodedTransactionsByBatchNumber returns the encoded field of all
// transactions in the given batch.
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

// GetTransactionsByBatchNumber returns the transactions in the given batch.
func (p *PostgresStorage) GetTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (txs []types.Transaction, err error) {
	encodedTxs, err := p.GetEncodedTransactionsByBatchNumber(ctx, batchNumber, dbTx)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(encodedTxs); i++ {
		tx, err := DecodeTx(encodedTxs[i])
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}
	return
}

// GetTxsHashesByBatchNumber returns the hashes of the transactions in the
// given batch.
func (p *PostgresStorage) GetTxsHashesByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (encoded []common.Hash, err error) {
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getTransactionHashesByBatchNumberSQL, batchNumber)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]common.Hash, 0, len(rows.RawValues()))

	for rows.Next() {
		var hexHash string
		err := rows.Scan(&hexHash)
		if err != nil {
			return nil, err
		}

		txs = append(txs, common.HexToHash(hexHash))
	}
	return txs, nil
}

// ResetTrustedBatch resets the batches which the batch number is higher than the input.
func (p *PostgresStorage) ResetTrustedBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, resetTrustedBatchSQL, batchNumber)
	return err
}

// AddVirtualBatch adds a new virtual batch to the storage.
func (p *PostgresStorage) AddVirtualBatch(ctx context.Context, virtualBatch *VirtualBatch, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addVirtualBatchSQL, virtualBatch.BatchNumber, virtualBatch.TxHash.String(), virtualBatch.Coinbase.String(), virtualBatch.BlockNumber)
	return err
}

func (p *PostgresStorage) storeGenesisBatch(ctx context.Context, batch Batch, dbTx pgx.Tx) error {
	if batch.BatchNumber != 0 {
		return fmt.Errorf("unexpected batch number. Got %d, should be 0", batch.BatchNumber)
	}
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(
		ctx,
		addGenesisBatchSQL,
		batch.BatchNumber,
		batch.GlobalExitRoot.String(),
		batch.LocalExitRoot.String(),
		batch.StateRoot.String(),
		batch.Timestamp.UTC(),
		batch.Coinbase.String(),
		batch.BatchL2Data,
	)

	return err
}

// openBatch adds a new batch into the state, with the necessary data to start processing transactions within it.
// It's meant to be used by sequencers, since they don't necessarely know what transactions are going to be added
// in this batch yet. In other words it's the creation of a WIP batch.
// Note that this will add a batch with batch number N + 1, where N it's the greates batch number on the state.
func (p *PostgresStorage) openBatch(ctx context.Context, batchContext ProcessingContext, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(
		ctx, openBatchSQL,
		batchContext.BatchNumber,
		batchContext.GlobalExitRoot.String(),
		batchContext.Timestamp.UTC(),
		batchContext.Coinbase.String(),
	)
	return err
}

func (p *PostgresStorage) closeBatch(ctx context.Context, receipt ProcessingReceipt, rawTxs []byte, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, closeBatchSQL, receipt.StateRoot.String(), receipt.LocalExitRoot.String(), rawTxs, receipt.BatchNumber)
	return err
}

// UpdateGERInOpenBatch update ger in open batch
func (p *PostgresStorage) UpdateGERInOpenBatch(ctx context.Context, ger common.Hash, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	var (
		batchNumber   uint64
		isBatchHasTxs bool
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastBatchNumberSQL).Scan(&batchNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrStateNotSynchronized
	}

	const isBatchHasTxsQuery = `SELECT EXISTS (SELECT 1 FROM state.l2block WHERE batch_num = $1)`
	err = e.QueryRow(ctx, isBatchHasTxsQuery, batchNumber).Scan(&isBatchHasTxs)
	if err != nil {
		return err
	}

	if isBatchHasTxs {
		return errors.New("batch has txs, can't change GER")
	}

	const updateGER = `
			UPDATE 
    			state.batch
			SET global_exit_root = $1, timestamp = $2
			WHERE batch_num = $3
				AND state_root IS NULL`
	_, err = e.Exec(ctx, updateGER, ger.String(), time.Now().UTC(), batchNumber)
	return err
}

// IsBatchClosed indicates if the batch referenced by batchNum is closed or not
func (p *PostgresStorage) IsBatchClosed(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (bool, error) {
	q := p.getExecQuerier(dbTx)
	var isClosed bool
	err := q.QueryRow(ctx, isBatchClosedSQL, batchNum).Scan(&isClosed)
	return isClosed, err
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

// GetL2BlockByNumber gets a l2 block by its number
func (p *PostgresStorage) GetL2BlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*types.Block, error) {
	header := &types.Header{}
	uncles := []*types.Header{}
	receivedAt := time.Time{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getL2BlockByNumberSQL, blockNumber).
		Scan(&header, &uncles, &receivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	transactions, err := p.GetTxsByBlockNumber(ctx, header.Number.Uint64(), dbTx)
	if errors.Is(err, pgx.ErrNoRows) {
		transactions = []*types.Transaction{}
	} else if err != nil {
		return nil, err
	}

	block := types.NewBlockWithHeader(header).WithBody(transactions, uncles)
	block.ReceivedAt = receivedAt
	return block, nil
}

// GetTransactionByHash gets a transaction accordingly to the provided transaction hash
func (p *PostgresStorage) GetTransactionByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error) {
	var encoded string
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByHashSQL, transactionHash.String()).Scan(&encoded)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := DecodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionReceipt gets a transaction receipt accordingly to the provided transaction hash
func (p *PostgresStorage) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Receipt, error) {
	var txHash, encodedTx, contractAddress, l2BlockHash string
	var l2BlockNum uint64

	const getReceiptSQL = `
		SELECT r.tx_hash
		     , r.type
			 , r.post_state
			 , r.status
			 , r.cumulative_gas_used
			 , r.gas_used
			 , r.contract_address
			 , t.encoded
			 , t.l2_block_num
			 , b.block_hash
	      FROM state.receipt r
		 INNER JOIN state.transaction t
		    ON t.hash = r.tx_hash
		 INNER JOIN state.l2block b
		    ON b.block_num = t.l2_block_num
		 WHERE r.tx_hash = $1`

	receipt := types.Receipt{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getReceiptSQL, transactionHash.String()).
		Scan(&txHash,
			&receipt.Type,
			&receipt.PostState,
			&receipt.Status,
			&receipt.CumulativeGasUsed,
			&receipt.GasUsed,
			&contractAddress,
			&encodedTx,
			&l2BlockNum,
			&l2BlockHash,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	receipt.TxHash = common.HexToHash(txHash)
	receipt.ContractAddress = common.HexToAddress(contractAddress)

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

// GetTransactionByL2BlockHashAndIndex gets a transaction accordingly to the block hash and transaction index provided.
// since we only have a single transaction per l2 block, any index different from 0 will return a not found result
func (p *PostgresStorage) GetTransactionByL2BlockHashAndIndex(ctx context.Context, blockHash common.Hash, index uint64, dbTx pgx.Tx) (*types.Transaction, error) {
	var encoded string
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByL2BlockHashAndIndexSQL, blockHash.String(), index).Scan(&encoded)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := DecodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionByL2BlockNumberAndIndex gets a transaction accordingly to the block number and transaction index provided.
// since we only have a single transaction per l2 block, any index different from 0 will return a not found result
func (p *PostgresStorage) GetTransactionByL2BlockNumberAndIndex(ctx context.Context, blockNumber uint64, index uint64, dbTx pgx.Tx) (*types.Transaction, error) {
	var encoded string
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByL2BlockNumberAndIndexSQL, blockNumber, index).Scan(&encoded)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := DecodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetL2BlockTransactionCountByHash returns the number of transactions related to the provided block hash
func (p *PostgresStorage) GetL2BlockTransactionCountByHash(ctx context.Context, blockHash common.Hash, dbTx pgx.Tx) (uint64, error) {
	var count uint64
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getL2BlockTransactionCountByHashSQL, blockHash.String()).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetL2BlockTransactionCountByNumber returns the number of transactions related to the provided block number
func (p *PostgresStorage) GetL2BlockTransactionCountByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	var count uint64
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getL2BlockTransactionCountByNumberSQL, blockNumber).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// getTransactionLogs returns the logs of a transaction by transaction hash
func (p *PostgresStorage) getTransactionLogs(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) ([]*types.Log, error) {
	q := p.getExecQuerier(dbTx)

	const getTransactionLogsSQL = `
	SELECT t.l2_block_num, b.block_hash, l.tx_hash, l.log_index, l.address, l.data, l.topic0, l.topic1, l.topic2, l.topic3
	FROM state.log l
	INNER JOIN state.transaction t ON t.hash = l.tx_hash
	INNER JOIN state.l2block b ON b.block_num = t.l2_block_num 
	WHERE t.hash = $1`
	rows, err := q.Query(ctx, getTransactionLogsSQL, transactionHash.String())
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	return scanLogs(rows)
}

func scanLogs(rows pgx.Rows) ([]*types.Log, error) {
	defer rows.Close()

	logs := make([]*types.Log, 0, len(rows.RawValues()))

	for rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		var log types.Log
		var blockHash, txHash, logAddress, logData, topic0 string
		var topic1, topic2, topic3 *string

		err := rows.Scan(&log.BlockNumber, &blockHash, &txHash, &log.Index,
			&logAddress, &logData, &topic0, &topic1, &topic2, &topic3)
		if err != nil {
			return nil, err
		}

		log.BlockHash = common.HexToHash(blockHash)
		log.TxHash = common.HexToHash(txHash)
		log.Address = common.HexToAddress(logAddress)
		log.TxIndex = uint(0)
		log.Data, err = hex.DecodeHex(logData)
		if err != nil {
			return nil, err
		}
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

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return logs, nil
}

// AddL2Block adds a new L2 block to the State Store
func (p *PostgresStorage) AddL2Block(ctx context.Context, batchNumber uint64, l2Block *types.Block, receipts []*types.Receipt, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)

	const addL2BlockSQL = `
        INSERT INTO state.l2block (block_num, block_hash, header, uncles, parent_hash, state_root, received_at, batch_num)
                           VALUES (       $1,         $2,     $3,     $4,          $5,         $6,          $7,        $8)`

	var header = "{}"
	if l2Block.Header() != nil {
		headerBytes, err := json.Marshal(l2Block.Header())
		if err != nil {
			return err
		}
		header = string(headerBytes)
	}

	var uncles = "[]"
	if l2Block.Uncles() != nil {
		unclesBytes, err := json.Marshal(l2Block.Uncles())
		if err != nil {
			return err
		}
		uncles = string(unclesBytes)
	}

	if _, err := e.Exec(ctx, addL2BlockSQL,
		l2Block.Number().Uint64(), l2Block.Hash().String(), header, uncles,
		l2Block.ParentHash().String(), l2Block.Root().String(),
		l2Block.ReceivedAt, batchNumber); err != nil {
		return err
	}

	for _, tx := range l2Block.Transactions() {
		binary, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		encoded := hex.EncodeToHex(binary)

		binary, err = tx.MarshalJSON()
		if err != nil {
			return err
		}
		decoded := string(binary)
		_, err = e.Exec(ctx, addTransactionSQL, tx.Hash().String(), encoded, decoded, l2Block.Number().Uint64())
		if err != nil {
			return err
		}
	}

	for _, receipt := range receipts {
		err := p.AddReceipt(ctx, receipt, dbTx)
		if err != nil {
			return err
		}

		for _, log := range receipt.Logs {
			err := p.AddLog(ctx, log, dbTx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetLastConsolidatedL2BlockNumber gets the last l2 block verified
func (p *PostgresStorage) GetLastConsolidatedL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var lastConsolidatedBlockNumber uint64
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getLastConsolidatedBlockNumberSQL).Scan(&lastConsolidatedBlockNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return lastConsolidatedBlockNumber, nil
}

// GetLastL2BlockNumber gets the last l2 block number
func (p *PostgresStorage) GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var lastBlockNumber uint64
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getLastL2BlockNumber).Scan(&lastBlockNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrStateNotSynchronized
	} else if err != nil {
		return 0, err
	}

	return lastBlockNumber, nil
}

// GetLastL2BlockHeader gets the last l2 block number
func (p *PostgresStorage) GetLastL2BlockHeader(ctx context.Context, dbTx pgx.Tx) (*types.Header, error) {
	header := &types.Header{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getLastVirtualBlockHeaderSQL).Scan(&header)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	return header, nil
}

// GetLastL2Block retrieves the latest L2 Block from the State data base
func (p *PostgresStorage) GetLastL2Block(ctx context.Context, dbTx pgx.Tx) (*types.Block, error) {
	var (
		headerStr  string
		unclesStr  string
		receivedAt time.Time
	)
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getLastL2BlockSQL).Scan(&headerStr, &unclesStr, &receivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	header := &types.Header{}
	uncles := []*types.Header{}

	if err := json.Unmarshal([]byte(headerStr), header); err != nil {
		return nil, err
	}
	if unclesStr != "[]" {
		if err := json.Unmarshal([]byte(unclesStr), &uncles); err != nil {
			return nil, err
		}
	}

	transactions, err := p.GetTxsByBlockNumber(ctx, header.Number.Uint64(), dbTx)
	if errors.Is(err, pgx.ErrNoRows) {
		transactions = []*types.Transaction{}
	} else if err != nil {
		return nil, err
	}

	block := types.NewBlockWithHeader(header).WithBody(transactions, uncles)
	block.ReceivedAt = receivedAt
	return block, nil
}

// GetLastVerifiedBatchNumberSeenOnEthereum gets last verified batch number seen on ethereum
func (p *PostgresStorage) GetLastVerifiedBatchNumberSeenOnEthereum(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	const getLastVerifiedBatchSeenSQL = "SELECT last_batch_num_verified FROM state.sync_info LIMIT 1"
	var batchNumber uint64
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastVerifiedBatchSeenSQL).Scan(&batchNumber)
	if err != nil {
		return 0, err
	}
	return batchNumber, nil
}

// GetLastVerifiedBatch gets last verified batch
func (p *PostgresStorage) GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*VerifiedBatch, error) {
	const query = "SELECT block_num, batch_num, tx_hash, aggregator FROM state.verified_batch ORDER BY batch_num DESC LIMIT 1"
	var (
		verifiedBatch VerifiedBatch
		txHash, agg   string
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, query).Scan(&verifiedBatch.BlockNumber, &verifiedBatch.BatchNumber, &txHash, &agg)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	verifiedBatch.Aggregator = common.HexToAddress(agg)
	verifiedBatch.TxHash = common.HexToHash(txHash)
	return &verifiedBatch, nil
}

// GetStateRootByBatchNumber get state root by batch number
func (p *PostgresStorage) GetStateRootByBatchNumber(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (common.Hash, error) {
	const query = "SELECT state_root FROM state.batch WHERE batch_num = $1"
	var stateRootStr string
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, query, batchNum).Scan(&stateRootStr)
	if errors.Is(err, pgx.ErrNoRows) {
		return common.Hash{}, ErrNotFound
	} else if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(stateRootStr), nil
}

// GetLocalExitRootByBatchNumber get local exit root by batch number
func (p *PostgresStorage) GetLocalExitRootByBatchNumber(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (common.Hash, error) {
	const query = "SELECT local_exit_root FROM state.batch WHERE batch_num = $1"
	var localExitRootStr string
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, query, batchNum).Scan(&localExitRootStr)
	if errors.Is(err, pgx.ErrNoRows) {
		return common.Hash{}, ErrNotFound
	} else if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(localExitRootStr), nil
}

// GetBlockNumVirtualBatchByBatchNum get block num of virtual batch by block num
func (p *PostgresStorage) GetBlockNumVirtualBatchByBatchNum(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (uint64, error) {
	const query = "SELECT block_num FROM state.virtual_batch WHERE batch_num = $1"
	var blockNum uint64
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, query, batchNum).Scan(&blockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}
	return blockNum, nil
}

// GetL2BlockByHash gets a l2 block from its hash
func (p *PostgresStorage) GetL2BlockByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (*types.Block, error) {
	header := &types.Header{}
	uncles := []*types.Header{}
	receivedAt := time.Time{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getL2BlockByHashSQL, hash.String()).
		Scan(&header, &uncles, &receivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	transactions, err := p.GetTxsByBlockNumber(ctx, header.Number.Uint64(), dbTx)
	if errors.Is(err, pgx.ErrNoRows) {
		transactions = []*types.Transaction{}
	} else if err != nil {
		return nil, err
	}

	block := types.NewBlockWithHeader(header).WithBody(transactions, uncles)
	block.ReceivedAt = receivedAt
	return block, nil
}

// GetTxsByBlockNumber returns all the txs in a given block
func (p *PostgresStorage) GetTxsByBlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) ([]*types.Transaction, error) {
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, getTxsByBlockNumSQL, blockNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	txs := make([]*types.Transaction, 0, len(rows.RawValues()))
	var encoded string
	for rows.Next() {
		if err = rows.Scan(&encoded); err != nil {
			return nil, err
		}

		tx, err := DecodeTx(encoded)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

// GetL2BlockHeaderByHash gets the block header by block number
func (p *PostgresStorage) GetL2BlockHeaderByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (*types.Header, error) {
	header := &types.Header{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getL2BlockHeaderByHashSQL, hash.String()).Scan(&header)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return header, nil
}

// GetL2BlockHeaderByNumber gets the block header by block number
func (p *PostgresStorage) GetL2BlockHeaderByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*types.Header, error) {
	header := &types.Header{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getL2BlockHeaderByNumberSQL, blockNumber).Scan(&header)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return header, nil
}

// GetL2BlockHashesSince gets the block hashes added since the provided date
func (p *PostgresStorage) GetL2BlockHashesSince(ctx context.Context, since time.Time, dbTx pgx.Tx) ([]common.Hash, error) {
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, getL2BlockHashesSinceSQL, since)
	if errors.Is(err, pgx.ErrNoRows) {
		return []common.Hash{}, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	blockHashes := make([]common.Hash, 0, len(rows.RawValues()))

	for rows.Next() {
		var blockHash string
		err := rows.Scan(&blockHash)
		if err != nil {
			return nil, err
		}

		blockHashes = append(blockHashes, common.HexToHash(blockHash))
	}

	return blockHashes, nil
}

// IsL2BlockConsolidated checks if the block ID is consolidated
func (p *PostgresStorage) IsL2BlockConsolidated(ctx context.Context, blockNumber int, dbTx pgx.Tx) (bool, error) {
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, isL2BlockConsolidated, blockNumber)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	isConsolidated := rows.Next()

	if rows.Err() != nil {
		return false, rows.Err()
	}

	return isConsolidated, nil
}

// IsL2BlockVirtualized checks if the block  ID is virtualized
func (p *PostgresStorage) IsL2BlockVirtualized(ctx context.Context, blockNumber int, dbTx pgx.Tx) (bool, error) {
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, isL2BlockVirtualized, blockNumber)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	isVirtualized := rows.Next()

	if rows.Err() != nil {
		return false, rows.Err()
	}

	return isVirtualized, nil
}

// GetLogs returns the logs that match the filter
func (p *PostgresStorage) GetLogs(ctx context.Context, fromBlock uint64, toBlock uint64, addresses []common.Address, topics [][]common.Hash, blockHash *common.Hash, since *time.Time, dbTx pgx.Tx) ([]*types.Log, error) {
	const getLogsByBlockHashSQL = `
	  SELECT t.l2_block_num, b.block_hash, l.tx_hash, l.log_index, l.address, l.data, l.topic0, l.topic1, l.topic2, l.topic3
		FROM state.log l
	   INNER JOIN state.transaction t ON t.hash = l.tx_hash
	   INNER JOIN state.l2block b ON b.block_num = t.l2_block_num
	   WHERE b.block_hash = $1`
	const getLogsByFilterSQL = `
	  SELECT t.l2_block_num, b.block_hash, l.tx_hash, l.log_index, l.address, l.data, l.topic0, l.topic1, l.topic2, l.topic3
	    FROM state.log l
	   INNER JOIN state.transaction t ON t.hash = l.tx_hash
	   INNER JOIN state.l2block b ON b.block_num = t.l2_block_num
	   WHERE b.block_num BETWEEN $1 AND $2 AND (l.address = any($3) OR $3 IS NULL)
	     AND (l.topic0 = any($4) OR $4 IS NULL)
		 AND (l.topic1 = any($5) OR $5 IS NULL)
		 AND (l.topic2 = any($6) OR $6 IS NULL)
		 AND (l.topic3 = any($7) OR $7 IS NULL)
		 AND (b.received_at >= $8 OR $8 IS NULL)
		ORDER BY b.block_num ASC`

	var err error
	var rows pgx.Rows
	q := p.getExecQuerier(dbTx)
	if blockHash != nil {
		rows, err = q.Query(ctx, getLogsByBlockHashSQL, blockHash.String())
	} else {
		args := []interface{}{fromBlock, toBlock}

		if len(addresses) > 0 {
			args = append(args, p.addressesToHex(addresses))
		} else {
			args = append(args, nil)
		}

		for i := 0; i < maxTopics; i++ {
			if len(topics) > i && len(topics[i]) > 0 {
				args = append(args, p.hashesToHex(topics[i]))
			} else {
				args = append(args, nil)
			}
		}

		args = append(args, since)

		rows, err = q.Query(ctx, getLogsByFilterSQL, args...)
	}

	if err != nil {
		return nil, err
	}
	return scanLogs(rows)
}

// GetSyncingInfo returns information regarding the syncing status of the node
func (p *PostgresStorage) GetSyncingInfo(ctx context.Context, dbTx pgx.Tx) (SyncingInfo, error) {
	var info SyncingInfo
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getSyncingInfoSQL).
		Scan(&info.InitialSyncingBlock, &info.LastBlockNumberSeen, &info.LastBlockNumberConsolidated,
			&info.InitialSyncingBatch, &info.LastBatchNumberSeen, &info.LastBatchNumberConsolidated)
	if err != nil {
		return SyncingInfo{}, nil
	}

	lastBlockNumber, err := p.GetLastL2BlockNumber(ctx, dbTx)
	if err != nil {
		return SyncingInfo{}, nil
	}
	info.CurrentBlockNumber = lastBlockNumber

	lastBatchNumber, err := p.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		return SyncingInfo{}, nil
	}
	info.CurrentBatchNumber = lastBatchNumber

	return info, err
}

func (p *PostgresStorage) addressesToHex(addresses []common.Address) []string {
	converted := make([]string, 0, len(addresses))

	for _, address := range addresses {
		converted = append(converted, address.String())
	}

	return converted
}

func (p *PostgresStorage) hashesToHex(hashes []common.Hash) []string {
	converted := make([]string, 0, len(hashes))

	for _, hash := range hashes {
		converted = append(converted, hash.String())
	}

	return converted
}

// AddReceipt adds a new receipt to the State Store
func (p *PostgresStorage) AddReceipt(ctx context.Context, receipt *types.Receipt, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	const addReceiptSQL = `
        INSERT INTO state.receipt (tx_hash, type, post_state, status, cumulative_gas_used, gas_used, block_num, tx_index, contract_address)
                           VALUES (     $1,   $2,         $3,     $4,                  $5,       $6,        $7,       $8,               $9)`
	_, err := e.Exec(ctx, addReceiptSQL, receipt.TxHash.String(), receipt.Type, receipt.PostState, receipt.Status, receipt.CumulativeGasUsed, receipt.GasUsed, receipt.BlockNumber.Uint64(), receipt.TransactionIndex, receipt.ContractAddress.String())
	return err
}

// AddLog adds a new log to the State Store
func (p *PostgresStorage) AddLog(ctx context.Context, l *types.Log, dbTx pgx.Tx) error {
	const addLogSQL = `INSERT INTO state.log (tx_hash, log_index, address, data, topic0, topic1, topic2, topic3)
	                                  VALUES (     $1,        $2,      $3,   $4,     $5,     $6,     $7,     $8)`

	var topicsAsHex [maxTopics]*string
	for i := 0; i < len(l.Topics); i++ {
		topicHex := l.Topics[i].String()
		topicsAsHex[i] = &topicHex
	}

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addLogSQL,
		l.TxHash.String(), l.Index, l.Address.String(), hex.EncodeToHex(l.Data),
		topicsAsHex[0], topicsAsHex[1], topicsAsHex[2], topicsAsHex[3])
	return err
}

// GetExitRootByGlobalExitRoot returns the mainnet and rollup exit root given
// a global exit root number.
func (p *PostgresStorage) GetExitRootByGlobalExitRoot(ctx context.Context, ger common.Hash, dbTx pgx.Tx) (*GlobalExitRoot, error) {
	var (
		exitRoot  GlobalExitRoot
		globalNum uint64
		err       error
	)

	const sql = "SELECT block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root FROM state.exit_root WHERE global_exit_root = $1 ORDER BY block_num DESC LIMIT 1"

	e := p.getExecQuerier(dbTx)
	err = e.QueryRow(ctx, sql, ger).Scan(&exitRoot.BlockNumber, &globalNum, &exitRoot.MainnetExitRoot, &exitRoot.RollupExitRoot, &exitRoot.GlobalExitRoot)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	exitRoot.GlobalExitRootNum = new(big.Int).SetUint64(globalNum)
	return &exitRoot, nil
}

// AddGeneratedProof adds a generated proof to the storage
func (p *PostgresStorage) AddGeneratedProof(ctx context.Context, proof *Proof, dbTx pgx.Tx) error {
	const addGeneratedProofSQL = "INSERT INTO state.proof (batch_num, proof, proof_id, input_prover, prover) VALUES ($1, $2, $3, $4, $5)"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addGeneratedProofSQL, proof.BatchNumber, proof.Proof, proof.ProofID, proof.InputProver, proof.Prover)
	return err
}

// UpdateGeneratedProof updates a generated proof in the storage
func (p *PostgresStorage) UpdateGeneratedProof(ctx context.Context, proof *Proof, dbTx pgx.Tx) error {
	const addGeneratedProofSQL = "UPDATE state.proof SET proof = $2, proof_id = $3, input_prover = $4, prover = $5 WHERE batch_num = $1"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addGeneratedProofSQL, proof.BatchNumber, proof.Proof, proof.ProofID, proof.InputProver, proof.Prover)
	return err
}

// GetGeneratedProofByBatchNumber gets a generated proof from the storage
func (p *PostgresStorage) GetGeneratedProofByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*Proof, error) {
	var (
		proof *Proof = &Proof{}
		err   error
	)

	const getGeneratedProofSQL = "SELECT batch_num, proof, proof_id, input_prover, prover FROM state.proof WHERE batch_num = $1"
	e := p.getExecQuerier(dbTx)
	err = e.QueryRow(ctx, getGeneratedProofSQL, batchNumber).Scan(&proof.BatchNumber, &proof.Proof, &proof.ProofID, &proof.InputProver, &proof.Prover)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return proof, err
}

// DeleteGeneratedProof deletes a generated proof from the storage
func (p *PostgresStorage) DeleteGeneratedProof(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	const deleteGeneratedProofSQL = "DELETE FROM state.proof WHERE batch_num = $1"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, deleteGeneratedProofSQL, batchNumber)
	return err
}

// DeleteUngeneratedProofs deletes ungenerated proofs from state.proof table
// This method is meant to be use during aggregator boot-up sequence
func (p *PostgresStorage) DeleteUngeneratedProofs(ctx context.Context, dbTx pgx.Tx) error {
	const deleteUngeneratedProofsSQL = "DELETE FROM state.proof WHERE proof is null"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, deleteUngeneratedProofsSQL)
	return err
}

// GetWIPProofByProver gets a generated proof from its prover URI
func (p *PostgresStorage) GetWIPProofByProver(ctx context.Context, prover string, dbTx pgx.Tx) (*Proof, error) {
	var (
		proof *Proof = &Proof{}
		err   error
	)

	const getGeneratedProofSQL = "SELECT batch_num, proof, proof_id, input_prover, prover FROM state.proof WHERE prover = $1 and proof is null"
	e := p.getExecQuerier(dbTx)
	err = e.QueryRow(ctx, getGeneratedProofSQL, prover).Scan(&proof.BatchNumber, &proof.Proof, &proof.ProofID, &proof.InputProver, &proof.Prover)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return proof, err
}
