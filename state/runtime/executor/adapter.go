package executor

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	"github.com/hermeznetwork/hermez-core/state/runtime/executor/pb"
)

// Adapter exposes the Executor methods required by the state and translates them into
// gRPC calls using its client member
type Adapter struct {
	grpcClient pb.ExecutorServiceClient
	stateRoot  []byte
	LastBatch  *state.Batch
}

// NewAdapter is the constructor of Adapter
func NewAdapter(client pb.ExecutorServiceClient, stateRoot []byte) *Adapter {
	return &Adapter{
		grpcClient: client,
		stateRoot:  stateRoot,
	}
}

// ProcessBatch processes all transactions inside a batch
func (a *Adapter) ProcessBatch(ctx context.Context, batch *state.Batch) error {
	var receipts []*state.Receipt

	request := &pb.ProcessBatchRequest{
		StateRoot:   a.stateRoot,
		BatchL2Data: batch.RawTxsData,
	}

	result, err := a.grpcClient.ProcessBatch(ctx, request)
	if err != nil {
		return err
	}

	for i, response := range result.Responses {
		receipt := a.generateReceipt(batch, response, result.CumulativeGasUsed, i)
		receipts = append(receipts, receipt)
	}

	batch.Receipts = receipts

	// Set batch Header
	a.populateBatchHeader(batch, result)

	return nil
}

// ProcessTransaction processes a transaction
func (a *Adapter) ProcessTransaction(ctx context.Context, tx *types.Transaction) *runtime.ExecutionResult {
	request := &pb.ProcessTransactionRequest{
		StateRoot: a.stateRoot,
		TxL2Data:  tx.Data(),
	}

	result, err := a.grpcClient.ProcessTransaction(ctx, request)
	if err != nil {
		log.Debugf("Error on ProcessTransaction: %v", err)
		return nil
	}

	return &runtime.ExecutionResult{
		ReturnValue:   result.ReturnValue,
		GasLeft:       result.GasLeft,
		GasUsed:       result.GasUsed,
		Err:           fmt.Errorf(result.Error),
		CreateAddress: common.HexToAddress(result.CreateAddress),
		StateRoot:     result.StateRoot,
	}
}

// EstimateGas for a transaction
func (a *Adapter) EstimateGas(ctx context.Context, tx *types.Transaction) (uint64, error) {
	request := &pb.ProcessTransactionRequest{
		StateRoot: a.stateRoot,
		TxL2Data:  tx.Data(),
	}

	result, err := a.grpcClient.EstimateGas(ctx, request)
	if err != nil {
		log.Debugf("Error on ProcessTransaction: %v", err)
		return 0, err
	}

	return result.GasUsed, fmt.Errorf(result.Error)
}

// DebugTransaction generates Geth Style traces
func (a *Adapter) DebugTransaction(ctx context.Context, tx *types.Transaction, tracer string) {
	request := &pb.ProcessTransactionRequest{
		StateRoot: a.stateRoot,
		TxL2Data:  tx.Data(),
		Tracer:    tracer,
	}

	result, err := a.grpcClient.DebugTransaction(ctx, request)
	if err != nil {
		log.Debugf("Error on ProcessTransaction: %v", err)
	}

	log.Debugf("$v", result)
}

// DebugBatch debugs all transactions inside a batch
func (a *Adapter) DebugBatch(ctx context.Context, batch *state.Batch) {
	request := &pb.ProcessBatchRequest{
		StateRoot:   a.stateRoot,
		BatchL2Data: batch.RawTxsData,
	}

	result, err := a.grpcClient.DebugBatch(ctx, request)
	if err != nil {
		log.Debugf("Error on DebugBatch: %v", err)
	}

	log.Debugf("$v", result)
}

func (a *Adapter) generateReceipt(batch *state.Batch, response *pb.ProcessTransactionResponse, cumulativeGasUsed uint64, index int) *state.Receipt {
	receipt := &state.Receipt{}
	receipt.Type = uint8(response.Type)
	receipt.PostState = response.StateRoot

	if response.Error == "" {
		receipt.Status = types.ReceiptStatusSuccessful
	} else {
		receipt.Status = types.ReceiptStatusFailed
	}

	receipt.CumulativeGasUsed = cumulativeGasUsed
	receipt.BlockNumber = batch.Number()
	receipt.BlockHash = batch.Hash()
	receipt.GasUsed = response.GasUsed
	receipt.TxHash = common.BytesToHash(response.TxHash)
	receipt.TransactionIndex = uint(index)
	receipt.ContractAddress = common.HexToAddress(response.CreateAddress)
	/*
		receipt.To = (*common.Address)(common.HexToAddress(response.))
		if senderAddress != nil {
			receipt.From = *senderAddress
		}
	*/

	return receipt
}

func (a *Adapter) populateBatchHeader(batch *state.Batch, response *pb.ProcessBatchResponse) {
	parentHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	if a.LastBatch != nil {
		parentHash = a.LastBatch.Hash()
	}

	rr := make([]*types.Receipt, 0, len(batch.Receipts))
	for _, receipt := range batch.Receipts {
		r := receipt.Receipt
		r.Logs = a.getLogs(receipt.TxHash, response)
		rr = append(rr, &r)
	}
	block := types.NewBlock(batch.Header, batch.Transactions, batch.Uncles, rr, &trie.StackTrie{})

	batch.Header.ParentHash = parentHash
	batch.Header.UncleHash = block.UncleHash()
	batch.Header.Coinbase = batch.Sequencer
	batch.Header.Root = common.BytesToHash(response.Responses[len(response.Responses)-1].StateRoot)
	batch.Header.TxHash = block.TxHash()
	batch.Header.ReceiptHash = block.ReceiptHash()
	batch.Header.Bloom = block.Bloom()
	batch.Header.Difficulty = new(big.Int).SetUint64(0)
	batch.Header.GasLimit = 30000000
	batch.Header.GasUsed = response.CumulativeGasUsed
	batch.Header.Time = uint64(time.Now().Unix())
	batch.Header.MixDigest = block.MixDigest()
	batch.Header.Nonce = block.Header().Nonce
}

func (a *Adapter) getLogs(txHash common.Hash, response *pb.ProcessBatchResponse) []*types.Log {
	var logs []*pb.Log

	for _, response := range response.Responses {
		if common.BytesToHash(response.TxHash) == txHash {
			logs = response.Logs
		}
	}

	returnedLogs := make([]*types.Log, 0, len(logs))

	for _, log := range logs {
		txLog := &types.Log{
			Address: common.HexToAddress(log.Address),
			Topics:  a.getTopics(log.Topics),
			Data:    log.Data,
			TxHash:  common.BytesToHash(log.TxHash),
			TxIndex: uint(log.TxIndex),
			Index:   uint(log.Index),
			Removed: log.Removed,
		}

		returnedLogs = append(returnedLogs, txLog)
	}

	return returnedLogs
}

func (a *Adapter) getTopics(topics [][]byte) []common.Hash {
	formatedTopics := make([]common.Hash, 0, len(topics))
	for _, topic := range topics {
		formatedTopics = append(formatedTopics, common.BytesToHash(topic))
	}
	return formatedTopics
}
