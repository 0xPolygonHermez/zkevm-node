package pgstatestorage_test

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/hashdb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var (
	testState  *state.State
	stateTree  *merkletree.StateTree
	stateDb    *pgxpool.Pool
	err        error
	stateDBCfg = dbutils.NewStateConfigFromEnv()
	ctx        = context.Background()
	stateCfg   = state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          5,
			Version:         "",
		}},
	}
	forkID                             uint64 = 5
	executorClient                     executor.ExecutorServiceClient
	mtDBServiceClient                  hashdb.HashDBServiceClient
	executorClientConn, mtDBClientConn *grpc.ClientConn
	batchResources                     = state.BatchResources{
		ZKCounters: state.ZKCounters{
			UsedKeccakHashes: 1,
		},
		Bytes: 1,
	}
	closingReason = state.GlobalExitRootDeadlineClosingReason
)

func initOrResetDB() {
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	initOrResetDB()

	stateDb, err = db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	defer stateDb.Close()

	zkProverURI := testutils.GetEnv("ZKPROVER_URI", "localhost")

	executorServerConfig := executor.Config{URI: fmt.Sprintf("%s:50071", zkProverURI), MaxGRPCMessageSize: 100000000}
	var executorCancel context.CancelFunc
	executorClient, executorClientConn, executorCancel = executor.NewExecutorClient(ctx, executorServerConfig)
	s := executorClientConn.GetState()
	log.Infof("executorClientConn state: %s", s.String())
	defer func() {
		executorCancel()
		executorClientConn.Close()
	}()

	mtDBServerConfig := merkletree.Config{URI: fmt.Sprintf("%s:50061", zkProverURI)}
	var mtDBCancel context.CancelFunc
	mtDBServiceClient, mtDBClientConn, mtDBCancel = merkletree.NewMTDBServiceClient(ctx, mtDBServerConfig)
	s = mtDBClientConn.GetState()
	log.Infof("stateDbClientConn state: %s", s.String())
	defer func() {
		mtDBCancel()
		mtDBClientConn.Close()
	}()

	stateTree = merkletree.NewStateTree(mtDBServiceClient)

	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		panic(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	testState = state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateCfg, stateDb), executorClient, stateTree, eventLog)

	result := m.Run()

	os.Exit(result)
}

var (
	pgStateStorage *pgstatestorage.PostgresStorage
	block          = &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
)

func setup() {
	cfg := state.Config{
		MaxLogsCount:      10000,
		MaxLogsBlockRange: 10000,
	}
	pgStateStorage = pgstatestorage.NewPostgresStorage(cfg, stateDb)
}

func TestGetBatchByL2BlockNumber(t *testing.T) {
	setup()
	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	err = testState.AddBlock(ctx, block, dbTx)
	assert.NoError(t, err)

	batchNumber := uint64(1)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES ($1)", batchNumber)
	assert.NoError(t, err)

	time := time.Now()
	blockNumber := big.NewInt(1)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      0,
		GasPrice: big.NewInt(0),
	})

	receipt := &types.Receipt{
		Type:              uint8(tx.Type()),
		PostState:         state.ZeroHash.Bytes(),
		CumulativeGasUsed: 0,
		EffectiveGasPrice: big.NewInt(0),
		BlockNumber:       blockNumber,
		GasUsed:           tx.Gas(),
		TxHash:            tx.Hash(),
		TransactionIndex:  0,
		Status:            types.ReceiptStatusSuccessful,
	}

	header := &types.Header{
		Number:     big.NewInt(1),
		ParentHash: state.ZeroHash,
		Coinbase:   state.ZeroAddress,
		Root:       state.ZeroHash,
		GasUsed:    1,
		GasLimit:   10,
		Time:       uint64(time.Unix()),
	}
	transactions := []*types.Transaction{tx}

	receipts := []*types.Receipt{receipt}

	// Create block to be able to calculate its hash
	l2Block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
	receipt.BlockHash = l2Block.Hash()

	storeTxsEGPData := []state.StoreTxEGPData{}
	for range transactions {
		storeTxsEGPData = append(storeTxsEGPData, state.StoreTxEGPData{EGPLog: nil, EffectivePercentage: state.MaxEffectivePercentage})
	}

	err = pgStateStorage.AddL2Block(ctx, batchNumber, l2Block, receipts, storeTxsEGPData, dbTx)
	require.NoError(t, err)
	result, err := pgStateStorage.BatchNumberByL2BlockNumber(ctx, l2Block.Number().Uint64(), dbTx)
	require.NoError(t, err)
	assert.Equal(t, batchNumber, result)
	require.NoError(t, dbTx.Commit(ctx))
}

