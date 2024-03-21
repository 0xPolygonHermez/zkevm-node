package state

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgx/v4"
)

type L1InfoTreeExitRootStorageEntryV2Feijoa L1InfoTreeExitRootStorageEntry

type StateL1InfoTreeV2 struct {
	storageL1InfoTreeV2 storageL1InfoTreeV2
	l1InfoTreeV2        *l1infotree.L1InfoTree
}

func (s *StateL1InfoTreeV2) buildL1InfoTreeV2CacheIfNeed(ctx context.Context, dbTx pgx.Tx) error {
	if s.l1InfoTreeV2 != nil {
		return nil
	}
	log.Debugf("Building L1InfoTree cache")
	allLeaves, err := s.storageL1InfoTreeV2.GetAllL1InfoRootEntriesV2Feijoa(ctx, dbTx)
	if err != nil {
		log.Error("error getting all leaves. Error: ", err)
		return fmt.Errorf("error getting all leaves. Error: %w", err)
	}
	var leaves [][32]byte
	for _, leaf := range allLeaves {
		leaves = append(leaves, leaf.Hash())
	}
	mt, err := l1infotree.NewL1InfoTree(uint8(32), leaves) //nolint:gomnd
	if err != nil {
		log.Error("error creating L1InfoTree. Error: ", err)
		return fmt.Errorf("error creating L1InfoTree. Error: %w", err)
	}
	s.l1InfoTreeV2 = mt
	return nil
}
