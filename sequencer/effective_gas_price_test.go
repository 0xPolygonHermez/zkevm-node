package sequencer

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcGasPriceEffectivePercentage(t *testing.T) {
	testCases := []struct {
		name          string
		breakEven     *big.Int
		gasPrice      *big.Int
		expectedValue uint8
		err           error
	}{
		{
			name:          "Nil breakEven or gasPrice",
			gasPrice:      big.NewInt(1),
			expectedValue: uint8(0),
		},
		{
			name:          "Zero breakEven or gasPrice",
			breakEven:     big.NewInt(1),
			gasPrice:      big.NewInt(0),
			expectedValue: uint8(0),
		},
		{
			name:          "Both positive, gasPrice less than breakEven",
			breakEven:     big.NewInt(22000000000),
			gasPrice:      big.NewInt(11000000000),
			expectedValue: uint8(255),
		},
		{
			name:          "Both positive, gasPrice more than breakEven",
			breakEven:     big.NewInt(19800000000),
			gasPrice:      big.NewInt(22000000000),
			expectedValue: uint8(230),
		},
		{
			name:          "100% (255) effective percentage 1",
			gasPrice:      big.NewInt(22000000000),
			breakEven:     big.NewInt(22000000000),
			expectedValue: 255,
		},
		{
			name:          "100% (255) effective percentage 2",
			gasPrice:      big.NewInt(22000000000),
			breakEven:     big.NewInt(21999999999),
			expectedValue: 255,
		},
		{
			name:          "100% (255) effective percentage 3",
			gasPrice:      big.NewInt(22000000000),
			breakEven:     big.NewInt(21900000000),
			expectedValue: 254,
		},
		{
			name:          "50% (127) effective percentage",
			gasPrice:      big.NewInt(22000000000),
			breakEven:     big.NewInt(11000000000),
			expectedValue: 127,
		},
		{
			name:          "(40) effective percentage",
			gasPrice:      big.NewInt(1000),
			breakEven:     big.NewInt(157),
			expectedValue: 40,
		},
		{
			name:          "(1) effective percentage",
			gasPrice:      big.NewInt(1000),
			breakEven:     big.NewInt(1),
			expectedValue: 0,
		},
		{
			name:          "(2) effective percentage",
			gasPrice:      big.NewInt(1000),
			breakEven:     big.NewInt(4),
			expectedValue: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, _ := CalculateEffectiveGasPricePercentage(tc.gasPrice, tc.breakEven)
			assert.Equal(t, tc.err, err)
			if actual != 0 {
				assert.Equal(t, tc.expectedValue, actual)
			} else {
				assert.Zero(t, tc.expectedValue)
			}
		})
	}
}
