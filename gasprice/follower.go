package gasprice

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

// FollowerGasPrice struct.
type FollowerGasPrice struct {
	cfg  Config
	pool pool
	ctx  context.Context
	eth  ethermanInterface
}

// newFollowerGasPriceSuggester inits l2 follower gas price suggester which is based on the l1 gas price.
func newFollowerGasPriceSuggester(ctx context.Context, cfg Config, pool pool, ethMan ethermanInterface) *FollowerGasPrice {
	gps := &FollowerGasPrice{
		cfg:  cfg,
		pool: pool,
		ctx:  ctx,
		eth:  ethMan,
	}
	gps.UpdateGasPriceAvg()
	return gps
}

// UpdateGasPriceAvg updates the gas price.
func (f *FollowerGasPrice) UpdateGasPriceAvg() {
	ctx := context.Background()
	// Get L1 gasprice
	gp := f.eth.GetL1GasPrice(f.ctx)
	if big.NewInt(0).Cmp(gp) == 0 {
		log.Warn("gas price 0 received. Skipping update...")
		return
	}
	// Apply factor to calculate l2 gasPrice
	factor := big.NewFloat(0).SetFloat64(f.cfg.Factor)
	res := new(big.Float).Mul(factor, big.NewFloat(0).SetInt(gp))

	// Store l2 gasPrice calculated
	result := new(big.Int)
	res.Int(result)
	minGasPrice := big.NewInt(0).SetUint64(f.cfg.DefaultGasPriceWei)
	if minGasPrice.Cmp(result) == 1 { // minGasPrice > result
		log.Warn("setting minGasPrice for L2")
		result = minGasPrice
	}
	var truncateValue *big.Int
	log.Debug("Full L2 gas price value: ", result, ". Length: ", len(result.String()))
	numLength := len(result.String())
	if numLength > 3 { //nolint:gomnd
		aux := "%0" + strconv.Itoa(numLength-3) + "d" //nolint:gomnd
		var ok bool
		value := result.String()[:3] + fmt.Sprintf(aux, 0)
		truncateValue, ok = new(big.Int).SetString(value, encoding.Base10)
		if !ok {
			log.Error("error converting: ", truncateValue)
		}
	} else {
		truncateValue = result
	}
	log.Debug("Storing truncated L2 gas price: ", truncateValue)
	if truncateValue != nil {
		err := f.pool.SetGasPrice(ctx, truncateValue.Uint64())
		if err != nil {
			log.Errorf("failed to update gas price in poolDB, err: %v", err)
		}
	} else {
		log.Error("nil value detected. Skipping...")
	}
}
