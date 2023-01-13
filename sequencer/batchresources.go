package sequencer

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/state"
)

// batchResources is a struct that contains the ZKEVM resources used by a batch/tx
type batchResources struct {
	zKCounters state.ZKCounters
	bytes      uint64
	gas        uint64 // TODO: Delete gas?
}

// sub subtracts the batch resources from other
func (r *batchResources) sub(other batchResources) error {
	// Gas
	if other.gas > r.gas {
		return fmt.Errorf("%w. Resource: Gas", ErrBatchRemainingResourcesUnderflow)
	}
	// Bytes
	if other.bytes > r.bytes {
		return fmt.Errorf("%w. Resource: Bytes", ErrBatchRemainingResourcesUnderflow)
	}
	bytesBackup := r.bytes
	gasBackup := r.gas
	r.bytes -= other.bytes
	r.gas -= other.gas
	err := r.zKCounters.Sub(other.zKCounters)
	if err != nil {
		return fmt.Errorf("%w. %s", ErrBatchRemainingResourcesUnderflow, err)
	}
	r.bytes = bytesBackup
	r.gas = gasBackup

	return err
}
