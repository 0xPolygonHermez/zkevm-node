// Package ethtxmanager handles ethereum transactions:  It makes
// calls to send and to aggregate batch, checks possible errors, like wrong nonce or gas limit too low
// and make correct adjustments to request according to it. Also it tracks transaction receipt and status
// of tx in case tx is rejected and send signals to sequencer/aggregator to resend sequence/batch
package ethtxmanager

import (
	"github.com/hermeznetwork/hermez-core/ethermanv2"
)

// Client for eth tx manager
type Client struct {
	cfg Config
}

// New creates new eth tx manager
func New(cfg Config) {

}

// SequenceBatches send request to ethereum
func (c *Client) SequenceBatches(sequences []*ethermanv2.Sequence) error {
	return nil
}
