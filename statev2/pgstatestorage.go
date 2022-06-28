package statev2

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	addGlobalExitRootSQL                   = "INSERT INTO statev2.exit_root (block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root) VALUES ($1, $2, $3, $4, $5)"
	getLatestExitRootSQL                   = "SELECT block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root FROM statev2.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
	getLatestExitRootBlockNumSQL           = "SELECT block_num FROM statev2.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
	addForcedBatchSQL                      = "INSERT INTO statev2.forced_batch (block_num, forced_batch_num, global_exit_root, timestamp, raw_txs_data, sequencer, batch_num) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	getForcedBatchSQL                      = "SELECT block_num, batch_num, forced_batch_num, global_exit_root, timestamp, raw_txs_data, sequencer FROM statev2.forced_batch WHERE forced_batch_num = $1"
	addBlockSQL                            = "INSERT INTO statev2.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)"
	resetSQL                               = "DELETE FROM statev2.block WHERE block_num > $1"
	addVerifiedBatchSQL                    = "INSERT INTO statev2.verified_batch (block_num, batch_num, tx_hash, aggregator) VALUES ($1, $2, $3, $4)"
	getVerifiedBatchSQL                    = "SELECT block_num, batch_num, tx_hash, aggregator FROM statev2.verified_batch WHERE batch_num = $1"
	getLastBatchSQL                        = "SELECT batch_num, global_exit_root, timestamp from statev2.batch ORDER BY batch_num DESC LIMIT 1"
	getLastBatchTimeSQL                    = "SELECT timestamp FROM statev2.batch ORDER BY batch_num DESC LIMIT 1"
	getLastVirtualBatchNumSQL              = "SELECT batch_num FROM statev2.virtual_batch ORDER BY batch_num DESC LIMIT 1"
	getLastVirtualBatchBlockNumSQL         = "SELECT block_num FROM statev2.virtual_batch ORDER BY batch_num DESC LIMIT 1"
	getLastBlockNumSQL                     = "SELECT block_num FROM statev2.block ORDER BY block_num DESC LIMIT 1"
	getBlockTimeByNumSQL                   = "SELECT received_at FROM statev2.block WHERE block_num = $1"
	getBatchByNumberSQL                    = "SELECT batch_num, global_exit_root, timestamp from statev2.batch WHERE batch_num = $1"
	getEncodedTransactionsByBatchNumberSQL = "SELECT encoded from statev2.transaction where batch_num = $1"
	getLastBatchSeenSQL                    = "SELECT last_batch_num_seen FROM statev2.sync_info LIMIT 1"
	updateLastBatchSeenSQL                 = "UPDATE statev2.sync_info SET last_batch_num_seen = $1"
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

// getQuerier determines which queryer to use, dbTx or the main pgxpool.
func (s *PostgresStorage) getQuerier(dbTx pgx.Tx) querier {
	if dbTx != nil {
		return dbTx
	}
	return s
}

// getRowQuerier determines which rowQuerier to use, dbTx or the main pgxpool.
func (s *PostgresStorage) getRowQuerier(dbTx pgx.Tx) rowQuerier {
	if dbTx != nil {
		return dbTx
	}
	return s
}

// Reset resets the state to a block
func (s *PostgresStorage) Reset(ctx context.Context, block *Block, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, resetSQL, block.BlockNumber); err != nil {
		return err
	}
	//Remove consolidations
	//TODO
	return nil
}

// AddBlock adds a new block to the State Store
func (s *PostgresStorage) AddBlock(ctx context.Context, block *Block, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.String(), block.ParentHash.String(), block.ReceivedAt)
	return err
}

// AddGlobalExitRoot adds a new ExitRoot to the db
func (s *PostgresStorage) AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, addGlobalExitRootSQL, exitRoot.BlockNumber, exitRoot.GlobalExitRootNum.String(), exitRoot.MainnetExitRoot, exitRoot.RollupExitRoot, exitRoot.GlobalExitRoot)
	return err
}

