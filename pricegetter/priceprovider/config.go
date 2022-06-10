package priceprovider

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
