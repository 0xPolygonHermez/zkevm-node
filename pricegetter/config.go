package pricegetter

import (
	"fmt"
	"math/big"
	"time"

	"github.com/hermeznetwork/hermez-core/pricegetter/priceprovider"
)

// Duration is a wrapper type that parses time duration from text.
type Duration struct {
	time.Duration `validate:"required"`
}

// UnmarshalText unmarshalls time duration from text.
func (d *Duration) UnmarshalText(data []byte) error {
	duration, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	d.Duration = duration
	return nil
}

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
	UpdateFrequency Duration `mapstructure:"UpdateFrequency"`

	// DefaultPrice is used only for the default type
	DefaultPrice TokenPrice `mapstructure:"DefaultPrice"`
}
