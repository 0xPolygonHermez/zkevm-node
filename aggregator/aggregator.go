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
	ethMan etherman.EtherMan,
	syncClient SynchronizerClient,
	prover ProverClient,
) (Aggregator, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return Aggregator{
		State:          state,
		BatchProcessor: bp,
		EtherMan:       ethMan,
		Synchronizer:   syncClient,
		Prover:         prover,

		ctx:    ctx,
		cancel: cancel,
	}, nil
}

type txsWithProof struct {
	txs   []*types.Transaction
	proof *state.Proof
}

func (agr *Aggregator) Start() {
	// reads from batchesChan
	// get txs from state by batchNum
	// check if it's profitable or not
	// send zki + txs to the prover
	// send proof + txs to the SC
	go func() {
		txsByBatchNum := make(map[uint64]txsWithProof)
		for {
			select {
			case <-agr.ctx.Done():
				return
			case event := <-agr.Synchronizer.VirtualBatchEventChan:
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
				proof, err := agr.Prover.SendTxs(txs)
				if err != nil {
					continue
				}
				txsByBatchNum[event.BatchNum] = txsWithProof{txs: txs, proof: proof}
			case event := <-agr.Synchronizer.ConsolidatedBatchEventChan:
				previousBatchNum := event.BatchNum - 1
				// send txs and proof to the eth contract
				txsWithProof, ok := txsByBatchNum[previousBatchNum]
				if ok {

					batch := state.Batch{Transactions: txsWithProof.txs}
					if _, err := agr.EtherMan.ConsolidateBatch(batch, *txsWithProof.proof); err != nil {
						continue
					}
					delete(txsByBatchNum, previousBatchNum)
				}
			}
		}
	}()
}

func (agr *Aggregator) isProfitable(txs []*types.Transaction) bool {
	// get strategy from the config and check
	return true
}

func (agr *Aggregator) Stop() {
	agr.cancel()
}
