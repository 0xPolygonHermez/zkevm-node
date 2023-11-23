package l1infotree

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// L1InfoTree holds the test vector for the merkle tree
type VectorL1InfoTree struct {
	PreviousLeafValues []common.Hash `json:"previousLeafValues"`
	CurrentRoot        common.Hash   `json:"currentRoot"`
	NewLeafValue       common.Hash   `json:"newLeafValue"`
	NewRoot            common.Hash   `json:"newRoot"`
}

func TestComputeTreeRoot(t *testing.T) {
	data, err := os.ReadFile("../test/vectors/src/merkle-tree/l1-info-tree/root-vectors.json")
	require.NoError(t, err)
	var mtTestVectors []VectorL1InfoTree
	err = json.Unmarshal(data, &mtTestVectors)
	require.NoError(t, err)
	for _, testVector := range mtTestVectors {
		input := testVector.PreviousLeafValues
		mt := NewL1InfoTree(uint8(32))
		require.NoError(t, err)

		var leaves [][32]byte
		for _, v := range input {
			leaves = append(leaves, v)
		}

		if len(leaves) != 0 {
			root, err := mt.BuildL1InfoRoot(leaves)
			require.NoError(t, err)
			require.Equal(t, testVector.CurrentRoot, root)
		}

		leaves = append(leaves, testVector.NewLeafValue)
		newRoot, err := mt.BuildL1InfoRoot(leaves)
		require.NoError(t, err)
		require.Equal(t, testVector.NewRoot, newRoot)
	}
}
