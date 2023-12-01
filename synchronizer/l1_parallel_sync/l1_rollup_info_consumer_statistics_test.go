package l1_parallel_sync

import (
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/stretchr/testify/assert"
)

func TestL1RollupInfoConsumerStatistics(t *testing.T) {
	cfg := ConfigConsumer{
		ApplyAfterNumRollupReceived: 10,
		AceptableInacctivityTime:    5 * time.Second,
	}
	stats := l1RollupInfoConsumerStatistics{
		cfg: cfg,
	}

	stats.onStart()
	stats.onStartStep()

	// Test onFinishProcessIncommingRollupInfoData
	rollupInfo := rollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 1,
			toBlock:   10,
		},
		blocks: []etherman.Block{},
	}
	executionTime := 2 * time.Second
	stats.onStartProcessIncommingRollupInfoData(rollupInfo)
	stats.onFinishProcessIncommingRollupInfoData(rollupInfo, executionTime, error(nil))
	assert.Equal(t, stats.timePreviousProcessingDuration, executionTime)
	assert.Equal(t, stats.numProcessedRollupInfo, uint64(1))
	assert.Equal(t, stats.numProcessedBlocks, uint64(len(rollupInfo.blocks)))

	stats.onStart()
	stats.onStartStep()

	msg := stats.onStartProcessIncommingRollupInfoData(rollupInfo)
	assert.Contains(t, msg, "wasted_time_waiting_for_data")
	assert.Contains(t, msg, "last_process_time")
	assert.Contains(t, msg, "block_per_second")
	assert.NotContains(t, msg, "WASTED_TIME_EXCEED")
	assert.NotContains(t, msg, "WARNING_WASTED_TIME")
}

func TestL1RollupInfoConsumerStatisticsWithExceedTimeButNoWarningGenerated(t *testing.T) {
	cfg := ConfigConsumer{
		ApplyAfterNumRollupReceived: 10,
		AceptableInacctivityTime:    0 * time.Second,
	}
	stats := l1RollupInfoConsumerStatistics{
		cfg: cfg,
	}

	stats.onStart()
	stats.onStartStep()

	// Test onFinishProcessIncommingRollupInfoData
	rollupInfo := rollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 1,
			toBlock:   10,
		},
		blocks: []etherman.Block{},
	}
	executionTime := 2 * time.Second
	err := error(nil)
	stats.onStartProcessIncommingRollupInfoData(rollupInfo)
	stats.onFinishProcessIncommingRollupInfoData(rollupInfo, executionTime, err)

	stats.onStartStep()
	msg := stats.onStartProcessIncommingRollupInfoData(rollupInfo)
	assert.Contains(t, msg, "wasted_time_waiting_for_data")
	assert.Contains(t, msg, "last_process_time")
	assert.Contains(t, msg, "block_per_second")
	assert.Contains(t, msg, "WASTED_TIME_EXCEED")
	assert.NotContains(t, msg, "WARNING_WASTED_TIME")
}

func TestL1RollupInfoConsumerStatisticsWithExceedTimeButAndWarningGenerated(t *testing.T) {
	cfg := ConfigConsumer{
		ApplyAfterNumRollupReceived: 1,
		AceptableInacctivityTime:    0 * time.Second,
	}
	stats := l1RollupInfoConsumerStatistics{
		cfg: cfg,
	}

	stats.onStart()
	stats.onStartStep()

	// Test onFinishProcessIncommingRollupInfoData
	rollupInfo := rollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 1,
			toBlock:   10,
		},
		blocks: []etherman.Block{},
	}
	executionTime := 2 * time.Second
	err := error(nil)
	stats.onStartProcessIncommingRollupInfoData(rollupInfo)
	stats.onFinishProcessIncommingRollupInfoData(rollupInfo, executionTime, err)
	stats.onStartProcessIncommingRollupInfoData(rollupInfo)
	stats.onFinishProcessIncommingRollupInfoData(rollupInfo, executionTime, err)

	stats.onStartStep()
	msg := stats.onStartProcessIncommingRollupInfoData(rollupInfo)
	assert.Contains(t, msg, "wasted_time_waiting_for_data")
	assert.Contains(t, msg, "last_process_time")
	assert.Contains(t, msg, "block_per_second")
	assert.Contains(t, msg, "WASTED_TIME_EXCEED")
	assert.Contains(t, msg, "WARNING_WASTED_TIME")
}
