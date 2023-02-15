package sequencer

import (
	"fmt"
)

var (
	// ErrBatchResourceBytesUnderflow happens when the batch runs out of bytes
	ErrBatchResourceBytesUnderflow = NewBatchRemainingResourcesUnderflowError(nil, "Bytes")
)

// BatchRemainingResourcesUnderflowError happens when the execution of a batch runs out of counters
type BatchRemainingResourcesUnderflowError struct {
	Message      string
	Code         int
	Err          error
	ResourceName string
}

// Error returns the error message
func (b BatchRemainingResourcesUnderflowError) Error() string {
	return constructErrorMsg(b.ResourceName)
}

// NewBatchRemainingResourcesUnderflowError creates a new BatchRemainingResourcesUnderflowError
func NewBatchRemainingResourcesUnderflowError(err error, resourceName string) error {
	return &BatchRemainingResourcesUnderflowError{
		Message:      constructErrorMsg(resourceName),
		Code:         1,
		Err:          err,
		ResourceName: resourceName,
	}
}

func constructErrorMsg(resourceName string) string {
	return fmt.Sprintf("underflow of remaining resources for current batch. Resource %s", resourceName)
}
