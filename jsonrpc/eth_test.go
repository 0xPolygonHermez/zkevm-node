package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBlockNumber(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult uint64
		ExpectedError  interface{}
		SetupMocks     func(m *mocks)
	}

	testCases := []testCase{
		{
			Name:           "get block number successfully",
			ExpectedError:  nil,
			ExpectedResult: 10,
			SetupMocks: func(m *mocks) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get block number",
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last block number from state"),
			ExpectedResult: 0,
			SetupMocks: func(m *mocks) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			result, err := c.BlockNumber(context.Background())
			assert.Equal(t, testCase.ExpectedResult, result)

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

func TestCall(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		name           string
		from           common.Address
		to             *common.Address
		gas            uint64
		gasPrice       *big.Int
		value          *big.Int
		data           []byte
		blockNumber    *big.Int
		expectedResult []byte
		expectedError  interface{}
		setupMocks     func(Config, *mocks, *testCase)
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
			expectedError:  nil,
			setupMocks: func(c Config, m *mocks, testCase *testCase) {
				blockNumber := uint64(1)
				nonce := uint64(7)
				m.DbTx.On("Commit", context.Background()).Return(nil).Once()
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Once()
				m.State.On("GetLastL2BlockNumber", context.Background(), m.DbTx).Return(blockNumber, nil).Once()
				txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
					return tx != nil &&
						tx.Gas() == testCase.gas &&
						tx.To().Hex() == testCase.to.Hex() &&
						tx.GasPrice().Uint64() == testCase.gasPrice.Uint64() &&
						tx.Value().Uint64() == testCase.value.Uint64() &&
						hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data) &&
						tx.Nonce() == nonce
				})
				m.State.On("GetNonce", context.Background(), testCase.from, blockNumber, m.DbTx).Return(nonce, nil).Once()
				var nilBlockNumber *uint64
				m.State.On("ProcessUnsignedTransaction", context.Background(), txMatchBy, testCase.from, nilBlockNumber, true, m.DbTx).Return(&runtime.ExecutionResult{ReturnValue: testCase.expectedResult}).Once()
			},
		},
		{
			name:           "Transaction without from and gas from latest block",
			to:             addressPtr(common.HexToAddress("0x2")),
			gasPrice:       big.NewInt(0),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: []byte("hello world"),
			expectedError:  nil,
			setupMocks: func(c Config, m *mocks, testCase *testCase) {
				blockNumber := uint64(1)
				block := types.NewBlockWithHeader(&types.Header{Root: common.Hash{}, GasLimit: s.Config.MaxCumulativeGasUsed})
				m.DbTx.On("Commit", context.Background()).Return(nil).Once()
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Once()
				m.State.On("GetLastL2BlockNumber", context.Background(), m.DbTx).Return(blockNumber, nil).Once()
				m.State.On("GetLastL2Block", context.Background(), m.DbTx).Return(block, nil).Once()
				txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
					hasTx := tx != nil
					gasMatch := tx.Gas() == block.Header().GasLimit
					toMatch := tx.To().Hex() == testCase.to.Hex()
					gasPriceMatch := tx.GasPrice().Uint64() == testCase.gasPrice.Uint64()
					valueMatch := tx.Value().Uint64() == testCase.value.Uint64()
					dataMatch := hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data)
					return hasTx && gasMatch && toMatch && gasPriceMatch && valueMatch && dataMatch
				})
				var nilBlockNumber *uint64
				m.State.On("ProcessUnsignedTransaction", context.Background(), txMatchBy, common.HexToAddress(c.DefaultSenderAddress), nilBlockNumber, true, m.DbTx).Return(&runtime.ExecutionResult{ReturnValue: testCase.expectedResult}).Once()
			},
		},
		{
			name:           "Transaction without from and gas from pending block",
			to:             addressPtr(common.HexToAddress("0x2")),
			gasPrice:       big.NewInt(0),
			value:          big.NewInt(2),
			data:           []byte("data"),
			blockNumber:    big.NewInt(-1),
			expectedResult: []byte("hello world"),
			expectedError:  nil,
			setupMocks: func(c Config, m *mocks, testCase *testCase) {
				blockNumber := uint64(1)
				block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1), Root: common.Hash{}, GasLimit: s.Config.MaxCumulativeGasUsed})
				m.DbTx.On("Commit", context.Background()).Return(nil).Once()
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Once()
				m.State.On("GetLastL2BlockNumber", context.Background(), m.DbTx).Return(blockNumber, nil).Once()
				m.State.On("GetLastL2Block", context.Background(), m.DbTx).Return(block, nil).Once()
				txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
					hasTx := tx != nil
					gasMatch := tx.Gas() == block.Header().GasLimit
					toMatch := tx.To().Hex() == testCase.to.Hex()
					gasPriceMatch := tx.GasPrice().Uint64() == testCase.gasPrice.Uint64()
					valueMatch := tx.Value().Uint64() == testCase.value.Uint64()
					dataMatch := hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data)
					return hasTx && gasMatch && toMatch && gasPriceMatch && valueMatch && dataMatch
				})
				var nilBlockNumber *uint64
				m.State.On("ProcessUnsignedTransaction", context.Background(), txMatchBy, common.HexToAddress(c.DefaultSenderAddress), nilBlockNumber, true, m.DbTx).Return(&runtime.ExecutionResult{ReturnValue: testCase.expectedResult}).Once()
			},
		},
		{
			name:           "Transaction without from and gas and failed to get latest block header",
			to:             addressPtr(common.HexToAddress("0x2")),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: nil,
			expectedError:  newRPCError(defaultErrorCode, "failed to get block header"),
			setupMocks: func(c Config, m *mocks, testCase *testCase) {
				m.DbTx.On("Rollback", context.Background()).Return(nil).Once()
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Once()
				m.State.On("GetLastL2Block", context.Background(), m.DbTx).Return(nil, errors.New("failed to get last block")).Once()
			},
		},
		{
			name:           "Transaction without from and gas and failed to get pending block header",
			to:             addressPtr(common.HexToAddress("0x2")),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			blockNumber:    big.NewInt(-1),
			expectedResult: nil,
			expectedError:  newRPCError(defaultErrorCode, "failed to get block header"),
			setupMocks: func(c Config, m *mocks, testCase *testCase) {
				m.DbTx.On("Rollback", context.Background()).Return(nil).Once()
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Once()
				m.State.On("GetLastL2Block", context.Background(), m.DbTx).Return(nil, errors.New("failed to get last block")).Once()
			},
		},
		{
			name:           "Transaction with gas but failed to get last block number",
			from:           common.HexToAddress("0x1"),
			to:             addressPtr(common.HexToAddress("0x2")),
			gas:            uint64(24000),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: nil,
			expectedError:  newRPCError(defaultErrorCode, "failed to get the last block number from state"),
			setupMocks: func(c Config, m *mocks, testCase *testCase) {
				m.DbTx.On("Rollback", context.Background()).Return(nil).Once()
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Once()
				m.State.On("GetLastL2BlockNumber", context.Background(), m.DbTx).Return(uint64(0), errors.New("failed to get last block number")).Once()
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
			expectedResult: nil,
			expectedError:  newRPCError(defaultErrorCode, "failed to process unsigned transaction"),
			setupMocks: func(c Config, m *mocks, testCase *testCase) {
				blockNumber := uint64(1)
				nonce := uint64(7)
				m.DbTx.On("Rollback", context.Background()).Return(nil).Once()
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Once()
				m.State.On("GetLastL2BlockNumber", context.Background(), m.DbTx).Return(blockNumber, nil).Once()
				txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
					hasTx := tx != nil
					gasMatch := tx.Gas() == testCase.gas
					toMatch := tx.To().Hex() == testCase.to.Hex()
					gasPriceMatch := tx.GasPrice().Uint64() == testCase.gasPrice.Uint64()
					valueMatch := tx.Value().Uint64() == testCase.value.Uint64()
					dataMatch := hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data)
					nonceMatch := tx.Nonce() == nonce
					return hasTx && gasMatch && toMatch && gasPriceMatch && valueMatch && dataMatch && nonceMatch
				})
				m.State.On("GetNonce", context.Background(), testCase.from, blockNumber, m.DbTx).Return(nonce, nil).Once()
				var nilBlockNumber *uint64
				m.State.On("ProcessUnsignedTransaction", context.Background(), txMatchBy, testCase.from, nilBlockNumber, true, m.DbTx).Return(&runtime.ExecutionResult{Err: errors.New("failed to process unsigned transaction")}).Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			msg := ethereum.CallMsg{From: testCase.from, To: testCase.to, Gas: testCase.gas, GasPrice: testCase.gasPrice, Value: testCase.value, Data: testCase.data}

			testCase.setupMocks(s.Config, m, testCase)

			result, err := c.CallContract(context.Background(), msg, testCase.blockNumber)
			assert.Equal(t, testCase.expectedResult, result)
			if err != nil || testCase.expectedError != nil {
				if expectedErr, ok := testCase.expectedError.(*RPCError); ok {
					rpcErr := err.(rpcError)
					assert.Equal(t, expectedErr.ErrorCode(), rpcErr.ErrorCode())
					assert.Equal(t, expectedErr.Error(), rpcErr.Error())
				} else {
					assert.Equal(t, testCase.expectedError, err)
				}
			}
		})
	}
}

func TestChainID(t *testing.T) {
	s, _, c := newSequencerMockedServer(t)
	defer s.Stop()

	chainID, err := c.ChainID(context.Background())
	require.NoError(t, err)

	assert.Equal(t, s.Config.ChainID, chainID.Uint64())
}

