package etherman

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
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
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var newBatcheventSignatureHash = crypto.Keccak256Hash([]byte("SendBatch(uint256,address)"))

type EtherMan struct {
	EtherClient *ethclient.Client
	PoE         *proofofefficiency.Proofofefficiency
	SmcAddreses []common.Address
}

func NewEtherman(url string, poeAddr common.Address) (*EtherMan, error) {
	//TODO
	//Connect to ethereum node
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		log.Printf("error connecting to %s: %+v", url, err)
		return nil, err
	}
	poe, err := proofofefficiency.NewProofofefficiency(poeAddr, ethClient)
	if err != nil {
		return nil, err
	}
	var smcAddreses []common.Address
	smcAddreses = append(smcAddreses, poeAddr)
	return &EtherMan{EtherClient: ethClient, PoE: poe}, nil
}

// EthBlockByNumber function retrieves the ethereum block information by ethereum block number
func (etherMan *EtherMan) EthBlockByNumber(blockNum int64) (types.Block, error) {
	//TODO
	return types.Block{}, nil
}

// GetBatchesByBlock function retrieves the batches information that are included in a specific ethereum block
func (etherMan *EtherMan) GetBatchesByBlock(blockNum int64) ([]state.Batch, error) {
	//TODO
	return []state.Batch{}, nil
}

// GetBatchesFromBlockTo function retrieves the batches information that are included in all this ethereum blocks
//from block x to block y
func (etherMan *EtherMan) GetBatchesFromBlockTo(fromBlock uint64, toBlock uint64) ([]state.Batch, error) {
	//First filter query
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: etherMan.SmcAddreses,
	}
	batchEvent, err := etherMan.readEvents(fromBlock, toBlock, query)
	if err != nil {
		return nil, err
	}
	fmt.Println(batchEvent)
	var batch []state.Batch 	//TODO Fix the type problem
	return batch, nil
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

type batchEvent struct {
	TxHash    common.Hash
	BlockNum  uint64
	BlockHash common.Hash
	Batch struct {
		BatchNum  uint64
		Sequencer common.Address
		Txs []types.LegacyTx
	}
}

func (etherMan *EtherMan) readEvents(fromBlock uint64, toBlock uint64, query ethereum.FilterQuery) (batchEvent, error) {
	logs, err := etherMan.EtherClient.FilterLogs(context.Background(), query)
	if err != nil {
		return batchEvent{}, err
	}
	var newBatch batchEvent
	for _, vLog := range logs {
		switch vLog.Topics[0] {
		case newBatcheventSignatureHash:
			//Indexed parameters using topics
			newBatch.Batch.BatchNum = new(big.Int).SetBytes(vLog.Topics[1][:]).Uint64()
			newBatch.Batch.Sequencer = common.BytesToAddress(vLog.Topics[2].Bytes())
			newBatch.TxHash = vLog.TxHash
			newBatch.BlockNum = vLog.BlockNumber
			newBatch.BlockHash = vLog.BlockHash
			//Now, We have to read the tx for this batch.
			tx, isPending, err := etherMan.EtherClient.TransactionByHash(context.Background(), newBatch.TxHash)
			if err != nil || isPending {
				return batchEvent{}, err
			}
			//Get sequencer chainId
			seq, err := etherMan.PoE.Sequencers(&bind.CallOpts{Pending: false}, newBatch.Batch.Sequencer)
			if err != nil {
				return batchEvent{}, err
			}
			chainId := seq.ChainID
			txData := tx.Data()
			txs, err := decodeTxs(txData, chainId)
			if err != nil {
				return batchEvent{}, err
			}
			newBatch.Batch.Txs = txs
		}
	}
	return batchEvent{}, nil
}

func decodeTxs(txsData []byte, chainId *big.Int) ([]types.LegacyTx, error) {
	// First split txs
	var pos int64
	var txs []types.LegacyTx
	for pos<int64(len(txsData)) {
		lenght := txsData[pos+1:pos+2]
		str := hex.EncodeToString(lenght)
		num, err := strconv.ParseInt(str, 16, 64)
		if err != nil {
			fmt.Println("error: skipping tx. Err: ", err)
			continue
		}
		data := txsData[pos : pos+num+2]
		pos = pos + num + 2

		//Decode tx
		var tx types.LegacyTx
		rlp.DecodeBytes(data, &tx)
		isValid, fromAddr, err := checkSignature(tx, chainId)
		if err != nil {
			fmt.Println("error: skipping tx. ", err)
			continue
		} else if !isValid {
			fmt.Println("Signature invalid: ",isValid, "FromAddr: ", fromAddr)
			continue
		}
		
		fmt.Println("Signature verified: ",isValid, "FromAddr: ", fromAddr)

		txs = append(txs, tx)
	}
    return txs, nil
}

func checkSignature(tx types.LegacyTx, chainId *big.Int) (bool, common.Address, error) {
	decodedChainId := deriveChainId(tx.V)
	if decodedChainId.Cmp(chainId) != 0 {
		return false, common.Address{}, fmt.Errorf("error: incorrect chainId. Decoded chainId: %d and chainId retrieved from smc: %d",
		decodedChainId, chainId)
	}
	// Formula: v = chainId * 2 + 36 or 35
	x := tx.V.Int64()-(chainId.Int64()*2)
	var vField byte
	if x == 35 {
		vField = byte(0)
	} else if x == 36 {
		vField = byte(1)
	} else {
		return false, common.Address{}, fmt.Errorf("Error invalid signature v value: %d", tx.V)
	}
	if !crypto.ValidateSignatureValues(vField, tx.R, tx.S, false) {
		fmt.Println("Invalid Signature values")
		return false, common.Address{}, fmt.Errorf("Error invalid Signature values")
	}

	//Get fromSender address
	var signature []byte
	signature = append(signature, tx.R.Bytes()[:]...)
	signature = append(signature, tx.S.Bytes()[:]...)
	signature = append(signature, vField)
	fmt.Println("signature: ",hexutil.Encode(signature))

	auxTx := struct {
		Nonce    uint64
		Gasprice *big.Int
		Gas      uint64
		To       *common.Address
		Value    *big.Int
		Data     []byte
		ChainID  *big.Int
		// aux      *big.Int
		// aux1     *big.Int

	}{
		Nonce: tx.Nonce,
		Gasprice: tx.GasPrice,
		Gas: tx.Gas,
		To: tx.To,
		Value: tx.Value,
		Data: tx.Data,
		ChainID: chainId,
	}
	fmt.Println(auxTx)
	auxTxB, err := json.Marshal(auxTx)
	if err != nil {
		fmt.Printf("error marshaling: %v\n", err)
		return false, common.Address{}, fmt.Errorf("error marshaling: %v\n", err)
	}

	hash := crypto.Keccak256Hash(auxTxB)
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
    if err != nil {
        log.Fatal(err)
		return false, common.Address{}, fmt.Errorf("error: %v\n", err)
    }
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
    if err != nil {
        log.Fatal(err)
		return false, common.Address{}, fmt.Errorf("error: %v\n", err)
    }

	fromAddr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	// sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)

	// fmt.Println("PubK1: ", sigPublicKey)
	// fmt.Println("PubK2: ", sigPublicKeyBytes)
	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	isValid := crypto.VerifySignature(sigPublicKey, hash.Bytes(), signatureNoRecoverID)

	return isValid, fromAddr, nil
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