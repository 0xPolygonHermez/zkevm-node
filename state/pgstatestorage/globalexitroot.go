package pgstatestorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// AddGlobalExitRoot adds a new ExitRoot to the db
func (p *PostgresStorage) AddGlobalExitRoot(ctx context.Context, exitRoot *state.GlobalExitRoot, dbTx pgx.Tx) error {
	const addGlobalExitRootSQL = "INSERT INTO state.exit_root (block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root) VALUES ($1, $2, $3, $4, $5)"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addGlobalExitRootSQL, exitRoot.BlockNumber, exitRoot.Timestamp, exitRoot.MainnetExitRoot, exitRoot.RollupExitRoot, exitRoot.GlobalExitRoot)
	return err
}

// GetLatestGlobalExitRoot get the latest global ExitRoot synced.
func (p *PostgresStorage) GetLatestGlobalExitRoot(ctx context.Context, maxBlockNumber uint64, dbTx pgx.Tx) (state.GlobalExitRoot, time.Time, error) {
	const getLatestExitRootSQL = "SELECT block_num, mainnet_exit_root, rollup_exit_root, global_exit_root FROM state.exit_root WHERE block_num <= $1 ORDER BY id DESC LIMIT 1"

	var (
		exitRoot   state.GlobalExitRoot
		err        error
		receivedAt time.Time
	)

	e := p.getExecQuerier(dbTx)
	err = e.QueryRow(ctx, getLatestExitRootSQL, maxBlockNumber).Scan(&exitRoot.BlockNumber, &exitRoot.MainnetExitRoot, &exitRoot.RollupExitRoot, &exitRoot.GlobalExitRoot)

	if errors.Is(err, pgx.ErrNoRows) {
		return state.GlobalExitRoot{}, time.Time{}, state.ErrNotFound
	} else if err != nil {
		return state.GlobalExitRoot{}, time.Time{}, err
	}

	err = e.QueryRow(ctx, getBlockTimeByNumSQL, exitRoot.BlockNumber).Scan(&receivedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return state.GlobalExitRoot{}, time.Time{}, state.ErrNotFound
	} else if err != nil {
		return state.GlobalExitRoot{}, time.Time{}, err
	}
	return exitRoot, receivedAt, nil
}

// GetNumberOfBlocksSinceLastGERUpdate gets number of blocks since last global exit root update
func (p *PostgresStorage) GetNumberOfBlocksSinceLastGERUpdate(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var (
		lastBlockNum         uint64
		lastExitRootBlockNum uint64
		err                  error
	)
	const getLatestExitRootBlockNumSQL = "SELECT block_num FROM state.exit_root ORDER BY id DESC LIMIT 1"

	e := p.getExecQuerier(dbTx)
	err = e.QueryRow(ctx, getLastBlockNumSQL).Scan(&lastBlockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	err = p.QueryRow(ctx, getLatestExitRootBlockNumSQL).Scan(&lastExitRootBlockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
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
	const getMainnetExitRoot = "SELECT block_num, mainnet_exit_root FROM state.exit_root WHERE global_exit_root = $1"

	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getMainnetExitRoot, ger.Bytes()).Scan(&blockNum, &mainnetExitRoot)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, common.Hash{}, state.ErrNotFound
	} else if err != nil {
		return 0, common.Hash{}, err
	}

	return blockNum, mainnetExitRoot, nil
}

// UpdateGERInOpenBatch update ger in open batch
func (p *PostgresStorage) UpdateGERInOpenBatch(ctx context.Context, ger common.Hash, dbTx pgx.Tx) error {
	if dbTx == nil {
		return state.ErrDBTxNil
	}

	var (
		batchNumber   uint64
		isBatchHasTxs bool
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastBatchNumberSQL).Scan(&batchNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return state.ErrStateNotSynchronized
	}

	const isBatchHasTxsQuery = `SELECT EXISTS (SELECT 1 FROM state.l2block WHERE batch_num = $1)`
	err = e.QueryRow(ctx, isBatchHasTxsQuery, batchNumber).Scan(&isBatchHasTxs)
	if err != nil {
		return err
	}

	if isBatchHasTxs {
		return errors.New("batch has txs, can't change globalExitRoot")
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

// GetLatestGer is used to get the latest ger
func (p *PostgresStorage) GetLatestGer(ctx context.Context, maxBlockNumber uint64) (state.GlobalExitRoot, time.Time, error) {
	ger, receivedAt, err := p.GetLatestGlobalExitRoot(ctx, maxBlockNumber, nil)
	if err != nil && errors.Is(err, state.ErrNotFound) {
		return state.GlobalExitRoot{}, time.Time{}, nil
	} else if err != nil {
		return state.GlobalExitRoot{}, time.Time{}, fmt.Errorf("failed to get latest global exit root, err: %w", err)
	} else {
		return ger, receivedAt, nil
	}
}

// GetExitRootByGlobalExitRoot returns the mainnet and rollup exit root given
// a global exit root number.
func (p *PostgresStorage) GetExitRootByGlobalExitRoot(ctx context.Context, ger common.Hash, dbTx pgx.Tx) (*state.GlobalExitRoot, error) {
	var (
		exitRoot state.GlobalExitRoot
		err      error
	)

	const sql = "SELECT block_num, mainnet_exit_root, rollup_exit_root, global_exit_root FROM state.exit_root WHERE global_exit_root = $1 ORDER BY id DESC LIMIT 1"

	e := p.getExecQuerier(dbTx)
	err = e.QueryRow(ctx, sql, ger).Scan(&exitRoot.BlockNumber, &exitRoot.MainnetExitRoot, &exitRoot.RollupExitRoot, &exitRoot.GlobalExitRoot)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &exitRoot, nil
}
