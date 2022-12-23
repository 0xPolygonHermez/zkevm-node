package ethtxmanager

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresStorage hold txs to be managed
type PostgresStorage struct {
	cfg Config
	db  *pgxpool.Pool
}

// NewPostgresStorage creates a new instance of storage that use
// postgres to store data
func NewPostgresStorage(cfg Config, dbCfg db.Config) (storageInterface, error) {
	ethTxManagerDb, err := db.NewSQLDB(dbCfg)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		cfg: cfg,
		db:  ethTxManagerDb,
	}, nil
}

// Add persist a monitored tx
func (s *PostgresStorage) Add(ctx context.Context, mTx monitoredTx) error {
	panic("not implemented yet")
	// if _, ok := s.monitoredTxs[mTx.id]; ok {
	// 	return ErrAlreadyExists
	// }
	// s.monitoredTxs[mTx.id] = mTx
	// return nil
}

// Get loads a persisted monitored tx
func (s *PostgresStorage) Get(ctx context.Context, id string) (monitoredTx, error) {
	panic("not implemented yet")
	// if mTx, ok := s.monitoredTxs[id]; ok {
	// 	return mTx, nil
	// }
	// return monitoredTx{}, ErrNotFound
}

// GetByStatus loads all monitored tx that match the provided status
func (s *PostgresStorage) GetByStatus(ctx context.Context, statuses ...MonitoredTxStatus) ([]monitoredTx, error) {
	panic("not implemented yet")
	// statusesMap := make(map[MonitoredTxStatus]bool, len(statuses))
	// for _, status := range statuses {
	// 	statusesMap[status] = true
	// }

	// mTxs := make([]monitoredTx, 0, len(s.monitoredTxs))
	// for _, mTx := range s.monitoredTxs {
	// 	if _, found := statusesMap[mTx.status]; len(statuses) == 0 || found {
	// 		mTxs = append(mTxs, mTx)
	// 	}
	// }
	// return mTxs, nil
}

// Update a persisted monitored tx
func (s *PostgresStorage) Update(ctx context.Context, mTx monitoredTx) error {
	panic("not implemented yet")
	// s.monitoredTxs[mTx.id] = mTx
	// return nil
}
