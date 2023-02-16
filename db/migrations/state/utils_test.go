package migrations_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/gobuffalo/packr/v2"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"
)

/*
	Considerations tricks and tips for migration file testing:

	- Functionality of the DB is tested by the rest of the packages, migration tests only have to check persistence across migrations (both UP and DOWN)
	- It's recommended to use real data (from testnet/mainnet), but modifying NULL fields to check that those are migrated properly
	- It's recommended to use some SQL tool (such as DBeaver) that generates insert queries from existing rows
	- Any new migration file could be tested using the existing `migrationTester` interface. Check `0002_test.go` for an example
*/

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stderr"},
	})
}

type migrationTester interface {
	// InsertData used to insert data in the affected tables of the migration that is being tested
	// data will be inserted with the schema as it was previous the migration that is being tested
	InsertData(*sql.DB) error
	// RunAssertsAfterMigrationUp this function will be called after running the migration is being tested
	// and should assert that the data inserted in the function InsertData is persisted properly
	RunAssertsAfterMigrationUp(*testing.T, *sql.DB)
	// RunAssertsAfterMigrationDown this function will be called after reverting the migration that is being tested
	// and should assert that the data inserted in the function InsertData is persisted properly
	RunAssertsAfterMigrationDown(*testing.T, *sql.DB)
}

var (
	stateDBCfg      = dbutils.NewStateConfigFromEnv()
	packrMigrations = map[string]*packr.Box{
		db.StateMigrationName: packr.New(db.StateMigrationName, "./migrations/state"),
		db.PoolMigrationName:  packr.New(db.PoolMigrationName, "./migrations/pool"),
	}
)

func runMigrationTest(t *testing.T, migrationNumber int, miter migrationTester) {
	// Initialize an empty DB
	d, err := initCleanSQLDB()
	require.NoError(t, err)
	require.NoError(t, runMigrationsDown(d, 0, db.StateMigrationName))
	// Run migrations until migration to test
	require.NoError(t, runMigrationsUp(d, migrationNumber-1, db.StateMigrationName))
	// Insert data into table(s) affected by migration
	require.NoError(t, miter.InsertData(d))
	// Run migration that is being tested
	require.NoError(t, runMigrationsUp(d, 1, db.StateMigrationName))
	// Check that data is persisted properly after migration up
	miter.RunAssertsAfterMigrationUp(t, d)
	// Revert migration to test
	require.NoError(t, runMigrationsDown(d, 1, db.StateMigrationName))
	// Check that data is persisted properly after migration down
	miter.RunAssertsAfterMigrationDown(t, d)
}

func initCleanSQLDB() (*sql.DB, error) {
	// run migrations
	if err := db.RunMigrationsDown(stateDBCfg, db.StateMigrationName); err != nil {
		return nil, err
	}
	c, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", stateDBCfg.User, stateDBCfg.Password, stateDBCfg.Host, stateDBCfg.Port, stateDBCfg.Name))
	if err != nil {
		return nil, err
	}
	sqlDB := stdlib.OpenDB(*c)
	return sqlDB, nil
}

func runMigrationsUp(d *sql.DB, n int, packrName string) error {
	box, ok := packrMigrations[packrName]
	if !ok {
		return fmt.Errorf("packr box not found with name: %v", packrName)
	}

	var migrations = &migrate.PackrMigrationSource{Box: box}
	nMigrations, err := migrate.ExecMax(d, "postgres", migrations, migrate.Up, n)
	if err != nil {
		return err
	}
	if nMigrations != n {
		return fmt.Errorf("Unexpected amount of migrations: expected: %d, actual: %d", n, nMigrations)
	}
	return nil
}

func runMigrationsDown(d *sql.DB, n int, packrName string) error {
	box, ok := packrMigrations[packrName]
	if !ok {
		return fmt.Errorf("packr box not found with name: %v", packrName)
	}

	var migrations = &migrate.PackrMigrationSource{Box: box}
	nMigrations, err := migrate.ExecMax(d, "postgres", migrations, migrate.Down, n)
	if err != nil {
		return err
	}
	if nMigrations != n {
		return fmt.Errorf("Unexpected amount of migrations: expected: %d, actual: %d", n, nMigrations)
	}
	return nil
}
