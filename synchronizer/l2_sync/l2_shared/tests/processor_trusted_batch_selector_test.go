package test_l2_shared

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	"github.com/stretchr/testify/require"
)

// Use case 1:
// - Running incaberry mode no forkid7 yet
// expected:
// -

func TestExecutorSelectorIncaberryBatchNoForkId7(t *testing.T) {
	mockIncaberry := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mock1Etrog := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mockState := mock_syncinterfaces.NewStateFullInterface(t)
	mockState.EXPECT().GetForkIDInMemory(uint64(7)).Return(nil)
	sut := l2_shared.NewSyncTrustedStateExecutorSelector(mockIncaberry, mock1Etrog, mockState)

	executor, maxBatch := sut.GetExecutor(1, 200)
	require.Equal(t, mockIncaberry, executor)
	require.Equal(t, uint64(200), maxBatch)
}

func TestExecutorSelectorIncaberryBatchForkId7(t *testing.T) {
	mockIncaberry := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mock1Etrog := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mockState := mock_syncinterfaces.NewStateFullInterface(t)
	interval := state.ForkIDInterval{
		FromBatchNumber: 100,
		ToBatchNumber:   200,
		ForkId:          7,
	}
	mockState.EXPECT().GetForkIDInMemory(uint64(7)).Return(&interval)
	sut := l2_shared.NewSyncTrustedStateExecutorSelector(mockIncaberry, mock1Etrog, mockState)

	executor, maxBatch := sut.GetExecutor(1, 200)
	require.Equal(t, mockIncaberry, executor)
	require.Equal(t, uint64(99), maxBatch)
}

func TestExecutorSelectorEtrogBatchForkId7(t *testing.T) {
	mockIncaberry := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mock1Etrog := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mockState := mock_syncinterfaces.NewStateFullInterface(t)
	interval := state.ForkIDInterval{
		FromBatchNumber: 100,
		ToBatchNumber:   300,
		ForkId:          7,
	}
	mockState.EXPECT().GetForkIDInMemory(uint64(7)).Return(&interval)
	sut := l2_shared.NewSyncTrustedStateExecutorSelector(mockIncaberry, mock1Etrog, mockState)

	executor, maxBatch := sut.GetExecutor(100, 200)
	require.Equal(t, mockIncaberry, executor)
	require.Equal(t, uint64(200), maxBatch)
}
