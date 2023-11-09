package pgstatestorage

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// AddForkID adds a new forkID to the storage
func (p *PostgresStorage) AddForkID(ctx context.Context, forkID state.ForkIDInterval, dbTx pgx.Tx) error {
	const addForkIDSQL = "INSERT INTO state.fork_id (from_batch_num, to_batch_num, fork_id, version, block_num) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (fork_id) DO UPDATE SET block_num = $5 WHERE state.fork_id.fork_id = $3;"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addForkIDSQL, forkID.FromBatchNumber, forkID.ToBatchNumber, forkID.ForkId, forkID.Version, forkID.BlockNumber)
	return err
}

// GetForkIDs get all the forkIDs stored
func (p *PostgresStorage) GetForkIDs(ctx context.Context, dbTx pgx.Tx) ([]state.ForkIDInterval, error) {
	const getForkIDsSQL = "SELECT from_batch_num, to_batch_num, fork_id, version, block_num FROM state.fork_id ORDER BY from_batch_num ASC"
	q := p.getExecQuerier(dbTx)

	rows, err := q.Query(ctx, getForkIDsSQL)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	forkIDs := make([]state.ForkIDInterval, 0, len(rows.RawValues()))

	for rows.Next() {
		var forkID state.ForkIDInterval
		if err := rows.Scan(
			&forkID.FromBatchNumber,
			&forkID.ToBatchNumber,
			&forkID.ForkId,
			&forkID.Version,
			&forkID.BlockNumber,
		); err != nil {
			return forkIDs, err
		}
		forkIDs = append(forkIDs, forkID)
	}
	return forkIDs, err
}

// UpdateForkID updates the forkID stored in db
func (p *PostgresStorage) UpdateForkID(ctx context.Context, forkID state.ForkIDInterval, dbTx pgx.Tx) error {
	const updateForkIDSQL = "UPDATE state.fork_id SET to_batch_num = $1 WHERE fork_id = $2"
	e := p.getExecQuerier(dbTx)
	if _, err := e.Exec(ctx, updateForkIDSQL, forkID.ToBatchNumber, forkID.ForkId); err != nil {
		return err
	}
	return nil
}
