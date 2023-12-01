package l1_parallel_sync

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGivenObjectWithDataWhenResetThenDontForgetLastBlockOnL1AndgetNextRangeReturnsNil(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.OnStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})

	s.Reset(1234)

	// lose lastBlockOnL1 so it returns a nil
	br := s.GetNextRange()
	require.Equal(t, *br, blockRange{fromBlock: 1235, toBlock: 1245})
}

func TestGivenObjectWithDataWhenResetAndSetLastBlockOnL1ThenGetNextRangeReturnsNextRange(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.OnStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})

	s.Reset(1234)
	s.setLastBlockOnL1(1982)
	// lose lastBlockOnL1 so it returns a nil
	br := s.GetNextRange()
	require.Equal(t, *br, blockRange{fromBlock: 1235, toBlock: 1245})
}

// Only could be 1 request to latest block
func TestGivenSychronizationWithThereAreARequestToLatestBlockWhenAskForNewBlockRangeItResponseNil(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.OnStartedNewWorker(blockRange{fromBlock: 1820, toBlock: latestBlockNumber})
	s.setLastBlockOnL1(1983)
	// Only could be 1 request to latest block
	br := s.GetNextRange()
	require.Nil(t, br)
	s.OnFinishWorker(blockRange{fromBlock: 1820, toBlock: latestBlockNumber}, true, uint64(1984))
	// We have a new segment to ask for because the last block have moved to 1984
	br = s.GetNextRange()
	require.Equal(t, blockRange{fromBlock: 1985, toBlock: latestBlockNumber}, *br)
}

func TestGivenSychronizationIAliveWhenWeAreInLatestBlockThenResponseNoNewBlockRange(t *testing.T) {
	s := newSyncStatus(1819, 10)
	s.setLastBlockOnL1(1823)
	br := s.GetNextRange()
	require.Equal(t, blockRange{fromBlock: 1820, toBlock: latestBlockNumber}, *br)
	s.OnStartedNewWorker(blockRange{fromBlock: 1820, toBlock: latestBlockNumber})
	s.setLastBlockOnL1(1824)
	// Only could be 1 request to latest block
	br = s.GetNextRange()
	require.Nil(t, br)
	s.OnFinishWorker(blockRange{fromBlock: 1820, toBlock: latestBlockNumber}, true, invalidBlockNumber)
	// We have a new segment to ask for because the last block have moved to 1984
	br = s.GetNextRange()
	require.Equal(t, blockRange{fromBlock: 1820, toBlock: latestBlockNumber}, *br)
	s.OnStartedNewWorker(blockRange{fromBlock: 1820, toBlock: latestBlockNumber})
	s.OnFinishWorker(blockRange{fromBlock: 1820, toBlock: latestBlockNumber}, true, 1830)
	// We have the latest block 1830, so we don't need to ask for something els until we update the last block on L1 (setLastBlockOnL1)
	br = s.GetNextRange()
	require.Nil(t, br)
}
func TestGivenThereAreALatestBlockErrorRangeIfMoveLastBlockBeyoundChunkThenDiscardErrorBR(t *testing.T) {
	s := newSyncStatus(1819, 10)
	s.setLastBlockOnL1(1823)
	br := s.GetNextRange()
	require.Equal(t, blockRange{fromBlock: 1820, toBlock: latestBlockNumber}, *br)
	s.OnStartedNewWorker(blockRange{fromBlock: 1820, toBlock: latestBlockNumber})
	s.setLastBlockOnL1(1824)
	// Only could be 1 request to latest block
	br = s.GetNextRange()
	require.Nil(t, br)
	s.OnFinishWorker(blockRange{fromBlock: 1820, toBlock: latestBlockNumber}, false, invalidBlockNumber)
	s.setLastBlockOnL1(1850)
	// We have a new segment to ask for because the last block have moved to 1984
	br = s.GetNextRange()
	require.Equal(t, blockRange{fromBlock: 1820, toBlock: 1830}, *br)
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
		{"less_chunk", 145, 150, 100, false, blockRange{fromBlock: 146, toBlock: latestBlockNumber}},
		{"1wide_range", 149, 150, 100, false, blockRange{fromBlock: 150, toBlock: latestBlockNumber}},
	}
	for _, tc := range tcs {
		s := newSyncStatus(tc.lastStoredBlock, tc.chuncks)
		s.setLastBlockOnL1(tc.lastL1Block)
		br := s.GetNextRange()
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
	res := s.OnFinishWorker(blockRange{fromBlock: 1618, toBlock: 1628}, true, uint64(1628))
	require.False(t, res)
	br := s.GetNextRange()
	require.Equal(t, blockRange{fromBlock: 1618, toBlock: 1628}, *br)
}

