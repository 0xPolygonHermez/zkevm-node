package sequencer

import "github.com/hermeznetwork/hermez-core/etherman"

type Config struct {
	PrivateKey                  string
	SendBatchFrequencyInSeconds uint
	Etherman                    etherman.Config
}
