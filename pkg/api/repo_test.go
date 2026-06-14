package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClient_ListRepos(t *testing.T) {
	// TODO: mock HTTP test
	t.Skip("requires mock setup")
}

func TestClient_GetRepo(t *testing.T) {
	t.Skip("requires mock setup")
}

func TestClient_ListUserRepos(t *testing.T) {
	t.Skip("requires mock setup")
}

func TestClient_CreateRepo(t *testing.T) {
	t.Skip("requires mock setup")
}

func TestClient_GetRepoCommitHistorySummaryUsesHeaderCount(t *testing.T) {
	client := NewClient("gitee.com", "test-token")
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v5/repos/gitee/demo-repo/commits" {
				t.Fatalf("unexpected path: %s", req.URL.Path)
			}
			if req.URL.Query().Get("sha") != "main" {
				t.Fatalf("expected sha=main, got %q", req.URL.RawQuery)
			}
			if req.URL.Query().Get("per_page") != "100" {
				t.Fatalf("expected per_page=100, got %q", req.URL.RawQuery)
			}

			header := make(http.Header)
			header.Set("X-Total-Count", "12")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`[{
					"sha":"abc",
					"commit":{"message":"feat: keep history","committer":{"date":"2026-06-14T16:00:00Z"}}
				}]`)),
				Header:  header,
				Request: req,
			}, nil
		}),
	}

	summary, err := client.GetRepoCommitHistorySummary("gitee", "demo-repo", "main")
	if err != nil {
		t.Fatalf("GetRepoCommitHistorySummary() returned error: %v", err)
	}
	if summary.Count != 12 {
		t.Fatalf("summary.Count = %d, want 12", summary.Count)
	}
	if summary.Latest == nil || summary.Latest.Commit.Message != "feat: keep history" {
		t.Fatalf("expected latest commit message, got %+v", summary.Latest)
	}
}
