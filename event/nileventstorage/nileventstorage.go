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

// LogEvent logs an event following the defined interface
func (p *NilEventStorage) LogEvent(ctx context.Context, ev *event.Event) error {
	LogEvent(ev)
	return nil
}

// LogEvent actually logs the event
func LogEvent(ev *event.Event) {
	switch ev.Level {
	case event.Level_Emergency, event.Level_Alert, event.Level_Critical, event.Level_Error:
		log.Errorf("Event: %+v", ev)
	case event.Level_Warning, event.Level_Notice:
		log.Warnf("Event: %+v", ev)
	case event.Level_Info:
		log.Infof("Event: %+v", ev)
	case event.Level_Debug:
		log.Debugf("Event: %+v", ev)
	}
}
