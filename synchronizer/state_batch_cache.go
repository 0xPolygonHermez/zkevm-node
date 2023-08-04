package synchronizer

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// StateBatchCacher is a decorator over interface stateInterface that create a cache
type StateBatchCacher interface {
	Set(batch *state.Batch)
	Get(batchNumber uint64) *state.Batch
	CleanCache()

	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	ProcessBatch(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	UpdateBatchL2Data(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error
}

// SynchronizerStateBatchCache it's a kind of decorator over stateInterface
// it keeps
type SynchronizerStateBatchCache struct {
	stateInterface
	Capacity int

	Mutex         sync.Mutex
	Cache         []*state.Batch
	LastStateRoot *common.Hash
	Head          int
	Tail          int
}

// NewSynchronizerStateBatchCache create a new struct
func NewSynchronizerStateBatchCache(statei stateInterface, capacity int) *SynchronizerStateBatchCache {
	return &SynchronizerStateBatchCache{
		stateInterface: statei,
		Capacity:       capacity,
		Cache:          make([]*state.Batch, capacity),
		LastStateRoot:  nil,
		Head:           -1,
		Tail:           -1,
	}
}

// UpdateBatchL2Data is a decorator for state.UpdateBatchL2Data
func (s *SynchronizerStateBatchCache) UpdateBatchL2Data(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error {
	res := s.stateInterface.UpdateBatchL2Data(ctx, batchNumber, batchL2Data, dbTx)
	if res == nil {
		cached_batch := s.getBatchByNumberInCache(batchNumber)
		if cached_batch != nil {
			// Update the cache
			updateBatchL2Data(cached_batch, batchL2Data)
		}
	}
	return res
}

// OpenBatch is a decorator for state.OpenBatch
func (s *SynchronizerStateBatchCache) OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error {
	res := s.stateInterface.OpenBatch(ctx, processingContext, dbTx)
	if res == nil {
		batchL2Data := processingContext.BatchL2Data
		if batchL2Data == nil {
			emptybatchL2Data := []byte{}
			batchL2Data = &emptybatchL2Data
		}
		newBatch := state.Batch{
			BatchNumber:    processingContext.BatchNumber,
			GlobalExitRoot: processingContext.GlobalExitRoot,
			Timestamp:      processingContext.Timestamp,
			Coinbase:       processingContext.Coinbase,
			ForcedBatchNum: processingContext.ForcedBatchNum,
			BatchL2Data:    *batchL2Data,
		}
		s.Set(&newBatch)
	}
	return res
}

// CloseBatch is a decorator for state.CloseBatch
func (s *SynchronizerStateBatchCache) CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error {
	err := s.stateInterface.CloseBatch(ctx, receipt, dbTx)
	if err == nil {
		cached_batch := s.getBatchByNumberInCache(receipt.BatchNumber)
		if cached_batch != nil {
			// Update the cache
			updateBatchWithReceipt(cached_batch, receipt)
		}
	}
	return err
}

func updateBatchWithReceipt(batch *state.Batch, receipt state.ProcessingReceipt) {
	if batch == nil {
		batch.BatchL2Data = receipt.BatchL2Data
		batch.StateRoot = receipt.StateRoot
		batch.LocalExitRoot = receipt.LocalExitRoot
		batch.AccInputHash = receipt.AccInputHash
		// receipt.batch_resources are not read from db by GetBatchByNumber
		// receipt.closing_reason is not read from db by GetBatchByNumber
		// batch.GlobalExitRoot is not in receipt
	}
}

func updateBatchL2Data(batch *state.Batch, batchL2Data []byte) {
	if batch != nil {
		batch.BatchL2Data = batchL2Data
	}
}

// ProcessBatch is a decorator for state.ProcessBatch
func (s *SynchronizerStateBatchCache) ProcessBatch(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error) {
	//s.CleanCache()
	// Does not modify state Database, it calls to executor
	return s.stateInterface.ProcessBatch(ctx, request, updateMerkleTree)
}

// ResetTrustedState is a decorator for state.ResetTrustedState
func (s *SynchronizerStateBatchCache) ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	s.CleanCache()
	return s.stateInterface.ResetTrustedState(ctx, batchNumber, dbTx)
}

// CleanCache remove all items from the cache
func (s *SynchronizerStateBatchCache) CleanCache() {
	s.Head = -1
	s.Tail = -1
}

// Set store a batch in the cache, if exist a previous one overwrite it
func (s *SynchronizerStateBatchCache) Set(batch *state.Batch) {
	if batch == nil {
		return
	}
	idx, err := s.getIndexByBatchNumber(batch.BatchNumber)
	if err == nil {
		s.Cache[idx] = batch
		return
	}
	s.emplace(batch)
}

func (s *SynchronizerStateBatchCache) emplace(batch *state.Batch) {
	newTail := (s.Tail + 1) % s.Capacity
	s.Cache[newTail] = batch
	s.Tail = newTail
	if s.Head == -1 {
		s.Head = 0
	} else if s.Tail == s.Head {
		s.Head = (s.Head + 1) % s.Capacity
	}
}

// Get a cached element by batchNumber
func (s *SynchronizerStateBatchCache) Get(batchNumber uint64) *state.Batch {
	idx, err := s.getIndexByBatchNumber(batchNumber)
	if err == nil {
		return s.Cache[idx]
	}
	return nil
}

// GetBatchByNumber get cached version or call to State.GetBatchByNumber
func (s *SynchronizerStateBatchCache) GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error) {
	//batch, err := s.state.GetBatchByNumber(ctx, batchNumber, dbTx)
	batch := s.getBatchByNumberInCache(batchNumber)
	if batch != nil {
		return batch, nil
	}

	if s.stateInterface == nil {
		return nil, errors.New("to allow retrieve [" + strconv.FormatUint(batchNumber, 10) + "] form db please pass a State != nil")
	}
	batch, err := s.stateInterface.GetBatchByNumber(ctx, batchNumber, dbTx)
	if batch != nil && err == nil {
		s.Set(batch)
	}

	return batch, err
}

func (s *SynchronizerStateBatchCache) getIndexByBatchNumber(batchNumber uint64) (int, error) {
	if !s.isEmpty() {
		for i := s.Head; i < s.Head+s.numElements(); i++ {
			idx := i % s.Capacity
			if s.Cache[idx].BatchNumber == batchNumber {
				return idx, nil
			}
		}
	}
	return -1, errors.New("not found batch in cache")
}

// Returns a pointer to the batch if it is in the cache, otherwise nil
func (s *SynchronizerStateBatchCache) getBatchByNumberInCache(batchNumber uint64) *state.Batch {
	idx, err := s.getIndexByBatchNumber(batchNumber)
	if err == nil {
		return s.Cache[idx]
	}
	return nil
}

func (s *SynchronizerStateBatchCache) numElements() int {
	if s.isEmpty() {
		return 0
	}

	if s.Head < s.Tail {
		return s.Tail - s.Head + 1
	}
	if s.Head > s.Tail {
		return s.Head - s.Tail + 1
	}
	return 1
}

func (s *SynchronizerStateBatchCache) isEmpty() bool {
	return s.Head < 0 || s.Tail < 0
}
