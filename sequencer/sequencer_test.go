package sequencer

import (
	"context"
	"errors"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestFinalizer_getLastBatchNumAndOldStateRoot(t *testing.T) {
	s := setupSequencer()
	dbManagerMock := new(DbManagerMock)
	testCases := []struct {
		name              string
		mockBatches       []*state.Batch
		mockError         error
		expectedBatchNum  uint64
		expectedStateRoot common.Hash
		expectedError     error
	}{
		{
			name: "Success with two batches",
			mockBatches: []*state.Batch{
				{BatchNumber: 2, StateRoot: common.BytesToHash([]byte("stateRoot2"))},
				{BatchNumber: 1, StateRoot: common.BytesToHash([]byte("stateRoot1"))},
			},
			mockError:         nil,
			expectedBatchNum:  2,
			expectedStateRoot: common.BytesToHash([]byte("stateRoot1")),
			expectedError:     nil,
		},
		{
			name: "Success with one batch",
			mockBatches: []*state.Batch{
				{BatchNumber: 1, StateRoot: common.BytesToHash([]byte("stateRoot1"))},
			},
			mockError:         nil,
			expectedBatchNum:  1,
			expectedStateRoot: common.BytesToHash([]byte("stateRoot1")),
			expectedError:     nil,
		},
		{
			name:              "Error while getting batches",
			mockBatches:       nil,
			mockError:         errors.New("database err"),
			expectedBatchNum:  0,
			expectedStateRoot: common.Hash{},
			expectedError:     errors.New("failed to get last 2 batches, err: database err"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			dbManagerMock.On("GetLastNBatches", context.Background(), uint(2)).Return(tc.mockBatches, tc.mockError).Once()

			// act
			batchNum, stateRoot, err := s.getLastBatchNumAndOldStateRoot(context.Background(), dbManagerMock)

			// assert
			assert.Equal(t, tc.expectedBatchNum, batchNum)
			assert.Equal(t, tc.expectedStateRoot, stateRoot)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			dbManagerMock.AssertExpectations(t)
		})
	}
}

func setupSequencer() *Sequencer {
	return &Sequencer{}
}
