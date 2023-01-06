package ethtxmanager

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestTx(t *testing.T) {
	to := common.HexToAddress("0x2")
	nonce := uint64(1)
	value := big.NewInt(2)
	data := []byte("data")
	gas := uint64(3)
	gasPrice := big.NewInt(4)

	mTx := monitoredTx{
		to:       &to,
		nonce:    nonce,
		value:    value,
		data:     data,
		gas:      gas,
		gasPrice: gasPrice,
	}

	tx := mTx.Tx()

	assert.Equal(t, &to, tx.To())
	assert.Equal(t, nonce, tx.Nonce())
	assert.Equal(t, value, tx.Value())
	assert.Equal(t, data, tx.Data())
	assert.Equal(t, gas, tx.Gas())
	assert.Equal(t, gasPrice, tx.GasPrice())
}
