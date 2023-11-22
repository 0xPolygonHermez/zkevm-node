package pgstatestorage

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// AddL1InfoRootToExitRoot adds a new entry in ExitRoot and returns index of L1InfoTree and error
func (p *PostgresStorage) AddL1InfoRootToExitRoot(ctx context.Context, exitRoot *state.L1InfoTreeExitRootStorageEntry, dbTx pgx.Tx) (uint64, error) {
	const addGlobalExitRootSQL = `INSERT INTO state.exit_root
				(block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index)
		VALUES ($1,         $2,        $3,                $4,               $5,               $6,              $7,
			COALESCE((SELECT MAX(l1_info_tree_index) + 1 FROM state.exit_root), 1))
		RETURNING exit_root.l1_info_tree_index;`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, addGlobalExitRootSQL,
		exitRoot.BlockNumber, exitRoot.Timestamp, exitRoot.MainnetExitRoot, exitRoot.RollupExitRoot,
		exitRoot.GlobalExitRoot.GlobalExitRoot, exitRoot.PreviousBlockHash, exitRoot.L1InfoTreeRoot)
	var l1InfoTreeIndex uint64
	err := row.Scan(&l1InfoTreeIndex)
	return l1InfoTreeIndex, err
}

func (p *PostgresStorage) GetAllL1InfoRootEntries(ctx context.Context, dbTx pgx.Tx) ([]state.L1InfoTreeExitRootStorageEntry, error) {
	const getL1InfoRootSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root
		FROM state.exit_root ORDER BY l1_info_tree_index`

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getL1InfoRootSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []state.L1InfoTreeExitRootStorageEntry
	for rows.Next() {
		var entry state.L1InfoTreeExitRootStorageEntry
		err := rows.Scan(&entry.BlockNumber, &entry.Timestamp, &entry.MainnetExitRoot, &entry.RollupExitRoot, &entry.GlobalExitRoot, &entry.PreviousBlockHash, &entry.L1InfoTreeRoot)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
