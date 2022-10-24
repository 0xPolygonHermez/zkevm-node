package ethgasstation

import (
	"context"
	"math/big"
	"testing"
	"io/ioutil"
	"bytes"
	"net/http"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})
}

func TestGetGasPrice(t *testing.T) {
	ctx := context.Background()
	c := NewEthGasStationService()
	httpM := new(httpMock)
	c.Http = httpM
	data := []byte(`{"baseFee":10,"blockNumber":15817089,"blockTime":11.88,"gasPrice":{"fast":11,"instant":66,"standard":10},"nextBaseFee":10,"priorityFee":{"fast":2,"instant":2,"standard":1}}`)
	body := ioutil.NopCloser(bytes.NewReader(data))
	httpM.On("Get", "https://api.ethgasstation.info/api/fee-estimate").Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

	gp, err := c.GetGasPrice(ctx)
	require.NoError(t, err)
	log.Debug("EthGasStation GasPrice: ", gp)
	assert.Equal(t, big.NewInt(66000000000), gp)
}
