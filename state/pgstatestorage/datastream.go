package pgstatestorage

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// GetDSGenesisBlock returns the genesis block
func (p *PostgresStorage) GetDSGenesisBlock(ctx context.Context, dbTx pgx.Tx) (*state.DSL2Block, error) {
	const genesisL2BlockSQL = `SELECT 0 as batch_num, l2b.block_num, l2b.received_at, '0x0000000000000000000000000000000000000000' as global_exit_root, l2b.header->>'miner' AS coinbase, 0 as fork_id, l2b.block_hash, l2b.state_root
							FROM state.l2block l2b
							WHERE l2b.block_num  = 0`

	e := p.getExecQuerier(dbTx)

	row := e.QueryRow(ctx, genesisL2BlockSQL)

	l2block, err := scanL2Block(row)
	if err != nil {
		return nil, err
	}

	return l2block, nil
}

// GetDSL2Blocks returns the L2 blocks
func (p *PostgresStorage) GetDSL2Blocks(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]*state.DSL2Block, error) {
	const l2BlockSQL = `SELECT l2b.batch_num, l2b.block_num, l2b.received_at, b.global_exit_root, l2b.header->>'miner' AS coinbase, f.fork_id, l2b.block_hash, l2b.state_root
						FROM state.l2block l2b, state.batch b, state.fork_id f
						WHERE l2b.batch_num = $1 AND l2b.batch_num = b.batch_num AND l2b.batch_num between f.from_batch_num AND f.to_batch_num
						ORDER BY l2b.block_num ASC`
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, l2BlockSQL, batchNumber)
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

// GetDSL2Transactions returns the L2 transactions
func (p *PostgresStorage) GetDSL2Transactions(ctx context.Context, minL2Block, maxL2Block uint64, dbTx pgx.Tx) ([]*state.DSL2Transaction, error) {
	const l2TxSQL = `SELECT t.effective_percentage, t.encoded
					 FROM state.transaction t
					 WHERE l2_block_num BETWEEN $1 AND $2
					 ORDER BY t.l2_block_num ASC`

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, l2TxSQL, minL2Block, maxL2Block)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	l2Txs := make([]*state.DSL2Transaction, 0, len(rows.RawValues()))

	for rows.Next() {
		l2Tx, err := scanDSL2Transaction(rows)
		if err != nil {
			return nil, err
		}
		l2Txs = append(l2Txs, l2Tx)
	}

	return l2Txs, nil
}

func scanDSL2Transaction(row pgx.Row) (*state.DSL2Transaction, error) {
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

// GetDSBatch returns the batch with the given number in DS format
func (p *PostgresStorage) GetDSBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.DSBatch, error) {
	const getBatchByNumberSQL = `
		SELECT b.batch_num, b.global_exit_root, b.local_exit_root, b.acc_input_hash, b.state_root, b.timestamp, b.coinbase, b.raw_txs_data, b.forced_batch_num, f.fork_id
		  FROM state.batch b, state.fork_id f
		 WHERE b.state_root is not null AND batch_num = $1 AND batch_num between f.from_batch_num AND f.to_batch_num`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getBatchByNumberSQL, batchNumber)
	batch, err := scanDSBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	return &batch, nil
}

func scanDSBatch(row pgx.Row) (state.DSBatch, error) {
	batch := state.DSBatch{}
	var (
		gerStr      string
		lerStr      *string
		aihStr      *string
		stateStr    *string
		coinbaseStr string
	)
	err := row.Scan(
		&batch.BatchNumber,
		&gerStr,
		&lerStr,
		&aihStr,
		&stateStr,
		&batch.Timestamp,
		&coinbaseStr,
		&batch.BatchL2Data,
		&batch.ForcedBatchNum,
		&batch.ForkID,
	)
	if err != nil {
		return batch, err
	}
	batch.GlobalExitRoot = common.HexToHash(gerStr)
	if lerStr != nil {
		batch.LocalExitRoot = common.HexToHash(*lerStr)
	}
	if stateStr != nil {
		batch.StateRoot = common.HexToHash(*stateStr)
	}
	if aihStr != nil {
		batch.AccInputHash = common.HexToHash(*aihStr)
	}

	batch.Coinbase = common.HexToAddress(coinbaseStr)
	return batch, nil
}
