package aggregator

import (
	"context"

	"github.com/hermeznetwork/hermez-core/proverclient/pb"
)

// Prime field. It is the prime number used as the order in our elliptic curve
const fr = "21888242871839275222246405745257275088548364400416034343698204186575808495617"

// Aggregator represents an aggregator
type Aggregator struct {
	cfg Config

	State          stateInterface
	EtherMan       etherman
	ZkProverClient pb.ZKProverServiceClient

	ProfitabilityChecker aggregatorTxProfitabilityChecker

	ctx    context.Context
	cancel context.CancelFunc
}

// NewAggregator creates a new aggregator
func NewAggregator(
	cfg Config,
	state stateInterface,
	ethMan etherman,
	zkProverClient pb.ZKProverServiceClient,
) (Aggregator, error) {
	ctx, cancel := context.WithCancel(context.Background())

	var profitabilityChecker aggregatorTxProfitabilityChecker
	switch cfg.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = NewTxProfitabilityCheckerBase(state, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration, cfg.TxProfitabilityMinReward.Int)
	case ProfitabilityAcceptAll:
		profitabilityChecker = NewTxProfitabilityCheckerAcceptAll(state, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration)
	}
	a := Aggregator{
		cfg: cfg,

		State:                state,
		EtherMan:             ethMan,
		ZkProverClient:       zkProverClient,
		ProfitabilityChecker: profitabilityChecker,

		ctx:    ctx,
		cancel: cancel,
	}

	return a, nil
}

// Start starts the aggregator
func (a *Aggregator) Start() {
	// TODO: the old aggregator was consolidating batches, but now the concept is different, so we need to review it
}

// Stop stops the aggregator
func (a *Aggregator) Stop() {
	a.cancel()
}
