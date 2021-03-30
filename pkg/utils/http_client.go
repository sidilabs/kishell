package utils

import (
	"io"
	"net/http"
	"time"
)

// Default HTTP client.
type DefaultHttpClient struct {
	Timeout time.Duration
}


// HTTP client contract.
type HttpClient interface {
	Call(req *http.Request) (*http.Response, error)
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
}

// Makes the http call.
func (c *DefaultHttpClient) Call(req *http.Request) (*http.Response, error) {
	client := http.Client{
		Timeout: c.Timeout,
	}
	return client.Do(req)
}

// Creates new HTTP requests.
func (c *DefaultHttpClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return &http.Request{}, nil
}
