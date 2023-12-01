package state

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// L1InfoTreeLeaf leaf of the L1InfoTree
type L1InfoTreeLeaf struct {
	GlobalExitRoot
	PreviousBlockHash common.Hash
}

// L1InfoTreeIndexType type of the index of the leafs of L1InfoTree
//
//	the leaf starts at 0
type L1InfoTreeIndexType uint32

// L1InfoTreeExitRootStorageEntry entry of the Database
type L1InfoTreeExitRootStorageEntry struct {
	L1InfoTreeLeaf
	L1InfoTreeRoot  common.Hash
	L1InfoTreeIndex L1InfoTreeIndexType
}

var (
	// TODO: Put the real hash of Leaf 0, pending of deploying contracts
	leaf0Hash = [32]byte{} //nolint:gomnd
)

// Hash returns the hash of the leaf
func (l *L1InfoTreeLeaf) Hash() common.Hash {
	timestamp := uint64(l.Timestamp.Unix())
	return l1infotree.HashLeafData(l.GlobalExitRoot.GlobalExitRoot, l.PreviousBlockHash, timestamp)
}

// AddL1InfoTreeLeaf adds a new leaf to the L1InfoTree and returns the entry and error
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
	mt := l1infotree.NewL1InfoTree(uint8(32)) //nolint:gomnd
	var leaves [][32]byte
	// Insert the Leaf0 that is not used but compute for the Merkle Tree
	leaves = append(leaves, leaf0Hash)
	for _, leaf := range allLeaves {
		leaves = append(leaves, leaf.Hash())
	}
	root, err := mt.BuildL1InfoRoot(leaves)
	return root, err
}
