package state_test

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
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/pgstatestorage"
	"github.com/hermeznetwork/hermez-core/state/runtime/evm"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/hermeznetwork/hermez-core/test/vectors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	stateDb                                                *pgxpool.Pool
	testState                                              state.State
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
	DefaultChainID: 1000,
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
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	testState = state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateDb), tree.NewStateTree(mt, scCodeStore, nil))

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
		BatchNumber:        batchNumber1,
		BlockNumber:        blockNumber1,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: consolidatedTxHash,
		Header:             nil,
		Uncles:             nil,
		RawTxsData:         nil,
		MaticCollateral:    maticCollateral,
		ReceivedAt:         time.Now(),
	}
	batch2 = &state.Batch{
		BatchNumber:        batchNumber2,
		BlockNumber:        blockNumber1,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: consolidatedTxHash,
		Header:             nil,
		Uncles:             nil,
		RawTxsData:         nil,
		MaticCollateral:    maticCollateral,
		ReceivedAt:         time.Now(),
	}
	batch3 = &state.Batch{
		BatchNumber:        batchNumber3,
		BlockNumber:        blockNumber2,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: common.Hash{},
		Header:             nil,
		Uncles:             nil,
		Transactions:       nil,
		RawTxsData:         nil,
		MaticCollateral:    maticCollateral,
		ReceivedAt:         time.Now(),
	}
	batch4 = &state.Batch{
		BatchNumber:        batchNumber4,
		BlockNumber:        blockNumber2,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: common.Hash{},
		Header:             nil,
		Uncles:             nil,
		Transactions:       nil,
		RawTxsData:         nil,
		MaticCollateral:    maticCollateral,
		ReceivedAt:         time.Now(),
	}

	_, err = stateDb.Exec(ctx, "DELETE FROM state.batch")
	if err != nil {
		panic(err)
	}

	batches := []*state.Batch{batch1, batch2, batch3, batch4}

	bp, err := testState.NewGenesisBatchProcessor(nil)
	if err != nil {
		panic(err)
	}

	for _, b := range batches {
		err := bp.ProcessBatch(b)
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
	lastBlock, err := testState.GetLastBlock(ctx)
	assert.NoError(t, err)
	assert.Equal(t, block2.BlockNumber, lastBlock.BlockNumber)
}

func TestBasicState_GetPreviousBlock(t *testing.T) {
	previousBlock, err := testState.GetPreviousBlock(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, block1.BlockNumber, previousBlock.BlockNumber)
}

func TestBasicState_GetBlockByHash(t *testing.T) {
	block, err := testState.GetBlockByHash(ctx, hash1)
	assert.NoError(t, err)
	assert.Equal(t, block1.BlockHash, block.BlockHash)
	assert.Equal(t, block1.BlockNumber, block.BlockNumber)
}

func TestBasicState_GetBlockByNumber(t *testing.T) {
	block, err := testState.GetBlockByNumber(ctx, blockNumber2)
	assert.NoError(t, err)
	assert.Equal(t, block2.BlockNumber, block.BlockNumber)
	assert.Equal(t, block2.BlockHash, block.BlockHash)
}

func TestBasicState_GetLastVirtualBatch(t *testing.T) {
	lastBatch, err := testState.GetLastBatch(ctx, true)
	assert.NoError(t, err)
	assert.Equal(t, batch4.Hash(), lastBatch.Hash())
	assert.Equal(t, batch4.BatchNumber, lastBatch.BatchNumber)
}

func TestBasicState_GetLastBatch(t *testing.T) {
	lastBatch, err := testState.GetLastBatch(ctx, false)
	assert.NoError(t, err)
	assert.Equal(t, batch2.Hash(), lastBatch.Hash())
	assert.Equal(t, batch2.BatchNumber, lastBatch.BatchNumber)
	assert.Equal(t, maticCollateral, lastBatch.MaticCollateral)
}

func TestBasicState_GetPreviousBatch(t *testing.T) {
	previousBatch, err := testState.GetPreviousBatch(ctx, false, 1)
	assert.NoError(t, err)
	assert.Equal(t, batch1.Hash(), previousBatch.Hash())
	assert.Equal(t, batch1.BatchNumber, previousBatch.BatchNumber)
	assert.Equal(t, maticCollateral, previousBatch.MaticCollateral)
}

