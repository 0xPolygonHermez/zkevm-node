package hex

import (
	"encoding/hex"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeBig(t *testing.T) {
	b := big.NewInt(math.MaxInt64)
	e := EncodeBig(b)
	d := DecodeBig(e)
	assert.Equal(t, b.Uint64(), d.Uint64())
}

// Define a struct for test cases
type TestCase struct {
	input  string
	output []byte
	err    error
}

// Unit test function
func TestDecodeHex(t *testing.T) {
	testCases := []TestCase{
		{"0", []byte{0}, nil},
		{"00", []byte{0}, nil},
		{"0x0", []byte{0}, nil},
		{"0x00", []byte{0}, nil},
		{"1", []byte{1}, nil},
		{"01", []byte{1}, nil},
		{"", []byte{}, hex.ErrLength},
		{"0x", []byte{}, hex.ErrLength},
		{"zz", []byte{}, hex.InvalidByteError('z')},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			output, err := DecodeHex(tc.input)
			if tc.err != nil {
				require.Error(t, tc.err, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, output, tc.output)
		})
	}
}
