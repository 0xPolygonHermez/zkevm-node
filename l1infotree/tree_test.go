package l1infotree_test

import (
	"encoding/json"
	"os"
	"testing"

	packagesut "github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/test/vectors"
	"github.com/stretchr/testify/require"
)

func TestComputeTreeRoot(t *testing.T) {
	data, err := os.ReadFile("../test/vectors/src/merkle-tree/l1-info-tree/root-vectors.json")
	require.NoError(t, err)
	var mtTestVectors []vectors.L1InfoTree
	err = json.Unmarshal(data, &mtTestVectors)
	require.NoError(t, err)
	for _, testVector := range mtTestVectors {
		input := testVector.PreviousLeafValues
		mt := packagesut.NewL1InfoTree(uint8(32))
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
