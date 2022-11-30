package ethtxmanager

import (
	"context"
	"math/big"
	"testing"

	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/stretchr/testify/assert"
)

func TestIncreaseGasLimit(t *testing.T) {
	actual := increaseGasLimit(100, 1)
	assert.Equal(t, uint64(101), actual)
}

func TestIncreaseGasPrice(t *testing.T) {
	actual := increaseGasPrice(big.NewInt(100), 1)
	assert.Equal(t, big.NewInt(101), actual)
}

func TestSequenceBatchesWithROEthman(t *testing.T) {
	ethManRO, _, _, _, _ := ethman.NewSimulatedEtherman(ethman.Config{}, nil)
	txMan := New(Config{MaxSendBatchTxRetries: 2}, ethManRO, nil) // 3 executions in total

	err := txMan.SequenceBatches(context.Background(), []ethmanTypes.Sequence{})

	assert.ErrorIs(t, err, ethman.ErrIsReadOnlyMode)
}

func TestVerifyBatchWithROEthman(t *testing.T) {
	ethManRO, _, _, _, _ := ethman.NewSimulatedEtherman(ethman.Config{}, nil)
	txMan := New(Config{MaxVerifyBatchTxRetries: 2}, ethManRO, nil) // 3 executions in total

	err := txMan.VerifyBatch(context.Background(), 42, nil)

	assert.ErrorIs(t, err, ethman.ErrIsReadOnlyMode)
}
