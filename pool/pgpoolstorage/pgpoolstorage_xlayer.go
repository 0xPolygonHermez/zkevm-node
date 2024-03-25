package pgpoolstorage

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// GetAllAddressesWhitelisted get all addresses whitelisted
func (p *PostgresPoolStorage) GetAllAddressesWhitelisted(ctx context.Context) ([]common.Address, error) {
	sql := `SELECT addr FROM pool.whitelisted`

	rows, err := p.db.Query(ctx, sql)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var addrs []common.Address
	for rows.Next() {
		var addr string
		err := rows.Scan(&addr)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, common.HexToAddress(addr))
	}

	return addrs, nil
}

// CREATE TABLE pool.innertx (
// hash VARCHAR(128) PRIMARY KEY NOT NULL,
// innertx text,
// created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

// AddInnerTx add inner tx
func (p *PostgresPoolStorage) AddInnerTx(ctx context.Context, txHash common.Hash, innerTx []byte) error {
	sql := `INSERT INTO pool.innertx(hash, innertx) VALUES ($1, $2)`

	_, err := p.db.Exec(ctx, sql, txHash.Hex(), innerTx)
	if err != nil {
		return err
	}

	return nil
}

// GetInnerTx get inner tx
func (p *PostgresPoolStorage) GetInnerTx(ctx context.Context, txHash common.Hash) (string, error) {
	sql := `SELECT innertx FROM pool.innertx WHERE hash = $1`

	var innerTx string
	err := p.db.QueryRow(ctx, sql, txHash.Hex()).Scan(&innerTx)
	if err != nil {
		return "", err
	}

	return innerTx, nil
}
