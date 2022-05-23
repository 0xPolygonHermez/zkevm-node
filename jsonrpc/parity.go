package jsonrpc

// Parity contains implementations for the "parity" RPC endpoints
type Parity struct{}

// PendingTransactions creates a response for parity_pendingTransactions request.
// See https://openethereum.github.io/JSONRPC-parity-module#parity_pendingtransactions
func (p *Parity) PendingTransactions() (interface{}, error) {
	return []interface{}{}, nil
}
