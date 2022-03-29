package tree

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	poseidon "github.com/iden3/go-iden3-crypto/goldenposeidon"
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
	data, err := os.ReadFile("test/vectors/src/merkle-tree/smt-raw.json")
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

	stores := map[string]Store{
		"postgres": NewPostgresStore(mtDb),
		"memory":   NewMemStore(),
	}

	ctx := context.Background()

	for storeName, store := range stores {
		for ti, testVector := range testVectors {
			t.Run(fmt.Sprintf("Test vector %d on %s store", ti, storeName), func(t *testing.T) {
				root := scalarToh4(new(big.Int))
				mt := NewMerkleTree(store, testVector.Arity)
				log.Debugf("expectedRoot: %v", testVector.ExpectedRoot)
				for i := 0; i < len(testVector.Keys); i++ {
					k, ok := new(big.Int).SetString(testVector.Keys[i], 10)
					require.True(t, ok)
					v, ok := new(big.Int).SetString(testVector.Values[i], 10)
					require.True(t, ok)
					vH8 := scalar2fea(v)

					updateProof, err := mt.Set(ctx, root, scalarToh4(k), vH8)
					require.NoError(t, err)

					root = updateProof.NewRoot
				}
				assert.Equal(t, testVector.ExpectedRoot, h4ToString(root))
			})
		}
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

	root := new(big.Int)
	store := NewPostgresStore(mtDb)
	mt := NewMerkleTree(store, DefaultMerkleTreeArity)

	k1, success := new(big.Int).SetString("03ae74d1bbdff41d14f155ec79bb389db716160c1766a49ee9c9707407f80a11", 16)
	require.True(t, success)

	v1, success := new(big.Int).SetString("200000000000000000000", 10)
	require.True(t, success)

	k1h4 := scalarToh4(k1)
	v1H8 := scalar2fea(v1)

	updateProof, err := mt.Set(ctx, scalarToh4(root), k1h4, v1H8)
	require.NoError(t, err)
	root = h4ToScalar(updateProof.NewRoot)

	v1Proof, err := mt.Get(ctx, scalarToh4(root), k1h4)
	require.NoError(t, err)

	assert.Equal(t, v1, fea2scalar(v1Proof.Value))

	k2, success := new(big.Int).SetString("0540ae2a259cb9179561cffe6a0a3852a2c1806ad894ed396a2ef16e1f10e9c7", 16)
	require.True(t, success)

	v2, success := new(big.Int).SetString("100000000000000000000", 10)
	require.True(t, success)

	k2h4 := scalarToh4(k2)
	v2H8 := scalar2fea(v2)

	updateProof, err = mt.Set(ctx, scalarToh4(root), k2h4, v2H8)
	require.NoError(t, err)
	root = h4ToScalar(updateProof.NewRoot)

	v2Proof, err := mt.Get(ctx, scalarToh4(root), k2h4)
	require.NoError(t, err)

	assert.Equal(t, v2, fea2scalar(v2Proof.Value))

	v1ProofNew, err := mt.Get(ctx, scalarToh4(root), k1h4)
	require.NoError(t, err)

	assert.Equal(t, v1, fea2scalar(v1ProofNew.Value))
}

func TestHashBytecode(t *testing.T) {
	data, err := os.ReadFile("test/vectors/src/merkle-tree/smt-hash-bytecode.json")
	require.NoError(t, err)

	var testVectors []struct {
		Bytecode     string
		ExpectedHash string
	}
	err = json.Unmarshal(data, &testVectors)
	require.NoError(t, err)

	dbCfg := dbutils.NewConfigFromEnv()

	err = dbutils.InitOrReset(dbCfg)
	require.NoError(t, err)

	mtDb, err := db.NewSQLDB(dbCfg)
	require.NoError(t, err)

	defer mtDb.Close()

	store := NewPostgresStore(mtDb)
	mt := NewMerkleTree(store, DefaultMerkleTreeArity)

	for i, testVector := range testVectors {
		testVector := testVector
		t.Run(fmt.Sprintf("Test vector %d", i), func(t *testing.T) {
			inputBytes, err := hex.DecodeString(testVector.Bytecode)
			require.NoError(t, err)

			actual, err := mt.scHashFunction(inputBytes)
			require.NoError(t, err)

			if h4ToString(actual[:]) != testVector.ExpectedHash {
				t.Errorf("Hash bytecode failed, want %q, got %q", testVector.ExpectedHash, h4ToString(actual[:]))
			}
		})
	}
}

