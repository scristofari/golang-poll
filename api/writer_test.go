package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Bytes to print.
var (
	text = []byte("text")
)

// Mock a DefaultWriter.
type DefaultWriter struct {
	http.ResponseWriter
}

// Define WriteHeader Method which will be executed.
func (w *DefaultWriter) WriteHeader(status int) {}

// Define WriteHeader Method which will be executed.
func (w *DefaultWriter) Write(bytes []byte) (int, error) {
	return len(text), nil
}

func TestLogResponseWriter(t *testing.T) {
	w := new(DefaultWriter)

	writer := LogResponseWriter{ResponseWriter: w}
	writer.WriteHeader(http.StatusOK)
	writer.Write(text)

	assert := assert.New(t)
	assert.Equal(http.StatusOK, writer.Status(), "Get status code")
	assert.Equal(len(text), writer.Size(), "Get size")

	r := &http.Request{
		Proto:  "HTTP/1.1",
		Method: "GET",
		URL: &url.URL{
			Path: "/api/v1/polls",
		},
	}
	assert.Equal("GET /api/v1/polls HTTP/1.1 200 4", writer.String(r), "Log access request/response")
}