func TestAddAndGetSequences(t *testing.T) {
	initOrResetDB()

	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, dbTx)
	assert.NoError(t, err)

	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (0)")
	require.NoError(t, err)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (1)")
	require.NoError(t, err)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (2)")
	require.NoError(t, err)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (3)")
	require.NoError(t, err)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (4)")
	require.NoError(t, err)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (5)")
	require.NoError(t, err)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (6)")
	require.NoError(t, err)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (7)")
	require.NoError(t, err)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (8)")
	require.NoError(t, err)

	sequence := state.Sequence{
		FromBatchNumber: 0,
		ToBatchNumber:   3,
	}
	err = testState.AddSequence(ctx, sequence, dbTx)
	require.NoError(t, err)

	sequence2 := state.Sequence{
		FromBatchNumber: 3,
		ToBatchNumber:   7,
	}
	err = testState.AddSequence(ctx, sequence2, dbTx)
	require.NoError(t, err)

	sequence3 := state.Sequence{
		FromBatchNumber: 7,
		ToBatchNumber:   7,
	}
	err = testState.AddSequence(ctx, sequence3, dbTx)
	require.NoError(t, err)

	// Insert it again to test on conflict
	sequence3.ToBatchNumber = 8
	err = testState.AddSequence(ctx, sequence3, dbTx)
	require.NoError(t, err)

	sequences, err := testState.GetSequences(ctx, 0, dbTx)
	require.NoError(t, err)
	require.Equal(t, 3, len(sequences))
	require.Equal(t, uint64(0), sequences[0].FromBatchNumber)
	require.Equal(t, uint64(3), sequences[1].FromBatchNumber)
	require.Equal(t, uint64(7), sequences[2].FromBatchNumber)
	require.Equal(t, uint64(3), sequences[0].ToBatchNumber)
	require.Equal(t, uint64(7), sequences[1].ToBatchNumber)
	require.Equal(t, uint64(8), sequences[2].ToBatchNumber)

	sequences, err = testState.GetSequences(ctx, 3, dbTx)
	require.NoError(t, err)
	require.Equal(t, 2, len(sequences))
	require.Equal(t, uint64(3), sequences[0].FromBatchNumber)
	require.Equal(t, uint64(7), sequences[1].FromBatchNumber)
	require.Equal(t, uint64(7), sequences[0].ToBatchNumber)
	require.Equal(t, uint64(8), sequences[1].ToBatchNumber)

	require.NoError(t, dbTx.Commit(ctx))
}

func TestAddGlobalExitRoot(t *testing.T) {
	// Init database instance
	initOrResetDB()

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
	globalExitRoot := state.GlobalExitRoot{
		BlockNumber:     1,
		MainnetExitRoot: common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		RollupExitRoot:  common.HexToHash("0x30a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
		GlobalExitRoot:  common.HexToHash("0x40a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
	}
	err = testState.AddGlobalExitRoot(ctx, &globalExitRoot, tx)
	require.NoError(t, err)
	exit, _, err := testState.GetLatestGlobalExitRoot(ctx, math.MaxInt64, tx)
	require.NoError(t, err)
	err = tx.Commit(ctx)
	require.NoError(t, err)
	assert.Equal(t, globalExitRoot.BlockNumber, exit.BlockNumber)
	assert.Equal(t, globalExitRoot.MainnetExitRoot, exit.MainnetExitRoot)
	assert.Equal(t, globalExitRoot.RollupExitRoot, exit.RollupExitRoot)
	assert.Equal(t, globalExitRoot.GlobalExitRoot, exit.GlobalExitRoot)
}

func TestVerifiedBatch(t *testing.T) {
	initOrResetDB()

	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, dbTx)
	assert.NoError(t, err)
	//require.NoError(t, tx.Commit(ctx))

	lastBlock, err := testState.GetLastBlock(ctx, dbTx)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), lastBlock.BlockNumber)

	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (1)")

	require.NoError(t, err)
	virtualBatch := state.VirtualBatch{
		BlockNumber: 1,
		BatchNumber: 1,
		TxHash:      common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
	}
	err = testState.AddVirtualBatch(ctx, &virtualBatch, dbTx)
	require.NoError(t, err)
	expectedVerifiedBatch := state.VerifiedBatch{
		BlockNumber: 1,
		BatchNumber: 1,
		StateRoot:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f2"),
		Aggregator:  common.HexToAddress("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		TxHash:      common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		IsTrusted:   true,
	}
	err = testState.AddVerifiedBatch(ctx, &expectedVerifiedBatch, dbTx)
	require.NoError(t, err)

	// Step to create done, retrieve it

	actualVerifiedBatch, err := testState.GetVerifiedBatch(ctx, 1, dbTx)
	require.NoError(t, err)
	require.Equal(t, expectedVerifiedBatch, *actualVerifiedBatch)

	require.NoError(t, dbTx.Commit(ctx))
}

