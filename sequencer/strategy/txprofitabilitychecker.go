package strategy

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman"
)

// TxProfitabilityCheckerType for different profitability checkers types
type TxProfitabilityCheckerType string

const (
	// ProfitabilityBase type that checks sum of costs of txs against min reward
	ProfitabilityBase = "base"
	// ProfitabilityAcceptAll validate batch anyway and don't check anything
	ProfitabilityAcceptAll = "acceptall"
)

// TxProfitabilityChecker interface for different profitability checkers
type TxProfitabilityChecker interface {
	IsProfitable(context.Context, []*types.Transaction) (bool, error)
}

// TxProfitabilityCheckerBase struct
type TxProfitabilityCheckerBase struct {
	EthMan etherman.EtherMan

	MinReward *big.Int
}

func NewTxProfitabilityCheckerBase(ethMan etherman.EtherMan, minReward *big.Int) TxProfitabilityChecker {
	return &TxProfitabilityCheckerBase{
		EthMan:    ethMan,
		MinReward: minReward,
	}
}

// IsProfitable checks for txs cost against the main reward
func (pc *TxProfitabilityCheckerBase) IsProfitable(ctx context.Context, txs []*types.Transaction) (bool, error) {
	txsGasCost := big.NewInt(0)
	for _, tx := range txs {
		txsGasCost.Add(txsGasCost, tx.Cost())
	}
	// sending amount of txs as matic reward there, bcs to calculate gas cost for tx this value is not important
	// checker can't get it before, bcs final matic amount is dependent on value from gas estimation
	gasCostForSendingBatch, err := pc.EthMan.EstimateGasForSendBatch(ctx, txs, big.NewInt(int64(len(txs))))
	if err != nil {
		return false, fmt.Errorf("failed to estimate gas for sending batch, err: %v", err)
	}

	reward := new(big.Int).Sub(txsGasCost, gasCostForSendingBatch)
	if reward.Cmp(big.NewInt(0)) < 0 {
		return false, nil
	}
	if reward.Cmp(pc.MinReward) < 0 {
		return false, nil
	}

	// TODO: take from config, how much matic tokens in percentage from whole reward sequencer is willing to pay to aggregator
	// to get his batch consolidated, then calculate, how much is reward minus minReward in matic tokens and multiply it
	// by percentage in config. This value we can also return from function

	return false, nil
}

// TxProfitabilityCheckerAcceptAll always returns true
type TxProfitabilityCheckerAcceptAll struct{}

// IsProfitable always returns true
func (pc *TxProfitabilityCheckerAcceptAll) IsProfitable(ctx context.Context, txs []*types.Transaction) (bool, error) {
	return true, nil
}
