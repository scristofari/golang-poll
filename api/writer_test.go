package api

import (
	"net/http"
	"net/url"
	"testing"
	"time"

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

// Define Write Method which will be executed.
func (w *DefaultWriter) Write(bytes []byte) (int, error) {
	return len(text), nil
}

func TestLogResponseWriter(t *testing.T) {
	w := new(DefaultWriter)

	writer := LogResponseWriter{ResponseWriter: w}
	writer.WriteHeader(http.StatusOK)
	writer.Write(text)
	writer.start = time.Now()

	r := &http.Request{
		Proto:  "HTTP/1.1",
		Method: "GET",
		URL: &url.URL{
			Path: "/api/v1/polls",
		},
	}

	assert := assert.New(t)
	assert.Equal("GET /api/v1/polls HTTP/1.1 200 4 0.00ms", writer.String(r), "Log access request/response")
}
