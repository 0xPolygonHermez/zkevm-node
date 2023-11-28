package state

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestGetL2Hash computes the l2 hash of a transaction for testing purposes
func TestGetL2Hash(tx types.Transaction, sender common.Address) (common.Hash, error) {
	return getL2Hash(tx, sender)
}

// GetL2Hash computes the l2 hash of a transaction
func GetL2Hash(tx types.Transaction) (common.Hash, error) {
	sender, err := GetSender(tx)
	if err != nil {
		// This is normal for unsigned transactions
		log.Debugf("error getting sender: %v", err)
	}

	return getL2Hash(tx, sender)
}

func getL2Hash(tx types.Transaction, sender common.Address) (common.Hash, error) {
	var input string
	input += formatL2TxHashParam(fmt.Sprintf("%x", tx.Nonce()))
	input += formatL2TxHashParam(fmt.Sprintf("%x", tx.GasPrice()))
	input += formatL2TxHashParam(fmt.Sprintf("%x", tx.Gas()))
	input += pad20Bytes(formatL2TxHashParam(fmt.Sprintf("%x", tx.To())))
	input += formatL2TxHashParam(fmt.Sprintf("%x", tx.Value()))
	if len(tx.Data()) > 0 {
		input += formatL2TxHashParam(fmt.Sprintf("%x", tx.Data()))
	}
	if sender != ZeroAddress {
		input += pad20Bytes(formatL2TxHashParam(fmt.Sprintf("%x", sender)))
	}

	h4Hash, err := merkletree.HashContractBytecode(common.Hex2Bytes(input))
	if err != nil {
		return common.Hash{}, err
	}

	return common.HexToHash(merkletree.H4ToString(h4Hash)), nil
}

// pad20Bytes pads the given address with 0s to make it 20 bytes long
func pad20Bytes(address string) string {
	const addressLength = 40

	if len(address) < addressLength {
		address = strings.Repeat("0", addressLength-len(address)) + address
	}
	return address
}

func formatL2TxHashParam(param string) string {
	param = strings.TrimLeft(param, "0x")

	if param == "00" || param == "" {
		return "00"
	}

	if len(param)%2 != 0 {
		param = "0" + param
	}

	return param
}

// GetSender gets the sender from the transaction's signature
func GetSender(tx types.Transaction) (common.Address, error) {
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(&tx)
	if err != nil {
		return common.Address{}, err
	}
	return sender, nil
}

// RlpFieldsToLegacyTx parses the rlp fields slice into a type.LegacyTx
// in this specific order:
//
// required fields:
// [0] Nonce    uint64
// [1] GasPrice *big.Int
// [2] Gas      uint64
// [3] To       *common.Address
// [4] Value    *big.Int
// [5] Data     []byte
//
// optional fields:
// [6] V        *big.Int
// [7] R        *big.Int
// [8] S        *big.Int
func RlpFieldsToLegacyTx(fields [][]byte, v, r, s []byte) (tx *types.LegacyTx, err error) {
	const (
		fieldsSizeWithoutChainID = 6
		fieldsSizeWithChainID    = 7
	)

	if len(fields) < fieldsSizeWithoutChainID {
		return nil, types.ErrTxTypeNotSupported
	}

	nonce := big.NewInt(0).SetBytes(fields[0]).Uint64()
	gasPrice := big.NewInt(0).SetBytes(fields[1])
	gas := big.NewInt(0).SetBytes(fields[2]).Uint64()
	var to *common.Address

	if fields[3] != nil && len(fields[3]) != 0 {
		tmp := common.BytesToAddress(fields[3])
		to = &tmp
	}
	value := big.NewInt(0).SetBytes(fields[4])
	data := fields[5]

	txV := big.NewInt(0).SetBytes(v)
	if len(fields) >= fieldsSizeWithChainID {
		chainID := big.NewInt(0).SetBytes(fields[6])

		// a = chainId * 2
		// b = v - 27
		// c = a + 35
		// v = b + c
		//
		// same as:
		// v = v-27+chainId*2+35
		a := new(big.Int).Mul(chainID, big.NewInt(double))
		b := new(big.Int).Sub(new(big.Int).SetBytes(v), big.NewInt(ether155V))
		c := new(big.Int).Add(a, big.NewInt(etherPre155V))
		txV = new(big.Int).Add(b, c)
	}

	txR := big.NewInt(0).SetBytes(r)
	txS := big.NewInt(0).SetBytes(s)

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       to,
		Value:    value,
		Data:     data,
		V:        txV,
		R:        txR,
		S:        txS,
	}, nil
}

