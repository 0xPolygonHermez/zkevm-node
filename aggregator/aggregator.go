package aggregator

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
)

// Aggregator represents an aggregator
type Aggregator struct {
	cfg Config

	State          state.State
	BatchProcessor state.BatchProcessor
	EtherMan       etherman.EtherMan
	Prover         ProverClient

	ctx    context.Context
	cancel context.CancelFunc
}

// NewAggregator creates a new aggregator
func NewAggregator(
	cfg Config,
	state state.State,
	bp state.BatchProcessor,
	ethMan etherman.EtherMan,
	prover ProverClient,
) (Aggregator, error) {
	ctx, cancel := context.WithCancel(context.Background())

	a := Aggregator{
		cfg: cfg,

		State:          state,
		BatchProcessor: bp,
		EtherMan:       ethMan,
		Prover:         prover,

		ctx:    ctx,
		cancel: cancel,
	}

	return a, nil
}

// Start starts the aggregator
func (a *Aggregator) Start() {
	for {
		time.Sleep(a.cfg.IntervalToConsolidateState)

		// 1. find next batch to consolidate
		lastConsolidatedBatch, err := a.State.GetLastBatch(false)
		if err != nil {
			log.Error(err)
			continue
		}

		batchToConsolidate, err := a.State.GetBatchByNumber(lastConsolidatedBatch.BatchNumber + 1)
		if err != nil {
			log.Error(err)
			continue
		}

		// 2. check if it's profitable or not
		// check is it profitable to aggregate txs or not
		if !a.isProfitable(batchToConsolidate.Transactions()) {
			log.Info("Batch %d is not profitable", batchToConsolidate.Number().Uint64())
			continue
		}

		// // 3. send zki + txs to the prover
		// proof, err := a.Prover.SendTxs(batchToConsolidate.Transactions())
		// if err != nil {
		// 	log.Error(err)
		// 	continue
		// }

		// // 4. send proof + txs to the SC
		// h, err := a.EtherMan.ConsolidateBatch(batchToConsolidate, *proof)
		// if err != nil {
		// 	log.Error(err)
		// 	continue
		// }

		// log.Infof("Batch %d consolidated: %s", batchToConsolidate.Number().Uint64(), h.Hex())
	}
}

func (a *Aggregator) isProfitable(txs []*types.Transaction) bool {
	// get strategy from the config and check
	return true
}

// Stop stops the aggregator
func (a *Aggregator) Stop() {
	a.cancel()
}
