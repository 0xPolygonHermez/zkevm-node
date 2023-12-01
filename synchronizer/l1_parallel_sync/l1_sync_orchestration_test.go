package l1_parallel_sync

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mocksOrgertration struct {
	producer *l1RollupProducerInterfaceMock
	consumer *l1RollupConsumerInterfaceMock
}

func TestGivenOrquestrationWhenHappyPathThenReturnsBlockAndNoErrorAndProducerIsRunning(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	sut, mocks := setupOrchestrationTest(t, ctxTimeout)
	mocks.producer.On("Reset", mock.Anything).Return()
	mocks.producer.On("Start", mock.Anything).Return(func(context.Context) error {
		time.Sleep(time.Second * 2)
		return nil
	})
	block := state.Block{}
	mocks.consumer.On("Reset", mock.Anything).Return()
	mocks.consumer.On("GetLastEthBlockSynced").Return(block, true)
	mocks.consumer.On("Start", mock.Anything, mock.Anything).Return(func(context.Context, *state.Block) error {
		time.Sleep(time.Millisecond * 100)
		return nil
	})
	sut.Reset(123)
	returnedBlock, err := sut.Start(&block)
	require.NoError(t, err)
	require.Equal(t, block, *returnedBlock)
	require.Equal(t, true, sut.producerRunning)
	require.Equal(t, false, sut.consumerRunning)
}

func setupOrchestrationTest(t *testing.T, ctx context.Context) (*L1SyncOrchestration, mocksOrgertration) {
	producer := newL1RollupProducerInterfaceMock(t)
	consumer := newL1RollupConsumerInterfaceMock(t)

	return NewL1SyncOrchestration(ctx, producer, consumer), mocksOrgertration{
		producer: producer,
		consumer: consumer,
	}
}
