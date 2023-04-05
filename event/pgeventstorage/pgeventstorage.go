package pgeventstorage

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresEventStorage is an implementation of the event storage interface
// that uses a postgres database to store the data
type PostgresEventStorage struct {
	db *pgxpool.Pool
}

// NewPostgresEventStorage creates and initializes an instance of PostgresEventStorage
func NewPostgresEventStorage(cfg db.Config) (*PostgresEventStorage, error) {
	poolDB, err := db.NewSQLDB(cfg)
	if err != nil {
		return nil, err
	}

	return &PostgresEventStorage{
		db: poolDB,
	}, nil
}

// Close closes the database connection
func (p *PostgresEventStorage) Close() error {
	p.db.Close()
	return nil
}

// LogEvent logs an event to the database
func (p *PostgresEventStorage) LogEvent(ctx context.Context, ev *event.Event) error {
	const insertEventSQL = "INSERT INTO event (received_at, ip_address, source, component, level, event_id, description, data, json) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	var ipAddressPtr *string
	if ev.IPAddress != "" {
		ipAddressPtr = &ev.IPAddress
	}

	nileventstorage.LogEvent(ev)
	_, err := p.db.Exec(ctx, insertEventSQL, ev.ReceivedAt, ipAddressPtr, ev.Source, ev.Component, ev.Level, ev.EventID, ev.Description, ev.Data, ev.Json)
	return err
}
