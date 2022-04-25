package state_test

import (
	"context"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/testutils"
	"github.com/hermeznetwork/hermez-core/test/vectors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

var (
	stateDb                                                *pgxpool.Pool
	testState                                              *state.State
	block1, block2                                         *state.Block
	addr                                                   common.Address = common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	hash1, hash2                                           common.Hash
	blockNumber1, blockNumber2                             uint64 = 1, 2
	batchNumber1, batchNumber2, batchNumber3, batchNumber4 uint64 = 1, 2, 3, 4
	batch1, batch2, batch3, batch4                         *state.Batch
	consolidatedTxHash                                     common.Hash = common.HexToHash("0x125714bb4db48757007fff2671b37637bbfd6d47b3a4757ebbd0c5222984f905")
	txHash                                                 common.Hash
	ctx                                                           = context.Background()
	lastBatchNumberSeen                                    uint64 = 1
	maticCollateral                                               = big.NewInt(1000000000000000000)
)

var cfg = dbutils.NewConfigFromEnv()

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

	if err := dbutils.InitOrReset(cfg); err != nil {
		panic(err)
	}

	stateDb, err = db.NewSQLDB(cfg)
	if err != nil {
		panic(err)
	}
	defer stateDb.Close()
	hash1 = common.HexToHash("0x65b4699dda5f7eb4519c730e6a48e73c90d2b1c8efcd6a6abdfd28c3b8e7d7d9")
	hash2 = common.HexToHash("0x613aabebf4fddf2ad0f034a8c73aa2f9c5a6fac3a07543023e0a6ee6f36e5795")

	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	testState = state.NewState(stateCfg, state.NewPostgresStorage(stateDb), tree.NewStateTree(mt, scCodeStore))

	setUpBlocks()
	setUpBatches()
	setUpTransactions()

	result := m.Run()

	os.Exit(result)
}

func setUpBlocks() {
	var err error
	block1 = &state.Block{
		BlockNumber: blockNumber1,
		BlockHash:   hash1,
		ParentHash:  hash1,
		ReceivedAt:  time.Now(),
	}
	block2 = &state.Block{
		BlockNumber: blockNumber2,
		BlockHash:   hash2,
		ParentHash:  hash1,
		ReceivedAt:  time.Now(),
	}

	_, err = stateDb.Exec(ctx, "DELETE FROM state.block")
	if err != nil {
		panic(err)
	}

	_, err = stateDb.Exec(ctx, "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)",
		block1.BlockNumber, block1.BlockHash.Bytes(), block1.ParentHash.Bytes(), block1.ReceivedAt)
	if err != nil {
		panic(err)
	}

	_, err = stateDb.Exec(ctx, "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)",
		block2.BlockNumber, block2.BlockHash.Bytes(), block2.ParentHash.Bytes(), block2.ReceivedAt)
	if err != nil {
		panic(err)
	}
}

func setUpBatches() {
	var err error

	batch1 = &state.Batch{
		BlockNumber:        blockNumber1,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: consolidatedTxHash,
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(batchNumber1)},
		Uncles:             nil,
		RawTxsData:         nil,
		MaticCollateral:    maticCollateral,
		ReceivedAt:         time.Now(),
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}
	batch2 = &state.Batch{
		BlockNumber:        blockNumber1,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: consolidatedTxHash,
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(batchNumber2)},
		Uncles:             nil,
		RawTxsData:         nil,
		MaticCollateral:    maticCollateral,
		ReceivedAt:         time.Now(),
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}
	batch3 = &state.Batch{
		BlockNumber:        blockNumber2,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: common.Hash{},
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(batchNumber3)},
		Uncles:             nil,
		Transactions:       nil,
		RawTxsData:         nil,
		MaticCollateral:    maticCollateral,
		ReceivedAt:         time.Now(),
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}
	batch4 = &state.Batch{
		BlockNumber:        blockNumber2,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: common.Hash{},
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(batchNumber4)},
		Uncles:             nil,
		Transactions:       nil,
		RawTxsData:         nil,
		MaticCollateral:    maticCollateral,
		ReceivedAt:         time.Now(),
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}

	_, err = stateDb.Exec(ctx, "DELETE FROM state.batch")
	if err != nil {
		panic(err)
	}

	batches := []*state.Batch{batch1, batch2, batch3, batch4}

	bp, err := testState.NewGenesisBatchProcessor(nil, "")
	if err != nil {
		panic(err)
	}

	for _, b := range batches {
		err := bp.ProcessBatch(ctx, b)
		if err != nil {
			panic(err)
		}
	}
}

func setUpTransactions() {
	tx1Inner := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	txHash = tx1Inner.Hash()
	b, err := tx1Inner.MarshalBinary()
	if err != nil {
		panic(err)
	}
	encoded := hex.EncodeToHex(b)

	b, err = tx1Inner.MarshalJSON()
	if err != nil {
		panic(err)
	}
	decoded := string(b)
	sql := "INSERT INTO state.transaction (hash, from_address, encoded, decoded, batch_num) VALUES($1, $2, $3, $4, $5)"
	if _, err := stateDb.Exec(ctx, sql, txHash, addr, encoded, decoded, batchNumber1); err != nil {
		panic(err)
	}
}

func TestBasicState_GetLastBlock(t *testing.T) {
	lastBlock, err := testState.GetLastBlock(ctx, "")
	assert.NoError(t, err)
	assert.Equal(t, block2.BlockNumber, lastBlock.BlockNumber)
}

func TestBasicState_GetPreviousBlock(t *testing.T) {
	previousBlock, err := testState.GetPreviousBlock(ctx, 1, "")
	assert.NoError(t, err)
	assert.Equal(t, block1.BlockNumber, previousBlock.BlockNumber)
}

func TestBasicState_GetBlockByHash(t *testing.T) {
	block, err := testState.GetBlockByHash(ctx, hash1, "")
	assert.NoError(t, err)
	assert.Equal(t, block1.BlockHash, block.BlockHash)
	assert.Equal(t, block1.BlockNumber, block.BlockNumber)
}

func TestBasicState_GetBlockByNumber(t *testing.T) {
	block, err := testState.GetBlockByNumber(ctx, blockNumber2, "")
	assert.NoError(t, err)
	assert.Equal(t, block2.BlockNumber, block.BlockNumber)
	assert.Equal(t, block2.BlockHash, block.BlockHash)
}

func TestBasicState_GetLastVirtualBatch(t *testing.T) {
	lastBatch, err := testState.GetLastBatch(ctx, true, "")
	assert.NoError(t, err)
	assert.Equal(t, batch4.Hash(), lastBatch.Hash())
	assert.Equal(t, batch4.Number().Uint64(), lastBatch.Number().Uint64())
}

func TestBasicState_GetLastBatch(t *testing.T) {
	lastBatch, err := testState.GetLastBatch(ctx, false, "")
	assert.NoError(t, err)
	assert.Equal(t, batch2.Hash(), lastBatch.Hash())
	assert.Equal(t, batch2.Number().Uint64(), lastBatch.Number().Uint64())
	assert.Equal(t, maticCollateral, lastBatch.MaticCollateral)
}

