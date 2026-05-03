package llmstream

import (
	"context"
	"fmt"
	"net/http"
)

func (s *Client) ProxyStreamWithContext(
	ctx context.Context,
	w http.ResponseWriter,
	resp *http.Response,
	headers map[string]string,
) error {
	defer resp.Body.Close()

	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming not supported")
	}

	hopByHop := map[string]struct{}{
		"Connection":          {},
		"Keep-Alive":          {},
		"Proxy-Authenticate":  {},
		"Proxy-Authorization": {},
		"Te":                  {},
		"Trailer":             {},
		"Transfer-Encoding":   {},
		"Upgrade":             {},
	}

	for k, values := range resp.Header {
		if _, skip := hopByHop[k]; skip {
			continue
		}
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}

	for k, v := range headers {
		w.Header().Set(k, v)
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.WriteHeader(resp.StatusCode)

	return s.StreamWithContext(ctx, resp.Body, FlushWriter{
		W:       w,
		Flusher: flusher,
	})
}

func (h *Client) ProxyStream(
	w http.ResponseWriter,
	resp *http.Response,
	headers map[string]string,
) error {
	return h.ProxyStreamWithContext(context.Background(), w, resp, headers)
}
