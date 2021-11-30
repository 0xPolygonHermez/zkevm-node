package state

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

var (
	stateDb                                                *pgxpool.Pool
	state                                                  State
	block1, block2                                         *Block
	addr                                                   common.Address = common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	hash1, hash2                                           common.Hash
	hash3                                                  common.Hash = common.HexToHash("0x56ab2c03b9ffc32ed927c3665d6c21c431527e676c345d18f2841747a3a9af34")
	hash4                                                  common.Hash = common.HexToHash("0x8b86252fd1b94139154aee46b61f7610100d4075da3886d95ef3694aa016b4ab")
	blockNumber1, blockNumber2                             uint64      = 1, 2
	batchNumber1, batchNumber2, batchNumber3, batchNumber4 uint64      = 1, 2, 3, 4
	batch1, batch2, batch3, batch4                         *Batch
	consolidatedTxHash                                     common.Hash = common.HexToHash("0x125714bb4db48757007fff2671b37637bbfd6d47b3a4757ebbd0c5222984f905")
	txHash                                                 common.Hash
	ctx                                                    = context.Background()
)

// TODO: understand, from where should we get config for tests. This is temporary
var cfg = db.Config{
	Database: "polygon-hermez",
	User:     "hermez",
	Password: "polygon",
	Host:     "localhost",
	Port:     "5432",
}

func TestMain(m *testing.M) {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})

	dbutils.StartPostgreSQL(cfg.Database, cfg.User, cfg.Password, "") //nolint:gosec,errcheck
	defer dbutils.StopPostgreSQL()                                    //nolint:gosec,errcheck

	// init db
	var err error
	err = db.RunMigrations(cfg)
	if err != nil {
		panic(err)
	}

	stateDb, err = db.NewSQLDB(cfg)
	if err != nil {
		panic(err)
	}
	hash1 = common.HexToHash("0x65b4699dda5f7eb4519c730e6a48e73c90d2b1c8efcd6a6abdfd28c3b8e7d7d9")
	hash2 = common.HexToHash("0x613aabebf4fddf2ad0f034a8c73aa2f9c5a6fac3a07543023e0a6ee6f36e5795")

	setUpBlocks()
	setUpBatches()
	setUpTransactions()
	state = NewState(stateDb, nil)

	result := m.Run()

	cleanUp()

	stateDb.Close()
	os.Exit(result)
}

func cleanUp() {
	var err error
	_, err = stateDb.Exec(ctx, "DELETE FROM batch")
	if err != nil {
		panic(err)
	}
	_, err = stateDb.Exec(ctx, "DELETE FROM block")
	if err != nil {
		panic(err)
	}
	_, err = stateDb.Exec(ctx, "DELETE FROM transaction")
	if err != nil {
		panic(err)
	}
}

func setUpBlocks() {
	var err error
	block1 = &Block{
		BlockNumber: blockNumber1,
		BlockHash:   hash1,
		ParentHash:  hash1,
		ReceivedAt:  time.Now(),
	}
	block2 = &Block{
		BlockNumber: blockNumber2,
		BlockHash:   hash2,
		ParentHash:  hash1,
		ReceivedAt:  time.Now(),
	}

	_, err = stateDb.Exec(ctx, "DELETE FROM block")
	if err != nil {
		panic(err)
	}

	_, err = stateDb.Exec(ctx, "INSERT INTO block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)",
		block1.BlockNumber, block1.BlockHash.Bytes(), block1.ParentHash.Bytes(), block1.ReceivedAt)
	if err != nil {
		panic(err)
	}

	_, err = stateDb.Exec(ctx, "INSERT INTO block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)",
		block2.BlockNumber, block2.BlockHash.Bytes(), block2.ParentHash.Bytes(), block2.ReceivedAt)
	if err != nil {
		panic(err)
	}
}

