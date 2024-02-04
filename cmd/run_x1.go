package main

import (
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/config/apollo"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencesender"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func createSequenceSenderX1(cfg config.Config, pool *pool.Pool, etmStorage *ethtxmanager.PostgresStorage, st *state.State, eventLog *event.EventLog) *sequencesender.SequenceSender {
	etherman, err := newEtherman(cfg)
	if err != nil {
		log.Fatal(err)
	}

	_, privKey, err := etherman.LoadAuthFromKeyStoreX1(cfg.SequenceSender.DAPermitApiPrivateKey.Path, cfg.SequenceSender.DAPermitApiPrivateKey.Password)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("from pk %s, from sender %s", crypto.PubkeyToAddress(privKey.PublicKey), cfg.SequenceSender.SenderAddress.String())
	if cfg.SequenceSender.SenderAddress.Cmp(common.Address{}) == 0 {
		log.Fatal("Sequence sender address not found")
	}
	if privKey == nil {
		log.Fatal("DA permit api private key not found")
	}

	cfg.SequenceSender.ForkUpgradeBatchNumber = cfg.ForkUpgradeBatchNumber

	ethTxManager := ethtxmanager.New(cfg.EthTxManager, etherman, etmStorage, st)

	seqSender, err := sequencesender.New(cfg.SequenceSender, st, etherman, ethTxManager, eventLog, privKey)
	if err != nil {
		log.Fatal(err)
	}

	return seqSender
}

func initRunForX1(c *config.Config, components []string) {
	// Read configure from apollo
	apolloClient := apollo.NewClient(c)
	if apolloClient.LoadConfig() {
		log.Info("apollo config loaded")
	}

	pool.SetL2BridgeAddr(c.NetworkConfig.L2BridgeAddr)

	for _, component := range components {
		switch component {
		case RPC:
			jsonrpc.InitRateLimit(c.RPC.RateLimit)
		}
	}
}
