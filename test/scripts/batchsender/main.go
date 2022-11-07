package main

import (
	"bytes"
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func main() {
	var cfg config.Config
	viper.SetConfigType("toml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte(config.DefaultValues)))
	checkErr(err)
	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	checkErr(err)

	client, err := ethclient.Dial(operations.DefaultL1NetworkURL)
	checkErr(err)

	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL1ChainID)
	checkErr(err)

	ethMan, err := etherman.NewClient(cfg.Etherman, auth)
	checkErr(err)

	ctx := context.Background()

	amount := big.NewInt(10000)
	toAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	// toAddress := cfg.Etherman.PoEAddr

	log.Debug("estimating gas limit")
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &toAddress, Value: amount})
	checkErr(err)

	log.Debug("estimating gas price")
	gasPrice, err := client.SuggestGasPrice(ctx)
	checkErr(err)

	log.Debug("getting nonce")
	nonce, err := client.PendingNonceAt(ctx, auth.From)
	checkErr(err)

	tx := ethtypes.NewTransaction(nonce+uint64(1), toAddress, amount, gasLimit, gasPrice, nil)
	seqs := []types.Sequence{{
		Txs: []ethtypes.Transaction{*tx},
	}}
	_, err = ethMan.SequenceBatches(ctx, seqs, gasLimit*2, gasPrice, big.NewInt(int64(nonce)))
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