func TestWhenAllRequestAreSendThenGetNextRangeReturnsNil(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.OnStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})
	s.OnStartedNewWorker(blockRange{fromBlock: 1921, toBlock: 1982})
	br := s.GetNextRange()
	require.Nil(t, br)
}

func TestSecondRunWithPendingBlocksToRetrieve(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.OnStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
}

func TestGenerateNextRangeWithPreviousResult(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.OnStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
	require.Equal(t, s.processingRanges.len(), 1)
}

func TestGenerateNextRangeWithProcessedResult(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.OnStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	res := s.OnFinishWorker(blockRange{fromBlock: 101, toBlock: 111}, true, uint64(111))
	require.True(t, res)
	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
	require.Equal(t, s.processingRanges.len(), 0)
}

func TestGivenMultiplesWorkersWhenBrInMiddleFinishThenDontChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	//previousValue := s.lastBlockStoreOnStateDB
	s.OnStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.OnStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.OnStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	res := s.OnFinishWorker(blockRange{fromBlock: 112, toBlock: 122}, true, uint64(122))
	require.True(t, res)
	//require.Equal(t, previousValue, s.lastBlockStoreOnStateDB)

	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, blockRange{fromBlock: 134, toBlock: 144}, *br)
}

func TestGivenMultiplesWorkersWhenFirstFinishThenChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.OnStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.OnStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.OnStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	res := s.OnFinishWorker(blockRange{fromBlock: 101, toBlock: 111}, true, uint64(111))
	require.True(t, res)
	require.Equal(t, uint64(111), s.lastBlockStoreOnStateDB)

	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 134, toBlock: 144})
}

func TestGivenMultiplesWorkersWhenLastFinishThenDontChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	//previousValue := s.lastBlockStoreOnStateDB
	s.OnStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.OnStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.OnStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	res := s.OnFinishWorker(blockRange{fromBlock: 123, toBlock: 133}, true, uint64(133))
	require.True(t, res)
	//require.Equal(t, previousValue, s.lastBlockStoreOnStateDB)

	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, blockRange{fromBlock: 134, toBlock: 144}, *br)
}

func TestGivenMultiplesWorkersWhenLastFinishAndFinishAlsoNextOneThenDontChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(200)
	//previousValue := s.lastBlockStoreOnStateDB
	s.OnStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.OnStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.OnStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	res := s.OnFinishWorker(blockRange{fromBlock: 123, toBlock: 133}, true, uint64(133))
	require.True(t, res)
	s.OnStartedNewWorker(blockRange{fromBlock: 134, toBlock: 144})
	//require.Equal(t, previousValue, s.lastBlockStoreOnStateDB)

	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 145, toBlock: 155})
}

func TestGivenMultiplesWorkersWhenNextRangeThenTheRangeIsCappedToLastBlockOnL1(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(105)

	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 101, toBlock: latestBlockNumber})
}

func TestWhenRequestALatestBlockThereIsNoMoreBlocks(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(105)

	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 101, toBlock: latestBlockNumber})

	s.OnStartedNewWorker(*br)
	br = s.GetNextRange()
	require.Nil(t, br)
}

func TestWhenFinishALatestBlockIfNoNewLastBlockOnL1NothingToDo(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(105)

	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, blockRange{fromBlock: 101, toBlock: latestBlockNumber}, *br)

	s.OnStartedNewWorker(*br)
	noBR := s.GetNextRange()
	require.Nil(t, noBR)

	s.OnFinishWorker(*br, true, uint64(105))
	br = s.GetNextRange()
	require.Nil(t, br)
}

func TestWhenFinishALatestBlockIfThereAreNewLastBlockOnL1ThenThereIsANewRange(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(105)

	br := s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 101, toBlock: latestBlockNumber})

	s.OnStartedNewWorker(*br)
	noBR := s.GetNextRange()
	require.Nil(t, noBR)

	s.setLastBlockOnL1(106)
	s.OnFinishWorker(*br, true, invalidBlockNumber) // No block info in the answer
	br = s.GetNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 101, toBlock: latestBlockNumber})
}
