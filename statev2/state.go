package state

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/umbracle/ethgo/abi"
)

const (
	// TxTransferGas used for TXs that do not create a contract
	TxTransferGas uint64 = 21000
	// TxSmartContractCreationGas used for TXs that create a contract
	TxSmartContractCreationGas uint64 = 53000
	half                       uint64 = 2
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
	cfg  Config
	tree statetree
	*PostgresStorage

	mu    *sync.Mutex
	dbTxs map[string]bool
}

// NewState creates a new State
func NewState(cfg Config, storage *PostgresStorage, tree statetree) *State {
	return &State{
		cfg:             cfg,
		tree:            tree,
		PostgresStorage: storage,

		mu:    new(sync.Mutex),
		dbTxs: make(map[string]bool),
	}
}

// BeginStateTransaction starts a transaction block
func (s *State) BeginStateTransaction(ctx context.Context) (string, error) {
	const maxAttempts = 3
	var (
		txBundleID string
		found      bool
	)

	s.mu.Lock()
	for i := 0; i < maxAttempts; i++ {
		txBundleID = uuid.NewString()
		_, idExists := s.dbTxs[txBundleID]
		if !idExists {
			found = true
			break
		}
	}
	s.mu.Unlock()

	if !found {
		return "", fmt.Errorf("could not find unused uuid for db tx bundle")
	}

	if err := s.PostgresStorage.BeginDBTransaction(ctx, txBundleID); err != nil {
		return "", err
	}
	if err := s.tree.BeginDBTransaction(ctx, txBundleID); err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.dbTxs[txBundleID] = true

	return txBundleID, nil
}

// CommitState commits a state into db
func (s *State) CommitState(ctx context.Context, txBundleID string) error {
	if err := s.tree.Commit(ctx, txBundleID); err != nil {
		return err
	}

	if err := s.PostgresStorage.Commit(ctx, txBundleID); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.dbTxs, txBundleID)
	return nil
}

// RollbackState rollbacks a db state transaction
func (s *State) RollbackState(ctx context.Context, txBundleID string) error {
	if err := s.tree.Rollback(ctx, txBundleID); err != nil {
		return err
	}

	if err := s.PostgresStorage.Rollback(ctx, txBundleID); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.dbTxs, txBundleID)
	return nil
}

// ProcessSequence process sequence of the txs
// TODO: implement function
func (s *State) ProcessBatchAndStoreLastTx(ctx context.Context, txs []types.Transaction) *runtime.ExecutionResult {
	return &runtime.ExecutionResult{}
}

// ResetDB resets the state to block for the given DB tx bundle.
func (s *State) ResetDB(ctx context.Context, block *Block, txBundleID string) error {
	return s.PostgresStorage.Reset(ctx, block, txBundleID)
}

func constructErrorFromRevert(result *runtime.ExecutionResult) error {
	revertErrMsg, unpackErr := abi.UnpackRevertError(result.ReturnValue)
	if unpackErr != nil {
		return result.Err
	}

	return fmt.Errorf("%w: %s", result.Err, revertErrMsg)
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