package priceprovider

import "time"

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

// Type for the price provider
type Type string

const (
	// UniswapType for uniswap price provider
	UniswapType Type = "uniswap"
)

// Config represents the configuration of the pricegetter
type Config struct {
	// URL is Ethereum network url, if type is uniswap
	URL string `mapstructure:"URL"`

	// Type is price getter type
	Type Type `mapstructure:"Type"`
}
