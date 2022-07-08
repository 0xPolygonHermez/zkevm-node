package txselector_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/strategy/txselector"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestTxSorterByCostAndNonce_SortTxs_SameCost(t *testing.T) {
	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx4 := types.NewTransaction(uint64(3), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx5 := types.NewTransaction(uint64(4), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})

	txs := []pool.Transaction{{Transaction: *tx4}, {Transaction: *tx2}, {Transaction: *tx5}, {Transaction: *tx1}, {Transaction: *tx3}}
	txSorter := &txselector.TxSorterByCostAndNonce{}
	sortedTxs := txSorter.SortTxs(txs)
	nonce := uint64(0)
	for _, v := range sortedTxs {
		assert.Equal(t, nonce, v.Nonce())
		nonce++
	}
}

func TestTxSorterByCostAndNonce_SortTxs_DiffCost(t *testing.T) {
	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(12), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(14), []byte{})
	tx4 := types.NewTransaction(uint64(3), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(16), []byte{})
	tx5 := types.NewTransaction(uint64(4), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(18), []byte{})
	txs := []pool.Transaction{{Transaction: *tx4}, {Transaction: *tx2}, {Transaction: *tx5}, {Transaction: *tx1}, {Transaction: *tx3}}
	txSorter := &txselector.TxSorterByCostAndNonce{}
	sortedTxs := txSorter.SortTxs(txs)
	for i, v := range sortedTxs {
		assert.Equal(t, 0, v.Cost().Cmp(big.NewInt(int64(28-i*2)))) // it's gas price + tx value (10), start from biggest value
	}
}

func TestTxSorterByCostAndTime_SortTxs_SameCost(t *testing.T) {
	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx4 := types.NewTransaction(uint64(3), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx5 := types.NewTransaction(uint64(4), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})

	receivedAt := time.Now()

	txs := []pool.Transaction{
		{Transaction: *tx4, ReceivedAt: receivedAt.Add(4)},
		{Transaction: *tx2, ReceivedAt: receivedAt.Add(2)},
		{Transaction: *tx5, ReceivedAt: receivedAt.Add(5)},
		{Transaction: *tx1, ReceivedAt: receivedAt.Add(1)},
		{Transaction: *tx3, ReceivedAt: receivedAt.Add(3)},
	}

	txSorter := &txselector.TxSorterByCostAndTime{}
	sortedTxs := txSorter.SortTxs(txs)
	isSortingRight := true
	for i := 0; i < len(sortedTxs)-1; i++ {
		if sortedTxs[i].ReceivedAt.Unix() > sortedTxs[i+1].ReceivedAt.Unix() {
			isSortingRight = false
			break
		}
	}
	assert.True(t, isSortingRight)
}

func TestTxSorterByCostAndTime_SortTxs_DiffCost(t *testing.T) {
	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	tx2 := types.NewTransaction(uint64(1), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(12), []byte{})
	tx3 := types.NewTransaction(uint64(2), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(14), []byte{})
	tx4 := types.NewTransaction(uint64(3), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(16), []byte{})
	tx5 := types.NewTransaction(uint64(4), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(18), []byte{})
	txs := []pool.Transaction{{Transaction: *tx4}, {Transaction: *tx2}, {Transaction: *tx5}, {Transaction: *tx1}, {Transaction: *tx3}}
	txSorter := &txselector.TxSorterByCostAndTime{}
	sortedTxs := txSorter.SortTxs(txs)
	for i, v := range sortedTxs {
		assert.Equal(t, 0, v.Cost().Cmp(big.NewInt(int64(28-i*2)))) // it's gas price + tx value (10), start from biggest value
	}
}
