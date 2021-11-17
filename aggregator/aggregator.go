package aggregator

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state"
)

type Aggregator struct {
	State          state.State
	BatchProcessor state.BatchProcessor
	EtherMan       etherman.EtherMan
	Synchronizer   SynchronizerClient
	Prover         ProverClient

	ctx    context.Context
	cancel context.CancelFunc
}

func NewAggregator(
	cfg Config,
	state state.State,
	bp state.BatchProcessor,
	ethClient etherman.EtherMan,
	syncClient SynchronizerClient,
	prover ProverClient,
) (Aggregator, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return Aggregator{
		State:          state,
		BatchProcessor: bp,
		EtherMan:       ethClient,
		Synchronizer:   syncClient,
		Prover:         prover,

		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (agr *Aggregator) Start() {
	// reads from batchesChan
	// get txs from state by batchNum
	// check if it's profitable or not
	// send zki + txs to the prover
	// send proof + txs to the SC
	go func() {
		for {
			select {
			case <-agr.ctx.Done():
				return
			case event := <-agr.Synchronizer.SyncEventChan:
				// get txs to send
				txs, err := agr.State.GetTxsByBatchNum(event.BatchNum)
				if err != nil {
					return
				}
				// check is it profitable to aggregate txs or not
				if !agr.isProfitable(txs) {
					continue
				}
				// send txs and zki to the prover
				proof, err := agr.Prover.SendTxsAndZKI(txs, event.ZKI)
				if err != nil {
					continue
				}
				// send txs and proof to the eth contract
				batch := state.Batch{Txs: txs}
				if _, err = agr.EtherMan.ConsolidateBatch(batch, *proof); err != nil {
					continue
				}
			}
		}
	}()
}

func (agr *Aggregator) isProfitable(txs []types.Transaction) bool {
	// get strategy from the config and check
	return true
}

func (agr *Aggregator) Stop() {
	agr.cancel()
}
