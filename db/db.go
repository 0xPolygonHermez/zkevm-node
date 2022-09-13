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
func RunMigrationsUp(cfg Config, dir string) error {
	return runMigrations(cfg, dir, migrate.Up)
}

// RunMigrationsDown runs migrate-down for the given config.
func RunMigrationsDown(cfg Config, dir string) error {
	return runMigrations(cfg, dir, migrate.Down)
}

// runMigrations will execute pending migrations if needed to keep
// the database updated with the latest changes in either direction,
// up or down.
func runMigrations(cfg Config, dir string, direction migrate.MigrationDirection) error {
	c, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name))
	if err != nil {
		return err
	}
	db := stdlib.OpenDB(*c)

	var migrations = &migrate.PackrMigrationSource{Box: packr.New("zkevm-db-migrations", dir)}
	nMigrations, err := migrate.Exec(db, "postgres", migrations, direction)
	if err != nil {
		return err
	}

	log.Info("successfully ran ", nMigrations, " migrations")
	return nil
}
