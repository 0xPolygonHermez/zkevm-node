package tree_test

import (
	"context"
	"fmt"
	"math/big"
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
)

const (
	host = "0.0.0.0"
	port = 50051

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
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)

	stree := tree.NewStateTree(mt, scCodeStore)

	if mtSrv != nil {
		mtSrv.SetStree(stree)
	}
	return stree, nil
}

func initConn() (*grpc.ClientConn, context.CancelFunc, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
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

	cfg := &tree.Config{
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

	expectedBalance := big.NewInt(100)
	root, _, err := stree.SetBalance(common.HexToAddress(ethAddress), expectedBalance, nil)
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.GetBalance(ctx, &pb.GetBalanceRequest{
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

	expectedNonce := big.NewInt(200)
	root, _, err := stree.SetNonce(common.HexToAddress(ethAddress), expectedNonce)
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.GetNonce(ctx, &pb.GetNonceRequest{
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
	root, _, err := stree.SetCode(common.HexToAddress(ethAddress), code)
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.GetCode(ctx, &pb.GetCodeRequest{
		EthAddress: ethAddress,
		Root:       hex.EncodeToString(root),
	})
	require.NoError(t, err)

	assert.Equal(t, expectedCode, resp.Code, "Did not get the expected code")
}

func Test_MTServer_GetCodeHash(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	code, err := hex.DecodeString("dead")
	require.NoError(t, err)

	// code hash from test vectors
	expectedHash := "0244ec1a137a24c92404de9f9c39907be151026a4eb7f9cfea60a5740e8a73b7"
	root, _, err := stree.SetCode(common.HexToAddress(ethAddress), code)
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.GetCodeHash(ctx, &pb.GetCodeHashRequest{
		EthAddress: ethAddress,
		Root:       hex.EncodeToString(root),
	})
	require.NoError(t, err)

	assert.Equal(t, expectedHash, resp.Hash, "Did not get the expected code hash")
}

func Test_MTServer_GetStorageAt(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	expectedValue := big.NewInt(100)

	position := uint64(101)
	positionBI := new(big.Int).SetUint64(position)
	root, _, err := stree.SetStorageAt(common.HexToAddress(ethAddress), common.BigToHash(positionBI), expectedValue)
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
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

	expectedBalance := big.NewInt(100)
	root, _, err := stree.SetBalance(common.HexToAddress(ethAddress), expectedBalance)
	require.NoError(t, err)

	key, err := tree.GetKey(tree.LeafTypeBalance, common.HexToAddress(ethAddress), nil, tree.DefaultMerkleTreeArity, nil)
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
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

	newRoot, err := stree.GetCurrentRoot()
	require.NoError(t, err)

	actualBalance, err := stree.GetBalance(common.HexToAddress(ethAddress), newRoot)
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

	newRoot, err := stree.GetCurrentRoot()
	require.NoError(t, err)

	actualNonce, err := stree.GetNonce(common.HexToAddress(ethAddress), newRoot)
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

	newRoot, err := stree.GetCurrentRoot()
	require.NoError(t, err)

	actualCode, err := stree.GetCode(common.HexToAddress(ethAddress), newRoot)
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
		Position:   common.BytesToHash(positionBI.Bytes()).String(),
		Value:      expectedValue.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Success)

	newRoot, err := stree.GetCurrentRoot()
	require.NoError(t, err)

	actualStorageAt, err := stree.GetStorageAt(common.HexToAddress(ethAddress), common.BigToHash(positionBI), newRoot)
	require.NoError(t, err)

	assert.Equal(t, expectedValue.String(), actualStorageAt.String(), "Did not set the expected storage at")
}

func Test_MTServer_SetHashValue(t *testing.T) {
	require.NoError(t, dbutils.InitOrReset(dbutils.NewConfigFromEnv()))
	stree, err := initStree()
	require.NoError(t, err)

	initialBalance := big.NewInt(100)
	_, _, err = stree.SetBalance(common.HexToAddress(ethAddress), initialBalance)
	require.NoError(t, err)

	key, err := tree.GetKey(tree.LeafTypeBalance, common.HexToAddress(ethAddress), nil, tree.DefaultMerkleTreeArity, nil)
	require.NoError(t, err)

	expectedValue := big.NewInt(100)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.SetHashValue(ctx, &pb.SetHashValueRequest{
		Hash:  hex.EncodeToString(key),
		Value: expectedValue.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Success)

	newRoot, err := stree.GetCurrentRoot()
	require.NoError(t, err)

	actualValue, err := stree.GetBalance(common.HexToAddress(ethAddress), newRoot)
	require.NoError(t, err)

	assert.Equal(t, expectedValue.String(), actualValue.String(), "Did not set the expected hash value")
}
