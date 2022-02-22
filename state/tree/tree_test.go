package tree

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestBasicTree(t *testing.T) {
	dbCfg := dbutils.NewConfigFromEnv()

	err := dbutils.InitOrReset(dbCfg)
	require.NoError(t, err)

	mtDb, err := db.NewSQLDB(dbCfg)
	require.NoError(t, err)

	defer mtDb.Close()

	store := NewPostgresStore(mtDb)
	mt := NewMerkleTree(store, DefaultMerkleTreeArity, nil)
	scCodeStore := NewPostgresSCCodeStore(mtDb)
	tree := NewStateTree(mt, scCodeStore, nil)

	address := common.Address{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
	}

	// Balance
	bal, err := tree.GetBalance(address, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), bal)

	_, _, err = tree.SetBalance(address, big.NewInt(1))
	require.NoError(t, err)

	bal, err = tree.GetBalance(address, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1), bal)

	// Nonce
	nonce, err := tree.GetNonce(address, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), nonce)

	_, _, err = tree.SetNonce(address, big.NewInt(2))
	require.NoError(t, err)

	nonce, err = tree.GetNonce(address, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(2), nonce)

	// Code
	code, err := tree.GetCode(address, nil)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, code)

	scCode, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	_, _, err = tree.SetCode(address, scCode)
	require.NoError(t, err)

	code, err = tree.GetCode(address, nil)
	require.NoError(t, err)
	assert.Equal(t, scCode, code)

	// Storage
	position := common.Hash{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
		0x30, 0x31,
	}
	storage, err := tree.GetStorageAt(address, position, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), storage)

	_, _, err = tree.SetStorageAt(address, position, big.NewInt(4))
	require.NoError(t, err)

	storage, err = tree.GetStorageAt(address, position, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(4), storage)

	position2 := common.Hash{
		0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x11, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x21, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
		0x31, 0x31,
	}

	storage2, err := tree.GetStorageAt(address, position2, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), storage2)

	root, _, err := tree.SetStorageAt(address, position2, big.NewInt(5))
	require.NoError(t, err)

	storage2, err = tree.GetStorageAt(address, position2, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(5), storage2)

	storage, err = tree.GetStorageAt(address, position, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(4), storage)

	bal, err = tree.GetBalance(address, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1), bal)

	balX, _ := new(big.Int).SetString("200000000000000000000", 10)
	newRoot, _, err := tree.SetBalance(address, balX)
	require.NoError(t, err)

	bal, err = tree.GetBalance(address, newRoot)
	require.NoError(t, err)
	assert.Equal(t, balX, bal)

	bal, err = tree.GetBalance(address, root)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1), bal)

	code, err = tree.GetCode(address, root)
	require.NoError(t, err)
	assert.Equal(t, scCode, code)
}

type testAddState struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
	Nonce   string `json:"nonce"`
}

type testVectorGenesis struct {
	Arity        uint8          `json:"arity"`
	Addresses    []testAddState `json:"addresses"`
	ExpectedRoot string         `json:"expectedRoot"`
}

func TestMerkleTreeGenesis(t *testing.T) {
	data, err := os.ReadFile("test/vectors/smt/smt-genesis.json")
	require.NoError(t, err)

	var testVectors []testVectorGenesis
	err = json.Unmarshal(data, &testVectors)
	require.NoError(t, err)

	dbCfg := dbutils.NewConfigFromEnv()

	err = dbutils.InitOrReset(dbCfg)
	require.NoError(t, err)

	mtDb, err := db.NewSQLDB(dbCfg)
	require.NoError(t, err)

	defer mtDb.Close()
	store := NewPostgresStore(mtDb)
	scCodeStore := NewPostgresSCCodeStore(mtDb)

	for ti, testVector := range testVectors {
		t.Run(fmt.Sprintf("Test vector %d", ti), func(t *testing.T) {
			var root []byte
			var newRoot []byte
			mt := NewMerkleTree(store, testVector.Arity, nil)
			tree := NewStateTree(mt, scCodeStore, root)
			for _, addrState := range testVector.Addresses {
				// convert strings to big.Int
				addr := common.HexToAddress(addrState.Address)

				balance, success := new(big.Int).SetString(addrState.Balance, 10)
				require.True(t, success)

				nonce, success := new(big.Int).SetString(addrState.Nonce, 10)
				require.True(t, success)

				_, _, err = tree.SetBalance(addr, balance)
				require.NoError(t, err)

				newRoot, _, err = tree.SetNonce(addr, nonce)
				require.NoError(t, err)
			}

			assert.Equal(t, testVector.ExpectedRoot, new(big.Int).SetBytes(newRoot).String())
		})
	}
}
