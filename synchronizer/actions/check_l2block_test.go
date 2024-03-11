package actions_test

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/stretchr/testify/assert"
)

func TestCheckL2BlockHash_GetNextL2BlockToCheck(t *testing.T) {
	// Create an instance of CheckL2BlockHash
	values := []struct {
		initial  uint64
		modulus  uint64
		expected uint64
	}{
		{0, 10, 10},
		{1, 10, 10},
		{9, 10, 10},
		{10, 10, 20},
		{0, 0, 1},
		{1, 0, 2},
	}
	for _, data := range values {
		// Call the GetNextL2BlockToCheck method
		checkL2Block := actions.NewCheckL2BlockHash(nil, nil, data.initial, data.modulus)
		nextL2Block := checkL2Block.GetNextL2BlockToCheck()

		// Assert the expected result
		assert.Equal(t, data.expected, nextL2Block)
	}
}
