package jsonrpc

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config represents the configuration of the json rpc
type Config struct {
	// Host defines the network adapter that will be used to serve the HTTP requests
	Host string `mapstructure:"Host"`

	// Port defines the port to serve the endpoints via HTTP
	Port int `mapstructure:"Port"`

	// ReadTimeout is the HTTP server read timeout
	// check net/http.server.ReadTimeout and net/http.server.ReadHeaderTimeout
	ReadTimeout types.Duration `mapstructure:"ReadTimeout"`

	// WriteTimeout is the HTTP server write timeout
	// check net/http.server.WriteTimeout
	WriteTimeout types.Duration `mapstructure:"WriteTimeout"`

	// MaxRequestsPerIPAndSecond defines how much requests a single IP can
	// send within a single second
	MaxRequestsPerIPAndSecond float64 `mapstructure:"MaxRequestsPerIPAndSecond"`

	// SequencerNodeURI is used allow Non-Sequencer nodes
	// to relay transactions to the Sequencer node
	SequencerNodeURI string `mapstructure:"SequencerNodeURI"`

	// MaxCumulativeGasUsed is the max gas allowed per batch
	MaxCumulativeGasUsed uint64

	// WebSockets configuration
	WebSockets WebSocketsConfig `mapstructure:"WebSockets"`

	// EnableL2SuggestedGasPricePolling enables polling of the L2 gas price to block tx in the RPC with lower gas price.
	EnableL2SuggestedGasPricePolling bool `mapstructure:"EnableL2SuggestedGasPricePolling"`

	// TraceBatchUseHTTPS enables, in the debug_traceBatchByNum endpoint, the use of the HTTPS protocol (instead of HTTP)
	// to do the parallel requests to RPC.debug_traceTransaction endpoint
	TraceBatchUseHTTPS bool `mapstructure:"TraceBatchUseHTTPS"`

	// EnablePendingTransactionFilter enables pending transaction filter that can support query L2 pending transaction
	EnablePendingTransactionFilter bool `mapstructure:"EnablePendingTransactionFilter"`

	// Nacos configuration
	Nacos NacosConfig `mspstructure:"Nacos"`
}

// WebSocketsConfig has parameters to config the rpc websocket support
type WebSocketsConfig struct {
	// Enabled defines if the WebSocket requests are enabled or disabled
	Enabled bool `mapstructure:"Enabled"`

	// Host defines the network adapter that will be used to serve the WS requests
	Host string `mapstructure:"Host"`

	// Port defines the port to serve the endpoints via WS
	Port int `mapstructure:"Port"`
}

type NacosConfig struct {
	// URLs nacos server urls for discovery service of rest api, url is separated by ","
	URLs string `mapstructure:"URLs"`

	// NamespaceId nacos namepace id for discovery service of rest api
	NamespaceId string `mapstructure:"NamespaceId"`

	// ApplicationName rest application name in  nacos
	ApplicationName string `mapstructure:"ApplicationName"`

	// ExternalListenAddr Set the rest-server external ip and port, when it is launched by Docker
	ExternalListenAddr string `mapstructure:"ExternalListenAddr"`
}