func TestEstimateGas(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		name       string
		from       common.Address
		to         *common.Address
		gas        uint64
		gasPrice   *big.Int
		value      *big.Int
		data       []byte
		setupMocks func(Config, *mocks, *testCase)

		expectedResult uint64
	}

	testCases := []testCase{
		{
			name:           "Transaction with all information",
			from:           common.HexToAddress("0x1"),
			to:             addressPtr(common.HexToAddress("0x2")),
			gas:            uint64(24000),
			gasPrice:       big.NewInt(1),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: 100,
			setupMocks: func(c Config, m *mocks, testCase *testCase) {
				blockNumber := uint64(10)
				nonce := uint64(7)
				txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
					if tx == nil {
						return false
					}

					matchTo := tx.To().Hex() == testCase.to.Hex()
					matchGasPrice := tx.GasPrice().Uint64() == testCase.gasPrice.Uint64()
					matchValue := tx.Value().Uint64() == testCase.value.Uint64()
					matchData := hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data)
					matchNonce := tx.Nonce() == nonce
					return matchTo && matchGasPrice && matchValue && matchData && matchNonce
				})

				m.DbTx.On("Commit", context.Background()).Return(nil).Once()
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(blockNumber, nil).
					Once()

				m.State.
					On("GetNonce", context.Background(), testCase.from, blockNumber, m.DbTx).
					Return(nonce, nil).
					Once()

				var nilBlockNumber *uint64
				m.State.
					On("EstimateGas", txMatchBy, testCase.from, nilBlockNumber, m.DbTx).
					Return(testCase.expectedResult, nil).
					Once()
			},
		},
		{
			name:           "Transaction without from and gas",
			to:             addressPtr(common.HexToAddress("0x2")),
			gasPrice:       big.NewInt(0),
			value:          big.NewInt(2),
			data:           []byte("data"),
			expectedResult: 100,
			setupMocks: func(c Config, m *mocks, testCase *testCase) {
				blockNumber := uint64(9)
				nonce := uint64(0)
				txMatchBy := mock.MatchedBy(func(tx *types.Transaction) bool {
					if tx == nil {
						return false
					}

					matchTo := tx.To().Hex() == testCase.to.Hex()
					matchGasPrice := tx.GasPrice().Uint64() == testCase.gasPrice.Uint64()
					matchValue := tx.Value().Uint64() == testCase.value.Uint64()
					matchData := hex.EncodeToHex(tx.Data()) == hex.EncodeToHex(testCase.data)
					matchNonce := tx.Nonce() == nonce
					return matchTo && matchGasPrice && matchValue && matchData && matchNonce
				})

				m.DbTx.On("Commit", context.Background()).Return(nil).Once()
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(blockNumber, nil).
					Once()

				var nilBlockNumber *uint64
				m.State.
					On("EstimateGas", txMatchBy, common.HexToAddress(c.DefaultSenderAddress), nilBlockNumber, m.DbTx).
					Return(testCase.expectedResult, nil).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tc := testCase
			tc.setupMocks(s.Config, m, &tc)

			msg := ethereum.CallMsg{From: testCase.from, To: testCase.to, Gas: testCase.gas, GasPrice: testCase.gasPrice, Value: testCase.value, Data: testCase.data}
			result, err := c.EstimateGas(context.Background(), msg)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

func TestGasPrice(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
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
	s, m, c := newSequencerMockedServer(t)
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
			expectedError:   newRPCError(defaultErrorCode, "failed to get the last block number from state"),
			setupMocks: func(m *mocks, t *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number")).
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
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(lastBlockNumber, nil).
					Once()

				m.State.
					On("GetBalance", context.Background(), t.addr, lastBlockNumber, m.DbTx).
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
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(lastBlockNumber, nil).
					Once()

				m.State.
					On("GetBalance", context.Background(), t.addr, lastBlockNumber, m.DbTx).
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
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(lastBlockNumber, nil).
					Once()

				m.State.
					On("GetBalance", context.Background(), t.addr, lastBlockNumber, m.DbTx).
					Return(nil, errors.New("failed to get balance")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tc := testCase
			testCase.setupMocks(m, &tc)
			balance, err := c.BalanceAt(context.Background(), tc.addr, tc.blockNumber)
			assert.Equal(t, tc.expectedBalance, balance.Uint64())
			if err != nil || tc.expectedError != nil {
				rpcErr := err.(rpcError)
				assert.Equal(t, tc.expectedError.ErrorCode(), rpcErr.ErrorCode())
				assert.Equal(t, tc.expectedError.Error(), rpcErr.Error())
			}
		})
	}
}

func TestGetL2BlockByHash(t *testing.T) {
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
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockByHash", context.Background(), tc.Hash, m.DbTx).
					Return(nil, state.ErrNotFound)
			},
		},
		{
			Name:           "Failed get block from state",
			Hash:           common.HexToHash("0x234"),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get block by hash from state"),
			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockByHash", context.Background(), tc.Hash, m.DbTx).
					Return(nil, errors.New("failed to get block from state")).
					Once()
			},
		},
		{
			Name: "get block successfully",
			Hash: common.HexToHash("0x345"),
			ExpectedResult: types.NewBlock(
				&types.Header{Number: big.NewInt(1), UncleHash: types.EmptyUncleHash, Root: types.EmptyRootHash},
				[]*types.Transaction{types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})},
				nil,
				[]*types.Receipt{types.NewReceipt([]byte{}, false, uint64(0))},
				&trie.StackTrie{},
			),
			ExpectedError: nil,
			SetupMocks: func(m *mocks, tc *testCase) {
				block := types.NewBlock(types.CopyHeader(tc.ExpectedResult.Header()), tc.ExpectedResult.Transactions(), tc.ExpectedResult.Uncles(), []*types.Receipt{types.NewReceipt([]byte{}, false, uint64(0))}, &trie.StackTrie{})

				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockByHash", context.Background(), tc.Hash, m.DbTx).
					Return(block, nil).
					Once()
			},
		},
	}

	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			testCase.SetupMocks(m, &tc)

			result, err := c.BlockByHash(context.Background(), tc.Hash)

			if result != nil || tc.ExpectedResult != nil {
				assert.Equal(t, tc.ExpectedResult.Number().Uint64(), result.Number().Uint64())
				assert.Equal(t, len(tc.ExpectedResult.Transactions()), len(result.Transactions()))
				assert.Equal(t, tc.ExpectedResult.Hash(), result.Hash())
			}

			if err != nil || tc.ExpectedError != nil {
				if expectedErr, ok := tc.ExpectedError.(*RPCError); ok {
					rpcErr := err.(rpcError)
					assert.Equal(t, expectedErr.ErrorCode(), rpcErr.ErrorCode())
					assert.Equal(t, expectedErr.Error(), rpcErr.Error())
				} else {
					assert.Equal(t, tc.ExpectedError, err)
				}
			}
		})
	}
}

func TestGetL2BlockByNumber(t *testing.T) {
	type testCase struct {
		Name           string
		Number         *big.Int
		ExpectedResult *types.Block
		ExpectedError  interface{}
		SetupMocks     func(*mocks, *testCase)
	}

	testCases := []testCase{
		{
			Name:           "Block not found",
			Number:         big.NewInt(123),
			ExpectedResult: nil,
			ExpectedError:  ethereum.NotFound,
			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockByNumber", context.Background(), tc.Number.Uint64(), m.DbTx).
					Return(nil, state.ErrNotFound)
			},
		},
		{
			Name:   "get specific block successfully",
			Number: big.NewInt(345),
			ExpectedResult: types.NewBlock(
				&types.Header{Number: big.NewInt(1), UncleHash: types.EmptyUncleHash, Root: types.EmptyRootHash},
				[]*types.Transaction{types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})},
				nil,
				[]*types.Receipt{types.NewReceipt([]byte{}, false, uint64(0))},
				&trie.StackTrie{},
			),
			ExpectedError: nil,
			SetupMocks: func(m *mocks, tc *testCase) {
				block := types.NewBlock(types.CopyHeader(tc.ExpectedResult.Header()), tc.ExpectedResult.Transactions(),
					tc.ExpectedResult.Uncles(), []*types.Receipt{types.NewReceipt([]byte{}, false, uint64(0))}, &trie.StackTrie{})

				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockByNumber", context.Background(), tc.Number.Uint64(), m.DbTx).
					Return(block, nil).
					Once()
			},
		},
		{
			Name:   "get latest block successfully",
			Number: nil,
			ExpectedResult: types.NewBlock(
				&types.Header{Number: big.NewInt(2), UncleHash: types.EmptyUncleHash, Root: types.EmptyRootHash},
				[]*types.Transaction{types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})},
				nil,
				[]*types.Receipt{types.NewReceipt([]byte{}, false, uint64(0))},
				&trie.StackTrie{},
			),
			ExpectedError: nil,
			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(tc.ExpectedResult.Number().Uint64(), nil).
					Once()

				m.State.
					On("GetL2BlockByNumber", context.Background(), tc.ExpectedResult.Number().Uint64(), m.DbTx).
					Return(tc.ExpectedResult, nil).
					Once()
			},
		},
		{
			Name:           "get latest block fails to compute block number",
			Number:         nil,
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last block number from state"),
			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number")).
					Once()
			},
		},
		{
			Name:           "get latest block fails to load block by number",
			Number:         nil,
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "couldn't load block from state by number 1"),
			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(1), nil).
					Once()

				m.State.
					On("GetL2BlockByNumber", context.Background(), uint64(1), m.DbTx).
					Return(nil, errors.New("failed to load block by number")).
					Once()
			},
		},
		{
			Name:           "get pending block successfully",
			Number:         big.NewInt(-1),
			ExpectedResult: types.NewBlock(&types.Header{Number: big.NewInt(2)}, nil, nil, nil, &trie.StackTrie{}),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc *testCase) {
				lastBlockHeader := types.CopyHeader(tc.ExpectedResult.Header())
				lastBlockHeader.Number.Sub(lastBlockHeader.Number, big.NewInt(1))
				lastBlock := types.NewBlock(lastBlockHeader, nil, nil, nil, &trie.StackTrie{})

				expectedResultHeader := types.CopyHeader(tc.ExpectedResult.Header())
				expectedResultHeader.ParentHash = lastBlock.Hash()
				tc.ExpectedResult = types.NewBlock(expectedResultHeader, nil, nil, nil, &trie.StackTrie{})

				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2Block", context.Background(), m.DbTx).
					Return(lastBlock, nil).
					Once()
			},
		},
		{
			Name:           "get pending block fails",
			Number:         big.NewInt(-1),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "couldn't load last block from state to compute the pending block"),
			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2Block", context.Background(), m.DbTx).
					Return(nil, errors.New("failed to load last block")).
					Once()
			},
		},
	}

	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			testCase.SetupMocks(m, &tc)

			result, err := c.BlockByNumber(context.Background(), tc.Number)

			if result != nil || tc.ExpectedResult != nil {
				expectedResultJSON, _ := json.Marshal(tc.ExpectedResult.Header())
				resultJSON, _ := json.Marshal(result.Header())

				expectedResultJSONStr := string(expectedResultJSON)
				resultJSONStr := string(resultJSON)

				assert.JSONEq(t, expectedResultJSONStr, resultJSONStr)
				assert.Equal(t, tc.ExpectedResult.Number().Uint64(), result.Number().Uint64())
				assert.Equal(t, len(tc.ExpectedResult.Transactions()), len(result.Transactions()))
				assert.Equal(t, tc.ExpectedResult.Hash(), result.Hash())
			}

			if err != nil || tc.ExpectedError != nil {
				if expectedErr, ok := tc.ExpectedError.(*RPCError); ok {
					rpcErr := err.(rpcError)
					assert.Equal(t, expectedErr.ErrorCode(), rpcErr.ErrorCode())
					assert.Equal(t, expectedErr.Error(), rpcErr.Error())
				} else {
					assert.Equal(t, tc.ExpectedError, err)
				}
			}
		})
	}
}

