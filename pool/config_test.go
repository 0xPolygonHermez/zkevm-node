package pool

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
)

func TestIsWithinConstraints(t *testing.T) {
	cfg := BatchConstraintsCfg{
		MaxCumulativeGasUsed: 500,
		MaxKeccakHashes:      100,
		MaxPoseidonHashes:    200,
		MaxPoseidonPaddings:  150,
		MaxMemAligns:         1000,
		MaxArithmetics:       2000,
		MaxBinaries:          3000,
		MaxSteps:             4000,
	}

	testCases := []struct {
		desc     string
		counters state.ZKCounters
		expected bool
	}{
		{
			desc: "All constraints within limits",
			counters: state.ZKCounters{
				CumulativeGasUsed:    300,
				UsedKeccakHashes:     50,
				UsedPoseidonHashes:   100,
				UsedPoseidonPaddings: 75,
				UsedMemAligns:        500,
				UsedArithmetics:      1000,
				UsedBinaries:         2000,
				UsedSteps:            2000,
			},
			expected: true,
		},
		{
			desc: "All constraints exceed limits",
			counters: state.ZKCounters{
				CumulativeGasUsed:    600,
				UsedKeccakHashes:     150,
				UsedPoseidonHashes:   300,
				UsedPoseidonPaddings: 200,
				UsedMemAligns:        2000,
				UsedArithmetics:      3000,
				UsedBinaries:         4000,
				UsedSteps:            5000,
			},
			expected: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := cfg.IsWithinConstraints(tC.counters); got != tC.expected {
				t.Errorf("Expected %v, got %v", tC.expected, got)
			}
		})
	}
}
