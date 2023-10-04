package synchronizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGivenObjectWithDataWhenResetThenDontForgetLastBlockOnL1AndgetNextRangeReturnsNil(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})

	s.reset(1234)

	// lose lastBlockOnL1 so it returns a nil
	br := s.getNextRange()
	require.Equal(t, *br, blockRange{fromBlock: 1235, toBlock: 1245})
}

func TestGivenObjectWithDataWhenResetAndSetLastBlockOnL1ThenGetNextRangeReturnsNextRange(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})

	s.reset(1234)
	s.setLastBlockOnL1(1982)
	// lose lastBlockOnL1 so it returns a nil
	br := s.getNextRange()
	require.Equal(t, *br, blockRange{fromBlock: 1235, toBlock: 1245})
}

func TestFirstRunWithPendingBlocksToRetrieve(t *testing.T) {
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

func TestWhenReceiveAndNoStartedBlockRangeThenIgnore(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	res := s.onFinishWorker(blockRange{fromBlock: 1618, toBlock: 1628}, true)
	require.False(t, res)
	br := s.getNextRange()
	require.Equal(t, blockRange{fromBlock: 1618, toBlock: 1628}, *br)
}

func TestWhenAllRequestAreSendThenGetNextRangeReturnsNil(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})
	s.onStartedNewWorker(blockRange{fromBlock: 1921, toBlock: 1982})
	br := s.getNextRange()
	require.Nil(t, br)
}

func TestSecondRunWithPendingBlocksToRetrieve(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
}

func TestGenerateNextRangeWithPreviousResult(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
	require.Equal(t, s.processingRanges.len(), 1)
}

func TestGenerateNextRangeWithProcessedResult(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	res := s.onFinishWorker(blockRange{fromBlock: 101, toBlock: 111}, true)
	require.True(t, res)
	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
	require.Equal(t, s.processingRanges.len(), 0)
}

func TestGivenMultiplesWorkersWhenBrInMiddleFinishThenDontChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	previousValue := s.lastBlockStoreOnStateDB
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.onStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.onStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	res := s.onFinishWorker(blockRange{fromBlock: 112, toBlock: 122}, true)
	require.True(t, res)
	require.Equal(t, previousValue, s.lastBlockStoreOnStateDB)

	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 134, toBlock: 144})
}

func TestGivenMultiplesWorkersWhenFirstFinishThenChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.onStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.onStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	res := s.onFinishWorker(blockRange{fromBlock: 101, toBlock: 111}, true)
	require.True(t, res)
	require.Equal(t, uint64(111), s.lastBlockStoreOnStateDB)

	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 134, toBlock: 144})
}

func TestGivenMultiplesWorkersWhenLastFinishThenDontChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	previousValue := s.lastBlockStoreOnStateDB
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.onStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.onStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	res := s.onFinishWorker(blockRange{fromBlock: 123, toBlock: 133}, true)
	require.True(t, res)
	require.Equal(t, previousValue, s.lastBlockStoreOnStateDB)

	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 134, toBlock: 144})
}

func TestGivenMultiplesWorkersWhenLastFinishAndFinishAlsoNextOneThenDontChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(200)
	previousValue := s.lastBlockStoreOnStateDB
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.onStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.onStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	res := s.onFinishWorker(blockRange{fromBlock: 123, toBlock: 133}, true)
	require.True(t, res)
	s.onStartedNewWorker(blockRange{fromBlock: 134, toBlock: 144})
	require.Equal(t, previousValue, s.lastBlockStoreOnStateDB)

	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 145, toBlock: 155})
}

func TestGivenMultiplesWorkersWhenNextRangeThenTheRangeIsCappedToLastBlockOnL1(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(105)

	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 101, toBlock: 105})
}

func TestWhenAllRequestAreSendThenGetNextRangeReturnsNil2(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})
	s.onStartedNewWorker(blockRange{fromBlock: 1921, toBlock: 1982})
	br := s.getNextRange()
	require.Nil(t, br)
}
