package tree

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestSplit(t *testing.T) {
	v1str := "200000000000000000000"

	v1, success := new(big.Int).SetString(v1str, 10)
	require.True(t, success)

	v, err := scalar2fea(v1)
	require.NoError(t, err)

	v2 := fea2scalar(v)
	assert.Equal(t, v1, v2)

	vv, err := scalar2fea(v2)
	require.NoError(t, err)
	assert.Equal(t, v, vv)
}
