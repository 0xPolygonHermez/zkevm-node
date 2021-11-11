package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

// JSONRPC is an API backend
type HTTPServer struct {
	config  Config
	handler *JSONRpcHandler
}

// NewHTTPServer returns the JsonRPC http server
func NewHTTPServer(config Config) *HTTPServer {
	srv := &HTTPServer{
		config:  config,
		handler: newJSONRpcHandler(config.ChainID),
	}
	return srv
}

func (s *HTTPServer) Start() error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Println("http server started", "addr", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("failed to create tcp listener", "err", err)
		return err
	}

	mux := http.DefaultServeMux
	mux.HandleFunc("/", s.handleHTTPRpc)

	srv := http.Server{
		Handler: mux,
	}
	if err := srv.Serve(lis); err != nil {
		fmt.Println("closed http connection", "err", err)
		return err
	}
	return nil
}

func (s *HTTPServer) handleHTTPRpc(w http.ResponseWriter, req *http.Request) {
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

func (s *HTTPServer) isSingleRequest(data []byte) (bool, error) {
	x := bytes.TrimLeft(data, " \t\r\n")

	if len(x) == 0 {
		return false, NewInvalidRequestError("Invalid json request")
	}

	return x[0] == '{', nil
}

func (s *HTTPServer) handleSingleRequest(w http.ResponseWriter, data []byte) {
	request, err := s.parseRequest(data)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	response := s.handler.Handle(request)

	respBytes, _ := response.Bytes()
	w.Write(respBytes)
}

func (s *HTTPServer) handleBatchRequest(w http.ResponseWriter, data []byte) {
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

func (s *HTTPServer) parseRequest(data []byte) (Request, error) {
	var req Request

	if err := json.Unmarshal(data, &req); err != nil {
		return Request{}, NewInvalidRequestError("Invalid json request")
	}

	return req, nil
}

func (s *HTTPServer) parseRequests(data []byte) ([]Request, error) {
	var requests []Request

	if err := json.Unmarshal(data, &requests); err != nil {
		return nil, NewInvalidRequestError("Invalid json request")
	}

	return requests, nil
}
