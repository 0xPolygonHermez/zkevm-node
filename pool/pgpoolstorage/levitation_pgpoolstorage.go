package pgpoolstorage

import (
	"encoding/json"
	"fmt"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// LevitationPoolStorage uses levitation chain to store data
type LevitationPoolStorage struct {
}

// NewLevitationPoolStorage creates and initializes an instance of LevitationProolStorage
func NewLevitationPoolStorage() (*LevitationPoolStorage, error) {
	return &LevitationPoolStorage{}, nil
}

type TxData struct {
	Hash                 string        `json:"hash"`
	Encoded              string        `json:"encoded"`
	Decoded              string        `json:"decoded"`
	TxStatus             pool.TxStatus `json:"txStatus"`
	GasPrice             uint64        `json:"gasPrice"`
	Nonce                uint64        `json:"nonce"`
	CumulativeGasUsed    uint64        `json:"cumulativeGasUsed"`
	UsedKeccakHashes     uint32        `json:"usedKeccakHashes"`
	UsedPoseidonHashes   uint32        `json:"usedPoseidonHashes"`
	UsedPoseidonPaddings uint32        `json:"usedPoseidonPaddings"`
	UsedMemAligns        uint32        `json:"usedMemAligns"`
	UsedArithmetics      uint32        `json:"usedArithmetics"`
	UsedBinaries         uint32        `json:"usedBinaries"`
	UsedSteps            uint32        `json:"usedSteps"`
	ReceivedAt           time.Time     `json:"receivedAt"`
	FromAddress          string        `json:"fromAddress"`
	IsWIP                bool          `json:"isWIP"`
	IP                   string        `json:"iP"`
}

// AddTx adds a transaction to the pool table with the provided status
func (p *LevitationPoolStorage) AddTx(hash string, encoded string, decoded string, txStatus pool.TxStatus, gasPrice uint64,
	nonce uint64, cumulativeGasUsed uint64, usedKeccakHashes uint32, usedPoseidonHashes uint32,
	usedPoseidonPaddings uint32, usedMemAligns uint32, usedArithmetics uint32, usedBinaries uint32, usedSteps uint32,
	receivedAt time.Time, fromAddress string, isWIP bool, iP string) error {

	// Initialize our Data struct
	data := TxData{
		Hash:                 hash,
		Encoded:              encoded,
		Decoded:              decoded,
		TxStatus:             txStatus,
		GasPrice:             gasPrice,
		Nonce:                nonce,
		CumulativeGasUsed:    cumulativeGasUsed,
		UsedKeccakHashes:     usedKeccakHashes,
		UsedPoseidonHashes:   usedPoseidonHashes,
		UsedPoseidonPaddings: usedPoseidonPaddings,
		UsedMemAligns:        usedMemAligns,
		UsedArithmetics:      usedArithmetics,
		UsedBinaries:         usedBinaries,
		UsedSteps:            usedSteps,
		ReceivedAt:           receivedAt,
		FromAddress:          fromAddress,
		IsWIP:                isWIP,
		IP:                   iP,
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error occurred during marshalling: %s", err)
	}

	err2, done := p.writeJsonToFile(err, receivedAt, jsonData, "/tmp/to_pending_queue")
	if done {
		return err2
	}

	return nil
}

func (p *LevitationPoolStorage) writeJsonToFile(err error, receivedAt time.Time, jsonData []byte,
	dirName string) (error, bool) {

	err = os.Mkdir(dirName, 0666)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("Error creating directory: #{err}")
			return err, true
		}
	}
	// Generate a random number between 100000 and 999999
	strRand := strconv.Itoa(rand.Intn(900000) + 100000)
	filename := dirName + "/" + receivedAt.Format(time.RFC3339) + "-" + strRand

	// Write the JSON data to the file
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error occurred during writing to file: %s", err)
		return err, true
	}

	fmt.Println("JSON data written to", filename)
	return nil, false
}

type GasPriceData struct {
	L2GasPrice uint64    `json:"l2GasPrice"`
	L1GasPrice uint64    `json:"l1GasPrice"`
	TimeStamp  time.Time `json:"timeStamp"`
}

// LevitationSetGasPrices sets the latest l2 and l1 gas prices
func (p *LevitationPoolStorage) LevitationSetGasPrices(l2GasPrice, l1GasPrice uint64, timeStamp time.Time) error {

	// Initialize our Data struct
	data := GasPriceData{
		L2GasPrice: l2GasPrice,
		L1GasPrice: l1GasPrice,
		TimeStamp:  timeStamp,
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error occurred during marshalling: %s", err)
	}

	err2, done := p.writeJsonToFile(err, timeStamp, jsonData, "/tmp/gas_price_updates")
	if done {
		return err2
	}

	return nil

}
