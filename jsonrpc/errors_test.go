package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPCErrorConstants(t *testing.T) {
	assert.Equal(t, -32000, defaultErrorCode)
	assert.Equal(t, -32600, invalidRequestErrorCode)
	assert.Equal(t, -32601, notFoundErrorCode)
	assert.Equal(t, -32602, invalidParamsErrorCode)
	assert.Equal(t, -32700, parserErrorCode)
}

func TestRPCErrorMethods(t *testing.T) {
	const code, msg = 1, "err"

	var err rpcError = newRPCError(code, msg)

	assert.Equal(t, code, err.ErrorCode())
	assert.Equal(t, msg, err.Error())
}
