package jsonrpc

import (
	"fmt"
)

type detailedError interface {
	Error() string
	Code() int
}

type invalidParamsError struct {
	err string
}

func (e *invalidParamsError) Error() string {
	return e.err
}

func (e *invalidParamsError) Code() int {
	return -32602
}

type invalidRequestError struct {
	err string
}

func (e *invalidRequestError) Error() string {
	return e.err
}

func (e *invalidRequestError) Code() int {
	return -32600
}

type methodNotFoundError struct {
	err string
}

func (e *methodNotFoundError) Error() string {
	return e.err
}

func (e *methodNotFoundError) Code() int {
	return -32601
}

type genericError struct {
	err  string
	code int
}

func newGenericError(err string, code int) *genericError {
	return &genericError{err: err, code: code}
}

func (e *genericError) Error() string {
	return e.err
}

func (e *genericError) Code() int {
	return e.code
}

// NewMethodNotFoundError used when the RPC method does not exist or is not available
func newMethodNotFoundError(method string) *methodNotFoundError {
	e := &methodNotFoundError{fmt.Sprintf("the method %s does not exist/is not available", method)}
	return e
}

// NewInvalidRequestError used when the request is invalid
func newInvalidRequestError(msg string) *invalidRequestError {
	e := &invalidRequestError{msg}
	return e
}

// NewInvalidParamsError is used when the request has invalid parameters
func newInvalidParamsError(msg string) *invalidParamsError {
	e := &invalidParamsError{msg}
	return e
}
