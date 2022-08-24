package pgpoolstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	// ErrNotFound indicates an object has not been found for the search criteria used
	ErrNotFound = errors.New("object not found")
)

// PostgresPoolStorage is an implementation of the Pool interface
// that uses a postgres database to store the data
type PostgresPoolStorage struct {
	db *pgxpool.Pool
}

// NewPostgresPoolStorage creates and initializes an instance of PostgresPoolStorage
func NewPostgresPoolStorage(cfg db.Config) (*PostgresPoolStorage, error) {
	poolDB, err := db.NewSQLDB(cfg)
	if err != nil {
		return nil, err
	}

	return &PostgresPoolStorage{
		db: poolDB,
	}, nil
}

// AddTx adds a transaction to the pool table with the provided state
func (p *PostgresPoolStorage) AddTx(ctx context.Context, tx pool.Transaction) error {
	hash := tx.Hash().Hex()

	b, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	encoded := hex.EncodeToHex(b)

	b, err = tx.MarshalJSON()
	if err != nil {
		return err
	}
	decoded := string(b)

	gasPrice := tx.GasPrice().Uint64()
	nonce := tx.Nonce()
	sql := `
		INSERT INTO pool.txs 
		(
			hash,
			encoded,
			decoded,
			state,
			gas_price,
			nonce,
			is_claims,
			cumulative_gas_used,
			used_keccak_hashes,
			used_poseidon_hashes,
			used_poseidon_paddings,
			used_mem_aligns,
			used_arithmetics,
			used_binaries,
			used_steps,
			received_at,
			from_address
		) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	// Get FromAddress from the JSON data
	data, err := state.GetSender(tx.Transaction)
	if err != nil {
		return err
	}
	fromAddress := data.String()

	if _, err := p.db.Exec(ctx, sql,
		hash,
		encoded,
		decoded,
		tx.State,
		gasPrice,
		nonce,
		tx.IsClaims,
		tx.CumulativeGasUsed,
		tx.UsedKeccakHashes,
		tx.UsedPoseidonHashes,
		tx.UsedPoseidonPaddings,
		tx.UsedMemAligns,
		tx.UsedArithmetics,
		tx.UsedBinaries,
		tx.UsedSteps,
		tx.ReceivedAt,
		fromAddress); err != nil {
		return err
	}
	return nil
}

// MarkReorgedTxsAsPending updated reorged txs state from selected to pending
func (p *PostgresPoolStorage) MarkReorgedTxsAsPending(ctx context.Context) error {
	const updateReorgedTxsToPending = "UPDATE pool.txs pt SET state = $1 WHERE state = $2 AND NOT EXISTS (SELECT hash FROM state.transaction WHERE hash = pt.hash)"
	if _, err := p.db.Exec(ctx, updateReorgedTxsToPending, pool.TxStatePending, pool.TxStateSelected); err != nil {
		return err
	}
	return nil
}

// GetTxsByState returns an array of transactions filtered by state
// limit parameter is used to limit amount txs from the db,
// if limit = 0, then there is no limit
func (p *PostgresPoolStorage) GetTxsByState(ctx context.Context, state pool.TxState, isClaims bool, limit uint64) ([]pool.Transaction, error) {
	var (
		rows pgx.Rows
		err  error
		sql  string
	)
	if limit == 0 {
		sql = "SELECT encoded, state, received_at FROM pool.txs WHERE state = $1 ORDER BY gas_price DESC"
		rows, err = p.db.Query(ctx, sql, state.String())
	} else {
		sql = "SELECT encoded, state, received_at FROM pool.txs WHERE state = $1 AND is_claims = $2 ORDER BY gas_price DESC LIMIT $3"
		rows, err = p.db.Query(ctx, sql, state.String(), isClaims, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]pool.Transaction, 0, len(rows.RawValues()))
	for rows.Next() {
		tx, err := scanTx(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}

	return txs, nil
}

// GetPendingTxHashesSince returns the pending tx since the given time.
func (p *PostgresPoolStorage) GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error) {
	sql := "SELECT hash FROM pool.txs WHERE state = $1 AND received_at >= $2"
	rows, err := p.db.Query(ctx, sql, pool.TxStatePending, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hashes := make([]common.Hash, 0, len(rows.RawValues()))
	for rows.Next() {
		var hash string
		if err := rows.Scan(&hash); err != nil {
			return nil, err
		}
		hashes = append(hashes, common.HexToHash(hash))
	}

	return hashes, nil
}

// GetTopPendingTxByProfitabilityAndZkCounters gets top pending tx by profitability and zk counter
func (p *PostgresPoolStorage) GetTopPendingTxByProfitabilityAndZkCounters(ctx context.Context, maxZkCounters pool.ZkCounters) (*pool.Transaction, error) {
	sql := `
		SELECT 
			encoded, 
			state,
			cumulative_gas_used,
			used_keccak_hashes,
			used_poseidon_hashes,
			used_poseidon_paddings, 
			used_mem_aligns,
			used_arithmetics,
			used_binaries,
			used_steps,
			received_at,
			nonce
		FROM
			pool.txs p1
		WHERE 
			state = $1 AND 
			cumulative_gas_used < $2 AND 
			used_keccak_hashes < $3 AND 
			used_poseidon_hashes < $4 AND 
			used_poseidon_paddings < $5 AND
			used_mem_aligns < $6 AND 
			used_arithmetics < $7 AND
			used_binaries < $8 AND 
			used_steps < $9 AND
			nonce = (
				SELECT MIN(p2.nonce)
				FROM pool.txs p2
				WHERE p1.from_address = p2.from_address AND
				state = $10
			)
		GROUP BY 
			from_address, p1.hash
		ORDER BY
			nonce 
		DESC
		LIMIT 1
	`
	var (
		encoded, state    string
		receivedAt        time.Time
		cumulativeGasUsed int64

		usedKeccakHashes, usedPoseidonHashes, usedPoseidonPaddings,
		usedMemAligns, usedArithmetics, usedBinaries, usedSteps int32
		nonce uint64
	)
	err := p.db.QueryRow(ctx, sql,
		pool.TxStatePending,
		maxZkCounters.CumulativeGasUsed,
		maxZkCounters.UsedKeccakHashes,
		maxZkCounters.UsedPoseidonHashes,
		maxZkCounters.UsedPoseidonPaddings,
		maxZkCounters.UsedMemAligns,
		maxZkCounters.UsedArithmetics,
		maxZkCounters.UsedBinaries,
		maxZkCounters.UsedSteps,
		pool.TxStatePending).
		Scan(&encoded,
			&state,
			&cumulativeGasUsed,
			&usedKeccakHashes,
			&usedPoseidonHashes,
			&usedPoseidonPaddings,
			&usedMemAligns,
			&usedArithmetics,
			&usedBinaries,
			&usedSteps,
			&receivedAt,
			&nonce)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx := new(pool.Transaction)
	b, err := hex.DecodeHex(encoded)
	if err != nil {
		return nil, err
	}
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	tx.State = pool.TxState(state)
	tx.ReceivedAt = receivedAt
	tx.ZkCounters = pool.ZkCounters{
		CumulativeGasUsed:    cumulativeGasUsed,
		UsedKeccakHashes:     usedKeccakHashes,
		UsedPoseidonHashes:   usedPoseidonHashes,
		UsedPoseidonPaddings: usedPoseidonPaddings,
		UsedMemAligns:        usedMemAligns,
		UsedArithmetics:      usedArithmetics,
		UsedBinaries:         usedBinaries,
		UsedSteps:            usedSteps,
	}

	return tx, nil
}

// CountTransactionsByState get number of transactions
// accordingly to the provided state
func (p *PostgresPoolStorage) CountTransactionsByState(ctx context.Context, state pool.TxState) (uint64, error) {
	sql := "SELECT COUNT(*) FROM pool.txs WHERE state = $1"
	var counter uint64
	err := p.db.QueryRow(ctx, sql, state.String()).Scan(&counter)
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// UpdateTxState updates a transaction state accordingly to the
// provided state and hash
func (p *PostgresPoolStorage) UpdateTxState(ctx context.Context, hash common.Hash, newState pool.TxState) error {
	sql := "UPDATE pool.txs SET state = $1 WHERE hash = $2"
	if _, err := p.db.Exec(ctx, sql, newState, hash.Hex()); err != nil {
		return err
	}
	return nil
}

// UpdateTxsState updates transactions state accordingly to the provided state and hashes
func (p *PostgresPoolStorage) UpdateTxsState(ctx context.Context, hashes []common.Hash, newState pool.TxState) error {
	hh := make([]string, 0, len(hashes))
	for _, h := range hashes {
		hh = append(hh, h.Hex())
	}

	sql := "UPDATE pool.txs SET state = $1 WHERE hash = ANY ($2)"
	if _, err := p.db.Exec(ctx, sql, newState, hh); err != nil {
		return err
	}
	return nil
}

// DeleteTxsByHashes deletes txs by their hashes
func (p *PostgresPoolStorage) DeleteTxsByHashes(ctx context.Context, hashes []common.Hash) error {
	hh := make([]string, 0, len(hashes))
	for _, h := range hashes {
		hh = append(hh, h.Hex())
	}

	query := "DELETE FROM pool.txs WHERE hash = ANY ($1)"
	if _, err := p.db.Exec(ctx, query, hh); err != nil {
		return err
	}
	return nil
}

// SetGasPrice allows an external component to define the gas price
func (p *PostgresPoolStorage) SetGasPrice(ctx context.Context, gasPrice uint64) error {
	sql := "INSERT INTO pool.gas_price (price, timestamp) VALUES ($1, $2)"
	if _, err := p.db.Exec(ctx, sql, gasPrice, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

// GetGasPrice returns the current gas price
func (p *PostgresPoolStorage) GetGasPrice(ctx context.Context) (uint64, error) {
	sql := "SELECT price FROM pool.gas_price ORDER BY item_id DESC LIMIT 1"
	rows, err := p.db.Query(ctx, sql)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	defer rows.Close()

	gasPrice := uint64(0)

	for rows.Next() {
		err := rows.Scan(&gasPrice)
		if err != nil {
			return 0, err
		}
	}

	return gasPrice, nil
}

// IsTxPending determines if the tx associated to the given hash is pending or
// not.
func (p *PostgresPoolStorage) IsTxPending(ctx context.Context, hash common.Hash) (bool, error) {
	var exists bool
	req := "SELECT exists (SELECT 1 FROM pool.txs WHERE hash = $1 AND state = $2)"
	err := p.db.QueryRow(ctx, req, hash.Hex(), pool.TxStatePending).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}

// GetTxsByFromAndNonce get all the transactions from the pool with the same from and nonce
func (p *PostgresPoolStorage) GetTxsByFromAndNonce(ctx context.Context, from common.Address, nonce uint64) ([]pool.Transaction, error) {
	sql := `SELECT encoded, state, received_at
	          FROM pool.txs
			 WHERE from_address = $1
			   AND decoded->>'nonce' = $2`
	rows, err := p.db.Query(ctx, sql, from.String(), hex.EncodeUint64(nonce))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]pool.Transaction, 0, len(rows.RawValues()))
	for rows.Next() {
		tx, err := scanTx(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}

	return txs, nil
}

func scanTx(rows pgx.Rows) (*pool.Transaction, error) {
	var (
		encoded, state string
		receivedAt     time.Time
	)

	if err := rows.Scan(&encoded, &state, &receivedAt); err != nil {
		return nil, err
	}

	tx := new(pool.Transaction)

	b, err := hex.DecodeHex(encoded)
	if err != nil {
		return nil, err
	}

	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	tx.State = pool.TxState(state)
	tx.ReceivedAt = receivedAt

	return tx, nil
}
