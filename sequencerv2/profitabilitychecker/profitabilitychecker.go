package profitabilitychecker

import (
	"context"
	"math/big"

	"github.com/hermeznetwork/hermez-core/ethermanv2/types"
	"github.com/hermeznetwork/hermez-core/pricegetter"
)

// Checker checks profitability to send sequences
type Checker struct {
	EthMan      etherman
	PriceGetter pricegetter.Client
}

// New creates new profitability checker
func New(
	etherMan etherman,
	priceGetter priceGetter) *Checker {
	return &Checker{
		EthMan:      etherMan,
		PriceGetter: priceGetter,
	}
}

// IsSequenceProfitable check is sequence profitable by comparing tx gas cost with fee
func (c *Checker) IsSequenceProfitable(ctx context.Context, sequence types.Sequence) (bool, error) {
	fee, err := c.EthMan.GetFee()
	if err != nil {
		return false, err
	}

	reward := big.NewInt(0)
	for _, tx := range sequence.Txs {
		reward.Add(reward, tx.Cost())
	}

	price, err := c.PriceGetter.GetPrice(ctx)
	if err != nil {
		return false, err
	}

	priceInt := new(big.Int)
	price.Int(priceInt)
	reward.Mul(reward, priceInt)

	if reward.Cmp(fee) < 0 {
		return false, nil
	}

	return true, nil
}

// IsSendSequencesProfitable checks profitability to send sequences to the ethereum
func (c *Checker) IsSendSequencesProfitable(estimatedGas *big.Int, sequences []types.Sequence) bool {
	gasCostSequences := big.NewInt(0)
	for _, seq := range sequences {
		for _, tx := range seq.Txs {
			gasCostSequences.Add(gasCostSequences, tx.Cost())
			if gasCostSequences.Cmp(estimatedGas) > 0 {
				return true
			}
		}
	}

	return false
}
