package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// JSONRPC is an API backend
type Server struct {
	config  Config
	handler *JSONRPCHandler
}

// NewServer returns the JsonRPC server
func NewServer(config Config, p pool.Pool, s state.State) *Server {
	ethEndpoints := &Eth{chainID: config.ChainID, pool: p, state: s}
	netEndpoints := &Net{chainID: config.ChainID}

	handler := newJSONRpcHandler(ethEndpoints, netEndpoints)

	srv := &Server{
		config:  config,
		handler: handler,
	}
	return srv
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Infof("http server started: %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener: %v", err)
		return err
	}

	mux := http.DefaultServeMux
	mux.HandleFunc("/", s.handle)

	srv := http.Server{
		Handler: mux,
	}
	if err := srv.Serve(lis); err != nil {
		log.Errorf("closed http connection: %v", err)
		return err
	}
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
		w.Write([]byte("Hermez JSON-RPC"))
		return
	}

	if req.Method != "POST" {
		w.Write([]byte("method " + req.Method + " not allowed"))
		return
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	single, err := s.isSingleRequest(data)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if single {
		s.handleSingleRequest(w, data)
	} else {
		s.handleBatchRequest(w, data)
	}
}

func (s *Server) isSingleRequest(data []byte) (bool, error) {
	x := bytes.TrimLeft(data, " \t\r\n")

	if len(x) == 0 {
		return false, NewInvalidRequestError("Invalid json request")
	}

	return x[0] == '{', nil
}

func (s *Server) handleSingleRequest(w http.ResponseWriter, data []byte) {
	request, err := s.parseRequest(data)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	response := s.handler.Handle(request)

	respBytes, _ := response.Bytes()
	w.Write(respBytes)
}

func (s *Server) handleBatchRequest(w http.ResponseWriter, data []byte) {
	requests, err := s.parseRequests(data)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	responses := make([]Response, 0, len(requests))

	for _, request := range requests {
		response := s.handler.Handle(request)
		responses = append(responses, response)
	}

	respBytes, _ := json.Marshal(responses)
	w.Write(respBytes)
}

func (s *Server) parseRequest(data []byte) (Request, error) {
	var req Request

	if err := json.Unmarshal(data, &req); err != nil {
		return Request{}, NewInvalidRequestError("Invalid json request")
	}

	return req, nil
}

func (s *Server) parseRequests(data []byte) ([]Request, error) {
	var requests []Request

	if err := json.Unmarshal(data, &requests); err != nil {
		return nil, NewInvalidRequestError("Invalid json request")
	}

	return requests, nil
}
