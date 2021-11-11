package aggregator

import "github.com/ethereum/go-ethereum/core/types"

type Aggregator struct {
	State          state.State
	BatchProcessor state.BatchProcessor
	EthClient      eth.Client
	Synchronizer   SynchronizerClient
}

func NewAggregator(cfg Config) (Aggregator, error) {
	state := state.NewState()
	bp := state.NewBatchProcessor(cfg.StartingHash, cfg.WithProofCalulation)
	ethClient := eth.NewClient()
	synchronizerClient := NewSynchronizerClient()
	return Aggregator{
		State:          state,
		BatchProcessor: bp,
		EthClient:      ethClient,
		Synchronizer:   synchronizerClient,
	}, nil
}

func (agr *Aggregator) generateAndSendProofs() {
	// reads from batchesChan
	// TODO: get txs from MT by batchNum
	// check if it's profitable or not
	// send proof + txs to the prover
	// send proof + txs to the SC
}

func (agr *Aggregator) isProfitable(txs []types.Transaction) bool {
	// get strategy from the config and check
}

func (agr *Aggregator) Run() {

}
