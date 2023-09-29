package synchronizer

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
	mocks.producer.On("ResetAndStop", mock.Anything).Return()
	mocks.producer.On("Start", mock.Anything).Return(func(context.Context) error {
		time.Sleep(time.Second * 2)
		return nil
	})
	block := state.Block{}
	mocks.consumer.On("GetLastEthBlockSynced").Return(block, true)
	mocks.consumer.On("Start", mock.Anything).Return(nil)
	sut.reset(123)
	returnedBlock, err := sut.start()
	require.NoError(t, err)
	require.Equal(t, block, *returnedBlock)
	require.Equal(t, true, sut.producerRunning)
	require.Equal(t, false, sut.consumerRunning)
}

func setupOrchestrationTest(t *testing.T, ctx context.Context) (*l1SyncOrchestration, mocksOrgertration) {
	producer := newL1RollupProducerInterfaceMock(t)
	consumer := newL1RollupConsumerInterfaceMock(t)

	return newL1SyncOrchestration(ctx, producer, consumer), mocksOrgertration{
		producer: producer,
		consumer: consumer,
	}
}
