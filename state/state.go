package state

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/helper"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/evm"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation"
	"github.com/hermeznetwork/hermez-core/state/tree"
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
		return "", fmt.Errorf("Could not find unused uuid for db tx bundle")
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

// NewBatchProcessor creates a new batch processor
func (s *State) NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, stateRoot []byte, txBundleID string) (*BatchProcessor, error) {
	// Get Sequencer's Chain ID
	chainID := s.cfg.DefaultChainID
	sq, err := s.GetSequencer(ctx, sequencerAddress, txBundleID)
	if err == nil {
		chainID = sq.ChainID.Uint64()
	}

	lastBatch, err := s.GetLastBatchByStateRoot(ctx, stateRoot, txBundleID)
	if err != ErrNotFound && err != nil {
		return nil, err
	}

	host := Host{State: s, stateRoot: stateRoot, transactionContext: transactionContext{difficulty: new(big.Int)}, txBundleID: txBundleID}
	host.setRuntime(evm.NewEVM())
	blockNumber, err := s.GetLastBlockNumber(ctx, txBundleID)
	if err != nil {
		return nil, err
	}
	host.forks = runtime.AllForksEnabled.At(blockNumber)

	batchProcessor := &BatchProcessor{SequencerAddress: sequencerAddress, SequencerChainID: chainID, LastBatch: lastBatch, MaxCumulativeGasUsed: s.cfg.MaxCumulativeGasUsed, Host: host}
	batchProcessor.Host.setRuntime(evm.NewEVM())

	return batchProcessor, nil
}

// NewGenesisBatchProcessor creates a new batch processor
func (s *State) NewGenesisBatchProcessor(genesisStateRoot []byte, txBundleID string) (*BatchProcessor, error) {
	host := Host{State: s, stateRoot: genesisStateRoot, transactionContext: transactionContext{difficulty: new(big.Int)}, txBundleID: txBundleID}
	host.setRuntime(evm.NewEVM())
	host.forks = runtime.AllForksEnabled.At(0)
	return &BatchProcessor{Host: host}, nil
}

// GetStateRoot returns the root of the state tree
func (s *State) GetStateRoot(ctx context.Context, virtual bool, txBundleID string) ([]byte, error) {
	batch, err := s.GetLastBatch(ctx, virtual, txBundleID)
	if err != nil {
		return nil, err
	}

	if batch.Header == nil {
		return nil, ErrInvalidBatchHeader
	}

	return batch.Header.Root[:], nil
}

// GetStateRootByBatchNumber returns state root by batch number from the MT
func (s *State) GetStateRootByBatchNumber(ctx context.Context, batchNumber uint64, txBundleID string) ([]byte, error) {
	batch, err := s.GetBatchByNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return nil, err
	}

	if batch.Header == nil {
		return nil, ErrInvalidBatchHeader
	}

	return batch.Header.Root[:], nil
}

// GetBalance from a given address
func (s *State) GetBalance(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) (*big.Int, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return nil, err
	}

	return s.tree.GetBalance(ctx, address, root, txBundleID)
}

// GetCode from a given address
func (s *State) GetCode(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) ([]byte, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return nil, err
	}

	return s.tree.GetCode(ctx, address, root, txBundleID)
}

// EstimateGas for a transaction
func (s *State) EstimateGas(transaction *types.Transaction, txBundleID string) (uint64, error) {
	ctx := context.Background()
	sequencerAddress := common.Address{}
	lastBatch, err := s.GetLastBatch(ctx, true, txBundleID)
	if err != nil {
		log.Errorf("failed to get last batch from the state, err: %v", err)
		return 0, err
	}
	bp, err := s.NewBatchProcessor(ctx, sequencerAddress, lastBatch.Header.Root[:], txBundleID)
	if err != nil {
		log.Errorf("failed to get create a new batch processor, err: %v", err)
		return 0, err
	}
	result := bp.estimateGas(ctx, transaction)
	return result.GasUsed, result.Err
}

