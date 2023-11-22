package state

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
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
	L1InfoTreeRoot common.Hash
}

// TODO: Implement Hash function
func (l *L1InfoTreeLeaf) Hash() common.Hash {
	return common.Hash{}
}

func (s *State) AddL1InfoTreeLeaf(ctx context.Context, L1InfoTreeLeaf *L1InfoTreeLeaf, dbTx pgx.Tx) error {
	allLeaves, err := s.GetAllL1InfoRootEntries(ctx, dbTx)
	if err != nil {
		return err
	}
	root, err := buildL1InfoTree(allLeaves)
	if err != nil {
		return err
	}
	entry := L1InfoTreeExitRootStorageEntry{
		L1InfoTreeLeaf: *L1InfoTreeLeaf,
		L1InfoTreeRoot: root,
	}
	return s.AddL1InfoRootToExitRoot(ctx, &entry, dbTx)
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
