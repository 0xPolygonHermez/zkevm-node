package main

import "github.com/ethereum/go-ethereum/common"

type reprocessingOutputer interface {
	start(fromBatchNumber uint64, toBatchNumber uint64, l2ChainId uint64)
	startProcessingBatch(current_batch_number uint64)
	numOfTransactionsInBatch(numOfTrs int)
	addTransactionError(trxIndex int, err error)
	isWrittenOnHashDB(isWritten bool, flushid uint64)
	finishProcessingBatch(stateRoot common.Hash, err error)
	end(err error)
}
