package api

import (
	"bytes"
	"encoding/json"
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

func TestNotFound(t *testing.T) {
	r, _ := http.NewRequest("GET", server.URL+"notfound", nil)
	w := httptest.NewRecorder()
	notFoundHandler(w, r)
	assert := assert.New(t)
	assert.Equal(http.StatusNotFound, w.Code, "Status must be 404")
}

func TestListPolls(t *testing.T) {
	request, err := http.NewRequest("GET", pollsURL, nil)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}

	assert := assert.New(t)
	assert.Equal(res.StatusCode, http.StatusOK, "Bad request status !")
	defer res.Body.Close()
	assert.NotEmpty(res.Body, "Must return a response !")
	assert.Equal("*", res.Header.Get("Access-Control-Allow-Origin"), "CORS header *")
	assert.Equal("application/json", res.Header.Get("Content-Type"), "Content Type header")
}

func TestPostPoll(t *testing.T) {
	pollTest := `{"name":"test poll","question":"Is it true or false ?","answers":[{"label":"true","votes":0},{"label":"false","votes":0}]}`
	body := bytes.NewBuffer([]byte(pollTest))
	request, err := http.NewRequest("POST", pollsURL, body)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	poll := new(Poll)
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()

	if err := decoder.Decode(poll); err != nil {
		assert.Error(t, err, "Decoding fail, not poll ?")
	}

	assert := assert.New(t)
	assert.Equal(res.StatusCode, http.StatusCreated, "Bad request status !")
	assert.NotEmpty(poll.ID, "Poll not saved ?")
	assert.NotEmpty(poll.CreatedAt, "Created")
	assert.NotEmpty(poll.UpdatedAt, "Updated")
	assert.Equal("Test Poll", poll.Name, "Name")
	assert.Equal(2, len(poll.Answers), "answers")
}
