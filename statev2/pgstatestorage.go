package statev2

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const maxLogTopics = 4

const (
	addGlobalExitRootSQL                   = "INSERT INTO statev2.exit_root (block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root) VALUES ($1, $2, $3, $4, $5)"
	getLatestExitRootSQL                   = "SELECT block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root, global_exit_root FROM statev2.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
	getLatestExitRootBlockNumSQL           = "SELECT block_num FROM statev2.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
	getForcedBatchSQL                      = "SELECT block_num, batch_num, forced_batch_num, global_exit_root, timestamp, raw_txs_data, sequencer FROM statev2.forced_batch WHERE forced_batch_num = $1"
	addBlockSQL                            = "INSERT INTO statev2.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)"
	resetSQL                               = "DELETE FROM statev2.block WHERE block_num > $1"
	resetTrustedStateSQL                   = "DELETE FROM statev2.batch WHERE batch_num > $1"
	addVerifiedBatchSQL                    = "INSERT INTO statev2.verified_batch (block_num, batch_num, tx_hash, aggregator) VALUES ($1, $2, $3, $4)"
	getVerifiedBatchSQL                    = "SELECT block_num, batch_num, tx_hash, aggregator FROM statev2.verified_batch WHERE batch_num = $1"
	getLastNBatchesSQL                     = "SELECT batch_num, global_exit_root, timestamp from statev2.batch ORDER BY batch_num DESC LIMIT $1"
	getLastBatchTimeSQL                    = "SELECT timestamp FROM statev2.batch ORDER BY batch_num DESC LIMIT 1"
	getLastVirtualBatchNumSQL              = "SELECT batch_num FROM statev2.virtual_batch ORDER BY batch_num DESC LIMIT 1"
	getLastVirtualBatchBlockNumSQL         = "SELECT block_num FROM statev2.virtual_batch ORDER BY batch_num DESC LIMIT 1"
	getLastBlockNumSQL                     = "SELECT block_num FROM statev2.block ORDER BY block_num DESC LIMIT 1"
	getBlockTimeByNumSQL                   = "SELECT received_at FROM statev2.block WHERE block_num = $1"
	getBatchByNumberSQL                    = "SELECT batch_num, global_exit_root, timestamp from statev2.batch WHERE batch_num = $1"
	getEncodedTransactionsByBatchNumberSQL = "SELECT encoded from statev2.transaction where batch_num = $1"
	getLastBatchSeenSQL                    = "SELECT last_batch_num_seen FROM statev2.sync_info LIMIT 1"
	updateLastBatchSeenSQL                 = "UPDATE statev2.sync_info SET last_batch_num_seen = $1"
	getL2BlockByNumberSQL                  = "SELECT l2_block_num, encoded, header, uncles, received_at from statev2.transaction WHERE batch_num = $1"
)

// PostgresStorage implements the Storage interface
type PostgresStorage struct {
	*pgxpool.Pool
}

// NewPostgresStorage creates a new StateDB
func NewPostgresStorage(db *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{
		db,
	}
}

// getExecQuerier determines which execQuerier to use, dbTx or the main pgxpool
func (p *PostgresStorage) getExecQuerier(dbTx pgx.Tx) execQuerier {
	if dbTx != nil {
		return dbTx
	}
	return p
}

// Reset resets the state to a block
func (p *PostgresStorage) Reset(ctx context.Context, block *Block, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	if _, err := e.Exec(ctx, resetSQL, block.BlockNumber); err != nil {
		return err
	}
	// TODO: Remove consolidations
	return nil
}

func (p *PostgresStorage) ResetTrustedState(ctx context.Context, batchNum uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	if _, err := e.Exec(ctx, resetTrustedStateSQL, batchNum); err != nil {
		return err
	}
	// TODO: Remove consolidations
	return nil
}

// AddBlock adds a new block to the State Store
func (p *PostgresStorage) AddBlock(ctx context.Context, block *Block, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.String(), block.ParentHash.String(), block.ReceivedAt)
	return err
}

