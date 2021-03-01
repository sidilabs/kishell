package options

import (
  "bytes"
  "fmt"
  "github.com/go-http-utils/headers"
  "github.com/sidilabs/kishell/pkg/config"
  "github.com/stretchr/testify/mock"
  "io/ioutil"
  "net/http"
  "testing"
)

func TestMatchAll(t *testing.T) {
  err := testIt(t,
    `{"name":"Test Name","full_name":"test full name","owner":{"login": "octocat"}}`,
    &http.Response{
      Header: http.Header{
        headers.ContentType: []string { "application/json" },
      },
      StatusCode: 200,
    },
  )
  if err != nil {
    t.Fatal("Searching for documents must succeed", err)
  }
}

func TestInvalidJsonResponse(t *testing.T) {
  err := testIt(t,
    `{"name":"Test Name","full_name":"test full name","owner":{"login": "octocat"}`,
    &http.Response{
      Header: http.Header{
        headers.ContentType: []string { "application/json" },
      },
      StatusCode: 200,
    },
  )
  if err == nil {
    t.Fatal("Searching for document with invalid json has failed", err)
  }
}

func TestInvalidContentTypeResponse(t *testing.T) {
  err := testIt(t,
    `{"name":"Test Name","full_name":"test full name","owner":{"login": "octocat"}`,
    &http.Response{
      Header: http.Header{
        headers.ContentType: []string { "application/xml" },
      },
      StatusCode: 200,
    },
  )
  if err == nil {
    t.Fatal("Searching for document with invalid content type response header failed", err)
  }
}

func testIt(t *testing.T, source string, response *http.Response) error {
  serverName := "ut-server"
  configuration := new(ConfigurationMock)
  configuration.On("CheckEmpty").Return(nil)
  configuration.On("GetCurrentServer").Return(config.Server{})
  configuration.On("FindServer", serverName).Return(config.Server{
   Protocol: "http",
   Hostname: "ut.server",
  }, true)
  configuration.On("GetCurrentRole").Return(config.Role{})

  json := fmt.Sprintf(`{ "responses": [ { "hits": { "hits": [ { "_source": %s } ] } } ] }`, source)
  response.Body = ioutil.NopCloser(bytes.NewReader([]byte(json)))

  httpClient := new(MockHttpClient)
  httpClient.Request = http.Request{
    Header: map[string][]string{},
  }
  httpClient.On("NewRequest", mock.Anything, mock.Anything, mock.Anything).Return(&httpClient.Request, nil)
  httpClient.On("Call", &httpClient.Request).Return(response, nil)

  context := Context{
    Debug:         true,
    Configuration: configuration,
  }

  cmd := SearchCmd{
    Server: serverName,
    Limit: 30,
    httpClient: httpClient,
    Older: "30m",
  }

  err := cmd.Run(&context)
  configuration.AssertExpectations(t)
  httpClient.AssertExpectations(t)
  return err
}
