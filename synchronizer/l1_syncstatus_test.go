package synchronizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FirstRunWithPendingBlocksToRetrieve(t *testing.T) {
	tcs := []struct {
		description           string
		lastStoredBlock       uint64
		lastL1Block           uint64
		chuncks               uint64
		expectedBlockRangeNil bool
		expectedBlockRange    blockRange
	}{
		{"normal", 100, 150, 10, false, blockRange{fromBlock: 101, toBlock: 111}},
		{"sync", 150, 150, 50, true, blockRange{}},
		{"less_chunk", 145, 150, 100, false, blockRange{fromBlock: 146, toBlock: 150}},
		{"1wide_range", 149, 150, 100, false, blockRange{fromBlock: 150, toBlock: 150}},
	}
	for _, tc := range tcs {
		s := newSyncStatus(tc.lastStoredBlock, tc.chuncks)
		s.setLastBlockOnL1(tc.lastL1Block)
		br := s.getNextRange()
		if tc.expectedBlockRangeNil {
			require.Nil(t, br, tc.description)
		} else {
			require.NotNil(t, br, tc.description)
			require.Equal(t, *br, tc.expectedBlockRange, tc.description)
		}
	}

}

func Test_SecondRunWithPendingBlocksToRetrieve(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.processedRanges.addBlockRange(blockRange{fromBlock: 101, toBlock: 111})
	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
}

func Test_generateNextRangeWithStoreResult(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.processedRanges.addBlockRange(blockRange{fromBlock: 101, toBlock: 111})
	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
	require.Equal(t, s.processedRanges.len(), 1)
}
