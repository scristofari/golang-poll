package api

import "net/http"

type LogResponseWriter struct {
	status int
	size   int
	http.ResponseWriter
}

func (w *LogResponseWriter) Status() int {
	return w.status
}

func (w *LogResponseWriter) Size() int {
	return w.size
}

func (w *LogResponseWriter) Write(data []byte) (int, error) {
	var err error
	w.size, err = w.ResponseWriter.Write(data)
	return w.size, err
}

func (w *LogResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
