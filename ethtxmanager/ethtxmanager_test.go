package ethtxmanager

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncreaseGasPrice(t *testing.T) {
	actual := increaseGasPrice(big.NewInt(100), 1)
	assert.Equal(t, big.NewInt(101), actual)
}

func TestIncreaseGasLimit(t *testing.T) {
	actual := increaseGasLimit(100, 1)
	assert.Equal(t, uint64(101), actual)
}
