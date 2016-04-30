package api

import (
	"fmt"
	"net/http"
)

// New ResponseWriter
// Needs :
// - status
// - size
type LogResponseWriter struct {
	status int
	size   int
	http.ResponseWriter
}

// Get the status
func (w *LogResponseWriter) Status() int {
	return w.status
}

// Get the size
func (w *LogResponseWriter) Size() int {
	return w.size
}

// Store the size of the repsonse
func (w *LogResponseWriter) Write(data []byte) (int, error) {
	var err error
	w.size, err = w.ResponseWriter.Write(data)
	return w.size, err
}

// Store the status of the response
func (w *LogResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Get the size
func (w *LogResponseWriter) String(r *http.Request) string {
	return fmt.Sprintf("%s %s %s %d %d", r.Method, r.URL.String(), r.Proto, w.Status(), w.Size())
}
