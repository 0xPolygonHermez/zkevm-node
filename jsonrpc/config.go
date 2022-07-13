package jsonrpc

// Config represents the configuration of the json rpc
type Config struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`

	MaxRequestsPerIPAndSecond float64 `mapstructure:"MaxRequestsPerIPAndSecond"`

	// SequencerNodeURI is used allow Non-Sequencer nodes
	// to relay transactions to the Sequencer node
	SequencerNodeURI string `mapstructure:"URI"`
}
