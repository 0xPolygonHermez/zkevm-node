package jsonrpc

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBlockNumber(t *testing.T) {
	s, m, c := newMockedServer(t)
	defer s.Stop()

	testCases := []struct {
		name                string
		blockNumber         uint64
		error               error
		expectedBlockNumber uint64
	}{
		{"block number is zero", 0, nil, 0},
		{"block number is not zero", 5, nil, 5},
		{"failed to get block number", 5, errors.New("failed to get block number"), 0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m.State.
				On("GetLastBatchNumber", context.Background(), "").
				Return(testCase.blockNumber, testCase.error).
				Once()

			bn, err := c.BlockNumber(context.Background())
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedBlockNumber, bn)
		})
	}
}

func TestCall(t *testing.T) {
	s, m, c := newMockedServer(t)
	defer s.Stop()

	type testCase struct {
		name           string
		from           common.Address
		to             *common.Address
		gas            uint64
		gasPrice       *big.Int
		value          *big.Int
		data           []byte
		expectedResult []byte
		setupMocks     func(*mocks, *testCase)
	}

	testCases := []*testCase{
		{
			name:           "Transaction with all information",
			from:           common.HexToAddress("0x1"),
			to:             addressPtr(common.HexToAddress("0x2")),
			gas:            uint64(24000),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: []byte("hello world"),
			setupMocks: func(m *mocks, testCase *testCase) {
				batchNumber := uint64(1)
				txBundleID := ""
				batch := &state.Batch{Header: &types.Header{Root: common.Hash{}, GasLimit: 123456}}
				m.State.On("GetLastBatchNumber", context.Background(), txBundleID).Return(batchNumber, nil).Once()
				m.State.On("GetBatchByNumber", context.Background(), batchNumber, txBundleID).Return(batch, nil).Once()
				m.State.On("NewBatchProcessor", context.Background(), s.SequencerAddress, batch.Header.Root[:], txBundleID).Return(m.BatchProcessor, nil).Once()
				txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
					return tx != nil &&
						tx.Gas() == testCase.gas &&
						tx.To().Hex() == testCase.to.Hex() &&
						tx.GasPrice().Uint64() == testCase.gasPrice.Uint64() &&
						tx.Value().Uint64() == testCase.value.Uint64() &&
						hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data)
				})
				m.BatchProcessor.On("ProcessUnsignedTransaction", context.Background(), txMatchBy, testCase.from, s.SequencerAddress).Return(&runtime.ExecutionResult{ReturnValue: testCase.expectedResult}).Once()
			},
		},
		{
			name:           "Transaction without from and gas",
			to:             addressPtr(common.HexToAddress("0x2")),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: []byte("hello world"),
			setupMocks: func(m *mocks, testCase *testCase) {
				batchNumber := uint64(1)
				txBundleID := ""
				batch := &state.Batch{Header: &types.Header{Root: common.Hash{}, GasLimit: 123456}}
				m.State.On("GetLastBatchNumber", context.Background(), txBundleID).Return(batchNumber, nil).Once()
				m.State.On("GetBatchByNumber", context.Background(), batchNumber, txBundleID).Return(batch, nil).Once()
				m.State.On("NewBatchProcessor", context.Background(), s.SequencerAddress, batch.Header.Root[:], txBundleID).Return(m.BatchProcessor, nil).Once()
				txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
					return tx != nil &&
						tx.Gas() == batch.Header.GasLimit &&
						tx.To().Hex() == testCase.to.Hex() &&
						tx.GasPrice().Uint64() == testCase.gasPrice.Uint64() &&
						tx.Value().Uint64() == testCase.value.Uint64() &&
						hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data)
				})
				m.BatchProcessor.On("ProcessUnsignedTransaction", context.Background(), txMatchBy, testCase.from, s.SequencerAddress).Return(&runtime.ExecutionResult{ReturnValue: testCase.expectedResult}).Once()
				m.State.On("GetLastBatch", context.Background(), false, txBundleID).Return(batch, nil).Once()
			},
		},
		{
			name:           "Transaction without from and gas and failed to get last batch",
			to:             addressPtr(common.HexToAddress("0x2")),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: []byte{},
			setupMocks: func(m *mocks, testCase *testCase) {
				txBundleID := ""
				m.State.On("GetLastBatch", context.Background(), false, txBundleID).Return(nil, errors.New("failed to get last batch")).Once()
			},
		},
		{
			name:           "Transaction with gas but failed to get last batch number",
			from:           common.HexToAddress("0x1"),
			to:             addressPtr(common.HexToAddress("0x2")),
			gas:            uint64(24000),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: []byte{},
			setupMocks: func(m *mocks, testCase *testCase) {
				txBundleID := ""
				m.State.On("GetLastBatchNumber", context.Background(), txBundleID).Return(uint64(0), errors.New("failed to get last batch number")).Once()
			},
		},
		{
			name:           "Transaction with all information but failed to get batch by number",
			from:           common.HexToAddress("0x1"),
			to:             addressPtr(common.HexToAddress("0x2")),
			gas:            uint64(24000),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: []byte{},
			setupMocks: func(m *mocks, testCase *testCase) {
				batchNumber := uint64(1)
				txBundleID := ""
				m.State.On("GetLastBatchNumber", context.Background(), txBundleID).Return(batchNumber, nil).Once()
				m.State.On("GetBatchByNumber", context.Background(), batchNumber, txBundleID).Return(nil, errors.New("failed to get batch by number")).Once()
			},
		},
		{
			name:           "Transaction with all information but failed to create batch processor",
			from:           common.HexToAddress("0x1"),
			to:             addressPtr(common.HexToAddress("0x2")),
			gas:            uint64(24000),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: []byte{},
			setupMocks: func(m *mocks, testCase *testCase) {
				batchNumber := uint64(1)
				txBundleID := ""
				batch := &state.Batch{Header: &types.Header{Root: common.Hash{}, GasLimit: 123456}}
				m.State.On("GetLastBatchNumber", context.Background(), txBundleID).Return(batchNumber, nil).Once()
				m.State.On("GetBatchByNumber", context.Background(), batchNumber, txBundleID).Return(batch, nil).Once()
				m.State.On("NewBatchProcessor", context.Background(), s.SequencerAddress, batch.Header.Root[:], txBundleID).Return(nil, errors.New("failed to create batch processor")).Once()
			},
		},
		{
			name:           "Transaction with all information but failed to process unsigned transaction",
			from:           common.HexToAddress("0x1"),
			to:             addressPtr(common.HexToAddress("0x2")),
			gas:            uint64(24000),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: []byte{},
			setupMocks: func(m *mocks, testCase *testCase) {
				batchNumber := uint64(1)
				txBundleID := ""
				batch := &state.Batch{Header: &types.Header{Root: common.Hash{}, GasLimit: 123456}}
				m.State.On("GetLastBatchNumber", context.Background(), txBundleID).Return(batchNumber, nil).Once()
				m.State.On("GetBatchByNumber", context.Background(), batchNumber, txBundleID).Return(batch, nil).Once()
				m.State.On("NewBatchProcessor", context.Background(), s.SequencerAddress, batch.Header.Root[:], txBundleID).Return(m.BatchProcessor, nil).Once()
				txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
					return tx != nil &&
						tx.Gas() == testCase.gas &&
						tx.To().Hex() == testCase.to.Hex() &&
						tx.GasPrice().Uint64() == testCase.gasPrice.Uint64() &&
						tx.Value().Uint64() == testCase.value.Uint64() &&
						hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data)
				})
				m.BatchProcessor.On("ProcessUnsignedTransaction", context.Background(), txMatchBy, testCase.from, s.SequencerAddress).Return(&runtime.ExecutionResult{Err: errors.New("failed to process unsigned transaction")}).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			msg := ethereum.CallMsg{From: testCase.from, To: testCase.to, Gas: testCase.gas, GasPrice: testCase.gasPrice, Value: testCase.value, Data: testCase.data}

			testCase.setupMocks(m, testCase)

			result, err := c.CallContract(context.Background(), msg, nil)
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

