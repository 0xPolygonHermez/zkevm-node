package sequencer

import "errors"

var (
	// ErrExpiredTransaction happens when the transaction is expired
	ErrExpiredTransaction = errors.New("transaction expired")
	// ErrBreakEvenGasPriceEmpty happens when the breakEven or gasPrice is nil or zero
	ErrBreakEvenGasPriceEmpty = errors.New("breakEven and gasPrice cannot be nil or zero")
	// ErrEffectiveGasPriceReprocess happens when the effective gas price requires reexecution
	ErrEffectiveGasPriceReprocess = errors.New("effective gas price requires reprocessing the transaction")
	// ErrZeroL1GasPrice is returned if the L1 gas price is 0.
	ErrZeroL1GasPrice = errors.New("L1 gas price 0")
	// ErrDuplicatedNonce is returned when adding a new tx to the worker and there is an existing tx
	// with the same nonce and higher gasPrice (in this case we keep the existing tx)
	ErrDuplicatedNonce = errors.New("duplicated nonce")
	// ErrReplacedTransaction is returned when an existing tx is replaced by a new tx with the same nonce and higher gasPrice
	ErrReplacedTransaction = errors.New("replaced transaction")
	// ErrGetBatchByNumber happens when we get an error trying to get a batch by number (GetBatchByNumber)
	ErrGetBatchByNumber = errors.New("get batch by number error")
	// ErrDecodeBatchL2Data happens when we get an error trying to decode BatchL2Data (DecodeTxs)
	ErrDecodeBatchL2Data = errors.New("decoding BatchL2Data error")
	// ErrProcessBatch happens when we get an error trying to process (executor) a batch
	ErrProcessBatch = errors.New("processing batch error")
	// ErrProcessBatchOOC happens when we get an OOC when processing (executor) a batch
	ErrProcessBatchOOC = errors.New("processing batch OOC")
	// ErrStateRootNoMatch happens when the SR returned for a full batch processing (sanity check) doesn't match
	// the SR calculated when filling a batch tx by tx
	ErrStateRootNoMatch = errors.New("state root no match")
	// ErrExecutorError happens when we got an executor error when processing a batch
	ErrExecutorError = errors.New("executor error")
)
