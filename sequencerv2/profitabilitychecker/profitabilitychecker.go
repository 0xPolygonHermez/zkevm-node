package profitabilitychecker

import (
	"context"
	"math/big"

	"github.com/hermeznetwork/hermez-core/ethermanv2/types"
	"github.com/hermeznetwork/hermez-core/pricegetter"
)

// Checker checks profitability to send sequences
type Checker struct {
	Config      Config
	EthMan      etherman
	PriceGetter pricegetter.Client
}

// New creates new profitability checker
func New(
	cfg Config,
	etherMan etherman,
	priceGetter priceGetter) *Checker {
	return &Checker{
		Config:      cfg,
		EthMan:      etherMan,
		PriceGetter: priceGetter,
	}
}

// IsSequenceProfitable check if sequence is profitable by comparing L1 tx gas cost and collateral with fee rewards
func (c *Checker) IsSequenceProfitable(ctx context.Context, sequence types.Sequence) (bool, error) {
	if c.Config.SendBatchesEvenWhenNotProfitable {
		return true, nil
	}
	// fee - it's collateral for batch, get from SC in matic
	fee, err := c.EthMan.GetSendSequenceFee()
	if err != nil {
		return false, err
	}

	// this reward is in ethereum wei
	reward := big.NewInt(0)
	for _, tx := range sequence.Txs {
		reward.Add(reward, new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(tx.Gas())))
	}

	// get price of matic (1 eth = x matic)
	price, err := c.PriceGetter.GetEthToMaticPrice(ctx)
	if err != nil {
		return false, err
	}

	// convert reward in eth to reward in matic
	priceInt := new(big.Int)
	price.Int(priceInt)
	reward.Mul(reward, priceInt)

	if reward.Cmp(fee) < 0 {
		return false, nil
	}

	return true, nil
}

// IsSendSequencesProfitable checks profitability to send sequences to the ethereum
func (c *Checker) IsSendSequencesProfitable(estimatedGas uint64, sequences []types.Sequence) bool {
	if c.Config.SendBatchesEvenWhenNotProfitable {
		return true
	}

	gasCostSequences := big.NewInt(0)
	for _, seq := range sequences {
		for _, tx := range seq.Txs {
			gasCostSequences.Add(gasCostSequences, new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(tx.Gas())))
			if gasCostSequences.Cmp(big.NewInt(0).SetUint64(estimatedGas)) > 0 {
				return true
			}
		}
	}

	return false
}
