package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// FilterTypeLog represents a filter of type log.
	FilterTypeLog = "log"
	// FilterTypeBlock represents a filter of type block.
	FilterTypeBlock = "block"
	// FilterTypePendingTx represent a filter of type pending Tx.
	FilterTypePendingTx = "pendingTx"
)

// ErrNotFound represent a not found error.
var ErrNotFound = errors.New("object not found")

// PostgresStorage uses a postgres database to store the data
// related to the json rpc server
type PostgresStorage struct {
	db *pgxpool.Pool
}

// NewPostgresStorage creates and initializes an instance of PostgresStorage
func NewPostgresStorage(cfg db.Config) (*PostgresStorage, error) {
	poolDB, err := db.NewSQLDB(cfg)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: poolDB,
	}, nil
}

// NewLogFilter persists a new log filter
func (s *PostgresStorage) NewLogFilter(filter LogFilter) (uint64, error) {
	parametersBytes, err := json.Marshal(&filter)
	if err != nil {
		return 0, err
	}

	parameters := string(parametersBytes)

	return s.insertFilter(Filter{
		Type:       FilterTypeLog,
		Parameters: parameters,
	})
}

// NewBlockFilter persists a new block log filter
func (s *PostgresStorage) NewBlockFilter() (uint64, error) {
	return s.insertFilter(Filter{
		Type:       FilterTypeBlock,
		Parameters: "{}",
	})
}

// NewPendingTransactionFilter persists a new pending transaction filter
func (s *PostgresStorage) NewPendingTransactionFilter() (uint64, error) {
	return s.insertFilter(Filter{
		Type:       FilterTypePendingTx,
		Parameters: "{}",
	})
}

// insertFilter persists the filter to the db
func (s *PostgresStorage) insertFilter(filter Filter) (uint64, error) {
	lastPoll := time.Now().UTC()
	sql := `INSERT INTO rpc.filters (filter_type, parameters, last_poll) VALUES($1, $2, $3) RETURNING "id"`

	var id uint64
	err := s.db.QueryRow(context.Background(), sql, filter.Type, filter.Parameters, lastPoll).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetFilter gets a filter by its id
func (s *PostgresStorage) GetFilter(filterID uint64) (*Filter, error) {
	filter := &Filter{}
	sql := `SELECT id, filter_type, parameters, last_poll FROM rpc.filters WHERE id = $1`
	err := s.db.QueryRow(context.Background(), sql, filterID).Scan(&filter.ID, &filter.Type, &filter.Parameters, &filter.LastPoll)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return filter, nil
}

// UpdateFilterLastPoll updates the last poll to now
func (s *PostgresStorage) UpdateFilterLastPoll(filterID uint64) error {
	sql := "UPDATE rpc.filters SET last_poll = $2 WHERE id = $1"
	_, err := s.db.Exec(context.Background(), sql, filterID, time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}

// UninstallFilter deletes a filter by its id
func (s *PostgresStorage) UninstallFilter(filterID uint64) (bool, error) {
	sql := "DELETE FROM rpc.filters WHERE id = $1"
	res, err := s.db.Exec(context.Background(), sql, filterID)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}
