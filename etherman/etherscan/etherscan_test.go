package etherscan

import (
	"context"
	"math/big"
	"testing"
	"net/http"
	"bytes"
	"io/ioutil"

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
	apiKey := ""
	c := NewEtherscanService(apiKey)
	httpM := new(httpMock)
	c.Http = httpM
	data := []byte(`{"status":"1","message":"OK","result":{"LastBlock":"15816910","SafeGasPrice":"10","ProposeGasPrice":"11","FastGasPrice":"55","suggestBaseFee":"9.849758735","gasUsedRatio":"0.779364333333333,0.2434028,0.610012833333333,0.1246597,0.995500566666667"}}`)
	body := ioutil.NopCloser(bytes.NewReader(data))
	httpM.On("Get", "https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=").Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

	gp, err := c.GetGasPrice(ctx)
	require.NoError(t, err)
	log.Debug("Etherscan GasPrice: ", gp)
	assert.Equal(t, big.NewInt(55000000000), gp)
}
