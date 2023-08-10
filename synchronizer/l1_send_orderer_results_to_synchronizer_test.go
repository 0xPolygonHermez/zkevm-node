package synchronizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SOR_TrivialCaseThatArrivesNextBlock(t *testing.T) {
	ch := make(chan getRollupInfoByBlockRangeResult, 100)
	lastBlock := uint64(100)
	sut := newSendResultsToSynchronizer(ch, lastBlock)
	result := getRollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 101,
			toBlock:   110,
		},
	}
	sut.addResultAndSendToConsumer(&result)
	require.Equal(t, 0, len(sut.pendingResults))
	require.Equal(t, uint64(110), sut.lastBlockOnSynchronizer)
}

func Test_SOR_ReceivedABlockThatIsNotNextOne(t *testing.T) {
	ch := make(chan getRollupInfoByBlockRangeResult, 100)
	lastBlock := uint64(100)
	sut := newSendResultsToSynchronizer(ch, lastBlock)
	result := getRollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 111,
			toBlock:   120,
		},
	}
	sut.addResultAndSendToConsumer(&result)
	require.Equal(t, 1, len(sut.pendingResults))
	require.Equal(t, lastBlock, sut.lastBlockOnSynchronizer)
}

func Test_SOR_ThereAreSomePendingBlocksAndArriveTheMissingOne(t *testing.T) {
	ch := make(chan getRollupInfoByBlockRangeResult, 100)
	lastBlock := uint64(100)
	sut := newSendResultsToSynchronizer(ch, lastBlock)
	sut.addResultAndSendToConsumer(&getRollupInfoByBlockRangeResult{blockRange: blockRange{fromBlock: 111, toBlock: 120}})
	sut.addResultAndSendToConsumer(&getRollupInfoByBlockRangeResult{blockRange: blockRange{fromBlock: 121, toBlock: 130}})
	sut.addResultAndSendToConsumer(&getRollupInfoByBlockRangeResult{blockRange: blockRange{fromBlock: 131, toBlock: 140}})
	require.Equal(t, 3, len(sut.pendingResults))
	require.Equal(t, lastBlock, sut.lastBlockOnSynchronizer)

	sut.addResultAndSendToConsumer(&getRollupInfoByBlockRangeResult{blockRange: blockRange{fromBlock: 101, toBlock: 110}})
	require.Equal(t, 0, len(sut.pendingResults))
	require.Equal(t, uint64(140), sut.lastBlockOnSynchronizer)
}
