package core

import "net/http"

// HTTPLogWriter writes logs
type HTTPLogWriter struct {
	http.ResponseWriter
	status int
	length int
	body   []byte
}

// WriteHeader writes to an response's header
func (w *HTTPLogWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Write writes a response
func (w *HTTPLogWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(body)
	w.length += n
	w.body = body
	return n, err
}
