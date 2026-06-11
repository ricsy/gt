package api

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("gitee.com", "test-token")
	if client == nil {
		t.Error("NewClient returned nil")
	}
}

func TestNewClientUsesDefaultTimeout(t *testing.T) {
	client := NewClient("gitee.com", "test-token")

	if client.HTTPClient.Timeout != DefaultTimeout {
		t.Fatalf("expected default timeout %s, got %s", DefaultTimeout, client.HTTPClient.Timeout)
	}
}

func TestNewClientWithTimeout(t *testing.T) {
	timeout := 5 * time.Second

	client := NewClientWithTimeout("gitee.com", "test-token", timeout)

	if client.HTTPClient.Timeout != timeout {
		t.Fatalf("expected timeout %s, got %s", timeout, client.HTTPClient.Timeout)
	}
}

func TestDoAllowsEmptySuccessfulResponseBody(t *testing.T) {
	client := NewClient("gitee.com", "test-token")
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNoContent,
				Body:       io.NopCloser(bytes.NewReader(nil)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}),
	}

	var response struct {
		Login string `json:"login"`
	}

	if err := client.Do(http.MethodGet, "/repos/owner/repo/collaborators/user", nil, &response); err != nil {
		t.Fatalf("client.Do() returned error for empty successful body: %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
