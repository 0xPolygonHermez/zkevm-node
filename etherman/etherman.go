package etherman

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
)

var (
	newBatchEventSignatureHash    = crypto.Keccak256Hash([]byte("SendBatch(uint256,address)"))
	consolidateBatchSignatureHash = crypto.Keccak256Hash([]byte("VerifyBatch(uint256,address)"))
	newSequencerSignatureHash     = crypto.Keccak256Hash([]byte("SetSequencer(address,string)"))
)

// EtherMan represents an Ethereum Manager
type EtherMan interface {
	EthBlockByNumber(ctx context.Context, blockNum int64) (*types.Block, error)
	GetBatchesByBlock(ctx context.Context, blockNum uint64, blockHash *common.Hash) ([]state.Block, error)
	GetBatchesByBlockRange(ctx context.Context, fromBlock uint64, toBlock uint64) ([]state.Block, error)
	SendBatch(batch state.Batch) (common.Hash, error)
	ConsolidateBatch(batch state.Batch, proof state.Proof) (common.Hash, error)
}

// ClientEtherMan is a simple implementation of EtherMan
type ClientEtherMan struct {
	EtherClient *ethclient.Client
	PoE         *proofofefficiency.Proofofefficiency
	SCAddresses []common.Address

	key *keystore.Key
}