func TestBasicState_GetPreviousBatch(t *testing.T) {
	previousBatch, err := testState.GetPreviousBatch(ctx, false, 1, "")
	assert.NoError(t, err)
	assert.Equal(t, batch1.Hash(), previousBatch.Hash())
	assert.Equal(t, batch1.Number().Uint64(), previousBatch.Number().Uint64())
	assert.Equal(t, maticCollateral, previousBatch.MaticCollateral)
}

func TestBasicState_GetBatchByHash(t *testing.T) {
	batch, err := testState.GetBatchByHash(ctx, batch1.Hash(), "")
	assert.NoError(t, err)
	assert.Equal(t, batch1.Hash(), batch.Hash())
	assert.Equal(t, batch1.Number().Uint64(), batch.Number().Uint64())
	assert.Equal(t, maticCollateral, batch1.MaticCollateral)
}

func TestBasicState_GetBatchByNumber(t *testing.T) {
	batch, err := testState.GetBatchByNumber(ctx, batch1.Number().Uint64(), "")
	assert.NoError(t, err)
	assert.Equal(t, batch1.Number().Uint64(), batch.Number().Uint64())
	assert.Equal(t, batch1.Hash(), batch.Hash())
}

func TestBasicState_GetLastBatchNumber(t *testing.T) {
	batchNumber, err := testState.GetLastBatchNumber(ctx, "")
	assert.NoError(t, err)
	assert.Equal(t, batch4.Number().Uint64(), batchNumber)
}

func TestBasicState_ConsolidateBatch(t *testing.T) {
	batchNumber := uint64(5)
	batch := &state.Batch{
		BlockNumber:        blockNumber2,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: common.Hash{},
		Header: &types.Header{
			Number: big.NewInt(0).SetUint64(batchNumber),
		},
		Uncles:          nil,
		Transactions:    nil,
		RawTxsData:      nil,
		MaticCollateral: maticCollateral,
		ReceivedAt:      time.Now(),
		ChainID:         big.NewInt(1000),
		GlobalExitRoot:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}

	bp, err := testState.NewGenesisBatchProcessor(nil, "")
	assert.NoError(t, err)

	err = bp.ProcessBatch(ctx, batch)
	assert.NoError(t, err)

	insertedBatch, err := testState.GetBatchByNumber(ctx, batchNumber, "")
	assert.NoError(t, err)
	assert.Equal(t, common.Hash{}, insertedBatch.ConsolidatedTxHash)
	assert.Equal(t, big.NewInt(1000), insertedBatch.ChainID)
	assert.NotEqual(t, common.Hash{}, insertedBatch.GlobalExitRoot)

	err = testState.ConsolidateBatch(ctx, batchNumber, consolidatedTxHash, time.Now(), batch.Aggregator, "")
	assert.NoError(t, err)

	insertedBatch, err = testState.GetBatchByNumber(ctx, batchNumber, "")
	assert.NoError(t, err)
	assert.Equal(t, consolidatedTxHash, insertedBatch.ConsolidatedTxHash)

	_, err = stateDb.Exec(ctx, "DELETE FROM state.batch WHERE batch_num = $1", batchNumber)
	assert.NoError(t, err)
}

func TestBasicState_GetTransactionCount(t *testing.T) {
	count, err := testState.GetTransactionCount(ctx, addr, "")
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), count)
}

func TestBasicState_GetTxsByBatchNum(t *testing.T) {
	txs, err := testState.GetTxsByBatchNum(ctx, batchNumber1, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(txs))
}

func TestBasicState_GetTransactionByHash(t *testing.T) {
	tx, err := testState.GetTransactionByHash(ctx, txHash, "")
	assert.NoError(t, err)
	assert.Equal(t, txHash, tx.Hash())
}

func TestBasicState_AddBlock(t *testing.T) {
	lastBN, err := testState.GetLastBlockNumber(ctx, "")
	assert.NoError(t, err)

	block1 := &state.Block{
		BlockNumber: lastBN + 1,
		BlockHash:   hash1,
		ParentHash:  hash1,
		ReceivedAt:  time.Now(),
	}
	block2 := &state.Block{
		BlockNumber: lastBN + 2,
		BlockHash:   hash2,
		ParentHash:  hash1,
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block1, "")
	assert.NoError(t, err)
	err = testState.AddBlock(ctx, block2, "")
	assert.NoError(t, err)

	block3, err := testState.GetBlockByNumber(ctx, block1.BlockNumber, "")
	assert.NoError(t, err)
	assert.Equal(t, block1.BlockHash, block3.BlockHash)
	assert.Equal(t, block1.ParentHash, block3.ParentHash)

	block4, err := testState.GetBlockByNumber(ctx, block2.BlockNumber, "")
	assert.NoError(t, err)
	assert.Equal(t, block2.BlockHash, block4.BlockHash)
	assert.Equal(t, block2.ParentHash, block4.ParentHash)

	_, err = stateDb.Exec(ctx, "DELETE FROM state.block WHERE block_num = $1", block1.BlockNumber)
	assert.NoError(t, err)
	_, err = stateDb.Exec(ctx, "DELETE FROM state.block WHERE block_num = $1", block2.BlockNumber)
	assert.NoError(t, err)
}

func TestBasicState_AddSequencer(t *testing.T) {
	lastBN, err := testState.GetLastBlockNumber(ctx, "")
	assert.NoError(t, err)
	sequencer1 := state.Sequencer{
		Address:     common.HexToAddress("0xab5801a7d398351b8be11c439e05c5b3259aec9b"),
		URL:         "http://www.adrresss1.com",
		ChainID:     big.NewInt(1234),
		BlockNumber: lastBN,
	}
	sequencer2 := state.Sequencer{
		Address:     common.HexToAddress("0xab5801a7d398351b8be11c439e05c5b3259aec9c"),
		URL:         "http://www.adrresss2.com",
		ChainID:     big.NewInt(5678),
		BlockNumber: lastBN,
	}

	sequencer5 := state.Sequencer{
		Address:     common.HexToAddress("0xab5801a7d398351b8be11c439e05c5b3259aec9c"),
		URL:         "http://www.adrresss3.com",
		ChainID:     big.NewInt(5678),
		BlockNumber: lastBN,
	}

	err = testState.AddSequencer(ctx, sequencer1, "")
	assert.NoError(t, err)

	sequencer3, err := testState.GetSequencer(ctx, sequencer1.Address, "")
	assert.NoError(t, err)
	assert.Equal(t, sequencer1.ChainID, sequencer3.ChainID)

	err = testState.AddSequencer(ctx, sequencer2, "")
	assert.NoError(t, err)

	sequencer4, err := testState.GetSequencer(ctx, sequencer2.Address, "")
	assert.NoError(t, err)
	assert.Equal(t, sequencer2, *sequencer4)

	// Update Sequencer
	err = testState.AddSequencer(ctx, sequencer5, "")
	assert.NoError(t, err)

	sequencer6, err := testState.GetSequencer(ctx, sequencer5.Address, "")
	assert.NoError(t, err)
	assert.Equal(t, sequencer5, *sequencer6)
	assert.Equal(t, sequencer5.URL, sequencer6.URL)

	_, err = stateDb.Exec(ctx, "DELETE FROM state.sequencer WHERE chain_id = $1", sequencer1.ChainID.Uint64())
	assert.NoError(t, err)
	_, err = stateDb.Exec(ctx, "DELETE FROM state.sequencer WHERE chain_id = $1", sequencer2.ChainID.Uint64())
	assert.NoError(t, err)
}

