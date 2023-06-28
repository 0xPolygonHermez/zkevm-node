package jsonrpc

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/mocks"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	maxRequestsPerIPAndSecond        = 1000
	chainID                   uint64 = 1000
)

type mockedServer struct {
	Config    Config
	Server    *Server
	ServerURL string
}

type mocksWrapper struct {
	Pool    *mocks.PoolMock
	State   *mocks.StateMock
	Storage *storageMock
	DbTx    *mocks.DBTxMock
}

func newMockedServer(t *testing.T, cfg Config) (*mockedServer, *mocksWrapper, *ethclient.Client) {
	pool := mocks.NewPoolMock(t)
	st := mocks.NewStateMock(t)
	etherman := mocks.NewEthermanMock(t)
	storage := newStorageMock(t)
	dbTx := mocks.NewDBTxMock(t)
	apis := map[string]bool{
		APIEth:    true,
		APINet:    true,
		APIDebug:  true,
		APIZKEVM:  true,
		APITxPool: true,
		APIWeb3:   true,
	}

	var newL2BlockEventHandler state.NewL2BlockEventHandler = func(e state.NewL2BlockEvent) {}
	st.On("RegisterNewL2BlockEventHandler", mock.IsType(newL2BlockEventHandler)).Once()
	st.On("PrepareWebSocket").Once()

	services := []Service{}
	if _, ok := apis[APIEth]; ok {
		services = append(services, Service{
			Name:    APIEth,
			Service: NewEthEndpoints(cfg, chainID, pool, st, etherman, storage),
		})
	}

	if _, ok := apis[APINet]; ok {
		services = append(services, Service{
			Name:    APINet,
			Service: NewNetEndpoints(chainID),
		})
	}

	if _, ok := apis[APIZKEVM]; ok {
		services = append(services, Service{
			Name:    APIZKEVM,
			Service: NewZKEVMEndpoints(st),
		})
	}

	if _, ok := apis[APITxPool]; ok {
		services = append(services, Service{
			Name:    APITxPool,
			Service: &TxPoolEndpoints{},
		})
	}

	if _, ok := apis[APIDebug]; ok {
		services = append(services, Service{
			Name:    APIDebug,
			Service: NewDebugEndpoints(st, etherman),
		})
	}

	if _, ok := apis[APIWeb3]; ok {
		services = append(services, Service{
			Name:    APIWeb3,
			Service: &Web3Endpoints{},
		})
	}
	server := NewServer(cfg, chainID, pool, st, storage, services)

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

	mks := &mocksWrapper{
		Pool:    pool,
		State:   st,
		Storage: storage,
		DbTx:    dbTx,
	}

	return msv, mks, ethClient
}

func getDefaultConfig() Config {
	cfg := Config{
		Host:                      "0.0.0.0",
		Port:                      9123,
		MaxRequestsPerIPAndSecond: maxRequestsPerIPAndSecond,
		MaxCumulativeGasUsed:      300000,
	}
	return cfg
}

func newSequencerMockedServer(t *testing.T) (*mockedServer, *mocksWrapper, *ethclient.Client) {
	cfg := getDefaultConfig()
	return newMockedServer(t, cfg)
}

func newNonSequencerMockedServer(t *testing.T, sequencerNodeURI string) (*mockedServer, *mocksWrapper, *ethclient.Client) {
	cfg := getDefaultConfig()
	cfg.Port = 9124
	cfg.SequencerNodeURI = sequencerNodeURI
	return newMockedServer(t, cfg)
}

func (s *mockedServer) Stop() {
	err := s.Server.Stop()
	if err != nil {
		panic(err)
	}
}

func (s *mockedServer) JSONRPCCall(method string, parameters ...interface{}) (types.Response, error) {
	return client.JSONRPCCall(s.ServerURL, method, parameters...)
}

func (s *mockedServer) ChainID() uint64 {
	return chainID
}
