package tree_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/state/tree/pb"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/operations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	host = "0.0.0.0"
	port = 50060

	ethAddress = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
)

var (
	address = fmt.Sprintf("%s:%d", host, port)
	mtSrv   *tree.Server
	conn    *grpc.ClientConn
	cancel  context.CancelFunc
	err     error
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

func TestMain(m *testing.M) {
	initialize()
	defer teardown()

	os.Exit(m.Run())
}

func initialize() {
	mtSrv, err = initMTServer()
	if err != nil {
		panic(err)
	}
	go mtSrv.Start()

	conn, cancel, err = initConn()
	if err != nil {
		panic(err)
	}

	err = operations.WaitGRPCHealthy(address)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	cancel()
	mtSrv.Stop()

	dbCfg := dbutils.NewConfigFromEnv()
	err = dbutils.InitOrReset(dbCfg)
	if err != nil {
		panic(err)
	}
}

func initStree() (*tree.StateTree, error) {
	dbCfg := dbutils.NewConfigFromEnv()
	err = dbutils.InitOrReset(dbCfg)
	if err != nil {
		return nil, err
	}

	stateDb, err := db.NewSQLDB(dbCfg)
	if err != nil {
		return nil, err
	}
	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)

	stree := tree.NewStateTree(mt, scCodeStore)

	if mtSrv != nil {
		mtSrv.SetStree(stree)
	}
	return stree, nil
}

func initConn() (*grpc.ClientConn, context.CancelFunc, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	conn, err := grpc.DialContext(ctx, address, opts...)
	return conn, cancel, err
}

func initMTServer() (*tree.Server, error) {
	stree, err := initStree()
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()

	cfg := &tree.ServerConfig{
		Host: host,
		Port: port,
	}
	mtSrv = tree.NewServer(cfg, stree)
	pb.RegisterMTServiceServer(s, mtSrv)

	return mtSrv, nil
}

func Test_MTServer_GetBalance(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	ctx := context.Background()
	expectedBalance := big.NewInt(100)
	root, _, err := stree.SetBalance(ctx, common.HexToAddress(ethAddress), expectedBalance, nil)
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	resp, err := client.GetBalance(ctx, &pb.CommonGetRequest{
		EthAddress: ethAddress,
		Root:       hex.EncodeToString(root),
	})
	require.NoError(t, err)

	assert.Equal(t, expectedBalance.String(), resp.Balance, "Did not get the expected balance")
}

func Test_MTServer_GetNonce(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	ctx := context.Background()
	expectedNonce := big.NewInt(100)
	root, _, err := stree.SetNonce(ctx, common.HexToAddress(ethAddress), expectedNonce, nil)
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	resp, err := client.GetNonce(ctx, &pb.CommonGetRequest{
		EthAddress: ethAddress,
		Root:       hex.EncodeToString(root),
	})
	require.NoError(t, err)

	assert.Equal(t, expectedNonce.Uint64(), resp.Nonce, "Did not get the expected nonce")
}

func Test_MTServer_GetCode(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	expectedCode := "dead"
	code, err := hex.DecodeString(expectedCode)
	require.NoError(t, err)
	ctx := context.Background()
	root, _, err := stree.SetCode(ctx, common.HexToAddress(ethAddress), code, nil)
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	resp, err := client.GetCode(ctx, &pb.CommonGetRequest{
		EthAddress: ethAddress,
		Root:       hex.EncodeToString(root),
	})
	require.NoError(t, err)

	assert.Equal(t, expectedCode, resp.Code, "Did not get the expected code")
}

func Test_MTServer_GetCodeHash(t *testing.T) {
	data, err := os.ReadFile("test/vectors/src/merkle-tree/smt-hash-bytecode.json")
	require.NoError(t, err)

	var testVectors []struct {
		Bytecode     string
		ExpectedHash string
	}
	err = json.Unmarshal(data, &testVectors)
	require.NoError(t, err)

	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	for _, testVector := range testVectors {
		code, err := hex.DecodeString(testVector.Bytecode)
		require.NoError(t, err)
		ctx := context.Background()

		expectedHash := testVector.ExpectedHash
		root, _, err := stree.SetCode(ctx, common.HexToAddress(ethAddress), code, nil)
		require.NoError(t, err)

		client := pb.NewMTServiceClient(conn)
		resp, err := client.GetCodeHash(ctx, &pb.CommonGetRequest{
			EthAddress: ethAddress,
			Root:       hex.EncodeToString(root),
		})
		require.NoError(t, err)

		assert.Equal(t, expectedHash, resp.Hash, "Did not get the expected code hash")
	}
}

func Test_MTServer_GetStorageAt(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	expectedValue := big.NewInt(100)

	ctx := context.Background()
	position := uint64(101)
	positionBI := new(big.Int).SetUint64(position)
	root, _, err := stree.SetStorageAt(ctx, common.HexToAddress(ethAddress), positionBI, expectedValue, nil)
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	resp, err := client.GetStorageAt(ctx, &pb.GetStorageAtRequest{
		EthAddress: ethAddress,
		Root:       hex.EncodeToString(root),
		Position:   position,
	})
	require.NoError(t, err)

	assert.Equal(t, expectedValue.String(), resp.Value, "Did not get the expected storage at")
}

