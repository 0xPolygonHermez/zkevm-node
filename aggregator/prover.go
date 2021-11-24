package aggregator

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

// ProverClient is used to interact with the prover
type ProverClient struct {
}

// NewProverClient creates a new prover client
func NewProverClient() ProverClient {
	return ProverClient{}
}

// SendTxs sends txs to the prover
func (p *ProverClient) SendTxs(txs []*types.Transaction) (*state.Proof, error) {
	return nil, nil
}
