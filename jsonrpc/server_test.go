package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

type mockedServer struct {
	DefaultChainID   uint64
	ChainID          uint64
	SequencerAddress common.Address

	Server    *Server
	ServerURL string
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
		DefaultChainID:   defaultChainID,
		ChainID:          chainID,
		SequencerAddress: sequencerAddress,

		Server:    server,
		ServerURL: serverURL,
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

func (s *mockedServer) JSONRPCCall(method string, parameters ...interface{}) (Response, error) {
	params, err := json.Marshal(parameters)
	if err != nil {
		return Response{}, err
	}

	req := Request{
		JSONRPC: "2.0",
		ID:      float64(1),
		Method:  method,
		Params:  params,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return Response{}, err
	}

	reqBodyReader := bytes.NewReader(reqBody)
	httpReq, err := http.NewRequest(http.MethodPost, s.ServerURL, reqBodyReader)
	if err != nil {
		return Response{}, err
	}

	httpReq.Header.Add("Content-type", "application/json")

	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return Response{}, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("Invalid status code, expected: %v, found: %v", http.StatusOK, httpRes.StatusCode)
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return Response{}, err
	}
	defer httpRes.Body.Close()

	var res Response
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return Response{}, err
	}

	return res, nil
}
