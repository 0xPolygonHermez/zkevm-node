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
		INSERT INTO pool.transaction 
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
			from_address,
			is_wip,
			ip,
			deposit_count
		) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
			ON CONFLICT (hash) DO UPDATE SET 
			encoded = $2,
			decoded = $3,
			status = $4,
			gas_price = $5,
			nonce = $6,
			is_claims = $7,
			cumulative_gas_used = $8,
			used_keccak_hashes = $9, 
			used_poseidon_hashes = $10,
			used_poseidon_paddings = $11,
			used_mem_aligns = $12,
			used_arithmetics = $13,
			used_binaries = $14,
			used_steps = $15,
			received_at = $16,
			from_address = $17,
			is_wip = $18,
			ip = $19,
			deposit_count = $20
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
		fromAddress,
		tx.IsWIP,
		tx.IP,
		tx.DepositCount); err != nil {
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
		sql = "SELECT encoded, status, received_at, is_wip, ip, deposit_count FROM pool.transaction WHERE status = $1 ORDER BY gas_price DESC"
		rows, err = p.db.Query(ctx, sql, status.String())
	} else {
		sql = "SELECT encoded, status, received_at, is_wip, ip, deposit_count FROM pool.transaction WHERE status = $1 AND is_claims = $2 ORDER BY gas_price DESC LIMIT $3"
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

// GetNonWIPTxsByStatus returns an array of transactions filtered by status
// limit parameter is used to limit amount txs from the db,
// if limit = 0, then there is no limit
func (p *PostgresPoolStorage) GetNonWIPTxsByStatus(ctx context.Context, status pool.TxStatus, isClaims bool, limit uint64) ([]pool.Transaction, error) {
	var (
		rows pgx.Rows
		err  error
		sql  string
	)
	if limit == 0 {
		sql = "SELECT encoded, status, received_at, is_wip, ip, deposit_count FROM pool.transaction WHERE is_wip IS FALSE and status = $1 ORDER BY gas_price DESC"
		rows, err = p.db.Query(ctx, sql, status.String())
	} else {
		sql = "SELECT encoded, status, received_at, is_wip, ip, deposit_count FROM pool.transaction WHERE is_wip IS FALSE and status = $1 AND is_claims = $2 ORDER BY gas_price DESC LIMIT $3"
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
	sql := "SELECT hash FROM pool.transaction WHERE status = $1 AND received_at >= $2"
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
			is_wip,
			ip
		FROM
			pool.transaction p1
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
				is_wip,
				ip
			FROM
				pool.transaction p1
			WHERE
				status = $1 AND
				gas_price >= $2 AND 
				is_claims = $3
			ORDER BY 
				nonce ASC
			LIMIT $4
			) as tmp
		ORDER BY nonce ASC
		`
	}

	var (
		encoded, status, ip string
		receivedAt          time.Time
		cumulativeGasUsed   uint64

		usedKeccakHashes, usedPoseidonHashes, usedPoseidonPaddings,
		usedMemAligns, usedArithmetics, usedBinaries, usedSteps uint32
		nonce uint64
		isWIP bool
	)

	args := []interface{}{filterStatus, minGasPrice, isClaims, limit}

	rows, err := p.db.Query(ctx, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pool.ErrNotFound
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
			&isWIP,
			&ip,
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
		tx.ZKCounters = state.ZKCounters{
			CumulativeGasUsed:    cumulativeGasUsed,
			UsedKeccakHashes:     usedKeccakHashes,
			UsedPoseidonHashes:   usedPoseidonHashes,
			UsedPoseidonPaddings: usedPoseidonPaddings,
			UsedMemAligns:        usedMemAligns,
			UsedArithmetics:      usedArithmetics,
			UsedBinaries:         usedBinaries,
			UsedSteps:            usedSteps,
		}
		tx.IsWIP = isWIP
		tx.IP = ip

		txs = append(txs, tx)
	}

	return txs, nil
}

// CountTransactionsByStatus get number of transactions
// accordingly to the provided status
func (p *PostgresPoolStorage) CountTransactionsByStatus(ctx context.Context, status pool.TxStatus) (uint64, error) {
	sql := "SELECT COUNT(*) FROM pool.transaction WHERE status = $1"
	var counter uint64
	err := p.db.QueryRow(ctx, sql, status.String()).Scan(&counter)
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// UpdateTxStatus updates a transaction status accordingly to the
// provided status and hash
func (p *PostgresPoolStorage) UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus pool.TxStatus, isWIP bool) error {
	sql := "UPDATE pool.transaction SET status = $1, is_wip = $2 WHERE hash = $3"
	if _, err := p.db.Exec(ctx, sql, newStatus, isWIP, hash.Hex()); err != nil {
		return err
	}
	return nil
}

// UpdateTxsStatus updates transactions status accordingly to the provided status and hashes
func (p *PostgresPoolStorage) UpdateTxsStatus(ctx context.Context, hashes []string, newStatus pool.TxStatus) error {
	sql := "UPDATE pool.transaction SET status = $1 WHERE hash = ANY ($2)"
	if _, err := p.db.Exec(ctx, sql, newStatus, hashes); err != nil {
		return err
	}
	return nil
}

// DeleteTransactionsByHashes deletes txs by their hashes
func (p *PostgresPoolStorage) DeleteTransactionsByHashes(ctx context.Context, hashes []common.Hash) error {
	hh := make([]string, 0, len(hashes))
	for _, h := range hashes {
		hh = append(hh, h.Hex())
	}

	query := "DELETE FROM pool.transaction WHERE hash = ANY ($1)"
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

// MinGasPriceSince returns the min gas price after given timestamp
func (p *PostgresPoolStorage) MinGasPriceSince(ctx context.Context, timestamp time.Time) (uint64, error) {
	sql := "SELECT COALESCE(MIN(price), 0) FROM pool.gas_price WHERE \"timestamp\" >= $1 LIMIT 1"
	var gasPrice uint64
	err := p.db.QueryRow(ctx, sql, timestamp).Scan(&gasPrice)
	if gasPrice == 0 || errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return gasPrice, nil
}

// IsTxPending determines if the tx associated to the given hash is pending or
// not.
func (p *PostgresPoolStorage) IsTxPending(ctx context.Context, hash common.Hash) (bool, error) {
	var exists bool
	req := "SELECT exists (SELECT 1 FROM pool.transaction WHERE hash = $1 AND status = $2)"
	err := p.db.QueryRow(ctx, req, hash.Hex(), pool.TxStatusPending).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}

// GetTxsByFromAndNonce get all the transactions from the pool with the same from and nonce
func (p *PostgresPoolStorage) GetTxsByFromAndNonce(ctx context.Context, from common.Address, nonce uint64) ([]pool.Transaction, error) {
	sql := `SELECT encoded, status, received_at, is_wip, ip, deposit_count
	          FROM pool.transaction
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
			  FROM pool.transaction
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

// GetNonce gets the nonce to the provided address accordingly to the txs in the pool
func (p *PostgresPoolStorage) GetNonce(ctx context.Context, address common.Address) (uint64, error) {
	sql := `SELECT MAX(nonce)
              FROM pool.transaction
             WHERE from_address = $1
               AND status IN  ($2, $3)`
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
		encoded, status, ip string
		receivedAt          time.Time
		isWIP               bool
	)

	sql := `SELECT encoded, status, received_at, is_wip, ip
	          FROM pool.transaction
			 WHERE hash = $1`
	err := p.db.QueryRow(ctx, sql, hash.String()).Scan(&encoded, &status, &receivedAt, &isWIP, &ip)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pool.ErrNotFound
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

	poolTx := &pool.Transaction{
		ReceivedAt:  receivedAt,
		Status:      pool.TxStatus(status),
		Transaction: *tx,
		IsWIP:       isWIP,
		IP:          ip,
	}

	return poolTx, nil
}

func scanTx(rows pgx.Rows) (*pool.Transaction, error) {
	var (
		encoded, status, ip string
		receivedAt          time.Time
		isWIP               bool
		depositCount        *uint64
	)

	if err := rows.Scan(&encoded, &status, &receivedAt, &isWIP, &ip, &depositCount); err != nil {
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
	tx.IsWIP = isWIP
	tx.IP = ip
	tx.DepositCount = depositCount

	return tx, nil
}

// DeleteTransactionByHash deletes tx by its hash
func (p *PostgresPoolStorage) DeleteTransactionByHash(ctx context.Context, hash common.Hash) error {
	query := "DELETE FROM pool.transaction WHERE hash = $1"
	if _, err := p.db.Exec(ctx, query, hash); err != nil {
		return err
	}
	return nil
}

// GetTxZkCountersByHash gets a transaction zkcounters by its hash
func (p *PostgresPoolStorage) GetTxZkCountersByHash(ctx context.Context, hash common.Hash) (*state.ZKCounters, error) {
	var zkCounters state.ZKCounters

	sql := `SELECT cumulative_gas_used, used_keccak_hashes, used_poseidon_hashes, used_poseidon_paddings, used_mem_aligns,
			used_arithmetics, used_binaries, used_steps FROM pool.transaction WHERE hash = $1`
	err := p.db.QueryRow(ctx, sql, hash.String()).Scan(&zkCounters.CumulativeGasUsed, &zkCounters.UsedKeccakHashes,
		&zkCounters.UsedPoseidonHashes, &zkCounters.UsedPoseidonPaddings,
		&zkCounters.UsedMemAligns, &zkCounters.UsedArithmetics, &zkCounters.UsedBinaries, &zkCounters.UsedSteps)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pool.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &zkCounters, nil
}

// MarkWIPTxsAsPending updates WIP status to non WIP
func (p *PostgresPoolStorage) MarkWIPTxsAsPending(ctx context.Context) error {
	const query = `UPDATE pool.transaction SET is_wip = false WHERE is_wip = true`
	if _, err := p.db.Exec(ctx, query); err != nil {
		return err
	}
	return nil
}

// UpdateTxWIPStatus updates a transaction wip status accordingly to the
// provided WIP status and hash
func (p *PostgresPoolStorage) UpdateTxWIPStatus(ctx context.Context, hash common.Hash, isWIP bool) error {
	sql := "UPDATE pool.transaction SET is_wip = $1 WHERE hash = $2"
	if _, err := p.db.Exec(ctx, sql, isWIP, hash.Hex()); err != nil {
		return err
	}
	return nil
}

// DepositCountExists checks if already exists a transaction in the pool with the
// provided deposit count
func (p *PostgresPoolStorage) DepositCountExists(ctx context.Context, depositCount uint64) (bool, error) {
	var exists bool
	req := "SELECT exists (SELECT 1 FROM pool.transaction WHERE deposit_count = $1)"
	err := p.db.QueryRow(ctx, req, depositCount).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}
