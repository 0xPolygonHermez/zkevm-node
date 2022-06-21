package statev2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/fakevm"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation"
	"github.com/hermeznetwork/hermez-core/state/runtime/instrumentation/tracers"
	"github.com/hermeznetwork/hermez-core/statev2/runtime/executor/pb"
	"github.com/holiman/uint256"
	"github.com/jackc/pgx/v4"
)

const (
	// TxTransferGas used for TXs that do not create a contract
	TxTransferGas uint64 = 21000
	// TxSmartContractCreationGas used for TXs that create a contract
	TxSmartContractCreationGas uint64 = 53000
	// Size of the memory in bytes reserved by the zkEVM
	zkEVMReservedMemorySize int = 128
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
	// ErrParsingExecutorTrace indicates an error occurred while parsing the executor trace
	ErrParsingExecutorTrace = fmt.Errorf("error while parsing executor trace")
)

// State is a implementation of the state
type State struct {
	cfg Config
	*PostgresStorage
	executorClient *pb.ExecutorServiceClient
}

// NewState creates a new State
func NewState(cfg Config, storage *PostgresStorage, executorClient *pb.ExecutorServiceClient) *State {
	return &State{
		cfg:             cfg,
		PostgresStorage: storage,
		executorClient:  executorClient,
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
func (s *State) CommitStateTransaction(ctx context.Context, dbTx pgx.Tx) error {
	err := dbTx.Commit(ctx)
	return err
}

// Rollback rollbacks a state transaction
func (s *State) RollbackStateTransaction(ctx context.Context, dbTx pgx.Tx) error {
	err := dbTx.Rollback(ctx)
	return err
}

// ResetDB resets the state to a block for the given DB tx
func (s *State) ResetDB(ctx context.Context, block *Block, dbTx pgx.Tx) error {
	return s.PostgresStorage.Reset(ctx, block, dbTx)
}

// ResetTrustedState resets the state db to a batch by its number
func (s *State) ResetTrustedState(ctx context.Context, batchNum uint64, dbTx pgx.Tx) error {
	// TODO: Implement
	// This method will need to update a field in the forced_batch table
	return nil
}

func (s *State) AddGlobalExitRoot(ctx context.Context, exitRoot *GlobalExitRoot, dbTx pgx.Tx) error {
	return s.PostgresStorage.AddGlobalExitRoot(ctx, exitRoot, dbTx)
}

func (s *State) GetLatestGlobalExitRoot(ctx context.Context, dbTx pgx.Tx) (*GlobalExitRoot, error) {
	return s.PostgresStorage.GetLatestGlobalExitRoot(ctx, dbTx)
}

func (s *State) AddForcedBatch(ctx context.Context, forcedBatch *ForcedBatch, dbTx pgx.Tx) error {
	return s.PostgresStorage.AddForcedBatch(ctx, forcedBatch, dbTx)
}

func (s *State) GetForcedBatch(ctx context.Context, dbTx pgx.Tx, forcedBatchNumber uint64) (*ForcedBatch, error) {
	return s.PostgresStorage.GetForcedBatch(ctx, dbTx, forcedBatchNumber)
}

// AddBlock adds a new block to the State Store.
func (s *State) AddBlock(ctx context.Context, block *Block, dbTx pgx.Tx) error {
	return s.PostgresStorage.AddBlock(ctx, block, dbTx)
}

// GetBalance from a given address
func (s *State) GetBalance(ctx context.Context, address common.Address, blockNumber uint64, dbdbTx pgx.Tx) (*big.Int, error) {
	// TODO: implement
	return nil, nil
}

// GetCode from a given address
func (s *State) GetCode(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) ([]byte, error) {
	// TODO: implement
	return nil, nil
}

// EstimateGas for a transaction
func (s *State) EstimateGas(transaction *types.Transaction, senderAddress common.Address) (uint64, error) {
	// TODO: implement
	return 0, nil
}

// GetNonce returns the nonce of the given account at the given block number
func (s *State) GetNonce(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	// TODO: implement
	return 0, nil
}

// GetStorageAt from a given address
func (s *State) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, batchNumber uint64, dbTx pgx.Tx) (*big.Int, error) {
	// TODO: implement
	return new(big.Int), nil
}

// StoreBatchHeader is used by the Trusted Sequencer to create a new batch
func (s *State) StoreBatchHeader(ctx context.Context, batch Batch) error {
	// TODO: implement
	return nil
}

// ProcessBatch is used by the Trusted Sequencer to add transactions to the batch
func (s *State) ProcessBatch(ctx context.Context, batchNumber uint64, txs []types.Transaction) (*ProcessBatchResponse, error) {
	// TODO: implement
	// check batchNumber is the latest in db
	return nil, nil
}

// StoreTransactions is used by the Trusted Sequencer to add processed transactions into the data base
func (s *State) StoreTransactions(batchNum uint64, processedTxs []*ProcessTransactionResponse) error {
	// TODO: implement
	return nil
}

// ProcessAndStoreWIPBatch is used by the Synchronizer to add a work-in-progress batch into the data base
func (s *State) ProcessAndStoreWIPBatch(ctx context.Context, batch Batch) error {
	// TODO: implement
	return nil
}

// ProcessAndStoreClosedBatch is used by the Synchronizer to a add closed batch into the data base
func (s *State) ProcessAndStoreClosedBatch(ctx context.Context, batch Batch) error {
	// TODO: implement
	return nil
}

// GetLastBatch gets latest batch (closed or not) on the data base
func (s *State) GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*Batch, error) {
	// TODO: implement
	return nil, nil
}

