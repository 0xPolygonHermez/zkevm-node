package clients

import (
	"context"
	"math/big"
	"strconv"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/shared"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/scripts/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Init() (context.Context, *pgpoolstorage.PostgresPoolStorage, *ethclient.Client, *bind.TransactOpts, uint64) {
	ctx := common.Ctx
	pl, err := pgpoolstorage.NewPostgresPoolStorage(db.Config{
		Name:      common.PoolDbName,
		User:      common.PoolDbUser,
		Password:  common.PoolDbPass,
		Host:      common.PoolDbHost,
		Port:      common.PoolDbPort,
		EnableLog: false,
		MaxConns:  10,
	})
	if err != nil {
		panic(err)
	}

	l2Client, err := ethclient.Dial(common.L2NetworkRPCURL)
	if err != nil {
		panic(err)
	}
	// PrivateKey is the private key of the sender
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(common.PrivateKey, "0x"))
	if err != nil {
		panic(err)
	}
	chainId, err := strconv.ParseUint(common.ChainId, 10, 64)
	if err != nil {
		panic(err)
	}
	// Auth is the auth of the sender
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(chainId))
	if err != nil {
		panic(err)
	}

	// Print Info before send
	senderBalance, err := l2Client.BalanceAt(ctx, auth.From, nil)
	if err != nil {
		panic(err)
	}
	senderNonce, err := l2Client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		panic(err)
	}

	// Print Initial Stats
	log.Infof("Receiver Addr: %v", shared.To.String())
	log.Infof("Sender Addr: %v", auth.From.String())
	log.Infof("Sender Balance: %v", senderBalance.String())
	log.Infof("Sender Nonce: %v", senderNonce)

	gasPrice, err := l2Client.SuggestGasPrice(ctx)
	if err != nil {
		panic(err)
	}
	auth.GasPrice = gasPrice

	return ctx, pl, l2Client, auth, senderNonce
}
