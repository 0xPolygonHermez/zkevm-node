package tree

import (
	"github.com/hermeznetwork/hermez-core/db"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicTree(t *testing.T) {

	dbCfg := db.NewConfigFromEnv()

	err := db.RunMigrations(dbCfg)
	require.NoError(t, err)

	mtDb, err := db.NewSQLDB(dbCfg)
	require.NoError(t, err)

	defer mtDb.Close()

	tree := NewBasicTree(mtDb)

	address := common.Address{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
	}

	// Balance
	bal, err := tree.GetBalance(address, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), bal)

	_, _, err = tree.SetBalance(address, big.NewInt(1))
	require.NoError(t, err)

	bal, err = tree.GetBalance(address, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1), bal)

	// Nonce
	nonce, err := tree.GetNonce(address, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), nonce)

	_, _, err = tree.SetNonce(address, big.NewInt(2))
	require.NoError(t, err)

	nonce, err = tree.GetNonce(address, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(2), nonce)

	// Code
	//code, err := tree.GetCode(address, nil)
	//require.NoError(t, err)
	//assert.Equal(t, nil, code)

	// Storage
	position := common.Hash{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
		0x30, 0x31,
	}
	storage, err := tree.GetStorageAt(address, position, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), storage)

	_, _, err = tree.SetStorageAt(address, position, big.NewInt(4))
	require.NoError(t, err)

	storage, err = tree.GetStorageAt(address, position, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(4), storage)

	position2 := common.Hash{
		0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x11, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x21, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
		0x31, 0x31,
	}

	storage2, err := tree.GetStorageAt(address, position2, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), storage2)

	_, _, err = tree.SetStorageAt(address, position2, big.NewInt(5))
	require.NoError(t, err)

	storage2, err = tree.GetStorageAt(address, position2, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(5), storage2)

	storage, err = tree.GetStorageAt(address, position, nil)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(4), storage)
}
