package pgstatestorage

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// GetSyncingInfo returns information regarding the syncing status of the node
func (p *PostgresStorage) GetSyncingInfo(ctx context.Context, dbTx pgx.Tx) (state.SyncingInfo, error) {
	var info state.SyncingInfo
	const getSyncingInfoSQL = `
	select coalesce(MIN(initial_blocks.block_num), 0) as init_sync_block
	FROM state.sync_info sy
	 INNER JOIN state.l2block initial_blocks
		ON initial_blocks.batch_num = sy.init_sync_batch;
	`
	q := p.getExecQuerier(dbTx)
	row := q.QueryRow(ctx, getSyncingInfoSQL)

	err := row.Scan(&info.InitialSyncingBlock)
	if err != nil {
		return state.SyncingInfo{}, nil
	}
	const getSyncTableSQL = `
	select last_batch_num_seen, last_batch_num_consolidated, init_sync_batch from state.sync_info;
	`
	row = q.QueryRow(ctx, getSyncTableSQL)
	err = row.Scan(&info.LastBatchNumberSeen, &info.LastBatchNumberConsolidated, &info.InitialSyncingBatch)
	if err != nil {
		return state.SyncingInfo{}, nil
	}
	lastBlockNumber, err := p.GetLastL2BlockNumber(ctx, dbTx)
	if err != nil {
		return state.SyncingInfo{}, nil
	}
	info.CurrentBlockNumber = lastBlockNumber

	lastBatchNumber, err := p.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		return state.SyncingInfo{}, nil
	}
	info.CurrentBatchNumber = lastBatchNumber
	info.IsSynchronizing = info.LastBatchNumberSeen > lastBatchNumber
	return info, err
}
