package synchronizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type mocksOrgertration struct {
	producer *l1RollupProducerInterfaceMock
	consumer *l1RollupConsumerInterfaceMock
}

func Test_A(t *testing.T) {
	t.Skip("TODO")
	sut, _ := setupOrchestrationTest(t)
	_, err := sut.start(123)
	require.NoError(t, err)
}

func setupOrchestrationTest(t *testing.T) (*l1SyncOrchestration, mocksOrgertration) {
	producer := newL1RollupProducerInterfaceMock(t)
	consumer := newL1RollupConsumerInterfaceMock(t)
	return newL1SyncOrchestration(producer, consumer), mocksOrgertration{
		producer: producer,
		consumer: consumer,
	}
}
