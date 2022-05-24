package dbutils

import (
	"context"
	"os"

	"github.com/hermeznetwork/hermez-core/db"
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

	// reset db droping migrations table and schemas
	if _, err := dbPool.Exec(context.Background(), "DROP TABLE IF EXISTS gorp_migrations CASCADE;"); err != nil {
		return err
	}
	if _, err := dbPool.Exec(context.Background(), "DROP SCHEMA IF EXISTS state CASCADE;"); err != nil {
		return err
	}
	if _, err := dbPool.Exec(context.Background(), "DROP SCHEMA IF EXISTS pool CASCADE;"); err != nil {
		return err
	}
	if _, err := dbPool.Exec(context.Background(), "DROP SCHEMA IF EXISTS rpc CASCADE;"); err != nil {
		return err
	}

	// run migrations
	return db.RunMigrations(cfg)
}

// NewConfigFromEnv creates config from standard postgres environment variables,
// see https://www.postgresql.org/docs/11/libpq-envars.html for details
func NewConfigFromEnv() db.Config {
	const maxDBPoolConns = 50

	return db.Config{
		User:      getEnv("PGUSER", "test_user"),
		Password:  getEnv("PGPASSWORD", "test_password"),
		Name:      getEnv("PGDATABASE", "test_db"),
		Host:      getEnv("PGHOST", "localhost"),
		Port:      getEnv("PGPORT", "5432"),
		EnableLog: true,
		MaxConns:  maxDBPoolConns,
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
