package state

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// State is the interface of the Hermez state
type State interface {
	NewBatchProcessor(sequencerAddress common.Address, lastBatchNumber uint64) (BatchProcessor, error)
	NewGenesisBatchProcessor(genesisStateRoot []byte) (BatchProcessor, error)
	GetStateRoot(ctx context.Context, virtual bool) ([]byte, error)
	GetBalance(address common.Address, batchNumber uint64) (*big.Int, error)
	EstimateGas(transaction *types.Transaction) uint64
	GetLastBlock(ctx context.Context) (*Block, error)
	GetPreviousBlock(ctx context.Context, offset uint64) (*Block, error)
	GetBlockByHash(ctx context.Context, hash common.Hash) (*Block, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error)
	GetLastBlockNumber(ctx context.Context) (uint64, error)
	GetNonce(address common.Address, batchNumber uint64) (uint64, error)
	GetLastBatch(ctx context.Context, isVirtual bool) (*Batch, error)
	GetPreviousBatch(ctx context.Context, isVirtual bool, offset uint64) (*Batch, error)
	GetBatchByHash(ctx context.Context, hash common.Hash) (*Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64) (*Batch, error)
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error)
	GetTransactionCount(ctx context.Context, address common.Address) (uint64, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*types.Receipt, error)
	Reset(blockNumber uint64) error
	ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash) error
	GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error)
	AddSequencer(ctx context.Context, seq Sequencer) error
	GetSequencer(ctx context.Context, address common.Address) (*Sequencer, error)
	SetGenesis(ctx context.Context, genesis Genesis) error
	AddBlock(ctx context.Context, block *Block) error
	SetLastBatchNumberSeenOnEthereum(batchNumber uint64) error
	GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error)
	GetStateRootByBatchNumber(batchNumber uint64) ([]byte, error)
}

const (
	getLastBlockSQL                 = "SELECT * FROM state.block ORDER BY block_num DESC LIMIT 1"
	getPreviousBlockSQL             = "SELECT * FROM state.block ORDER BY block_num DESC LIMIT 1 OFFSET $1"
	getBlockByHashSQL               = "SELECT * FROM state.block WHERE block_hash = $1"
	getBlockByNumberSQL             = "SELECT * FROM state.block WHERE block_num = $1"
	getLastBlockNumberSQL           = "SELECT MAX(block_num) FROM state.block"
	getLastVirtualBatchSQL          = "SELECT * FROM state.batch ORDER BY batch_num DESC LIMIT 1"
	getLastConsolidatedBatchSQL     = "SELECT * FROM state.batch WHERE consolidated_tx_hash != $1 ORDER BY batch_num DESC LIMIT 1"
	getPreviousVirtualBatchSQL      = "SELECT * FROM state.batch ORDER BY batch_num DESC LIMIT 1 OFFSET $1"
	getPreviousConsolidatedBatchSQL = "SELECT * FROM state.batch WHERE consolidated_tx_hash != $1 ORDER BY batch_num DESC LIMIT 1 OFFSET $2"
	getBatchByHashSQL               = "SELECT * FROM state.batch WHERE batch_hash = $1"
	getBatchByNumberSQL             = "SELECT * FROM state.batch WHERE batch_num = $1"
	getLastBatchNumberSQL           = "SELECT COALESCE(MAX(batch_num), 0) FROM state.batch"
	getTransactionByHashSQL         = "SELECT transaction.encoded FROM state.transaction WHERE hash = $1"
	getTransactionCountSQL          = "SELECT COUNT(*) FROM state.transaction WHERE from_address = $1"
	consolidateBatchSQL             = "UPDATE state.batch SET consolidated_tx_hash = $1 WHERE batch_num = $2"
	getTxsByBatchNumSQL             = "SELECT transaction.encoded FROM state.transaction WHERE batch_num = $1"
	addBlockSQL                     = "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)"
	addSequencerSQL                 = "INSERT INTO state.sequencer (address, url, chain_id, block_num) VALUES ($1, $2, $3, $4)"
	getSequencerSQL                 = "SELECT * FROM state.sequencer WHERE address = $1"
	getReceiptSQL                   = "SELECT * FROM state.receipt WHERE tx_hash = $1"
)

