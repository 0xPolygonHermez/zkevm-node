package types

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDurationUnmarshal(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult *Duration
		expectedErr    error
	}

	testCases := []testCase{
		{
			name:           "valid duration",
			input:          "60s",
			expectedResult: &Duration{Duration: time.Minute},
		},
		{
			name:        "int value",
			input:       "60",
			expectedErr: fmt.Errorf("time: missing unit in duration \"60\""),
		},
		{
			name:        "no duration value",
			input:       "abc",
			expectedErr: fmt.Errorf("time: invalid duration \"abc\""),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var d Duration
			input, err := json.Marshal(testCase.input)
			require.NoError(t, err)
			err = json.Unmarshal(input, &d)

			if testCase.expectedResult != nil {
				require.Equal(t, (*testCase.expectedResult).Nanoseconds(), d.Nanoseconds())
			}

			if err != nil {
				require.Equal(t, testCase.expectedErr.Error(), err.Error())
			}
		})
	}
}
