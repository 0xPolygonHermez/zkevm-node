package sequencer

import "errors"

var (
	// ErrBatchRemainingResourcesOverflow error returned when the resources for the current batch are overflown
	ErrBatchRemainingResourcesOverflow = errors.New("overflow of remaining resources for current batch %s")
)
