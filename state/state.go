package state

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/evm"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/jackc/pgx/v4"
)

const (
	// TxGas used for TXs that do not create a contract
	TxGas uint64 = 21000
	// TxGasContractCreation user for transactions that create a contract
	TxGasContractCreation uint64 = 53000
)

// State is the interface of the Hermez state
type State interface {
	NewBatchProcessor(sequencerAddress common.Address, lastBatchNumber uint64) (BatchProcessor, error)
	NewGenesisBatchProcessor(genesisStateRoot []byte) (BatchProcessor, error)
	GetStateRoot(ctx context.Context, virtual bool) ([]byte, error)
	GetBalance(address common.Address, batchNumber uint64) (*big.Int, error)
	GetCode(address common.Address, batchNumber uint64) ([]byte, error)
	EstimateGas(transaction *types.Transaction) uint64
	GetNonce(address common.Address, batchNumber uint64) (uint64, error)
	SetGenesis(ctx context.Context, genesis Genesis) error
	GetStateRootByBatchNumber(batchNumber uint64) ([]byte, error)
	Storage
}

// Storage is the interface of the Hermez state methods that access database
type Storage interface {
	BeginDBTransaction(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetLastBlock(ctx context.Context) (*Block, error)
	GetPreviousBlock(ctx context.Context, offset uint64) (*Block, error)
	GetBlockByHash(ctx context.Context, hash common.Hash) (*Block, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error)
	GetLastBlockNumber(ctx context.Context) (uint64, error)
	GetLastBatch(ctx context.Context, isVirtual bool) (*Batch, error)
	GetPreviousBatch(ctx context.Context, isVirtual bool, offset uint64) (*Batch, error)
	GetBatchByHash(ctx context.Context, hash common.Hash) (*Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64) (*Batch, error)
	GetBatchHeader(ctx context.Context, batchNumber uint64) (*types.Header, error)
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	GetLastConsolidatedBatchNumber(ctx context.Context) (uint64, error)
	GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error)
	GetTransactionCount(ctx context.Context, address common.Address) (uint64, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*Receipt, error)
	Reset(ctx context.Context, blockNumber uint64) error
	ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash, consolidatedAt time.Time, aggregator common.Address) error
	GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error)
	AddSequencer(ctx context.Context, seq Sequencer) error
	GetSequencer(ctx context.Context, address common.Address) (*Sequencer, error)
	AddBlock(ctx context.Context, block *Block) error
	SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64) error
	GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error)
	AddBatch(ctx context.Context, batch *Batch) error
	AddTransaction(ctx context.Context, tx *types.Transaction, batchNumber uint64, index uint) error
	AddReceipt(ctx context.Context, receipt *Receipt) error
	SetLastBatchNumberConsolidatedOnEthereum(ctx context.Context, batchNumber uint64) error
	GetLastBatchNumberConsolidatedOnEthereum(ctx context.Context) (uint64, error)
}

var (
	// ErrInvalidBatchHeader indicates the batch header is invalid
	ErrInvalidBatchHeader = errors.New("invalid batch header")
	// ErrStateNotSynchronized indicates the state database may be empty
	ErrStateNotSynchronized = errors.New("state not synchronized")
	// ErrNotFound indicates an object has not been found for the search criteria used
	ErrNotFound = errors.New("object not found")
	// ErrNilDBTransaction indicates the db transaction has not been properly initialized
	ErrNilDBTransaction = errors.New("database transaction not properly initialized")
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
		return nil, err
	}

	s.tree.SetCurrentRoot(stateRoot)

	// Get Sequencer's Chain ID
	chainID := s.cfg.DefaultChainID
	sq, err := s.GetSequencer(context.Background(), sequencerAddress)
	if err == nil {
		chainID = sq.ChainID.Uint64()
	}

	lastBatch, err := s.GetBatchByNumber(context.Background(), lastBatchNumber)
	if err != ErrNotFound && err != nil {
		return nil, err
	}

	batchProcessor := &BasicBatchProcessor{State: s, stateRoot: stateRoot, SequencerAddress: sequencerAddress, SequencerChainID: chainID, LastBatch: lastBatch}
	batchProcessor.setRuntime(evm.NewEVM())
	blockNumber, err := s.GetLastBlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	batchProcessor.forks = runtime.AllForksEnabled.At(blockNumber)
	return batchProcessor, nil
}

// NewGenesisBatchProcessor creates a new batch processor
func (s *BasicState) NewGenesisBatchProcessor(genesisStateRoot []byte) (BatchProcessor, error) {
	s.tree.SetCurrentRoot(genesisStateRoot)

	return &BasicBatchProcessor{State: s, stateRoot: genesisStateRoot}, nil
}

// GetStateRoot returns the root of the state tree
func (s *BasicState) GetStateRoot(ctx context.Context, virtual bool) ([]byte, error) {
	batch, err := s.GetLastBatch(ctx, virtual)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
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

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
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
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	return s.tree.GetBalance(address, root)
}

// GetCode from a given address
func (s *BasicState) GetCode(address common.Address, batchNumber uint64) ([]byte, error) {
	root, err := s.GetStateRootByBatchNumber(batchNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}

	return s.tree.GetCode(address, root)
}

// EstimateGas for a transaction
func (s *BasicState) EstimateGas(transaction *types.Transaction) uint64 {
	// TODO: Calculate once we have txs that interact with SCs
	return TxGas
}

// SetGenesis populates state with genesis information
func (s *BasicState) SetGenesis(ctx context.Context, genesis Genesis) error {
	// Generate Genesis Block
	block := &Block{
		BlockNumber: genesis.Block.NumberU64(),
		BlockHash:   genesis.Block.Hash(),
		ParentHash:  genesis.Block.ParentHash(),
		ReceivedAt:  genesis.Block.ReceivedAt,
	}

	// Add Block
	err := s.AddBlock(ctx, block)
	if err != nil {
		return err
	}

	// reset tree current root
	s.tree.SetCurrentRoot(nil)

	var root common.Hash

	if genesis.Balances != nil { // Genesis Balances
		for address, balance := range genesis.Balances {
			newRoot, _, err := s.tree.SetBalance(address, balance)
			if err != nil {
				return err
			}
			root.SetBytes(newRoot)
		}
	} else { // Genesis Smart Contracts
		for address, sc := range genesis.SmartContracts {
			newRoot, _, err := s.tree.SetCode(address, sc)
			if err != nil {
				return err
			}
			root.SetBytes(newRoot)
		}
	}

	// Generate Genesis Batch
	receivedAt := genesis.Block.ReceivedAt
	batch := &Batch{
		Header: &types.Header{
			Number: big.NewInt(0),
		},
		BlockNumber:        genesis.Block.NumberU64(),
		ConsolidatedTxHash: common.HexToHash("0x1"),
		ConsolidatedAt:     &receivedAt,
		MaticCollateral:    big.NewInt(0),
		ReceivedAt:         time.Now(),
		ChainID:            new(big.Int).SetUint64(genesis.L2ChainID),
		GlobalExitRoot:     common.Hash{},
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