/*
func TestStateTransition(t *testing.T) {
	// Load test vectors
	var stateTransitionTestCases []vectors.StateTransitionTestCase
	root := filepath.Clean("../test/vectors/src/state-transition/no-data")
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		tmpStateTransitionTestCases, err := vectors.LoadStateTransitionTestCases(path)
		if err != nil {
			return err
		}

		stateTransitionTestCases = append(stateTransitionTestCases, tmpStateTransitionTestCases...)
		return nil
	})
	require.NoError(t, err)

	for _, testCase := range stateTransitionTestCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()
			// Init database instance
			err = dbutils.InitOrReset(cfg)
			require.NoError(t, err)

			// Create State db
			stateDb, err = db.NewSQLDB(cfg)
			require.NoError(t, err)

			// Create State tree
			store := tree.NewPostgresStore(stateDb)
			mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
			scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
			stateTree := tree.NewStateTree(mt, scCodeStore)

			// Create state
			st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)
			genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
			genesisBlock.ReceivedAt = time.Now()
			genesis := state.Genesis{
				Block:    genesisBlock,
				Balances: make(map[common.Address]*big.Int),
			}

			for _, gacc := range testCase.GenesisAccounts {
				balance := gacc.Balance.Int
				genesis.Balances[common.HexToAddress(gacc.Address)] = &balance
			}

			for gaddr := range genesis.Balances {
				balance, err := stateTree.GetBalance(ctx, gaddr, nil)
				require.NoError(t, err)
				assert.Equal(t, big.NewInt(0), balance)
			}

			err = st.SetGenesis(ctx, genesis)
			require.NoError(t, err)

			root, err := st.GetStateRootByBatchNumber(ctx, 0)
			require.NoError(t, err)

			for gaddr, gbalance := range genesis.Balances {
				balance, err := stateTree.GetBalance(ctx, gaddr, root)
				require.NoError(t, err)
				assert.Equal(t, gbalance, balance)
			}

			var txs []*types.Transaction

			// Check Old roots
			assert.Equal(t, testCase.ExpectedOldRoot, hex.EncodeToHex(root))

			// Check if sequencer is in the DB
			_, err = st.GetSequencer(ctx, common.HexToAddress(testCase.SequencerAddress))
			if err == state.ErrNotFound {
				sq := state.Sequencer{
					Address:     common.HexToAddress(testCase.SequencerAddress),
					URL:         "",
					ChainID:     new(big.Int).SetUint64(testCase.ChainIDSequencer),
					BlockNumber: 0,
				}

				err = st.AddSequencer(ctx, sq)
				require.NoError(t, err)
			}

			// Create Transaction
			for _, vectorTx := range testCase.Txs {
				if string(vectorTx.RawTx) != "" && vectorTx.Overwrite.S == "" {
					var tx types.LegacyTx
					bytes, _ := hex.DecodeString(strings.TrimPrefix(string(vectorTx.RawTx), "0x"))

					err = rlp.DecodeBytes(bytes, &tx)
					if err == nil {
						txs = append(txs, types.NewTx(&tx))
					}
					require.NoError(t, err)
				}
			}

			// Create Batch
			batch := &state.Batch{
				BlockNumber:        uint64(0),
				Sequencer:          common.HexToAddress(testCase.SequencerAddress),
				Aggregator:         addr,
				ConsolidatedTxHash: common.Hash{},
				Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
				Uncles:             nil,
				Transactions:       txs,
				RawTxsData:         nil,
				MaticCollateral:    big.NewInt(1),
				ChainID:            big.NewInt(1000),
				GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
			}

			// Create Batch Processor
			bp, err := st.NewBatchProcessor(ctx, common.HexToAddress(testCase.SequencerAddress), common.Hex2Bytes(strings.TrimPrefix(testCase.ExpectedOldRoot, "0x")))
			require.NoError(t, err)

			err = bp.ProcessBatch(ctx, batch)
			require.NoError(t, err)

			// Check Transaction and Receipts
			transactions, err := testState.GetTxsByBatchNum(ctx, batch.Number().Uint64())
			require.NoError(t, err)

			if len(transactions) > 0 {
				// Check get transaction by batch number and index
				transaction, err := testState.GetTransactionByBatchNumberAndIndex(ctx, batch.Number().Uint64(), 0)
				require.NoError(t, err)
				assert.Equal(t, transaction.Hash(), transactions[0].Hash())

				// Check get transaction by hash and index
				transaction, err = testState.GetTransactionByBatchHashAndIndex(ctx, batch.Hash(), 0)
				require.NoError(t, err)
				assert.Equal(t, transaction.Hash(), transactions[0].Hash())
			}

			root, err = st.GetStateRootByBatchNumber(ctx, batch.Number().Uint64())
			require.NoError(t, err)

			// Check new roots
			assert.Equal(t, testCase.ExpectedNewRoot, hex.EncodeToHex(root))

			for key, vectorLeaf := range testCase.ExpectedNewLeafs {
				newBalance, err := stateTree.GetBalance(ctx, common.HexToAddress(key), root)
				require.NoError(t, err)
				assert.Equal(t, vectorLeaf.Balance.String(), newBalance.String())

				newNonce, err := stateTree.GetNonce(ctx, common.HexToAddress(key), root)
				require.NoError(t, err)
				leafNonce, _ := big.NewInt(0).SetString(vectorLeaf.Nonce, 10)
				assert.Equal(t, leafNonce.String(), newNonce.String())
			}
		})
	}
}
*/
func TestStateTransitionSC(t *testing.T) {
	// Load test vector
	stateTransitionTestCases, err := vectors.LoadStateTransitionTestCases("../test/vectors/src/state-transition-sc.json")
	require.NoError(t, err)

	for _, testCase := range stateTransitionTestCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()
			// Init database instance
			err = dbutils.InitOrReset(cfg)
			require.NoError(t, err)

			// Create State db
			stateDb, err = db.NewSQLDB(cfg)
			require.NoError(t, err)

			// Create State tree
			store := tree.NewPostgresStore(stateDb)
			mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
			scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
			stateTree := tree.NewStateTree(mt, scCodeStore)

			// Create state
			st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)
			genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
			genesisBlock.ReceivedAt = time.Now()
			genesis := state.Genesis{
				Block:          genesisBlock,
				SmartContracts: make(map[common.Address][]byte),
			}

			for _, gsc := range testCase.GenesisSmartContracts {
				genesis.SmartContracts[common.HexToAddress(gsc.Address)] = []byte(gsc.Code)
			}

			err = st.SetGenesis(ctx, genesis, "")
			require.NoError(t, err)
		})
	}
}

