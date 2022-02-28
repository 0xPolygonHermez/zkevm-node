package etherman

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/matic"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/state"
)

var (
	newBatchEventSignatureHash             = crypto.Keccak256Hash([]byte("SendBatch(uint32,address,uint32,bytes32)"))
	consolidateBatchSignatureHash          = crypto.Keccak256Hash([]byte("VerifyBatch(uint32,address)"))
	newSequencerSignatureHash              = crypto.Keccak256Hash([]byte("RegisterSequencer(address,string,uint32)"))
	ownershipTransferredSignatureHash      = crypto.Keccak256Hash([]byte("OwnershipTransferred(address,address)"))

	// ErrNotFound is used when the object is not found
	ErrNotFound = errors.New("Not found")
)

// EventOrder is the the type used to identify the events order
type EventOrder string

const (
	// BatchesOrder identifies a batch event
	BatchesOrder EventOrder = "Batches"
	//NewSequencersOrder identifies a newSequencer event
	NewSequencersOrder EventOrder = "NewSequencers"
)

type ethClienter interface {
	ethereum.ChainReader
	ethereum.LogFilterer
	ethereum.TransactionReader
}

// Client is a simple implementation of EtherMan.
type Client struct {
	EtherClient           ethClienter
	PoE                   *proofofefficiency.Proofofefficiency
	Matic                 *matic.Matic
	SCAddresses           []common.Address

	auth *bind.TransactOpts
}

// NewClient creates a new etherman.
func NewClient(cfg Config, auth *bind.TransactOpts, PoEAddr common.Address, maticAddr common.Address) (*Client, error) {
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
	matic, err := matic.NewMatic(maticAddr, ethClient)
	if err != nil {
		return nil, err
	}
	var scAddresses []common.Address
	scAddresses = append(scAddresses, PoEAddr)

	return &Client{EtherClient: ethClient, PoE: poe, Matic: matic, SCAddresses: scAddresses, auth: auth}, nil
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number.
func (etherMan *Client) EthBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	block, err := etherMan.EtherClient.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		if errors.Is(err, ethereum.NotFound) || err.Error() == "block does not exist in blockchain" {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return block, nil
}

// GetRollupInfoByBlock function retrieves the Rollup information that are included in a specific ethereum block.
func (etherMan *Client) GetRollupInfoByBlock(ctx context.Context, blockNumber uint64, blockHash *common.Hash) ([]state.Block, map[common.Hash][]Order, error) {
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
	blocks, order, err := etherMan.readEvents(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	return blocks, order, nil
}

// GetRollupInfoByBlockRange function retrieves the Rollup information that are included in all this ethereum blocks
// from block x to block y.
func (etherMan *Client) GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]state.Block, map[common.Hash][]Order, error) {
	// First filter query
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		Addresses: etherMan.SCAddresses,
	}
	if toBlock != nil {
		query.ToBlock = new(big.Int).SetUint64(*toBlock)
	}
	blocks, order, err := etherMan.readEvents(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	return blocks, order, nil
}

// SendBatch function allows the sequencer send a new batch proposal to the rollup.
func (etherMan *Client) SendBatch(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*types.Transaction, error) {
	return etherMan.sendBatch(ctx, etherMan.auth, txs, maticAmount)
}

const (
	ether155V = 27
)