func TestAddAccumulatedInputHash(t *testing.T) {
	initOrResetDB()

	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, dbTx)
	assert.NoError(t, err)

	_, err = testState.Exec(ctx, `INSERT INTO state.batch
	(batch_num, global_exit_root, local_exit_root, state_root, timestamp, coinbase, raw_txs_data)
	VALUES(1, '0x0000000000000000000000000000000000000000000000000000000000000000', '0x0000000000000000000000000000000000000000000000000000000000000000', '0xbf34f9a52a63229e90d1016011655bc12140bba5b771817b88cbf340d08dcbde', '2022-12-19 08:17:45.000', '0x0000000000000000000000000000000000000000', NULL);
	`)
	require.NoError(t, err)

	accInputHash := common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f2")
	batchNum := uint64(1)
	err = testState.AddAccumulatedInputHash(ctx, batchNum, accInputHash, dbTx)
	require.NoError(t, err)

	b, err := testState.GetBatchByNumber(ctx, batchNum, dbTx)
	require.NoError(t, err)
	assert.Equal(t, b.BatchNumber, batchNum)
	assert.Equal(t, b.AccInputHash, accInputHash)
	require.NoError(t, dbTx.Commit(ctx))
}

func TestForcedBatch(t *testing.T) {
	// Init database instance
	initOrResetDB()

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
	rtx := "29e885edaf8e4b51e1d2e05f9da28000000000000000000000000000000000000000000000000000000161d2fb4f6b1d53827d9b80a23cf2d7d9f1"
	raw, err := hex.DecodeString(rtx)
	assert.NoError(t, err)
	forcedBatch := state.ForcedBatch{
		BlockNumber:       1,
		ForcedBatchNumber: 1,
		Sequencer:         common.HexToAddress("0x2536C2745Ac4A584656A830f7bdCd329c94e8F30"),
		RawTxsData:        raw,
		ForcedAt:          time.Now(),
		GlobalExitRoot:    common.HexToHash("0x40a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0"),
	}
	err = testState.AddForcedBatch(ctx, &forcedBatch, tx)
	require.NoError(t, err)
	fb, err := testState.GetForcedBatch(ctx, 1, tx)
	require.NoError(t, err)
	err = tx.Commit(ctx)
	require.NoError(t, err)
	assert.Equal(t, forcedBatch.BlockNumber, fb.BlockNumber)
	assert.Equal(t, forcedBatch.ForcedBatchNumber, fb.ForcedBatchNumber)
	assert.Equal(t, forcedBatch.Sequencer, fb.Sequencer)
	assert.Equal(t, forcedBatch.RawTxsData, fb.RawTxsData)
	assert.Equal(t, rtx, common.Bytes2Hex(fb.RawTxsData))
	assert.Equal(t, forcedBatch.ForcedAt.Unix(), fb.ForcedAt.Unix())
	assert.Equal(t, forcedBatch.GlobalExitRoot, fb.GlobalExitRoot)
}
func TestCleanupLockedProofs(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	initOrResetDB()
	ctx := context.Background()
	batchNumber := uint64(42)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES ($1), ($2), ($3)", batchNumber, batchNumber+1, batchNumber+2)
	require.NoError(err)
	const addGeneratedProofSQL = "INSERT INTO state.proof (batch_num, batch_num_final, proof, proof_id, input_prover, prover, prover_id, generating_since, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	// proof with `generating_since` older than interval
	now := time.Now().Round(time.Microsecond)
	oneHourAgo := now.Add(-time.Hour).Round(time.Microsecond)
	olderProofID := "olderProofID"
	olderProof := state.Proof{
		ProofID:          &olderProofID,
		BatchNumber:      batchNumber,
		BatchNumberFinal: batchNumber,
		GeneratingSince:  &oneHourAgo,
	}
	_, err := testState.Exec(ctx, addGeneratedProofSQL, olderProof.BatchNumber, olderProof.BatchNumberFinal, olderProof.Proof, olderProof.ProofID, olderProof.InputProver, olderProof.Prover, olderProof.ProverID, olderProof.GeneratingSince, oneHourAgo, oneHourAgo)
	require.NoError(err)
	// proof with `generating_since` newer than interval
	newerProofID := "newerProofID"
	newerProof := state.Proof{
		ProofID:          &newerProofID,
		BatchNumber:      batchNumber + 1,
		BatchNumberFinal: batchNumber + 1,
		GeneratingSince:  &now,
		CreatedAt:        oneHourAgo,
		UpdatedAt:        now,
	}
	_, err = testState.Exec(ctx, addGeneratedProofSQL, newerProof.BatchNumber, newerProof.BatchNumberFinal, newerProof.Proof, newerProof.ProofID, newerProof.InputProver, newerProof.Prover, newerProof.ProverID, newerProof.GeneratingSince, oneHourAgo, now)
	require.NoError(err)
	// proof with `generating_since` nil (currently not generating)
	olderNotGenProofID := "olderNotGenProofID"
	olderNotGenProof := state.Proof{
		ProofID:          &olderNotGenProofID,
		BatchNumber:      batchNumber + 2,
		BatchNumberFinal: batchNumber + 2,
		CreatedAt:        oneHourAgo,
		UpdatedAt:        oneHourAgo,
	}
	_, err = testState.Exec(ctx, addGeneratedProofSQL, olderNotGenProof.BatchNumber, olderNotGenProof.BatchNumberFinal, olderNotGenProof.Proof, olderNotGenProof.ProofID, olderNotGenProof.InputProver, olderNotGenProof.Prover, olderNotGenProof.ProverID, olderNotGenProof.GeneratingSince, oneHourAgo, oneHourAgo)
	require.NoError(err)

	_, err = testState.CleanupLockedProofs(ctx, "1m", nil)

	require.NoError(err)
	rows, err := testState.Query(ctx, "SELECT batch_num, batch_num_final, proof, proof_id, input_prover, prover, prover_id, generating_since, created_at, updated_at FROM state.proof")
	require.NoError(err)
	proofs := make([]state.Proof, 0, len(rows.RawValues()))
	for rows.Next() {
		var proof state.Proof
		err := rows.Scan(
			&proof.BatchNumber,
			&proof.BatchNumberFinal,
			&proof.Proof,
			&proof.ProofID,
			&proof.InputProver,
			&proof.Prover,
			&proof.ProverID,
			&proof.GeneratingSince,
			&proof.CreatedAt,
			&proof.UpdatedAt,
		)
		require.NoError(err)
		proofs = append(proofs, proof)
	}
	assert.Len(proofs, 2)
	assert.Contains(proofs, olderNotGenProof)
	assert.Contains(proofs, newerProof)
}

