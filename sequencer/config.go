package sequencer

import "math/big"

type Config struct {
	sendBatchFrequencyInSeconds uint
	baseFee *big.Int
}
