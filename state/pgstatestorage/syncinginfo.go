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
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getSyncingInfoSQL).
		Scan(&info.InitialSyncingBlock, &info.LastBlockNumberSeen, &info.LastBlockNumberConsolidated,
			&info.InitialSyncingBatch, &info.LastBatchNumberSeen, &info.LastBatchNumberConsolidated)
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

	return info, err
}
