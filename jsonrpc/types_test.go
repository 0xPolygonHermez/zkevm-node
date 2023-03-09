package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgHashUnmarshalFromShortString(t *testing.T) {
	str := "0x01"
	arg := argHash{}
	err := arg.UnmarshalText([]byte(str))
	require.NoError(t, err)

	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", arg.Hash().String())
}

func TestArgAddressUnmarshalFromShortString(t *testing.T) {
	str := "0x01"
	arg := argAddress{}
	err := arg.UnmarshalText([]byte(str))
	require.NoError(t, err)

	assert.Equal(t, "0x0000000000000000000000000000000000000001", arg.Address().String())
}
