package statev2

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/jackc/pgx/v4"
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
	*PostgresStorage
}

// NewState creates a new State
func NewState(storage *PostgresStorage) *State {
	return &State{
		PostgresStorage: storage,
	}
}

// BeginDBTransaction starts a state transaction
func (s *State) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// Commit commits a state transaction
func (s *State) CommitStateTransaction(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	return err
}

// Rollback rollbacks a state transaction
func (s *State) RollbackStateTransaction(ctx context.Context, tx pgx.Tx) error {
	err := tx.Rollback(ctx)
	return err
}

// ResetDB resets the state to block for the given DB tx .
func (s *State) ResetDB(ctx context.Context, block *Block, tx pgx.Tx) error {
	return s.PostgresStorage.Reset(ctx, block, tx)
}

func (s *State) AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, tx pgx.Tx) error {
	return s.PostgresStorage.AddGlobalExitRoot(ctx, exitRoot, tx)
}

func (s *State) GetLatestGlobalExitRoot(ctx context.Context, tx pgx.Tx) (*GlobalExitRoot, error) {
	return s.PostgresStorage.GetLatestGlobalExitRoot(ctx, tx)
}

func (s *State) AddForcedBatch(ctx context.Context, forcedBatch *ForcedBatch, tx pgx.Tx) error {
	return s.PostgresStorage.AddForcedBatch(ctx, forcedBatch, tx)
}

func (s *State) GetForcedBatch(ctx context.Context, tx pgx.Tx, forcedBatchNumber uint64) (*ForcedBatch, error) {
	return s.PostgresStorage.GetForcedBatch(ctx, tx, forcedBatchNumber)
}

// AddBlock adds a new block to the State Store.
func (s *State) AddBlock(ctx context.Context, block *Block, tx pgx.Tx) error {
	return s.PostgresStorage.AddBlock(ctx, block, tx)
}

// ProcessSequence process sequence of the txs
// TODO: implement function
func (s *State) ProcessBatchAndStoreLastTx(ctx context.Context, txs []types.Transaction) *runtime.ExecutionResult {
	return &runtime.ExecutionResult{}
}

// GetLastL1InteractionTime get time from last l1 interaction time
// TODO: implement function
func (s *State) GetLastL1InteractionTime(ctx context.Context) (time.Time, error) {
	return time.Now(), nil
}

// GetNumberOfBlocksSinceLastGERUpdate get time from last time get
// TODO: implement function
func (s *State) GetNumberOfBlocksSinceLastGERUpdate(ctx context.Context) (uint32, error) {
	return 0, nil
}

// GetLastBatchTime get last batch time
// TODO: implement function
func (s *State) GetLastBatchTime(ctx context.Context) (time.Time, error) {
	return time.Now(), nil
}

// AddVerifiedBatch adds a new VerifiedBatch to the db
func (s *State) AddVerifiedBatch(ctx context.Context, verifiedBatch *VerifiedBatch, tx pgx.Tx) error {
	return s.PostgresStorage.AddVerifiedBatch(ctx, verifiedBatch, tx)
}

// GetVerifiedBatch get an L1 verifiedBatch.
func (s *State) GetVerifiedBatch(ctx context.Context, tx pgx.Tx, batchNumber uint64) (*VerifiedBatch, error) {
	return s.PostgresStorage.GetVerifiedBatch(ctx, tx, batchNumber)
}
