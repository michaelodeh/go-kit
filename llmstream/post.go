package llmstream

import (
	"context"
	"io"
	"net/http"
)

func (h *Client) PostWithContext(ctx context.Context, url string, payload io.Reader, headers map[string]string) (*http.Response, error) {
	return h.RequestWithContext(ctx, url, http.MethodPost, payload, headers)
}

func (h *Client) Post(url string, payload io.Reader, headers map[string]string) (*http.Response, error) {
	return h.Request(url, http.MethodPost, payload, headers)
}