// AddGlobalExitRoot adds a new ExitRoot to the db
func (p *PostgresStorage) AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addGlobalExitRootSQL, exitRoot.BlockNumber, exitRoot.GlobalExitRootNum.String(), exitRoot.MainnetExitRoot, exitRoot.RollupExitRoot, exitRoot.GlobalExitRoot)
	return err
}

// GetLatestExitRoot get the latest ExitRoot synced.
func (p *PostgresStorage) GetLatestGlobalExitRoot(ctx context.Context, dbTx pgx.Tx) (*GlobalExitRoot, error) {
	var (
		exitRoot  GlobalExitRoot
		globalNum uint64
		err       error
	)

	e := p.getExecQuerier(dbTx)
	err = e.QueryRow(ctx, getLatestExitRootSQL).Scan(&exitRoot.BlockNumber, &globalNum, &exitRoot.MainnetExitRoot, &exitRoot.RollupExitRoot, &exitRoot.GlobalExitRoot)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	exitRoot.GlobalExitRootNum = new(big.Int).SetUint64(globalNum)
	return &exitRoot, nil
}

// GetNumberOfBlocksSinceLastGERUpdate gets number of blocks since last global exit root update
func (p *PostgresStorage) GetNumberOfBlocksSinceLastGERUpdate(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var (
		lastBlockNum         uint64
		lastExitRootBlockNum uint64
		err                  error
	)
	e := p.getExecQuerier(dbTx)
	err = e.QueryRow(ctx, getLastBlockNumSQL).Scan(&lastBlockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	err = p.QueryRow(ctx, getLatestExitRootBlockNumSQL).Scan(&lastExitRootBlockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return lastBlockNum - lastExitRootBlockNum, nil
}

func (p *PostgresStorage) GetLastSendSequenceTime(ctx context.Context, dbTx pgx.Tx) (time.Time, error) {
	var (
		blockNum  uint64
		timestamp time.Time
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastVirtualBatchBlockNumSQL).Scan(&blockNum)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrNotFound
	} else if err != nil {
		return time.Time{}, err
	}

	err = p.QueryRow(ctx, getBlockTimeByNumSQL, blockNum).Scan(&timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrNotFound
	} else if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

// GetForcedBatch get an L1 forcedBatch.
func (p *PostgresStorage) GetForcedBatch(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (*ForcedBatch, error) {
	var (
		forcedBatch    ForcedBatch
		globalExitRoot string
		rawTxs         string
		seq            string
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getForcedBatchSQL, forcedBatchNumber).Scan(&forcedBatch.BlockNumber, &forcedBatch.ForcedBatchNumber, &globalExitRoot, &forcedBatch.ForcedAt, &rawTxs, &seq)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	forcedBatch.RawTxsData, err = hex.DecodeString(rawTxs)
	if err != nil {
		return nil, err
	}
	forcedBatch.Sequencer = common.HexToAddress(seq)
	forcedBatch.GlobalExitRoot = common.HexToHash(globalExitRoot)
	return &forcedBatch, nil
}

// AddVerifiedBatch adds a new VerifiedBatch to the db
func (p *PostgresStorage) AddVerifiedBatch(ctx context.Context, verifiedBatch *VerifiedBatch, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addVerifiedBatchSQL, verifiedBatch.BlockNumber, verifiedBatch.BatchNumber, verifiedBatch.TxHash.String(), verifiedBatch.Aggregator.String())
	return err
}

// GetVerifiedBatch get an L1 verifiedBatch.
func (p *PostgresStorage) GetVerifiedBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*VerifiedBatch, error) {
	var (
		verifiedBatch VerifiedBatch
		txHash        string
		agg           string
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getVerifiedBatchSQL, batchNumber).Scan(&verifiedBatch.BlockNumber, &verifiedBatch.BatchNumber, &txHash, &agg)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	verifiedBatch.Aggregator = common.HexToAddress(agg)
	verifiedBatch.TxHash = common.HexToHash(txHash)
	return &verifiedBatch, nil
}

func (p *PostgresStorage) GetLastNBatches(ctx context.Context, numBatches uint, dbTx pgx.Tx) ([]*Batch, error) {
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getLastNBatchesSQL, numBatches)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	batches := make([]*Batch, 0, len(rows.RawValues()))

	for rows.Next() {
		var (
			batch  Batch
			gerStr string
		)
		err := rows.Scan(&batch.BatchNumber, &gerStr, &batch.Timestamp)
		if err != nil {
			return nil, err
		}

		batch.GlobalExitRootNum = new(big.Int).SetBytes(common.FromHex(gerStr))
		batches = append(batches, &batch)
	}

	return batches, nil
}

