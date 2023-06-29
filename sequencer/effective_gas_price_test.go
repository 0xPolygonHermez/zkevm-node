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
			name:     "Nil breakEven or gasPrice",
			gasPrice: big.NewInt(1),
			err:      ErrBreakEvenGasPriceEmpty,
		},
		{
			name:      "Zero breakEven or gasPrice",
			breakEven: big.NewInt(1),
			gasPrice:  big.NewInt(0),
			err:       ErrBreakEvenGasPriceEmpty,
		},
		{
			name:          "Both positive, gasPrice less than breakEven",
			breakEven:     big.NewInt(22000000000),
			gasPrice:      big.NewInt(11000000000),
			expectedValue: uint8(255),
			err:           nil,
		},
		{
			name:          "Both positive, gasPrice more than breakEven",
			breakEven:     big.NewInt(19800000000),
			gasPrice:      big.NewInt(22000000000),
			expectedValue: uint8(231),
			err:           nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := CalculateEffectiveGasPricePercentage(tc.gasPrice, tc.breakEven)
			assert.Equal(t, tc.err, err)
			if actual != 0 {
				assert.Equal(t, tc.expectedValue, actual)
			} else {
				assert.Zero(t, tc.expectedValue)
			}
		})
	}
}
