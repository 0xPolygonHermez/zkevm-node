package ethtxmanager

import (
	"github.com/hermeznetwork/hermez-core/ethermanv2"
)

type Client struct {
	cfg Config
}

func (c *Client) SequenceBatches(sequences []*ethermanv2.Sequence) error {
	return nil
}
