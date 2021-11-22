package aggregator

import "github.com/hermeznetwork/hermez-core/etherman"

type Config struct {
	// PrivateKey is used to sign l1 tx
	PrivateKey string
	Etherman   etherman.Config
}
