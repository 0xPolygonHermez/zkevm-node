package pool

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	database "github.com/hermeznetwork/hermez-core/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresPool struct {
	db *pgxpool.Pool
}

type decodedTx struct {
	nonce uint64 `json:""`
}

func newPostgresPool(cfg Config) (*postgresPool, error) {
	sqlDB, err := database.NewSQLDB(cfg.Database, cfg.User, cfg.Password, cfg.Host, cfg.Port)
	if err != nil {
		return nil, err
	}

	return &postgresPool{
		db: sqlDB,
	}, nil
}

func (p *postgresPool) AddTx(ctx context.Context, tx types.Transaction) error {
	// hash
	hash := tx.Hash()

	// encoded
	sw := &strings.Builder{}
	if err := tx.EncodeRLP(sw); err != nil {
		return err
	}
	encoded := sw.String()

	// decoded
	decodedTx := decodeTx(tx)
	decoded, err := json.Marshal(decodedTx)
	if err != nil {
		return err
	}

	// state
	state := TxStatePending

	// save
	sql := "INSERT INTO pool.txs(hash, encoded, decoded, state) VALUES(?,?,?,?)"
	if _, err := p.db.Exec(ctx, sql, hash, encoded, decoded, state); err != nil {
		return err
	}
	return nil
}

func (p *postgresPool) GetPendingTxs(ctx context.Context) ([]Transaction, error) {
	panic("not implemented yet")
}

func (p *postgresPool) UpdateTxState(ctx context.Context, hash common.Hash, newState TxState) error {
	panic("not implemented yet")
}

func (p *postgresPool) CleanUpInvalidAndNonSelectedTxs(ctx context.Context) error {
	panic("not implemented yet")
}

func (p *postgresPool) GetGasPrice(ctx context.Context) (uint64, error) {
	panic("not implemented yet")
}

func decodeTx(tx types.Transaction) decodedTx {
	return decodedTx{
		nonce: tx.Nonce(),
	}
}
