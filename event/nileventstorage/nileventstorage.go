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
func (p *NilEventStorage) LogEvent(ctx context.Context, ev *event.Event) error {
	switch ev.Level {
	case event.Level_Emergency:
		log.Error("Event: %+v", ev)
	case event.Level_Alert:
		log.Error("Event: %+v", ev)
	case event.Level_Critical:
		log.Error("Event: %+v", ev)
	case event.Level_Error:
		log.Error("Event: %+v", ev)
	case event.Level_Warning:
		log.Warn("Event: %+v", ev)
	case event.Level_Notice:
		log.Info("Event: %+v", ev)
	case event.Level_Info:
		log.Info("Event: %+v", ev)
	case event.Level_Debug:
		log.Debug("Event: %+v", ev)
	}
	return nil
}
