package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ricsy/gt/pkg/config"
)

const (
	authHeaderPrefix = "token "
	// DefaultTimeout is the default timeout for API requests.
	DefaultTimeout = 30 * time.Second
)

type Client struct {
	host       string
	token      string
	HTTPClient *http.Client
}

func NewClient(host, token string) *Client {
	return NewClientWithTimeout(host, token, DefaultTimeout)
}

// NewClientWithTimeout creates a new Gitee API client with the provided timeout.
func NewClientWithTimeout(host, token string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	return &Client{
		host:       host,
		token:      token,
		HTTPClient: &http.Client{Timeout: timeout},
	}
}

// BoolPtr returns a pointer to a bool value
func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}

// Token returns the client's auth token (needed by api command for raw requests).
func (c *Client) Token() string {
	return c.token
}

func (c *Client) Do(method, path string, body interface{}, response interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	url := config.ApiUrl(c.host) + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", authHeaderPrefix+c.token)
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

	// 某些存在性检查接口会返回 204 No Content，此时不应继续按 JSON 反序列化。
	if len(respBody) == 0 {
		return nil
	}

	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// DoFromEndpoint performs an HTTP request from an Endpoint, auto-extracting method
func (c *Client) DoFromEndpoint(e Endpoint, pathArgs []interface{}, body interface{}, response interface{}) error {
	path := e.Build(pathArgs...)
	return c.Do(string(e.Method), path, body, response)
}