func TestLastSeenBatch(t *testing.T) {
	// Create State db
	mtDb, err := db.NewSQLDB(cfg)
	require.NoError(t, err)

	store := tree.NewPostgresStore(mtDb)

	// Create State tree
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)

	// Create state
	scCodeStore := tree.NewPostgresSCCodeStore(mtDb)
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), tree.NewStateTree(mt, scCodeStore))
	ctx := context.Background()

	// Clean Up to reset Genesis
	_, err = stateDb.Exec(ctx, "DELETE FROM state.block")
	if err != nil {
		panic(err)
	}

	err = st.SetLastBatchNumberSeenOnEthereum(ctx, lastBatchNumberSeen, "")
	require.NoError(t, err)
	bn, err := st.GetLastBatchNumberSeenOnEthereum(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, lastBatchNumberSeen, bn)

	err = st.SetLastBatchNumberSeenOnEthereum(ctx, lastBatchNumberSeen+1, "")
	require.NoError(t, err)
	bn, err = st.GetLastBatchNumberSeenOnEthereum(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, lastBatchNumberSeen+1, bn)
}

/*
func TestReceipts(t *testing.T) {
	// Load test vector
	stateTransitionTestCases, err := vectors.LoadStateTransitionTestCases("../test/vectors/src/receipt-test-vectors/receipt-vector.json")
	require.NoError(t, err)

	for _, testCase := range stateTransitionTestCases {
		t.Run(testCase.Description, func(t *testing.T) {
			ctx := context.Background()
			// Init database instance
			err = dbutils.InitOrReset(cfg)
			require.NoError(t, err)

			// Create State db
			stateDb, err = db.NewSQLDB(cfg)
			require.NoError(t, err)

			// Create State tree
			store := tree.NewPostgresStore(stateDb)
			mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
			scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
			stateTree := tree.NewStateTree(mt, scCodeStore)

			// Create state
			st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)

			genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
			genesisBlock.ReceivedAt = time.Now()
			genesis := state.Genesis{
				Block:    genesisBlock,
				Balances: make(map[common.Address]*big.Int),
			}

			for _, gacc := range testCase.GenesisAccounts {
				balance := gacc.Balance.Int
				genesis.Balances[common.HexToAddress(gacc.Address)] = &balance
			}

			for gaddr := range genesis.Balances {
				balance, err := stateTree.GetBalance(ctx, gaddr, nil)
				require.NoError(t, err)
				assert.Equal(t, big.NewInt(0), balance)
			}

			err = st.SetGenesis(ctx, genesis)
			require.NoError(t, err)

			root, err := st.GetStateRootByBatchNumber(ctx, 0)
			require.NoError(t, err)

			for gaddr, gbalance := range genesis.Balances {
				balance, err := stateTree.GetBalance(ctx, gaddr, root)
				require.NoError(t, err)
				assert.Equal(t, gbalance, balance)
			}

			var txs []*types.Transaction

			// Check Old roots
			assert.Equal(t, testCase.ExpectedOldRoot, hex.EncodeToHex(root))

			// Check if sequencer is in the DB
			_, err = st.GetSequencer(ctx, common.HexToAddress(testCase.SequencerAddress))
			if err == state.ErrNotFound {
				sq := state.Sequencer{
					Address:     common.HexToAddress(testCase.SequencerAddress),
					URL:         "",
					ChainID:     new(big.Int).SetUint64(testCase.ChainIDSequencer),
					BlockNumber: 0,
				}

				err = st.AddSequencer(ctx, sq)
				require.NoError(t, err)
			}

			// Create Transaction
			for _, vectorTx := range testCase.Txs {
				if string(vectorTx.RawTx) != "" && vectorTx.Reason == "" {
					var tx types.LegacyTx
					bytes, err := hex.DecodeHex(vectorTx.RawTx)
					require.NoError(t, err)

					err = rlp.DecodeBytes(bytes, &tx)
					require.NoError(t, err)
					txs = append(txs, types.NewTx(&tx))
				}
			}

			// Create Batch
			batch := &state.Batch{
				BlockNumber:        uint64(0),
				Sequencer:          common.HexToAddress(testCase.SequencerAddress),
				Aggregator:         addr,
				ConsolidatedTxHash: common.Hash{},
				Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
				Uncles:             nil,
				Transactions:       txs,
				RawTxsData:         nil,
				MaticCollateral:    big.NewInt(1),
				ChainID:            big.NewInt(1000),
				GlobalExitRoot:     common.HexToHash(testCase.GlobalExitRoot),
			}

			// Create Batch Processor
			bp, err := st.NewBatchProcessor(ctx, common.HexToAddress(testCase.SequencerAddress), root)
			require.NoError(t, err)

			err = bp.ProcessBatch(ctx, batch)
			require.NoError(t, err)

			// Check Transaction and Receipts
			transactions, err := testState.GetTxsByBatchNum(ctx, batch.Number().Uint64())
			require.NoError(t, err)

			if len(transactions) > 0 {
				// Check get transaction by batch number and index
				transaction, err := testState.GetTransactionByBatchNumberAndIndex(ctx, batch.Number().Uint64(), 0)
				require.NoError(t, err)
				assert.Equal(t, transaction.Hash(), transactions[0].Hash())

				// Check get transaction by hash and index
				transaction, err = testState.GetTransactionByBatchHashAndIndex(ctx, batch.Hash(), 0)
				require.NoError(t, err)
				assert.Equal(t, transaction.Hash(), transactions[0].Hash())
			}

			root, err = st.GetStateRootByBatchNumber(ctx, batch.Number().Uint64())
			require.NoError(t, err)

			// Check new roots
			assert.Equal(t, testCase.ExpectedNewRoot, hex.EncodeToHex(root))

			for key, vectorLeaf := range testCase.ExpectedNewLeafs {
				newBalance, err := stateTree.GetBalance(ctx, common.HexToAddress(key), root)
				require.NoError(t, err)
				assert.Equal(t, vectorLeaf.Balance.String(), newBalance.String())

				newNonce, err := stateTree.GetNonce(ctx, common.HexToAddress(key), root)
				require.NoError(t, err)
				leafNonce, _ := big.NewInt(0).SetString(vectorLeaf.Nonce, 10)
				assert.Equal(t, leafNonce.String(), newNonce.String())
			}

			// Get Receipts from vector
			for _, testReceipt := range testCase.Receipts {
				receipt, err := testState.GetTransactionReceipt(ctx, common.HexToHash(testReceipt.Receipt.TransactionHash))
				require.NoError(t, err)
				assert.Equal(t, common.HexToHash(testReceipt.Receipt.TransactionHash), receipt.TxHash)

				// Compare against test receipt
				assert.Equal(t, testReceipt.Receipt.TransactionHash, receipt.TxHash.String())
				assert.Equal(t, testReceipt.Receipt.TransactionIndex, receipt.TransactionIndex)
				assert.Equal(t, batch.Number().Uint64(), receipt.BlockNumber.Uint64())
				assert.Equal(t, testReceipt.Receipt.From, receipt.From.String())
				assert.Equal(t, testReceipt.Receipt.To, receipt.To.String())
				assert.Equal(t, testReceipt.Receipt.CumulativeGastUsed, receipt.CumulativeGasUsed)
				assert.Equal(t, testReceipt.Receipt.GasUsedForTx, receipt.GasUsed)
				assert.Equal(t, testReceipt.Receipt.Status, receipt.Status)

				// BLOCKHASH -> BatchHash
				// This assertion is wrong due to a missalignment between the node team and the protocol team
				// We are commenting this line for now in order to unblock the development and we have created
				// the issue #290 in order to track this fix: https://github.com/hermeznetwork/hermez-core/issues/290
				// assert.Equal(t, common.HexToHash(testReceipt.Receipt.BlockHash), receipt.BlockHash)
			}
		})
	}
}
*/
func TestLastConsolidatedBatch(t *testing.T) {
	// Create State db
	mtDb, err := db.NewSQLDB(cfg)
	require.NoError(t, err)

	store := tree.NewPostgresStore(mtDb)

	// Create State tree
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)

	// Create state
	scCodeStore := tree.NewPostgresSCCodeStore(mtDb)
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), tree.NewStateTree(mt, scCodeStore))
	ctx := context.Background()

	// Clean Up to reset Genesis
	_, err = stateDb.Exec(ctx, "DELETE FROM state.block")
	if err != nil {
		panic(err)
	}

	err = st.SetLastBatchNumberConsolidatedOnEthereum(ctx, lastBatchNumberSeen, "")
	require.NoError(t, err)
	bn, err := st.GetLastBatchNumberConsolidatedOnEthereum(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, lastBatchNumberSeen, bn)

	err = st.SetLastBatchNumberConsolidatedOnEthereum(ctx, lastBatchNumberSeen+1, "")
	require.NoError(t, err)
	bn, err = st.GetLastBatchNumberConsolidatedOnEthereum(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, lastBatchNumberSeen+1, bn)
}