func TestBasicState_GetBatchByHash(t *testing.T) {
	batch, err := testState.GetBatchByHash(ctx, batch1.Hash())
	assert.NoError(t, err)
	assert.Equal(t, batch1.Hash(), batch.Hash())
	assert.Equal(t, batch1.BatchNumber, batch.BatchNumber)
	assert.Equal(t, maticCollateral, batch1.MaticCollateral)
}

func TestBasicState_GetBatchByNumber(t *testing.T) {
	batch, err := testState.GetBatchByNumber(ctx, batch1.BatchNumber)
	assert.NoError(t, err)
	assert.Equal(t, batch1.BatchNumber, batch.BatchNumber)
	assert.Equal(t, batch1.Hash(), batch.Hash())
}

func TestBasicState_GetLastBatchNumber(t *testing.T) {
	batchNumber, err := testState.GetLastBatchNumber(ctx)
	assert.NoError(t, err)
	assert.Equal(t, batch4.BatchNumber, batchNumber)
}

func TestBasicState_ConsolidateBatch(t *testing.T) {
	batchNumber := uint64(5)
	batch := &state.Batch{
		BatchNumber:        batchNumber,
		BlockNumber:        blockNumber2,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: common.Hash{},
		Header:             nil,
		Uncles:             nil,
		Transactions:       nil,
		RawTxsData:         nil,
		MaticCollateral:    maticCollateral,
		ReceivedAt:         time.Now(),
	}

	bp, err := testState.NewGenesisBatchProcessor(nil)
	assert.NoError(t, err)

	err = bp.ProcessBatch(batch)
	assert.NoError(t, err)

	insertedBatch, err := testState.GetBatchByNumber(ctx, batchNumber)
	assert.NoError(t, err)
	assert.Equal(t, common.Hash{}, insertedBatch.ConsolidatedTxHash)

	err = testState.ConsolidateBatch(ctx, batchNumber, consolidatedTxHash, time.Now())
	assert.NoError(t, err)

	insertedBatch, err = testState.GetBatchByNumber(ctx, batchNumber)
	assert.NoError(t, err)
	assert.Equal(t, consolidatedTxHash, insertedBatch.ConsolidatedTxHash)

	_, err = stateDb.Exec(ctx, "DELETE FROM state.batch WHERE batch_num = $1", batchNumber)
	assert.NoError(t, err)
}

func TestBasicState_GetTransactionCount(t *testing.T) {
	count, err := testState.GetTransactionCount(ctx, addr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), count)
}

func TestBasicState_GetTxsByBatchNum(t *testing.T) {
	txs, err := testState.GetTxsByBatchNum(ctx, batchNumber1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(txs))
}

func TestBasicState_GetTransactionByHash(t *testing.T) {
	tx, err := testState.GetTransactionByHash(ctx, txHash)
	assert.NoError(t, err)
	assert.Equal(t, txHash, tx.Hash())
}

func TestBasicState_AddBlock(t *testing.T) {
	lastBN, err := testState.GetLastBlockNumber(ctx)
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
	err = testState.AddBlock(ctx, block1)
	assert.NoError(t, err)
	err = testState.AddBlock(ctx, block2)
	assert.NoError(t, err)

	block3, err := testState.GetBlockByNumber(ctx, block1.BlockNumber)
	assert.NoError(t, err)
	assert.Equal(t, block1.BlockHash, block3.BlockHash)
	assert.Equal(t, block1.ParentHash, block3.ParentHash)

	block4, err := testState.GetBlockByNumber(ctx, block2.BlockNumber)
	assert.NoError(t, err)
	assert.Equal(t, block2.BlockHash, block4.BlockHash)
	assert.Equal(t, block2.ParentHash, block4.ParentHash)

	_, err = stateDb.Exec(ctx, "DELETE FROM state.block WHERE block_num = $1", block1.BlockNumber)
	assert.NoError(t, err)
	_, err = stateDb.Exec(ctx, "DELETE FROM state.block WHERE block_num = $1", block2.BlockNumber)
	assert.NoError(t, err)
}

