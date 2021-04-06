package utils

import (
	"io"
	"net/http"
	"time"
)

// A DefaultHTTPClient represents properties to be used during a HTTP call.
type DefaultHTTPClient struct {
	Timeout time.Duration
}

// A HTTPClient represents a contract to make http requests.
type HTTPClient interface {
	Call(req *http.Request) (*http.Response, error)
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
}

// Call places the actual the http request.
func (c *DefaultHTTPClient) Call(req *http.Request) (*http.Response, error) {
	client := http.Client{
		Timeout: c.Timeout,
	}
	return client.Do(req)
}

// NewRequest creates new HTTP requests.
func (c *DefaultHTTPClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}
