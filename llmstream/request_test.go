package llmstream

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRequestWithContextReturnsNon2xxResponses(t *testing.T) {
	client := NewClient(&http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Header: http.Header{
					"Content-Type": {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"error":"upstream unavailable"}`)),
			}, nil
		}),
	})
	resp, err := client.Post("http://llm.local/api", strings.NewReader(`{}`), map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status code = %d, want %d", resp.StatusCode, http.StatusBadGateway)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}
	if got, want := string(body), `{"error":"upstream unavailable"}`; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
