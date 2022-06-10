package state

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state/store"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	maxTopics = 4
)

const (
	getLastBlockSQL                        = "SELECT * FROM state.block ORDER BY block_num DESC LIMIT 1"
	getPreviousBlockSQL                    = "SELECT * FROM state.block ORDER BY block_num DESC LIMIT 1 OFFSET $1"
	getBlockByHashSQL                      = "SELECT * FROM state.block WHERE block_hash = $1"
	getBlockByNumberSQL                    = "SELECT * FROM state.block WHERE block_num = $1"
	getLastBlockNumberSQL                  = "SELECT COALESCE(MAX(block_num), 0) FROM state.block"
	getLastVirtualBatchSQL                 = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root, rollup_exit_root FROM state.batch ORDER BY batch_num DESC LIMIT 1"
	getLastConsolidatedBatchSQL            = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root, rollup_exit_root FROM state.batch WHERE consolidated_tx_hash != $1 ORDER BY batch_num DESC LIMIT 1"
	getPreviousVirtualBatchSQL             = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root, rollup_exit_root FROM state.batch ORDER BY batch_num DESC LIMIT 1 OFFSET $1"
	getPreviousConsolidatedBatchSQL        = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root, rollup_exit_root FROM state.batch WHERE consolidated_tx_hash != $1 ORDER BY batch_num DESC LIMIT 1 OFFSET $2"
	getBatchByHashSQL                      = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root, rollup_exit_root FROM state.batch WHERE batch_hash = $1"
	getBatchByNumberSQL                    = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root, rollup_exit_root FROM state.batch WHERE batch_num = $1"
	getBatchHashesSinceSQL                 = "SELECT header->>'hash' as hash FROM state.batch WHERE received_at >= $1"
	getLastBatchByStateRootSQL             = "SELECT block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, consolidated_at, chain_id, global_exit_root, rollup_exit_root FROM state.batch WHERE header->>'stateRoot' = $1 ORDER BY batch_num DESC LIMIT 1"
	getLastVirtualBatchNumberSQL           = "SELECT COALESCE(MAX(batch_num), 0) FROM state.batch"
	getLastConsolidatedBatchNumberSQL      = "SELECT COALESCE(MAX(batch_num), 0) FROM state.batch WHERE consolidated_tx_hash != $1"
	getTransactionByHashSQL                = "SELECT transaction.encoded FROM state.transaction WHERE hash = $1"
	getTransactionByBatchHashAndIndexSQL   = "SELECT transaction.encoded FROM state.transaction inner join state.batch on (state.transaction.batch_num = state.batch.batch_num) WHERE state.batch.batch_hash = $1 and state.transaction.tx_index = $2"
	getTransactionByBatchNumberAndIndexSQL = "SELECT transaction.encoded FROM state.transaction WHERE batch_num = $1 AND tx_index = $2"
	getTransactionCountSQL                 = "SELECT COUNT(*) FROM state.transaction WHERE from_address = $1"
	getTransactionCountByBatchHashSQL      = "SELECT COUNT(*) FROM state.transaction t, state.batch b WHERE t.batch_num = b.batch_num AND batch_hash = $1"
	getTransactionCountByBatchNumberSQL    = "SELECT COUNT(*) FROM state.transaction WHERE batch_num = $1"
	consolidateBatchSQL                    = "UPDATE state.batch SET consolidated_tx_hash = $1, consolidated_at = $3, aggregator = $4 WHERE batch_num = $2"
	getTxsByBatchNumSQL                    = "SELECT transaction.encoded FROM state.transaction WHERE batch_num = $1"
	addBlockSQL                            = "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)"
	addSequencerSQL                        = "INSERT INTO state.sequencer (address, url, chain_id, block_num) VALUES ($1, $2, $3, $4) ON CONFLICT (chain_id) DO UPDATE SET address = EXCLUDED.address, url = EXCLUDED.url, block_num = EXCLUDED.block_num"
	updateLastBatchSeenSQL                 = "UPDATE state.misc SET last_batch_num_seen = $1"
	getLastBatchSeenSQL                    = "SELECT last_batch_num_seen FROM state.misc LIMIT 1"
	getSyncingInfoSQL                      = "SELECT last_batch_num_seen, last_batch_num_consolidated, init_sync_batch FROM state.misc LIMIT 1"
	updateLastBatchConsolidatedSQL         = "UPDATE state.misc SET last_batch_num_consolidated = $1"
	updateInitBatchSQL                     = "UPDATE state.misc SET init_sync_batch = $1"
	getLastBatchConsolidatedSQL            = "SELECT last_batch_num_consolidated FROM state.misc LIMIT 1"
	getSequencerSQL                        = "SELECT * FROM state.sequencer WHERE address = $1"
	getReceiptSQL                          = "SELECT * FROM state.receipt WHERE tx_hash = $1"
	resetSQL                               = "DELETE FROM state.block WHERE block_num > $1"
	resetConsolidationSQL                  = "UPDATE state.batch SET aggregator = decode('0000000000000000000000000000000000000000', 'hex'), consolidated_tx_hash = decode('0000000000000000000000000000000000000000000000000000000000000000', 'hex'), consolidated_at = null WHERE consolidated_at > $1"
	addBatchSQL                            = "INSERT INTO state.batch (batch_num, batch_hash, block_num, sequencer, aggregator, consolidated_tx_hash, header, uncles, raw_txs_data, matic_collateral, received_at, chain_id, global_exit_root, rollup_exit_root) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)"
	addTransactionSQL                      = "INSERT INTO state.transaction (hash, from_address, encoded, decoded, batch_num, tx_index) VALUES($1, $2, $3, $4, $5, $6)"
	addReceiptSQL                          = "INSERT INTO state.receipt (type, post_state, status, cumulative_gas_used, gas_used, batch_num, batch_hash, tx_hash, tx_index, tx_from, tx_to, contract_address)	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
	addLogSQL                              = "INSERT INTO state.log (log_index, transaction_index, transaction_hash, batch_hash, batch_num, address, data, topic0, topic1, topic2, topic3) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	getTransactionLogsSQL                  = "SELECT * FROM state.log WHERE transaction_hash = $1"
	getLogsSQLByBatchHash                  = "SELECT * FROM state.log WHERE batch_hash = $1"
	getLogsByFilter                        = "SELECT state.log.* FROM state.log INNER JOIN state.batch ON state.log.batch_num = state.batch.batch_num WHERE state.log.batch_num BETWEEN $1 AND $2 AND (address = any($3) OR $3 IS NULL) AND (topic0 = any($4) OR $4 IS NULL) AND (topic1 = any($5) OR $5 IS NULL) AND (topic2 = any($6) OR $6 IS NULL) AND (topic3 = any($7) OR $7 IS NULL) AND (state.batch.received_at >= $8 OR $8 IS NULL)"
	addGlobalExitRootSQL                   = "INSERT INTO state.exit_root (block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root) VALUES ($1, $2, $3, $4)"
	getLatestExitRootSQL                   = "SELECT block_num, global_exit_root_num, mainnet_exit_root, rollup_exit_root FROM state.exit_root ORDER BY global_exit_root_num DESC LIMIT 1"
)

