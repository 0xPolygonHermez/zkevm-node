package dbutils

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
)

// InitOrReset will initializes the db running the migrations or
// will reset all the known data and rerun the migrations
func InitOrReset(cfg db.Config, dir string) error {
	// connect to database
	dbPool, err := db.NewSQLDB(cfg)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	// run migrations
	if err := db.RunMigrationsDown(cfg, dir); err != nil {
		return err
	}
	return db.RunMigrationsUp(cfg, dir)
}

// NewStateConfigFromEnv return a config for state db
func NewStateConfigFromEnv() db.Config {
	return newConfigFromEnv("state", "5432")
}

// NewPoolConfigFromEnv return a config for pool db
func NewPoolConfigFromEnv() db.Config {
	return newConfigFromEnv("pool", "5433")
}

// NewRPCConfigFromEnv return a config for RPC db
func NewRPCConfigFromEnv() db.Config {
	return newConfigFromEnv("rpc", "5434")
}

// newConfigFromEnv creates config from standard postgres environment variables,
// see https://www.postgresql.org/docs/11/libpq-envars.html for details
func newConfigFromEnv(prefix, port string) db.Config {
	const maxDBPoolConns = 50

	return db.Config{
		User:      testutils.GetEnv("PGUSER", fmt.Sprintf("%v_user", prefix)),
		Password:  testutils.GetEnv("PGPASSWORD", fmt.Sprintf("%v_password", prefix)),
		Name:      testutils.GetEnv("PGDATABASE", fmt.Sprintf("%v_db", prefix)),
		Host:      testutils.GetEnv("PGHOST", "localhost"),
		Port:      testutils.GetEnv("PGPORT", port),
		EnableLog: true,
		MaxConns:  maxDBPoolConns,
	}
}
