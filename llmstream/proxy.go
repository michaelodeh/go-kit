package llmstream

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (h *Client) ProxyStreamWithContext(
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

	// filter hop-by-hop headers
	hopByHop := map[string]struct{}{
		"Connection": {}, "Keep-Alive": {}, "Proxy-Authenticate": {},
		"Proxy-Authorization": {}, "Te": {}, "Trailer": {},
		"Transfer-Encoding": {}, "Upgrade": {},
	}

	// 1. upstream headers
	for k, v := range resp.Header {
		if _, skip := hopByHop[k]; skip {
			continue
		}
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	// 2. override headers
	for key, value := range headers {
		w.Header().Set(key, value)
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.WriteHeader(resp.StatusCode)

	buf := make([]byte, 32*1024)

	go func() {
		<-ctx.Done()
		resp.Body.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := resp.Body.Read(buf)

		if n > 0 {
			if _, err := w.Write(buf[:n]); err != nil {
				return err
			}
			flusher.Flush()
		}

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (h *Client) ProxyStream(
	w http.ResponseWriter,
	resp *http.Response,
	headers map[string]string,
) error {
	return h.ProxyStreamWithContext(context.Background(), w, resp, headers)
}
