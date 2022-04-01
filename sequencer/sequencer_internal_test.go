package sequencer

import (
	"context"
	"errors"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/pool/pgpoolstorage"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txprofitabilitychecker"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txselector"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

type stateTestInterface interface {
	stateInterface
	// following methods used for tests
	NewGenesisBatchProcessor(genesisStateRoot []byte) (*state.BasicBatchProcessor, error)
	GetNonce(ctx context.Context, address common.Address, batchNumber uint64) (uint64, error)
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64) (*big.Int, error)
	SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64) error
	SetGenesis(ctx context.Context, genesis state.Genesis) error
}

var (
	stateDB   *pgxpool.Pool
	testState stateTestInterface
	seqCfg    Config
	pl        *pool.Pool

	genesisHash common.Hash
	txs         []*types.Transaction

	addr                       = common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	consolidatedTxHash         = common.HexToHash("0x125714bb4db48757007fff2671b37637bbfd6d47b3a4757ebbd0c5222984f905")
	maticCollateral            = big.NewInt(1000000000000000000)
	maticAmount                = big.NewInt(1000000000000000001)
	lastBatchNumberSeen uint64 = 1
	senderPrivateKey           = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
)

var dbCfg = dbutils.NewConfigFromEnv()

var stateCfg = state.Config{
	DefaultChainID:       1000,
	MaxCumulativeGasUsed: 800000,
}

func setUpBlock(ctx context.Context, t *testing.T) {
	blockHash := common.HexToHash("0x65b4699dda5f7eb4519c730e6a48e73c90d2b1c8efcd6a6abdfd28c3b8e7d7d9")

	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   blockHash,
		ParentHash:  genesisHash,
		ReceivedAt:  time.Now(),
	}

	_, err := stateDB.Exec(ctx, "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)",
		block.BlockNumber, block.BlockHash.Bytes(), block.ParentHash.Bytes(), block.ReceivedAt)
	if err != nil {
		require.NoError(t, err)
	}
}

func setUpBatch(ctx context.Context, t *testing.T) {
	receivedAt := time.Now().Add(time.Duration(-5) * time.Minute)
	consolidatedAt := time.Now()
	batch := &state.Batch{
		BlockNumber:        1,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: consolidatedTxHash,
		Header:             &types.Header{Number: big.NewInt(1)},
		Uncles:             nil,
		RawTxsData:         nil,
		MaticCollateral:    maticCollateral,
		ReceivedAt:         receivedAt,
		ConsolidatedAt:     &consolidatedAt,
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.Hash{},
	}
	bp, err := testState.NewGenesisBatchProcessor(nil)
	if err != nil {
		require.NoError(t, err)
	}

	err = bp.ProcessBatch(ctx, batch)
	if err != nil {
		require.NoError(t, err)
	}
}

func cleanUpBatches(ctx context.Context, t *testing.T) {
	_, err := stateDB.Exec(ctx, "DELETE FROM state.batch WHERE block_num = $1", 1)
	if err != nil {
		require.NoError(t, err)
	}
}

func cleanUpBlocks(ctx context.Context, t *testing.T) {
	_, err := stateDB.Exec(ctx, "DELETE FROM state.block WHERE block_num = $1", 1)
	if err != nil {
		require.NoError(t, err)
	}
}

