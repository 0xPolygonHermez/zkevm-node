package state

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// GlobalExitRoot struct
type L1InfoTreeLeaf struct {
	GlobalExitRoot
	PreviousBlockHash common.Hash
}

type L1InfoTreeExitRootStorageEntry struct {
	L1InfoTreeLeaf
	L1InfoTreeRoot  common.Hash
	L1InfoTreeIndex uint64
}

// TODO: Implement Hash function
func (l *L1InfoTreeLeaf) Hash() common.Hash {
	return common.Hash{}
}

func (s *State) AddL1InfoTreeLeaf(ctx context.Context, L1InfoTreeLeaf *L1InfoTreeLeaf, dbTx pgx.Tx) (*L1InfoTreeExitRootStorageEntry, error) {
	allLeaves, err := s.GetAllL1InfoRootEntries(ctx, dbTx)
	if err != nil {
		log.Error("error getting all leaves. Error: ", err)
		return nil, err
	}
	root, err := buildL1InfoTree(allLeaves)
	if err != nil {
		log.Error("error building L1InfoTree. Error: ", err)
		return nil, err
	}
	entry := L1InfoTreeExitRootStorageEntry{
		L1InfoTreeLeaf: *L1InfoTreeLeaf,
		L1InfoTreeRoot: root,
	}
	index, err := s.AddL1InfoRootToExitRoot(ctx, &entry, dbTx)
	if err != nil {
		log.Error("error adding L1InfoRoot to ExitRoot. Error: ", err)
		return nil, err
	}
	entry.L1InfoTreeIndex = index
	return &entry, nil
}

func buildL1InfoTree(allLeaves []L1InfoTreeExitRootStorageEntry) (common.Hash, error) {
	mt := l1infotree.NewL1InfoTree(uint8(32))
	var leaves [][32]byte
	for _, leaf := range allLeaves {
		leaves = append(leaves, leaf.Hash())
	}
	root, err := mt.BuildL1InfoRoot(leaves)
	return root, err
}
