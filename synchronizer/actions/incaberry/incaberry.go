package incaberry

import "github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"

var (
	// ForkIDIncaberry is the forkId for incaberry
	ForkIDIncaberry = actions.ForkIdType(6) // nolint:gomnd
	// ForksIdToIncaberry support all forkIds till incaberry
	ForksIdToIncaberry = []actions.ForkIdType{1, 2, 3, 4, 5, ForkIDIncaberry}
)
