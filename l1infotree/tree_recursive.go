package l1infotree

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	firstLeafHistoricL1InfoTree = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

type L1InfoTreeRecursive struct {
	historicL1InfoTree *L1InfoTree
	l1InfoTreeDataHash *common.Hash
	leaves             [][32]byte
}

func NewL1InfoTreeRecursive(height uint8) (*L1InfoTreeRecursive, error) {
	historic, err := NewL1InfoTree(height, nil)
	if err != nil {
		return nil, err
	}
	// Insert first leaf, all zeros (no changes in tree, just to skip leaf with index 0)
	//historic.AddLeaf(0, common.HexToHash(firstLeafHistoricL1InfoTree))

	return &L1InfoTreeRecursive{
		historicL1InfoTree: historic,
	}, nil
}

func (mt *L1InfoTreeRecursive) AddLeaf(index uint32, leaf [32]byte) (common.Hash, error) {

	previousRoot := mt.GetRoot()
	_, err := mt.historicL1InfoTree.AddLeaf(index, previousRoot)
	if err != nil {
		return common.Hash{}, err
	}
	mt.leaves = append(mt.leaves, leaf)
	leafHash := common.Hash(leaf)
	mt.l1InfoTreeDataHash = &leafHash
	return mt.GetRoot(), nil

}

func (mt *L1InfoTreeRecursive) GetRoot() common.Hash {
	if mt.l1InfoTreeDataHash == nil {
		return common.HexToHash(firstLeafHistoricL1InfoTree)
	}
	return crypto.Keccak256Hash(mt.historicL1InfoTree.GetRoot().Bytes(), mt.l1InfoTreeDataHash.Bytes())

}

func (mt *L1InfoTreeRecursive) ComputeMerkleProofFromLeaves(gerIndex uint32, leaves [][32]byte) ([][32]byte, common.Hash, error) {
	return mt.historicL1InfoTree.ComputeMerkleProof(gerIndex, leaves)
}

func (mt *L1InfoTreeRecursive) ComputeMerkleProof(gerIndex uint32) ([][32]byte, common.Hash, error) {
	return mt.historicL1InfoTree.ComputeMerkleProof(gerIndex, mt.leaves)
}
