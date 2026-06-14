package api

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ricsy/gt/pkg/api/response"
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

func TestDoFromEndpointEncodesGetOptionsAsQuery(t *testing.T) {
	client := NewClient("gitee.com", "test-token")
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v5/notifications/messages" {
				t.Fatalf("expected path /notifications/messages, got %s", req.URL.Path)
			}
			if req.URL.Query().Get("unread") != "true" {
				t.Fatalf("expected unread=true query, got %q", req.URL.RawQuery)
			}
			if req.URL.Query().Get("page") != "2" {
				t.Fatalf("expected page=2 query, got %q", req.URL.RawQuery)
			}
			if req.Body == nil {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"total_count":0,"list":[]}`)),
					Header:     make(http.Header),
					Request:    req,
				}, nil
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("failed to read request body: %v", err)
			}
			if len(body) != 0 {
				t.Fatalf("expected empty GET body, got %q", string(body))
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"total_count":0,"list":[]}`)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}),
	}

	_, err := client.ListMessages(response.ListMessagesOptions{
		Unread:  BoolPtr(true),
		Page:    2,
		PerPage: 0,
	})
	if err != nil {
		t.Fatalf("client.ListMessages() returned error: %v", err)
	}
}

func TestGetNotificationCountEncodesUnreadAsQuery(t *testing.T) {
	client := NewClient("gitee.com", "test-token")
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v5/notifications/count" {
				t.Fatalf("expected path /notifications/count, got %s", req.URL.Path)
			}
			if req.URL.Query().Get("unread") != "true" {
				t.Fatalf("expected unread=true query, got %q", req.URL.RawQuery)
			}
			if req.Body != nil {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}
				if len(body) != 0 {
					t.Fatalf("expected empty GET body, got %q", string(body))
				}
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"total_count":1,"notification_count":0,"message_count":1}`)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}),
	}

	count, err := client.GetNotificationCount(BoolPtr(true))
	if err != nil {
		t.Fatalf("client.GetNotificationCount() returned error: %v", err)
	}
	if count.MessageCount != 1 {
		t.Fatalf("expected message count 1, got %d", count.MessageCount)
	}
}

func TestDoWithHeadersReturnsResponseHeaders(t *testing.T) {
	client := NewClient("gitee.com", "test-token")
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			header := make(http.Header)
			header.Set("X-Total-Count", "12")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`[]`)),
				Header:     header,
				Request:    req,
			}, nil
		}),
	}

	var response []struct{}
	headers, err := client.DoWithHeaders(http.MethodGet, "/repos/owner/repo/commits", nil, &response)
	if err != nil {
		t.Fatalf("client.DoWithHeaders() returned error: %v", err)
	}
	if headers.Get("X-Total-Count") != "12" {
		t.Fatalf("expected X-Total-Count header, got %q", headers.Get("X-Total-Count"))
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