// NewEtherman creates a new etherman
func NewEtherman(cfg Config) (EtherMan, error) {
	//Connect to ethereum node
	ethClient, err := ethclient.Dial(cfg.URL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", cfg.URL, err)
		return nil, err
	}
	//Create smc clients
	poe, err := proofofefficiency.NewProofofefficiency(cfg.PoEAddress, ethClient)
	if err != nil {
		return nil, err
	}
	var scAddresses []common.Address
	scAddresses = append(scAddresses, cfg.PoEAddress)

	var key *keystore.Key
	if cfg.PrivateKeyPath != "" || cfg.PrivateKeyPassword != "" {
		key, err = decryptKeystore(cfg.PrivateKeyPath, cfg.PrivateKeyPassword)
		if err != nil {
			return nil, err
		}
	}
	return &ClientEtherMan{EtherClient: ethClient, PoE: poe, SCAddresses: scAddresses, key: key}, nil
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number
func (etherMan *ClientEtherMan) EthBlockByNumber(ctx context.Context, blockNum int64) (*types.Block, error) {
	block, err := etherMan.EtherClient.BlockByNumber(ctx, big.NewInt(blockNum))
	if err != nil {
		return &types.Block{}, err
	}
	return block, nil
}

// GetBatchesByBlock function retrieves the batches information that are included in a specific ethereum block
func (etherMan *ClientEtherMan) GetBatchesByBlock(ctx context.Context, blockNum uint64, blockHash *common.Hash) ([]state.Block, error) {
	//First filter query
	var blockNumBigInt *big.Int
	if blockHash == nil {
		blockNumBigInt = new(big.Int).SetUint64(blockNum)
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
//from block x to block y
func (etherMan *ClientEtherMan) GetBatchesByBlockRange(ctx context.Context, fromBlock uint64, toBlock uint64) ([]state.Block, error) {
	//First filter query
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: etherMan.SCAddresses,
	}
	blocks, err := etherMan.readEvents(ctx, query)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// SendBatch function allows the sequencer send a new batch proposal to the rollup
func (etherMan *ClientEtherMan) SendBatch(batch state.Batch) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}

// ConsolidateBatch function allows the agregator send the proof for a batch and consolidate it
func (etherMan *ClientEtherMan) ConsolidateBatch(batch state.Batch, proof state.Proof) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}

func (etherMan *ClientEtherMan) readEvents(ctx context.Context, query ethereum.FilterQuery) ([]state.Block, error) {
	logs, err := etherMan.EtherClient.FilterLogs(ctx, query)
	if err != nil {
		return []state.Block{}, err
	}
	blocks := make(map[common.Hash]state.Block)
	for _, vLog := range logs {
		block, err := etherMan.processEvent(ctx, vLog)
		if err != nil {
			log.Warn("error processing event: ", err, vLog)
			continue
		}
		if b, exists := blocks[block.BlockHash]; exists {
			b.Batches = append(blocks[block.BlockHash].Batches, block.Batches...)
			b.NewSequencers = append(blocks[block.BlockHash].NewSequencers, block.NewSequencers...)
			blocks[block.BlockHash] = b
		} else {
			blocks[block.BlockHash] = *block
		}
	}
	var blockArr []state.Block
	for _, b := range blocks {
		blockArr = append(blockArr, b)
	}
	return blockArr, nil
}

func (etherMan *ClientEtherMan) processEvent(ctx context.Context, vLog types.Log) (*state.Block, error) {
	switch vLog.Topics[0] {
	case newBatchEventSignatureHash:
		var block state.Block
		//Indexed parameters using topics
		var batch state.Batch
		batch.BatchNumber = new(big.Int).SetBytes(vLog.Topics[1][:]).Uint64()
		batch.Sequencer = common.BytesToAddress(vLog.Topics[2].Bytes())
		var head types.Header
		head.TxHash = vLog.TxHash
		batch.Header = &head
		block.BlockNumber = vLog.BlockNumber
		batch.BlockNumber = vLog.BlockNumber
		block.BlockHash = vLog.BlockHash
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err == nil {
			block.ParentHash = fullBlock.ParentHash()
		}
		//Read the tx for this batch.
		tx, isPending, err := etherMan.EtherClient.TransactionByHash(ctx, batch.Header.TxHash)
		if err != nil {
			return nil, err
		} else if isPending {
			return nil, fmt.Errorf("error: tx is still pending. TxHash: %s", tx.Hash().String())
		}
		batch.RawTxsData = tx.Data()
		txs, err := decodeTxs(tx.Data())
		if err != nil {
			return nil, err
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
		if err == nil {
			block.ParentHash = fullBlock.ParentHash()
		}
		block.Batches = append(block.Batches, batch)
		return &block, nil
	case newSequencerSignatureHash:
		seq, err := etherMan.PoE.ParseSetSequencer(vLog)
		if err != nil {
			return nil, err
		}
		var block state.Block
		var sequencer state.Sequencer
		sequencer.Sequencer = seq.SequencerAddress
		sequencer.URL = seq.SequencerURL
		block.BlockHash = vLog.BlockHash
		block.BlockNumber = vLog.BlockNumber
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err == nil {
			block.ParentHash = fullBlock.ParentHash()
		}
		//Get sequencer chainId
		se, err := etherMan.PoE.Sequencers(&bind.CallOpts{Pending: false}, seq.SequencerAddress)
		if err != nil {
			return nil, err
		}
		sequencer.ChainID = se.ChainID
		block.NewSequencers = append(block.NewSequencers, sequencer)
		return &block, nil
	}
	return nil, fmt.Errorf("Event not registered")
}

const headerByteLength = 2

func decodeTxs(txsData []byte) ([]*types.Transaction, error) {
	// First split txs
	// The first two bytes are the header. The information related to the length of the tx is stored in the second byte.
	// So, first read the second byte to check the tx length. Then, copy from the current position to the last
	// byte of the tx if exists (if not will be completed with zeros). Now, I try to decode the tx, If it is possible,
	// everything is fine. If not, print error and try to get the next tx.

	//Extract coded txs.
	// load contract ABI
	abi, err := abi.JSON(strings.NewReader(proofofefficiency.ProofofefficiencyABI))
	if err != nil {
		log.Fatal(err)
	}

	// recover Method from signature and ABI
	method, err := abi.MethodById(txsData[:4])
	if err != nil {
		log.Fatal(err)
	}

	// unpack method inputs
	data, err := method.Inputs.Unpack(txsData[4:])
	if err != nil {
		log.Fatal(err)
	}

	txsData = data[0].([]byte)

	//Process coded txs
	var pos int64
	var txs []*types.Transaction
	for pos < int64(len(txsData)) {
		length := txsData[pos+1 : pos+2]
		str := hex.EncodeToString(length)
		num, err := strconv.ParseInt(str, 16, 64)
		if err != nil {
			log.Warn("error: skipping tx. Err: ", err)
			continue
		}

		data := txsData[pos : pos+num+2]
		pos = pos + num + headerByteLength

		//Decode tx
		var tx types.LegacyTx
		err = rlp.DecodeBytes(data, &tx)
		if err != nil {
			log.Info("error decoding tx bytes: ", err, data)
			continue
		}
		txs = append(txs, types.NewTransaction(tx.Nonce, *tx.To, tx.Value, tx.Gas, tx.GasPrice, tx.Data))
	}
	return txs, nil
}
