package elderberry

import "github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"

var (
	// ForkIDElderberry is the forkId for Elderberry
	ForkIDElderberry = actions.ForkIdType(8) //nolint:gomnd
	// ForksIdOnlyElderberry support only elderberry forkId
	ForksIdOnlyElderberry = []actions.ForkIdType{ForkIDElderberry}
)
