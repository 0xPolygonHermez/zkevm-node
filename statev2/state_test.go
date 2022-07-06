package statev2_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/merkletree"
	state "github.com/hermeznetwork/hermez-core/statev2"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/executor"
	executorclientpb "github.com/hermeznetwork/hermez-core/statev2/runtime/executor/pb"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/testutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	ether155V = 27
)

var (
	testState    *state.State
	hash1, hash2 common.Hash
	stateDb      *pgxpool.Pool
	err          error
	cfg          = dbutils.NewConfigFromEnv()
	ctx          = context.Background()
	stateCfg     = state.Config{
		MaxCumulativeGasUsed: 800000,
	}
	executorServerConfig = executor.Config{URI: "54.170.178.97:50071"}
	executorClient       executorclientpb.ExecutorServiceClient
	clientConn           *grpc.ClientConn
	stateDBServerConfig  = merkletree.Config{URI: "54.170.178.97:50061"}
)

func TestMain(m *testing.M) {
	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	stateDb, err = db.NewSQLDB(cfg)
	if err != nil {
		panic(err)
	}
	defer stateDb.Close()

	executorClient, clientConn = executor.NewExecutorClient(executorServerConfig)
	defer clientConn.Close()

	stateDbServiceClient, stateClientConn := merkletree.NewStateDBServiceClient(stateDBServerConfig)
	defer stateClientConn.Close()

	stateTree := merkletree.NewStateTree(stateDbServiceClient)

	hash1 = common.HexToHash("0x65b4699dda5f7eb4519c730e6a48e73c90d2b1c8efcd6a6abdfd28c3b8e7d7d9")
	hash2 = common.HexToHash("0x613aabebf4fddf2ad0f034a8c73aa2f9c5a6fac3a07543023e0a6ee6f36e5795")
	testState = state.NewState(stateCfg, state.NewPostgresStorage(stateDb), executorClient, stateTree)

	result := m.Run()

	os.Exit(result)
}

func TestAddBlock(t *testing.T) {
	// Init database instance
	err := dbutils.InitOrReset(cfg)
	require.NoError(t, err)
	// ctx := context.Background()
	fmt.Println("db: ", stateDb)
	tx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, tx)
	assert.NoError(t, err)
	// Add the second block
	block.BlockNumber = 2
	err = testState.AddBlock(ctx, block, tx)
	assert.NoError(t, err)
	err = tx.Commit(ctx)
	require.NoError(t, err)
	// Get the last block
	lastBlock, err := testState.GetLastBlock(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), lastBlock.BlockNumber)
	assert.Equal(t, block.BlockHash, lastBlock.BlockHash)
	assert.Equal(t, block.ParentHash, lastBlock.ParentHash)
	// Get the previous block
	prevBlock, err := testState.GetPreviousBlock(ctx, 1, nil)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), prevBlock.BlockNumber)
}

