package state

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestConcurrentSync(t *testing.T) {
	/*
		Case 1:

		Trusted sync init tx
		L1 sync init tx
		Trusted sync reads latest state root, so the dbTx get's invalidated if changed
		Trusted sync open batch N+1
		Trusted sync closes batch N+1
		L1 sync reorgs from batch N-1
		L1 sync commits
		Trusted sync commits
		Expected: Last batch is N-1, trusted sync commit fails due to change on the readed value

		Case 2:

		Trusted sync init tx
		L1 sync init tx
		Trusted sync reads latest state root, so the dbTx get's invalidated if changed
		Trusted sync open batch N+1
		Trusted sync closes batch N+1
		L1 sync reorgs from batch N-1
		Trusted sync commits
		L1 sync commits
		Expected: Last batch is N+1, L1 sync commit fails due to change on the readed value
	*/

	/*
		Solution using isolation: https://www.postgresql.org/docs/current/transaction-iso.html
		PROS:
		* works and easy to implement, just needs to:
			- instantiate the dbTx like this: BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
			- and handle the potential error due to failed to serialize
		* Isolation levels is SQL standard, nothing too crazy
		CONS:
		* potential performance impact
		* could create many retries attempts for the virtual state sync, but this should only be a problem in case of trusted reorg
		* trusted sync needs to do an artifical select to invalidate the insert if the selected content has changed

		locks: https://www.postgresql.org/docs/current/explicit-locking.html
	*/

	// Setup
	stateDBCfg := dbutils.NewStateConfigFromEnv()
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}
	stateDB, err := db.NewSQLDB(stateDBCfg)
	state := NewPostgresStorage(stateDB)
	require.NoError(t, err)
	defer stateDB.Close()

	ctx := context.Background()
	setupDBTx, err := stateDB.Begin(ctx)
	require.NoError(t, err)

	// Insert N trusted batches
	nTrusteBatches := 5
	for i := 1; i <= nTrusteBatches; i++ {
		require.NoError(t, state.openBatch(ctx, ProcessingContext{BatchNumber: uint64(i)}, setupDBTx))
		require.NoError(t, state.closeBatch(ctx, ProcessingReceipt{
			BatchNumber: uint64(i),
			StateRoot:   common.BigToHash(big.NewInt(int64(i))),
		}, setupDBTx))
	}
	require.NoError(t, setupDBTx.Commit(ctx))

	// Using isolation
	//	Case 1
	trustedDBTx, l1DBTx := doDBInteractionsUsingIsolation(t, nTrusteBatches, state, stateDB)
	expectedLastBatch := uint64(nTrusteBatches - 1)
	// L1 sync commits first
	assertsIsolation(t, expectedLastBatch, l1DBTx, trustedDBTx, state, stateDB)
	nTrusteBatches = int(expectedLastBatch)

	//	Case 2
	trustedDBTx, l1DBTx = doDBInteractionsUsingIsolation(t, nTrusteBatches, state, stateDB)
	expectedLastBatch = uint64(nTrusteBatches + 1)
	// Trusted sync commits first
	assertsIsolation(t, expectedLastBatch, trustedDBTx, l1DBTx, state, stateDB)
	nTrusteBatches = int(expectedLastBatch)

	/*
		Solution using locks: https://www.postgresql.org/docs/current/explicit-locking.html and https://www.postgresql.org/docs/current/sql-lock.html
		PROS: TBD
		CONS: TBD
	*/

	// Using isolation
	//	Case 1
	var wg sync.WaitGroup
	trustedDBTx, l1DBTx = doDBInteractionsUsingLocks(t, nTrusteBatches, true, state, stateDB, &wg)
	expectedLastBatch = uint64(nTrusteBatches - 1)
	// L1 sync commits first
	assertsLocks(t, expectedLastBatch, l1DBTx, trustedDBTx, state, stateDB, &wg)
	nTrusteBatches = int(expectedLastBatch)

	//	Case 2
	trustedDBTx, l1DBTx = doDBInteractionsUsingLocks(t, nTrusteBatches, false, state, stateDB, &wg)
	expectedLastBatch = uint64(nTrusteBatches - 1)
	// Trusted sync commits first
	assertsLocks(t, expectedLastBatch, trustedDBTx, l1DBTx, state, stateDB, &wg)
}

