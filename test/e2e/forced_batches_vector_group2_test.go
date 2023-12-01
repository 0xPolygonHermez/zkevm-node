package e2e

import (
	"testing"
)

func TestForcedBatchesVectorFilesGroup2(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	LaunchTestForcedBatchesVectorFilesGroup(t, "./../vectors/src/state-transition/forced-tx/group2")
}