// StoreTransactions is used by the sequencer to add processed transactions into
// an open batch. If the batch already has txs, the processedTxs must be a super
// set of the existing ones, preserving order.
func (s *State) StoreTransactions(ctx context.Context, batchNumber uint64, processedBlocks []*ProcessBlockResponse, txsEGPLog []*EffectiveGasPriceLog, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	// Check if last batch is closed. Note that it's assumed that only the latest batch can be open
	isBatchClosed, err := s.IsBatchClosed(ctx, batchNumber, dbTx)
	if err != nil {
		return err
	}
	if isBatchClosed {
		return ErrBatchAlreadyClosed
	}

	for _, processedBlock := range processedBlocks {
		processedTxs := processedBlock.TransactionResponses
		// check existing txs vs parameter txs
		/*
			existingTxs, err := s.GetTxsHashesByBatchNumber(ctx, batchNumber, dbTx)
			if err != nil {
				return err
			}*/

		// TODO: Refactor
		/*
			if err := CheckSupersetBatchTransactions(existingTxs, processedTxs); err != nil {
				return err
			}
		*/

		processingContext, err := s.GetProcessingContext(ctx, batchNumber, dbTx)
		if err != nil {
			return err
		}

		// firstTxToInsert := len(existingTxs)

		firstTxToInsert := 0

		for i := firstTxToInsert; i < len(processedTxs); i++ {
			processedTx := processedTxs[i]
			// if the transaction has an intrinsic invalid tx error it means
			// the transaction has not changed the state, so we don't store it
			// and just move to the next
			if executor.IsIntrinsicError(executor.RomErrorCode(processedTx.RomError)) || errors.Is(processedTx.RomError, executor.RomErr(executor.RomError_ROM_ERROR_INVALID_RLP)) {
				continue
			}

			lastL2Block, err := s.GetLastL2Block(ctx, dbTx)
			if err != nil {
				return err
			}

			header := &types.Header{
				Number:     new(big.Int).SetUint64(lastL2Block.Number().Uint64() + 1),
				ParentHash: lastL2Block.Hash(),
				Coinbase:   processingContext.Coinbase,
				Root:       processedTx.StateRoot,
				GasUsed:    processedTx.GasUsed,
				GasLimit:   s.cfg.MaxCumulativeGasUsed,
				Time:       uint64(processingContext.Timestamp.Unix()),
			}
			transactions := []*types.Transaction{&processedTx.Tx}

			receipt := generateReceipt(header.Number, processedTx)
			if !CheckLogOrder(receipt.Logs) {
				return fmt.Errorf("error: logs received from executor are not in order")
			}
			receipts := []*types.Receipt{receipt}

			// Create block to be able to calculate its hash
			block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
			block.ReceivedAt = processingContext.Timestamp

			receipt.BlockHash = block.Hash()

			storeTxsEGPData := []StoreTxEGPData{{EGPLog: nil, EffectivePercentage: uint8(processedTx.EffectivePercentage)}}
			if txsEGPLog != nil {
				storeTxsEGPData[0].EGPLog = txsEGPLog[i]
			}

			// Store L2 block and its transaction
			if err := s.AddL2Block(ctx, batchNumber, block, receipts, storeTxsEGPData, dbTx); err != nil {
				return err
			}
		}
	}
	return nil
}

// PreProcessTransaction processes the transaction in order to calculate its zkCounters before adding it to the pool
func (s *State) PreProcessTransaction(ctx context.Context, tx *types.Transaction, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	sender, err := GetSender(*tx)
	if err != nil {
		return nil, err
	}

	response, err := s.internalProcessUnsignedTransaction(ctx, tx, sender, nil, false, dbTx)
	if err != nil {
		return response, err
	}

	return response, nil
}

