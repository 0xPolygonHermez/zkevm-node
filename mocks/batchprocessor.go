package mocks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

// BatchProcessorMock mocked version of BatchProcessor
type BatchProcessorMock struct{}

// ProcessBatch mock
func (bp *BatchProcessorMock) ProcessBatch(batch *state.Batch) error {
	return nil
}

// ProcessTransaction mock
func (bp *BatchProcessorMock) ProcessTransaction(tx *types.Transaction, sequencerAddress common.Address) error {
	return nil
}

// CheckTransaction mock
func (bp *BatchProcessorMock) CheckTransaction(tx *types.Transaction) (common.Address, *big.Int, *big.Int, error) {
	return common.Address{}, big.NewInt(0), big.NewInt(0), nil
}
