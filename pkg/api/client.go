package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
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
	_, respBody, err := c.doRequest(method, path, body)
	if err != nil {
		return err
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

// DoWithHeaders performs a request and returns response headers for callers that need pagination metadata.
func (c *Client) DoWithHeaders(method, path string, body interface{}, response interface{}) (http.Header, error) {
	headers, respBody, err := c.doRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	if len(respBody) == 0 {
		return headers, nil
	}

	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return headers, nil
}

func (c *Client) doRequest(method, path string, body interface{}) (http.Header, []byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	url := config.ApiUrl(c.host) + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", authHeaderPrefix+c.token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return resp.Header.Clone(), respBody, nil
}

// DoFromEndpoint performs an HTTP request from an Endpoint, auto-extracting method
func (c *Client) DoFromEndpoint(e Endpoint, pathArgs []interface{}, body interface{}, response interface{}) error {
	path := e.Build(pathArgs...)
	if e.Method == GET && body != nil {
		query, err := buildQueryFromRequest(body)
		if err != nil {
			return err
		}
		if query != "" {
			path += "?" + query
		}
		body = nil
	}
	return c.Do(string(e.Method), path, body, response)
}

// DoFromEndpointWithHeaders performs an endpoint request and exposes response headers for pagination-aware callers.
func (c *Client) DoFromEndpointWithHeaders(e Endpoint, pathArgs []interface{}, body interface{}, response interface{}) (http.Header, error) {
	path := e.Build(pathArgs...)
	if e.Method == GET && body != nil {
		query, err := buildQueryFromRequest(body)
		if err != nil {
			return nil, err
		}
		if query != "" {
			path += "?" + query
		}
		body = nil
	}
	return c.DoWithHeaders(string(e.Method), path, body, response)
}

// buildQueryFromRequest 将 GET 请求的结构化参数编码为查询字符串，避免被错误发送为请求体。
func buildQueryFromRequest(body interface{}) (string, error) {
	values := url.Values{}
	if err := appendQueryValues(values, reflect.ValueOf(body)); err != nil {
		return "", err
	}
	return values.Encode(), nil
}

// appendQueryValues 递归展开结构体和指针字段，仅编码非零值。
func appendQueryValues(values url.Values, value reflect.Value) error {
	if !value.IsValid() {
		return nil
	}

	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return fmt.Errorf("failed to encode GET query params: unsupported type %T", value.Interface())
	}

	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := valueType.Field(i)
		if !structField.IsExported() {
			continue
		}

		tagName, omitEmpty := parseJSONTag(structField.Tag.Get("json"))
		if tagName == "" || tagName == "-" {
			continue
		}

		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		}

		if omitEmpty && field.IsZero() {
			continue
		}

		values.Set(tagName, formatQueryValue(field))
	}

	return nil
}

func parseJSONTag(tag string) (string, bool) {
	if tag == "" {
		return "", false
	}

	parts := strings.Split(tag, ",")
	name := parts[0]
	omitEmpty := false
	for _, option := range parts[1:] {
		if option == "omitempty" {
			omitEmpty = true
			break
		}
	}
	return name, omitEmpty
}

func formatQueryValue(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.String:
		return value.String()
	default:
		return fmt.Sprint(value.Interface())
	}
}
