package txprofitabilitychecker_test

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pricegetter"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txprofitabilitychecker"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBatch(ctx context.Context, isVirtual bool) (*state.Batch, error)
	NewGenesisBatchProcessor(genesisStateRoot []byte) (*state.BatchProcessor, error)
}

var (
	stateDB     *pgxpool.Pool
	testState   stateInterface
	priceGetter pricegetter.Client

	addr               common.Address = common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	consolidatedTxHash common.Hash    = common.HexToHash("0x125714bb4db48757007fff2671b37637bbfd6d47b3a4757ebbd0c5222984f905")
	maticCollateral                   = big.NewInt(1000000000000000000)
	txs                []*types.Transaction
	maticAmount        = big.NewInt(1000000000000000001)
)
var dbCfg = dbutils.NewConfigFromEnv()

var stateCfg = state.Config{
	DefaultChainID:       1000,
	MaxCumulativeGasUsed: 800000,
}

func TestMain(m *testing.M) {
	var err error

	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})

	if err := dbutils.InitOrReset(dbCfg); err != nil {
		panic(err)
	}

	stateDB, err = db.NewSQLDB(dbCfg)
	if err != nil {
		panic(err)
	}
	defer stateDB.Close()
	store := tree.NewPostgresStore(stateDB)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDB)
	testState = state.NewState(stateCfg, state.NewPostgresStorage(stateDB), tree.NewStateTree(mt, scCodeStore))
	tx := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	txs = []*types.Transaction{tx}

	defaultPrice := new(pricegetter.TokenPrice)
	_ = defaultPrice.UnmarshalText([]byte("2000"))
	priceGetter, err = pricegetter.NewClient(pricegetter.Config{
		Type:         pricegetter.DefaultType,
		DefaultPrice: *defaultPrice,
	})
	if err != nil {
		panic(err)
	}

	setUpBlock()
	setUpBatch()
	result := m.Run()
	os.Exit(result)
}

func setUpBlock() {
	var err error
	ctx := context.Background()
	hash1 := common.HexToHash("0x65b4699dda5f7eb4519c730e6a48e73c90d2b1c8efcd6a6abdfd28c3b8e7d7d9")

	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   hash1,
		ParentHash:  hash1,
		ReceivedAt:  time.Now(),
	}

	_, err = stateDB.Exec(ctx, "DELETE FROM state.block")
	if err != nil {
		panic(err)
	}

	_, err = stateDB.Exec(ctx, "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)",
		block.BlockNumber, block.BlockHash.Bytes(), block.ParentHash.Bytes(), block.ReceivedAt)
	if err != nil {
		panic(err)
	}
}

func setUpBatch() {
	var err error
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
	ctx := context.Background()
	_, err = stateDB.Exec(ctx, "DELETE FROM state.batch")
	if err != nil {
		panic(err)
	}

	bp, err := testState.NewGenesisBatchProcessor(nil)
	if err != nil {
		panic(err)
	}

	err = bp.ProcessBatch(ctx, batch)
	if err != nil {
		panic(err)
	}
}

func TestBase_IsProfitable_FailByMinReward(t *testing.T) {
	minReward := new(big.Int).Mul(big.NewInt(1000), big.NewInt(encoding.TenToThePowerOf18))
	ethMan := new(etherman)
	txProfitabilityChecker := txprofitabilitychecker.NewTxProfitabilityCheckerBase(ethMan, testState, priceGetter, minReward, time.Duration(60), 50)
	ctx := context.Background()

	ethMan.On("EstimateSendBatchCost", ctx, txs, maticAmount).Return(big.NewInt(1), nil)
	isProfitable, _, err := txProfitabilityChecker.IsProfitable(ctx, txs)
	ethMan.AssertExpectations(t)
	assert.NoError(t, err)
	assert.False(t, isProfitable)
}

func TestBase_IsProfitable_SendBatchAnyway(t *testing.T) {
	minReward := big.NewInt(0)
	ethMan := new(etherman)
	txProfitabilityChecker := txprofitabilitychecker.NewTxProfitabilityCheckerBase(ethMan, testState, priceGetter, minReward, time.Duration(1), 50)

	ctx := context.Background()

	ethMan.On("EstimateSendBatchCost", ctx, txs, maticAmount).Return(maticAmount, nil)

	isProfitable, reward, err := txProfitabilityChecker.IsProfitable(ctx, txs)
	ethMan.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 0, reward.Cmp(minReward))
	assert.True(t, isProfitable)
}

func TestBase_IsProfitable_GasCostTooBigForSendingTx(t *testing.T) {
	minReward := big.NewInt(0)
	ethMan := new(etherman)
	txProfitabilityChecker := txprofitabilitychecker.NewTxProfitabilityCheckerBase(ethMan, testState, priceGetter, minReward, time.Duration(60), 50)

	ctx := context.Background()

	ethMan.On("EstimateSendBatchCost", ctx, txs, maticAmount).Return(maticAmount, nil)
	isProfitable, _, err := txProfitabilityChecker.IsProfitable(ctx, txs)
	ethMan.AssertExpectations(t)
	assert.NoError(t, err)
	assert.False(t, isProfitable)
}

func TestBase_IsProfitable(t *testing.T) {
	ethMan := new(etherman)
	txProfitabilityChecker := txprofitabilitychecker.NewTxProfitabilityCheckerBase(ethMan, testState, priceGetter, big.NewInt(0), time.Duration(60), 50)

	ctx := context.Background()

	ethMan.On("EstimateSendBatchCost", ctx, txs, maticAmount).Return(big.NewInt(10), nil)
	ethMan.On("GetCurrentSequencerCollateral").Return(big.NewInt(1), nil)

	isProfitable, reward, err := txProfitabilityChecker.IsProfitable(ctx, txs)
	ethMan.AssertExpectations(t)
	assert.NoError(t, err)
	assert.True(t, isProfitable)
	// gasCostForSendingBatch is 10, tx cost is 20, reward will be 10 eth, aggregator takes 50%, so his reward is 5 eth = 10000 matic
	assert.Equal(t, 0, big.NewInt(10000).Cmp(reward))
}
