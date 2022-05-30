package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetailedErrors(t *testing.T) {

	testCases := []struct {
		name          string
		err           detailedError
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Invalid Request Error",
			err:           newInvalidRequestError("err message 1"),
			expectedCode:  -32600,
			expectedError: "err message 1",
		},
		{
			name:          "Method Not Found Error",
			err:           newMethodNotFoundError("example_method_name"),
			expectedCode:  -32601,
			expectedError: "the method example_method_name does not exist/is not available",
		},
		{
			name:          "Invalid Params Error",
			err:           newInvalidParamsError("err message 2"),
			expectedCode:  -32602,
			expectedError: "err message 2",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedCode, testCase.err.Code())
			assert.Equal(t, testCase.expectedError, testCase.err.Error())
		})
	}
}
