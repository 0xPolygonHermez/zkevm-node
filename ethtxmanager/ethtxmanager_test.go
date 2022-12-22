package ethtxmanager

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/db"
	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	etherman, _, _, _, _ := ethman.NewSimulatedEtherman(ethman.Config{}, nil)

	cfg := Config{}
	dbCfg := db.Config{}
	storage, err := NewPostgresStorage(cfg, dbCfg)
	require.NoError(t, err)

	ethTxManagerClient := New(cfg, etherman, storage)

	id := "unique_id"
	from := common.HexToAddress("")
	to := common.HexToAddress("")
	value := big.NewInt(0)
	data := []byte{}

	ctx := context.Background()

	require.NoError(t, ethTxManagerClient.Add(ctx, id, from, &to, value, data))

	status, err := ethTxManagerClient.Status(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, status)
}
