package sequencer

import (
	"math"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TxTracker is a struct that contains all the tx data needed to be managed by the worker
type TxTracker struct {
	Hash           common.Hash
	HashStr        string
	From           common.Address
	FromStr        string
	Nonce          uint64
	Gas            uint64 // To check if it fits into a batch
	GasPrice       *big.Int
	Cost           *big.Int       // Cost = Amount + Benefit
	Benefit        *big.Int       // GasLimit * GasPrice
	IsClaim        bool           // Needed to calculate efficiency
	BatchResources batchResources // To check if it fits into a batch
	Efficiency     float64
	RawTx          []byte
	ReceivedAt     time.Time // To check if it has been in the efficiency list for too long
	IP             string    // IP of the tx sender
}

// newTxTracker creates and inits a TxTracker
func newTxTracker(tx types.Transaction, isClaim bool, counters state.ZKCounters, constraints batchConstraints, weights batchResourceWeights, ip string) (*TxTracker, error) {
	addr, err := state.GetSender(tx)
	if err != nil {
		return nil, err
	}

	txTracker := &TxTracker{
		Hash:       tx.Hash(),
		From:       addr,
		Nonce:      tx.Nonce(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		Cost:       tx.Cost(),
		ReceivedAt: time.Now(),
		IP:         ip,
	}

	txTracker.IsClaim = isClaim
	txTracker.BatchResources.zKCounters = counters
	txTracker.BatchResources.bytes = tx.Size()
	txTracker.HashStr = txTracker.Hash.String()
	txTracker.FromStr = txTracker.From.String()
	txTracker.Benefit = new(big.Int).Mul(new(big.Int).SetUint64(txTracker.Gas), txTracker.GasPrice)
	txTracker.calculateEfficiency(constraints, weights)
	txTracker.RawTx, err = state.EncodeTransactions([]types.Transaction{tx})
	if err != nil {
		return nil, err
	}

	return txTracker, nil
}

// updateZKCounters updates the counters of the tx and recalculates the tx efficiency
func (tx *TxTracker) updateZKCounters(counters state.ZKCounters, constraints batchConstraints, weights batchResourceWeights) {
	tx.BatchResources.zKCounters = counters
	tx.calculateEfficiency(constraints, weights)
}

// calculateEfficiency calculates the tx efficiency
func (tx *TxTracker) calculateEfficiency(constraints batchConstraints, weights batchResourceWeights) {
	const perThousand = 1000 // TODO: Add this as config parameter

	totalWeight := float64(weights.WeightArithmetics + weights.WeightBatchBytesSize + weights.WeightBinaries + weights.WeightCumulativeGasUsed +
		weights.WeightKeccakHashes + weights.WeightMemAligns + weights.WeightPoseidonHashes + weights.WeightPoseidonPaddings + weights.WeightSteps)

	// TODO: Optmize tx.Efficiency calculation (precalculate constansts values)
	// TODO: Evaluate avoid type conversion (performance impact?)
	resourceCost := (float64(tx.BatchResources.zKCounters.CumulativeGasUsed)/float64(constraints.MaxCumulativeGasUsed))*float64(weights.WeightCumulativeGasUsed)/totalWeight +
		(float64(tx.BatchResources.zKCounters.UsedArithmetics)/float64(constraints.MaxArithmetics))*float64(weights.WeightArithmetics)/totalWeight +
		(float64(tx.BatchResources.zKCounters.UsedBinaries)/float64(constraints.MaxBinaries))*float64(weights.WeightBinaries)/totalWeight +
		(float64(tx.BatchResources.zKCounters.UsedKeccakHashes)/float64(constraints.MaxKeccakHashes))*float64(weights.WeightKeccakHashes)/totalWeight +
		(float64(tx.BatchResources.zKCounters.UsedMemAligns)/float64(constraints.MaxMemAligns))*float64(weights.WeightMemAligns)/totalWeight +
		(float64(tx.BatchResources.zKCounters.UsedPoseidonHashes)/float64(constraints.MaxPoseidonHashes))*float64(weights.WeightPoseidonHashes)/totalWeight +
		(float64(tx.BatchResources.zKCounters.UsedPoseidonPaddings)/float64(constraints.MaxPoseidonPaddings))*float64(weights.WeightPoseidonPaddings)/totalWeight +
		(float64(tx.BatchResources.zKCounters.UsedSteps)/float64(constraints.MaxSteps))*float64(weights.WeightSteps)/totalWeight +
		(float64(tx.BatchResources.bytes)/float64(constraints.MaxBatchBytesSize))*float64(weights.WeightBatchBytesSize)/totalWeight //Meto config

	resourceCost = resourceCost * perThousand
	var eff *big.Float
	if tx.IsClaim {
		eff = big.NewFloat(math.MaxFloat64)
	} else {
		ben := big.NewFloat(0).SetInt(tx.Benefit)
		rc := big.NewFloat(0).SetFloat64(resourceCost)
		eff = big.NewFloat(0).Quo(ben, rc)
	}

	var accuracy big.Accuracy
	tx.Efficiency, accuracy = eff.Float64()
	if accuracy != big.Exact {
		log.Errorf("CalculateEfficiency accuracy warning (%s). Calculated=%s Assigned=%f", accuracy.String(), eff.String(), tx.Efficiency)
	}
}
