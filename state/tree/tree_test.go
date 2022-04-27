package tree

import (
	"context"
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
	mt := NewMerkleTree(store, DefaultMerkleTreeArity)
	scCodeStore := NewPostgresSCCodeStore(mtDb)
	tree := NewStateTree(mt, scCodeStore)

	address := common.Address{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
	}

	ctx := context.Background()
	// Balance
	bal, err := tree.GetBalance(ctx, address, nil, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), bal)

	root, _, err := tree.SetBalance(ctx, address, big.NewInt(1), nil, "")
	require.NoError(t, err)

	bal, err = tree.GetBalance(ctx, address, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1), bal)

	// Nonce
	nonce, err := tree.GetNonce(ctx, address, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), nonce)

	root, _, err = tree.SetNonce(ctx, address, big.NewInt(2), root, "")
	require.NoError(t, err)

	nonce, err = tree.GetNonce(ctx, address, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(2), nonce)

	// Code
	code, err := tree.GetCode(ctx, address, root, "")
	require.NoError(t, err)
	assert.Equal(t, []byte{}, code)

	scCode, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	root, _, err = tree.SetCode(ctx, address, scCode, root, "")
	require.NoError(t, err)

	code, err = tree.GetCode(ctx, address, root, "")
	require.NoError(t, err)
	assert.Equal(t, scCode, code)

	// Storage
	position := common.Hash{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
		0x30, 0x31,
	}
	positionBI := new(big.Int).SetBytes(position.Bytes())

	storage, err := tree.GetStorageAt(ctx, address, positionBI, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), storage)

	root, _, err = tree.SetStorageAt(ctx, address, positionBI, big.NewInt(4), root, "")
	require.NoError(t, err)

	storage, err = tree.GetStorageAt(ctx, address, positionBI, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(4), storage)

	position2 := common.Hash{
		0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x11, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x21, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
		0x31, 0x31,
	}
	position2BI := new(big.Int).SetBytes(position2.Bytes())

	storage2, err := tree.GetStorageAt(ctx, address, position2BI, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), storage2)

	root, _, err = tree.SetStorageAt(ctx, address, position2BI, big.NewInt(5), root, "")
	require.NoError(t, err)

	storage2, err = tree.GetStorageAt(ctx, address, position2BI, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(5), storage2)

	storage, err = tree.GetStorageAt(ctx, address, positionBI, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(4), storage)

	bal, err = tree.GetBalance(ctx, address, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1), bal)

	balX, _ := new(big.Int).SetString("200000000000000000000", 10)
	newRoot, _, err := tree.SetBalance(ctx, address, balX, root, "")
	require.NoError(t, err)

	bal, err = tree.GetBalance(ctx, address, newRoot, "")
	require.NoError(t, err)
	assert.Equal(t, balX, bal)

	bal, err = tree.GetBalance(ctx, address, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1), bal)

	code, err = tree.GetCode(ctx, address, root, "")
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
	data, err := os.ReadFile("test/vectors/src/merkle-tree/smt-genesis.json")
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

	ctx := context.Background()
	for ti, testVector := range testVectors {
		t.Run(fmt.Sprintf("Test vector %d", ti), func(t *testing.T) {
			var root []byte
			mt := NewMerkleTree(store, testVector.Arity)
			tree := NewStateTree(mt, scCodeStore)
			for _, addrState := range testVector.Addresses {
				// convert strings to big.Int
				addr := common.HexToAddress(addrState.Address)

				balance, success := new(big.Int).SetString(addrState.Balance, 10)
				require.True(t, success)

				nonce, success := new(big.Int).SetString(addrState.Nonce, 10)
				require.True(t, success)

				root, _, err = tree.SetBalance(ctx, addr, balance, root, "")
				require.NoError(t, err)

				root, _, err = tree.SetNonce(ctx, addr, nonce, root, "")
				require.NoError(t, err)
			}

			assert.Equal(t, testVector.ExpectedRoot, new(big.Int).SetBytes(root).String())
		})
	}
}