func TestStateErrors(t *testing.T) {
	// Create State db
	mtDb, err := db.NewSQLDB(cfg)
	require.NoError(t, err)

	store := tree.NewPostgresStore(mtDb)

	// Create State tree
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)

	// Create state
	scCodeStore := tree.NewPostgresSCCodeStore(mtDb)
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), tree.NewStateTree(mt, scCodeStore))
	ctx := context.Background()

	// Clean Up to reset Genesis
	_, err = stateDb.Exec(ctx, "DELETE FROM state.block")
	require.NoError(t, err)

	_, err = st.GetStateRoot(ctx, true, "")
	require.Equal(t, state.ErrStateNotSynchronized, err)

	_, err = st.GetBalance(ctx, addr, 0, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetNonce(ctx, addr, 0, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetStateRootByBatchNumber(ctx, 0, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetLastBlock(ctx, "")
	require.Equal(t, state.ErrStateNotSynchronized, err)

	_, err = st.GetPreviousBlock(ctx, 0, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetBlockByHash(ctx, hash1, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetBlockByNumber(ctx, 0, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetLastBlockNumber(ctx, "")
	require.NoError(t, err)

	_, err = st.GetLastBatch(ctx, true, "")
	require.Equal(t, state.ErrStateNotSynchronized, err)

	_, err = st.GetPreviousBatch(ctx, true, 0, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetBatchByHash(ctx, batch1.Hash(), "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetBatchByNumber(ctx, 0, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetLastBatchNumber(ctx, "")
	require.NoError(t, err)

	_, err = st.GetLastConsolidatedBatchNumber(ctx, "")
	require.NoError(t, err)

	_, err = st.GetTransactionByBatchHashAndIndex(ctx, batch1.Hash(), 0, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetTransactionByBatchNumberAndIndex(ctx, batch1.Number().Uint64(), 0, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetTransactionByHash(ctx, txHash, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetTransactionReceipt(ctx, txHash, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetTxsByBatchNum(ctx, batchNumber1, "")
	require.NoError(t, err)

	_, err = st.GetSequencer(ctx, batch1.Sequencer, "")
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetLastBatchNumberSeenOnEthereum(ctx, "")
	require.NoError(t, err)

	_, err = st.GetLastBatchNumberConsolidatedOnEthereum(ctx, "")
	require.NoError(t, err)
}

func TestSCExecution(t *testing.T) {
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var sequencerBalance = 400000
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")

	// Init database instance
	err := dbutils.InitOrReset(cfg)
	require.NoError(t, err)

	// Create State db
	stateDb, err = db.NewSQLDB(cfg)
	require.NoError(t, err)

	// Create State tree
	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	stateTree := tree.NewStateTree(mt, scCodeStore)

	// Create state
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	genesis := state.Genesis{
		Block:    genesisBlock,
		Balances: make(map[common.Address]*big.Int),
	}

	genesis.Balances[sequencerAddress] = new(big.Int).SetInt64(int64(sequencerBalance))
	err = st.SetGenesis(ctx, genesis, "")
	require.NoError(t, err)

	// Register Sequencer
	sequencer := state.Sequencer{
		Address:     sequencerAddress,
		URL:         "http://www.address.com",
		ChainID:     chainIDSequencer,
		BlockNumber: genesisBlock.Header().Number.Uint64(),
	}

	err = testState.AddSequencer(ctx, sequencer, "")
	assert.NoError(t, err)

	var txs []*types.Transaction

	txSCDeploy := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(1),
		Data:     common.Hex2Bytes("608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea2646970667358221220404e37f487a89a932dca5e77faaf6ca2de3b991f93d230604b1b8daaef64766264736f6c63430008070033"),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTxSCDeploy, err := auth.Signer(auth.From, txSCDeploy)
	require.NoError(t, err)

	txs = append(txs, signedTxSCDeploy)

	// Set stored value to 2
	txStoreValue := types.NewTransaction(1, scAddress, new(big.Int), state.TxTransferGas, new(big.Int).SetUint64(1), common.Hex2Bytes("6057361d0000000000000000000000000000000000000000000000000000000000000002"))
	signedTxStoreValue, err := auth.Signer(auth.From, txStoreValue)
	require.NoError(t, err)

	txs = append(txs, signedTxStoreValue)

	// Retrieve stored value
	txRetrieveValue := types.NewTransaction(2, scAddress, new(big.Int), state.TxTransferGas, new(big.Int).SetUint64(1), common.Hex2Bytes("2e64cec1"))
	signedTxRetrieveValue, err := auth.Signer(auth.From, txRetrieveValue)
	require.NoError(t, err)

	txs = append(txs, signedTxRetrieveValue)

	// Create Batch
	batch := &state.Batch{
		BlockNumber:        uint64(0),
		Sequencer:          sequencerAddress,
		Aggregator:         sequencerAddress,
		ConsolidatedTxHash: common.Hash{},
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
		Uncles:             nil,
		Transactions:       txs,
		RawTxsData:         nil,
		MaticCollateral:    big.NewInt(1),
		ReceivedAt:         time.Now(),
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}

	// Create Batch Processor
	stateRoot, err := testState.GetStateRoot(ctx, true, "")
	require.NoError(t, err)
	bp, err := st.NewBatchProcessor(ctx, sequencerAddress, stateRoot, "")
	require.NoError(t, err)

	err = bp.ProcessBatch(ctx, batch)
	require.NoError(t, err)

	receipt, err := testState.GetTransactionReceipt(ctx, signedTxStoreValue.Hash(), "")
	require.NoError(t, err)
	assert.Equal(t, uint64(5420), receipt.GasUsed)

	receipt2, err := testState.GetTransactionReceipt(ctx, signedTxRetrieveValue.Hash(), "")
	require.NoError(t, err)
	assert.Equal(t, uint64(1115), receipt2.GasUsed)

	// Check GetCode
	lastBatch, err := testState.GetLastBatch(ctx, true, "")
	assert.NoError(t, err)
	code, err := st.GetCode(ctx, scAddress, lastBatch.Number().Uint64(), "")
	assert.NoError(t, err)
	assert.NotEqual(t, "", code)
}

/*
func TestSCCall(t *testing.T) {
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var sequencerBalance = 4000000
	scCounterByteCode, err := testutils.ReadBytecode("Counter/Counter.bin")
	require.NoError(t, err)
	var scCounterAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")
	scInteractionByteCode, err := testutils.ReadBytecode("Interaction/Interaction.bin")
	require.NoError(t, err)
	var scInteractionAddress = common.HexToAddress("0x85e844b762A271022b692CF99cE5c59BA0650Ac8")
	var stateRoot = "0x236a5c853ae354e96f6d52b8b40bf46d4348b1ea10364a9de93b68c7b5e40444"
	var expectedFinalRoot = "20664951758840745974444121049737642565652818998385629128471095142198722727287"

	// Init database instance
	err = dbutils.InitOrReset(cfg)
	require.NoError(t, err)

	// Create State db
	stateDb, err = db.NewSQLDB(cfg)
	require.NoError(t, err)

	// Create State tree
	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	stateTree := tree.NewStateTree(mt, scCodeStore)

	// Create state
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	genesis := state.Genesis{
		Block:    genesisBlock,
		Balances: make(map[common.Address]*big.Int),
	}

	genesis.Balances[sequencerAddress] = new(big.Int).SetInt64(int64(sequencerBalance))
	err = st.SetGenesis(ctx, genesis)
	require.NoError(t, err)

	// Register Sequencer
	sequencer := state.Sequencer{
		Address:     sequencerAddress,
		URL:         "http://www.address.com",
		ChainID:     chainIDSequencer,
		BlockNumber: genesisBlock.Header().Number.Uint64(),
	}

	err = st.AddSequencer(ctx, sequencer)
	assert.NoError(t, err)

	var txs []*types.Transaction

	// Deploy counter.sol
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(1),
		Data:     common.Hex2Bytes(scCounterByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)
	txs = append(txs, signedTx)

	// Deploy interaction.sol
	tx1 := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(1),
		Data:     common.Hex2Bytes(scInteractionByteCode),
	})

	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)

	txs = append(txs, signedTx1)

	// Call setCounterAddr method from Interaction SC to set Counter SC Address
	tx2 := types.NewTransaction(2, scInteractionAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("ec39b429000000000000000000000000"+strings.TrimPrefix(scCounterAddress.String(), "0x")))
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)
	txs = append(txs, signedTx2)

	// Increment Counter calling Counter SC
	tx3 := types.NewTransaction(3, scCounterAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("d09de08a"))
	signedTx3, err := auth.Signer(auth.From, tx3)
	require.NoError(t, err)
	txs = append(txs, signedTx3)

	// Retrieve counter value calling Interaction SC (this is the real test as Interaction SC will call Counter SC)
	tx4 := types.NewTransaction(4, scInteractionAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("a87d942c"))
	signedTx4, err := auth.Signer(auth.From, tx4)
	require.NoError(t, err)
	txs = append(txs, signedTx4)

	// Increment Counter calling again
	tx5 := types.NewTransaction(5, scCounterAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("d09de08a"))
	signedTx5, err := auth.Signer(auth.From, tx5)
	require.NoError(t, err)
	txs = append(txs, signedTx5)

	// Retrieve counter value again
	tx6 := types.NewTransaction(6, scInteractionAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("a87d942c"))
	signedTx6, err := auth.Signer(auth.From, tx6)
	require.NoError(t, err)
	txs = append(txs, signedTx6)

	// Create Batch
	batch := &state.Batch{
		BlockNumber:        uint64(0),
		Sequencer:          sequencerAddress,
		Aggregator:         sequencerAddress,
		ConsolidatedTxHash: common.Hash{},
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
		Uncles:             nil,
		Transactions:       txs,
		RawTxsData:         nil,
		MaticCollateral:    big.NewInt(1),
		ReceivedAt:         time.Now(),
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}

	// Create Batch Processor
	bp, err := st.NewBatchProcessor(ctx, sequencerAddress, common.Hex2Bytes(strings.TrimPrefix(stateRoot, "0x")))
	require.NoError(t, err)

	err = bp.ProcessBatch(ctx, batch)
	require.NoError(t, err)

	receipt, err := st.GetTransactionReceipt(ctx, signedTx6.Hash())
	require.NoError(t, err)
	assert.Equal(t, expectedFinalRoot, new(big.Int).SetBytes(receipt.PostState).String())
}
*/
func TestGenesisStorage(t *testing.T) {
	var address = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	// Init database instance
	err := dbutils.InitOrReset(cfg)
	require.NoError(t, err)

	// Create State db
	stateDb, err = db.NewSQLDB(cfg)
	require.NoError(t, err)

	// Create State tree
	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	stateTree := tree.NewStateTree(mt, scCodeStore)

	// Create state
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	genesis := state.Genesis{
		Block:   genesisBlock,
		Storage: make(map[common.Address]map[*big.Int]*big.Int),
	}

	values := make(map[*big.Int]*big.Int)

	for i := 0; i < 10; i++ {
		values[new(big.Int).SetInt64(int64(i))] = new(big.Int).SetInt64(int64(i))
	}

	genesis.Storage[address] = values
	err = st.SetGenesis(ctx, genesis, "")
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		value, err := st.GetStorageAt(ctx, address, new(big.Int).SetInt64(int64(i)), 0, "")
		assert.NoError(t, err)
		assert.NotEqual(t, int64(i), value)
	}
}
func TestSCSelfDestruct(t *testing.T) {
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var sequencerBalance = 120000
	// /tests/contracts/destruct.sol
	scByteCode, err := testutils.ReadBytecode("Destruct/Destruct.bin")
	require.NoError(t, err)
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")

	// Init database instance
	err = dbutils.InitOrReset(cfg)
	require.NoError(t, err)

	// Create State db
	stateDb, err = db.NewSQLDB(cfg)
	require.NoError(t, err)

	// Create State tree
	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	stateTree := tree.NewStateTree(mt, scCodeStore)

	// Create state
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	genesis := state.Genesis{
		Block:    genesisBlock,
		Balances: make(map[common.Address]*big.Int),
	}

	genesis.Balances[sequencerAddress] = new(big.Int).SetInt64(int64(sequencerBalance))
	err = st.SetGenesis(ctx, genesis, "")
	require.NoError(t, err)

	// Register Sequencer
	sequencer := state.Sequencer{
		Address:     sequencerAddress,
		URL:         "http://www.address.com",
		ChainID:     chainIDSequencer,
		BlockNumber: genesisBlock.Header().Number.Uint64(),
	}

	err = st.AddSequencer(ctx, sequencer, "")
	assert.NoError(t, err)

	var txs []*types.Transaction

	// Deploy destruct.sol
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int).SetUint64(0),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(1),
		Data:     common.Hex2Bytes(scByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)
	txs = append(txs, signedTx)

	// Call close method from SC to destroy it
	tx1 := types.NewTransaction(1, scAddress, new(big.Int), 40000, new(big.Int).SetUint64(1), common.Hex2Bytes("43d726d6"))
	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)
	txs = append(txs, signedTx1)

	// Create Batch
	batch := &state.Batch{
		BlockNumber:        uint64(0),
		Sequencer:          sequencerAddress,
		Aggregator:         sequencerAddress,
		ConsolidatedTxHash: common.Hash{},
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
		Uncles:             nil,
		Transactions:       txs,
		RawTxsData:         nil,
		MaticCollateral:    big.NewInt(1),
		ReceivedAt:         time.Now(),
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}

	// Create Batch Processor
	bp, err := st.NewBatchProcessor(ctx, sequencerAddress, common.Hex2Bytes("0x"), "")
	require.NoError(t, err)

	err = bp.ProcessBatch(ctx, batch)
	require.NoError(t, err)

	// Get SC bytecode
	code, err := st.GetCode(ctx, scAddress, batch.Number().Uint64(), "")
	require.NoError(t, err)
	assert.Equal(t, []byte{}, code)
}

func TestEmitLog(t *testing.T) {
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var sequencerBalance = 1200000
	// /tests/contracts/emitLog.sol
	scByteCode, err := testutils.ReadBytecode("EmitLog/EmitLog.bin")
	require.NoError(t, err)
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")

	// Init database instance
	err = dbutils.InitOrReset(cfg)
	require.NoError(t, err)

	// Create State db
	stateDb, err = db.NewSQLDB(cfg)
	require.NoError(t, err)

	// Create State tree
	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	stateTree := tree.NewStateTree(mt, scCodeStore)

	// Create state
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	genesis := state.Genesis{
		Block:    genesisBlock,
		Balances: make(map[common.Address]*big.Int),
	}

	genesis.Balances[sequencerAddress] = new(big.Int).SetInt64(int64(sequencerBalance))
	err = st.SetGenesis(ctx, genesis, "")
	require.NoError(t, err)

	// Register Sequencer
	sequencer := state.Sequencer{
		Address:     sequencerAddress,
		URL:         "http://www.address.com",
		ChainID:     chainIDSequencer,
		BlockNumber: genesisBlock.Header().Number.Uint64(),
	}

	err = st.AddSequencer(ctx, sequencer, "")
	assert.NoError(t, err)

	var txs []*types.Transaction

	// Deploy SC
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int).SetUint64(0),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(1),
		Data:     common.Hex2Bytes(scByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	// tx to call emitLog
	hashCall := sha3.NewLegacyKeccak256()
	_, err = hashCall.Write([]byte("emitLogs()"))
	require.NoError(t, err)
	dataCall := hashCall.Sum(nil)[:4]
	txCall := types.NewTransaction(1, scAddress, new(big.Int), uint64(sequencerBalance), new(big.Int).SetUint64(1), dataCall)
	signedTxCall, err := auth.Signer(auth.From, txCall)
	require.NoError(t, err)

	txs = append(txs, signedTx, signedTxCall)

	// Create Batch
	batch := &state.Batch{
		BlockNumber:        uint64(0),
		Sequencer:          sequencerAddress,
		Aggregator:         sequencerAddress,
		ConsolidatedTxHash: common.Hash{},
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
		Uncles:             nil,
		Transactions:       txs,
		RawTxsData:         nil,
		MaticCollateral:    big.NewInt(1),
		ReceivedAt:         time.Now(),
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}

	// Create Batch Processor
	stateRoot, err := testState.GetStateRoot(ctx, true, "")
	require.NoError(t, err)
	bp, err := st.NewBatchProcessor(ctx, sequencerAddress, stateRoot, "")
	require.NoError(t, err)

	err = bp.ProcessBatch(ctx, batch)
	require.NoError(t, err)

	// Check logs
	receipt, err := st.GetTransactionReceipt(ctx, signedTxCall.Hash(), "")
	require.NoError(t, err)
	require.Equal(t, 10, len(receipt.Logs))
	for _, l := range receipt.Logs {
		assert.Equal(t, scAddress, l.Address)
	}

	hash := batch.Hash()
	logs, err := st.GetLogs(ctx, 0, 0, nil, nil, &hash, "")
	require.NoError(t, err)
	require.Equal(t, 10, len(logs))
	for _, l := range logs {
		assert.Equal(t, scAddress, l.Address)
	}

	logs, err = st.GetLogs(ctx, 0, 5, nil, nil, nil, "")
	require.NoError(t, err)
	require.Equal(t, 10, len(logs))
	for _, l := range logs {
		assert.Equal(t, scAddress, l.Address)
	}

	logs, err = st.GetLogs(ctx, 5, 5, nil, nil, nil, "")
	require.NoError(t, err)
	assert.Equal(t, 0, len(logs))

	addresses := []common.Address{}
	addresses = append(addresses, scAddress)
	logs, err = st.GetLogs(ctx, 0, 5, addresses, nil, nil, "")
	require.NoError(t, err)
	require.Equal(t, 10, len(logs))
	for _, l := range logs {
		assert.Equal(t, scAddress, l.Address)
	}

	type topicsTestCase struct {
		topics           [][]common.Hash
		expectedLogCount int
	}

	topicsTestCases := []topicsTestCase{
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x5e7df75d54e493185612379c616118a4c9ac802de621b010c96f74d22df4b30a")},
			},
			expectedLogCount: 2,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac62")},
			},
			expectedLogCount: 2,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac62")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
			},
			expectedLogCount: 2,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0xbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb09")},
			},
			expectedLogCount: 2,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0xbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb09")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0xbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb09")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0xbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb09")},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7")},
			},
			expectedLogCount: 2,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7")},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")},
			},
			expectedLogCount: 2,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7")},
				{},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7")},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{common.HexToHash("0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7")},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
			},
			expectedLogCount: 5,
		},
		{
			topics: [][]common.Hash{
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000004")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")},
			},
			expectedLogCount: 4,
		},
		{
			topics: [][]common.Hash{
				{},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{},
				{},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{},
				{},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")},
			},
			expectedLogCount: 1,
		},
		{
			topics: [][]common.Hash{
				{},
				{},
				{},
				{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003")},
			},
			expectedLogCount: 2,
		},
	}

	for i, testCase := range topicsTestCases {
		name := strconv.Itoa(i)
		t.Run(name, func(t *testing.T) {
			logs, err = st.GetLogs(ctx, 0, 5, nil, testCase.topics, nil, "")
			require.NoError(t, err)
			require.Equal(t, testCase.expectedLogCount, len(logs))
			for _, l := range logs {
				assert.Equal(t, scAddress, l.Address)
			}
		})
	}
}
func TestEstimateGas(t *testing.T) {
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var sequencerBalance = 400000
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")

	// Init database instance
	err := dbutils.InitOrReset(cfg)
	require.NoError(t, err)

	// Create State db
	stateDb, err = db.NewSQLDB(cfg)
	require.NoError(t, err)

	// Create State tree
	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	stateTree := tree.NewStateTree(mt, scCodeStore)

	// Create state
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	genesis := state.Genesis{
		Block:    genesisBlock,
		Balances: make(map[common.Address]*big.Int),
	}

	genesis.Balances[sequencerAddress] = new(big.Int).SetInt64(int64(sequencerBalance))
	err = st.SetGenesis(ctx, genesis, "")
	require.NoError(t, err)

	// Register Sequencer
	sequencer := state.Sequencer{
		Address:     sequencerAddress,
		URL:         "http://www.address.com",
		ChainID:     chainIDSequencer,
		BlockNumber: genesisBlock.Header().Number.Uint64(),
	}

	err = st.AddSequencer(ctx, sequencer, "")
	assert.NoError(t, err)

	var txs []*types.Transaction

	txSCDeploy := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(1),
		Data:     common.Hex2Bytes("608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea2646970667358221220404e37f487a89a932dca5e77faaf6ca2de3b991f93d230604b1b8daaef64766264736f6c63430008070033"),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTxSCDeploy, err := auth.Signer(auth.From, txSCDeploy)
	require.NoError(t, err)

	txs = append(txs, signedTxSCDeploy)

	// Estimate Gas
	gasEstimation, err := st.EstimateGas(signedTxSCDeploy, "")
	require.NoError(t, err)
	assert.Equal(t, uint64(376040), gasEstimation)

	// Create Batch
	batch := &state.Batch{
		BlockNumber:        uint64(0),
		Sequencer:          sequencerAddress,
		Aggregator:         sequencerAddress,
		ConsolidatedTxHash: common.Hash{},
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
		Uncles:             nil,
		Transactions:       txs,
		RawTxsData:         nil,
		MaticCollateral:    big.NewInt(1),
		ReceivedAt:         time.Now(),
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}

	// Create Batch Processor
	stateRoot, err := testState.GetStateRoot(ctx, true, "")
	require.NoError(t, err)
	bp, err := st.NewBatchProcessor(ctx, sequencerAddress, stateRoot, "")
	require.NoError(t, err)

	err = bp.ProcessBatch(ctx, batch)
	require.NoError(t, err)

	// Set stored value to 2
	txStoreValue := types.NewTransaction(1, scAddress, new(big.Int), state.TxTransferGas, new(big.Int).SetUint64(1), common.Hex2Bytes("6057361d0000000000000000000000000000000000000000000000000000000000000002"))
	signedTxStoreValue, err := auth.Signer(auth.From, txStoreValue)
	require.NoError(t, err)

	// Estimate Gas
	gasEstimation, err = st.EstimateGas(signedTxStoreValue, "")
	require.NoError(t, err)
	assert.Equal(t, uint64(107320), gasEstimation)

	txs = []*types.Transaction{}
	txs = append(txs, signedTxStoreValue)
	batch.Header.Number = big.NewInt(0).SetUint64(2)
	batch.Transactions = txs

	err = bp.ProcessBatch(ctx, batch)
	require.NoError(t, err)

	// Transfer
	txTransfer := types.NewTransaction(1, sequencerAddress, new(big.Int).SetInt64(10000), state.TxTransferGas, new(big.Int).SetUint64(1), nil)
	signedTxTransfer, err := auth.Signer(auth.From, txTransfer)
	require.NoError(t, err)

	// Estimate Gas
	gasEstimation, err = st.EstimateGas(signedTxTransfer, "")
	require.NoError(t, err)
	assert.Equal(t, uint64(state.TxTransferGas), gasEstimation)
}