func TestVirtualBatch(t *testing.T) {
	initOrResetDB()

	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, dbTx)
	assert.NoError(t, err)
	//require.NoError(t, tx.Commit(ctx))

	lastBlock, err := testState.GetLastBlock(ctx, dbTx)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), lastBlock.BlockNumber)

	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (1)")

	require.NoError(t, err)
	addr := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	virtualBatch := state.VirtualBatch{
		BlockNumber:   1,
		BatchNumber:   1,
		Coinbase:      addr,
		SequencerAddr: addr,
		TxHash:        common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
	}
	err = testState.AddVirtualBatch(ctx, &virtualBatch, dbTx)
	require.NoError(t, err)

	actualVirtualBatch, err := testState.GetVirtualBatch(ctx, 1, dbTx)
	require.NoError(t, err)
	require.Equal(t, virtualBatch, *actualVirtualBatch)
	require.NoError(t, dbTx.Commit(ctx))
}

func TestForkIDs(t *testing.T) {
	initOrResetDB()

	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	block1 := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f0"),
		ReceivedAt:  time.Now(),
	}
	block2 := &state.Block{
		BlockNumber: 2,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f2"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block1, dbTx)
	assert.NoError(t, err)
	err = testState.AddBlock(ctx, block2, dbTx)
	assert.NoError(t, err)

	forkID1 := state.ForkIDInterval{
		FromBatchNumber: 0,
		ToBatchNumber:   10,
		ForkId:          1,
		Version:         "version 1",
		BlockNumber:     1,
	}
	forkID2 := state.ForkIDInterval{
		FromBatchNumber: 11,
		ToBatchNumber:   20,
		ForkId:          2,
		Version:         "version 2",
		BlockNumber:     1,
	}
	forkID3 := state.ForkIDInterval{
		FromBatchNumber: 21,
		ToBatchNumber:   100,
		ForkId:          3,
		Version:         "version 3",
		BlockNumber:     2,
	}
	forks := []state.ForkIDInterval{forkID1, forkID2, forkID3}
	for _, fork := range forks {
		err = testState.AddForkID(ctx, fork, dbTx)
		require.NoError(t, err)
		// Insert twice to test on conflict do nothing
		err = testState.AddForkID(ctx, fork, dbTx)
		require.NoError(t, err)
	}

	forkIDs, err := testState.GetForkIDs(ctx, dbTx)
	require.NoError(t, err)
	require.Equal(t, 3, len(forkIDs))
	for i, forkId := range forkIDs {
		require.Equal(t, forks[i].BlockNumber, forkId.BlockNumber)
		require.Equal(t, forks[i].ForkId, forkId.ForkId)
		require.Equal(t, forks[i].FromBatchNumber, forkId.FromBatchNumber)
		require.Equal(t, forks[i].ToBatchNumber, forkId.ToBatchNumber)
		require.Equal(t, forks[i].Version, forkId.Version)
	}
	forkID3.ToBatchNumber = 18446744073709551615
	err = testState.UpdateForkID(ctx, forkID3, dbTx)
	require.NoError(t, err)

	forkIDs, err = testState.GetForkIDs(ctx, dbTx)
	require.NoError(t, err)
	require.Equal(t, 3, len(forkIDs))
	require.Equal(t, forkID3.ToBatchNumber, forkIDs[len(forkIDs)-1].ToBatchNumber)
	require.Equal(t, forkID3.ForkId, forkIDs[len(forkIDs)-1].ForkId)

	forkID3.BlockNumber = 101
	err = testState.AddForkID(ctx, forkID3, dbTx)
	require.NoError(t, err)
	forkIDs, err = testState.GetForkIDs(ctx, dbTx)
	require.NoError(t, err)
	require.Equal(t, 3, len(forkIDs))
	require.Equal(t, forkID3.ToBatchNumber, forkIDs[len(forkIDs)-1].ToBatchNumber)
	require.Equal(t, forkID3.ForkId, forkIDs[len(forkIDs)-1].ForkId)
	require.Equal(t, forkID3.BlockNumber, forkIDs[len(forkIDs)-1].BlockNumber)

	forkID3.BlockNumber = 2
	err = testState.AddForkID(ctx, forkID3, dbTx)
	require.NoError(t, err)
	forkIDs, err = testState.GetForkIDs(ctx, dbTx)
	require.NoError(t, err)
	require.Equal(t, 3, len(forkIDs))
	require.Equal(t, forkID3.ToBatchNumber, forkIDs[len(forkIDs)-1].ToBatchNumber)
	require.Equal(t, forkID3.ForkId, forkIDs[len(forkIDs)-1].ForkId)
	require.Equal(t, forkID3.BlockNumber, forkIDs[len(forkIDs)-1].BlockNumber)

	require.NoError(t, dbTx.Commit(ctx))
}

