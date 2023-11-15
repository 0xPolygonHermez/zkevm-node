package pgstatestorage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// GetL2BlockByNumber gets a l2 block by its number
func (p *PostgresStorage) GetL2BlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*types.Block, error) {
	const query = "SELECT header, uncles, received_at FROM state.l2block b WHERE b.block_num = $1"

	q := p.getExecQuerier(dbTx)
	row := q.QueryRow(ctx, query, blockNumber)
	header, uncles, receivedAt, err := p.scanL2BlockInfo(ctx, row, dbTx)
	if err != nil {
		return nil, err
	}

	transactions, err := p.GetTxsByBlockNumber(ctx, header.Number.Uint64(), dbTx)
	if errors.Is(err, pgx.ErrNoRows) {
		transactions = []*types.Transaction{}
	} else if err != nil {
		return nil, err
	}

	block := types.NewBlockWithHeader(header).WithBody(transactions, uncles)
	block.ReceivedAt = receivedAt

	return block, nil
}

// GetL2BlocksByBatchNumber get all blocks associated to a batch
// accordingly to the provided batch number
func (p *PostgresStorage) GetL2BlocksByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]types.Block, error) {
	const query = `
        SELECT bl.header, bl.uncles, bl.received_at
          FROM state.l2block bl
		 INNER JOIN state.batch ba
		    ON ba.batch_num = bl.batch_num
         WHERE ba.batch_num = $1`

	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, query, batchNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	type l2BlockInfo struct {
		header     *types.Header
		uncles     []*types.Header
		receivedAt time.Time
	}

	l2BlockInfos := []l2BlockInfo{}
	for rows.Next() {
		header, uncles, receivedAt, err := p.scanL2BlockInfo(ctx, rows, dbTx)
		if err != nil {
			return nil, err
		}
		l2BlockInfos = append(l2BlockInfos, l2BlockInfo{
			header:     header,
			uncles:     uncles,
			receivedAt: receivedAt,
		})
	}

	l2Blocks := make([]types.Block, 0, len(rows.RawValues()))
	for _, l2BlockInfo := range l2BlockInfos {
		transactions, err := p.GetTxsByBlockNumber(ctx, l2BlockInfo.header.Number.Uint64(), dbTx)
		if errors.Is(err, pgx.ErrNoRows) {
			transactions = []*types.Transaction{}
		} else if err != nil {
			return nil, err
		}

		block := types.NewBlockWithHeader(l2BlockInfo.header).WithBody(transactions, l2BlockInfo.uncles)
		block.ReceivedAt = l2BlockInfo.receivedAt

		l2Blocks = append(l2Blocks, *block)
	}

	return l2Blocks, nil
}

func (p *PostgresStorage) scanL2BlockInfo(ctx context.Context, rows pgx.Row, dbTx pgx.Tx) (header *types.Header, uncles []*types.Header, receivedAt time.Time, err error) {
	header = &types.Header{}
	uncles = []*types.Header{}
	receivedAt = time.Time{}

	err = rows.Scan(&header, &uncles, &receivedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, time.Time{}, state.ErrNotFound
	} else if err != nil {
		return nil, nil, time.Time{}, err
	}

	return header, uncles, receivedAt, nil
}

// GetLastL2BlockCreatedAt gets the timestamp of the last l2 block
func (p *PostgresStorage) GetLastL2BlockCreatedAt(ctx context.Context, dbTx pgx.Tx) (*time.Time, error) {
	var createdAt time.Time
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, "SELECT created_at FROM state.l2block b order by b.block_num desc LIMIT 1").Scan(&createdAt)
	if err != nil {
		return nil, err
	}
	return &createdAt, nil
}

