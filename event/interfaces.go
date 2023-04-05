package event

import (
	"context"
)

// Storage is the interface for the event storage
type Storage interface {
	// LogEvent logs an event
	LogEvent(ctx context.Context, event *Event) error
}
