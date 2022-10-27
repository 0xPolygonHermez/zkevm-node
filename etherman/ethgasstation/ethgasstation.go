package ethgasstation

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/etherman/restclient"
)

type ethGasStationResponse struct {
	BaseFee     uint64                `json:"baseFee"`
	BlockNumber uint64                `json:"blockNumber"`
	GasPrice    gasPriceEthGasStation `json:"gasPrice"`
}

// gasPriceEthGasStation definition
type gasPriceEthGasStation struct {
	Standard uint64 `json:"standard"`
	Instant  uint64 `json:"instant"`
	Fast     uint64 `json:"fast"`
}

// Client for ethGasStation
type Client struct {
	Http restclient.HttpI
}

// NewEthGasStationService is the constructor that creates an ethGasStationService
func NewEthGasStationService() *Client {
	return &Client{
		Http: restclient.NewClient(),
	}
}

// GetGasPrice retrieves the gas price estimation from ethGasStation
func (e *Client) GetGasPrice(ctx context.Context) (*big.Int, error) {
	var resBody ethGasStationResponse
	url := "https://api.ethgasstation.info/api/fee-estimate"
	res, err := e.Http.Get(url)
	if err != nil {
		return big.NewInt(0), err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return big.NewInt(0), err
	}
	if res.StatusCode != http.StatusOK {
		return big.NewInt(0), fmt.Errorf("http response is %d", res.StatusCode)
	}
	// Unmarshal result
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return big.NewInt(0), fmt.Errorf("Reading body failed: %w", err)
	}
	fgp := big.NewInt(0).SetUint64(resBody.GasPrice.Instant)
	return new(big.Int).Mul(fgp, big.NewInt(encoding.Gwei)), nil
}
