package dbutils

import (
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
