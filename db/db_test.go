package db

import (
	"context"
	"testing"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cfg = Config{
	Database: "testing",
	User:     "hermez",
	Password: "password",
	Host:     "localhost",
	Port:     "5432",
}

func TestDBService(t *testing.T) {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})

	err := dbutils.StartPostgreSQL(cfg.Database, cfg.User, cfg.Password, "")
	require.NoError(t, err)

	err = RunMigrations(cfg)
	require.NoError(t, err)

	db, err := NewSQLDB(cfg)
	require.NoError(t, err)

	var result uint
	err = db.QueryRow(context.Background(), "select count(*) from block").Scan(&result)
	require.NoError(t, err)
	assert.Equal(t, result, uint(0))

	db.Close() //nolint:gosec,errcheck

	// Stop DB Server
	err = dbutils.StopPostgreSQL()
	require.NoError(t, err)
}
