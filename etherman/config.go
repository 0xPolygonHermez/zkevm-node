package etherman

import "github.com/ethereum/go-ethereum/common"

// Config represents the configuration of the etherman
type Config struct {
	URL        string         `env:"HERMEZCORE_ETHERMAN_URL"`
	PoEAddress common.Address `env:"HERMEZCORE_ETHERMAN_POEADDRESS"`

	PrivateKeyPath     string `env:"HERMEZCORE_ETHERMAN_PRIVATEKEY_PATH"`
	PrivateKeyPassword string `env:"HERMEZCORE_ETHERMAN_PRIVATEKEY_PASSWORD"`
}
