package tree

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testVectorKey struct {
	LeafType    LeafType `json:"leafType"`
	EthAddr     string   `json:"ethAddr"`
	Arity       uint8    `json:"arity"`
	ExpectedKey string   `json:"expectedKey"`
}

func init() {
	// Change dir to project root
	// This is important because we have relative paths to files containing test vectors
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestMerkleTreeKey(t *testing.T) {
	data, err := os.ReadFile("test/vectors/smt/smt-key.json")
	require.NoError(t, err)

	var testVectors []testVectorKey
	err = json.Unmarshal(data, &testVectors)
	require.NoError(t, err)

	for ti, testVector := range testVectors {
		t.Run(fmt.Sprintf("Test vector %d", ti), func(t *testing.T) {
			key, err := GetKey(testVector.LeafType, common.HexToAddress(testVector.EthAddr), nil, testVector.Arity, nil)
			require.NoError(t, err)
			expected, _ := new(big.Int).SetString(testVector.ExpectedKey, 10)
			assert.Equal(t, hex.EncodeToString(expected.Bytes()), hex.EncodeToString(key))
		})
	}
}
