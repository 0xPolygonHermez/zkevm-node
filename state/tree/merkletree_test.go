package tree

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"path"
	"runtime"
	"testing"
)

type testVectorKey struct {
	LeafType    LeafType         `json:"leafType"`
	EthAddr     []common.Address `json:"ethAddr"`
	Arity       uint8            `json:"arity"`
	ExpectedKey string           `json:"expectedKey"`
}

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
	data, err := os.ReadFile("state/tree/test-vector-data/smt-raw.json")
	require.NoError(t, err)

	var testVectors []testVectorRaw
	err = json.Unmarshal(data, &testVectors)
	require.NoError(t, err)

	dbCfg := dbutils.NewConfigFromEnv()

	err = db.RunMigrations(dbCfg)
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
