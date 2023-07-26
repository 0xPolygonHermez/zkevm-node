package sequencer

import (
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TxTracker is a struct that contains all the tx data needed to be managed by the worker
type TxTracker struct {
	Hash                              common.Hash
	HashStr                           string
	From                              common.Address
	FromStr                           string
	Nonce                             uint64
	Gas                               uint64 // To check if it fits into a batch
	GasPrice                          *big.Int
	Cost                              *big.Int             // Cost = Amount + Benefit
	Benefit                           *big.Int             // GasLimit * GasPrice
	BatchResources                    state.BatchResources // To check if it fits into a batch
	Efficiency                        float64
	RawTx                             []byte
	ReceivedAt                        time.Time // To check if it has been in the efficiency list for too long
	IP                                string    // IP of the tx sender
	FailedReason                      *string   // FailedReason is the reason why the tx failed, if it failed
	Constraints                       batchConstraintsFloat64
	WeightMultipliers                 batchResourceWeightMultipliers
	ResourceCostMultiplier            float64
	TotalWeight                       float64
	BreakEvenGasPrice                 *big.Int
	GasPriceEffectivePercentage       uint8
	EffectiveGasPriceProcessCount     uint8
	IsEffectiveGasPriceFinalExecution bool
	L1GasPrice                        uint64
}

// batchResourceWeightMultipliers is a struct that contains the weight multipliers for each resource
type batchResourceWeightMultipliers struct {
	cumulativeGasUsed float64
	arithmetics       float64
	binaries          float64
	keccakHashes      float64
	memAligns         float64
	poseidonHashes    float64
	poseidonPaddings  float64
	steps             float64
	batchBytesSize    float64
}

// batchConstraints represents the Constraints for a batch in float64
type batchConstraintsFloat64 struct {
	maxTxsPerBatch       float64
	maxBatchBytesSize    float64
	maxCumulativeGasUsed float64
	maxKeccakHashes      float64
	maxPoseidonHashes    float64
	maxPoseidonPaddings  float64
	maxMemAligns         float64
	maxArithmetics       float64
	maxBinaries          float64
	maxSteps             float64
}

// newTxTracker creates and inti a TxTracker
func newTxTracker(tx types.Transaction, counters state.ZKCounters, constraints batchConstraintsFloat64, weights batchResourceWeights, resourceCostMultiplier float64, ip string) (*TxTracker, error) {
	addr, err := state.GetSender(tx)
	if err != nil {
		return nil, err
	}

	totalWeight := float64(weights.WeightArithmetics + weights.WeightBatchBytesSize + weights.WeightBinaries + weights.WeightCumulativeGasUsed +
		weights.WeightKeccakHashes + weights.WeightMemAligns + weights.WeightPoseidonHashes + weights.WeightPoseidonPaddings + weights.WeightSteps)

	rawTx, err := state.EncodeTransactionWithoutEffectivePercentage(tx)
	if err != nil {
		return nil, err
	}
	txTracker := &TxTracker{
		Hash:     tx.Hash(),
		HashStr:  tx.Hash().String(),
		From:     addr,
		FromStr:  addr.String(),
		Nonce:    tx.Nonce(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Cost:     tx.Cost(),
		Benefit:  new(big.Int).Mul(new(big.Int).SetUint64(tx.Gas()), tx.GasPrice()),
		BatchResources: state.BatchResources{
			Bytes:      tx.Size(),
			ZKCounters: counters,
		},
		Efficiency:                        0,
		RawTx:                             rawTx,
		ReceivedAt:                        time.Now(),
		IP:                                ip,
		Constraints:                       constraints,
		WeightMultipliers:                 calculateWeightMultipliers(weights, totalWeight),
		ResourceCostMultiplier:            resourceCostMultiplier,
		TotalWeight:                       totalWeight,
		BreakEvenGasPrice:                 new(big.Int).SetUint64(0),
		EffectiveGasPriceProcessCount:     0,
		IsEffectiveGasPriceFinalExecution: false,
	}
	txTracker.calculateEfficiency(constraints, weights)

	return txTracker, nil
}

// updateZKCounters updates the counters of the tx and recalculates the tx efficiency

func (tx *TxTracker) updateZKCounters(counters state.ZKCounters, constraints batchConstraintsFloat64, weights batchResourceWeights) {
	tx.BatchResources.ZKCounters = counters
	tx.calculateEfficiency(constraints, weights)
}

// calculateEfficiency calculates the tx efficiency
func (tx *TxTracker) calculateEfficiency(constraints batchConstraintsFloat64, weights batchResourceWeights) {
	totalWeight := float64(weights.WeightArithmetics + weights.WeightBatchBytesSize + weights.WeightBinaries + weights.WeightCumulativeGasUsed +
		weights.WeightKeccakHashes + weights.WeightMemAligns + weights.WeightPoseidonHashes + weights.WeightPoseidonPaddings + weights.WeightSteps)

	// TODO: Optmize tx.Efficiency calculation (precalculate constansts values)
	// TODO: Evaluate avoid type conversion (performance impact?)
	resourceCost := (float64(tx.BatchResources.ZKCounters.CumulativeGasUsed)/constraints.maxCumulativeGasUsed)*float64(weights.WeightCumulativeGasUsed)/totalWeight +
		(float64(tx.BatchResources.ZKCounters.UsedArithmetics)/constraints.maxArithmetics)*float64(weights.WeightArithmetics)/totalWeight +
		(float64(tx.BatchResources.ZKCounters.UsedBinaries)/constraints.maxBinaries)*float64(weights.WeightBinaries)/totalWeight +
		(float64(tx.BatchResources.ZKCounters.UsedKeccakHashes)/constraints.maxKeccakHashes)*float64(weights.WeightKeccakHashes)/totalWeight +
		(float64(tx.BatchResources.ZKCounters.UsedMemAligns)/constraints.maxMemAligns)*float64(weights.WeightMemAligns)/totalWeight +
		(float64(tx.BatchResources.ZKCounters.UsedPoseidonHashes)/constraints.maxPoseidonHashes)*float64(weights.WeightPoseidonHashes)/totalWeight +
		(float64(tx.BatchResources.ZKCounters.UsedPoseidonPaddings)/constraints.maxPoseidonPaddings)*float64(weights.WeightPoseidonPaddings)/totalWeight +
		(float64(tx.BatchResources.ZKCounters.UsedSteps)/constraints.maxSteps)*float64(weights.WeightSteps)/totalWeight +
		(float64(tx.BatchResources.Bytes)/constraints.maxBatchBytesSize)*float64(weights.WeightBatchBytesSize)/totalWeight //Meto config

	resourceCost = resourceCost * tx.ResourceCostMultiplier

	var eff *big.Float

	ben := big.NewFloat(0).SetInt(tx.Benefit)
	rc := big.NewFloat(0).SetFloat64(resourceCost)
	eff = big.NewFloat(0).Quo(ben, rc)

	var accuracy big.Accuracy
	tx.Efficiency, accuracy = eff.Float64()
	log.Infof("CalculateEfficiency(%f) for tx(%s)", tx.Efficiency, tx.Hash.String())
	if accuracy != big.Exact {
		log.Errorf("CalculateEfficiency accuracy warning (%s). Calculated=%s Assigned=%f", accuracy.String(), eff.String(), tx.Efficiency)
	}
}

// calculateWeightMultipliers calculates the weight multipliers for each resource
func calculateWeightMultipliers(weights batchResourceWeights, totalWeight float64) batchResourceWeightMultipliers {
	return batchResourceWeightMultipliers{
		cumulativeGasUsed: float64(weights.WeightCumulativeGasUsed) / totalWeight,
		arithmetics:       float64(weights.WeightArithmetics) / totalWeight,
		binaries:          float64(weights.WeightBinaries) / totalWeight,
		keccakHashes:      float64(weights.WeightKeccakHashes) / totalWeight,
		memAligns:         float64(weights.WeightMemAligns) / totalWeight,
		poseidonHashes:    float64(weights.WeightPoseidonHashes) / totalWeight,
		poseidonPaddings:  float64(weights.WeightPoseidonPaddings) / totalWeight,
		steps:             float64(weights.WeightSteps) / totalWeight,
		batchBytesSize:    float64(weights.WeightBatchBytesSize) / totalWeight,
	}
}
