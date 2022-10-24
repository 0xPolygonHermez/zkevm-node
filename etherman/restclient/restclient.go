package restclient

// This package has been created to mock the http requests in unit tests
import (
	"net/http"
)

// HttpI is the restClient interface
type HttpI interface {
	Get(url string) (*http.Response, error)
}

// Client for restHttp
type Client struct {
}

// NewClient is the constructor that creates an etherscanService
func NewClient() *Client {
	return &Client{}
}

// Get sends a Get request to the URL with the body
func (c *Client) Get(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(request)
}
