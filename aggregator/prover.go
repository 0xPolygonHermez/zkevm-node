package aggregator

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

type ProverClient struct {
}

func NewProverClient() ProverClient {
	return ProverClient{}
}

func (p *ProverClient) SendTxs(txs []*types.LegacyTx) (*state.Proof, error) {
	return nil, nil
}
