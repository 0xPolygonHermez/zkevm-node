package pricegetter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter/priceprovider"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter/priceprovider/uniswap"
)

// Client for the pricegetter
type Client interface {
	// Start price getter client
	Start(ctx context.Context)
	// GetEthToMaticPrice getting eth to matic price
	GetEthToMaticPrice(ctx context.Context) (*big.Float, error)
}

// NewClient inits price getter client
func NewClient(cfg Config) (Client, error) {
	var (
		priceProvider priceprovider.PriceProvider
		err           error
	)
	switch cfg.PriceProvider.Type {
	case priceprovider.UniswapType:
		priceProvider, err = uniswap.NewPriceProvider(cfg.PriceProvider.URL)
		if err != nil {
			return nil, err
		}
	}
	switch cfg.Type {
	case SyncType:
		return &syncClient{priceProvider: priceProvider}, nil
	case AsyncType:
		return &asyncClient{
			cfg:           cfg,
			priceProvider: priceProvider,
			price:         nil,
			lastUpdated:   time.Time{},
		}, nil
	case DefaultType:
		return &defaultClient{defaultPrice: cfg.DefaultPrice.Float}, nil
	}
	return nil, fmt.Errorf("pricegetter type is not specified")
}

// defaultClient using default price set by config
type defaultClient struct {
	defaultPrice *big.Float
}

// GetEthToMaticPrice getting default price
func (c *defaultClient) GetEthToMaticPrice(ctx context.Context) (*big.Float, error) {
	return c.defaultPrice, nil
}

// Start function for default client
func (c *defaultClient) Start(ctx context.Context) {}

// syncClient using synchronous request
type syncClient struct {
	priceProvider priceprovider.PriceProvider
}

// GetEthToMaticPrice getting price from the price provider
func (c *syncClient) GetEthToMaticPrice(ctx context.Context) (*big.Float, error) {
	return c.priceProvider.GetEthToMaticPrice(ctx)
}

// Start starting sync client
func (c *syncClient) Start(ctx context.Context) {}

// asyncClient
type asyncClient struct {
	cfg           Config
	priceProvider priceprovider.PriceProvider
	price         *big.Float
	lastUpdated   time.Time
}

// SyncPrice syncing price with the price provider every n second
func (c *asyncClient) syncPrice(ctx context.Context) {
	ticker := time.NewTicker(c.cfg.UpdateFrequency.Duration)
	defer ticker.Stop()
	var err error
	for {
		c.price, err = c.priceProvider.GetEthToMaticPrice(ctx)
		if err != nil {
			log.Errorf("failed to get price matic price, err: %v", err)
		} else {
			c.lastUpdated = time.Now()
		}
		select {
		case <-ticker.C:
			// nothing
		case <-ctx.Done():
			return
		}
	}
}

// GetEthToMaticPrice get price, that is syncing every n second
func (c *asyncClient) GetEthToMaticPrice(ctx context.Context) (*big.Float, error) {
	return c.price, nil
}

func (c *asyncClient) Start(ctx context.Context) {
	go c.syncPrice(ctx)
}
