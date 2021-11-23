package etherman

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	newBatcheventSignatureHash = crypto.Keccak256Hash([]byte("SendBatch(uint256,address)"))
	consolidateBatchSignatureHash = crypto.Keccak256Hash([]byte("VerifyBatch(uint256,aggregator)"))
)

type EtherMan struct {
	EtherClient *ethclient.Client
	PoE         *proofofefficiency.Proofofefficiency
	SmcAddreses []common.Address
}

func NewEtherman(url string, poeAddr common.Address) (*EtherMan, error) {
	//Connect to ethereum node
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", url, err)
		return nil, err
	}
	//Create smc clients
	poe, err := proofofefficiency.NewProofofefficiency(poeAddr, ethClient)
	if err != nil {
		return nil, err
	}
	var smcAddreses []common.Address
	smcAddreses = append(smcAddreses, poeAddr)
	return &EtherMan{EtherClient: ethClient, PoE: poe, SmcAddreses: smcAddreses}, nil
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number
func (etherMan *EtherMan) EthBlockByNumber(ctx context.Context, blockNum int64) (*types.Block, error) {
	block, err := etherMan.EtherClient.BlockByNumber(ctx, big.NewInt(blockNum))
	if err != nil {
		return &types.Block{}, nil
	}
	return block, nil
}

// GetBatchesByBlock function retrieves the batches information that are included in a specific ethereum block
func (etherMan *EtherMan) GetBatchesByBlock(blockNum uint64, blockHash *common.Hash) ([]state.Block, error) {
	//First filter query
	var blockNumBigInt *big.Int
	if blockHash == nil {
		blockNumBigInt = big.NewInt(int64(blockNum))
	}
	query := ethereum.FilterQuery{
		BlockHash: blockHash,
		FromBlock: blockNumBigInt,
		ToBlock:   blockNumBigInt,
		Addresses: etherMan.SmcAddreses,
	}
	blocks, err := etherMan.readEvents(query)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// GetBatchesFromBlockTo function retrieves the batches information that are included in all this ethereum blocks
//from block x to block y
func (etherMan *EtherMan) GetBatchesFromBlockTo(fromBlock uint64, toBlock uint64) ([]state.Block, error) {
	//First filter query
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: etherMan.SmcAddreses,
	}
	blocks, err := etherMan.readEvents(query)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// SendBatch function allows the sequencer send a new batch proposal to the rollup
func (etherMan *EtherMan) SendBatch(batch state.Batch) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}

// ConsolidateBatch function allows the agregator send the proof for a batch and consolidate it
func (etherMan *EtherMan) ConsolidateBatch(batch state.Batch, proof state.Proof) (common.Hash, error) {
	//TODO
	return common.Hash{}, nil
}

func (etherMan *EtherMan) readEvents(query ethereum.FilterQuery) ([]state.Block, error) {
	logs, err := etherMan.EtherClient.FilterLogs(context.Background(), query)
	if err != nil {
		return []state.Block{}, err
	}
	var blocks []state.Block
	for _, vLog := range logs {
		block, err := etherMan.processEvent(vLog)
		if err != nil {
			log.Warn("error processing event: ", err, vLog)
			continue
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (etherMan *EtherMan) processEvent(vLog types.Log) (state.Block, error) {
	switch vLog.Topics[0] {
	case newBatcheventSignatureHash:
		var block state.Block
		//Indexed parameters using topics
		var batch state.Batch
		batch.BatchNumber = new(big.Int).SetBytes(vLog.Topics[1][:]).Uint64()
		batch.Sequencer = common.BytesToAddress(vLog.Topics[2].Bytes())
		batch.Header.TxHash = vLog.TxHash
		block.BlockNum = vLog.BlockNumber
		block.BlockHash = vLog.BlockHash
		//Now, We have to read the tx for this batch.
		tx, isPending, err := etherMan.EtherClient.TransactionByHash(context.Background(), batch.Header.TxHash)
		if err != nil || isPending {
			return state.Block{}, err
		}
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
		block.Batches = append(block.Batches, batch)
		return block, nil
	}
	return state.Block{}, nil
}

func decodeTxs(txsData []byte, chainId *big.Int) ([]*types.LegacyTx, error) {
	// First split txs
	var pos int64
	var txs []*types.LegacyTx
	for pos<int64(len(txsData)) {
		lenght := txsData[pos+1:pos+2]
		str := hex.EncodeToString(lenght)
		num, err := strconv.ParseInt(str, 16, 64)
		if err != nil {
			log.Warn("error: skipping tx. Err: ", err)
			continue
		}
		data := txsData[pos : pos+num+2]
		pos = pos + num + 2

		//Decode tx
		var tx types.LegacyTx
		rlp.DecodeBytes(data, &tx)
		isValid, err := checkSignature(tx, chainId)
		if err != nil {
			log.Warn("error: skipping tx. ", err)
			continue
		} else if !isValid {
			log.Debug("Signature invalid: ",isValid)
			continue
		}
		txs = append(txs, &tx)
	}
    return txs, nil
}

func checkSignature(tx types.LegacyTx, chainId *big.Int) (bool, error) {
	decodedChainId := deriveChainId(tx.V)
	if decodedChainId.Cmp(chainId) != 0 {
		return false, fmt.Errorf("error: incorrect chainId. Decoded chainId: %d and chainId retrieved from smc: %d",
		decodedChainId, chainId)
	}
	// Formula: v = chainId * 2 + 36 or 35; x = 35 or 36
	v := new(big.Int).SetBytes(tx.V.Bytes())
	r := new(big.Int).SetBytes(tx.R.Bytes())
	s := new(big.Int).SetBytes(tx.S.Bytes())
	x := v.Int64()-(chainId.Int64()*2)
	var vField byte
	if x == 35 {
		vField = byte(0)
	} else if x == 36 {
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

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}