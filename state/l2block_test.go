package state

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestL2BlockHash(t *testing.T) {
	// create a geth header and block
	header := &types.Header{Number: big.NewInt(1)}
	ethBlock := types.NewBlockWithHeader(header)

	// create a l2 header and l2 block from geth header
	l2Header := NewL2Header(header)
	l2Block := NewL2BlockWithHeader(l2Header)

	// compare geth and l2 block hashes, they must match
	assert.Equal(t, ethBlock.Hash().String(), l2Block.Hash().String())
}
