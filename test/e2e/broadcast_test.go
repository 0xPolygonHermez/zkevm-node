package e2e

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	serverAddress     = "localhost:61090"
	totalBatches      = 2
	totalTxsLastBatch = 5
	encodedFmt        = "encoded-%d"
	forcedBatchNumber = 18
)

var (
	ctx = context.Background()
	cfg = dbutils.NewConfigFromEnv()
)

func TestBroadcast(t *testing.T) {
	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	if testing.Short() {
		t.Skip()
	}

	require.NoError(t, operations.StartComponent("broadcast"))
	defer func() {
		require.NoError(t, operations.StopComponent("broadcast"))
	}()
	st, err := initState()
	require.NoError(t, err)

	require.NoError(t, populateDB(ctx, st))

	conn, cancel, err := initConn()
	require.NoError(t, err)
	defer func() {
		cancel()
		require.NoError(t, conn.Close())
	}()

	client := pb.NewBroadcastServiceClient(conn)

	lastBatch, err := client.GetLastBatch(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	require.Equal(t, totalBatches, int(lastBatch.BatchNumber))

	batch, err := client.GetBatch(ctx, &pb.GetBatchRequest{
		BatchNumber: uint64(totalBatches),
	})
	require.NoError(t, err)
	require.Equal(t, totalBatches, int(batch.BatchNumber))

	require.Equal(t, totalTxsLastBatch, len(batch.Transactions))

	for i, tx := range batch.Transactions {
		require.Equal(t, fmt.Sprintf(encodedFmt, i+1), tx.Encoded)
	}
	require.EqualValues(t, forcedBatchNumber, batch.ForcedBatchNumber)
}

func initState() (*state.State, error) {
	dbConfig := dbutils.NewConfigFromEnv()
	err := dbutils.InitOrReset(dbConfig)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.NewSQLDB(dbConfig)
	if err != nil {
		return nil, err
	}
	stateDb := state.NewPostgresStorage(sqlDB)

	executorClient, _, _ := executor.NewExecutorClient(ctx, executor.Config{URI: "localhost:8080"})

	mtDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, merkletree.Config{URI: "localhost:8080"})
	stateTree := merkletree.NewStateTree(mtDBClient)
	return state.NewState(state.Config{}, stateDb, executorClient, stateTree), nil
}

func initConn() (*grpc.ClientConn, context.CancelFunc, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	conn, err := grpc.DialContext(ctx, serverAddress, opts...)
	return conn, cancel, err
}

func populateDB(ctx context.Context, st *state.State) error {
	const addBatch = "INSERT INTO state.batch (batch_num, global_exit_root, timestamp, sequencer, local_exit_root, state_root) VALUES ($1, $2, $3, $4, $5, $6)"
	const addTransaction = "INSERT INTO state.transaction (batch_num, encoded, hash, received_at, l2_block_num) VALUES ($1, $2, $3, $4, $5)"
	const addForcedBatch = "INSERT INTO state.forced_batch (forced_batch_num, global_exit_root, raw_txs_data, sequencer, timestamp, batch_num, block_num) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	const blockNumber = 1

	var parentHash common.Hash
	var l2Block types.Block

	for i := 1; i <= totalBatches; i++ {
		if _, err := st.PostgresStorage.Exec(ctx, addBatch, i, common.Hash{}.String(), time.Now(), common.HexToAddress("").String(), common.Hash{}.String(), common.Hash{}.String()); err != nil {
			return err
		}
	}

	for i := 1; i <= totalTxsLastBatch; i++ {
		if i == 1 {
			parentHash = state.ZeroHash
		} else {
			parentHash = l2Block.Hash()
		}

		// Store L2 Genesis Block
		header := new(types.Header)
		header.Number = new(big.Int).SetUint64(uint64(i - 1))
		header.ParentHash = parentHash
		l2Block := types.NewBlockWithHeader(header)
		l2Block.ReceivedAt = time.Now()

		if err := st.PostgresStorage.AddL2Block(ctx, uint64(i), l2Block, []*types.Receipt{}, nil); err != nil {
			return err
		}

		if _, err := st.PostgresStorage.Exec(ctx, addTransaction, totalBatches, fmt.Sprintf(encodedFmt, i), fmt.Sprintf("hash-%d", i), time.Now(), l2Block.Number()); err != nil {
			return err
		}
	}
	if _, err := st.PostgresStorage.Exec(ctx, addBlock, blockNumber, time.Now(), ""); err != nil {
		return err
	}
	_, err := st.PostgresStorage.Exec(ctx, addForcedBatch, forcedBatchNumber, common.Hash{}.String(), "", common.HexToAddress("").String(), time.Now(), totalBatches, blockNumber)
	return err
}