func merkleTreeAddN(b *testing.B, store Store, n int, hashFunction HashFunction) {
	//b.ResetTimer()

	mt := NewMerkleTree(store, DefaultMerkleTreeArity)

	ctx := context.Background()
	root := scalarToh4(new(big.Int))

	for j := 0; j < n; j++ {
		key := new(big.Int).SetUint64(uint64(j))
		value := new(big.Int).SetUint64(uint64(j))
		keyH4 := scalarToh4(key)
		valueH8 := scalar2fea(value)

		proof, err := mt.Set(ctx, root, keyH4, valueH8)
		require.NoError(b, err)

		root = proof.NewRoot
	}
}

type benchStore interface {
	Store
	Reset() error
}

func BenchmarkMerkleTreeAdd(b *testing.B) {
	nLeaves := []int{
		10,
		//100,
		//1000,
	}
	hashFunctions := map[string]HashFunction{
		"poseidon": poseidon.Hash,
		"sha256":   sha256Hash,
	}

	log.Init(log.Config{
		Level:   "error",
		Outputs: []string{"stdout"},
	})

	dbCfg := dbutils.NewConfigFromEnv()

	err := dbutils.InitOrReset(dbCfg)
	require.NoError(b, err)

	mtDb, err := db.NewSQLDB(dbCfg)
	require.NoError(b, err)
	defer mtDb.Close()

	cache, err := NewStoreCache()
	require.NoError(b, err)

	dir, err := ioutil.TempDir("", "badgerRistretoDB")
	require.NoError(b, err)
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			log.Errorf("Could not remove temporary dir %q: %v", dir, err)
		}
	}()
	badgerDb, err := NewBadgerDB(dir)
	require.NoError(b, err)

	stores := map[string]benchStore{
		"postgres":        NewPostgresStore(mtDb),
		"memory":          NewMemStore(),
		"pgRistretto":     NewPgRistrettoStore(mtDb, cache),
		"badgerRistretto": NewBadgerRistrettoStore(badgerDb, cache),
	}

	for _, n := range nLeaves {
		for hName, h := range hashFunctions {
			for storeName, store := range stores {
				b.Run(fmt.Sprintf("n=%d,store=%s,hash=%s", n, storeName, hName), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						require.NoError(b, store.Reset())
						merkleTreeAddN(b, store, n, h)
					}
				})
			}
		}
	}
}

func sha256Hash(inp [poseidon.NROUNDSF]uint64, cap [poseidon.CAPLEN]uint64) ([poseidon.CAPLEN]uint64, error) {
	return [poseidon.CAPLEN]uint64{}, nil
}

const (
	maxBenchmarkItems = 50
)

func toKey(i int) []byte {
	return []byte(fmt.Sprintf("item:%d", i))
}

func BenchmarkMerkleTreeGet(b *testing.B) {
	dbCfg := dbutils.NewConfigFromEnv()
	err := dbutils.InitOrReset(dbCfg)
	require.NoError(b, err)

	stateDb, err := db.NewSQLDB(dbCfg)
	require.NoError(b, err)

	cache, err := NewStoreCache()
	require.NoError(b, err)

	dir, err := ioutil.TempDir("", "badgerRistretoDB")
	require.NoError(b, err)
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			log.Errorf("Could not remove temporary dir %q: %v", dir, err)
		}
	}()
	badgerDb, err := NewBadgerDB(dir)
	require.NoError(b, err)

	stores := map[string]benchStore{
		"pg":              NewPostgresStore(stateDb),
		"pgRistretto":     NewPgRistrettoStore(stateDb, cache),
		"badgerRistretto": NewBadgerRistrettoStore(badgerDb, cache),
	}
	for name, store := range stores {
		require.NoError(b, store.Reset())
		b.Run(fmt.Sprintf("store=%s", name), func(b *testing.B) {
			ctx := context.Background()
			for i := 0; i < maxBenchmarkItems; i++ {
				require.NoError(b, store.Set(ctx, toKey(i), toKey(i)))
			}
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				for i := 0; i < maxBenchmarkItems; i++ {
					_, err := store.Get(ctx, toKey(i))
					require.NoError(b, err)
				}
			}
		})
	}
}
