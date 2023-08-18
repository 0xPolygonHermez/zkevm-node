package sequencer

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type workerAddTxTestCase struct {
	name                 string
	from                 common.Address
	txHash               common.Hash
	nonce                uint64
	cost                 *big.Int
	counters             state.ZKCounters
	usedBytes            uint64
	gasPrice             *big.Int
	expectedTxSortedList []common.Hash
}

type workerAddrQueueInfo struct {
	from    common.Address
	nonce   *big.Int
	balance *big.Int
}

func processWorkerAddTxTestCases(ctx context.Context, t *testing.T, worker *Worker, testCases []workerAddTxTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tx := TxTracker{}

			tx.Hash = testCase.txHash
			tx.HashStr = testCase.txHash.String()
			tx.From = testCase.from
			tx.FromStr = testCase.from.String()
			tx.Nonce = testCase.nonce
			tx.Cost = testCase.cost
			tx.BatchResources.Bytes = testCase.usedBytes
			tx.GasPrice = testCase.gasPrice
			tx.updateZKCounters(testCase.counters)
			t.Logf("%s=%d", testCase.name, tx.GasPrice)

			_, err := worker.AddTxTracker(ctx, &tx)
			if err != nil {
				return
			}

			el := worker.txSortedList
			if el.len() != len(testCase.expectedTxSortedList) {
				t.Fatalf("Error txSortedList.len(%d) != expectedTxSortedList.len(%d)", el.len(), len(testCase.expectedTxSortedList))
			}
			for i := 0; i < el.len(); i++ {
				if el.getByIndex(i).HashStr != string(testCase.expectedTxSortedList[i].String()) {
					t.Fatalf("Error txSortedList(%d). Expected=%s, Actual=%s", i, testCase.expectedTxSortedList[i].String(), el.getByIndex(i).HashStr)
				}
			}
		})
	}
}

func TestWorkerAddTx(t *testing.T) {
	var nilErr error

	stateMock := NewStateMock(t)
	worker := initWorker(stateMock)

	ctx := context.Background()

	stateMock.On("GetLastStateRoot", ctx, nil).Return(common.Hash{0}, nilErr)

	addrQueueInfo := []workerAddrQueueInfo{
		{from: common.Address{1}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{2}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{3}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{4}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
	}

	for _, aq := range addrQueueInfo {
		stateMock.On("GetNonceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.nonce, nilErr)
		stateMock.On("GetBalanceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.balance, nilErr)
	}

	addTxsTC := []workerAddTxTestCase{
		{
			name: "Adding from:0x01, tx:0x01/gp:10", from: common.Address{1}, txHash: common.Hash{1}, nonce: 1, gasPrice: new(big.Int).SetInt64(10),
			cost:      new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 1, UsedKeccakHashes: 1, UsedPoseidonHashes: 1, UsedPoseidonPaddings: 1, UsedMemAligns: 1, UsedArithmetics: 1, UsedBinaries: 1, UsedSteps: 1},
			usedBytes: 1,
			expectedTxSortedList: []common.Hash{
				{1},
			},
		},
		{
			name: "Adding from:0x02, tx:0x02/gp:4", from: common.Address{2}, txHash: common.Hash{2}, nonce: 1, gasPrice: new(big.Int).SetInt64(4),
			cost:      new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 1, UsedKeccakHashes: 1, UsedPoseidonHashes: 1, UsedPoseidonPaddings: 1, UsedMemAligns: 1, UsedArithmetics: 1, UsedBinaries: 1, UsedSteps: 1},
			usedBytes: 1,
			expectedTxSortedList: []common.Hash{
				{1}, {2},
			},
		},
		{
			name: "Readding from:0x02, tx:0x02/gp:20", from: common.Address{2}, txHash: common.Hash{2}, nonce: 1, gasPrice: new(big.Int).SetInt64(20),
			cost:      new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 5, UsedKeccakHashes: 5, UsedPoseidonHashes: 5, UsedPoseidonPaddings: 5, UsedMemAligns: 5, UsedArithmetics: 5, UsedBinaries: 5, UsedSteps: 5},
			usedBytes: 5,
			expectedTxSortedList: []common.Hash{
				{2}, {1},
			},
		},
		{
			name: "Readding from:0x03, tx:0x03/gp:25", from: common.Address{3}, txHash: common.Hash{3}, nonce: 1, gasPrice: new(big.Int).SetInt64(25),
			cost:      new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 2, UsedKeccakHashes: 2, UsedPoseidonHashes: 2, UsedPoseidonPaddings: 2, UsedMemAligns: 2, UsedArithmetics: 2, UsedBinaries: 2, UsedSteps: 2},
			usedBytes: 2,
			expectedTxSortedList: []common.Hash{
				{3}, {2}, {1},
			},
		},
		{
			name: "Adding from:0x04, tx:0x04/gp:100", from: common.Address{4}, txHash: common.Hash{4}, nonce: 1, gasPrice: new(big.Int).SetInt64(100),
			cost:      new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 1, UsedKeccakHashes: 1, UsedPoseidonHashes: 1, UsedPoseidonPaddings: 1, UsedMemAligns: 1, UsedArithmetics: 1, UsedBinaries: 1, UsedSteps: 1},
			usedBytes: 1,
			expectedTxSortedList: []common.Hash{
				{4}, {3}, {2}, {1},
			},
		},
	}

	processWorkerAddTxTestCases(ctx, t, worker, addTxsTC)
}

