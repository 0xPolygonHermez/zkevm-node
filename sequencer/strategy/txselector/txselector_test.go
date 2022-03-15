package txselector_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/crypto"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txselector"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/stretchr/testify/assert"
)

func TestBase_SelectTxs(t *testing.T) {
	seqAddress := common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")

	bp := new(batchProcessor)

	txSelector := txselector.NewTxSelectorBase(txselector.Config{
		TxSorterType: "bycostandnonce",
	})

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(12), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(14), []byte{})
	tx4 := types.NewTransaction(uint64(100), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(16), []byte{})
	txs := []pool.Transaction{{Transaction: *tx2}, {Transaction: *tx1}, {Transaction: *tx4}, {Transaction: *tx3}}

	ctx := context.Background()
	bp.On("ProcessTransaction", ctx, tx1, seqAddress).Return(&runtime.ExecutionResult{Err: state.ErrInvalidBalance})
	bp.On("ProcessTransaction", ctx, tx2, seqAddress).Return(&runtime.ExecutionResult{})
	bp.On("ProcessTransaction", ctx, tx3, seqAddress).Return(&runtime.ExecutionResult{Err: crypto.ErrInvalidSig})
	bp.On("ProcessTransaction", ctx, tx4, seqAddress).Return(&runtime.ExecutionResult{Err: state.ErrNonceIsBiggerThanAccountNonce})

	selectedTxs, selectedTxsHashes, invalidTxsHashes, err := txSelector.SelectTxs(ctx, bp, txs, seqAddress)
	bp.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(selectedTxs))
	assert.Equal(t, 1, len(selectedTxsHashes))
	assert.Equal(t, 2, len(invalidTxsHashes))
}

func TestBase_SelectTxs_ExceededGasLimit(t *testing.T) {
	seqAddress := common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")

	bp := new(batchProcessor)

	txSelector := txselector.NewTxSelectorBase(txselector.Config{
		TxSorterType: "bycostandnonce",
	})

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(16), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(14), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(12), []byte{})
	tx4 := types.NewTransaction(uint64(3), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	txs := []pool.Transaction{{Transaction: *tx2}, {Transaction: *tx1}, {Transaction: *tx4}, {Transaction: *tx3}}

	ctx := context.Background()
	bp.On("ProcessTransaction", ctx, tx1, seqAddress).Return(&runtime.ExecutionResult{})
	bp.On("ProcessTransaction", ctx, tx2, seqAddress).Return(&runtime.ExecutionResult{})
	bp.On("ProcessTransaction", ctx, tx3, seqAddress).Return(&runtime.ExecutionResult{Err: state.ErrInvalidCumulativeGas})
	bp.AssertNotCalled(t, "ProcessTransaction", tx4, seqAddress)

	selectedTxs, selectedTxsHashes, invalidTxsHashes, err := txSelector.SelectTxs(ctx, bp, txs, seqAddress)
	bp.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(selectedTxs))
	assert.Equal(t, 2, len(selectedTxsHashes))
	assert.Equal(t, 0, len(invalidTxsHashes))
}

func TestAcceptAll_SelectTxs(t *testing.T) {
	seqAddress := common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")

	bp := new(batchProcessor)

	txSelector := txselector.NewTxSelectorAcceptAll()

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(12), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(14), []byte{})
	txs := []pool.Transaction{{Transaction: *tx2}, {Transaction: *tx1}, {Transaction: *tx3}}

	ctx := context.Background()

	selectedTxs, selectedTxsHashes, invalidTxsHashes, err := txSelector.SelectTxs(ctx, bp, txs, seqAddress)
	bp.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(selectedTxs))
	assert.Equal(t, 3, len(selectedTxsHashes))
	assert.Equal(t, 0, len(invalidTxsHashes))
}
