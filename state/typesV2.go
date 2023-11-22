package state

import (
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ProcessBatchResponseV2 represents the response of a batch process for forkID >= ETROG
type ProcessBatchResponseV2 struct {
	NewStateRoot         common.Hash
	NewAccInputHash      common.Hash
	NewLocalExitRoot     common.Hash
	NewBatchNumber       uint64
	UsedZkCounters       ZKCounters
	BlockResponses       []*ProcessBlockResponseV2
	ExecutorError        error
	ReadWriteAddresses   map[common.Address]*InfoReadWrite
	IsRomLevelError      bool
	IsExecutorLevelError bool
	IsRomOOCError        bool
	FlushID              uint64
	StoredFlushID        uint64
	ProverID             string
	GasUsed              uint64
	SMTKeys              []merkletree.Key
	ProgramKeys          []merkletree.Key
	ForkID               uint64
}

// ProcessBlockResponseV2 represents the response of a block process for forkID >= ETROG
type ProcessBlockResponseV2 struct {
	ParentHash           common.Hash
	Coinbase             common.Address
	GasLimit             uint64
	BlockNumber          uint64
	Timestamp            uint64
	GlobalExitRoot       common.Hash
	BlockHashL1          common.Hash
	GasUsed              uint64
	BlockInfoRoot        common.Hash
	BlockHash            common.Hash
	TransactionResponses []*ProcessTransactionResponseV2
	Logs                 []*types.Log
}

// ProcessTransactionResponseV2 represents the response of a tx process  for forkID >= ETROG
type ProcessTransactionResponseV2 struct {
	// TxHash is the hash of the transaction
	TxHash common.Hash
	// Type indicates legacy transaction
	// It will be always 0 (legacy) in the executor
	Type uint32
	// ReturnValue is the returned data from the runtime (function result or data supplied with revert opcode)
	ReturnValue []byte
	// GasLeft is the total gas left as result of execution
	GasLeft uint64
	// GasUsed is the total gas used as result of execution or gas estimation
	GasUsed uint64
	// GasRefunded is the total gas refunded as result of execution
	GasRefunded uint64
	// RomError represents any error encountered during the execution
	RomError error
	// CreateAddress is the new SC Address in case of SC creation
	CreateAddress common.Address
	// StateRoot is the State Root
	StateRoot common.Hash
	// Logs emitted by LOG opcode
	Logs []*types.Log
	// ChangesStateRoot indicates if this tx affects the state
	ChangesStateRoot bool
	// Tx is the whole transaction object
	Tx types.Transaction
	// FullTrace contains the call trace.
	FullTrace instrumentation.FullTrace
	// EffectiveGasPrice effective gas price used for the tx
	EffectiveGasPrice string
	//EffectivePercentage effective percentage used for the tx
	EffectivePercentage uint32
	//HasGaspriceOpcode flag to indicate if opcode 'GASPRICE' has been called
	HasGaspriceOpcode bool
	//HasBalanceOpcode flag to indicate if opcode 'BALANCE' has been called
	HasBalanceOpcode bool
}
