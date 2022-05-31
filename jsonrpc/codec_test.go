package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/hermeznetwork/hermez-core/hex"
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
	s := newStateMock(t)

	type testCase struct {
		name                string
		bn                  *BlockNumber
		expectedBlockNumber uint64
		expectedError       error
		setupMocks          func(s *stateMock, t *testCase)
	}

	testCases := []testCase{
		{
			name:                "BlockNumber nil",
			bn:                  nil,
			expectedBlockNumber: 0,
			expectedError:       nil,
			setupMocks:          func(s *stateMock, t *testCase) {},
		},
		{
			name:                "BlockNumber LatestBlockNumber",
			bn:                  bnPtr(LatestBlockNumber),
			expectedBlockNumber: 50,
			expectedError:       nil,
			setupMocks: func(s *stateMock, t *testCase) {
				s.
					On("GetLastBatchNumber", context.Background(), "").
					Return(uint64(50), nil).
					Once()
			},
		},
		{
			name:                "BlockNumber PendingBlockNumber",
			bn:                  bnPtr(PendingBlockNumber),
			expectedBlockNumber: 30,
			expectedError:       nil,
			setupMocks: func(s *stateMock, t *testCase) {
				s.
					On("GetLastBatchNumber", context.Background(), "").
					Return(uint64(30), nil).
					Once()
			},
		},
		{
			name:                "BlockNumber EarliestBlockNumber",
			bn:                  bnPtr(EarliestBlockNumber),
			expectedBlockNumber: 0,
			expectedError:       nil,
			setupMocks:          func(s *stateMock, t *testCase) {},
		},
		{
			name:                "BlockNumber Positive Number",
			bn:                  bnPtr(BlockNumber(int64(10))),
			expectedBlockNumber: 10,
			expectedError:       nil,
			setupMocks:          func(s *stateMock, t *testCase) {},
		},
		{
			name:                "BlockNumber Negative Number <= -4",
			bn:                  bnPtr(BlockNumber(int64(-4))),
			expectedBlockNumber: 0,
			expectedError:       fmt.Errorf("invalid block number: -4"),
			setupMocks:          func(s *stateMock, t *testCase) {},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tc := testCase
			testCase.setupMocks(s, &tc)
			result, err := testCase.bn.getNumericBlockNumber(context.Background(), s)
			assert.Equal(t, testCase.expectedBlockNumber, result)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestResponseMarshal(t *testing.T) {
	const jsonRPCValue = "jsonrpc"
	const idValue = 1
	result, err := json.Marshal(struct {
		A string `json:"A"`
	}{"A"})
	require.NoError(t, err)
	errorObjValue := newRPCError(123, "m")

	expectedBytes := []byte(fmt.Sprintf("{\"jsonrpc\":\"%v\",\"id\":%v,\"result\":{\"A\":\"A\"},\"error\":{\"code\":123,\"message\":\"m\"}}", jsonRPCValue, idValue))

	req := Request{
		JSONRPC: jsonRPCValue,
		ID:      idValue,
	}

	resp := NewResponse(req, &result, errorObjValue)

	bytes, err := json.Marshal(resp)
	require.NoError(t, err)
	assert.Equal(t, hex.EncodeToString(expectedBytes), hex.EncodeToString(bytes), fmt.Sprintf("expected:%v\nfound:%v", string(expectedBytes), string(bytes)))
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

func bnPtr(bn BlockNumber) *BlockNumber {
	return &bn
}
