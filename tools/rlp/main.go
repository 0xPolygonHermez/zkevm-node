package main

import (
	"math/big"
	"strings"
	"strconv"

	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/encoding"
)

const (
	ether155V = 27
)

func main() {
	// Short txs
	callData := "0x06d6490f00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000000000000000000000000e1ee80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a00008082019180806e209c61ca92c2b980d6197e7ac9ccc3f547bf13be6455dfe682aa5dda9655ef16819a7edcc3fefec81ca97c7a6f3d10ec774440e409adbba693ce8b698d41f11cef80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff89056bc75e2d63100000808203e98080fe1e96b35c836fbebac887681150c5fc9fdae862d747aaaf8c30373c0becf7691ff0c900aaaac6d1565a603f69b5a45f222ed205f0a36fdc6e4e4c5a7b88d45b1b00000000000000000000000000000000000000000000000000000000000000"
	bytesCallData, err := hex.DecodeHex(callData)
	if err != nil {
		log.Fatal("error decoding callData: ", err)
	}
	txs, rawTxs, err := decodeFullCallDataToTxs(bytesCallData)
	log.Debug("RawTxs: ", hex.EncodeToHex(rawTxs))
	for _, tx := range txs {
		log.Debug("#######################################################################")
		log.Debug("#######################################################################")
		log.Debugf("Decoded tx: %+v", tx)
		log.Debug("ChainID: ", tx.ChainId())
		log.Debug("Cost: ", tx.Cost())
		log.Debug("Data: ", hex.EncodeToString(tx.Data()))
		log.Debug("Gas: ", tx.Gas())
		log.Debug("GasPrice: ", tx.GasPrice())
		log.Debug("Hash: ", tx.Hash())
		log.Debug("Nonce: ", tx.Nonce())
		v, r, s := tx.RawSignatureValues()
		log.Debug("V: ", v)
		log.Debug("R: ", r)
		log.Debug("S: ", s)
		log.Debug("To: ", tx.To())
		log.Debug("Type: ", tx.Type())
		log.Debug("Value: ", tx.Value())
	}

	//Long Tx. SmartContract deployment
	callData = "0x06d6490f00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000000000000000000000001c9f90185808502540be400832dc6c08080b90170608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea2646970667358221220a92848e725904714428852952fc1c595af6ecdf7f8cfcc7c961cba7c7208044d64736f6c634300080700338203e980803f34a8855378618502c371823a1a6a5d244d9681d6aaf9b35338925ac875c64a4a613f018d4d5842364afdb4a9c1e123983448d82c6d857aba49dd0495bcb2bd1b0000000000000000000000000000000000000000000000"
	bytesCallData, err = hex.DecodeHex(callData)
	if err != nil {
		log.Fatal("error decoding callData: ", err)
	}
	txs, rawTxs, err = decodeFullCallDataToTxs(bytesCallData)
	log.Debug("RawTxs: ", hex.EncodeToHex(rawTxs))
	for _, tx := range txs {
		log.Debug("#######################################################################")
		log.Debug("#######################################################################")
		log.Debugf("Decoded tx: %+v", tx)
		log.Debug("ChainID: ", tx.ChainId())
		log.Debug("Cost: ", tx.Cost())
		log.Debug("Data: ", hex.EncodeToString(tx.Data()))
		log.Debug("Gas: ", tx.Gas())
		log.Debug("GasPrice: ", tx.GasPrice())
		log.Debug("Hash: ", tx.Hash())
		log.Debug("Nonce: ", tx.Nonce())
		v, r, s := tx.RawSignatureValues()
		log.Debug("V: ", v)
		log.Debug("R: ", r)
		log.Debug("S: ", s)
		log.Debug("To: ", tx.To())
		log.Debug("Type: ", tx.Type())
		log.Debug("Value: ", tx.Value())
	}
}

func decodeFullCallDataToTxs(txsData []byte) ([]*types.Transaction, []byte, error) {
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

	txs, err := decodeRawTxs(txsData)
	return txs, txsData, err
}


func decodeRawTxs(txsData []byte) ([]*types.Transaction, error) {
	// Process coded txs
	var pos int64
	var txs []*types.Transaction
	const (
		headerByteLength = 2
		sLength          = 32
		rLength          = 32
		vLength          = 1
		c0               = 192 // 192 is c0. This value is defined by the rlp protocol
		ff               = 255 // max value of rlp header
		shortRlp         = 55 // length of the short rlp codification
		f7               = 247 // 192 + 55 = c0 + shortRlp
		etherNewV        = 35
		mul2             = 2
	)
	txDataLength := len(txsData)
	for pos < int64(txDataLength) {
		num, err := strconv.ParseInt(hex.EncodeToString(txsData[pos : pos+1]), hex.Base, encoding.BitSize64)
		if err != nil {
			log.Debug("error parsing header length: ", err)
			return []*types.Transaction{}, err
		}
		// First byte is the length and must be ignored
		len := num - c0 - 1

		if len > shortRlp { // If rlp is bigger than lenght 55
			// numH is the length of the bytes that give the length of the rlp
			numH, err := strconv.ParseInt(hex.EncodeToString(txsData[pos : pos+1]), hex.Base, encoding.BitSize64)
			if err != nil {
				log.Debug("error parsing length of the bytes: ", err)
				return []*types.Transaction{}, err
			}
			// n is the length of the rlp data without the header (1 byte) for example "0xf7"
			n, err := strconv.ParseInt(hex.EncodeToString(txsData[pos+1 : pos+1+numH-f7]), hex.Base, encoding.BitSize64) // +1 is the header. For example 0xf7
			if err != nil {
				log.Debug("error parsing length: ", err)
				return []*types.Transaction{}, err
			}
			len = n+1 // +1 is the header. For example 0xf7
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
			log.Debug("error decoding tx bytes: ", err, ". fullDataTx: ", hex.EncodeToString(fullDataTx), "\n tx: ", hex.EncodeToString(txInfo))
			return []*types.Transaction{}, err
		}

		//tx.V = v-27+chainId*2+35
		tx.V = new(big.Int).Add(new(big.Int).Sub(new(big.Int).SetBytes(v), big.NewInt(ether155V)), new(big.Int).Add(new(big.Int).Mul(tx.V, big.NewInt(mul2)), big.NewInt(etherNewV)))
		tx.R = new(big.Int).SetBytes(r)
		tx.S = new(big.Int).SetBytes(s)

		txs = append(txs, types.NewTx(&tx))
	}
	return txs, nil
}