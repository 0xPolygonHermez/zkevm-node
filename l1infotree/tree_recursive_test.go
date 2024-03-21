package l1infotree_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	L1InfoRootRecursiveHeight = uint8(32)
	EmptyL1InfoRootRecursive  = "0x0000000000000000000000000000000000000000000000000000000000000000"

	root1            = "0xcc4105312818e9b7f692c9c807ea09699f4f290e5e31671a8e0c2c937f1c43f0"
	filenameTestData = "../test/vectors/src/merkle-tree/l1-info-tree-recursive/smt-full-output.json"
)

type vectorTestData struct {
	GlobalExitRoot         common.Hash   `json:"globalExitRoot"`
	BlockHash              common.Hash   `json:"blockHash"`
	Timestamp              uint64        `json:"timestamp"`
	SmtProofPreviousIndex  []common.Hash `json:"smtProofPreviousIndex"`
	Index                  uint32        `json:"index"`
	PreviousL1InfoTreeRoot common.Hash   `json:"previousL1InfoTreeRoot"`
	L1DataHash             common.Hash   `json:"l1DataHash"`
	L1InfoTreeRoot         common.Hash   `json:"l1InfoTreeRoot"`
	HistoricL1InfoRoot     common.Hash   `json:"historicL1InfoRoot"`
}

func readData(t *testing.T) []vectorTestData {
	data, err := os.ReadFile(filenameTestData)
	require.NoError(t, err)
	var mtTestVectors []vectorTestData
	err = json.Unmarshal(data, &mtTestVectors)
	require.NoError(t, err)
	return mtTestVectors
}

func TestBuildTreeVectorData(t *testing.T) {
	data := readData(t)
	sut, err := l1infotree.NewL1InfoTreeRecursive(L1InfoRootRecursiveHeight)
	require.NoError(t, err)
	for _, testVector := range data {
		// Add leaf
		leafData := l1infotree.HashLeafData(testVector.GlobalExitRoot, testVector.BlockHash, testVector.Timestamp)
		leafDataHash := common.BytesToHash(leafData[:])
		root, err := sut.AddLeaf(testVector.Index-1, leafData)
		require.NoError(t, err)
		require.Equal(t, testVector.L1InfoTreeRoot.String(), root.String(), "Roots do not match leaf", testVector.Index)
		require.Equal(t, testVector.L1DataHash.String(), leafDataHash.String(), "leafData do not match leaf", testVector.Index)
	}
}

func TestEmptyL1InfoRootRecursive(t *testing.T) {
	// empty
	sut, err := l1infotree.NewL1InfoTreeRecursive(L1InfoRootRecursiveHeight)
	require.NoError(t, err)
	require.NotNil(t, sut)
	root := sut.GetRoot()
	require.Equal(t, EmptyL1InfoRootRecursive, root.String())
}
func TestProofsTreeVectorData(t *testing.T) {
	data := readData(t)
	sut, err := l1infotree.NewL1InfoTreeRecursive(L1InfoRootRecursiveHeight)
	require.NoError(t, err)
	for _, testVector := range data {
		// Add leaf
		leafData := l1infotree.HashLeafData(testVector.GlobalExitRoot, testVector.BlockHash, testVector.Timestamp)

		_, err := sut.AddLeaf(testVector.Index-1, leafData)
		require.NoError(t, err)
		mp, _, err := sut.ComputeMerkleProof(testVector.Index)
		require.NoError(t, err)
		for i, v := range mp {
			c := common.Hash(v)
			if c.String() != testVector.SmtProofPreviousIndex[i].String() {
				log.Info("MerkleProof: index ", testVector.Index, " mk:", i, " v:", c.String(), " expected:", testVector.SmtProofPreviousIndex[i].String())
			}
		}
	}
}
