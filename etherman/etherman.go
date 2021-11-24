package etherman

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
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
	consolidateBatchSignatureHash = crypto.Keccak256Hash([]byte("VerifyBatch(uint256,aggregator)"))
)

// EtherMan represents an Ethereum Manager
type EtherMan interface {
	EthBlockByNumber(ctx context.Context, blockNum int64) (*types.Block, error)
	GetBatchesByBlock(ctx context.Context, blockNum uint64, blockHash *common.Hash) ([]state.Block, error)
	GetBatchesFromBlockTo(ctx context.Context, fromBlock uint64, toBlock uint64) ([]state.Block, error)
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
	poe, err := proofofefficiency.NewProofofefficiency(cfg.PoeAddress, ethClient)
	if err != nil {
		return nil, err
	}
	var scAddresses []common.Address
	scAddresses = append(scAddresses, cfg.PoeAddress)

	var key *keystore.Key
	if cfg.PrivateKeyPath != "" || cfg.PrivateKeyPassword != "" {
		key, err = decryptKeystore(cfg.PrivateKeyPath, cfg.PrivateKeyPassword)
		if err != nil {
			return nil, err
		}
	}
	return &ClientEtherMan{EtherClient: ethClient, PoE: poe, SCAddresses: scAddresses, key: key}, nil
}

// TestClientEtherMan is a simple implementation of EtherMan
type TestClientEtherMan struct {
	EtherClient *backends.SimulatedBackend
	PoE         *proofofefficiency.Proofofefficiency
	SCAddresses []common.Address

	key *keystore.Key
}