func TestMain(m *testing.M) {
	var err error

	if err := dbutils.InitOrReset(dbCfg); err != nil {
		panic(err)
	}

	stateDB, err = db.NewSQLDB(dbCfg)
	if err != nil {
		panic(err)
	}
	defer stateDB.Close()

	store := tree.NewPostgresStore(stateDB)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDB)
	testState = state.NewState(stateCfg, state.NewPostgresStorage(stateDB), tree.NewStateTree(mt, scCodeStore))

	intervalToProposeBatch := new(Duration)
	_ = intervalToProposeBatch.UnmarshalText([]byte("5s"))
	intervalAfterWhichBatchSentAnyway := new(Duration)
	_ = intervalAfterWhichBatchSentAnyway.UnmarshalText([]byte("60s"))
	minReward := new(txprofitabilitychecker.TokenAmountWithDecimals)
	_ = minReward.UnmarshalText([]byte("1.1"))

	s, err := pgpoolstorage.NewPostgresPoolStorage(dbCfg)
	if err != nil {
		panic(err)
	}
	pl = pool.NewPool(s, testState)
	seqCfg = Config{
		IntervalToProposeBatch:            *intervalToProposeBatch,
		SyncedBlockDif:                    1,
		IntervalAfterWhichBatchSentAnyway: *intervalAfterWhichBatchSentAnyway,
		Strategy: strategy.Strategy{
			TxSelector: txselector.Config{
				Type:         "base",
				TxSorterType: "bycostandnonce",
			},
			TxProfitabilityChecker: txprofitabilitychecker.Config{
				Type:                         "acceptall",
				MinReward:                    *minReward,
				RewardPercentageToAggregator: 50,
			},
		},
		AllowNonRegistered: true,
		DefaultChainID:     1000,
	}

	err = testState.SetLastBatchNumberSeenOnEthereum(context.Background(), lastBatchNumberSeen)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	_, err = stateDB.Exec(ctx, "DELETE FROM state.block")
	if err != nil {
		panic(err)
	}

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	genesisHash = genesisBlock.Hash()
	balance, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	genesis := state.Genesis{
		Block: genesisBlock,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): balance,
		},
	}
	err = testState.SetGenesis(ctx, genesis)
	if err != nil {
		panic(err)
	}
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1000))
	if err != nil {
		panic(err)
	}

	for i := 0; i < 4; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(21000), big.NewInt(10), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		if err != nil {
			panic(err)
		}
		ctx := context.Background()
		if err := pl.AddTx(ctx, *signedTx); err != nil {
			panic(err)
		}
		txs = append(txs, signedTx)
	}

	result := m.Run()
	os.Exit(result)
}

func TestSequencerIsSynced(t *testing.T) {
	eth := new(ethermanMock)
	eth.On("GetAddress").Return(addr)
	ctx := context.Background()
	setUpBlock(ctx, t)
	setUpBatch(ctx, t)

	seq, err := NewSequencer(seqCfg, pl, testState, eth)
	require.NoError(t, err)

	synced := seq.isSynced()
	require.True(t, synced)

	cleanUpBatches(ctx, t)
	cleanUpBlocks(ctx, t)
}

func TestSequencerIsNotSynced(t *testing.T) {
	eth := new(ethermanMock)
	ctx := context.Background()

	setUpBlock(ctx, t)
	setUpBatch(ctx, t)

	eth.On("GetAddress").Return(addr)

	seq, err := NewSequencer(seqCfg, pl, testState, eth)
	require.NoError(t, err)

	err = testState.SetLastBatchNumberSeenOnEthereum(ctx, 5)
	require.NoError(t, err)

	synced := seq.isSynced()
	require.False(t, synced)

	err = testState.SetLastBatchNumberSeenOnEthereum(ctx, lastBatchNumberSeen)
	require.NoError(t, err)

	cleanUpBatches(ctx, t)
	cleanUpBlocks(ctx, t)
}

func TestSequencerGetPendingTxs(t *testing.T) {
	eth := new(ethermanMock)
	ctx := context.Background()
	eth.On("GetAddress").Return(addr)

	seq, err := NewSequencer(seqCfg, pl, testState, eth)
	require.NoError(t, err)

	for i := 0; i < 2; i++ {
		err = pl.UpdateTxState(ctx, txs[i].Hash(), pool.TxStateSelected)
		require.NoError(t, err)
	}

	pendTxs, ok := seq.getPendingTxs()
	require.True(t, ok)
	require.Equal(t, 2, len(pendTxs))

	for i := 2; i < 4; i++ {
		err = pl.UpdateTxState(ctx, txs[i].Hash(), pool.TxStateSelected)
		require.NoError(t, err)
	}

	pendTxs, ok = seq.getPendingTxs()
	require.False(t, ok)
	require.Equal(t, 0, len(pendTxs))

	setTxsToPendingState(ctx, t)
}

func TestSequencerSelectTxs(t *testing.T) {
	eth := new(ethermanMock)
	ctx := context.Background()
	eth.On("GetAddress").Return(addr)

	seq, err := NewSequencer(seqCfg, pl, testState, eth)
	require.NoError(t, err)

	txs, err := pl.GetPendingTxs(ctx, 0)
	require.NoError(t, err)
	selTxs, selTxsHashes, ok := seq.selectTxs(txs)
	require.True(t, ok)
	require.Equal(t, 4, len(selTxs))
	require.Equal(t, 4, len(selTxsHashes))
}

