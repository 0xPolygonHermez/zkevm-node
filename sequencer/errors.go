package sequencer

import (
	"fmt"
)

var (
	ErrBatchResourceBytesUnderflow = NewBatchRemainingResourcesUnderflowError(nil, "Bytes")
)

type BatchRemainingResourcesUnderflowError struct {
	Message      string
	Code         int
	Err          error
	ResourceName string
}

func (b BatchRemainingResourcesUnderflowError) Error() string {
	return constructErrorMsg(b.ResourceName)
}

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
