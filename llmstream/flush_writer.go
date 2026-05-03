package llmstream

import (
	"net/http"
)

type FlushWriter struct {
	W       http.ResponseWriter
	Flusher http.Flusher
}

func (f FlushWriter) Write(p []byte) (int, error) {
	n, err := f.W.Write(p)
	if err == nil {
		f.Flusher.Flush()
	}
	return n, err
}