func TestGetLastVerifiedL2BlockNumberUntilL1Block(t *testing.T) {
	initOrResetDB()
	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	defer func() { require.NoError(t, dbTx.Commit(ctx)) }()

	// prepare data
	addr := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	hash := common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1")
	for i := 1; i <= 10; i++ {
		blockNumber := uint64(i)

		// add l1 block
		err = testState.AddBlock(ctx, state.NewBlock(blockNumber), dbTx)
		require.NoError(t, err)

		batchNumber := uint64(i * 10)

		// add batch
		_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES ($1)", batchNumber)
		require.NoError(t, err)

		// add l2 block
		l2Block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(0).SetUint64(blockNumber + uint64(10))})

		storeTxsEGPData := []state.StoreTxEGPData{}
		for range l2Block.Transactions() {
			storeTxsEGPData = append(storeTxsEGPData, state.StoreTxEGPData{EGPLog: nil, EffectivePercentage: uint8(0)})
		}

		err = testState.AddL2Block(ctx, batchNumber, l2Block, []*types.Receipt{}, storeTxsEGPData, dbTx)
		require.NoError(t, err)

		virtualBatch := state.VirtualBatch{BlockNumber: blockNumber, BatchNumber: batchNumber, Coinbase: addr, SequencerAddr: addr, TxHash: hash}
		err = testState.AddVirtualBatch(ctx, &virtualBatch, dbTx)
		require.NoError(t, err)

		verifiedBatch := state.VerifiedBatch{BlockNumber: blockNumber, BatchNumber: batchNumber, TxHash: hash}
		err = testState.AddVerifiedBatch(ctx, &verifiedBatch, dbTx)
		require.NoError(t, err)
	}

	type testCase struct {
		name                string
		l1BlockNumber       uint64
		expectedBatchNumber uint64
	}

	testCases := []testCase{
		{name: "l1 block number smaller than block number for the last verified batch", l1BlockNumber: 1, expectedBatchNumber: 11},
		{name: "l1 block number equal to block number for the last verified batch", l1BlockNumber: 10, expectedBatchNumber: 20},
		{name: "l1 block number bigger than number for the last verified batch", l1BlockNumber: 20, expectedBatchNumber: 20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			batchNumber, err := testState.GetLastVerifiedL2BlockNumberUntilL1Block(ctx, uint64(tc.l1BlockNumber), dbTx)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedBatchNumber, batchNumber)
		})
	}
}

