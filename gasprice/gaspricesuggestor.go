package gasprice

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// L2GasPriceSuggestorInterface interface for gas price suggestor.
type L2GasPriceSuggestorInterface interface {
	UpdateGasPriceAvg()
}

// GasPriceSuggestor struct for gas price suggestor.
type GasPriceSuggestor struct {
	cfg Config
	ctx context.Context

	methods L2GasPriceSuggestorInterface
}

// NewL2GasPriceSuggestor init.
func NewL2GasPriceSuggestor(ctx context.Context, cfg Config, pool pool, ethMan *etherman.Client, state *state.State) *GasPriceSuggestor {
	var gpricer L2GasPriceSuggestorInterface
	switch cfg.Type {
	case LastNBatchesType:
		log.Info("Lastnbatches type selected")
		gpricer = newSuggestorLastNL2Blocks(ctx, cfg, state, pool)
	case FollowerType:
		log.Info("Follower type selected")
		gpricer = newFollowerGasPriceSuggestor(ctx, cfg, pool, ethMan)
	case DefaultType:
		log.Info("Default type selected")
		gpricer = newDefaultSuggestor(ctx, cfg, pool)
	default:
		log.Fatal("unknown l2 gas price suggestor type ", cfg.Type, ". Please specify a valid one: 'lastnbatches', 'follower' or 'default'")
	}
	gps := &GasPriceSuggestor{
		cfg:     cfg,
		ctx:     ctx,
		methods: gpricer,
	}
	return gps
}

// Start function runs the GasPriceSuggestor
func (g GasPriceSuggestor) Start() error {
	for {
		select {
		case <-g.ctx.Done():
			return nil
		case <-time.After(g.cfg.UpdatePeriod.Duration):
			g.methods.UpdateGasPriceAvg()
		}
	}
}