func TestUnsetCode(t *testing.T) {
	dbCfg := dbutils.NewConfigFromEnv()

	err := dbutils.InitOrReset(dbCfg)
	require.NoError(t, err)

	mtDb, err := db.NewSQLDB(dbCfg)
	require.NoError(t, err)

	defer mtDb.Close()

	store := NewPostgresStore(mtDb)
	mt := NewMerkleTree(store, DefaultMerkleTreeArity)
	scCodeStore := NewPostgresSCCodeStore(mtDb)
	tree := NewStateTree(mt, scCodeStore)

	address := common.Address{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
	}

	ctx := context.Background()
	// populate the tree
	bal, err := tree.GetBalance(ctx, address, nil, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), bal)

	oldRoot, _, err := tree.SetBalance(ctx, address, big.NewInt(1), nil, "")
	require.NoError(t, err)

	// set and unset code
	code, err := tree.GetCode(ctx, address, oldRoot, "")
	require.NoError(t, err)
	assert.Equal(t, []byte{}, code)

	scCode, err := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	require.NoError(t, err)

	root, _, err := tree.SetCode(ctx, address, scCode, oldRoot, "")
	require.NoError(t, err)

	code, err = tree.GetCode(ctx, address, root, "")
	require.NoError(t, err)
	assert.Equal(t, scCode, code)

	newRoot, _, err := tree.SetCode(ctx, address, []byte{}, root, "")
	require.NoError(t, err)

	code, err = tree.GetCode(ctx, address, newRoot, "")
	require.NoError(t, err)
	assert.Equal(t, []byte{}, code)
}

func TestUnsetStorageAtPosition(t *testing.T) {
	dbCfg := dbutils.NewConfigFromEnv()

	err := dbutils.InitOrReset(dbCfg)
	require.NoError(t, err)

	mtDb, err := db.NewSQLDB(dbCfg)
	require.NoError(t, err)

	defer mtDb.Close()

	store := NewPostgresStore(mtDb)
	mt := NewMerkleTree(store, DefaultMerkleTreeArity)
	scCodeStore := NewPostgresSCCodeStore(mtDb)
	tree := NewStateTree(mt, scCodeStore)

	address := common.Address{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
	}

	// Storage
	position := common.Hash{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
		0x30, 0x31,
	}
	positionBI := new(big.Int).SetBytes(position.Bytes())

	ctx := context.Background()
	storage, err := tree.GetStorageAt(ctx, address, positionBI, nil, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), storage)

	oldRoot, _, err := tree.SetStorageAt(ctx, address, positionBI, big.NewInt(4), nil, "")
	require.NoError(t, err)

	storage, err = tree.GetStorageAt(ctx, address, positionBI, oldRoot, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(4), storage)

	position2 := common.Hash{
		0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x11, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x21, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
		0x31, 0x31,
	}
	position2BI := new(big.Int).SetBytes(position2.Bytes())

	storage2, err := tree.GetStorageAt(ctx, address, position2BI, oldRoot, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), storage2)

	root, _, err := tree.SetStorageAt(ctx, address, position2BI, big.NewInt(5), oldRoot, "")
	require.NoError(t, err)

	storage2, err = tree.GetStorageAt(ctx, address, position2BI, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(5), storage2)

	storage, err = tree.GetStorageAt(ctx, address, positionBI, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(4), storage)

	_, _, err = tree.SetStorageAt(ctx, address, position2BI, big.NewInt(0), root, "")
	require.NoError(t, err)

	storage, err = tree.GetStorageAt(ctx, address, position2BI, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(5), storage)

	storage2, err = tree.GetStorageAt(ctx, address, position2BI, root, "")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(5), storage2)
}

func TestSetGetNode(t *testing.T) {
	dbCfg := dbutils.NewConfigFromEnv()

	err := dbutils.InitOrReset(dbCfg)
	require.NoError(t, err)

	mtDb, err := db.NewSQLDB(dbCfg)
	require.NoError(t, err)

	defer mtDb.Close()

	store := NewPostgresStore(mtDb)
	mt := NewMerkleTree(store, DefaultMerkleTreeArity)
	scCodeStore := NewPostgresSCCodeStore(mtDb)
	tree := NewStateTree(mt, scCodeStore)

	ctx := context.Background()

	key := big.NewInt(15)
	value := big.NewInt(10)

	require.NoError(t, tree.SetNodeData(ctx, key, value))

	actualValue, err := tree.GetNodeData(ctx, key)
	require.NoError(t, err)

	expectedValue := value
	require.Equal(t, expectedValue, actualValue)
}

func TestGetCodeHash(t *testing.T) {
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
	scCodeStore := NewPostgresSCCodeStore(mtDb)
	tree := NewStateTree(mt, scCodeStore)

	ethAddress := common.HexToAddress("0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9")
	for _, testVector := range testVectors {
		code, err := hex.DecodeString(testVector.Bytecode)
		require.NoError(t, err)
		ctx := context.Background()

		expectedHash := testVector.ExpectedHash
		root, _, err := tree.SetCode(ctx, ethAddress, code, nil, "")
		require.NoError(t, err)

		resp, err := tree.GetCodeHash(ctx, ethAddress, root, "")
		require.NoError(t, err)

		require.Equal(t, expectedHash, hex.EncodeToHex(resp), "Did not get the expected code hash")
	}
}
