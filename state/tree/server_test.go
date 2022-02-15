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
	mtSrv   *tree.Server
	address = fmt.Sprintf("%s:%d", host, port)
)

func initStree() (*tree.StateTree, error) {
	dbCfg := dbutils.NewConfigFromEnv()
	err := dbutils.InitOrReset(dbCfg)
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

	return tree.NewStateTree(mt, scCodeStore), nil
}

func initConn() (*grpc.ClientConn, context.CancelFunc, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	conn, err := grpc.DialContext(ctx, address, opts...)
	return conn, cancel, err
}

func initMTServer() (*tree.Server, *tree.StateTree, error) {
	stree, err := initStree()
	if err != nil {
		return nil, nil, err
	}

	s := grpc.NewServer()

	cfg := &tree.Config{
		Host: host,
		Port: port,
	}
	mtSrv = tree.NewServer(cfg, stree)
	pb.RegisterMTServiceServer(s, mtSrv)

	return mtSrv, stree, nil
}

func Test_MTServer_GetBalance(t *testing.T) {
	mtSrv, stree, err := initMTServer()
	require.NoError(t, err)
	go mtSrv.Start()
	defer mtSrv.Stop()

	conn, cancel, err := initConn()
	require.NoError(t, err)
	defer func() {
		cancel()
		require.NoError(t, conn.Close())
	}()

	err = operations.WaitGRPCHealthy(address)
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
	mtSrv, stree, err := initMTServer()
	require.NoError(t, err)
	go mtSrv.Start()
	defer mtSrv.Stop()

	conn, cancel, err := initConn()
	require.NoError(t, err)
	defer func() {
		cancel()
		require.NoError(t, conn.Close())
	}()

	err = operations.WaitGRPCHealthy(address)
	require.NoError(t, err)

	expectedNonce := big.NewInt(100)
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
	mtSrv, stree, err := initMTServer()
	require.NoError(t, err)
	go mtSrv.Start()
	defer mtSrv.Stop()

	conn, cancel, err := initConn()
	require.NoError(t, err)
	defer func() {
		cancel()
		require.NoError(t, conn.Close())
	}()

	err = operations.WaitGRPCHealthy(address)
	require.NoError(t, err)

	expectedCode := "dead"
	root, _, err := stree.SetCode(common.HexToAddress(ethAddress), []byte(expectedCode))
	require.NoError(t, err)

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.GetCode(ctx, &pb.GetCodeRequest{
		EthAddress: ethAddress,
		Root:       hex.EncodeToString(root),
	})
	require.NoError(t, err)

	assert.Equal(t, string(expectedCode), resp.Code, "Did not get the expected code")
}
