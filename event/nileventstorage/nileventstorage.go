package nileventstorage

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

// NilEventStorage is an implementation of the event storage interface
// that just logs but does not store the data
type NilEventStorage struct {
}

// NewNilEventStorage creates and initializes an instance of NewNilEventStorage
func NewNilEventStorage() (*NilEventStorage, error) {
	return &NilEventStorage{}, nil
}

// LogEvent logs an event
func (p *NilEventStorage) LogEvent(ctx context.Context, event *event.Event) error {
	log.Debugf("Event: %v", event)
	return nil
}