func doDBInteractionsUsingIsolation(t *testing.T, nTrusteBatches int, state *PostgresStorage, stateDB *pgxpool.Pool) (trustedDBTx, l1DBTx pgx.Tx) {
	ctx := context.Background()
	var err error
	trustedDBTx, err = stateDB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	require.NoError(t, err)
	l1DBTx, err = stateDB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	require.NoError(t, err)

	// Trusted sync interactions
	var actualRoot string
	err = trustedDBTx.QueryRow(ctx, "SELECT state_root FROM state.batch ORDER BY batch_num DESC LIMIT 1").Scan(&actualRoot)
	require.NoError(t, err)
	require.Equal(t, common.BigToHash(big.NewInt(int64(nTrusteBatches))).Hex(), actualRoot)
	require.NoError(t, state.openBatch(ctx, ProcessingContext{BatchNumber: uint64(nTrusteBatches + 1)}, trustedDBTx))
	require.NoError(t, state.closeBatch(ctx, ProcessingReceipt{BatchNumber: uint64(nTrusteBatches + 1)}, trustedDBTx))

	// L1 sync interactions
	const resetSQL = "DELETE FROM state.batch WHERE batch_num > $1"
	_, err = l1DBTx.Exec(ctx, resetSQL, nTrusteBatches-1)
	require.NoError(t, err)

	return
}

func assertsIsolation(t *testing.T, expectedLastBatchNum uint64, firstCommiter, secondCommiter pgx.Tx, state *PostgresStorage, stateDB *pgxpool.Pool) {
	ctx := context.Background()
	require.NoError(t, firstCommiter.Commit(ctx))
	// https://github.com/jackc/pgx/wiki/Error-Handling
	err := secondCommiter.Commit(ctx)
	require.NotNil(t, err)
	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr))
	require.Equal(t, "40001", pgErr.Code)
	bn, err := state.GetLastBatchNumber(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedLastBatchNum, bn)
}

func doDBInteractionsUsingLocks(t *testing.T, nTrusteBatches int, l1GoesFirst bool, state *PostgresStorage, stateDB *pgxpool.Pool, wg *sync.WaitGroup) (trustedDBTx, l1DBTx pgx.Tx) {
	wg.Add(2)
	log.Debug("INIT doDBInteractionsUsingLocks")
	ctx := context.Background()
	var err error
	trustedDBTx, err = stateDB.Begin(ctx)
	require.NoError(t, err)
	l1DBTx, err = stateDB.Begin(ctx)
	require.NoError(t, err)

	// Lock
	const lockQuery = "SELECT batch_num FROM state.batch ORDER BY batch_num DESC LIMIT 1 FOR UPDATE"
	trustedInteractions := func() {
		// Trusted sync interactions
		var actualBatchNum int
		require.NoError(t, trustedDBTx.QueryRow(ctx, lockQuery).Scan(&actualBatchNum))
		log.Warnf("trustedInteractions reads actualBatchNum: %d", actualBatchNum)
		if actualBatchNum == nTrusteBatches {
			log.Warn("Trusted Interactions executing")
			// L1 Sync has not reorged yet, let's insert next trusted batch
			require.NoError(t, state.openBatch(ctx, ProcessingContext{BatchNumber: uint64(nTrusteBatches + 1)}, trustedDBTx))
			require.NoError(t, state.closeBatch(ctx, ProcessingReceipt{BatchNumber: uint64(nTrusteBatches + 1)}, trustedDBTx))
			log.Warn("Trusted Interactions DONE executing")
		}
		wg.Done()
	}
	l1Interactions := func() {
		// L1 sync interactions
		var actualBatchNum int
		require.NoError(t, l1DBTx.QueryRow(ctx, lockQuery).Scan(&actualBatchNum))
		log.Warnf("l1Interactions reads actualBatchNum: %d", actualBatchNum)
		log.Warn("L1 Interactions executing")
		/*
			NOTE: If the trusted sync locks first, actualBatchNum will still be nTrusteBatches
			instead of nTrusteBatches + 1. This is because nTrusteBatches doesn't get modified.
			However this is not a problem in reality, on the countrary, this is the desired behaviour:
			1. Trusted sync adds the next batch atomically
			2. L1 Sync waits until Trusted sync is done
			3. Trusted sync finishes inserting nTrusteBatches + 1
			4. L1 Sync gets unlocked and deletes nTrusteBatches and nTrusteBatches + 1

			This is not how the test using isolation behaves, as the test finish when the error gets detected, but right after that
			the L1 Sync should re-try the reorg query and achieve the same result
		*/
		const resetSQL = "DELETE FROM state.batch WHERE batch_num > $1"
		_, err = l1DBTx.Exec(ctx, resetSQL, nTrusteBatches-1)
		require.NoError(t, err)
		log.Warn("L1 Interactions DONE executing")
		wg.Done()
	}
	if l1GoesFirst {
		l1Interactions()
		go trustedInteractions()
	} else {
		trustedInteractions()
		go l1Interactions()
	}
	return
}

func assertsLocks(t *testing.T, expectedLastBatchNum uint64, firstCommiter, secondCommiter pgx.Tx, state *PostgresStorage, stateDB *pgxpool.Pool, wg *sync.WaitGroup) {
	ctx := context.Background()
	require.NoError(t, firstCommiter.Commit(ctx))
	wg.Wait()
	require.NoError(t, secondCommiter.Commit(ctx))
	bn, err := state.GetLastBatchNumber(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, expectedLastBatchNum, bn)
}
