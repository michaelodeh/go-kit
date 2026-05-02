package llmstream

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (h *Client) handleRequest(req *http.Request, headers map[string]string) (*http.Response, error) {

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (h *Client) RequestWithContext(ctx context.Context, url string, method string, body io.Reader, headers map[string]string) (*http.Response, error) {

	if url == "" {
		return nil, fmt.Errorf("url not set")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return h.handleRequest(req, headers)
}

func (h *Client) Request(url string, method string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return h.RequestWithContext(context.Background(), url, method, body, headers)
}