// GetBatchByNumber gets a batch from data base by its number
func (s *State) GetBatchByNumber(ctx context.Context, batchNumber uint64, tx pgx.Tx) (*Batch, error) {
	// TODO: implement
	return nil, nil
}

func (s *State) GetTrustedBatchByNumber(ctx context.Context, batchNumber uint64, tx pgx.Tx) (*Batch, error) {
	// TODO: implement
	return nil, nil
}

// GetEncodedTransactionsByBatchNumber gets the txs for a given batch in encoded form
func (s *State) GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, tx pgx.Tx) (encoded []string, err error) {
	// TODO: implement
	return nil, nil
}

// ProcessSequence process sequence of the txs
// TODO: implement function
func (s *State) ProcessBatchAndStoreLastTx(ctx context.Context, txs []types.Transaction) *runtime.ExecutionResult {
	return &runtime.ExecutionResult{}
}

// GetLastSendSequenceTime get time from last l1 interaction time
// TODO: implement function
func (s *State) GetLastSendSequenceTime(ctx context.Context) (time.Time, error) {
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
func (s *State) AddVerifiedBatch(ctx context.Context, verifiedBatch *VerifiedBatch, dbTx pgx.Tx) error {
	return s.PostgresStorage.AddVerifiedBatch(ctx, verifiedBatch, dbTx)
}

// GetVerifiedBatch get an L1 verifiedBatch
func (s *State) GetVerifiedBatch(ctx context.Context, dbTx pgx.Tx, batchNumber uint64) (*VerifiedBatch, error) {
	return s.PostgresStorage.GetVerifiedBatch(ctx, dbTx, batchNumber)
}

// DebugTransaction reexecutes a tx to generate its trace
func (s *State) DebugTransaction(ctx context.Context, transactionHash common.Hash, tracer string) (*runtime.ExecutionResult, error) {
	// TODO: Implement
	return new(runtime.ExecutionResult), nil
}

func (s *State) ParseTheTraceUsingTheTracer(env *fakevm.FakeEVM, trace instrumentation.ExecutorTrace, jsTracer tracers.Tracer) (json.RawMessage, error) {
	var previousDepth int
	var previousOpcode string
	var stateRoot []byte

	contextGas, ok := new(big.Int).SetString(trace.Context.Gas, encoding.Base10)
	if !ok {
		log.Debugf("error while parsing contextGas")
		return nil, ErrParsingExecutorTrace
	}
	value, ok := new(big.Int).SetString(trace.Context.Value, encoding.Base10)
	if !ok {
		log.Debugf("error while parsing value")
		return nil, ErrParsingExecutorTrace
	}

	jsTracer.CaptureTxStart(contextGas.Uint64())
	jsTracer.CaptureStart(env, common.HexToAddress(trace.Context.From), common.HexToAddress(trace.Context.To), trace.Context.Type == "CREATE", common.Hex2Bytes(strings.TrimLeft(trace.Context.Input, "0x")), contextGas.Uint64(), value)

	stack := fakevm.Newstack()
	memory := fakevm.NewMemory()

	bigStateRoot, ok := new(big.Int).SetString(trace.Context.OldStateRoot, 0)
	if !ok {
		log.Debugf("error while parsing context oldStateRoot")
		return nil, ErrParsingExecutorTrace
	}
	stateRoot = bigStateRoot.Bytes()
	env.StateDB.SetStateRoot(stateRoot)

	for i, step := range trace.Steps {
		gas, ok := new(big.Int).SetString(step.Gas, encoding.Base10)
		if !ok {
			log.Debugf("error while parsing step gas")
			return nil, ErrParsingExecutorTrace
		}

		gasCost, ok := new(big.Int).SetString(step.GasCost, encoding.Base10)
		if !ok {
			log.Debugf("error while parsing step gasCost")
			return nil, ErrParsingExecutorTrace
		}

		value, ok := new(big.Int).SetString(step.Contract.Value, encoding.Base10)
		if !ok {
			log.Debugf("error while parsing step value")
			return nil, ErrParsingExecutorTrace
		}

		op, ok := new(big.Int).SetString(step.Op, 0)
		if !ok {
			log.Debugf("error while parsing step op")
			return nil, ErrParsingExecutorTrace
		}

		scope := &fakevm.ScopeContext{
			Contract: vm.NewContract(fakevm.NewAccount(common.HexToAddress(step.Contract.Caller)), fakevm.NewAccount(common.HexToAddress(step.Contract.Address)), value, gas.Uint64()),
			Memory:   memory,
			Stack:    stack,
		}

		codeAddr := common.HexToAddress(step.Contract.Address)
		scope.Contract.CodeAddr = &codeAddr

		opcode := vm.OpCode(op.Uint64()).String()

		if previousOpcode == "CALL" && step.Pc != 0 {
			jsTracer.CaptureExit(common.Hex2Bytes(step.ReturnData), gasCost.Uint64(), fmt.Errorf(step.Error))
		}

		if opcode != "CALL" || trace.Steps[i+1].Pc == 0 {
			if step.Error != "" {
				err := fmt.Errorf(step.Error)
				jsTracer.CaptureFault(step.Pc, vm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, step.Depth, err)
			} else {
				jsTracer.CaptureState(step.Pc, vm.OpCode(op.Uint64()), gas.Uint64(), gasCost.Uint64(), scope, common.Hex2Bytes(strings.TrimLeft(step.ReturnData, "0x")), step.Depth, nil)
			}
		}

		if opcode == "CREATE" || opcode == "CREATE2" || opcode == "CALL" || opcode == "CALLCODE" || opcode == "DELEGATECALL" || opcode == "STATICCALL" || opcode == "SELFDESTRUCT" {
			jsTracer.CaptureEnter(vm.OpCode(op.Uint64()), common.HexToAddress(step.Contract.Caller), common.HexToAddress(step.Contract.Address), common.Hex2Bytes(strings.TrimLeft(step.Contract.Input, "0x")), gas.Uint64(), value)
			if step.OpCode == "SELFDESTRUCT" {
				jsTracer.CaptureExit(common.Hex2Bytes(step.ReturnData), gasCost.Uint64(), fmt.Errorf(step.Error))
			}
		}

		// Set Memory
		if len(step.Memory) > 0 {
			memory.Resize(uint64(fakevm.MemoryItemSize*len(step.Memory) + zkEVMReservedMemorySize))
			for offset, memoryContent := range step.Memory {
				memory.Set(uint64((offset*fakevm.MemoryItemSize)+zkEVMReservedMemorySize), uint64(fakevm.MemoryItemSize), common.Hex2Bytes(memoryContent))
			}
		} else {
			memory = fakevm.NewMemory()
		}

		// Set Stack
		stack = fakevm.Newstack()
		for _, stackContent := range step.Stack {
			valueBigInt, ok := new(big.Int).SetString(stackContent, 0)
			if !ok {
				log.Debugf("error while parsing stack valueBigInt")
				return nil, ErrParsingExecutorTrace
			}
			value, _ := uint256.FromBig(valueBigInt)
			stack.Push(value)
		}

		// Returning from a call or create
		if previousDepth > step.Depth {
			jsTracer.CaptureExit(common.Hex2Bytes(step.ReturnData), gasCost.Uint64(), fmt.Errorf(step.Error))
		}

		// Set StateRoot
		bigStateRoot, ok := new(big.Int).SetString(step.StateRoot, 0)
		if !ok {
			log.Debugf("error while parsing step stateRoot")
			return nil, ErrParsingExecutorTrace
		}

		stateRoot = bigStateRoot.Bytes()
		env.StateDB.SetStateRoot(stateRoot)
		previousDepth = step.Depth
		previousOpcode = step.OpCode
	}

	gasUsed, ok := new(big.Int).SetString(trace.Context.GasUsed, encoding.Base10)
	if !ok {
		log.Debugf("error while parsing gasUsed")
		return nil, ErrParsingExecutorTrace
	}

	jsTracer.CaptureTxEnd(gasUsed.Uint64())
	jsTracer.CaptureEnd(common.Hex2Bytes(trace.Context.Output), gasUsed.Uint64(), time.Duration(trace.Context.Time), nil)

	return jsTracer.GetResult()
}

func (s *State) GetLastBatch(ctx context.Context, tx pgx.Tx) (*Batch, error) {
	return s.PostgresStorage.GetLastBatch(ctx, tx)
}

func (s *State) GetBatchByNumber(ctx context.Context, batchNumber uint64, tx pgx.Tx) (*Batch, error) {
	return s.PostgresStorage.GetBatchByNumber(ctx, batchNumber, tx)
}

func (s *State) GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, tx pgx.Tx) ([]string, error) {
	return s.PostgresStorage.GetEncodedTransactionsByBatchNumber(ctx, batchNumber, tx)
}