// GetLatestExitRoot get the latest ExitRoot synced.
func (s *PostgresStorage) GetLatestGlobalExitRoot(ctx context.Context, tx pgx.Tx) (*GlobalExitRoot, error) {
	var (
		exitRoot  GlobalExitRoot
		globalNum uint64
		err       error
	)
	if tx == nil {
		err = s.QueryRow(ctx, getLatestExitRootSQL).Scan(&exitRoot.BlockNumber, &globalNum, &exitRoot.MainnetExitRoot, &exitRoot.RollupExitRoot, &exitRoot.GlobalExitRoot)
	} else {
		err = tx.QueryRow(ctx, getLatestExitRootSQL).Scan(&exitRoot.BlockNumber, &globalNum, &exitRoot.MainnetExitRoot, &exitRoot.RollupExitRoot, &exitRoot.GlobalExitRoot)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	exitRoot.GlobalExitRootNum = new(big.Int).SetUint64(globalNum)
	return &exitRoot, nil
}

// GetNumberOfBlocksSinceLastGERUpdate gets number of blocks since last global exit root update
func (s *PostgresStorage) GetNumberOfBlocksSinceLastGERUpdate(ctx context.Context) (uint64, error) {
	var (
		lastBlockNum         uint64
		lastExitRootBlockNum uint64
		err                  error
	)
	err = s.QueryRow(ctx, getLastBlockNumSQL).Scan(&lastBlockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	err = s.QueryRow(ctx, getLatestExitRootBlockNumSQL).Scan(&lastExitRootBlockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return lastBlockNum - lastExitRootBlockNum, nil
}

func (s *PostgresStorage) GetLastSendSequenceTime(ctx context.Context) (time.Time, error) {
	var (
		blockNum  uint64
		timestamp time.Time
	)
	err := s.QueryRow(ctx, getLastVirtualBatchBlockNumSQL).Scan(&blockNum)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrNotFound
	} else if err != nil {
		return time.Time{}, err
	}

	err = s.QueryRow(ctx, getBlockTimeByNumSQL, blockNum).Scan(&timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrNotFound
	} else if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

// AddForcedBatch adds a new ForcedBatch to the db
func (s *PostgresStorage) AddForcedBatch(ctx context.Context, forcedBatch *ForcedBatch, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, addForcedBatchSQL, forcedBatch.BlockNumber, forcedBatch.ForcedBatchNumber, forcedBatch.GlobalExitRoot.String(), forcedBatch.ForcedAt, forcedBatch.RawTxsData, forcedBatch.Sequencer.String(), forcedBatch.BatchNumber)
	return err
}

// GetForcedBatch get an L1 forcedBatch.
func (s *PostgresStorage) GetForcedBatch(ctx context.Context, tx pgx.Tx, forcedBatchNumber uint64) (*ForcedBatch, error) {
	var (
		forcedBatch    ForcedBatch
		globalExitRoot string
		rawTxs         string
		seq            string
	)
	err := tx.QueryRow(ctx, getForcedBatchSQL, forcedBatchNumber).Scan(&forcedBatch.BlockNumber, &forcedBatch.BatchNumber, &forcedBatch.ForcedBatchNumber, &globalExitRoot, &forcedBatch.ForcedAt, &rawTxs, &seq)
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
func (s *PostgresStorage) AddVerifiedBatch(ctx context.Context, verifiedBatch *VerifiedBatch, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, addVerifiedBatchSQL, verifiedBatch.BlockNumber, verifiedBatch.BatchNumber, verifiedBatch.TxHash.String(), verifiedBatch.Aggregator.String())
	return err
}

// GetVerifiedBatch get an L1 verifiedBatch.
func (s *PostgresStorage) GetVerifiedBatch(ctx context.Context, tx pgx.Tx, batchNumber uint64) (*VerifiedBatch, error) {
	var (
		verifiedBatch VerifiedBatch
		txHash        string
		agg           string
	)
	err := tx.QueryRow(ctx, getVerifiedBatchSQL, batchNumber).Scan(&verifiedBatch.BlockNumber, &verifiedBatch.BatchNumber, &txHash, &agg)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	verifiedBatch.Aggregator = common.HexToAddress(agg)
	verifiedBatch.TxHash = common.HexToHash(txHash)
	return &verifiedBatch, nil
}

func (s *PostgresStorage) GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*Batch, error) {
	var (
		batch  Batch
		gerStr string
	)
	q := s.getRowQuerier(dbTx)

	err := q.QueryRow(ctx, getLastBatchSQL).Scan(&batch.BatchNumber, &gerStr, &batch.Timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	batch.GlobalExitRootNum = new(big.Int).SetBytes(common.FromHex(gerStr))
	return &batch, nil
}

// GetLastBatchTime gets last trusted batch time
func (s *PostgresStorage) GetLastBatchTime(ctx context.Context) (time.Time, error) {
	var timestamp time.Time
	err := s.QueryRow(ctx, getLastBatchTimeSQL).Scan(&timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrStateNotSynchronized
	} else if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

// GetLastVirtualBatchNum gets last virtual batch num
func (s *PostgresStorage) GetLastVirtualBatchNum(ctx context.Context) (uint64, error) {
	var batchNum uint64
	err := s.QueryRow(ctx, getLastVirtualBatchNumSQL).Scan(&batchNum)

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
func (s *PostgresStorage) SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64) error {
	_, err := s.Exec(ctx, updateLastBatchSeenSQL, batchNumber)
	return err
}

// GetLastBatchNumberSeenOnEthereum returns the last batch number stored
// in the state that represents the last batch number that affected the
// roll-up in the Ethereum network.
func (s *PostgresStorage) GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error) {
	var batchNumber uint64
	err := s.QueryRow(ctx, getLastBatchSeenSQL).Scan(&batchNumber)

	if err != nil {
		return 0, err
	}

	return batchNumber, nil
}

func (s *PostgresStorage) GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*Batch, error) {
	var (
		batch  Batch
		gerStr string
	)
	q := s.getRowQuerier(dbTx)

	err := q.QueryRow(ctx, getBatchByNumberSQL, batchNumber).Scan(&batch.BatchNumber, &gerStr, &batch.Timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	batch.GlobalExitRootNum = new(big.Int).SetBytes(common.FromHex(gerStr))
	return &batch, nil
}

func (s *PostgresStorage) GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (encoded []string, err error) {
	q := s.getQuerier(dbTx)

	rows, err := q.Query(ctx, getEncodedTransactionsByBatchNumberSQL, batchNumber)
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
