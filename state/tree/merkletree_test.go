package tree

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testVectorRaw struct {
	Arity        uint8    `json:"arity"`
	Keys         []string `json:"keys"`
	Values       []string `json:"values"`
	ExpectedRoot string   `json:"expectedRoot"`
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

func TestMerkleTreeRaw(t *testing.T) {
	data, err := os.ReadFile("test/vectors/smt/smt-raw.json")
	require.NoError(t, err)

	var testVectors []testVectorRaw
	err = json.Unmarshal(data, &testVectors)
	require.NoError(t, err)

	dbCfg := dbutils.NewConfigFromEnv()

	err = dbutils.InitOrReset(dbCfg)
	require.NoError(t, err)

	mtDb, err := db.NewSQLDB(dbCfg)
	require.NoError(t, err)

	defer mtDb.Close()

	ctx := context.Background()

	for ti, testVector := range testVectors {
		t.Run(fmt.Sprintf("Test vector %d", ti), func(t *testing.T) {
			root := big.NewInt(0)
			mt := NewMerkleTree(mtDb, testVector.Arity, nil)
			for i := 0; i < len(testVector.Keys); i++ {
				// convert strings to big.Int
				k, success := new(big.Int).SetString(testVector.Keys[i], 10)
				require.True(t, success)

				v, success := new(big.Int).SetString(testVector.Values[i], 10)
				require.True(t, success)

				updateProof, err := mt.Set(ctx, root, k, v)
				require.NoError(t, err)
				root = updateProof.NewRoot
			}
			expected, _ := new(big.Int).SetString(testVector.ExpectedRoot, 10)

			r := root.Bytes()

			assert.Equal(t, hex.EncodeToString(expected.Bytes()), hex.EncodeToString(r))
		})
	}
}

func TestMerkleTree(t *testing.T) {
	dbCfg := dbutils.NewConfigFromEnv()

	err := dbutils.InitOrReset(dbCfg)
	require.NoError(t, err)

	mtDb, err := db.NewSQLDB(dbCfg)
	require.NoError(t, err)

	defer mtDb.Close()

	ctx := context.Background()

	root := big.NewInt(0)
	mt := NewMerkleTree(mtDb, 4, nil)

	k1, success := new(big.Int).SetString("03ae74d1bbdff41d14f155ec79bb389db716160c1766a49ee9c9707407f80a11", 16)
	require.True(t, success)

	v1, success := new(big.Int).SetString("200000000000000000000", 10)
	require.True(t, success)

	updateProof, err := mt.Set(ctx, root, k1, v1)
	require.NoError(t, err)
	root = updateProof.NewRoot

	v1Proof, err := mt.Get(ctx, root, k1)
	require.NoError(t, err)

	assert.Equal(t, v1, v1Proof.Value)

	k2, success := new(big.Int).SetString("0540ae2a259cb9179561cffe6a0a3852a2c1806ad894ed396a2ef16e1f10e9c7", 16)
	require.True(t, success)

	v2, success := new(big.Int).SetString("100000000000000000000", 10)
	require.True(t, success)

	updateProof, err = mt.Set(ctx, root, k2, v2)
	require.NoError(t, err)
	root = updateProof.NewRoot

	v2Proof, err := mt.Get(ctx, root, k2)
	require.NoError(t, err)

	assert.Equal(t, v2, v2Proof.Value)

	v1ProofNew, err := mt.Get(ctx, root, k1)
	require.NoError(t, err)

	assert.Equal(t, v1, v1ProofNew.Value)
}
