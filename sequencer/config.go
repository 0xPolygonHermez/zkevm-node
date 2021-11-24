package sequencer

import "github.com/hermeznetwork/hermez-core/etherman"

// Config represents the configuration of a sequencer
type Config struct {
	SendBatchFrequencyInSeconds uint
	Etherman                    etherman.Config
}