func TestBasicState_AddSequencer(t *testing.T) {
	lastBN, err := testState.GetLastBlockNumber(ctx)
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

	err = testState.AddSequencer(ctx, sequencer1)
	assert.NoError(t, err)

	sequencer3, err := testState.GetSequencer(ctx, sequencer1.Address)
	assert.NoError(t, err)
	assert.Equal(t, sequencer1.ChainID, sequencer3.ChainID)

	err = testState.AddSequencer(ctx, sequencer2)
	assert.NoError(t, err)

	sequencer4, err := testState.GetSequencer(ctx, sequencer2.Address)
	assert.NoError(t, err)
	assert.Equal(t, sequencer2, *sequencer4)

	// Update Sequencer
	err = testState.AddSequencer(ctx, sequencer5)
	assert.NoError(t, err)

	sequencer6, err := testState.GetSequencer(ctx, sequencer5.Address)
	assert.NoError(t, err)
	assert.Equal(t, sequencer5, *sequencer6)
	assert.Equal(t, sequencer5.URL, sequencer6.URL)

	_, err = stateDb.Exec(ctx, "DELETE FROM state.sequencer WHERE chain_id = $1", sequencer1.ChainID.Uint64())
	assert.NoError(t, err)
	_, err = stateDb.Exec(ctx, "DELETE FROM state.sequencer WHERE chain_id = $1", sequencer2.ChainID.Uint64())
	assert.NoError(t, err)
}

func TestStateTransition(t *testing.T) {
	// Load test vector
	stateTransitionTestCases, err := vectors.LoadStateTransitionTestCases("../test/vectors/state-transition.json")
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
			mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)
			scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
			stateTree := tree.NewStateTree(mt, scCodeStore, nil)

			// Create state
			st := state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateDb), stateTree)

			genesis := state.Genesis{
				Balances: make(map[common.Address]*big.Int),
			}

			for _, gacc := range testCase.GenesisAccounts {
				balance := gacc.Balance.Int
				genesis.Balances[common.HexToAddress(gacc.Address)] = &balance
			}

			for gaddr := range genesis.Balances {
				balance, err := stateTree.GetBalance(gaddr, nil)
				require.NoError(t, err)
				assert.Equal(t, big.NewInt(0), balance)
			}

			err = st.SetGenesis(ctx, genesis)
			require.NoError(t, err)

			root, err := st.GetStateRootByBatchNumber(0)
			require.NoError(t, err)

			for gaddr, gbalance := range genesis.Balances {
				balance, err := stateTree.GetBalance(gaddr, root)
				require.NoError(t, err)
				assert.Equal(t, gbalance, balance)
			}

			var txs []*types.Transaction

			// Check Old roots
			assert.Equal(t, testCase.ExpectedOldRoot, new(big.Int).SetBytes(root).String())

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
				BatchNumber:        1,
				BlockNumber:        uint64(0),
				Sequencer:          common.HexToAddress(testCase.SequencerAddress),
				Aggregator:         addr,
				ConsolidatedTxHash: common.Hash{},
				Header:             nil,
				Uncles:             nil,
				Transactions:       txs,
				RawTxsData:         nil,
				MaticCollateral:    big.NewInt(1),
			}

			// Create Batch Processor
			bp, err := st.NewBatchProcessor(common.HexToAddress(testCase.SequencerAddress), 0)
			require.NoError(t, err)

			err = bp.ProcessBatch(batch)
			require.NoError(t, err)

			// Check Transaction and Receipts
			transactions, err := testState.GetTxsByBatchNum(ctx, batch.BatchNumber)
			require.NoError(t, err)

			if len(transactions) > 0 {
				// Check get transaction by batch number and index
				transaction, err := testState.GetTransactionByBatchNumberAndIndex(ctx, batch.BatchNumber, 0)
				require.NoError(t, err)
				assert.Equal(t, transaction.Hash(), transactions[0].Hash())

				// Check get transaction by hash and index
				transaction, err = testState.GetTransactionByBatchHashAndIndex(ctx, batch.Hash(), 0)
				require.NoError(t, err)
				assert.Equal(t, transaction.Hash(), transactions[0].Hash())
			}

			root, err = st.GetStateRootByBatchNumber(batch.BatchNumber)
			require.NoError(t, err)

			// Check new roots
			assert.Equal(t, testCase.ExpectedNewRoot, new(big.Int).SetBytes(root).String())

			for key, vectorLeaf := range testCase.ExpectedNewLeafs {
				newBalance, err := stateTree.GetBalance(common.HexToAddress(key), root)
				require.NoError(t, err)
				assert.Equal(t, vectorLeaf.Balance.String(), newBalance.String())

				newNonce, err := stateTree.GetNonce(common.HexToAddress(key), root)
				require.NoError(t, err)
				leafNonce, _ := big.NewInt(0).SetString(vectorLeaf.Nonce, 10)
				assert.Equal(t, leafNonce.String(), newNonce.String())
			}
		})
	}
}

