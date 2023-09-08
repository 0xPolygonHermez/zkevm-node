package synchronizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_When_Reset_Then_Forget_LastBlockOnL1_getNextRange_ReturnsNil(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})

	s.reset(1234)

	// lose lastBlockOnL1 so it returns a nil
	br := s.getNextRange()
	require.Nil(t, br)
}

func Test_When_Reset_Then_If_Set_LastBlockOnL1_getNextRange_ReturnsNextRange(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})

	s.reset(1234)
	s.setLastBlockOnL1(1982)
	// lose lastBlockOnL1 so it returns a nil
	br := s.getNextRange()
	require.Equal(t, *br, blockRange{fromBlock: 1235, toBlock: 1245})
}

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

func Test_When_ReceiveAndNoStartedBlockRange_Then_Ignore(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.onFinishWorker(blockRange{fromBlock: 1618, toBlock: 1628}, true)
	br := s.getNextRange()
	require.Equal(t, blockRange{fromBlock: 1618, toBlock: 1628}, *br)
}

func Test_When_AllRequestAreSend_Then_getNextRange_ReturnsNil(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})
	s.onStartedNewWorker(blockRange{fromBlock: 1921, toBlock: 1982})
	br := s.getNextRange()
	require.Nil(t, br)
}

func Test_SecondRunWithPendingBlocksToRetrieve(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
}

func Test_generateNextRangeWithPreviousResult(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
	require.Equal(t, s.processingRanges.len(), 1)
}

func Test_generateNextRangeWithProcessedResult(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.onFinishWorker(blockRange{fromBlock: 101, toBlock: 111}, true)
	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 112, toBlock: 122})
	require.Equal(t, s.processingRanges.len(), 0)
}

func Test_Given_MultiplesWorkers_When_BrInMiddleFinish_Then_DontChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	previousValue := s.lastBlockStoreOnStateDB
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.onStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.onStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	s.onFinishWorker(blockRange{fromBlock: 112, toBlock: 122}, true)
	require.Equal(t, previousValue, s.lastBlockStoreOnStateDB)

	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 134, toBlock: 144})
}

func Test_Given_MultiplesWorkers_When_FirstFinish_Then_ChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)

	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.onStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.onStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	s.onFinishWorker(blockRange{fromBlock: 101, toBlock: 111}, true)
	require.Equal(t, uint64(111), s.lastBlockStoreOnStateDB)

	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 134, toBlock: 144})
}

func Test_Given_MultiplesWorkers_When_LastFinish_Then_DontChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(150)
	previousValue := s.lastBlockStoreOnStateDB
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.onStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.onStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	s.onFinishWorker(blockRange{fromBlock: 123, toBlock: 133}, true)
	require.Equal(t, previousValue, s.lastBlockStoreOnStateDB)

	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 134, toBlock: 144})
}

func Test_Given_MultiplesWorkers_When_LastFinishAndFinishAlsoNextOne_Then_DontChangeLastBlock(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(200)
	previousValue := s.lastBlockStoreOnStateDB
	s.onStartedNewWorker(blockRange{fromBlock: 101, toBlock: 111})
	s.onStartedNewWorker(blockRange{fromBlock: 112, toBlock: 122})
	s.onStartedNewWorker(blockRange{fromBlock: 123, toBlock: 133})
	s.onFinishWorker(blockRange{fromBlock: 123, toBlock: 133}, true)
	s.onStartedNewWorker(blockRange{fromBlock: 134, toBlock: 144})
	require.Equal(t, previousValue, s.lastBlockStoreOnStateDB)

	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 145, toBlock: 155})
}

func Test_Given_MultiplesWorkers_When_NextRange_Then_TheRangeIsCappedToLastBlockOnL1(t *testing.T) {
	s := newSyncStatus(100, 10)
	s.setLastBlockOnL1(105)

	br := s.getNextRange()
	require.NotNil(t, br)
	require.Equal(t, *br, blockRange{fromBlock: 101, toBlock: 105})
}

// Test status

func Test_When_Create_Then_status_is_idle(t *testing.T) {
	s := newSyncStatus(1617, 10)
	require.Equal(t, syncStatusIdle, s.getStatus())
}

func Test_When_SetLastBlockOnL1GreaterCurrentBlock_Then_status_is_working(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(2000)
	require.Equal(t, syncStatusWorking, s.getStatus())
}

func Test_When_SetLastBlockOnL1BelowCurrentBlock_Then_status_is_working(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(100)
	require.Equal(t, syncStatusSynchronized, s.getStatus())
}

func Test_When_StartedWorker_Then_status_is_working(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})
	require.Equal(t, syncStatusWorking, s.getStatus())
}

func Test_When_EndWorkerIfNoEqualLastBlockL1_Then_status_is_working(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(2000)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})
	s.onFinishWorker(blockRange{fromBlock: 1820, toBlock: 1920}, true)
	require.Equal(t, syncStatusWorking, s.getStatus())
}

func Test_When_EndWorkerEqualToLastBlockL1_Then_status_is_sync(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(2000)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 2000})
	s.onFinishWorker(blockRange{fromBlock: 1820, toBlock: 2000}, true)
	require.Equal(t, syncStatusSynchronized, s.getStatus())
}

func Test_When_EndWorkerGreaterToLastBlockL1_Then_status_is_sync(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(2000)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 2001})
	s.onFinishWorker(blockRange{fromBlock: 1820, toBlock: 2001}, true)
	require.Equal(t, syncStatusSynchronized, s.getStatus())
}

func Test_When_AllRequestAreSend_Then_getNextRange_ReturnsNil2(t *testing.T) {
	s := newSyncStatus(1617, 10)
	s.setLastBlockOnL1(1982)
	s.onStartedNewWorker(blockRange{fromBlock: 1820, toBlock: 1920})
	s.onStartedNewWorker(blockRange{fromBlock: 1921, toBlock: 1982})
	br := s.getNextRange()
	require.Nil(t, br)
}
