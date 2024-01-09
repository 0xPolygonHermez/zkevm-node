package l1_parallel_sync

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/stretchr/testify/require"
)

func TestProducerStatisticsPercent(t *testing.T) {
	sut := newRollupInfoProducerStatistics(100, &common.MockTimerProvider{})
	sut.updateLastBlockNumber(200)
	require.Equal(t, float64(0.0), sut.getPercent())

	sut.onResponseRollupInfo(responseRollupInfoByBlockRange{
		generic: genericResponse{
			err:      nil,
			duration: 0,
		},
		result: &rollupInfoByBlockRangeResult{
			blockRange: blockRange{
				fromBlock: 101,
				toBlock:   200,
			},
		},
	})

	require.Equal(t, float64(100.0), sut.getPercent())

	sut.reset(100)
	require.Equal(t, float64(0.0), sut.getPercent())
}
