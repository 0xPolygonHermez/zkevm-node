package db

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/gobuffalo/packr/v2"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

var packrMigrations = map[string]*packr.Box{
	"zkevm-state-db": packr.New("zkevm-state-db", "./migrations/state"),
	"zkevm-pool-db":  packr.New("zkevm-pool-db", "./migrations/pool"),
	"zkevm-rpc-db":   packr.New("zkevm-rpc-db", "./migrations/rpc"),
}

// NewSQLDB creates a new SQL DB
func NewSQLDB(cfg Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=%d", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.MaxConns))
	if err != nil {
		log.Errorf("Unable to parse DB config: %v\n", err)
		return nil, err
	}
	if cfg.EnableLog {
		config.ConnConfig.Logger = logger{}
	}
	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Errorf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	return conn, nil
}

// RunMigrationsUp runs migrate-up for the given config.
func RunMigrationsUp(cfg Config, name string) error {
	return runMigrations(cfg, name, migrate.Up)
}

// RunMigrationsDown runs migrate-down for the given config.
func RunMigrationsDown(cfg Config, name string) error {
	return runMigrations(cfg, name, migrate.Down)
}

// runMigrations will execute pending migrations if needed to keep
// the database updated with the latest changes in either direction,
// up or down.
func runMigrations(cfg Config, packrName string, direction migrate.MigrationDirection) error {
	c, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name))
	if err != nil {
		return err
	}
	db := stdlib.OpenDB(*c)

	box, ok := packrMigrations[packrName]
	if !ok {
		return fmt.Errorf("packr box not found with name: %v", packrName)
	}

	var migrations = &migrate.PackrMigrationSource{Box: box}
	nMigrations, err := migrate.Exec(db, "postgres", migrations, direction)
	if err != nil {
		return err
	}

	log.Info("successfully ran ", nMigrations, " migrations")
	return nil
}