// GetL2BlockTransactionCountByHash returns the number of transactions related to the provided block hash
func (p *PostgresStorage) GetL2BlockTransactionCountByHash(ctx context.Context, blockHash common.Hash, dbTx pgx.Tx) (uint64, error) {
	var count uint64
	const getL2BlockTransactionCountByHashSQL = "SELECT COUNT(*) FROM state.transaction t INNER JOIN state.l2block b ON b.block_num = t.l2_block_num WHERE b.block_hash = $1"

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getL2BlockTransactionCountByHashSQL, blockHash.String()).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetL2BlockTransactionCountByNumber returns the number of transactions related to the provided block number
func (p *PostgresStorage) GetL2BlockTransactionCountByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	var count uint64
	const getL2BlockTransactionCountByNumberSQL = "SELECT COUNT(*) FROM state.transaction t WHERE t.l2_block_num = $1"

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getL2BlockTransactionCountByNumberSQL, blockNumber).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// AddL2Block adds a new L2 block to the State Store
func (p *PostgresStorage) AddL2Block(ctx context.Context, batchNumber uint64, l2Block *types.Block, receipts []*types.Receipt, txsEGPData []state.StoreTxEGPData, dbTx pgx.Tx) error {
	log.Debugf("[AddL2Block] adding l2 block: %v", l2Block.NumberU64())
	start := time.Now()

	e := p.getExecQuerier(dbTx)

	const addTransactionSQL = "INSERT INTO state.transaction (hash, encoded, decoded, l2_block_num, effective_percentage, egp_log) VALUES($1, $2, $3, $4, $5, $6)"
	const addL2BlockSQL = `
        INSERT INTO state.l2block (block_num, block_hash, header, uncles, parent_hash, state_root, received_at, batch_num, created_at)
                           VALUES (       $1,         $2,     $3,     $4,          $5,         $6,          $7,        $8,         $9)`

	var header = "{}"
	if l2Block.Header() != nil {
		headerBytes, err := json.Marshal(l2Block.Header())
		if err != nil {
			return err
		}
		header = string(headerBytes)
	}

	var uncles = "[]"
	if l2Block.Uncles() != nil {
		unclesBytes, err := json.Marshal(l2Block.Uncles())
		if err != nil {
			return err
		}
		uncles = string(unclesBytes)
	}

	if _, err := e.Exec(ctx, addL2BlockSQL,
		l2Block.Number().Uint64(), l2Block.Hash().String(), header, uncles,
		l2Block.ParentHash().String(), l2Block.Root().String(),
		l2Block.ReceivedAt, batchNumber, time.Now().UTC()); err != nil {
		return err
	}

	for idx, tx := range l2Block.Transactions() {
		egpLog := ""
		if txsEGPData != nil {
			egpLogBytes, err := json.Marshal(txsEGPData[idx].EGPLog)
			if err != nil {
				return err
			}
			egpLog = string(egpLogBytes)
		}

		binary, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		encoded := hex.EncodeToHex(binary)

		binary, err = tx.MarshalJSON()
		if err != nil {
			return err
		}
		decoded := string(binary)
		_, err = e.Exec(ctx, addTransactionSQL, tx.Hash().String(), encoded, decoded, l2Block.Number().Uint64(), txsEGPData[idx].EffectivePercentage, egpLog)
		if err != nil {
			return err
		}
	}

	for _, receipt := range receipts {
		err := p.AddReceipt(ctx, receipt, dbTx)
		if err != nil {
			return err
		}

		for _, log := range receipt.Logs {
			err := p.AddLog(ctx, log, dbTx)
			if err != nil {
				return err
			}
		}
	}
	log.Debugf("[AddL2Block] l2 block %v took %v to be added", l2Block.NumberU64(), time.Since(start))
	return nil
}

// GetLastVirtualizedL2BlockNumber gets the last l2 block virtualized
func (p *PostgresStorage) GetLastVirtualizedL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var lastVirtualizedBlockNumber uint64
	const getLastVirtualizedBlockNumberSQL = `
    SELECT b.block_num
      FROM state.l2block b
     INNER JOIN state.virtual_batch vb
        ON vb.batch_num = b.batch_num
     ORDER BY b.block_num DESC LIMIT 1`

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getLastVirtualizedBlockNumberSQL).Scan(&lastVirtualizedBlockNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return lastVirtualizedBlockNumber, nil
}

// GetLastConsolidatedL2BlockNumber gets the last l2 block verified
func (p *PostgresStorage) GetLastConsolidatedL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var lastConsolidatedBlockNumber uint64
	const getLastConsolidatedBlockNumberSQL = `
    SELECT b.block_num
      FROM state.l2block b
     INNER JOIN state.verified_batch vb
        ON vb.batch_num = b.batch_num
     ORDER BY b.block_num DESC LIMIT 1`

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getLastConsolidatedBlockNumberSQL).Scan(&lastConsolidatedBlockNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return lastConsolidatedBlockNumber, nil
}

// GetLastVerifiedL2BlockNumberUntilL1Block gets the last block number that was verified in
// or before the provided l1 block number. This is used to identify if a l2 block is safe or finalized.
func (p *PostgresStorage) GetLastVerifiedL2BlockNumberUntilL1Block(ctx context.Context, l1FinalizedBlockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	var blockNumber uint64
	const query = `
    SELECT b.block_num
      FROM state.l2block b
	 INNER JOIN state.verified_batch vb
	    ON vb.batch_num = b.batch_num
	 WHERE vb.block_num <= $1
     ORDER BY b.block_num DESC LIMIT 1`

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, query, l1FinalizedBlockNumber).Scan(&blockNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

// GetLastL2BlockNumber gets the last l2 block number
func (p *PostgresStorage) GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var lastBlockNumber uint64
	const getLastL2BlockNumber = "SELECT block_num FROM state.l2block ORDER BY block_num DESC LIMIT 1"

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getLastL2BlockNumber).Scan(&lastBlockNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrStateNotSynchronized
	} else if err != nil {
		return 0, err
	}

	return lastBlockNumber, nil
}

