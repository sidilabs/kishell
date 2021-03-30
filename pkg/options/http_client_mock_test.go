package options

import (
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"time"
)

type MockHttpClient struct {
	mock.Mock
	Timeout time.Duration
	Request http.Request
}

func (c *MockHttpClient) Call(req *http.Request) (*http.Response, error) {
	args := c.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (c *MockHttpClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	args := c.Called(method, url, body)
	return args.Get(0).(*http.Request), args.Error(1)
}
