package synchronizer

import (
	"testing"
)

type mocksOrgertration struct {
	producer *l1RollupProducerInterfaceMock
	consumer *l1RollupConsumerInterfaceMock
}

func Test_A(t *testing.T) {
	t.Skip("TODO")
	sut, _ := setupOrchestrationTest(t)
	sut.start(123)
}

func setupOrchestrationTest(t *testing.T) (*l1SyncOrchestration, mocksOrgertration) {
	producer := newL1RollupProducerInterfaceMock(t)
	consumer := newL1RollupConsumerInterfaceMock(t)
	return newL1SyncOrchestration(producer, consumer), mocksOrgertration{
		producer: producer,
		consumer: consumer,
	}
}