func TestGetUncleByBlockHashAndIndex(t *testing.T) {
	s, _, _ := newSequencerMockedServer(t)
	defer s.Stop()

	res, err := s.JSONRPCCall("eth_getUncleByBlockHashAndIndex", common.HexToHash("0x123").Hex(), "0x1")
	require.NoError(t, err)

	assert.Equal(t, float64(1), res.ID)
	assert.Equal(t, "2.0", res.JSONRPC)
	assert.Nil(t, res.Error)

	var result interface{}
	err = json.Unmarshal(res.Result, &result)
	require.NoError(t, err)

	assert.Nil(t, result)
}

func TestGetUncleByBlockNumberAndIndex(t *testing.T) {
	s, _, _ := newSequencerMockedServer(t)
	defer s.Stop()

	res, err := s.JSONRPCCall("eth_getUncleByBlockNumberAndIndex", "0x123", "0x1")
	require.NoError(t, err)

	assert.Equal(t, float64(1), res.ID)
	assert.Equal(t, "2.0", res.JSONRPC)
	assert.Nil(t, res.Error)

	var result interface{}
	err = json.Unmarshal(res.Result, &result)
	require.NoError(t, err)

	assert.Nil(t, result)
}

func TestGetUncleCountByBlockHash(t *testing.T) {
	s, _, _ := newSequencerMockedServer(t)
	defer s.Stop()

	res, err := s.JSONRPCCall("eth_getUncleCountByBlockHash", common.HexToHash("0x123"))
	require.NoError(t, err)

	assert.Equal(t, float64(1), res.ID)
	assert.Equal(t, "2.0", res.JSONRPC)
	assert.Nil(t, res.Error)

	var result argUint64
	err = json.Unmarshal(res.Result, &result)
	require.NoError(t, err)

	assert.Equal(t, uint64(0), uint64(result))
}

func TestGetUncleCountByBlockNumber(t *testing.T) {
	s, _, _ := newSequencerMockedServer(t)
	defer s.Stop()

	res, err := s.JSONRPCCall("eth_getUncleCountByBlockNumber", "0x123")
	require.NoError(t, err)

	assert.Equal(t, float64(1), res.ID)
	assert.Equal(t, "2.0", res.JSONRPC)
	assert.Nil(t, res.Error)

	var result argUint64
	err = json.Unmarshal(res.Result, &result)
	require.NoError(t, err)

	assert.Equal(t, uint64(0), uint64(result))
}

func TestGetCode(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		Addr           common.Address
		BlockNumber    *big.Int
		ExpectedResult []byte
		ExpectedError  interface{}

		SetupMocks func(m *mocks, tc *testCase)
	}

	testCases := []testCase{
		{
			Name:           "failed to identify the block",
			Addr:           common.HexToAddress("0x123"),
			BlockNumber:    nil,
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last block number from state"),

			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number")).
					Once()
			},
		},
		{
			Name:           "failed to get code",
			Addr:           common.HexToAddress("0x123"),
			BlockNumber:    big.NewInt(1),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get code"),

			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetCode", context.Background(), tc.Addr, tc.BlockNumber.Uint64(), m.DbTx).
					Return(nil, errors.New("failed to get code")).
					Once()
			},
		},
		{
			Name:           "code not found",
			Addr:           common.HexToAddress("0x123"),
			BlockNumber:    big.NewInt(1),
			ExpectedResult: []byte{},
			ExpectedError:  nil,

			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetCode", context.Background(), tc.Addr, tc.BlockNumber.Uint64(), m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "get code successfully",
			Addr:           common.HexToAddress("0x123"),
			BlockNumber:    big.NewInt(1),
			ExpectedResult: []byte{1, 2, 3},
			ExpectedError:  nil,

			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetCode", context.Background(), tc.Addr, tc.BlockNumber.Uint64(), m.DbTx).
					Return(tc.ExpectedResult, nil).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, &tc)
			result, err := c.CodeAt(context.Background(), tc.Addr, tc.BlockNumber)
			assert.Equal(t, tc.ExpectedResult, result)

			if err != nil || tc.ExpectedError != nil {
				if expectedErr, ok := tc.ExpectedError.(*RPCError); ok {
					rpcErr := err.(rpcError)
					assert.Equal(t, expectedErr.ErrorCode(), rpcErr.ErrorCode())
					assert.Equal(t, expectedErr.Error(), rpcErr.Error())
				} else {
					assert.Equal(t, tc.ExpectedError, err)
				}
			}
		})
	}
}

func TestGetStorageAt(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		Addr           common.Address
		Key            common.Hash
		BlockNumber    *big.Int
		ExpectedResult []byte
		ExpectedError  interface{}

		SetupMocks func(m *mocks, tc *testCase)
	}

	testCases := []testCase{
		{
			Name:           "failed to identify the block",
			Addr:           common.HexToAddress("0x123"),
			Key:            common.HexToHash("0x123"),
			BlockNumber:    nil,
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last block number from state"),

			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number")).
					Once()
			},
		},
		{
			Name:           "failed to get storage at",
			Addr:           common.HexToAddress("0x123"),
			Key:            common.HexToHash("0x123"),
			BlockNumber:    big.NewInt(1),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get storage value from state"),

			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetStorageAt", context.Background(), tc.Addr, tc.Key.Big(), tc.BlockNumber.Uint64(), m.DbTx).
					Return(nil, errors.New("failed to get storage at")).
					Once()
			},
		},
		{
			Name:           "code not found",
			Addr:           common.HexToAddress("0x123"),
			Key:            common.HexToHash("0x123"),
			BlockNumber:    big.NewInt(1),
			ExpectedResult: common.Hash{}.Bytes(),
			ExpectedError:  nil,

			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetStorageAt", context.Background(), tc.Addr, tc.Key.Big(), tc.BlockNumber.Uint64(), m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "get code successfully",
			Addr:           common.HexToAddress("0x123"),
			Key:            common.HexToHash("0x123"),
			BlockNumber:    big.NewInt(1),
			ExpectedResult: common.BigToHash(big.NewInt(123)).Bytes(),
			ExpectedError:  nil,

			SetupMocks: func(m *mocks, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetStorageAt", context.Background(), tc.Addr, tc.Key.Big(), tc.BlockNumber.Uint64(), m.DbTx).
					Return(big.NewInt(123), nil).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, &tc)
			result, err := c.StorageAt(context.Background(), tc.Addr, tc.Key, tc.BlockNumber)
			assert.Equal(t, tc.ExpectedResult, result)

			if err != nil || tc.ExpectedError != nil {
				if expectedErr, ok := tc.ExpectedError.(*RPCError); ok {
					rpcErr := err.(rpcError)
					assert.Equal(t, expectedErr.ErrorCode(), rpcErr.ErrorCode())
					assert.Equal(t, expectedErr.Error(), rpcErr.Error())
				} else {
					assert.Equal(t, tc.ExpectedError, err)
				}
			}
		})
	}
}

func TestGetCompilers(t *testing.T) {
	s, _, _ := newSequencerMockedServer(t)
	defer s.Stop()

	res, err := s.JSONRPCCall("eth_getCompilers")
	require.NoError(t, err)

	assert.Equal(t, float64(1), res.ID)
	assert.Equal(t, "2.0", res.JSONRPC)
	assert.Nil(t, res.Error)

	var result []interface{}
	err = json.Unmarshal(res.Result, &result)
	require.NoError(t, err)

	assert.Equal(t, 0, len(result))
}

