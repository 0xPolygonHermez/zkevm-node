package dbutils

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
)

// InitOrResetState will initializes the State db running the migrations or
// will reset all the known data and rerun the migrations
func InitOrResetState(cfg db.Config) error {
	return initOrReset(cfg, "zkevm-state-db")
}

// InitOrResetPool will initializes the Pool db running the migrations or
// will reset all the known data and rerun the migrations
func InitOrResetPool(cfg db.Config) error {
	return initOrReset(cfg, "zkevm-pool-db")
}

// InitOrResetRPC will initializes the RPC db running the migrations or
// will reset all the known data and rerun the migrations
func InitOrResetRPC(cfg db.Config) error {
	return initOrReset(cfg, "zkevm-rpc-db")
}

// initOrReset will initializes the db running the migrations or
// will reset all the known data and rerun the migrations
func initOrReset(cfg db.Config, name string) error {
	// connect to database
	dbPool, err := db.NewSQLDB(cfg)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	// run migrations
	if err := db.RunMigrationsDown(cfg, name); err != nil {
		return err
	}
	return db.RunMigrationsUp(cfg, name)
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
