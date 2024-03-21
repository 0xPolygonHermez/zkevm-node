package pgstatestorage

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

const (
	l1InfoTreeIndexFieldNameV2Feijoa = "l1_info_tree_index_feijoa"
)

// AddL1InfoRootToExitRoot adds a new entry in ExitRoot and returns index of L1InfoTree and error
func (p *PostgresStorage) AddL1InfoRootToExitRootV2Feijoa(ctx context.Context, exitRoot *state.L1InfoTreeExitRootStorageEntryV2Feijoa, dbTx pgx.Tx) error {
	exitRootOld := state.L1InfoTreeExitRootStorageEntry(*exitRoot)
	return p.addL1InfoRootToExitRootVx(ctx, &exitRootOld, dbTx, l1InfoTreeIndexFieldNameV2Feijoa)
}

func (p *PostgresStorage) GetAllL1InfoRootEntriesV2Feijoa(ctx context.Context, dbTx pgx.Tx) ([]state.L1InfoTreeExitRootStorageEntryV2Feijoa, error) {
	res, err := p.GetAllL1InfoRootEntriesVx(ctx, dbTx, l1InfoTreeIndexFieldNameV2Feijoa)
	if err != nil {
		return nil, err
	}
	var entriesV2Feijoa []state.L1InfoTreeExitRootStorageEntryV2Feijoa
	for _, entry := range res {
		entriesV2Feijoa = append(entriesV2Feijoa, state.L1InfoTreeExitRootStorageEntryV2Feijoa(entry))
	}
	return entriesV2Feijoa, nil
}

func (p *PostgresStorage) GetLatestL1InfoRootV2Feijoa(ctx context.Context, maxBlockNumber uint64, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntryV2Feijoa, error) {
	res, err := p.GetLatestL1InfoRootVx(ctx, maxBlockNumber, dbTx, l1InfoTreeIndexFieldNameV2Feijoa)
	if err != nil {
		return state.L1InfoTreeExitRootStorageEntryV2Feijoa{}, err
	}
	return state.L1InfoTreeExitRootStorageEntryV2Feijoa(res), nil
}

func (p *PostgresStorage) GetLatestIndexV2Feijoa(ctx context.Context, dbTx pgx.Tx) (uint32, error) {
	return p.GetLatestIndexVx(ctx, dbTx, l1InfoTreeIndexFieldNameV2Feijoa)
}