func TestAddGlobalExitRoot(t *testing.T) {
	// Init database instance
	err := dbutils.InitOrReset(cfg)
	require.NoError(t, err)
	ctx := context.Background()
	fmt.Println("db: ", stateDb)
	tx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, tx)
	assert.NoError(t, err)
	globalExitRoot := state.GlobalExitRoot{
		BlockNumber:       1,
		GlobalExitRootNum: big.NewInt(2),
		MainnetExitRoot:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		RollupExitRoot:    common.HexToHash("0x30a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
		GlobalExitRoot:    common.HexToHash("0x40a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
	}
	err = testState.AddGlobalExitRoot(ctx, &globalExitRoot, tx)
	require.NoError(t, err)
	exit, err := testState.GetLatestGlobalExitRoot(ctx, tx)
	require.NoError(t, err)
	err = tx.Commit(ctx)
	require.NoError(t, err)
	assert.Equal(t, globalExitRoot.BlockNumber, exit.BlockNumber)
	assert.Equal(t, globalExitRoot.GlobalExitRootNum, exit.GlobalExitRootNum)
	assert.Equal(t, globalExitRoot.MainnetExitRoot, exit.MainnetExitRoot)
	assert.Equal(t, globalExitRoot.RollupExitRoot, exit.RollupExitRoot)
	assert.Equal(t, globalExitRoot.GlobalExitRoot, exit.GlobalExitRoot)
}

func TestAddForcedBatch(t *testing.T) {
	// Init database instance
	err := dbutils.InitOrReset(cfg)
	require.NoError(t, err)
	ctx := context.Background()
	tx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, tx)
	assert.NoError(t, err)
	b := common.Hex2Bytes("0x617b3a3528F9")
	assert.NoError(t, err)
	var bN uint64 = 3
	forcedBatch := state.ForcedBatch{
		BlockNumber:       1,
		ForcedBatchNumber: 2,
		BatchNumber:       &bN,
		GlobalExitRoot:    common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Sequencer:         common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
		RawTxsData:        b,
		ForcedAt:          time.Now(),
	}
	err = testState.AddForcedBatch(ctx, &forcedBatch, tx)
	require.NoError(t, err)
	fb, err := testState.GetForcedBatch(ctx, 2, tx)
	require.NoError(t, err)
	err = tx.Commit(ctx)
	require.NoError(t, err)
	assert.Equal(t, forcedBatch.BlockNumber, fb.BlockNumber)
	assert.Equal(t, forcedBatch.BatchNumber, fb.BatchNumber)
	assert.Equal(t, forcedBatch.ForcedBatchNumber, fb.ForcedBatchNumber)
	assert.NotEqual(t, time.Time{}, fb.ForcedAt)
	assert.Equal(t, forcedBatch.GlobalExitRoot, fb.GlobalExitRoot)
	assert.Equal(t, forcedBatch.RawTxsData, fb.RawTxsData)
	// Test GetNextForcedBatches
	tx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	forcedBatch = state.ForcedBatch{
		BlockNumber:       1,
		ForcedBatchNumber: 3,
		BatchNumber:       nil,
		GlobalExitRoot:    common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Sequencer:         common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
		RawTxsData:        b,
		ForcedAt:          time.Now(),
	}
	err = testState.AddForcedBatch(ctx, &forcedBatch, tx)
	require.NoError(t, err)
	batches, err := testState.GetNextForcedBatches(ctx, 1, tx)
	require.NoError(t, err)
	require.NoError(t, tx.Commit(ctx))
	assert.Equal(t, forcedBatch.BlockNumber, batches[0].BlockNumber)
	assert.Equal(t, forcedBatch.BatchNumber, batches[0].BatchNumber)
	assert.Equal(t, forcedBatch.ForcedBatchNumber, batches[0].ForcedBatchNumber)
	assert.NotEqual(t, time.Time{}, batches[0].ForcedAt)
	assert.Equal(t, forcedBatch.GlobalExitRoot, batches[0].GlobalExitRoot)
	assert.Equal(t, forcedBatch.RawTxsData, batches[0].RawTxsData)
	// Test AddBatchNumberInForcedBatch
	tx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	err = testState.AddBatchNumberInForcedBatch(ctx, 3, 2, tx)
	require.NoError(t, err)
	fb, err = testState.GetForcedBatch(ctx, 3, tx)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), *fb.BatchNumber)
	require.NoError(t, tx.Commit(ctx))
}

func TestAddVirtualBatch(t *testing.T) {
	// Init database instance
	err := dbutils.InitOrReset(cfg)
	require.NoError(t, err)
	ctx := context.Background()
	tx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, tx)
	assert.NoError(t, err)
	batch := state.Batch{
		BatchNumber:    1,
		GlobalExitRoot: common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Coinbase:       common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
		Timestamp:      time.Now(),
		BatchL2Data:    common.Hex2Bytes("0x617b3a3528F9"),
	}
	err = testState.StoreBatchHeader(ctx, batch, tx)
	require.NoError(t, err)
	virtualBatch := state.VirtualBatch{
		BlockNumber: 1,
		BatchNumber: 1,
		TxHash:      common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Sequencer:   common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
	}
	err = testState.AddVirtualBatch(ctx, &virtualBatch, tx)
	require.NoError(t, err)
	require.NoError(t, tx.Commit(ctx))
}

