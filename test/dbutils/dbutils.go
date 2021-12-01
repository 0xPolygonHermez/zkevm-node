package dbutils

import (
	"context"

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

	// run migrations
	if err := db.RunMigrations(cfg); err != nil {
		return err
	}
	return nil
}
