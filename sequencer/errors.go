package sequencer

import "errors"

var (
	// ErrBatchRemainingResourcesUnderflow error returned when there is underflow of the resources for the current batch
	ErrBatchRemainingResourcesUnderflow = errors.New("underflow of remaining resources for current batch %s")
)
