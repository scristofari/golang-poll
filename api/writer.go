package api

import (
	"fmt"
	"net/http"
	"time"
)

// New ResponseWriter
// Needs :
// - status
// - size
// - start
type LogResponseWriter struct {
	status int
	size   int
	start  time.Time
	http.ResponseWriter
}

// Create a new middleware logger
func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{
		ResponseWriter: w,
		start:          time.Now(),
	}
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
	ms := fmt.Sprintf("%.2fms", time.Since(w.start).Seconds()*1000)
	return fmt.Sprintf("%s %s %s %d %d %s", r.Method, r.URL.String(), r.Proto, w.status, w.size, ms)
}
