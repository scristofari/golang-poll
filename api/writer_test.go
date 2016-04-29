package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Bytes to print
var (
	text = []byte("test")
)

// Mock a DefaultWriter
type WriterMock struct {
	mock.Mock
	http.ResponseWriter
}

// Define WriteHeader Method to be executed
func (w *WriterMock) WriteHeader(status int) {}

// Define WriteHeader Method to be executed
func (w *WriterMock) Write(bytes []byte) (int, error) {
	args := w.Called(bytes)
	return args.Int(0), nil
}

func TestLogWriter(t *testing.T) {
	w := new(WriterMock)
	w.On("WriteHeader", http.StatusOK)
	w.On("Write", text).Return(len(text))

	writer := LogResponseWriter{ResponseWriter: w}
	writer.WriteHeader(http.StatusOK)
	writer.Write(text)

	assert := assert.New(t)
	assert.Equal(http.StatusOK, writer.Status(), "Get status code")
	assert.Equal(len(text), writer.Size(), "Get size")
}
