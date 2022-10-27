package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
)

type etherscanResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Result  gasPriceEtherscan `json:"result"`
}

// gasPriceEtherscan definition
type gasPriceEtherscan struct {
	LastBlock       string `json:"LastBlock"`
	SafeGasPrice    string `json:"SafeGasPrice"`
	ProposeGasPrice string `json:"ProposeGasPrice"`
	FastGasPrice    string `json:"FastGasPrice"`
}

// Config structure
type Config struct {
	ApiKey string `mapstructure:"ApiKey"`
}

// Client for etherscan
type Client struct {
	config Config
	Http   HttpI
}

// NewEtherscanService is the constructor that creates an etherscanService
func NewEtherscanService(apikey string) *Client {
	return &Client{
		config: Config{
			ApiKey: apikey,
		},
		Http: http.DefaultClient,
	}
}

// GetGasPrice retrieves the gas price estimation from etherscan
func (e *Client) GetGasPrice(ctx context.Context) (*big.Int, error) {
	var resBody etherscanResponse
	url := "https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=" + e.config.ApiKey
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
	fgp, _ := big.NewInt(0).SetString(resBody.Result.FastGasPrice, encoding.Base10)
	return new(big.Int).Mul(fgp, big.NewInt(encoding.Gwei)), nil
}
