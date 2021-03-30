package utils

import (
	"io"
	"net/http"
	"time"
)

type DefaultHttpClient struct {
	Timeout time.Duration
}

type HttpClient interface {
	Call(req *http.Request) (*http.Response, error)
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
}

func (c *DefaultHttpClient) Call(req *http.Request) (*http.Response, error) {
	client := http.Client{
		Timeout: c.Timeout,
	}
	return client.Do(req)
}

func (c *DefaultHttpClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return &http.Request{}, nil
}
