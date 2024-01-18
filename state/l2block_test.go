package state_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestL2BlockHash(t *testing.T) {
	// create a geth header and block
	header := &types.Header{Number: big.NewInt(1)}
	ethBlock := types.NewBlockWithHeader(header)

	// create a l2 header and l2 block from geth header
	l2Header := state.NewL2Header(header)
	l2Block := state.NewL2BlockWithHeader(l2Header)

	// compare geth and l2 block hashes, they must match
	assert.Equal(t, ethBlock.Hash().String(), l2Block.Hash().String())
}

func TestGERPropagation(t *testing.T) {
	var dbTx pgx.Tx

	ctx := context.Background()
	batchNumber := uint64(1)

	type testCase struct {
		name        string
		block       *state.L2Block
		blockAdded  *state.L2Block
		expectedGER common.Hash
		prepare     func(*testCase) *mocks.StorageMock
	}
	testCases := []testCase{
		{
			name:        "add first block with GER different from zero",
			expectedGER: common.HexToHash("0x1"),
			prepare: func(tc *testCase) *mocks.StorageMock {
				storageMock := mocks.NewStorageMock(t)

				h := state.NewL2Header(&types.Header{Number: big.NewInt(1)})
				h.GlobalExitRoot = common.HexToHash("0x1")
				tc.block = state.NewL2Block(h, nil, nil, nil, &trie.StackTrie{})

				storageMock.
					On("AddL2Block", ctx, batchNumber, tc.block, []*types.Receipt(nil), []state.StoreTxEGPData(nil), dbTx).
					Return(nil).
					Run(func(args mock.Arguments) {
						tc.blockAdded = args.Get(2).(*state.L2Block)
					}).
					Once()

				return storageMock
			},
		},
		{
			name:        "add first block with GER zero",
			expectedGER: common.HexToHash("0x0"),
			prepare: func(tc *testCase) *mocks.StorageMock {
				storageMock := mocks.NewStorageMock(t)

				h := state.NewL2Header(&types.Header{Number: big.NewInt(1)})
				h.GlobalExitRoot = common.HexToHash("0x0")
				tc.block = state.NewL2Block(h, nil, nil, nil, &trie.StackTrie{})

				storageMock.
					On("GetLastL2Block", ctx, dbTx).
					Return(nil, state.ErrStateNotSynchronized).
					Once()

				storageMock.
					On("AddL2Block", ctx, batchNumber, tc.block, []*types.Receipt(nil), []state.StoreTxEGPData(nil), dbTx).
					Return(nil).
					Run(func(args mock.Arguments) {
						tc.blockAdded = args.Get(2).(*state.L2Block)
					}).
					Once()

				return storageMock
			},
		},
		{
			name:        "add block with GER different from zero",
			expectedGER: common.HexToHash("0x1"),
			prepare: func(tc *testCase) *mocks.StorageMock {
				storageMock := mocks.NewStorageMock(t)

				h := state.NewL2Header(&types.Header{Number: big.NewInt(2)})
				h.GlobalExitRoot = common.HexToHash("0x1")
				tc.block = state.NewL2Block(h, nil, nil, nil, &trie.StackTrie{})

				storageMock.
					On("AddL2Block", ctx, batchNumber, tc.block, []*types.Receipt(nil), []state.StoreTxEGPData(nil), dbTx).
					Return(nil).
					Run(func(args mock.Arguments) {
						tc.blockAdded = args.Get(2).(*state.L2Block)
					}).
					Once()

				return storageMock
			},
		},
		{
			name:        "add block with GER zero",
			expectedGER: common.HexToHash("0x1"),
			prepare: func(tc *testCase) *mocks.StorageMock {
				storageMock := mocks.NewStorageMock(t)

				h := state.NewL2Header(&types.Header{Number: big.NewInt(1)})
				h.GlobalExitRoot = common.HexToHash("0x1")
				latestBlock := state.NewL2Block(h, nil, nil, nil, &trie.StackTrie{})

				storageMock.
					On("GetLastL2Block", ctx, dbTx).
					Return(latestBlock, nil).
					Once()

				h = state.NewL2Header(&types.Header{Number: big.NewInt(2)})
				h.GlobalExitRoot = common.HexToHash("0x0")
				tc.block = state.NewL2Block(h, nil, nil, nil, &trie.StackTrie{})

				storageMock.
					On("AddL2Block", ctx, batchNumber, tc.block, []*types.Receipt(nil), []state.StoreTxEGPData(nil), dbTx).
					Return(nil).
					Run(func(args mock.Arguments) {
						tc.blockAdded = args.Get(2).(*state.L2Block)
					}).
					Once()

				return storageMock
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			storageMock := tc.prepare(&tc)
			st := state.NewState(state.Config{}, storageMock, nil, nil, nil, nil)

			err := st.AddL2Block(ctx, batchNumber, tc.block, nil, nil, dbTx)
			require.NoError(t, err)

			expectedGER := tc.expectedGER.String()
			ger := tc.blockAdded.GlobalExitRoot().String()

			assert.Equal(t, expectedGER, ger)
		})
	}
}
