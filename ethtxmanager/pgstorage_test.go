package ethtxmanager

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddGetAndUpdate(t *testing.T) {
	dbCfg := dbutils.NewStateConfigFromEnv()
	require.NoError(t, dbutils.InitOrResetState(dbCfg))

	storage, err := NewPostgresStorage(dbCfg)
	require.NoError(t, err)

	id := "id"
	from := common.HexToAddress("0x1")
	to := common.HexToAddress("0x2")
	nonce := uint64(1)
	value := big.NewInt(2)
	data := []byte("data")
	gas := uint64(3)
	gasPrice := big.NewInt(4)
	status := MonitoredTxStatusCreated
	history := map[common.Hash]bool{common.HexToHash("0x3"): true, common.HexToHash("0x4"): true}

	mTx := monitoredTx{
		id: id, from: from, to: &to, nonce: nonce, value: value, data: data,
		gas: gas, gasPrice: gasPrice, status: status, history: history,
	}
	err = storage.Add(context.Background(), mTx, nil)
	require.NoError(t, err)

	returnedMtx, err := storage.Get(context.Background(), id, nil)
	require.NoError(t, err)

	assert.Equal(t, id, returnedMtx.id)
	assert.Equal(t, from.String(), returnedMtx.from.String())
	assert.Equal(t, to.String(), returnedMtx.to.String())
	assert.Equal(t, nonce, returnedMtx.nonce)
	assert.Equal(t, value, returnedMtx.value)
	assert.Equal(t, data, returnedMtx.data)
	assert.Equal(t, gas, returnedMtx.gas)
	assert.Equal(t, gasPrice, returnedMtx.gasPrice)
	assert.Equal(t, status, returnedMtx.status)
	assert.Equal(t, history, returnedMtx.history)
	assert.Greater(t, time.Now().UTC().Round(time.Microsecond), returnedMtx.createdAt)
	assert.Less(t, time.Time{}, returnedMtx.createdAt)
	assert.Greater(t, time.Now().UTC().Round(time.Microsecond), returnedMtx.updatedAt)
	assert.Less(t, time.Time{}, returnedMtx.updatedAt)

	from = common.HexToAddress("0x11")
	to = common.HexToAddress("0x22")
	nonce = uint64(11)
	value = big.NewInt(22)
	data = []byte("data data")
	gas = uint64(33)
	gasPrice = big.NewInt(44)
	status = MonitoredTxStatusFailed
	history = map[common.Hash]bool{common.HexToHash("0x33"): true, common.HexToHash("0x44"): true}

	mTx = monitoredTx{
		id: id, from: from, to: &to, nonce: nonce, value: value, data: data,
		gas: gas, gasPrice: gasPrice, status: status, history: history,
	}
	err = storage.Update(context.Background(), mTx, nil)
	require.NoError(t, err)

	returnedMtx, err = storage.Get(context.Background(), id, nil)
	require.NoError(t, err)

	assert.Equal(t, id, returnedMtx.id)
	assert.Equal(t, from.String(), returnedMtx.from.String())
	assert.Equal(t, to.String(), returnedMtx.to.String())
	assert.Equal(t, nonce, returnedMtx.nonce)
	assert.Equal(t, value, returnedMtx.value)
	assert.Equal(t, data, returnedMtx.data)
	assert.Equal(t, gas, returnedMtx.gas)
	assert.Equal(t, gasPrice, returnedMtx.gasPrice)
	assert.Equal(t, status, returnedMtx.status)
	assert.Equal(t, history, returnedMtx.history)
	assert.Greater(t, time.Now().UTC().Round(time.Microsecond), returnedMtx.createdAt)
	assert.Less(t, time.Time{}, returnedMtx.createdAt)
	assert.Greater(t, time.Now().UTC().Round(time.Microsecond), returnedMtx.updatedAt)
	assert.Less(t, time.Time{}, returnedMtx.updatedAt)
}

func TestAddAndGetByStatus(t *testing.T) {
	dbCfg := dbutils.NewStateConfigFromEnv()
	require.NoError(t, dbutils.InitOrResetState(dbCfg))

	storage, err := NewPostgresStorage(dbCfg)
	require.NoError(t, err)

	to := common.HexToAddress("0x2")
	baseMtx := monitoredTx{
		from: common.HexToAddress("0x1"), to: &to, nonce: uint64(1), value: big.NewInt(2), data: []byte("data"),
		gas: uint64(3), gasPrice: big.NewInt(4), history: map[common.Hash]bool{common.HexToHash("0x3"): true, common.HexToHash("0x4"): true},
	}

	type mTxReplaceInfo struct {
		id     string
		status MonitoredTxStatus
	}

	mTxsReplaceInfo := []mTxReplaceInfo{
		{id: "created1", status: MonitoredTxStatusCreated},
		{id: "sent1", status: MonitoredTxStatusSent},
		{id: "failed1", status: MonitoredTxStatusFailed},
		{id: "confirmed1", status: MonitoredTxStatusConfirmed},
		{id: "created2", status: MonitoredTxStatusCreated},
		{id: "sent2", status: MonitoredTxStatusSent},
		{id: "failed2", status: MonitoredTxStatusFailed},
		{id: "confirmed2", status: MonitoredTxStatusConfirmed},
	}

	for _, replaceInfo := range mTxsReplaceInfo {
		baseMtx.id = replaceInfo.id
		baseMtx.status = replaceInfo.status
		baseMtx.createdAt = baseMtx.createdAt.Add(time.Microsecond)
		baseMtx.updatedAt = baseMtx.updatedAt.Add(time.Microsecond)
		err = storage.Add(context.Background(), baseMtx, nil)
		require.NoError(t, err)
	}

	mTxs, err := storage.GetByStatus(context.Background(), []MonitoredTxStatus{MonitoredTxStatusConfirmed}, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, len(mTxs))
	assert.Equal(t, "confirmed1", mTxs[0].id)
	assert.Equal(t, "confirmed2", mTxs[1].id)

	mTxs, err = storage.GetByStatus(context.Background(), []MonitoredTxStatus{MonitoredTxStatusSent, MonitoredTxStatusCreated}, nil)
	require.NoError(t, err)
	assert.Equal(t, 4, len(mTxs))
	assert.Equal(t, "created1", mTxs[0].id)
	assert.Equal(t, "sent1", mTxs[1].id)
	assert.Equal(t, "created2", mTxs[2].id)
	assert.Equal(t, "sent2", mTxs[3].id)

	mTxs, err = storage.GetByStatus(context.Background(), []MonitoredTxStatus{}, nil)
	require.NoError(t, err)
	assert.Equal(t, 8, len(mTxs))
	assert.Equal(t, "created1", mTxs[0].id)
	assert.Equal(t, "sent1", mTxs[1].id)
	assert.Equal(t, "failed1", mTxs[2].id)
	assert.Equal(t, "confirmed1", mTxs[3].id)
	assert.Equal(t, "created2", mTxs[4].id)
	assert.Equal(t, "sent2", mTxs[5].id)
	assert.Equal(t, "failed2", mTxs[6].id)
	assert.Equal(t, "confirmed2", mTxs[7].id)
}
