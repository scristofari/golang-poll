package test

import (
	"encoding/json"
	"fmt"
	"golang-poll/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert := assert.New(t)
	assert.Equal(res.StatusCode, http.StatusOK, "Bad request status !")
	assert.NotEmpty(res.Body, "Must return a response !")

	dec := json.NewDecoder(res.Body)
	defer res.Body.Close()
	result := new(api.ResultList)
	assert.Nil(dec.Decode(&result), "The response must be of type json")
	assert.IsType(&api.ResultList{}, result, "Must be a struct of type ResultList")
}
