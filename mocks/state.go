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
	panic("not implemented yet")
}

func (s *StateMock) GetBalance(address common.Address, batchNumber uint64) (*big.Int, error) {
	return big.NewInt(balance), nil
}

func (s *StateMock) EstimageGas(address common.Address) uint64 {
	return estimatedGas
}

func (s *StateMock) GetLastBlock() (*types.Block, error) {
	return block, nil
}

func (s *StateMock) GetLastBatch(isVirtual bool) (*state.Batch, error) {
	return batch, nil
}

func (s *StateMock) GetBatchByHash(hash common.Hash, withTxDetails, isVirtual bool) (*state.Batch, error) {
	return batch, nil
}

func (s *StateMock) GetBatchByNumber(number uint64, withTxDetails, isVirtual bool) (*state.Batch, error) {
	return batch, nil
}

func (s *StateMock) GetTransactionByBatchHashAndIndex(hash common.Hash, index uint64) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetTransactionByBatchNumberAndIndex(number uint64, index uint64) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetTransaction(hash common.Hash) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	return txReceipt, nil
}

func (s *StateMock) GetNonce(address common.Address, batchNumber uint64) (uint64, error) {
	return txNonce, nil
}

func (s *StateMock) Reset(batchnum uint64) error {
	return nil
}

func (s *StateMock) ConsolidateBatch(batch state.Batch) error {
	return nil
}
