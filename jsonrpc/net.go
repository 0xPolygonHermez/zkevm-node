package jsonrpc

// Net is the net jsonrpc endpoint
type Net struct{}

// Version returns the current network id
func (n *Net) Version() (interface{}, error) {
	return "0x99999999", nil // 2576980377
}