func TestChainID(t *testing.T) {
	s, _, c := newMockedServer(t)
	defer s.Stop()

	chainID, err := c.ChainID(context.Background())
	require.NoError(t, err)

	assert.Equal(t, s.ChainID, chainID.Uint64())
}

func TestEstimateGas(t *testing.T) {
	s, m, c := newMockedServer(t)
	defer s.Stop()

	testCases := []struct {
		name     string
		from     common.Address
		to       *common.Address
		gas      uint64
		gasPrice *big.Int
		value    *big.Int
		data     []byte

		expectedResult uint64
	}{
		{
			name:           "Transaction with all information",
			from:           common.HexToAddress("0x1"),
			to:             addressPtr(common.HexToAddress("0x2")),
			gas:            uint64(24000),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: 100,
		},
		{
			name:           "Transaction without from and gas",
			to:             addressPtr(common.HexToAddress("0x2")),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: 100,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			msg := ethereum.CallMsg{From: testCase.from, To: testCase.to, Gas: testCase.gas, GasPrice: testCase.gasPrice, Value: testCase.value, Data: testCase.data}

			txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
				if tx == nil {
					return false
				}

				return tx.To().Hex() == testCase.to.Hex() &&
					tx.GasPrice().Uint64() == testCase.gasPrice.Uint64() &&
					tx.Value().Uint64() == testCase.value.Uint64() &&
					hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data)
			})

			m.State.
				On("EstimateGas", txMatchBy, testCase.from).
				Return(testCase.expectedResult, nil).
				Once()

			result, err := c.EstimateGas(context.Background(), msg)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

