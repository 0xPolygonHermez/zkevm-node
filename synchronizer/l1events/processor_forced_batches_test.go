package l1events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessorForcedBatchesName(t *testing.T) {
	sut := NewProcessForcedBatches(nil)
	name := sut.Name()
	require.Equal(t, "ProcessForcedBatches", name)
}
