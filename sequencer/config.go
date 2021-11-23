package sequencer

import "github.com/hermeznetwork/hermez-core/etherman"

type Config struct {
	SendBatchFrequencyInSeconds uint
	Etherman                    etherman.Config
}
