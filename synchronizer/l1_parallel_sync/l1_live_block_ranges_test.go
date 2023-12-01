package l1_parallel_sync

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsertBR(t *testing.T) {
	sut := newLiveBlockRanges()
	err := sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	require.NoError(t, err)
	require.Equal(t, sut.len(), 1)
}
func TestInsertOverlappedBR(t *testing.T) {
	sut := newLiveBlockRanges()
	err := sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	require.NoError(t, err)
	err = sut.addBlockRange(blockRange{fromBlock: 5, toBlock: 15})
	require.Error(t, err)
	require.Equal(t, sut.len(), 1)
}

func TestInsertDuplicatedBR(t *testing.T) {
	sut := newLiveBlockRanges()
	err := sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	require.NoError(t, err)
	err = sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	require.Error(t, err)
	require.Equal(t, sut.len(), 1)
}

func TestRemoveExistingBR(t *testing.T) {
	sut := newLiveBlockRanges()
	err := sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	require.NoError(t, err)
	err = sut.addBlockRange(blockRange{fromBlock: 11, toBlock: 20})
	require.NoError(t, err)
	err = sut.removeBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	require.NoError(t, err)
	require.Equal(t, sut.len(), 1)
}

func TestInsertWrongBR1(t *testing.T) {
	sut := newLiveBlockRanges()
	err := sut.addBlockRange(blockRange{})
	require.Error(t, err)
	require.Equal(t, sut.len(), 0)
}
func TestInsertWrongBR2(t *testing.T) {
	sut := newLiveBlockRanges()
	err := sut.addBlockRange(blockRange{fromBlock: 10, toBlock: 5})
	require.Error(t, err)
	require.Equal(t, sut.len(), 0)
}

func TestGetSuperBlockRangeEmpty(t *testing.T) {
	sut := newLiveBlockRanges()
	res := sut.getSuperBlockRange()
	require.Nil(t, res)
}

func TestGetSuperBlockRangeWithData(t *testing.T) {
	sut := newLiveBlockRanges()
	err := sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	require.NoError(t, err)
	err = sut.addBlockRange(blockRange{fromBlock: 11, toBlock: 20})
	require.NoError(t, err)
	err = sut.addBlockRange(blockRange{fromBlock: 21, toBlock: 109})
	require.NoError(t, err)
	err = sut.addBlockRange(blockRange{fromBlock: 110, toBlock: 200})
	require.NoError(t, err)

	res := sut.getSuperBlockRange()
	require.NotNil(t, res)
	require.Equal(t, *res, blockRange{fromBlock: 1, toBlock: 200})
}
