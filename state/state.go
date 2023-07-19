package state

import (
	"context"
	"math/big"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	// ZeroHash is the hash 0x0000000000000000000000000000000000000000000000000000000000000000
	ZeroHash = common.Hash{}
	// ZeroAddress is the address 0x0000000000000000000000000000000000000000
	ZeroAddress = common.Address{}
)

// State is an implementation of the state
type State struct {
	cfg Config
	*PostgresStorage
	executorClient executor.ExecutorServiceClient
	tree           *merkletree.StateTree
	eventLog       *event.EventLog

	lastL2BlockSeen         types.Block
	newL2BlockEvents        chan NewL2BlockEvent
	newL2BlockEventHandlers []NewL2BlockEventHandler
}

// NewState creates a new State
func NewState(cfg Config, storage *PostgresStorage, executorClient executor.ExecutorServiceClient, stateTree *merkletree.StateTree, eventLog *event.EventLog) *State {
	var once sync.Once
	once.Do(func() {
		metrics.Register()
	})

	state := &State{
		cfg:                     cfg,
		PostgresStorage:         storage,
		executorClient:          executorClient,
		tree:                    stateTree,
		eventLog:                eventLog,
		newL2BlockEvents:        make(chan NewL2BlockEvent),
		newL2BlockEventHandlers: []NewL2BlockEventHandler{},
	}

	return state
}

// BeginStateTransaction starts a state transaction
func (s *State) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetBalance from a given address
func (s *State) GetBalance(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error) {
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	return s.tree.GetBalance(ctx, address, root.Bytes())
}

// GetCode from a given address
func (s *State) GetCode(ctx context.Context, address common.Address, root common.Hash) ([]byte, error) {
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	return s.tree.GetCode(ctx, address, root.Bytes())
}

// GetNonce returns the nonce of the given account at the given block number
func (s *State) GetNonce(ctx context.Context, address common.Address, root common.Hash) (uint64, error) {
	if s.tree == nil {
		return 0, ErrStateTreeNil
	}
	nonce, err := s.tree.GetNonce(ctx, address, root.Bytes())
	if err != nil {
		return 0, err
	}
	return nonce.Uint64(), nil
}

// GetStorageAt from a given address
func (s *State) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root common.Hash) (*big.Int, error) {
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	return s.tree.GetStorageAt(ctx, address, position, root.Bytes())
}

// GetLastStateRoot returns the latest state root
func (s *State) GetLastStateRoot(ctx context.Context, dbTx pgx.Tx) (common.Hash, error) {
	lastBlockHeader, err := s.GetLastL2BlockHeader(ctx, dbTx)
	if err != nil {
		return common.Hash{}, err
	}
	return lastBlockHeader.Root, nil
}

// GetBalanceByStateRoot gets balance from the MT Service using the provided state root
func (s *State) GetBalanceByStateRoot(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error) {
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	balance, err := s.tree.GetBalance(ctx, address, root.Bytes())
	if err != nil && balance == nil {
		balance = big.NewInt(0)
	}
	return balance, err
}

// GetNonceByStateRoot gets nonce from the MT Service using the provided state root
func (s *State) GetNonceByStateRoot(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error) {
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	return s.tree.GetNonce(ctx, address, root.Bytes())
}

// GetTree returns State inner tree
func (s *State) GetTree() *merkletree.StateTree {
	return s.tree
}

// FlushMerkleTree persists updates in the Merkle tree
func (s *State) FlushMerkleTree(ctx context.Context) error {
	if s.tree == nil {
		return ErrStateTreeNil
	}
	return s.tree.Flush(ctx)
}

// GetStoredFlushID returns the stored flush ID and Prover ID
func (s *State) GetStoredFlushID(ctx context.Context) (uint64, string, error) {
	if s.executorClient == nil {
		return 0, "", ErrExecutorNil
	}
	res, err := s.executorClient.GetFlushStatus(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, "", err
	}

	return res.StoredFlushId, res.ProverId, nil
}