func TestGetLastVerifiedBatchNumberUntilL1Block(t *testing.T) {
	initOrResetDB()
	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	defer func() { require.NoError(t, dbTx.Commit(ctx)) }()

	// prepare data
	addr := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	hash := common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1")
	for i := 1; i <= 10; i++ {
		blockNumber := uint64(i)

		// add l1 block
		err = testState.AddBlock(ctx, state.NewBlock(blockNumber), dbTx)
		require.NoError(t, err)

		batchNumber := uint64(i * 10)

		// add batch
		_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES ($1)", batchNumber)
		require.NoError(t, err)

		virtualBatch := state.VirtualBatch{BlockNumber: blockNumber, BatchNumber: batchNumber, Coinbase: addr, SequencerAddr: addr, TxHash: hash}
		err = testState.AddVirtualBatch(ctx, &virtualBatch, dbTx)
		require.NoError(t, err)

		verifiedBatch := state.VerifiedBatch{BlockNumber: blockNumber, BatchNumber: batchNumber, TxHash: hash}
		err = testState.AddVerifiedBatch(ctx, &verifiedBatch, dbTx)
		require.NoError(t, err)
	}

	type testCase struct {
		name                string
		l1BlockNumber       uint64
		expectedBatchNumber uint64
	}

	testCases := []testCase{
		{name: "l1 block number smaller than block number for the last verified batch", l1BlockNumber: 1, expectedBatchNumber: 10},
		{name: "l1 block number equal to block number for the last verified batch", l1BlockNumber: 10, expectedBatchNumber: 100},
		{name: "l1 block number bigger than number for the last verified batch", l1BlockNumber: 20, expectedBatchNumber: 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			batchNumber, err := testState.GetLastVerifiedBatchNumberUntilL1Block(ctx, uint64(tc.l1BlockNumber), dbTx)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedBatchNumber, batchNumber)
		})
	}
}

