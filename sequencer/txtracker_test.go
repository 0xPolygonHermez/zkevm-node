package sequencer

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/stretchr/testify/assert"
)

type efficiencyCalcTestCase struct {
	Name           string
	benefit        int64
	counters       state.ZKCounters
	usedBytes      uint64
	expectedResult float64
}

func TestTxTrackerEfficiencyCalculation(t *testing.T) {
	// Init ZKEVM resourceCostWeight values
	rcWeigth := batchResourceWeights{}
	rcWeigth.WeightCumulativeGasUsed = 1
	rcWeigth.WeightArithmetics = 1
	rcWeigth.WeightBinaries = 1
	rcWeigth.WeightKeccakHashes = 1
	rcWeigth.WeightMemAligns = 1
	rcWeigth.WeightPoseidonHashes = 1
	rcWeigth.WeightPoseidonPaddings = 1
	rcWeigth.WeightSteps = 1
	rcWeigth.WeightBatchBytesSize = 2

	// Init ZKEVM resourceCostMax values
	rcMax := batchConstraintsFloat64{}
	rcMax.maxCumulativeGasUsed = 10
	rcMax.maxArithmetics = 10
	rcMax.maxBinaries = 10
	rcMax.maxKeccakHashes = 10
	rcMax.maxMemAligns = 10
	rcMax.maxPoseidonHashes = 10
	rcMax.maxPoseidonPaddings = 10
	rcMax.maxSteps = 10
	rcMax.maxBatchBytesSize = 10

	totalWeight := float64(rcWeigth.WeightArithmetics + rcWeigth.WeightBatchBytesSize + rcWeigth.WeightBinaries + rcWeigth.WeightCumulativeGasUsed +
		rcWeigth.WeightKeccakHashes + rcWeigth.WeightMemAligns + rcWeigth.WeightPoseidonHashes + rcWeigth.WeightPoseidonPaddings + rcWeigth.WeightSteps)

	testCases := []efficiencyCalcTestCase{
		{
			Name:           "Using all of the resources",
			benefit:        1000000,
			counters:       state.ZKCounters{CumulativeGasUsed: 10, UsedKeccakHashes: 10, UsedPoseidonHashes: 10, UsedPoseidonPaddings: 10, UsedMemAligns: 10, UsedArithmetics: 10, UsedBinaries: 10, UsedSteps: 10},
			usedBytes:      10,
			expectedResult: 1000.00,
		},
		{
			Name:           "Using half of the resources",
			benefit:        1000000,
			counters:       state.ZKCounters{CumulativeGasUsed: 5, UsedKeccakHashes: 5, UsedPoseidonHashes: 5, UsedPoseidonPaddings: 5, UsedMemAligns: 5, UsedArithmetics: 5, UsedBinaries: 5, UsedSteps: 5},
			usedBytes:      5,
			expectedResult: 2000.00,
		},
		{
			Name:           "Using all the bytes and half of the remain resources",
			benefit:        1000000,
			counters:       state.ZKCounters{CumulativeGasUsed: 5, UsedKeccakHashes: 5, UsedPoseidonHashes: 5, UsedPoseidonPaddings: 5, UsedMemAligns: 5, UsedArithmetics: 5, UsedBinaries: 5, UsedSteps: 5},
			usedBytes:      10,
			expectedResult: 1666.67,
		},
		{
			Name:           "Using all the steps and half of the remain resources",
			benefit:        1000000,
			counters:       state.ZKCounters{CumulativeGasUsed: 5, UsedKeccakHashes: 5, UsedPoseidonHashes: 5, UsedPoseidonPaddings: 5, UsedMemAligns: 5, UsedArithmetics: 5, UsedBinaries: 5, UsedSteps: 10},
			usedBytes:      5,
			expectedResult: 1818.18,
		},
		{
			Name:           "Using 10% of all the resources",
			benefit:        1000000,
			counters:       state.ZKCounters{CumulativeGasUsed: 1, UsedKeccakHashes: 1, UsedPoseidonHashes: 1, UsedPoseidonPaddings: 1, UsedMemAligns: 1, UsedArithmetics: 1, UsedBinaries: 1, UsedSteps: 1},
			usedBytes:      1,
			expectedResult: 10000.00,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tx := TxTracker{}
			tx.Benefit = new(big.Int).SetInt64(testCase.benefit)

			tx.BatchResources.Bytes = testCase.usedBytes
			tx.updateZKCounters(testCase.counters, rcMax, rcWeigth)
			tx.weightMultipliers = calculateWeightMultipliers(rcWeigth, totalWeight)
			tx.resourceCostMultiplier = 1000
			tx.updateZKCounters(testCase.counters, rcMax, rcWeigth)
			t.Logf("%s=%s", testCase.Name, fmt.Sprintf("%.2f", tx.Efficiency))
			assert.Equal(t, fmt.Sprintf("%.2f", testCase.expectedResult), fmt.Sprintf("%.2f", tx.Efficiency), "Efficiency calculation err. Expected=%s, Actual=%s", fmt.Sprintf("%.2f", testCase.expectedResult), fmt.Sprintf("%.2f", tx.Efficiency))
		})
	}
}