func TestSequencerSelectTxsInvTxs(t *testing.T) {
	eth := new(ethermanMock)
	ctx := context.Background()
	eth.On("GetAddress").Return(addr)

	seq, err := NewSequencer(seqCfg, pl, testState, eth)
	require.NoError(t, err)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	if err != nil {
		require.NoError(t, err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1000))
	if err != nil {
		require.NoError(t, err)
	}

	tx := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(5), uint64(21000), big.NewInt(10), []byte{})
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		require.NoError(t, err)
	}
	if err := pl.AddTx(ctx, *signedTx); err != nil {
		require.NoError(t, err)
	}
	txs = append(txs, signedTx)

	txs, err := pl.GetPendingTxs(ctx, 0)
	require.NoError(t, err)
	selTxs, selTxsHashes, ok := seq.selectTxs(txs)
	require.True(t, ok)
	require.Equal(t, 4, len(selTxs))
	require.Equal(t, 4, len(selTxsHashes))

	rows, err := stateDB.Query(ctx, "SELECT state FROM pool.txs WHERE hash = $1", signedTx.Hash().Hex())

	if err != nil {
		require.NoError(t, err)
	}
	defer rows.Close()

	var state string
	rows.Next()
	if err := rows.Scan(&state); err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, pool.TxStateInvalid, pool.TxState(state))

	_, err = stateDB.Exec(ctx, "DELETE FROM pool.txs WHERE hash = $1", tx.Hash().Hex())
	require.NoError(t, err)
}

func TestSequencerSendBatchEthereum(t *testing.T) {
	eth := new(ethermanMock)
	ctx := context.Background()
	eth.On("GetAddress").Return(addr)

	seq, err := NewSequencer(seqCfg, pl, testState, eth)
	require.NoError(t, err)

	txs, err := pl.GetPendingTxs(ctx, 0)
	require.NoError(t, err)
	selTxs, selTxsHashes, ok := seq.selectTxs(txs)
	require.True(t, ok)
	require.Equal(t, 4, len(selTxs))
	require.Equal(t, 4, len(selTxsHashes))

	aggrReward := big.NewInt(1)
	eth.On("EstimateSendBatchCost", ctx, txs, maticAmount).Return(big.NewInt(10), nil)
	eth.On("GetCurrentSequencerCollateral").Return(aggrReward, nil)

	tx := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(5), uint64(21000), big.NewInt(10), []byte{})

	eth.On("SendBatch", seq.ctx, selTxs, aggrReward).Return(tx, nil)

	cuttedTxs, cuttedTxsHash, ok := seq.sendBatchToEthereum(selTxs, selTxsHashes)
	require.True(t, ok)
	require.Nil(t, cuttedTxs)
	require.Nil(t, cuttedTxsHash)

	var count int
	err = stateDB.QueryRow(ctx, "SELECT COUNT(*) FROM pool.txs WHERE state = $1", pool.TxStateSelected).Scan(&count)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, 4, count)
	setTxsToPendingState(ctx, t)
}

func TestSequencerSendBatchEthereumCut(t *testing.T) {
	eth := new(ethermanMock)
	ctx := context.Background()
	eth.On("GetAddress").Return(addr)

	seq, err := NewSequencer(seqCfg, pl, testState, eth)
	require.NoError(t, err)

	txs, err := pl.GetPendingTxs(ctx, 0)
	require.NoError(t, err)
	selTxs, selTxsHashes, ok := seq.selectTxs(txs)
	require.True(t, ok)
	require.Equal(t, 4, len(selTxs))
	require.Equal(t, 4, len(selTxsHashes))
	aggrReward := big.NewInt(1)
	eth.On("EstimateSendBatchCost", ctx, txs, maticAmount).Return(big.NewInt(10), nil)
	eth.On("GetCurrentSequencerCollateral").Return(aggrReward, nil)
	eth.On("SendBatch", seq.ctx, selTxs, aggrReward).Return(nil, errors.New("gas required exceeds allowance"))

	cutTxs, cutTxsHash, ok := seq.sendBatchToEthereum(selTxs, selTxsHashes)
	require.False(t, ok)
	require.Equal(t, 3, len(cutTxs))
	require.Equal(t, 3, len(cutTxsHash))
}

func setTxsToPendingState(ctx context.Context, t *testing.T) {
	for _, tx := range txs {
		require.NoError(t, pl.UpdateTxState(ctx, tx.Hash(), pool.TxStatePending))
	}
}