func setUpBatches() {
	var err error

	batch1 = &Batch{
		BatchNumber:        batchNumber1,
		BatchHash:          hash1,
		BlockNumber:        blockNumber1,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: consolidatedTxHash,
		Header:             nil,
		Uncles:             nil,
		RawTxsData:         nil,
	}
	batch2 = &Batch{
		BatchNumber:        batchNumber2,
		BatchHash:          hash2,
		BlockNumber:        blockNumber1,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: consolidatedTxHash,
		Header:             nil,
		Uncles:             nil,
		RawTxsData:         nil,
	}
	batch3 = &Batch{
		BatchNumber:        batchNumber3,
		BatchHash:          hash3,
		BlockNumber:        blockNumber2,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: common.Hash{},
		Header:             nil,
		Uncles:             nil,
		Transactions:       nil,
		RawTxsData:         nil,
	}
	batch4 = &Batch{
		BatchNumber:        batchNumber4,
		BatchHash:          hash4,
		BlockNumber:        blockNumber2,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: common.Hash{},
		Header:             nil,
		Uncles:             nil,
		Transactions:       nil,
		RawTxsData:         nil,
	}

	_, err = stateDb.Exec(ctx, "DELETE FROM batch")
	if err != nil {
		panic(err)
	}

	batches := []*Batch{batch1, batch2, batch3, batch4}

	for _, b := range batches {
		_, err = stateDb.Exec(ctx, "INSERT INTO batch (batch_num, batch_hash, block_num, sequencer, aggregator, consolidated_tx_hash) VALUES ($1, $2, $3, $4, $5, $6)",
			b.BatchNumber, b.BatchHash, b.BlockNumber, b.Sequencer, b.Aggregator, b.ConsolidatedTxHash)
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
	sql := "INSERT INTO transaction (hash, from_address, encoded, decoded, batch_num) VALUES($1, $2, $3, $4, $5)"
	if _, err := stateDb.Exec(ctx, sql, txHash, addr, encoded, decoded, batchNumber1); err != nil {
		panic(err)
	}
}

func TestBasicState_GetLastBlock(t *testing.T) {
	lastBlock, err := state.GetLastBlock(ctx)
	assert.NoError(t, err)
	assert.Equal(t, block2.BlockNumber, lastBlock.BlockNumber)
}

func TestBasicState_GetPreviousBlock(t *testing.T) {
	previousBlock, err := state.GetPreviousBlock(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, block1.BlockNumber, previousBlock.BlockNumber)
}

func TestBasicState_GetBlockByHash(t *testing.T) {
	block, err := state.GetBlockByHash(ctx, hash1)
	assert.NoError(t, err)
	assert.Equal(t, block1.BlockHash, block.BlockHash)
	assert.Equal(t, block1.BlockNumber, block.BlockNumber)
}

func TestBasicState_GetBlockByNumber(t *testing.T) {
	block, err := state.GetBlockByNumber(ctx, blockNumber2)
	assert.NoError(t, err)
	assert.Equal(t, block2.BlockNumber, block.BlockNumber)
	assert.Equal(t, block2.BlockHash, block.BlockHash)
}

func TestBasicState_GetLastVirtualBatch(t *testing.T) {
	lastBatch, err := state.GetLastBatch(ctx, true)
	assert.NoError(t, err)
	assert.Equal(t, batch4.BatchHash, lastBatch.BatchHash)
	assert.Equal(t, batch4.BatchNumber, lastBatch.BatchNumber)
}

func TestBasicState_GetLastBatch(t *testing.T) {
	lastBatch, err := state.GetLastBatch(ctx, false)
	assert.NoError(t, err)
	assert.Equal(t, batch2.BatchHash, lastBatch.BatchHash)
	assert.Equal(t, batch2.BatchNumber, lastBatch.BatchNumber)
}

func TestBasicState_GetPreviousBatch(t *testing.T) {
	previousBatch, err := state.GetPreviousBatch(ctx, false, 1)
	assert.NoError(t, err)
	assert.Equal(t, batch1.BatchHash, previousBatch.BatchHash)
	assert.Equal(t, batch1.BatchNumber, previousBatch.BatchNumber)
}

func TestBasicState_GetBatchByHash(t *testing.T) {
	batch, err := state.GetBatchByHash(ctx, batch1.BatchHash)
	assert.NoError(t, err)
	assert.Equal(t, batch1.BatchHash, batch.BatchHash)
	assert.Equal(t, batch1.BatchNumber, batch.BatchNumber)
}

func TestBasicState_GetBatchByNumber(t *testing.T) {
	batch, err := state.GetBatchByNumber(ctx, batch1.BatchNumber)
	assert.NoError(t, err)
	assert.Equal(t, batch1.BatchNumber, batch.BatchNumber)
	assert.Equal(t, batch1.BatchHash, batch.BatchHash)
}

func TestBasicState_GetLastBatchNumber(t *testing.T) {
	batchNumber, err := state.GetLastBatchNumber(ctx)
	assert.NoError(t, err)
	assert.Equal(t, batch4.BatchNumber, batchNumber)
}

func TestBasicState_ConsolidateBatch(t *testing.T) {
	batchNumber := uint64(5)
	batch := &Batch{
		BatchNumber:        batchNumber,
		BatchHash:          common.HexToHash("0xaca7af32007b3d33d9d2342221093cd2fdae39ac29c170923c0519f0ca9b35bd"),
		BlockNumber:        blockNumber2,
		Sequencer:          addr,
		Aggregator:         addr,
		ConsolidatedTxHash: common.Hash{},
		Header:             nil,
		Uncles:             nil,
		Transactions:       nil,
		RawTxsData:         nil,
	}

	_, err := stateDb.Exec(ctx, "INSERT INTO batch (batch_num, batch_hash, block_num, sequencer, aggregator, consolidated_tx_hash) VALUES ($1, $2, $3, $4, $5, $6)",
		batch.BatchNumber, batch.BatchHash, batch.BlockNumber, batch.Sequencer, batch.Aggregator, batch.ConsolidatedTxHash)
	assert.NoError(t, err)

	insertedBatch, err := state.GetBatchByNumber(ctx, batchNumber)
	assert.NoError(t, err)
	assert.Equal(t, common.Hash{}, insertedBatch.ConsolidatedTxHash)

	err = state.ConsolidateBatch(ctx, batchNumber, consolidatedTxHash)
	assert.NoError(t, err)

	insertedBatch, err = state.GetBatchByNumber(ctx, batchNumber)
	assert.NoError(t, err)
	assert.Equal(t, consolidatedTxHash, insertedBatch.ConsolidatedTxHash)

	_, err = stateDb.Exec(ctx, "DELETE FROM batch WHERE batch_num = $1", batchNumber)
	assert.NoError(t, err)
}

func TestBasicState_GetTransactionCount(t *testing.T) {
	count, err := state.GetTransactionCount(ctx, addr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), count)
}

func TestBasicState_GetTxsByBatchNum(t *testing.T) {
	txs, err := state.GetTxsByBatchNum(ctx, batchNumber1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(txs))
}

func TestBasicState_GetTransactionByHash(t *testing.T) {
	tx, err := state.GetTransactionByHash(ctx, txHash)
	assert.NoError(t, err)
	assert.Equal(t, txHash, tx.Hash())
}