func TestGasPrice(t *testing.T) {
	s, m, c := newMockedServer(t)
	defer s.Stop()

	testCases := []struct {
		name             string
		gasPrice         *big.Int
		error            error
		expectedGasPrice uint64
	}{
		{"GasPrice nil", nil, nil, 0},
		{"GasPrice with value", big.NewInt(50), nil, 50},
		{"failed to get gas price", big.NewInt(50), errors.New("failed to get gas price"), 0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m.GasPriceEstimator.
				On("GetAvgGasPrice", context.Background()).
				Return(testCase.gasPrice, testCase.error).
				Once()

			gasPrice, err := c.SuggestGasPrice(context.Background())
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedGasPrice, gasPrice.Uint64())
		})
	}
}

func TestGetBalance(t *testing.T) {
	s, m, c := newMockedServer(t)
	defer s.Stop()

	type testCase struct {
		name            string
		addr            common.Address
		balance         *big.Int
		blockNumber     *big.Int
		expectedBalance uint64
		expectedError   *RPCError
		setupMocks      func(m *mocks, t *testCase)
	}

	testCases := []testCase{
		{
			name:            "get balance but failed to get latest block number",
			addr:            common.HexToAddress("0x123"),
			balance:         big.NewInt(1000),
			blockNumber:     nil,
			expectedBalance: 0,
			expectedError:   newRPCError(invalidParamsErrorCode, "invalid argument 1: failed to get last batch number"),
			setupMocks: func(m *mocks, t *testCase) {
				m.State.
					On("GetLastBatchNumber", context.Background(), "").
					Return(uint64(0), errors.New("failed to get last batch number")).
					Once()
			},
		},
		{
			name:            "get balance for block nil",
			addr:            common.HexToAddress("0x123"),
			balance:         big.NewInt(1000),
			blockNumber:     nil,
			expectedBalance: 1000,
			expectedError:   nil,
			setupMocks: func(m *mocks, t *testCase) {
				const lastBlockNumber = uint64(10)
				m.State.
					On("GetLastBatchNumber", context.Background(), "").
					Return(lastBlockNumber, nil).
					Once()

				m.State.
					On("GetBalance", context.Background(), t.addr, lastBlockNumber, "").
					Return(t.balance, nil).
					Once()
			},
		},
		{
			name:            "get balance for not found result",
			addr:            common.HexToAddress("0x123"),
			balance:         big.NewInt(1000),
			blockNumber:     nil,
			expectedBalance: 0,
			expectedError:   nil,
			setupMocks: func(m *mocks, t *testCase) {
				const lastBlockNumber = uint64(10)
				m.State.
					On("GetLastBatchNumber", context.Background(), "").
					Return(lastBlockNumber, nil).
					Once()

				m.State.
					On("GetBalance", context.Background(), t.addr, lastBlockNumber, "").
					Return(big.NewInt(0), state.ErrNotFound).
					Once()
			},
		},
		{
			name:            "get balance with state failure",
			addr:            common.HexToAddress("0x123"),
			balance:         big.NewInt(1000),
			blockNumber:     nil,
			expectedBalance: 0,
			expectedError:   newRPCError(defaultErrorCode, "failed to get balance from state"),
			setupMocks: func(m *mocks, t *testCase) {
				const lastBlockNumber = uint64(10)
				m.State.
					On("GetLastBatchNumber", context.Background(), "").
					Return(lastBlockNumber, nil).
					Once()

				m.State.
					On("GetBalance", context.Background(), t.addr, lastBlockNumber, "").
					Return(nil, errors.New("failed to get balance")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tc := testCase
			testCase.setupMocks(m, &tc)
			balance, err := c.BalanceAt(context.Background(), testCase.addr, testCase.blockNumber)
			assert.Equal(t, testCase.expectedBalance, balance.Uint64())
			if err != nil || testCase.expectedError != nil {
				rpcErr := err.(rpcError)
				assert.Equal(t, testCase.expectedError.ErrorCode(), rpcErr.ErrorCode())
				assert.Equal(t, testCase.expectedError.Error(), rpcErr.Error())
			}
		})
	}
}

func TestGetBlockByHash(t *testing.T) {
	type testCase struct {
		Name           string
		Hash           common.Hash
		ExpectedResult *types.Block
		ExpectedError  interface{}
		SetupMocks     func(*mocks, *testCase)
	}

	testCases := []testCase{
		{
			Name:           "Block not found",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  ethereum.NotFound,
			SetupMocks: func(m *mocks, tc *testCase) {
				m.State.
					On("GetBatchByHash", context.Background(), tc.Hash, "").
					Return(nil, state.ErrNotFound)
			},
		},
		{
			Name:           "Failed get block from state",
			Hash:           common.HexToHash("0x234"),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get block from state"),
			SetupMocks: func(m *mocks, tc *testCase) {
				m.State.
					On("GetBatchByHash", context.Background(), tc.Hash, "").
					Return(nil, errors.New("failed to get block from state")).
					Once()
			},
		},
		{
			Name: "get block successfully",
			Hash: common.HexToHash("0x234"),
			ExpectedResult: types.NewBlock(&types.Header{
				Number:    big.NewInt(1),
				UncleHash: types.EmptyUncleHash,
				Root:      types.EmptyRootHash,
			}, []*types.Transaction{
				types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{}),
			}, nil, []*types.Receipt{
				types.NewReceipt([]byte{}, false, uint64(0)),
			}, &trie.StackTrie{}),
			ExpectedError: nil,
			SetupMocks: func(m *mocks, tc *testCase) {
				batch := &state.Batch{
					Header: &types.Header{
						Number:    big.NewInt(1),
						UncleHash: types.EmptyUncleHash,
						Root:      types.EmptyRootHash,
					},
					Transactions: []*types.Transaction{
						types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{}),
					},
					Receipts: []*state.Receipt{
						{
							Receipt: *types.NewReceipt([]byte{}, false, uint64(0)),
						},
					},
				}

				m.State.
					On("GetBatchByHash", context.Background(), tc.Hash, "").
					Return(batch, nil).
					Once()
			},
		},
	}

	s, m, c := newMockedServer(t)
	defer s.Stop()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			testCase.SetupMocks(m, &testCase)

			result, err := c.BlockByHash(context.Background(), testCase.Hash)

			if result != nil || testCase.ExpectedResult != nil {
				assert.Equal(t, testCase.ExpectedResult.Number().Uint64(), result.Number().Uint64())
				assert.Equal(t, len(testCase.ExpectedResult.Transactions()), len(result.Transactions()))
			}

			if err != nil || testCase.ExpectedError != nil {
				if expectedErr, ok := testCase.ExpectedError.(*RPCError); ok {
					rpcErr := err.(rpcError)
					assert.Equal(t, expectedErr.ErrorCode(), rpcErr.ErrorCode())
					assert.Equal(t, expectedErr.Error(), rpcErr.Error())
				} else {
					assert.Equal(t, testCase.ExpectedError, err)
				}
			}
		})
	}
}

func TestGetTransactionReceipt(t *testing.T) {
	s, m, c := newMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name string
		Hash common.Hash

		ExpectedResult *types.Receipt
		ExpectedError  error

		SetupMocks func(m *mocks, t testCase)
	}

	testCases := []testCase{
		{
			Name:           "tx receipt not found",
			Hash:           common.HexToHash("0x1"),
			ExpectedResult: nil,
			ExpectedError:  ethereum.NotFound,
			SetupMocks: func(m *mocks, t testCase) {
				m.State.
					On("GetTransactionReceipt", context.Background(), t.Hash, "").
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			testCase.SetupMocks(m, testCase)
			r, err := c.TransactionReceipt(context.Background(), testCase.Hash)
			assert.Equal(t, testCase.ExpectedResult, r)
			assert.Equal(t, testCase.ExpectedError, err)
		})
	}
}

func addressPtr(i common.Address) *common.Address {
	return &i
}
