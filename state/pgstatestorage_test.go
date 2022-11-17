package state_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	pgStateStorage *state.PostgresStorage
	block          = &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
)

func setup() {
	pgStateStorage = state.NewPostgresStorage(stateDb)
}

func TestGetBatchByL2BlockNumber(t *testing.T) {
	setup()
	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	err = testState.AddBlock(ctx, block, dbTx)
	assert.NoError(t, err)

	batchNumber := uint64(1)
	_, err = testState.PostgresStorage.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES ($1)", batchNumber)
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

	err = pgStateStorage.AddL2Block(ctx, batchNumber, l2Block, receipts, dbTx)
	require.NoError(t, err)
	result, err := pgStateStorage.GetBatchNumberOfL2Block(ctx, l2Block.Number().Uint64(), dbTx)
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

	_, err = testState.PostgresStorage.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (0)")
	require.NoError(t, err)
	_, err = testState.PostgresStorage.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (1)")
	require.NoError(t, err)
	_, err = testState.PostgresStorage.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (2)")
	require.NoError(t, err)
	_, err = testState.PostgresStorage.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (3)")
	require.NoError(t, err)
	_, err = testState.PostgresStorage.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (4)")
	require.NoError(t, err)
	_, err = testState.PostgresStorage.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (5)")
	require.NoError(t, err)
	_, err = testState.PostgresStorage.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (6)")
	require.NoError(t, err)
	_, err = testState.PostgresStorage.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (7)")
	require.NoError(t, err)
	_, err = testState.PostgresStorage.Exec(ctx, "INSERT INTO state.batch (batch_num) VALUES (8)")
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
		ToBatchNumber:   8,
	}
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
