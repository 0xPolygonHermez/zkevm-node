package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
)

// EventLog is the main struct for the event log
type EventLog struct {
	cfg     Config
	storage Storage
}

// NewEventLog creates and initializes an instance of EventLog
func NewEventLog(cfg Config, storage Storage) *EventLog {
	return &EventLog{
		cfg:     cfg,
		storage: storage,
	}
}

// LogEvent is used to store an event for runtime debugging
func (e *EventLog) LogEvent(ctx context.Context, event *Event) error {
	return e.storage.LogEvent(ctx, event)
}

// LogExecutorError is used to store Executor error for runtime debugging
func (e *EventLog) LogExecutorError(ctx context.Context, responseError executor.ExecutorError, processBatchRequest *executor.ProcessBatchRequest) {
	timestamp := time.Now()
	log.Errorf("error found in the executor: %v at %v", responseError, timestamp)
	payload, err := json.Marshal(processBatchRequest)
	if err != nil {
		log.Errorf("error marshaling payload: %v", err)
	} else {
		event := &Event{
			ReceivedAt:  timestamp,
			Source:      Source_Node,
			Component:   Component_Executor,
			Level:       Level_Error,
			EventID:     EventID_ExecutorError,
			Description: responseError.String(),
			Json:        string(payload),
		}
		err = e.storage.LogEvent(ctx, event)
		if err != nil {
			log.Errorf("error storing event: %v", err)
		}
	}
}
