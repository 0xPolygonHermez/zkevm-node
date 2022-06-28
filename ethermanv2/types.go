package ethermanv2

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/ethermanv2/smartcontracts/proofofefficiency"
)

// Block struct
type Block struct {
	BlockNumber           uint64
	BlockHash             common.Hash
	ParentHash            common.Hash
	GlobalExitRoots       []GlobalExitRoot
	ForcedBatches         []ForcedBatch
	SequencedBatches      [][]SequencedBatch
	VerifiedBatches       []VerifiedBatch
	SequencedForceBatches []SequencedForceBatch
	ReceivedAt            time.Time
}

// GlobalExitRoot struct
type GlobalExitRoot struct {
	BlockNumber       uint64
	GlobalExitRootNum *big.Int
	MainnetExitRoot   common.Hash
	RollupExitRoot    common.Hash
	GlobalExitRoot    common.Hash
}

// SequencedBatches represents virtual batches
type SequencedBatch struct {
	BatchNumber uint64
	Sequencer   common.Address
	TxHash      common.Hash
	proofofefficiency.ProofOfEfficiencyBatchData
}

// ForcedBatch represents a ForcedBatch
type ForcedBatch struct {
	BlockNumber       uint64
	ForcedBatchNumber uint64
	Sequencer         common.Address
	GlobalExitRoot    common.Hash
	RawTxsData        []byte
	ForcedAt          time.Time
}

// VerifyBatch represents a VerifyBatch
type VerifiedBatch struct {
	BlockNumber uint64
	BatchNumber uint64
	Aggregator  common.Address
	TxHash      common.Hash
}

// SequencedForceBatch is a sturct to track the ForceSequencedBatches event.
type SequencedForceBatch struct {
	LastBatchSequenced uint64
	ForceBatchNumber   uint64
	Sequencer          common.Address
}
