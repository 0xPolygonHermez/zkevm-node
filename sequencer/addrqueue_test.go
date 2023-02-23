package sequencer

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type notReadyTx struct {
	nonce uint64
	hash  common.Hash
}

type addrQueueAddTxTestCase struct {
	name               string
	hash               common.Hash
	nonce              uint64
	gasPrice           *big.Int
	cost               *big.Int
	expectedReadyTx    common.Hash
	expectedNotReadyTx []notReadyTx
}

var addr addrQueue

func newTestTxTracker(hash common.Hash, nonce uint64, gasPrice *big.Int, cost *big.Int) *TxTracker {
	tx := TxTracker{Hash: hash, Nonce: nonce, GasPrice: gasPrice, Cost: cost}
	tx.HashStr = tx.Hash.String()
	return &tx
}

func processAddTxTestCases(t *testing.T, testCases []addrQueueAddTxTestCase) {
	var emptyHash common.Hash = common.Hash{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx := newTestTxTracker(tc.hash, tc.nonce, tc.gasPrice, tc.cost)
			newReadyTx, _ := addr.addTx(tx)
			if tc.expectedReadyTx.String() == emptyHash.String() {
				if !(addr.readyTx == nil) {
					t.Fatalf("Error readyTx. Expected=%s, Actual=%s", tc.expectedReadyTx, "")
				}
				if !(newReadyTx == nil) {
					t.Fatalf("Error newReadyTx. Expected=nil, Actual=%s", newReadyTx.HashStr)
				}
			} else {
				if !(addr.readyTx.Hash == tc.expectedReadyTx) {
					t.Fatalf("Error readyTx. Expected=%s, Actual=%s", tc.expectedReadyTx, addr.readyTx.HashStr)
				}
			}

			for _, nr := range tc.expectedNotReadyTx {
				txTmp, found := addr.notReadyTxs[nr.nonce]
				if !(found) {
					t.Fatalf("Error notReadyTx nonce=%d don't exists", nr.nonce)
				}
				if !(txTmp.Hash == nr.hash) {
					t.Fatalf("Error notReadyTx nonce=%d. Expected=%s, Actual=%s", nr.nonce, nr.hash.String(), txTmp.HashStr)
				}
			}
		})
	}
}

func TestAddrQueue(t *testing.T) {
	addr = addrQueue{fromStr: "0x99999", currentNonce: 1, currentBalance: new(big.Int).SetInt64(10), notReadyTxs: make(map[uint64]*TxTracker)}

	addTxTestCases := []addrQueueAddTxTestCase{
		{
			name: "Add not ready tx 0x02", hash: common.Hash{0x2}, nonce: 2, gasPrice: new(big.Int).SetInt64(2), cost: new(big.Int).SetInt64(5),
			expectedReadyTx: common.Hash{},
			expectedNotReadyTx: []notReadyTx{
				{nonce: 2, hash: common.Hash{0x2}},
			},
		},
		{
			name: "Add ready tx 0x01", hash: common.Hash{0x1}, nonce: 1, gasPrice: new(big.Int).SetInt64(2), cost: new(big.Int).SetInt64(5),
			expectedReadyTx: common.Hash{1},
			expectedNotReadyTx: []notReadyTx{
				{nonce: 2, hash: common.Hash{0x2}},
			},
		},
		{
			name: "Add not ready tx 0x04", hash: common.Hash{0x4}, nonce: 4, gasPrice: new(big.Int).SetInt64(2), cost: new(big.Int).SetInt64(5),
			expectedReadyTx: common.Hash{1},
			expectedNotReadyTx: []notReadyTx{
				{nonce: 2, hash: common.Hash{0x2}},
				{nonce: 4, hash: common.Hash{0x4}},
			},
		},
		{
			name: "Replace tx with nonce 4 for tx 0x44 with best GasPrice", hash: common.Hash{0x44}, nonce: 4, gasPrice: new(big.Int).SetInt64(3), cost: new(big.Int).SetInt64(5),
			expectedReadyTx: common.Hash{1},
			expectedNotReadyTx: []notReadyTx{
				{nonce: 2, hash: common.Hash{0x2}},
				{nonce: 4, hash: common.Hash{0x44}},
			},
		},
	}

	processAddTxTestCases(t, addTxTestCases)

	t.Run("Delete readyTx 0x01", func(t *testing.T) {
		tc := addTxTestCases[1]
		tx := newTestTxTracker(tc.hash, tc.nonce, tc.gasPrice, tc.cost)
		deltx := addr.deleteTx(tx.Hash)
		if !(addr.readyTx == nil) {
			t.Fatalf("Error readyTx not nil. Expected=%s, Actual=%s", "", addr.readyTx.HashStr)
		}
		if !(deltx.HashStr == tx.HashStr) {
			t.Fatalf("Error returning deletedReadyTx. Expected=%s, Actual=%s", tx.HashStr, deltx.HashStr)
		}
	})

	processAddTxTestCases(t, []addrQueueAddTxTestCase{
		{
			name: "Add tx with nonce = currentNonce but with cost > currentBalance", hash: common.Hash{0x11}, nonce: 1, gasPrice: new(big.Int).SetInt64(2), cost: new(big.Int).SetInt64(15),
			expectedReadyTx: common.Hash{},
			expectedNotReadyTx: []notReadyTx{
				{nonce: 1, hash: common.Hash{0x11}},
				{nonce: 2, hash: common.Hash{0x2}},
				{nonce: 4, hash: common.Hash{0x44}},
			},
		},
	})

	t.Run("Update currentBalance = 15, set tx 0x11 as ready", func(t *testing.T) {
		tmpHash := common.Hash{0x11}
		addr.updateCurrentNonceBalance(&addr.currentNonce, new(big.Int).SetInt64(15))
		if !(addr.readyTx != nil && addr.readyTx.Hash.String() == tmpHash.String()) {
			t.Fatalf("Error readyTx. Expected=%s, Actual=%s", tmpHash, "")
		}

		tx, found := addr.notReadyTxs[1]

		if found {
			t.Fatalf("Error notReadyTx nonce=%d. Expected=%s, Actual=%s", addr.currentNonce, "", tx.Hash.String())
		}
	})

	t.Run("Update currentNonce = 4, set tx 0x04 as ready", func(t *testing.T) {
		tmpHash := common.Hash{0x44}
		newNonce := uint64(4)
		addr.updateCurrentNonceBalance(&newNonce, new(big.Int).SetInt64(15))
		if !(addr.readyTx != nil && addr.readyTx.Hash.String() == tmpHash.String()) {
			t.Fatalf("Error readyTx. Expected=%s, Actual=%s", tmpHash, addr.readyTx.Hash.String())
		}

		if len(addr.notReadyTxs) > 0 {
			t.Fatalf("Error notReadyTx not empty. Expected=%d, Actual=%d", 0, len(addr.notReadyTxs))
		}
	})
}
