package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const (
	apiPathRepos       = "/repos/%s/%s"
	apiPathOrgs        = "/orgs/%s"
	apiPathUser        = "/user"
	apiPathUserOrgs    = "/user/orgs"
	apiPathUserRepos   = "/user/repos"
	apiPathUsersRepos  = "/users/%s/repos"
	apiPathIssues      = "/repos/%s/%s/issues"
	apiPathIssueCreate = "/repos/%s/issues"
	apiPathIssueUpdate = "/repos/%s/issues/%s"
	apiPathPRs         = "/repos/%s/%s/pulls"
	apiPathReleases    = "/repos/%s/%s/releases"
)

// Client is the Gitee API client
type Client struct {
	Host       string
	Token      string
	HTTPClient *http.Client
}

var (
	defaultHTTPClient     *http.Client
	defaultHTTPClientOnce sync.Once
)

func getDefaultHTTPClient() *http.Client {
	defaultHTTPClientOnce.Do(func() {
		defaultHTTPClient = &http.Client{}
	})
	return defaultHTTPClient
}

// NewClient creates a new Gitee API client
func NewClient(host, token string) *Client {
	return &Client{
		Host:       host,
		Token:      token,
		HTTPClient: getDefaultHTTPClient(),
	}
}

// Do performs an HTTP request to the Gitee API
func (c *Client) Do(method, path string, body interface{}, response interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	url := fmt.Sprintf("https://%s/api/v5%s", c.Host, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}
