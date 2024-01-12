package state

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

var (
	EmptyL1InfoRoot = common.HexToHash("0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
)

// L1InfoTreeLeaf leaf of the L1InfoTree
type L1InfoTreeLeaf struct {
	GlobalExitRoot
	PreviousBlockHash common.Hash
}

// L1InfoTreeExitRootStorageEntry entry of the Database
type L1InfoTreeExitRootStorageEntry struct {
	L1InfoTreeLeaf
	L1InfoTreeRoot  common.Hash
	L1InfoTreeIndex uint32
}

// Hash returns the hash of the leaf
func (l *L1InfoTreeLeaf) Hash() common.Hash {
	timestamp := uint64(l.Timestamp.Unix())
	return l1infotree.HashLeafData(l.GlobalExitRoot.GlobalExitRoot, l.PreviousBlockHash, timestamp)
}

// AddL1InfoTreeLeaf adds a new leaf to the L1InfoTree and returns the entry and error
func (s *State) AddL1InfoTreeLeaf(ctx context.Context, l1InfoTreeLeaf *L1InfoTreeLeaf, dbTx pgx.Tx) (*L1InfoTreeExitRootStorageEntry, error) {
	var newIndex uint32
	gerIndex, err := s.GetLatestIndex(ctx, dbTx)
	if err != nil && !errors.Is(err, ErrNotFound) {
		log.Error("error getting latest l1InfoTree index. Error: ", err)
		return nil, err
	} else if err == nil {
		newIndex = gerIndex + 1
	}
	log.Debug("latestIndex: ", gerIndex)
	root, err := s.l1InfoTree.AddLeaf(newIndex, l1InfoTreeLeaf.Hash())
	if err != nil {
		log.Error("error add new leaf to the L1InfoTree. Error: ", err)
		return nil, err
	}
	entry := L1InfoTreeExitRootStorageEntry{
		L1InfoTreeLeaf:  *l1InfoTreeLeaf,
		L1InfoTreeRoot:  root,
		L1InfoTreeIndex: newIndex,
	}
	err = s.AddL1InfoRootToExitRoot(ctx, &entry, dbTx)
	if err != nil {
		log.Error("error adding L1InfoRoot to ExitRoot. Error: ", err)
		return nil, err
	}
	return &entry, nil
}

// GetCurrentL1InfoRoot Return current L1InfoRoot
func (s *State) GetCurrentL1InfoRoot() common.Hash {
	root, _, _ := s.l1InfoTree.GetCurrentRootCountAndSiblings()
	return root
}