// ReplayTransaction gets trace by rexecuting a transaction
func (s *State) ReplayTransaction(transactionHash common.Hash) *runtime.ExecutionResult {
	ctx := context.Background()

	txBundleID, err := s.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("trace transaction: failed to begin db transaction, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return &runtime.ExecutionResult{Err: err}
	}

	tx, err := s.GetTransactionByHash(ctx, transactionHash, txBundleID)
	if err != nil {
		log.Errorf("trace transaction: failed to get transaction by hash, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return &runtime.ExecutionResult{Err: err}
	}

	receipt, err := s.GetTransactionReceipt(ctx, transactionHash, txBundleID)
	if err != nil {
		log.Errorf("trace transaction: failed to get receipt by tx hash, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return &runtime.ExecutionResult{Err: err}
	}

	batch, err := s.GetBatchByHash(ctx, receipt.BlockHash, txBundleID)
	if err != nil {
		log.Errorf("trace transaction: failed to get batch by hash, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return &runtime.ExecutionResult{Err: err}
	}

	var stateRoot []byte

	if receipt.TransactionIndex > 0 {
		previousTX, err := s.GetTransactionByBatchHashAndIndex(ctx, receipt.BlockHash, uint64(receipt.TransactionIndex-1), txBundleID)
		if err != nil {
			log.Errorf("trace transaction: failed to get previous tx, err: %v", err)
			rbErr := s.RollbackState(ctx, txBundleID)
			if rbErr != nil {
				log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
			}
			return &runtime.ExecutionResult{Err: err}
		}

		previousReceipt, err := s.GetTransactionReceipt(ctx, previousTX.Hash(), txBundleID)
		if err != nil {
			log.Errorf("trace transaction: failed to get receipt by previous tx hash, err: %v", err)
			rbErr := s.RollbackState(ctx, txBundleID)
			if rbErr != nil {
				log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
			}
			return &runtime.ExecutionResult{Err: err}
		}

		stateRoot = previousReceipt.PostState
	} else {
		previousBatch, err := s.GetBatchByHash(ctx, batch.Header.ParentHash, txBundleID)
		if err == ErrNotFound {
			previousBatch, err = s.GetLastBatch(ctx, true, txBundleID)
			if err != nil {
				log.Errorf("trace transaction: failed to get last batch, err: %v", err)
				rbErr := s.RollbackState(ctx, txBundleID)
				if rbErr != nil {
					log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
				}
				return &runtime.ExecutionResult{Err: err}
			}
		} else if err != nil {
			log.Errorf("trace transaction: failed to get batch by hash, err: %v", err)
			rbErr := s.RollbackState(ctx, txBundleID)
			if rbErr != nil {
				log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
			}
			return &runtime.ExecutionResult{Err: err}
		}

		stateRoot = previousBatch.Header.Root.Bytes()
	}

	sequencerAddress := batch.Header.Coinbase

	log.Debugf("replay root: %v", common.Bytes2Hex(stateRoot))

	bp, err := s.NewBatchProcessor(ctx, sequencerAddress, stateRoot, txBundleID)
	if err != nil {
		log.Errorf("trace transaction: failed to create a new batch processor, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return &runtime.ExecutionResult{Err: err}
	}

	// Activate EVM Instrumentation
	bp.Host.runtimes = []runtime.Runtime{}
	evm := evm.NewEVM()
	evm.EnableInstrumentation()
	bp.Host.setRuntime(evm)

	result := bp.processTransaction(ctx, tx, receipt.From, sequencerAddress)

	// Trace
	trace := instrumentation.Trace{}

	if bp.IsContractCreation(tx) {
		trace.Type = "create"
	} else if bp.IsSmartContractExecution(ctx, tx) {
		trace.Type = "call"
	}

	senderAddress, err := helper.GetSender(*tx)
	if err == nil {
		trace.Action = instrumentation.TraceAction{From: senderAddress.String(), To: tx.To().String(), Value: tx.Value().Uint64(), Gas: tx.Gas(), Input: tx.Data(), CallType: trace.Type}
	}

	if result.Err != nil {
		error := result.Err.Error()
		trace.Error = &error
	} else {
		trace.Result = &instrumentation.TraceResult{GasUsed: result.GasUsed, Output: result.ReturnValue}
	}

	result.Trace = trace

	// Rollback
	err = s.RollbackState(ctx, txBundleID)
	if err != nil {
		log.Errorf("trace transaction: failed to rollback transaction, err: %v", err)
		result.Err = err
		return result
	}

	return result
}

// ReplayBatchTransactions gets trace by rexecuting all the transactions of a specific batch
func (s *State) ReplayBatchTransactions(batchNumber uint64) ([]*runtime.ExecutionResult, error) {
	ctx := context.Background()

	txBundleID, err := s.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("trace transaction: failed to begin db transaction, err: %v", err)
		return nil, err
	}

	batch, err := s.GetBatchByNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		log.Errorf("trace transaction: failed to get batch by hash, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return nil, err
	}

	var stateRoot []byte

	previousBatch, err := s.GetBatchByHash(ctx, batch.Header.ParentHash, txBundleID)
	if err == ErrNotFound {
		previousBatch, err = s.GetLastBatch(ctx, true, txBundleID)
		if err != nil {
			log.Errorf("trace transaction: failed to get last batch, err: %v", err)
			rbErr := s.RollbackState(ctx, txBundleID)
			if rbErr != nil {
				log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
			}
			return nil, err
		}
	} else if err != nil {
		log.Errorf("trace transaction: failed to get batch by hash, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return nil, err
	}

	stateRoot = previousBatch.Header.Root.Bytes()

	sequencerAddress := batch.Header.Coinbase

	log.Debugf("replay root: %v", common.Bytes2Hex(stateRoot))

	bp, err := s.NewBatchProcessor(ctx, sequencerAddress, stateRoot, txBundleID)
	if err != nil {
		log.Errorf("trace transaction: failed to create a new batch processor, err: %v", err)
		rbErr := s.RollbackState(ctx, txBundleID)
		if rbErr != nil {
			log.Errorf("trace transaction: failed to rollback db transaction on error, err: %v, rollback err: %v", err, rbErr)
		}
		return nil, err
	}

	// Activate EVM Instrumentation
	bp.Host.runtimes = []runtime.Runtime{}
	evm := evm.NewEVM()
	evm.EnableInstrumentation()
	bp.Host.setRuntime(evm)

	results := make([]*runtime.ExecutionResult, 0, len(batch.Transactions))

	for _, tx := range batch.Transactions {
		from, err := helper.GetSender(*tx)
		if err != nil {
			results = append(results, &runtime.ExecutionResult{
				Err: err,
			})
		}

		result := bp.processTransaction(ctx, tx, from, sequencerAddress)

		// Trace
		trace := instrumentation.Trace{}

		if bp.IsContractCreation(tx) {
			trace.Type = "create"
		} else if bp.IsSmartContractExecution(ctx, tx) {
			trace.Type = "call"
		}

		senderAddress, err := helper.GetSender(*tx)
		if err == nil {
			trace.Action = instrumentation.TraceAction{From: senderAddress.String(), To: tx.To().String(), Value: tx.Value().Uint64(), Gas: tx.Gas(), Input: tx.Data(), CallType: trace.Type}
		}

		if result.Err != nil {
			error := result.Err.Error()
			trace.Error = &error
		} else {
			trace.Result = &instrumentation.TraceResult{GasUsed: result.GasUsed, Output: result.ReturnValue}
		}

		result.Trace = trace

		results = append(results, result)
	}

	err = s.RollbackState(ctx, txBundleID)
	if err != nil {
		log.Errorf("trace transaction: failed to rollback transaction, err: %v", err)
		return nil, err
	}

	return results, nil
}

// SetGenesis populates state with genesis information
func (s *State) SetGenesis(ctx context.Context, genesis Genesis, txBundleID string) error {
	// Generate Genesis Block
	block := &Block{
		BlockNumber: genesis.Block.NumberU64(),
		BlockHash:   genesis.Block.Hash(),
		ParentHash:  genesis.Block.ParentHash(),
		ReceivedAt:  genesis.Block.ReceivedAt,
	}

	// Add Block
	err := s.PostgresStorage.AddBlock(ctx, block, txBundleID)
	if err != nil {
		return err
	}

	var (
		root    common.Hash
		newRoot []byte
	)

	if genesis.Balances != nil {
		for address, balance := range genesis.Balances {
			newRoot, _, err = s.tree.SetBalance(ctx, address, balance, newRoot, txBundleID)
			if err != nil {
				return err
			}
		}
		root.SetBytes(newRoot)
	}

	if genesis.SmartContracts != nil {
		for address, sc := range genesis.SmartContracts {
			newRoot, _, err = s.tree.SetCode(ctx, address, sc, newRoot, txBundleID)
			if err != nil {
				return err
			}
		}
		root.SetBytes(newRoot)
	}

	if len(genesis.Storage) > 0 {
		for address, storage := range genesis.Storage {
			for key, value := range storage {
				newRoot, _, err = s.tree.SetStorageAt(ctx, address, key, value, newRoot, txBundleID)
				if err != nil {
					return err
				}
			}
		}
		root.SetBytes(newRoot)
	}

	if genesis.Nonces != nil {
		for address, nonce := range genesis.Nonces {
			newRoot, _, err = s.tree.SetNonce(ctx, address, nonce, newRoot, txBundleID)
			if err != nil {
				return err
			}
		}
		root.SetBytes(newRoot)
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
	bp, err := s.NewGenesisBatchProcessor(root[:], txBundleID)
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
func (s *State) GetNonce(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) (uint64, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return 0, err
	}

	n, err := s.tree.GetNonce(ctx, address, root, txBundleID)
	if errors.Is(err, tree.ErrNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return n.Uint64(), nil
}

// GetStorageAt from a given address
func (s *State) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, batchNumber uint64, txBundleID string) (*big.Int, error) {
	root, err := s.GetStateRootByBatchNumber(ctx, batchNumber, txBundleID)
	if err != nil {
		return nil, err
	}

	return s.tree.GetStorageAt(ctx, address, position, root, txBundleID)
}

// AddBlock adds a new block to the State Store.
func (s *State) AddBlock(ctx context.Context, block *Block, txBundleID string) error {
	return s.PostgresStorage.AddBlock(ctx, block, txBundleID)
}

// ConsolidateBatch consolidates a batch for the given DB tx bundle.
func (s *State) ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash, consolidatedAt time.Time, aggregator common.Address, txBundleID string) error {
	return s.PostgresStorage.ConsolidateBatch(ctx, batchNumber, consolidatedTxHash, consolidatedAt, aggregator, txBundleID)
}

// ResetDB resets the state to block for the given DB tx bundle.
func (s *State) ResetDB(ctx context.Context, block *Block, txBundleID string) error {
	return s.PostgresStorage.Reset(ctx, block, txBundleID)
}
