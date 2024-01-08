package l2_shared

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
			Coinbase:    common.HexToAddress("0x123"),
			StateRoot:   common.HexToHash("0x123"),
			WIP:         true,
		},
		statePreviousBatch: &state.Batch{
			BatchNumber: 122,
			Coinbase:    common.HexToAddress("0x123"),
			StateRoot:   common.HexToHash("0x122"),
			WIP:         false,
		},
		trustedNodeBatch: &types.Batch{
			Number:    123,
			Coinbase:  common.HexToAddress("0x123"),
			StateRoot: common.HexToHash("0x123-1"),
			Closed:    true,
		},
	}

}

func TestGetModeForProcessBatch(t *testing.T) {
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
