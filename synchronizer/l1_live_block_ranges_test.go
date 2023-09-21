package synchronizer

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

func TestGetTagByBlockRangeThatExists(t *testing.T) {
	sut := newLiveBlockRangesWithTag[string]()
	err := sut.addBlockRangeWithTag(blockRange{fromBlock: 1, toBlock: 10}, "a")
	require.NoError(t, err)
	brForTagB := blockRange{fromBlock: 11, toBlock: 20}
	err = sut.addBlockRangeWithTag(brForTagB, "b")
	require.NoError(t, err)
	err = sut.addBlockRangeWithTag(blockRange{fromBlock: 21, toBlock: 109}, "c")
	require.NoError(t, err)
	err = sut.addBlockRangeWithTag(blockRange{fromBlock: 110, toBlock: 200}, "d")
	require.NoError(t, err)

	tag, err := sut.getTagByBlockRange(brForTagB)
	require.NoError(t, err)
	require.Equal(t, "b", tag)
}

func TestGetTagByBlockRangeThatDontExists(t *testing.T) {
	sut := newLiveBlockRangesWithTag[string]()
	err := sut.addBlockRangeWithTag(blockRange{fromBlock: 1, toBlock: 10}, "a")
	require.NoError(t, err)
	brForTagB := blockRange{fromBlock: 11, toBlock: 20}
	err = sut.addBlockRangeWithTag(brForTagB, "b")
	require.NoError(t, err)
	err = sut.addBlockRangeWithTag(blockRange{fromBlock: 21, toBlock: 109}, "c")
	require.NoError(t, err)
	err = sut.addBlockRangeWithTag(blockRange{fromBlock: 110, toBlock: 200}, "d")
	require.NoError(t, err)

	_, err = sut.getTagByBlockRange(blockRange{fromBlock: 12210, toBlock: 22200})
	require.Error(t, err)
}

func TestFilterTagByBlockRangeThatDontExists(t *testing.T) {
	sut := newLiveBlockRangesWithTag[int]()
	err := sut.addBlockRangeWithTag(blockRange{fromBlock: 1, toBlock: 10}, 10)
	require.NoError(t, err)

	err = sut.addBlockRangeWithTag(blockRange{fromBlock: 11, toBlock: 20}, 20)
	require.NoError(t, err)
	brForTag1 := blockRange{fromBlock: 1100, toBlock: 2000}
	err = sut.addBlockRangeWithTag(brForTag1, 30)
	require.NoError(t, err)
	brForTag2 := blockRange{fromBlock: 410, toBlock: 500}
	err = sut.addBlockRangeWithTag(brForTag2, 30)
	require.NoError(t, err)

	res := sut.filterBlockRangesByTag(func(br blockRange, tag int) bool { return tag >= 30 })
	require.Equal(t, len(res), 2)
	require.Equal(t, []blockRange{brForTag1, brForTag2}, res)
}
