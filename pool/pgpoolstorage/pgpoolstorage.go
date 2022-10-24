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
	"github.com/ethereum/go-ethereum/core/types"
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

// AddTx adds a transaction to the pool table with the provided status
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
			status,
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
		tx.Status,
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

// GetTxsByStatus returns an array of transactions filtered by status
// limit parameter is used to limit amount txs from the db,
// if limit = 0, then there is no limit
func (p *PostgresPoolStorage) GetTxsByStatus(ctx context.Context, status pool.TxStatus, isClaims bool, limit uint64) ([]pool.Transaction, error) {
	var (
		rows pgx.Rows
		err  error
		sql  string
	)
	if limit == 0 {
		sql = "SELECT encoded, status, received_at FROM pool.txs WHERE status = $1 ORDER BY gas_price DESC"
		rows, err = p.db.Query(ctx, sql, status.String())
	} else {
		sql = "SELECT encoded, status, received_at FROM pool.txs WHERE status = $1 AND is_claims = $2 ORDER BY gas_price DESC LIMIT $3"
		rows, err = p.db.Query(ctx, sql, status.String(), isClaims, limit)
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
	sql := "SELECT hash FROM pool.txs WHERE status = $1 AND received_at >= $2"
	rows, err := p.db.Query(ctx, sql, pool.TxStatusPending, since)
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

// GetTxs gets txs with the lowest nonce
func (p *PostgresPoolStorage) GetTxs(ctx context.Context, filterStatus pool.TxStatus, isClaims bool, minGasPrice, limit uint64) ([]*pool.Transaction, error) {
	query := `
		SELECT
			encoded,
			status,
			cumulative_gas_used,
			used_keccak_hashes,
			used_poseidon_hashes,
			used_poseidon_paddings,
			used_mem_aligns,
			used_arithmetics,
			used_binaries,
			used_steps,
			received_at,
			nonce,
			failed_counter
		FROM
			pool.txs p1
		WHERE 
			status = $1 AND
			gas_price >= $2 AND
			is_claims = $3
		ORDER BY 
			nonce ASC
		LIMIT $4
	`

	if filterStatus == pool.TxStatusFailed {
		query = `
		SELECT * FROM (
			SELECT
				encoded,
				status,
				cumulative_gas_used,
				used_keccak_hashes,
				used_poseidon_hashes,
				used_poseidon_paddings,
				used_mem_aligns,
				used_arithmetics,
				used_binaries,
				used_steps,
				received_at,
				nonce,
				failed_counter
			FROM
				pool.txs p1
			WHERE
				status = $1 AND
				gas_price >= $2 AND 
				is_claims = $3
			ORDER BY 
				failed_counter ASC
			LIMIT $4
			) as tmp
		ORDER BY nonce ASC
		`
	}

	var (
		encoded, status   string
		receivedAt        time.Time
		cumulativeGasUsed int64

		usedKeccakHashes, usedPoseidonHashes, usedPoseidonPaddings,
		usedMemAligns, usedArithmetics, usedBinaries, usedSteps int32
		nonce, failedCounter uint64
	)

	args := []interface{}{filterStatus, minGasPrice, isClaims, limit}

	rows, err := p.db.Query(ctx, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	txs := make([]*pool.Transaction, 0, len(rows.RawValues()))
	for rows.Next() {
		err := rows.Scan(
			&encoded,
			&status,
			&cumulativeGasUsed,
			&usedKeccakHashes,
			&usedPoseidonHashes,
			&usedPoseidonPaddings,
			&usedMemAligns,
			&usedArithmetics,
			&usedBinaries,
			&usedSteps,
			&receivedAt,
			&nonce,
			&failedCounter,
		)

		if err != nil {
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

		tx.Status = pool.TxStatus(status)
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
		tx.FailedCounter = failedCounter

		txs = append(txs, tx)
	}

	return txs, nil
}

// CountTransactionsByStatus get number of transactions
// accordingly to the provided status
func (p *PostgresPoolStorage) CountTransactionsByStatus(ctx context.Context, status pool.TxStatus) (uint64, error) {
	sql := "SELECT COUNT(*) FROM pool.txs WHERE status = $1"
	var counter uint64
	err := p.db.QueryRow(ctx, sql, status.String()).Scan(&counter)
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// UpdateTxStatus updates a transaction status accordingly to the
// provided status and hash
func (p *PostgresPoolStorage) UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus pool.TxStatus) error {
	sql := "UPDATE pool.txs SET status = $1 WHERE hash = $2"
	if _, err := p.db.Exec(ctx, sql, newStatus, hash.Hex()); err != nil {
		return err
	}
	return nil
}

// UpdateTxsStatus updates transactions status accordingly to the provided status and hashes
func (p *PostgresPoolStorage) UpdateTxsStatus(ctx context.Context, hashes []string, newStatus pool.TxStatus) error {
	sql := "UPDATE pool.txs SET status = $1 WHERE hash = ANY ($2)"
	if _, err := p.db.Exec(ctx, sql, newStatus, hashes); err != nil {
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
	req := "SELECT exists (SELECT 1 FROM pool.txs WHERE hash = $1 AND status = $2)"
	err := p.db.QueryRow(ctx, req, hash.Hex(), pool.TxStatusPending).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}

// GetTxsByFromAndNonce get all the transactions from the pool with the same from and nonce
func (p *PostgresPoolStorage) GetTxsByFromAndNonce(ctx context.Context, from common.Address, nonce uint64) ([]pool.Transaction, error) {
	sql := `SELECT encoded, status, received_at
	          FROM pool.txs
			 WHERE from_address = $1
			   AND nonce = $2`
	rows, err := p.db.Query(ctx, sql, from.String(), nonce)
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

// GetTxFromAddressFromByHash gets tx from address by hash
func (p *PostgresPoolStorage) GetTxFromAddressFromByHash(ctx context.Context, hash common.Hash) (common.Address, uint64, error) {
	query := `SELECT from_address, nonce
			  FROM pool.txs
			  WHERE hash = $1
	`

	var (
		fromAddr string
		nonce    uint64
	)
	err := p.db.QueryRow(ctx, query, hash.String()).Scan(&fromAddr, &nonce)
	if err != nil {
		return common.Address{}, 0, err
	}

	return common.HexToAddress(fromAddr), nonce, nil
}

// IncrementFailedCounter increment for failed txs failed counter
func (p *PostgresPoolStorage) IncrementFailedCounter(ctx context.Context, hashes []string) error {
	sql := "UPDATE pool.txs SET failed_counter = failed_counter + 1 WHERE hash = ANY ($1)"
	if _, err := p.db.Exec(ctx, sql, hashes); err != nil {
		return err
	}
	return nil
}

// GetNonce gets the nonce to the provided address accordingly to the txs in the pool
func (p *PostgresPoolStorage) GetNonce(ctx context.Context, address common.Address) (uint64, error) {
	sql := `SELECT MAX(nonce)
              FROM pool.txs
             WHERE from_address = $1
               AND (status = $2 OR status = $3)`
	rows, err := p.db.Query(ctx, sql, address.String(), pool.TxStatusPending, pool.TxStatusSelected)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	defer rows.Close()

	var nonce *uint64
	for rows.Next() {
		err := rows.Scan(&nonce)
		if err != nil {
			return 0, err
		} else if rows.Err() != nil {
			return 0, rows.Err()
		}
	}

	if nonce == nil {
		n := uint64(0)
		nonce = &n
	} else {
		n := *nonce + 1
		nonce = &n
	}

	return *nonce, nil
}

// GetTxByHash gets a transaction in the pool by its hash
func (p *PostgresPoolStorage) GetTxByHash(ctx context.Context, hash common.Hash) (*pool.Transaction, error) {
	var (
		encoded, status string
		receivedAt      time.Time
	)

	sql := `SELECT encoded, status, received_at
	          FROM pool.txs
			 WHERE hash = $1`
	err := p.db.QueryRow(ctx, sql, hash.String()).Scan(&encoded, &status, &receivedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	b, err := hex.DecodeHex(encoded)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return &pool.Transaction{
		ReceivedAt:  receivedAt,
		Status:      pool.TxStatus(status),
		Transaction: *tx,
	}, nil
}

func scanTx(rows pgx.Rows) (*pool.Transaction, error) {
	var (
		encoded, status string
		receivedAt      time.Time
	)

	if err := rows.Scan(&encoded, &status, &receivedAt); err != nil {
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

	tx.Status = pool.TxStatus(status)
	tx.ReceivedAt = receivedAt

	return tx, nil
}