// GetLastBatchTime gets last trusted batch time
func (p *PostgresStorage) GetLastBatchTime(ctx context.Context, dbTx pgx.Tx) (time.Time, error) {
	var timestamp time.Time
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastBatchTimeSQL).Scan(&timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrStateNotSynchronized
	} else if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

// GetLastVirtualBatchNum gets last virtual batch num
func (p *PostgresStorage) GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var batchNum uint64
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastVirtualBatchNumSQL).Scan(&batchNum)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}
	return batchNum, nil
}

// SetLastBatchNumberSeenOnEthereum sets the last batch number that affected
// the roll-up in order to allow the components to know if the state
// is synchronized or not
func (p *PostgresStorage) SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, updateLastBatchSeenSQL, batchNumber)
	return err
}

// GetLastBatchNumberSeenOnEthereum returns the last batch number stored
// in the state that represents the last batch number that affected the
// roll-up in the Ethereum network.
func (p *PostgresStorage) GetLastBatchNumberSeenOnEthereum(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var batchNumber uint64
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastBatchSeenSQL).Scan(&batchNumber)

	if err != nil {
		return 0, err
	}

	return batchNumber, nil
}

func (p *PostgresStorage) GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*Batch, error) {
	var (
		batch  Batch
		gerStr string
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getBatchByNumberSQL, batchNumber).Scan(&batch.BatchNumber, &gerStr, &batch.Timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	batch.GlobalExitRootNum = new(big.Int).SetBytes(common.FromHex(gerStr))
	return &batch, nil
}

func (p *PostgresStorage) GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (encoded []string, err error) {
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getEncodedTransactionsByBatchNumberSQL, batchNumber)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]string, 0, len(rows.RawValues()))

	for rows.Next() {
		var encoded string
		err := rows.Scan(&encoded)
		if err != nil {
			return nil, err
		}

		txs = append(txs, encoded)
	}
	return txs, nil
}

func (p *PostgresStorage) GetL2BlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*L2Block, error) {
	var block L2Block
	var encoded string
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getL2BlockByNumberSQL, blockNumber).Scan(&block.BlockNumber, &encoded, &block.Header, &block.Uncles, &block.ReceivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	block.Transactions = make([]*types.Transaction, 1)

	tx, err := decodeTx(encoded)
	if err != nil {
		return nil, err
	}

	block.Transactions[0] = tx

	return &block, nil
}

func (p *PostgresStorage) GetTransactionByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error) {
	const getTransactionByHashSQL = `
		SELECT transaction.encoded 
		  FROM statev2.transaction
		 WHERE hash = $1
	`

	var encoded string
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByHashSQL, transactionHash).Scan(&encoded)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := decodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *PostgresStorage) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Receipt, error) {
	const getReceiptSQL = `
		SELECT r.tx_hash
			 , r.type
			 , r.post_state
			 , r.status
			 , r.cumulative_gas_used
			 , r.gas_used
			 , r.contract_address
			 , t.encoded
			 , t.l2_block_num
			 , b.block_hash
		  FROM statev2.receipt r
		 INNER JOIN statev2.transaction t
		    ON t.hash = r.tx_hash
		 INNER JOIN statev2.l2block b
			ON b.block_num = t.l2_block_num
		 WHERE r.tx_hash = $1
	`

	var encodedTx string
	var l2BlockNum uint64
	var l2BlockHash string

	receipt := types.Receipt{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getReceiptSQL, transactionHash).
		Scan(&receipt.TxHash,
			&receipt.Type,
			&receipt.PostState,
			&receipt.Status,
			&receipt.CumulativeGasUsed,
			&receipt.GasUsed,
			&receipt.ContractAddress,
			&encodedTx,
			&l2BlockNum,
			&l2BlockHash,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	logs, err := p.getTransactionLogs(ctx, transactionHash, dbTx)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}

	receipt.BlockNumber = big.NewInt(0).SetUint64(l2BlockNum)
	receipt.BlockHash = common.HexToHash(l2BlockHash)
	receipt.TransactionIndex = 0

	receipt.Logs = logs
	receipt.Bloom = types.CreateBloom(types.Receipts{&receipt})

	return &receipt, nil
}

