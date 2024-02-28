package db

import (
	"context"
	"embed"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

const (
	// StateMigrationName is the name of the migration for state DB
	StateMigrationName = "zkevm-state-db"
	// PoolMigrationName is the name of the migration for pool DB
	PoolMigrationName = "zkevm-pool-db"
)

var Migrations = map[string]string{
	StateMigrationName: "./migrations/state",
	PoolMigrationName:  "./migrations/pool",
}

//go:embed migrations
var f embed.FS

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
	log.Info("running migrations up")
	return runMigrations(cfg, name, migrate.Up)
}

// CheckMigrations runs migrate-up for the given config.
func CheckMigrations(cfg Config, name string) error {
	return checkMigrations(cfg, name, migrate.Up)
}

// RunMigrationsDown runs migrate-down for the given config.
func RunMigrationsDown(cfg Config, name string) error {
	log.Info("running migrations down")
	return runMigrations(cfg, name, migrate.Down)
}

// runMigrations will execute pending migrations if needed to keep
// the database updated with the latest changes in either direction,
// up or down.
func runMigrations(cfg Config, name string, direction migrate.MigrationDirection) error {
	c, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name))
	if err != nil {
		return err
	}
	db := stdlib.OpenDB(*c)

	path, ok := migrations[name]
	if !ok {
		return fmt.Errorf("migration not found with name: %v", name)
	}

	var migrations = &migrate.EmbedFileSystemMigrationSource{
		FileSystem: f,
		Root:       path,
	}
	nMigrations, err := migrate.Exec(db, "postgres", migrations, direction)
	if err != nil {
		return err
	}

	log.Info("successfully ran ", nMigrations, " migrations")
	return nil
}

func checkMigrations(cfg Config, name string, direction migrate.MigrationDirection) error {
	c, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name))
	if err != nil {
		return err
	}
	db := stdlib.OpenDB(*c)

	path, ok := migrations[name]
	if !ok {
		return fmt.Errorf("migration not found with name: %v", name)
	}

	migrationSource := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: f,
		Root:       path,
	}
	migrations, err := migrationSource.FindMigrations()
	if err != nil {
		log.Errorf("error getting migrations from source: %v", err)
		return err
	}

	var expected int
	for _, migration := range migrations {
		if len(migration.Up) != 0 {
			expected++
		}
	}

	var actual int
	query := `SELECT COUNT(1) FROM public.gorp_migrations`
	err = db.QueryRow(query).Scan(&actual)
	if err != nil {
		log.Error("error getting migrations count: ", err)
		return err
	}
	if expected == actual {
		log.Infof("Found %d migrations as expected", actual)
	} else {
		return fmt.Errorf("error the component needs to run %d migrations before starting. DB only contains %d migrations", expected, actual)
	}
	return nil
}
