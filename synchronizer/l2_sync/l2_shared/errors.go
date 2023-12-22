/*
// https://www.digitalocean.com/community/tutorials/how-to-add-extra-information-to-errors-in-go

This error DeSyncPermissionlessAndTrustedNodeError have a field L1BlockNumber that contains the block number where the discrepancy is.
*/
package l2_shared

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
)

// DeSyncPermissionlessAndTrustedNodeError is an error type that contains the Block where is the discrepancy
type DeSyncPermissionlessAndTrustedNodeError struct {
	L1BlockNumber uint64
	Err           error
}

// NewDeSyncPermissionlessAndTrustedNodeError returns a new instance of DeSyncPermissionlessAndTrustedNodeError
func NewDeSyncPermissionlessAndTrustedNodeError(L1BlockNumber uint64) *DeSyncPermissionlessAndTrustedNodeError {
	return &DeSyncPermissionlessAndTrustedNodeError{
		L1BlockNumber: L1BlockNumber,
		Err:           syncinterfaces.ErrFatalDesyncFromL1,
	}
}

func (e *DeSyncPermissionlessAndTrustedNodeError) Error() string {
	return fmt.Sprintf("DeSyncPermissionlessAndTrustedNode. Block:%d Err: %s", e.L1BlockNumber, e.Err)
}

func (e *DeSyncPermissionlessAndTrustedNodeError) Unwrap() error {
	return e.Err
}
