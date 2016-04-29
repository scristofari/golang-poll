package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	server   *httptest.Server
	pollsURL string
)

func init() {
	server = httptest.NewServer(Handlers())
	pollsURL = fmt.Sprintf("%s/api/v1/polls", server.URL)
}

func TestListPolls(t *testing.T) {
	request, err := http.NewRequest("GET", pollsURL, nil)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}

	assert := assert.New(t)
	assert.Equal(res.StatusCode, http.StatusOK, "Bad request status !")
	assert.NotEmpty(res.Body, "Must return a response !")
	assert.Equal("*", res.Header.Get("Access-Control-Allow-Origin"), "CORS header *")
	assert.Equal("application/json", res.Header.Get("Content-Type"), "Content Type header")
}
