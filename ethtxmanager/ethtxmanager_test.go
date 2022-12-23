package ethtxmanager

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	cfg := Config{
		FrequencyForResendingFailedTxs: types.NewDuration(time.Second),
		WaitTxToBeMined:                types.NewDuration(1 * time.Minute),
	}
	dbCfg := dbutils.NewStateConfigFromEnv()
	require.NoError(t, dbutils.InitOrResetState(dbCfg))

	etherman := newSimulatedEtherman(t)
	storage, err := NewPostgresStorage(dbCfg)
	require.NoError(t, err)

	ethTxManagerClient := New(cfg, etherman, storage)

	id := "unique_id"
	from := common.HexToAddress("")
	var to *common.Address
	var value *big.Int
	var data []byte = nil

	ctx := context.Background()

	require.NoError(t, ethTxManagerClient.Add(ctx, id, from, to, value, data, nil))

	status, err := ethTxManagerClient.Status(ctx, id, nil)
	require.NoError(t, err)
	require.NotNil(t, status)
}

// This function prepare the blockchain, the wallet with funds and deploy the smc
func newSimulatedEtherman(t *testing.T) *etherman.Client {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)
	ethman, _, _, _, err := etherman.NewSimulatedEtherman(etherman.Config{}, auth)
	require.NoError(t, err)
	return ethman
}
