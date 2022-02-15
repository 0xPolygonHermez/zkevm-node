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

func initMTServer(stree *tree.StateTree) (*tree.Server, error) {
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
	stree, err := initStree()
	if err != nil {
		t.Fatalf("Could not initialize state tree, %v", err)
	}

	mtSrv, err := initMTServer(stree)
	if err != nil {
		t.Fatalf("Could not initialize MTServer, %v", err)
	}
	go mtSrv.Start()
	defer mtSrv.Stop()

	conn, cancel, err := initConn()
	if err != nil {
		t.Fatalf("Failed to initialize grpc connection: %v", err)
	}
	defer func() {
		cancel()
		if err := conn.Close(); err != nil {
			t.Fatalf("Failed to close conn: %v", err)
		}
	}()

	err = operations.WaitGRPCHealthy(address)
	if err != nil {
		t.Fatalf("gRPC server did not come up on time: %v", err)
	}

	expectedBalance := big.NewInt(100)
	root, _, err := stree.SetBalance(common.HexToAddress(ethAddress), expectedBalance, nil)
	if err != nil {
		t.Fatalf("could not set balance: %v", err)
	}

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.GetBalance(ctx, &pb.GetBalanceRequest{
		EthAddress: ethAddress,
		Root:       hex.EncodeToString(root),
	})
	if err != nil {
		t.Fatalf("GetBalance failed: %v", err)
	}

	assert.Equal(t, expectedBalance.String(), resp.Balance, "Did not get the expected balance")
}

func Test_MTServer_GetNonce(t *testing.T) {
	stree, err := initStree()
	if err != nil {
		t.Fatalf("Could not initialize state tree, %v", err)
	}

	mtSrv, err := initMTServer(stree)
	if err != nil {
		t.Fatalf("Could not initialize MTServer, %v", err)
	}
	go mtSrv.Start()
	defer mtSrv.Stop()

	conn, cancel, err := initConn()
	if err != nil {
		t.Fatalf("Failed to initialize grpc connection: %v", err)
	}
	defer func() {
		cancel()
		if err := conn.Close(); err != nil {
			t.Fatalf("Failed to close conn: %v", err)
		}
	}()

	err = operations.WaitGRPCHealthy(address)
	if err != nil {
		t.Fatalf("gRPC server did not come up on time: %v", err)
	}

	expectedNonce := big.NewInt(100)
	root, _, err := stree.SetNonce(common.HexToAddress(ethAddress), expectedNonce)
	if err != nil {
		t.Fatalf("could not set balance: %v", err)
	}

	client := pb.NewMTServiceClient(conn)
	ctx := context.Background()
	resp, err := client.GetNonce(ctx, &pb.GetNonceRequest{
		EthAddress: ethAddress,
		Root:       hex.EncodeToString(root),
	})
	if err != nil {
		t.Fatalf("GetNonce failed: %v", err)
	}
	assert.Equal(t, expectedNonce.Uint64(), resp.Nonce, "Did not get the expected nonce")
}
