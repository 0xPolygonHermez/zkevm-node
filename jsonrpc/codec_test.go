package jsonrpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

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

func TestSuccessResponseMarshal(t *testing.T) {
	const jsonRPCValue = "jsonrpc"
	const idValue = 1
	result, err := json.Marshal(struct {
		A string `json:"A"`
	}{"A"})
	require.NoError(t, err)
	errorObjValue := &ErrorObject{Code: 123, Message: "m", Data: "abc"}

	expectedBytes := []byte(fmt.Sprintf("{\"jsonrpc\":\"%v\",\"id\":%v,\"result\":{\"A\":\"A\"},\"error\":{\"code\":123,\"message\":\"m\",\"data\":\"abc\"}}", jsonRPCValue, idValue))

	resp := Response{
		JSONRPC: jsonRPCValue,
		ID:      idValue,
		Result:  result,
		Error:   errorObjValue,
	}

	bytes, err := json.Marshal(resp)
	require.NoError(t, err)
	assert.Equal(t, hex.EncodeToString(expectedBytes), hex.EncodeToString(bytes), fmt.Sprintf("expected:%v\nfound:%v", string(expectedBytes), string(bytes)))
}
