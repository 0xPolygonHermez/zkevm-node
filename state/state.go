package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/js"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/holiman/uint256"
	"github.com/jackc/pgx/v4"
)

const (
	// Size of the memory in bytes reserved by the zkEVM
	zkEVMReservedMemorySize int  = 128
	two                     uint = 2
	cTrue                        = 1
	cFalse                       = 0
)

var (
	once sync.Once
)

// CallerLabel is used to point which entity is the caller of a given function
type CallerLabel string

const (
	// SequencerCallerLabel is used when sequencer is calling the function
	SequencerCallerLabel CallerLabel = "sequencer"
	// SynchronizerCallerLabel is used when synchronizer is calling the function
	SynchronizerCallerLabel CallerLabel = "synchronizer"
	// DiscardCallerLabel is used we want to skip measuring the execution time
	DiscardCallerLabel CallerLabel = "discard"
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
	executorClient pb.ExecutorServiceClient
	tree           *merkletree.StateTree

	lastL2BlockSeen         types.Block
	newL2BlockEvents        chan NewL2BlockEvent
	newL2BlockEventHandlers []NewL2BlockEventHandler
}

// NewState creates a new State
func NewState(cfg Config, storage *PostgresStorage, executorClient pb.ExecutorServiceClient, stateTree *merkletree.StateTree) *State {
	once.Do(func() {
		metrics.Register()
	})

	s := &State{
		cfg:                     cfg,
		PostgresStorage:         storage,
		executorClient:          executorClient,
		tree:                    stateTree,
		newL2BlockEvents:        make(chan NewL2BlockEvent),
		newL2BlockEventHandlers: []NewL2BlockEventHandler{},
	}

	return s
}

// PrepareWebSocket allows the RPC to prepare ws
func (s *State) PrepareWebSocket() {
	lastL2Block, err := s.PostgresStorage.GetLastL2Block(context.Background(), nil)
	if errors.Is(err, ErrStateNotSynchronized) {
		lastL2Block = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(0)})
	} else if err != nil {
		log.Fatalf("failed to load the last l2 block: %v", err)
	}
	s.lastL2BlockSeen = *lastL2Block
	go s.monitorNewL2Blocks()
	go s.handleEvents()
}

// UpdateForkIDIntervals updates the forkID intervals
func (s *State) UpdateForkIDIntervals(intervals []ForkIDInterval) {
	log.Infof("Updating forkIDs. Setting %d forkIDs", len(intervals))
	s.cfg.ForkIDIntervals = intervals
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
func (s *State) GetBalance(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (*big.Int, error) {
	l2Block, err := s.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	return s.tree.GetBalance(ctx, address, l2Block.Root().Bytes())
}

// GetCode from a given address
func (s *State) GetCode(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) ([]byte, error) {
	l2Block, err := s.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	return s.tree.GetCode(ctx, address, l2Block.Root().Bytes())
}

// GetNonce returns the nonce of the given account at the given block number
func (s *State) GetNonce(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	l2Block, err := s.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		return 0, err
	}

	nonce, err := s.tree.GetNonce(ctx, address, l2Block.Root().Bytes())
	if err != nil {
		return 0, err
	}
	return nonce.Uint64(), nil
}

// GetLastStateRoot returns the latest state root
func (s *State) GetLastStateRoot(ctx context.Context, dbTx pgx.Tx) (common.Hash, error) {
	lastBlockHeader, err := s.GetLastL2BlockHeader(ctx, dbTx)
	if err != nil {
		return common.Hash{}, err
	}
	return lastBlockHeader.Root, nil
}

// GetStorageAt from a given address
func (s *State) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, blockNumber uint64, dbTx pgx.Tx) (*big.Int, error) {
	l2Block, err := s.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	return s.tree.GetStorageAt(ctx, address, position, l2Block.Root().Bytes())
}

