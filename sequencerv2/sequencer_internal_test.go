package sequencerv2

import (
	"context"
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
	cfgTypes "github.com/hermeznetwork/hermez-core/config/types"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	ethmanTypes "github.com/hermeznetwork/hermez-core/ethermanv2/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/pool/pgpoolstorage"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type stateTestInterface interface {
	stateInterface
	GetNonce(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) (uint64, error)
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) (*big.Int, error)
	NewGenesisBatchProcessor(genesisStateRoot []byte, txBundleID string) (*state.BatchProcessor, error)
	SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64, txBundleID string) error
	SetGenesis(ctx context.Context, genesis state.Genesis, txBundleID string) error
}

var (
	dbCfg     = dbutils.NewConfigFromEnv()
	stateDB   *pgxpool.Pool
	testState stateTestInterface
	seqCfg    Config
	pl        *pool.Pool

	genesisHash common.Hash
	txs         []*types.Transaction

	senderPrivateKey           = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	addr                       = common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	consolidatedTxHash         = common.HexToHash("0x125714bb4db48757007fff2671b37637bbfd6d47b3a4757ebbd0c5222984f905")
	maticCollateral            = big.NewInt(1000000000000000000)
	lastBatchNumberSeen uint64 = 1
)

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
	bp, err := testState.NewGenesisBatchProcessor(nil, "")
	if err != nil {
		require.NoError(t, err)
	}

	err = bp.ProcessBatch(ctx, batch)
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
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDB)
	testState = state.NewState(stateCfg, state.NewPostgresStorage(stateDB), tree.NewStateTree(mt, scCodeStore))

	seqCfg = Config{WaitPeriodPoolIsEmpty: cfgTypes.NewDuration(time.Second)}

	s, err := pgpoolstorage.NewPostgresPoolStorage(dbCfg)
	if err != nil {
		panic(err)
	}
	pl = pool.NewPool(s, testState, stateCfg.L2GlobalExitRootManagerAddr)

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

	err = testState.SetLastBatchNumberSeenOnEthereum(context.Background(), lastBatchNumberSeen, "")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	_, err = stateDB.Exec(ctx, "DELETE FROM state.block")
	if err != nil {
		panic(err)
	}

	err = testState.SetGenesis(ctx, genesis, "")
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

	for i := 0; i < 25; i++ {
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

func TestSequencerBaseFlow(t *testing.T) {
	txManager := new(txmanagerMock)

	ctx := context.Background()
	setUpBlock(ctx, t)
	setUpBatch(ctx, t)

	seq, err := New(seqCfg, pl, testState, txManager)
	require.NoError(t, err)
	ticker := time.NewTicker(seqCfg.WaitPeriodPoolIsEmpty.Duration)

	pendTxs, err := pl.GetPendingTxs(ctx, false, 30)
	require.NoError(t, err)
	require.Equal(t, 25, len(pendTxs))
	for i := 0; i < maxTxsInSequence; i++ {
		seq.tryToProcessTx(ctx, ticker)
	}
	sequencesToSent := make([]ethmanTypes.Sequence, 5)
	for i := 0; i < maxSequencesLength; i++ {
		for k := maxSequencesLength * i; k < maxSequencesLength*(i+1); k++ {
			sequencesToSent[i].Txs = append(sequencesToSent[i].Txs, pendTxs[k].Transaction)
		}
	}
	txManager.On("SequenceBatches", mock.MatchedBy(func(sequences []ethmanTypes.Sequence) bool {
		res := true
		for i := 0; i < len(sequences); i++ {
			for k := 0; k < len(sequences[i].Txs); k++ {
				res = res && sequences[i].Txs[k].Hash() == sequencesToSent[i].Txs[k].Hash()
			}
		}
		return res
	})).Return(nil, nil)
	pendTxs, err = pl.GetPendingTxs(ctx, false, 30)
	require.NoError(t, err)
	require.Equal(t, 20, len(pendTxs))
	require.Equal(t, maxTxsInSequence, len(seq.sequenceInProgress.Txs))
	require.Equal(t, 0, len(seq.sequencesInProgress))

	seq.tryToProcessTx(ctx, ticker)
	pendTxs, err = pl.GetPendingTxs(ctx, false, 30)
	require.NoError(t, err)
	require.Equal(t, 19, len(pendTxs))
	require.Equal(t, 1, len(seq.sequenceInProgress.Txs))
	require.Equal(t, 1, len(seq.sequencesInProgress))
	require.Equal(t, maxTxsInSequence, len(seq.sequencesInProgress[0].Txs))

	for i := 0; i < 20; i++ {
		seq.tryToProcessTx(ctx, ticker)
	}
	pendTxs, err = pl.GetPendingTxs(ctx, false, 30)
	require.NoError(t, err)
	require.Equal(t, 0, len(pendTxs))
	require.Equal(t, 0, len(seq.sequenceInProgress.Txs))
	require.Equal(t, 0, len(seq.sequencesInProgress))

	txManager.AssertNumberOfCalls(t, "SequenceBatches", 1)
	cleanUpBatches(ctx, t)
	cleanUpBlocks(ctx, t)
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
