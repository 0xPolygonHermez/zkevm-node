package synchronizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Insert_BR(t *testing.T) {
	sut := NewLiveBlockRanges()
	sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	require.Equal(t, sut.len(), 1)
}
func Test_Insert_Overlapped_BR(t *testing.T) {
	sut := NewLiveBlockRanges()
	sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	err := sut.addBlockRange(blockRange{fromBlock: 5, toBlock: 15})
	require.Error(t, err)
	require.Equal(t, sut.len(), 1)
}
func Test_Insert_Duplicated_BR(t *testing.T) {
	sut := NewLiveBlockRanges()
	sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	err := sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	require.Error(t, err)
	require.Equal(t, sut.len(), 1)
}

func Test_Insert_NoConsecutiveBlock_BR(t *testing.T) {
	sut := NewLiveBlockRanges()
	sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	err := sut.addBlockRange(blockRange{fromBlock: 12, toBlock: 20})
	require.Error(t, err)
	require.Equal(t, sut.len(), 1)
}

func Test_Remove_Existing_BR(t *testing.T) {
	sut := NewLiveBlockRanges()
	sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	sut.addBlockRange(blockRange{fromBlock: 11, toBlock: 20})
	err := sut.removeBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	require.NoError(t, err)
	require.Equal(t, sut.len(), 1)
}

func Test_Insert_Wrong_BR1(t *testing.T) {
	sut := NewLiveBlockRanges()
	err := sut.addBlockRange(blockRange{})
	require.Error(t, err)
	require.Equal(t, sut.len(), 0)
}
func Test_Insert_Wrong_BR2(t *testing.T) {
	sut := NewLiveBlockRanges()
	err := sut.addBlockRange(blockRange{fromBlock: 10, toBlock: 5})
	require.Error(t, err)
	require.Equal(t, sut.len(), 0)
}

func Test_GetSuperBlockRange_Emtpy(t *testing.T) {
	sut := NewLiveBlockRanges()
	res := sut.GetSuperBlockRange()
	require.Nil(t, res)
}

func Test_GetSuperBlockRange_WithData(t *testing.T) {
	sut := NewLiveBlockRanges()
	sut.addBlockRange(blockRange{fromBlock: 1, toBlock: 10})
	sut.addBlockRange(blockRange{fromBlock: 11, toBlock: 20})
	sut.addBlockRange(blockRange{fromBlock: 21, toBlock: 109})
	sut.addBlockRange(blockRange{fromBlock: 110, toBlock: 200})

	res := sut.GetSuperBlockRange()
	require.NotNil(t, res)
	require.Equal(t, *res, blockRange{fromBlock: 1, toBlock: 200})
}
