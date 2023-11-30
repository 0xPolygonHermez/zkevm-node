package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
	s, m, _ := newSequencerMockedServer(t)
	defer s.Stop()

	type testCase struct {
		Name           string
		ExpectedResult TxPoolStatusResponse
		ExpectedError  types.Error
		SetupMocks     func(m *mocksWrapper, tc testCase)
	}

	testCases := []testCase{
		{
			Name: "Successfully count pending and queued txs from txpool",
			ExpectedResult: TxPoolStatusResponse{
				Pending: hexutil.Uint(10),
				Queued:  hexutil.Uint(2),
			},
			ExpectedError: nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.Pool.
					On("CountPendingTransactions", context.Background()).
					Return(uint64(10), nil).
					Once()
				m.Pool.
					On("CountQueuedTransactions", context.Background()).
					Return(uint64(2), nil).
					Once()
			},
		},
		{
			Name: "Successfully count pending and queued txs from empty txpool",
			ExpectedResult: TxPoolStatusResponse{
				Pending: hexutil.Uint(0),
				Queued:  hexutil.Uint(0),
			},
			ExpectedError: nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.Pool.
					On("CountPendingTransactions", context.Background()).
					Return(uint64(0), nil).
					Once()
				m.Pool.
					On("CountQueuedTransactions", context.Background()).
					Return(uint64(0), nil).
					Once()
			},
		},
		{
			Name:           "Failed to count pending txs from txpool",
			ExpectedResult: TxPoolStatusResponse{},
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "Failed to count pending txs from pool"),
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.Pool.
					On("CountPendingTransactions", context.Background()).
					Return(uint64(0), errors.New("failed to get pending tx from txpool")).
					Once()
			},
		},
		{
			Name:           "Failed to count queued txs from txpool",
			ExpectedResult: TxPoolStatusResponse{},
			ExpectedError:  types.NewRPCError(types.DefaultErrorCode, "Failed to count queued txs from pool"),
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.Pool.
					On("CountPendingTransactions", context.Background()).
					Return(uint64(0), nil).
					Once()
				m.Pool.
					On("CountQueuedTransactions", context.Background()).
					Return(uint64(0), errors.New("failed to get queued tx from txpool")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m, tc)
			res, err := s.JSONRPCCall("txpool_status")

			require.NoError(t, err)
			assert.Equal(t, float64(1), res.ID)
			assert.Equal(t, "2.0", res.JSONRPC)

			if res.Result != nil {
				var result TxPoolStatusResponse
				err = json.Unmarshal(res.Result, &result)
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, result)
			}

			if res.Error != nil || testCase.ExpectedError != nil {
				assert.Equal(t, testCase.ExpectedError.ErrorCode(), res.Error.Code)
				assert.Equal(t, testCase.ExpectedError.Error(), res.Error.Message)
			}
		})
	}
}
