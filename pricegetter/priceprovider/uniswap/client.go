package uniswap

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	x96                        = new(big.Float).SetInt(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(96), nil)) // nolint:gomnd
	uniswapAddressEthMaticPool = common.HexToAddress("0x290a6a7460b308ee3f19023d2d00de604bcf5b42")
)

// PriceProvider price proved which takes price from the eth/matic pool from uniswap
type PriceProvider struct {
	Uni *Uniswap
}

// NewPriceProvider init uniswap price provider
func NewPriceProvider(URL string) (*PriceProvider, error) {
	ethClient, err := ethclient.Dial(URL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", URL, err)
		return nil, err
	}

	uni, err := NewUniswap(uniswapAddressEthMaticPool, ethClient)
	if err != nil {
		return nil, err
	}
	return &PriceProvider{Uni: uni}, nil
}

// SqrtPriceX96ToPrice convert uniswap v3 sqrt price in x96 format to big.Float
// calculation taken from here - https://docs.uniswap.org/sdk/guides/fetching-prices#understanding-sqrtprice
func sqrtPriceX96ToPrice(sqrtPriceX96 *big.Int) (price *big.Float) {
	d := big.NewFloat(0).Quo(new(big.Float).SetInt(sqrtPriceX96), x96)
	p := big.NewFloat(0).Mul(d, d)

	price = big.NewFloat(0).Quo(big.NewFloat(1), p)
	return
}

// GetEthToMaticPrice price from the uniswap pool contract
func (c *PriceProvider) GetEthToMaticPrice(ctx context.Context) (*big.Float, error) {
	slot, err := c.Uni.Slot0(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("failed to get matic price from uniswap: %v", err)
	}

	return sqrtPriceX96ToPrice(slot.SqrtPriceX96), nil
}
