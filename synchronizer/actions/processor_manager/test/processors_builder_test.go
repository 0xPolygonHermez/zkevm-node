package processor_manager_test

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/processor_manager"
	"github.com/stretchr/testify/assert"
)

func TestL1EventProcessorsBuilder_Register(t *testing.T) {
	// Create a new instance of L1EventProcessorsBuilder
	builder := processor_manager.NewL1EventProcessorsBuilder()

	// Create a mock L1EventProcessor
	mockProcessor := &ProcessorStub{
		name:             "mockProcessor",
		supportedEvents:  []etherman.EventOrder{"event1", "event2"},
		supportedForkIds: []actions.ForkIdType{1, 2},
	}
	// Register the mock processor
	builder.Register(mockProcessor)
	result := builder.Build()
	// Verify that the processor is registered for all supported fork IDs and events
	for _, forkID := range mockProcessor.SupportedForkIds() {
		for _, event := range mockProcessor.SupportedEvents() {
			processor := result.Get(forkID, event)
			assert.Equal(t, mockProcessor, processor, "Registered processor should match the mock processor")
		}
	}
}
