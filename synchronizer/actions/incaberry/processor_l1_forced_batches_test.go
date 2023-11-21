package incaberry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessorForcedBatchesName(t *testing.T) {
	sut := NewProcessL1ForcedBatches(nil)
	name := sut.Name()
	require.Equal(t, "ProcessL1ForcedBatches", name)
}
