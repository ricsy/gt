package api

import (
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
