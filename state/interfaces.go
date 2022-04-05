package state

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state/tree"
)

// Consumer interfaces required by the package.

// merkletree contains the methods required to interact with the Merkle tree.
type merkletree interface {
	SupportsDBTransactions() bool
	BeginDBTransaction(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetBalance(ctx context.Context, address common.Address, root []byte) (*big.Int, error)
	GetNonce(ctx context.Context, address common.Address, root []byte) (*big.Int, error)
	GetCode(ctx context.Context, address common.Address, root []byte) ([]byte, error)
	GetCodeHash(ctx context.Context, address common.Address, root []byte) ([]byte, error)
	GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root []byte) (*big.Int, error)

	SetBalance(ctx context.Context, address common.Address, balance *big.Int, root []byte) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetNonce(ctx context.Context, address common.Address, nonce *big.Int, root []byte) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetCode(ctx context.Context, address common.Address, code []byte, root []byte) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetStorageAt(ctx context.Context, address common.Address, key *big.Int, value *big.Int, root []byte) (newRoot []byte, proof *tree.UpdateProof, err error)
}

// storage is the interface of the Hermez state methods that access database.
type storage interface {
	BeginDBTransaction(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetLastBlock(ctx context.Context) (*Block, error)
	GetPreviousBlock(ctx context.Context, offset uint64) (*Block, error)
	GetBlockByHash(ctx context.Context, hash common.Hash) (*Block, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error)
	GetLastBlockNumber(ctx context.Context) (uint64, error)
	GetLastBatch(ctx context.Context, isVirtual bool) (*Batch, error)
	GetPreviousBatch(ctx context.Context, isVirtual bool, offset uint64) (*Batch, error)
	GetBatchByHash(ctx context.Context, hash common.Hash) (*Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64) (*Batch, error)
	GetLastBatchByStateRoot(ctx context.Context, stateRoot []byte) (*Batch, error)
	GetBatchHeader(ctx context.Context, batchNumber uint64) (*types.Header, error)
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	GetLastConsolidatedBatchNumber(ctx context.Context) (uint64, error)
	GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error)
	GetTransactionCount(ctx context.Context, address common.Address) (uint64, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*Receipt, error)
	Reset(ctx context.Context, block *Block) error
	ConsolidateBatch(ctx context.Context, batchNumber uint64, consolidatedTxHash common.Hash, consolidatedAt time.Time, aggregator common.Address) error
	GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error)
	AddSequencer(ctx context.Context, seq Sequencer) error
	GetSequencer(ctx context.Context, address common.Address) (*Sequencer, error)
	AddBlock(ctx context.Context, block *Block) error
	SetLastBatchNumberSeenOnEthereum(ctx context.Context, batchNumber uint64) error
	GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error)
	AddBatch(ctx context.Context, batch *Batch) error
	AddTransaction(ctx context.Context, tx *types.Transaction, batchNumber uint64, index uint) error
	AddReceipt(ctx context.Context, receipt *Receipt) error
	AddLog(ctx context.Context, log types.Log) error
	GetLogs(ctx context.Context, fromBatch uint64, toBatch uint64, addresses []common.Address, topics [][]common.Hash, batchHash *common.Hash) ([]*types.Log, error)
	SetLastBatchNumberConsolidatedOnEthereum(ctx context.Context, batchNumber uint64) error
	GetLastBatchNumberConsolidatedOnEthereum(ctx context.Context) (uint64, error)
}
