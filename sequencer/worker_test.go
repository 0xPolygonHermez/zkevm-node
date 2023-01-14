package sequencer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type workerAddTxTestCase struct {
	Name                   string
	Hash                   common.Hash
	Nonce                  uint64
	GasPrice               *big.Int
	Cost                   *big.Int
	expectedEfficiencyList []TxTracker
}

/*
func TestWorker(t *testing.T) {
	cfg := Config{}
	worker := NewWorker(cfg)

	addTxs := []struct {

	}
}
*/
