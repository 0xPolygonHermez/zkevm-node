package types

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockNumberMarshalJSON(t *testing.T) {
	testCases := []struct {
		jsonValue           string
		expectedBlockNumber int64
		expectedError       error
	}{
		{"latest", int64(LatestBlockNumber), nil},
		{"pending", int64(PendingBlockNumber), nil},
		{"earliest", int64(EarliestBlockNumber), nil},
		{"safe", int64(SafeBlockNumber), nil},
		{"finalized", int64(FinalizedBlockNumber), nil},
		{"", int64(LatestBlockNumber), nil},
		{"0", int64(0), nil},
		{"10", int64(10), nil},
		{"0x2", int64(2), nil},
		{"0xA", int64(10), nil},
		{"abc", int64(0), &strconv.NumError{Err: strconv.ErrSyntax, Func: "ParseUint", Num: "abc"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.jsonValue, func(t *testing.T) {
			data, err := json.Marshal(testCase.jsonValue)
			require.NoError(t, err)
			bn := BlockNumber(int64(0))
			err = json.Unmarshal(data, &bn)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedBlockNumber, int64(bn))
		})
	}
}

func TestGetNumericBlockNumber(t *testing.T) {
	s := mocks.NewStateMock(t)
	e := mocks.NewEthermanMock(t)

	type testCase struct {
		name                string
		bn                  *BlockNumber
		expectedBlockNumber uint64
		expectedError       Error
		setupMocks          func(s *mocks.StateMock, d *mocks.DBTxMock, t *testCase)
	}

	testCases := []testCase{
		{
			name:                "BlockNumber nil",
			bn:                  nil,
			expectedBlockNumber: 40,
			expectedError:       nil,
			setupMocks: func(s *mocks.StateMock, d *mocks.DBTxMock, t *testCase) {
				s.
					On("GetLastL2BlockNumber", context.Background(), d).
					Return(uint64(40), nil).
					Once()
			},
		},
		{
			name:                "BlockNumber LatestBlockNumber",
			bn:                  bnPtr(LatestBlockNumber),
			expectedBlockNumber: 50,
			expectedError:       nil,
			setupMocks: func(s *mocks.StateMock, d *mocks.DBTxMock, t *testCase) {
				s.
					On("GetLastL2BlockNumber", context.Background(), d).
					Return(uint64(50), nil).
					Once()
			},
		},
		{
			name:                "BlockNumber PendingBlockNumber",
			bn:                  bnPtr(PendingBlockNumber),
			expectedBlockNumber: 30,
			expectedError:       nil,
			setupMocks: func(s *mocks.StateMock, d *mocks.DBTxMock, t *testCase) {
				s.
					On("GetLastL2BlockNumber", context.Background(), d).
					Return(uint64(30), nil).
					Once()
			},
		},
		{
			name:                "BlockNumber EarliestBlockNumber",
			bn:                  bnPtr(EarliestBlockNumber),
			expectedBlockNumber: 0,
			expectedError:       nil,
			setupMocks:          func(s *mocks.StateMock, d *mocks.DBTxMock, t *testCase) {},
		},
		{
			name:                "BlockNumber SafeBlockNumber",
			bn:                  bnPtr(SafeBlockNumber),
			expectedBlockNumber: 40,
			expectedError:       nil,
			setupMocks: func(s *mocks.StateMock, d *mocks.DBTxMock, t *testCase) {
				liSafeBlock := uint64(30)
				e.
					On("GetSafeBlockNumber", context.Background()).
					Return(liSafeBlock, nil).
					Once()

				s.
					On("GetSafeL2BlockNumber", context.Background(), liSafeBlock, d).
					Return(uint64(40), nil).
					Once()
			},
		},
		{
			name:                "BlockNumber FinalizedBlockNumber",
			bn:                  bnPtr(FinalizedBlockNumber),
			expectedBlockNumber: 60,
			expectedError:       nil,
			setupMocks: func(s *mocks.StateMock, d *mocks.DBTxMock, t *testCase) {
				liFinalizedBlock := uint64(50)
				e.
					On("GetFinalizedBlockNumber", context.Background()).
					Return(liFinalizedBlock, nil).
					Once()

				s.
					On("GetFinalizedL2BlockNumber", context.Background(), liFinalizedBlock, d).
					Return(uint64(60), nil).
					Once()
			},
		},
		{
			name:                "BlockNumber Positive Number",
			bn:                  bnPtr(BlockNumber(int64(10))),
			expectedBlockNumber: 10,
			expectedError:       nil,
			setupMocks:          func(s *mocks.StateMock, d *mocks.DBTxMock, t *testCase) {},
		},
		{
			name:                "BlockNumber Negative Number <= -6",
			bn:                  bnPtr(BlockNumber(int64(-6))),
			expectedBlockNumber: 0,
			expectedError:       NewRPCError(InvalidParamsErrorCode, "invalid block number: -6"),
			setupMocks:          func(s *mocks.StateMock, d *mocks.DBTxMock, t *testCase) {},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tc := testCase
			dbTx := mocks.NewDBTxMock(t)
			testCase.setupMocks(s, dbTx, &tc)
			result, rpcErr := testCase.bn.GetNumericBlockNumber(context.Background(), s, e, dbTx)
			assert.Equal(t, testCase.expectedBlockNumber, result)
			if rpcErr != nil || testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.ErrorCode(), rpcErr.ErrorCode())
				assert.Equal(t, testCase.expectedError.Error(), rpcErr.Error())
			}
		})
	}
}

