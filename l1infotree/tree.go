package l1infotree

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// L1InfoTree provides methods to compute L1InfoTree
type L1InfoTree struct {
	height     uint8
	zeroHashes [][32]byte
}

// NewL1InfoTree creates new L1InfoTree.
func NewL1InfoTree(height uint8) *L1InfoTree {
	return &L1InfoTree{
		zeroHashes: generateZeroHashes(height),
		height:     height,
	}
}

func buildIntermediate(leaves [][32]byte) ([][][]byte, [][32]byte) {
	var (
		nodes  [][][]byte
		hashes [][32]byte
	)
	for i := 0; i < len(leaves); i += 2 {
		var left, right int = i, i + 1
		hash := Hash(leaves[left], leaves[right])
		nodes = append(nodes, [][]byte{hash[:], leaves[left][:], leaves[right][:]})
		hashes = append(hashes, hash)
	}
	return nodes, hashes
}

// BuildL1InfoRoot computes the root given the leaves of the tree
func (mt *L1InfoTree) BuildL1InfoRoot(leaves [][32]byte) (common.Hash, error) {
	var (
		nodes [][][][]byte
		ns    [][][]byte
	)
	if len(leaves) == 0 {
		leaves = append(leaves, mt.zeroHashes[0])
	}

	for h := uint8(0); h < mt.height; h++ {
		if len(leaves)%2 == 1 {
			leaves = append(leaves, mt.zeroHashes[h])
		}
		ns, leaves = buildIntermediate(leaves)
		nodes = append(nodes, ns)
	}
	if len(ns) != 1 {
		return common.Hash{}, fmt.Errorf("error: more than one root detected: %+v", nodes)
	}

	return common.BytesToHash(ns[0][0]), nil
}

// ComputeMerkleProof computes the merkleProof and root given the leaves of the tree
func (mt *L1InfoTree) ComputeMerkleProof(gerIndex uint, leaves [][32]byte) ([][32]byte, common.Hash, error) {
	var ns [][][]byte
	if len(leaves) == 0 {
		leaves = append(leaves, mt.zeroHashes[0])
	}
	var siblings [][32]byte
	index := gerIndex
	for h := uint8(0); h < mt.height; h++ {
		if len(leaves)%2 == 1 {
			leaves = append(leaves, mt.zeroHashes[h])
		}
		if index%2 == 1 { //If it is odd
			siblings = append(siblings, leaves[index-1])
		} else { // It is even
			if len(leaves) > 1 {
				siblings = append(siblings, leaves[index+1])
			}
		}
		var (
			nsi    [][][]byte
			hashes [][32]byte
		)
		for i := 0; i < len(leaves); i += 2 {
			var left, right int = i, i + 1
			hash := Hash(leaves[left], leaves[right])
			nsi = append(nsi, [][]byte{hash[:], leaves[left][:], leaves[right][:]})
			hashes = append(hashes, hash)
			// Find the index of the leave in the next level of the tree.
			// Divide the index by 2 to find the position in the upper level
			index = uint(float64(index) / 2) //nolint:gomnd
		}
		ns = nsi
		leaves = hashes
	}
	if len(ns) != 1 {
		return nil, common.Hash{}, fmt.Errorf("error: more than one root detected: %+v", ns)
	}

	return siblings, common.BytesToHash(ns[0][0]), nil
}
