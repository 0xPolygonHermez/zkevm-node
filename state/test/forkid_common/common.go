package test

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/hashdb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
)

const (
	Ether155V = 27
)

var (
	stateTree                          *merkletree.StateTree
	stateDb                            *pgxpool.Pool
	err                                error
	StateDBCfg                         = dbutils.NewStateConfigFromEnv()
	ctx                                = context.Background()
	ExecutorClient                     executor.ExecutorServiceClient
	mtDBServiceClient                  hashdb.HashDBServiceClient
	executorClientConn, mtDBClientConn *grpc.ClientConn
	Genesis                            = state.Genesis{
		FirstBatchData: &state.BatchData{
			Transactions:   "0xf8c380808401c9c380942a3dd3eb832af982ec71669e178424b10dca2ede80b8a4d3476afe000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a40d5f56745a118d0906a34e69aec8c0db1cb8fa000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005ca1ab1e0000000000000000000000000000000000000000000000000000000005ca1ab1e1bff",
			GlobalExitRoot: common.Hash{},
			Timestamp:      1697640780,
			Sequencer:      common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		},
	}
)

func InitTestState(stateCfg state.Config) *state.State {
	InitOrResetDB(StateDBCfg)

	stateDb, err = db.NewSQLDB(StateDBCfg)
	if err != nil {
		panic(err)
	}
	// defer stateDb.Close()

	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "zkevm-prover")

	executorServerConfig := executor.Config{URI: fmt.Sprintf("%s:50071", zkProverURI), MaxGRPCMessageSize: 100000000}
	// var executorCancel context.CancelFunc
	ExecutorClient, executorClientConn, _ = executor.NewExecutorClient(ctx, executorServerConfig)
	s := executorClientConn.GetState()
	log.Infof("executorClientConn state: %s", s.String())
	/*
		defer func() {
			executorCancel()
			executorClientConn.Close()
		}()
	*/

	mtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:50061", zkProverURI)}
	// var mtDBCancel context.CancelFunc
	mtDBServiceClient, mtDBClientConn, _ = merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
	s = mtDBClientConn.GetState()
	log.Infof("stateDbClientConn state: %s", s.String())
	/*
		defer func() {
			mtDBCancel()
			mtDBClientConn.Close()
		}()
	*/

	stateTree = merkletree.NewStateTree(mtDBServiceClient)

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		panic(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)
	mt, err := l1infotree.NewL1InfoTree(32, [][32]byte{})
	if err != nil {
		panic(err)
	}
	return state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateCfg, stateDb), ExecutorClient, stateTree, eventLog, mt)
}

func InitOrResetDB(cfg db.Config) {
	if err := dbutils.InitOrResetState(cfg); err != nil {
		panic(err)
	}
}