// EstimateGas for a transaction
func (s *State) EstimateGas(transaction *types.Transaction, senderAddress common.Address, l2BlockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	const ethTransferGas = 21000

	var lowEnd uint64
	var highEnd uint64

	ctx := context.Background()

	lastBatches, l2BlockStateRoot, err := s.PostgresStorage.GetLastNBatchesByL2BlockNumber(ctx, &l2BlockNumber, two, dbTx)
	if err != nil {
		return 0, err
	}

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	lowEnd, err = core.IntrinsicGas(transaction.Data(), transaction.AccessList(), s.isContractCreation(transaction), true, false, false)
	if err != nil {
		return 0, err
	}

	if lowEnd == ethTransferGas && transaction.To() != nil {
		code, err := s.tree.GetCode(ctx, *transaction.To(), l2BlockStateRoot.Bytes())
		if err != nil {
			log.Warnf("error while getting transaction.to() code %v", err)
		} else if len(code) == 0 {
			return lowEnd, nil
		}
	}

	if transaction.Gas() != 0 && transaction.Gas() > lowEnd {
		highEnd = transaction.Gas()
	} else {
		highEnd = s.cfg.MaxCumulativeGasUsed
	}

	var availableBalance *big.Int

	if senderAddress != ZeroAddress {
		senderBalance, err := s.tree.GetBalance(ctx, senderAddress, l2BlockStateRoot.Bytes())
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				senderBalance = big.NewInt(0)
			} else {
				return 0, err
			}
		}

		availableBalance = new(big.Int).Set(senderBalance)

		if transaction.Value() != nil {
			if transaction.Value().Cmp(availableBalance) > 0 {
				return 0, ErrInsufficientFunds
			}

			availableBalance.Sub(availableBalance, transaction.Value())
		}
	}

	if transaction.GasPrice().BitLen() != 0 && // Gas price has been set
		availableBalance != nil && // Available balance is found
		availableBalance.Cmp(big.NewInt(0)) > 0 { // Available balance > 0
		gasAllowance := new(big.Int).Div(availableBalance, transaction.GasPrice())

		// Check the gas allowance for this account, make sure high end is capped to it
		if gasAllowance.IsUint64() && highEnd > gasAllowance.Uint64() {
			log.Debugf("Gas estimation high-end capped by allowance [%d]", gasAllowance.Uint64())
			highEnd = gasAllowance.Uint64()
		}
	}

	// Run the transaction with the specified gas value.
	// Returns a status indicating if the transaction failed, if it was reverted and the accompanying error
	testTransaction := func(gas uint64, shouldOmitErr bool) (bool, bool, uint64, error) {
		var gasUsed uint64
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    transaction.Nonce(),
			To:       transaction.To(),
			Value:    transaction.Value(),
			Gas:      gas,
			GasPrice: transaction.GasPrice(),
			Data:     transaction.Data(),
		})

		batchL2Data, err := EncodeUnsignedTransaction(*tx, s.cfg.ChainID, nil)
		if err != nil {
			log.Errorf("error encoding unsigned transaction ", err)
			return false, false, gasUsed, err
		}

		forkID := GetForkIDByBatchNumber(s.cfg.ForkIDIntervals, lastBatch.BatchNumber)
		// Create a batch to be sent to the executor
		processBatchRequest := &pb.ProcessBatchRequest{
			OldBatchNum:      lastBatch.BatchNumber,
			BatchL2Data:      batchL2Data,
			From:             senderAddress.String(),
			OldStateRoot:     l2BlockStateRoot.Bytes(),
			GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
			OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
			EthTimestamp:     uint64(lastBatch.Timestamp.Unix()),
			Coinbase:         lastBatch.Coinbase.String(),
			UpdateMerkleTree: cFalse,
			ChainId:          s.cfg.ChainID,
			ForkId:           forkID,
		}

		log.Debugf("EstimateGas[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
		// log.Debugf("EstimateGas[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
		log.Debugf("EstimateGas[processBatchRequest.From]: %v", processBatchRequest.From)
		log.Debugf("EstimateGas[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
		log.Debugf("EstimateGas[processBatchRequest.globalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
		log.Debugf("EstimateGas[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
		log.Debugf("EstimateGas[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
		log.Debugf("EstimateGas[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
		log.Debugf("EstimateGas[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
		log.Debugf("EstimateGas[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
		log.Debugf("EstimateGas[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)

		txExecutionOnExecutorTime := time.Now()
		processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
		log.Debugf("executor time: %vms", time.Since(txExecutionOnExecutorTime).Milliseconds())
		if err != nil {
			log.Errorf("error estimating gas: %v", err)
			return false, false, gasUsed, err
		}
		gasUsed = processBatchResponse.Responses[0].GasUsed
		if processBatchResponse.Error != executor.EXECUTOR_ERROR_NO_ERROR {
			err = executor.ExecutorErr(processBatchResponse.Error)
			s.LogExecutorError(processBatchResponse.Error, processBatchRequest)
			return false, false, gasUsed, err
		}

		// Check if an out of gas error happened during EVM execution
		if processBatchResponse.Responses[0].Error != pb.RomError(executor.ROM_ERROR_NO_ERROR) {
			err := executor.RomErr(processBatchResponse.Responses[0].Error)

			if (isGasEVMError(err) || isGasApplyError(err)) && shouldOmitErr {
				// Specifying the transaction failed, but not providing an error
				// is an indication that a valid error occurred due to low gas,
				// which will increase the lower bound for the search
				return true, false, gasUsed, nil
			}

			if isEVMRevertError(err) {
				// The EVM reverted during execution, attempt to extract the
				// error message and return it
				return true, true, gasUsed, constructErrorFromRevert(err, processBatchResponse.Responses[0].ReturnValue)
			}

			return true, false, gasUsed, err
		}

		return false, false, gasUsed, nil
	}

	txExecutions := []time.Duration{}
	var totalExecutionTime time.Duration

	// Check if the highEnd is a good value to make the transaction pass
	failed, reverted, gasUsed, err := testTransaction(highEnd, false)
	log.Debugf("Estimate gas. Trying to execute TX with %v gas", highEnd)
	if failed {
		if reverted {
			return 0, err
		}

		// The transaction shouldn't fail, for whatever reason, at highEnd
		return 0, fmt.Errorf(
			"unable to apply transaction even for the highest gas limit %d: %w",
			highEnd,
			err,
		)
	}

	if lowEnd < gasUsed {
		lowEnd = gasUsed
	}

	// Start the binary search for the lowest possible gas price
	for (lowEnd < highEnd) && (highEnd-lowEnd) > 4096 {
		txExecutionStart := time.Now()
		mid := (lowEnd + highEnd) / uint64(two)

		log.Debugf("Estimate gas. Trying to execute TX with %v gas", mid)

		failed, reverted, _, testErr := testTransaction(mid, true)
		executionTime := time.Since(txExecutionStart)
		totalExecutionTime += executionTime
		txExecutions = append(txExecutions, executionTime)
		if testErr != nil && !reverted {
			// Reverts are ignored in the binary search, but are checked later on
			// during the execution for the optimal gas limit found
			return 0, testErr
		}

		if failed {
			// If the transaction failed => increase the gas
			lowEnd = mid + 1
		} else {
			// If the transaction didn't fail => make this ok value the high end
			highEnd = mid
		}
	}

	executions := int64(len(txExecutions))
	if executions > 0 {
		log.Debugf("EstimateGas executed the TX %v times", executions)
		averageExecutionTime := totalExecutionTime.Milliseconds() / executions
		log.Debugf("EstimateGas tx execution average time is %v milliseconds", averageExecutionTime)
	} else {
		log.Error("Estimate gas. Tx not executed")
	}
	return highEnd, nil
}

// Checks if executor level valid gas errors occurred
func isGasApplyError(err error) bool {
	return errors.Is(err, ErrNotEnoughIntrinsicGas)
}

// Checks if EVM level valid gas errors occurred
func isGasEVMError(err error) bool {
	return errors.Is(err, runtime.ErrOutOfGas)
}

// Checks if the EVM reverted during execution
func isEVMRevertError(err error) bool {
	return errors.Is(err, runtime.ErrExecutionReverted)
}

// OpenBatch adds a new batch into the state, with the necessary data to start processing transactions within it.
// It's meant to be used by sequencers, since they don't necessarily know what transactions are going to be added
// in this batch yet. In other words it's the creation of a WIP batch.
// Note that this will add a batch with batch number N + 1, where N it's the greatest batch number on the state.
func (s *State) OpenBatch(ctx context.Context, processingContext ProcessingContext, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}
	// Check if the batch that is being opened has batch num + 1 compared to the latest batch
	lastBatchNum, err := s.PostgresStorage.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		return err
	}
	if lastBatchNum+1 != processingContext.BatchNumber {
		return fmt.Errorf("%w number %d, should be %d", ErrUnexpectedBatch, processingContext.BatchNumber, lastBatchNum+1)
	}
	// Check if last batch is closed
	isLastBatchClosed, err := s.PostgresStorage.IsBatchClosed(ctx, lastBatchNum, dbTx)
	if err != nil {
		return err
	}
	if !isLastBatchClosed {
		return ErrLastBatchShouldBeClosed
	}
	// Check that timestamp is equal or greater compared to previous batch
	prevTimestamp, err := s.GetLastBatchTime(ctx, dbTx)
	if err != nil {
		return err
	}
	if prevTimestamp.Unix() > processingContext.Timestamp.Unix() {
		return ErrTimestampGE
	}
	return s.PostgresStorage.openBatch(ctx, processingContext, dbTx)
}

// ProcessSequencerBatch is used by the sequencers to process transactions into an open batch
func (s *State) ProcessSequencerBatch(ctx context.Context, batchNumber uint64, batchL2Data []byte, caller CallerLabel, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	log.Debugf("*******************************************")
	log.Debugf("ProcessSequencerBatch start")

	processBatchResponse, err := s.processBatch(ctx, batchNumber, batchL2Data, caller, dbTx)
	if err != nil {
		return nil, err
	}

	txs, _, err := DecodeTxs(batchL2Data)
	if err != nil && !errors.Is(err, InvalidData) {
		return nil, err
	}
	result, err := s.convertToProcessBatchResponse(txs, processBatchResponse)
	if err != nil {
		return nil, err
	}
	log.Debugf("ProcessSequencerBatch end")
	log.Debugf("*******************************************")
	return result, nil
}

// ProcessBatch processes a batch
func (s *State) ProcessBatch(ctx context.Context, request ProcessRequest, updateMerkleTree bool) (*ProcessBatchResponse, error) {
	log.Debugf("*******************************************")
	log.Debugf("ProcessBatch start")

	updateMT := uint32(cFalse)
	if updateMerkleTree {
		updateMT = cTrue
	}

	forkID := GetForkIDByBatchNumber(s.cfg.ForkIDIntervals, request.BatchNumber)
	// Create Batch
	processBatchRequest := &pb.ProcessBatchRequest{
		OldBatchNum:      request.BatchNumber - 1,
		Coinbase:         request.Coinbase.String(),
		BatchL2Data:      request.Transactions,
		OldStateRoot:     request.OldStateRoot.Bytes(),
		GlobalExitRoot:   request.GlobalExitRoot.Bytes(),
		OldAccInputHash:  request.OldAccInputHash.Bytes(),
		EthTimestamp:     request.Timestamp,
		UpdateMerkleTree: updateMT,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
	}
	res, err := s.sendBatchRequestToExecutor(ctx, processBatchRequest, request.Caller)
	if err != nil {
		return nil, err
	}

	txs, _, err := DecodeTxs(request.Transactions)
	if err != nil && !errors.Is(err, InvalidData) {
		return nil, err
	}

	var result *ProcessBatchResponse
	result, err = s.convertToProcessBatchResponse(txs, res)
	if err != nil {
		return nil, err
	}

	log.Debugf("ProcessBatch end")
	log.Debugf("*******************************************")

	return result, nil
}

// ExecuteBatch is used by the synchronizer to reprocess batches to compare generated state root vs stored one
// It is also used by the sequencer in order to calculate used zkCounter of a WIPBatch
func (s *State) ExecuteBatch(ctx context.Context, batch Batch, updateMerkleTree bool, dbTx pgx.Tx) (*pb.ProcessBatchResponse, error) {
	if dbTx == nil {
		return nil, ErrDBTxNil
	}

	// Get previous batch to get state root and local exit root
	previousBatch, err := s.PostgresStorage.GetBatchByNumber(ctx, batch.BatchNumber-1, dbTx)
	if err != nil {
		return nil, err
	}

	forkId := s.GetForkIdByBatchNumber(batch.BatchNumber)

	updateMT := uint32(cFalse)
	if updateMerkleTree {
		updateMT = cTrue
	}

	// Create Batch
	processBatchRequest := &pb.ProcessBatchRequest{
		OldBatchNum:     batch.BatchNumber - 1,
		Coinbase:        batch.Coinbase.String(),
		BatchL2Data:     batch.BatchL2Data,
		OldStateRoot:    previousBatch.StateRoot.Bytes(),
		GlobalExitRoot:  batch.GlobalExitRoot.Bytes(),
		OldAccInputHash: previousBatch.AccInputHash.Bytes(),
		EthTimestamp:    uint64(batch.Timestamp.Unix()),
		// Changed for new sequencer strategy
		UpdateMerkleTree: updateMT,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkId,
	}

	// Send Batch to the Executor
	log.Debugf("ExecuteBatch[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
	log.Debugf("ExecuteBatch[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
	log.Debugf("ExecuteBatch[processBatchRequest.From]: %v", processBatchRequest.From)
	log.Debugf("ExecuteBatch[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
	log.Debugf("ExecuteBatch[processBatchRequest.GlobalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
	log.Debugf("ExecuteBatch[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
	log.Debugf("ExecuteBatch[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
	log.Debugf("ExecuteBatch[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
	log.Debugf("ExecuteBatch[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
	log.Debugf("ExecuteBatch[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
	log.Debugf("ExecuteBatch[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)

	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		log.Error("error executing batch: ", err)
		return nil, err
	} else if processBatchResponse != nil && processBatchResponse.Error != executor.EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.LogExecutorError(processBatchResponse.Error, processBatchRequest)
	}

	return processBatchResponse, err
}

func (s *State) processBatch(
	ctx context.Context,
	batchNumber uint64,
	batchL2Data []byte,
	caller CallerLabel,
	dbTx pgx.Tx,
) (*pb.ProcessBatchResponse, error) {
	if dbTx == nil {
		return nil, ErrDBTxNil
	}
	lastBatches, err := s.PostgresStorage.GetLastNBatches(ctx, two, dbTx)
	if err != nil {
		return nil, err
	}

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	isBatchClosed, err := s.PostgresStorage.IsBatchClosed(ctx, batchNumber, dbTx)
	if err != nil {
		return nil, err
	}
	if isBatchClosed {
		return nil, ErrBatchAlreadyClosed
	}

	// Check provided batch number is the latest in db
	if lastBatch.BatchNumber != batchNumber {
		return nil, ErrInvalidBatchNumber
	}
	forkID := GetForkIDByBatchNumber(s.cfg.ForkIDIntervals, lastBatch.BatchNumber)
	// Create Batch
	processBatchRequest := &pb.ProcessBatchRequest{
		OldBatchNum:      lastBatch.BatchNumber - 1,
		Coinbase:         lastBatch.Coinbase.String(),
		BatchL2Data:      batchL2Data,
		OldStateRoot:     previousBatch.StateRoot.Bytes(),
		GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
		OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
		EthTimestamp:     uint64(lastBatch.Timestamp.Unix()),
		UpdateMerkleTree: cTrue,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
	}

	res, err := s.sendBatchRequestToExecutor(ctx, processBatchRequest, caller)

	return res, err
}

func (s *State) sendBatchRequestToExecutor(ctx context.Context, processBatchRequest *pb.ProcessBatchRequest, caller CallerLabel) (*pb.ProcessBatchResponse, error) {
	// Send Batch to the Executor
	if caller != DiscardCallerLabel {
		log.Debugf("processBatch[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
		log.Debugf("processBatch[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
		log.Debugf("processBatch[processBatchRequest.From]: %v", processBatchRequest.From)
		log.Debugf("processBatch[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
		log.Debugf("processBatch[processBatchRequest.GlobalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
		log.Debugf("processBatch[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
		log.Debugf("processBatch[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
		log.Debugf("processBatch[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
		log.Debugf("processBatch[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
		log.Debugf("processBatch[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
		log.Debugf("processBatch[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
	}
	now := time.Now()
	res, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		log.Errorf("Error s.executorClient.ProcessBatch: %v", err)
		log.Errorf("Error s.executorClient.ProcessBatch: %s", err.Error())
		log.Errorf("Error s.executorClient.ProcessBatch response: %v", res)
	} else if res.Error != executor.EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(res.Error)
		s.LogExecutorError(res.Error, processBatchRequest)
	}
	elapsed := time.Since(now)
	if caller != DiscardCallerLabel {
		metrics.ExecutorProcessingTime(string(caller), elapsed)
	}
	log.Infof("Batch: %d took %v to be processed by the executor ", processBatchRequest.OldBatchNum+1, elapsed)

	return res, err
}

// StoreTransactions is used by the sequencer to add processed transactions into
// an open batch. If the batch already has txs, the processedTxs must be a super
// set of the existing ones, preserving order.
func (s *State) StoreTransactions(ctx context.Context, batchNumber uint64, processedTxs []*ProcessTransactionResponse, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	// check existing txs vs parameter txs
	existingTxs, err := s.GetTxsHashesByBatchNumber(ctx, batchNumber, dbTx)
	if err != nil {
		return err
	}
	if err := CheckSupersetBatchTransactions(existingTxs, processedTxs); err != nil {
		return err
	}

	// Check if last batch is closed. Note that it's assumed that only the latest batch can be open
	isBatchClosed, err := s.PostgresStorage.IsBatchClosed(ctx, batchNumber, dbTx)
	if err != nil {
		return err
	}
	if isBatchClosed {
		return ErrBatchAlreadyClosed
	}

	processingContext, err := s.GetProcessingContext(ctx, batchNumber, dbTx)
	if err != nil {
		return err
	}

	firstTxToInsert := len(existingTxs)

	for i := firstTxToInsert; i < len(processedTxs); i++ {
		processedTx := processedTxs[i]
		// if the transaction has an intrinsic invalid tx error it means
		// the transaction has not changed the state, so we don't store it
		// and just move to the next
		if executor.IsIntrinsicError(executor.RomErrorCode(processedTx.RomError)) {
			continue
		}

		lastL2Block, err := s.GetLastL2Block(ctx, dbTx)
		if err != nil {
			return err
		}

		header := &types.Header{
			Number:     new(big.Int).SetUint64(lastL2Block.Number().Uint64() + 1),
			ParentHash: lastL2Block.Hash(),
			Coinbase:   processingContext.Coinbase,
			Root:       processedTx.StateRoot,
			GasUsed:    processedTx.GasUsed,
			GasLimit:   s.cfg.MaxCumulativeGasUsed,
			Time:       uint64(processingContext.Timestamp.Unix()),
		}
		transactions := []*types.Transaction{&processedTx.Tx}

		receipt := generateReceipt(header.Number, processedTx)
		receipts := []*types.Receipt{receipt}

		// Create block to be able to calculate its hash
		block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
		block.ReceivedAt = processingContext.Timestamp

		receipt.BlockHash = block.Hash()

		// Store L2 block and its transaction
		if err := s.AddL2Block(ctx, batchNumber, block, receipts, dbTx); err != nil {
			return err
		}
	}
	return nil
}

func (s *State) isBatchClosable(ctx context.Context, receipt ProcessingReceipt, dbTx pgx.Tx) error {
	// Check if the batch that is being closed is the last batch
	lastBatchNum, err := s.PostgresStorage.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		return err
	}
	if lastBatchNum != receipt.BatchNumber {
		return fmt.Errorf("%w number %d, should be %d", ErrUnexpectedBatch, receipt.BatchNumber, lastBatchNum)
	}
	// Check if last batch is closed
	isLastBatchClosed, err := s.PostgresStorage.IsBatchClosed(ctx, lastBatchNum, dbTx)
	if err != nil {
		return err
	}
	if isLastBatchClosed {
		return ErrBatchAlreadyClosed
	}

	return nil
}

// CloseBatch is used by sequencer to close the current batch
func (s *State) CloseBatch(ctx context.Context, receipt ProcessingReceipt, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	err := s.isBatchClosable(ctx, receipt, dbTx)
	if err != nil {
		return err
	}

	return s.PostgresStorage.closeBatch(ctx, receipt, dbTx)
}

// ProcessAndStoreClosedBatch is used by the Synchronizer to add a closed batch into the data base
func (s *State) ProcessAndStoreClosedBatch(
	ctx context.Context,
	processingCtx ProcessingContext,
	encodedTxs []byte,
	dbTx pgx.Tx,
	caller CallerLabel,
) (common.Hash, error) {
	// Decode transactions
	decodedTransactions, _, err := DecodeTxs(encodedTxs)
	if err != nil && !errors.Is(err, InvalidData) {
		log.Debugf("error decoding transactions: %v", err)
		return common.Hash{}, err
	}

	// Open the batch and process the txs
	if dbTx == nil {
		return common.Hash{}, ErrDBTxNil
	}
	if err := s.OpenBatch(ctx, processingCtx, dbTx); err != nil {
		return common.Hash{}, err
	}
	processed, err := s.processBatch(ctx, processingCtx.BatchNumber, encodedTxs, caller, dbTx)
	if err != nil {
		return common.Hash{}, err
	}

	// Sanity check
	if len(decodedTransactions) != len(processed.Responses) {
		log.Errorf("number of decoded (%d) and processed (%d) transactions do not match", len(decodedTransactions), len(processed.Responses))
	}

	// Filter unprocessed txs and decode txs to store metadata
	// note that if the batch is not well encoded it will result in an empty batch (with no txs)
	for i := 0; i < len(processed.Responses); i++ {
		if !isProcessed(processed.Responses[i].Error) {
			if executor.IsROMOutOfCountersError(processed.Responses[i].Error) {
				processed.Responses = []*pb.ProcessTransactionResponse{}
				break
			}

			// Remove unprocessed tx
			if i == len(processed.Responses)-1 {
				processed.Responses = processed.Responses[:i]
				decodedTransactions = decodedTransactions[:i]
			} else {
				processed.Responses = append(processed.Responses[:i], processed.Responses[i+1:]...)
				decodedTransactions = append(decodedTransactions[:i], decodedTransactions[i+1:]...)
			}
			i--
		}
	}

	processedBatch, err := s.convertToProcessBatchResponse(decodedTransactions, processed)
	if err != nil {
		return common.Hash{}, err
	}

	if len(processedBatch.Responses) > 0 {
		// Store processed txs into the batch
		err = s.StoreTransactions(ctx, processingCtx.BatchNumber, processedBatch.Responses, dbTx)
		if err != nil {
			return common.Hash{}, err
		}
	}

	// Close batch
	return common.BytesToHash(processed.NewStateRoot), s.closeBatch(ctx, ProcessingReceipt{
		BatchNumber:   processingCtx.BatchNumber,
		StateRoot:     processedBatch.NewStateRoot,
		LocalExitRoot: processedBatch.NewLocalExitRoot,
		AccInputHash:  processedBatch.NewAccInputHash,
		BatchL2Data:   encodedTxs,
	}, dbTx)
}

// GetLastBatch gets latest batch (closed or not) on the data base
func (s *State) GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*Batch, error) {
	batches, err := s.PostgresStorage.GetLastNBatches(ctx, 1, dbTx)
	if err != nil {
		return nil, err
	}
	if len(batches) == 0 {
		return nil, ErrNotFound
	}
	return batches[0], nil
}

// DebugTransaction re-executes a tx to generate its trace
func (s *State) DebugTransaction(ctx context.Context, transactionHash common.Hash, traceConfig TraceConfig, dbTx pgx.Tx) (*runtime.ExecutionResult, error) {
	result := new(runtime.ExecutionResult)

	// Get the transaction
	tx, err := s.GetTransactionByHash(ctx, transactionHash, dbTx)
	if err != nil {
		return nil, err
	}

	// Get batch including the transaction
	batch, err := s.GetBatchByTxHash(ctx, transactionHash, dbTx)
	if err != nil {
		return nil, err
	}

	// The previous batch to get OldStateRoot and globalExitRoot
	pBatch, err := s.GetBatchByNumber(ctx, batch.BatchNumber-1, dbTx)
	if err != nil {
		return nil, err
	}

	batchL2Data := batch.BatchL2Data
	if batchL2Data == nil {
		txs, err := s.GetTransactionsByBatchNumber(ctx, batch.BatchNumber, dbTx)
		if err != nil {
			return nil, err
		}

		for _, tx := range txs {
			log.Debugf(tx.Hash().String())
		}

		batchL2Data, err = EncodeTransactions(txs)
		if err != nil {
			return nil, err
		}
	}

	forkId := s.GetForkIdByBatchNumber(batch.BatchNumber)

	// Create Batch
	traceConfigRequest := &pb.TraceConfig{
		TxHashToGenerateCallTrace:    transactionHash.Bytes(),
		TxHashToGenerateExecuteTrace: transactionHash.Bytes(),
	}

	if traceConfig.DisableStorage {
		traceConfigRequest.DisableStorage = cTrue
	}
	if traceConfig.DisableStack {
		traceConfigRequest.DisableStack = cTrue
	}
	if traceConfig.EnableMemory {
		traceConfigRequest.EnableMemory = cTrue
	}
	if traceConfig.EnableReturnData {
		traceConfigRequest.EnableReturnData = cTrue
	}

	processBatchRequest := &pb.ProcessBatchRequest{
		OldBatchNum:      batch.BatchNumber - 1,
		BatchL2Data:      batchL2Data,
		OldStateRoot:     pBatch.StateRoot.Bytes(),
		GlobalExitRoot:   batch.GlobalExitRoot.Bytes(),
		OldAccInputHash:  pBatch.AccInputHash.Bytes(),
		EthTimestamp:     uint64(batch.Timestamp.Unix()),
		Coinbase:         batch.Coinbase.String(),
		UpdateMerkleTree: cFalse,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkId,
		TraceConfig:      traceConfigRequest,
	}

	// Send Batch to the Executor
	startTime := time.Now()
	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		return nil, err
	} else if processBatchResponse.Error != executor.EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.LogExecutorError(processBatchResponse.Error, processBatchRequest)
		return nil, err
	}
	endTime := time.Now()

	// //save process batch response file
	// b, err := json.Marshal(processBatchResponse)
	// if err != nil {
	// 	return nil, err
	// }
	// filePath := "./processBatchResponse.json"
	// err = os.WriteFile(filePath, b, 0644)
	// if err != nil {
	// 	return nil, err
	// }

	for _, response := range processBatchResponse.Responses {
		log.Debugf(string(response.TxHash))
	}

	txs, _, err := DecodeTxs(batchL2Data)
	if err != nil && !errors.Is(err, InvalidData) {
		return nil, err
	}

	for _, tx := range txs {
		log.Debugf(tx.Hash().String())
	}

	convertedResponse, err := s.convertToProcessBatchResponse(txs, processBatchResponse)
	if err != nil {
		return nil, err
	}

	var response *ProcessTransactionResponse

	// Get the response for the tx
	for _, response = range convertedResponse.Responses {
		log.Debugf(response.TxHash.String())
		if response.TxHash == transactionHash {
			break
		}
	}

	// Sanity check
	if response.TxHash != transactionHash {
		return nil, fmt.Errorf("tx hash not found in executor response")
	}

	result.CreateAddress = response.CreateAddress
	result.GasLeft = response.GasLeft
	result.GasUsed = response.GasUsed
	result.ReturnValue = response.ReturnValue
	result.StateRoot = response.StateRoot.Bytes()
	result.StructLogs = response.ExecutionTrace

	if traceConfig.Tracer == nil || *traceConfig.Tracer == "" {
		return result, nil
	}

	// Parse the executor-like trace using the FakeEVM
	jsTracer, err := js.NewJsTracer(*traceConfig.Tracer, new(tracers.Context))
	if err != nil {
		log.Errorf("debug transaction: failed to create jsTracer, err: %v", err)
		return nil, fmt.Errorf("failed to create jsTracer, err: %v", err)
	}

	context := instrumentation.Context{}

	// Fill trace context
	if tx.To() == nil {
		context.Type = "CREATE"
		context.To = result.CreateAddress.Hex()
	} else {
		context.Type = "CALL"
		context.To = tx.To().Hex()
	}

	senderAddress, err := GetSender(*tx)
	if err != nil {
		return nil, err
	}

	context.From = senderAddress.String()
	context.Input = "0x" + hex.EncodeToString(tx.Data())
	context.Gas = strconv.FormatUint(tx.Gas(), encoding.Base10)
	context.Value = tx.Value().String()
	context.Output = "0x" + hex.EncodeToString(result.ReturnValue)
	context.GasPrice = tx.GasPrice().String()
	context.OldStateRoot = batch.StateRoot.String()
	context.Time = uint64(endTime.Sub(startTime))
	context.GasUsed = strconv.FormatUint(result.GasUsed, encoding.Base10)

	result.ExecutorTrace.Context = context

	gasPrice, ok := new(big.Int).SetString(context.GasPrice, encoding.Base10)
	if !ok {
		log.Errorf("debug transaction: failed to parse gasPrice")
		return nil, fmt.Errorf("failed to parse gasPrice")
	}

	env := fakevm.NewFakeEVM(vm.BlockContext{BlockNumber: big.NewInt(1)}, vm.TxContext{GasPrice: gasPrice}, params.TestChainConfig, fakevm.Config{Debug: true, Tracer: jsTracer})
	fakeDB := &FakeDB{State: s, stateRoot: batch.StateRoot.Bytes()}
	env.SetStateDB(fakeDB)

	traceResult, err := s.ParseTheTraceUsingTheTracer(env, result.ExecutorTrace, jsTracer)
	if err != nil {
		log.Errorf("debug transaction: failed parse the trace using the tracer: %v", err)
		return nil, fmt.Errorf("failed parse the trace using the tracer: %v", err)
	}

	result.ExecutorTraceResult = traceResult

	return result, nil
}

// ParseTheTraceUsingTheTracer parses the given trace with the given tracer.
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

// PreProcessTransaction processes the transaction in order to calculate its zkCounters before adding it to the pool
func (s *State) PreProcessTransaction(ctx context.Context, tx *types.Transaction, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	sender, err := GetSender(*tx)
	if err != nil {
		return nil, err
	}

	lastL2BlockNumber, err := s.GetLastL2BlockNumber(ctx, dbTx)
	if err != nil {
		return nil, err
	}

	response, err := s.internalProcessUnsignedTransaction(ctx, tx, sender, lastL2BlockNumber, false, dbTx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ProcessUnsignedTransaction processes the given unsigned transaction.
func (s *State) ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, l2BlockNumber uint64, noZKEVMCounters bool, dbTx pgx.Tx) (*runtime.ExecutionResult, error) {
	result := new(runtime.ExecutionResult)
	response, err := s.internalProcessUnsignedTransaction(ctx, tx, senderAddress, l2BlockNumber, noZKEVMCounters, dbTx)
	if err != nil {
		return nil, err
	}

	r := response.Responses[0]
	result.ReturnValue = r.ReturnValue
	result.GasLeft = r.GasLeft
	result.GasUsed = r.GasUsed
	result.CreateAddress = r.CreateAddress
	result.StateRoot = r.StateRoot.Bytes()

	if errors.Is(r.RomError, runtime.ErrExecutionReverted) {
		result.Err = constructErrorFromRevert(r.RomError, r.ReturnValue)
	} else {
		result.Err = r.RomError
	}

	return result, nil
}

// ProcessUnsignedTransaction processes the given unsigned transaction.
func (s *State) internalProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, l2BlockNumber uint64, noZKEVMCounters bool, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	lastBatches, _, err := s.PostgresStorage.GetLastNBatchesByL2BlockNumber(ctx, &l2BlockNumber, two, dbTx)
	if err != nil {
		return nil, err
	}

	l2Block, err := s.GetL2BlockByNumber(ctx, l2BlockNumber, dbTx)
	if err != nil {
		log.Errorf("error getting l2 block", err)
		return nil, err
	}

	nonce, err := s.tree.GetNonce(ctx, senderAddress, l2Block.Root().Bytes())
	if err != nil {
		return nil, err
	}
	forcedNonce := nonce.Uint64()

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	batchL2Data, err := EncodeUnsignedTransaction(*tx, s.cfg.ChainID, &forcedNonce)
	if err != nil {
		log.Errorf("error encoding unsigned transaction ", err)
		return nil, err
	}

	forkID := GetForkIDByBatchNumber(s.cfg.ForkIDIntervals, lastBatch.BatchNumber)
	// Create Batch
	processBatchRequest := &pb.ProcessBatchRequest{
		OldBatchNum:      lastBatch.BatchNumber,
		BatchL2Data:      batchL2Data,
		From:             senderAddress.String(),
		OldStateRoot:     l2Block.Root().Bytes(),
		GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
		OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
		EthTimestamp:     uint64(lastBatch.Timestamp.Unix()),
		Coinbase:         lastBatch.Coinbase.String(),
		UpdateMerkleTree: cFalse,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
	}

	if noZKEVMCounters {
		processBatchRequest.NoCounters = cTrue
	}

	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.From]: %v", processBatchRequest.From)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.globalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)

	// Send Batch to the Executor
	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		// Log this error as an executor unspecified error
		s.LogExecutorError(pb.ExecutorError_EXECUTOR_ERROR_UNSPECIFIED, processBatchRequest)
		log.Errorf("error processing unsigned transaction ", err)
		return nil, err
	} else if processBatchResponse.Error != executor.EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.LogExecutorError(processBatchResponse.Error, processBatchRequest)
		return nil, err
	}

	response, err := s.convertToProcessBatchResponse([]types.Transaction{*tx}, processBatchResponse)
	if err != nil {
		return nil, err
	}

	if processBatchResponse.Responses[0].Error != pb.RomError(executor.ROM_ERROR_NO_ERROR) {
		err := executor.RomErr(processBatchResponse.Responses[0].Error)
		if !isEVMRevertError(err) {
			return response, err
		}
	}

	return response, nil
}

// GetTree returns State inner tree
func (s *State) GetTree() *merkletree.StateTree {
	return s.tree
}

// SetGenesis populates state with genesis information
func (s *State) SetGenesis(ctx context.Context, block Block, genesis Genesis, dbTx pgx.Tx) ([]byte, error) {
	var (
		root    common.Hash
		newRoot []byte
		err     error
	)
	if dbTx == nil {
		return newRoot, ErrDBTxNil
	}

	for _, action := range genesis.Actions {
		address := common.HexToAddress(action.Address)
		switch action.Type {
		case int(merkletree.LeafTypeBalance):
			balance, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return newRoot, err
			}
			newRoot, _, err = s.tree.SetBalance(ctx, address, balance, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeNonce):
			nonce, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return newRoot, err
			}
			newRoot, _, err = s.tree.SetNonce(ctx, address, nonce, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeCode):
			code, err := hex.DecodeHex(action.Bytecode)
			if err != nil {
				return newRoot, fmt.Errorf("could not decode SC bytecode for address %q: %v", address, err)
			}
			newRoot, _, err = s.tree.SetCode(ctx, address, code, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeStorage):
			// Parse position and value
			positionBI, err := encoding.DecodeBigIntHexOrDecimal(action.StoragePosition)
			if err != nil {
				return newRoot, err
			}
			valueBI, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return newRoot, err
			}
			// Store
			newRoot, _, err = s.tree.SetStorageAt(ctx, address, positionBI, valueBI, newRoot)
			if err != nil {
				return newRoot, err
			}
		case int(merkletree.LeafTypeSCLength):
			log.Debug("Skipped genesis action of type merkletree.LeafTypeSCLength, these actions will be handled as part of merkletree.LeafTypeCode actions")
		default:
			return newRoot, fmt.Errorf("unknown genesis action type %q", action.Type)
		}
	}

	root.SetBytes(newRoot)

	// flush state db
	err = s.tree.Flush(ctx)
	if err != nil {
		log.Errorf("error flushing state tree after genesis: %v", err)
		return newRoot, err
	}

	// store L1 block related to genesis batch
	err = s.AddBlock(ctx, &block, dbTx)
	if err != nil {
		return newRoot, err
	}

	// store genesis batch
	batch := Batch{
		BatchNumber:    0,
		Coinbase:       ZeroAddress,
		BatchL2Data:    nil,
		StateRoot:      root,
		LocalExitRoot:  ZeroHash,
		Timestamp:      block.ReceivedAt,
		Transactions:   []types.Transaction{},
		GlobalExitRoot: ZeroHash,
		ForcedBatchNum: nil,
	}

	err = s.storeGenesisBatch(ctx, batch, dbTx)
	if err != nil {
		return newRoot, err
	}

	// mark the genesis batch as virtualized
	virtualBatch := &VirtualBatch{
		BatchNumber: batch.BatchNumber,
		TxHash:      ZeroHash,
		Coinbase:    ZeroAddress,
		BlockNumber: block.BlockNumber,
	}
	err = s.AddVirtualBatch(ctx, virtualBatch, dbTx)
	if err != nil {
		return newRoot, err
	}

	// mark the genesis batch as verified/consolidated
	verifiedBatch := &VerifiedBatch{
		BatchNumber: batch.BatchNumber,
		TxHash:      ZeroHash,
		Aggregator:  ZeroAddress,
		BlockNumber: block.BlockNumber,
	}
	err = s.AddVerifiedBatch(ctx, verifiedBatch, dbTx)
	if err != nil {
		return newRoot, err
	}

	// store L2 genesis block
	header := &types.Header{
		Number:     big.NewInt(0),
		ParentHash: ZeroHash,
		Coinbase:   ZeroAddress,
		Root:       root,
		Time:       uint64(block.ReceivedAt.Unix()),
	}
	rootHex := root.Hex()
	log.Info("Genesis root ", rootHex)

	receipts := []*types.Receipt{}
	l2Block := types.NewBlock(header, []*types.Transaction{}, []*types.Header{}, receipts, &trie.StackTrie{})
	l2Block.ReceivedAt = block.ReceivedAt

	return newRoot, s.AddL2Block(ctx, batch.BatchNumber, l2Block, receipts, dbTx)
}

// CheckSupersetBatchTransactions verifies that processedTransactions is a
// superset of existingTxs and that the existing txs have the same order,
// returns a non-nil error if that is not the case.
func CheckSupersetBatchTransactions(existingTxHashes []common.Hash, processedTxs []*ProcessTransactionResponse) error {
	if len(existingTxHashes) > len(processedTxs) {
		return ErrExistingTxGreaterThanProcessedTx
	}
	for i, existingTxHash := range existingTxHashes {
		if existingTxHash != processedTxs[i].TxHash {
			return ErrOutOfOrderProcessedTx
		}
	}
	return nil
}

// isContractCreation checks if the tx is a contract creation
func (s *State) isContractCreation(tx *types.Transaction) bool {
	return tx.To() == nil && len(tx.Data()) > 0
}

// DetermineProcessedTransactions splits the given tx process responses
// returning a slice with only processed and a map unprocessed txs
// respectively.
func DetermineProcessedTransactions(responses []*ProcessTransactionResponse) (
	[]*ProcessTransactionResponse, []string, map[string]*ProcessTransactionResponse, []string) {
	processedTxResponses := []*ProcessTransactionResponse{}
	processedTxsHashes := []string{}
	unprocessedTxResponses := map[string]*ProcessTransactionResponse{}
	unprocessedTxsHashes := []string{}
	for _, response := range responses {
		if response.IsProcessed {
			processedTxResponses = append(processedTxResponses, response)
			processedTxsHashes = append(processedTxsHashes, response.TxHash.String())
		} else {
			log.Infof("Tx %s has not been processed", response.TxHash)
			unprocessedTxResponses[response.TxHash.String()] = response
			unprocessedTxsHashes = append(unprocessedTxsHashes, response.TxHash.String())
		}
	}
	return processedTxResponses, processedTxsHashes, unprocessedTxResponses, unprocessedTxsHashes
}

// WaitSequencingTxToBeSynced waits for a sequencing transaction to be synced into the state
func (s *State) WaitSequencingTxToBeSynced(parentCtx context.Context, tx *types.Transaction, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	for {
		virtualized, err := s.IsSequencingTXSynced(ctx, tx.Hash(), nil)
		if err != nil && err != ErrNotFound {
			log.Errorf("error waiting sequencing tx %s to be synced: %v", tx.Hash().String(), err)
			return err
		} else if ctx.Err() != nil {
			log.Errorf("error waiting sequencing tx %s to be synced: %v", tx.Hash().String(), err)
			return ctx.Err()
		} else if virtualized {
			break
		}

		time.Sleep(time.Second)
	}

	log.Debug("Sequencing txh successfully synced: ", tx.Hash().String())
	return nil
}

// WaitVerifiedBatchToBeSynced waits for a sequenced batch to be synced into the state
func (s *State) WaitVerifiedBatchToBeSynced(parentCtx context.Context, batchNumber uint64, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	for {
		batch, err := s.GetVerifiedBatch(ctx, batchNumber, nil)
		if err != nil && err != ErrNotFound {
			log.Errorf("error waiting verified batch [%d] to be synced: %v", batchNumber, err)
			return err
		} else if ctx.Err() != nil {
			log.Errorf("error waiting verified batch [%d] to be synced: %v", batchNumber, err)
			return ctx.Err()
		} else if batch != nil {
			break
		}

		time.Sleep(time.Second)
	}

	log.Debug("Verified batch successfully synced: ", batchNumber)
	return nil
}

func (s *State) monitorNewL2Blocks() {
	waitNextCycle := func() {
		time.Sleep(1 * time.Second)
	}

	for {
		if len(s.newL2BlockEventHandlers) == 0 {
			waitNextCycle()
			continue
		}

		lastL2Block, err := s.GetLastL2Block(context.Background(), nil)
		if errors.Is(err, ErrStateNotSynchronized) {
			waitNextCycle()
			continue
		} else if err != nil {
			log.Errorf("failed to get last l2 block while monitoring new blocks: %v", err)
			waitNextCycle()
			continue
		}

		// not updates until now
		if lastL2Block == nil || s.lastL2BlockSeen.NumberU64() >= lastL2Block.NumberU64() {
			waitNextCycle()
			continue
		}

		for bn := s.lastL2BlockSeen.NumberU64() + uint64(1); bn <= lastL2Block.NumberU64(); bn++ {
			block, err := s.GetL2BlockByNumber(context.Background(), bn, nil)
			if err != nil {
				log.Errorf("failed to l2 block while monitoring new blocks: %v", err)
				break
			}

			s.newL2BlockEvents <- NewL2BlockEvent{
				Block: *block,
			}
			log.Infof("new l2 blocks detected, Number %v, Hash %v", block.NumberU64(), block.Hash().String())
			s.lastL2BlockSeen = *block
		}

		// interval to check for new l2 blocks
		waitNextCycle()
	}
}

func (s *State) handleEvents() {
	for newL2BlockEvent := range s.newL2BlockEvents {
		if len(s.newL2BlockEventHandlers) == 0 {
			continue
		}

		wg := sync.WaitGroup{}
		for _, handler := range s.newL2BlockEventHandlers {
			wg.Add(1)
			go func(h NewL2BlockEventHandler) {
				defer func() {
					wg.Done()
					if r := recover(); r != nil {
						log.Errorf("failed and recovered in NewL2BlockEventHandler: %v", r)
					}
				}()
				h(newL2BlockEvent)
			}(handler)
		}
		wg.Wait()
	}
}

// NewL2BlockEventHandler represent a func that will be called by the
// state when a NewL2BlockEvent is triggered
type NewL2BlockEventHandler func(e NewL2BlockEvent)

// NewL2BlockEvent is a struct provided from the state to the NewL2BlockEventHandler
// when a new l2 block is detected with data related to this new l2 block.
type NewL2BlockEvent struct {
	Block types.Block
}

// RegisterNewL2BlockEventHandler add the provided handler to the list of handlers
// that will be triggered when a new l2 block event is triggered
func (s *State) RegisterNewL2BlockEventHandler(h NewL2BlockEventHandler) {
	log.Info("new l2 block event handler registered")
	s.newL2BlockEventHandlers = append(s.newL2BlockEventHandlers, h)
}

// StoreTransaction is used by the sequencer to add process a transaction
func (s *State) StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *ProcessTransactionResponse, coinbase common.Address, timestamp uint64, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	// Check if last batch is closed. Note that it's assumed that only the latest batch can be open
	/*
			isBatchClosed, err := s.PostgresStorage.IsBatchClosed(ctx, batchNumber, dbTx)
			if err != nil {
				return err
			}
			if isBatchClosed {
				return ErrBatchAlreadyClosed
			}

		processingContext, err := s.GetProcessingContext(ctx, batchNumber, dbTx)
		if err != nil {
			return err
		}
	*/
	// if the transaction has an intrinsic invalid tx error it means
	// the transaction has not changed the state, so we don't store it
	if executor.IsIntrinsicError(executor.RomErrorCode(processedTx.RomError)) {
		return nil
	}

	lastL2Block, err := s.GetLastL2Block(ctx, dbTx)
	if err != nil {
		return err
	}

	header := &types.Header{
		Number:     new(big.Int).SetUint64(lastL2Block.Number().Uint64() + 1),
		ParentHash: lastL2Block.Hash(),
		Coinbase:   coinbase,
		Root:       processedTx.StateRoot,
		GasUsed:    processedTx.GasUsed,
		GasLimit:   s.cfg.MaxCumulativeGasUsed,
		Time:       timestamp,
	}
	transactions := []*types.Transaction{&processedTx.Tx}

	receipt := generateReceipt(header.Number, processedTx)
	receipts := []*types.Receipt{receipt}

	// Create block to be able to calculate its hash
	block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
	block.ReceivedAt = time.Unix(int64(timestamp), 0)

	receipt.BlockHash = block.Hash()

	// Store L2 block and its transaction
	if err := s.AddL2Block(ctx, batchNumber, block, receipts, dbTx); err != nil {
		return err
	}

	return nil
}

// GetBalanceByStateRoot gets balance from the MT Service using the provided state root
func (s *State) GetBalanceByStateRoot(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error) {
	balance, err := s.tree.GetBalance(ctx, address, root.Bytes())
	if err != nil && balance == nil {
		balance = big.NewInt(0)
	}
	return balance, err
}

// GetNonceByStateRoot gets nonce from the MT Service using the provided state root
func (s *State) GetNonceByStateRoot(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error) {
	return s.tree.GetNonce(ctx, address, root.Bytes())
}

// LogExecutorError is used to store Executor error for runtime debugging
func (s *State) LogExecutorError(responseError pb.ExecutorError, processBatchRequest *pb.ProcessBatchRequest) {
	timestamp := time.Now()
	log.Errorf("error found in the executor: %v at %v", responseError, timestamp)
	payload, err := json.Marshal(processBatchRequest)
	if err != nil {
		log.Errorf("error marshaling payload: %v", err)
	} else {
		debugInfo := &DebugInfo{
			ErrorType: DebugInfoErrorType_EXECUTOR_ERROR + " " + responseError.String(),
			Timestamp: timestamp,
			Payload:   string(payload),
		}
		err = s.AddDebugInfo(context.Background(), debugInfo, nil)
		if err != nil {
			log.Errorf("error storing payload: %v", err)
		}
	}
}

// GetForkIdByBatchNumber returns the fork id for the given batch number
func (s *State) GetForkIdByBatchNumber(batchNumber uint64) uint64 {
	return GetForkIDByBatchNumber(s.cfg.ForkIDIntervals, batchNumber)
}

// FlushMerkleTree persists updates in the Merkle tree
func (s *State) FlushMerkleTree(ctx context.Context) error {
	return s.tree.Flush(ctx)
}