func TestSyncing(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult *ethereum.SyncProgress
		ExpectedError  rpcError
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "failed to get syncing information",
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get syncing info from state"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetSyncingInfo", context.Background(), m.DbTx).
					Return(state.SyncingInfo{}, errors.New("failed to get syncing info from state")).
					Once()
			},
		},
		{
			Name:           "get syncing information successfully while syncing",
			ExpectedResult: &ethereum.SyncProgress{StartingBlock: 1, CurrentBlock: 2, HighestBlock: 3},
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetSyncingInfo", context.Background(), m.DbTx).
					Return(state.SyncingInfo{InitialSyncingBlock: 1, CurrentBlockNumber: 2, LastBlockNumberSeen: 3, LastBlockNumberConsolidated: 3}, nil).
					Once()
			},
		},
		{
			Name:           "get syncing information successfully when synced",
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetSyncingInfo", context.Background(), m.DbTx).
					Return(state.SyncingInfo{InitialSyncingBlock: 1, CurrentBlockNumber: 1, LastBlockNumberSeen: 1, LastBlockNumberConsolidated: 1}, nil).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			testCase.SetupMocks(m, testCase)
			result, err := c.SyncProgress(context.Background())

			if result != nil || testCase.ExpectedResult != nil {
				assert.Equal(t, testCase.ExpectedResult.StartingBlock, result.StartingBlock)
				assert.Equal(t, testCase.ExpectedResult.CurrentBlock, result.CurrentBlock)
				assert.Equal(t, testCase.ExpectedResult.HighestBlock, result.HighestBlock)
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

func TestGetTransactionL2onByBlockHashAndIndex(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name  string
		Hash  common.Hash
		Index uint

		ExpectedResult *types.Transaction
		ExpectedError  interface{}
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "Get Tx Successfully",
			Hash:           common.HexToHash("0x999"),
			Index:          uint(1),
			ExpectedResult: types.NewTransaction(1, common.HexToAddress("0x111"), big.NewInt(2), 3, big.NewInt(4), []byte{5, 6, 7, 8}),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				tx := tc.ExpectedResult
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByL2BlockHashAndIndex", context.Background(), tc.Hash, uint64(tc.Index), m.DbTx).
					Return(tx, nil).
					Once()

				receipt := types.NewReceipt([]byte{}, false, 0)
				receipt.BlockHash = common.Hash{}
				receipt.BlockNumber = big.NewInt(1)
				receipt.TransactionIndex = tc.Index

				m.State.
					On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
					Return(receipt, nil).
					Once()
			},
		},
		{
			Name:           "Tx not found",
			Hash:           common.HexToHash("0x999"),
			Index:          uint(1),
			ExpectedResult: nil,
			ExpectedError:  ethereum.NotFound,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByL2BlockHashAndIndex", context.Background(), tc.Hash, uint64(tc.Index), m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "Get Tx fail to get tx from state",
			Hash:           common.HexToHash("0x999"),
			Index:          uint(1),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get transaction"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByL2BlockHashAndIndex", context.Background(), tc.Hash, uint64(tc.Index), m.DbTx).
					Return(nil, errors.New("failed to get transaction by block and index from state")).
					Once()
			},
		},
		{
			Name:           "Tx found but receipt not found",
			Hash:           common.HexToHash("0x999"),
			Index:          uint(1),
			ExpectedResult: nil,
			ExpectedError:  ethereum.NotFound,
			SetupMocks: func(m *mocks, tc testCase) {
				tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), []byte{})
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByL2BlockHashAndIndex", context.Background(), tc.Hash, uint64(tc.Index), m.DbTx).
					Return(tx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "Get Tx fail to get tx receipt from state",
			Hash:           common.HexToHash("0x999"),
			Index:          uint(1),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get transaction receipt"),
			SetupMocks: func(m *mocks, tc testCase) {
				tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), []byte{})
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByL2BlockHashAndIndex", context.Background(), tc.Hash, uint64(tc.Index), m.DbTx).
					Return(tx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
					Return(nil, errors.New("failed to get transaction receipt from state")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)

			result, err := c.TransactionInBlock(context.Background(), tc.Hash, tc.Index)

			if result != nil || testCase.ExpectedResult != nil {
				assert.Equal(t, testCase.ExpectedResult.Hash(), result.Hash())
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

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name        string
		BlockNumber string
		Index       uint

		ExpectedResult *types.Transaction
		ExpectedError  rpcError
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "Get Tx Successfully",
			BlockNumber:    "0x1",
			Index:          uint(0),
			ExpectedResult: types.NewTransaction(1, common.HexToAddress("0x111"), big.NewInt(2), 3, big.NewInt(4), []byte{5, 6, 7, 8}),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				tx := tc.ExpectedResult
				blockNumber, _ := encoding.DecodeUint64orHex(&tc.BlockNumber)
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByL2BlockNumberAndIndex", context.Background(), blockNumber, uint64(tc.Index), m.DbTx).
					Return(tx, nil).
					Once()

				receipt := types.NewReceipt([]byte{}, false, 0)
				receipt.BlockHash = common.Hash{}
				receipt.BlockNumber = big.NewInt(1)
				receipt.TransactionIndex = tc.Index
				m.State.
					On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
					Return(receipt, nil).
					Once()
			},
		},
		{
			Name:           "failed to identify block number",
			BlockNumber:    "latest",
			Index:          uint(0),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last block number from state"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number")).
					Once()
			},
		},
		{
			Name:           "Tx not found",
			BlockNumber:    "0x1",
			Index:          uint(0),
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				blockNumber, _ := encoding.DecodeUint64orHex(&tc.BlockNumber)
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByL2BlockNumberAndIndex", context.Background(), blockNumber, uint64(tc.Index), m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "Get Tx fail to get tx from state",
			BlockNumber:    "0x1",
			Index:          uint(0),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get transaction"),
			SetupMocks: func(m *mocks, tc testCase) {
				blockNumber, _ := encoding.DecodeUint64orHex(&tc.BlockNumber)
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByL2BlockNumberAndIndex", context.Background(), blockNumber, uint64(tc.Index), m.DbTx).
					Return(nil, errors.New("failed to get transaction by block and index from state")).
					Once()
			},
		},
		{
			Name:           "Tx found but receipt not found",
			BlockNumber:    "0x1",
			Index:          uint(0),
			ExpectedResult: nil,
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), []byte{})

				blockNumber, _ := encoding.DecodeUint64orHex(&tc.BlockNumber)
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByL2BlockNumberAndIndex", context.Background(), blockNumber, uint64(tc.Index), m.DbTx).
					Return(tx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "Get Tx fail to get tx receipt from state",
			BlockNumber:    "0x1",
			Index:          uint(0),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get transaction receipt"),
			SetupMocks: func(m *mocks, tc testCase) {
				tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), []byte{})

				blockNumber, _ := encoding.DecodeUint64orHex(&tc.BlockNumber)
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByL2BlockNumberAndIndex", context.Background(), blockNumber, uint64(tc.Index), m.DbTx).
					Return(tx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
					Return(nil, errors.New("failed to get transaction receipt from state")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)

			res, err := s.JSONRPCCall("eth_getTransactionByBlockNumberAndIndex", tc.BlockNumber, tc.Index)
			require.NoError(t, err)
			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result interface{}
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)

				if result != nil || testCase.ExpectedResult != nil {
					var tx types.Transaction
					err = json.Unmarshal(res.Result, &tx)
					require.NoError(t, err)
					assert.Equal(t, testCase.ExpectedResult.Hash(), tx.Hash())
				}
			}

			if res.Error != nil || testCase.ExpectedError != nil {
				assert.Equal(t, testCase.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, testCase.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestGetTransactionByHash(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name            string
		Hash            common.Hash
		ExpectedPending bool
		ExpectedResult  *types.Transaction
		ExpectedError   interface{}
		SetupMocks      func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:            "Get TX Successfully from state",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{}),
			ExpectedError:   nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(tc.ExpectedResult, nil).
					Once()

				receipt := types.NewReceipt([]byte{}, false, 0)
				receipt.BlockHash = common.Hash{}
				receipt.BlockNumber = big.NewInt(1)

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, m.DbTx).
					Return(receipt, nil).
					Once()
			},
		},
		{
			Name:            "Get TX Successfully from pool",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: true,
			ExpectedResult:  types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{}),
			ExpectedError:   nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()

				m.Pool.
					On("GetTxByHash", context.Background(), tc.Hash).
					Return(&pool.Transaction{Transaction: *tc.ExpectedResult}, nil).
					Once()
			},
		},
		{
			Name:            "TX Not Found",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  nil,
			ExpectedError:   ethereum.NotFound,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()

				m.Pool.
					On("GetTxByHash", context.Background(), tc.Hash).
					Return(nil, pgpoolstorage.ErrNotFound).
					Once()
			},
		},
		{
			Name:            "TX failed to load from the state",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  nil,
			ExpectedError:   newRPCError(defaultErrorCode, "failed to load transaction by hash from state"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(nil, errors.New("failed to load transaction by hash from state")).
					Once()
			},
		},
		{
			Name:            "TX failed to load from the pool",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  nil,
			ExpectedError:   newRPCError(defaultErrorCode, "failed to load transaction by hash from pool"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()

				m.Pool.
					On("GetTxByHash", context.Background(), tc.Hash).
					Return(nil, errors.New("failed to load transaction by hash from pool")).
					Once()
			},
		},
		{
			Name:            "TX receipt Not Found",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  nil,
			ExpectedError:   newRPCError(defaultErrorCode, "transaction receipt not found"),
			SetupMocks: func(m *mocks, tc testCase) {
				tx := &types.Transaction{}
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(tx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:            "TX receipt failed to load",
			Hash:            common.HexToHash("0x123"),
			ExpectedPending: false,
			ExpectedResult:  nil,
			ExpectedError:   newRPCError(defaultErrorCode, "failed to load transaction receipt from state"),
			SetupMocks: func(m *mocks, tc testCase) {
				tx := &types.Transaction{}
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(tx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, m.DbTx).
					Return(nil, errors.New("failed to load transaction receipt from state")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)

			result, pending, err := c.TransactionByHash(context.Background(), testCase.Hash)
			assert.Equal(t, testCase.ExpectedPending, pending)

			if result != nil || testCase.ExpectedResult != nil {
				assert.Equal(t, testCase.ExpectedResult.Hash(), result.Hash())
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

func TestGetBlockTransactionCountByHash(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		BlockHash      common.Hash
		ExpectedResult uint
		ExpectedError  interface{}
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "Count txs successfully",
			BlockHash:      common.HexToHash("0x123"),
			ExpectedResult: uint(10),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockTransactionCountByHash", context.Background(), tc.BlockHash, m.DbTx).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "Failed to count txs by hash",
			BlockHash:      common.HexToHash("0x123"),
			ExpectedResult: 0,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to count transactions"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockTransactionCountByHash", context.Background(), tc.BlockHash, m.DbTx).
					Return(uint64(0), errors.New("failed to count txs")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)
			result, err := c.TransactionCount(context.Background(), tc.BlockHash)

			assert.Equal(t, testCase.ExpectedResult, result)

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

func TestGetBlockTransactionCountByNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		BlockNumber    string
		ExpectedResult uint
		ExpectedError  rpcError
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "Count txs successfully for latest block",
			BlockNumber:    "latest",
			ExpectedResult: uint(10),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				blockNumber := uint64(10)
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(blockNumber, nil).
					Once()

				m.State.
					On("GetL2BlockTransactionCountByNumber", context.Background(), blockNumber, m.DbTx).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "Count txs successfully for pending block",
			BlockNumber:    "pending",
			ExpectedResult: uint(10),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.Pool.
					On("CountPendingTransactions", context.Background()).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get last block number",
			BlockNumber:    "latest",
			ExpectedResult: 0,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last block number from state"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number")).
					Once()
			},
		},
		{
			Name:           "failed to count tx",
			BlockNumber:    "latest",
			ExpectedResult: 0,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to count transactions"),
			SetupMocks: func(m *mocks, tc testCase) {
				blockNumber := uint64(10)
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(blockNumber, nil).
					Once()

				m.State.
					On("GetL2BlockTransactionCountByNumber", context.Background(), blockNumber, m.DbTx).
					Return(uint64(0), errors.New("failed to count")).
					Once()
			},
		},
		{
			Name:           "failed to count pending tx",
			BlockNumber:    "pending",
			ExpectedResult: 0,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to count pending transactions"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.Pool.
					On("CountPendingTransactions", context.Background()).
					Return(uint64(0), errors.New("failed to count")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)
			res, err := s.JSONRPCCall("eth_getBlockTransactionCountByNumber", tc.BlockNumber)

			require.NoError(t, err)
			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result argUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, uint(result))
			}

			if res.Error != nil || testCase.ExpectedError != nil {
				assert.Equal(t, testCase.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, testCase.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestGetTransactionCount(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		Address        string
		BlockNumber    string
		ExpectedResult uint
		ExpectedError  rpcError
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "Count txs successfully",
			Address:        common.HexToAddress("0x123").Hex(),
			BlockNumber:    "latest",
			ExpectedResult: uint(10),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				blockNumber := uint64(10)
				address := common.HexToAddress(tc.Address)
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(blockNumber, nil).
					Once()

				m.State.
					On("GetNonce", context.Background(), address, blockNumber, m.DbTx).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "Count txs nonce not found",
			Address:        common.HexToAddress("0x123").Hex(),
			BlockNumber:    "latest",
			ExpectedResult: 0,
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				blockNumber := uint64(10)
				address := common.HexToAddress(tc.Address)
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(blockNumber, nil).
					Once()

				m.State.
					On("GetNonce", context.Background(), address, blockNumber, m.DbTx).
					Return(uint64(0), state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "failed to get last block number",
			Address:        common.HexToAddress("0x123").Hex(),
			BlockNumber:    "latest",
			ExpectedResult: 0,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last block number from state"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number")).
					Once()
			},
		},
		{
			Name:           "failed to get nonce",
			Address:        common.HexToAddress("0x123").Hex(),
			BlockNumber:    "latest",
			ExpectedResult: 0,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to count transactions"),
			SetupMocks: func(m *mocks, tc testCase) {
				blockNumber := uint64(10)
				address := common.HexToAddress(tc.Address)
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(blockNumber, nil).
					Once()

				m.State.
					On("GetNonce", context.Background(), address, blockNumber, m.DbTx).
					Return(uint64(0), errors.New("failed to get nonce")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)
			res, err := s.JSONRPCCall("eth_getTransactionCount", tc.Address, tc.BlockNumber)

			require.NoError(t, err)
			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result argUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, uint(result))
			}

			if res.Error != nil || testCase.ExpectedError != nil {
				assert.Equal(t, testCase.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, testCase.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestGetTransactionReceipt(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		Hash           common.Hash
		ExpectedResult *types.Receipt
		ExpectedError  interface{}
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "Get TX receipt Successfully",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: types.NewReceipt([]byte{}, false, 0),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				tx := types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})
				privateKey, err := crypto.HexToECDSA(strings.TrimPrefix("0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e", "0x"))
				require.NoError(t, err)
				auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
				require.NoError(t, err)

				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(t, err)

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(signedTx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, m.DbTx).
					Return(tc.ExpectedResult, nil).
					Once()
			},
		},
		{
			Name:           "Get TX receipt but tx not found",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  ethereum.NotFound,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "Get TX receipt but failed to get tx",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get tx from state"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(nil, errors.New("failed to get tx")).
					Once()
			},
		},
		{
			Name:           "TX receipt Not Found",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  ethereum.NotFound,
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				tx := types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})
				privateKey, err := crypto.HexToECDSA(strings.TrimPrefix("0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e", "0x"))
				require.NoError(t, err)
				auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
				require.NoError(t, err)

				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(t, err)

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(signedTx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, m.DbTx).
					Return(nil, state.ErrNotFound).
					Once()
			},
		},
		{
			Name:           "TX receipt failed to load",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get tx receipt from state"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				tx := types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})
				privateKey, err := crypto.HexToECDSA(strings.TrimPrefix("0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e", "0x"))
				require.NoError(t, err)
				auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
				require.NoError(t, err)

				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(t, err)

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(signedTx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, m.DbTx).
					Return(nil, errors.New("failed to get tx receipt from state")).
					Once()
			},
		},
		{
			Name:           "Get TX but failed to build response Successfully",
			Hash:           common.HexToHash("0x123"),
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to build the receipt response"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				tx := types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})

				m.State.
					On("GetTransactionByHash", context.Background(), tc.Hash, m.DbTx).
					Return(tx, nil).
					Once()

				m.State.
					On("GetTransactionReceipt", context.Background(), tc.Hash, m.DbTx).
					Return(types.NewReceipt([]byte{}, false, 0), nil).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)

			result, err := c.TransactionReceipt(context.Background(), testCase.Hash)

			if result != nil || testCase.ExpectedResult != nil {
				assert.Equal(t, testCase.ExpectedResult.TxHash, result.TxHash)
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

func TestSendRawTransactionViaGeth(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name          string
		Tx            *types.Transaction
		ExpectedError interface{}
		SetupMocks    func(t *testing.T, m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:          "Send TX successfully",
			Tx:            types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), uint64(1), big.NewInt(1), []byte{}),
			ExpectedError: nil,
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				txMatchByHash := mock.MatchedBy(func(tx types.Transaction) bool {
					h1 := tx.Hash().Hex()
					h2 := tc.Tx.Hash().Hex()
					return h1 == h2
				})

				m.Pool.
					On("AddTx", context.Background(), txMatchByHash).
					Return(nil).
					Once()
			},
		},
		{
			Name:          "Send TX failed to add to the pool",
			Tx:            types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), uint64(1), big.NewInt(1), []byte{}),
			ExpectedError: newRPCError(defaultErrorCode, "failed to add TX to the pool"),
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				txMatchByHash := mock.MatchedBy(func(tx types.Transaction) bool {
					h1 := tx.Hash().Hex()
					h2 := tc.Tx.Hash().Hex()
					return h1 == h2
				})

				m.Pool.
					On("AddTx", context.Background(), txMatchByHash).
					Return(errors.New("failed to add TX to the pool")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(t, m, tc)

			err := c.SendTransaction(context.Background(), tc.Tx)

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

func TestSendRawTransactionJSONRPCCall(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		Input          string
		ExpectedResult *common.Hash
		ExpectedError  rpcError
		Prepare        func(t *testing.T, tc *testCase)
		SetupMocks     func(t *testing.T, m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name: "Send TX successfully",
			Prepare: func(t *testing.T, tc *testCase) {
				tx := types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), uint64(1), big.NewInt(1), []byte{})

				txBinary, err := tx.MarshalBinary()
				require.NoError(t, err)

				rawTx := hex.EncodeToHex(txBinary)
				require.NoError(t, err)

				tc.Input = rawTx
				tc.ExpectedResult = hashPtr(tx.Hash())
				tc.ExpectedError = nil
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				m.Pool.
					On("AddTx", context.Background(), mock.IsType(types.Transaction{})).
					Return(nil).
					Once()
			},
		},
		{
			Name: "Send TX failed to add to the pool",
			Prepare: func(t *testing.T, tc *testCase) {
				tx := types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), uint64(1), big.NewInt(1), []byte{})

				txBinary, err := tx.MarshalBinary()
				require.NoError(t, err)

				rawTx := hex.EncodeToHex(txBinary)
				require.NoError(t, err)

				tc.Input = rawTx
				tc.ExpectedResult = nil
				tc.ExpectedError = newRPCError(defaultErrorCode, "failed to add TX to the pool")
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				m.Pool.
					On("AddTx", context.Background(), mock.IsType(types.Transaction{})).
					Return(errors.New("failed to add TX to the pool")).
					Once()
			},
		},
		{
			Name: "Send invalid tx input",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.Input = "0x1234"
				tc.ExpectedResult = nil
				tc.ExpectedError = newRPCError(invalidParamsErrorCode, "invalid tx input")
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.Prepare(t, &tc)
			tc.SetupMocks(t, m, tc)

			res, err := s.JSONRPCCall("eth_sendRawTransaction", tc.Input)
			require.NoError(t, err)

			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil || tc.ExpectedResult != nil {
				var result common.Hash
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, *tc.ExpectedResult, result)
			}
			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestSendRawTransactionViaGethForNonSequencerNode(t *testing.T) {
	sequencerServer, sequencerMocks, _ := newSequencerMockedServer(t)
	defer sequencerServer.Stop()
	nonSequencerServer, _, nonSequencerClient := newNonSequencerMockedServer(t, sequencerServer.ServerURL)
	defer nonSequencerServer.Stop()

	type testCase struct {
		Name          string
		Tx            *types.Transaction
		ExpectedError interface{}
		SetupMocks    func(t *testing.T, m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:          "Send TX successfully",
			Tx:            types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), uint64(1), big.NewInt(1), []byte{}),
			ExpectedError: nil,
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				txMatchByHash := mock.MatchedBy(func(tx types.Transaction) bool {
					h1 := tx.Hash().Hex()
					h2 := tc.Tx.Hash().Hex()
					return h1 == h2
				})

				m.Pool.
					On("AddTx", context.Background(), txMatchByHash).
					Return(nil).
					Once()
			},
		},
		{
			Name:          "Send TX failed to add to the pool",
			Tx:            types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), uint64(1), big.NewInt(1), []byte{}),
			ExpectedError: newRPCError(defaultErrorCode, "failed to add TX to the pool"),
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				txMatchByHash := mock.MatchedBy(func(tx types.Transaction) bool {
					h1 := tx.Hash().Hex()
					h2 := tc.Tx.Hash().Hex()
					return h1 == h2
				})

				m.Pool.
					On("AddTx", context.Background(), txMatchByHash).
					Return(errors.New("failed to add TX to the pool")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(t, sequencerMocks, tc)

			err := nonSequencerClient.SendTransaction(context.Background(), tc.Tx)

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

func TestSendRawTransactionViaGethForNonSequencerNodeFailsToRelayTxToSequencerNode(t *testing.T) {
	nonSequencerServer, _, nonSequencerClient := newNonSequencerMockedServer(t, "http://wrong.url")
	defer nonSequencerServer.Stop()

	type testCase struct {
		Name          string
		Tx            *types.Transaction
		ExpectedError interface{}
	}

	testCases := []testCase{
		{
			Name:          "Send TX failed to relay tx to the sequencer node",
			Tx:            types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(1), uint64(1), big.NewInt(1), []byte{}),
			ExpectedError: newRPCError(defaultErrorCode, "failed to relay tx to the sequencer node"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase

			err := nonSequencerClient.SendTransaction(context.Background(), tc.Tx)

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

func TestProtocolVersion(t *testing.T) {
	s, _, _ := newSequencerMockedServer(t)
	defer s.Stop()

	res, err := s.JSONRPCCall("eth_protocolVersion")
	require.NoError(t, err)

	assert.Equal(t, float64(1), res.ID)
	assert.Equal(t, "2.0", res.JSONRPC)
	assert.Nil(t, res.Error)

	var result string
	err = json.Unmarshal(res.Result, &result)
	require.NoError(t, err)

	assert.Equal(t, "0x0", result)
}

func TestNewFilter(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		LogFilter      *LogFilter
		ExpectedResult argUint64
		ExpectedError  rpcError
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "New filter created successfully",
			LogFilter:      &LogFilter{},
			ExpectedResult: argUint64(1),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.Storage.
					On("NewLogFilter", *tc.LogFilter).
					Return(uint64(1), nil).
					Once()
			},
		},
		{
			Name:           "failed to create new filter",
			LogFilter:      &LogFilter{},
			ExpectedResult: argUint64(0),
			ExpectedError:  newRPCError(defaultErrorCode, "failed to create new log filter"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.Storage.
					On("NewLogFilter", *tc.LogFilter).
					Return(uint64(0), errors.New("failed to add new filter")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)

			res, err := s.JSONRPCCall("eth_newFilter", tc.LogFilter)
			require.NoError(t, err)

			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result argUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, result)
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestNewBlockFilter(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult argUint64
		ExpectedError  rpcError
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "New block filter created successfully",
			ExpectedResult: argUint64(1),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.Storage.
					On("NewBlockFilter").
					Return(uint64(1), nil).
					Once()
			},
		},
		{
			Name:           "failed to create new block filter",
			ExpectedResult: argUint64(0),
			ExpectedError:  newRPCError(defaultErrorCode, "failed to create new block filter"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.Storage.
					On("NewBlockFilter").
					Return(uint64(0), errors.New("failed to add new block filter")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)

			res, err := s.JSONRPCCall("eth_newBlockFilter")
			require.NoError(t, err)

			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result argUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, result)
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestNewPendingTransactionFilter(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult argUint64
		ExpectedError  rpcError
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "New pending transaction filter created successfully",
			ExpectedResult: argUint64(1),
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.Storage.
					On("NewPendingTransactionFilter").
					Return(uint64(1), nil).
					Once()
			},
		},
		{
			Name:           "failed to create new pending transaction filter",
			ExpectedResult: argUint64(0),
			ExpectedError:  newRPCError(defaultErrorCode, "failed to create new pending transaction filter"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.Storage.
					On("NewPendingTransactionFilter").
					Return(uint64(0), errors.New("failed to add new pending transaction filter")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)

			res, err := s.JSONRPCCall("eth_newPendingTransactionFilter")
			require.NoError(t, err)

			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result argUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, result)
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestUninstallFilter(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		FilterID       argUint64
		ExpectedResult bool
		ExpectedError  rpcError
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name:           "Uninstalls filter successfully",
			FilterID:       argUint64(1),
			ExpectedResult: true,
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.Storage.
					On("UninstallFilter", uint64(tc.FilterID)).
					Return(true, nil).
					Once()
			},
		},
		{
			Name:           "filter already uninstalled",
			FilterID:       argUint64(1),
			ExpectedResult: false,
			ExpectedError:  nil,
			SetupMocks: func(m *mocks, tc testCase) {
				m.Storage.
					On("UninstallFilter", uint64(tc.FilterID)).
					Return(false, nil).
					Once()
			},
		},
		{
			Name:           "failed to uninstall filter",
			FilterID:       argUint64(1),
			ExpectedResult: false,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to uninstall filter"),
			SetupMocks: func(m *mocks, tc testCase) {
				m.Storage.
					On("UninstallFilter", uint64(tc.FilterID)).
					Return(false, errors.New("failed to uninstall filter")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)

			res, err := s.JSONRPCCall("eth_uninstallFilter", tc.FilterID)
			require.NoError(t, err)

			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result bool
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, result)
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestGetLogs(t *testing.T) {
	s, m, c := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		Filter         ethereum.FilterQuery
		ExpectedResult []types.Log
		ExpectedError  interface{}
		Prepare        func(t *testing.T, tc *testCase)
		SetupMocks     func(m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name: "Get logs successfully",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.Filter = ethereum.FilterQuery{
					FromBlock: big.NewInt(1), ToBlock: big.NewInt(2),
					Addresses: []common.Address{common.HexToAddress("0x111")},
					Topics:    [][]common.Hash{{common.HexToHash("0x222")}},
				}
				tc.ExpectedResult = []types.Log{{
					Address: common.Address{}, Topics: []common.Hash{}, Data: []byte{},
					BlockNumber: uint64(1), TxHash: common.Hash{}, TxIndex: uint(1),
					BlockHash: common.Hash{}, Index: uint(1), Removed: false,
				}}
				tc.ExpectedError = nil
			},
			SetupMocks: func(m *mocks, tc testCase) {
				var since *time.Time
				logs := make([]*types.Log, 0, len(tc.ExpectedResult))
				for _, log := range tc.ExpectedResult {
					l := log
					logs = append(logs, &l)
				}

				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLogs", context.Background(), tc.Filter.FromBlock.Uint64(), tc.Filter.ToBlock.Uint64(), tc.Filter.Addresses, tc.Filter.Topics, tc.Filter.BlockHash, since, m.DbTx).
					Return(logs, nil).
					Once()
			},
		},
		{
			Name: "Get logs fails to get logs from state",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.Filter = ethereum.FilterQuery{
					FromBlock: big.NewInt(1), ToBlock: big.NewInt(2),
					Addresses: []common.Address{common.HexToAddress("0x111")},
					Topics:    [][]common.Hash{{common.HexToHash("0x222")}},
				}
				tc.ExpectedResult = nil
				tc.ExpectedError = newRPCError(defaultErrorCode, "failed to get logs from state")
			},
			SetupMocks: func(m *mocks, tc testCase) {
				var since *time.Time
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLogs", context.Background(), tc.Filter.FromBlock.Uint64(), tc.Filter.ToBlock.Uint64(), tc.Filter.Addresses, tc.Filter.Topics, tc.Filter.BlockHash, since, m.DbTx).
					Return(nil, errors.New("failed to get logs from state")).
					Once()
			},
		},
		{
			Name: "Get logs fails to identify from block",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.Filter = ethereum.FilterQuery{
					FromBlock: big.NewInt(-1), ToBlock: big.NewInt(2),
					Addresses: []common.Address{common.HexToAddress("0x111")},
					Topics:    [][]common.Hash{{common.HexToHash("0x222")}},
				}
				tc.ExpectedResult = nil
				tc.ExpectedError = newRPCError(defaultErrorCode, "failed to get the last block number from state")
			},
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number from state")).
					Once()
			},
		},
		{
			Name: "Get logs fails to identify to block",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.Filter = ethereum.FilterQuery{
					FromBlock: big.NewInt(1), ToBlock: big.NewInt(-1),
					Addresses: []common.Address{common.HexToAddress("0x111")},
					Topics:    [][]common.Hash{{common.HexToHash("0x222")}},
				}
				tc.ExpectedResult = nil
				tc.ExpectedError = newRPCError(defaultErrorCode, "failed to get the last block number from state")
			},
			SetupMocks: func(m *mocks, tc testCase) {
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLastL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last block number from state")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.Prepare(t, &tc)
			tc.SetupMocks(m, tc)

			result, err := c.FilterLogs(context.Background(), tc.Filter)

			if result != nil || tc.ExpectedResult != nil {
				assert.ElementsMatch(t, tc.ExpectedResult, result)
			}

			if err != nil || tc.ExpectedError != nil {
				if expectedErr, ok := tc.ExpectedError.(*RPCError); ok {
					rpcErr := err.(rpcError)
					assert.Equal(t, expectedErr.ErrorCode(), rpcErr.ErrorCode())
					assert.Equal(t, expectedErr.Error(), rpcErr.Error())
				} else {
					assert.Equal(t, tc.ExpectedError, err)
				}
			}
		})
	}
}