// NewTestEtherman creates a new test etherman
func NewTestEtherman(cfg Config, etherCLient *backends.SimulatedBackend, poe *proofofefficiency.Proofofefficiency) (EtherMan, error) {
	var scAddresses []common.Address
	scAddresses = append(scAddresses, cfg.PoeAddress)

	var (
		key *keystore.Key
		err error
	)
	if cfg.PrivateKeyPath != "" || cfg.PrivateKeyPassword != "" {
		key, err = decryptKeystore(cfg.PrivateKeyPath, cfg.PrivateKeyPassword)
		if err != nil {
			return nil, err
		}
	}
	return &TestClientEtherMan{EtherClient: etherCLient, PoE: poe, SCAddresses: scAddresses, key: key}, nil
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number
func (etherMan *ClientEtherMan) EthBlockByNumber(ctx context.Context, blockNum int64) (*types.Block, error) {
	block, err := etherMan.EtherClient.BlockByNumber(ctx, big.NewInt(blockNum))
	if err != nil {
		return &types.Block{}, nil
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

// GetBatchesFromBlockTo function retrieves the batches information that are included in all this ethereum blocks
//from block x to block y
func (etherMan *ClientEtherMan) GetBatchesFromBlockTo(ctx context.Context, fromBlock uint64, toBlock uint64) ([]state.Block, error) {
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
	var blocks []state.Block
	for _, vLog := range logs {
		block, err := etherMan.processEvent(ctx, vLog)
		if err != nil {
			log.Warn("error processing event: ", err, vLog)
			continue
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (etherMan *ClientEtherMan) processEvent(ctx context.Context, vLog types.Log) (state.Block, error) {
	switch vLog.Topics[0] {
	case newBatchEventSignatureHash:
		var block state.Block
		//Indexed parameters using topics
		var batch state.Batch
		batch.BatchNumber = new(big.Int).SetBytes(vLog.Topics[1][:]).Uint64()
		batch.Sequencer = common.BytesToAddress(vLog.Topics[2].Bytes())
		batch.Header.TxHash = vLog.TxHash
		block.BlockNum = vLog.BlockNumber
		block.BlockHash = vLog.BlockHash
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err == nil {
			block.ParentHash = fullBlock.ParentHash()
		}
		//Now, We have to read the tx for this batch.
		tx, isPending, err := etherMan.EtherClient.TransactionByHash(ctx, batch.Header.TxHash)
		if err != nil || isPending {
			return state.Block{}, err
		}
		batch.RawTxsData = tx.Data()
		//Get sequencer chainId
		seq, err := etherMan.PoE.Sequencers(&bind.CallOpts{Pending: false}, batch.Sequencer)
		if err != nil {
			return state.Block{}, err
		}
		txs, err := decodeTxs(tx.Data(), seq.ChainID)
		if err != nil {
			return state.Block{}, err
		}
		batch.Transactions = txs
		block.Batches = append(block.Batches, batch)
		return block, nil
	case consolidateBatchSignatureHash:
		var block state.Block
		var batch state.Batch
		batch.BatchNumber = new(big.Int).SetBytes(vLog.Topics[1][:]).Uint64()
		batch.Aggregator = common.BytesToAddress(vLog.Topics[2].Bytes())
		batch.ConsolidatedTxHash = vLog.TxHash
		block.BlockNum = vLog.BlockNumber
		block.BlockHash = vLog.BlockHash
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err == nil {
			block.ParentHash = fullBlock.ParentHash()
		}
		block.Batches = append(block.Batches, batch)
		return block, nil
	}
	return state.Block{}, nil
}

const (
	formulaDiv       = 2
	fomulaConst      = 35
	fomulaConst2     = 36
	maxVLength       = 64
	headerByteLength = 2
)

func decodeTxs(txsData []byte, chainID *big.Int) ([]*types.Transaction, error) {
	// First split txs
	// The first two bytes are the header. The information related to the length of the tx is stored in the second byte.
	// So, first we've to read the second byte to check the tx length. Then, copy from the current position to the last
	// byte of the tx if exists (if not will be completed with zeros). Now, I try to decode the tx, If it is possible,
	// everything is fine. If not, print error and try to get the next tx.
	var pos int64 = 100
	var txs []*types.Transaction
	for pos < int64(len(txsData)) {
		length := txsData[pos+1 : pos+2]
		str := hex.EncodeToString(length)
		num, err := strconv.ParseInt(str, 16, 64)
		if err != nil {
			log.Warn("error: skipping tx. Err: ", err)
			continue
		}
		if num == 0 {
			break
		}

		data := txsData[pos : pos+num+2]
		pos = pos + num + headerByteLength

		//Decode tx
		var tx types.LegacyTx
		err = rlp.DecodeBytes(data, &tx)
		if err != nil {
			log.Error("error decoding tx bytes: ", err, data)
			continue
		}
		isValid, err := checkSignature(tx, chainID)
		if err != nil {
			log.Warn("error: skipping tx. ", err, tx)
			continue
		} else if !isValid {
			log.Debug("Signature invalid for tx: ", tx)
			continue
		}
		txs = append(txs, types.NewTransaction(tx.Nonce, *tx.To, tx.Value, tx.Gas, tx.GasPrice, tx.Data))
	}
	return txs, nil
}

func checkSignature(tx types.LegacyTx, chainID *big.Int) (bool, error) {
	decodedChainID := deriveChainID(tx.V)
	if decodedChainID.Cmp(chainID) != 0 {
		return false, fmt.Errorf("error: incorrect chainId. Decoded chainId: %d and chainId retrieved from smc: %d",
			decodedChainID, chainID)
	}
	// Formula: v = chainId * 2 + 36 or 35; x = 35 or 36
	v := new(big.Int).SetBytes(tx.V.Bytes())
	r := new(big.Int).SetBytes(tx.R.Bytes())
	s := new(big.Int).SetBytes(tx.S.Bytes())
	x := v.Int64() - (chainID.Int64() * 2)
	var vField byte
	if x == fomulaConst {
		vField = byte(0)
	} else if x == fomulaConst2 {
		vField = byte(1)
	} else {
		return false, fmt.Errorf("Error invalid signature v value: %d", tx.V)
	}
	if !crypto.ValidateSignatureValues(vField, r, s, false) {
		log.Warn("Invalid Signature values: ", tx)
		return false, nil
	}
	return true, nil
}

// deriveChainID derives the chain id from the given v parameter
func deriveChainID(v *big.Int) *big.Int {
	if v.BitLen() <= maxVLength {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - fomulaConst) / formulaDiv)
	}
	v = new(big.Int).Sub(v, big.NewInt(fomulaConst))
	return v.Div(v, big.NewInt(formulaDiv))
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number
func (etherMan *TestClientEtherMan) EthBlockByNumber(ctx context.Context, blockNum int64) (*types.Block, error) {
	block, err := etherMan.EtherClient.BlockByNumber(ctx, big.NewInt(blockNum))
	if err != nil {
		return &types.Block{}, nil
	}
	return block, nil
}

// GetBatchesByBlock function retrieves the batches information that are included in a specific ethereum block
func (etherMan *TestClientEtherMan) GetBatchesByBlock(ctx context.Context, blockNum uint64, blockHash *common.Hash) ([]state.Block, error) {
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

// GetBatchesFromBlockTo function retrieves the batches information that are included in all this ethereum blocks
//from block x to block y
func (etherMan *TestClientEtherMan) GetBatchesFromBlockTo(ctx context.Context, fromBlock uint64, toBlock uint64) ([]state.Block, error) {
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
func (etherMan *TestClientEtherMan) SendBatch(batch state.Batch) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}

// ConsolidateBatch function allows the agregator send the proof for a batch and consolidate it
func (etherMan *TestClientEtherMan) ConsolidateBatch(batch state.Batch, proof state.Proof) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}

func (etherMan *TestClientEtherMan) readEvents(ctx context.Context, query ethereum.FilterQuery) ([]state.Block, error) {
	logs, err := etherMan.EtherClient.FilterLogs(ctx, query)
	if err != nil {
		return []state.Block{}, err
	}
	var blocks []state.Block
	for _, vLog := range logs {
		block, err := etherMan.processEvent(ctx, vLog)
		if err != nil {
			log.Warnf("error processing event: %s %+v", err, vLog)
			continue
		}
		if block != nil {
			blocks = append(blocks, *block)
		}
	}
	return blocks, nil
}

func (etherMan *TestClientEtherMan) processEvent(ctx context.Context, vLog types.Log) (*state.Block, error) {
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
		block.BlockNum = vLog.BlockNumber
		block.BlockHash = vLog.BlockHash
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err == nil {
			block.ParentHash = fullBlock.ParentHash()
		}
		//Now, We have to read the tx for this batch.
		tx, isPending, err := etherMan.EtherClient.TransactionByHash(ctx, batch.Header.TxHash)
		if err != nil || isPending {
			return nil, err
		}
		batch.RawTxsData = tx.Data()
		txs, err := decodeTxs(tx.Data(), big.NewInt(1))
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
		block.BlockNum = vLog.BlockNumber
		block.BlockHash = vLog.BlockHash
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err == nil {
			block.ParentHash = fullBlock.ParentHash()
		}
		block.Batches = append(block.Batches, batch)
		return &block, nil
	}
	return nil, fmt.Errorf("Event not registered")
}