// ProcessUnsignedTransaction processes the given unsigned transaction.
func (s *State) ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, noZKEVMCounters bool, dbTx pgx.Tx) (*runtime.ExecutionResult, error) {
	result := new(runtime.ExecutionResult)
	response, err := s.internalProcessUnsignedTransaction(ctx, tx, senderAddress, l2BlockNumber, noZKEVMCounters, dbTx)
	if err != nil {
		return nil, err
	}

	r := response.BlockResponses[0].TransactionResponses[0]
	result.ReturnValue = r.ReturnValue
	result.GasLeft = r.GasLeft
	result.GasUsed = r.GasUsed
	result.CreateAddress = r.CreateAddress
	result.StateRoot = r.StateRoot.Bytes()

	if errors.Is(r.RomError, runtime.ErrExecutionReverted) {
		result.Err = constructErrorFromRevert(r.RomError, r.ReturnValue)
	} else {
		result.Err = r.RomError
	}

	return result, nil
}

// ProcessUnsignedTransaction processes the given unsigned transaction.
func (s *State) internalProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, noZKEVMCounters bool, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	var attempts = 1

	if s.executorClient == nil {
		return nil, ErrExecutorNil
	}
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	lastBatches, l2BlockStateRoot, err := s.GetLastNBatchesByL2BlockNumber(ctx, l2BlockNumber, 2, dbTx) // nolint: gomnd
	if err != nil {
		return nil, err
	}

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	stateRoot := l2BlockStateRoot
	timestamp := uint64(lastBatch.Timestamp.Unix())
	if l2BlockNumber != nil {
		l2Block, err := s.GetL2BlockByNumber(ctx, *l2BlockNumber, dbTx)
		if err != nil {
			return nil, err
		}
		stateRoot = l2Block.Root()

		latestL2BlockNumber, err := s.GetLastL2BlockNumber(ctx, dbTx)
		if err != nil {
			return nil, err
		}

		if *l2BlockNumber == latestL2BlockNumber {
			timestamp = uint64(time.Now().Unix())
		}
	}

	forkID := s.GetForkIDByBatchNumber(lastBatch.BatchNumber)
	loadedNonce, err := s.tree.GetNonce(ctx, senderAddress, stateRoot.Bytes())
	if err != nil {
		return nil, err
	}
	nonce := loadedNonce.Uint64()

	batchL2Data, err := EncodeUnsignedTransaction(*tx, s.cfg.ChainID, &nonce, forkID)
	if err != nil {
		log.Errorf("error encoding unsigned transaction ", err)
		return nil, err
	}

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      lastBatch.BatchNumber,
		BatchL2Data:      batchL2Data,
		From:             senderAddress.String(),
		OldStateRoot:     stateRoot.Bytes(),
		GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
		OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
		EthTimestamp:     timestamp,
		Coinbase:         lastBatch.Coinbase.String(),
		UpdateMerkleTree: cFalse,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	if noZKEVMCounters {
		processBatchRequest.NoCounters = cTrue
	}

	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.From]: %v", processBatchRequest.From)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.globalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)

	// Send Batch to the Executor
	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		if status.Code(err) == codes.ResourceExhausted || (processBatchResponse != nil && processBatchResponse.Error == executor.ExecutorError(executor.ExecutorError_EXECUTOR_ERROR_DB_ERROR)) {
			log.Errorf("error processing unsigned transaction ", err)
			for attempts < s.cfg.MaxResourceExhaustedAttempts {
				time.Sleep(s.cfg.WaitOnResourceExhaustion.Duration)
				log.Errorf("retrying to process unsigned transaction")
				processBatchResponse, err = s.executorClient.ProcessBatch(ctx, processBatchRequest)
				if status.Code(err) == codes.ResourceExhausted || (processBatchResponse != nil && processBatchResponse.Error == executor.ExecutorError(executor.ExecutorError_EXECUTOR_ERROR_DB_ERROR)) {
					log.Errorf("error processing unsigned transaction ", err)
					attempts++
					continue
				}
				break
			}
		}

		if err != nil {
			if status.Code(err) == codes.ResourceExhausted || (processBatchResponse != nil && processBatchResponse.Error == executor.ExecutorError(executor.ExecutorError_EXECUTOR_ERROR_DB_ERROR)) {
				log.Error("reporting error as time out")
				return nil, runtime.ErrGRPCResourceExhaustedAsTimeout
			}
			// Log the error
			event := &event.Event{
				ReceivedAt:  time.Now(),
				Source:      event.Source_Node,
				Level:       event.Level_Error,
				EventID:     event.EventID_ExecutorError,
				Description: fmt.Sprintf("error processing unsigned transaction %s: %v", tx.Hash(), err),
			}

			err2 := s.eventLog.LogEvent(context.Background(), event)
			if err2 != nil {
				log.Errorf("error logging event %v", err2)
			}
			log.Errorf("error processing unsigned transaction ", err)
			return nil, err
		}
	}

	if err == nil && processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.eventLog.LogExecutorError(ctx, processBatchResponse.Error, processBatchRequest)
		return nil, err
	}

	response, err := s.convertToProcessBatchResponse(processBatchResponse)
	if err != nil {
		return nil, err
	}

	if processBatchResponse.Responses[0].Error != executor.RomError_ROM_ERROR_NO_ERROR {
		err := executor.RomErr(processBatchResponse.Responses[0].Error)
		if !isEVMRevertError(err) {
			return response, err
		}
	}

	return response, nil
}

