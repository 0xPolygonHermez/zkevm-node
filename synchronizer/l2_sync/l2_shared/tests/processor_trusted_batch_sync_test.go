package test_l2_shared

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	commonSync "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	mock_l2_shared "github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestCacheEmpty(t *testing.T) {
	mockExecutor := mock_l2_shared.NewSyncTrustedBatchExecutor(t)
	mockTimer := &commonSync.MockTimerProvider{}
	sut := l2_shared.NewProcessorTrustedBatchSync(mockExecutor, mockTimer)

	current, previous := sut.GetCurrentAndPreviousBatchFromCache(&l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil, nil},
	})
	require.Nil(t, current)
	require.Nil(t, previous)
	current, previous = sut.GetCurrentAndPreviousBatchFromCache(&l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil},
	})
	require.Nil(t, current)
	require.Nil(t, previous)

	current, previous = sut.GetCurrentAndPreviousBatchFromCache(&l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{},
	})
	require.Nil(t, current)
	require.Nil(t, previous)
}

func TestCacheJustCurrent(t *testing.T) {
	mockExecutor := mock_l2_shared.NewSyncTrustedBatchExecutor(t)
	mockTimer := &commonSync.MockTimerProvider{}
	batchA := state.Batch{
		BatchNumber: 123,
		Coinbase:    common.HexToAddress("0x123"),
	}
	status := l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{&batchA},
	}
	sut := l2_shared.NewProcessorTrustedBatchSync(mockExecutor, mockTimer)

	current, previous := sut.GetCurrentAndPreviousBatchFromCache(&status)
	require.Nil(t, previous)
	require.Equal(t, &batchA, current)
	require.NotEqual(t, &batchA, &current)
}

func TestCacheJustPrevious(t *testing.T) {
	mockExecutor := mock_l2_shared.NewSyncTrustedBatchExecutor(t)
	mockTimer := &commonSync.MockTimerProvider{}
	batchA := state.Batch{
		BatchNumber: 123,
		Coinbase:    common.HexToAddress("0x123"),
	}
	status := l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil, &batchA},
	}
	sut := l2_shared.NewProcessorTrustedBatchSync(mockExecutor, mockTimer)

	current, previous := sut.GetCurrentAndPreviousBatchFromCache(&status)
	require.Nil(t, current)
	require.Equal(t, &batchA, previous)
	require.NotEqual(t, &batchA, &previous)
}

type TestDataForProcessorTrustedBatchSync struct {
	mockTimer          *commonSync.MockTimerProvider
	mockExecutor       *mock_l2_shared.SyncTrustedBatchExecutor
	sut                *l2_shared.ProcessorTrustedBatchSync
	trustedNodeBatch   *types.Batch
	stateCurrentBatch  *state.Batch
	statePreviousBatch *state.Batch
}

func newTestDataForProcessorTrustedBatchSync(t *testing.T) *TestDataForProcessorTrustedBatchSync {
	mockExecutor := mock_l2_shared.NewSyncTrustedBatchExecutor(t)
	mockTimer := &commonSync.MockTimerProvider{}
	return &TestDataForProcessorTrustedBatchSync{
		mockTimer:    mockTimer,
		mockExecutor: mockExecutor,
		sut:          l2_shared.NewProcessorTrustedBatchSync(mockExecutor, mockTimer),
		stateCurrentBatch: &state.Batch{
			BatchNumber: 123,
			Coinbase:    common.HexToAddress("0x1230"),
			StateRoot:   common.HexToHash("0x123410"),
			WIP:         true,
		},
		statePreviousBatch: &state.Batch{
			BatchNumber: 122,
			Coinbase:    common.HexToAddress("0x1230"),
			StateRoot:   common.HexToHash("0x1220"),
			WIP:         false,
		},
		trustedNodeBatch: &types.Batch{
			Number:    123,
			Coinbase:  common.HexToAddress("0x1230"),
			StateRoot: common.HexToHash("0x123410"),
			Closed:    true,
		},
	}
}

func TestGetModeForProcessBatchIncremental(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.trustedNodeBatch.Closed = true
	testData.trustedNodeBatch.BatchL2Data = []byte("test")
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.IncrementalProcessMode, processData.Mode, "current batch is WIP and have a intermediate state root")
	require.Equal(t, true, processData.BatchMustBeClosed, "the trustedNode batch is closed")
	require.Equal(t, testData.stateCurrentBatch.StateRoot, processData.OldStateRoot, "the old state root is the intermediate state root (the current batch state root)")
}

func TestGetModeForProcessBatchFullProcessMode(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.trustedNodeBatch.Closed = true
	testData.trustedNodeBatch.BatchL2Data = []byte("test") // We add some data
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, nil, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.FullProcessMode, processData.Mode, "there is no local batch, so it needs to full process")
	require.Equal(t, true, processData.BatchMustBeClosed, "the trustedNode batch is closed")
	require.Equal(t, testData.statePreviousBatch.StateRoot, processData.OldStateRoot, "the old state root is the previous batch SR")
}

