//nolint
package mocks

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

type StateMock struct{}

func NewState() state.State {
	return &StateMock{}
}

func (s *StateMock) NewBatchProcessor(lastBatchNumber uint64, withProofCalculation bool) state.BatchProcessor {
	return &state.BasicBatchProcessor{}
}

func (s *StateMock) GetStateRoot(ctx context.Context, virtual bool) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (s *StateMock) GetBalance(address common.Address, batchNumber uint64) (*big.Int, error) {
	return big.NewInt(balance), nil
}

func (s *StateMock) EstimateGas(transaction *types.Transaction) uint64 {
	return estimatedGas
}

func (s *StateMock) GetLastBlock(ctx context.Context) (*state.Block, error) {
	return block, nil
}

func (s *StateMock) GetPreviousBlock(ctx context.Context, offset uint64) (*state.Block, error) {
	return block, nil
}

func (s *StateMock) GetBlockByHash(ctx context.Context, hash common.Hash) (*state.Block, error) {
	return block, nil
}

func (s *StateMock) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*state.Block, error) {
	return block, nil
}

func (s *StateMock) GetLastBlockNumber(ctx context.Context) (uint64, error) {
	return blockNumber, nil
}

func (s *StateMock) GetLastBatch(ctx context.Context, isVirtual bool) (*state.Batch, error) {
	return batch, nil
}

func (s *StateMock) GetTransaction(ctx context.Context, hash common.Hash) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetNonce(address common.Address, batchNumber uint64) (uint64, error) {
	return txNonce, nil
}

func (s *StateMock) GetPreviousBatch(ctx context.Context, isVirtual bool, offset uint64) (*state.Batch, error) {
	return batch, nil
}

func (s *StateMock) GetBatchByHash(ctx context.Context, hash common.Hash) (*state.Batch, error) {
	return batch, nil
}

func (s *StateMock) GetBatchByNumber(ctx context.Context, batchNumber uint64) (*state.Batch, error) {
	return batch, nil
}

func (s *StateMock) GetLastBatchNumber(ctx context.Context) (uint64, error) {
	return batchNumber, nil
}

func (s *StateMock) GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error) {
	return tx, nil
}

func (s *StateMock) GetTransactionCount(ctx context.Context, address common.Address) (uint64, error) {
	return txNonce, nil
}

func (s *StateMock) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*types.Receipt, error) {
	return txReceipt, nil
}

func (s *StateMock) Reset(blockNumber uint64) error {
	return nil
}

func (s *StateMock) ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash) error {
	return nil
}

func (s *StateMock) GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error) {
	return []*types.Transaction{
		tx,
	}, nil
}

func (s *StateMock) AddSequencer(ctx context.Context, seq state.Sequencer) error {
	return nil
}

func (s *StateMock) GetSequencerByChainID(ctx context.Context, chainID *big.Int) (*state.Sequencer, error) {
	return nil, nil
}

func (s *StateMock) SetGenesis(ctx context.Context, genesis state.Genesis) error {
	return nil
}

func (s *StateMock) AddBlock(ctx context.Context, block *state.Block) error {
	return nil
}

func (s *StateMock) SetLastBatchNumberSeenOnEthereum(batchNumber uint64) error {
	return nil
}

func (s *StateMock) GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error) {
	return 0, nil
}
