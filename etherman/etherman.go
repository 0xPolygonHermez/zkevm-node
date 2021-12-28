package etherman

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/state"
)

var (
	newBatchEventSignatureHash        = crypto.Keccak256Hash([]byte("SendBatch(uint32,address)"))
	consolidateBatchSignatureHash     = crypto.Keccak256Hash([]byte("VerifyBatch(uint32,address)"))
	newSequencerSignatureHash         = crypto.Keccak256Hash([]byte("RegisterSequencer(address,string,uint32)"))
	ownershipTransferredSignatureHash = crypto.Keccak256Hash([]byte("OwnershipTransferred(address,address)"))
)

var (
	// ErrNotFound is used when the object is not found
	ErrNotFound = errors.New("Not found")
)

// EtherMan represents an Ethereum Manager
type EtherMan interface {
	EthBlockByNumber(ctx context.Context, blockNum uint64) (*types.Block, error)
	GetBatchesByBlock(ctx context.Context, blockNum uint64, blockHash *common.Hash) ([]state.Block, error)
	GetBatchesByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]state.Block, error)
	SendBatch(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*types.Transaction, error)
	ConsolidateBatch(batchNum *big.Int, proof *proverclient.Proof) (*types.Transaction, error)
	RegisterSequencer(url string) (*types.Transaction, error)
	GetAddress() common.Address
	GetDefaultChainID() (*big.Int, error)
	EstimateSendBatchCost(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*big.Int, error)
	GetLatestProposedBatchNumber() (uint64, error)
	GetSequencerCollateral(batchNumber uint64) (*big.Int, error)
}

type ethClienter interface {
	ethereum.ChainReader
	ethereum.LogFilterer
	ethereum.TransactionReader
}

// ClientEtherMan is a simple implementation of EtherMan
type ClientEtherMan struct {
	EtherClient ethClienter
	PoE         *proofofefficiency.Proofofefficiency
	SCAddresses []common.Address

	auth *bind.TransactOpts
}

// NewEtherman creates a new etherman
func NewEtherman(cfg Config, auth *bind.TransactOpts, PoEAddr common.Address) (*ClientEtherMan, error) {
	// Connect to ethereum node
	ethClient, err := ethclient.Dial(cfg.URL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", cfg.URL, err)
		return nil, err
	}
	// Create smc clients
	poe, err := proofofefficiency.NewProofofefficiency(PoEAddr, ethClient)
	if err != nil {
		return nil, err
	}
	var scAddresses []common.Address
	scAddresses = append(scAddresses, PoEAddr)

	return &ClientEtherMan{EtherClient: ethClient, PoE: poe, SCAddresses: scAddresses, auth: auth}, nil
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number
func (etherMan *ClientEtherMan) EthBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	block, err := etherMan.EtherClient.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		if errors.Is(err, ethereum.NotFound) || err.Error() == "block does not exist in blockchain" {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return block, nil
}

// GetBatchesByBlock function retrieves the batches information that are included in a specific ethereum block
func (etherMan *ClientEtherMan) GetBatchesByBlock(ctx context.Context, blockNumber uint64, blockHash *common.Hash) ([]state.Block, error) {
	// First filter query
	var blockNumBigInt *big.Int
	if blockHash == nil {
		blockNumBigInt = new(big.Int).SetUint64(blockNumber)
	}
	query := ethereum.FilterQuery{
		BlockHash: blockHash,
		FromBlock: blockNumBigInt,
		ToBlock:   blockNumBigInt,
		Addresses: etherMan.SCAddresses,
	}
	blocks, err := etherMan.readEvents(ctx, query)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// GetBatchesByBlockRange function retrieves the batches information that are included in all this ethereum blocks
// from block x to block y
func (etherMan *ClientEtherMan) GetBatchesByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]state.Block, error) {
	// First filter query
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		Addresses: etherMan.SCAddresses,
	}
	if toBlock != nil {
		query.ToBlock = new(big.Int).SetUint64(*toBlock)
	}
	blocks, err := etherMan.readEvents(ctx, query)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// SendBatch function allows the sequencer send a new batch proposal to the rollup
func (etherMan *ClientEtherMan) SendBatch(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*types.Transaction, error) {
	return etherMan.sendBatch(ctx, etherMan.auth, txs, maticAmount)
}

