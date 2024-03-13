package etherman

import "github.com/ethereum/go-ethereum/common"

type dataAvailabilityProvider interface {
	GetBatchL2Data(batchNum []uint64, hash []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error)
}