func TestSyncInfo(t *testing.T) {
	// Init database instance
	initOrResetDB()

	ctx := context.Background()
	tx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	// Test update on conflict
	err = testState.SetInitSyncBatch(ctx, 1, tx)
	require.NoError(t, err)
	err = testState.SetInitSyncBatch(ctx, 1, tx)
	require.NoError(t, err)
	err = testState.SetLastBatchInfoSeenOnEthereum(ctx, 10, 8, tx)
	require.NoError(t, err)
	err = testState.SetInitSyncBatch(ctx, 1, tx)
	require.NoError(t, err)
	err = testState.SetLastBatchInfoSeenOnEthereum(ctx, 10, 8, tx)
	require.NoError(t, err)
	err = testState.SetLastBatchInfoSeenOnEthereum(ctx, 10, 8, tx)
	require.NoError(t, err)

	err = tx.Commit(ctx)
	require.NoError(t, err)
}

func TestGetBatchByNumber(t *testing.T) {
	initOrResetDB()

	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	_, err = testState.Exec(ctx, `INSERT INTO state.batch
	(batch_num, global_exit_root, local_exit_root, state_root, timestamp, coinbase, raw_txs_data)
	VALUES(1, '0x0000000000000000000000000000000000000000000000000000000000000000', '0x0000000000000000000000000000000000000000000000000000000000000000', '0xbf34f9a52a63229e90d1016011655bc12140bba5b771817b88cbf340d08dcbde', '2022-12-19 08:17:45.000', '0x0000000000000000000000000000000000000000', NULL);
	`)
	require.NoError(t, err)

	batchNum := uint64(1)
	b, err := testState.GetBatchByNumber(ctx, batchNum, dbTx)
	require.NoError(t, err)
	assert.Equal(t, b.BatchNumber, batchNum)

	batchNum = uint64(2)
	b, err = testState.GetBatchByNumber(ctx, batchNum, dbTx)
	require.Error(t, state.ErrNotFound, err)
	assert.Nil(t, b)

	require.NoError(t, dbTx.Commit(ctx))
}

func TestGetLogs(t *testing.T) {
	initOrResetDB()

	ctx := context.Background()

	cfg := state.Config{
		MaxLogsCount:      8,
		MaxLogsBlockRange: 10,
	}

	testState = state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(cfg, stateDb), executorClient, stateTree, nil)

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	err = testState.AddBlock(ctx, block, dbTx)
	assert.NoError(t, err)

	batchNumber := uint64(1)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES ($1)", batchNumber)
	assert.NoError(t, err)

	time := time.Now()
	blockNumber := big.NewInt(1)

	for i := 0; i < 3; i++ {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    uint64(i),
			To:       nil,
			Value:    new(big.Int),
			Gas:      0,
			GasPrice: big.NewInt(0),
		})

		logs := []*types.Log{}
		for j := 0; j < 4; j++ {
			logs = append(logs, &types.Log{TxHash: tx.Hash(), Index: uint(j)})
		}

		receipt := &types.Receipt{
			Type:              uint8(tx.Type()),
			PostState:         state.ZeroHash.Bytes(),
			CumulativeGasUsed: 0,
			EffectiveGasPrice: big.NewInt(0),
			BlockNumber:       blockNumber,
			GasUsed:           tx.Gas(),
			TxHash:            tx.Hash(),
			TransactionIndex:  0,
			Status:            types.ReceiptStatusSuccessful,
			Logs:              logs,
		}

		transactions := []*types.Transaction{tx}
		receipts := []*types.Receipt{receipt}

		header := &types.Header{
			Number:     big.NewInt(int64(i) + 1),
			ParentHash: state.ZeroHash,
			Coinbase:   state.ZeroAddress,
			Root:       state.ZeroHash,
			GasUsed:    1,
			GasLimit:   10,
			Time:       uint64(time.Unix()),
		}

		l2Block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
		for _, receipt := range receipts {
			receipt.BlockHash = l2Block.Hash()
		}

		storeTxsEGPData := []state.StoreTxEGPData{}
		for range transactions {
			storeTxsEGPData = append(storeTxsEGPData, state.StoreTxEGPData{EGPLog: nil, EffectivePercentage: state.MaxEffectivePercentage})
		}

		err = testState.AddL2Block(ctx, batchNumber, l2Block, receipts, storeTxsEGPData, dbTx)
		require.NoError(t, err)
	}

	type testCase struct {
		name          string
		from          uint64
		to            uint64
		logCount      int
		expectedError error
	}

	testCases := []testCase{
		{
			name:          "invalid block range",
			from:          2,
			to:            1,
			logCount:      0,
			expectedError: state.ErrInvalidBlockRange,
		},
		{
			name:          "block range bigger than allowed",
			from:          1,
			to:            12,
			logCount:      0,
			expectedError: state.ErrMaxLogsBlockRangeLimitExceeded,
		},
		{
			name:          "log count bigger than allowed",
			from:          1,
			to:            3,
			logCount:      0,
			expectedError: state.ErrMaxLogsCountLimitExceeded,
		},
		{
			name:          "logs returned successfully",
			from:          1,
			to:            2,
			logCount:      8,
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			logs, err := testState.GetLogs(ctx, testCase.from, testCase.to, []common.Address{}, [][]common.Hash{}, nil, nil, dbTx)

			assert.Equal(t, testCase.logCount, len(logs))
			assert.Equal(t, testCase.expectedError, err)
		})
	}
	require.NoError(t, dbTx.Commit(ctx))
}

