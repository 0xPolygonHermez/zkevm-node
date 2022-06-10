package sequencer_test

import (
	"context"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	cfgTypes "github.com/hermeznetwork/hermez-core/config/types"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/pool/pgpoolstorage"
	"github.com/hermeznetwork/hermez-core/sequencer"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/tree"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

var senderPrivateKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"

var dbCfg = dbutils.NewConfigFromEnv()

var queueCfg = sequencer.PendingTxsQueueConfig{
	TxPendingInQueueCheckingFrequency: cfgTypes.NewDuration(1 * time.Second),
	GetPendingTxsFrequency:            cfgTypes.NewDuration(1 * time.Second),
}

func TestQueue_AddAndPopTx(t *testing.T) {
	if err := dbutils.InitOrReset(dbCfg); err != nil {
		panic(err)
	}

	sqlDB, err := db.NewSQLDB(dbCfg)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	st := newState(sqlDB)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	balance, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	genesis := state.Genesis{
		Block: genesisBlock,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): balance,
		},
	}
	err = st.SetGenesis(context.Background(), genesis, "")
	if err != nil {
		panic(err)
	}

	s, err := pgpoolstorage.NewPostgresPoolStorage(dbCfg)
	if err != nil {
		panic(err)
	}

	p := pool.NewPool(s, st, common.Address{})

	const txsCount = 10

	ctx := context.Background()

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		panic(err)
	}

	pendQueue := sequencer.NewPendingTxsQueue(queueCfg, p)
	go pendQueue.KeepPendingTxsQueue(ctx)
	go pendQueue.CleanPendTxsChan(ctx)
	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(i)), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		if err != nil {
			panic(err)
		}
		if err := p.AddTx(ctx, *signedTx); err != nil {
			panic(err)
		}
	}
	tx := pendQueue.PopPendingTx()
	assert.Equal(t, uint64(19), tx.GasPrice().Uint64())
	assert.Equal(t, 9, pendQueue.GetPendingTxsQueueLength())
	tx = pendQueue.PopPendingTx()
	assert.Equal(t, uint64(18), tx.GasPrice().Uint64())
	assert.Equal(t, 8, pendQueue.GetPendingTxsQueueLength())

	newTx := types.NewTransaction(uint64(txsCount), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(txsCount)), []byte{})
	signedTx, err := auth.Signer(auth.From, newTx)
	if err != nil {
		panic(err)
	}
	if err := p.AddTx(ctx, *signedTx); err != nil {
		panic(err)
	}

	time.Sleep(queueCfg.TxPendingInQueueCheckingFrequency.Duration * 2)
	assert.Equal(t, 9, pendQueue.GetPendingTxsQueueLength())
}

func TestQueue_AddOneTx(t *testing.T) {
	if err := dbutils.InitOrReset(dbCfg); err != nil {
		panic(err)
	}

	sqlDB, err := db.NewSQLDB(dbCfg)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	st := newState(sqlDB)

	genesisBlock := types.NewBlock(&types.Header{Number: big.NewInt(0)}, []*types.Transaction{}, []*types.Header{}, []*types.Receipt{}, &trie.StackTrie{})
	genesisBlock.ReceivedAt = time.Now()
	balance, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	genesis := state.Genesis{
		Block: genesisBlock,
		Balances: map[common.Address]*big.Int{
			common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"): balance,
		},
	}
	err = st.SetGenesis(context.Background(), genesis, "")
	if err != nil {
		panic(err)
	}

	s, err := pgpoolstorage.NewPostgresPoolStorage(dbCfg)
	if err != nil {
		panic(err)
	}

	p := pool.NewPool(s, st, common.Address{})

	const txsCount = 1

	ctx := context.Background()

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(senderPrivateKey, "0x"))
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		panic(err)
	}

	pendQueue := sequencer.NewPendingTxsQueue(queueCfg, p)
	go pendQueue.KeepPendingTxsQueue(ctx)
	go pendQueue.CleanPendTxsChan(ctx)
	// insert pending transactions
	for i := 0; i < txsCount; i++ {
		tx := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(i)), []byte{})
		signedTx, err := auth.Signer(auth.From, tx)
		if err != nil {
			panic(err)
		}
		if err := p.AddTx(ctx, *signedTx); err != nil {
			panic(err)
		}
	}
	tx := pendQueue.PopPendingTx()
	assert.Equal(t, uint64(10), tx.GasPrice().Uint64())
	assert.Equal(t, 0, pendQueue.GetPendingTxsQueueLength())
	tx = pendQueue.PopPendingTx()
	assert.Nil(t, tx)

	newTx := types.NewTransaction(uint64(txsCount), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(txsCount)), []byte{})
	signedTx, err := auth.Signer(auth.From, newTx)
	if err != nil {
		panic(err)
	}
	if err := p.AddTx(ctx, *signedTx); err != nil {
		panic(err)
	}
	time.Sleep(queueCfg.TxPendingInQueueCheckingFrequency.Duration * 2)
	assert.Equal(t, 1, pendQueue.GetPendingTxsQueueLength())
}

func newState(sqlDB *pgxpool.Pool) *state.State {
	store := tree.NewPostgresStore(sqlDB)
	mt := tree.NewMerkleTree(store, tree.DefaultMerkleTreeArity)
	scCodeStore := tree.NewPostgresSCCodeStore(sqlDB)
	tr := tree.NewStateTree(mt, scCodeStore)

	stateCfg := state.Config{
		DefaultChainID:       1000,
		MaxCumulativeGasUsed: 800000,
	}

	stateDB := state.NewPostgresStorage(sqlDB)
	st := state.NewState(stateCfg, stateDB, tr)

	return st
}
