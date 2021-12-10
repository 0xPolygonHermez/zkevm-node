package aggregator

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

// Aggregator represents an aggregator
type Aggregator struct {
	cfg Config

	State          state.State
	BatchProcessor state.BatchProcessor
	EtherMan       etherman.EtherMan
	ZkProverClient proverclient.ZKProverClient

	ctx    context.Context
	cancel context.CancelFunc
}

// NewAggregator creates a new aggregator
func NewAggregator(
	cfg Config,
	state state.State,
	ethMan etherman.EtherMan,
	zkProverClient proverclient.ZKProverClient,
) (Aggregator, error) {
	ctx, cancel := context.WithCancel(context.Background())
	a := Aggregator{
		cfg: cfg,

		State:          state,
		EtherMan:       ethMan,
		ZkProverClient: zkProverClient,

		ctx:    ctx,
		cancel: cancel,
	}

	return a, nil
}

// Start starts the aggregator
func (a *Aggregator) Start() {
	// init connection to the prover
	var opts []grpc.CallOption
	getProofClient, err := a.ZkProverClient.GenProof(a.ctx, opts...)
	if err != nil {
		log.Errorf("failed to connect to the prover, err: %v", err)
		return
	}

	// this is a batches, that were sent to ethereum to consolidate
	batchesSent := make(map[uint64]bool)

	for {
		select {
		case <-time.After(a.cfg.IntervalToConsolidateState):
			// 1. check, if state is synced
			lastSyncedBatchNum, err := a.State.GetLastBatchNumber(a.ctx)
			if err != nil {
				log.Warnf("failed to get last synced batch, err: %v", err)
				continue
			}
			lastEthBatchNum, err := a.State.GetLastBatchNumberSeenOnEthereum(a.ctx)
			if err != nil {
				log.Warnf("failed to get last eth batch, err: %v", err)
				continue
			}
			if lastSyncedBatchNum < lastEthBatchNum {
				log.Infow("waiting for the state to be synced, lastSyncedBatchNum: %d, lastEthBatchNum: %d", lastSyncedBatchNum, lastEthBatchNum)
				continue
			}

			// 2. find next batch to consolidate
			lastConsolidatedBatch, err := a.State.GetLastBatch(a.ctx, false)
			if err != nil {
				log.Warnf("failed to get last consolidated batch, err: %v", err)
				continue
			}
			delete(batchesSent, lastConsolidatedBatch.BatchNumber)

			batchToConsolidate, err := a.State.GetBatchByNumber(a.ctx, lastConsolidatedBatch.BatchNumber+1)

			if err != nil {
				if err == pgx.ErrNoRows {
					log.Infof("there is no batches to consolidate")
					continue
				}
				log.Warnf("failed to get batch to consolidate, err: %v", err)
				continue
			}

			if batchesSent[batchToConsolidate.BatchNumber] {
				log.Infof("batch with number %d was already sent, but not yet consolidated by synchronizer",
					batchToConsolidate.BatchNumber)
				continue
			}

			// 3. check if it's profitable or not
			// check is it profitable to aggregate txs or not
			if !a.isProfitable(batchToConsolidate.Transactions) {
				log.Info("Batch %d is not profitable", batchToConsolidate.BatchNumber)
				continue
			}

			// 4. send zki + txs to the prover
			stateRootConsolidated, err := a.State.GetStateRootByBatchNumber(lastConsolidatedBatch.BatchNumber)
			if err != nil {
				log.Warnf("failed to get current state root, err: %v", err)
				continue
			}

			stateRootToConsolidate, err := a.State.GetStateRootByBatchNumber(batchToConsolidate.BatchNumber)
			if err != nil {
				log.Warnf("failed to get state root to consolidate, err: %v", err)
				continue
			}

			// TODO: change this, once we have exit root
			fakeLastGlobalExitRoot, _ := new(big.Int).SetString("1234123412341234123412341234123412341234123412341234123412341234", 16)
			chainID := uint64(1337) //nolint:gomnd
			batch := &proverclient.Batch{
				Message:            "calculate",
				CurrentStateRoot:   stateRootConsolidated,
				NewStateRoot:       stateRootToConsolidate,
				L2Txs:              batchToConsolidate.RawTxsData,
				LastGlobalExitRoot: fakeLastGlobalExitRoot.Bytes(),
				SequencerAddress:   batchToConsolidate.Sequencer.String(),
				// TODO: consider to put chain id to batch, so there is no need to request block
				ChainId: chainID,
			}

			err = getProofClient.Send(batch)
			if err != nil {
				log.Warnf("failed to send batch to the prover, batchNumber: %v, err: %v", batchToConsolidate.BatchNumber, err)
				continue
			}
			proof, err := getProofClient.Recv()
			if err != nil {
				log.Warnf("failed to get proof from the prover, batchNumber: %v, err: %v", batchToConsolidate.BatchNumber, err)
				continue
			}
			// 4. send proof + txs to the SC
			batchNum := new(big.Int).SetUint64(batchToConsolidate.BatchNumber)
			h, err := a.EtherMan.ConsolidateBatch(batchNum, proof)
			if err != nil {
				log.Warnf("failed to send request to consolidate batch to ethereum, batch number: %d, err: %v",
					batchToConsolidate.BatchNumber, err)
				continue
			}
			batchesSent[batchToConsolidate.BatchNumber] = true

			log.Infof("Batch %d consolidated: %s", batchToConsolidate.BatchNumber, h.Hash())

		case <-a.ctx.Done():
			return
		}
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
