package l1infotree

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	firstLeafHistoricL1InfoTree = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

// L1InfoTreeRecursive is a recursive implementation of the L1InfoTree of Feijoa
type L1InfoTreeRecursive struct {
	historicL1InfoTree *L1InfoTree
	l1InfoTreeDataHash *common.Hash
	leaves             [][32]byte
}

// NewL1InfoTreeRecursive creates a new empty L1InfoTreeRecursive
func NewL1InfoTreeRecursive(height uint8) (*L1InfoTreeRecursive, error) {
	historic, err := NewL1InfoTree(height, nil)
	if err != nil {
		return nil, err
	}

	return &L1InfoTreeRecursive{
		historicL1InfoTree: historic,
	}, nil
}

// NewL1InfoTreeRecursiveFromLeaves creates a new L1InfoTreeRecursive from leaves
func NewL1InfoTreeRecursiveFromLeaves(height uint8, leaves [][32]byte) (*L1InfoTreeRecursive, error) {
	historic, err := NewL1InfoTree(height, nil)
	if err != nil {
		return nil, err
	}

	res := &L1InfoTreeRecursive{
		historicL1InfoTree: historic,
	}
	for _, leaf := range leaves {
		_, err := res.AddLeaf(uint32(len(res.leaves)), leaf)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// AddLeaf adds a new leaf to the L1InfoTreeRecursive
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

// GetRoot returns the root of the L1InfoTreeRecursive
func (mt *L1InfoTreeRecursive) GetRoot() common.Hash {
	if mt.l1InfoTreeDataHash == nil {
		return common.HexToHash(firstLeafHistoricL1InfoTree)
	}
	return crypto.Keccak256Hash(mt.historicL1InfoTree.GetRoot().Bytes(), mt.l1InfoTreeDataHash.Bytes())
}

// ComputeMerkleProofFromLeaves computes the Merkle proof from the leaves
func (mt *L1InfoTreeRecursive) ComputeMerkleProofFromLeaves(gerIndex uint32, leaves [][32]byte) ([][32]byte, common.Hash, error) {
	return mt.historicL1InfoTree.ComputeMerkleProof(gerIndex, leaves)
}

// ComputeMerkleProof computes the Merkle proof from the current leaves
func (mt *L1InfoTreeRecursive) ComputeMerkleProof(gerIndex uint32) ([][32]byte, common.Hash, error) {
	return mt.historicL1InfoTree.ComputeMerkleProof(gerIndex, mt.leaves)
}