func TestStateTransitionSC(t *testing.T) {
	// Load test vector
	stateTransitionTestCases, err := vectors.LoadStateTransitionTestCases("../test/vectors/state-transition-sc.json")
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
			mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)
			scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
			stateTree := tree.NewStateTree(mt, scCodeStore, nil)

			// Create state
			st := state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateDb), stateTree)

			genesis := state.Genesis{
				SmartContracts: make(map[common.Address][]byte),
			}

			for _, gsc := range testCase.GenesisSmartContracts {
				genesis.SmartContracts[common.HexToAddress(gsc.Address)] = []byte(gsc.Code)
			}

			err = st.SetGenesis(ctx, genesis)
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
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)

	// Create state
	scCodeStore := tree.NewPostgresSCCodeStore(mtDb)
	st := state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateDb), tree.NewStateTree(mt, scCodeStore, nil))
	ctx := context.Background()

	// Clean Up to reset Genesis
	_, err = stateDb.Exec(ctx, "DELETE FROM state.block")
	if err != nil {
		panic(err)
	}

	err = st.SetLastBatchNumberSeenOnEthereum(ctx, lastBatchNumberSeen)
	require.NoError(t, err)
	bn, err := st.GetLastBatchNumberSeenOnEthereum(ctx)
	require.NoError(t, err)
	assert.Equal(t, lastBatchNumberSeen, bn)

	err = st.SetLastBatchNumberSeenOnEthereum(ctx, lastBatchNumberSeen+1)
	require.NoError(t, err)
	bn, err = st.GetLastBatchNumberSeenOnEthereum(ctx)
	require.NoError(t, err)
	assert.Equal(t, lastBatchNumberSeen+1, bn)
}

