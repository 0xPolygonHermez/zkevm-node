package pricegetter

import (
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/pricegetter/priceprovider"
)

// TokenPrice is a wrapper type that parses token amount to big float
type TokenPrice struct {
	*big.Float `validate:"required"`
}

// UnmarshalText unmarshal token amount from float string to big int
func (t *TokenPrice) UnmarshalText(data []byte) error {
	amount, ok := new(big.Float).SetString(string(data))
	if !ok {
		return fmt.Errorf("failed to unmarshal string to float")
	}
	t.Float = amount

	return nil
}

// Type for the pricegetter
type Type string

const (
	// SyncType synchronous request to price provider
	SyncType Type = "sync"
	// AsyncType update price every n second
	AsyncType Type = "async"
	// DefaultType get default price from the config
	DefaultType Type = "default"
)

// Config represents the configuration of the pricegetter
type Config struct {
	// Type is price getter type
	Type Type `mapstructure:"Type"`

	// PriceProvider config
	PriceProvider priceprovider.Config `mapstructure:"PriceProvider"`

	// UpdateFrequency is price updating frequency, used only for the async type
	UpdateFrequency types.Duration `mapstructure:"UpdateFrequency"`

	// DefaultPrice is used only for the default type
	DefaultPrice TokenPrice `mapstructure:"DefaultPrice"`
}
