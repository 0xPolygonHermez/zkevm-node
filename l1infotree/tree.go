package l1infotree

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
)

// L1InfoTree provides methods to compute L1InfoTree
type L1InfoTree struct {
	height      uint8
	zeroHashes  [][32]byte
	count       uint32
	siblings    [][32]byte
	currentRoot common.Hash
}

// NewL1InfoTree creates new L1InfoTree.
func NewL1InfoTree(height uint8, initialLeaves [][32]byte) (*L1InfoTree, error) {
	mt := &L1InfoTree{
		zeroHashes: generateZeroHashes(height),
		height:     height,
		count:      uint32(len(initialLeaves)),
	}
	var err error
	mt.siblings, mt.currentRoot, err = mt.initSiblings(initialLeaves)
	if err != nil {
		log.Error("error initializing si siblings. Error: ", err)
		return nil, err
	}
	log.Debug("Initial count: ", mt.count)
	log.Debug("Initial root: ", mt.currentRoot)
	return mt, nil
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
func (mt *L1InfoTree) ComputeMerkleProof(gerIndex uint32, leaves [][32]byte) ([][32]byte, common.Hash, error) {
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
				if index >= uint32(len(leaves)) {
					// siblings = append(siblings, mt.zeroHashes[h])
					siblings = append(siblings, leaves[index-1])
				} else {
					siblings = append(siblings, leaves[index+1])
				}
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
		}
		// Find the index of the leave in the next level of the tree.
		// Divide the index by 2 to find the position in the upper level
		index = uint32(float64(index) / 2) //nolint:gomnd
		ns = nsi
		leaves = hashes
	}
	if len(ns) != 1 {
		return nil, common.Hash{}, fmt.Errorf("error: more than one root detected: %+v", ns)
	}

	return siblings, common.BytesToHash(ns[0][0]), nil
}

// AddLeaf adds new leaves to the tree and computes the new root
func (mt *L1InfoTree) AddLeaf(index uint32, leaf [32]byte) (common.Hash, error) {
	if index != mt.count {
		return common.Hash{}, fmt.Errorf("mismatched leaf count: %d, expected: %d", index, mt.count)
	}
	cur := leaf
	isFilledSubTree := true

	for h := uint8(0); h < mt.height; h++ {
		if index&(1<<h) > 0 {
			var child [32]byte
			copy(child[:], cur[:])
			parent := Hash(mt.siblings[h], child)
			cur = parent
		} else {
			if isFilledSubTree {
				// we will update the sibling when the sub tree is complete
				copy(mt.siblings[h][:], cur[:])
				// we have a left child in this layer, it means the right child is empty so the sub tree is not completed
				isFilledSubTree = false
			}
			var child [32]byte
			copy(child[:], cur[:])
			parent := Hash(child, mt.zeroHashes[h])
			cur = parent
			// the sibling of 0 bit should be the zero hash, since we are in the last node of the tree
		}
	}
	mt.currentRoot = cur
	mt.count++
	return cur, nil
}

// initSiblings returns the siblings of the node at the given index.
// it is used to initialize the siblings array in the beginning.
func (mt *L1InfoTree) initSiblings(initialLeaves [][32]byte) ([][32]byte, common.Hash, error) {
	if mt.count != uint32(len(initialLeaves)) {
		return nil, [32]byte{}, fmt.Errorf("error: mt.count and initialLeaves length mismatch")
	}
	if mt.count == 0 {
		var siblings [][32]byte
		for h := 0; h < int(mt.height); h++ {
			var left [32]byte
			copy(left[:], mt.zeroHashes[h][:])
			siblings = append(siblings, left)
		}
		root, err := mt.BuildL1InfoRoot(initialLeaves)
		if err != nil {
			log.Error("error calculating initial root: ", err)
			return nil, [32]byte{}, err
		}
		return siblings, root, nil
	}

	return mt.ComputeMerkleProof(mt.count, initialLeaves)
}

// GetCurrentRootCountAndSiblings returns the latest root, count and sibblings
func (mt *L1InfoTree) GetCurrentRootCountAndSiblings() (common.Hash, uint32, [][32]byte) {
	return mt.currentRoot, mt.count, mt.siblings
}