var (
	// ErrInvalidBatchHeader indicates the batch header is invalid
	ErrInvalidBatchHeader = errors.New("invalid batch header")
)

// BasicState is a implementation of the state
type BasicState struct {
	db   *pgxpool.Pool
	Tree tree.ReadWriter
}

// NewState creates a new State
func NewState(db *pgxpool.Pool, tree tree.ReadWriter) State {
	return &BasicState{db: db, Tree: tree}
}

// NewBatchProcessor creates a new batch processor
func (s *BasicState) NewBatchProcessor(sequencerAddress common.Address, lastBatchNumber uint64) (BatchProcessor, error) {
	// init correct state root from previous batch
	stateRoot, err := s.GetStateRootByBatchNumber(lastBatchNumber)
	if err != nil {
		return nil, err
	}

	s.Tree.SetCurrentRoot(stateRoot)

	var chainID = uint64(0)

	// Get Sequencer's Chain ID
	sq, err := s.GetSequencer(context.Background(), sequencerAddress)
	if err != nil {
		return nil, err
	}

	chainID = sq.ChainID.Uint64()

	return &BasicBatchProcessor{State: s, stateRoot: stateRoot, SequencerAddress: sequencerAddress, SequencerChainID: chainID}, nil
}

// NewGenesisBatchProcessor creates a new batch processor
func (s *BasicState) NewGenesisBatchProcessor(genesisStateRoot []byte) (BatchProcessor, error) {
	s.Tree.SetCurrentRoot(genesisStateRoot)

	return &BasicBatchProcessor{State: s, stateRoot: genesisStateRoot}, nil
}

// GetStateRoot returns the root of the state tree
func (s *BasicState) GetStateRoot(ctx context.Context, virtual bool) ([]byte, error) {
	batch, err := s.GetLastBatch(ctx, virtual)
	if err != nil {
		return nil, err
	}

	if batch.Header == nil {
		return nil, ErrInvalidBatchHeader
	}

	return batch.Header.Root[:], nil
}

// GetStateRootByBatchNumber returns state root by batch number from the MT
func (s *BasicState) GetStateRootByBatchNumber(batchNumber uint64) ([]byte, error) {
	ctx := context.Background()
	batch, err := s.GetBatchByNumber(ctx, batchNumber)
	if err != nil {
		return nil, err
	}

	if batch.Header == nil {
		return nil, ErrInvalidBatchHeader
	}

	return batch.Header.Root[:], nil
}

// GetBalance from a given address
func (s *BasicState) GetBalance(address common.Address, batchNumber uint64) (*big.Int, error) {
	root, err := s.GetStateRootByBatchNumber(batchNumber)
	if err != nil {
		return nil, err
	}

	return s.Tree.GetBalance(address, root)
}

// EstimateGas for a transaction
func (s *BasicState) EstimateGas(transaction *types.Transaction) uint64 {
	// TODO: Calculate once we have txs that interact with SCs
	return 21000 //nolint:gomnd
}

