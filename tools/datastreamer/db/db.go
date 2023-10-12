package db

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// StateDB implements the StateDB interface
type StateDB struct {
	*pgxpool.Pool
}

// NewStateDB creates a new StateDB
func NewStateDB(db *pgxpool.Pool) *StateDB {
	return &StateDB{
		db,
	}
}

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

// GetGenesisBlock returns the genesis block
func (db *StateDB) GetGenesisBlock(ctx context.Context) (*state.DSL2Block, error) {
	const genesisL2BlockSQL = `SELECT 0 as batch_num, l2b.block_num, l2b.created_at, '0x0000000000000000000000000000000000000000' as global_exit_root, l2b.header->>'miner' AS coinbase, 0 as fork_id, l2b.block_hash, l2b.state_root
							FROM state.l2block l2b
							WHERE l2b.block_num  = 0`

	row := db.QueryRow(ctx, genesisL2BlockSQL)

	l2block, err := scanL2Block(row)
	if err != nil {
		return nil, err
	}

	return l2block, nil
}

// GetL2Blocks returns the L2 blocks
func (db *StateDB) GetL2Blocks(ctx context.Context, limit, offset uint64) ([]*state.DSL2Block, error) {
	const l2BlockSQL = `SELECT l2b.batch_num, l2b.block_num, l2b.received_at, b.global_exit_root, l2b.header->>'miner' AS coinbase, f.fork_id, l2b.block_hash, l2b.state_root
						FROM state.l2block l2b, state.batch b, state.fork_id f
						WHERE l2b.batch_num = b.batch_num AND l2b.batch_num between f.from_batch_num AND f.to_batch_num
						ORDER BY l2b.block_num ASC limit $1 offset $2`

	rows, err := db.Query(ctx, l2BlockSQL, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	l2blocks := make([]*state.DSL2Block, 0, len(rows.RawValues()))

	for rows.Next() {
		l2block, err := scanL2Block(rows)
		if err != nil {
			return nil, err
		}
		l2blocks = append(l2blocks, l2block)
	}

	return l2blocks, nil
}

func scanL2Block(row pgx.Row) (*state.DSL2Block, error) {
	l2Block := state.DSL2Block{}
	var (
		gerStr       string
		coinbaseStr  string
		timestamp    time.Time
		blockHashStr string
		stateRootStr string
	)
	if err := row.Scan(
		&l2Block.BatchNumber,
		&l2Block.L2BlockNumber,
		&timestamp,
		&gerStr,
		&coinbaseStr,
		&l2Block.ForkID,
		&blockHashStr,
		&stateRootStr,
	); err != nil {
		return &l2Block, err
	}
	l2Block.GlobalExitRoot = common.HexToHash(gerStr)
	l2Block.Coinbase = common.HexToAddress(coinbaseStr)
	l2Block.Timestamp = timestamp.Unix()
	l2Block.BlockHash = common.HexToHash(blockHashStr)
	l2Block.StateRoot = common.HexToHash(stateRootStr)

	return &l2Block, nil
}

// GetL2Transactions returns the L2 transactions
func (db *StateDB) GetL2Transactions(ctx context.Context, minL2Block, maxL2Block uint64) ([]*state.DSL2Transaction, error) {
	const l2TxSQL = `SELECT t.effective_percentage, t.encoded
					 FROM state.transaction t
					 WHERE l2_block_num BETWEEN $1 AND $2
					 ORDER BY t.l2_block_num ASC`

	rows, err := db.Query(ctx, l2TxSQL, minL2Block, maxL2Block)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	l2Txs := make([]*state.DSL2Transaction, 0, len(rows.RawValues()))

	for rows.Next() {
		l2Tx, err := scanL2Transaction(rows)
		if err != nil {
			return nil, err
		}
		l2Txs = append(l2Txs, l2Tx)
	}

	return l2Txs, nil
}

func scanL2Transaction(row pgx.Row) (*state.DSL2Transaction, error) {
	l2Transaction := state.DSL2Transaction{}
	encoded := []byte{}
	if err := row.Scan(
		&l2Transaction.EffectiveGasPricePercentage,
		&encoded,
	); err != nil {
		return nil, err
	}
	tx, err := state.DecodeTx(string(encoded))
	if err != nil {
		return nil, err
	}

	binaryTxData, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	l2Transaction.Encoded = binaryTxData
	l2Transaction.EncodedLength = uint32(len(l2Transaction.Encoded))
	l2Transaction.IsValid = 1
	return &l2Transaction, nil
}
