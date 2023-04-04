package pgeventstorage

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
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
func (p *PostgresEventStorage) LogEvent(ctx context.Context, event *event.Event) error {
	const insertEventSQL = "INSERT INTO event (received_at, ip_address, source, component, level, eventid, description, data, json) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	log.Debugf("Event: %v", event)
	_, err := p.db.Exec(ctx, insertEventSQL, event.ReceivedAt, event.IPAddress, event.Source, event.Component, event.Level, event.EventID, event.Description, event.Data, event.Json)
	return err
}
