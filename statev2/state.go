package statev2

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

const (
	// TxTransferGas used for TXs that do not create a contract
	TxTransferGas uint64 = 21000
	// TxSmartContractCreationGas used for TXs that create a contract
	TxSmartContractCreationGas uint64 = 53000
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
	// ErrAlreadyInitializedDBTransaction indicates the db transaction was already initialized
	ErrAlreadyInitializedDBTransaction = errors.New("database transaction already initialized")
	// ErrNotEnoughIntrinsicGas indicates the gas is not enough to cover the intrinsic gas cost
	ErrNotEnoughIntrinsicGas = fmt.Errorf("not enough gas supplied for intrinsic gas costs")
)

// State is a implementation of the state
type State struct {
	cfg Config
	*PostgresStorage

	mu    *sync.Mutex
	dbTxs map[string]bool
}

// NewState creates a new State
func NewState(cfg Config, storage *PostgresStorage) *State {
	return &State{
		cfg:             cfg,
		PostgresStorage: storage,

		mu:    new(sync.Mutex),
		dbTxs: make(map[string]bool),
	}
}

// ResetDB resets the state to block for the given DB tx bundle.
func (s *State) ResetDB(ctx context.Context, block *Block, txBundleID string) error {
	return s.PostgresStorage.Reset(ctx, block, txBundleID)
}

func (s *State) AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, txBundleID string) error {
	return s.PostgresStorage.AddGlobalExitRoot(ctx, exitRoot, txBundleID)
}

func (s *State) GetLatestGlobalExitRoot(ctx context.Context, txBundleID string) (*GlobalExitRoot, error) {
	return s.PostgresStorage.GetLatestGlobalExitRoot(ctx, txBundleID)
}

func (s *State) AddForcedBatch(ctx context.Context, forcedBatch *ForcedBatch, txBundleID string) error {
	return s.PostgresStorage.AddForcedBatch(ctx, forcedBatch, txBundleID)
}

func (s *State) GetForcedBatch(ctx context.Context, txBundleID string, forcedBatchNumber uint64) (*ForcedBatch, error) {
	return s.PostgresStorage.GetForcedBatch(ctx, txBundleID, forcedBatchNumber)
}

// AddBlock adds a new block to the State Store.
func (s *State) AddBlock(ctx context.Context, block *Block, txBundleID string) error {
	return s.PostgresStorage.AddBlock(ctx, block, txBundleID)
}