func TestGetModeForProcessBatchReprocessMode(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.trustedNodeBatch.Closed = true
	testData.trustedNodeBatch.BatchL2Data = []byte("test") // We add some data
	testData.stateCurrentBatch.StateRoot = state.ZeroHash
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.ReprocessProcessMode, processData.Mode, "local batch doesnt have stateRoot but exists, so  so it needs to be reprocess")
	require.Equal(t, true, processData.BatchMustBeClosed, "the trustedNode batch is closed")
	require.Equal(t, testData.statePreviousBatch.StateRoot, processData.OldStateRoot, "the old state root is the previous batch SR")
}

func TestGetModeForProcessBatchNothing(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.stateCurrentBatch.WIP = true
	testData.trustedNodeBatch.Closed = true
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.NothingProcessMode, processData.Mode, "current batch and trusted batch are the same, just need to be closed")
	require.Equal(t, true, processData.BatchMustBeClosed, "the trustedNode batch is closed")
	require.Equal(t, state.ZeroHash, processData.OldStateRoot, "no OldStateRoot, because you dont need to process anything")

	testData.stateCurrentBatch.WIP = false
	testData.trustedNodeBatch.Closed = true
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.NothingProcessMode, processData.Mode, "current batch and trusted batch are the same, just need to be closed")
	require.Equal(t, false, processData.BatchMustBeClosed, "the trustedNode batch is closed but the state batch is also closed, so nothing to do")

	testData.stateCurrentBatch.WIP = false
	testData.trustedNodeBatch.Closed = false
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.NothingProcessMode, processData.Mode, "current batch and trusted batch are the same, just need to be closed")
	require.Equal(t, false, processData.BatchMustBeClosed, "nothing to do")
}

func TestGetNextStatusClear(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	previousStatus := l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{testData.statePreviousBatch, testData.statePreviousBatch},
	}
	processResponse := l2_shared.NewProcessResponse()

	processResponse.ClearCache = true
	res, err := testData.sut.GetNextStatus(previousStatus, &processResponse, false, "test")
	require.NoError(t, err)
	require.True(t, res.IsEmpty())

	processResponse.ClearCache = false
	res, err = testData.sut.GetNextStatus(l2_shared.TrustedState{}, &processResponse, false, "test")
	require.NoError(t, err)
	require.True(t, res.IsEmpty())

	processResponse.ClearCache = false
	res, err = testData.sut.GetNextStatus(l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil, nil},
	}, &processResponse, false, "test")
	require.NoError(t, err)
	require.True(t, res.IsEmpty())

	processResponse.ClearCache = false
	processResponse.UpdateBatchWithProcessBatchResponse = true
	res, err = testData.sut.GetNextStatus(l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil, nil},
	}, &processResponse, false, "test")
	require.NoError(t, err)
	require.True(t, res.IsEmpty())
}

func TestGetNextStatusUpdate(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	previousStatus := l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{testData.statePreviousBatch, testData.statePreviousBatch},
	}
	processBatchResp := l2_shared.NewProcessResponse()
	newBatch := state.Batch{
		BatchNumber: 123,
		Coinbase:    common.HexToAddress("0x123467"),
		StateRoot:   common.HexToHash("0x123456"),
		WIP:         true,
	}
	processBatchResp.UpdateCurrentBatch(&newBatch)
	res, err := testData.sut.GetNextStatus(previousStatus, &processBatchResp, false, "test")
	require.NoError(t, err)
	require.False(t, res.IsEmpty())
	require.Equal(t, *res.LastTrustedBatches[0], newBatch)

	res, err = testData.sut.GetNextStatus(previousStatus, &processBatchResp, true, "test")
	require.NoError(t, err)
	require.False(t, res.IsEmpty())
	require.Nil(t, res.LastTrustedBatches[0])
	require.Equal(t, newBatch, *res.LastTrustedBatches[1])

	ProcessBatchResponse := &state.ProcessBatchResponse{
		NewStateRoot:     common.HexToHash("0x123-2"),
		NewAccInputHash:  common.HexToHash("0x123-3"),
		NewLocalExitRoot: common.HexToHash("0x123-4"),
		NewBatchNumber:   123,
	}
	processBatchResp.UpdateCurrentBatchWithExecutionResult(&newBatch, ProcessBatchResponse)
	res, err = testData.sut.GetNextStatus(previousStatus, &processBatchResp, true, "test")
	require.NoError(t, err)
	require.False(t, res.IsEmpty())
	require.Nil(t, res.LastTrustedBatches[0])
	require.Equal(t, processBatchResp.ProcessBatchResponse.NewStateRoot, res.LastTrustedBatches[1].StateRoot)
}
