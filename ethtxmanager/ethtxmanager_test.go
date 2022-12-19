package ethtxmanager

import (
	"context"
	"testing"

	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/stretchr/testify/assert"
)

func TestSequenceBatchesWithROEthman(t *testing.T) {
	ethManRO, _, _, _, _ := ethman.NewSimulatedEtherman(ethman.Config{}, nil)
	txMan := New(Config{}, ethManRO, nil)

	err := txMan.SequenceBatches(context.Background(), []ethmanTypes.Sequence{})

	assert.ErrorIs(t, err, ethman.ErrIsReadOnlyMode)
}

func TestVerifyBatchesWithROEthman(t *testing.T) {
	ethManRO, _, _, _, _ := ethman.NewSimulatedEtherman(ethman.Config{}, nil)
	txMan := New(Config{}, ethManRO, nil)

	_, err := txMan.VerifyBatches(context.Background(), 41, 42, nil)

	assert.ErrorIs(t, err, ethman.ErrIsReadOnlyMode)
}