func TestResponseMarshal(t *testing.T) {
	testCases := []struct {
		Name    string
		JSONRPC string
		ID      interface{}
		Result  interface{}
		Error   Error

		ExpectedJSON string
	}{
		{
			Name:    "Error is nil",
			JSONRPC: "2.0",
			ID:      1,
			Result: struct {
				A string `json:"A"`
			}{"A"},
			Error: nil,

			ExpectedJSON: "{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"A\":\"A\"}}",
		},
		{
			Name:    "Result is nil and Error is not nil",
			JSONRPC: "2.0",
			ID:      1,
			Result:  nil,
			Error:   NewRPCError(123, "m"),

			ExpectedJSON: "{\"jsonrpc\":\"2.0\",\"id\":1,\"error\":{\"code\":123,\"message\":\"m\"}}",
		},
		{
			Name:    "Result is not nil and Error is not nil",
			JSONRPC: "2.0",
			ID:      1,
			Result: struct {
				A string `json:"A"`
			}{"A"},
			Error: NewRPCError(123, "m"),

			ExpectedJSON: "{\"jsonrpc\":\"2.0\",\"id\":1,\"error\":{\"code\":123,\"message\":\"m\"}}",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			req := Request{
				JSONRPC: testCase.JSONRPC,
				ID:      testCase.ID,
			}
			var result []byte
			if testCase.Result != nil {
				r, err := json.Marshal(testCase.Result)
				require.NoError(t, err)
				result = r
			}

			res := NewResponse(req, result, testCase.Error)
			bytes, err := json.Marshal(res)
			require.NoError(t, err)
			assert.Equal(t, string(testCase.ExpectedJSON), string(bytes))
		})
	}
}

func TestIndexUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		input         []byte
		expectedIndex int64
		expectedError error
	}{
		{
			input:         []byte("\"0x86\""),
			expectedIndex: 134,
			expectedError: nil,
		},
		{
			input:         []byte("\"abc\""),
			expectedIndex: 0,
			expectedError: &strconv.NumError{},
		},
	}

	for _, testCase := range testCases {
		var i Index
		err := json.Unmarshal(testCase.input, &i)
		assert.Equal(t, int64(testCase.expectedIndex), int64(i))
		assert.IsType(t, testCase.expectedError, err)
	}
}

func TestBlockNumberStringOrHex(t *testing.T) {
	testCases := []struct {
		bn             *BlockNumber
		expectedResult string
	}{
		{bn: bnPtr(BlockNumber(-3)), expectedResult: "pending"},
		{bn: bnPtr(BlockNumber(-2)), expectedResult: "latest"},
		{bn: bnPtr(BlockNumber(-1)), expectedResult: "earliest"},
		{bn: bnPtr(BlockNumber(0)), expectedResult: "0x0"},
		{bn: bnPtr(BlockNumber(100)), expectedResult: "0x64"},
	}

	for _, testCase := range testCases {
		result := testCase.bn.StringOrHex()
		assert.Equal(t, testCase.expectedResult, result)
	}
}

func TestBlockNumberOrHashMarshaling(t *testing.T) {
	type testCase struct {
		json           string
		expectedResult *BlockNumberOrHash
		expectedError  error
	}

	testCases := []testCase{
		// success
		{`{"blockNumber":"1"}`, &BlockNumberOrHash{number: bnPtr(BlockNumber(uint64(1)))}, nil},
		{`{"blockHash":"0x1"}`, &BlockNumberOrHash{hash: argHashPtr(common.HexToHash("0x1"))}, nil},
		{`{"blockHash":"0x1", "requireCanonical":true}`, &BlockNumberOrHash{hash: argHashPtr(common.HexToHash("0x1")), requireCanonical: true}, nil},
		// float wrong value
		{`{"blockNumber":1.0}`, &BlockNumberOrHash{}, fmt.Errorf("invalid blockNumber")},
		{`{"blockHash":1.0}`, &BlockNumberOrHash{}, fmt.Errorf("invalid blockHash")},
		{`{"blockHash":"0x1", "requireCanonical":1.0}`, &BlockNumberOrHash{}, fmt.Errorf("invalid requireCanonical")},
		// int wrong value
		{`{"blockNumber":1}`, &BlockNumberOrHash{}, fmt.Errorf("invalid blockNumber")},
		{`{"blockHash":1}`, &BlockNumberOrHash{}, fmt.Errorf("invalid blockHash")},
		{`{"blockHash":"0x1", "requireCanonical":1}`, &BlockNumberOrHash{}, fmt.Errorf("invalid requireCanonical")},
		// string wrong value
		{`{"blockNumber":"aaa"}`, &BlockNumberOrHash{}, fmt.Errorf("invalid blockNumber")},
		{`{"blockHash":"ggg"}`, &BlockNumberOrHash{}, fmt.Errorf("invalid blockHash")},
		{`{"blockHash":"0x1", "requireCanonical":"aaa"}`, &BlockNumberOrHash{}, fmt.Errorf("invalid requireCanonical")},
	}

	for _, testCase := range testCases {
		var result *BlockNumberOrHash
		err := json.Unmarshal([]byte(testCase.json), &result)

		assert.NotNil(t, result)
		assert.Equal(t, testCase.expectedResult.number, result.number)
		assert.Equal(t, testCase.expectedResult.hash, result.hash)
		assert.Equal(t, testCase.expectedResult.requireCanonical, testCase.expectedResult.requireCanonical)

		if testCase.expectedError == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, testCase.expectedError.Error(), err.Error())
		}
	}
}

func bnPtr(bn BlockNumber) *BlockNumber {
	return &bn
}

func argHashPtr(hash common.Hash) *ArgHash {
	h := ArgHash(hash)
	return &h
}