func TestReceipts(t *testing.T) {
	// Load test vector
	stateTransitionTestCases, err := vectors.LoadStateTransitionTestCases("../test/vectors/receipt-vector.json")
	if err != nil {
		t.Error(err)
		return
	}

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
			mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)
			scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
			stateTree := tree.NewStateTree(mt, scCodeStore, nil)

			// Create state
			st := state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateDb), stateTree)

			genesis := state.Genesis{
				Balances: make(map[common.Address]*big.Int),
			}

			for _, gacc := range testCase.GenesisAccounts {
				balance := gacc.Balance.Int
				genesis.Balances[common.HexToAddress(gacc.Address)] = &balance
			}

			for gaddr := range genesis.Balances {
				balance, err := stateTree.GetBalance(gaddr, nil)
				require.NoError(t, err)
				assert.Equal(t, big.NewInt(0), balance)
			}

			err = st.SetGenesis(ctx, genesis)
			require.NoError(t, err)

			root, err := st.GetStateRootByBatchNumber(0)
			require.NoError(t, err)

			for gaddr, gbalance := range genesis.Balances {
				balance, err := stateTree.GetBalance(gaddr, root)
				require.NoError(t, err)
				assert.Equal(t, gbalance, balance)
			}

			var txs []*types.Transaction

			// Check Old roots
			assert.Equal(t, testCase.ExpectedOldRoot, new(big.Int).SetBytes(root).String())

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
				if string(vectorTx.RawTx) != "" && vectorTx.Overwrite.S == "" && vectorTx.Reason == "" {
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
				BatchNumber:        1,
				BlockNumber:        uint64(0),
				Sequencer:          common.HexToAddress(testCase.SequencerAddress),
				Aggregator:         addr,
				ConsolidatedTxHash: common.Hash{},
				Header:             nil,
				Uncles:             nil,
				Transactions:       txs,
				RawTxsData:         nil,
				MaticCollateral:    big.NewInt(1),
			}

			// Create Batch Processor
			bp, err := st.NewBatchProcessor(common.HexToAddress(testCase.SequencerAddress), 0)
			require.NoError(t, err)

			err = bp.ProcessBatch(batch)
			require.NoError(t, err)

			// Check Transaction and Receipts
			transactions, err := testState.GetTxsByBatchNum(ctx, batch.BatchNumber)
			require.NoError(t, err)

			if len(transactions) > 0 {
				// Check get transaction by batch number and index
				transaction, err := testState.GetTransactionByBatchNumberAndIndex(ctx, batch.BatchNumber, 0)
				require.NoError(t, err)
				assert.Equal(t, transaction.Hash(), transactions[0].Hash())

				// Check get transaction by hash and index
				transaction, err = testState.GetTransactionByBatchHashAndIndex(ctx, batch.Hash(), 0)
				require.NoError(t, err)
				assert.Equal(t, transaction.Hash(), transactions[0].Hash())
			}

			root, err = st.GetStateRootByBatchNumber(batch.BatchNumber)
			require.NoError(t, err)

			// Check new roots
			assert.Equal(t, testCase.ExpectedNewRoot, new(big.Int).SetBytes(root).String())

			for key, vectorLeaf := range testCase.ExpectedNewLeafs {
				newBalance, err := stateTree.GetBalance(common.HexToAddress(key), root)
				require.NoError(t, err)
				assert.Equal(t, vectorLeaf.Balance.String(), newBalance.String())

				newNonce, err := stateTree.GetNonce(common.HexToAddress(key), root)
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
				assert.Equal(t, testReceipt.Receipt.BlockNumber, receipt.BlockNumber.Uint64())
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

func TestLastConsolidatedBatch(t *testing.T) {
	// Create State db
	mtDb, err := db.NewSQLDB(cfg)
	require.NoError(t, err)

	store := tree.NewPostgresStore(mtDb)

	// Create State tree
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)

	// Create state
	scCodeStore := tree.NewPostgresSCCodeStore(mtDb)
	st := state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateDb), tree.NewStateTree(mt, scCodeStore, nil))
	ctx := context.Background()

	// Clean Up to reset Genesis
	_, err = stateDb.Exec(ctx, "DELETE FROM state.block")
	if err != nil {
		panic(err)
	}

	err = st.SetLastBatchNumberConsolidatedOnEthereum(ctx, lastBatchNumberSeen)
	require.NoError(t, err)
	bn, err := st.GetLastBatchNumberConsolidatedOnEthereum(ctx)
	require.NoError(t, err)
	assert.Equal(t, lastBatchNumberSeen, bn)

	err = st.SetLastBatchNumberConsolidatedOnEthereum(ctx, lastBatchNumberSeen+1)
	require.NoError(t, err)
	bn, err = st.GetLastBatchNumberConsolidatedOnEthereum(ctx)
	require.NoError(t, err)
	assert.Equal(t, lastBatchNumberSeen+1, bn)
}

