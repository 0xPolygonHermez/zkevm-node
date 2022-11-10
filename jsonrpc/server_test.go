package jsonrpc

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	host                      = "localhost"
	maxRequestsPerIPAndSecond = 1000
)

type mockedServer struct {
	Config    Config
	Server    *Server
	ServerURL string
}

type mocks struct {
	Pool              *poolMock
	State             *stateMock
	GasPriceEstimator *gasPriceEstimatorMock
	Storage           *storageMock
	DbTx              *dbTxMock
}

func newMockedServer(t *testing.T, cfg Config) (*mockedServer, *mocks, *ethclient.Client) {
	pool := newPoolMock(t)
	state := newStateMock(t)
	gasPriceEstimator := newGasPriceEstimatorMock(t)
	storage := newStorageMock(t)
	dbTx := newDbTxMock(t)
	apis := map[string]bool{
		APIEth:    true,
		APINet:    true,
		APIDebug:  true,
		APIZKEVM:  true,
		APITxPool: true,
		APIWeb3:   true,
	}

	server := NewServer(cfg, pool, state, gasPriceEstimator, storage, apis)

	go func() {
		err := server.Start()
		if err != nil {
			panic(err)
		}
	}()

	serverURL := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	for {
		fmt.Println("waiting server to get ready...") // fmt is used here to avoid race condition with logs
		res, err := http.Get(serverURL)               //nolint:gosec
		if err == nil && res.StatusCode == http.StatusOK {
			fmt.Println("server ready!") // fmt is used here to avoid race condition with logs
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	ethClient, err := ethclient.Dial(serverURL)
	require.NoError(t, err)

	msv := &mockedServer{
		Config:    cfg,
		Server:    server,
		ServerURL: serverURL,
	}

	mks := &mocks{
		Pool:              pool,
		State:             state,
		GasPriceEstimator: gasPriceEstimator,
		Storage:           storage,
		DbTx:              dbTx,
	}

	return msv, mks, ethClient
}

func getDefaultConfig() Config {
	cfg := Config{
		Host:                      host,
		Port:                      8123,
		MaxRequestsPerIPAndSecond: maxRequestsPerIPAndSecond,
		DefaultSenderAddress:      "0x1111111111111111111111111111111111111111",
		MaxCumulativeGasUsed:      300000,
		ChainID:                   1000,
	}
	return cfg
}

func newSequencerMockedServer(t *testing.T) (*mockedServer, *mocks, *ethclient.Client) {
	cfg := getDefaultConfig()
	return newMockedServer(t, cfg)
}

func newNonSequencerMockedServer(t *testing.T, sequencerNodeURI string) (*mockedServer, *mocks, *ethclient.Client) {
	cfg := getDefaultConfig()
	cfg.Port = 8124
	cfg.SequencerNodeURI = sequencerNodeURI
	return newMockedServer(t, cfg)
}

func (s *mockedServer) Stop() {
	err := s.Server.Stop()
	if err != nil {
		panic(err)
	}
}

func (s *mockedServer) JSONRPCCall(method string, parameters ...interface{}) (Response, error) {
	return JSONRPCCall(s.ServerURL, method, parameters...)
}
