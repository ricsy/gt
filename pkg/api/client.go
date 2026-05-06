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
	defaultTimeout   = 30 * time.Second
)

type Client struct {
	host       string
	token      string
	HTTPClient *http.Client
}

func NewClient(host, token string) *Client {
	return &Client{
		host:       host,
		token:      token,
		HTTPClient: &http.Client{Timeout: defaultTimeout},
	}
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
