package pricegetter

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pricegetter/priceprovider"
	"github.com/hermeznetwork/hermez-core/pricegetter/priceprovider/uniswap"
)

// Client for the pricegetter
type Client interface {
	// SyncPrice sync price in endless for loop, used only with Async mode
	SyncPrice(ctx context.Context)
	// GetPrice getting eth to matic price
	GetPrice(ctx context.Context) (*big.Float, error)
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

// GetPrice getting default price
func (c *defaultClient) GetPrice(ctx context.Context) (*big.Float, error) {
	return c.defaultPrice, nil
}

// SyncPrice not needed for the default type
func (c *defaultClient) SyncPrice(ctx context.Context) {}

// syncClient using synchronous request
type syncClient struct {
	priceProvider priceprovider.PriceProvider
}

// GetPrice getting price from the price provider
func (c *syncClient) GetPrice(ctx context.Context) (*big.Float, error) {
	return c.priceProvider.GetPrice(ctx)
}

// SyncPrice not used with sync type
func (c *syncClient) SyncPrice(ctx context.Context) {}

// asyncClient
type asyncClient struct {
	cfg           Config
	priceProvider priceprovider.PriceProvider
	price         *big.Float
	lastUpdated   time.Time
}

// SyncPrice syncing price with the price provider every n second
func (c *asyncClient) SyncPrice(ctx context.Context) {
	ticker := time.NewTicker(c.cfg.UpdateFrequency.Duration)
	defer ticker.Stop()
	var err error
	for {
		c.price, err = c.priceProvider.GetPrice(ctx)
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

// GetPrice get price, that is syncing every n second
func (c *asyncClient) GetPrice(ctx context.Context) (*big.Float, error) {
	return c.price, nil
}
