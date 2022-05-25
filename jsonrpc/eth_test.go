package jsonrpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockNumber(t *testing.T) {
	const expectedBlockNumber = uint64(10)

	server, ethClient := newMockedServer(t)
	defer server.Server.Stop()

	server.State.
		On("GetLastBatchNumber", context.Background(), "").
		Return(expectedBlockNumber, nil)

	bn, err := ethClient.BlockNumber(context.Background())
	require.NoError(t, err)

	assert.Equal(t, expectedBlockNumber, bn)
}

func TestChainID(t *testing.T) {
	server, ethClient := newMockedServer(t)
	defer server.Server.Stop()

	chainID, err := ethClient.ChainID(context.Background())
	require.NoError(t, err)

	assert.Equal(t, server.ChainID, chainID.Uint64())
}
