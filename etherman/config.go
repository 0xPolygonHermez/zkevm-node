package etherman

import "github.com/ethereum/go-ethereum/common"

// Address is a wrapper type that parses time duration from text.
type Address struct {
	common.Address `validate:"required"`
}

// UnmarshalText unmarshalls address from text.
func (d *Address) UnmarshalText(data []byte) error {
	addr := common.HexToAddress(string(data))
	d.Address = addr
	return nil
}

// Config represents the configuration of the etherman
type Config struct {
	URL        string  `mapstructure:"URL"`
	PoEAddress Address `mapstructure:"PoEAddress"`

	PrivateKeyPath     string `mapstructure:"PrivateKeyPath"`
	PrivateKeyPassword string `mapstructure:"PrivateKeyPassword"`
}
