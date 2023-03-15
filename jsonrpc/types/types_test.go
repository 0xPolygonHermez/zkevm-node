package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgHashUnmarshalFromShortString(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult string
		expectedError  error
	}
	testCases := []testCase{
		{
			name:           "valid hex value starting with 0x",
			input:          "0x1",
			expectedResult: "0x0000000000000000000000000000000000000000000000000000000000000001",
			expectedError:  nil,
		},
		{
			name:           "valid hex value starting without 0x",
			input:          "1",
			expectedResult: "0x0000000000000000000000000000000000000000000000000000000000000001",
			expectedError:  nil,
		},
		{
			name:           "valid full hash value",
			input:          "0x05b21ee5f65c28a0af8e71290fc33625a1279a8b3d6357ce3ca60f22dbf59e63",
			expectedResult: "0x05b21ee5f65c28a0af8e71290fc33625a1279a8b3d6357ce3ca60f22dbf59e63",
			expectedError:  nil,
		},
		{
			name:           "invalid hex value starting with 0x",
			input:          "0xG",
			expectedResult: "0x0000000000000000000000000000000000000000000000000000000000000000",
			expectedError:  fmt.Errorf("invalid hash, it needs to be a hexadecimal value"),
		},
		{
			name:           "invalid hex value starting without 0x",
			input:          "G",
			expectedResult: "0x0000000000000000000000000000000000000000000000000000000000000000",
			expectedError:  fmt.Errorf("invalid hash, it needs to be a hexadecimal value"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			arg := ArgHash{}
			err := arg.UnmarshalText([]byte(testCase.input))
			require.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, arg.Hash().String())
		})
	}
}

func TestArgAddressUnmarshalFromShortString(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult string
		expectedError  error
	}
	testCases := []testCase{
		{
			name:           "valid hex value starting with 0x",
			input:          "0x1",
			expectedResult: "0x0000000000000000000000000000000000000001",
			expectedError:  nil,
		},
		{
			name:           "valid hex value starting without 0x",
			input:          "1",
			expectedResult: "0x0000000000000000000000000000000000000001",
			expectedError:  nil,
		},
		{
			name:           "valid full address value",
			input:          "0x748964F22eFd023eB78A246A7AC2506e84CC4545",
			expectedResult: "0x748964F22eFd023eB78A246A7AC2506e84CC4545",
			expectedError:  nil,
		},
		{
			name:           "invalid hex value starting with 0x",
			input:          "0xG",
			expectedResult: "0x0000000000000000000000000000000000000000",
			expectedError:  fmt.Errorf("invalid address, it needs to be a hexadecimal value"),
		},
		{
			name:           "invalid hex value starting without 0x",
			input:          "G",
			expectedResult: "0x0000000000000000000000000000000000000000",
			expectedError:  fmt.Errorf("invalid address, it needs to be a hexadecimal value"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			arg := ArgAddress{}
			err := arg.UnmarshalText([]byte(testCase.input))
			require.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, arg.Address().String())
		})
	}
}
