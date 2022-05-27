package jsonrpc

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientVersion(t *testing.T) {
	web3Endpoints := Web3{}

	result, err := web3Endpoints.ClientVersion()
	require.NoError(t, err)

	assert.Equal(t, "Polygon Hermez/v1.5.0", result)
}

func TestSha3(t *testing.T) {
	web3Endpoints := Web3{}

	helloWorld := argBig{}
	err := helloWorld.UnmarshalText([]byte("0x68656c6c6f20776f726c64"))
	require.NoError(t, err)

	resultInterface, err := web3Endpoints.Sha3(helloWorld)
	require.NoError(t, err)

	resultArgBytes := resultInterface.(argBytes)
	resultBytes := []byte(resultArgBytes)

	resultHex := common.Bytes2Hex(resultBytes)
	assert.Equal(t, resultHex, "47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad")
}
