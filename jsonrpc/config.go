package jsonrpc

// Config represents the configuration of the json rpc
type Config struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`

	MaxRequestsPerIPAndSecond float64 `mapstructure:"MaxRequestsPerIPAndSecond"`

	// TrustedSequencerURI is used allow Permission less nodes
	// to relay transactions to the Trusted node
	TrustedNodeURI string `mapstructure:"URI"`
}
