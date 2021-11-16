package aggregator

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type ProverClient struct {
}

func NewProverClient() ProverClient {
	return ProverClient{}
}

func (p *ProverClient) SendTxsAndProof(txs []types.Transaction, zki *big.Int) (*big.Int, error) {
	return nil, nil
}