func (p *PostgresStorage) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index uint64, dbTx pgx.Tx) (*types.Transaction, error) {
	const query = `
		SELECT t.encoded
		  FROM statev2.transaction t
		 INNER JOIN statev2.l2block b
		    ON t.l2_block_num = b.batch_num
		 WHERE b.block_hash = $1
		   AND 0 = $2
	`

	var encoded string
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, query, blockHash.Hex(), index).Scan(&encoded)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := decodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *PostgresStorage) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNumber uint64, index uint64, dbTx pgx.Tx) (*types.Transaction, error) {
	const query = `
		SELECT t.encoded
		  FROM statev2.transaction t
		 WHERE t.l2_block_num = $1
		   AND 0 = $2
	`

	var encoded string
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, query, blockNumber, index).Scan(&encoded)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := decodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *PostgresStorage) GetBlockTransactionCountByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (uint64, error) {
	const query = `
		SELECT COUNT(*)
		  FROM statev2.transaction t
		 INNER JOIN statev2.l2block b
			ON b.block_num = t.l2_block_num
		 WHERE b.block_hash = $1
	`

	var count uint64
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, query, hash.Hex()).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *PostgresStorage) GetBlockTransactionCountByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	const query = `
		SELECT COUNT(*)
		  FROM state.transaction t
		 WHERE t.l2_block_num = $1
	`
	var count uint64
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, query, blockNumber).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// getTransactionLogs returns the logs of a transaction by transaction hash
func (p *PostgresStorage) getTransactionLogs(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) ([]*types.Log, error) {
	const query = `
		SELECT t.l2_block_num
			 , b.block_hash
			 , l.tx_hash
			 , l.log_index
			 , l.address
			 , l.data
			 , l.topic0
			 , l.topic1
			 , l.topic2
			 , l.topic3
		  FROM state.log l
		 INNER JOIN statev2.transaction t
			ON t.hash = l.tx_hash
		 INNER JOIN statev2.l2block b
			ON b.block_num = t.l2_block_num
		 WHERE transaction_hash = $1
	`
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, query, transactionHash)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]*types.Log, 0, len(rows.RawValues()))

	for rows.Next() {
		var log types.Log
		var txHash, logAddress, logData, topic0 string
		var topic1, topic2, topic3 *string

		err := rows.Scan(
			&log.BlockNumber,
			&log.BlockHash,
			&txHash,
			&log.Index,
			&logAddress,
			&logData,
			&topic0,
			&topic1,
			&topic2,
			&topic3)
		if err != nil {
			return nil, err
		}

		log.TxHash = common.HexToHash(txHash)
		log.Address = common.HexToAddress(logAddress)
		log.TxIndex = uint(0)
		log.Data = []byte(logData)

		log.Topics = []common.Hash{common.HexToHash(topic0)}
		if topic1 != nil {
			log.Topics = append(log.Topics, common.HexToHash(*topic1))
		}

		if topic2 != nil {
			log.Topics = append(log.Topics, common.HexToHash(*topic2))
		}

		if topic3 != nil {
			log.Topics = append(log.Topics, common.HexToHash(*topic3))
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

func decodeTx(encodedTx string) (*types.Transaction, error) {
	b, err := hex.DecodeHex(encodedTx)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return tx, nil
}
