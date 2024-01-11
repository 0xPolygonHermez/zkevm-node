package etherman

import "github.com/ethereum/go-ethereum/common"

type dataAvailabilitier interface {
	GetBatchL2Data(batchNum uint64, hash common.Hash) ([]byte, error)
}