// isContractCreation checks if the tx is a contract creation
func (s *State) isContractCreation(tx *types.Transaction) bool {
	return tx.To() == nil && len(tx.Data()) > 0
}

// StoreTransaction is used by the sequencer and trusted state synchronizer to add process a transaction.
func (s *State) StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *ProcessTransactionResponse, coinbase common.Address, timestamp uint64, egpLog *EffectiveGasPriceLog, dbTx pgx.Tx) (*types.Header, error) {
	if dbTx == nil {
		return nil, ErrDBTxNil
	}

	// if the transaction has an intrinsic invalid tx error it means
	// the transaction has not changed the state, so we don't store it
	if executor.IsIntrinsicError(executor.RomErrorCode(processedTx.RomError)) {
		return nil, nil
	}

	lastL2Block, err := s.GetLastL2Block(ctx, dbTx)
	if err != nil {
		return nil, err
	}

	header := &types.Header{
		Number:     new(big.Int).SetUint64(lastL2Block.Number().Uint64() + 1),
		ParentHash: lastL2Block.Hash(),
		Coinbase:   coinbase,
		Root:       processedTx.StateRoot,
		GasUsed:    processedTx.GasUsed,
		GasLimit:   s.cfg.MaxCumulativeGasUsed,
		Time:       timestamp,
	}
	transactions := []*types.Transaction{&processedTx.Tx}

	receipt := generateReceipt(header.Number, processedTx)
	receipts := []*types.Receipt{receipt}

	// Create block to be able to calculate its hash
	block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
	block.ReceivedAt = time.Unix(int64(timestamp), 0)

	receipt.BlockHash = block.Hash()

	storeTxsEGPData := []StoreTxEGPData{{EGPLog: egpLog, EffectivePercentage: uint8(processedTx.EffectivePercentage)}}

	// Store L2 block and its transaction
	if err := s.AddL2Block(ctx, batchNumber, block, receipts, storeTxsEGPData, dbTx); err != nil {
		return nil, err
	}

	return block.Header(), nil
}

// CheckSupersetBatchTransactions verifies that processedTransactions is a
// superset of existingTxs and that the existing txs have the same order,
// returns a non-nil error if that is not the case.
func CheckSupersetBatchTransactions(existingTxHashes []common.Hash, processedTxs []*ProcessTransactionResponse) error {
	if len(existingTxHashes) > len(processedTxs) {
		return ErrExistingTxGreaterThanProcessedTx
	}
	for i, existingTxHash := range existingTxHashes {
		if existingTxHash != processedTxs[i].TxHash {
			return ErrOutOfOrderProcessedTx
		}
	}
	return nil
}

