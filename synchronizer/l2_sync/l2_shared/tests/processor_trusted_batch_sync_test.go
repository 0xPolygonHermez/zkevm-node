package l2_shared

import (
	"testing"

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