func TestGetFilterLogs(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		FilterID       argUint64
		ExpectedResult []types.Log
		ExpectedError  rpcError
		Prepare        func(t *testing.T, tc *testCase)
		SetupMocks     func(t *testing.T, m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name: "Get filter logs successfully",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				tc.ExpectedResult = []types.Log{{
					Address: common.Address{}, Topics: []common.Hash{}, Data: []byte{},
					BlockNumber: uint64(1), TxHash: common.Hash{}, TxIndex: uint(1),
					BlockHash: common.Hash{}, Index: uint(1), Removed: false,
				}}
				tc.ExpectedError = nil
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				var since *time.Time
				logs := make([]*types.Log, 0, len(tc.ExpectedResult))
				for _, log := range tc.ExpectedResult {
					l := log
					logs = append(logs, &l)
				}

				logFilter := LogFilter{
					FromBlock: BlockNumber(1), ToBlock: BlockNumber(2),
					Addresses: []common.Address{common.HexToAddress("0x111")},
					Topics:    [][]common.Hash{{common.HexToHash("0x222")}},
				}

				logFilterJSON, err := json.Marshal(&logFilter)
				require.NoError(t, err)

				parameters := string(logFilterJSON)

				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypeLog,
					LastPoll:   time.Now(),
					Parameters: parameters,
				}
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()

				m.State.
					On("GetLogs", context.Background(), uint64(logFilter.FromBlock), uint64(logFilter.ToBlock), logFilter.Addresses, logFilter.Topics, logFilter.BlockHash, since, m.DbTx).
					Return(logs, nil).
					Once()
			},
		},
		{
			Name: "Get filter logs filter not found",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				tc.ExpectedResult = nil
				tc.ExpectedError = nil
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(nil, ErrNotFound).
					Once()
			},
		},
		{
			Name: "Get filter logs failed to get filter",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				tc.ExpectedResult = nil
				tc.ExpectedError = newRPCError(defaultErrorCode, "failed to get filter from storage")
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(nil, errors.New("failed to get filter")).
					Once()
			},
		},
		{
			Name: "Get filter logs is a valid filter but its not a log filter",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				tc.ExpectedResult = nil
				tc.ExpectedError = nil
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypeBlock,
					LastPoll:   time.Now(),
					Parameters: "",
				}

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()
			},
		},
		{
			Name: "Get filter logs failed to parse filter parameters",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				tc.ExpectedResult = nil
				tc.ExpectedError = newRPCError(defaultErrorCode, "failed to read filter parameters")
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypeLog,
					LastPoll:   time.Now(),
					Parameters: "invalid parameters",
				}

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.Prepare(t, &tc)
			tc.SetupMocks(t, m, tc)

			res, err := s.JSONRPCCall("eth_getFilterLogs", hex.EncodeUint64(uint64(tc.FilterID)))
			require.NoError(t, err)
			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result interface{}
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)

				if result != nil || tc.ExpectedResult != nil {
					var logs []types.Log
					err = json.Unmarshal(res.Result, &logs)
					require.NoError(t, err)
					assert.ElementsMatch(t, tc.ExpectedResult, logs)
				}
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestGetFilterChanges(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name            string
		FilterID        argUint64
		ExpectedResults []interface{}
		ExpectedErrors  []rpcError
		Prepare         func(t *testing.T, tc *testCase)
		SetupMocks      func(t *testing.T, m *mocks, tc testCase)
	}

	testCases := []testCase{
		{
			Name: "Get block filter changes multiple times successfully",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(2)
				// first call
				tc.ExpectedResults = append(tc.ExpectedResults, []common.Hash{
					common.HexToHash("0x111"),
				})
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)

				// second call
				tc.ExpectedResults = append(tc.ExpectedResults, []common.Hash{
					common.HexToHash("0x222"),
					common.HexToHash("0x333"),
				})
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)

				// third call
				tc.ExpectedResults = append(tc.ExpectedResults, []common.Hash{})
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypeBlock,
					LastPoll:   time.Now(),
					Parameters: "{}",
				}

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()

				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetL2BlockHashesSince", context.Background(), filter.LastPoll, m.DbTx).
					Return(tc.ExpectedResults[0].([]common.Hash), nil).
					Once()

				m.Storage.
					On("UpdateFilterLastPoll", uint64(tc.FilterID)).
					Run(func(args mock.Arguments) {
						filter.LastPoll = time.Now()

						m.Storage.
							On("GetFilter", uint64(tc.FilterID)).
							Return(filter, nil).
							Once()

						m.DbTx.
							On("Commit", context.Background()).
							Return(nil).
							Once()

						m.State.
							On("BeginStateTransaction", context.Background()).
							Return(m.DbTx, nil).
							Once()

						m.State.
							On("GetL2BlockHashesSince", context.Background(), filter.LastPoll, m.DbTx).
							Return(tc.ExpectedResults[1].([]common.Hash), nil).
							Once()

						m.Storage.
							On("UpdateFilterLastPoll", uint64(tc.FilterID)).
							Run(func(args mock.Arguments) {
								filter.LastPoll = time.Now()

								m.Storage.
									On("GetFilter", uint64(tc.FilterID)).
									Return(filter, nil).
									Once()

								m.DbTx.
									On("Commit", context.Background()).
									Return(nil).
									Once()

								m.State.
									On("BeginStateTransaction", context.Background()).
									Return(m.DbTx, nil).
									Once()

								m.State.
									On("GetL2BlockHashesSince", context.Background(), filter.LastPoll, m.DbTx).
									Return(tc.ExpectedResults[2].([]common.Hash), nil).
									Once()

								m.Storage.
									On("UpdateFilterLastPoll", uint64(tc.FilterID)).
									Return(nil).
									Once()
							}).
							Return(nil).
							Once()
					}).
					Return(nil).
					Once()
			},
		},
		{
			Name: "Get pending transactions filter changes multiple times successfully",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(3)
				// first call
				tc.ExpectedResults = append(tc.ExpectedResults, []common.Hash{
					common.HexToHash("0x444"),
				})
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)

				// second call
				tc.ExpectedResults = append(tc.ExpectedResults, []common.Hash{
					common.HexToHash("0x555"),
					common.HexToHash("0x666"),
				})
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)

				// third call
				tc.ExpectedResults = append(tc.ExpectedResults, []common.Hash{})
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypePendingTx,
					LastPoll:   time.Now(),
					Parameters: "{}",
				}

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()

				m.Pool.
					On("GetPendingTxHashesSince", context.Background(), filter.LastPoll).
					Return(tc.ExpectedResults[0].([]common.Hash), nil).
					Once()

				m.Storage.
					On("UpdateFilterLastPoll", uint64(tc.FilterID)).
					Run(func(args mock.Arguments) {
						filter.LastPoll = time.Now()

						m.Storage.
							On("GetFilter", uint64(tc.FilterID)).
							Return(filter, nil).
							Once()

						m.Pool.
							On("GetPendingTxHashesSince", context.Background(), filter.LastPoll).
							Return(tc.ExpectedResults[1].([]common.Hash), nil).
							Once()

						m.Storage.
							On("UpdateFilterLastPoll", uint64(tc.FilterID)).
							Run(func(args mock.Arguments) {
								filter.LastPoll = time.Now()

								m.Storage.
									On("GetFilter", uint64(tc.FilterID)).
									Return(filter, nil).
									Once()

								m.Pool.
									On("GetPendingTxHashesSince", context.Background(), filter.LastPoll).
									Return(tc.ExpectedResults[2].([]common.Hash), nil).
									Once()

								m.Storage.
									On("UpdateFilterLastPoll", uint64(tc.FilterID)).
									Return(nil).
									Once()
							}).
							Return(nil).
							Once()
					}).
					Return(nil).
					Once()
			},
		},
		{
			Name: "Get log filter changes multiple times successfully",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				// first call
				tc.ExpectedResults = append(tc.ExpectedResults, []types.Log{{
					Address: common.Address{}, Topics: []common.Hash{}, Data: []byte{},
					BlockNumber: uint64(1), TxHash: common.Hash{}, TxIndex: uint(1),
					BlockHash: common.Hash{}, Index: uint(1), Removed: false,
				}})
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)

				// second call
				tc.ExpectedResults = append(tc.ExpectedResults, []types.Log{{
					Address: common.Address{}, Topics: []common.Hash{}, Data: []byte{},
					BlockNumber: uint64(1), TxHash: common.Hash{}, TxIndex: uint(1),
					BlockHash: common.Hash{}, Index: uint(1), Removed: false,
				}, {
					Address: common.Address{}, Topics: []common.Hash{}, Data: []byte{},
					BlockNumber: uint64(1), TxHash: common.Hash{}, TxIndex: uint(1),
					BlockHash: common.Hash{}, Index: uint(1), Removed: false,
				}})
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)

				// third call
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				logFilter := LogFilter{
					FromBlock: BlockNumber(1), ToBlock: BlockNumber(2),
					Addresses: []common.Address{common.HexToAddress("0x111")},
					Topics:    [][]common.Hash{{common.HexToHash("0x222")}},
				}

				logFilterJSON, err := json.Marshal(&logFilter)
				require.NoError(t, err)

				parameters := string(logFilterJSON)

				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypeLog,
					LastPoll:   time.Now(),
					Parameters: parameters,
				}
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()

				expectedLogs := tc.ExpectedResults[0].([]types.Log)
				logs := make([]*types.Log, 0, len(expectedLogs))
				for _, log := range expectedLogs {
					l := log
					logs = append(logs, &l)
				}

				m.State.
					On("GetLogs", context.Background(), uint64(logFilter.FromBlock), uint64(logFilter.ToBlock), logFilter.Addresses, logFilter.Topics, logFilter.BlockHash, &filter.LastPoll, m.DbTx).
					Return(logs, nil).
					Once()

				m.Storage.
					On("UpdateFilterLastPoll", uint64(tc.FilterID)).
					Run(func(args mock.Arguments) {
						filter.LastPoll = time.Now()

						m.DbTx.
							On("Commit", context.Background()).
							Return(nil).
							Once()

						m.State.
							On("BeginStateTransaction", context.Background()).
							Return(m.DbTx, nil).
							Once()

						m.Storage.
							On("GetFilter", uint64(tc.FilterID)).
							Return(filter, nil).
							Once()

						expectedLogs = tc.ExpectedResults[1].([]types.Log)
						logs = make([]*types.Log, 0, len(expectedLogs))
						for _, log := range expectedLogs {
							l := log
							logs = append(logs, &l)
						}

						m.State.
							On("GetLogs", context.Background(), uint64(logFilter.FromBlock), uint64(logFilter.ToBlock), logFilter.Addresses, logFilter.Topics, logFilter.BlockHash, &filter.LastPoll, m.DbTx).
							Return(logs, nil).
							Once()

						m.Storage.
							On("UpdateFilterLastPoll", uint64(tc.FilterID)).
							Run(func(args mock.Arguments) {
								filter.LastPoll = time.Now()
								m.DbTx.
									On("Commit", context.Background()).
									Return(nil).
									Once()

								m.State.
									On("BeginStateTransaction", context.Background()).
									Return(m.DbTx, nil).
									Once()

								m.Storage.
									On("GetFilter", uint64(tc.FilterID)).
									Return(filter, nil).
									Once()

								m.State.
									On("GetLogs", context.Background(), uint64(logFilter.FromBlock), uint64(logFilter.ToBlock), logFilter.Addresses, logFilter.Topics, logFilter.BlockHash, &filter.LastPoll, m.DbTx).
									Return([]*types.Log{}, nil).
									Once()

								m.Storage.
									On("UpdateFilterLastPoll", uint64(tc.FilterID)).
									Return(nil).
									Once()
							}).
							Return(nil).
							Once()
					}).
					Return(nil).
					Once()
			},
		},
		{
			Name: "Get filter changes when filter is not found",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				// first call
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(nil, ErrNotFound).
					Once()
			},
		},
		{
			Name: "Get filter changes fails to get filter",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				// first call
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, newRPCError(defaultErrorCode, "failed to get filter from storage"))
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(nil, errors.New("failed to get filter")).
					Once()
			},
		},
		{
			Name: "Get log filter changes fails to parse parameters",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				// first call
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, newRPCError(defaultErrorCode, "failed to read filter parameters"))
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypeLog,
					LastPoll:   time.Now(),
					Parameters: "invalid parameters",
				}

				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()
			},
		},
		{
			Name: "Get block filter changes fails to get block hashes",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(2)
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, newRPCError(defaultErrorCode, "failed to get block hashes"))
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypeBlock,
					LastPoll:   time.Now(),
					Parameters: "{}",
				}
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()

				m.State.
					On("GetL2BlockHashesSince", context.Background(), filter.LastPoll, m.DbTx).
					Return([]common.Hash{}, errors.New("failed to get hashes")).
					Once()
			},
		},
		{
			Name: "Get block filter changes fails to update the last time it was requested",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(2)
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, newRPCError(defaultErrorCode, "failed to update last time the filter changes were requested"))
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypeBlock,
					LastPoll:   time.Now(),
					Parameters: "{}",
				}

				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()

				m.State.
					On("GetL2BlockHashesSince", context.Background(), filter.LastPoll, m.DbTx).
					Return([]common.Hash{}, nil).
					Once()

				m.Storage.
					On("UpdateFilterLastPoll", uint64(tc.FilterID)).
					Return(errors.New("failed to update filter last poll")).
					Once()
			},
		},
		{
			Name: "Get pending transactions filter fails to get the hashes",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(3)
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, newRPCError(defaultErrorCode, "failed to get pending transaction hashes"))
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypePendingTx,
					LastPoll:   time.Now(),
					Parameters: "{}",
				}

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()

				m.Pool.
					On("GetPendingTxHashesSince", context.Background(), filter.LastPoll).
					Return([]common.Hash{}, errors.New("failed to get pending tx hashes")).
					Once()
			},
		},
		{
			Name: "Get pending transactions fails to update the last time it was requested",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(3)
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, newRPCError(defaultErrorCode, "failed to update last time the filter changes were requested"))
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypePendingTx,
					LastPoll:   time.Now(),
					Parameters: "{}",
				}

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()

				m.Pool.
					On("GetPendingTxHashesSince", context.Background(), filter.LastPoll).
					Return([]common.Hash{}, nil).
					Once()

				m.Storage.
					On("UpdateFilterLastPoll", uint64(tc.FilterID)).
					Return(errors.New("failed to update filter last poll")).
					Once()
			},
		},
		{
			Name: "Get log filter changes fails to get logs",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, newRPCError(defaultErrorCode, "failed to get logs from state"))
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				logFilter := LogFilter{
					FromBlock: BlockNumber(1), ToBlock: BlockNumber(2),
					Addresses: []common.Address{common.HexToAddress("0x111")},
					Topics:    [][]common.Hash{{common.HexToHash("0x222")}},
				}

				logFilterJSON, err := json.Marshal(&logFilter)
				require.NoError(t, err)

				parameters := string(logFilterJSON)

				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypeLog,
					LastPoll:   time.Now(),
					Parameters: parameters,
				}
				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()

				m.State.
					On("GetLogs", context.Background(), uint64(logFilter.FromBlock), uint64(logFilter.ToBlock), logFilter.Addresses, logFilter.Topics, logFilter.BlockHash, &filter.LastPoll, m.DbTx).
					Return(nil, errors.New("failed to get logs")).
					Once()
			},
		},
		{
			Name: "Get log filter changes fails to update the last time it was requested",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(1)
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, newRPCError(defaultErrorCode, "failed to update last time the filter changes were requested"))
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				logFilter := LogFilter{
					FromBlock: BlockNumber(1), ToBlock: BlockNumber(2),
					Addresses: []common.Address{common.HexToAddress("0x111")},
					Topics:    [][]common.Hash{{common.HexToHash("0x222")}},
				}

				logFilterJSON, err := json.Marshal(&logFilter)
				require.NoError(t, err)

				parameters := string(logFilterJSON)

				filter := &Filter{
					ID:         uint64(tc.FilterID),
					Type:       FilterTypeLog,
					LastPoll:   time.Now(),
					Parameters: parameters,
				}

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()

				m.DbTx.
					On("Rollback", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				m.State.
					On("GetLogs", context.Background(), uint64(logFilter.FromBlock), uint64(logFilter.ToBlock), logFilter.Addresses, logFilter.Topics, logFilter.BlockHash, &filter.LastPoll, m.DbTx).
					Return([]*types.Log{}, nil).
					Once()

				m.Storage.
					On("UpdateFilterLastPoll", uint64(tc.FilterID)).
					Return(errors.New("failed to update filter last poll")).
					Once()
			},
		},
		{
			Name: "Get filter changes for a unknown log type",
			Prepare: func(t *testing.T, tc *testCase) {
				tc.FilterID = argUint64(4)
				tc.ExpectedResults = append(tc.ExpectedResults, nil)
				tc.ExpectedErrors = append(tc.ExpectedErrors, nil)
			},
			SetupMocks: func(t *testing.T, m *mocks, tc testCase) {
				filter := &Filter{
					Type: "unknown type",
				}

				m.Storage.
					On("GetFilter", uint64(tc.FilterID)).
					Return(filter, nil).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.Prepare(t, &tc)
			tc.SetupMocks(t, m, tc)

			timesToCall := len(tc.ExpectedResults)

			for i := 0; i < timesToCall; i++ {
				res, err := s.JSONRPCCall("eth_getFilterChanges", tc.FilterID)
				require.NoError(t, err)
				assert.Equal(t, float64(1), res.ID)
				assert.Equal(t, "2.0", res.JSONRPC)

				if res.Result != nil {
					var result interface{}
					err = json.Unmarshal(res.Result, &result)
					require.NoError(t, err)

					if result != nil || tc.ExpectedResults[i] != nil {
						if logs, ok := tc.ExpectedResults[i].([]types.Log); ok {
							err = json.Unmarshal(res.Result, &logs)
							require.NoError(t, err)
							assert.ElementsMatch(t, tc.ExpectedResults[i], logs)
						}
						if hashes, ok := tc.ExpectedResults[i].([]common.Hash); ok {
							err = json.Unmarshal(res.Result, &hashes)
							require.NoError(t, err)
							assert.ElementsMatch(t, tc.ExpectedResults[i], hashes)
						}
					}
				}

				if res.Error != nil || tc.ExpectedErrors[i] != nil {
					assert.Equal(t, tc.ExpectedErrors[i].ErrorCode(), res.Error.Code)
					assert.Equal(t, tc.ExpectedErrors[i].Error(), res.Error.Message)
				}
			}
		})
	}
}

func addressPtr(i common.Address) *common.Address {
	return &i
}

func hashPtr(h common.Hash) *common.Hash {
	return &h
}
