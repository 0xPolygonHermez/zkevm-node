package ethgasstation

import (
	"net/http"
)

// HttpI is the http interface
type HttpI interface {
	Get(url string) (*http.Response, error)
}
