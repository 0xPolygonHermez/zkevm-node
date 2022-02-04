package txselector_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txselector"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBase_SelectTxs(t *testing.T) {
	seqAddress := common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")

	bp := new(mocks.BatchProcessor)

	txSelector := txselector.NewTxSelectorBase(txselector.Config{
		TxSorterType: "bycostandnonce",
	})

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(12), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(14), []byte{})
	txs := []pool.Transaction{{Transaction: *tx2}, {Transaction: *tx1}, {Transaction: *tx3}}

	bp.On("ProcessTransaction", tx1, seqAddress).Return(state.ErrInvalidBalance)
	bp.On("ProcessTransaction", tx2, seqAddress).Return(nil)
	bp.On("ProcessTransaction", tx3, seqAddress).Return(state.ErrInvalidSig)
	selectedTxs, selectedTxsHashes, invalidTxsHashes, err := txSelector.SelectTxs(bp, txs, seqAddress)
	bp.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, len(selectedTxs), 1)
	assert.Equal(t, len(selectedTxsHashes), 1)
	assert.Equal(t, len(invalidTxsHashes), 2)
}

func TestAcceptAll_SelectTxs(t *testing.T) {
	seqAddress := common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")

	bp := new(mocks.BatchProcessor)

	txSelector := txselector.NewTxSelectorAcceptAll()

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(12), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(14), []byte{})
	txs := []pool.Transaction{{Transaction: *tx2}, {Transaction: *tx1}, {Transaction: *tx3}}

	bp.On("ProcessTransaction", tx1, seqAddress).Return(state.ErrInvalidBalance)
	bp.On("ProcessTransaction", tx2, seqAddress).Return(state.ErrInvalidNonce)
	bp.On("ProcessTransaction", tx3, seqAddress).Return(state.ErrInvalidSig)
	selectedTxs, selectedTxsHashes, invalidTxsHashes, err := txSelector.SelectTxs(bp, txs, seqAddress)
	bp.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, len(selectedTxs), 3)
	assert.Equal(t, len(selectedTxsHashes), 3)
	assert.Equal(t, len(invalidTxsHashes), 0)
}
