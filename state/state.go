package state

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/db/statedb"
	"github.com/hermeznetwork/hermez-core/state/tree"
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
	GetLastConsolidatedBatchNumber(ctx context.Context) (uint64, error)
	GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error)
	GetTransactionCount(ctx context.Context, address common.Address) (uint64, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*types.Receipt, error)
	Reset(ctx context.Context, blockNumber uint64) error
	ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash) error
	GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error)
	AddSequencer(ctx context.Context, seq Sequencer) error
	GetSequencer(ctx context.Context, address common.Address) (*Sequencer, error)
	SetGenesis(ctx context.Context, genesis Genesis) error
	AddBlock(ctx context.Context, block *Block) error
	SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64) error
	GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error)
	GetStateRootByBatchNumber(batchNumber uint64) ([]byte, error)
}

var (
	// ErrInvalidBatchHeader indicates the batch header is invalid
	ErrInvalidBatchHeader = errors.New("invalid batch header")
)

// BasicState is a implementation of the state
type BasicState struct {
	cfg  Config
	db   statedb.StateDB
	tree tree.ReadWriter
}

// NewState creates a new State
func NewState(cfg Config, db *pgxpool.Pool, tree tree.ReadWriter) State {
	return &BasicState{cfg: cfg, db: statedb.newStateDB(db), tree: tree}
}

// NewBatchProcessor creates a new batch processor
func (s *BasicState) NewBatchProcessor(sequencerAddress common.Address, lastBatchNumber uint64) (BatchProcessor, error) {
	// init correct state root from previous batch
	stateRoot, err := s.GetStateRootByBatchNumber(lastBatchNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get state root for batch number %d, err: %v", lastBatchNumber, err)
	}

	s.tree.SetCurrentRoot(stateRoot)

	// Get Sequencer's Chain ID
	sq, err := s.db.GetSequencer(context.Background(), sequencerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequencer %s, err: %v", sequencerAddress, err)
	}

	return &BasicBatchProcessor{State: s, stateRoot: stateRoot, SequencerAddress: sequencerAddress, SequencerChainID: sq.ChainID.Uint64()}, nil
}

// NewGenesisBatchProcessor creates a new batch processor
func (s *BasicState) NewGenesisBatchProcessor(genesisStateRoot []byte) (BatchProcessor, error) {
	s.tree.SetCurrentRoot(genesisStateRoot)

	return &BasicBatchProcessor{State: s, stateRoot: genesisStateRoot}, nil
}

// GetStateRoot returns the root of the state tree
func (s *BasicState) GetStateRoot(ctx context.Context, virtual bool) ([]byte, error) {
	batch, err := s.db.GetLastBatch(ctx, virtual)
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

	return s.tree.GetBalance(address, root)
}

// EstimateGas for a transaction
func (s *BasicState) EstimateGas(transaction *types.Transaction) uint64 {
	// TODO: Calculate once we have txs that interact with SCs
	return 21000 //nolint:gomnd
}

// GetLastBlock gets the latest block
func (s *BasicState) GetLastBlock(ctx context.Context) (*Block, error) {
	return s.db.stateDB.GetLastBlock(ctx)
}

// GetPreviousBlock gets the offset previous block respect to latest
func (s *BasicState) GetPreviousBlock(ctx context.Context, offset uint64) (*Block, error) {
	return s.db.stateDB.GetPreviousBlock(ctx, offset)
}

// GetBlockByHash gets the block with the required hash
func (s *BasicState) GetBlockByHash(ctx context.Context, hash common.Hash) (*Block, error) {
	return s.db.stateDB.GetBlockByHash(ctx, hash)
}

// GetBlockByNumber gets the block with the required number
func (s *BasicState) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error) {
	return s.db.stateDB.GetBlockByNumber(ctx, blockNumber)
}

// GetLastBlockNumber gets the latest block number
func (s *BasicState) GetLastBlockNumber(ctx context.Context) (uint64, error) {
	return s.db.stateDB.GetLastBlockNumber(ctx)
}

// GetLastBatch gets the latest batch
func (s *BasicState) GetLastBatch(ctx context.Context, isVirtual bool) (*Batch, error) {
	return s.db.stateDB.GetLastBatch(ctx, isVirtual)
}

// GetPreviousBatch gets the offset previous batch respect to latest
func (s *BasicState) GetPreviousBatch(ctx context.Context, isVirtual bool, offset uint64) (*Batch, error) {
	return s.db.stateDB.GetPreviousBatch(ctx, isVirtual, offset)
}

