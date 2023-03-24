package jsonrpc

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config represents the configuration of the json rpc
type Config struct {
	Host         string         `mapstructure:"Host"`
	Port         int            `mapstructure:"Port"`
	ReadTimeout  types.Duration `mapstructure:"ReadTimeout"`
	WriteTimeout types.Duration `mapstructure:"WriteTimeout"`

	MaxRequestsPerIPAndSecond float64 `mapstructure:"MaxRequestsPerIPAndSecond"`

	// SequencerNodeURI is used allow Non-Sequencer nodes
	// to relay transactions to the Sequencer node
	SequencerNodeURI string `mapstructure:"SequencerNodeURI"`

	// BroadcastURI is the URL of the Trusted State broadcast service
	BroadcastURI string `mapstructure:"BroadcastURI"`

	// DefaultSenderAddress is the address that jRPC will use
	// to communicate with the state for eth_EstimateGas and eth_Call when
	// the From field is not specified because it is optional
	DefaultSenderAddress string `mapstructure:"DefaultSenderAddress"`

	// MaxCumulativeGasUsed is the max gas allowed per batch
	MaxCumulativeGasUsed uint64

	// ChainID is the L2 ChainID provided by the Network Config
	ChainID uint64

	// Websockets
	WebSockets WebSocketsConfig `mapstructure:"WebSockets"`

	// EnableL2SuggestedGasPricePolling enables polling of the L2 gas price to block tx in the RPC with lower gas price.
	EnableL2SuggestedGasPricePolling bool `mapstructure:"EnableL2SuggestedGasPricePolling"`
}

// WebSocketsConfig has parameters to config the rpc websocket support
type WebSocketsConfig struct {
	Enabled bool `mapstructure:"Enabled"`
	Port    int  `mapstructure:"Port"`
}
