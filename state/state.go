package state

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state/tree"
)

// State is the interface of the Hermez state
type State interface {
	NewBatchProcessor(sequencerAddress common.Address, lastBatchNumber uint64) (BatchProcessor, error)
	NewGenesisBatchProcessor(genesisStateRoot []byte) (BatchProcessor, error)
	GetStateRoot(ctx context.Context, virtual bool) ([]byte, error)
	GetBalance(address common.Address, batchNumber uint64) (*big.Int, error)
	EstimateGas(transaction *types.Transaction) uint64
	GetNonce(address common.Address, batchNumber uint64) (uint64, error)
	SetGenesis(ctx context.Context, genesis Genesis) error
	GetStateRootByBatchNumber(batchNumber uint64) ([]byte, error)
	Storage
}

// Storage is the interface of the Hermez state methods that access database
type Storage interface {
	GetLastBlock(ctx context.Context) (*Block, error)
	GetPreviousBlock(ctx context.Context, offset uint64) (*Block, error)
	GetBlockByHash(ctx context.Context, hash common.Hash) (*Block, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error)
	GetLastBlockNumber(ctx context.Context) (uint64, error)
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
	ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash, consolidatedAt time.Time) error
	GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error)
	AddSequencer(ctx context.Context, seq Sequencer) error
	GetSequencer(ctx context.Context, address common.Address) (*Sequencer, error)
	AddBlock(ctx context.Context, block *Block) error
	SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64) error
	GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error)
	AddBatch(ctx context.Context, batch *Batch) error
	AddTransaction(ctx context.Context, tx *types.Transaction, batchNumber uint64, index uint) error
	AddReceipt(ctx context.Context, receipt *types.Receipt) error
}

var (
	// ErrInvalidBatchHeader indicates the batch header is invalid
	ErrInvalidBatchHeader = errors.New("invalid batch header")
)

// BasicState is a implementation of the state
type BasicState struct {
	cfg  Config
	tree tree.ReadWriter
	Storage
}

// NewState creates a new State
func NewState(cfg Config, storage Storage, tree tree.ReadWriter) State {
	return &BasicState{cfg: cfg, tree: tree, Storage: storage}
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
	chainId := s.cfg.DefaultChainID
	sq, err := s.GetSequencer(context.Background(), sequencerAddress)
	if err == nil {
		chainId = sq.ChainID.Uint64()
	}

	return &BasicBatchProcessor{State: s, stateRoot: stateRoot, SequencerAddress: sequencerAddress, SequencerChainID: chainId}, nil
}

// NewGenesisBatchProcessor creates a new batch processor
func (s *BasicState) NewGenesisBatchProcessor(genesisStateRoot []byte) (BatchProcessor, error) {
	s.tree.SetCurrentRoot(genesisStateRoot)

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

	return s.tree.GetBalance(address, root)
}

// EstimateGas for a transaction
func (s *BasicState) EstimateGas(transaction *types.Transaction) uint64 {
	// TODO: Calculate once we have txs that interact with SCs
	return 21000 //nolint:gomnd
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
