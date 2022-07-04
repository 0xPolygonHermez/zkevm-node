package dbutils

import (
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/test/testutils"
)

// InitOrReset will initializes the db running the migrations or
// will reset all the known data and rerun the migrations
func InitOrReset(cfg db.Config) error {
	// connect to database
	dbPool, err := db.NewSQLDB(cfg)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	// run migrations
	if err := db.RunMigrationsDown(cfg); err != nil {
		return err
	}
	return db.RunMigrationsUp(cfg)
}

// NewConfigFromEnv creates config from standard postgres environment variables,
// see https://www.postgresql.org/docs/11/libpq-envars.html for details
func NewConfigFromEnv() db.Config {
	const maxDBPoolConns = 50

	return db.Config{
		User:      testutils.GetEnv("PGUSER", "test_user"),
		Password:  testutils.GetEnv("PGPASSWORD", "test_password"),
		Name:      testutils.GetEnv("PGDATABASE", "test_db"),
		Host:      testutils.GetEnv("PGHOST", "localhost"),
		Port:      testutils.GetEnv("PGPORT", "5434"),
		EnableLog: true,
		MaxConns:  maxDBPoolConns,
	}
}