func Test_MTServer_ReverseHash(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	ctx := context.Background()
	expectedBalance := big.NewInt(100)
	root, _, err := stree.SetBalance(ctx, common.HexToAddress(ethAddress), expectedBalance, nil)
	require.NoError(t, err)

	key, err := tree.KeyEthAddrBalance(common.HexToAddress(ethAddress))
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	resp, err := client.ReverseHash(ctx, &pb.ReverseHashRequest{
		Hash: hex.EncodeToString(key),
		Root: hex.EncodeToString(root),
	})
	require.NoError(t, err)

	assert.Equal(t, expectedBalance.String(), resp.MtNodeValue, "Did not get the expected MT node value")
}

func Test_MTServer_SetBalance(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	expectedBalance := big.NewInt(100)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.SetBalance(ctx, &pb.SetBalanceRequest{
		EthAddress: ethAddress,
		Balance:    expectedBalance.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Success)

	newRoot, err := hex.DecodeString(resp.NewRoot)
	require.NoError(t, err)

	actualBalance, err := stree.GetBalance(ctx, common.HexToAddress(ethAddress), newRoot)
	require.NoError(t, err)

	assert.Equal(t, expectedBalance.String(), actualBalance.String(), "Did not set the expected balance")
}

func Test_MTServer_SetNonce(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	expectedNonce := big.NewInt(556)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.SetNonce(ctx, &pb.SetNonceRequest{
		EthAddress: ethAddress,
		Nonce:      expectedNonce.Uint64(),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Success)

	newRoot, err := hex.DecodeString(resp.NewRoot)
	require.NoError(t, err)

	actualNonce, err := stree.GetNonce(ctx, common.HexToAddress(ethAddress), newRoot)
	require.NoError(t, err)

	assert.Equal(t, expectedNonce.String(), actualNonce.String(), "Did not set the expected nonce")
}

func Test_MTServer_SetCode(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	expectedCode := "dead"

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.SetCode(ctx, &pb.SetCodeRequest{
		EthAddress: ethAddress,
		Code:       expectedCode,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Success)

	newRoot, err := hex.DecodeString(resp.NewRoot)
	require.NoError(t, err)

	actualCode, err := stree.GetCode(ctx, common.HexToAddress(ethAddress), newRoot)
	require.NoError(t, err)

	assert.Equal(t, expectedCode, hex.EncodeToString(actualCode), "Did not set the expected code")
}

func Test_MTServer_SetStorageAt(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	expectedValue := big.NewInt(100)
	position := uint64(101)
	positionBI := new(big.Int).SetUint64(position)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.SetStorageAt(ctx, &pb.SetStorageAtRequest{
		EthAddress: ethAddress,
		Position:   positionBI.String(),
		Value:      expectedValue.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Success)

	newRoot, err := hex.DecodeString(resp.NewRoot)
	require.NoError(t, err)

	actualStorageAt, err := stree.GetStorageAt(ctx, common.HexToAddress(ethAddress), positionBI, newRoot)
	require.NoError(t, err)

	assert.Equal(t, expectedValue.String(), actualStorageAt.String(), "Did not set the expected storage at")
}

func Test_MTServer_SetHashValue(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	ctx := context.Background()

	key := big.NewInt(200)
	expectedValue := big.NewInt(100)

	client := pb.NewMTServiceClient(conn)
	resp, err := client.SetHashValue(ctx, &pb.HashValuePair{
		Hash:  key.String(),
		Value: expectedValue.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Success)

	actualValue, err := stree.GetNodeData(ctx, key)
	require.NoError(t, err)

	assert.Equal(t, expectedValue.String(), actualValue.String(), "Did not set the expected hash value")
}

func Test_MTServer_SetStateTransitionNodes(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	const (
		totalItems = 5
		maxValue   = 3000
		maxKey     = 1000
	)

	requests := []*pb.HashValuePair{}

	ctx := context.Background()

	for i := 0; i < totalItems; i++ {
		valueBI, err := rand.Int(rand.Reader, big.NewInt(maxValue))
		require.NoError(t, err)

		keyBI, err := rand.Int(rand.Reader, big.NewInt(maxKey))
		require.NoError(t, err)

		requests = append(requests, &pb.HashValuePair{
			Hash:  keyBI.String(),
			Value: valueBI.String(),
		})
	}

	client := pb.NewMTServiceClient(conn)
	resp, err := client.SetStateTransitionNodes(ctx, &pb.SetStateTransitionNodesRequest{
		WriteHashValues: requests,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Success)

	for _, keyValuePair := range requests {
		keyBI, success := new(big.Int).SetString(keyValuePair.Hash, 10)
		require.True(t, success)

		actualValue, err := stree.GetNodeData(ctx, keyBI)
		require.NoError(t, err)

		assert.Equal(t, keyValuePair.Value, actualValue.String(), "Did not set the expected hash value bulk")
	}
}