func (etherMan *ClientEtherMan) sendBatch(ctx context.Context, opts *bind.TransactOpts, txs []*types.Transaction, maticAmount *big.Int) (*types.Transaction, error) {
	if len(txs) == 0 {
		return nil, errors.New("invalid txs: is empty slice")
	}
	var data [][]byte
	for _, tx := range txs {
		a := new(bytes.Buffer)
		err := tx.EncodeRLP(a)
		if err != nil {
			return nil, err
		}
		log.Debug("Coded tx: ", hex.EncodeToString(a.Bytes()))
		data = append(data, a.Bytes())
	}
	b := new(bytes.Buffer)
	err := rlp.Encode(b, data)
	if err != nil {
		return nil, err
	}

	tx, err := etherMan.PoE.SendBatch(etherMan.auth, b.Bytes(), maticAmount)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// ConsolidateBatch function allows the aggregator send the proof for a batch and consolidate it
func (etherMan *ClientEtherMan) ConsolidateBatch(batchNumber *big.Int, proof *proverclient.Proof) (*types.Transaction, error) {
	newLocalExitRoot, err := byteSliceToFixedByteArray(proof.PublicInputs.NewLocalExitRoot)
	if err != nil {
		return nil, err
	}
	newStateRoot, err := byteSliceToFixedByteArray(proof.PublicInputs.NewStateRoot)
	if err != nil {
		return nil, err
	}

	proofA, err := strSliceToBigIntArray(proof.ProofA)
	if err != nil {
		return nil, err
	}

	proofB, err := proofSlcToIntArray(proof.ProofB)
	if err != nil {
		return nil, err
	}
	proofC, err := strSliceToBigIntArray(proof.ProofC)
	if err != nil {
		return nil, err
	}

	tx, err := etherMan.PoE.VerifyBatch(
		etherMan.auth,
		newLocalExitRoot,
		newStateRoot,
		uint32(batchNumber.Uint64()),
		proofA,
		proofB,
		proofC,
	)

	if err != nil {
		return nil, err
	}
	return tx, nil
}

// RegisterSequencer function allows to register a new sequencer in the rollup
func (etherMan *ClientEtherMan) RegisterSequencer(url string) (*types.Transaction, error) {
	tx, err := etherMan.PoE.RegisterSequencer(etherMan.auth, url)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (etherMan *ClientEtherMan) readEvents(ctx context.Context, query ethereum.FilterQuery) ([]state.Block, error) {
	logs, err := etherMan.EtherClient.FilterLogs(ctx, query)
	if err != nil {
		return []state.Block{}, err
	}
	blocks := make(map[common.Hash]state.Block)
	var blockKeys []common.Hash
	for _, vLog := range logs {
		block, err := etherMan.processEvent(ctx, vLog)
		if err != nil {
			log.Warnf("error processing event. Retrying... Error: %w. vLog: %+v", err, vLog)
			break
		}
		if block == nil {
			continue
		}
		if b, exists := blocks[block.BlockHash]; exists {
			b.Batches = append(blocks[block.BlockHash].Batches, block.Batches...)
			b.NewSequencers = append(blocks[block.BlockHash].NewSequencers, block.NewSequencers...)
			blocks[block.BlockHash] = b
		} else {
			blocks[block.BlockHash] = *block
			blockKeys = append(blockKeys, block.BlockHash)
		}
	}
	var blockArr []state.Block
	for _, hash := range blockKeys {
		blockArr = append(blockArr, blocks[hash])
	}
	return blockArr, nil
}

func (etherMan *ClientEtherMan) processEvent(ctx context.Context, vLog types.Log) (*state.Block, error) {
	switch vLog.Topics[0] {
	case newBatchEventSignatureHash:
		var block state.Block
		// Indexed parameters using topics
		var batch state.Batch
		batch.BatchNumber = new(big.Int).SetBytes(vLog.Topics[1][:]).Uint64()
		batch.Sequencer = common.BytesToAddress(vLog.Topics[2].Bytes())
		var head types.Header
		head.TxHash = vLog.TxHash
		head.Difficulty = big.NewInt(0)
		head.Number = new(big.Int).SetUint64(batch.BatchNumber)
		batch.Header = &head
		block.BlockNumber = vLog.BlockNumber
		batch.BlockNumber = vLog.BlockNumber
		maticCollateral, err := etherMan.GetSequencerCollateral(batch.BatchNumber)
		if err != nil {
			return nil, fmt.Errorf("error getting matic collateral for batch: %d. BlockNumber: %d. Error: %w", batch.BatchNumber, block.BlockNumber, err)
		}
		batch.MaticCollateral = maticCollateral
		block.BlockHash = vLog.BlockHash
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return nil, fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", block.BlockNumber, err)
		}
		block.ParentHash = fullBlock.ParentHash()
		// Read the tx for this batch.
		tx, isPending, err := etherMan.EtherClient.TransactionByHash(ctx, batch.Header.TxHash)
		if err != nil {
			return nil, err
		} else if isPending {
			return nil, fmt.Errorf("error: tx is still pending. TxHash: %s", tx.Hash().String())
		}
		batch.RawTxsData = tx.Data()
		txs, err := decodeTxs(tx.Data())
		if err != nil {
			log.Warn("No txs decoded in batch: ", batch.BatchNumber, ". This batch is inside block: ", batch.BlockNumber,
				". Error: ", err)
		}
		batch.Transactions = txs
		block.Batches = append(block.Batches, batch)
		return &block, nil
	case consolidateBatchSignatureHash:
		var block state.Block
		var batch state.Batch
		batch.BatchNumber = new(big.Int).SetBytes(vLog.Topics[1][:]).Uint64()
		batch.Aggregator = common.BytesToAddress(vLog.Topics[2].Bytes())
		batch.ConsolidatedTxHash = vLog.TxHash
		block.BlockNumber = vLog.BlockNumber
		block.BlockHash = vLog.BlockHash
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return nil, fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", block.BlockNumber, err)
		}
		block.ParentHash = fullBlock.ParentHash()
		block.Batches = append(block.Batches, batch)
		return &block, nil
	case newSequencerSignatureHash:
		seq, err := etherMan.PoE.ParseRegisterSequencer(vLog)
		if err != nil {
			return nil, err
		}
		var block state.Block
		var sequencer state.Sequencer
		sequencer.Address = seq.SequencerAddress
		sequencer.URL = seq.SequencerURL
		block.BlockHash = vLog.BlockHash
		block.BlockNumber = vLog.BlockNumber
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return nil, fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", block.BlockNumber, err)
		}
		block.ParentHash = fullBlock.ParentHash()
		sequencer.ChainID = new(big.Int).SetUint64(uint64(seq.ChainID))
		block.NewSequencers = append(block.NewSequencers, sequencer)
		return &block, nil
	case ownershipTransferredSignatureHash:
		log.Debug("Unhandled event: OwnershipTransferred: ", vLog)
		return nil, nil
	}
	log.Debug("Event not registered")
	return nil, nil
}

