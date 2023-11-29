package state

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestForceL2BlockHash(t *testing.T) {
	block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
	hashToForce := common.HexToHash("0x111222333")
	assert.NotEqual(t, hashToForce.String(), block.Hash().String())
	ForceL2BlockHash(block, hashToForce)
	assert.Equal(t, hashToForce.String(), block.Hash().String())
}