// GetBatchByHash gets the batch with the required hash
func (s *BasicState) GetBatchByHash(ctx context.Context, hash common.Hash) (*Batch, error) {
	return s.db.stateDB.GetBatchByHash(ctx, hash)
}

// GetBatchByNumber gets the batch with the required number
func (s *BasicState) GetBatchByNumber(ctx context.Context, batchNumber uint64) (*Batch, error) {
	return s.db.stateDB.getBatchByNumber(ctx, batchNumber)
}

// GetLastBatchNumber gets the latest batch number
func (s *BasicState) GetLastBatchNumber(ctx context.Context) (uint64, error) {
	return s.db.stateDB.GetLastBatchNumber(ctx)
}

// GetLastConsolidatedBatchNumber gets the latest consolidated batch number
func (s *BasicState) GetLastConsolidatedBatchNumber(ctx context.Context) (uint64, error) {
	return s.db.stateDB.GetLastConsolidatedBatchNumber(ctx)
}

// GetTransactionByBatchHashAndIndex gets a transaction from a batch by index
func (s *BasicState) GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error) {
	return s.db.stateDB.GetTransactionByBatchHashAndIndex(ctx, batchHash, index)
}

// GetTransactionByBatchNumberAndIndex gets a transaction from a batch by index
func (s *BasicState) GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error) {
	return s.db.stateDB.GetTransactionByBatchNumberAndIndex(ctx, batchNumber, index)
}

// GetTransactionByHash gets a transaction by its hash
func (s *BasicState) GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error) {
	return s.db.stateDB.GetTransactionByHash(ctx, transactionHash)
}

// GetTransactionCount returns the number of transactions sent from an address
func (s *BasicState) GetTransactionCount(ctx context.Context, fromAddress common.Address) (uint64, error) {
	return s.db.stateDB.GetTransactionCount(ctx, fromAddress)
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
func (s *BasicState) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*types.Receipt, error) {
	return s.db.stateDB.GetTransactionReceipt(ctx, transactionHash)
}

// Reset resets the state to a block
func (s *BasicState) Reset(ctx context.Context, blockNumber uint64) error {
	return s.db.stateDB.Reset(ctx, blockNumber)
}

// ConsolidateBatch changes the virtual status of a batch
func (s *BasicState) ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash) error {
	return s.db.stateDB.ConsolidateBatch(ctx, batchNumber, consolidatedTxHash)
}

// GetTxsByBatchNum returns all the txs in a given batch
func (s *BasicState) GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error) {
	return s.db.stateDB.GetTxsByBatchNum(ctx, batchNum)
}

// AddSequencer stores a new sequencer
func (s *BasicState) AddSequencer(ctx context.Context, seq Sequencer) error {
	return s.db.stateDB.AddSequencer(ctx, seq)
}

// GetSequencer gets a sequencer
func (s *BasicState) GetSequencer(ctx context.Context, address common.Address) (*Sequencer, error) {
	return s.db.stateDB.GetSequencer(ctx, address)
}

// AddBlock adds a new block to the State Store
func (s *BasicState) AddBlock(ctx context.Context, block *Block) error {
	return s.db.stateDB.AddBlock(ctx, block)
}

// SetLastBatchNumberSeenOnEthereum sets the last batch number that affected
// the roll-up in order to allow the components to know if the state
// is synchronized or not
func (s *BasicState) SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64) error {
	return s.db.stateDB.SetLastBatchNumberSeenOnEthereum(ctx, batchNumber)
}

// GetLastBatchNumberSeenOnEthereum returns the last batch number stored
// in the state that represents the last batch number that affected the
// roll-up in the Ethereum network.
func (s *BasicState) GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error) {
	return s.db.stateDB.GetLastBatchNumberSeenOnEthereum(ctx)
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
	err := s.db.AddBlock(ctx, block)
	if err != nil {
		return err
	}

	// reset tree current root
	s.tree.SetCurrentRoot(nil)

	var root common.Hash

	// Genesis Balances
	for address, balance := range genesis.Balances {
		newRoot, _, err := s.tree.SetBalance(address, balance)
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
		MaticCollateral:    big.NewInt(0),
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

// GetNonce returns the nonce of the given account at the given batch number
func (s *BasicState) GetNonce(address common.Address, batchNumber uint64) (uint64, error) {
	root, err := s.GetStateRootByBatchNumber(batchNumber)
	if err != nil {
		return 0, err
	}

	n, err := s.tree.GetNonce(address, root)
	if err != nil {
		return 0, err
	}

	return n.Uint64(), nil
}