var (
	ten = big.NewInt(encoding.Base10)
)

// PostgresStorage implements the Storage interface
type PostgresStorage struct {
	*store.Pg
}

// NewPostgresStorage creates a new StateDB
func NewPostgresStorage(db *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{
		Pg: store.NewPg(db),
	}
}

// GetLastBlock gets the latest block
func (s *PostgresStorage) GetLastBlock(ctx context.Context, txBundleID string) (*Block, error) {
	var block Block
	err := s.QueryRow(ctx, txBundleID, getLastBlockSQL).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetPreviousBlock gets the offset previous block respect to latest
func (s *PostgresStorage) GetPreviousBlock(ctx context.Context, offset uint64, txBundleID string) (*Block, error) {
	var block Block
	err := s.QueryRow(ctx, txBundleID, getPreviousBlockSQL, offset).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetBlockByHash gets the block with the required hash
func (s *PostgresStorage) GetBlockByHash(ctx context.Context, hash common.Hash, txBundleID string) (*Block, error) {
	var block Block
	err := s.QueryRow(ctx, txBundleID, getBlockByHashSQL, hash).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetBlockByNumber gets the block with the required number
func (s *PostgresStorage) GetBlockByNumber(ctx context.Context, blockNumber uint64, txBundleID string) (*Block, error) {
	var block Block
	err := s.QueryRow(ctx, txBundleID, getBlockByNumberSQL, blockNumber).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetLastBlockNumber gets the latest block number
func (s *PostgresStorage) GetLastBlockNumber(ctx context.Context, txBundleID string) (uint64, error) {
	var lastBlockNum uint64
	err := s.QueryRow(ctx, txBundleID, getLastBlockNumberSQL).Scan(&lastBlockNum)

	if reflect.TypeOf(err) == reflect.TypeOf(pgx.ScanArgError{}) {
		return 0, ErrStateNotSynchronized
	} else if err != nil {
		return 0, err
	}

	return lastBlockNum, nil
}

// GetLastBatch gets the latest batch
func (s *PostgresStorage) GetLastBatch(ctx context.Context, isVirtual bool, txBundleID string) (*Batch, error) {
	var (
		batch           Batch
		maticCollateral pgtype.Numeric

		chain uint64
		err   error
	)

	if isVirtual {
		err = s.QueryRow(ctx, txBundleID, getLastVirtualBatchSQL).Scan(&batch.BlockNumber,
			&batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
			&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
			&batch.ReceivedAt, &batch.ConsolidatedAt, &chain, &batch.GlobalExitRoot, &batch.RollupExitRoot)
	} else {
		err = s.QueryRow(ctx, txBundleID, getLastConsolidatedBatchSQL, common.Hash{}).Scan(
			&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
			&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
			&batch.ReceivedAt, &batch.ConsolidatedAt, &chain, &batch.GlobalExitRoot, &batch.RollupExitRoot)
	}
	batch.ChainID = new(big.Int).SetUint64(chain)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	batch.MaticCollateral = new(big.Int).Mul(maticCollateral.Int, big.NewInt(0).Exp(ten, big.NewInt(int64(maticCollateral.Exp)), nil))

	return &batch, nil
}

// GetPreviousBatch gets the offset previous batch respect to latest
func (s *PostgresStorage) GetPreviousBatch(ctx context.Context, isVirtual bool, offset uint64, txBundleID string) (*Batch, error) {
	var (
		batch           Batch
		maticCollateral pgtype.Numeric

		chain uint64
		err   error
	)

	if isVirtual {
		err = s.QueryRow(ctx, txBundleID, getPreviousVirtualBatchSQL, offset).Scan(
			&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
			&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
			&batch.ReceivedAt, &batch.ConsolidatedAt, &chain, &batch.GlobalExitRoot, &batch.RollupExitRoot)
	} else {
		err = s.QueryRow(ctx, txBundleID, getPreviousConsolidatedBatchSQL, common.Hash{}, offset).Scan(
			&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash, &batch.Header,
			&batch.Uncles, &batch.RawTxsData, &maticCollateral,
			&batch.ReceivedAt, &batch.ConsolidatedAt, &chain, &batch.GlobalExitRoot, &batch.RollupExitRoot)
	}
	batch.ChainID = new(big.Int).SetUint64(chain)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	batch.MaticCollateral = new(big.Int).Mul(maticCollateral.Int, big.NewInt(0).Exp(ten, big.NewInt(int64(maticCollateral.Exp)), nil))
	return &batch, nil
}

// GetBatchByHash gets the batch with the required hash
func (s *PostgresStorage) GetBatchByHash(ctx context.Context, hash common.Hash, txBundleID string) (*Batch, error) {
	var (
		batch           Batch
		maticCollateral pgtype.Numeric
		chain           uint64
	)
	err := s.QueryRow(ctx, txBundleID, getBatchByHashSQL, hash).Scan(
		&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
		&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
		&batch.ReceivedAt, &batch.ConsolidatedAt, &chain, &batch.GlobalExitRoot, &batch.GlobalExitRoot)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	batch.ChainID = new(big.Int).SetUint64(chain)
	batch.MaticCollateral = new(big.Int).Mul(maticCollateral.Int, big.NewInt(0).Exp(ten, big.NewInt(int64(maticCollateral.Exp)), nil))
	batch.Transactions, err = s.getBatchTransactions(ctx, batch, txBundleID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return &batch, nil
}

// GetBatchTransactionCountByHash return the number of transactions in the batch
func (s *PostgresStorage) GetBatchTransactionCountByHash(ctx context.Context, hash common.Hash, txBundleID string) (uint64, error) {
	var count uint64
	err := s.QueryRow(ctx, txBundleID, getTransactionCountByBatchHashSQL, hash.Bytes()).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetBatchTransactionCountByNumber return the number of transactions in the batch
func (s *PostgresStorage) GetBatchTransactionCountByNumber(ctx context.Context, batchNumber uint64, txBundleID string) (uint64, error) {
	var count uint64
	err := s.QueryRow(ctx, txBundleID, getTransactionCountByBatchNumberSQL, batchNumber).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetBatchByNumber gets the batch with the required number
func (s *PostgresStorage) GetBatchByNumber(ctx context.Context, batchNumber uint64, txBundleID string) (*Batch, error) {
	batch, err := s.getBatchWithoutTxsByNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return nil, err
	}

	batch.Transactions, err = s.getBatchTransactions(ctx, *batch, txBundleID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return batch, nil
}

// GetLastBatchByStateRoot gets the last batch with the required state root
func (s *PostgresStorage) GetLastBatchByStateRoot(ctx context.Context, stateRoot []byte, txBundleID string) (*Batch, error) {
	batch, err := s.getLastBatchWithoutTxsByStateRoot(ctx, stateRoot, txBundleID)
	if err != nil {
		return nil, err
	}

	batch.Transactions, err = s.getBatchTransactions(ctx, *batch, txBundleID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return batch, nil
}

// GetBatchHeader gets the batch header with the required number.
func (s *PostgresStorage) GetBatchHeader(ctx context.Context, batchNumber uint64, txBundleID string) (*types.Header, error) {
	batch, err := s.getBatchWithoutTxsByNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return nil, err
	}
	return batch.Header, nil
}

// GetLastBatchNumber gets the latest batch number
func (s *PostgresStorage) GetLastBatchNumber(ctx context.Context, txBundleID string) (uint64, error) {
	var lastBatchNumber uint64
	err := s.QueryRow(ctx, txBundleID, getLastVirtualBatchNumberSQL).Scan(&lastBatchNumber)

	if err != nil {
		return 0, err
	}

	return lastBatchNumber, nil
}

// GetLastConsolidatedBatchNumber gets the latest consolidated batch number
func (s *PostgresStorage) GetLastConsolidatedBatchNumber(ctx context.Context, txBundleID string) (uint64, error) {
	var lastBatchNumber uint64
	err := s.QueryRow(ctx, txBundleID, getLastConsolidatedBatchNumberSQL, common.Hash{}).Scan(&lastBatchNumber)

	if err != nil {
		return 0, err
	}

	return lastBatchNumber, nil
}

// GetTransactionByBatchHashAndIndex gets a transaction from a batch by index
func (s *PostgresStorage) GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64, txBundleID string) (*types.Transaction, error) {
	var encoded string
	err := s.QueryRow(ctx, txBundleID, getTransactionByBatchHashAndIndexSQL, batchHash.Bytes(), index).Scan(&encoded)

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

	return tx, nil
}

// GetTransactionByBatchNumberAndIndex gets a transaction from a batch by index
func (s *PostgresStorage) GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64, txBundleID string) (*types.Transaction, error) {
	var encoded string
	err := s.QueryRow(ctx, txBundleID, getTransactionByBatchNumberAndIndexSQL, batchNumber, index).Scan(&encoded)

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

	return tx, nil
}

// GetTransactionByHash gets a transaction by its hash
func (s *PostgresStorage) GetTransactionByHash(ctx context.Context, transactionHash common.Hash, txBundleID string) (*types.Transaction, error) {
	var encoded string
	err := s.QueryRow(ctx, txBundleID, getTransactionByHashSQL, transactionHash).Scan(&encoded)

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

	return tx, nil
}

// GetTransactionCount returns the number of transactions sent from an address
func (s *PostgresStorage) GetTransactionCount(ctx context.Context, fromAddress common.Address, txBundleID string) (uint64, error) {
	var count uint64
	err := s.QueryRow(ctx, txBundleID, getTransactionCountSQL, fromAddress).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
func (s *PostgresStorage) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, txBundleID string) (*Receipt, error) {
	var receipt Receipt
	var batchNumber uint64
	var to *[]byte

	err := s.QueryRow(ctx, txBundleID, getReceiptSQL, transactionHash).Scan(&receipt.Type, &receipt.PostState, &receipt.Status,
		&receipt.CumulativeGasUsed, &receipt.GasUsed, &batchNumber, &receipt.BlockHash, &receipt.TxHash, &receipt.TransactionIndex, &receipt.From, &to, &receipt.ContractAddress)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if to != nil {
		toAddr := common.BytesToAddress(*to)
		receipt.To = &toAddr
	}

	receipt.BlockNumber = new(big.Int).SetUint64(batchNumber)

	logs, err := s.getTransactionLogs(ctx, transactionHash, txBundleID)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}

	receipt.Logs = logs
	receipt.Bloom = types.CreateBloom(types.Receipts{&receipt.Receipt})

	return &receipt, nil
}

// getTransactionLogs returns the logs of a transaction by transaction hash
func (s *PostgresStorage) getTransactionLogs(ctx context.Context, transactionHash common.Hash, txBundleID string) ([]*types.Log, error) {
	rows, err := s.Query(ctx, txBundleID, getTransactionLogsSQL, transactionHash.Bytes())
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]*types.Log, 0, len(rows.RawValues()))

	for rows.Next() {
		var log types.Log
		var topicsAsBytes [maxTopics]*[]byte
		err := rows.Scan(&log.Index, &log.TxIndex, &log.TxHash, &log.BlockHash, &log.BlockNumber, &log.Address, &log.Data,
			&topicsAsBytes[0], &topicsAsBytes[1], &topicsAsBytes[2], &topicsAsBytes[3])
		if err != nil {
			return nil, err
		}

		log.Topics = []common.Hash{}
		for i := 0; i < maxTopics; i++ {
			if topicsAsBytes[i] != nil {
				topicHash := common.BytesToHash(*topicsAsBytes[i])
				log.Topics = append(log.Topics, topicHash)
			}
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

func (s *PostgresStorage) addressesToBytes(addresses []common.Address) [][]byte {
	converted := make([][]byte, 0, len(addresses))

	for _, address := range addresses {
		converted = append(converted, address.Bytes())
	}

	return converted
}

func (s *PostgresStorage) hashesToBytes(hashes []common.Hash) [][]byte {
	converted := make([][]byte, 0, len(hashes))

	for _, hash := range hashes {
		converted = append(converted, hash.Bytes())
	}

	return converted
}

// GetLogs returns the logs that match the filter
func (s *State) GetLogs(ctx context.Context, fromBatch uint64, toBatch uint64, addresses []common.Address, topics [][]common.Hash, batchHash *common.Hash, since *time.Time, txBundleID string) ([]*types.Log, error) {
	var err error
	var rows pgx.Rows
	if batchHash != nil {
		rows, err = s.Query(ctx, txBundleID, getLogsSQLByBatchHash, batchHash.Bytes())
	} else {
		args := []interface{}{fromBatch, toBatch}

		if len(addresses) > 0 {
			args = append(args, s.addressesToBytes(addresses))
		} else {
			args = append(args, nil)
		}

		for i := 0; i < maxTopics; i++ {
			if len(topics) > i && len(topics[i]) > 0 {
				args = append(args, s.hashesToBytes(topics[i]))
			} else {
				args = append(args, nil)
			}
		}

		args = append(args, since)

		rows, err = s.Query(ctx, txBundleID, getLogsByFilter, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]*types.Log, 0, len(rows.RawValues()))

	for rows.Next() {
		var log types.Log
		var topicsAsBytes [maxTopics]*[]byte
		err := rows.Scan(&log.Index, &log.TxIndex, &log.TxHash, &log.BlockHash, &log.BlockNumber, &log.Address, &log.Data,
			&topicsAsBytes[0], &topicsAsBytes[1], &topicsAsBytes[2], &topicsAsBytes[3])
		if err != nil {
			return nil, err
		}

		log.Topics = []common.Hash{}
		for i := 0; i < maxTopics; i++ {
			if topicsAsBytes[i] != nil {
				topicHash := common.BytesToHash(*topicsAsBytes[i])
				log.Topics = append(log.Topics, topicHash)
			}
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

// Reset resets the state to a block
func (s *PostgresStorage) Reset(ctx context.Context, block *Block, txBundleID string) error {
	if _, err := s.Exec(ctx, txBundleID, resetSQL, block.BlockNumber); err != nil {
		return err
	}
	//Remove consolidations
	if _, err := s.Exec(ctx, txBundleID, resetConsolidationSQL, block.ReceivedAt); err != nil {
		return err
	}
	return nil
}

// ConsolidateBatch changes the virtual status of a batch
func (s *PostgresStorage) ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash, consolidatedAt time.Time, aggregator common.Address, txBundleID string) error {
	_, err := s.Exec(ctx, txBundleID, consolidateBatchSQL, consolidatedTxHash, batchNumber, consolidatedAt, aggregator)
	return err
}

// GetTxsByBatchNum returns all the txs in a given batch
func (s *PostgresStorage) GetTxsByBatchNum(ctx context.Context, batchNum uint64, txBundleID string) ([]*types.Transaction, error) {
	rows, err := s.Query(ctx, txBundleID, getTxsByBatchNumSQL, batchNum)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	txs := make([]*types.Transaction, 0, len(rows.RawValues()))
	var (
		encoded string
		tx      *types.Transaction
		b       []byte
	)
	for rows.Next() {
		if err = rows.Scan(&encoded); err != nil {
			return nil, err
		}

		tx = new(types.Transaction)

		b, err = hex.DecodeHex(encoded)
		if err != nil {
			return nil, err
		}

		if err := tx.UnmarshalBinary(b); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

// AddSequencer stores a new sequencer
func (s *PostgresStorage) AddSequencer(ctx context.Context, seq Sequencer, txBundleID string) error {
	_, err := s.Exec(ctx, txBundleID, addSequencerSQL, seq.Address, seq.URL, seq.ChainID.Uint64(), seq.BlockNumber)
	return err
}

// GetSequencer gets a sequencer
func (s *PostgresStorage) GetSequencer(ctx context.Context, address common.Address, txBundleID string) (*Sequencer, error) {
	var seq Sequencer
	var cID uint64
	err := s.QueryRow(ctx, txBundleID, getSequencerSQL, address.Bytes()).Scan(&seq.Address, &seq.URL, &cID, &seq.BlockNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	seq.ChainID = big.NewInt(0).SetUint64(cID)

	return &seq, nil
}

// AddBlock adds a new block to the State Store
func (s *PostgresStorage) AddBlock(ctx context.Context, block *Block, txBundleID string) error {
	_, err := s.Exec(ctx, txBundleID, addBlockSQL, block.BlockNumber, block.BlockHash.Bytes(), block.ParentHash.Bytes(), block.ReceivedAt)
	return err
}

// SetLastBatchNumberSeenOnEthereum sets the last batch number that affected
// the roll-up in order to allow the components to know if the state
// is synchronized or not
func (s *PostgresStorage) SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64, txBundleID string) error {
	_, err := s.Exec(ctx, txBundleID, updateLastBatchSeenSQL, batchNumber)
	return err
}

// GetLastBatchNumberSeenOnEthereum returns the last batch number stored
// in the state that represents the last batch number that affected the
// roll-up in the Ethereum network.
func (s *PostgresStorage) GetLastBatchNumberSeenOnEthereum(ctx context.Context, txBundleID string) (uint64, error) {
	var batchNumber uint64
	err := s.QueryRow(ctx, txBundleID, getLastBatchSeenSQL).Scan(&batchNumber)

	if err != nil {
		return 0, err
	}

	return batchNumber, nil
}

// GetSyncingInfo returns information regarding the syncing status of the node
func (s *PostgresStorage) GetSyncingInfo(ctx context.Context, txBundleID string) (SyncingInfo, error) {
	var info SyncingInfo
	err := s.QueryRow(ctx, txBundleID, getSyncingInfoSQL).Scan(&info.LastBatchNumberSeen, &info.LastBatchNumberConsolidated, &info.InitialSyncingBatch)
	return info, err
}

// SetLastBatchNumberConsolidatedOnEthereum sets the last batch number that was consolidated on ethereum
func (s *PostgresStorage) SetLastBatchNumberConsolidatedOnEthereum(ctx context.Context, batchNumber uint64, txBundleID string) error {
	_, err := s.Exec(ctx, txBundleID, updateLastBatchConsolidatedSQL, batchNumber)
	return err
}

// SetInitSyncBatch sets the initial batch number where the synchronization started
func (s *PostgresStorage) SetInitSyncBatch(ctx context.Context, batchNumber uint64, txBundleID string) error {
	_, err := s.Exec(ctx, txBundleID, updateInitBatchSQL, batchNumber)
	return err
}

// GetLastBatchNumberConsolidatedOnEthereum sets the last batch number that was consolidated on ethereum
func (s *PostgresStorage) GetLastBatchNumberConsolidatedOnEthereum(ctx context.Context, txBundleID string) (uint64, error) {
	var batchNumber uint64
	err := s.QueryRow(ctx, txBundleID, getLastBatchConsolidatedSQL).Scan(&batchNumber)

	if err != nil {
		return 0, err
	}

	return batchNumber, nil
}

// AddBatch adds a new batch to the State Store
func (s *PostgresStorage) AddBatch(ctx context.Context, batch *Batch, txBundleID string) error {
	_, err := s.Exec(ctx, txBundleID, addBatchSQL, batch.Number().Uint64(), batch.Hash(), batch.BlockNumber, batch.Sequencer, batch.Aggregator,
		batch.ConsolidatedTxHash, batch.Header, batch.Uncles, batch.RawTxsData, batch.MaticCollateral.String(), batch.ReceivedAt, batch.ChainID.String(), batch.GlobalExitRoot, batch.RollupExitRoot)
	return err
}

// AddTransaction adds a new transqaction to the State Store
func (s *PostgresStorage) AddTransaction(ctx context.Context, tx *types.Transaction, batchNumber uint64, index uint, txBundleID string) error {
	binary, err := tx.MarshalBinary()
	if err != nil {
		panic(err)
	}
	encoded := hex.EncodeToHex(binary)

	binary, err = tx.MarshalJSON()
	if err != nil {
		panic(err)
	}
	decoded := string(binary)

	_, err = s.Exec(ctx, txBundleID, addTransactionSQL, tx.Hash().Bytes(), "", encoded, decoded, batchNumber, index)
	return err
}

// AddReceipt adds a new receipt to the State Store
func (s *PostgresStorage) AddReceipt(ctx context.Context, receipt *Receipt, txBundleID string) error {
	var to *[]byte
	if receipt.To != nil {
		b := receipt.To.Bytes()
		to = &b
	}

	_, err := s.Exec(ctx, txBundleID, addReceiptSQL, receipt.Type, receipt.PostState, receipt.Status, receipt.CumulativeGasUsed, receipt.GasUsed, receipt.BlockNumber.Uint64(), receipt.BlockHash.Bytes(), receipt.TxHash.Bytes(), receipt.TransactionIndex, receipt.From.Bytes(), to, receipt.ContractAddress.Bytes())
	return err
}

// AddLog adds a new log to the State Store
func (s *PostgresStorage) AddLog(ctx context.Context, l types.Log, txBundleID string) error {
	var topicsAsBytes [maxTopics]*[]byte
	for i := 0; i < len(l.Topics); i++ {
		topicBytes := l.Topics[i].Bytes()
		topicsAsBytes[i] = &topicBytes
	}

	_, err := s.Exec(ctx, txBundleID, addLogSQL, l.Index, l.TxIndex, l.TxHash.Bytes(),
		l.BlockHash.Bytes(), l.BlockNumber, l.Address.Bytes(), l.Data,
		topicsAsBytes[0], topicsAsBytes[1], topicsAsBytes[2], topicsAsBytes[3])
	return err
}

func (s *PostgresStorage) getBatchTransactions(ctx context.Context, batch Batch, txBundleID string) ([]*types.Transaction, error) {
	transactions, err := s.GetTxsByBatchNum(ctx, batch.Number().Uint64(), txBundleID)
	if errors.Is(err, pgx.ErrNoRows) {
		transactions = []*types.Transaction{}
	} else if err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetBatchHashesSince get all batch hashes since a timestamp
func (s *State) GetBatchHashesSince(ctx context.Context, since time.Time, txBundleID string) ([]common.Hash, error) {
	rows, err := s.Query(ctx, txBundleID, getBatchHashesSinceSQL, since)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	defer rows.Close()

	batchHashes := make([]common.Hash, 0, len(rows.RawValues()))

	for rows.Next() {
		var batchHash string
		err := rows.Scan(&batchHash)
		if err != nil {
			return nil, err
		}

		batchHashes = append(batchHashes, common.HexToHash(batchHash))
	}

	return batchHashes, nil
}

func (s *PostgresStorage) getBatchWithoutTxsByNumber(ctx context.Context, batchNumber uint64, txBundleID string) (*Batch, error) {
	var (
		batch           Batch
		maticCollateral pgtype.Numeric
		chain           uint64
	)
	err := s.QueryRow(ctx, txBundleID, getBatchByNumberSQL, batchNumber).Scan(
		&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
		&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
		&batch.ReceivedAt, &batch.ConsolidatedAt, &chain, &batch.GlobalExitRoot, &batch.RollupExitRoot)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	batch.ChainID = new(big.Int).SetUint64(chain)
	batch.MaticCollateral = new(big.Int).Mul(maticCollateral.Int, big.NewInt(0).Exp(ten, big.NewInt(int64(maticCollateral.Exp)), nil))

	return &batch, nil
}

func (s *PostgresStorage) getLastBatchWithoutTxsByStateRoot(ctx context.Context, stateRoot []byte, txBundleID string) (*Batch, error) {
	var (
		batch           Batch
		maticCollateral pgtype.Numeric
		chain           uint64
	)
	err := s.QueryRow(ctx, txBundleID, getLastBatchByStateRootSQL, hex.EncodeToHex(stateRoot)).Scan(
		&batch.BlockNumber, &batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
		&batch.Header, &batch.Uncles, &batch.RawTxsData, &maticCollateral,
		&batch.ReceivedAt, &batch.ConsolidatedAt, &chain, &batch.GlobalExitRoot, &batch.RollupExitRoot)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	batch.ChainID = new(big.Int).SetUint64(chain)
	batch.MaticCollateral = new(big.Int).Mul(maticCollateral.Int, big.NewInt(0).Exp(ten, big.NewInt(int64(maticCollateral.Exp)), nil))

	return &batch, nil
}

// AddExitRoot adds a new ExitRoot to the db
func (s *PostgresStorage) AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, txBundleID string) error {
	_, err := s.Exec(ctx, txBundleID, addGlobalExitRootSQL, exitRoot.BlockNumber, exitRoot.GlobalExitRootNum.String(), exitRoot.MainnetExitRoot, exitRoot.RollupExitRoot)
	return err
}

// GetLatestExitRoot get the latest ExitRoot from L1.
func (s *PostgresStorage) GetLatestGlobalExitRoot(ctx context.Context, txBundleID string) (*GlobalExitRoot, error) {
	var (
		exitRoot        GlobalExitRoot
		globalNum       uint64
	)
	err := s.QueryRow(ctx, txBundleID, getLatestExitRootSQL).Scan(&exitRoot.BlockNumber, &globalNum, &exitRoot.MainnetExitRoot, &exitRoot.RollupExitRoot)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	exitRoot.GlobalExitRootNum = new(big.Int).SetUint64(globalNum)
	return &exitRoot, nil
}
