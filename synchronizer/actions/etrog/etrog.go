package etrog

import "github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"

var (
	// ForkIDEtrog is the forkId for etrog
	ForkIDEtrog = actions.ForkIdType(7) //nolint:gomnd
	// ForksIdOnlyEtrog support only etrog forkId
	ForksIdOnlyEtrog = []actions.ForkIdType{ForkIDEtrog}
)
