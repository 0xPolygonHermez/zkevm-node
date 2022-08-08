package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

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

func TestIsL2BlockConsolidated(t *testing.T) {
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
					On("IsL2BlockConsolidated", context.Background(), 1, m.DbTx).
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
					On("IsL2BlockConsolidated", context.Background(), 1, m.DbTx).
					Return(true, errors.New("failed to check if the block is consolidated")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_isL2BlockConsolidated", 1)
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

func TestIsL2BlockVirtualized(t *testing.T) {
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
					On("IsL2BlockVirtualized", context.Background(), 1, m.DbTx).
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
					On("IsL2BlockVirtualized", context.Background(), 1, m.DbTx).
					Return(true, errors.New("failed to check if the block is virtualized")).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(m)

			res, err := s.JSONRPCCall("zkevm_isL2BlockVirtualized", 1)
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

func ptrUint64(n uint64) *uint64 {
	return &n
}
