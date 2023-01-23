package gasprice

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// L2GasPricer interface for gas price suggestor.
type L2GasPricer interface {
	UpdateGasPriceAvg()
}

// NewL2GasPriceSuggestor init.
func NewL2GasPriceSuggestor(ctx context.Context, cfg Config, pool pool, ethMan *etherman.Client, state *state.State) {
	var gpricer L2GasPricer
	switch cfg.Type {
	case LastNBatchesType:
		log.Info("Lastnbatches type selected")
		gpricer = newLastNL2BlocksGasPriceSuggestor(ctx, cfg, state, pool)
	case FollowerType:
		log.Info("Follower type selected")
		gpricer = newFollowerGasPriceSuggestor(ctx, cfg, pool, ethMan)
	case DefaultType:
		log.Info("Default type selected")
		gpricer = newDefaultGasPriceSuggestor(ctx, cfg, pool)
	default:
		log.Fatal("unknown l2 gas price suggestor type ", cfg.Type, ". Please specify a valid one: 'lastnbatches', 'follower' or 'default'")
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("Finishing l2 gas price suggestor...")
			return
		case <-time.After(cfg.UpdatePeriod.Duration):
			gpricer.UpdateGasPriceAvg()
		}
	}
}
