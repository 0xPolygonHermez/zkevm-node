package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsolidatedBlockNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult *uint64
		ExpectedError  rpcError
		SetupMocks     func(m *mocks)
	}

	testCases := []testCase{
		{
			Name:           "Get consolidated block number successfully",
			ExpectedResult: ptrUint64(10),
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
					On("GetLastConsolidatedL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get consolidated block number",
			ExpectedResult: nil,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get last consolidated block number from state"),
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
					On("GetLastConsolidatedL2BlockNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last consolidated block number")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_consolidatedBlockNumber")
			require.NoError(t, err)

			if res.Result != nil {
				var result argUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, *tc.ExpectedResult, uint64(result))
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestIsBlockConsolidated(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult bool
		ExpectedError  rpcError
		SetupMocks     func(m *mocks)
	}

	testCases := []testCase{
		{
			Name:           "Query status of block number successfully",
			ExpectedResult: true,
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
					On("IsL2BlockConsolidated", context.Background(), uint64(1), m.DbTx).
					Return(true, nil).
					Once()
			},
		},
		{
			Name:           "Failed to query the consolidation status",
			ExpectedResult: false,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to check if the block is consolidated"),
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
					On("IsL2BlockConsolidated", context.Background(), uint64(1), m.DbTx).
					Return(false, errors.New("failed to check if the block is consolidated")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_isBlockConsolidated", "0x1")
			require.NoError(t, err)

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

func TestIsBlockVirtualized(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult bool
		ExpectedError  rpcError
		SetupMocks     func(m *mocks)
	}

	testCases := []testCase{
		{
			Name:           "Query status of block number successfully",
			ExpectedResult: true,
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
					On("IsL2BlockVirtualized", context.Background(), uint64(1), m.DbTx).
					Return(true, nil).
					Once()
			},
		},
		{
			Name:           "Failed to query the virtualization status",
			ExpectedResult: false,
			ExpectedError:  newRPCError(defaultErrorCode, "failed to check if the block is virtualized"),
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
					On("IsL2BlockVirtualized", context.Background(), uint64(1), m.DbTx).
					Return(false, errors.New("failed to check if the block is virtualized")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_isBlockVirtualized", "0x1")
			require.NoError(t, err)

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

func TestBatchNumberByBlockNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()
	blockNumber := uint64(1)
	batchNumber := uint64(1)

	type testCase struct {
		Name           string
		ExpectedResult uint64
		ExpectedError  rpcError
		SetupMocks     func(m *mocks)
	}

	testCases := []testCase{
		{
			Name:           "Query status of batch number of l2 block by its number successfully",
			ExpectedResult: batchNumber,
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
					On("BatchNumberByL2BlockNumber", context.Background(), blockNumber, m.DbTx).
					Return(batchNumber, nil).
					Once()
			},
		},
		{
			Name:           "Failed to query the consolidation status",
			ExpectedResult: uint64(0),
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get batch number from block number"),
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
					On("BatchNumberByL2BlockNumber", context.Background(), blockNumber, m.DbTx).
					Return(uint64(0), errors.New("failed to get batch number of l2 batchNum")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_batchNumberByBlockNumber", hex.EncodeUint64(blockNumber))
			require.NoError(t, err)

			if res.Result != nil {
				var result uint64
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

func TestBatchNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult uint64
		ExpectedError  rpcError
		SetupMocks     func(m *mocks)
	}

	testCases := []testCase{
		{
			Name:           "get batch number successfully",
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
					On("GetLastBatchNumber", context.Background(), m.DbTx).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get batch number",
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last batch number from state"),
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
					On("GetLastBatchNumber", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last batch number")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_batchNumber")
			require.NoError(t, err)

			if res.Result != nil {
				var result argUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, uint64(result))
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestVirtualBatchNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult uint64
		ExpectedError  rpcError
		SetupMocks     func(m *mocks)
	}

	testCases := []testCase{
		{
			Name:           "get virtual batch number successfully",
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
					On("GetLastVirtualBatchNum", context.Background(), m.DbTx).
					Return(uint64(10), nil).
					Once()
			},
		},
		{
			Name:           "failed to get virtual batch number",
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last virtual batch number from state"),
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
					On("GetLastVirtualBatchNum", context.Background(), m.DbTx).
					Return(uint64(0), errors.New("failed to get last batch number")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_virtualBatchNumber")
			require.NoError(t, err)

			if res.Result != nil {
				var result argUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, uint64(result))
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestVerifiedBatchNumber(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult uint64
		ExpectedError  rpcError
		SetupMocks     func(m *mocks)
	}

	testCases := []testCase{
		{
			Name:           "get verified batch number successfully",
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
					On("GetLastVerifiedBatch", context.Background(), m.DbTx).
					Return(&state.VerifiedBatch{BatchNumber: uint64(10)}, nil).
					Once()
			},
		},
		{
			Name:           "failed to get verified batch number",
			ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last verified batch number from state"),
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
					On("GetLastVerifiedBatch", context.Background(), m.DbTx).
					Return(nil, errors.New("failed to get last batch number")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_verifiedBatchNumber")
			require.NoError(t, err)

			if res.Result != nil {
				var result argUint64
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedResult, uint64(result))
			}

			if res.Error != nil || tc.ExpectedError != nil {
				assert.Equal(t, tc.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, tc.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func TestGetBatchByNumber(t *testing.T) {
	type testCase struct {
		Name           string
		Number         *big.Int
		ExpectedResult *rpcBatch
		ExpectedError  rpcError
		SetupMocks     func(*mockedServer, *mocks, *testCase)
	}

	testCases := []testCase{
		// {
		// 	Name:           "Batch not found",
		// 	Number:         big.NewInt(123),
		// 	ExpectedResult: nil,
		// 	ExpectedError:  ethereum.NotFound,
		// 	SetupMocks: func(m *mocks, tc *testCase) {
		// 		m.DbTx.
		// 			On("Commit", context.Background()).
		// 			Return(nil).
		// 			Once()

		// 		m.State.
		// 			On("BeginStateTransaction", context.Background()).
		// 			Return(m.DbTx, nil).
		// 			Once()

		// 		m.State.
		// 			On("GetBatchByNumber", context.Background(), tc.Number.Uint64(), m.DbTx).
		// 			Return(nil, state.ErrNotFound)
		// 	},
		// },
		{
			Name:   "get specific batch successfully",
			Number: big.NewInt(345),
			ExpectedResult: &rpcBatch{
				Number:              1,
				Coinbase:            common.HexToAddress("0x1"),
				StateRoot:           common.HexToHash("0x2"),
				AccInputHash:        common.HexToHash("0x3"),
				GlobalExitRoot:      common.HexToHash("0x4"),
				Timestamp:           1,
				SendSequencesTxHash: ptrHash(common.HexToHash("0x10")),
				VerifyBatchTxHash:   ptrHash(common.HexToHash("0x20")),
			},
			ExpectedError: nil,
			SetupMocks: func(s *mockedServer, m *mocks, tc *testCase) {
				m.DbTx.
					On("Commit", context.Background()).
					Return(nil).
					Once()

				m.State.
					On("BeginStateTransaction", context.Background()).
					Return(m.DbTx, nil).
					Once()

				batch := &state.Batch{
					BatchNumber:    1,
					Coinbase:       common.HexToAddress("0x1"),
					StateRoot:      common.HexToHash("0x2"),
					AccInputHash:   common.HexToHash("0x3"),
					GlobalExitRoot: common.HexToHash("0x4"),
					Timestamp:      time.Unix(1, 0),
				}

				m.State.
					On("GetBatchByNumber", context.Background(), tc.Number.Uint64(), m.DbTx).
					Return(batch, nil).
					Once()

				virtualBatch := &state.VirtualBatch{
					TxHash: common.HexToHash("0x10"),
				}

				m.State.
					On("GetVirtualBatch", context.Background(), tc.Number.Uint64(), m.DbTx).
					Return(virtualBatch, nil).
					Once()

				verifiedBatch := &state.VerifiedBatch{
					TxHash: common.HexToHash("0x20"),
				}

				m.State.
					On("GetVerifiedBatch", context.Background(), tc.Number.Uint64(), m.DbTx).
					Return(verifiedBatch, nil).
					Once()

				txs := []*types.Transaction{
					signTx(types.NewTransaction(1001, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.Config.ChainID),
					signTx(types.NewTransaction(1002, common.HexToAddress("0x1000"), big.NewInt(1000), 1001, big.NewInt(1002), []byte("1003")), s.Config.ChainID),
				}

				batchTxs := make([]types.Transaction, 0, len(txs))

				tc.ExpectedResult.Transactions = []rpcTransactionOrHash{}

				for i, tx := range txs {
					blockNumber := big.NewInt(int64(i))
					blockHash := common.HexToHash(hex.EncodeUint64(uint64(i)))
					receipt := types.NewReceipt([]byte{}, false, uint64(0))
					receipt.TxHash = tx.Hash()
					receipt.TransactionIndex = uint(i)
					receipt.BlockNumber = blockNumber
					receipt.BlockHash = blockHash
					m.State.
						On("GetTransactionReceipt", context.Background(), tx.Hash(), m.DbTx).
						Return(receipt, nil).
						Once()

					from, _ := state.GetSender(*tx)
					V, R, S := tx.RawSignatureValues()

					tc.ExpectedResult.Transactions = append(tc.ExpectedResult.Transactions,
						rpcTransaction{
							Nonce:       argUint64(tx.Nonce()),
							GasPrice:    argBig(*tx.GasPrice()),
							Gas:         argUint64(tx.Gas()),
							To:          tx.To(),
							Value:       argBig(*tx.Value()),
							Input:       tx.Data(),
							Hash:        tx.Hash(),
							From:        from,
							BlockNumber: ptrArgUint64FromUint64(blockNumber.Uint64()),
							BlockHash:   ptrHash(receipt.BlockHash),
							TxIndex:     ptrArgUint64FromUint(receipt.TransactionIndex),
							ChainID:     argBig(*tx.ChainId()),
							Type:        argUint64(tx.Type()),
							V:           argBig(*V),
							R:           argBig(*R),
							S:           argBig(*S),
						},
					)

					batchTxs = append(batchTxs, *tx)
				}
				m.State.
					On("GetTransactionsByBatchNumber", context.Background(), tc.Number.Uint64(), m.DbTx).
					Return(batchTxs, nil).
					Once()
			},
		},
		// {
		// 	Name:   "get latest batch successfully",
		// 	Number: nil,
		// 	ExpectedResult: types.NewBatch(
		// 		&types.Header{Number: big.NewInt(2), UncleHash: types.EmptyUncleHash, Root: types.EmptyRootHash},
		// 		[]*types.Transaction{types.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})},
		// 		nil,
		// 		[]*types.Receipt{types.NewReceipt([]byte{}, false, uint64(0))},
		// 		&trie.StackTrie{},
		// 	),
		// 	ExpectedError: nil,
		// 	SetupMocks: func(m *mocks, tc *testCase) {
		// 		m.DbTx.
		// 			On("Commit", context.Background()).
		// 			Return(nil).
		// 			Once()

		// 		m.State.
		// 			On("BeginStateTransaction", context.Background()).
		// 			Return(m.DbTx, nil).
		// 			Once()

		// 		m.State.
		// 			On("GetLastBatchNumber", context.Background(), m.DbTx).
		// 			Return(tc.ExpectedResult.Number().Uint64(), nil).
		// 			Once()

		// 		m.State.
		// 			On("GetBatchByNumber", context.Background(), tc.ExpectedResult.Number().Uint64(), m.DbTx).
		// 			Return(tc.ExpectedResult, nil).
		// 			Once()
		// 	},
		// },
		// {
		// 	Name:           "get latest batch fails to compute batch number",
		// 	Number:         nil,
		// 	ExpectedResult: nil,
		// 	ExpectedError:  newRPCError(defaultErrorCode, "failed to get the last batch number from state"),
		// 	SetupMocks: func(m *mocks, tc *testCase) {
		// 		m.DbTx.
		// 			On("Rollback", context.Background()).
		// 			Return(nil).
		// 			Once()

		// 		m.State.
		// 			On("BeginStateTransaction", context.Background()).
		// 			Return(m.DbTx, nil).
		// 			Once()

		// 		m.State.
		// 			On("GetLastBatchNumber", context.Background(), m.DbTx).
		// 			Return(uint64(0), errors.New("failed to get last batch number")).
		// 			Once()
		// 	},
		// },
		// {
		// 	Name:           "get latest batch fails to load batch by number",
		// 	Number:         nil,
		// 	ExpectedResult: nil,
		// 	ExpectedError:  newRPCError(defaultErrorCode, "couldn't load batch from state by number 1"),
		// 	SetupMocks: func(m *mocks, tc *testCase) {
		// 		m.DbTx.
		// 			On("Rollback", context.Background()).
		// 			Return(nil).
		// 			Once()

		// 		m.State.
		// 			On("BeginStateTransaction", context.Background()).
		// 			Return(m.DbTx, nil).
		// 			Once()

		// 		m.State.
		// 			On("GetLastBatchNumber", context.Background(), m.DbTx).
		// 			Return(uint64(1), nil).
		// 			Once()

		// 		m.State.
		// 			On("GetBatchByNumber", context.Background(), uint64(1), m.DbTx).
		// 			Return(nil, errors.New("failed to load batch by number")).
		// 			Once()
		// 	},
		// },
		// {
		// 	Name:           "get pending batch successfully",
		// 	Number:         big.NewInt(-1),
		// 	ExpectedResult: types.NewBatch(&types.Header{Number: big.NewInt(2)}, nil, nil, nil, &trie.StackTrie{}),
		// 	ExpectedError:  nil,
		// 	SetupMocks: func(m *mocks, tc *testCase) {
		// 		lastBatchHeader := types.CopyHeader(tc.ExpectedResult.Header())
		// 		lastBatchHeader.Number.Sub(lastBatchHeader.Number, big.NewInt(1))
		// 		lastBatch := types.NewBatch(lastBatchHeader, nil, nil, nil, &trie.StackTrie{})

		// 		expectedResultHeader := types.CopyHeader(tc.ExpectedResult.Header())
		// 		expectedResultHeader.ParentHash = lastBatch.Hash()
		// 		tc.ExpectedResult = types.NewBatch(expectedResultHeader, nil, nil, nil, &trie.StackTrie{})

		// 		m.DbTx.
		// 			On("Commit", context.Background()).
		// 			Return(nil).
		// 			Once()

		// 		m.State.
		// 			On("BeginStateTransaction", context.Background()).
		// 			Return(m.DbTx, nil).
		// 			Once()

		// 		m.State.
		// 			On("GetLastBatch", context.Background(), m.DbTx).
		// 			Return(lastBatch, nil).
		// 			Once()
		// 	},
		// },
		// {
		// 	Name:           "get pending batch fails",
		// 	Number:         big.NewInt(-1),
		// 	ExpectedResult: nil,
		// 	ExpectedError:  newRPCError(defaultErrorCode, "couldn't load last batch from state to compute the pending batch"),
		// 	SetupMocks: func(m *mocks, tc *testCase) {
		// 		m.DbTx.
		// 			On("Rollback", context.Background()).
		// 			Return(nil).
		// 			Once()

		// 		m.State.
		// 			On("BeginStateTransaction", context.Background()).
		// 			Return(m.DbTx, nil).
		// 			Once()

		// 		m.State.
		// 			On("GetLastBatch", context.Background(), m.DbTx).
		// 			Return(nil, errors.New("failed to load last batch")).
		// 			Once()
		// 	},
		// },
	}

	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			testCase.SetupMocks(s, m, &tc)

			res, err := s.JSONRPCCall("zkevm_getBatchByNumber", hex.EncodeBig(tc.Number), false)
			require.NoError(t, err)
			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result interface{}
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)

				if result != nil || testCase.ExpectedResult != nil {
					var batch map[string]interface{}
					err = json.Unmarshal(res.Result, &batch)
					require.NoError(t, err)
					assert.Equal(t, tc.ExpectedResult.Number.Hex(), batch["number"].(string))
					assert.Equal(t, tc.ExpectedResult.Coinbase.String(), batch["coinbase"].(string))
					assert.Equal(t, tc.ExpectedResult.StateRoot.String(), batch["stateRoot"].(string))
					assert.Equal(t, tc.ExpectedResult.GlobalExitRoot.String(), batch["globalExitRoot"].(string))
					assert.Equal(t, tc.ExpectedResult.AccInputHash.String(), batch["accInputHash"].(string))
					assert.Equal(t, tc.ExpectedResult.Timestamp.Hex(), batch["timestamp"].(string))
					assert.Equal(t, tc.ExpectedResult.SendSequencesTxHash.String(), batch["sendSequencesTxHash"].(string))
					assert.Equal(t, tc.ExpectedResult.VerifyBatchTxHash.String(), batch["verifyBatchTxHash"].(string))
					batchTxs := batch["transactions"].([]interface{})
					for i, tx := range tc.ExpectedResult.Transactions {
						switch batchTxOrHash := batchTxs[i].(type) {
						case string:
							assert.Equal(t, tx.getHash().String(), batchTxOrHash)
						case map[string]interface{}:
							tx := tx.(rpcTransaction)
							assert.Equal(t, tx.Nonce.Hex(), batchTxOrHash["nonce"].(string))
							assert.Equal(t, tx.GasPrice.Hex(), batchTxOrHash["gasPrice"].(string))
							assert.Equal(t, tx.Gas.Hex(), batchTxOrHash["gas"].(string))
							assert.Equal(t, tx.To.String(), batchTxOrHash["to"].(string))
							assert.Equal(t, tx.Value.Hex(), batchTxOrHash["value"].(string))
							assert.Equal(t, tx.Input.Hex(), batchTxOrHash["input"].(string))
							assert.Equal(t, tx.V.Hex(), batchTxOrHash["v"].(string))
							assert.Equal(t, tx.R.Hex(), batchTxOrHash["r"].(string))
							assert.Equal(t, tx.S.Hex(), batchTxOrHash["s"].(string))
							assert.Equal(t, tx.Hash.String(), batchTxOrHash["hash"].(string))
							assert.Equal(t, strings.ToLower(tx.From.String()), strings.ToLower(batchTxOrHash["from"].(string)))
							assert.Equal(t, tx.BlockHash.String(), batchTxOrHash["blockHash"].(string))
							assert.Equal(t, tx.BlockNumber.Hex(), batchTxOrHash["blockNumber"].(string))
							assert.Equal(t, tx.TxIndex.Hex(), batchTxOrHash["transactionIndex"].(string))
							assert.Equal(t, tx.ChainID.Hex(), batchTxOrHash["chainId"].(string))
							assert.Equal(t, tx.Type.Hex(), batchTxOrHash["type"].(string))
						}
					}
				}
			}

			if res.Error != nil || testCase.ExpectedError != nil {
				assert.Equal(t, testCase.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, testCase.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}

func ptrUint64(n uint64) *uint64 {
	return &n
}

func ptrArgUint64FromUint(n uint) *argUint64 {
	tmp := argUint64(n)
	return &tmp
}

func ptrArgUint64FromUint64(n uint64) *argUint64 {
	tmp := argUint64(n)
	return &tmp
}

func ptrHash(h common.Hash) *common.Hash {
	return &h
}

func signTx(tx *types.Transaction, chainID uint64) *types.Transaction {
	privateKey, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(0).SetUint64(chainID))
	signedTx, _ := auth.Signer(auth.From, tx)
	return signedTx
}
