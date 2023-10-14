package main

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type trxErrorEntry struct {
	trxIndex int
	err      error
}

type reprocessingOutputPretty struct {
	fromBatchNumber                   uint64
	toBatchNumber                     uint64
	l2ChainId                         uint64
	timeStart                         time.Time
	trxErrors                         []trxErrorEntry
	thereisABatchProcessingInProgress bool
	currentBatchNumber                uint64
	estimatedTime                     estimatedTimeOfArrival
}

func (o *reprocessingOutputPretty) start(fromBatchNumber uint64, toBatchNumber uint64, l2ChainId uint64) {
	o.fromBatchNumber = fromBatchNumber
	o.toBatchNumber = toBatchNumber
	o.l2ChainId = l2ChainId
	o.timeStart = time.Now()
	o.trxErrors = make([]trxErrorEntry, 0)
	o.estimatedTime.start(int(toBatchNumber - fromBatchNumber))
	fmt.Printf("START: batches [%d to %d] l2ChainId=[%d]\n", fromBatchNumber, toBatchNumber, l2ChainId)
}

func (o *reprocessingOutputPretty) numOfTransactionsInBatch(numOfTrs int) {
	fmt.Printf(" ntx: %3d", numOfTrs)
}

func (o *reprocessingOutputPretty) startProcessingBatch(current_batch_number uint64) {
	fmt.Printf("\t batch %6d %2.2f%%: ...", current_batch_number, float64(current_batch_number-o.fromBatchNumber)*100/float64(o.toBatchNumber-o.fromBatchNumber))
	o.currentBatchNumber = current_batch_number
	o.thereisABatchProcessingInProgress = true
}
func (o *reprocessingOutputPretty) addTransactionError(trxIndex int, err error) {
	o.trxErrors = append(o.trxErrors, trxErrorEntry{trxIndex: trxIndex, err: err})
}

func (o *reprocessingOutputPretty) finishProcessingBatch(stateRoot common.Hash, err error) {
	estimatedTime, _, itemsPerSecond := o.estimatedTime.step(1)
	fmt.Printf(" ETA:%10s speed:%3.1f batch/s ", estimatedTime.Round(time.Second), itemsPerSecond)
	if err == nil {
		fmt.Printf(" StateRoot:%30s [OK]\n", stateRoot)
	} else {
		fmt.Printf(" StateRoot:%30s [ERROR] %v\n", "", err)
	}
	for _, trxError := range o.trxErrors {
		fmt.Printf("\t\t[ERROR] trx %d: %v\n", trxError.trxIndex, trxError.err)
	}
	o.trxErrors = make([]trxErrorEntry, 0)
	o.thereisABatchProcessingInProgress = false
}

func (o *reprocessingOutputPretty) isWrittenOnHashDB(isWritten bool, flushid uint64) {
	if isWritten {
		fmt.Printf(" WRITE (flush:%5d) ", flushid)
	}
}

func (o *reprocessingOutputPretty) end(err error) {
	if err != nil {
		if o.thereisABatchProcessingInProgress {
			o.finishProcessingBatch(common.Hash{}, err)
		} else {
			fmt.Printf("\n[ERROR] %v", err)
		}
	}
}
