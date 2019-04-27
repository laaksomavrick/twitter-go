package core

import "net/http"

// HttpLogWriter writes logs
type HttpLogWriter struct {
	http.ResponseWriter
	status int
	length int
	body   []byte
}

// WriteHeader writes to an response's header
func (w *HttpLogWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *HttpLogWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(body)
	w.length += n
	w.body = body
	return n, err
}
