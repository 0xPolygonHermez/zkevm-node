package etherman

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
)

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
	scAddresses = append(scAddresses, cfg.PoEAddress)

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

// GetBatchesByBlockRange function retrieves the batches information that are included in all this ethereum blocks
//from block x to block y
func (etherMan *TestClientEtherMan) GetBatchesByBlockRange(ctx context.Context, fromBlock uint64, toBlock uint64) ([]state.Block, error) {
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
		block.BlockNumber = vLog.BlockNumber
		batch.BlockNumber = vLog.BlockNumber
		block.BlockHash = vLog.BlockHash
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err == nil {
			block.ParentHash = fullBlock.ParentHash()
		}
		if err != nil {
			log.Warn("error getting hashParent. BlockNumber: ", block.BlockNumber, " error: ", err)
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
		if err != nil {
			log.Warn("error getting hashParent. BlockNumber: ", block.BlockNumber, " error: ", err)
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
		sequencer.Address = seq.SequencerAddress
		sequencer.URL = seq.SequencerURL
		block.BlockHash = vLog.BlockHash
		block.BlockNumber = vLog.BlockNumber
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err == nil {
			block.ParentHash = fullBlock.ParentHash()
		}
		if err != nil {
			log.Warn("error getting hashParent. BlockNumber: ", block.BlockNumber, " error: ", err)
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
