package aggregator

import (
	"github.com/ethereum/go-ethereum/core/types"
	eth "github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/state"
)

type Aggregator struct {
	State          state.State
	BatchProcessor state.BatchProcessor
	EthClient      eth.Client
	Synchronizer   SynchronizerClient
}

func NewAggregator(cfg Config, state state.State, bp state.BatchProcessor, ethClient eth.Client, syncClient SynchronizerClient) (Aggregator, error) {
	return Aggregator{
		State:          state,
		BatchProcessor: bp,
		EthClient:      ethClient,
		Synchronizer:   syncClient,
	}, nil
}

func (agr *Aggregator) generateAndSendProofs() {
	// reads from batchesChan
	// get txs from state by batchNum
	// check if it's profitable or not
	// send proof + txs to the prover
	// send proof + txs to the SC
}

func (agr *Aggregator) isProfitable(txs []types.Transaction) bool {
	// get strategy from the config and check
	panic("not implemented yet")
}

func (agr *Aggregator) Run() {

}
