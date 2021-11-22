package database

import (
	"context"
	"testing"

	dbutils "github.com/hermeznetwork/hermez-core/test/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	dbName     = "testing"
	dbUser     = "hermez"
	dbPassword = "password"
)

func TestDBService(t *testing.T) {
	// Start DB Server
	err := dbutils.StartPostgreSQL(dbName, dbUser, dbPassword, "./migrations/0001.sql")
	require.NoError(t, err)

	db, err := NewSQLDB(dbName, dbUser, dbPassword, dbutils.DBHost, dbutils.DBPort)
	require.NoError(t, err)

	var result uint
	err = db.QueryRow(context.Background(), "select count(*) from block").Scan(&result)
	require.NoError(t, err)
	assert.Equal(t, result, uint(0))

	db.Close()

	// Stop DB Server
	err = dbutils.StopPostgreSQL()
	require.NoError(t, err)
}
