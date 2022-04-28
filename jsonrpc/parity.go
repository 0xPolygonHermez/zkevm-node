package jsonrpc

type Parity struct{}

func (p *Parity) PendingTransactions() (interface{}, error) {
	return []interface{}{}, nil
}
