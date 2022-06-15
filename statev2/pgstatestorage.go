package statev2

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	addGlobalExitRootSQL = "INSERT INTO statev2.exit_root (block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root) VALUES ($1, $2, $3, $4)"
	getLatestExitRootSQL = "SELECT block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root FROM statev2.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
	addForcedBatchSQL    = "INSERT INTO statev2.forced_batch (block_num, forced_batch_num, global_exit_root, timestamp, raw_txs_data, sequencer) VALUES ($1, $2, $3, $4, $5, $6)"
	getForcedBatchSQL    = "SELECT block_num, forced_batch_num, global_exit_root, timestamp, raw_txs_data, sequencer FROM statev2.forced_batch WHERE forced_batch_num = $1"
	addBlockSQL          = "INSERT INTO statev2.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)"
	resetSQL             = "DELETE FROM statev2.block WHERE block_num > $1"
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

// Reset resets the state to a block
func (s *PostgresStorage) Reset(ctx context.Context, block *Block, txID pgx.Tx) error {
	if _, err := txID.Exec(ctx, resetSQL, block.BlockNumber); err != nil {
		return err
	}
	//Remove consolidations
	//TODO
	return nil
}

// AddBlock adds a new block to the State Store
func (s *PostgresStorage) AddBlock(ctx context.Context, block *Block, txID pgx.Tx) error {
	_, err := txID.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.String(), block.ParentHash.String(), block.ReceivedAt)
	return err
}

// AddGlobalExitRoot adds a new ExitRoot to the db
func (s *PostgresStorage) AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, txID pgx.Tx) error {
	_, err := txID.Exec(ctx, addGlobalExitRootSQL, exitRoot.BlockNumber, exitRoot.GlobalExitRootNum.String(), exitRoot.MainnetExitRoot, exitRoot.RollupExitRoot)
	return err
}

// GetLatestExitRoot get the latest ExitRoot synced.
func (s *PostgresStorage) GetLatestGlobalExitRoot(ctx context.Context, txID pgx.Tx) (*GlobalExitRoot, error) {
	var (
		exitRoot  GlobalExitRoot
		globalNum uint64
	)
	err := txID.QueryRow(ctx, getLatestExitRootSQL).Scan(&exitRoot.BlockNumber, &globalNum, &exitRoot.MainnetExitRoot, &exitRoot.RollupExitRoot)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	exitRoot.GlobalExitRootNum = new(big.Int).SetUint64(globalNum)
	return &exitRoot, nil
}

// AddForcedBatch adds a new ForcedBatch to the db
func (s *PostgresStorage) AddForcedBatch(ctx context.Context, forcedBatch *ForcedBatch, txID pgx.Tx) error {
	_, err := txID.Exec(ctx, addForcedBatchSQL, forcedBatch.BlockNumber, forcedBatch.ForcedBatchNumber, forcedBatch.GlobalExitRoot.String(), forcedBatch.ForcedAt, hex.EncodeToString(forcedBatch.RawTxsData), forcedBatch.Sequencer.String())
	return err
}

// GetForcedBatch get an L1 forcedBatch.
func (s *PostgresStorage) GetForcedBatch(ctx context.Context, txID pgx.Tx, forcedBatchNumber uint64) (*ForcedBatch, error) {
	var (
		forcedBatch    ForcedBatch
		globalExitRoot string
		rawTxs         string
		seq            string
	)
	err := txID.QueryRow(ctx, getForcedBatchSQL, forcedBatchNumber).Scan(&forcedBatch.BlockNumber, &forcedBatch.ForcedBatchNumber, &globalExitRoot, &forcedBatch.ForcedAt, &rawTxs, &seq)
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