func TestGetNativeBlockHashesInRange(t *testing.T) {
	initOrResetDB()

	ctx := context.Background()

	cfg := state.Config{
		MaxNativeBlockHashBlockRange: 10,
	}

	testState = state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(cfg, stateDb), executorClient, stateTree, nil)

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	err = testState.AddBlock(ctx, block, dbTx)
	assert.NoError(t, err)

	batchNumber := uint64(1)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES ($1)", batchNumber)
	assert.NoError(t, err)

	time := time.Now()
	blockNumber := big.NewInt(1)

	nativeBlockHashes := []common.Hash{}

	for i := 0; i < 10; i++ {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    uint64(i),
			To:       nil,
			Value:    new(big.Int),
			Gas:      0,
			GasPrice: big.NewInt(0),
		})

		receipt := &types.Receipt{
			Type:              uint8(tx.Type()),
			PostState:         state.ZeroHash.Bytes(),
			CumulativeGasUsed: 0,
			EffectiveGasPrice: big.NewInt(0),
			BlockNumber:       blockNumber,
			GasUsed:           tx.Gas(),
			TxHash:            tx.Hash(),
			TransactionIndex:  0,
			Status:            types.ReceiptStatusSuccessful,
		}

		transactions := []*types.Transaction{tx}
		receipts := []*types.Receipt{receipt}

		header := &types.Header{
			Number:     big.NewInt(int64(i) + 1),
			ParentHash: state.ZeroHash,
			Coinbase:   state.ZeroAddress,
			Root:       common.HexToHash(hex.EncodeBig(big.NewInt(int64(i)))),
			GasUsed:    1,
			GasLimit:   10,
			Time:       uint64(time.Unix()),
		}

		l2Block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
		for _, receipt := range receipts {
			receipt.BlockHash = l2Block.Hash()
		}

		storeTxsEGPData := []state.StoreTxEGPData{}
		for range transactions {
			storeTxsEGPData = append(storeTxsEGPData, state.StoreTxEGPData{EGPLog: nil, EffectivePercentage: state.MaxEffectivePercentage})
		}

		err = testState.AddL2Block(ctx, batchNumber, l2Block, receipts, storeTxsEGPData, dbTx)
		require.NoError(t, err)

		nativeBlockHashes = append(nativeBlockHashes, l2Block.Header().Root)
	}

	type testCase struct {
		name            string
		from            uint64
		to              uint64
		expectedResults []common.Hash
		expectedError   error
	}

	testCases := []testCase{
		{
			name:            "invalid block range",
			from:            2,
			to:              1,
			expectedResults: nil,
			expectedError:   state.ErrInvalidBlockRange,
		},
		{
			name:            "block range bigger than allowed",
			from:            1,
			to:              12,
			expectedResults: nil,
			expectedError:   state.ErrMaxNativeBlockHashBlockRangeLimitExceeded,
		},
		{
			name:            "hashes returned successfully",
			from:            4,
			to:              7,
			expectedResults: nativeBlockHashes[3:7],
			expectedError:   nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			results, err := testState.GetNativeBlockHashesInRange(ctx, testCase.from, testCase.to, dbTx)

			assert.ElementsMatch(t, testCase.expectedResults, results)
			assert.Equal(t, testCase.expectedError, err)
		})
	}

	require.NoError(t, dbTx.Commit(ctx))
}