// GetLastL2BlockHeader gets the last l2 block number
func (p *PostgresStorage) GetLastL2BlockHeader(ctx context.Context, dbTx pgx.Tx) (*types.Header, error) {
	const query = "SELECT b.header FROM state.l2block b ORDER BY b.block_num DESC LIMIT 1"
	header := &types.Header{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, query).Scan(&header)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	return header, nil
}

// GetLastL2Block retrieves the latest L2 Block from the State data base
func (p *PostgresStorage) GetLastL2Block(ctx context.Context, dbTx pgx.Tx) (*types.Block, error) {
	const query = "SELECT header, uncles, received_at FROM state.l2block b ORDER BY b.block_num DESC LIMIT 1"

	q := p.getExecQuerier(dbTx)
	row := q.QueryRow(ctx, query)
	header, uncles, receivedAt, err := p.scanL2BlockInfo(ctx, row, dbTx)
	if errors.Is(err, state.ErrNotFound) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	transactions, err := p.GetTxsByBlockNumber(ctx, header.Number.Uint64(), dbTx)
	if errors.Is(err, pgx.ErrNoRows) {
		transactions = []*types.Transaction{}
	} else if err != nil {
		return nil, err
	}

	block := types.NewBlockWithHeader(header).WithBody(transactions, uncles)
	block.ReceivedAt = receivedAt

	return block, nil
}

// GetL2BlockByHash gets a l2 block from its hash
func (p *PostgresStorage) GetL2BlockByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (*types.Block, error) {
	const query = "SELECT header, uncles, received_at FROM state.l2block b WHERE b.block_hash = $1"

	q := p.getExecQuerier(dbTx)
	row := q.QueryRow(ctx, query, hash.String())
	header, uncles, receivedAt, err := p.scanL2BlockInfo(ctx, row, dbTx)
	if err != nil {
		return nil, err
	}

	transactions, err := p.GetTxsByBlockNumber(ctx, header.Number.Uint64(), dbTx)
	if errors.Is(err, pgx.ErrNoRows) {
		transactions = []*types.Transaction{}
	} else if err != nil {
		return nil, err
	}

	block := types.NewBlockWithHeader(header).WithBody(transactions, uncles)
	block.ReceivedAt = receivedAt

	return block, nil
}

// GetL2BlockHeaderByHash gets the block header by block number
func (p *PostgresStorage) GetL2BlockHeaderByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (*types.Header, error) {
	const getL2BlockHeaderByHashSQL = "SELECT header FROM state.l2block b WHERE b.block_hash = $1"

	header := &types.Header{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getL2BlockHeaderByHashSQL, hash.String()).Scan(&header)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return header, nil
}

// GetL2BlockHeaderByNumber gets the block header by block number
func (p *PostgresStorage) GetL2BlockHeaderByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*types.Header, error) {
	const getL2BlockHeaderByNumberSQL = "SELECT header FROM state.l2block b WHERE b.block_num = $1"

	header := &types.Header{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getL2BlockHeaderByNumberSQL, blockNumber).Scan(&header)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return header, nil
}

// GetL2BlockHashesSince gets the block hashes added since the provided date
func (p *PostgresStorage) GetL2BlockHashesSince(ctx context.Context, since time.Time, dbTx pgx.Tx) ([]common.Hash, error) {
	const getL2BlockHashesSinceSQL = "SELECT block_hash FROM state.l2block WHERE created_at >= $1"

	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, getL2BlockHashesSinceSQL, since)
	if errors.Is(err, pgx.ErrNoRows) {
		return []common.Hash{}, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	blockHashes := make([]common.Hash, 0, len(rows.RawValues()))

	for rows.Next() {
		var blockHash string
		err := rows.Scan(&blockHash)
		if err != nil {
			return nil, err
		}

		blockHashes = append(blockHashes, common.HexToHash(blockHash))
	}

	return blockHashes, nil
}

// IsL2BlockConsolidated checks if the block ID is consolidated
func (p *PostgresStorage) IsL2BlockConsolidated(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (bool, error) {
	const isL2BlockConsolidated = "SELECT l2b.block_num FROM state.l2block l2b INNER JOIN state.verified_batch vb ON vb.batch_num = l2b.batch_num WHERE l2b.block_num = $1"

	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, isL2BlockConsolidated, blockNumber)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	isConsolidated := rows.Next()

	if rows.Err() != nil {
		return false, rows.Err()
	}

	return isConsolidated, nil
}

// IsL2BlockVirtualized checks if the block  ID is virtualized
func (p *PostgresStorage) IsL2BlockVirtualized(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (bool, error) {
	const isL2BlockVirtualized = "SELECT l2b.block_num FROM state.l2block l2b INNER JOIN state.virtual_batch vb ON vb.batch_num = l2b.batch_num WHERE l2b.block_num = $1"

	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, isL2BlockVirtualized, blockNumber)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	isVirtualized := rows.Next()

	if rows.Err() != nil {
		return false, rows.Err()
	}

	return isVirtualized, nil
}
