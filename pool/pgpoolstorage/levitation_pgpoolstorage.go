package pgpoolstorage

import (
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"time"
)

// LevitationPoolStorage uses levitation chain to store data
type LevitationPoolStorage struct {
}

// NewLevitationPoolStorage creates and initializes an instance of LevitationProolStorage
func NewLevitationPoolStorage() (*LevitationPoolStorage, error) {
	return &LevitationPoolStorage{}, nil
}

// AddTx adds a transaction to the pool table with the provided status
func (p *LevitationPoolStorage) AddTx(hash string, encoded string, decoded string, txStatus pool.TxStatus, gasPrice uint64,
	nonce uint64, cumulativeGasUsed uint64, usedKeccakHashes uint32, usedPoseidonHashes uint32,
	usedPoseidonPaddings uint32, usedMemAligns uint32, usedArithmetics uint32, usedBinaries uint32, usedSteps uint32,
	receivedAt time.Time, fromAddress string, isWIP bool, iP string) error {
	return nil
}

// LevitationSetGasPrices sets the latest l2 and l1 gas prices
func (p *LevitationPoolStorage) LevitationSetGasPrices(l2GasPrice, l1GasPrice uint64, timeStamp time.Time) error {
	return nil
}
