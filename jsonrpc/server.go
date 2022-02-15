package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/gaspriceestimator"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state"
)

// Server is an API backend to handle RPC requests
type Server struct {
	config  Config
	handler *Handler
	srv     *http.Server
}

// NewServer returns the JsonRPC server
func NewServer(
	config Config,
	defaultChainID uint64,
	sequencerAddress common.Address,
	p jsonRPCTxPool,
	s state.State,
	chainID uint64,
	gpe gaspriceestimator.GasPriceEstimator) *Server {
	chainIDSelector := newChainIDSelector(chainID)
	ethEndpoints := &Eth{
		chainIDSelector:  chainIDSelector,
		pool:             p,
		state:            s,
		gpe:              gpe,
		sequencerAddress: sequencerAddress,
	}
	netEndpoints := &Net{chainIDSelector: chainIDSelector}
	hezEndpoints := &Hez{defaultChainID: defaultChainID, state: s}
	txPoolEndpoints := &TxPool{pool: p}
	debugEndpoints := &Debug{state: s}

	handler := newJSONRpcHandler(ethEndpoints, netEndpoints, hezEndpoints, txPoolEndpoints, debugEndpoints)

	srv := &Server{
		config:  config,
		handler: handler,
	}
	return srv
}

// Start initializes the JSON RPC server to listen for request
func (s *Server) Start() error {
	if s.srv != nil {
		return fmt.Errorf("server already started")
	}

	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Infof("http server started: %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener: %v", err)
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)

	s.srv = &http.Server{
		Handler: mux,
	}
	if err := s.srv.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Infof("http server stopped")
			return nil
		}
		log.Errorf("closed http connection: %v", err)
		return err
	}
	return nil
}

// Stop shutdown the rpc server
func (s *Server) Stop() error {
	if s.srv == nil {
		return nil
	}

	if err := s.srv.Shutdown(context.Background()); err != nil {
		return err
	}

	if err := s.srv.Close(); err != nil {
		return err
	}

	s.srv = nil

	return nil
}

func (s *Server) handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if (*req).Method == "OPTIONS" {
		return
	}

	if req.Method == "GET" {
		_, err := w.Write([]byte("Hermez JSON-RPC"))
		if err != nil {
			log.Error(err)
		}
		return
	}

	if req.Method != "POST" {
		err := errors.New("method " + req.Method + " not allowed")
		handleError(w, err)
		return
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	single, err := s.isSingleRequest(data)
	if err != nil {
		handleError(w, err)
		return
	}

	if single {
		s.handleSingleRequest(w, data)
	} else {
		s.handleBatchRequest(w, data)
	}
}

func (s *Server) isSingleRequest(data []byte) (bool, detailedError) {
	x := bytes.TrimLeft(data, " \t\r\n")

	if len(x) == 0 {
		return false, newInvalidRequestError("Invalid json request")
	}

	return x[0] == '{', nil
}

func (s *Server) handleSingleRequest(w http.ResponseWriter, data []byte) {
	request, err := s.parseRequest(data)
	if err != nil {
		handleError(w, err)
		return
	}

	response := s.handler.Handle(request)

	respBytes, _ := response.Bytes()
	_, err = w.Write(respBytes)
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) handleBatchRequest(w http.ResponseWriter, data []byte) {
	requests, err := s.parseRequests(data)
	if err != nil {
		handleError(w, err)
		return
	}

	responses := make([]Response, 0, len(requests))

	for _, request := range requests {
		response := s.handler.Handle(request)
		responses = append(responses, response)
	}

	respBytes, _ := json.Marshal(responses)
	_, err = w.Write(respBytes)
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) parseRequest(data []byte) (Request, error) {
	var req Request

	if err := json.Unmarshal(data, &req); err != nil {
		return Request{}, newInvalidRequestError("Invalid json request")
	}

	return req, nil
}

func (s *Server) parseRequests(data []byte) ([]Request, error) {
	var requests []Request

	if err := json.Unmarshal(data, &requests); err != nil {
		return nil, newInvalidRequestError("Invalid json request")
	}

	return requests, nil
}

func handleError(w http.ResponseWriter, err error) {
	log.Error(err)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		log.Error(err)
	}
}
