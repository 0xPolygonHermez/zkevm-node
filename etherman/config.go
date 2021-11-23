package etherman

import "github.com/ethereum/go-ethereum/common"

// Config represents the configuration of the etherman
type Config struct {
	URL        string
	PoeAddress common.Address
}
