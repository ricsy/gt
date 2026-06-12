package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClient_MergePR(t *testing.T) {
	client := NewClient("gitee.com", "test-token")
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Fatalf("expected PUT, got %s", req.Method)
			}
			if req.URL.Path != "/api/v5/repos/owner/repo/pulls/7/merge" {
				t.Fatalf("expected merge path, got %s", req.URL.Path)
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if !strings.Contains(string(body), `"merge_method":"squash"`) {
				t.Fatalf("expected merge_method in body, got %s", string(body))
			}

			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}),
	}

	if err := client.MergePR("owner", "repo", 7, MergePRRequest{MergeMethod: "squash"}); err != nil {
		t.Fatalf("MergePR() returned error: %v", err)
	}
}

func TestClient_ReviewPR(t *testing.T) {
	client := NewClient("gitee.com", "test-token")
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", req.Method)
			}
			if req.URL.Path != "/api/v5/repos/owner/repo/pulls/7/review" {
				t.Fatalf("expected review path, got %s", req.URL.Path)
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if !strings.Contains(string(body), `"force":true`) {
				t.Fatalf("expected force=true in body, got %s", string(body))
			}

			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}),
	}

	if err := client.ReviewPR("owner", "repo", 7, ReviewPRRequest{Force: true}); err != nil {
		t.Fatalf("ReviewPR() returned error: %v", err)
	}
}

func TestClient_TestPR(t *testing.T) {
	client := NewClient("gitee.com", "test-token")
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", req.Method)
			}
			if req.URL.Path != "/api/v5/repos/owner/repo/pulls/7/test" {
				t.Fatalf("expected test path, got %s", req.URL.Path)
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if !strings.Contains(string(body), `"force":true`) {
				t.Fatalf("expected force=true in body, got %s", string(body))
			}

			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}),
	}

	if err := client.TestPR("owner", "repo", 7, TestPRRequest{Force: true}); err != nil {
		t.Fatalf("TestPR() returned error: %v", err)
	}
}
