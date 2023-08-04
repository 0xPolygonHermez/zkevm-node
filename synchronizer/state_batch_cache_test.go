package synchronizer

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	ttlOfStateDbCacheTest = 60 * time.Second
)

func Test_Given_CacheState_Act_StateInterface(t *testing.T) {
	stateMock := newStateMock(t)
	s := NewSynchronizerStateBatchCache(stateMock, 2, ttlOfStateDbCacheTest)
	require.Implements(t, (*stateInterface)(nil), s)
}

func Test_Given_CacheWithAnElement_When_Retieve_Then_ReturnCachedOne(t *testing.T) {
	stateMock := newStateMock(t)
	s := NewSynchronizerStateBatchCache(stateMock, 2, ttlOfStateDbCacheTest)
	batch := state.Batch{BatchNumber: 123}
	s.Set(&batch)
	retrieved_batch, err := s.GetBatchByNumber(context.TODO(), batch.BatchNumber, nil)
	require.NoError(t, err)
	require.Equal(t, &batch, retrieved_batch)
}

func Test_Given_CacheWithNoElement_When_Retrieve_Then_AskToRealStateDb(t *testing.T) {
	stateMock := newStateMock(t)
	batch := state.Batch{BatchNumber: 123}
	stateMock.
		On("GetBatchByNumber", mock.Anything, uint64(batch.BatchNumber), mock.Anything).
		Return(&batch, nil).
		Once()
	s := NewSynchronizerStateBatchCache(stateMock, 2, ttlOfStateDbCacheTest)

	retrieved_batch, err := s.GetBatchByNumber(context.TODO(), 123, nil)
	require.NoError(t, err)
	require.Equal(t, &batch, retrieved_batch)
}

func Test_Given_CacheWithNoElement_When_RetrieveTwice_Then_OnlyFirstTimeAskToRealStateDb(t *testing.T) {
	stateMock := newStateMock(t)
	batch := state.Batch{BatchNumber: 123}
	stateMock.
		On("GetBatchByNumber", mock.Anything, uint64(batch.BatchNumber), mock.Anything).
		Return(&batch, nil).
		Once()
	s := NewSynchronizerStateBatchCache(stateMock, 2, ttlOfStateDbCacheTest)

	_, err := s.GetBatchByNumber(context.TODO(), 123, nil)
	require.NoError(t, err)
	var retrieved_batch *state.Batch
	retrieved_batch, err = s.GetBatchByNumber(context.TODO(), 123, nil)
	require.NoError(t, err)
	require.Equal(t, &batch, retrieved_batch)
}

func Test_Given_CacheWithMoreBatchesInsertedThanCapacity_When_RetrieveOldBatch_Then_Fails(t *testing.T) {
	s := NewSynchronizerStateBatchCache(nil, 2, ttlOfStateDbCacheTest)
	for i := 0; i < 3; i++ {
		batch := state.Batch{BatchNumber: uint64(100 + i)}
		s.Set(&batch)
	}
	require.Equal(t, 2, s.numElements())
	var err error
	_, err = s.GetBatchByNumber(context.TODO(), 102, nil)
	require.NoError(t, err)
	_, err = s.GetBatchByNumber(context.TODO(), 101, nil)
	require.NoError(t, err)

	_, err = s.GetBatchByNumber(context.TODO(), 100, nil)
	require.Error(t, err)
}

func Test_Given_CacheWithMoreBatchesInsertedThanCapacity_When_Clean_Then_DoesntFoundAnyBatch(t *testing.T) {
	s := NewSynchronizerStateBatchCache(nil, 2, ttlOfStateDbCacheTest)
	for i := 0; i < 3; i++ {
		batch := state.Batch{BatchNumber: uint64(100 + i)}
		s.Set(&batch)
	}
	require.Equal(t, 2, s.numElements())
	s.CleanCache()

	var err error
	_, err = s.GetBatchByNumber(context.TODO(), 102, nil)
	require.Error(t, err)
	_, err = s.GetBatchByNumber(context.TODO(), 101, nil)
	require.Error(t, err)

	_, err = s.GetBatchByNumber(context.TODO(), 100, nil)
	require.Error(t, err)
}