func TestStateErrors(t *testing.T) {
	// Create State db
	mtDb, err := db.NewSQLDB(cfg)
	require.NoError(t, err)

	store := tree.NewPostgresStore(mtDb)

	// Create State tree
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)

	// Create state
	scCodeStore := tree.NewPostgresSCCodeStore(mtDb)
	st := state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateDb), tree.NewStateTree(mt, scCodeStore, nil))
	ctx := context.Background()

	// Clean Up to reset Genesis
	_, err = stateDb.Exec(ctx, "DELETE FROM state.block")
	require.NoError(t, err)

	_, err = st.GetStateRoot(ctx, true)
	require.Equal(t, state.ErrStateNotSynchronized, err)

	_, err = st.GetBalance(addr, 0)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetNonce(addr, 0)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetStateRootByBatchNumber(0)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetLastBlock(ctx)
	require.Equal(t, state.ErrStateNotSynchronized, err)

	_, err = st.GetPreviousBlock(ctx, 0)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetBlockByHash(ctx, hash1)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetBlockByNumber(ctx, 0)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetLastBlockNumber(ctx)
	require.NoError(t, err)

	_, err = st.GetLastBatch(ctx, true)
	require.Equal(t, state.ErrStateNotSynchronized, err)

	_, err = st.GetPreviousBatch(ctx, true, 0)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetBatchByHash(ctx, batch1.Hash())
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetBatchByNumber(ctx, 0)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetLastBatchNumber(ctx)
	require.NoError(t, err)

	_, err = st.GetLastConsolidatedBatchNumber(ctx)
	require.NoError(t, err)

	_, err = st.GetTransactionByBatchHashAndIndex(ctx, batch1.Hash(), 0)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetTransactionByBatchNumberAndIndex(ctx, batch1.BatchNumber, 0)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetTransactionByHash(ctx, txHash)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetTransactionReceipt(ctx, txHash)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetTxsByBatchNum(ctx, batchNumber1)
	require.NoError(t, err)

	_, err = st.GetSequencer(ctx, batch1.Sequencer)
	require.Equal(t, state.ErrNotFound, err)

	_, err = st.GetLastBatchNumberSeenOnEthereum(ctx)
	require.NoError(t, err)

	_, err = st.GetLastBatchNumberConsolidatedOnEthereum(ctx)
	require.NoError(t, err)
}

func TestSCExecution(t *testing.T) {
	var chainIDSequencer = new(big.Int).SetInt64(400)
	var sequencerAddress = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	var sequencerPvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"

	// Init database instance
	err := dbutils.InitOrReset(cfg)
	require.NoError(t, err)

	// Create State db
	stateDb, err = db.NewSQLDB(cfg)
	require.NoError(t, err)

	// Create State tree
	store := tree.NewPostgresStore(stateDb)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity, nil)
	scCodeStore := tree.NewPostgresSCCodeStore(stateDb)
	stateTree := tree.NewStateTree(mt, scCodeStore, nil)

	// Create state
	st := state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateDb), stateTree)

	genesis := state.Genesis{
		SmartContracts: make(map[common.Address][]byte),
	}

	genesis.SmartContracts[addr] = []byte{
		evm.PUSH1, 0x01, evm.PUSH1, 0x02, evm.ADD,
		evm.PUSH1, 0x00, evm.MSTORE8,
		evm.PUSH1, 0x01, evm.PUSH1, 0x00, evm.RETURN,
	}

	err = st.SetGenesis(ctx, genesis)
	require.NoError(t, err)

	var txs []*types.Transaction

	tx := types.NewTransaction(0, addr, big.NewInt(1000000000000000000), 10000, new(big.Int).SetUint64(1), nil)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(sequencerPvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDSequencer)
	require.NoError(t, err)
	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	txs = append(txs, signedTx)

	// Create Batch
	batch := &state.Batch{
		BatchNumber:        1,
		BlockNumber:        uint64(0),
		Sequencer:          sequencerAddress,
		Aggregator:         sequencerAddress,
		ConsolidatedTxHash: common.Hash{},
		Header:             nil,
		Uncles:             nil,
		Transactions:       txs,
		RawTxsData:         nil,
		MaticCollateral:    big.NewInt(1),
	}

	// Create Batch Processor
	bp, err := st.NewBatchProcessor(addr, 0)
	require.NoError(t, err)

	err = bp.ProcessBatch(batch)
	require.NoError(t, err)

	receipt, err := testState.GetTransactionReceipt(ctx, signedTx.Hash())
	require.NoError(t, err)
	// Compare against test receipt
	assert.Equal(t, uint64(24), receipt.GasUsed)
}
