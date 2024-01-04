package pgstatestorage

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// AddForcedBatch adds a new ForcedBatch to the db
func (p *PostgresStorage) AddForcedBatch(ctx context.Context, forcedBatch *state.ForcedBatch, tx pgx.Tx) error {
	const addForcedBatchSQL = "INSERT INTO state.forced_batch (forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase, block_num) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := tx.Exec(ctx, addForcedBatchSQL, forcedBatch.ForcedBatchNumber, forcedBatch.GlobalExitRoot.String(), forcedBatch.ForcedAt, hex.EncodeToString(forcedBatch.RawTxsData), forcedBatch.Sequencer.String(), forcedBatch.BlockNumber)
	return err
}

// GetForcedBatch get an L1 forcedBatch.
func (p *PostgresStorage) GetForcedBatch(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (*state.ForcedBatch, error) {
	var (
		forcedBatch    state.ForcedBatch
		globalExitRoot string
		rawTxs         string
		seq            string
	)
	const getForcedBatchSQL = "SELECT forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase, block_num FROM state.forced_batch WHERE forced_batch_num = $1"
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getForcedBatchSQL, forcedBatchNumber).Scan(&forcedBatch.ForcedBatchNumber, &globalExitRoot, &forcedBatch.ForcedAt, &rawTxs, &seq, &forcedBatch.BlockNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
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

// GetForcedBatchesSince gets L1 forced batches since forcedBatchNumber
func (p *PostgresStorage) GetForcedBatchesSince(ctx context.Context, forcedBatchNumber, maxBlockNumber uint64, dbTx pgx.Tx) ([]*state.ForcedBatch, error) {
	const getForcedBatchesSQL = "SELECT forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase, block_num FROM state.forced_batch WHERE forced_batch_num > $1 AND block_num <= $2 ORDER BY forced_batch_num ASC"
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, getForcedBatchesSQL, forcedBatchNumber, maxBlockNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return []*state.ForcedBatch{}, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	forcesBatches := make([]*state.ForcedBatch, 0, len(rows.RawValues()))

	for rows.Next() {
		forcedBatch, err := scanForcedBatch(rows)
		if err != nil {
			return nil, err
		}

		forcesBatches = append(forcesBatches, &forcedBatch)
	}

	return forcesBatches, nil
}

// GetNextForcedBatches gets the next forced batches from the queue.
func (p *PostgresStorage) GetNextForcedBatches(ctx context.Context, nextForcedBatches int, dbTx pgx.Tx) ([]state.ForcedBatch, error) {
	const getNextForcedBatchesSQL = `
		SELECT forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase, block_num 
		FROM state.forced_batch
		WHERE forced_batch_num > (Select coalesce(max(forced_batch_num),0) as forced_batch_num from state.batch INNER JOIN state.virtual_batch ON state.virtual_batch.batch_num = state.batch.batch_num)
		ORDER BY forced_batch_num ASC LIMIT $1;
	`
	q := p.getExecQuerier(dbTx)
	// Get the next forced batches
	rows, err := q.Query(ctx, getNextForcedBatchesSQL, nextForcedBatches)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	batches := make([]state.ForcedBatch, 0, len(rows.RawValues()))

	for rows.Next() {
		var (
			forcedBatch    state.ForcedBatch
			globalExitRoot string
			rawTxs         string
			seq            string
		)
		err := rows.Scan(&forcedBatch.ForcedBatchNumber, &globalExitRoot, &forcedBatch.ForcedAt, &rawTxs, &seq, &forcedBatch.BlockNumber)
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

// GetLastTrustedForcedBatchNumber get last trusted forced batch number
func (p *PostgresStorage) GetLastTrustedForcedBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	const getLastTrustedForcedBatchNumberSQL = "SELECT COALESCE(MAX(forced_batch_num), 0) FROM state.batch"
	var forcedBatchNumber uint64
	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getLastTrustedForcedBatchNumberSQL).Scan(&forcedBatchNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrStateNotSynchronized
	}
	return forcedBatchNumber, err
}

// GetBatchByForcedBatchNum returns the batch with the given forced batch number.
func (p *PostgresStorage) GetBatchByForcedBatchNum(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (*state.Batch, error) {
	const getForcedBatchByNumberSQL = `
		SELECT batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, batch_resources, wip
		  FROM state.batch
		 WHERE forced_batch_num = $1`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getForcedBatchByNumberSQL, forcedBatchNumber)
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	return &batch, nil
}
