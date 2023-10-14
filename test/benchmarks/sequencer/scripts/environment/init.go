package environment

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	maxConnections = 10
	bitSize        = 64
)

// Init sets up the environment for the benchmark
func Init() (*pgpoolstorage.PostgresPoolStorage, *ethclient.Client, *bind.TransactOpts) {
	ctx := context.Background()
	pl, err := pgpoolstorage.NewPostgresPoolStorage(db.Config{
		Name:      poolDbName,
		User:      poolDbUser,
		Password:  poolDbPass,
		Host:      poolDbHost,
		Port:      poolDbPort,
		EnableLog: false,
		MaxConns:  maxConnections,
	})
	if err != nil {
		panic(err)
	}

	l2Client, err := ethclient.Dial(L2NetworkRPCURL)
	if err != nil {
		panic(err)
	}
	// PrivateKey is the private key of the sender
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(PrivateKey, "0x"))
	if err != nil {
		panic(err)
	}
	chainId, err := strconv.ParseUint(L2ChainId, IntBase, bitSize)
	if err != nil {
		panic(err)
	}
	fmt.Printf("L2ChainId: %d\n", chainId)
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
	fmt.Printf("Receiver Addr: %v\n", params.To.String())
	fmt.Printf("Sender Addr: %v\n", auth.From.String())
	fmt.Printf("Sender Balance: %v\n", senderBalance.String())
	fmt.Printf("Sender Nonce: %v\n", senderNonce)

	gasPrice, err := l2Client.SuggestGasPrice(ctx)
	if err != nil {
		panic(err)
	}
	auth.GasPrice = gasPrice

	return pl, l2Client, auth
}