func decodeTxs(txsData []byte) ([]*types.Transaction, error) {
	// The first 4 bytes are the function hash bytes. These bytes has to be ripped.
	// After that, the unpack method is used to read the call data.
	// The txs data is encoded using rlp and contains encoded txs. So, decoding the txs data,
	// it is obteined an array of encoded txs. Each of these txs must to be decoded using rlp.

	// Extract coded txs.
	// Load contract ABI
	abi, err := abi.JSON(strings.NewReader(proofofefficiency.ProofofefficiencyABI))
	if err != nil {
		log.Fatal("error reading smart contract abi: ", err)
	}

	// Recover Method from signature and ABI
	method, err := abi.MethodById(txsData[:4])
	if err != nil {
		log.Fatal("error getting abi method: ", err)
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txsData[4:])
	if err != nil {
		log.Fatal("error reading call data: ", err)
	}

	txsData = data[0].([]byte)

	// Decode array of txs
	var codedTxs [][]byte
	err = rlp.DecodeBytes(txsData, &codedTxs)
	if err != nil {
		log.Debug("error decoding tx bytes: ", err, ". Data: ", hex.EncodeToString(txsData))
		return nil, err
	}

	// Process coded txs
	var txs []*types.Transaction
	for _, codedTx := range codedTxs {
		// Decode tx
		var tx types.LegacyTx
		err = rlp.DecodeBytes(codedTx, &tx)
		if err != nil {
			log.Debug("error decoding tx bytes: ", err, data)
			continue
		}
		txs = append(txs, types.NewTx(&tx))
	}
	return txs, nil
}

// GetAddress function allows to retrieve the wallet address
func (etherMan *ClientEtherMan) GetAddress() common.Address {
	return etherMan.auth.From
}

// GetDefaultChainID function allows to retrieve the default chainID from the smc
func (etherMan *ClientEtherMan) GetDefaultChainID() (*big.Int, error) {
	defaulChainID, err := etherMan.PoE.DEFAULTCHAINID(&bind.CallOpts{Pending: false})
	return new(big.Int).SetUint64(uint64(defaulChainID)), err
}

// EstimateSendBatchCost function estimate gas cost for sending batch to ethereum sc
func (etherMan *ClientEtherMan) EstimateSendBatchCost(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*big.Int, error) {
	noSendOpts := etherMan.auth
	noSendOpts.NoSend = true
	tx, err := etherMan.sendBatch(ctx, noSendOpts, txs, maticAmount)
	if err != nil {
		return nil, err
	}
	return tx.Cost(), nil
}

// GetLatestProposedBatchNumber function allows to retrieve the latest proposed batch in the smc
func (etherMan *ClientEtherMan) GetLatestProposedBatchNumber() (uint64, error) {
	latestBatch, err := etherMan.PoE.LastBatchSent(&bind.CallOpts{Pending: false})
	return uint64(latestBatch), err
}

// GetSequencerCollateral function allows to retrieve the sequencer collateral from the smc
func (etherMan *ClientEtherMan) GetSequencerCollateral(batchNumber uint64) (*big.Int, error) {
	batchInfo, err := etherMan.PoE.SentBatches(&bind.CallOpts{Pending: false}, uint32(batchNumber))
	return batchInfo.MaticCollateral, err
}
