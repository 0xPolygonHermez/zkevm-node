package state

// Config is state config
type Config struct {
	// Arity represents the maximum children the merkle tree node can have
	Arity uint8
	// DefaultChainID is the common ChainID to all the sequencers
	DefaultChainID uint64
}
