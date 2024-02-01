package l1_parallel_sync

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkerDecoratorLimitRetriesByTime_asyncRequestRollupInfoByBlockRange(t *testing.T) {
	// Create a new worker decorator with a minimum time between calls of 1 second
	workersMock := newWorkersInterfaceMock(t)
	decorator := newWorkerDecoratorLimitRetriesByTime(workersMock, time.Second)

	// Create a block range to use for testing
	blockRange := blockRange{1, 10}

	// Test the case where there is no previous call to the block range
	ctx := context.Background()
	workersMock.On("asyncRequestRollupInfoByBlockRange", ctx, requestRollupInfoByBlockRange{blockRange: blockRange, sleepBefore: noSleepTime, requestLastBlockIfNoBlocksInAnswer: requestLastBlockModeIfNoBlocksInAnswer}).Return(nil, nil).Once()
	_, err := decorator.asyncRequestRollupInfoByBlockRange(ctx, newRequestNoSleep(blockRange))
	assert.NoError(t, err)

	// Test the case where there is a previous call to the block range
	workersMock.On("asyncRequestRollupInfoByBlockRange", ctx, mock.MatchedBy(func(req requestRollupInfoByBlockRange) bool { return req.sleepBefore > 0 })).Return(nil, nil).Once()
	_, err = decorator.asyncRequestRollupInfoByBlockRange(ctx, newRequestNoSleep(blockRange))
	assert.NoError(t, err)
}

func TestWorkerDecoratorLimitRetriesByTimeIfRealWorkerReturnsAllBusyDoesntCountAsRetry(t *testing.T) {
	// Create a new worker decorator with a minimum time between calls of 1 second
	workersMock := newWorkersInterfaceMock(t)
	decorator := newWorkerDecoratorLimitRetriesByTime(workersMock, time.Second)

	// Create a block range to use for testing
	blockRange := blockRange{1, 10}

	// Test the case where there is no previous call to the block range
	ctx := context.Background()
	workersMock.On("asyncRequestRollupInfoByBlockRange", ctx, requestRollupInfoByBlockRange{blockRange: blockRange, sleepBefore: noSleepTime, requestLastBlockIfNoBlocksInAnswer: requestLastBlockModeIfNoBlocksInAnswer}).
		Return(nil, errAllWorkersBusy).
		Once()
	_, err := decorator.asyncRequestRollupInfoByBlockRange(ctx, newRequestNoSleep(blockRange))
	assert.Error(t, err)

	// Test the case where there is a previous call to the block range
	workersMock.On("asyncRequestRollupInfoByBlockRange", ctx, mock.MatchedBy(func(req requestRollupInfoByBlockRange) bool { return req.sleepBefore == 0 })).Return(nil, nil).
		Once()
	_, err = decorator.asyncRequestRollupInfoByBlockRange(ctx, newRequestNoSleep(blockRange))
	assert.NoError(t, err)
}
