package state

import "github.com/ethereum/go-ethereum/common"

type ProcessBatchRequest struct {
	BatchNumber      uint64
	StateRoot        common.Hash
	GlobalExitRoot   common.Hash
	OldAccInputHash  common.Hash
	TxData           []byte
	SequencerAddress common.Address
	Timestamp        uint64
}
