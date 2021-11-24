package db

import (
	"database/sql"

	"github.com/gobuffalo/packr/v2"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

var migrations *migrate.PackrMigrationSource

// NewSQLDB creates a new SQL DB
func NewSQLDB(cfg Config) (*sql.DB, error) {
	c, err := pgx.ParseConfig("postgres://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.Database)
	if err != nil {
		return nil, err
	}
	db := stdlib.OpenDB(*c)
	return db, nil
}

// RunMigrations will execute pending migrations if needed to keep
// the database updated with the latest changes
func RunMigrations(cfg Config) error {
	db, err := NewSQLDB(cfg)
	if err != nil {
		return err
	}

	migrations = &migrate.PackrMigrationSource{Box: packr.New("hermez-db-migrations", "./migrations")}
	nMigrations, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Info("successfully ran ", nMigrations, " migrations Up")
	return nil
}
