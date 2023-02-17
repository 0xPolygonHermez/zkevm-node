package state_test

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stderr"},
	})
}

func TestDecodeRandomBatchL2Data(t *testing.T) {
	randomData := []byte("Random data")
	txs, _, err := state.DecodeTxs(randomData)
	require.Error(t, err)
	assert.Equal(t, []types.Transaction{}, txs)
	t.Log("Txs decoded 1: ", txs)

	randomData = []byte("Esto es autentica basura")
	txs, _, err = state.DecodeTxs(randomData)
	require.Error(t, err)
	assert.Equal(t, []types.Transaction{}, txs)
	t.Log("Txs decoded 2: ", txs)

	randomData = []byte("beef")
	txs, _, err = state.DecodeTxs(randomData)
	require.Error(t, err)
	assert.Equal(t, []types.Transaction{}, txs)
	t.Log("Txs decoded 3: ", txs)
}
