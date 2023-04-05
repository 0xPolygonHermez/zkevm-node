package event_test

import (
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/pgeventstorage"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestStoreEvent(t *testing.T) {
	ctx := context.Background()

	eventDBCfg := dbutils.NewEventConfigFromEnv()
	eventStorage, err := pgeventstorage.NewPostgresEventStorage(eventDBCfg)
	require.NoError(t, err)

	ev := &event.Event{
		ReceivedAt:  time.Now(),
		IPAddress:   "127.0.0.1",
		Source:      event.Source_Node,
		Component:   event.Component_Sequencer,
		Level:       event.Level_Error,
		EventID:     event.EventID_ExecutorError,
		Description: "This is a test event",
		Data:        []byte("This is a test event"),
		Json:        eventDBCfg,
	}

	eventLog := event.NewEventLog(event.Config{}, eventStorage)
	defer eventStorage.Close() //nolint:gosec,errcheck

	err = eventLog.LogEvent(ctx, ev)
	require.NoError(t, err)
}
