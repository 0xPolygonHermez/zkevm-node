package jsonrpc

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/stretchr/testify/require"
)

type mockedServer struct {
	DefaultChainID   uint64
	ChainID          uint64
	SequencerAddress common.Address

	Server *Server
}

type mocks struct {
	Pool              *poolMock
	State             *stateMock
	BatchProcessor    *batchProcessorMock
	GasPriceEstimator *gasPriceEstimatorMock
	Storage           *storageMock
}

func newMockedServer(t *testing.T) (*mockedServer, *mocks, *ethclient.Client) {
	const (
		defaultChainID      = 1000
		chainID             = 1001
		sequencerAddressHex = "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"

		host                      = "localhost"
		port                      = 8123
		maxRequestsPerIPAndSecond = 1000
	)

	cfg := Config{
		Host:                      host,
		Port:                      port,
		MaxRequestsPerIPAndSecond: maxRequestsPerIPAndSecond,
	}

	sequencerAddress := common.HexToAddress(sequencerAddressHex)
	pool := newPoolMock(t)
	state := newStateMock(t)
	batchProcessor := newBatchProcessorMock(t)
	gasPriceEstimator := newGasPriceEstimatorMock(t)
	storage := newStorageMock(t)

	server := NewServer(cfg, defaultChainID, chainID, sequencerAddress,
		pool, state, gasPriceEstimator, storage)

	go func() {
		err := server.Start()
		if err != nil {
			panic(err)
		}
	}()

	serverURL := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	for {
		log.Debugf("waiting server to get ready...")
		res, err := http.Get(serverURL) //nolint:gosec
		if err == nil && res.StatusCode == http.StatusOK {
			log.Debugf("server ready!")
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	ethClient, err := ethclient.Dial(serverURL)
	require.NoError(t, err)

	msv := &mockedServer{
		DefaultChainID:   defaultChainID,
		ChainID:          chainID,
		SequencerAddress: sequencerAddress,

		Server: server,
	}

	mks := &mocks{
		Pool:              pool,
		State:             state,
		BatchProcessor:    batchProcessor,
		GasPriceEstimator: gasPriceEstimator,
		Storage:           storage,
	}

	return msv, mks, ethClient
}

func (s *mockedServer) Stop() {
	err := s.Server.Stop()
	if err != nil {
		panic(err)
	}
}
