package tree

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {
	v1str := "115792089237316195423570985008687907853269984665640564039457584007913129639935"

	v1, ok := new(big.Int).SetString(v1str, 10)
	require.True(t, ok)

	v := scalar2fea(v1)

	v2 := fea2scalar(v)
	require.Equal(t, v1, v2)

	vv := scalar2fea(v2)

	require.Equal(t, v, vv)
}

func Test_h4ToScalar(t *testing.T) {
	tcs := []struct {
		input    []uint64
		expected string
	}{
		{
			input:    []uint64{0, 0, 0, 0},
			expected: "0",
		},
		{
			input:    []uint64{0, 1, 2, 3},
			expected: "18831305206160042292187933003464876175252262292329349513216",
		},
	}

	for i, tc := range tcs {
		tc := tc
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actual := h4ToScalar(tc.input)
			expected, ok := new(big.Int).SetString(tc.expected, 10)
			require.True(t, ok)
			require.Equal(t, expected, actual)
		})
	}
}

func Test_scalarToh4(t *testing.T) {
	tcs := []struct {
		input    string
		expected []uint64
	}{
		{
			input:    "0",
			expected: []uint64{0, 0, 0, 0},
		},
		{
			input:    "18831305206160042292187933003464876175252262292329349513216",
			expected: []uint64{0, 1, 2, 3},
		},
	}

	for i, tc := range tcs {
		tc := tc
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			bi, ok := new(big.Int).SetString(tc.input, 10)
			require.True(t, ok)

			actual := scalarToh4(bi)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func Test_h4ToString(t *testing.T) {
	tcs := []struct {
		input    []uint64
		expected string
	}{
		{
			input:    []uint64{0, 0, 0, 0},
			expected: "0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			input:    []uint64{0, 1, 2, 3},
			expected: "0x0000000000000003000000000000000200000000000000010000000000000000",
		},
	}

	for i, tc := range tcs {
		tc := tc
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			actual := h4ToString(tc.input)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func Test_Conversions(t *testing.T) {
	tcs := []struct {
		input []uint64
	}{
		{
			input: []uint64{0, 0, 0, 0},
		},
		{
			input: []uint64{0, 1, 2, 3},
		},
	}

	for i, tc := range tcs {
		tc := tc
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			resScalar := h4ToScalar(tc.input)
			init := scalarToh4(resScalar)
			require.Equal(t, tc.input, init)
		})
	}
}

func Test_scalar2fea(t *testing.T) {
	tcs := []struct {
		input string
	}{
		{
			input: "0",
		},
		{
			input: "100",
		},
		{
			input: "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		},
	}

	for i, tc := range tcs {
		tc := tc
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			scalar, ok := new(big.Int).SetString(tc.input, 10)
			require.True(t, ok)

			res := scalar2fea(scalar)

			actual := fea2scalar(res)
			require.Equal(t, tc.input, actual.String())
		})
	}
}

func Test_fea2scalar(t *testing.T) {
	tcs := []struct {
		input []uint64
	}{
		{
			input: []uint64{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			input: []uint64{1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			input: []uint64{1, 0, 0, 0, 3693650181, 4001443757, 599269951, 1255793162},
		},
	}

	for i, tc := range tcs {
		tc := tc
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			res := fea2scalar(tc.input)

			actual := scalar2fea(res)

			require.Equal(t, tc.input, actual)
		})
	}
}
