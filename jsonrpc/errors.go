package jsonrpc

import "fmt"

const (
	defaultErrorCode        = -32000
	invalidRequestErrorCode = -32600
	notFoundErrorCode       = -32601
	invalidParamsErrorCode  = -32602
	parserErrorCode         = -32700
)

type rpcError interface {
	Error() string
	ErrorCode() int
}

// RPCError represents an RPC error.
type RPCError struct {
	err  string
	code int
}

func newRPCError(code int, err string, args ...interface{}) *RPCError {
	var errMessage string
	if len(args) > 0 {
		errMessage = fmt.Sprintf(err, args...)
	} else {
		errMessage = err
	}
	return &RPCError{code: code, err: errMessage}
}

// Error returns the error message.
func (e *RPCError) Error() string {
	return e.err
}

// ErrorCode returns the error code.
func (e *RPCError) ErrorCode() int {
	return e.code
}
