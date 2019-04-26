package core

import "net/http"

// LogWriter writes logs
type LogWriter struct {
	http.ResponseWriter
	status int
	length int
	body   []byte
}

// WriteHeader writes to an response's header
func (w *LogWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *LogWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(body)
	w.length += n
	w.body = body
	return n, err
}