// GetLastBlock gets the latest block
func (s *BasicState) GetLastBlock(ctx context.Context) (*Block, error) {
	var block Block
	err := s.db.QueryRow(ctx, getLastBlockSQL).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// GetPreviousBlock gets the offset previous block respect to latest
func (s *BasicState) GetPreviousBlock(ctx context.Context, offset uint64) (*Block, error) {
	var block Block
	err := s.db.QueryRow(ctx, getPreviousBlockSQL, offset).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// GetBlockByHash gets the block with the required hash
func (s *BasicState) GetBlockByHash(ctx context.Context, hash common.Hash) (*Block, error) {
	var block Block
	err := s.db.QueryRow(ctx, getBlockByHashSQL, hash).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// GetBlockByNumber gets the block with the required number
func (s *BasicState) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error) {
	var block Block
	err := s.db.QueryRow(ctx, getBlockByNumberSQL, blockNumber).Scan(&block.BlockNumber, &block.BlockHash, &block.ParentHash, &block.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// GetLastBlockNumber gets the latest block number
func (s *BasicState) GetLastBlockNumber(ctx context.Context) (uint64, error) {
	var lastBlockNum uint64
	err := s.db.QueryRow(ctx, getLastBlockNumberSQL).Scan(&lastBlockNum)
	if err != nil {
		return 0, err
	}
	return lastBlockNum, nil
}

// GetLastBatch gets the latest batch
func (s *BasicState) GetLastBatch(ctx context.Context, isVirtual bool) (*Batch, error) {
	var row pgx.Row

	if isVirtual {
		row = s.db.QueryRow(ctx, getLastVirtualBatchSQL)
	} else {
		row = s.db.QueryRow(ctx, getLastConsolidatedBatchSQL, common.Hash{})
	}

	var batch Batch
	err := row.Scan(
		&batch.BatchNumber, &batch.BatchHash, &batch.BlockNumber,
		&batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash,
		&batch.Header, &batch.Uncles, &batch.RawTxsData)

	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetPreviousBatch gets the offset previous batch respect to latest
func (s *BasicState) GetPreviousBatch(ctx context.Context, isVirtual bool, offset uint64) (*Batch, error) {
	var row pgx.Row
	if isVirtual {
		row = s.db.QueryRow(ctx, getPreviousVirtualBatchSQL, offset)
	} else {
		row = s.db.QueryRow(ctx, getPreviousConsolidatedBatchSQL, common.Hash{}, offset)
	}
	var batch Batch
	err := row.Scan(
		&batch.BatchNumber, &batch.BatchHash, &batch.BlockNumber,
		&batch.Sequencer, &batch.Aggregator, &batch.ConsolidatedTxHash, &batch.Header,
		&batch.Uncles, &batch.RawTxsData)

	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetBatchByHash gets the batch with the required hash
func (s *BasicState) GetBatchByHash(ctx context.Context, hash common.Hash) (*Batch, error) {
	var batch Batch
	err := s.db.QueryRow(ctx, getBatchByHashSQL, hash).Scan(
		&batch.BatchNumber, &batch.BatchHash, &batch.BlockNumber, &batch.Sequencer, &batch.Aggregator,
		&batch.ConsolidatedTxHash, &batch.Header, &batch.Uncles, &batch.RawTxsData)

	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetBatchByNumber gets the batch with the required number
func (s *BasicState) GetBatchByNumber(ctx context.Context, batchNumber uint64) (*Batch, error) {
	var batch Batch
	err := s.db.QueryRow(ctx, getBatchByNumberSQL, batchNumber).Scan(
		&batch.BatchNumber, &batch.BatchHash, &batch.BlockNumber, &batch.Sequencer, &batch.Aggregator,
		&batch.ConsolidatedTxHash, &batch.Header, &batch.Uncles, &batch.RawTxsData)
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetLastBatchNumber gets the latest batch number
func (s *BasicState) GetLastBatchNumber(ctx context.Context) (uint64, error) {
	var lastBatchNumber uint64
	err := s.db.QueryRow(ctx, getLastBatchNumberSQL).Scan(&lastBatchNumber)
	if err != nil {
		return 0, err
	}
	return lastBatchNumber, nil
}

// GetNonce returns the nonce of the given account at the given batch number
func (s *BasicState) GetNonce(address common.Address, batchNumber uint64) (uint64, error) {
	root, err := s.GetStateRootByBatchNumber(batchNumber)
	if err != nil {
		return 0, err
	}

	n, err := s.Tree.GetNonce(address, root)
	if err != nil {
		return 0, err
	}

	return n.Uint64(), nil
}

// GetTransactionByBatchHashAndIndex gets a transaction from a batch by index
func (s *BasicState) GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error) {
	panic("not implemented")
}

// GetTransactionByBatchNumberAndIndex gets a transaction from a batch by index
func (s *BasicState) GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error) {
	panic("not implemented")
}

// GetTransactionByHash gets a transaction by its hash
func (s *BasicState) GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error) {
	var encoded string
	if err := s.db.QueryRow(ctx, getTransactionByHashSQL, transactionHash).Scan(&encoded); err != nil {
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
func (s *BasicState) GetTransactionCount(ctx context.Context, fromAddress common.Address) (uint64, error) {
	var count uint64
	err := s.db.QueryRow(ctx, getTransactionCountSQL, fromAddress).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
func (s *BasicState) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*types.Receipt, error) {
	var receipt types.Receipt
	var blockNumber uint64
	err := s.db.QueryRow(ctx, getReceiptSQL, transactionHash).Scan(&receipt.Type, &receipt.PostState, &receipt.Status,
		&receipt.CumulativeGasUsed, &receipt.GasUsed, &blockNumber, &receipt.TxHash, &receipt.TransactionIndex)
	if err != nil {
		return nil, err
	}

	receipt.BlockNumber = new(big.Int).SetUint64(blockNumber)
	return &receipt, nil
}

// Reset resets the state to a block
func (s *BasicState) Reset(blockNumber uint64) error {
	panic("not implemented")
}

// ConsolidateBatch changes the virtual status of a batch
func (s *BasicState) ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash) error {
	if _, err := s.db.Exec(ctx, consolidateBatchSQL, consolidatedTxHash, batchNumber); err != nil {
		return err
	}
	return nil
}

// GetTxsByBatchNum returns all the txs in a given batch
func (s *BasicState) GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error) {
	rows, err := s.db.Query(ctx, getTxsByBatchNumSQL, batchNum)
	if err != nil {
		return nil, err
	}
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
func (s *BasicState) AddSequencer(ctx context.Context, seq Sequencer) error {
	_, err := s.db.Exec(ctx, addSequencerSQL, seq.Address, seq.URL, seq.ChainID.Uint64(), seq.BlockNumber)
	return err
}

// GetSequencer gets a sequencer
func (s *BasicState) GetSequencer(ctx context.Context, address common.Address) (*Sequencer, error) {
	var seq Sequencer
	var cID uint64
	err := s.db.QueryRow(ctx, getSequencerSQL, address.Bytes()).Scan(&seq.Address, &seq.URL, &cID, &seq.BlockNumber)
	if err != nil {
		return nil, err
	}

	seq.ChainID = big.NewInt(0).SetUint64(cID)

	return &seq, nil
}

// SetGenesis populates state with genesis information
func (s *BasicState) SetGenesis(ctx context.Context, genesis Genesis) error {
	// Generate Genesis Block
	block := &Block{
		BlockNumber: 0,
		BlockHash:   common.HexToHash("0x0000000000000"),
		ParentHash:  common.HexToHash("0x0000000000000"),
	}

	// Add Block
	err := s.AddBlock(ctx, block)
	if err != nil {
		return err
	}

	// reset tree current root
	s.Tree.SetCurrentRoot(nil)

	var root common.Hash

	// Genesis Balances
	for address, balance := range genesis.Balances {
		newRoot, _, err := s.Tree.SetBalance(address, balance)
		if err != nil {
			return err
		}
		root.SetBytes(newRoot)
	}

	// Generate Genesis Batch
	batch := &Batch{
		BatchNumber:        0,
		BlockNumber:        0,
		ConsolidatedTxHash: common.HexToHash("0x1"),
	}

	// Store batch into db
	bp, err := s.NewGenesisBatchProcessor(root[:])
	if err != nil {
		return err
	}
	err = bp.ProcessBatch(batch)
	if err != nil {
		return err
	}

	return nil
}

// AddBlock adds a new block to the State DB
func (s *BasicState) AddBlock(ctx context.Context, block *Block) error {
	_, err := s.db.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.Bytes(), block.ParentHash.Bytes(), block.ReceivedAt)
	return err
}

// SetLastBatchNumberSeenOnEthereum sets the last batch number that affected
// the roll-up in order to allow the components to know if the state
// is synchronized or not
func (s *BasicState) SetLastBatchNumberSeenOnEthereum(batchNumber uint64) error {
	return nil
}

// GetLastBatchNumberSeenOnEthereum returns the last batch number stored
// in the state that represents the last batch number that affected the
// roll-up in the Ethereum network.
func (s *BasicState) GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error) {
	return 0, nil
}
