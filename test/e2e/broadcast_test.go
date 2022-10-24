package e2e

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/constants"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
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
	ctx             = context.Background()
	stateDBCfg      = dbutils.NewStateConfigFromEnv()
	ger             = common.HexToHash("deadbeef")
	mainnetExitRoot = common.HexToHash("caffe")
	rollupExitRoot  = common.HexToHash("bead")
)

func TestBroadcast(t *testing.T) {
	initOrResetDB()

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

	client, conn, cancel := broadcast.NewClient(ctx, serverAddress)
	defer func() {
		cancel()
		require.NoError(t, conn.Close())
	}()

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

	require.Equal(t, mainnetExitRoot.String(), batch.MainnetExitRoot)
	require.Equal(t, rollupExitRoot.String(), batch.RollupExitRoot)
}

func initState() (*state.State, error) {
	initOrResetDB()
	sqlDB, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		return nil, err
	}
	stateDb := state.NewPostgresStorage(sqlDB)
	executorUri := testutils.GetEnv(constants.ENV_ZKPROVER_URI, "127.0.0.1:50071")
	merkleTreeUri := testutils.GetEnv(constants.ENV_MERKLETREE_URI, "127.0.0.1:50061")
	executorClient, _, _ := executor.NewExecutorClient(ctx, executor.Config{URI: executorUri})
	mtDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, merkletree.Config{URI: merkleTreeUri})
	stateTree := merkletree.NewStateTree(mtDBClient)
	return state.NewState(state.Config{}, stateDb, executorClient, stateTree), nil
}

func populateDB(ctx context.Context, st *state.State) error {
	const blockNumber = 1

	var parentHash common.Hash
	var l2Block types.Block

	const addBatch = "INSERT INTO state.batch (batch_num, global_exit_root, timestamp, coinbase, local_exit_root, state_root) VALUES ($1, $2, $3, $4, $5, $6)"
	for i := 1; i <= totalBatches; i++ {
		if _, err := st.PostgresStorage.Exec(ctx, addBatch, i, ger.String(), time.Now(), common.HexToAddress("").String(), common.Hash{}.String(), common.Hash{}.String()); err != nil {
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

		if err := st.PostgresStorage.AddL2Block(ctx, totalBatches, l2Block, []*types.Receipt{}, nil); err != nil {
			return err
		}

		const addTransaction = "INSERT INTO state.transaction (hash, encoded, l2_block_num) VALUES ($1, $2, $3)"
		if _, err := st.PostgresStorage.Exec(ctx, addTransaction, fmt.Sprintf("hash-%d", i), fmt.Sprintf(encodedFmt, i), l2Block.Number().Uint64()); err != nil {
			return err
		}
	}

	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := st.PostgresStorage.Exec(ctx, addBlock, blockNumber, time.Now(), ""); err != nil {
		return err
	}

	const addForcedBatch = "INSERT INTO state.forced_batch (forced_batch_num, global_exit_root, raw_txs_data, coinbase, timestamp, batch_num, block_num) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if _, err := st.PostgresStorage.Exec(ctx, addForcedBatch, forcedBatchNumber, ger.String(), "", common.HexToAddress("").String(), time.Now(), totalBatches, blockNumber); err != nil {
		return err
	}

	const addExitRoots = "INSERT INTO state.exit_root (block_num, global_exit_root, mainnet_exit_root, rollup_exit_root, global_exit_root_num) VALUES ($1, $2, $3, $4, $5)"
	_, err := st.PostgresStorage.Exec(ctx, addExitRoots, blockNumber, ger, mainnetExitRoot, rollupExitRoot, 1)
	return err
}

func initOrResetDB() {
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}
}
