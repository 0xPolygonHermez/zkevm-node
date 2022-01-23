package bridgetree

import (
	"bytes"
	"fmt"
)

// BridgeTree is Merkle Tree used in Bridge Contract for Deposits and Withdrawals
type BridgeTree struct {
	height     uint8
	dirty      bool
	tree       [][][32]byte
	zeroHashes [][32]byte
}

// NewBridgeTree creates new BridgeTree
func NewBridgeTree(height uint8) *BridgeTree {
	var zeroHashes [][32]byte
	if int(height) <= len(ZeroHashes) {
		zeroHashes = ZeroHashes
	} else {
		zeroHashes = generateZeroHashes(height)
	}
	tree := make([][][32]byte, height+1)
	for i := 0; i <= int(height); i++ {
		tree[i] = [][32]byte{}
	}
	return &BridgeTree{
		height:     height,
		dirty:      true,
		tree:       tree,
		zeroHashes: zeroHashes,
	}
}

// Add adds leaf to the tree
func (bt *BridgeTree) Add(leaf [32]byte) {
	bt.dirty = true
	bt.tree[0] = append(bt.tree[0], leaf)
}

func (bt *BridgeTree) calcBranches() {
	for i := 0; i < int(bt.height); i++ {
		child := bt.tree[i]
		for j := 0; j < len(child); j += 2 {
			leftNode := child[j]
			rightNode := bt.zeroHashes[i]
			if j+1 < len(child) {
				rightNode = child[j+1]
			}
			hashRes := hash(leftNode, rightNode)
			if j/2 < len(bt.tree[i+1]) {
				bt.tree[i+1][j/2] = hashRes
			} else if j/2 == len(bt.tree[i+1]) {
				bt.tree[i+1] = append(bt.tree[i+1], hashRes)
			} else {
				panic("we shouldn't be here")
			}
		}
	}
	bt.dirty = false
}

// GetProofTreeByIndex returns proof for the leaf with specified index in the tree
func (bt *BridgeTree) GetProofTreeByIndex(index int) [][32]byte {
	if bt.dirty {
		bt.calcBranches()
	}
	var proof [][32]byte
	currentIndex := index
	for i := 0; i < int(bt.height); i++ {
		if currentIndex%2 == 1 {
			currentIndex -= 1
		} else {
			currentIndex += 1
		}
		if currentIndex < len(bt.tree[i]) {
			proof = append(proof, bt.tree[i][currentIndex])
		} else {
			proof = append(proof, bt.zeroHashes[i])
		}
		currentIndex = currentIndex / 2
	}
	return proof
}

// GetProofTreeByValue returns proof for the leaf with specified value in the tree
func (bt *BridgeTree) GetProofTreeByValue(value [32]byte) ([][32]byte, error) {
	index := -1
	for i, val := range bt.tree[0] {
		if bytes.Compare(val[:], value[:]) == 0 {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, fmt.Errorf("value not found in the tree")
	}
	return bt.GetProofTreeByIndex(index), nil
}

// GetRoot returns merkle root of the tree
func (bt *BridgeTree) GetRoot() [32]byte {
	if bt.tree == nil || bt.tree[0] == nil || len(bt.tree[0]) == 0 {
		// No leafs in the tree, calculate root with all leafs to 0
		return bt.zeroHashes[bt.height]
	}
	if bt.dirty {
		bt.calcBranches()
	}
	return bt.tree[bt.height][0]
}