func TestExecuteTransaction(t *testing.T) {
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var sequencerBalance = 4000000
	var scCounterAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")
	var scInteractionAddress = common.HexToAddress("0x85e844b762A271022b692CF99cE5c59BA0650Ac8")

	scCounterByteCode, err := testutils.ReadBytecode("Counter/Counter.bin")
	require.NoError(t, err)
	scInteractionByteCode, err := testutils.ReadBytecode("Interaction/Interaction.bin")
	require.NoError(t, err)

	// Deploy counter.sol
	tx0 := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     common.Hex2Bytes(scCounterByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx0, err := auth.Signer(auth.From, tx0)
	require.NoError(t, err)

	// Deploy interaction.sol
	tx1 := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(0),
		Data:     common.Hex2Bytes(scInteractionByteCode),
	})

	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)

	// Call setCounterAddr method from Interaction SC to set Counter SC Address
	tx2 := types.NewTransaction(2, scInteractionAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("ec39b429000000000000000000000000"+strings.TrimPrefix(scCounterAddress.String(), "0x")))
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)

	// Increment Counter calling Counter SC
	tx3 := types.NewTransaction(3, scCounterAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("d09de08a"))
	signedTx3, err := auth.Signer(auth.From, tx3)
	require.NoError(t, err)

	// Retrieve counter value calling Interaction SC (this is the real test as Interaction SC will call Counter SC)
	tx4 := types.NewTransaction(4, scInteractionAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("a87d942c"))
	signedTx4, err := auth.Signer(auth.From, tx4)
	require.NoError(t, err)

	// Increment Counter calling again
	tx5 := types.NewTransaction(5, scCounterAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("d09de08a"))
	signedTx5, err := auth.Signer(auth.From, tx5)
	require.NoError(t, err)

	// Retrieve counter value again
	tx6 := types.NewTransaction(6, scInteractionAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("a87d942c"))
	signedTx6, err := auth.Signer(auth.From, tx6)
	require.NoError(t, err)

	batchL2Data, err := state.EncodeTransactions([]types.Transaction{*signedTx0, *signedTx1, *signedTx2, *signedTx3, *signedTx4, *signedTx5, *signedTx6})
	require.NoError(t, err)

	// Create Batch
	processBatchRequest := &executorclientpb.ProcessBatchRequest{
		BatchNum:             1,
		Coinbase:             sequencerAddress.String(),
		BatchL2Data:          batchL2Data,
		OldStateRoot:         common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		GlobalExitRoot:       common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		OldLocalExitRoot:     common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		EthTimestamp:         uint64(time.Now().Unix()),
		UpdateMerkleTree:     0,
		GenerateExecuteTrace: 0,
		GenerateCallTrace:    0,
	}

	processBatchResponse, err := executorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)

	file, _ := json.MarshalIndent(processBatchResponse, "", " ")
	err = ioutil.WriteFile("trace.json", file, 0644)
	require.NoError(t, err)

}

func TestExecutor(t *testing.T) {

	// Create Batch
	processBatchRequest := &executorclientpb.ProcessBatchRequest{
		BatchNum:             1,
		Coinbase:             common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D").String(),
		BatchL2Data:          common.Hex2Bytes("ee80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880801cee7e01dc62f69a12c3510c6d64de04ee6346d84b6a017f3e786c7d87f963e75d8cc91fa983cd6d9cf55fff80d73bd26cd333b0f098acc1e58edb1fd484ad731b"),
		OldStateRoot:         common.Hex2Bytes("2dc4db4293af236cb329700be43f08ace740a05088f8c7654736871709687e90"),
		GlobalExitRoot:       common.Hex2Bytes("090bcaf734c4f06c93954a827b45a6e8c67b8e0fd1e0a35a1c5982d6961828f9"),
		OldLocalExitRoot:     common.Hex2Bytes("17c04c3760510b48c6012742c540a81aba4bca2f78b9d14bfd2f123e2e53ea3e"),
		EthTimestamp:         uint64(1944498031),
		UpdateMerkleTree:     0,
		GenerateExecuteTrace: 0,
		GenerateCallTrace:    0,
	}

	processBatchResponse, err := executorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)

	log.Debugf("request:%v", processBatchRequest)
	log.Debugf("new_state_root=%v", common.BigToHash(new(big.Int).SetBytes(processBatchResponse.NewStateRoot)))

	file, _ := json.MarshalIndent(processBatchResponse, "", " ")
	err = ioutil.WriteFile("trace.json", file, 0644)
	require.NoError(t, err)
}
