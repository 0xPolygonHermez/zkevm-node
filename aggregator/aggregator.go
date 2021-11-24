package aggregator

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/synchronizer"
)

type Aggregator struct {
	State          state.State
	BatchProcessor state.BatchProcessor
	EtherMan       etherman.EtherMan
	Synchronizer   synchronizer.Synchronizer
	Prover         ProverClient

	ctx           context.Context
	cancel        context.CancelFunc
	txsByBatchNum map[uint64]txsWithProof
}

func NewAggregator(
	cfg Config,
	state state.State,
	bp state.BatchProcessor,
	ethMan etherman.EtherMan,
	sy synchronizer.Synchronizer,
	prover ProverClient,
) (Aggregator, error) {
	ctx, cancel := context.WithCancel(context.Background())

	a := Aggregator{
		State:          state,
		BatchProcessor: bp,
		EtherMan:       ethMan,
		Synchronizer:   sy,
		Prover:         prover,

		ctx:           ctx,
		cancel:        cancel,
		txsByBatchNum: make(map[uint64]txsWithProof),
	}

	sy.RegisterNewBatchProposalHandler(a.onNewBatchPropostal)
	sy.RegisterNewConsolidatedStateHandler(a.onNewConsolidatedState)

	return a, nil
}

type txsWithProof struct {
	txs   []*types.Transaction
	proof *state.Proof
}

func (a *Aggregator) Start() {
	// reads from batchesChan
	// get txs from state by batchNum
	// check if it's profitable or not
	// send zki + txs to the prover
	// send proof + txs to the SC
	a.Synchronizer.Sync()
}

func (a *Aggregator) onNewBatchPropostal(batchNumber uint64, root common.Hash) {
	// get txs to send
	txs, err := a.State.GetTxsByBatchNum(batchNumber)
	if err != nil {
		return
	}
	// check is it profitable to aggregate txs or not
	if !a.isProfitable(txs) {
		return
	}
	// send txs and zki to the prover
	proof, err := a.Prover.SendTxs(txs)
	if err != nil {
		return
	}
	a.txsByBatchNum[batchNumber] = txsWithProof{txs: txs, proof: proof}
}

func (a *Aggregator) onNewConsolidatedState(batchNumber uint64, root common.Hash) {
	previousBatchNum := batchNumber - 1
	// send txs and proof to the eth contract
	txsWithProof, ok := a.txsByBatchNum[previousBatchNum]
	if ok {
		batch := state.Batch{Transactions: txsWithProof.txs}
		if _, err := a.EtherMan.ConsolidateBatch(batch, *txsWithProof.proof); err != nil {
			return
		}
		delete(a.txsByBatchNum, previousBatchNum)
	}
}

func (a *Aggregator) isProfitable(txs []*types.Transaction) bool {
	// get strategy from the config and check
	return true
}

func (a *Aggregator) Stop() {
	a.cancel()
}