func TestStorageOnDeploy(t *testing.T) {
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	var sequencerBalance = 80000000
	var scAddress = common.HexToAddress("0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98")
	var expectedStoredValue = common.BigToHash(new(big.Int).SetInt64(1234))
	scByteCode, err := testutils.ReadBytecode("StorageOnDeploy/StorageOnDeploy.bin")
	require.NoError(t, err)

	// Init database instance
	err = dbutils.InitOrReset(cfg)
	require.NoError(t, err)

	// Create State db
	stateDb, err = db.NewSQLDB(cfg)
	require.NoError(t, err)

	// Create State tree
	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	stateTree := tree.NewStateTree(mt, scCodeStore)

	// Create state
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	genesis := state.Genesis{
		Block:    genesisBlock,
		Balances: make(map[common.Address]*big.Int),
	}

	genesis.Balances[sequencerAddress] = new(big.Int).SetInt64(int64(sequencerBalance))
	err = st.SetGenesis(ctx, genesis, "")
	require.NoError(t, err)

	// Register Sequencer
	sequencer := state.Sequencer{
		Address:     sequencerAddress,
		URL:         "http://www.address.com",
		ChainID:     chainIDSequencer,
		BlockNumber: genesisBlock.Header().Number.Uint64(),
	}

	err = st.AddSequencer(ctx, sequencer, "")
	assert.NoError(t, err)

	var txs []*types.Transaction

	txSCDeploy := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      uint64(sequencerBalance),
		GasPrice: new(big.Int).SetUint64(1),
		Data:     common.Hex2Bytes(scByteCode),
	})

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)

	signedTxSCDeploy, err := auth.Signer(auth.From, txSCDeploy)
	require.NoError(t, err)

	txs = append(txs, signedTxSCDeploy)

	// Create Batch
	batch := &state.Batch{
		BlockNumber:        uint64(0),
		Sequencer:          sequencerAddress,
		Aggregator:         sequencerAddress,
		ConsolidatedTxHash: common.Hash{},
		Header:             &types.Header{Number: big.NewInt(0).SetUint64(1)},
		Uncles:             nil,
		Transactions:       txs,
		RawTxsData:         nil,
		MaticCollateral:    big.NewInt(1),
		ReceivedAt:         time.Now(),
		ChainID:            big.NewInt(1000),
		GlobalExitRoot:     common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9fc"),
	}

	lastBatch, err := st.GetLastBatch(ctx, true, "")
	require.NoError(t, err)

	// Create Batch Processor
	bp, err := st.NewBatchProcessor(ctx, sequencerAddress, lastBatch.Header.Root[:], "")
	require.NoError(t, err)

	err = bp.ProcessBatch(ctx, batch)
	require.NoError(t, err)

	value := bp.Host.GetStorage(ctx, scAddress, new(big.Int).SetInt64(0))
	assert.Equal(t, expectedStoredValue, value)
}

func TestConcurrentDBTransactions(t *testing.T) {
	// Init database instance
	err := dbutils.InitOrReset(cfg)
	require.NoError(t, err)

	// Create State db
	stateDb, err = db.NewSQLDB(cfg)
	require.NoError(t, err)

	// Create State tree
	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	stateTree := tree.NewStateTree(mt, scCodeStore)

	// Create state
	st := state.NewState(stateCfg, state.NewPostgresStorage(stateDb), stateTree)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			txBundleID, err := st.BeginStateTransaction(ctx)
			require.NoError(t, err)

			require.NoError(t, st.Commit(ctx, txBundleID))
		}(i)
	}
	wg.Wait()
}
