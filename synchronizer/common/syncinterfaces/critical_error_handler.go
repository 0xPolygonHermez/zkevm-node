package syncinterfaces

import "context"

// CriticalErrorHandler is an interface for handling critical errors. Before that class this was called Halt()
type CriticalErrorHandler interface {
	// CriticalError is called when a critical error occurs. The error is passed in as a parameter.
	// this function could be blocking or non-blocking, depending on the implementation.
	CriticalError(ctx context.Context, err error)
}
