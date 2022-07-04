package statev2_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/merkletree"
	state "github.com/hermeznetwork/hermez-core/statev2"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/executor"
	executorclientpb "github.com/hermeznetwork/hermez-core/statev2/runtime/executor/pb"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
// ether155V = 27
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

func TestGetTxsHashesToDelete(t *testing.T) {
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
	batch1 := state.Batch{
		BatchNumber:    1,
		GlobalExitRoot: common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Coinbase:       common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
		Timestamp:      time.Now(),
		BatchL2Data:    common.Hex2Bytes("0x617b3a3528F9"),
	}

	err = testState.StoreBatchHeader(ctx, batch1, tx)
	require.NoError(t, err)
	virtualBatch1 := state.VirtualBatch{
		BlockNumber: 1,
		BatchNumber: 1,
		TxHash:      common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Sequencer:   common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
	}

	batch2 := state.Batch{
		BatchNumber:    2,
		GlobalExitRoot: common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Coinbase:       common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
		Timestamp:      time.Now(),
		BatchL2Data:    common.Hex2Bytes("0x617b3a3528F9"),
	}

	err = testState.StoreBatchHeader(ctx, batch2, tx)
	require.NoError(t, err)
	virtualBatch2 := state.VirtualBatch{
		BlockNumber: 1,
		BatchNumber: 2,
		TxHash:      common.HexToHash("0x132"),
		Sequencer:   common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
	}
	err = testState.AddVirtualBatch(ctx, &virtualBatch1, tx)
	require.NoError(t, err)
	err = testState.AddVirtualBatch(ctx, &virtualBatch2, tx)
	require.NoError(t, err)
	require.NoError(t, tx.Commit(ctx))

	_, err = testState.Exec(ctx, "INSERT INTO statev2.l2block (block_num, block_hash, received_at, batch_num) VALUES ($1, $2, $3, $4)", 1, "0x423", time.Now(), 1)
	require.NoError(t, err)
	l2Tx := types.NewTransaction(1, common.Address{}, big.NewInt(10), 21000, big.NewInt(1), []byte{})
	_, err = testState.Exec(ctx, "INSERT INTO statev2.transaction (l2_block_num, encoded, hash) VALUES ($1, $2, $3)",
		virtualBatch1.BatchNumber, fmt.Sprintf("encoded-%d", virtualBatch1.BatchNumber), l2Tx.Hash().Hex())
	require.NoError(t, err)

	txHashes, err := testState.GetTxsHashesToDelete(ctx, block.BlockNumber, nil)
	require.NoError(t, err)
	require.Equal(t, l2Tx.Hash().Hex(), txHashes[0].Hex())
}

/*
func TestExecuteTransaction(t *testing.T) {
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var sequencerBalance = 4000000
	scCounterByteCode, err := testutils.ReadBytecode("Counter/Counter.bin")
	require.NoError(t, err)

	// Deploy counter.sol
	tx := types.NewTx(&types.LegacyTx{
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

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	// Encode transaction
	v, r, s := signedTx.RawSignatureValues()
	sign := 1 - (v.Uint64() & 1)

	txCodedRlp, err := rlp.EncodeToBytes([]interface{}{
		signedTx.Nonce(),
		signedTx.GasPrice(),
		signedTx.Gas(),
		signedTx.To(),
		signedTx.Value(),
		signedTx.Data(),
		signedTx.ChainId(), uint(0), uint(0),
	})
	require.NoError(t, err)

	newV := new(big.Int).Add(big.NewInt(ether155V), big.NewInt(int64(sign)))
	newRPadded := fmt.Sprintf("%064s", r.Text(hex.Base))
	newSPadded := fmt.Sprintf("%064s", s.Text(hex.Base))
	newVPadded := fmt.Sprintf("%02s", newV.Text(hex.Base))
	batchL2Data, err := hex.DecodeString(hex.EncodeToString(txCodedRlp) + newRPadded + newSPadded + newVPadded)
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
		UpdateMerkleTree:     false,
		GenerateExecuteTrace: false,
		GenerateCallTrace:    false,
	}

	log.Debugf("%v", processBatchRequest)

	processBatchResponse, err := executorClient.ProcessBatch(ctx, processBatchRequest)
	require.NoError(t, err)

	file, _ := json.MarshalIndent(processBatchResponse, "", " ")
	err = ioutil.WriteFile("trace.json", file, 0644)
	require.NoError(t, err)
}
*/
/*
func TestGenesis(t *testing.T) {
	balances := map[common.Address]*big.Int{
		common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1000),
		common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(2000),
	}

	nonces := map[common.Address]*big.Int{
		common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA"): big.NewInt(1),
		common.HexToAddress("0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB"): big.NewInt(1),
	}

	smartContracts := map[common.Address][]byte{
		common.HexToAddress("0xae4bb80be56b819606589de61d5ec3b522eeb032"): common.Hex2Bytes("608060405234801561001057600080fd5b50600436106100675760003560e01c806333d6247d1161005057806333d6247d146100a85780633ed691ef146100bd578063a3c573eb146100d257600080fd5b806301fd90441461006c5780633381fe9014610088575b600080fd5b61007560015481565b6040519081526020015b60405180910390f35b6100756100963660046101c7565b60006020819052908152604090205481565b6100bb6100b63660046101c7565b610117565b005b43600090815260208190526040902054610075565b6002546100f29073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161007f565b60025473ffffffffffffffffffffffffffffffffffffffff1633146101c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603460248201527f476c6f62616c45786974526f6f744d616e616765724c323a3a7570646174654560448201527f786974526f6f743a204f4e4c595f425249444745000000000000000000000000606482015260840160405180910390fd5b600155565b6000602082840312156101d957600080fd5b503591905056fea2646970667358221220d6ed73b81f538d38669b0b750b93be08ca365978fae900eedc9ca93131c97ca664736f6c63430008090033"),
	}

	storage := map[common.Address]map[*big.Int]*big.Int{
		common.HexToAddress("0xae4bb80be56b819606589de61d5ec3b522eeb032"): {new(big.Int).SetBytes(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000002")): new(big.Int).SetBytes(common.Hex2Bytes("9d98deabc42dd696deb9e40b4f1cab7ddbf55988"))},
	}

	genesis := state.Genesis{
		Balances:       balances,
		Nonces:         nonces,
		SmartContracts: smartContracts,
		Storage:        storage,
	}
	err := testState.SetGenesis(ctx, genesis, nil)
	require.NoError(t, err)
}
*/
