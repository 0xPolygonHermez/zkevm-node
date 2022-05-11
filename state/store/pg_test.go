package store_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/store"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

const (
	txBundleID = "txBundleID"
	tableName  = "test"
)

var (
	dbCfg   = dbutils.NewConfigFromEnv()
	stateDB *pgxpool.Pool
)

func resetDB(tableName string) {
	pg := store.NewPg(stateDB)
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName)
	_, err := pg.Exec(context.Background(), "", sql)
	if err != nil {
		log.Errorf("Error reseting db: %v", err)
	}
}

func TestMain(m *testing.M) {
	var err error

	if err := dbutils.InitOrReset(dbCfg); err != nil {
		panic(err)
	}

	stateDB, err = db.NewSQLDB(dbCfg)
	if err != nil {
		panic(err)
	}

	defer stateDB.Close()

	result := m.Run()
	os.Exit(result)
}

func TestPGStoreCommitedTransaction(t *testing.T) {
	defer resetDB(tableName)
	subject := store.NewPg(stateDB)
	ctx := context.Background()

	require.NoError(t, subject.BeginDBTransaction(ctx, txBundleID))
	_, err := subject.Exec(ctx, txBundleID, "DROP TABLE IF EXISTS test; CREATE TABLE test (field VARCHAR ( 32 ));")
	require.NoError(t, err)

	_, err = subject.Exec(ctx, txBundleID, "INSERT INTO test (field) VALUES ($1)", "testValue")
	require.NoError(t, err)

	require.NoError(t, subject.Commit(ctx, txBundleID))

	row := subject.QueryRow(ctx, "", "SELECT field FROM test;")
	var r struct{ Field string }
	require.NoError(t, row.Scan(&r.Field))

	require.Equal(t, "testValue", r.Field)
}

func TestPGStoreRolledbackTransaction(t *testing.T) {
	defer resetDB(tableName)
	subject := store.NewPg(stateDB)
	ctx := context.Background()

	require.NoError(t, subject.BeginDBTransaction(ctx, txBundleID))
	_, err := subject.Exec(ctx, txBundleID, "DROP TABLE IF EXISTS test; CREATE TABLE test (field VARCHAR ( 32 ));")
	require.NoError(t, err)

	_, err = subject.Exec(ctx, txBundleID, "INSERT INTO test (field) VALUES ($1)", "testValue")
	require.NoError(t, err)

	require.NoError(t, subject.Rollback(ctx, txBundleID))

	row := subject.QueryRow(ctx, "", "SELECT field FROM test;")
	var r struct{ Field string }
	err = row.Scan(&r.Field)
	require.Error(t, err)
	require.Contains(t, err.Error(), `ERROR: relation "test" does not exist`)
}

func TestPGStoreNonExistentTransaction(t *testing.T) {
	subject := store.NewPg(stateDB)
	ctx := context.Background()

	require.NoError(t, subject.BeginDBTransaction(ctx, txBundleID))
	defer require.NoError(t, subject.Rollback(ctx, txBundleID))

	_, err := subject.Exec(ctx, "non-existent-tx-bundle-id", "DROP TABLE IF EXISTS test; CREATE TABLE test (field VARCHAR ( 32 ));")
	require.Error(t, err)
	require.Equal(t, `DB Tx bundle "non-existent-tx-bundle-id" does not exist`, err.Error())
}

func TestPGStoreConcurrentTransactions(t *testing.T) {
	defer resetDB(tableName)
	subject := store.NewPg(stateDB)
	ctx := context.Background()

	_, err := subject.Exec(ctx, "", "DROP TABLE IF EXISTS test; CREATE TABLE test (field VARCHAR ( 32 ));")
	require.NoError(t, err)

	var wg sync.WaitGroup
	const valueFmt = "testValue-%03d"

	totalWorkers := 100
	for i := 0; i < totalWorkers; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			txID := fmt.Sprintf("%s-%d", txBundleID, i)
			value := fmt.Sprintf(valueFmt, i)

			require.NoError(t, subject.BeginDBTransaction(ctx, txID))
			_, err := subject.Exec(ctx, txID, "INSERT INTO test (field) VALUES ($1)", value)
			require.NoError(t, err)

			require.NoError(t, subject.Commit(ctx, txID))
		}(i)
	}

	wg.Wait()

	rows, err := subject.Query(ctx, "", "SELECT field FROM test ORDER BY field;")
	require.NoError(t, err)

	var r struct{ Field string }

	count := 0
	for rows.Next() {
		require.NoError(t, rows.Scan(&r.Field))

		expectedValue := fmt.Sprintf(valueFmt, count)
		require.Equal(t, expectedValue, r.Field)
		count++
	}
	require.Equal(t, totalWorkers, count)
}
