package jsonrpc

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

type RPCError struct {
	err  string
	code int
}

func newRPCError(code int, err string) *RPCError {
	return &RPCError{code: code, err: err}
}

func (e *RPCError) Error() string {
	return e.err
}

func (e *RPCError) ErrorCode() int {
	return e.code
}
