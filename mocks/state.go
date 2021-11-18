package mocks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

type StateMock struct{}

func NewState() state.State {
	return &StateMock{}
}

func (s *StateMock) NewBatchProcessor(startingHash common.Hash, withProofCalculation bool) state.BatchProcessor {
	return &state.BasicBatchProcessor{}
}

func (s *StateMock) GetStateRoot(virtual bool) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (s *StateMock) GetBalance(address common.Address, batchNumber uint64) (*big.Int, error) {
	return big.NewInt(balance), nil
}

func (s *StateMock) EstimageGas(transaction types.Transaction) uint64 {
	return estimatedGas
}

func (s *StateMock) GetLastBlock() (*types.Block, error) {
	return block, nil
}

func (s *StateMock) GetPreviousBlock(offset uint64) (*types.Block, error) {
	return block, nil
}

func (s *StateMock) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	return block, nil
}

func (s *StateMock) GetBlockByNumber(blockNumber uint64) (*types.Block, error) {
	return block, nil
}

func (s *StateMock) GetLastBlockNumber() (uint64, error) {
	return blockNumber, nil
}

func (s *StateMock) GetLastBatch(isVirtual bool) (*state.Batch, error) {
	return batch, nil
}

func (s *StateMock) GetTransaction(hash common.Hash) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetNonce(address common.Address, batchNumber uint64) (uint64, error) {
	return txNonce, nil
}

func (s *StateMock) GetPreviousBatch(offset uint64) (*state.Batch, error) {
	return batch, nil
}

func (s *StateMock) GetBatchByHash(hash common.Hash) (*types.Block, error) {
	return block, nil
}

func (s *StateMock) GetBatchByNumber(batchNumber uint64) (*types.Block, error) {
	return block, nil
}

func (s *StateMock) GetLastBatchNumber() (uint64, error) {
	return batchNumber, nil
}

func (s *StateMock) GetTransactionByBatchHashAndIndex(batchHash common.Hash, index uint64) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetTransactionByBatchNumberAndIndex(batchNumber uint64, index uint64) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetTransactionByHash(transactionHash common.Hash) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetTransactionCount(address common.Address) (uint64, error) {
	return txNonce, nil
}

func (s *StateMock) GetTransactionReceipt(transactionHash common.Hash) (*types.Receipt, error) {
	return txReceipt, nil
}

func (s *StateMock) Reset(blockNumber uint64) error {
	return nil
}

func (s *StateMock) ConsolidateBatch(batchNumber uint64) error {
	return nil
}