// EstimateGas for a transaction
func (s *State) EstimateGas(transaction *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, dbTx pgx.Tx) (uint64, []byte, error) {
	const ethTransferGas = 21000

	var lowEnd uint64
	var highEnd uint64

	ctx := context.Background()

	lastBatches, l2BlockStateRoot, err := s.GetLastNBatchesByL2BlockNumber(ctx, l2BlockNumber, 2, dbTx) // nolint:gomnd
	if err != nil {
		return 0, nil, err
	}

	stateRoot := l2BlockStateRoot
	if l2BlockNumber != nil {
		l2Block, err := s.GetL2BlockByNumber(ctx, *l2BlockNumber, dbTx)
		if err != nil {
			return 0, nil, err
		}
		stateRoot = l2Block.Root()
	}

	loadedNonce, err := s.tree.GetNonce(ctx, senderAddress, stateRoot.Bytes())
	if err != nil {
		return 0, nil, err
	}
	nonce := loadedNonce.Uint64()

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	lowEnd, err = core.IntrinsicGas(transaction.Data(), transaction.AccessList(), s.isContractCreation(transaction), true, false, false)
	if err != nil {
		return 0, nil, err
	}

	if lowEnd == ethTransferGas && transaction.To() != nil {
		code, err := s.tree.GetCode(ctx, *transaction.To(), stateRoot.Bytes())
		if err != nil {
			log.Warnf("error while getting transaction.to() code %v", err)
		} else if len(code) == 0 {
			return lowEnd, nil, nil
		}
	}

	if transaction.Gas() != 0 && transaction.Gas() > lowEnd {
		highEnd = transaction.Gas()
	} else {
		highEnd = s.cfg.MaxCumulativeGasUsed
	}

	var availableBalance *big.Int

	if senderAddress != ZeroAddress {
		senderBalance, err := s.tree.GetBalance(ctx, senderAddress, stateRoot.Bytes())
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				senderBalance = big.NewInt(0)
			} else {
				return 0, nil, err
			}
		}

		availableBalance = new(big.Int).Set(senderBalance)

		if transaction.Value() != nil {
			if transaction.Value().Cmp(availableBalance) > 0 {
				return 0, nil, ErrInsufficientFunds
			}

			availableBalance.Sub(availableBalance, transaction.Value())
		}
	}

	if transaction.GasPrice().BitLen() != 0 && // Gas price has been set
		availableBalance != nil && // Available balance is found
		availableBalance.Cmp(big.NewInt(0)) > 0 { // Available balance > 0
		gasAllowance := new(big.Int).Div(availableBalance, transaction.GasPrice())

		// Check the gas allowance for this account, make sure high end is capped to it
		if gasAllowance.IsUint64() && highEnd > gasAllowance.Uint64() {
			log.Debugf("Gas estimation high-end capped by allowance [%d]", gasAllowance.Uint64())
			highEnd = gasAllowance.Uint64()
		}
	}

	// Run the transaction with the specified gas value.
	// Returns a status indicating if the transaction failed, if it was reverted and the accompanying error
	testTransaction := func(gas uint64, nonce uint64, shouldOmitErr bool) (failed, reverted bool, gasUsed uint64, returnValue []byte, err error) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       transaction.To(),
			Value:    transaction.Value(),
			Gas:      gas,
			GasPrice: transaction.GasPrice(),
			Data:     transaction.Data(),
		})

		forkID := s.GetForkIDByBatchNumber(lastBatch.BatchNumber)

		batchL2Data, err := EncodeUnsignedTransaction(*tx, s.cfg.ChainID, nil, forkID)
		if err != nil {
			log.Errorf("error encoding unsigned transaction ", err)
			return false, false, gasUsed, nil, err
		}

		// Create a batch to be sent to the executor
		processBatchRequest := &executor.ProcessBatchRequest{
			OldBatchNum:      lastBatch.BatchNumber,
			BatchL2Data:      batchL2Data,
			From:             senderAddress.String(),
			OldStateRoot:     stateRoot.Bytes(),
			GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
			OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
			EthTimestamp:     uint64(lastBatch.Timestamp.Unix()),
			Coinbase:         lastBatch.Coinbase.String(),
			UpdateMerkleTree: cFalse,
			ChainId:          s.cfg.ChainID,
			ForkId:           forkID,
			ContextId:        uuid.NewString(),
		}

		log.Debugf("EstimateGas[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
		// log.Debugf("EstimateGas[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
		log.Debugf("EstimateGas[processBatchRequest.From]: %v", processBatchRequest.From)
		log.Debugf("EstimateGas[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
		log.Debugf("EstimateGas[processBatchRequest.globalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
		log.Debugf("EstimateGas[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
		log.Debugf("EstimateGas[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
		log.Debugf("EstimateGas[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
		log.Debugf("EstimateGas[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
		log.Debugf("EstimateGas[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
		log.Debugf("EstimateGas[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
		log.Debugf("EstimateGas[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)

		txExecutionOnExecutorTime := time.Now()
		processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
		log.Debugf("executor time: %vms", time.Since(txExecutionOnExecutorTime).Milliseconds())
		if err != nil {
			log.Errorf("error estimating gas: %v", err)
			return false, false, gasUsed, nil, err
		}
		if processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
			err = executor.ExecutorErr(processBatchResponse.Error)
			s.eventLog.LogExecutorError(ctx, processBatchResponse.Error, processBatchRequest)
			return false, false, gasUsed, nil, err
		}
		gasUsed = processBatchResponse.Responses[0].GasUsed

		// Check if an out of gas error happened during EVM execution
		if processBatchResponse.Responses[0].Error != executor.RomError_ROM_ERROR_NO_ERROR {
			err := executor.RomErr(processBatchResponse.Responses[0].Error)

			if (isGasEVMError(err) || isGasApplyError(err)) && shouldOmitErr {
				// Specifying the transaction failed, but not providing an error
				// is an indication that a valid error occurred due to low gas,
				// which will increase the lower bound for the search
				return true, false, gasUsed, nil, nil
			}

			if isEVMRevertError(err) {
				// The EVM reverted during execution, attempt to extract the
				// error message and return it
				returnValue := processBatchResponse.Responses[0].ReturnValue
				return true, true, gasUsed, returnValue, constructErrorFromRevert(err, returnValue)
			}

			return true, false, gasUsed, nil, err
		}

		return false, false, gasUsed, nil, nil
	}

	txExecutions := []time.Duration{}
	var totalExecutionTime time.Duration

	// Check if the highEnd is a good value to make the transaction pass
	failed, reverted, gasUsed, returnValue, err := testTransaction(highEnd, nonce, false)
	log.Debugf("Estimate gas. Trying to execute TX with %v gas", highEnd)
	if failed {
		if reverted {
			return 0, returnValue, err
		}

		// The transaction shouldn't fail, for whatever reason, at highEnd
		return 0, nil, fmt.Errorf(
			"unable to apply transaction even for the highest gas limit %d: %w",
			highEnd,
			err,
		)
	}

	if lowEnd < gasUsed {
		lowEnd = gasUsed
	}

	// Start the binary search for the lowest possible gas price
	for (lowEnd < highEnd) && (highEnd-lowEnd) > 4096 {
		txExecutionStart := time.Now()
		mid := (lowEnd + highEnd) / 2 // nolint:gomnd

		log.Debugf("Estimate gas. Trying to execute TX with %v gas", mid)

		failed, reverted, _, _, testErr := testTransaction(mid, nonce, true)
		executionTime := time.Since(txExecutionStart)
		totalExecutionTime += executionTime
		txExecutions = append(txExecutions, executionTime)
		if testErr != nil && !reverted {
			// Reverts are ignored in the binary search, but are checked later on
			// during the execution for the optimal gas limit found
			return 0, nil, testErr
		}

		if failed {
			// If the transaction failed => increase the gas
			lowEnd = mid + 1
		} else {
			// If the transaction didn't fail => make this ok value the high end
			highEnd = mid
		}
	}

	executions := int64(len(txExecutions))
	if executions > 0 {
		log.Infof("EstimateGas executed TX %v %d times in %d milliseconds", transaction.Hash(), executions, totalExecutionTime.Milliseconds())
	} else {
		log.Error("Estimate gas. Tx not executed")
	}
	return highEnd, nil, nil
}

// Checks if executor level valid gas errors occurred
func isGasApplyError(err error) bool {
	return errors.Is(err, ErrNotEnoughIntrinsicGas)
}

// Checks if EVM level valid gas errors occurred
func isGasEVMError(err error) bool {
	return errors.Is(err, runtime.ErrOutOfGas)
}

// Checks if the EVM reverted during execution
func isEVMRevertError(err error) bool {
	return errors.Is(err, runtime.ErrExecutionReverted)
}