func TestWorkerGetBestTx(t *testing.T) {
	var nilErr error

	rc := state.BatchResources{
		ZKCounters: state.ZKCounters{CumulativeGasUsed: 10, UsedKeccakHashes: 10, UsedPoseidonHashes: 10, UsedPoseidonPaddings: 10, UsedMemAligns: 10, UsedArithmetics: 10, UsedBinaries: 10, UsedSteps: 10},
		Bytes:      10,
	}

	stateMock := NewStateMock(t)
	worker := initWorker(stateMock)

	ctx := context.Background()

	stateMock.On("GetLastStateRoot", ctx, nil).Return(common.Hash{0}, nilErr)

	addrQueueInfo := []workerAddrQueueInfo{
		{from: common.Address{1}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{2}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{3}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
		{from: common.Address{4}, nonce: new(big.Int).SetInt64(1), balance: new(big.Int).SetInt64(10)},
	}

	for _, aq := range addrQueueInfo {
		stateMock.On("GetNonceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.nonce, nilErr)
		stateMock.On("GetBalanceByStateRoot", ctx, aq.from, common.Hash{0}).Return(aq.balance, nilErr)
	}

	addTxsTC := []workerAddTxTestCase{
		{
			name: "Adding from:0x01, tx:0x01/gp:10", from: common.Address{1}, txHash: common.Hash{1}, nonce: 1, gasPrice: new(big.Int).SetInt64(10),
			cost:      new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 1, UsedKeccakHashes: 1, UsedPoseidonHashes: 1, UsedPoseidonPaddings: 1, UsedMemAligns: 1, UsedArithmetics: 1, UsedBinaries: 1, UsedSteps: 1},
			usedBytes: 1,
			expectedTxSortedList: []common.Hash{
				{1},
			},
		},
		{
			name: "Adding from:0x02, tx:0x02/gp:12", from: common.Address{2}, txHash: common.Hash{2}, nonce: 1, gasPrice: new(big.Int).SetInt64(12),
			cost:      new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 5, UsedKeccakHashes: 5, UsedPoseidonHashes: 5, UsedPoseidonPaddings: 5, UsedMemAligns: 5, UsedArithmetics: 5, UsedBinaries: 5, UsedSteps: 5},
			usedBytes: 5,
			expectedTxSortedList: []common.Hash{
				{2}, {1},
			},
		},
		{
			name: "Readding from:0x03, tx:0x03/gp:25", from: common.Address{3}, txHash: common.Hash{3}, nonce: 1, gasPrice: new(big.Int).SetInt64(25),
			cost:      new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 2, UsedKeccakHashes: 2, UsedPoseidonHashes: 2, UsedPoseidonPaddings: 2, UsedMemAligns: 2, UsedArithmetics: 2, UsedBinaries: 2, UsedSteps: 2},
			usedBytes: 2,
			expectedTxSortedList: []common.Hash{
				{3}, {2}, {1},
			},
		},
		{
			name: "Adding from:0x04, tx:0x04/gp:100", from: common.Address{4}, txHash: common.Hash{4}, nonce: 1, gasPrice: new(big.Int).SetInt64(100),
			cost:      new(big.Int).SetInt64(5),
			counters:  state.ZKCounters{CumulativeGasUsed: 4, UsedKeccakHashes: 4, UsedPoseidonHashes: 4, UsedPoseidonPaddings: 4, UsedMemAligns: 4, UsedArithmetics: 4, UsedBinaries: 4, UsedSteps: 4},
			usedBytes: 4,
			expectedTxSortedList: []common.Hash{
				{4}, {3}, {2}, {1},
			},
		},
	}

	processWorkerAddTxTestCases(ctx, t, worker, addTxsTC)

	expectedGetBestTx := []common.Hash{{4}, {3}, {1}}
	ct := 0

	for {
		tx := worker.GetBestFittingTx(rc)
		if tx != nil {
			if ct >= len(expectedGetBestTx) {
				t.Fatalf("Error getting more best tx than expected. Expected=%d, Actual=%d", len(expectedGetBestTx), ct+1)
			}
			if tx.HashStr != string(expectedGetBestTx[ct].String()) {
				t.Fatalf("Error GetBestFittingTx(%d). Expected=%s, Actual=%s", ct, expectedGetBestTx[ct].String(), tx.HashStr)
			}
			err := rc.Sub(tx.BatchResources)
			assert.NoError(t, err)

			touch := make(map[common.Address]*state.InfoReadWrite)
			var newNonce uint64 = tx.Nonce + 1
			touch[tx.From] = &state.InfoReadWrite{Address: tx.From, Nonce: &newNonce, Balance: new(big.Int).SetInt64(10)}
			worker.UpdateAfterSingleSuccessfulTxExecution(tx.From, touch)
			ct++
		} else {
			if ct < len(expectedGetBestTx) {
				t.Fatalf("Error expecting more best tx. Expected=%d, Actual=%d", len(expectedGetBestTx), ct)
			}
			break
		}
	}
}

func initWorker(stateMock *StateMock) *Worker {
	worker := NewWorker(stateMock)
	return worker
}
