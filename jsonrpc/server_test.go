package jsonrpc

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/mocks"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
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
	Pool     *mocks.PoolMock
	State    *mocks.StateMock
	Etherman *mocks.EthermanMock
	Storage  *storageMock
	DbTx     *mocks.DBTxMock
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
			Service: NewNetEndpoints(cfg, chainID),
		})
	}

	if _, ok := apis[APIZKEVM]; ok {
		services = append(services, Service{
			Name:    APIZKEVM,
			Service: NewZKEVMEndpoints(cfg, st, etherman),
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
			Service: NewDebugEndpoints(cfg, st, etherman),
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
		Pool:     pool,
		State:    st,
		Etherman: etherman,
		Storage:  storage,
		DbTx:     dbTx,
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

func newSequencerMockedServerWithCustomConfig(t *testing.T, cfg Config) (*mockedServer, *mocksWrapper, *ethclient.Client) {
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

func (s *mockedServer) JSONRPCBatchCall(calls ...client.BatchCall) ([]types.Response, error) {
	return client.JSONRPCBatchCall(s.ServerURL, calls...)
}

func (s *mockedServer) ChainID() uint64 {
	return chainID
}

func TestBatchRequests(t *testing.T) {
	type testCase struct {
		Name                 string
		BatchRequestsEnabled bool
		BatchRequestsLimit   uint
		NumberOfRequests     int
		ExpectedError        error
		SetupMocks           func(m *mocksWrapper, tc testCase)
	}

	block := ethTypes.NewBlock(
		&ethTypes.Header{Number: big.NewInt(2), UncleHash: ethTypes.EmptyUncleHash, Root: ethTypes.EmptyRootHash},
		[]*ethTypes.Transaction{ethTypes.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})},
		nil,
		[]*ethTypes.Receipt{ethTypes.NewReceipt([]byte{}, false, uint64(0))},
		&trie.StackTrie{},
	)

	testCases := []testCase{
		{
			Name:                 "batch requests disabled",
			BatchRequestsEnabled: false,
			BatchRequestsLimit:   0,
			NumberOfRequests:     10,
			ExpectedError:        types.ErrBatchRequestsDisabled,
			SetupMocks:           func(m *mocksWrapper, tc testCase) {},
		},
		{
			Name:                 "batch requests over the limit",
			BatchRequestsEnabled: true,
			BatchRequestsLimit:   5,
			NumberOfRequests:     6,
			ExpectedError:        types.ErrBatchRequestsLimitExceeded,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
			},
		},
		{
			Name:                 "batch requests unlimited",
			BatchRequestsEnabled: true,
			BatchRequestsLimit:   0,
			NumberOfRequests:     100,
			ExpectedError:        nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.DbTx.On("Commit", context.Background()).Return(nil).Times(tc.NumberOfRequests)
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Times(tc.NumberOfRequests)
				m.State.On("GetLastL2BlockNumber", context.Background(), m.DbTx).Return(block.Number().Uint64(), nil).Times(tc.NumberOfRequests)
				m.State.On("GetL2BlockByNumber", context.Background(), block.Number().Uint64(), m.DbTx).Return(block, nil).Times(tc.NumberOfRequests)
				m.State.On("GetTransactionReceipt", context.Background(), mock.Anything, m.DbTx).Return(ethTypes.NewReceipt([]byte{}, false, uint64(0)), nil)
			},
		},
		{
			Name:                 "batch requests equal the limit",
			BatchRequestsEnabled: true,
			BatchRequestsLimit:   5,
			NumberOfRequests:     5,
			ExpectedError:        nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.DbTx.On("Commit", context.Background()).Return(nil).Times(tc.NumberOfRequests)
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Times(tc.NumberOfRequests)
				m.State.On("GetLastL2BlockNumber", context.Background(), m.DbTx).Return(block.Number().Uint64(), nil).Times(tc.NumberOfRequests)
				m.State.On("GetL2BlockByNumber", context.Background(), block.Number().Uint64(), m.DbTx).Return(block, nil).Times(tc.NumberOfRequests)
				m.State.On("GetTransactionReceipt", context.Background(), mock.Anything, m.DbTx).Return(ethTypes.NewReceipt([]byte{}, false, uint64(0)), nil)
			},
		},
		{
			Name:                 "batch requests under the limit",
			BatchRequestsEnabled: true,
			BatchRequestsLimit:   5,
			NumberOfRequests:     4,
			ExpectedError:        nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.DbTx.On("Commit", context.Background()).Return(nil).Times(tc.NumberOfRequests)
				m.State.On("BeginStateTransaction", context.Background()).Return(m.DbTx, nil).Times(tc.NumberOfRequests)
				m.State.On("GetLastL2BlockNumber", context.Background(), m.DbTx).Return(block.Number().Uint64(), nil).Times(tc.NumberOfRequests)
				m.State.On("GetL2BlockByNumber", context.Background(), block.Number().Uint64(), m.DbTx).Return(block, nil).Times(tc.NumberOfRequests)
				m.State.On("GetTransactionReceipt", context.Background(), mock.Anything, m.DbTx).Return(ethTypes.NewReceipt([]byte{}, false, uint64(0)), nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase

			cfg := getDefaultConfig()
			cfg.BatchRequestsEnabled = tc.BatchRequestsEnabled
			cfg.BatchRequestsLimit = tc.BatchRequestsLimit
			s, m, _ := newSequencerMockedServerWithCustomConfig(t, cfg)
			defer s.Stop()

			tc.SetupMocks(m, tc)

			calls := []client.BatchCall{}

			for i := 0; i < tc.NumberOfRequests; i++ {
				calls = append(calls, client.BatchCall{
					Method:     "eth_getBlockByNumber",
					Parameters: []interface{}{"latest"},
				})
			}

			result, err := s.JSONRPCBatchCall(calls...)
			if testCase.ExpectedError == nil {
				assert.Equal(t, testCase.NumberOfRequests, len(result))
			} else {
				assert.Equal(t, 0, len(result))
				assert.Equal(t, testCase.ExpectedError.Error(), err.Error())
			}
		})
	}
}
