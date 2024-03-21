package state

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// L1InfoTreeExitRootStorageEntryV2Feijoa leaf of the L1InfoTreeRecurisve
type L1InfoTreeExitRootStorageEntryV2Feijoa L1InfoTreeExitRootStorageEntry

// StateL1InfoTreeV2 state for L1InfoTreeV2 Feijoa Recursive Tree
type StateL1InfoTreeV2 struct {
	storageL1InfoTreeV2 storageL1InfoTreeV2
	l1InfoTreeV2        *l1infotree.L1InfoTreeRecursive
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
	mt, err := l1infotree.NewL1InfoTreeRecursiveFromLeaves(uint8(32), leaves) //nolint:gomnd
	if err != nil {
		log.Error("error creating L1InfoTree. Error: ", err)
		return fmt.Errorf("error creating L1InfoTree. Error: %w", err)
	}
	s.l1InfoTreeV2 = mt
	return nil
}

// GetCurrentL1InfoRoot Return current L1InfoRoot
func (s *StateL1InfoTreeV2) GetCurrentL1InfoRoot(ctx context.Context, dbTx pgx.Tx) (common.Hash, error) {
	err := s.buildL1InfoTreeV2CacheIfNeed(ctx, dbTx)
	if err != nil {
		log.Error("error building L1InfoTree cache. Error: ", err)
		return ZeroHash, err
	}
	return s.l1InfoTreeV2.GetRoot(), nil
}
