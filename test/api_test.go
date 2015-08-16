package test

import (
	"fmt"
	"golang-poll/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server   *httptest.Server
	pollsUrl string
)

func init() {
	server = httptest.NewServer(api.Handlers())
	pollsUrl = fmt.Sprintf("%s/api/v1/polls", server.URL)
}

func TestListPolls(t *testing.T) {
	request, err := http.NewRequest("GET", pollsUrl, nil)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d, got: %d", http.StatusOK, res.StatusCode)
	}
}
