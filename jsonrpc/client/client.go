package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

const jsonRPCVersion = "2.0"

// Client defines typed wrappers for the zkEVM RPC API.
type Client struct {
	url string
}

// NewClient creates an instance of client
func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}

// JSONRPCCall executes a 2.0 JSON RPC HTTP Post Request to the provided URL with
// the provided method and parameters, which is compatible with the Ethereum
// JSON RPC Server.
func JSONRPCCall(url, method string, parameters ...interface{}) (types.Response, error) {
	params, err := json.Marshal(parameters)
	if err != nil {
		return types.Response{}, err
	}

	request := types.Request{
		JSONRPC: jsonRPCVersion,
		ID:      float64(1),
		Method:  method,
		Params:  params,
	}

	httpRes, err := sendJSONRPC_HTTPRequest(url, request)
	if err != nil {
		return types.Response{}, err
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return types.Response{}, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		return types.Response{}, fmt.Errorf("%v - %v", httpRes.StatusCode, string(resBody))
	}

	var res types.Response
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return types.Response{}, err
	}
	return res, nil
}

// BatchCall used in batch requests to send multiple methods and parameters at once
type BatchCall struct {
	Method     string
	Parameters []interface{}
}

// JSONRPCBatchCall executes a 2.0 JSON RPC HTTP Post Batch Request to the provided URL with
// the provided method and parameters groups, which is compatible with the Ethereum
// JSON RPC Server.
func JSONRPCBatchCall(url string, calls ...BatchCall) ([]types.Response, error) {
	requests := []types.Request{}

	for i, call := range calls {
		params, err := json.Marshal(call.Parameters)
		if err != nil {
			return nil, err
		}

		req := types.Request{
			JSONRPC: jsonRPCVersion,
			ID:      float64(i),
			Method:  call.Method,
			Params:  params,
		}

		requests = append(requests, req)
	}

	httpRes, err := sendJSONRPC_HTTPRequest(url, requests)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v - %v", httpRes.StatusCode, string(resBody))
	}

	var res []types.Response
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		errorMessage := string(resBody)
		return nil, fmt.Errorf(errorMessage)
	}

	return res, nil
}

func sendJSONRPC_HTTPRequest(url string, payload interface{}) (*http.Response, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	reqBodyReader := bytes.NewReader(reqBody)
	httpReq, err := http.NewRequest(http.MethodPost, url, reqBodyReader)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Content-type", "application/json")

	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return httpRes, nil
}
