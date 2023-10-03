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
			name:           "valid duration I",
			input:          "60s",
			expectedResult: &Duration{Duration: time.Minute},
		},
		{
			name:           "valid duration II",
			input:          "1m0s",
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

func TestDurationMarshalText(t *testing.T) {
	type testCase struct {
		name           string
		input          *Duration
		expectedResult string
	}

	testCases := []testCase{
		{
			name:           "valid duration",
			input:          &Duration{Duration: time.Minute},
			expectedResult: `"1m0s"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			byteDuration, err := json.Marshal(testCase.input)
			require.NoError(t, err)
			require.Equal(t, string(byteDuration), testCase.expectedResult)
		})
	}
}

func TestDurationMarshalJSON(t *testing.T) {
	type foo struct {
		F *Duration `json:"f"`
	}
	type bar struct {
		B Duration `json:"f"`
	}
	type testCase struct {
		name     string
		inputFoo foo
		inputBar bar
		expected string
	}
	testCases := []testCase{
		{
			name:     "valid duration",
			inputFoo: foo{F: &Duration{Duration: time.Minute}},
			inputBar: bar{B: Duration{Duration: time.Minute}},
			expected: `{"f":"1m0s"}`,
		},
		{
			name:     "valid duration",
			inputFoo: foo{F: &Duration{Duration: time.Hour + time.Minute}},
			inputBar: bar{B: Duration{Duration: time.Hour + time.Minute}},
			expected: `{"f":"1h1m0s"}`,
		},
		{
			name:     "valid duration",
			inputFoo: foo{F: &Duration{Duration: time.Hour}},
			inputBar: bar{B: Duration{Duration: time.Hour}},
			expected: `{"f":"1h0m0s"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fooBytes, err := json.Marshal(testCase.inputFoo)
			require.NoError(t, err)
			require.Equal(t, string(fooBytes), testCase.expected)
			barBytes, err := json.Marshal(testCase.inputBar)
			require.NoError(t, err)
			require.Equal(t, string(barBytes), testCase.expected)
		})
	}
}
