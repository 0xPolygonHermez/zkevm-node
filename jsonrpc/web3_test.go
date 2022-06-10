package jsonrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientVersion(t *testing.T) {
	s, _, _ := newMockedServer(t)
	defer s.Stop()

	reqBodyObj := Request{
		JSONRPC: "2.0",
		ID:      float64(1),
		Method:  "web3_clientVersion",
	}

	reqBody, err := json.Marshal(reqBodyObj)
	require.NoError(t, err)

	reqBodyReader := bytes.NewReader(reqBody)
	req, err := http.NewRequest(http.MethodPost, s.ServerURL, reqBodyReader)
	require.NoError(t, err)

	req.Header.Add("Content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	if res.StatusCode != http.StatusOK {
		log.Error("Invalid status code, expected: %v, found: %v", http.StatusOK, res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	defer res.Body.Close()

	var resBodyObj Response
	err = json.Unmarshal(resBody, &resBodyObj)
	require.NoError(t, err)

	assert.Equal(t, reqBodyObj.JSONRPC, resBodyObj.JSONRPC)
	assert.Equal(t, reqBodyObj.ID, resBodyObj.ID)
	assert.Nil(t, resBodyObj.Error)

	var result string
	err = json.Unmarshal(resBodyObj.Result, &result)
	require.NoError(t, err)

	assert.Equal(t, "Polygon Hermez/v1.5.0", result)
}

func TestSha3(t *testing.T) {
	s, _, _ := newMockedServer(t)
	defer s.Stop()

	params, err := json.Marshal([]string{"0x68656c6c6f20776f726c64"})
	require.NoError(t, err)

	reqBodyObj := Request{
		JSONRPC: "2.0",
		ID:      float64(1),
		Method:  "web3_sha3",
		Params:  params,
	}

	reqBody, err := json.Marshal(reqBodyObj)
	require.NoError(t, err)

	reqBodyReader := bytes.NewReader(reqBody)
	req, err := http.NewRequest(http.MethodPost, s.ServerURL, reqBodyReader)
	require.NoError(t, err)

	req.Header.Add("Content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	if res.StatusCode != http.StatusOK {
		log.Error("Invalid status code, expected: %v, found: %v", http.StatusOK, res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	defer res.Body.Close()

	var resBodyObj Response
	err = json.Unmarshal(resBody, &resBodyObj)
	require.NoError(t, err)

	assert.Equal(t, reqBodyObj.JSONRPC, resBodyObj.JSONRPC)
	assert.Equal(t, reqBodyObj.ID, resBodyObj.ID)
	assert.Nil(t, resBodyObj.Error)

	var result string
	err = json.Unmarshal(resBodyObj.Result, &result)
	require.NoError(t, err)

	assert.Equal(t, "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad", result)
}
