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
		expectedValue *big.Int
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
			gasPrice:      big.NewInt(55000000000),
			expectedValue: big.NewInt(255),
			err:           nil,
		},
		{
			name:          "Both positive, gasPrice more than breakEven",
			breakEven:     big.NewInt(19800000000),
			gasPrice:      big.NewInt(22000000000),
			expectedValue: big.NewInt(231),
			err:           nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := CalcGasPriceEffectivePercentage(tc.breakEven, tc.gasPrice)
			assert.Equal(t, tc.err, err)
			if actual != nil {
				assert.Zero(t, tc.expectedValue.Cmp(actual))
			} else {
				assert.Nil(t, tc.expectedValue)
			}
		})
	}
}
