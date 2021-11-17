package aggregator

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

type ProverClient struct {
}

func NewProverClient() ProverClient {
	return ProverClient{}
}

func (p *ProverClient) SendTxsAndZKI(txs []types.Transaction, zki *big.Int) (*state.Proof, error) {
	return nil, nil
}