func (etherMan *Client) sendBatch(ctx context.Context, opts *bind.TransactOpts, txs []*types.Transaction, maticAmount *big.Int) (*types.Transaction, error) {
	if len(txs) == 0 {
		return nil, errors.New("invalid txs: is empty slice")
	}
	var callDataHex string
	for _, tx := range txs {
		v, r, s := tx.RawSignatureValues()
		sign := 1 - (v.Uint64() & 1)

		txCodedRlp, err := rlp.EncodeToBytes([]interface{}{
			tx.Nonce(),
			tx.GasPrice(),
			tx.Gas(),
			tx.To(),
			tx.Value(),
			tx.Data(),
			tx.ChainId(), uint(0), uint(0),
		})
		if err != nil {
			log.Error("error encoding rlp tx: ", err)
			return nil, errors.New("error encoding rlp tx: " + err.Error())
		}
		newV := new(big.Int).Add(big.NewInt(ether155V), big.NewInt(int64(sign)))
		newRPadded := fmt.Sprintf("%064s", r.Text(hex.Base))
		newSPadded := fmt.Sprintf("%064s", s.Text(hex.Base))
		newVPadded := fmt.Sprintf("%02s", newV.Text(hex.Base))
		callDataHex = callDataHex + hex.EncodeToString(txCodedRlp) + newRPadded + newSPadded + newVPadded
	}
	callData, err := hex.DecodeString(callDataHex)
	if err != nil {
		log.Error("error coverting hex string to []byte. Error: ", err)
		return nil, errors.New("error coverting hex string to []byte. Error: " + err.Error())
	}
	tx, err := etherMan.PoE.SendBatch(etherMan.auth, callData, maticAmount)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// ConsolidateBatch function allows the aggregator send the proof for a batch and consolidate it
func (etherMan *Client) ConsolidateBatch(batchNumber *big.Int, proof *proverclient.Proof) (*types.Transaction, error) {
	publicInputs := proof.PublicInputsExtended.PublicInputs
	newLocalExitRoot, err := stringToFixedByteArray(publicInputs.NewLocalExitRoot)
	if err != nil {
		return nil, err
	}
	newStateRoot, err := stringToFixedByteArray(publicInputs.NewStateRoot)
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

// RegisterSequencer function allows to register a new sequencer in the rollup.
func (etherMan *Client) RegisterSequencer(url string) (*types.Transaction, error) {
	tx, err := etherMan.PoE.RegisterSequencer(etherMan.auth, url)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// Order contains the event order to let the synchronizer store the information following this order.
type Order struct {
	Name EventOrder
	Pos  int
}

func (etherMan *Client) readEvents(ctx context.Context, query ethereum.FilterQuery) ([]state.Block, map[common.Hash][]Order, error) {
	logs, err := etherMan.EtherClient.FilterLogs(ctx, query)
	if err != nil {
		return []state.Block{}, nil, err
	}
	blockOrder := make(map[common.Hash][]Order)
	blocks := make(map[common.Hash]state.Block)
	var blockKeys []common.Hash
	for _, vLog := range logs {
		block, err := etherMan.processEvent(ctx, vLog)
		if err != nil {
			log.Warnf("error processing event. Retrying... Error: %s. vLog: %+v", err.Error(), vLog)
			break
		}
		if block == nil {
			continue
		}
		if b, exists := blocks[block.BlockHash]; exists {
			if len(block.Batches) != 0 {
				b.Batches = append(blocks[block.BlockHash].Batches, block.Batches...)
				or := Order{
					Name: BatchesOrder,
					Pos:  len(b.Batches) - 1,
				}
				blockOrder[b.BlockHash] = append(blockOrder[b.BlockHash], or)
			}
			if len(block.NewSequencers) != 0 {
				b.NewSequencers = append(blocks[block.BlockHash].NewSequencers, block.NewSequencers...)
				or := Order{
					Name: NewSequencersOrder,
					Pos:  len(b.NewSequencers) - 1,
				}
				blockOrder[b.BlockHash] = append(blockOrder[b.BlockHash], or)
			}
			blocks[block.BlockHash] = b
		} else {
			if len(block.Batches) != 0 {
				or := Order{
					Name: BatchesOrder,
					Pos:  len(block.Batches) - 1,
				}
				blockOrder[block.BlockHash] = append(blockOrder[block.BlockHash], or)
			}
			if len(block.NewSequencers) != 0 {
				or := Order{
					Name: NewSequencersOrder,
					Pos:  len(block.NewSequencers) - 1,
				}
				blockOrder[block.BlockHash] = append(blockOrder[block.BlockHash], or)
			}
			blocks[block.BlockHash] = *block
			blockKeys = append(blockKeys, block.BlockHash)
		}
	}
	var blockArr []state.Block
	for _, hash := range blockKeys {
		blockArr = append(blockArr, blocks[hash])
	}
	return blockArr, blockOrder, nil
}

func (etherMan *Client) processEvent(ctx context.Context, vLog types.Log) (*state.Block, error) {
	switch vLog.Topics[0] {
	case newBatchEventSignatureHash:
		batchEvent, err := etherMan.PoE.ParseSendBatch(vLog)
		if err != nil {
			return nil, err
		}
		// Indexed parameters using topics
		var head types.Header
		head.TxHash = vLog.TxHash
		head.Difficulty = big.NewInt(0)
		head.Number = new(big.Int).SetUint64(uint64(batchEvent.NumBatch))

		var batch state.Batch
		batch.Sequencer = batchEvent.Sequencer
		batch.ChainID = new(big.Int).SetUint64(uint64(batchEvent.BatchChainID))
		batch.GlobalExitRoot = batchEvent.LastGlobalExitRoot
		batch.Header = &head
		batch.BlockNumber = vLog.BlockNumber
		maticCollateral, err := etherMan.GetSequencerCollateralByBatchNumber(batch.Number().Uint64())
		if err != nil {
			return nil, fmt.Errorf("error getting matic collateral for batch: %d. BlockNumber: %d. Error: %w", batch.Number().Uint64(), vLog.BlockNumber, err)
		}
		batch.MaticCollateral = maticCollateral
		// Read the tx for this batch.
		tx, isPending, err := etherMan.EtherClient.TransactionByHash(ctx, batch.Header.TxHash)
		if err != nil {
			return nil, err
		} else if isPending {
			return nil, fmt.Errorf("error: tx is still pending. TxHash: %s", tx.Hash().String())
		}
		txs, rawTxs, err := decodeTxs(tx.Data())
		batch.RawTxsData = rawTxs
		if err != nil {
			log.Warn("No txs decoded in batch: ", batch.Number(), ". This batch is inside block: ", batch.BlockNumber,
				". Error: ", err)
		}
		batch.Transactions = txs
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return nil, fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
		}
		batch.ReceivedAt = fullBlock.ReceivedAt

		var block state.Block
		block.BlockNumber = vLog.BlockNumber
		block.BlockHash = vLog.BlockHash
		block.ParentHash = fullBlock.ParentHash()
		block.ReceivedAt = fullBlock.ReceivedAt
		block.Batches = append(block.Batches, batch)
		return &block, nil
	case consolidateBatchSignatureHash:
		var head types.Header
		head.Number = new(big.Int).SetBytes(vLog.Topics[1][:])

		var batch state.Batch
		batch.Header = &head
		batch.BlockNumber = vLog.BlockNumber
		batch.Aggregator = common.BytesToAddress(vLog.Topics[2].Bytes())
		batch.ConsolidatedTxHash = vLog.TxHash
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return nil, fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", vLog.BlockNumber, err)
		}
		batch.ConsolidatedAt = &fullBlock.ReceivedAt

		var block state.Block
		block.BlockNumber = vLog.BlockNumber
		block.BlockHash = vLog.BlockHash
		block.ParentHash = fullBlock.ParentHash()
		block.ReceivedAt = fullBlock.ReceivedAt
		block.Batches = append(block.Batches, batch)

		log.Debug("Consolidated tx hash: ", vLog.TxHash, batch.ConsolidatedTxHash)
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
		sequencer.BlockNumber = vLog.BlockNumber
		fullBlock, err := etherMan.EtherClient.BlockByHash(ctx, vLog.BlockHash)
		if err != nil {
			return nil, fmt.Errorf("error getting hashParent. BlockNumber: %d. Error: %w", block.BlockNumber, err)
		}
		block.ParentHash = fullBlock.ParentHash()
		block.ReceivedAt = fullBlock.ReceivedAt
		sequencer.ChainID = new(big.Int).SetUint64(uint64(seq.ChainID))
		block.NewSequencers = append(block.NewSequencers, sequencer)
		return &block, nil
	case ownershipTransferredSignatureHash:
		ownership, err := etherMan.PoE.ParseOwnershipTransferred(vLog)
		if err != nil {
			return nil, err
		}
		emptyAddr := common.Address{}
		if ownership.PreviousOwner == emptyAddr {
			log.Debug("New rollup smc deployment detected. Deployment account: ", ownership.NewOwner)
		} else {
			log.Debug("Rollup smc OwnershipTransferred from account ", ownership.PreviousOwner, " to ", ownership.NewOwner)
		}
		return nil, nil
	}
	log.Warn("Event not registered: ", vLog)
	return nil, nil
}

func decodeTxs(txsData []byte) ([]*types.Transaction, []byte, error) {
	// The first 4 bytes are the function hash bytes. These bytes has to be ripped.
	// After that, the unpack method is used to read the call data.
	// The txs data is a chunk of concatenated rawTx. This rawTx is the encoded tx information in rlp + the signature information (v, r, s).
	//So, txs data will look like: txRLP+r+s+v+txRLP2+r2+s2+v2

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

	// Process coded txs
	var pos int64
	var txs []*types.Transaction
	const (
		headerByteLength = 1
		sLength          = 32
		rLength          = 32
		vLength          = 1
		c0               = 192 // 192 is c0. This value is defined by the rlp protocol
		ff               = 255 // max value of rlp header
		shortRlp         = 55  // length of the short rlp codification
		f7               = 247 // 192 + 55 = c0 + shortRlp
		etherNewV        = 35
		mul2             = 2
	)
	txDataLength := len(txsData)
	for pos < int64(txDataLength) {
		num, err := strconv.ParseInt(hex.EncodeToString(txsData[pos:pos+1]), hex.Base, encoding.BitSize64)
		if err != nil {
			log.Debug("error parsing header length: ", err)
			return []*types.Transaction{}, []byte{}, err
		}
		// First byte is the length and must be ignored
		len := num - c0
		if len > shortRlp { // If rlp is bigger than length 55
			// n is the length of the rlp data without the header (1 byte) for example "0xf7"
			n, err := strconv.ParseInt(hex.EncodeToString(txsData[pos+1:pos+1+num-f7]), hex.Base, encoding.BitSize64) // +1 is the header. For example 0xf7
			if err != nil {
				log.Debug("error parsing length: ", err)
				return []*types.Transaction{}, []byte{}, err
			}
			len = n + num - f7 // num - f7 is the header. For example 0xf7
		}

		fullDataTx := txsData[pos : pos+len+rLength+sLength+vLength+headerByteLength]
		txInfo := txsData[pos : pos+len+headerByteLength]
		r := txsData[pos+len+headerByteLength : pos+len+rLength+headerByteLength]
		s := txsData[pos+len+rLength+headerByteLength : pos+len+rLength+sLength+headerByteLength]
		v := txsData[pos+len+rLength+sLength+headerByteLength : pos+len+rLength+sLength+vLength+headerByteLength]

		pos = pos + len + rLength + sLength + vLength + headerByteLength

		// Decode tx
		var tx types.LegacyTx
		err = rlp.DecodeBytes(txInfo, &tx)
		if err != nil {
			log.Debug("error decoding tx bytes: ", err, ". fullDataTx: ", hex.EncodeToString(fullDataTx), "\n tx: ", hex.EncodeToString(txInfo), "\n Txs received: ", hex.EncodeToString(txsData))
			return []*types.Transaction{}, []byte{}, err
		}

		//tx.V = v-27+chainId*2+35
		tx.V = new(big.Int).Add(new(big.Int).Sub(new(big.Int).SetBytes(v), big.NewInt(ether155V)), new(big.Int).Add(new(big.Int).Mul(tx.V, big.NewInt(mul2)), big.NewInt(etherNewV)))
		tx.R = new(big.Int).SetBytes(r)
		tx.S = new(big.Int).SetBytes(s)

		txs = append(txs, types.NewTx(&tx))
	}
	return txs, txsData, nil
}

// GetAddress function allows to retrieve the wallet address
func (etherMan *Client) GetAddress() common.Address {
	return etherMan.auth.From
}

// GetDefaultChainID function allows to retrieve the default chainID from the smc
func (etherMan *Client) GetDefaultChainID() (*big.Int, error) {
	defaulChainID, err := etherMan.PoE.DEFAULTCHAINID(&bind.CallOpts{Pending: false})
	return new(big.Int).SetUint64(uint64(defaulChainID)), err
}

// GetCustomChainID function allows to retrieve the custom chainID from the latest
// status of the smart contract (not meant to be used by the synchronizer).
func (etherMan *Client) GetCustomChainID() (*big.Int, error) {
	address := etherMan.GetAddress()
	sequencer, err := etherMan.PoE.Sequencers(&bind.CallOpts{Pending: false}, address)
	return new(big.Int).SetUint64(uint64(sequencer.ChainID)), err
}

// EstimateSendBatchCost function estimate gas cost for sending batch to ethereum sc
func (etherMan *Client) EstimateSendBatchCost(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*big.Int, error) {
	noSendOpts := etherMan.auth
	noSendOpts.NoSend = true
	tx, err := etherMan.sendBatch(ctx, noSendOpts, txs, maticAmount)
	if err != nil {
		return nil, err
	}
	return tx.Cost(), nil
}

// GetLatestProposedBatchNumber function allows to retrieve the latest proposed batch in the smc
func (etherMan *Client) GetLatestProposedBatchNumber() (uint64, error) {
	latestBatch, err := etherMan.PoE.LastBatchSent(&bind.CallOpts{Pending: false})
	return uint64(latestBatch), err
}

// GetLatestConsolidatedBatchNumber function allows to retrieve the latest consolidated batch in the smc
func (etherMan *Client) GetLatestConsolidatedBatchNumber() (uint64, error) {
	latestBatch, err := etherMan.PoE.LastVerifiedBatch(&bind.CallOpts{Pending: false})
	return uint64(latestBatch), err
}

// GetSequencerCollateralByBatchNumber function allows to retrieve the sequencer collateral from the smc.
func (etherMan *Client) GetSequencerCollateralByBatchNumber(batchNumber uint64) (*big.Int, error) {
	batchInfo, err := etherMan.PoE.SentBatches(&bind.CallOpts{Pending: false}, uint32(batchNumber))
	return batchInfo.MaticCollateral, err
}

// ApproveMatic function allow to approve tokens in matic smc
func (etherMan *Client) ApproveMatic(maticAmount *big.Int, to common.Address) (*types.Transaction, error) {
	tx, err := etherMan.Matic.Approve(etherMan.auth, etherMan.SCAddresses[0], maticAmount)
	if err != nil {
		return nil, fmt.Errorf("error approving balance to send the batch. Error: %w", err)
	}
	return tx, nil
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (etherMan *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return etherMan.EtherClient.HeaderByNumber(ctx, number)
}

// GetCurrentSequencerCollateral function allows to retrieve the current sequencer collateral from the smc.
func (etherMan *Client) GetCurrentSequencerCollateral() (*big.Int, error) {
	return etherMan.PoE.CalculateSequencerCollateral(&bind.CallOpts{Pending: false})
}
