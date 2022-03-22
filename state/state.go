package state

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/evm"
	"github.com/hermeznetwork/hermez-core/state/tree"
)

const (
	// TxTransferGas used for TXs that do not create a contract
	TxTransferGas uint64 = 21000
)

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

// State is a implementation of the state
type State struct {
	cfg  Config
	tree merkletree
	storage
}

// NewState creates a new State
func NewState(cfg Config, storage storage, tree merkletree) *State {
	return &State{cfg: cfg, tree: tree, storage: storage}
}

// BeginStateTransaction starts a transaction block
func (s *State) BeginStateTransaction(ctx context.Context) error {
	err := s.storage.BeginDBTransaction(ctx)
	if err != nil {
		return err
	}
	return s.tree.BeginDBTransaction(ctx)
}

// CommitState commits a state into db
func (s *State) CommitState(ctx context.Context) error {
	err := s.storage.Commit(ctx)
	if err != nil {
		return err
	}
	return s.tree.Commit(ctx)
}

// RollbackState rollbacks a db state transaction
func (s *State) RollbackState(ctx context.Context) error {
	err := s.storage.Rollback(ctx)
	if err != nil {
		return err
	}
	return s.tree.Rollback(ctx)
}

// NewBatchProcessor creates a new batch processor
func (s *State) NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, lastBatchNumber uint64) (*BasicBatchProcessor, error) {
	// init correct state root from previous batch
	stateRoot, err := s.GetStateRootByBatchNumber(ctx, lastBatchNumber)
	if err != nil {
		return nil, err
	}

	// Get Sequencer's Chain ID
	chainID := s.cfg.DefaultChainID
	sq, err := s.GetSequencer(ctx, sequencerAddress)
	if err == nil {
		chainID = sq.ChainID.Uint64()
	}

	lastBatch, err := s.GetBatchByNumber(ctx, lastBatchNumber)
	if err != ErrNotFound && err != nil {
		return nil, err
	}

	transactionContext := transactionContext{difficulty: new(big.Int)}

	batchProcessor := &BasicBatchProcessor{State: s, stateRoot: stateRoot, SequencerAddress: sequencerAddress, SequencerChainID: chainID, LastBatch: lastBatch, MaxCumulativeGasUsed: s.cfg.MaxCumulativeGasUsed, transactionContext: transactionContext}
	batchProcessor.setRuntime(evm.NewEVM())
	blockNumber, err := s.GetLastBlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	batchProcessor.forks = runtime.AllForksEnabled.At(blockNumber)
	return batchProcessor, nil
}

// NewGenesisBatchProcessor creates a new batch processor
func (s *State) NewGenesisBatchProcessor(genesisStateRoot []byte) (*BasicBatchProcessor, error) {
	return &BasicBatchProcessor{State: s, stateRoot: genesisStateRoot}, nil
}

// GetStateRoot returns the root of the state tree
func (s *State) GetStateRoot(ctx context.Context, virtual bool) ([]byte, error) {
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
func (s *State) GetStateRootByBatchNumber(ctx context.Context, batchNumber uint64) ([]byte, error) {
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
func (s *State) GetBalance(ctx context.Context, address common.Address, batchNumber uint64) (*big.Int, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber)
	if err != nil {
		return nil, err
	}

	return s.tree.GetBalance(ctx, address, root)
}

// GetCode from a given address
func (s *State) GetCode(ctx context.Context, address common.Address, batchNumber uint64) ([]byte, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber)
	if err != nil {
		return nil, err
	}

	return s.tree.GetCode(ctx, address, root)
}

// EstimateGas for a transaction
func (s *State) EstimateGas(transaction *types.Transaction) (uint64, error) {
	ctx := context.Background()
	sequencerAddress := common.Address{}
	lastBatch, err := s.GetLastBatch(ctx, true)
	if err != nil {
		log.Errorf("failed to get last batch from the state, err: %v", err)
		return 0, err
	}
	bp, err := s.NewBatchProcessor(ctx, sequencerAddress, lastBatch.Number().Uint64())
	if err != nil {
		log.Errorf("failed to get create a new batch processor, err: %v", err)
		return 0, err
	}
	result := bp.estimateGas(ctx, transaction, sequencerAddress)
	return result.GasUsed, result.Err
}

// SetGenesis populates state with genesis information
func (s *State) SetGenesis(ctx context.Context, genesis Genesis) error {
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

	var (
		root    common.Hash
		newRoot []byte
	)

	if genesis.Balances != nil { // Genesis Balances
		for address, balance := range genesis.Balances {
			newRoot, _, err = s.tree.SetBalance(ctx, address, balance, newRoot)
			if err != nil {
				return err
			}
			root.SetBytes(newRoot)
		}
	} else { // Genesis Smart Contracts
		for address, sc := range genesis.SmartContracts {
			newRoot, _, err = s.tree.SetCode(ctx, address, sc, newRoot)
			if err != nil {
				return err
			}
			root.SetBytes(newRoot)
		}
	}

	for address, storage := range genesis.Storage {
		for key, value := range storage {
			newRoot, _, err = s.tree.SetStorageAt(ctx, address, key, value, newRoot)
			if err != nil {
				return err
			}
		}
	}
	root.SetBytes(newRoot)

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
	err = bp.ProcessBatch(ctx, batch)
	if err != nil {
		return err
	}

	return nil
}

// GetNonce returns the nonce of the given account at the given batch number
func (s *State) GetNonce(ctx context.Context, address common.Address, batchNumber uint64) (uint64, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber)
	if err != nil {
		return 0, err
	}

	n, err := s.tree.GetNonce(ctx, address, root)
	if errors.Is(err, tree.ErrNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return n.Uint64(), nil
}

// GetStorageAt from a given address
func (s *State) GetStorageAt(ctx context.Context, address common.Address, position common.Hash, batchNumber uint64) (*big.Int, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber)
	if err != nil {
		return nil, err
	}

	return s.tree.GetStorageAt(ctx, address, position, root)
}
