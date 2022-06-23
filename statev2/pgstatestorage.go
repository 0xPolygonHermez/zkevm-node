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
	addGlobalExitRootSQL = "INSERT INTO statev2.exit_root (block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root) VALUES ($1, $2, $3, $4, $5)"
	getLatestExitRootSQL = "SELECT block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root FROM statev2.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
	addForcedBatchSQL    = "INSERT INTO statev2.forced_batch (block_num, forced_batch_num, global_exit_root, timestamp, raw_txs_data, sequencer) VALUES ($1, $2, $3, $4, $5, $6)"
	getForcedBatchSQL    = "SELECT block_num, forced_batch_num, global_exit_root, timestamp, raw_txs_data, sequencer FROM statev2.forced_batch WHERE forced_batch_num = $1"
	addBlockSQL          = "INSERT INTO statev2.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)"
	resetSQL             = "DELETE FROM statev2.block WHERE block_num > $1"
	resetTrustedStateSQL = "DELETE FROM statev2.batch WHERE batch_num > $1"
	addVerifiedBatchSQL  = "INSERT INTO statev2.verified_batch (block_num, batch_num, tx_hash, aggregator) VALUES ($1, $2, $3, $4)"
	getVerifiedBatchSQL  = "SELECT block_num, batch_num, tx_hash, aggregator FROM statev2.verified_batch WHERE batch_num = $1"
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
func (s *PostgresStorage) Reset(ctx context.Context, block *Block, dbTx pgx.Tx) error {
	if _, err := dbTx.Exec(ctx, resetSQL, block.BlockNumber); err != nil {
		return err
	}
	// TODO: Remove consolidations
	return nil
}

func (p *PostgresStorage) ResetTrustedState(ctx context.Context, batchNum uint64, dbTx pgx.Tx) error {
	if _, err := dbTx.Exec(ctx, resetTrustedStateSQL, batchNum); err != nil {
		return err
	}
	// TODO: Remove consolidations
	return nil
}

// AddBlock adds a new block to the State Store
func (p *PostgresStorage) AddBlock(ctx context.Context, block *Block, dbTx pgx.Tx) error {
	_, err := dbTx.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.String(), block.ParentHash.String(), block.ReceivedAt)
	return err
}

// AddGlobalExitRoot adds a new ExitRoot to the db
func (p *PostgresStorage) AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, dbTx pgx.Tx) error {
	_, err := dbTx.Exec(ctx, addGlobalExitRootSQL, exitRoot.BlockNumber, exitRoot.GlobalExitRootNum.String(), exitRoot.MainnetExitRoot, exitRoot.RollupExitRoot, exitRoot.GlobalExitRoot)
	return err
}

// GetLatestExitRoot get the latest ExitRoot synced.
func (p *PostgresStorage) GetLatestGlobalExitRoot(ctx context.Context, dbTx pgx.Tx) (*GlobalExitRoot, error) {
	var (
		exitRoot  GlobalExitRoot
		globalNum uint64
	)
	err := dbTx.QueryRow(ctx, getLatestExitRootSQL).Scan(&exitRoot.BlockNumber, &globalNum, &exitRoot.MainnetExitRoot, &exitRoot.RollupExitRoot, &exitRoot.GlobalExitRoot)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	exitRoot.GlobalExitRootNum = new(big.Int).SetUint64(globalNum)
	return &exitRoot, nil
}

// AddForcedBatch adds a new ForcedBatch to the db
func (p *PostgresStorage) AddForcedBatch(ctx context.Context, forcedBatch *ForcedBatch, dbTx pgx.Tx) error {
	_, err := dbTx.Exec(ctx, addForcedBatchSQL, forcedBatch.BlockNumber, forcedBatch.ForcedBatchNumber, forcedBatch.GlobalExitRoot.String(), forcedBatch.ForcedAt, hex.EncodeToString(forcedBatch.RawTxsData), forcedBatch.Sequencer.String())
	return err
}

// GetForcedBatch get an L1 forcedBatch.
func (p *PostgresStorage) GetForcedBatch(ctx context.Context, dbTx pgx.Tx, forcedBatchNumber uint64) (*ForcedBatch, error) {
	var (
		forcedBatch    ForcedBatch
		globalExitRoot string
		rawTxs         string
		seq            string
	)
	err := dbTx.QueryRow(ctx, getForcedBatchSQL, forcedBatchNumber).Scan(&forcedBatch.BlockNumber, &forcedBatch.ForcedBatchNumber, &globalExitRoot, &forcedBatch.ForcedAt, &rawTxs, &seq)
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
	_, err := dbTx.Exec(ctx, addVerifiedBatchSQL, verifiedBatch.BlockNumber, verifiedBatch.BatchNumber, verifiedBatch.TxHash.String(), verifiedBatch.Aggregator.String())
	return err
}

// GetVerifiedBatch get an L1 verifiedBatch.
func (p *PostgresStorage) GetVerifiedBatch(ctx context.Context, dbTx pgx.Tx, batchNumber uint64) (*VerifiedBatch, error) {
	var (
		verifiedBatch VerifiedBatch
		txHash        string
		agg           string
	)
	err := dbTx.QueryRow(ctx, getVerifiedBatchSQL, batchNumber).Scan(&verifiedBatch.BlockNumber, &verifiedBatch.BatchNumber, &txHash, &agg)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	verifiedBatch.Aggregator = common.HexToAddress(agg)
	verifiedBatch.TxHash = common.HexToHash(txHash)
	return &verifiedBatch, nil
}
